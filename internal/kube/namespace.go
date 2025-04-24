package kube

import (
	"context"

	"github.com/ketches/ktx/internal/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListNamespaces returns a list of namespaces
func ListNamespaces(kubeClientset kubernetes.Interface) []string {
	namespaces, err := kubeClientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		output.Fatal("Failed to list namespaces: %s", err)
	}
	var ns []string
	for _, namespace := range namespaces.Items {
		ns = append(ns, namespace.Name)
	}
	return ns
}
