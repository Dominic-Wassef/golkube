package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {
	// Load environment and configuration files
	loadEnvironment()
	loadConfiguration()

	// Set the logging level based on the log-level flag
	logLevel := viper.GetString("log-level")
	setLogLevel(logLevel)

	// Initialize Kubernetes and Docker clients
	kubeClient, dockerRegistry := initializeClients()

	// Register all CLI commands
	registerCommands(kubeClient, dockerRegistry)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
