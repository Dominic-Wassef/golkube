package commands

import (
	"fmt"
	"golkube/pkg/kube"
	"golkube/pkg/registry"

	"github.com/spf13/cobra"
)

// RootCmd is the base command for the Golkube CLI
var RootCmd = &cobra.Command{
	Use:   "golkube",
	Short: "Golkube CLI for managing Docker and Kubernetes workflows",
	Long:  `Golkube is a tool for managing Docker and Kubernetes workflows, including deployments, image builds, and monitoring.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Golkube! Use --help to see available commands.")
	},
}

// RegisterCommands adds all CLI commands to RootCmd
func RegisterCommands(kubeClient *kube.KubeClient, dockerRegistry *registry.RegistryClient) {
	// Register Docker-related commands
	RegisterDockerCommands(dockerRegistry)

	// Register Kubernetes-related commands
	RegisterKubeCommands(kubeClient)

	// Register configuration commands
	RegisterConfigCommands()

	// Register utility commands like "monitor"
	RegisterUtilityCommands(kubeClient)

	// Register pipeline-related commands
	RegisterPipelineCommand()
}
