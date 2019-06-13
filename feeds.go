package main

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// Feed is a generalization for fetching
// picture data from remote resource.
type Feed interface {
	Start()
	Stop()
	Last() (Item, error)
	Feed() <-chan Item
}

// Item is a single fetched item.
type Item struct {
	Published time.Time
	Title     string
	Image     string
}

func (it Item) String() string {
	return fmt.Sprintf("%s\n%s", it.Title, it.Image)
}

// RSS works with RSS feeds.
type RSS struct {
	parser   *gofeed.Parser
	url      string
	feed     chan Item
	interval time.Duration
	last     time.Time

	stop chan struct{}
	wg   *sync.WaitGroup
}

// NewRSS returns new RSS feed.
func NewRSS(url string, interval time.Duration) *RSS {
	return &RSS{
		parser:   gofeed.NewParser(),
		url:      url,
		feed:     make(chan Item),
		interval: interval,
		stop:     make(chan struct{}),
		wg:       &sync.WaitGroup{},
	}
}

// Start starts collecting the feed.
func (f *RSS) Start() {
	f.wg.Add(1)

	go func() {
		defer f.wg.Done()
		defer close(f.feed)

		t := time.NewTicker(f.interval)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				item, err := f.Last()
				if err != nil {
					log.Printf("Failed to fetch item: %v", err)
				}
				if !item.Published.After(f.last) {
					continue
				}

				log.Printf("Got new item from %s", f.url)
				f.feed <- item
				f.last = item.Published
			case <-f.stop:
				return
			}
		}
	}()
}

// Stop stops collecting the feed.
func (f *RSS) Stop() {
	close(f.stop)
	f.wg.Wait()
}

// Last fetches url of the last RSS item.
func (f *RSS) Last() (Item, error) {
	feed, err := f.parser.ParseURL(f.url)
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

// Feed periodically checks xkcd.com for updates and
// sends them to channel.
func (f *RSS) Feed() <-chan Item {
	return f.feed
}

var regexpImageSrc = regexp.MustCompile(`src="([^\s]+)"`)

func getImage(s string) string {
	res := regexpImageSrc.FindStringSubmatch(s)
	if len(res) != 2 {
		return ""
	}
	return res[1]
}
