package main

import (
	_ "github.com/joho/godotenv/autoload" // load env vars from .env file
	"github.com/kelseyhightower/envconfig"
)

// config represents application configuration.
type config struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN" required:"true"`
}

// readConfig returns configuration populated from environment variables.
func readConfig() (*config, error) {
	cfg := &config{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
