package smtp

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/community-invite/internal/config"
)

func SendTestEmail(cfg *config.Config, body string) error {
	password := os.Getenv("SMTP_PASSWORD")
	if password == "" {
		password = cfg.Email.Password
	}
	if password == "" {
		return fmt.Errorf("SMTP password required - set SMTP_PASSWORD env var")
	}

	auth := smtp.PlainAuth("", cfg.Email.From, password, cfg.Email.SMTPHost)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n%s",
		cfg.Email.TestRecipient,
		cfg.Email.Subject,
		body,
	))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", cfg.Email.SMTPHost, cfg.Email.SMTPPort),
		auth,
		cfg.Email.From,
		[]string{cfg.Email.TestRecipient},
		msg,
	)
}

func SendBulkEmail(cfg *config.Config, recipients []string, body string) error {
	password := os.Getenv("SMTP_PASSWORD")
	if password == "" {
		password = cfg.Email.Password
	}
	if password == "" {
		return fmt.Errorf("SMTP password required")
	}

	auth := smtp.PlainAuth("", cfg.Email.From, password, cfg.Email.SMTPHost)

	msg := []byte(fmt.Sprintf("Subject: %s\r\n"+
		"From: %s\r\n"+
		"Cc: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n%s",
		cfg.Email.Subject,
		cfg.Email.From,
		cfg.Email.From, // CC field
		body,
	))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", cfg.Email.SMTPHost, cfg.Email.SMTPPort),
		auth,
		cfg.Email.From,
		append(recipients, cfg.Email.From), // Include CC
		msg,
	)
}

func validateSMTPConfig(cfg *config.Config) error {
	if cfg.Email.SMTPPort != 587 {
		return fmt.Errorf("only port 587 with STARTTLS is supported")
	}
	// Add certificate validation if needed
	return nil
}
