package kube

import (
	"context"
	"fmt"
	"time"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeClient provides a comprehensive Kubernetes client
type KubeClient struct {
	Clientset     *kubernetes.Clientset
	DynamicClient dynamic.Interface
	RESTConfig    *rest.Config
}

// NewKubeClient initializes a Kubernetes client with typed, dynamic, and REST clients.
func NewKubeClient(kubeconfigPath string) (*KubeClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create typed Kubernetes client: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic Kubernetes client: %w", err)
	}

	return &KubeClient{
		Clientset:     clientset,
		DynamicClient: dynamicClient,
		RESTConfig:    config,
	}, nil
}

// TestConnection validates connectivity to the Kubernetes API server.
func (kc *KubeClient) TestConnection() error {
	version, err := kc.Clientset.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to connect to the Kubernetes API server: %w", err)
	}
	fmt.Printf("Connected to Kubernetes cluster. Version: %s\n", version.GitVersion)
	return nil
}

// CreateResource creates a Kubernetes resource dynamically.
func (kc *KubeClient) CreateResource(resource *unstructured.Unstructured, gvr schema.GroupVersionResource, namespace string) error {
	_, err := kc.DynamicClient.Resource(gvr).Namespace(namespace).Create(context.TODO(), resource, v1.CreateOptions{})
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return fmt.Errorf("resource already exists: %w", err)
		}
		return fmt.Errorf("failed to create resource: %w", err)
	}
	fmt.Printf("Resource %s created successfully in namespace %s\n", resource.GetName(), namespace)
	return nil
}

// UpdateResource updates an existing Kubernetes resource dynamically.
func (kc *KubeClient) UpdateResource(resource *unstructured.Unstructured, gvr schema.GroupVersionResource, namespace string) error {
	_, err := kc.DynamicClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), resource, v1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}
	fmt.Printf("Resource %s updated successfully in namespace %s\n", resource.GetName(), namespace)
	return nil
}

// DeleteResource deletes a Kubernetes resource dynamically.
func (kc *KubeClient) DeleteResource(name string, gvr schema.GroupVersionResource, namespace string) error {
	err := kc.DynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return fmt.Errorf("resource not found: %w", err)
		}
		return fmt.Errorf("failed to delete resource: %w", err)
	}
	fmt.Printf("Resource %s deleted successfully from namespace %s\n", name, namespace)
	return nil
}

// GetResource retrieves a Kubernetes resource dynamically.
func (kc *KubeClient) GetResource(name string, gvr schema.GroupVersionResource, namespace string) (*unstructured.Unstructured, error) {
	resource, err := kc.DynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return nil, fmt.Errorf("resource not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}
	return resource, nil
}

// ListResources retrieves a list of Kubernetes resources dynamically.
func (kc *KubeClient) ListResources(gvr schema.GroupVersionResource, namespace string) ([]*unstructured.Unstructured, error) {
	resourceList, err := kc.DynamicClient.Resource(gvr).Namespace(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	var resources []*unstructured.Unstructured
	for _, resource := range resourceList.Items {
		copy := resource
		resources = append(resources, &copy)
	}
	return resources, nil
}

// WaitForResource waits for a resource to be in a desired state.
func (kc *KubeClient) WaitForResource(name string, gvr schema.GroupVersionResource, namespace string, condition func(*unstructured.Unstructured) bool, timeout time.Duration) error {
	return wait.PollImmediate(2*time.Second, timeout, func() (done bool, err error) {
		resource, err := kc.GetResource(name, gvr, namespace)
		if err != nil {
			if k8sErrors.IsNotFound(err) {
				return false, nil // Keep waiting
			}
			return false, err // Fail on other errors
		}
		return condition(resource), nil
	})
}
