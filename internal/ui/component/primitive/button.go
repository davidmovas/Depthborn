package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// Button renders clickable button
func Button(props InteractiveProps, label string) component.Component {
	textComp := Text(TextProps{
		StyleProps:  props.StyleProps,
		LayoutProps: props.LayoutProps,
		Content:     label,
	})

	if props.FocusProps.OnClick != nil || len(props.FocusProps.Hotkeys) > 0 {
		autoFocus := false
		if props.AutoFocus != nil {
			autoFocus = *props.AutoFocus
		}

		focusStyle := props.FocusStyle
		if focusStyle == nil && props.Style != nil {
			focusStyle = props.Style
		}

		return component.MakeFocusable(textComp, component.FocusableConfig{
			Position:  props.Position,
			Hotkeys:   props.Hotkeys,
			CanFocus:  true,
			AutoFocus: autoFocus,
			IsInput:   false,
			OnActivateCallback: func() bool {
				if props.OnClick != nil {
					props.OnClick()
					return true
				}
				return false
			},
			OnFocusCallback: props.OnFocus,
			OnBlurCallback:  props.OnBlur,
			FocusedStyle: func(content string) string {
				if focusStyle != nil {
					return focusStyle.Render(content)
				}
				return content
			},
		})
	}

	return textComp
}
