package cmd

import (
	"fmt"

	"github.com/morri-son/community-invite/internal/config"
	"github.com/morri-son/community-invite/internal/render"
	"github.com/morri-son/community-invite/internal/smtp"
	"github.com/spf13/cobra"
)

func NewTestmailCmd() *cobra.Command {
	var toEmail string

	cmd := &cobra.Command{
		Use:   "testmail",
		Short: "Send test email to verification recipient",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			// Override recipient if flag is set
			if toEmail != "" {
				cfg.Email.TestRecipient = toEmail
			}

			body, err := render.HTMLBody(cfg)
			if err != nil {
				return err
			}

			if err := smtp.SendTestEmail(cfg, body); err != nil {
				return fmt.Errorf("SMTP error: %w", err)
			}

			fmt.Printf("Test email sent successfully to %s\n", cfg.Email.TestRecipient)
			return nil
		},
	}

	cmd.Flags().StringVarP(&toEmail, "to", "t", "", "Override test recipient email address")
	return cmd
}
