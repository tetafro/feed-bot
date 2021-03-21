package feed

import (
	"fmt"
	"time"
)

// Item is a single fetched item.
type Item struct {
	Published time.Time
	Title     string
	Image     string
}

func (i Item) String() string {
	return fmt.Sprintf("[%s] %s: %s",
		i.Published.Format("2006-01-02 15:04"),
		i.Title, i.Image)
}
