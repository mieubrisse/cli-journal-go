package entry_item

import (
	"github.com/mieubrisse/cli-journal-go/components/filterable_checklist_item"
	"time"
)

type Component interface {
	filterable_checklist_item.Component

	GetTimestamp() time.Time
	GetName() string
	GetTags() []string
}
