package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// ContainerConfig holds the configuration for managing a Docker container
type ContainerConfig struct {
	Name        string
	Image       string
	Env         []string
	Cmd         []string
	Ports       map[string]string // hostPort:containerPort
	NetworkName string
	Volumes     map[string]string // hostPath:containerPath
}

// StartContainer creates and starts a Docker container
func StartContainer(config ContainerConfig) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	// Pull the image if it doesn't already exist
	_, err = cli.ImagePull(context.Background(), config.Image, types.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image %s: %w", config.Image, err)
	}

	// Configure ports
	portBindings := make(map[nat.Port][]nat.PortBinding)
	exposedPorts := make(nat.PortSet)
	for hostPort, containerPort := range config.Ports {
		port := nat.Port(containerPort)
		exposedPorts[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{
			{
				HostPort: hostPort,
			},
		}
	}

	// Configure the container
	containerConfig := &container.Config{
		Image:        config.Image,
		Env:          config.Env,
		Cmd:          config.Cmd,
		ExposedPorts: exposedPorts,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        convertVolumes(config.Volumes),
	}

	networkConfig := &network.NetworkingConfig{}

	// Create the container
	containerResp, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, networkConfig, nil, config.Name)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	if err := cli.ContainerStart(context.Background(), containerResp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("Container %s started successfully\n", config.Name)
	return containerResp.ID, nil
}

// StopContainer stops a running Docker container
func StopContainer(containerID string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker Client: %w", err)
	}
	defer cli.Close()

	if err := cli.ContainerStop(context.Background(), containerID, nil); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	fmt.Printf("Container %s stopped successfully\n", containerID)
	return nil
}

// RestartContainer restarts a running Docker container
func RestartContainer(containerID string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	if err := cli.ContainerRestart(context.Background(), containerID, nil); err != nil {
		return fmt.Errorf("failed to restart container : %w", err)
	}

	fmt.Printf("Container %s restarted successfully\n", containerID)
	return nil
}

// InspectContainer inspects a Docker container and returns detailed information
func InspectContainer(containerID string) (types.ContainerJSON, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("failed to inspect container: %w", err)
	}

	return containerJSON, nil
}

func convertVolumes(volumes map[string]string) []string {
	var result []string
	for hostPath, containerPath := range volumes {
		result = append(result, fmt.Sprintf("%s:%s", hostPath, containerPath))
	}
	return result
}
