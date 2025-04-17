package render

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/morri-son/community-invite/internal/config"
)

const templatesDir = "templates"

type TemplateData struct {
	Date    time.Time
	Agenda  []config.AgendaItem
	Subject string
	From    string
}

func GenerateFiles(cfg *config.Config, outputDir string) error {
	data := TemplateData{
		Date:   cfg.Date,
		Agenda: cfg.Agenda,
	}

	// Find email target for generation
	var emailTarget *config.Target
	for _, t := range cfg.Targets {
		if t.Type == "email" {
			emailTarget = &t
			break
		}
	}

	if emailTarget != nil {
		data.Subject = emailTarget.Subject
		data.From = emailTarget.From

		// Generate HTML
		htmlContent, err := renderTemplate("email-template.html", data)
		if err != nil {
			return fmt.Errorf("failed to generate HTML: %w", err)
		}
		if err := os.WriteFile(filepath.Join(outputDir, "mail.html"), []byte(htmlContent), 0644); err != nil {
			return err
		}

		// Generate EML
		emlContent, err := generateEML(*emailTarget, data)
		if err != nil {
			return fmt.Errorf("failed to generate EML: %w", err)
		}
		if err := os.WriteFile(filepath.Join(outputDir, "mail.eml"), []byte(emlContent), 0644); err != nil {
			return err
		}
	}

	// Generate Slack
	slackContent, err := renderTemplate("slack-template.md", data)
	if err != nil {
		return fmt.Errorf("failed to generate Slack: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "slack.md"), []byte(slackContent), 0644); err != nil {
		return err
	}

	return nil
}

func generateEML(target config.Target, data TemplateData) (string, error) {
	htmlContent, err := renderTemplate("email-template.html", data)
	if err != nil {
		return "", err
	}

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
		target.From,
		strings.Join(target.Recipients, ", "),
		target.Subject,
		encodedHTML.String(),
	), nil
}

func renderTemplate(templateName string, data TemplateData) (string, error) {
	tmplPath := filepath.Join(templatesDir, templateName)
	tmpl, err := template.New(filepath.Base(tmplPath)).ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("template parsing failed: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}

func HTMLBody(target config.Target, data TemplateData) (string, error) {
	return renderTemplate(target.Template, data)
}

func SlackMessage(target config.Target, data TemplateData) (string, error) {
	return renderTemplate(target.Template, data)
}
