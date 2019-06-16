package main

import (
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// Fetcher fetches items from the storage.
type Fetcher interface {
	Fetch() (Item, error)
}

// RSSFetcher fetches items from RSS feed.
type RSSFetcher struct {
	addr   string
	parser *gofeed.Parser
}

// NewRSSFetcher creates new RSSFetcher.
func NewRSSFetcher(addr string) *RSSFetcher {
	return &RSSFetcher{
		addr:   addr,
		parser: gofeed.NewParser(),
	}
}

// Fetch fetches last item from RSS feed.
func (f *RSSFetcher) Fetch() (Item, error) {
	feed, err := f.parser.ParseURL(f.addr)
	if err != nil {
		return Item{}, errors.Wrap(err, "parse url")
	}
	if len(feed.Items) == 0 {
		return Item{}, errors.New("empty feed")
	}
	last := feed.Items[0]

	item := parse(last)

	return item, nil
}

func parse(in *gofeed.Item) Item {
	item := Item{
		Title:     in.Title,
		Published: time.Now(),
	}

	// Published
	if in.PublishedParsed != nil {
		item.Published = *in.PublishedParsed
	}
	if in.UpdatedParsed != nil {
		item.Published = *in.UpdatedParsed
	}

	// Image
	if in.Description != "" {
		res := regexpImageSrc.FindStringSubmatch(in.Description)
		if len(res) == 2 {
			item.Image = res[1]
		}
	}
	if in.Content != "" {
		res := regexpImageSrc.FindStringSubmatch(in.Content)
		if len(res) == 2 {
			item.Image = res[1]
		}
	}

	// Fix missing protocol
	if strings.HasPrefix(item.Image, "//") {
		item.Image = "https:" + item.Image
	}

	return item
}

var regexpImageSrc = regexp.MustCompile(`src="([^\s]+)"`)
