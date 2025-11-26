package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type BadgeVariant string

const (
	BadgeDefault BadgeVariant = "default"
	BadgePrimary BadgeVariant = "primary"
	BadgeSuccess BadgeVariant = "success"
	BadgeWarning BadgeVariant = "warning"
	BadgeError   BadgeVariant = "error"
	BadgeInfo    BadgeVariant = "info"
	BadgeOutline BadgeVariant = "outline"
)

type BadgeProps struct {
	BaseProps
	Content string
	Variant BadgeVariant
}

func Badge(props BadgeProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		baseStyle := style.Merge(
			style.Px(1),
			style.Py(0),
			style.Br(),
		)

		switch props.Variant {
		case BadgePrimary:
			baseStyle = style.Merge(baseStyle,
				style.Bg(style.Blue100),
				style.Fg(style.Blue800),
				style.BrColor(baseStyle, style.Blue500),
			)
		case BadgeSuccess:
			baseStyle = style.Merge(baseStyle,
				style.Bg(style.Green100),
				style.Fg(style.Green800),
				style.BrColor(baseStyle, style.Green500),
			)
		case BadgeWarning:
			baseStyle = style.Merge(baseStyle,
				style.Bg(style.Yellow100),
				style.Fg(style.Yellow800),
				style.BrColor(baseStyle, style.Yellow500),
			)
		case BadgeError:
			baseStyle = style.Merge(baseStyle,
				style.Bg(style.Red100),
				style.Fg(style.Red800),
				style.BrColor(baseStyle, style.Red500),
			)
		case BadgeInfo:
			baseStyle = style.Merge(baseStyle,
				style.Bg(style.Cyan100),
				style.Fg(style.Cyan800),
				style.BrColor(baseStyle, style.Cyan500),
			)
		case BadgeOutline:
			baseStyle = style.Merge(baseStyle,
				style.Bg(style.Transparent),
				style.Fg(style.Grey600),
				style.BrColor(baseStyle, style.Grey400),
			)
		default:
			baseStyle = style.Merge(baseStyle,
				style.Bg(style.Grey100),
				style.Fg(style.Grey700),
				style.BrColor(baseStyle, style.Grey400),
			)
		}

		s := ApplyAllProps(baseStyle, props.LayoutProps, props.StyleProps, ContentProps{})
		return s.Render(" " + props.Content + " ")
	})
}
