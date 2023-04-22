package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"regexp"
)

var acceptableFormFieldRegex = regexp.MustCompile("^[a-zA-Z0-9.-]+$")

func main() {
	topLevelModel := New()

	p := tea.NewProgram(topLevelModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
