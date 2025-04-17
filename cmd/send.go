package cmd

import (
	"fmt"

	"github.com/morri-son/community-invite/internal/config"
	"github.com/morri-son/community-invite/internal/render"
	"github.com/morri-son/community-invite/internal/slack"
	"github.com/morri-son/community-invite/internal/smtp"
	"github.com/spf13/cobra"
)

func NewSendCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send generated invitations to all configured targets",
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

	return slack.SendMessage(target.ChannelID, message, target.Workspace)
}
