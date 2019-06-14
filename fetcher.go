package main

import (
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// Fetcher fetches items from the storage.
type Fetcher interface {
	Fetch() (Item, error)
}

const (
	xkcdFeedAddr        = "https://xkcd.com/atom.xml"
	commitstripFeedAddr = "http://www.commitstrip.com/en/feed/"
	explosmFeedAddr     = "http://explosm-feed.antonymale.co.uk/comics_feed"
)

// XKCDFetcher fetches items from https://xkcd.com/.
type XKCDFetcher struct {
	addr   string
	parser *gofeed.Parser
}

// NewXKCDFetcher creates new XKCDFetcher.
func NewXKCDFetcher() *XKCDFetcher {
	return &XKCDFetcher{
		addr:   xkcdFeedAddr,
		parser: gofeed.NewParser(),
	}
}

// Fetch fetches last item from RSS feed.
func (f *XKCDFetcher) Fetch() (Item, error) {
	feed, err := f.parser.ParseURL(f.addr)
	if err != nil {
		return Item{}, errors.Wrap(err, "parse url")
	}
	if len(feed.Items) == 0 {
		return Item{}, errors.New("empty feed")
	}
	last := feed.Items[0]

	item := Item{
		Title:     last.Title,
		Image:     getImage(last.Description),
		Published: time.Now(),
	}
	if last.PublishedParsed != nil {
		item.Published = *last.PublishedParsed
	}
	if last.UpdatedParsed != nil {
		item.Published = *last.UpdatedParsed
	}

	return item, nil
}

// CommitStripFetcher fetches items from https://commitstrip.com/.
type CommitStripFetcher struct {
	addr   string
	parser *gofeed.Parser
}

// NewCommitStripFetcher creates new CommitStripFetcher.
func NewCommitStripFetcher() *CommitStripFetcher {
	return &CommitStripFetcher{
		addr:   commitstripFeedAddr,
		parser: gofeed.NewParser(),
	}
}

// Fetch fetches last item from RSS feed.
func (f *CommitStripFetcher) Fetch() (Item, error) {
	feed, err := f.parser.ParseURL(f.addr)
	if err != nil {
		return Item{}, errors.Wrap(err, "parse url")
	}
	if len(feed.Items) == 0 {
		return Item{}, errors.New("empty feed")
	}
	last := feed.Items[0]

	item := Item{
		Title:     last.Title,
		Image:     getImage(last.Content),
		Published: time.Now(),
	}
	if last.PublishedParsed != nil {
		item.Published = *last.PublishedParsed
	}
	if last.UpdatedParsed != nil {
		item.Published = *last.UpdatedParsed
	}

	return item, nil
}

// ExplosmFetcher fetches items from http://explosm.net/.
type ExplosmFetcher struct {
	addr   string
	parser *gofeed.Parser
}

// NewExplosmFetcher creates new ExplosmFetcher.
func NewExplosmFetcher() *ExplosmFetcher {
	return &ExplosmFetcher{
		addr:   explosmFeedAddr,
		parser: gofeed.NewParser(),
	}
}

// Fetch fetches last item from RSS feed.
func (f *ExplosmFetcher) Fetch() (Item, error) {
	feed, err := f.parser.ParseURL(f.addr)
	if err != nil {
		return Item{}, errors.Wrap(err, "parse url")
	}
	if len(feed.Items) == 0 {
		return Item{}, errors.New("empty feed")
	}
	last := feed.Items[0]

	item := Item{
		Title:     last.Title,
		Image:     "https:" + getImage(last.Description),
		Published: time.Now(),
	}
	if last.PublishedParsed != nil {
		item.Published = *last.PublishedParsed
	}
	if last.UpdatedParsed != nil {
		item.Published = *last.UpdatedParsed
	}

	return item, nil
}

var regexpImageSrc = regexp.MustCompile(`src="([^\s]+)"`)

func getImage(s string) string {
	res := regexpImageSrc.FindStringSubmatch(s)
	if len(res) != 2 {
		return ""
	}
	return res[1]
}
