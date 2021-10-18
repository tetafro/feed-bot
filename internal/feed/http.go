package feed

import (
	"net/http"
	"time"
)

// Default HTTP client timeout.
const timeout = 3 * time.Second

func newHTTPClient() *http.Client {
	return &http.Client{Timeout: timeout}
}
