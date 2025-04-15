package main

import (
	"fmt"
	"os"

	"github.com/morri-son/community-invite/internal/config"
	"github.com/morri-son/community-invite/internal/render"
	"github.com/morri-son/community-invite/internal/slack"
	"github.com/morri-son/community-invite/internal/smtp"
	"github.com/spf13/cobra"
)

// Package-level variable for config file path
var (
	cfgFile string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "community-invite",
		Short: "Generate and send community meeting invitations",
		Long:  "CLI tool for managing OCM community meeting communications",
	}

	// Initialize the persistent flag
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "Config file (default is config.yaml)")

	rootCmd.AddCommand(NewGenerateCmd())
	rootCmd.AddCommand(NewTestmailCmd())
	rootCmd.AddCommand(NewSendCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewGenerateCmd() *cobra.Command {
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

			fmt.Println("Files generated successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFolder, "output-folder", "o", "/tmp", "Output directory")
	return cmd
}

func NewTestmailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testmail",
		Short: "Send test email to verification recipient",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			body, err := render.HTMLBody(cfg)
			if err != nil {
				return err
			}

			if err := smtp.SendTestEmail(cfg, body); err != nil {
				return fmt.Errorf("SMTP error: %w", err)
			}

			fmt.Println("Test email sent successfully")
			return nil
		},
	}
	return cmd
}

func NewSendCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send invitations to all targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			for _, target := range cfg.Targets {
				switch target.Type {
				case "email":
					if err := handleEmailTarget(cfg, target, dryRun); err != nil {
						return err
					}
				case "slack":
					if err := handleSlackTarget(cfg, target, dryRun); err != nil {
						return err
					}
				default:
					return fmt.Errorf("unknown target type: %s", target.Type)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate sending without actual delivery")
	return cmd
}

func handleEmailTarget(cfg *config.Config, target config.Target, dryRun bool) error {
	body, err := render.HTMLBody(cfg)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Printf("[DRY-RUN] Would send email to: %v\n", target.Recipients)
		return nil
	}

	return smtp.SendBulkEmail(cfg, target.Recipients, body)
}

func handleSlackTarget(cfg *config.Config, target config.Target, dryRun bool) error {
	message, err := render.SlackMessage(cfg)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Printf("[DRY-RUN] Would post to Slack channel %s\n", target.ChannelID)
		return nil
	}

	// Added workspace parameter from target configuration
	return slack.SendMessage(target.ChannelID, message, target.Workspace)
}
