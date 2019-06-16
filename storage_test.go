package main

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
	defer os.Remove(f)

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

func TestFileStorage_SaveChats(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f)

	fs, err := NewFileStorage(f)
	assert.NoError(t, err)

	err = fs.SaveChats([]int64{1, 2, 3})
	assert.NoError(t, err)

	b, err := ioutil.ReadFile(fs.file)
	assert.NoError(t, err)
	assert.Equal(t, "{\n  \"chats\": [\n    1,\n    2,\n    3\n  ]\n}", string(b))
}

func TestFileStorage_GetChats(t *testing.T) {
	f := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("feed-bot-testing-%d", time.Now().Nanosecond()),
	)
	defer os.Remove(f)

	t.Run("valid data", func(t *testing.T) {
		err := ioutil.WriteFile(f, []byte("{\n  \"chats\": [\n    1,\n    2,\n    3\n  ]\n}"), 0600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		nn, err := fs.GetChats()
		assert.NoError(t, err)
		assert.Equal(t, []int64{1, 2, 3}, nn)
	})
	t.Run("invalid data", func(t *testing.T) {
		err := ioutil.WriteFile(f, []byte("1,a,3"), 0600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		_, err = fs.GetChats()
		assert.Error(t, err)
	})
	t.Run("empty file", func(t *testing.T) {
		err := ioutil.WriteFile(f, []byte(""), 0600)
		assert.NoError(t, err)

		fs, err := NewFileStorage(f)
		assert.NoError(t, err)

		_, err = fs.GetChats()
		assert.Error(t, err)
	})
}
