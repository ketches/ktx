package kube

import (
	"context"

	"github.com/poneding/ktx/internal/output"
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
