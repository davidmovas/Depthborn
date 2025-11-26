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

		if props.Width != nil {
			width = *props.Width
		}
		if props.Height != nil {
			height = *props.Height
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
	w := width
	return Spacer(LayoutProps{Width: &w})
}

// VSpacer renders vertical spacer
func VSpacer(height int) component.Component {
	h := height
	return Spacer(LayoutProps{Height: &h})
}
