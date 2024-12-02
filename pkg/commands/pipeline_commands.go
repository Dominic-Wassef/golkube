package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// registerPipelineCommand registers all pipeline-related commands.
func RegisterPipelineCommand() {
	// Parent command for pipeline operations
	pipelineCmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Manage and execute pipeline stages",
	}

	// Subcommand to execute a pipeline
	pipelineExecuteCmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute the pipeline as defined in the configuration",
		Run: func(cmd *cobra.Command, args []string) {
			pipelineFile := viper.GetString("pipeline.config_file")
			if pipelineFile == "" {
				log.Fatalf("Pipeline configuration file is not specified in the configuration")
			}

			// Read and parse the pipeline file
			config, err := os.ReadFile(pipelineFile)
			if err != nil {
				log.Fatalf("Failed to read pipeline file: %v", err)
			}

			var pipeline struct {
				Stages []struct {
					Name     string   `json:"name"`
					Commands []string `json:"commands"`
				} `json:"stages"`
			}

			if err := yaml.Unmarshal(config, &pipeline); err != nil {
				log.Fatalf("Failed to parse pipeline file: %v", err)
			}

			// Execute each stage in the pipeline
			for _, stage := range pipeline.Stages {
				fmt.Printf("Executing stage: %s\n", stage.Name)
				for _, command := range stage.Commands {
					fmt.Printf("Running command: %s\n", command)
					cmd := exec.Command("sh", "-c", command)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatalf("Command failed: %s\nError: %v", command, err)
					}
				}
				fmt.Printf("Stage %s completed successfully.\n", stage.Name)
			}
			fmt.Println("Pipeline execution completed successfully.")
		},
	}

	// Add the execute subcommand to the pipeline command
	pipelineCmd.AddCommand(pipelineExecuteCmd)

	// Add the pipeline command to the root command
	RootCmd.AddCommand(pipelineCmd)
}
