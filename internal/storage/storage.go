// Package storage is responsible for storing data.
package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// FileStorage is a storage that uses plain text file for storing data.
type FileStorage struct {
	file string
	mx   *sync.Mutex
}

// NewFileStorage creates new file storage.
func NewFileStorage(file string) (*FileStorage, error) {
	fs := &FileStorage{
		file: file,
		mx:   &sync.Mutex{},
	}

	_, err := os.Stat(fs.file)
	if os.IsNotExist(err) {
		// Init file with empty struct
		data := fileStorageData{
			Chats: []int64{},
			Feeds: map[string]time.Time{},
		}
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, errors.Wrap(err, "encode init data")
		}
		if err := ioutil.WriteFile(fs.file, b, 0o600); err != nil {
			return nil, errors.Wrap(err, "init data file")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "check data file")
	}

	return fs, nil
}

// GetChats gets list of chat IDs.
func (s *FileStorage) GetChats() ([]int64, error) {
	data, err := s.read()
	if err != nil {
		return nil, errors.Wrap(err, "read data")
	}
	return data.Chats, nil
}

// SaveChats saves list of chat IDs.
func (s *FileStorage) SaveChats(chats []int64) error {
	data, err := s.read()
	if err != nil {
		return errors.Wrap(err, "read data")
	}
	data.Chats = chats
	if err = s.save(data); err != nil {
		return errors.Wrap(err, "save data")
	}
	return nil
}

// GetLastUpdate gets last update time of the feed.
func (s *FileStorage) GetLastUpdate(feed string) (time.Time, error) {
	data, err := s.read()
	if err != nil {
		return time.Time{}, errors.Wrap(err, "read data")
	}
	t, ok := data.Feeds[feed]
	if !ok {
		return time.Time{}, nil
	}
	return t, nil
}

// SaveLastUpdate saves last feed update.
func (s *FileStorage) SaveLastUpdate(feed string, t time.Time) error {
	data, err := s.read()
	if err != nil {
		return errors.Wrap(err, "read data")
	}
	data.Feeds[feed] = t
	if err = s.save(data); err != nil {
		return errors.Wrap(err, "save data")
	}
	return nil
}

// read reads and parses whole data file.
func (s *FileStorage) read() (fileStorageData, error) {
	b, err := ioutil.ReadFile(s.file)
	if err != nil {
		return fileStorageData{}, errors.Wrap(err, "read file")
	}

	var data fileStorageData
	if err = json.Unmarshal(b, &data); err != nil {
		return fileStorageData{}, errors.Wrap(err, "decode data")
	}

	return data, nil
}

// save rewrites whole data file.
func (s *FileStorage) save(data fileStorageData) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "encode data")
	}
	err = ioutil.WriteFile(s.file, b, 0o600)
	if err != nil {
		return errors.Wrap(err, "write data to file")
	}
	return nil
}

// fileStorageData is an internal representation of data in the file.
type fileStorageData struct {
	Chats []int64              `json:"chats"`
	Feeds map[string]time.Time `json:"feeds"`
}
