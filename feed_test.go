package main

import (
	"testing"

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
