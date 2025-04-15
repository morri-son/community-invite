package main

import (
	"fmt"

	"github.com/community-invite/internal/config"
	"github.com/community-invite/internal/render"
	"github.com/community-invite/internal/slack"
	"github.com/community-invite/internal/smtp"
	"github.com/spf13/cobra"
)

var cfgFile string

func SendCmd() *cobra.Command {
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
	cmd.Flags().StringVar(&cfgFile, "config", "config.yaml", "Path to the configuration file")
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

	return slack.SendMessage(target.ChannelID, message)
}
