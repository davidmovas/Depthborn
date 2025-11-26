package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type ButtonProps struct {
	// Required fields
	Label string

	// Optional fields
	Hotkeys   []string
	Position  *component.FocusPosition
	OnClick   func()
	OnFocus   func(string) string
	AutoFocus *bool
	Class     *string

	// Styling
	Style      *style.Style
	FocusStyle *style.Style
}

// Button creates a focusable button component
func Button(props ButtonProps) component.Component {
	hotkeys := props.Hotkeys

	autoFocus := false
	if props.AutoFocus != nil {
		autoFocus = *props.AutoFocus
	}

	var baseComp component.Component
	if props.Style != nil {
		baseComp = Text(props.Label, *props.Style)
	} else {
		baseComp = Text(props.Label)
	}

	var focusedStyleFunc func(string) string
	if props.OnFocus != nil {
		focusedStyleFunc = props.OnFocus
	} else if props.FocusStyle != nil {
		focusedStyleFunc = func(content string) string {
			return props.FocusStyle.Render(content)
		}
	} else if props.Style != nil {
		focusedStyleFunc = func(content string) string {
			return props.Style.Render(content)
		}
	} else {
		focusedStyleFunc = func(content string) string {
			return content
		}
	}

	return component.MakeFocusable(baseComp, component.FocusableConfig{
		Position:  props.Position,
		Hotkeys:   hotkeys,
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

		FocusedStyle: focusedStyleFunc,
	})
}
