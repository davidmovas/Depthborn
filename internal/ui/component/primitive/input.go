package primitive

import (
	"fmt"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

type InputConfig struct {
	Placeholder string
	Value       string
	Focused     bool
	MaxLength   int
	Label       string
}

// Input creates text input component
// Note: Actual input handling is done by the screen/renderer
// This just renders the input field visually
func Input(config InputConfig) component.Component {
	return &input{config: config}
}

type input struct {
	config InputConfig
}

func (i *input) Render(ctx *component.Context) string {
	value := i.config.Value
	placeholder := i.config.Placeholder

	// Show placeholder if empty
	display := value
	if display == "" && placeholder != "" {
		display = placeholder
		// Gray out placeholder
		display = TextStyled(display, TextStyle{Color: "gray"}).Render(ctx)
	}

	// Add cursor if focused
	if i.config.Focused {
		display += "█" // block cursor
	}

	// Add label if provided
	if i.config.Label != "" {
		label := i.config.Label + ": "
		return label + display
	}

	return display
}

type TextAreaConfig struct {
	Value   string
	Focused bool
	Lines   int // number of visible lines
	Label   string
}

// TextArea creates multi-line text input
func TextArea(config TextAreaConfig) component.Component {
	return &textArea{config: config}
}

type textArea struct {
	config TextAreaConfig
}

func (ta *textArea) Render(ctx *component.Context) string {
	// For simplicity, just show the text
	// In real implementation, would handle scrolling, line wrapping, etc.

	value := ta.config.Value
	if ta.config.Focused {
		value += "█"
	}

	if ta.config.Label != "" {
		return fmt.Sprintf("%s:\n%s", ta.config.Label, value)
	}

	return value
}

type SelectConfig struct {
	Options  []string
	Selected int
	Label    string
	Focused  bool
}

// Select creates dropdown/select component
func Select(config SelectConfig) component.Component {
	return &selectComponent{config: config}
}

type selectComponent struct {
	config SelectConfig
}

func (s *selectComponent) Render(ctx *component.Context) string {
	if len(s.config.Options) == 0 {
		return "[No options]"
	}

	selected := s.config.Selected
	if selected < 0 || selected >= len(s.config.Options) {
		selected = 0
	}

	option := s.config.Options[selected]

	// Format: Label: [Option ▼]
	display := fmt.Sprintf("[%s ▼]", option)

	if s.config.Focused {
		display = TextStyled(display, TextStyle{Color: "cyan"}).Render(ctx)
	}

	if s.config.Label != "" {
		return s.config.Label + ": " + display
	}

	return display
}

type CheckboxConfig struct {
	Label   string
	Checked bool
	Focused bool
}

// Checkbox creates checkbox component
func Checkbox(config CheckboxConfig) component.Component {
	return &checkbox{config: config}
}

type checkbox struct {
	config CheckboxConfig
}

func (c *checkbox) Render(ctx *component.Context) string {
	box := "[ ]"
	if c.config.Checked {
		box = "[✓]"
	}

	if c.config.Focused {
		box = TextStyled(box, TextStyle{Color: "cyan"}).Render(ctx)
	}

	return box + " " + c.config.Label
}
