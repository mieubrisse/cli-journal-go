package text_block

import (
	"github.com/mieubrisse/cli-journal-go/components/framework"
)

type TextBlockComponent interface {
	framework.Component

	// TODO remove this??
	GetContents() string
}
