package controller

import (
	"fmt"
	"github.com/sirupsen/logrus"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type event struct {
	key       string
	eventType string
}

type controller struct {
	client   kubernetes.Interface
	informer cache.SharedIndexInformer
	queue    workqueue.RateLimitingInterface
}

func Start() {
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatal(err)
	}

	factory := informers.NewSharedInformerFactory(kc, 0)
	informer := factory.Core().V1().Pods().Informer()

	c := newController(kc, informer)
	stopCh := make(chan struct{})
	defer close(stopCh)

	c.Run(stopCh)
}

func newController(kc kubernetes.Interface, informer cache.SharedIndexInformer) *controller {
	q := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			var event event
			var err error
			event.key, err = cache.MetaNamespaceKeyFunc(obj)
			event.eventType = "create"
			if err == nil {
				q.Add(event)
			}
			logrus.Infof("Event received of type [%s] for [%s]", event.eventType, event.key)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			var event event
			var err error
			event.key, err = cache.MetaNamespaceKeyFunc(oldObj)
			event.eventType = "update"
			if err == nil {
				q.Add(event)
			}
			logrus.Infof("Event received of type [%s] for [%s]", event.eventType, event.key)
		},
		DeleteFunc: func(obj interface{}) {
			var event event
			var err error
			event.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			event.eventType = "delete"
			if err == nil {
				q.Add(event)
			}
			logrus.Infof("Event received of type [%s] for [%s]", event.eventType, event.key)
		},
	})

	return &controller{
		client:   kc,
		informer: informer,
		queue:    q,
	}
}

func (c *controller) Run(stopper <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	logrus.Info("Starting pod-informer")

	go c.informer.Run(stopper)

	logrus.Info("Synchronizing events...")

	if !cache.WaitForCacheSync(stopper, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for cache"))
		logrus.Info("Synchronization failed")
		return
	}

	logrus.Info("synchronization completed")

	wait.Until(c.runProcessing, time.Second, stopper)
}

func (c *controller) runProcessing() {
	for c.processNextItem() {

	}
}

func (c *controller) processNextItem() bool {
	e, term := c.queue.Get()

	if term {
		return false
	}

	err := c.processItem(e.(event))

	if err == nil {
		c.queue.Forget(e)
		return true
	}

	return true
}

func (c *controller) processItem(e event) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(e.key)

	if err != nil {
		logrus.Info("Processing error", err)
		return fmt.Errorf("error fetching object with key %s with error %s", e.key, err)
	}

	logrus.Infof("Processed one event of type [%s] for object [%s]", e.eventType, obj)

	return nil
}
