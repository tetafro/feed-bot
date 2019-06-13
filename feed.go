package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Item is a single fetched item.
type Item struct {
	Published time.Time
	Title     string
	Image     string
}

func (it Item) String() string {
	return fmt.Sprintf("%s\n%s", it.Title, it.Image)
}

// Feed works with data feeds..
type Feed struct {
	fetcher  Fetcher
	feed     chan Item
	interval time.Duration
	last     time.Time

	stop chan struct{}
	wg   *sync.WaitGroup
}

// NewFeed returns new feed.
func NewFeed(f Fetcher, interval time.Duration) *Feed {
	return &Feed{
		fetcher:  f,
		feed:     make(chan Item),
		interval: interval,
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
				item, err := f.fetcher.Fetch()
				if err != nil {
					log.Printf("Failed to fetch item: %v", err)
				}
				if !item.Published.After(f.last) {
					continue
				}

				f.feed <- item
				f.last = item.Published
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
