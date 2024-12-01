package kube

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

// MonitorConfig holds the configuration for monitoring resources
type MonitorConfig struct {
	Namespace     string
	ResourceType  schema.GroupVersionResource
	LabelSelector string
	Timeout       time.Duration
}

// MonitorEvents streams resource events (e.g., Pods, Deployments) in real-time
func (kc *KubeClient) MonitorEvents(config MonitorConfig) error {
	client := kc.DynamicClient.Resource(config.ResourceType).Namespace(config.Namespace)

	// Watch for resource events
	watcher, err := client.Watch(context.TODO(), metav1.ListOptions{
		LabelSelector: config.LabelSelector,
	})
	if err != nil {
		return fmt.Errorf("failed to watch resource events: %w", err)
	}
	defer watcher.Stop()

	fmt.Printf("Monitoring events for resource type %s in namespace %s\n", config.ResourceType.Resource, config.Namespace)

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Added:
			fmt.Println("[ADDED]", describeResource(event.Object))
		case watch.Modified:
			fmt.Println("[MODIFIED]", describeResource(event.Object))
		case watch.Deleted:
			fmt.Println("[DELETED]", describeResource(event.Object))
		default:
			fmt.Println("[UNKNOWN EVENT TYPE]")
		}
	}

	return nil
}

// MonitorPodHealth continuously checks the status of all Pods in a namespace
func (kc *KubeClient) MonitorPodHealth(namespace, labelSelector string, interval time.Duration) error {
	podsClient := kc.Clientset.CoreV1().Pods(namespace)

	fmt.Printf("Monitoring pod health in namespace %s\n", namespace)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		pods, err := podsClient.List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			fmt.Printf("Error listing pods: %v\n", err)
			continue
		}

		for _, pod := range pods.Items {
			status := summarizePodStatus(&pod)
			fmt.Printf("Pod %s: %s\n", pod.Name, status)
		}
	}

	// Ensure a return statement for the function
	return nil
}

// describeResource provides a brief description of a resource from a runtime.Object
func describeResource(obj runtime.Object) string {
	// Use the correct `meta.Accessor` for accessing resource metadata
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return "Unable to describe resource"
	}
	return fmt.Sprintf("Name: %s, Namespace: %s, Labels: %v", metaObj.GetName(), metaObj.GetNamespace(), metaObj.GetLabels())
}

// summarizePodStatus returns a brief status summary for a Pod
func summarizePodStatus(pod *corev1.Pod) string {
	status := "Unknown"
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			status = "Ready"
			break
		} else if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionFalse {
			status = "Not Ready"
		}
	}
	return status
}
