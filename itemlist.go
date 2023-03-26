package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"strings"
)

type Mode int

const (
	navigationMode Mode = iota

	filterMode Mode = iota

	// TODO add mode
)

var termenvOutput = termenv.DefaultOutput()

var cursorItemGreen = termenvOutput.Color("#4e9a06")
var selectedItemsYellow = termenvOutput.Color("#eed202")

type contentItem struct {
	name string
	tags []string
}

type contentList struct {
	mode Mode

	filterText string

	content             []contentItem
	cursorIdx           int
	selectedItemIndexes map[int]bool

	height int
	width  int
}

func (c *contentList) Init() tea.Cmd {
	return nil
}

func (c *contentList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return c, tea.Quit
		}

		switch c.mode {
		case navigationMode:
			return c, c.handleNavigationKeypress(msg)
		case filterMode:
			return c, c.handleFilteringKeypress(msg)
		}
	case tea.WindowSizeMsg:
		c.height = msg.Height
		c.width = msg.Width
	}

	return c, nil
}

func (c *contentList) View() string {
	resultBuilder := strings.Builder{}

	header := ""
	if c.mode == filterMode {
		header = " > " + c.filterText
	}
	resultBuilder.WriteString(header + "\n")

	resultBuilder.WriteString("\n")

	for idx, item := range c.content {
		maybeCheckmark := " "
		if _, found := c.selectedItemIndexes[idx]; found {
			maybeCheckmark = "✔️"
		}

		colorizedItemName := termenvOutput.String(item.name)
		if idx == c.cursorIdx {
			colorizedItemName = colorizedItemName.Foreground(cursorItemGreen).Bold()
		}

		row := fmt.Sprintf(" %s  %s", maybeCheckmark, colorizedItemName)

		resultBuilder.WriteString(row + "\n")
	}
	resultBuilder.WriteString("\n")

	footerRow := ""
	if len(c.selectedItemIndexes) > 0 {
		selectedItemsText := fmt.Sprintf(" %d items selected", len(c.selectedItemIndexes))
		footerRow = termenvOutput.String(selectedItemsText).Foreground(selectedItemsYellow).String()
	}
	resultBuilder.WriteString(footerRow + "\n")

	return resultBuilder.String()
}

func (c *contentList) handleNavigationKeypress(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "j":
		newCursorIdx := c.cursorIdx + 1
		if newCursorIdx < len(c.content) {
			c.cursorIdx = newCursorIdx
		}
	case "k":
		newCursorIdx := c.cursorIdx - 1
		if newCursorIdx >= 0 {
			c.cursorIdx = newCursorIdx
		}
	case "x":
		_, found := c.selectedItemIndexes[c.cursorIdx]
		if found {
			delete(c.selectedItemIndexes, c.cursorIdx)
		} else {
			c.selectedItemIndexes[c.cursorIdx] = true
		}
	case "a":
		if len(c.selectedItemIndexes) < len(c.content) {
			for idx := range c.content {
				c.selectedItemIndexes[idx] = true
			}
		} else {
			c.selectedItemIndexes = make(map[int]bool)
		}
	case "/":
		c.mode = filterMode
	}

	return nil
}

func (c *contentList) handleFilteringKeypress(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {

	case "backspace":
		if len(c.filterText) > 0 {
			c.filterText = c.filterText[0 : len(c.filterText)-1]
		}
	case "esc":
		c.filterText = ""
		c.mode = navigationMode
	default:
		c.filterText = c.filterText + string(msg.Runes)
	}

	return nil
}
