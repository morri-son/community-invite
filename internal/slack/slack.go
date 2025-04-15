package slack

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func SendMessage(channelID string, message string) error {
	token := os.Getenv("SLACK_API_TOKEN")
	if token == "" {
		return fmt.Errorf("SLACK_API_TOKEN environment variable required")
	}

	api := slack.New(token)
	_, _, err := api.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
	)
	return err
}
