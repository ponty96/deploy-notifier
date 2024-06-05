package k8s

import (
	"context"
	"os"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	clientSet *kubernetes.Clientset
}

type K8sEvent struct {
	Name            string
	Namespace       string
	ResourceVersion string
	Action          string
}

type K8sQueryFilter struct {
	Namespace string
	Context   string
}

func InitK8sClient(contextName string, kubeConfigPath string) K8sClient {
	// context := os.Getenv("K8S_CONTEXT")
	// kubeConfigPath := os.Getenv("KUBECONFIG")

	config, err := buildClusterConfig(contextName, kubeConfigPath)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return K8sClient{clientSet: clientset}
}

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

func buildClusterConfig(context string, kubeconfig string) (*rest.Config, error) {
	userHomeDir, err := os.UserHomeDir()
	if kubeconfig == "" && err != nil {
		kubeconfig = userHomeDir + "/.kube/config"
	}

	if kubeconfig != "" {
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			&clientcmd.ConfigOverrides{CurrentContext: context},
		).ClientConfig()
	} else {
		zap.L().Sugar().Warn("No kubeconfig path provided, using in-cluster config")
		// if we can't get the home dir and no config is passed, use in-cluster config
		return rest.InClusterConfig()
	}
}
