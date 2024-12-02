package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"golkube/pkg/docker"
	"golkube/pkg/registry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RegisterDockerCommands registers all Docker-related commands under the "docker" namespace.
func RegisterDockerCommands(dockerRegistry *registry.RegistryClient) {
	dockerCmd := findOrCreateDockerCommand()
	dockerCmd.AddCommand(buildDockerCmd())
	dockerCmd.AddCommand(pushDockerCmd())
	dockerCmd.AddCommand(runDockerCmd())
	dockerCmd.AddCommand(stopDockerCmd())
	dockerCmd.AddCommand(restartDockerCmd())
	dockerCmd.AddCommand(inspectDockerCmd())
}

// Helper function to find or create the "docker" command.
func findOrCreateDockerCommand() *cobra.Command {
	for _, cmd := range RootCmd.Commands() {
		if cmd.Use == "docker" {
			return cmd
		}
	}
	dockerCmd := &cobra.Command{
		Use:   "docker",
		Short: "Manage Docker resources",
	}
	RootCmd.AddCommand(dockerCmd)
	return dockerCmd
}

// Helper function to validate Docker credentials.
func validateDockerCredentials() (string, string) {
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if username == "" || password == "" {
		log.Fatalf("Error: DOCKER_USERNAME and DOCKER_PASSWORD environment variables are required.")
	}
	return username, password
}

// Build Docker image.
func buildDockerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build a Docker image",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := validateDockerCredentials()
			log.Printf("Using Docker credentials: username=%s", username)

			// Use the proper tag with namespace
			tag := fmt.Sprintf("%s/golkube:latest", username)
			contextDir := viper.GetString("build.context")
			dockerfile := viper.GetString("build.dockerfile")

			log.Printf("Building Docker image with tag '%s', contextDir '%s', and dockerfile '%s'", tag, contextDir, dockerfile)

			// Build the Docker image
			buildCmd := exec.Command("docker", "build", "-t", tag, "-f", dockerfile, contextDir)
			buildCmd.Stdout = os.Stdout
			buildCmd.Stderr = os.Stderr
			if err := buildCmd.Run(); err != nil {
				log.Fatalf("Docker build failed: %v", err)
			}

			log.Println("Image built successfully.")

			// Push the image automatically after build
			log.Printf("Pushing image '%s' to Docker Hub...", tag)
			pushCmd := exec.Command("docker", "push", tag)
			pushCmd.Stdout = os.Stdout
			pushCmd.Stderr = os.Stderr
			if err := pushCmd.Run(); err != nil {
				log.Fatalf("Docker push failed: %v", err)
			}
			log.Printf("Image '%s' successfully pushed to Docker Hub.", tag)
		},
	}
}

// Push Docker image to registry.
func pushDockerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Push a Docker image to a registry",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := validateDockerCredentials()

			// Ensure the correct tag with namespace
			imageTag := fmt.Sprintf("%s/golkube:latest", username)
			log.Printf("Pushing Docker image with tag '%s'...", imageTag)

			// Push the Docker image
			pushCmd := exec.Command("docker", "push", imageTag)
			pushCmd.Stdout = os.Stdout
			pushCmd.Stderr = os.Stderr
			if err := pushCmd.Run(); err != nil {
				log.Fatalf("Docker push failed: %v", err)
			}
			log.Printf("Image '%s' successfully pushed to Docker Hub.", imageTag)
		},
	}
}

// Run a Docker container.
func runDockerCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a Docker container",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			image, _ := cmd.Flags().GetString("image")
			ports, _ := cmd.Flags().GetStringToString("ports")
			env, _ := cmd.Flags().GetStringArray("env")

			for hostPort, containerPort := range ports {
				if hostPort == "" || containerPort == "" {
					log.Fatalf("Invalid port mapping: %s=%s. Ports must be in the format hostPort=containerPort", hostPort, containerPort)
				}
			}

			containerConfig := docker.ContainerConfig{
				Name:  name,
				Image: image,
				Ports: ports,
				Env:   env,
			}

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

	return runCmd
}

// Stop a running Docker container
func stopDockerCmd() *cobra.Command {
	return &cobra.Command{
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
}

// Restart a Docker container
func restartDockerCmd() *cobra.Command {
	return &cobra.Command{
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
}

// Inspect a Docker container
func inspectDockerCmd() *cobra.Command {
	return &cobra.Command{
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
}
