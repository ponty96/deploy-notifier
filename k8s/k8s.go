package k8s

import (
	"context"

	"go.uber.org/zap"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type K8s struct {
	client      K8sClient
	controllers []K8sController
}

type Resource struct {
	Deployment            bool `json:"deployment"`
	ReplicationController bool `json:"rc"`
	ReplicaSet            bool `json:"rs"`
	DaemonSet             bool `json:"ds"`
	StatefulSet           bool `json:"statefulset"`
	Services              bool `json:"svc"`
	Pod                   bool `json:"po"`
	Job                   bool `json:"job"`
	Node                  bool `json:"node"`
	ClusterRole           bool `json:"clusterrole"`
	ClusterRoleBinding    bool `json:"clusterrolebinding"`
	ServiceAccount        bool `json:"sa"`
	PersistentVolume      bool `json:"pv"`
	Namespace             bool `json:"ns"`
	Secret                bool `json:"secret"`
	ConfigMap             bool `json:"configmap"`
	Ingress               bool `json:"ing"`
	HPA                   bool `json:"hpa"`
	Event                 bool `json:"event"`
	CoreEvent             bool `json:"coreevent"`
}

type K8sConfig struct {
	ContextName string
	KubeConfig  string
	ResourceTM  Resource
	Namespace   string
}

func Setup(k8sCfonfig K8sConfig) {
	// Connect to k8s
	client := InitK8sClient(k8sCfonfig.ContextName, k8sCfonfig.KubeConfig)
	if k8sCfonfig.ResourceTM.Pod {
		zap.L().Sugar().Infof("ResourceTM Pods")
		// informerFactory := informers.NewSharedInformerFactory(client.clientSet, time.Minute*10)
		// zap.L().Sugar().Infof("InformerFactory: %v", informerFactory)
		// informer := informerFactory.Core().V1().Pods().Informer()
		// zap.L().Sugar().Infof("informer: %v", informer)
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					options.FieldSelector = ""
					return client.clientSet.CoreV1().Pods(k8sCfonfig.Namespace).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					options.FieldSelector = ""
					return client.clientSet.CoreV1().Pods(k8sCfonfig.Namespace).Watch(context.Background(), options)
				},
			},
			&api_v1.Event{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := NewController(client, eventHandler, informer)
		stopAllPodsCh := make(chan struct{})
		defer close(stopAllPodsCh)

		zap.L().Sugar().Infof("Starting controller %v", c)
		go c.Run(stopAllPodsCh)
	}
}

func eventHandler(event K8sEvent) {
	zap.L().Sugar().Infof("event: %v", event.Key)
}
