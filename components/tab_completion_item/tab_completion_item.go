package tab_completion_item

import "github.com/charmbracelet/lipgloss"

type TabCompletionItem struct {
	contents string

	// TODO add matched char index

	width  int
	height int
}

func New(contents string) *TabCompletionItem {
	return &TabCompletionItem{
		contents: contents,
		width:    0,
		height:   0,
	}
}

func (item TabCompletionItem) GetContents() string {
	return item.contents
}

func (item TabCompletionItem) View() string {
	return lipgloss.NewStyle().
		MaxWidth(item.width).
		MaxHeight(item.height).
		Render(item.contents)
}

func (item *TabCompletionItem) Resize(width int, height int) {
	item.width = width
	item.height = height
}

func (item TabCompletionItem) GetWidth() int {
	return item.width
}

func (item TabCompletionItem) GetHeight() int {
	return item.height
}
