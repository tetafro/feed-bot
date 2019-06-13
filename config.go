package main

import (
	"time"

	_ "github.com/joho/godotenv/autoload" // load env vars from .env file
	"github.com/kelseyhightower/envconfig"
)

// config represents application configuration.
type config struct {
	TelegramToken  string        `envconfig:"TELEGRAM_TOKEN" required:"true"`
	Feeds          []string      `envconfig:"FEEDS" required:"true"`
	UpdateInterval time.Duration `envconfig:"UPDATE_INTERVAL" default:"1h"`
}

// readConfig returns configuration populated from environment variables.
func readConfig() (*config, error) {
	cfg := &config{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
