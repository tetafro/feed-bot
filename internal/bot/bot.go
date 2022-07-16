// Package bot provides main application entity. It is responsible for wiring
// together all components: data feeds, chats, state storage.
package bot

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/tetafro/feed-bot/internal/feed"
)

// Notifier describes how clients are notified about new items.
type Notifier interface {
	Notify(context.Context, feed.Item) error
}

// Feed is a source of data.
type Feed interface {
	Fetch() ([]feed.Item, error)
}

// Bot fetches new items from data feeds, and sends it to all clients.
type Bot struct {
	notifier Notifier
	feeds    []Feed
	interval time.Duration
}

// New creates new bot.
func New(n Notifier, feeds []Feed, interval time.Duration) *Bot {
	return &Bot{notifier: n, feeds: feeds, interval: interval}
}

// Run starts listening for updates.
func (b *Bot) Run(ctx context.Context) {
	items := make(chan feed.Item)

	var wg sync.WaitGroup
	wg.Add(len(b.feeds))
	for _, f := range b.feeds {
		f := f
		go func() {
			b.runFetches(ctx, f, items)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(items)
	}()

	for item := range items {
		if err := b.notifier.Notify(ctx, item); err != nil {
			log.Printf("Failed to send notification: %v", err)
		}
	}
}

func (b *Bot) runFetches(ctx context.Context, f Feed, out chan feed.Item) {
	t := time.NewTicker(b.interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			items, err := f.Fetch()
			if err != nil {
				log.Printf("Failed to fetch items: %v", err)
				continue
			}
			for _, item := range items {
				out <- item
			}
		case <-ctx.Done():
			return
		}
	}
}
