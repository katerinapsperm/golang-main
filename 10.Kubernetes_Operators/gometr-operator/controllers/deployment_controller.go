/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gometrv1 "slurm.io/gometr-operator/api/v1"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=gometr.slurm.io,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gometr.slurm.io,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gometr.slurm.io,resources=deployments/finalizers,verbs=update

const (
	repeatingInterval = 120
	ownerKey          = ".metadata.controller"
)

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Deployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	var gometrDeployment gometrv1.Deployment

	if err := r.Get(ctx, req.NamespacedName, &gometrDeployment); err != nil {
		log.Log.Error(err, "unable to fetch gometr resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Log.Info("Got gometr deployment")

	var deployment *appsv1.Deployment
	deploymentLookupKey := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      gometrDeployment.Spec.DeploymentName,
	}

	if err := r.Get(ctx, deploymentLookupKey, deployment); client.IgnoreNotFound(err) != nil {
		deployment = nil
	}

	createdDeployments := 0
	if deployment == nil {
		deploymentTemplate := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gometrDeployment.Spec.DeploymentName,
				Namespace: req.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &(gometrDeployment.Spec.Replicas),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "gometr",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "gometr",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  gometrDeployment.Spec.DeploymentName,
								Image: gometrDeployment.Spec.Image,
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										Protocol:      corev1.ProtocolTCP,
										ContainerPort: gometrDeployment.Spec.Port,
									},
								},
							},
						},
					},
				},
			},
		}
		err := ctrl.SetControllerReference(&gometrDeployment, deploymentTemplate, r.Scheme)
		if err != nil {
			log.Log.Error(err, "unable to wire deployment")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, deploymentTemplate); err != nil {
			log.Log.Error(err, "unable to create deployment")
			return ctrl.Result{RequeueAfter: time.Duration(repeatingInterval) * time.Second}, err
		}
		log.Log.Info("Sucessfully created a deployment")
		createdDeployments++
	} else {
		createdDeployments = 1
		deployment.Spec.Replicas = &gometrDeployment.Spec.Replicas
		deployment.Spec.Template.Spec.Containers[0].Image = gometrDeployment.Spec.Image
		if err := r.Update(ctx, deployment); err != nil {
			log.Log.Error(err, "unable to update deployment")
			return ctrl.Result{RequeueAfter: time.Duration(repeatingInterval) * time.Second}, err
		}
	}

	var services corev1.ServiceList
	if err := r.List(ctx, &services, client.InNamespace(req.Namespace), client.MatchingFields{ownerKey: req.Name}); err != nil {
		log.Log.Error(err, "unable to get services list")
		return ctrl.Result{}, err
	}

	createdServices := len(services.Items)
	if len(services.Items) == 0 {
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gometrDeployment.Spec.ServiceName,
				Namespace: req.Namespace,
			},
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeClusterIP,
				Selector: map[string]string{
					"app": "gometr",
				},
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						TargetPort: intstr.IntOrString{IntVal: gometrDeployment.Spec.Port},
						Port:       gometrDeployment.Spec.Port,
					},
				},
			},
		}
		err := ctrl.SetControllerReference(&gometrDeployment, service, r.Scheme)
		if err != nil {
			log.Log.Error(err, "unable to wire service")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, service); err != nil {
			log.Log.Error(err, "unable to create service")
			return ctrl.Result{RequeueAfter: time.Duration(repeatingInterval) * time.Second}, nil
		}
		log.Log.Info("Successfully created a service")
		createdServices++
	} else {
		service := services.Items[0]
		service.ObjectMeta.Name = gometrDeployment.Spec.ServiceName
		service.Spec.Ports[0].NodePort = gometrDeployment.Spec.Port
		service.Spec.Ports[0].Port = gometrDeployment.Spec.Port
		if err := r.Update(ctx, &service); err != nil {
			log.Log.Error(err, "unable to update service")
			return ctrl.Result{RequeueAfter: time.Duration(repeatingInterval) * time.Second}, nil
		}
	}

	gometrDeployment.Status.Deployments = createdDeployments
	gometrDeployment.Status.Services = createdServices

	err := r.Update(ctx, &gometrDeployment)

	if err != nil {
		log.Log.Error(err, "unable to update gometr status")
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Service{}, ownerKey, func(rawObj client.Object) []string {
		service := rawObj.(*corev1.Service)
		owner := metav1.GetControllerOf(service)

		if owner == nil {
			return nil
		}

		if owner.Kind != "Deployment" {
			return nil
		}

		return []string{owner.Name}
	})

	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&gometrv1.Deployment{}).
		Complete(r)
}
