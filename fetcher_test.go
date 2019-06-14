package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestXKCDFetcher(t *testing.T) {
	server := httptest.NewServer(testXKCDServer{})
	defer server.Close()

	f := NewXKCDFetcher()
	f.addr = server.URL

	item, err := f.Fetch()
	assert.NoError(t, err)

	expected := Item{
		Title:     "Title",
		Image:     "https://example.com/0001.png",
		Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal(t, expected, item)
}

func TestCommitStripFetcher(t *testing.T) {
	server := httptest.NewServer(testCommitStripServer{})
	defer server.Close()

	f := NewCommitStripFetcher()
	f.addr = server.URL

	item, err := f.Fetch()
	assert.NoError(t, err)

	expected := Item{
		Title:     "Title",
		Image:     "https://example.com/0001.png",
		Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal(t, expected, item)
}

func TestExplosmFetcher(t *testing.T) {
	server := httptest.NewServer(testExplosmServer{})
	defer server.Close()

	f := NewExplosmFetcher()
	f.addr = server.URL

	item, err := f.Fetch()
	assert.NoError(t, err)

	expected := Item{
		Title:     "Title",
		Image:     "https://example.com/0001.png",
		Published: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal(t, expected, item)
}

type testXKCDServer struct{}

func (testXKCDServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

type testCommitStripServer struct{}

func (testCommitStripServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
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
		</rss>`
	w.Write([]byte(xml)) // nolint: errcheck
}

type testExplosmServer struct{}

func (testExplosmServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
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
		</rss>`
	w.Write([]byte(xml)) // nolint: errcheck
}
