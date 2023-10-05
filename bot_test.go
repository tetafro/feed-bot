package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	b := NewBot(&testNotifier{}, []Feed{&testFeed{}}, 5*time.Second)
	assert.NotNil(t, b.notifier)
	assert.Len(t, b.feeds, 1)
	assert.Equal(t, b.interval, 5*time.Second)
}

func TestBot_Run(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
		defer cancel()

		n := &testNotifier{}
		f1 := &testFeed{
			items: []Item{{Link: "One"}, {Link: "Two"}},
		}
		f2 := &testFeed{
			items: []Item{{Link: "Three"}, {Link: "Four"}},
		}
		b := NewBot(n, []Feed{f1, f2}, 1*time.Millisecond)

		b.Run(ctx)

		expected := []Item{
			{Link: "One"}, {Link: "Two"}, {Link: "Three"}, {Link: "Four"},
		}
		assert.ElementsMatch(t, expected, n.items)
	})

	t.Run("no data", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		n := &testNotifier{}
		b := NewBot(n, []Feed{&testFeed{}}, 1*time.Millisecond)

		cancel()
		b.Run(ctx)

		assert.Len(t, n.items, 0)
	})

	t.Run("feed error", func(t *testing.T) {
		defer log.SetOutput(log.Writer())
		defer log.SetFlags(log.Flags())

		var buf bytes.Buffer
		log.SetOutput(&buf)
		log.SetFlags(0)

		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
		defer cancel()

		n := &testNotifier{}
		f1 := &testFeed{
			items: []Item{{Link: "One"}, {Link: "Two"}},
		}
		f2 := &testFeed{err: errors.New("fail")}
		b := NewBot(n, []Feed{f1, f2}, 1*time.Millisecond)

		b.Run(ctx)

		expected := []Item{
			{Link: "One"}, {Link: "Two"},
		}
		assert.ElementsMatch(t, expected, n.items)
		assert.Equal(t, "Failed to fetch items [test-name]: fail\n", buf.String())
	})
}

type testNotifier struct {
	items []Item
}

func (n *testNotifier) Notify(_ context.Context, item Item) error {
	n.items = append(n.items, item)
	return nil
}

type testFeed struct {
	items []Item
	err   error
	done  bool
}

func (f *testFeed) Name() string {
	return "test-name"
}

func (f *testFeed) Fetch() ([]Item, error) {
	if f.done {
		return nil, nil
	}
	f.done = true
	return f.items, f.err
}
