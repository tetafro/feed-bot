package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFeed(t *testing.T) {
	t.Run("fetch new item", func(t *testing.T) {
		item := Item{
			Title:     "Title",
			Image:     "https://example.com/0001.png",
			Published: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		fetcher := &mockFetcher{item: item}
		st := &mockStorage{}

		f := NewFeed(st, "http://example.com", fetcher, 25*time.Millisecond)
		f.Start()

		timer := time.NewTimer(50 * time.Millisecond)
		select {
		case <-timer.C:
			t.Fatal("Got no item")
		case got := <-f.Feed():
			assert.Equal(t, item, got)
		}

		timer.Stop()
		f.Stop()
	})
	t.Run("fetch old item", func(t *testing.T) {
		item := Item{
			Title:     "Title",
			Image:     "https://example.com/0001.png",
			Published: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		fetcher := &mockFetcher{item: item}
		st := &mockStorage{}

		f := NewFeed(st, "http://example.com", fetcher, 25*time.Millisecond)
		f.Start()

		timer := time.NewTimer(50 * time.Millisecond)
		select {
		case <-f.Feed():
			t.Fatal("Got unexpected item")
		case <-timer.C:
		}

		timer.Stop()
		f.Stop()
	})
}

type mockFetcher struct {
	item Item
	err  error
}

func (m *mockFetcher) Fetch(_ string) (Item, error) {
	return m.item, m.err
}
