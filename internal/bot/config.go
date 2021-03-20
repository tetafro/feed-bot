package bot

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
)

// Config represents application configuration.
type Config struct {
	TelegramToken  string        `json:"telegram_token"`
	UpdateInterval time.Duration `json:"update_interval"`
	DataFile       string        `json:"data_file"`
	Feeds          []string      `json:"feeds"`
}

const (
	defaultUpdateInterval = 1 * time.Hour
	defaultDataFile       = "./data.json"
)

// UnmarshalJSON unmarshals config using additional data type `time.Duration`.
func (c *Config) UnmarshalJSON(data []byte) error {
	type alias Config
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

// ReadConfig returns configuration populated from environment variables.
func ReadConfig(file string) (Config, error) {
	data, err := ioutil.ReadFile(file) // nolint: gosec
	if err != nil {
		return Config{}, errors.Wrap(err, "read file")
	}

	cfg := Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, errors.Wrap(err, "unmarshal file")
	}

	if cfg.TelegramToken == "" {
		return Config{}, errors.New("empty token")
	}
	if cfg.UpdateInterval == 0 {
		cfg.UpdateInterval = defaultUpdateInterval
	}
	if cfg.DataFile == "" {
		cfg.DataFile = defaultDataFile
	}

	return cfg, nil
}
