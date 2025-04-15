package main

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/morri-son/community-invite/internal/config"
	"github.com/morri-son/community-invite/internal/render"
	"github.com/spf13/cobra"
)

func GenerateCmd() *cobra.Command {
	var outputFolder string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate invitation files",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
				return fmt.Errorf("output folder does not exist: %s", outputFolder)
			}

			if err := render.GenerateFiles(cfg, outputFolder); err != nil {
				return fmt.Errorf("generation failed: %w", err)
			}

			fmt.Printf("Files generated successfully in %s:\n- %s\n- %s\n- %s\n",
				outputFolder,
				filepath.Join(outputFolder, "mail.html"),
				filepath.Join(outputFolder, "mail.eml"),
				filepath.Join(outputFolder, "slack.md"))
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFolder, "output-folder", "o", "/tmp", "Output directory")
	return cmd
}
