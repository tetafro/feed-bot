package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetImage(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "image inside tag",
			in:   `<img src="http://example.com/image.png">`,
			out:  "http://example.com/image.png",
		},
		{
			name: "image inside tag with other attributes",
			in:   `<img src="http://example.com/image.png" alt="text">`,
			out:  "http://example.com/image.png",
		},
		{
			name: "broken input",
			in:   `<img src="http://example.com/image.png>`,
			out:  "",
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(*testing.T) {
			assert.Equal(t, tt.out, getImage(tt.in))
		})
	}
}

func TestFeed(t *testing.T) {
	item := Item{
		Title:     "Title",
		Image:     "https://example.com/0001.png",
		Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	fetcher := &fetcherMock{item: item}

	f := NewFeed(fetcher, 25*time.Millisecond)
	go f.Start()

	// Got item for the first time
	timer := time.NewTimer(50 * time.Millisecond)
	select {
	case <-timer.C:
		t.Fatal("Got no item")
	case got := <-f.Feed():
		assert.Equal(t, item, got)
	}

	// Got the same item again (no output)
	timer = time.NewTimer(50 * time.Millisecond)
	select {
	case <-f.Feed():
		t.Fatal("Got unexpected item")
	case <-timer.C:
	}

	f.Stop()
}

type fetcherMock struct {
	item Item
	err  error
}

func (f *fetcherMock) Fetch() (Item, error) {
	return f.item, f.err
}
