package feed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestItem_String(t *testing.T) {
	item := Item{
		Published: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
		Link:      "http://example.com/content/",
	}

	expected := "[2020-01-01 10:00] http://example.com/content/"
	assert.Equal(t, expected, item.String())
}
