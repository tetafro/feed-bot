package storage

// State is a representation of application state.
type State struct {
	Chats []int64          `yaml:"chats"`
	Feeds map[string]int64 `yaml:"feeds"`
}
