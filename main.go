package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/mieubrisse/cli-journal-go/components/app"
	"github.com/mieubrisse/cli-journal-go/components/filter_input"
	"github.com/mieubrisse/cli-journal-go/components/filterable_content_list"
	"github.com/mieubrisse/cli-journal-go/data_structures/content_item"
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

	filterInput := filter_input.New(textinput.New())

	contentList := filterable_content_list.New(content)
	contentList.Focus()

	topLevelModel := app.New(filterInput, contentList)

	p := tea.NewProgram(topLevelModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
