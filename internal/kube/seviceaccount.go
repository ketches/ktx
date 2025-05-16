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
	"k8s.io/client-go/kubernetes"
)

// GetServiceAccount returns a service account
func GetServiceAccount(kubeClientset kubernetes.Interface, serviceAccountName, namespace string) *v1.ServiceAccount {
	serviceAccount, err := kubeClientset.CoreV1().ServiceAccounts(namespace).Get(context.Background(), serviceAccountName, metav1.GetOptions{})
	if err != nil {
		output.Fatal("Failed to get service account %s.", serviceAccountName)
	}
	return serviceAccount
}

// ListServiceAccounts returns a list of service accounts
func ListServiceAccounts(kubeClientset kubernetes.Interface, namespace string) []string {
	serviceAccounts, err := kubeClientset.CoreV1().ServiceAccounts(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		output.Fatal("Failed to list service accounts: %s", err)
	}
	var sa []string
	for _, serviceAccount := range serviceAccounts.Items {
		sa = append(sa, serviceAccount.Name)
	}
	return sa
}
