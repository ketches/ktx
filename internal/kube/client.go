/*
Copyright 2025 The Ketches Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kube

import (
	"github.com/ketches/ktx/internal/output"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientOrDie creates a new kubernetes client from the given
// kubeconfig file and context, and panics if it fails.
func ClientOrDie(kubeConfigFile, ctx string) kubernetes.Interface {
	clientset, err := Client(kubeConfigFile, ctx)
	if err != nil {
		output.Fatal("Failed to create kubernetes client: %s from file %s", err, kubeConfigFile)
	}
	return clientset
}

// Client creates a new kubernetes client from the given
// kubeconfig file and context.
func Client(kubeConfigFile, ctx string) (kubernetes.Interface, error) {
	client, err := kubernetes.NewForConfig(configOrDie(kubeConfigFile, ctx))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// DiscoveryClient creates a new kubernetes discovery client
// from the given kubeconfig file and context.
func DiscoveryClient(kubeConfigFile, ctx string) (discovery.DiscoveryInterface, error) {
	client, err := Client(kubeConfigFile, ctx)
	if err != nil {
		return nil, err
	}
	return client.Discovery(), nil
}

// configOrDie creates a new kubernetes rest configOrDie from the
// given kubeconfig file, and panics if it fails.
func configOrDie(kubeConfigFile, ctx string) *rest.Config {
	config, err := config(kubeConfigFile, ctx)
	if err != nil {
		output.Fatal("Failed to build kubernetes rest config: %s from file %s", err, kubeConfigFile)
	}
	return config
}

// config creates a new kubernetes rest config from the given
// kubeconfig file and context.
func config(kubeConfigFile, ctx string) (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = kubeConfigFile

	configOverrides := &clientcmd.ConfigOverrides{}
	if len(ctx) > 0 {
		configOverrides.CurrentContext = ctx
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}
