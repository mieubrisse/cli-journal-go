package filterable_list_item

import (
	"github.com/mieubrisse/cli-journal-go/components"
)

type Component interface {
	components.Component

	SetHighlighted(isHighlighted bool)
	GetValue() string
}
