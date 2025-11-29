package primitive

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// InputProps configures the Input component.
type InputProps struct {
	BaseProps
	FocusProps

	// Current value
	Value string

	// Placeholder text when empty
	Placeholder string

	// Max length (0 = unlimited)
	MaxLength int

	// Input type: "text", "password", "number", "secret"
	Type string

	// Validation
	Validator func(string) error
	ErrorText string

	// Callbacks
	OnChange func(string)
	OnSubmit func(string)

	// Visual
	Prefix string // e.g., "$ " for money input
	Suffix string // e.g., " KB" for size input
}

// Input creates a single-line text input component with minimalist styling.
// Default: rounded border with gray border color, cyan text
// Focused: cyan border + subtle gray background
func Input(props InputProps) component.Component {
	width := props.Width
	if width == 0 {
		width = 20
	}

	inputType := props.Type
	if inputType == "" {
		inputType = "text"
	}

	baseComp := component.Func(func(ctx *component.Context) string {
		value := props.Value
		placeholder := props.Placeholder

		displayValue := value
		if inputType == "secret" || inputType == "password" {
			displayValue = strings.Repeat("*", len(value))
		}

		if props.Prefix != "" {
			displayValue = props.Prefix + displayValue
		}
		if props.Suffix != "" {
			displayValue = displayValue + props.Suffix
		}

		content := displayValue
		if content == "" && placeholder != "" {
			content = placeholder
		}

		// Truncate if too long
		if lipgloss.Width(content) > width {
			content = content[:width-1] + "..."
		}

		// Pad to fill width
		if lipgloss.Width(content) < width {
			content = content + strings.Repeat(" ", width-lipgloss.Width(content))
		}

		// Minimalist input style: rounded border, cyan text
		inputStyle := lipgloss.NewStyle().
			Width(width).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(style.Grey600).
			Foreground(style.Interactive)

		if props.Disabled {
			inputStyle = inputStyle.
				Foreground(style.Grey500).
				BorderForeground(style.Grey700)
		} else if value == "" && placeholder != "" {
			inputStyle = inputStyle.Foreground(style.Grey500)
		}

		if props.ErrorText != "" {
			inputStyle = inputStyle.BorderForeground(style.Error)
		}

		inputStyle = ApplyLayoutProps(inputStyle, props.LayoutProps)
		inputStyle = ApplyStyleProps(inputStyle, props.StyleProps)

		result := inputStyle.Render(content)

		// Add error message
		if props.ErrorText != "" {
			errorStyle := lipgloss.NewStyle().
				Foreground(style.Error).
				Italic(true)
			result += "\n" + errorStyle.Render("! "+props.ErrorText)
		}

		return result
	})

	// Return simple component if disabled
	if props.Disabled {
		return baseComp
	}

	// Wrap with focusable
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
			if props.OnSubmit != nil {
				props.OnSubmit(props.Value)
				return true
			}
			if props.OnClick != nil {
				props.OnClick()
				return true
			}
			return false
		},
		OnKeyPress: func(key string) bool {
			if props.OnChange == nil {
				return false
			}

			currentValue := props.Value

			// Handle backspace
			if key == "backspace" {
				if len(currentValue) > 0 {
					// Remove last rune (handle unicode properly)
					runes := []rune(currentValue)
					props.OnChange(string(runes[:len(runes)-1]))
					return true
				}
				return false
			}

			// Handle delete (clear all)
			if key == "ctrl+u" {
				props.OnChange("")
				return true
			}

			// Ignore control keys
			if strings.HasPrefix(key, "ctrl+") || strings.HasPrefix(key, "alt+") ||
				key == "esc" || key == "enter" || key == "tab" || key == "shift+tab" ||
				key == "up" || key == "down" || key == "left" || key == "right" {
				return false
			}

			// Check max length
			if props.MaxLength > 0 && len(currentValue) >= props.MaxLength {
				return false
			}

			// Handle space
			if key == "space" {
				key = " "
			}

			// Only accept printable characters (single rune)
			if len(key) == 1 {
				props.OnChange(currentValue + key)
				return true
			}

			return false
		},
		FocusedStyle: func(content string) string {
			if props.FocusStyle != nil {
				return props.FocusStyle.Render(content)
			}
			// Minimalist focus: cyan border + subtle background
			focusStyle := lipgloss.NewStyle().
				Width(width).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(style.FocusBorder).
				Background(style.FocusBg).
				Foreground(style.InteractiveFocus).
				Bold(true)
			// Re-render with focus style (extract text from content)
			return focusStyle.Render(strings.TrimSpace(content))
		},
	})
}
