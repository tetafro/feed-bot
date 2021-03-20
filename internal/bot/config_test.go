package bot

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		cfg := &Config{}
		err := cfg.UnmarshalJSON([]byte(json))
		assert.NoError(t, err)

		expected := &Config{
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
		cfg := &Config{}
		err := cfg.UnmarshalJSON([]byte(json))
		assert.Error(t, err)
	})
	t.Run("malformed duration", func(t *testing.T) {
		json := `{
			"telegram_token": "abc",
			"update_interval": "10S",
			"feeds": []
		}`
		cfg := &Config{}
		err := cfg.UnmarshalJSON([]byte(json))
		assert.Error(t, err)
	})
}

func TestReadConfig(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f) // nolint: errcheck

	t.Run("valid config", func(t *testing.T) {
		data := []byte(`{
			"telegram_token": "123456789:AAAAAAAAAAAAAAAAAAAAAAAAAAAAA-AAAAA",
			"update_interval": "3h",
			"data_file": "./data.json",
			"feeds": ["https://example.com/rss.xml"]
		}`)
		err := ioutil.WriteFile(f, data, 0o600)
		assert.NoError(t, err)

		conf, err := ReadConfig(f)
		assert.NoError(t, err)

		expected := Config{
			TelegramToken:  "123456789:AAAAAAAAAAAAAAAAAAAAAAAAAAAAA-AAAAA",
			UpdateInterval: 3 * time.Hour,
			DataFile:       "./data.json",
			Feeds:          []string{"https://example.com/rss.xml"},
		}
		assert.Equal(t, expected, conf)
	})
	t.Run("invalid config", func(t *testing.T) {
		data := []byte(`]`)
		err := ioutil.WriteFile(f, data, 0o600)
		assert.NoError(t, err)

		_, err = ReadConfig(f)
		assert.EqualError(t, err,
			"unmarshal file: invalid character ']' looking for beginning of value")
	})
	t.Run("non-existing file", func(t *testing.T) {
		_, err := ReadConfig("abc.json")
		assert.EqualError(t, err,
			"read file: open abc.json: no such file or directory")
	})
}
