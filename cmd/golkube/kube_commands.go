package main

import (
	"fmt"
	"log"

	"golkube/pkg/kube"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// registerKubeSubCommands registers Kubernetes-related commands under the "kube" command.
func registerKubeSubCommands(kubeCmd *cobra.Command, kubeClient *kube.KubeClient) {
	// Resource management commands
	kubeCmd.AddCommand(createResourceCmd(kubeClient))
	kubeCmd.AddCommand(updateResourceCmd(kubeClient))
	kubeCmd.AddCommand(deleteResourceCmd(kubeClient))
	kubeCmd.AddCommand(getResourceCmd(kubeClient))
	kubeCmd.AddCommand(listResourcesCmd(kubeClient))
	kubeCmd.AddCommand(waitResourceCmd(kubeClient))

	// ConfigMap management commands
	kubeCmd.AddCommand(createConfigMapCmd(kubeClient))
	kubeCmd.AddCommand(updateConfigMapCmd(kubeClient))
	kubeCmd.AddCommand(listConfigMapsCmd(kubeClient))
	kubeCmd.AddCommand(deleteConfigMapCmd(kubeClient))

	// Deployment management commands
	kubeCmd.AddCommand(createDeploymentCmd(kubeClient))
	kubeCmd.AddCommand(updateDeploymentCmd(kubeClient))
	kubeCmd.AddCommand(listDeploymentsCmd(kubeClient))
	kubeCmd.AddCommand(deleteDeploymentCmd(kubeClient))
}

// Command implementations

// Create Kubernetes resource dynamically
func createResourceCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "create-resource",
		Short: "Create a Kubernetes resource dynamically",
		Run: func(cmd *cobra.Command, args []string) {
			namespace := viper.GetString("kubernetes.namespace")
			gvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}
			resource := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name": "example-deployment",
					},
					"spec": map[string]interface{}{
						"replicas": int32(1),
						"selector": map[string]interface{}{
							"matchLabels": map[string]string{
								"app": "example-deployment",
							},
						},
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"labels": map[string]string{
									"app": "example-deployment",
								},
							},
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "nginx",
										"image": "nginx:latest",
										"ports": []interface{}{
											map[string]interface{}{
												"containerPort": int32(80),
											},
										},
									},
								},
							},
						},
					},
				},
			}

			err := kubeClient.CreateResource(resource, gvr, namespace)
			if err != nil {
				log.Fatalf("Error creating resource: %v", err)
			}
			fmt.Println("Resource created successfully.")
		},
	}
}

// Update an existing Kubernetes resource dynamically
func updateResourceCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "update-resource",
		Short: "Update an existing Kubernetes resource dynamically",
		Run: func(cmd *cobra.Command, args []string) {
			namespace := viper.GetString("kubernetes.namespace")
			gvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}
			resource := &unstructured.Unstructured{}
			resource.SetName("example-deployment")

			err := kubeClient.UpdateResource(resource, gvr, namespace)
			if err != nil {
				log.Fatalf("Error updating resource: %v", err)
			}
			fmt.Println("Resource updated successfully.")
		},
	}
}

// Delete a Kubernetes resource dynamically
func deleteResourceCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-resource",
		Short: "Delete a Kubernetes resource dynamically",
		Run: func(cmd *cobra.Command, args []string) {
			namespace := viper.GetString("kubernetes.namespace")
			gvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}
			name := "example-deployment"

			err := kubeClient.DeleteResource(name, gvr, namespace)
			if err != nil {
				log.Fatalf("Error deleting resource: %v", err)
			}
			fmt.Println("Resource deleted successfully.")
		},
	}
}

// Retrieve a Kubernetes resource dynamically
func getResourceCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "get-resource",
		Short: "Retrieve a Kubernetes resource dynamically",
		Run: func(cmd *cobra.Command, args []string) {
			namespace := viper.GetString("kubernetes.namespace")
			gvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}
			name := "example-deployment"

			resource, err := kubeClient.GetResource(name, gvr, namespace)
			if err != nil {
				log.Fatalf("Error retrieving resource: %v", err)
			}
			fmt.Printf("Retrieved resource: %+v\n", resource)
		},
	}
}

// List Kubernetes resources dynamically
func listResourcesCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "list-resources",
		Short: "List Kubernetes resources dynamically",
		Run: func(cmd *cobra.Command, args []string) {
			namespace := viper.GetString("kubernetes.namespace")
			gvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}

			resources, err := kubeClient.ListResources(gvr, namespace)
			if err != nil {
				log.Fatalf("Error listing resources: %v", err)
			}
			fmt.Printf("Resources: %+v\n", resources)
		},
	}
}

// Wait for a Kubernetes resource to reach a desired state
func waitResourceCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "wait-resource",
		Short: "Wait for a Kubernetes resource to reach a desired state",
		Run: func(cmd *cobra.Command, args []string) {
			name := viper.GetString("resource.name")
			namespace := viper.GetString("kubernetes.namespace")
			timeout := viper.GetDuration("resource.wait.timeout")

			gvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}

			condition := func(resource *unstructured.Unstructured) bool {
				available, found, _ := unstructured.NestedBool(resource.Object, "status", "available")
				return found && available
			}

			err := kubeClient.WaitForResource(name, gvr, namespace, condition, timeout)
			if err != nil {
				log.Fatalf("Error waiting for resource: %v", err)
			}
			fmt.Printf("Resource %s is now in the desired state.\n", name)
		},
	}
}

func createConfigMapCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "create-configmap",
		Short: "Create a Kubernetes ConfigMap",
		Run: func(cmd *cobra.Command, args []string) {
			config := kube.ConfigMapConfig{
				Name:      "example-configmap",
				Namespace: viper.GetString("kubernetes.namespace"),
				Labels: map[string]string{
					"app": "example",
				},
				Data: map[string]string{
					"key": "value",
				},
			}
			err := kubeClient.CreateConfigMap(config)
			if err != nil {
				log.Fatalf("Error creating ConfigMap: %v", err)
			}
			fmt.Println("ConfigMap created successfully.")
		},
	}
}

func updateConfigMapCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "update-configmap",
		Short: "Update an existing Kubernetes ConfigMap",
		Run: func(cmd *cobra.Command, args []string) {
			config := kube.ConfigMapConfig{
				Name:      "example-configmap",
				Namespace: viper.GetString("kubernetes.namespace"),
				Annotations: map[string]string{
					"updated": "true",
				},
				Data: map[string]string{
					"key": "new-value",
				},
			}
			err := kubeClient.UpdateConfigMap(config)
			if err != nil {
				log.Fatalf("Error updating ConfigMap: %v", err)
			}
			fmt.Println("ConfigMap updated successfully.")
		},
	}
}

func listConfigMapsCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "list-configmaps",
		Short: "List all Kubernetes ConfigMaps in a namespace",
		Run: func(cmd *cobra.Command, args []string) {
			namespace := viper.GetString("kubernetes.namespace")
			labelSelector := ""
			configMaps, err := kubeClient.ListConfigMaps(namespace, labelSelector)
			if err != nil {
				log.Fatalf("Error listing ConfigMaps: %v", err)
			}
			for _, cm := range configMaps {
				fmt.Printf("ConfigMap: %s, Data: %+v\n", cm.Name, cm.Data)
			}
		},
	}
}

func deleteConfigMapCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-configmap",
		Short: "Delete a Kubernetes ConfigMap",
		Run: func(cmd *cobra.Command, args []string) {
			name := "example-configmap"
			namespace := viper.GetString("kubernetes.namespace")
			err := kubeClient.DeleteConfigMap(name, namespace)
			if err != nil {
				log.Fatalf("Error deleting ConfigMap: %v", err)
			}
			fmt.Println("ConfigMap deleted successfully.")
		},
	}
}

// Deployment Commands

func createDeploymentCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "create-deployment",
		Short: "Create a Kubernetes Deployment",
		Run: func(cmd *cobra.Command, args []string) {
			config := kube.DeploymentConfig{
				Name:          "example-deployment",
				Namespace:     viper.GetString("kubernetes.namespace"),
				Replicas:      2,
				Image:         "nginx:latest",
				ContainerName: "nginx-container",
				ContainerPort: 80,
				Labels: map[string]string{
					"app": "example-deployment",
				},
				Annotations: map[string]string{
					"description": "Example deployment",
				},
			}
			err := kubeClient.CreateDeployment(config)
			if err != nil {
				log.Fatalf("Error creating Deployment: %v", err)
			}
			fmt.Println("Deployment created successfully.")
		},
	}
}

func updateDeploymentCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "update-deployment",
		Short: "Update an existing Kubernetes Deployment",
		Run: func(cmd *cobra.Command, args []string) {
			config := kube.DeploymentConfig{
				Name:      "example-deployment",
				Namespace: viper.GetString("kubernetes.namespace"),
				Replicas:  3,
				Image:     "nginx:stable",
			}
			err := kubeClient.UpdateDeployment(config)
			if err != nil {
				log.Fatalf("Error updating Deployment: %v", err)
			}
			fmt.Println("Deployment updated successfully.")
		},
	}
}

func listDeploymentsCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "list-deployments",
		Short: "List all Kubernetes Deployments in a namespace",
		Run: func(cmd *cobra.Command, args []string) {
			namespace := viper.GetString("kubernetes.namespace")
			labelSelector := ""
			deployments, err := kubeClient.ListDeployments(namespace, labelSelector)
			if err != nil {
				log.Fatalf("Error listing Deployments: %v", err)
			}
			for _, dep := range deployments {
				fmt.Printf("Deployment: %s, Replicas: %d\n", dep.Name, *dep.Spec.Replicas)
			}
		},
	}
}

func deleteDeploymentCmd(kubeClient *kube.KubeClient) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-deployment",
		Short: "Delete a Kubernetes Deployment",
		Run: func(cmd *cobra.Command, args []string) {
			name := "example-deployment"
			namespace := viper.GetString("kubernetes.namespace")
			err := kubeClient.DeleteDeployment(name, namespace)
			if err != nil {
				log.Fatalf("Error deleting Deployment: %v", err)
			}
			fmt.Println("Deployment deleted successfully.")
		},
	}
}
