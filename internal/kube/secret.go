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
	"context"

	"github.com/ketches/ktx/internal/output"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
)

// GetSecret returns a secret
func GetSecret(kubeClientset kubernetes.Interface, secretName, namespace string) *v1.Secret {
	secret, err := kubeClientset.CoreV1().Secrets(namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		output.Fatal("Failed to get secret: %s", err)
	}
	return secret
}

// CreateServiceAccountTokenSecret creates a secret for a service account
func CreateServiceAccountTokenSecret(kubeClientset kubernetes.Interface, serviceAccountName, namespace string) *v1.Secret {
	secret, err := kubeClientset.CoreV1().Secrets(namespace).Create(context.Background(), &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName + "-token-" + rand.String(5),
			Namespace: namespace,
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": serviceAccountName,
			},
		},
		Type: v1.SecretTypeServiceAccountToken,
	}, metav1.CreateOptions{})
	if err != nil {
		output.Fatal("Failed to create secret for service account %s from namespace %s", serviceAccountName, namespace)
	}

	return secret
}

// CreateSecret creates a secret
func CreateSecret(kubeClientset kubernetes.Interface, secret *v1.Secret, namespace string) *v1.Secret {
	secret, err := kubeClientset.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		output.Fatal("Failed to create secret %s from namespace %s", secret.Name, namespace)
	}

	return secret
}
