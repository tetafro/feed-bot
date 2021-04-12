// Package feed is responsible for getting data from external source (RSS).
package feed

import (
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// Storage describes persistent datastorage.
type Storage interface {
	GetLastUpdate(feed string) time.Time
	SaveLastUpdate(feed string, t time.Time) error
}

// RSSFeed reads data from RSS feed.
type RSSFeed struct {
	url     string
	storage Storage
	parser  *gofeed.Parser
}

// NewRSSFeed returns new RSS feed.
func NewRSSFeed(s Storage, url string) *RSSFeed {
	return &RSSFeed{
		url:     url,
		parser:  gofeed.NewParser(),
		storage: s,
	}
}

// Fetch fetches last item from RSS feed.
func (f *RSSFeed) Fetch() ([]Item, error) {
	last := f.storage.GetLastUpdate(f.url)
	if last.IsZero() {
		// First access
		if err := f.storage.SaveLastUpdate(f.url, time.Now()); err != nil {
			return nil, errors.Wrap(err, "save last update time")
		}
		return nil, nil
	}

	feed, err := f.parser.ParseURL(f.url)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}
	if len(feed.Items) == 0 {
		return nil, nil
	}

	var items []Item // nolint: prealloc
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
		return nil, errors.Wrap(err, "save last update time")
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
	}
	if in.UpdatedParsed != nil {
		item.Published = *in.UpdatedParsed
	}

	return item
}
