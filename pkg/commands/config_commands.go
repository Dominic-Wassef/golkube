package commands

import (
	"fmt"
	"log"
	"os"

	golConfig "golkube/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// registerConfigCommands registers all commands related to configuration management.
func RegisterConfigCommands() {
	// Parent command for configuration-related operations
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration files",
	}

	// Subcommand to show the current configuration
	configCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show the current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			configFile := viper.ConfigFileUsed()
			data, err := os.ReadFile(configFile)
			if err != nil {
				log.Fatalf("Error reading configuration file: %v", err)
			}
			fmt.Printf("Current Configuration:\n%s\n", string(data))
		},
	})

	// Subcommand to validate the current configuration
	configCmd.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			configFile := viper.ConfigFileUsed()
			config, err := golConfig.ParseConfig(configFile)
			if err != nil {
				log.Fatalf("Configuration validation failed: %v", err)
			}
			fmt.Printf("Configuration is valid:\n%+v\n", config)
		},
	})

	// Subcommand to save the current configuration to a file
	configCmd.AddCommand(&cobra.Command{
		Use:   "save <output-file>",
		Short: "Save the current configuration to a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			outputFile := args[0]
			configFile := viper.ConfigFileUsed()
			config, err := golConfig.ParseConfig(configFile)
			if err != nil {
				log.Fatalf("Error parsing configuration: %v", err)
			}
			err = golConfig.SaveConfig(outputFile, config)
			if err != nil {
				log.Fatalf("Error saving configuration: %v", err)
			}
			fmt.Printf("Configuration saved successfully to %s\n", outputFile)
		},
	})

	// Add the config command to the root command
	RootCmd.AddCommand(configCmd)
}
