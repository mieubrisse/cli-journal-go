package entry_item

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/global_styles"
	"github.com/mieubrisse/cli-journal-go/helpers"
	"strings"
	"time"
)

type componentSize int

const (
	contentTimestampFormat = "2006-01-02 15:04:05"

	checkmarkChar = '•'

	// Used when a line is too small
	continuationChar = '…'

	maxNameWidth = 45

	minimumNameAndTagWidth = 5

	wide componentSize = iota
	medium
	narrow
	sliver
)

// Minimum width, in characters, for the component to be classed as each size
var componentSizeThresholds = map[componentSize]int{
	wide:   150,
	medium: 120,
	narrow: 80,
	sliver: 0,
}
var desiredCheckmarkWidthsByComponentSize = map[componentSize]int{
	wide:   5,
	medium: 4,
	narrow: 3,
	sliver: 2,
}
var timestampWidthsByComponentSize = map[componentSize]int{
	wide:   len(contentTimestampFormat) + 4,
	medium: len(contentTimestampFormat) + 2,
	narrow: 0,
	sliver: 0,
}

type implementation struct {
	// TODO something about the value

	timestamp time.Time
	name      string
	tags      []string // Maybe make this a map??

	isHiglighted bool
	isSelected   bool

	width  int
	height int
}

func New(timestamp time.Time, name string, tags []string) Component {
	return &implementation{
		timestamp:    timestamp,
		name:         name,
		tags:         tags,
		isHiglighted: false,
		isSelected:   false,
		width:        0,
		height:       0,
	}
}

func (impl implementation) View() string {
	return impl.render()
}

func (impl implementation) Resize(width int, height int) {
	impl.width = width
	impl.height = height
}

func (impl implementation) GetTimestamp() time.Time {
	return impl.timestamp
}

func (impl implementation) GetName() string {
	return impl.name
}

func (impl implementation) GetTags() []string {
	return impl.tags
}

func (impl implementation) GetWidth() int {
	return impl.width
}

func (impl implementation) GetHeight() int {
	return impl.height
}

func (impl implementation) GetValue() string {
	panic("Implement me!")
}

func (impl implementation) SetHighlighted(isHighlighted bool) {
	impl.isHiglighted = isHighlighted
}

func (impl *implementation) SetSelection(isSelected bool) {
	impl.isSelected = isSelected
}

func (impl implementation) IsSelected() bool {
	return impl.isSelected
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
// TODO allow this to do wrapping
func (impl implementation) render() string {
	baseLineStyle := lipgloss.NewStyle()
	if impl.isHiglighted {
		baseLineStyle = baseLineStyle.Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
	}

	// Calculate the widths for the various components
	biggestThresholdPassed := sliver
	for trialComponentSize, threshold := range componentSizeThresholds {
		if impl.width > threshold && threshold > componentSizeThresholds[biggestThresholdPassed] {
			biggestThresholdPassed = trialComponentSize
		}
	}
	actualComponentSize := biggestThresholdPassed
	desiredCheckmarkWidth, found := desiredCheckmarkWidthsByComponentSize[actualComponentSize]
	if !found {
		panic("No checkmark width for terminal size")
	}
	timestampWidth, found := timestampWidthsByComponentSize[actualComponentSize]
	if !found {
		panic("No timestamp width for terminal size")
	}

	widthRemaining := helpers.GetMaxInt(0, impl.width-desiredCheckmarkWidth-timestampWidth)
	// Safety valve: if we don't have at least 10 characters, don't even bother

	nameWidth := helpers.GetMinInt(
		maxNameWidth,
		int(0.6*float64(widthRemaining)),
	)
	tagsWidth := helpers.GetMaxInt(0, widthRemaining-nameWidth)

	// Checkmark string
	checkmarkStr := ""
	if impl.isSelected {
		checkmarkStr = string(checkmarkChar)
	}
	checkmarkStr = baseLineStyle.Copy().
		Foreground(global_styles.Orange).
		Width(desiredCheckmarkWidth).
		AlignHorizontal(lipgloss.Center).
		Render(checkmarkStr)

	// Timestamp (disabled if too small)
	timestampStr := ""
	if timestampWidth > 0 {
		timestampStr = impl.timestamp.Format(contentTimestampFormat)
		timestampStr = baseLineStyle.Copy().
			Foreground(global_styles.Cyan).
			Width(timestampWidth).
			AlignHorizontal(lipgloss.Left).
			Render(timestampStr)
	}

	// Name
	nameStr := ""
	if nameWidth > minimumNameAndTagWidth {
		nameStr = impl.name
		nameLen := len(nameStr)
		if nameLen > nameWidth-1 {
			nameStr = nameStr[:nameWidth-2] + string(continuationChar)
		}
		nameStr = baseLineStyle.Copy().
			Foreground(global_styles.White).
			Width(nameWidth).
			AlignHorizontal(lipgloss.Left).
			Render(nameStr)
	}

	tagsStr := ""
	if tagsWidth > minimumNameAndTagWidth {
		tagsStr = strings.Join(impl.tags, " ")
		tagsLen := len(tagsStr)
		if tagsLen > tagsWidth-1 {
			tagsStr = tagsStr[:tagsWidth-2] + string(continuationChar)
		}
		tagsStr = baseLineStyle.Copy().
			Foreground(global_styles.Red).
			Width(tagsWidth).
			AlignHorizontal(lipgloss.Left).
			Render(tagsStr)
	}

	line := lipgloss.JoinHorizontal(
		lipgloss.Top,
		checkmarkStr,
		timestampStr,
		nameStr,
		tagsStr,
	)

	return baseLineStyle.Copy().
		Width(impl.width).
		MaxWidth(impl.width).
		Render(line)
}
