package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFileStorage(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f) // nolint: errcheck

	fs1, err := NewFileStorage(f)
	assert.NoError(t, err)
	assert.NotNil(t, fs1)

	_, err = os.Stat(f)
	assert.False(t, os.IsNotExist(err))

	// Try again with the same file
	fs2, err := NewFileStorage(f)
	assert.NoError(t, err)
	assert.NotNil(t, fs2)
}

func TestFileStorage_GetChats(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f) // nolint: errcheck

	t.Run("valid data", func(t *testing.T) {
		err := ioutil.WriteFile(f, []byte("{\n  \"chats\": [\n    1,\n    2,\n    3\n  ]\n}"), 0o600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		nn, err := fs.GetChats()
		assert.NoError(t, err)
		assert.Equal(t, []int64{1, 2, 3}, nn)
	})
	t.Run("invalid data", func(t *testing.T) {
		err := ioutil.WriteFile(f, []byte("1,a,3"), 0o600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		_, err = fs.GetChats()
		assert.Error(t, err)
	})
	t.Run("empty file", func(t *testing.T) {
		err := ioutil.WriteFile(f, []byte(""), 0o600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		_, err = fs.GetChats()
		assert.Error(t, err)
	})
}

func TestFileStorage_SaveChats(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f) // nolint: errcheck

	fs, err := NewFileStorage(f)
	assert.NoError(t, err)

	err = fs.SaveChats([]int64{1, 2, 3})
	assert.NoError(t, err)

	b, err := ioutil.ReadFile(fs.file)
	assert.NoError(t, err)

	expected := "{\n" +
		"  \"chats\": [\n" +
		"    1,\n" +
		"    2,\n" +
		"    3\n" +
		"  ],\n" +
		"  \"feeds\": {}\n" +
		"}"
	assert.Equal(t, expected, string(b))
}

func TestFileStorage_GetLastUpdate(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f) // nolint: errcheck

	t.Run("valid data", func(t *testing.T) {
		data := []byte(`{
			"chats": [],
			"feeds": {
				"my-feed": "2000-01-01T00:00:00Z"
			}
		}`)
		err := ioutil.WriteFile(f, data, 0o600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		ts, err := fs.GetLastUpdate("my-feed")
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), ts)
	})
	t.Run("unknown feed data", func(t *testing.T) {
		data := []byte(`{
			"chats": [],
			"feeds": {
				"my-feed": "2000-01-01T00:00:00Z"
			}
		}`)
		err := ioutil.WriteFile(f, data, 0o600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		ts, err := fs.GetLastUpdate("unknown")
		assert.NoError(t, err)
		assert.True(t, ts.IsZero())
	})
	t.Run("invalid data", func(t *testing.T) {
		err := ioutil.WriteFile(f, []byte("abc"), 0o600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		_, err = fs.GetLastUpdate("my-feed")
		assert.Error(t, err)
	})
	t.Run("empty file", func(t *testing.T) {
		err := ioutil.WriteFile(f, []byte(""), 0o600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		_, err = fs.GetLastUpdate("my-feed")
		assert.Error(t, err)
	})
}

func TestFileStorage_SaveLastUpdate(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f) // nolint: errcheck

	fs, err := NewFileStorage(f)
	assert.NoError(t, err)

	err = fs.SaveLastUpdate("my-feed", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)

	b, err := ioutil.ReadFile(fs.file)
	assert.NoError(t, err)
	assert.Equal(t, "{\n  \"chats\": [],\n  \"feeds\": {\n    \"my-feed\": \"2000-01-01T00:00:00Z\"\n  }\n}", string(b))
}
