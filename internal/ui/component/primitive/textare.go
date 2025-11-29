package primitive

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type TextAreaProps struct {
	BaseProps
	FocusProps

	// Current value
	Value string

	// Placeholder text when empty
	Placeholder string

	// Max length (0 = unlimited)
	MaxLength int

	// Validation
	Validator func(string) error
	ErrorText string

	// Callbacks
	OnChange func(string)
	OnSubmit func(string)

	// Visual
	ShowLineNumbers bool

	// Disabled state
	Disabled bool
}

// TextArea creates a multi-line text input component
func TextArea(props TextAreaProps) component.Component {
	width := 40
	height := 5
	if props.Width > 0 {
		width = props.Width
	}
	if props.Height > 0 {
		height = props.Height
	}

	baseComp := component.Func(func(ctx *component.Context) string {
		value := props.Value
		placeholder := props.Placeholder

		lines := strings.Split(value, "\n")

		if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
			if placeholder != "" {
				lines = []string{placeholder}
			} else {
				lines = []string{""}
			}
		}

		if len(lines) > height {
			lines = lines[:height]
		} else {
			for len(lines) < height {
				lines = append(lines, "")
			}
		}

		wrappedLines := make([]string, 0, len(lines))
		for _, line := range lines {
			if lipgloss.Width(line) > width-2 {
				line = line[:width-3] + "…"
			}

			if lipgloss.Width(line) < width-2 {
				line = line + strings.Repeat(" ", width-2-lipgloss.Width(line))
			}

			wrappedLines = append(wrappedLines, line)
		}

		content := strings.Join(wrappedLines, "\n")

		if props.ShowLineNumbers {
			numberedLines := make([]string, len(wrappedLines))
			for i, line := range wrappedLines {
				numStyle := lipgloss.NewStyle().
					Foreground(style.Grey500).
					Width(3).
					Align(lipgloss.Right)
				numberedLines[i] = numStyle.Render(fmt.Sprintf("%d", i+1)) + " │ " + line
			}
			content = strings.Join(numberedLines, "\n")
		}

		textareaStyle := lipgloss.NewStyle().
			Width(width).
			Height(height).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(style.Grey400).
			Padding(1)

		if props.Disabled {
			textareaStyle = textareaStyle.
				Foreground(style.Grey400).
				Background(style.Grey100)
		} else if value == "" && placeholder != "" {
			textareaStyle = textareaStyle.Foreground(style.Grey500)
		}

		if props.ErrorText != "" {
			textareaStyle = textareaStyle.BorderForeground(style.Error)
		}

		textareaStyle = ApplyLayoutProps(textareaStyle, props.LayoutProps)
		textareaStyle = ApplyStyleProps(textareaStyle, props.StyleProps)

		result := textareaStyle.Render(content)

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
			focusedStyle := lipgloss.NewStyle().
				BorderForeground(style.Primary).
				BorderStyle(lipgloss.RoundedBorder())

			lines := strings.Split(content, "\n")

			textareaLines := lines
			if props.ErrorText != "" && len(lines) > 0 {
				textareaLines = lines[:len(lines)-1]
			}

			textareaContent := strings.Join(textareaLines, "\n")
			focusedTextarea := focusedStyle.Render(textareaContent)

			if props.ErrorText != "" && len(lines) > 0 {
				focusedTextarea += "\n" + lines[len(lines)-1]
			}

			return focusedTextarea
		}

		if props.FocusStyle != nil {
			focusStyle = func(content string) string {
				return props.FocusStyle.Render(content)
			}
		}

		return component.MakeFocusable(baseComp, component.FocusableConfig{
			Position:  props.Position,
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
			FocusedStyle: focusStyle,
		})
	}

	return baseComp
}
