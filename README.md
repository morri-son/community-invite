# Community Call Invitation System

A secure CLI tool for managing community meeting communications.

## Features

- 📄 Generate HTML/EML/Slack Markdown from templates
- 📨 Test email delivery with SMTP
- 📦 Bulk send to recipients and Slack channels
- ⚙️ Configurable via YAML file
- 🔬 Dry-run mode for safety

## Prerequisites

- Go 1.24+ installed
- Google Account (for email sending)
- **Slack workspace with bot token** (xoxb-*)
- Basic terminal/command prompt skills
- Template Files for email and Slack messages in the ./templates directory named
  - email-template.html
  - slack-template.txt

## Project structure

```bash
community-invite/                # Root directory for the community invitation CLI tool
│
├── templates/                  # Directory containing all message templates
│   │                           # Templates use Go's templating syntax with {{.Variable}} placeholders
│   ├── email-template.html     # HTML template for email invitations (contains styling and layout)
│   └── slack-template.txt      # Markdown-formatted template for Slack messages (supports Slack formatting)
│
├── cmd/                       
│   ├── root.go                 # Root command and shared configuration
│   │                           # Each file corresponds to a subcommand (generate, testmail, send)
│   ├── generate.go             # Implements the 'generate' command for creating local template files
│   ├── testmail.go             # Implements the 'testmail' command for sending test emails
│   └── send.go                 # Implements the 'send' command for production distribution
│
├── internal/                   # Internal implementation packages (not exposed externally)
│   │
│   ├── config/                 # Configuration loading and validation
│   │   └── config.go           # Defines Config struct and YAML parsing logic
│   │
│   ├── smtp/                   # SMTP email sending functionality
│   │   └── smtp.go             # Handles secure email delivery with STARTTLS
│   │
│   ├── slack/                  # Slack API integration
│   │   └── slack.go            # Manages Slack message posting via Web API
│   │
│   └── render/                 # Template rendering engine
│       └── render.go           # Processes templates with dynamic data injection
│
├── go.mod                      # Go module definition (lists dependencies and module path)
├── go.sum                      # Go checksum file (ensures dependency integrity)
├── config.yaml                 # Example configuration file (user-editable settings)
└── main.go                     # Application entry point (CLI setup and command routing)
```

## Installation

```bash
go install github.com/morri-son/community-invite@latest
````

## Usage

```bash
# Generate files
./community-invite generate --output-folder ./out

# Send test email
SMTP_PASSWORD="secret" ./community-invite testmail

# Send test email with override recipient and verbose logging
SMTP_PASSWORD="secret" ./community-invite testmail --to dev@example.com --verbose

# Send to all targets (dry-run)
SLACK_API_TOKEN="xoxb-123" ./community-invite send --dry-run
```

## Configuration (default `config.yaml`)

Place the file in either:

- The working directory as config.yaml, or
- Specify path via --config flag:

```bash
./community-invite generate --config /path/to/config.yaml
```

Passwords for Slack API and SMTP SHOULD be configured as environment variables.
But it is also possible to set them in the config.yaml file.

Below is an example of a configuration file:

```yaml
date: 2024-04-09T15:00:00Z # Must be present and in RFC3339 format

# Agenda Items for the meeting. Requires type, title, and presenter.
agenda:
  - type: "Dev Update"
    title: "Project roadmap review"
    presenter: "Alice"
  - type: "Technical Discussion"
    title: "CI/CD pipeline improvements"
    presenter: "Bob"
  - type: "Open Format"
    title: "Open discussion"
    presenter: "All"

# Configuration for SMTP server
email:
  subject: "OCM Community Call Invitation"
  from: "opencomponentmodel@gmail.com"
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  password: "your-gmail-app-password" # SHOULD be omitted. Instead set SMTP_PASSWORD in the environment.
 

# Configuration for message targets
targets:
  - type: email
    recipients: ["ocm-announce@googlegroups.com"]
    template: "templates/email-template.html"
 - type: slack
    channel: "#ocm-announce"
    template: "templates/slack-template.txt"
    channelID: "C1234567890"
    workspace: "kubernetes.slack.com"
    api_token: "xoxb-1234567890-1234567890123-1234567890123-abc123" # SHOULD be omitted. Instead set SMTP_PASSWORD in the environment.
```

### Validation Rules

- date: Must be RFC3339 format
- agenda: Minimum 1 item required
- `email.from/test_recipient`: Valid email format
- targets: At least one email or slack target required

### Security Best Practices

- NEVER commit secrets to config.yaml
- Use SMTP_PASSWORD and SLACK_API_TOKEN environment variables
- Regularly rotate credentials

---

## Testing instructions

**Generate Files:**

```bash
./community-invite generate --output-folder ./test-out
```

**Test Email:**

```bash
SMTP_PASSWORD="yourpass" ./community-invite testmail
```

**Dry-Run Send:**

```bash
SLACK_API_TOKEN="xoxb-123" ./community-invite send --dry-run
```
