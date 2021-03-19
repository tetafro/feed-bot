package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
)

// config represents application configuration.
type config struct {
	TelegramToken  string        `json:"telegram_token"`
	UpdateInterval time.Duration `json:"update_interval"`
	DataFile       string        `json:"data_file"`
	Feeds          []string      `json:"feeds"`
}

const (
	defaultUpdateInterval = 1 * time.Hour
	defaultDataFile       = "./data.json"
)

func (c *config) UnmarshalJSON(data []byte) error {
	type alias config
	a := &struct {
		UpdateInterval string `json:"update_interval"`
		*alias
	}{
		alias: (*alias)(c),
	}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	var err error
	c.UpdateInterval, err = time.ParseDuration(a.UpdateInterval)
	if err != nil {
		return err
	}

	return nil
}

// readConfig returns configuration populated from environment variables.
func readConfig(file string) (config, error) {
	data, err := ioutil.ReadFile(file) // nolint: gosec
	if err != nil {
		return config{}, errors.Wrap(err, "read file")
	}

	cfg := config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return config{}, errors.Wrap(err, "unmarshal file")
	}

	if cfg.TelegramToken == "" {
		return config{}, errors.New("empty token")
	}
	if cfg.UpdateInterval == 0 {
		cfg.UpdateInterval = defaultUpdateInterval
	}
	if cfg.DataFile == "" {
		cfg.DataFile = defaultDataFile
	}

	return cfg, nil
}
