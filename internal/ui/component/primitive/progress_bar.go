package primitive

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type ProgressBarProps struct {
	BaseProps // StyleProps + LayoutProps

	// Progress value (0.0 to 1.0)
	Value float64

	// Visual style: "bar", "blocks", "dots", "gradient"
	Variant string

	// Color scheme
	Color       style.Color
	EmptyColor  style.Color
	BorderColor style.Color

	// Label options
	ShowLabel      bool   // Show "50%" text
	ShowPercentage bool   // Show percentage inside bar
	Label          string // Custom label instead of percentage

	// Animation (for indeterminate progress)
	Indeterminate bool
	AnimFrame     int // Current animation frame
}

// ProgressBar creates a progress bar component
func ProgressBar(props ProgressBarProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		width := 30
		if props.Width != nil {
			width = *props.Width
		}

		// Clamp value between 0 and 1
		value := props.Value
		if value < 0 {
			value = 0
		}
		if value > 1 {
			value = 1
		}

		variant := props.Variant
		if variant == "" {
			variant = "bar"
		}

		// Colors
		fillColor := props.Color
		if fillColor == nil {
			fillColor = style.Primary
		}

		emptyColor := props.EmptyColor
		if emptyColor == nil {
			emptyColor = style.Grey300
		}

		borderColor := props.BorderColor
		if borderColor == nil {
			borderColor = style.Grey400
		}

		var bar string

		switch variant {
		case "blocks":
			bar = renderBlocksProgress(value, width, fillColor, emptyColor)
		case "dots":
			bar = renderDotsProgress(value, width, fillColor, emptyColor)
		case "gradient":
			bar = renderGradientProgress(value, width)
		default:
			bar = renderBarProgress(value, width, fillColor, emptyColor, props.Indeterminate, props.AnimFrame)
		}

		// Add percentage label
		if props.ShowPercentage || props.ShowLabel {
			labelText := props.Label
			if labelText == "" && props.ShowPercentage {
				labelText = fmt.Sprintf("%.0f%%", value*100)
			}

			if labelText != "" {
				labelStyle := lipgloss.NewStyle().
					Foreground(style.Grey700).
					Bold(true)
				bar = bar + " " + labelStyle.Render(labelText)
			}
		}

		// Apply custom style
		if props.Style != nil {
			bar = props.Style.Render(bar)
		}

		return bar
	})
}

func renderBarProgress(value float64, width int, fillColor, emptyColor style.Color, indeterminate bool, frame int) string {
	if indeterminate {
		barWidth := width - 2
		chunkSize := barWidth / 5
		pos := (frame % (barWidth + chunkSize)) - chunkSize

		bar := make([]string, barWidth)
		for i := 0; i < barWidth; i++ {
			if i >= pos && i < pos+chunkSize {
				bar[i] = lipgloss.NewStyle().Foreground(fillColor).Render("█")
			} else {
				bar[i] = lipgloss.NewStyle().Foreground(emptyColor).Render("░")
			}
		}

		content := strings.Join(bar, "")
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(emptyColor).
			Render(content)
	}

	filled := int(math.Round(float64(width-2) * value))
	empty := (width - 2) - filled

	filledBar := lipgloss.NewStyle().Foreground(fillColor).Render(strings.Repeat("█", filled))
	emptyBar := lipgloss.NewStyle().Foreground(emptyColor).Render(strings.Repeat("░", empty))

	content := filledBar + emptyBar

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(emptyColor).
		Render(content)
}

func renderBlocksProgress(value float64, width int, fillColor, emptyColor style.Color) string {
	blocks := width / 2
	filled := int(math.Round(float64(blocks) * value))

	result := ""
	for i := 0; i < blocks; i++ {
		if i < filled {
			result += lipgloss.NewStyle().Foreground(fillColor).Render("█ ")
		} else {
			result += lipgloss.NewStyle().Foreground(emptyColor).Render("▯ ")
		}
	}

	return strings.TrimSpace(result)
}

func renderDotsProgress(value float64, width int, fillColor, emptyColor style.Color) string {
	dots := width / 2
	filled := int(math.Round(float64(dots) * value))

	result := ""
	for i := 0; i < dots; i++ {
		if i < filled {
			result += lipgloss.NewStyle().Foreground(fillColor).Render("● ")
		} else {
			result += lipgloss.NewStyle().Foreground(emptyColor).Render("○ ")
		}
	}

	return strings.TrimSpace(result)
}

func renderGradientProgress(value float64, width int) string {
	filled := int(math.Round(float64(width-2) * value))
	empty := (width - 2) - filled

	gradientChars := []string{"█", "▓", "▒", "░"}

	result := ""
	for i := 0; i < filled; i++ {
		idx := int(float64(i) / float64(filled) * float64(len(gradientChars)-1))
		if idx >= len(gradientChars) {
			idx = len(gradientChars) - 1
		}
		result += lipgloss.NewStyle().Foreground(style.Primary).Render(gradientChars[idx])
	}

	result += strings.Repeat("░", empty)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.Grey400).
		Render(result)
}

type CircularProgressProps struct {
	StyleProps

	// Progress value (0.0 to 1.0)
	Value float64

	// Size (diameter in characters)
	Size int

	// Show percentage in center
	ShowPercentage bool

	// Colors
	Color      style.Color
	EmptyColor style.Color
}

// CircularProgress creates a circular progress indicator
func CircularProgress(props CircularProgressProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		value := props.Value
		if value < 0 {
			value = 0
		}
		if value > 1 {
			value = 1
		}

		size := props.Size
		if size == 0 {
			size = 5
		}

		color := props.Color
		if color == nil {
			color = style.Primary
		}

		emptyColor := props.EmptyColor
		if emptyColor == nil {
			emptyColor = style.Grey300
		}

		percentage := int(value * 100)

		if size <= 5 {
			// Small circular indicator
			var char string
			if value >= 0.875 {
				char = "◉"
			} else if value >= 0.625 {
				char = "◕"
			} else if value >= 0.375 {
				char = "◔"
			} else if value >= 0.125 {
				char = "◑"
			} else {
				char = "○"
			}

			result := lipgloss.NewStyle().
				Foreground(color).
				Bold(true).
				Render(char)

			if props.ShowPercentage {
				percentText := lipgloss.NewStyle().
					Foreground(style.Grey700).
					Render(fmt.Sprintf(" %d%%", percentage))
				result += percentText
			}

			return result
		}

		// Larger circular progress (ASCII art)
		lines := []string{
			"  ███  ",
			" █   █ ",
			"█     █",
			" █   █ ",
			"  ███  ",
		}

		// Color based on progress
		styledLines := make([]string, len(lines))
		for i, line := range lines {
			styledLines[i] = lipgloss.NewStyle().Foreground(color).Render(line)
		}

		result := strings.Join(styledLines, "\n")

		// Add percentage below if requested
		if props.ShowPercentage {
			percentText := lipgloss.NewStyle().
				Foreground(style.Grey700).
				Bold(true).
				Align(lipgloss.Center).
				Width(lipgloss.Width(lines[0])).
				Render(fmt.Sprintf("%d%%", percentage))
			result += "\n" + percentText
		}

		if props.Style != nil {
			result = props.Style.Render(result)
		}

		return result
	})
}
