package kube

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceConfig holds the configuration for creating/updating a Kubernetes Service
type ServiceConfig struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
	Selector    map[string]string
	Type        corev1.ServiceType
	Ports       []corev1.ServicePort
}

// CreateService creates a Service based on the provided ServiceConfig
func (kc *KubeClient) CreateService(config ServiceConfig) error {
	servicesClient := kc.Clientset.CoreV1().Services(config.Namespace)

	// Define the Service spec
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        config.Name,
			Namespace:   config.Namespace,
			Labels:      config.Labels,
			Annotations: config.Annotations,
		},
		Spec: corev1.ServiceSpec{
			Selector: config.Selector,
			Type:     config.Type,
			Ports:    config.Ports,
		},
	}

	// Create the Service
	_, err := servicesClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	fmt.Printf("Service %s created successfully in namespace %s\n", config.Name, config.Namespace)
	return nil
}

// UpdateService updates an existing Service based on the provided ServiceConfig
func (kc *KubeClient) UpdateService(config ServiceConfig) error {
	servicesClient := kc.Clientset.CoreV1().Services(config.Namespace)

	// Fetch the existing Service
	existingService, err := servicesClient.Get(context.TODO(), config.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to fetch service: %w", err)
	}

	// Update fields
	existingService.Spec.Selector = config.Selector
	existingService.Spec.Ports = config.Ports
	existingService.Annotations = config.Annotations

	// Update the Service
	_, err = servicesClient.Update(context.TODO(), existingService, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	fmt.Printf("Service %s updated successfully in namespace %s\n", config.Name, config.Namespace)
	return nil
}

// ListServices lists all Services in the specified namespace
func (kc *KubeClient) ListServices(namespace, labelSelector string) ([]corev1.Service, error) {
	servicesClient := kc.Clientset.CoreV1().Services(namespace)

	// Fetch Services
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}
	services, err := servicesClient.List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	return services.Items, nil
}

// DeleteService deletes a Service by name in the specified namespace
func (kc *KubeClient) DeleteService(name, namespace string) error {
	servicesClient := kc.Clientset.CoreV1().Services(namespace)

	// Delete the Service
	err := servicesClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	fmt.Printf("Service %s deleted successfully from namespace %s\n", name, namespace)
	return nil
}
