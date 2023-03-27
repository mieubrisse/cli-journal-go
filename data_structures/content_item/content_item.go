package content_item

import "time"

type ContentItem struct {
	Timestamp time.Time
	Name      string
	Tags      []string
}
