package main

import (
	"os"

	"deploy-notifier/k8s"

	"go.uber.org/zap"
)

func main() {
	// do something here to set environment depending on an environment variable
	// or command-line flag

	if os.Getenv("ENVIRONMENT") == "production" {
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	} else {
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	}

	zap.L().Info("Hello from Zap logger!")

	go func() {
		httpServer := HTTPServer{}
		httpServer.serveHTTP()
	}()

	context := os.Getenv("K8S_CONTEXT")
	kubeConfigPath := os.Getenv("KUBECONFIG")

	k8sConfig := k8s.K8sConfig{
		ContextName: context,
		KubeConfig:  kubeConfigPath,
		ResourceTM: k8s.Resource{
			Pod: true,
		},
		Namespace: "default",
	}
	k8s.Setup(k8sConfig)
	// k8sClient.GetEventsForNamespace(k8s.K8sQueryFilter{Namespace: "default", Context: ""})
}
