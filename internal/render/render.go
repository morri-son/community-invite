package render

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/community-invite/internal/config"
)

// Template paths (relative to project root)
const (
	templatesDir  = "templates"
	emailTemplate = "email-template.html"
	slackTemplate = "slack-template.txt"
)

// TemplateData holds variables accessible in templates
type TemplateData struct {
	Date   time.Time
	Agenda []config.AgendaItem
}

// GenerateFiles creates all output files in the specified directory
func GenerateFiles(cfg *config.Config, outputDir string) error {
	data := TemplateData{
		Date:   cfg.Date,
		Agenda: cfg.Agenda,
	}

	// 1. Generate HTML email
	htmlPath := filepath.Join(outputDir, "mail.html")
	if err := renderTemplate(emailTemplate, htmlPath, data); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// 2. Generate EML file
	emlPath := filepath.Join(outputDir, "mail.eml")
	if err := generateEML(cfg, data, emlPath); err != nil {
		return fmt.Errorf("failed to generate EML: %w", err)
	}

	// 3. Generate Slack Markdown
	slackPath := filepath.Join(outputDir, "slack.md")
	if err := renderTemplate(slackTemplate, slackPath, data); err != nil {
		return fmt.Errorf("failed to generate Slack message: %w", err)
	}

	return nil
}

// generateEML creates an Outlook-compatible .eml file
func generateEML(cfg *config.Config, data TemplateData, outputPath string) error {
	// Render HTML body
	htmlBody, err := renderTemplateToString(emailTemplate, data)
	if err != nil {
		return err
	}

	// Construct RFC5322 message
	emlContent := fmt.Sprintf(`From: %s
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8

%s`,
		cfg.Email.From,
		cfg.Email.TestRecipient,
		cfg.Email.Subject,
		htmlBody,
	)

	return os.WriteFile(outputPath, []byte(emlContent), 0644)
}

// renderTemplate loads a template and writes rendered output to file
func renderTemplate(templateName, outputPath string, data TemplateData) error {
	// Get absolute template path
	tmplPath := filepath.Join(templatesDir, templateName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("template parsing failed (%s): %w", tmplPath, err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("file creation failed (%s): %w", outputPath, err)
	}
	defer outFile.Close()

	// Execute template
	if err := tmpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	return nil
}

// renderTemplateToString renders template to a string
func renderTemplateToString(templateName string, data TemplateData) (string, error) {
	tmplPath := filepath.Join(templatesDir, templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// HTMLBody renders just the email HTML content (for SMTP sending)
func HTMLBody(cfg *config.Config) (string, error) {
	data := TemplateData{
		Date:   cfg.Date,
		Agenda: cfg.Agenda,
	}
	return renderTemplateToString(emailTemplate, data)
}

// SlackMessage renders the Slack markdown content
func SlackMessage(cfg *config.Config) (string, error) {
	data := TemplateData{
		Date:   cfg.Date,
		Agenda: cfg.Agenda,
	}
	return renderTemplateToString(slackTemplate, data)
}
