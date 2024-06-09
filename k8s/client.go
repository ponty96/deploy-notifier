package k8s

import (
	"os"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	clientSet *kubernetes.Clientset
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
