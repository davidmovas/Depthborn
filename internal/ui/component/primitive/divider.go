package primitive

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type DividerProps struct {
	BaseProps
	Vertical *bool
	Label    *string
}

func Divider(props DividerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		isVertical := false
		if props.Vertical != nil {
			isVertical = *props.Vertical
		}

		char := "─"
		if isVertical {
			char = "│"
		}

		width := 1
		if props.LayoutProps.Width != nil {
			width = *props.LayoutProps.Width
		}

		line := strings.Repeat(char, width)

		if props.Label != nil && !isVertical {
			labelLen := len(*props.Label)
			if labelLen+4 < width {
				leftLen := (width - labelLen - 2) / 2
				rightLen := width - labelLen - 2 - leftLen
				line = strings.Repeat(char, leftLen) + " " + *props.Label + " " + strings.Repeat(char, rightLen)
			}
		}

		baseStyle := style.Fg(style.Grey400)
		s := ApplyAllProps(baseStyle, props.LayoutProps, props.StyleProps, ContentProps{})

		return s.Render(line)
	})
}
