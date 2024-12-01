package kube

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapConfig holds the configuration for creating/updating a Kubernetes ConfigMap
type ConfigMapConfig struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
	Data        map[string]string
	BinaryData  map[string][]byte
}

// CreateConfigMap creates a ConfigMap based on the provided ConfigMapConfig
func (kc *KubeClient) CreateConfigMap(config ConfigMapConfig) error {
	configMapsClient := kc.Clientset.CoreV1().ConfigMaps(config.Namespace)

	// Define the ConfigMap spec
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        config.Name,
			Namespace:   config.Namespace,
			Labels:      config.Labels,
			Annotations: config.Annotations,
		},
		Data:       config.Data,
		BinaryData: config.BinaryData,
	}

	// Create the ConfigMap
	_, err := configMapsClient.Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create configmap: %w", err)
	}

	fmt.Printf("ConfigMap %s created successfully in namespace %s\n", config.Name, config.Namespace)
	return nil
}

// UpdateConfigMap updates an existing ConfigMap based on the provided ConfigMapConfig
func (kc *KubeClient) UpdateConfigMap(config ConfigMapConfig) error {
	configMapsClient := kc.Clientset.CoreV1().ConfigMaps(config.Namespace)

	// Fetch the existing ConfigMap
	existingConfigMap, err := configMapsClient.Get(context.TODO(), config.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to fetch configmap: %w", err)
	}

	// Update fields
	existingConfigMap.Data = config.Data
	existingConfigMap.BinaryData = config.BinaryData
	existingConfigMap.Annotations = config.Annotations

	// Update the ConfigMap
	_, err = configMapsClient.Update(context.TODO(), existingConfigMap, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update configmap: %w", err)
	}

	fmt.Printf("ConfigMap %s updated successfully in namespace %s\n", config.Name, config.Namespace)
	return nil
}

// ListConfigMaps lists all ConfigMaps in the specified namespace
func (kc *KubeClient) ListConfigMaps(namespace, labelSelector string) ([]corev1.ConfigMap, error) {
	configMapsClient := kc.Clientset.CoreV1().ConfigMaps(namespace)

	// Fetch ConfigMaps
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}
	configMaps, err := configMapsClient.List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list configmaps: %w", err)
	}

	return configMaps.Items, nil
}

// DeleteConfigMap deletes a ConfigMap by name in the specified namespace
func (kc *KubeClient) DeleteConfigMap(name, namespace string) error {
	configMapsClient := kc.Clientset.CoreV1().ConfigMaps(namespace)

	// Delete the ConfigMap
	err := configMapsClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete configmap: %w", err)
	}

	fmt.Printf("ConfigMap %s deleted successfully from namespace %s\n", name, namespace)
	return nil
}
