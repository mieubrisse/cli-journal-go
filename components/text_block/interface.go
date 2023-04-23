package text_block

import (
	"github.com/mieubrisse/cli-journal-go/components"
)

type Component interface {
	components.Component

	// TODO remove this??
	GetContents() string
}
