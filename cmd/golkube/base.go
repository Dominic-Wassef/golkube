package main

import (
	"fmt"
	"golkube/pkg/kube"
	"golkube/pkg/registry"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for the Golkube CLI
var rootCmd = &cobra.Command{
	Use:   "golkube",
	Short: "Golkube CLI for managing Docker and Kubernetes workflows",
	Long:  `Golkube is a tool for managing Docker and Kubernetes workflows, including deployments, image builds, and monitoring.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Golkube! Use --help to see available commands.")
	},
}

// Helper function to find or create the "kube" command
func findOrCreateKubeCommand() *cobra.Command {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "kube" {
			return cmd
		}
	}
	kubeCmd := &cobra.Command{
		Use:   "kube",
		Short: "Manage Kubernetes resources dynamically",
	}
	rootCmd.AddCommand(kubeCmd)
	return kubeCmd
}

// registerCommands adds all CLI commands to rootCmd
func registerCommands(kubeClient *kube.KubeClient, dockerRegistry *registry.RegistryClient) {
	// Ensure the "kube" command is registered first
	kubeCmd := findOrCreateKubeCommand()

	// Register Docker-related commands and subcommands
	registerDockerCommands(dockerRegistry)
	registerRunCommand()
	registerStopCommand()
	registerRestartCommand()
	registerInspectCommand()

	// Register Kubernetes-related subcommands and deployment commands under the kubeCmd
	registerKubeSubCommands(kubeCmd, kubeClient)

	// Register configuration commands
	registerConfigCommands()

	// Register utility commands like "monitor"
	registerUtilityCommands(kubeClient)

	// Register pipeline-related commands
	registerPipelineCommand()
}
