package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config represents application configuration.
type Config struct {
	TelegramToken  string        `yaml:"telegram_token"`
	TelegramChat   string        `yaml:"telegram_chat"`
	UpdateInterval time.Duration `yaml:"update_interval"`
	DataFile       string        `yaml:"data_file"`
	Feeds          []string      `yaml:"feeds"`

	// Debug sets the log level and prints RSS items instead of sending
	// them to Telegram.
	Debug bool `yaml:"debug"`
}

const (
	defaultUpdateInterval = 1 * time.Hour
	defaultDataFile       = "./data.yaml"
)

// ReadConfig returns configuration populated from the config file.
func ReadConfig(file string) (Config, error) {
	data, err := os.ReadFile(file) //nolint:gosec
	if err != nil {
		return Config{}, fmt.Errorf("read file: %w", err)
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return Config{}, fmt.Errorf("unmarshal file: %w", err)
	}

	if !conf.Debug {
		if conf.TelegramToken == "" {
			return Config{}, errors.New("empty telegram token")
		}
		if conf.TelegramChat == "" {
			return Config{}, errors.New("empty telegram chat")
		}
	}
	if conf.UpdateInterval == 0 {
		conf.UpdateInterval = defaultUpdateInterval
	}
	if conf.DataFile == "" {
		conf.DataFile = defaultDataFile
	}
	if len(conf.Feeds) == 0 {
		return Config{}, errors.New("empty feeds list")
	}

	return conf, nil
}
