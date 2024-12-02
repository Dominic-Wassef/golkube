package commands

import (
	"golkube/pkg/kube"
	"golkube/pkg/registry"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// loadEnvironment loads environment variables from a `.env` file
func LoadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: No .env file found or failed to load it: %v", err)
	}
}

// loadConfiguration loads the application's configuration file
func LoadConfiguration() {
	RootCmd.PersistentFlags().String("config", "configs/default.yaml", "Path to configuration file")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	configFile := viper.GetString("config")
	if configFile == "" {
		configFile = "configs/default.yaml"
	}
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	RootCmd.PersistentFlags().String("namespace", viper.GetString("kubernetes.namespace"), "Kubernetes namespace")
	viper.BindPFlag("kubernetes.namespace", RootCmd.PersistentFlags().Lookup("namespace"))

	RootCmd.PersistentFlags().String("kubeconfig", viper.GetString("kubernetes.kubeconfig"), "Path to kubeconfig file")
	viper.BindPFlag("kubeconfig", RootCmd.PersistentFlags().Lookup("kubeconfig"))
}

// setLogLevel configures the logging level
func SetLogLevel(level string) {
	switch level {
	case "debug":
		log.SetFlags(log.LstdFlags | log.Lshortfile) // Include file and line number
		log.Println("Log level set to DEBUG")
	case "info":
		log.SetFlags(log.LstdFlags)
		log.Println("Log level set to INFO")
	case "warn":
		log.SetFlags(0) // Minimal logging
		log.Println("Log level set to WARN")
	case "error":
		log.SetFlags(0)
		log.Println("Log level set to ERROR")
	}
}

// initializeClients sets up Kubernetes and Docker clients
func InitializeClients() (*kube.KubeClient, *registry.RegistryClient) {
	kubeconfig := viper.GetString("kubernetes.kubeconfig")
	kubeClient, err := kube.NewKubeClient(kubeconfig)
	if err != nil {
		log.Fatalf("Error initializing Kubernetes client: %v", err)
	}

	err = kubeClient.TestConnection()
	if err != nil {
		log.Fatalf("Kubernetes API connection test failed: %v", err)
	}

	dockerRegistry, err := registry.NewRegistryClient()
	if err != nil {
		log.Fatalf("Error initializing Docker registry client: %v", err)
	}

	return kubeClient, dockerRegistry
}
