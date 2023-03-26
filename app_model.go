package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/filterable_content_list"
	"github.com/mieubrisse/cli-journal-go/global_styles"
)

type Mode int

const (
	navigationMode Mode = iota

	filterMode Mode = iota

	// TODO add mode

	numListElems = 20
)

var appStyle = lipgloss.NewStyle().Padding(1)
var componentStyle = lipgloss.NewStyle().Margin(2)

type appModel struct {
	mode Mode

	// filterInput tea.Model
	filterInput textinput.Model

	contentList filterable_content_list.Model

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
			switch msg.String() {
			case "/":
				model.mode = filterMode

				// TODO handle a command coming out the other side?
				model.contentList.Blur()

				// This will tell the input that it should display the cursor
				cmd := model.filterInput.Focus()
				return model, cmd
			case "c":
				model.filterInput.SetValue("")
				model.contentList.UpdateNameFilterText(model.filterInput.Value())
			}

			var cmd tea.Cmd
			model.contentList, cmd = model.contentList.Update(msg)
			return model, cmd
		case filterMode:
			// Back out of filter mode
			if msg.String() == "esc" || msg.String() == "enter" {
				model.filterInput.Blur()
				model.contentList.Focus()
				model.mode = navigationMode
				return model, nil
			}

			var cmd tea.Cmd
			model.filterInput, cmd = model.filterInput.Update(msg)

			// Make sure to tell the content list about the new filter update
			model.contentList.UpdateNameFilterText(model.filterInput.Value())

			return model, cmd
		}
	case tea.WindowSizeMsg:
		// TODO probably, pass this downwards
		model.height = msg.Height
		model.width = msg.Width
	}

	return model, nil
}

func (model appModel) View() string {
	sections := []string{}

	// Ideally, the textinput would have the ability to mutate background color based on focus state, but it doesn't
	// so we have to hack it in here
	filterInputTextStyle := lipgloss.NewStyle()

	if model.mode == filterMode {
		filterInputTextStyle = filterInputTextStyle.
			Bold(true).
			Background(global_styles.FocusedComponentBackgroundColor)
	}
	model.filterInput.TextStyle = filterInputTextStyle
	model.filterInput.PromptStyle = filterInputTextStyle
	sections = append(sections, model.filterInput.View())

	sections = append(sections, model.contentList.View())

	contents := componentStyle.Render(
		model.filterInput.View(),
		model.contentList.View(),
	)

	return appStyle.Render(contents)
}
