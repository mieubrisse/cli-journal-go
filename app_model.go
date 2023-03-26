package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"regexp"
	"strings"
)

type Mode int

const (
	navigationMode Mode = iota

	filterMode Mode = iota

	// TODO add mode

	numListElems = 20
)

var termenvOutput = termenv.DefaultOutput()

var cursorItemGreen = termenvOutput.Color("#4e9a06")
var selectedItemsYellow = termenvOutput.Color("#eed202")

type contentItem struct {
	name string
	tags []string
}

type appModel struct {
	mode Mode

	// filterInput tea.Model
	filterInput textinput.Model

	content             []contentItem
	cursorIdx           int
	selectedItemIndexes map[int]bool

	height int
	width  int
}

func (model appModel) Init() tea.Cmd {
	return nil
}

// NOTE: This returns a model because BubbleTea expects models to be passed by-value, so the way to "update" the model
// is to return a new instance of it
func (model appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return model, tea.Quit
		}

		switch model.mode {
		case navigationMode:
			return model.handleNavigationKeypress(msg)
		case filterMode:
			// Back out of filter mode
			if msg.String() == "esc" {
				model.filterInput.Blur()
				model.mode = navigationMode
				return model, nil
			}

			var cmd tea.Cmd
			model.filterInput, cmd = model.filterInput.Update(msg)
			return model, cmd
		}
	case tea.WindowSizeMsg:
		model.height = msg.Height
		model.width = msg.Width
	}

	return model, nil
}

func (model appModel) View() string {
	resultBuilder := strings.Builder{}

	// Filtering header
	header := ""
	if model.mode == filterMode || len(model.filterInput.Value()) > 0 {
		header = " " + model.filterInput.View()
	}
	resultBuilder.WriteString(header + "\n")

	resultBuilder.WriteString("\n")

	escapedFilterText := regexp.QuoteMeta(model.filterInput.Value())
	nameSearchTerms := strings.Fields(escapedFilterText)
	// The (?i) makes the search case-insensitive
	regexStr := "(?i)" + strings.Join(nameSearchTerms, ".*")
	// Okay to use MustCompile here because we quote the user's input so it should be safe
	matcher := regexp.MustCompile(regexStr)
	for idx, item := range model.content {
		if !matcher.MatchString(item.name) {
			continue
		}

		maybeCheckmark := " "
		if _, found := model.selectedItemIndexes[idx]; found {
			maybeCheckmark = "✔️"
		}

		colorizedItemName := termenvOutput.String(item.name)
		if idx == model.cursorIdx {
			colorizedItemName = colorizedItemName.Foreground(cursorItemGreen).Bold()
		}

		row := fmt.Sprintf(" %s  %s", maybeCheckmark, colorizedItemName)

		resultBuilder.WriteString(row + "\n")
	}
	resultBuilder.WriteString("\n")

	footerRow := ""
	if len(model.selectedItemIndexes) > 0 {
		selectedItemsText := fmt.Sprintf(" %d items selected", len(model.selectedItemIndexes))
		footerRow = termenvOutput.String(selectedItemsText).Foreground(selectedItemsYellow).String()
	}
	resultBuilder.WriteString(footerRow + "\n")

	return resultBuilder.String()
}

func (model appModel) handleNavigationKeypress(msg tea.KeyMsg) (appModel, tea.Cmd) {
	switch msg.String() {
	case "j":
		newCursorIdx := model.cursorIdx + 1
		if newCursorIdx < len(model.content) {
			model.cursorIdx = newCursorIdx
		}
	case "k":
		newCursorIdx := model.cursorIdx - 1
		if newCursorIdx >= 0 {
			model.cursorIdx = newCursorIdx
		}
	case "x":
		_, found := model.selectedItemIndexes[model.cursorIdx]
		if found {
			delete(model.selectedItemIndexes, model.cursorIdx)
		} else {
			model.selectedItemIndexes[model.cursorIdx] = true
		}
	case "a":
		if len(model.selectedItemIndexes) < len(model.content) {
			for idx := range model.content {
				model.selectedItemIndexes[idx] = true
			}
		} else {
			model.selectedItemIndexes = make(map[int]bool)
		}
	case "c":
		// Clear the filter
		model.filterInput.SetValue("")
	case "/":
		model.mode = filterMode

		// This will tell the input that it should display the cursor
		cmd := model.filterInput.Focus()
		return model, cmd
	}

	return model, nil
}
