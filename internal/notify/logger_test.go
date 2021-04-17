package notify

import (
	"bytes"
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tetafro/feed-bot/internal/feed"
)

func TestLogNotifier_Notify(t *testing.T) {
	defer log.SetOutput(log.Writer())
	defer log.SetFlags(log.Flags())

	// Write logs to file, disable timestamps
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)

	n := NewLogNotifier()
	n.Notify(context.Background(), feed.Item{
		Published: time.Date(2000, 1, 1, 10, 0, 0, 0, time.UTC),
		Link:      "http://example.com/feed/",
	})

	expected := "[notify] New item: [2000-01-01 10:00] http://example.com/feed/\n"
	assert.Equal(t, expected, buf.String())
}
