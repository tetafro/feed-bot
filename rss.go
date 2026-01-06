package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

// Default HTTP client timeout.
const timeout = 3 * time.Second

// Storage describes persistent datastorage.
type Storage interface {
	GetLastUpdate(feed string) time.Time
	SaveLastUpdate(feed string, t time.Time) error
}

// Item is a single fetched item.
type Item struct {
	Published time.Time
	Link      string
}

func (i Item) String() string {
	return fmt.Sprintf("[%s] %s",
		i.Published.Format("2006-01-02 15:04"),
		i.Link)
}

// RSSFeed reads data from RSS feed.
type RSSFeed struct {
	url     string
	storage Storage
	parser  *gofeed.Parser
}

// NewRSSFeed returns new RSS feed.
func NewRSSFeed(s Storage, url string) *RSSFeed {
	p := gofeed.NewParser()
	p.Client = &http.Client{Timeout: timeout}
	return &RSSFeed{
		url:     url,
		parser:  p,
		storage: s,
	}
}

// Name returns feed's URL for identification.
func (f *RSSFeed) Name() string {
	return f.url
}

// Fetch fetches last item from RSS feed.
func (f *RSSFeed) Fetch() ([]Item, error) {
	last := f.storage.GetLastUpdate(f.url)
	if last.IsZero() {
		// First access
		if err := f.storage.SaveLastUpdate(f.url, time.Now()); err != nil {
			return nil, fmt.Errorf("save last update time: %w", err)
		}
		return nil, nil
	}

	feed, err := f.parser.ParseURL(f.url)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	if len(feed.Items) == 0 {
		return nil, nil
	}

	var items []Item //nolint: prealloc
	for _, fitem := range feed.Items {
		item := parse(fitem)
		if !item.Published.After(last) {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, nil
	}

	if err := f.storage.SaveLastUpdate(f.url, items[0].Published); err != nil {
		return nil, fmt.Errorf("save last update time: %w", err)
	}
	return items, nil
}

func parse(in *gofeed.Item) Item {
	item := Item{
		Link:      in.Link,
		Published: time.Now(),
	}

	if in.PublishedParsed != nil {
		item.Published = *in.PublishedParsed
	} else if in.UpdatedParsed != nil {
		item.Published = *in.UpdatedParsed
	}

	return item
}
