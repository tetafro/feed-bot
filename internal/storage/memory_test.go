package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMemStorage(t *testing.T) {
	mem := NewMemStorage()
	assert.NotNil(t, mem.mx)
}

func TestMemStorage_GetChats(t *testing.T) {
	mem := NewMemStorage()

	assert.Len(t, mem.GetChats(), 0)

	chats := []int64{1, 2, 3}
	mem.state.Chats = chats
	assert.Equal(t, chats, mem.GetChats())
}

func TestMemStorage_SaveChats(t *testing.T) {
	mem := NewMemStorage()

	chats := []int64{1, 2, 3}
	assert.NoError(t, mem.SaveChats(chats))
	assert.Equal(t, chats, mem.state.Chats)
}

func TestMemStorage_GetLastUpdate(t *testing.T) {
	mem := NewMemStorage()

	ts := time.Now()
	mem.state.Feeds["feed1"] = ts

	assert.True(t, mem.GetLastUpdate("feed2").IsZero())
	assert.Equal(t, ts, mem.GetLastUpdate("feed1"))
}

func TestMemStorage_SaveLastUpdate(t *testing.T) {
	mem := NewMemStorage()

	ts1 := time.Now()
	ts2 := ts1.Add(1 * time.Second)

	assert.NoError(t, mem.SaveLastUpdate("feed", ts1))
	assert.Equal(t, ts1, mem.state.Feeds["feed"])
	assert.NoError(t, mem.SaveLastUpdate("feed", ts2))
	assert.Equal(t, ts2, mem.state.Feeds["feed"])
}
