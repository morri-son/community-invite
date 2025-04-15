package render

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"os"
	"path/filepath"
	"strings"
	"text/template" // Use text/template instead of html/template
	"time"

	"github.com/community-invite/internal/config"
)

const templatesDir = "templates"

type TemplateData struct {
	Date   time.Time
	Agenda []config.AgendaItem
}

func GenerateFiles(cfg *config.Config, outputDir string) error {
	data := TemplateData{
		Date:   cfg.Date,
		Agenda: cfg.Agenda,
	}

	// Generate HTML
	htmlContent, err := renderTemplate("email-template.html", data)
	if err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "mail.html"), []byte(htmlContent), 0644); err != nil {
		return err
	}

	// Generate EML
	emlContent, err := generateEML(cfg, data)
	if err != nil {
		return fmt.Errorf("failed to generate EML: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "mail.eml"), []byte(emlContent), 0644); err != nil {
		return err
	}

	// Generate Slack
	slackContent, err := renderTemplate("slack-template.txt", data)
	if err != nil {
		return fmt.Errorf("failed to generate Slack: %w", err)
	}
	return os.WriteFile(filepath.Join(outputDir, "slack.md"), []byte(slackContent), 0644)
}

func generateEML(cfg *config.Config, data TemplateData) (string, error) {
	var recipients []string
	for _, target := range cfg.Targets {
		if target.Type == "email" && len(target.Recipients) > 0 {
			recipients = target.Recipients
			break
		}
	}
	if len(recipients) == 0 {
		return "", fmt.Errorf("no email recipients found")
	}

	htmlContent, err := renderTemplate("email-template.html", data)
	if err != nil {
		return "", err
	}

	// Encode HTML content using quoted-printable
	var encodedHTML bytes.Buffer
	writer := quotedprintable.NewWriter(&encodedHTML)
	if _, err := writer.Write([]byte(htmlContent)); err != nil {
		return "", fmt.Errorf("failed to encode HTML: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	return fmt.Sprintf(`From: %s
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: quoted-printable

%s`,
		cfg.Email.From,
		strings.Join(recipients, ", "),
		cfg.Email.Subject,
		encodedHTML.String(),
	), nil
}

func renderTemplate(templateName string, data TemplateData) (string, error) {
	tmplPath := filepath.Join(templatesDir, templateName)

	// Use text/template instead of html/template
	tmpl, err := template.New(filepath.Base(tmplPath)).
		ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("template parsing failed: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}

func HTMLBody(cfg *config.Config) (string, error) {
	data := TemplateData{
		Date:   cfg.Date,
		Agenda: cfg.Agenda,
	}
	return renderTemplate("email-template.html", data)
}

func SlackMessage(cfg *config.Config) (string, error) {
	data := TemplateData{
		Date:   cfg.Date,
		Agenda: cfg.Agenda,
	}
	return renderTemplate("slack-template.txt", data)
}
