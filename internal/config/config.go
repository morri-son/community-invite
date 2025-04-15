package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Date    time.Time    `mapstructure:"date"`
	Agenda  []AgendaItem `mapstructure:"agenda"`
	Email   EmailConfig  `mapstructure:"email"`
	Targets []Target     `mapstructure:"targets"`
}

type AgendaItem struct {
	Type      string `mapstructure:"type"`
	Title     string `mapstructure:"title"`
	Presenter string `mapstructure:"presenter"`
}

type EmailConfig struct {
	Subject       string `mapstructure:"subject"`
	From          string `mapstructure:"from"`
	SMTPHost      string `mapstructure:"smtp_host"`
	SMTPPort      int    `mapstructure:"smtp_port"`
	TestRecipient string `mapstructure:"test_recipient"`
	Password      string `mapstructure:"password"` // Only for fallback
}

type Target struct {
	Type       string   `mapstructure:"type"`
	Recipients []string `mapstructure:"recipients"`
	ChannelID  string   `mapstructure:"channelID"`
	Template   string   `mapstructure:"template"`
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

	if cfg.Email.Password != "" {
		fmt.Fprintln(os.Stderr, "WARNING: Storing passwords in config.yaml is insecure! Use SMTP_PASSWORD environment variable instead.")
	}

	return &cfg, nil
}
