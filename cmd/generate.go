package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func GenerateCmd() *cobra.Command {
	var outputFolder string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate invitation files",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			cfg, err := LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			// Validate output folder
			if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
				return fmt.Errorf("output folder does not exist: %s", outputFolder)
			}

			// Render and save files
			if err := GenerateFiles(cfg, outputFolder); err != nil {
				return fmt.Errorf("generation failed: %w", err)
			}

			fmt.Println("Files generated successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFolder, "output-folder", "o", "/tmp", "Output directory")
	return cmd
}
