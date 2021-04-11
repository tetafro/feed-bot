package feed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestItem_String(t *testing.T) {
	item := Item{
		Published: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
		Title:     "Title",
		Image:     "http://example.com/image.png",
	}

	expected := "[2020-01-01 10:00] Title http://example.com/image.png"
	assert.Equal(t, expected, item.String())
}
