package cmd

import (
	"fmt"

	"github.com/morri-son/community-invite/internal/config"
	"github.com/morri-son/community-invite/internal/render"
	"github.com/morri-son/community-invite/internal/smtp"
	"github.com/spf13/cobra"
)

func NewTestmailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testmail",
		Short: "Send test email to verification recipient",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			// Find first email target
			var emailTarget *config.Target
			for _, t := range cfg.Targets {
				if t.Type == "email" {
					emailTarget = &t
					break
				}
			}
			if emailTarget == nil {
				return fmt.Errorf("no email target found in config")
			}

			data := render.TemplateData{
				Date:    cfg.Date,
				Agenda:  cfg.Agenda,
				Subject: emailTarget.Subject,
				From:    emailTarget.From,
			}

			body, err := render.HTMLBody(*emailTarget, data)
			if err != nil {
				return err
			}

			// Use email target's From address as recipient
			testTarget := *emailTarget
			testTarget.Recipients = []string{emailTarget.From}

			if err := smtp.SendBulkEmail(testTarget, body); err != nil {
				return fmt.Errorf("SMTP error: %w", err)
			}

			fmt.Println("Test email sent successfully")
			return nil
		},
	}
	return cmd
}
