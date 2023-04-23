package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go/app_components/app_model"
	"github.com/mieubrisse/cli-journal-go/app_components/entry_item"
	"os"
	"regexp"
	"time"
)

var acceptableFormFieldRegex = regexp.MustCompile("^[a-zA-Z0-9.-]+$")

func main() {
	// TODO set up more items and deal with pagination
	content := []entry_item.Component{
		entry_item.New(time.Now(), "scenarios.yml", []string{"general-reference/wealthdraft"}),
		entry_item.New(time.Now(), "projections.yml", []string{"project-support/wealthdraft"}),
		entry_item.New(time.Now(), "starlark-exploration.md", []string{"project-support/starlark"}),
		entry_item.New(time.Now(), "journalling-about-things.md", []string{}),
	}

	topLevelModel := app_model.New(content)

	p := tea.NewProgram(topLevelModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
