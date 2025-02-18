package kube

import (
	"github.com/poneding/ktx/internal/output"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client creates a new kubernetes client from the given kubeconfig file.
func Client(kubeConfigFile string) kubernetes.Interface {
	clientset, err := kubernetes.NewForConfig(config(kubeConfigFile))
	if err != nil {
		output.Fatal("Failed to create kubernetes client: %s from file %s", err, kubeConfigFile)
	}
	return clientset
}

// config creates a new kubernetes rest config from the given kubeconfig file.
func config(kubeConfigFile string) *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	if err != nil {
		output.Fatal("Failed to build kubernetes rest config: %s from file %s", err, kubeConfigFile)
	}
	return config
}
