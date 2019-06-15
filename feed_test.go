package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFeed(t *testing.T) {
	item := Item{
		Title:     "Title",
		Image:     "https://example.com/0001.png",
		Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	fetcher := &mockFetcher{item: item}

	f := NewFeed(fetcher, 25*time.Millisecond)
	go f.Start()

	// Got item for the first time
	timer := time.NewTimer(50 * time.Millisecond)
	select {
	case <-timer.C:
		t.Fatal("Got no item")
	case got := <-f.Feed():
		assert.Equal(t, item, got)
	}

	// Got the same item again (no output)
	timer = time.NewTimer(50 * time.Millisecond)
	select {
	case <-f.Feed():
		t.Fatal("Got unexpected item")
	case <-timer.C:
	}

	f.Stop()
}

type mockFetcher struct {
	item Item
	err  error
}

func (m *mockFetcher) Fetch() (Item, error) {
	return m.item, m.err
}
