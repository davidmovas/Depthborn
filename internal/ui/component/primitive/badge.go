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
	CommonProps
	Content string
	Variant *BadgeVariant
}

func Badge(props BadgeProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		variant := BadgeDefault
		if props.Variant != nil {
			variant = *props.Variant
		}

		badgeStyle := style.Merge(
			style.PaddingX1,
			style.PaddingY0,
			style.BorderRounded,
			style.TextDim,
		)

		switch variant {
		case BadgePrimary:
			badgeStyle = style.Merge(badgeStyle, style.Bg(style.Blue100), style.Fg(style.Blue800), style.BrColor(badgeStyle, style.Blue500))
		case BadgeSuccess:
			badgeStyle = style.Merge(badgeStyle, style.Bg(style.Green100), style.Fg(style.Green800), style.BrColor(badgeStyle, style.Green500))
		case BadgeWarning:
			badgeStyle = style.Merge(badgeStyle, style.Bg(style.Yellow100), style.Fg(style.Yellow800), style.BrColor(badgeStyle, style.Yellow500))
		case BadgeError:
			badgeStyle = style.Merge(badgeStyle, style.Bg(style.Red100), style.Fg(style.Red800), style.BrColor(badgeStyle, style.Red500))
		case BadgeInfo:
			badgeStyle = style.Merge(badgeStyle, style.Bg(style.Cyan100), style.Fg(style.Cyan800), style.BrColor(badgeStyle, style.Cyan500))
		case BadgeOutline:
			badgeStyle = style.Merge(badgeStyle, style.Bg(style.Transparent), style.Fg(style.Grey600), style.BrColor(badgeStyle, style.Grey400))
		default:
			badgeStyle = style.Merge(badgeStyle, style.Bg(style.Grey100), style.Fg(style.Grey700), style.BrColor(badgeStyle, style.Grey400))
		}

		if props.Style != nil {
			badgeStyle = badgeStyle.Inherit(*props.Style)
		}

		return badgeStyle.Render(" " + props.Content + " ")
	})
}
