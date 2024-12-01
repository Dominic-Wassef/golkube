package main

import (
	"fmt"
	"log"
	"os"

	"golkube/pkg/docker"
	"golkube/pkg/registry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// registerDockerCommands registers all Docker-related commands to the root command.
func registerDockerCommands(dockerRegistry *registry.RegistryClient) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "build",
		Short: "Build a Docker image",
		Run: func(cmd *cobra.Command, args []string) {
			buildConfig := docker.BuildConfig{
				Tag:         viper.GetString("build.tags.0"),
				ContextDir:  viper.GetString("build.context"),
				Dockerfile:  viper.GetString("build.dockerfile"),
				RegistryURL: viper.GetString("registry.url"),
				Push:        true,
			}
			err := docker.BuildImage(buildConfig)
			if err != nil {
				log.Fatalf("Error building image: %v", err)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "push",
		Short: "Push a Docker image to a registry",
		Run: func(cmd *cobra.Command, args []string) {
			imageTag := viper.GetString("build.tags.0")
			registryURL := viper.GetString("registry.url")
			username := os.Getenv("DOCKER_USERNAME")
			password := os.Getenv("DOCKER_PASSWORD")

			if username == "" || password == "" {
				log.Fatal("DOCKER_USERNAME and DOCKER_PASSWORD environment variables are required")
			}

			err := dockerRegistry.PushImage(imageTag, registryURL, username, password)
			if err != nil {
				log.Fatalf("Error pushing Docker image: %v", err)
			}
			fmt.Printf("Image %s pushed successfully to %s\n", imageTag, registryURL)
		},
	})

	registerRunCommand()
	registerStopCommand()
	registerRestartCommand()
	registerInspectCommand()
}

func registerRunCommand() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a Docker container",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			image, _ := cmd.Flags().GetString("image")
			ports, _ := cmd.Flags().GetStringToString("ports")
			env, _ := cmd.Flags().GetStringArray("env")

			// Validate port mappings
			for hostPort, containerPort := range ports {
				if hostPort == "" || containerPort == "" {
					log.Fatalf("Invalid port mapping: %s=%s. Ports must be in the format hostPort=containerPort", hostPort, containerPort)
				}
			}

			// Prepare Docker container config
			containerConfig := docker.ContainerConfig{
				Name:  name,
				Image: image,
				Ports: ports,
				Env:   env,
			}

			// Start the container
			containerID, err := docker.StartContainer(containerConfig)
			if err != nil {
				log.Fatalf("Error starting container: %v", err)
			}
			fmt.Printf("Container started successfully with ID: %s\n", containerID)
		},
	}

	runCmd.Flags().String("name", "default-container", "Name of the container")
	runCmd.Flags().String("image", "", "Docker image to use (required)")
	runCmd.Flags().StringToString("ports", nil, "Ports to map in format hostPort=containerPort")
	runCmd.Flags().StringArray("env", nil, "Environment variables to set in the container")
	runCmd.MarkFlagRequired("image")

	rootCmd.AddCommand(runCmd)
}

func registerStopCommand() {
	stopCmd := &cobra.Command{
		Use:   "stop <container-id>",
		Short: "Stop a running Docker container",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Error: container-id is required")
			}
			containerID := args[0]
			err := docker.StopContainer(containerID)
			if err != nil {
				log.Fatalf("Error stopping container: %v", err)
			}
			fmt.Printf("Container %s stopped successfully\n", containerID)
		},
	}
	rootCmd.AddCommand(stopCmd)
}

func registerRestartCommand() {
	restartCmd := &cobra.Command{
		Use:   "restart <container-id>",
		Short: "Restart a Docker container",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Error: container-id is required")
			}
			containerID := args[0]
			err := docker.RestartContainer(containerID)
			if err != nil {
				log.Fatalf("Error restarting container: %v", err)
			}
			fmt.Printf("Container %s restarted successfully\n", containerID)
		},
	}
	rootCmd.AddCommand(restartCmd)
}

func registerInspectCommand() {
	inspectCmd := &cobra.Command{
		Use:   "inspect <container-id>",
		Short: "Inspect a Docker container",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Error: container-id is required")
			}
			containerID := args[0]
			containerJSON, err := docker.InspectContainer(containerID)
			if err != nil {
				log.Fatalf("Error inspecting container: %v", err)
			}
			fmt.Printf("Container details:\n%+v\n", containerJSON)
		},
	}
	rootCmd.AddCommand(inspectCmd)
}
