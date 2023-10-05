package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	log := logrus.New()
	log.Out = io.Discard
	b := NewBot(&testNotifier{}, []Feed{&testFeed{}}, 5*time.Second, log)
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
		f1 := &testFeed{
			items: []Item{{Link: "One"}, {Link: "Two"}},
		}
		f2 := &testFeed{
			items: []Item{{Link: "Three"}, {Link: "Four"}},
		}
		b := NewBot(n, []Feed{f1, f2}, 1*time.Millisecond, log)

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
		b := NewBot(n, []Feed{&testFeed{}}, 1*time.Millisecond, log)

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
		f1 := &testFeed{
			items: []Item{{Link: "One"}, {Link: "Two"}},
		}
		f2 := &testFeed{err: errors.New("fail")}
		b := NewBot(n, []Feed{f1, f2}, 1*time.Millisecond, log)

		b.Run(ctx)

		expected := []Item{
			{Link: "One"}, {Link: "Two"},
		}
		assert.ElementsMatch(t, expected, n.items)
		assert.Contains(t, buf.String(), "Failed to fetch items [test-name]: fail")
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
