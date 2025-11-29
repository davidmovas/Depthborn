package primitive

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type SelectOption struct {
	Label string
	Value string
}

type SelectProps struct {
	BaseProps
	FocusProps

	// Options
	Options []SelectOption

	// Selected value
	Value string

	// Placeholder when nothing selected
	Placeholder string

	// Callback
	OnChange func(string)

	// Visual
	ShowArrow bool
}

// Select creates a dropdown selection component
func Select(props SelectProps) component.Component {
	width := props.Width
	if width <= 0 {
		width = 20
	}

	showArrow := props.ShowArrow
	if !props.Disabled {
		showArrow = true
	}

	baseComp := component.Func(func(ctx *component.Context) string {
		selectedLabel := props.Placeholder
		if selectedLabel == "" {
			selectedLabel = "Select..."
		}

		for _, opt := range props.Options {
			if opt.Value == props.Value {
				selectedLabel = opt.Label
				break
			}
		}

		content := selectedLabel
		if showArrow {
			content = content + " ▼"
		}

		if lipgloss.Width(content) > width {
			content = content[:width-1] + "…"
		}

		if lipgloss.Width(content) < width {
			content = content + strings.Repeat(" ", width-lipgloss.Width(content))
		}

		selectStyle := lipgloss.NewStyle().
			Width(width).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(style.Grey400)

		if props.Disabled {
			selectStyle = selectStyle.
				Foreground(style.Grey400).
				Background(style.Grey100)
		} else if props.Value == "" {
			selectStyle = selectStyle.Foreground(style.Grey500)
		}

		selectStyle = ApplyLayoutProps(selectStyle, props.LayoutProps)
		selectStyle = ApplyStyleProps(selectStyle, props.StyleProps)

		return selectStyle.Render(content)
	})

	if props.Disabled {
		return baseComp
	}

	return component.MakeFocusable(baseComp, component.FocusableConfig{
		ID:        props.ID,
		Position:  props.Position,
		Hotkeys:   props.Hotkeys,
		Disabled:  props.Disabled,
		AutoFocus: props.AutoFocus,
		IsInput:   true,
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
			focusedStyle := lipgloss.NewStyle().
				BorderForeground(style.Primary).
				BorderStyle(lipgloss.RoundedBorder())
			return focusedStyle.Render(content)
		},
	})
}
