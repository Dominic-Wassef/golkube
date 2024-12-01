package main

import (
	"log"

	"golkube/pkg/kube"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// registerUtilityCommands registers all utility commands to the root command.
func registerUtilityCommands(kubeClient *kube.KubeClient) {
	// Command to monitor Kubernetes resources
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor Kubernetes resources",
		Run: func(cmd *cobra.Command, args []string) {
			namespace := viper.GetString("kubernetes.namespace")
			monitorInterval := viper.GetDuration("monitoring.interval")

			// Monitoring Kubernetes resources for pod health
			err := kubeClient.MonitorPodHealth(namespace, "", monitorInterval)
			if err != nil {
				log.Fatalf("Error monitoring resources: %v", err)
			}
		},
	}

	// Add the monitor command to the root command
	rootCmd.AddCommand(monitorCmd)
}
