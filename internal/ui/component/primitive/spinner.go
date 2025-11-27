package primitive

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type SpinnerProps struct {
	StyleProps

	// Spinner style: "dots", "line", "arc", "circle", "bounce"
	Variant string

	// Current animation frame
	Frame int

	// Color
	Color style.Color

	// Label
	Label string
}

// Spinner creates an animated loading spinner
func Spinner(props SpinnerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		variant := props.Variant
		if variant == "" {
			variant = "dots"
		}

		color := props.Color
		if color == nil {
			color = style.Primary
		}

		frame := props.Frame % getSpinnerFrameCount(variant)

		var spinner string

		switch variant {
		case "line":
			frames := []string{"—", "\\", "|", "/"}
			spinner = frames[frame]
		case "arc":
			frames := []string{"◜", "◝", "◞", "◟"}
			spinner = frames[frame]
		case "circle":
			frames := []string{"◐", "◓", "◑", "◒"}
			spinner = frames[frame]
		case "bounce":
			frames := []string{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"}
			spinner = frames[frame]
		default: // dots
			frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			spinner = frames[frame]
		}

		spinnerStyled := lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Render(spinner)

		if props.Label != "" {
			labelStyle := lipgloss.NewStyle().Foreground(style.Grey700)
			spinnerStyled += " " + labelStyle.Render(props.Label)
		}

		if props.Style != nil {
			spinnerStyled = props.Style.Render(spinnerStyled)
		}

		return spinnerStyled
	})
}

func getSpinnerFrameCount(variant string) int {
	switch variant {
	case "line":
		return 4
	case "arc":
		return 4
	case "circle":
		return 4
	case "bounce":
		return 8
	default:
		return 10
	}
}
