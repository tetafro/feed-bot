package feed

import (
	"fmt"
	"time"
)

// Item is a single fetched item.
type Item struct {
	Published time.Time
	Link      string
}

func (i Item) String() string {
	return fmt.Sprintf("[%s] %s",
		i.Published.Format("2006-01-02 15:04"),
		i.Link)
}
