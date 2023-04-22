package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/filter_pane"
	"github.com/mieubrisse/cli-journal-go/components/filterable_content_list"
	"github.com/mieubrisse/cli-journal-go/components/form"
	"github.com/mieubrisse/cli-journal-go/data_structures/content_item"
	"github.com/mieubrisse/cli-journal-go/filterable_item_list"
	"github.com/mieubrisse/cli-journal-go/global_styles"
	"github.com/mieubrisse/cli-journal-go/helpers"
	"github.com/mieubrisse/vim-bubble/vim"
	"time"
)

const (
	maxCreateContentModalWidth  = 50
	maxCreateContentModalHeight = 3

	filterPaneHeight = 6
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

var filtersLabelLine = lipgloss.NewStyle().
	Foreground(global_styles.Cyan).
	Bold(true).
	Render("FILTERS")

type Model struct {
	createContentForm form.Model

	filterPane filter_pane.Model

	filterTabCompletionPane filterable_item_list.Model[tabCompletionItem]

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
	filterPane filter_pane.Model,
) Model {
	completionPane := filterable_item_list.New[tabCompletionItem]([]tabCompletionItem{
		{completion: "foo"},
		{completion: "bar bang blork whoa this is so long"},
	})

	return Model{
		createContentForm:       createContentForm,
		filterPane:              filterPane,
		filterTabCompletionPane: completionPane,
		contentList:             contentList,
		height:                  0,
		width:                   0,
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
			case "\\":
				model.contentList = model.contentList.Blur()

				// TODO switch to by-value
				model.filterPane.Focus()
				model.filterPane = model.filterPane.SetMode(vim.InsertMode)

				return model, nil
				/*
					case "#":
						model.contentList = model.contentList.Blur()

						// This will tell the input that it should display the cursor
						var cmd tea.Cmd
						model.tagFilterInput, cmd = model.tagFilterInput.Focus()

						return model, cmd

				*/
			case "c":
				// Clear all filters
				/*
					model.nameFilterInput = model.nameFilterInput.SetValue("")
					model.tagFilterInput = model.tagFilterInput.SetValue("")

				*/
				model.filterPane = model.filterPane.Clear()
				nameFilterLines, tagFilterLines := model.filterPane.GetFilterLines()
				model.contentList = model.contentList.SetFilters(nameFilterLines, tagFilterLines)

				return model, nil
			case "n":
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

		} else if model.filterPane.Focused() {
			if msg.String() == "\\" {
				// TODO switch to by-value
				model.filterPane.Blur()

				model.contentList = model.contentList.Focus()

				model = model.resizeContentListAndFilterPane()

				return model, nil
			}

			var cmd tea.Cmd
			model.filterPane, cmd = model.filterPane.Update(msg)

			// Make sure to let the content list know about the changes
			nameFilterList, tagFilterList := model.filterPane.GetFilterLines()
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

	filterView := lipgloss.JoinHorizontal(
		lipgloss.Center,
		model.filterPane.View(),
		model.filterTabCompletionPane.View(),
	)

	sections := []string{
		model.contentList.View(),
		filtersLabelLine,
		filterView,
	}

	contents := lipgloss.JoinVertical(lipgloss.Left, sections...)

	horizontalPad, verticalPad := getPadsForSize(model.width, model.height)
	result := lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		MaxWidth(model.width).
		MaxHeight(model.height).
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
	displaySpaceWidth := helpers.GetMaxInt(0, model.width-2*horizontalPad)
	displaySpaceHeight := helpers.GetMaxInt(0, model.height-2*verticalPad)

	filterPaneWidth := int(0.5 * float64(displaySpaceWidth))
	model.filterPane = model.filterPane.Resize(filterPaneWidth, filterPaneHeight)

	completionPaneWidth := displaySpaceWidth - filterPaneWidth
	model.filterTabCompletionPane = model.filterTabCompletionPane.Resize(completionPaneWidth, filterPaneHeight)

	/*
		model.nameFilterInput = model.nameFilterInput.Resize(displaySpaceWidth, 1)
		model.tagFilterInput = model.tagFilterInput.Resize(displaySpaceWidth, 1)

		contentListHeight := helpers.GetMaxInt(
			0,
			displaySpaceHeight-model.nameFilterInput.GetHeight()-model.tagFilterInput.GetHeight(),
		)
	*/

	// Leave one blank line for filters label
	contentListHeight := helpers.GetMaxInt(0, displaySpaceHeight-filterPaneHeight-1)

	model.contentList = model.contentList.Resize(displaySpaceWidth, contentListHeight)

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

func (model Model) resizeContentListAndFilterPane() Model {
	return model
}

func clampInt(value int, min int, max int) int {
	if max < min {
		max, min = min, max
	}

	return helpers.GetMaxInt(min, helpers.GetMinInt(max, value))
}
