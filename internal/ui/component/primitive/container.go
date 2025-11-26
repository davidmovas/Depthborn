package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// Box is a generic container with layout control
func Box(props ContainerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		content := RenderChildren(ctx, props.Children)

		s := style.New()
		s = ApplyAllProps(s, props.LayoutProps, props.StyleProps, props.ContentProps)

		return s.Render(content)
	})
}

// Panel renders styled panel container
func Panel(props ContainerProps) component.Component {
	baseStyle := style.Merge(
		style.Bg(style.Grey100),
		style.P(1),
		style.Br(),
	)

	if props.StyleProps.Style != nil {
		baseStyle = baseStyle.Inherit(*props.StyleProps.Style)
	}
	props.StyleProps.Style = &baseStyle

	return Box(props)
}

// Card renders card container
func Card(props ContainerProps) component.Component {
	baseStyle := style.Merge(
		style.Bg(style.White),
		style.P(2),
		style.Br(),
		style.BrColor(style.New(), style.Grey300),
	)

	if props.StyleProps.Style != nil {
		baseStyle = baseStyle.Inherit(*props.StyleProps.Style)
	}
	props.StyleProps.Style = &baseStyle

	return Box(props)
}
