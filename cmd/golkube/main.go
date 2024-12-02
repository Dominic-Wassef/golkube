package main

import (
	"golkube/pkg/commands"
	"log"

	"github.com/spf13/viper"
)

func main() {
	// Load environment and configuration files
	commands.LoadEnvironment()
	commands.LoadConfiguration()

	// Set the logging level based on the log-level flag
	logLevel := viper.GetString("log-level")
	commands.SetLogLevel(logLevel)

	// Initialize Kubernetes and Docker clients
	kubeClient, dockerRegistry := commands.InitializeClients()

	// Register all CLI commands
	commands.RegisterCommands(kubeClient, dockerRegistry)

	// Execute the root command
	if err := commands.RootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
