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

func TestFileStorage_Save(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())

	fs := FileStorage{file: f.Name()}

	err = fs.Save([]int64{1, 2, 3})
	assert.NoError(t, err)

	b, err := ioutil.ReadFile(fs.file)
	assert.NoError(t, err)
	assert.Equal(t, "1,2,3", string(b))
}

func TestFileStorage_Get(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())

	t.Run("valid data", func(t *testing.T) {
		err = ioutil.WriteFile(f.Name(), []byte("1,2,3"), 0666)
		assert.NoError(t, err)

		fs := FileStorage{file: f.Name()}

		nn, err := fs.Get()
		assert.NoError(t, err)
		assert.Equal(t, []int64{1, 2, 3}, nn)
	})
	t.Run("invalid data", func(t *testing.T) {
		err = ioutil.WriteFile(f.Name(), []byte("1,a,3"), 0666)
		assert.NoError(t, err)

		fs := FileStorage{file: f.Name()}

		_, err := fs.Get()
		assert.Error(t, err)
	})
	t.Run("empty file", func(t *testing.T) {
		err = ioutil.WriteFile(f.Name(), []byte(""), 0666)
		assert.NoError(t, err)

		fs := FileStorage{file: f.Name()}

		nn, err := fs.Get()
		assert.NoError(t, err)
		assert.Nil(t, nn)
	})
}
