package tab_completion_item

import (
	"github.com/mieubrisse/cli-journal-go/components/filterable_list_item"
	"github.com/mieubrisse/cli-journal-go/components/text_block"
)

type implementation struct {
	innerComponent text_block.TextBlockComponent

	contents string

	width  int
	height int
}

func New(contents string) filterable_list_item.Component {
	inner := text_block.New(contents)
	return &implementation{
		innerComponent: inner,
		contents:       contents,
		width:          0,
		height:         0,
	}
}

func (impl implementation) View() string {
	return impl.innerComponent.View()
}

func (impl *implementation) Resize(width int, height int) {
	impl.innerComponent.Resize(width, height)
	impl.width = width
	impl.height = height
}

func (impl implementation) GetWidth() int {
	return impl.width
}

func (impl implementation) GetHeight() int {
	return impl.height
}

func (impl implementation) GetValue() string {
	return impl.contents
}
