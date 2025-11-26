package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// Link renders clickable link
func Link(props InteractiveProps, label string) component.Component {
	baseStyle := style.Merge(
		style.Fg(style.Blue500),
		style.Underline,
	)

	if props.StyleProps.Style != nil {
		baseStyle = baseStyle.Inherit(*props.StyleProps.Style)
	}
	props.StyleProps.Style = &baseStyle

	focusStyle := baseStyle.Foreground(style.Blue700)
	props.FocusStyle = &focusStyle

	return Button(props, label)
}
