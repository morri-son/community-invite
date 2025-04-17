package smtp

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/morri-son/community-invite/internal/config"
)

func SendBulkEmail(target config.Target, body string) error {
	auth := smtp.PlainAuth(
		"",
		target.From,
		os.Getenv("SMTP_PASSWORD"), // Updated to SMTP_PASSWORD
		target.SMTPHost,
	)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n%s",
		target.From,
		strings.Join(target.Recipients, ", "),
		target.Subject,
		body,
	)

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", target.SMTPHost, target.SMTPPort),
		auth,
		target.From,
		target.Recipients,
		[]byte(msg),
	)
}
