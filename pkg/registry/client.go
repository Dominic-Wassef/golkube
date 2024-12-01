package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// RegistryClient provides methods to interact with Docker registries
type RegistryClient struct {
	Client *client.Client
}

// NewRegistryClient initializes a new Docker RegistryClient
func NewRegistryClient() (*RegistryClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	return &RegistryClient{Client: cli}, nil
}

// PushImage pushes a Docker image to the specified registry
func (rc *RegistryClient) PushImage(imageTag, registryURL, username, password string) error {
	authConfig := types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: registryURL,
	}
	encodedAuth, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("failed to encode authentication: %w", err)
	}

	resp, err := rc.Client.ImagePush(context.Background(), imageTag, types.ImagePushOptions{
		RegistryAuth: string(encodedAuth),
	})
	if err != nil {
		return fmt.Errorf("image push failed: %w", err)
	}
	defer resp.Close()

	if _, err := io.Copy(io.Discard, resp); err != nil {
		return fmt.Errorf("failed to read push response: %w", err)
	}

	fmt.Printf("Image %s pushed successfully to %s\n", imageTag, registryURL)
	return nil
}

// PullImage pulls a Docker image from the specified registry
func (rc *RegistryClient) PullImage(imageTag, registryURL, username, password string) error {
	authConfig := types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: registryURL,
	}
	encodedAuth, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("failed to encode authentication: %w", err)
	}

	resp, err := rc.Client.ImagePull(context.Background(), imageTag, types.ImagePullOptions{
		RegistryAuth: string(encodedAuth),
	})
	if err != nil {
		return fmt.Errorf("image pull failed: %w", err)
	}
	defer resp.Close()

	if _, err := io.Copy(io.Discard, resp); err != nil {
		return fmt.Errorf("failed to read pull response: %w", err)
	}

	fmt.Printf("Image %s pulled successfully from %s\n", imageTag, registryURL)
	return nil
}

// ListTags lists all tags for an image from a Docker registry
func (rc *RegistryClient) ListTags(repository, registryURL, username, password string) ([]string, error) {
	authConfig := types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: registryURL,
	}
	encodedAuth, err := json.Marshal(authConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode authentication: %w", err)
	}

	authHeader := fmt.Sprintf("Basic %s", encodedAuth)

	url := fmt.Sprintf("%s/v2/%s/tags/list", registryURL, repository)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tagsResp struct {
		Tags []string `json:"tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return nil, fmt.Errorf("failed to decode tags response: %w", err)
	}

	return tagsResp.Tags, nil
}
