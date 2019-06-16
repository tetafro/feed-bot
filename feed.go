package main

import (
	"log"
	"sync"
	"time"
)

// Feed works with data feeds..
type Feed struct {
	url      string
	fetcher  Fetcher
	feed     chan Item
	interval time.Duration
	storage  Storage

	stop chan struct{}
	wg   *sync.WaitGroup
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
		stop:     make(chan struct{}),
		wg:       &sync.WaitGroup{},
	}
}

// Start starts collecting the feed.
func (f *Feed) Start() {
	f.wg.Add(1)

	go func() {
		defer f.wg.Done()
		defer close(f.feed)

		t := time.NewTicker(f.interval)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				item, err := f.fetcher.Fetch(f.url)
				if err != nil {
					log.Printf("Failed to fetch item: %v", err)
				}

				last, err := f.storage.GetLastUpdate(f.url)
				if err != nil {
					log.Printf("Failed to get last update time: %v", err)
					continue
				}

				if !item.Published.After(last) {
					continue
				}

				f.feed <- item
				if err := f.storage.SaveLastUpdate(f.url, time.Now().UTC()); err != nil {
					log.Printf("Failed to save last update time: %v", err)
				}
			case <-f.stop:
				return
			}
		}
	}()
}

// Stop stops collecting the feed.
func (f *Feed) Stop() {
	close(f.stop)
	f.wg.Wait()
}

// Feed periodically checks for updates and sends them to channel.
func (f *Feed) Feed() <-chan Item {
	return f.feed
}
