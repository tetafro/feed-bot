package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f) //nolint: errcheck

	t.Run("valid config", func(t *testing.T) {
		data := []byte("telegram_token: \"123456789:AAAAAAAAAAAAAAAAAAAAAAAAAAAAA-AAAAA\"\n" +
			"update_interval: 3h\n" +
			"data_file: ./data.yaml\n" +
			"feeds: [\"https://example.com/rss.xml\"]\n")
		assert.NoError(t, os.WriteFile(f, data, 0o600))

		conf, err := ReadConfig(f)
		assert.NoError(t, err)

		expected := Config{
			TelegramToken:  "123456789:AAAAAAAAAAAAAAAAAAAAAAAAAAAAA-AAAAA",
			UpdateInterval: 3 * time.Hour,
			DataFile:       "./data.yaml",
			Feeds:          []string{"https://example.com/rss.xml"},
			LogLevel:       "info",
		}
		assert.Equal(t, expected, conf)
	})

	t.Run("debug config", func(t *testing.T) {
		data := []byte("telegram_token: \"123456789:AAAAAAAAAAAAAAAAAAAAAAAAAAAAA-AAAAA\"\n" +
			"feeds: [\"https://example.com/rss.xml\"]\n" +
			"data_file: ./data.yaml\n" +
			"log_level: \"debug\"\n")
		assert.NoError(t, os.WriteFile(f, data, 0o600))

		conf, err := ReadConfig(f)
		assert.NoError(t, err)

		expected := Config{
			TelegramToken:  "123456789:AAAAAAAAAAAAAAAAAAAAAAAAAAAAA-AAAAA",
			UpdateInterval: defaultUpdateInterval,
			DataFile:       "./data.yaml",
			Feeds:          []string{"https://example.com/rss.xml"},
			LogLevel:       "debug",
		}
		assert.Equal(t, expected, conf)
	})

	t.Run("invalid config", func(t *testing.T) {
		data := []byte(`]`)
		assert.NoError(t, os.WriteFile(f, data, 0o600))

		_, err := ReadConfig(f)
		assert.EqualError(t, err,
			"unmarshal file: yaml: did not find expected node content")
	})

	t.Run("non-existing file", func(t *testing.T) {
		_, err := ReadConfig("abc.yaml")
		assert.True(t, errors.Is(err, os.ErrNotExist))
	})
}
