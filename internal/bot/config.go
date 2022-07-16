package bot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Config represents application configuration.
type Config struct {
	TelegramToken  string        `yaml:"telegram_token"`
	TelegramChat   int64         `yaml:"telegram_chat"`
	UpdateInterval time.Duration `yaml:"update_interval"`
	DataFile       string        `yaml:"data_file"`
	Feeds          []string      `yaml:"feeds"`
	LogNotifier    bool          `yaml:"log_notifier"`
}

const (
	defaultUpdateInterval = 1 * time.Hour
	defaultDataFile       = "./data.yaml"
)

// ReadConfig returns configuration populated from environment variables.
func ReadConfig(file string) (Config, error) {
	data, err := ioutil.ReadFile(file) // nolint: gosec
	if err != nil {
		return Config{}, fmt.Errorf("read file: %w", err)
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return Config{}, fmt.Errorf("unmarshal file: %w", err)
	}

	if !conf.LogNotifier && conf.TelegramToken == "" {
		return Config{}, errors.New("empty token")
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
