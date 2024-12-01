package kube

import (
	"context"
	"fmt"
	"io"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodConfig holds the configuration for creating/updating a Kubernetes Pod
type PodConfig struct {
	Name             string
	Namespace        string
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
	Env              []corev1.EnvVar
	VolumeMounts     []corev1.VolumeMount
	Volumes          []corev1.Volume
	ServiceAccount   string
}

// CreatePod creates a Pod based on the provided PodConfig
func (kc *KubeClient) CreatePod(config PodConfig) error {
	podsClient := kc.Clientset.CoreV1().Pods(config.Namespace)

	// Define the Pod spec
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        config.Name,
			Namespace:   config.Namespace,
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
			Volumes:                       config.Volumes,
			RestartPolicy:                 config.RestartPolicy,
			TerminationGracePeriodSeconds: &config.TerminationGrace,
		},
	}

	// Create the Pod
	_, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create pod: %w", err)
	}

	fmt.Printf("Pod %s created successfully in namespace %s\n", config.Name, config.Namespace)
	return nil
}

// ListPods lists all Pods in the specified namespace
func (kc *KubeClient) ListPods(namespace, labelSelector string) ([]corev1.Pod, error) {
	podsClient := kc.Clientset.CoreV1().Pods(namespace)

	// Fetch Pods
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}
	pods, err := podsClient.List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	return pods.Items, nil
}

// DeletePod deletes a Pod by name in the specified namespace
func (kc *KubeClient) DeletePod(name, namespace string) error {
	podsClient := kc.Clientset.CoreV1().Pods(namespace)

	// Delete the Pod
	err := podsClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod: %w", err)
	}

	fmt.Printf("Pod %s deleted successfully from namespace %s\n", name, namespace)
	return nil
}

// StreamPodLogs streams logs from a specific Pod container in real-time
func (kc *KubeClient) StreamPodLogs(podName, namespace, containerName string) error {
	podsClient := kc.Clientset.CoreV1().Pods(namespace)

	logOptions := &corev1.PodLogOptions{
		Container: containerName,
		Follow:    true,
	}

	// Stream logs
	stream, err := podsClient.GetLogs(podName, logOptions).Stream(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to stream logs from pod %s: %w", podName, err)
	}
	defer stream.Close()

	fmt.Printf("Streaming logs for pod %s, container %s:\n", podName, containerName)
	_, err = io.Copy(os.Stdout, stream)
	if err != nil {
		return fmt.Errorf("error while streaming logs: %w", err)
	}

	return nil
}
