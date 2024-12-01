package kube

import (
	"context"
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/util/retry"
)

// StreamEvents streams Kubernetes resource events in real-time
func (kc *KubeClient) StreamEvents(namespace string) error {
	client := kc.Clientset.CoreV1().Events(namespace)

	// Watch for resource events
	watcher, err := client.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to start event watcher: %w", err)
	}
	defer watcher.Stop()

	fmt.Printf("Streaming events in namespace %s\n", namespace)

	// Process events from the watcher channel
	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Added:
			describeEvent("ADDED", event.Object)
		case watch.Modified:
			describeEvent("MODIFIED", event.Object)
		case watch.Deleted:
			describeEvent("DELETED", event.Object)
		default:
			log.Printf("Unknown event type: %v\n", event.Type)
		}
	}
	return nil
}

// EventWatcherConfig holds the configuration for watching Kubernetes events
type EventWatcherConfig struct {
	Namespace     string
	ResourceType  schema.GroupVersionResource
	LabelSelector string
	FieldSelector string
	RetryInterval time.Duration
	RetryTimeout  time.Duration
	OnAdd         func(runtime.Object)
	OnModify      func(runtime.Object)
	OnDelete      func(runtime.Object)
}

// WatchEvents continuously streams Kubernetes resource events
func (kc *KubeClient) WatchEvents(config EventWatcherConfig) error {
	client := kc.DynamicClient.Resource(config.ResourceType).Namespace(config.Namespace)

	// Retry mechanism for transient errors
	err := retry.OnError(retry.DefaultRetry, func(err error) bool {
		log.Printf("Retrying due to error: %v\n", err)
		return true // Always retry for this example; can be customized
	}, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), config.RetryTimeout)
		defer cancel()

		// Start watching events
		watcher, err := client.Watch(ctx, metav1.ListOptions{
			LabelSelector: config.LabelSelector,
			FieldSelector: config.FieldSelector,
		})
		if err != nil {
			return fmt.Errorf("failed to watch events for %s: %w", config.ResourceType.Resource, err)
		}
		defer watcher.Stop()

		fmt.Printf("Watching events for %s in namespace %s\n", config.ResourceType.Resource, config.Namespace)

		// Process events
		for event := range watcher.ResultChan() {
			processEvent(event, config)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("event watcher failed: %w", err)
	}
	return nil
}

// processEvent processes a single Kubernetes event
func processEvent(event watch.Event, config EventWatcherConfig) {
	switch event.Type {
	case watch.Added:
		if config.OnAdd != nil {
			config.OnAdd(event.Object)
		}
	case watch.Modified:
		if config.OnModify != nil {
			config.OnModify(event.Object)
		}
	case watch.Deleted:
		if config.OnDelete != nil {
			config.OnDelete(event.Object)
		}
	default:
		log.Printf("Unknown event type: %v\n", event.Type)
	}
}

// describeEvent provides a detailed log for a Kubernetes resource event
func describeEvent(eventType string, obj runtime.Object) {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		log.Printf("[%s] Unable to access object metadata: %v", eventType, err)
		return
	}
	fmt.Printf("[%s] Name: %s, Namespace: %s, Labels: %v\n",
		eventType, metaObj.GetName(), metaObj.GetNamespace(), metaObj.GetLabels())
}
