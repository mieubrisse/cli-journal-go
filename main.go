package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/mieubrisse/cli-journal-go/content_item"
	"github.com/mieubrisse/cli-journal-go/filterable_content_list"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	content := []content_item.ContentItem{
		{
			Name: "Foo",
			Tags: nil,
		},
		{
			Name: "Bar",
			Tags: nil,
		},
		{
			Name: "Bang",
			Tags: nil,
		},
	}

	filterInput := textinput.New()

	contentList := filterable_content_list.New(content)

	topLevelModel := &appModel{
		mode:        navigationMode,
		filterInput: filterInput,
		contentList: contentList,
		height:      0,
		width:       0,
	}

	p := tea.NewProgram(topLevelModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
