package primitive

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// Spacer renders empty space
func Spacer(props LayoutProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		width := 1
		height := 1

		if props.Width > 0 {
			width = props.Width
		}
		if props.Height > 0 {
			height = props.Height
		}

		content := strings.Repeat(" ", width)
		if height > 1 {
			lines := make([]string, height)
			for i := range lines {
				lines[i] = content
			}
			content = strings.Join(lines, "\n")
		}

		return content
	})
}

// HSpacer renders horizontal spacer
func HSpacer(width int) component.Component {
	return Spacer(LayoutProps{Width: width})
}

// VSpacer renders vertical spacer
func VSpacer(height int) component.Component {
	return Spacer(LayoutProps{Height: height})
}
