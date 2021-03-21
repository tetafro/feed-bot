package storage

import "time"

// State is a representation of application state.
type State struct {
	Chats []int64              `json:"chats"`
	Feeds map[string]time.Time `json:"feeds"`
}
