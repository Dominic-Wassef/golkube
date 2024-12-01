package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// BuildConfig holds the configuration for building a Docker image
type BuildConfig struct {
	Tag         string
	ContextDir  string
	Dockerfile  string
	BuildArgs   map[string]*string
	Push        bool
	RegistryURL string
}

// BuildImage builds a Docker image based on the given configuration
func BuildImage(config BuildConfig) error {
	if config.Tag == "" {
		return fmt.Errorf("image tag cannot be empty")
	}
	if _, err := os.Stat(config.ContextDir); os.IsNotExist(err) {
		return fmt.Errorf("context directory does not exist: %s", config.ContextDir)
	}

	// Validate Registry URL
	if config.RegistryURL != "" {
		if _, err := url.ParseRequestURI(config.RegistryURL); err != nil {
			return fmt.Errorf("invalid registry URL: %w", err)
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	// Create build context
	buildContext, err := createBuildContext(config.ContextDir, config.Dockerfile)
	if err != nil {
		return fmt.Errorf("failed to create build context: %w", err)
	}

	// Build the image
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{config.Tag},
		Dockerfile: config.Dockerfile,
		BuildArgs:  config.BuildArgs,
		Remove:     true,
	}

	fmt.Printf("Building Docker image %s...\n", config.Tag)
	resp, err := cli.ImageBuild(context.Background(), buildContext, buildOptions)
	if err != nil {
		return fmt.Errorf("image build failed: %w", err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return fmt.Errorf("failed to display build output: %w", err)
	}

	fmt.Printf("Image %s built successfully\n", config.Tag)

	// Push the image if required
	if config.Push {
		dockerHubTag := fmt.Sprintf("%s/%s", os.Getenv("DOCKER_USERNAME"), config.Tag)
		if err := TagImage(cli, config.Tag, dockerHubTag); err != nil {
			return fmt.Errorf("failed to tag image: %w", err)
		}

		if err := PushImage(cli, dockerHubTag, config.RegistryURL); err != nil {
			return fmt.Errorf("failed to push image: %w", err)
		}
	}

	return nil
}

// createBuildContext creates a tarball of the Docker build context
func createBuildContext(contextDir, dockerfile string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Ensure the Dockerfile is included in the tarball, even if it's not in the root of the context
	dockerfilePath := filepath.Join(contextDir, dockerfile)
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("dockerfile does not exist: %s", dockerfilePath)
	}

	err := filepath.Walk(contextDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(contextDir, file)
		if err != nil {
			return err
		}

		// Explicitly set the Dockerfile path to be "Dockerfile" if it matches the provided file
		if file == dockerfilePath {
			relPath = "Dockerfile"
		}

		header, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			if _, err := tw.Write(data); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}

// TagImage tags a Docker image
func TagImage(cli *client.Client, sourceTag, targetTag string) error {
	fmt.Printf("Tagging image %s as %s...\n", sourceTag, targetTag)
	return cli.ImageTag(context.Background(), sourceTag, targetTag)
}

// PushImage pushes a Docker image to the specified registry
func PushImage(cli *client.Client, tag, registryURL string) error {
	authConfig := types.AuthConfig{
		ServerAddress: registryURL,
		Username:      os.Getenv("DOCKER_USERNAME"),
		Password:      os.Getenv("DOCKER_PASSWORD"),
	}

	if authConfig.Username == "" || authConfig.Password == "" {
		return fmt.Errorf("docker credentials not provided; set DOCKER_USERNAME and DOCKER_PASSWORD")
	}

	encodedAuth, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("failed to encode authentication: %w", err)
	}

	fmt.Printf("Pushing Docker image %s...\n", tag)
	resp, err := cli.ImagePush(context.Background(), tag, types.ImagePushOptions{
		RegistryAuth: string(encodedAuth),
	})
	if err != nil {
		return fmt.Errorf("image push failed: %w", err)
	}
	defer resp.Close()

	if _, err := io.Copy(os.Stdout, resp); err != nil {
		return fmt.Errorf("failed to display push output: %w", err)
	}

	fmt.Printf("Image %s pushed successfully to %s\n", tag, registryURL)
	return nil
}
