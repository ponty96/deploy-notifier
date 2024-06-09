package k8s

import (
	"go.uber.org/zap"
	"k8s.io/client-go/informers"
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
		informerFactory := informers.NewSharedInformerFactory(client.clientSet, 0)
		zap.L().Sugar().Infof("InformerFactory: %v", informerFactory)
		informer := informerFactory.Core().V1().Pods().Informer()
		zap.L().Sugar().Infof("informer: %v", informer)

		c := NewController(client, eventHandler, informer)
		stopAllPodsCh := make(chan struct{})
		defer close(stopAllPodsCh)

		go c.Run(stopAllPodsCh)
	}
}

func eventHandler(event K8sEvent) {
	zap.L().Sugar().Infof("event: %v", event.Key)
}
