package kube

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentConfig holds the configuration for creating/updating a Kubernetes Deployment
type DeploymentConfig struct {
	Name             string
	Namespace        string
	Replicas         int32
	Image            string
	ContainerName    string
	ContainerPort    int32
	Labels           map[string]string
	Annotations      map[string]string
	NodeSelector     map[string]string
	Affinity         *corev1.Affinity
	Tolerations      []corev1.Toleration
	LivenessProbe    *corev1.Probe
	ReadinessProbe   *corev1.Probe
	Resources        corev1.ResourceRequirements
	RestartPolicy    corev1.RestartPolicy
	TerminationGrace int64
	ImagePullSecrets []corev1.LocalObjectReference
	Env              []corev1.EnvVar
	VolumeMounts     []corev1.VolumeMount
	Volumes          []corev1.Volume
	ServiceAccount   string
}

// CreateDeployment creates a Deployment based on the provided DeploymentConfig
func (kc *KubeClient) CreateDeployment(config DeploymentConfig) error {
	deploymentsClient := kc.Clientset.AppsV1().Deployments(config.Namespace)

	// Define the Deployment spec
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        config.Name,
			Namespace:   config.Namespace,
			Labels:      config.Labels,
			Annotations: config.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &config.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: config.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      config.Labels,
					Annotations: config.Annotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  config.ContainerName,
							Image: config.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: config.ContainerPort,
								},
							},
							Env:            config.Env,
							Resources:      config.Resources,
							VolumeMounts:   config.VolumeMounts,
							LivenessProbe:  config.LivenessProbe,
							ReadinessProbe: config.ReadinessProbe,
						},
					},
					NodeSelector:                  config.NodeSelector,
					Affinity:                      config.Affinity,
					Tolerations:                   config.Tolerations,
					ServiceAccountName:            config.ServiceAccount,
					ImagePullSecrets:              config.ImagePullSecrets,
					Volumes:                       config.Volumes,
					RestartPolicy:                 config.RestartPolicy,
					TerminationGracePeriodSeconds: &config.TerminationGrace,
				},
			},
		},
	}

	// Create the Deployment
	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	fmt.Printf("Deployment %s created successfully in namespace %s\n", config.Name, config.Namespace)
	return nil
}

// UpdateDeployment updates an existing Deployment based on the provided DeploymentConfig
func (kc *KubeClient) UpdateDeployment(config DeploymentConfig) error {
	deploymentsClient := kc.Clientset.AppsV1().Deployments(config.Namespace)

	// Fetch the existing Deployment
	existingDeployment, err := deploymentsClient.Get(context.TODO(), config.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to fetch deployment: %w", err)
	}

	// Update fields
	existingDeployment.Spec.Replicas = &config.Replicas
	existingDeployment.Spec.Template.Spec.Containers[0].Image = config.Image
	existingDeployment.Spec.Template.Spec.Containers[0].Resources = config.Resources

	// Update the Deployment
	_, err = deploymentsClient.Update(context.TODO(), existingDeployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	fmt.Printf("Deployment %s updated successfully in namespace %s\n", config.Name, config.Namespace)
	return nil
}

// ListDeployments lists all Deployments in the specified namespace
func (kc *KubeClient) ListDeployments(namespace string, labelSelector string) ([]appsv1.Deployment, error) {
	deploymentsClient := kc.Clientset.AppsV1().Deployments(namespace)

	// Fetch deployments
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}
	deployments, err := deploymentsClient.List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	return deployments.Items, nil
}

// DeleteDeployment deletes a Deployment by name in the specified namespace
func (kc *KubeClient) DeleteDeployment(name, namespace string) error {
	deploymentsClient := kc.Clientset.AppsV1().Deployments(namespace)

	// Delete the Deployment
	err := deploymentsClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}

	fmt.Printf("Deployment %s deleted successfully from namespace %s\n", name, namespace)
	return nil
}
