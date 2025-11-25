package primitive

import (
	"fmt"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

type ButtonConfig struct {
	Label    string
	Key      string // Key to press (e.g. "A", "Enter", "Esc")
	OnClick  func()
	Disabled bool
	Style    ButtonStyle
}

type ButtonStyle string

const (
	ButtonStyleDefault   ButtonStyle = "default"
	ButtonStylePrimary   ButtonStyle = "primary"
	ButtonStyleDanger    ButtonStyle = "danger"
	ButtonStyleSecondary ButtonStyle = "secondary"
)

// Button creates button component
// In terminal UI, buttons are just text labels with key hints
// The actual key handling is done in Screen.HandleInput()
//
// Example:
//
//	Button("Attack", "A", func() { /* handle click */ })
func Button(config ButtonConfig) component.Component {
	return &button{config: config}
}

type button struct {
	config ButtonConfig
}

func (b *button) Render(ctx *component.Context) string {
	label := b.config.Label
	key := b.config.Key

	// Format: [A] Attack
	if b.config.Disabled {
		// Gray out disabled buttons
		return fmt.Sprintf("[-] %s", label)
	}

	// Apply style coloring
	var color string
	switch b.config.Style {
	case ButtonStylePrimary:
		color = "cyan"
	case ButtonStyleDanger:
		color = "red"
	case ButtonStyleSecondary:
		color = "gray"
	default:
		color = "white"
	}

	// Render with key hint
	keyPart := fmt.Sprintf("[%s]", key)

	// Use styled text for coloring
	styledLabel := TextStyled(label, TextStyle{Color: color})
	styledKey := TextStyled(keyPart, TextStyle{Color: "yellow", Bold: true})

	return styledKey.Render(ctx) + " " + styledLabel.Render(ctx)
}

// SimpleButton creates button without key hint (for list selection, etc.)
func SimpleButton(label string, onClick func()) component.Component {
	return Button(ButtonConfig{
		Label:   label,
		OnClick: onClick,
	})
}

// PrimaryButton creates primary styled button
func PrimaryButton(label string, key string, onClick func()) component.Component {
	return Button(ButtonConfig{
		Label:   label,
		Key:     key,
		OnClick: onClick,
		Style:   ButtonStylePrimary,
	})
}

// DangerButton creates danger styled button (for destructive actions)
func DangerButton(label string, key string, onClick func()) component.Component {
	return Button(ButtonConfig{
		Label:   label,
		Key:     key,
		OnClick: onClick,
		Style:   ButtonStyleDanger,
	})
}
