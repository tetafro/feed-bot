package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestRSSFetcher(t *testing.T) {
	server := httptest.NewServer(testServer{})
	defer server.Close()

	f := NewRSSFetcher()

	item, err := f.Fetch(server.URL)
	assert.NoError(t, err)

	expected := Item{
		Title:     "Title",
		Image:     "https://example.com/0001.png",
		Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal(t, expected, item)
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		out  Item
	}{
		{
			name: "xkcd",
			in: `<feed xml:lang="en">
					<title>example.com</title>
					<link href="https://example.com/"/>
					<id>https://example.com/</id>
					<updated>2000-01-01T00:00:00Z</updated>
					<entry>
						<title>Title</title>
						<link href="https://example.com/0001/"/>
						<updated>2000-01-01T00:00:00Z</updated>
						<id>https://example.com/0001/</id>
						<summary type="html">
							<img src="https://example.com/0001.png"/>
						</summary>
					</entry>
				</feed>`,
			out: Item{
				Title:     "Title",
				Image:     "https://example.com/0001.png",
				Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "commit strip",
			in: `<?xml version="1.0" encoding="UTF-8"?>
				<rss version="2.0" xmlns:content="http://example.com/">
					<channel>
						<title>example.com</title>
						<atom:link href="http://example.com/feed/?"/>
						<link>http://example.com</link>
						<item>
							<title>Title</title>
							<link>http://example.com/0001/</link>
							<pubDate>Sat, 01 Jan 2000 00:00:00 +0000</pubDate>
							<content:encoded>
								<![CDATA[<p><img src="https://example.com/0001.png"/></p>]]>
							</content:encoded>
						</item>
					</channel>
				</rss>`,
			out: Item{
				Title:     "Title",
				Image:     "https://example.com/0001.png",
				Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "explosm",
			in: `<?xml version="1.0" encoding="UTF-8"?>
				<rss version="2.0">
					<channel>
						<title>Title</title>
						<link>https://example.com/</link>
						<item>
							<title>Title</title>
							<link>http://example.com/0001/</link>
							<description><img src="//example.com/0001.png"></description>
							<pubDate>Sat, 01 Jan 2000 00:00:00 +0000</pubDate>
						</item>
					</channel>
				</rss>`,
			out: Item{
				Title:     "Title",
				Image:     "https://example.com/0001.png",
				Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "smbc",
			in: `<?xml version="1.0" encoding="UTF-8"?>
				<rss version="2.0" xmlns:content="http://example.com/">
					<channel>
						<title>example.com</title>
						<atom:link href="http://example.com/feed/?"/>
						<link>http://example.com</link>
						<item>
							<title>Title</title>
							<link>http://example.com/0001/</link>
							<pubDate>Sat, 01 Jan 2000 00:00:00 +0000</pubDate>
							<content:encoded>
								<![CDATA[<p><img src="https://example.com/0001.png"/></p>]]>
							</content:encoded>
						</item>
					</channel>
				</rss>`,
			out: Item{
				Title:     "Title",
				Image:     "https://example.com/0001.png",
				Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.in)
			rss, err := gofeed.NewParser().Parse(r)
			assert.NoError(t, err)
			assert.Equal(t, tt.out, parse(rss.Items[0]))
		})
	}
}

type testServer struct{}

func (testServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	xml := `<feed xml:lang="en">
			<title>example.com</title>
			<link href="https://example.com/"/>
			<id>https://example.com/</id>
			<updated>2000-01-01T00:00:00Z</updated>
			<entry>
				<title>Title</title>
				<link href="https://example.com/0001/"/>
				<updated>2000-01-01T00:00:00Z</updated>
				<id>https://example.com/0001/</id>
				<summary type="html">
					<img src="https://example.com/0001.png"/>
				</summary>
			</entry>
		</feed>`
	w.Write([]byte(xml)) // nolint: errcheck
}
