package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// Config represents the structure of the configuration file
type Config struct {
	Namespace     string                 `yaml:"namespace" json:"namespace" validate:"required"`
	Image         string                 `yaml:"image" json:"image" validate:"required"`
	Replicas      int                    `yaml:"replicas" json:"replicas" validate:"gte=1"`
	Ports         []int                  `yaml:"ports" json:"ports" validate:"dive,gt=0"`
	Env           map[string]string      `yaml:"env" json:"env"`
	Notifications map[string]interface{} `yaml:"notifications" json:"notifications"`
}

// ParseConfig parses a YAML/JSON configuration file into a Config struct
func ParseConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// validateConfig validates the configuration using struct tags and custom rules
func validateConfig(config *Config) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

// SaveConfig saves the Config struct back to a YAML file
func SaveConfig(filePath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}
