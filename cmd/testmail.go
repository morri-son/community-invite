// cmd/testmail.go
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/morri-son/community-invite/internal/config"
	"github.com/morri-son/community-invite/internal/render"
	"github.com/morri-son/community-invite/internal/smtp"
	"github.com/spf13/cobra"
)

func NewTestmailCmd() *cobra.Command {
	var (
		toEmail string
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   "testmail",
		Short: "Send test email to verification recipient",
		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()

			if verbose {
				fmt.Fprintf(os.Stderr, "Starting test email procedure at %s\n", startTime.Format("15:04:05"))
				fmt.Fprintf(os.Stderr, "Using config file: %s\n", cfgFile)
			}

			// Load configuration
			if verbose {
				fmt.Fprintln(os.Stderr, "Loading configuration...")
			}
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			// Override recipient if specified
			if toEmail != "" {
				if verbose {
					fmt.Fprintf(os.Stderr, "Overriding test recipient from '%s' to '%s'\n",
						cfg.Email.TestRecipient,
						toEmail)
				}
				cfg.Email.TestRecipient = toEmail
			}

			// Render HTML body
			if verbose {
				fmt.Fprintln(os.Stderr, "Rendering email template...")
			}
			body, err := render.HTMLBody(cfg)
			if err != nil {
				return fmt.Errorf("rendering failed: %w", err)
			}
			if verbose {
				fmt.Fprintln(os.Stderr, "Template rendered successfully")
				fmt.Fprintf(os.Stderr, "Email subject: %s\n", cfg.Email.Subject)
				fmt.Fprintf(os.Stderr, "Body size: %d bytes\n", len(body))
			}

			// Send email
			if verbose {
				fmt.Fprintln(os.Stderr, "Initializing SMTP client...")
				fmt.Fprintf(os.Stderr, "SMTP Server: %s:%d\n", cfg.Email.SMTPHost, cfg.Email.SMTPPort)
				fmt.Fprintf(os.Stderr, "Sending test email to: %s\n", cfg.Email.TestRecipient)
			}

			if err := smtp.SendTestEmail(cfg, body); err != nil {
				return fmt.Errorf("SMTP error: %w", err)
			}

			if verbose {
				duration := time.Since(startTime)
				fmt.Fprintf(os.Stderr, "Procedure completed in %v\n", duration.Round(time.Millisecond))
			}

			fmt.Printf("Test email sent successfully to %s\n", cfg.Email.TestRecipient)
			return nil
		},
	}

	cmd.Flags().StringVarP(&toEmail, "to", "t", "", "Override test recipient email address")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	return cmd
}
