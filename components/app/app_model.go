package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/filter_pane"
	"github.com/mieubrisse/cli-journal-go/components/filterable_content_list"
	"github.com/mieubrisse/cli-journal-go/components/form"
	"github.com/mieubrisse/cli-journal-go/data_structures/content_item"
	"github.com/mieubrisse/cli-journal-go/helpers"
	"time"
)

const (
	maxCreateContentModalWidth  = 50
	maxCreateContentModalHeight = 3

	minFilterInputHeight = 5
	maxFilterInputHeight = 10
)

// "Constants"
var horizontalPadThresholdsByTerminalWidth = map[int]int{
	0:   0,
	60:  1,
	120: 2,
}
var verticalPadThresholdsByTerminalHeight = map[int]int{
	0:  0,
	40: 1,
}

type Model struct {
	createContentForm form.Model

	filterInput filter_pane.Model

	/*
		nameFilterInput text_input.Model

		tagFilterInput text_input.Model
	*/

	contentList filterable_content_list.Model

	height int
	width  int
}

func New(
	createContentForm form.Model,
	contentList filterable_content_list.Model,
	filterInput filter_pane.Model,
) Model {
	return Model{
		createContentForm: createContentForm,
		filterInput:       filterInput,
		contentList:       contentList,
		height:            0,
		width:             0,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

// NOTE: This returns a model because BubbleTea expects models to be passed by-value, so the way to "update" the model
// is to return a new instance of it
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return model, tea.Quit
		}

		if model.contentList.Focused() {
			switch msg.String() {
			case "/":
				model.contentList = model.contentList.Blur()

				model.filterInput.Focus()

				return model, nil
				/*
					case "#":
						model.contentList = model.contentList.Blur()

						// This will tell the input that it should display the cursor
						var cmd tea.Cmd
						model.tagFilterInput, cmd = model.tagFilterInput.Focus()

						return model, cmd

				*/
			case "esc":
				// Clear all filters
				/*
					model.nameFilterInput = model.nameFilterInput.SetValue("")
					model.tagFilterInput = model.tagFilterInput.SetValue("")

				*/
				model.filterInput = model.filterInput.Clear()
				nameFilterLines, tagFilterLines := model.filterInput.GetFilterLines()
				model.contentList.SetFilters(nameFilterLines, tagFilterLines)

				return model, nil
			case "c":
				model.contentList = model.contentList.Blur()

				var cmd tea.Cmd
				model.createContentForm, cmd = model.createContentForm.Focus()
				return model, cmd
			case "d":
				model.contentList = model.contentList.Blur()

				var cmd tea.Cmd
				model.createContentForm, cmd = model.createContentForm.Focus()
				return model, cmd
			}

			var cmd tea.Cmd
			model.contentList, cmd = model.contentList.Update(msg)
			return model, cmd

		} else if model.filterInput.Focused() {
			if model.filterInput.IsInNormalMode() && msg.String() == "esc" {
				model.filterInput.Blur()
				model.contentList = model.contentList.Focus()
			}

			var cmd tea.Cmd
			model.filterInput, cmd = model.filterInput.Update(msg)

			// Make sure to let the content list know about the changes
			nameFilterList, tagFilterList := model.filterInput.GetFilterLines()
			model.contentList = model.contentList.SetFilters(nameFilterList, tagFilterList)

			return model, cmd
		} else if model.createContentForm.Focused() {
			switch msg.String() {
			case "esc":
				// Back out of the create content modal
				model.createContentForm = model.createContentForm.Clear().
					Blur()
				model.contentList = model.contentList.Focus()
				return model, nil
			case "enter":
				// Create the new content piece
				content := content_item.ContentItem{
					Timestamp: time.Now(),
					Name:      model.createContentForm.GetValue(),
					Tags:      []string{},
				}
				model.contentList = model.contentList.AddItem(content).
					Focus()
				model.createContentForm = model.createContentForm.Clear().
					Blur()

				return model, nil
			}

			var cmd tea.Cmd
			model.createContentForm, cmd = model.createContentForm.Update(msg)
			return model, cmd
		}
	case tea.WindowSizeMsg:
		return model.Resize(msg.Width, msg.Height), nil
	}

	return model, nil
}

func (model Model) View() string {
	/*
		sections := []string{
			model.contentList.View(),
			model.nameFilterInput.View(),
			model.tagFilterInput.View(),
		}
	*/

	sections := []string{
		model.contentList.View(),
	}
	if model.filterInput.Focused() {
		sections = append(sections, model.filterInput.View())
	}

	contents := lipgloss.JoinVertical(lipgloss.Left, sections...)

	horizontalPad, verticalPad := getPadsForSize(model.width, model.height)
	result := lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		Padding(verticalPad, horizontalPad, verticalPad, horizontalPad).
		Render(contents)

	if model.createContentForm.Focused() {
		createContentFormStr := model.createContentForm.View()
		/*
			createContentModalStr = lipgloss.NewStyle().
				Border()
		*/

		createContentFormStr = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Render(createContentFormStr)

		result = helpers.OverlayString(result, createContentFormStr)
	}

	return result
}

func (model Model) Resize(width int, height int) Model {
	model.width = width
	model.height = height

	horizontalPad, verticalPad := getPadsForSize(model.width, model.height)
	componentSpaceWidth := helpers.GetMaxInt(0, model.width-2*horizontalPad)
	componentSpaceHeight := helpers.GetMaxInt(0, model.height-2*verticalPad)

	// Resize filter pane
	filterText := model.filterInput.GetValue()
	filterTextHeight := lipgloss.Height(filterText)
	filterPaneHeight := clampInt(filterTextHeight, minFilterInputHeight, maxFilterInputHeight)
	model.filterInput.Resize(width, filterPaneHeight)

	/*
		model.nameFilterInput = model.nameFilterInput.Resize(componentSpaceWidth, 1)
		model.tagFilterInput = model.tagFilterInput.Resize(componentSpaceWidth, 1)

		contentListHeight := helpers.GetMaxInt(
			0,
			componentSpaceHeight-model.nameFilterInput.GetHeight()-model.tagFilterInput.GetHeight(),
		)
	*/

	contentListHeight := componentSpaceHeight
	if model.filterInput.Focused() {
		contentListHeight = helpers.GetMaxInt(
			0,
			// TODO add space for a buffer line
			componentSpaceHeight-model.filterInput.GetHeight(),
		)
	}

	model.contentList = model.contentList.Resize(componentSpaceWidth, contentListHeight)

	createContentModalWidth := helpers.GetMinInt(model.width, maxCreateContentModalWidth)
	createContentModalHeight := helpers.GetMinInt(model.height, maxCreateContentModalHeight)
	model.createContentForm = model.createContentForm.Resize(createContentModalWidth, createContentModalHeight)

	return model
}

// =================================== Private Helper Functions ===================================
func getPadsForSize(width int, height int) (int, int) {
	actualHorizontalPad := 0
	for threshold, trialHorizontalPad := range horizontalPadThresholdsByTerminalWidth {
		if width > threshold && actualHorizontalPad < trialHorizontalPad {
			actualHorizontalPad = trialHorizontalPad
		}

	}

	actualVerticalPad := 0
	for threshold, trialVerticalPad := range verticalPadThresholdsByTerminalHeight {
		if height > threshold && actualVerticalPad < trialVerticalPad {
			actualVerticalPad = trialVerticalPad
		}
	}

	return actualHorizontalPad, actualVerticalPad
}

func clampInt(value int, min int, max int) int {
	if max < min {
		max, min = min, max
	}

	return helpers.GetMaxInt(min, helpers.GetMinInt(max, value))
}
