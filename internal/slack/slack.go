package slack

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func SendMessage(channelID string, message string, workspaceURL string) error {
	token := os.Getenv("SLACK_API_TOKEN")
	if token == "" {
		return fmt.Errorf("SLACK_API_TOKEN environment variable required")
	}

	// Initialize Slack client with optional custom workspace URL
	var options []slack.Option
	if workspaceURL != "" {
		options = append(options, slack.OptionAPIURL(workspaceURL+"/api/")) // Append /api/ for Slack compatibility
	}

	api := slack.New(token, options...)
	_, _, err := api.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
	)
	return err
}
