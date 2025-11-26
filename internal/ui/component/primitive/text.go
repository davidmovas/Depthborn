package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// Text renders styled text
func Text(props TextProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		s := style.New()
		s = ApplyAllProps(s, props.LayoutProps, props.StyleProps, props.ContentProps)
		return s.Render(props.Content)
	})
}

// Heading renders heading text
func Heading(props TextProps) component.Component {
	props.StyleProps.Style = style.S(
		style.Bold,
		style.Fg(style.Grey900),
		*props.StyleProps.Style,
	)
	return Text(props)
}

// Label renders label text
func Label(props TextProps) component.Component {
	props.StyleProps.Style = style.S(
		style.Fg(style.Grey700),
		*props.StyleProps.Style,
	)
	return Text(props)
}

// Code renders code-style text
func Code(props TextProps) component.Component {
	props.StyleProps.Style = style.S(
		style.Bg(style.Grey200),
		style.Fg(style.Grey800),
		style.Px(1),
		*props.StyleProps.Style,
	)
	return Text(props)
}
