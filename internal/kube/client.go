package kube

import (
	"github.com/ketches/ktx/internal/output"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client creates a new kubernetes client from the given kubeconfig file and context.
func Client(kubeConfigFile, ctx string) kubernetes.Interface {
	clientset, err := kubernetes.NewForConfig(config(kubeConfigFile, ctx))
	if err != nil {
		output.Fatal("Failed to create kubernetes client: %s from file %s", err, kubeConfigFile)
	}
	return clientset
}

// config creates a new kubernetes rest config from the given kubeconfig file.
func config(kubeConfigFile, ctx string) *rest.Config {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = kubeConfigFile

	configOverrides := &clientcmd.ConfigOverrides{}
	if len(ctx) > 0 {
		configOverrides.CurrentContext = ctx
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := clientConfig.ClientConfig()
	if err != nil {
		output.Fatal("Failed to build kubernetes rest config: %s from file %s", err, kubeConfigFile)
	}
	return config
}
