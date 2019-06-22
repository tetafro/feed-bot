package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshal json", func(t *testing.T) {
		json := `{
			"telegram_token": "abc",
			"update_interval": "10s",
			"data_file": "/tmp/data.json",
			"feeds": [
				"https://example-1.com",
				"https://example-2.com"
			]
		}`
		cfg := &config{}
		err := cfg.UnmarshalJSON([]byte(json))
		assert.NoError(t, err)

		expected := &config{
			TelegramToken:  "abc",
			UpdateInterval: 10 * time.Second,
			DataFile:       "/tmp/data.json",
			Feeds: []string{
				"https://example-1.com",
				"https://example-2.com",
			},
		}
		assert.Equal(t, expected, cfg)
	})
	t.Run("malformed json", func(t *testing.T) {
		json := `{"telegram_token": "abc",`
		cfg := &config{}
		err := cfg.UnmarshalJSON([]byte(json))
		assert.Error(t, err)
	})
	t.Run("malformed duration", func(t *testing.T) {
		json := `{
			"telegram_token": "abc",
			"update_interval": "10S",
			"feeds": []
		}`
		cfg := &config{}
		err := cfg.UnmarshalJSON([]byte(json))
		assert.Error(t, err)
	})
}
