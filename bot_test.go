package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	log := logrus.New()
	log.Out = io.Discard
	b := NewBot(&testNotifier{}, &testFetcher{}, []string{"test"}, 5*time.Second, log)
	assert.NotNil(t, b.notifier)
	assert.Len(t, b.feeds, 1)
	assert.Equal(t, b.interval, 5*time.Second)
}

func TestBot_Run(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
		defer cancel()

		log := logrus.New()
		log.Out = io.Discard

		n := &testNotifier{}
		f := &testFetcher{
			items: map[string][]Item{
				"f1": {{Link: "One"}, {Link: "Two"}},
				"f2": {{Link: "Three"}, {Link: "Four"}},
			},
		}
		b := NewBot(n, f, []string{"f1", "f2"}, 1*time.Millisecond, log)

		b.Run(ctx)

		expected := []Item{
			{Link: "One"}, {Link: "Two"}, {Link: "Three"}, {Link: "Four"},
		}
		assert.ElementsMatch(t, expected, n.items)
	})

	t.Run("no data", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		log := logrus.New()
		log.Out = io.Discard

		n := &testNotifier{}
		f := &testFetcher{}
		b := NewBot(n, f, []string{"f1", "f2"}, 1*time.Millisecond, log)

		cancel()
		b.Run(ctx)

		assert.Len(t, n.items, 0)
	})

	t.Run("feed error", func(t *testing.T) {
		var buf bytes.Buffer
		log := logrus.New()
		log.SetOutput(&buf)

		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
		defer cancel()

		n := &testNotifier{}
		f := &testFetcher{err: errors.New("fail")}
		b := NewBot(n, f, []string{"f1"}, 1*time.Millisecond, log)

		b.Run(ctx)

		assert.Contains(t, buf.String(), "Failed to fetch items [f1]: fail")
	})
}

type testNotifier struct {
	items []Item
}

func (n *testNotifier) Notify(_ context.Context, item Item) error {
	n.items = append(n.items, item)
	return nil
}

type testFetcher struct {
	items map[string][]Item
	err   error
	done  map[string]bool
	mx    sync.Mutex
}

func (f *testFetcher) Fetch(url string) ([]Item, error) {
	f.mx.Lock()
	defer f.mx.Unlock()

	if f.done == nil {
		f.done = map[string]bool{}
	}
	if f.done[url] {
		return nil, nil
	}
	f.done[url] = true
	return f.items[url], f.err
}
