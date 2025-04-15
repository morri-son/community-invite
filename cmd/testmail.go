package main

import (
	"fmt"

	"github.com/community-invite/internal/config"
	"github.com/community-invite/internal/render"
	"github.com/community-invite/internal/smtp"
	"github.com/spf13/cobra"
)

func TestmailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testmail",
		Short: "Send test email to verification recipient",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			// Render HTML body
			body, err := render.HTMLBody(cfg)
			if err != nil {
				return err
			}

			// Send via SMTP
			if err := smtp.SendTestEmail(cfg, body); err != nil {
				return fmt.Errorf("SMTP error: %w", err)
			}

			fmt.Println("Test email sent successfully")
			return nil
		},
	}
	return cmd
}
