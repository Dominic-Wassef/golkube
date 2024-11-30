package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

// rootCmd is the base command for the CLI
var rootCmd = &cobra.Command{
	Use:   "golkube",
	Short: "golkube: A Kubernetes and Docker Management Tool",
	Long: `golkube is a tool for managing Kubernetes resources and Docker workflows.
It provides exhaustive configuration, dynamic resource generation, validation, and robust error handling.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to golkube! Use the --help flag to explore available commands.")
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// initializeConfig loads the configuration file
func initializeConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("default")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./configs")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
}

// init function initializes the CLI
func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to the configuration file (default is ./configs/default.yaml)")
	rootCmd.PersistentFlags().String("namespace", "default", "Kubernetes namespace to use")
	rootCmd.PersistentFlags().String("kubeconfig", "~/.kube/config", "Path to the kubeconfig file")

	// Bind global flags to Viper
	viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))

	// Initialize configuration
	cobra.OnInitialize(initializeConfig)
}

// main function
func main() {
	Execute()
}
