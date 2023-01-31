package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFileStorage(t *testing.T) {
	t.Run("valid state", func(t *testing.T) {
		file := filepath.Join(
			os.TempDir(),
			fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
		)
		defer os.Remove(file) //nolint: errcheck

		data := []byte("feeds:\n" +
			"http://example.com/feed: 2021-03-20T05:00:00Z\n")
		assert.NoError(t, os.WriteFile(file, data, 0o600))

		fs, err := NewFileStorage(file)
		assert.NoError(t, err)
		assert.Equal(t, file, fs.file)
		assert.NotNil(t, fs.mx)

		b, err := os.ReadFile(fs.file)
		assert.NoError(t, err)
		assert.Equal(t, data, b)
	})

	t.Run("no state", func(t *testing.T) {
		file := filepath.Join(
			os.TempDir(),
			fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
		)
		defer os.Remove(file) //nolint: errcheck

		fs, err := NewFileStorage(file)
		assert.NoError(t, err)
		assert.Equal(t, file, fs.file)
		assert.NotNil(t, fs.mx)

		b, err := os.ReadFile(fs.file)
		assert.NoError(t, err)

		expected := "feeds: {}\n"
		assert.Equal(t, expected, string(b))
	})

	t.Run("empty state", func(t *testing.T) {
		file := filepath.Join(
			os.TempDir(),
			fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
		)
		defer os.Remove(file) //nolint: errcheck

		assert.NoError(t, os.WriteFile(file, []byte(""), 0o600))

		fs, err := NewFileStorage(file)
		assert.NoError(t, err)
		assert.Equal(t, file, fs.file)
		assert.NotNil(t, fs.mx)

		b, err := os.ReadFile(fs.file)
		assert.NoError(t, err)

		assert.Equal(t, "", string(b))
	})

	t.Run("invalid state", func(t *testing.T) {
		file := filepath.Join(
			os.TempDir(),
			fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
		)
		defer os.Remove(file) //nolint: errcheck

		data := []byte("]")
		assert.NoError(t, os.WriteFile(file, data, 0o600))

		_, err := NewFileStorage(file)
		assert.EqualError(t, err, "decode data: yaml: did not find expected node content")
	})
}

func TestFileStorage_GetLastUpdate(t *testing.T) {
	fs := &FileStorage{
		state: state{Feeds: map[string]time.Time{}},
		mx:    &sync.Mutex{},
	}

	ts := time.Now()
	fs.state.Feeds["feed1"] = ts

	assert.True(t, fs.GetLastUpdate("feed2").IsZero())
	assert.Equal(t, ts, fs.GetLastUpdate("feed1"))
}

func TestFileStorage_SaveLastUpdate(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f) //nolint: errcheck

	fs, err := NewFileStorage(f)
	assert.NoError(t, err)

	ts1 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ts2 := ts1.Add(time.Second)

	assert.NoError(t, fs.SaveLastUpdate("feed", ts1))
	assertFile(t, fs.file,
		"feeds:\n"+
			"  feed: 2000-01-01T00:00:00Z\n")

	assert.NoError(t, fs.SaveLastUpdate("feed", ts2))
	assertFile(t, fs.file,
		"feeds:\n"+
			"  feed: 2000-01-01T00:00:01Z\n")
}

func assertFile(t *testing.T, file, content string) {
	b, err := os.ReadFile(file) //nolint: gosec
	assert.NoError(t, err)
	assert.Equal(t, content, string(b))
}
