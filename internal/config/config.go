package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Date    time.Time    `mapstructure:"date"`
	Agenda  []AgendaItem `mapstructure:"agenda"`
	Targets []Target     `mapstructure:"targets"`
}

type AgendaItem struct {
	Type      string `mapstructure:"type"`
	Title     string `mapstructure:"title"`
	Presenter string `mapstructure:"presenter"`
}

type Target struct {
	Type string `mapstructure:"type"`
	// Email fields
	Subject    string   `mapstructure:"subject,omitempty"`
	From       string   `mapstructure:"from,omitempty"`
	SMTPHost   string   `mapstructure:"smtp_host,omitempty"`
	SMTPPort   int      `mapstructure:"smtp_port,omitempty"`
	Recipients []string `mapstructure:"recipients,omitempty"`
	// Slack fields
	ClientID  string `mapstructure:"client_id,omitempty"`
	ChannelID string `mapstructure:"channel_id,omitempty"`
	Workspace string `mapstructure:"workspace,omitempty"`
	// Common
	Template string `mapstructure:"template"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed reading config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal failed: %w", err)
	}

	if cfg.Date.IsZero() {
		return nil, fmt.Errorf("invalid date in config")
	}

	return &cfg, nil
}
