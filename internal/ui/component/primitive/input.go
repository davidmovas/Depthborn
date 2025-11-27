package primitive

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

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

	// Disabled state
	Disabled bool
}

// Input creates a single-line text input component
func Input(props InputProps) component.Component {
	width := 20
	if props.Width != nil {
		width = *props.Width
	}

	inputType := props.Type
	if inputType == "" {
		inputType = "text"
	}

	baseComp := component.Func(func(ctx *component.Context) string {
		value := props.Value
		placeholder := props.Placeholder

		displayValue := value
		if inputType == "secret" && len(value) > 0 {
			displayValue = strings.Repeat("•", len(value))
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

		if lipgloss.Width(content) > width {
			content = content[:width-1] + "…"
		}

		if lipgloss.Width(content) < width {
			content = content + strings.Repeat(" ", width-lipgloss.Width(content))
		}

		inputStyle := lipgloss.NewStyle().
			Width(width).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(style.Grey400)

		if props.Disabled {
			inputStyle = inputStyle.
				Foreground(style.Grey400).
				Background(style.Grey100)
		} else if value == "" && placeholder != "" {
			inputStyle = inputStyle.Foreground(style.Grey500)
		}

		if props.ErrorText != "" {
			inputStyle = inputStyle.BorderForeground(style.Error)
		}

		inputStyle = ApplyLayoutProps(inputStyle, props.LayoutProps)
		inputStyle = ApplyStyleProps(inputStyle, props.StyleProps)

		result := inputStyle.Render(content)

		if props.ErrorText != "" {
			errorStyle := lipgloss.NewStyle().
				Foreground(style.Error).
				Italic(true)
			result += "\n" + errorStyle.Render("⚠ "+props.ErrorText)
		}

		return result
	})

	if !props.Disabled {
		focusStyle := func(content string) string {
			lines := strings.Split(content, "\n")

			// Style input line with focus
			focusedStyle := lipgloss.NewStyle().
				BorderForeground(style.Primary).
				BorderStyle(lipgloss.RoundedBorder())

			lines[0] = focusedStyle.Render(strings.TrimSpace(lines[0]))

			return strings.Join(lines, "\n")
		}

		if props.FocusStyle != nil {
			focusStyle = func(content string) string {
				return props.FocusStyle.Render(content)
			}
		}

		return component.MakeFocusable(baseComp, component.FocusableConfig{
			Position:        props.Position,
			CanFocus:        true,
			AutoFocus:       props.AutoFocus != nil && *props.AutoFocus,
			IsInput:         true,
			OnFocusCallback: props.OnFocus,
			OnBlurCallback:  props.OnBlur,
			OnActivateCallback: func() bool {
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
			FocusedStyle: focusStyle,
		})
	}

	return baseComp
}
