package text_block

import "github.com/mieubrisse/cli-journal-go/components"

type TextBlockComponent interface {
	components.Component

	// TODO remove this??
	GetContents() string
}
