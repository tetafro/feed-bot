package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Storage describes persistent datastorage.
type Storage interface {
	GetChats() ([]int64, error)
	SaveChats([]int64) error
}

// FileStorage is a storage that uses plain text file for storing data.
type FileStorage struct {
	file string
}

// NewFileStorage creates new file storage.
func NewFileStorage(file string) (*FileStorage, error) {
	fs := &FileStorage{file: file}

	_, err := os.Stat(fs.file)
	if os.IsNotExist(err) {
		// Init file with empty struct
		b, err := json.MarshalIndent(fileStorageData{}, "", "  ")
		if err != nil {
			return nil, errors.Wrap(err, "encode init data")
		}
		if err := ioutil.WriteFile(fs.file, b, 0600); err != nil {
			return nil, errors.Wrap(err, "init data file")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "check data file")
	}

	return fs, nil
}

// GetChats gets list of chat IDs from file.
func (s *FileStorage) GetChats() ([]int64, error) {
	data, err := s.read()
	if err != nil {
		return nil, errors.Wrap(err, "read data")
	}
	return data.Chats, nil
}

// SaveChats saves list of chat IDs to file.
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
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "encode data")
	}
	err = ioutil.WriteFile(s.file, b, 0600)
	if err != nil {
		return errors.Wrap(err, "write data to file")
	}
	return nil
}

// fileStorageData is an internal representation of data in the file.
type fileStorageData struct {
	Chats []int64 `json:"chats"`
}
