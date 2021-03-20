// Package feed is responsible for getting data from external source (RSS).
package feed

import (
	"context"
	"log"
	"time"
)

// Storage describes persistent datastorage.
type Storage interface {
	GetLastUpdate(feed string) (time.Time, error)
	SaveLastUpdate(feed string, t time.Time) error
}

// Feed works with data feeds.
type Feed struct {
	url      string
	fetcher  Fetcher
	feed     chan Item
	interval time.Duration
	storage  Storage
}

// Item is a single fetched item.
type Item struct {
	Published time.Time
	Title     string
	Image     string
}

// NewFeed returns new feed.
func NewFeed(s Storage, url string, f Fetcher, interval time.Duration) *Feed {
	return &Feed{
		url:      url,
		fetcher:  f,
		feed:     make(chan Item),
		interval: interval,
		storage:  s,
	}
}

// Run starts collecting the feed.
func (f *Feed) Run(ctx context.Context) {
	defer close(f.feed)

	t := time.NewTicker(f.interval)
	defer t.Stop()

	f.fetch() // instant first fetch
	for {
		select {
		case <-t.C:
			f.fetch()
		case <-ctx.Done():
			return
		}
	}
}

// Feed periodically checks for updates and sends them to channel.
func (f *Feed) Feed() <-chan Item {
	return f.feed
}

func (f *Feed) fetch() {
	item, err := f.fetcher.Fetch(f.url)
	if err != nil {
		log.Printf("Failed to fetch item: %v", err)
	}

	last, err := f.storage.GetLastUpdate(f.url)
	if err != nil {
		log.Printf("Failed to get last update time: %v", err)
		return
	}

	// Item is not new
	if !item.Published.After(last) {
		return
	}

	f.feed <- item
	if err := f.storage.SaveLastUpdate(f.url, time.Now().UTC()); err != nil {
		log.Printf("Failed to save last update time: %v", err)
	}
}
