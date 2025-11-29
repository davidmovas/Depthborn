package primitive

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type DividerProps struct {
	BaseProps
	Vertical bool
	Label    string
}

func Divider(props DividerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		char := "─"
		if props.Vertical {
			char = "│"
		}

		width := props.LayoutProps.Width
		if width <= 0 {
			width = 1
		}

		line := strings.Repeat(char, width)

		if props.Label != "" && !props.Vertical {
			labelLen := len(props.Label)
			if labelLen+4 < width {
				leftLen := (width - labelLen - 2) / 2
				rightLen := width - labelLen - 2 - leftLen
				line = strings.Repeat(char, leftLen) + " " + props.Label + " " + strings.Repeat(char, rightLen)
			}
		}

		baseStyle := style.Fg(style.Grey400)
		s := ApplyAllProps(baseStyle, props.LayoutProps, props.StyleProps, ContentProps{})

		return s.Render(line)
	})
}
