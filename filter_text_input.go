package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type filterTextInput struct {
	content string
}

func (input filterTextInput) Init() tea.Cmd {
	return nil
}

func (input filterTextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			if len(input.content) > 0 {
				input.content = input.content[0 : len(input.content)-1]
			}
		default:
			input.content = input.content + string(msg.Runes)
		}
	}
	return input, nil
}

func (input filterTextInput) View() string {
	return input.content
}
