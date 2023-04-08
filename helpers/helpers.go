package helpers

import (
	"github.com/acarl005/stripansi"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var overlayBackgroundStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))

// We don't want the modal to be too far down the screen if it's really small, so for small modals this is the maximum distance
// from the top of the screen that the modal will be displayed (expressed as a % of the total height of the background)
const maxModalTopOffsetPercentage = 0.3

func GetMaxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetMinInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// TODO FIX BUG IN THIS - we can't just overlay using string length, because non-printing chars mess this up (e.g. color chars)
// TODO add a way to dim the background
func OverlayString(background string, overlay string) string {
	backgroundWidth, backgroundHeight := lipgloss.Size(background)
	overlayWidth, overlayHeight := lipgloss.Size(overlay)

	maxFirstOverlaidLineIdx := int(float64(backgroundHeight) * maxModalTopOffsetPercentage)

	// The index of the first line that will suffer replacement
	// We don't want it too far down the screen, so we have a cap
	firstOverlaidLineIdx := GetMinInt(
		(backgroundHeight/2)-(overlayHeight/2),
		maxFirstOverlaidLineIdx,
	)

	// The index of the line that switches back to being background again
	resumeLineIdx := firstOverlaidLineIdx + overlayHeight

	// The index of the first column that will be replaced
	cutpointIdx := (backgroundWidth / 2) - (overlayWidth / 2)

	// The index of the column that switches back to being background again
	resumepointIdx := backgroundWidth - cutpointIdx

	overlayLines := strings.Split(overlay, "\n")
	backgroundLines := strings.Split(background, "\n")

	resultLines := []string{}
	for idx, backgroundLine := range backgroundLines {
		if idx < firstOverlaidLineIdx || idx >= resumeLineIdx {
			resultLines = append(resultLines, backgroundLine)
			continue
		}

		overlayLine := overlayLines[idx-firstOverlaidLineIdx]

		resultLine := backgroundLine[:cutpointIdx] + overlayLine + backgroundLine[resumepointIdx:]
		resultLines = append(resultLines, resultLine)
	}

	return strings.Join(resultLines, "\n")
}

// TODO use this, somehow, to do background-dimming
// Removes all non-graphic characters
func ClearFormatting(input string) string {
	return stripansi.Strip(input)
}
