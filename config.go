package main

import (
	_ "github.com/joho/godotenv/autoload" // load env vars from .env file
	"github.com/kelseyhightower/envconfig"
)

// config represents application configuration.
type config struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN" required:"true"`
}

// mustConfig returns configuration populated from environment variables.
func mustConfig() *config {
	cfg := &config{}
	envconfig.MustProcess("", cfg)
	return cfg
}
