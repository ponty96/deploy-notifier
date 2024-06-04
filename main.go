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

	k8sClient := k8s.InitK8sClient("", "")
	k8sClient.GetEventsForNamespace(k8s.K8sQueryFilter{Namespace: "default", Context: ""})
}
