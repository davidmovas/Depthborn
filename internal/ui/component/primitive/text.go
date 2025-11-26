package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

func Text(content string, styles ...style.Style) component.Component {
	return component.Func(func(ctx *component.Context) string {
		return style.Render(content, styles...)
	})
}

func Heading(props TextProps) component.Component {
	headingStyle := style.New().
		Bold(true).
		Foreground(style.Grey900)

	if props.Style != nil {
		headingStyle = headingStyle.Inherit(*props.Style)
	}

	return Text(props.Content, headingStyle)
}

func Label(props TextProps) component.Component {
	labelStyle := style.New().
		Foreground(style.Grey700)

	if props.Style != nil {
		labelStyle = labelStyle.Inherit(*props.Style)
	}

	return Text(props.Content, labelStyle)
}

func Code(props TextProps) component.Component {
	codeStyle := style.New().
		Background(style.Grey200).
		Foreground(style.Grey800).
		Padding(0, 1)

	if props.Style != nil {
		codeStyle = codeStyle.Inherit(*props.Style)
	}

	return Text(props.Content, codeStyle)
}

func Link(props TextProps) component.Component {
	linkStyle := style.New().
		Foreground(style.Blue500).
		Underline(true)

	if props.Style != nil {
		linkStyle = linkStyle.Inherit(*props.Style)
	}

	baseComp := Text(props.Content, linkStyle)

	if props.OnClick != nil {
		return component.MakeFocusable(baseComp, component.FocusableConfig{
			CanFocus:  true,
			AutoFocus: props.AutoFocus != nil && *props.AutoFocus,
			Position:  props.Position,
			OnActivateCallback: func() bool {
				props.OnClick()
				return true
			},
			FocusedStyle: func(content string) string {
				if props.FocusStyle != nil {
					return props.FocusStyle.Render(content)
				}
				return linkStyle.Foreground(style.Blue700).Render(content)
			},
		})
	}

	return baseComp
}
