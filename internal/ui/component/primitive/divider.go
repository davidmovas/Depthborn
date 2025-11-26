package primitive

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type DividerProps struct {
	CommonProps
	Length *int
	Char   *string
	Label  *string
}

func Divider(props DividerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		length := 20
		if props.Length != nil {
			length = *props.Length
		}

		char := "─"
		if props.Char != nil {
			char = *props.Char
		}

		dividerStyle := style.Merge(
			style.Fg(style.Grey400),
			style.Dim,
		)

		if props.Label != nil {
			return renderLabeledDivider(*props.Label, length, char, dividerStyle, props)
		}

		content := strings.Repeat(char, length)

		if props.Style != nil {
			dividerStyle = dividerStyle.Inherit(*props.Style)
		}

		return dividerStyle.Render(content)
	})
}

func renderLabeledDivider(label string, length int, char string, baseStyle style.Style, props DividerProps) string {
	labelStyle := style.Merge(
		style.Fg(style.Grey600),
		style.Bold,
		style.PaddingX1,
	)

	renderedLabel := labelStyle.Render(" " + label + " ")
	labelWidth := style.CalculateWidth(renderedLabel)

	sideLength := (length - labelWidth) / 2
	if sideLength < 1 {
		sideLength = 1
	}

	leftLine := strings.Repeat(char, sideLength)
	rightLine := strings.Repeat(char, sideLength)

	content := leftLine + renderedLabel + rightLine

	dividerStyle := baseStyle
	if props.Style != nil {
		dividerStyle = dividerStyle.Inherit(*props.Style)
	}

	return dividerStyle.Render(content)
}

func VDivider(props DividerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		height := 5
		if props.Length != nil {
			height = *props.Length
		}

		char := "│"
		if props.Char != nil {
			char = *props.Char
		}

		lines := make([]string, height)
		for i := range lines {
			lines[i] = char
		}

		content := strings.Join(lines, "\n")

		dividerStyle := style.Merge(
			style.Fg(style.Grey400),
			style.Dim,
		)

		if props.Style != nil {
			dividerStyle = dividerStyle.Inherit(*props.Style)
		}

		return dividerStyle.Render(content)
	})
}
