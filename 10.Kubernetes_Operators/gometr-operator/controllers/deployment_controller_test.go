package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	v1 "slurm.io/gometr-operator/api/v1"
	"time"
)

var _ = Describe("Gometr controller", func() {
	const (
		GometrName      = "gometr-sample"
		GometrNamespace = "default"
		ServiceName     = "gometr-service"
		DeploymentName  = "gometr-deployment"
		GometrReplicas  = int32(3)

		timeout  = time.Second * 50
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)
	Context("When creating gometr deployment", func() {
		It("Should create gometr workload", func() {
			By("By creating new Gometr")
			ctx := context.Background()
			gometrDep := &v1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "custom-apps.slurm.io/v1",
					Kind:       "Gometr",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      GometrName,
					Namespace: GometrNamespace,
				},
				Spec: v1.DeploymentSpec{
					Port:           8000,
					ServiceName:    ServiceName,
					Replicas:       GometrReplicas,
					Image:          "nortsx/gometr:v1.0",
					DeploymentName: DeploymentName,
				},
			}

			Expect(k8sClient.Create(ctx, gometrDep)).Should(Succeed())

			lookupKey := types.NamespacedName{
				Namespace: GometrNamespace,
				Name:      GometrName,
			}

			createdGometrDeployment := &v1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupKey, createdGometrDeployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdGometrDeployment.Spec.Replicas).Should(Equal(GometrReplicas))
			Expect(createdGometrDeployment.Spec.ServiceName).Should(Equal(ServiceName))

			By("Checking that deployment is created")

			deploymentLookupKey := types.NamespacedName{Name: DeploymentName, Namespace: GometrNamespace}
			createdDeployment := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, deploymentLookupKey, createdDeployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(*createdDeployment.Spec.Replicas).Should(Equal(GometrReplicas))
			By("Checking that service is created")
			serviceLookupKey := types.NamespacedName{Name: ServiceName, Namespace: GometrNamespace}
			createdService := &corev1.Service{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, serviceLookupKey, createdService)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(len(createdService.Spec.Ports)).Should(Equal(1))
		})
	})
})
