package app

import "github.com/mieubrisse/cli-journal-go/data_structures/content_item"

type UpdateContentMsg struct {
	newContent []content_item.ContentItem
}

func (msg UpdateContentMsg) GetNewContent() []content_item.ContentItem {
	return msg.newContent
}
