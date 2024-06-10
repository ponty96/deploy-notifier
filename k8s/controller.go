package k8s

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type K8sQueryFilter struct {
	Namespace string
	Context   string
}

type K8sEvent struct {
	Action          string
	APIVersion      string
	Name            string
	Namespace       string
	ResourceVersion string
	ResourceType    string
	// Object          runtime.Object
	Key string
}

type K8sController struct {
	informer  cache.SharedIndexInformer
	workqueue workqueue.RateLimitingInterface
	client    K8sClient
	handler   func(event K8sEvent)
}

const maxRetries = 3

// func (k8sClient *K8sClient) GetEventsForNamespace(queryFilter K8sQueryFilter, eventHandler func(event K8sEvent))
func (k8sClient *K8sClient) GetEventsForNamespace(queryFilter K8sQueryFilter) {
	events, err := k8sClient.clientSet.CoreV1().Events(queryFilter.Namespace).List(context.TODO(), metav1.ListOptions{TypeMeta: metav1.TypeMeta{Kind: "Pod"}})
	if err != nil {
		panic(err.Error())
	}
	for _, event := range events.Items {
		myEvent := K8sEvent{
			Name:            event.ObjectMeta.Name,
			Namespace:       event.ObjectMeta.Namespace,
			ResourceVersion: event.InvolvedObject.ResourceVersion,
			Action:          event.Action,
		}
		zap.L().Sugar().Infof("PodName: %v", myEvent.Name)
		zap.L().Sugar().Infof("Namespace: %v", myEvent.Action)
	}
}

func NewController(
	client K8sClient,
	eventHandler func(event K8sEvent),
	informer cache.SharedIndexInformer,
) *K8sController {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var newEvent K8sEvent
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			zap.L().Sugar().Info("I got into the AddFunc:")
			pod := obj.(*v1.Pod)
			zap.L().Sugar().Info("I got into the AddFunc:")
			// var ok bool
			newEvent.Action = "create"
			newEvent.Namespace = pod.ObjectMeta.Namespace // namespace retrived in processItem incase namespace value is empty
			newEvent.Key, err = cache.MetaNamespaceKeyFunc(obj)
			newEvent.ResourceVersion = pod.ObjectMeta.ResourceVersion
			newEvent.ResourceType = pod.TypeMeta.Kind
			newEvent.APIVersion = pod.TypeMeta.APIVersion
			newEvent.Name = pod.ObjectMeta.Name

			zap.L().Sugar().Infof("Pod name is:  %v", newEvent.Name)

			if err == nil {
				queue.Add(newEvent)
			}
		},
		// UpdateFunc: func(old, new interface{}) {
		// 	var ok bool
		// 	newEvent.Action = "update"
		// 	newEvent.Namespace = obj.ObjectMeta.Namespace // namespace retrived in processItem incase namespace value is empty
		// 	newEvent.key, err = cache.MetaNamespaceKeyFunc(old)
		// 	newEvent.ResourceVersion = obj.InvolvedObject.ResourceVersion
		// 	newEvent.ResourceType = obj.TypeMeta.Kind
		// 	newEvent.APIVersion = obj.TypeMeta.APIVersion

		// 	if err == nil {
		// 		queue.Add(newEvent)
		// 	}
		// },
		// DeleteFunc: func(obj interface{}) {
		// 	var ok bool
		// 	newEvent.Action = "delete"
		// 	newEvent.Namespace = obj.ObjectMeta.Namespace // namespace retrived in processItem incase namespace value is empty
		// 	newEvent.key, err = cache.MetaNamespaceKeyFunc(old)
		// 	newEvent.ResourceVersion = obj.InvolvedObject.ResourceVersion
		// 	newEvent.ResourceType = obj.TypeMeta.Kind
		// 	newEvent.APIVersion = obj.TypeMeta.APIVersion

		// 	if err == nil {
		// 		queue.Add(newEvent)
		// 	}
		// },
	})

	return &K8sController{
		client:    client,
		informer:  informer,
		workqueue: queue,
		handler:   eventHandler,
	}
}

func (c *K8sController) Run(stopCh <-chan struct{}) {
	zap.L().Sugar().Info("Run")
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	zap.L().Sugar().Info("Controller is now running")
	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced is required for the cache.Controller interface.
func (c *K8sController) HasSynced() bool {
	return c.informer.HasSynced()
}
func (c *K8sController) processNextItem() bool {
	k8sEvent, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}
	defer c.workqueue.Done(k8sEvent)

	err := c.syncHandler(k8sEvent)
	if err == nil {
		c.workqueue.Forget(k8sEvent)
	} else if c.workqueue.NumRequeues(k8sEvent) < maxRetries {
		c.workqueue.AddRateLimited(k8sEvent)
	} else {
		zap.L().Sugar().Errorf("Dropping event %q from the queue: %v", k8sEvent, err)
		c.workqueue.Forget(k8sEvent)
		runtime.HandleError(err)
	}

	return true
}

func (c *K8sController) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

// this is the function where we call the handler
func (c *K8sController) syncHandler(k8sEvent interface{}) error {
	// event, ok := k8sEvent.(K8sEvent)

	// c.handler(event)
	// c.workqueue.Forget(k8sEvent)
	zap.L().Sugar().Infof("Event: %v", k8sEvent)

	return nil
}
