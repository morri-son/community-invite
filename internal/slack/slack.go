package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	ClientID     string
	ClientSecret string
	ChannelID    string
	Workspace    string
	HTTPClient   *http.Client
}

func NewClient(clientID, channelID, workspace string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
		ChannelID:    channelID,
		Workspace:    workspace,
		HTTPClient:   &http.Client{},
	}
}

func (c *Client) SendMessage(message string) error {
	// First get OAuth token
	token, err := c.getOAuthToken()
	if err != nil {
		return err
	}

	payload := map[string]string{
		"channel": c.ChannelID,
		"text":    message,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(
		"POST",
		"https://slack.com/api/chat.postMessage",
		bytes.NewBuffer(jsonData),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("slack API error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API error: %s", resp.Status)
	}

	return nil
}

func (c *Client) getOAuthToken() (string, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := http.PostForm("https://slack.com/api/oauth.v2.access", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("slack OAuth error: %s", result.Error)
	}

	return result.AccessToken, nil
}
