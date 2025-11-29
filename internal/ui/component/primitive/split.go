package primitive

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type SplitProps struct {
	ContainerProps
	Left  component.Component
	Right component.Component
	Ratio *float64 // 0.0 to 1.0, default 0.5
	Gap   *int
}

// Split renders two-pane split layout
func Split(props SplitProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		if props.Left == nil || props.Right == nil {
			return ""
		}

		ratio := 0.5
		if props.Ratio != nil {
			ratio = *props.Ratio
			if ratio < 0 {
				ratio = 0
			}
			if ratio > 1 {
				ratio = 1
			}
		}

		gap := 1
		if props.Gap != nil {
			gap = *props.Gap
		}

		width := 80
		if props.LayoutProps.Width > 0 {
			width = props.LayoutProps.Width
		}

		leftWidth := int(float64(width-gap) * ratio)
		rightWidth := width - leftWidth - gap

		leftStyle := style.W(leftWidth)
		rightStyle := style.W(rightWidth)

		leftContent := leftStyle.Render(props.Left.Render(ctx))
		rightContent := rightStyle.Render(props.Right.Render(ctx))

		gapStr := strings.Repeat(" ", gap)
		result := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftContent,
			gapStr,
			rightContent,
		)

		s := style.New()
		s = ApplyAllProps(s, props.LayoutProps, props.StyleProps, props.ContentProps)

		return s.Render(result)
	})
}
