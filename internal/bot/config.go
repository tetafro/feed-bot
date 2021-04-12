package bot

import (
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config represents application configuration.
type Config struct {
	TelegramToken  string        `yaml:"telegram_token"`
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
		return Config{}, errors.Wrap(err, "read file")
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return Config{}, errors.Wrap(err, "unmarshal file")
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
