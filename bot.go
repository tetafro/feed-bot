package main

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Notifier describes how clients are notified about new items.
type Notifier interface {
	Notify(context.Context, Item) error
}

// Fetcher fetches items from the given source.
type Fetcher interface {
	Fetch(url string) ([]Item, error)
}

// Bot fetches new items from data feeds, and sends it to all clients.
type Bot struct {
	notifier Notifier
	fetcher  Fetcher
	feeds    []string
	interval time.Duration
	log      *logrus.Logger
}

// NewBot creates new bot.
func NewBot(n Notifier, f Fetcher, feeds []string, interval time.Duration, log *logrus.Logger) *Bot {
	return &Bot{notifier: n, fetcher: f, feeds: feeds, interval: interval, log: log}
}

// Run starts listening for updates.
func (b *Bot) Run(ctx context.Context) {
	items := make(chan Item)

	var wg sync.WaitGroup
	wg.Add(len(b.feeds))
	for _, f := range b.feeds {
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
			b.log.Errorf("Failed to send notification: %v", err)
		}
	}
}

func (b *Bot) runFetches(ctx context.Context, f string, out chan Item) {
	// Run first fetch when started
	b.fetch(f, out)

	t := time.NewTicker(b.interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			b.fetch(f, out)
		case <-ctx.Done():
			return
		}
	}
}

func (b *Bot) fetch(f string, out chan Item) {
	items, err := b.fetcher.Fetch(f)
	if err != nil {
		b.log.Errorf("Failed to fetch items [%s]: %v", f, err)
		return
	}
	for _, item := range items {
		out <- item
	}
}
