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
	var (
		dryRun    bool
		sendAll   bool
		sendMail  bool
		sendSlack bool
	)

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send invitations to specified targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			targets := filterTargets(cfg.Targets, sendAll, sendMail, sendSlack)
			if len(targets) == 0 {
				return fmt.Errorf("no matching targets found")
			}

			baseData := render.TemplateData{
				Date:   cfg.Date,
				Agenda: cfg.Agenda,
			}

			for _, target := range targets {
				data := baseData
				data.Subject = target.Subject
				data.From = target.From

				switch target.Type {
				case "email":
					if err := handleEmailTarget(target, data, dryRun); err != nil {
						return err
					}
				case "slack":
					if err := handleSlackTarget(target, data, dryRun); err != nil {
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
	cmd.Flags().BoolVar(&sendAll, "all", false, "Send to all configured targets")
	cmd.Flags().BoolVar(&sendMail, "mail", false, "Send to email targets only")
	cmd.Flags().BoolVar(&sendSlack, "slack", false, "Send to Slack targets only")

	return cmd
}

func filterTargets(targets []config.Target, all, mail, slack bool) []config.Target {
	var result []config.Target
	for _, target := range targets {
		if all || (mail && target.Type == "email") || (slack && target.Type == "slack") {
			result = append(result, target)
		}
	}
	return result
}

func handleEmailTarget(target config.Target, data render.TemplateData, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY-RUN] Would send email to: %v\n", target.Recipients)
		return nil
	}

	body, err := render.HTMLBody(target, data)
	if err != nil {
		return err
	}

	return smtp.SendBulkEmail(target, body)
}

func handleSlackTarget(target config.Target, data render.TemplateData, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY-RUN] Would post to Slack channel %s\n", target.ChannelID)
		return nil
	}

	message, err := render.SlackMessage(target, data)
	if err != nil {
		return err
	}

	client := slack.NewClient(
		target.ClientID,
		target.ChannelID,
		target.Workspace,
	)
	return client.SendMessage(message)
}
