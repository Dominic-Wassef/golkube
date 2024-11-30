package main

import (
	"fmt"
	"golkube/pkg/kube"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deployCmd handles Kubernetes deployments
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy resources to Kubernetes",
	Long: `The deploy command applies Kubernetes manifests or programmatically
creates resources such as Deployments, Services, ConfigMaps, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		namespace := viper.GetString("namespace")
		kubeconfig := viper.GetString("kubeconfig")
		client, err := kube.NewKubeClient(kubeconfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
			os.Exit(1)
		}

		// Example Deployment Configuration
		err = client.CreateDeployment(kube.DeploymentConfig{
			Name:          "example-deployment",
			Namespace:     namespace,
			Replicas:      3,
			Image:         "nginx:latest",
			ContainerPort: 80,
			Labels: map[string]string{
				"app": "example",
			},
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating deployment: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Deployment created successfully!")
	},
}

// monitorCmd monitors Kubernetes resources
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor Kubernetes resources in real-time",
	Long:  `Streams logs and metrics from Kubernetes resources such as pods and deployments.`,
	Run: func(cmd *cobra.Command, args []string) {
		namespace := viper.GetString("namespace")
		kubeconfig := viper.GetString("kubeconfig")
		client, err := kube.NewKubeClient(kubeconfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
			os.Exit(1)
		}

		err = client.StreamEvents(namespace)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error monitoring events: %v\n", err)
			os.Exit(1)
		}
	},
}

// buildCmd handles Docker image builds
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and tag a Docker image",
	Long:  `Builds a Docker image from the specified context and Dockerfile, with optional tagging.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Docker build logic (to be implemented)
		fmt.Println("Building Docker image...")
	},
}

func init() {
	// Add subcommands to rootCmd
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(monitorCmd)
	rootCmd.AddCommand(buildCmd)
}
