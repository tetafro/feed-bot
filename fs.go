package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// FileStorage is a storage that uses plain a text file for storing data.
type FileStorage struct {
	file  string
	state state
	mx    *sync.Mutex
}

// state is a representation of application state.
type state struct {
	Feeds map[string]time.Time `yaml:"feeds"`
}

// NewFileStorage creates new file storage.
func NewFileStorage(file string) (*FileStorage, error) {
	s := &FileStorage{
		file:  file,
		state: state{Feeds: map[string]time.Time{}},
		mx:    &sync.Mutex{},
	}

	// Read or init
	b, err := os.ReadFile(s.file)
	if os.IsNotExist(err) {
		if err := s.save(); err != nil {
			return nil, fmt.Errorf("init file: %w", err)
		}
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// Unmarshal
	if err = yaml.Unmarshal(b, &s.state); err != nil {
		return nil, fmt.Errorf("decode data: %w", err)
	}
	if s.state.Feeds == nil {
		s.state.Feeds = map[string]time.Time{}
	}
	return s, nil
}

// GetLastUpdate gets last update time of the feed.
func (s *FileStorage) GetLastUpdate(feed string) time.Time {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.state.Feeds[feed]
}

// SaveLastUpdate saves last feed update.
func (s *FileStorage) SaveLastUpdate(feed string, t time.Time) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.state.Feeds[feed] = t
	return s.save()
}

// save rewrites whole current state in file.
func (s *FileStorage) save() error {
	b, err := yaml.Marshal(s.state)
	if err != nil {
		return fmt.Errorf("encode data: %w", err)
	}
	err = os.WriteFile(s.file, b, 0o600)
	if err != nil {
		return fmt.Errorf("write data to file: %w", err)
	}
	return nil
}
