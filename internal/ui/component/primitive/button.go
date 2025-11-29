package primitive

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// Button renders a clickable button with minimalist styling.
// Default: cyan text color (noticeable but not bright)
// Focused: bold + brighter cyan + subtle gray background
func Button(props InteractiveProps, label string) component.Component {
	// Apply default interactive styling if no custom style provided
	effectiveStyle := props.StyleProps
	if effectiveStyle.Style == nil {
		// Default minimalist button style: cyan text, no background
		defaultStyle := lipgloss.NewStyle().
			Foreground(style.Interactive).
			PaddingLeft(1).
			PaddingRight(1)
		effectiveStyle.Style = &defaultStyle
	}

	// Create base text component
	textComp := Text(TextProps{
		StyleProps:  effectiveStyle,
		LayoutProps: props.LayoutProps,
		Content:     label,
	})

	// If no interaction, return simple text
	if props.OnClick == nil && len(props.Hotkeys) == 0 && len(props.Actions) == 0 {
		return textComp
	}

	// Convert primitive.HotkeyAction to component.HotkeyAction
	var actions []component.HotkeyAction
	for _, a := range props.Actions {
		actions = append(actions, component.HotkeyAction{
			Key:    a.Key,
			Action: a.Action,
		})
	}

	// Wrap with focusable
	return component.MakeFocusable(textComp, component.FocusableConfig{
		ID:        props.ID,
		Position:  props.Position,
		Hotkeys:   props.Hotkeys,
		Actions:   actions,
		Disabled:  props.Disabled,
		AutoFocus: props.AutoFocus,
		IsInput:   false,
		OnFocus:   props.OnFocus,
		OnBlur:    props.OnBlur,
		OnActivate: func() bool {
			if props.OnClick != nil {
				props.OnClick()
				return true
			}
			return false
		},
		FocusedStyle: func(content string) string {
			if props.FocusStyle != nil {
				return props.FocusStyle.Render(content)
			}
			// Minimalist focus style: bold + brighter text + subtle background
			focusStyle := lipgloss.NewStyle().
				Foreground(style.InteractiveFocus).
				Background(style.FocusBg).
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1)
			return focusStyle.Render(content)
		},
	})
}
