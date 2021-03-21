package bot

import (
	"bytes"
	"context"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/tetafro/feed-bot/internal/feed"
)

func TestNewBot(t *testing.T) {
	b := New(&testNotifier{}, []Feed{&testFeed{}}, 5*time.Second)
	assert.NotNil(t, b.notifier)
	assert.Len(t, b.feeds, 1)
	assert.Equal(t, b.interval, 5*time.Second)
}

func TestBot_Run(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
		defer cancel()

		n := &testNotifier{}
		f1 := &testFeed{
			items: []feed.Item{{Title: "One"}, {Title: "Two"}},
		}
		f2 := &testFeed{
			items: []feed.Item{{Title: "Three"}, {Title: "Four"}},
		}
		b := New(n, []Feed{f1, f2}, 50*time.Millisecond)

		b.Run(ctx)

		expected := []feed.Item{
			{Title: "One"}, {Title: "Two"}, {Title: "Three"}, {Title: "Four"},
		}
		assert.ElementsMatch(t, expected, n.items)
	})

	t.Run("no data", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		n := &testNotifier{}
		b := New(n, []Feed{&testFeed{}}, 25*time.Millisecond)

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

		ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
		defer cancel()

		n := &testNotifier{}
		f1 := &testFeed{
			items: []feed.Item{{Title: "One"}, {Title: "Two"}},
		}
		f2 := &testFeed{err: errors.New("fail")}
		b := New(n, []Feed{f1, f2}, 50*time.Millisecond)

		b.Run(ctx)

		expected := []feed.Item{
			{Title: "One"}, {Title: "Two"},
		}
		assert.ElementsMatch(t, expected, n.items)
		assert.Equal(t, "Failed to fetch items: fail\n", buf.String())
	})
}

type testNotifier struct {
	items []feed.Item
}

func (n *testNotifier) Notify(_ context.Context, item feed.Item) {
	n.items = append(n.items, item)
}

type testFeed struct {
	items []feed.Item
	err   error
}

func (f *testFeed) Fetch() ([]feed.Item, error) {
	return f.items, f.err
}
