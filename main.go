package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go/components/app"
	"github.com/mieubrisse/cli-journal-go/components/filterable_content_list"
	"github.com/mieubrisse/cli-journal-go/components/form"
	"github.com/mieubrisse/cli-journal-go/components/text_input"
	"github.com/mieubrisse/cli-journal-go/data_structures/content_item"
	"os"
	"regexp"
	"time"
)

var acceptableFormFieldRegex = regexp.MustCompile("^[a-zA-Z0-9.-]+$")

func main() {
	// TODO set up more items and deal with pagination
	content := []content_item.ContentItem{
		{
			Timestamp: time.Now(),
			Name:      "scenarios.yml",
			Tags: []string{
				"general-reference/wealthdraft",
			},
		},
		{
			Timestamp: time.Now(),
			Name:      "projections.yml",
			Tags: []string{
				"project-support/wealthdraft",
			},
		},
		{
			Timestamp: time.Now(),
			Name:      "starlark-exploration.md",
			Tags:      []string{"project-support/starlark"},
		},
		{
			Timestamp: time.Now(),
			Name:      "journalling-about-frustrations.md",
			Tags:      []string{},
		},
	}

	createContentFormInput := text_input.New("Name: ")
	createContentForm := form.New(
		"Create Content",
		createContentFormInput,
		func(text string) bool {
			return acceptableFormFieldRegex.MatchString(text)
		},
	)

	nameFilterInput := text_input.New("🔎 ")
	tagFilterInput := text_input.New("🏷️  ")

	contentList := filterable_content_list.New(content)
	contentList.Focus()

	topLevelModel := app.New(createContentForm, nameFilterInput, tagFilterInput, contentList)

	p := tea.NewProgram(topLevelModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
