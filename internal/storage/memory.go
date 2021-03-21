package storage

import (
	"sync"
	"time"
)

// MemStorage is a storage that stores data in memory.
type MemStorage struct {
	state State
	mx    *sync.Mutex
}

// NewMemStorage creates new in-memory storage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		state: State{
			Chats: []int64{},
			Feeds: map[string]time.Time{},
		},
		mx: &sync.Mutex{},
	}
}

// GetChats gets list of chat IDs.
func (s *MemStorage) GetChats() []int64 {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.state.Chats
}

// SaveChats saves list of chat IDs.
func (s *MemStorage) SaveChats(chats []int64) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.state.Chats = chats
	return nil
}

// GetLastUpdate gets last update time of the feed.
func (s *MemStorage) GetLastUpdate(feed string) time.Time {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.state.Feeds[feed]
}

// SaveLastUpdate saves last feed update.
func (s *MemStorage) SaveLastUpdate(feed string, t time.Time) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.state.Feeds[feed] = t
	return nil
}
