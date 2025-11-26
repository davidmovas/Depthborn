package primitive

import (
	"fmt"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// ButtonConfig configures button component
type ButtonConfig struct {
	// Button label
	Label string

	// Hotkeys for activation (can be multiple)
	Hotkeys []string

	// Single hotkey (shorthand for Hotkeys: []string{Key})
	Key string

	// 2D position for navigation (optional)
	Position *component.FocusPosition

	// Click handler
	OnClick func()

	// Custom ID (REQUIRED - must be stable, e.g. "btn_new_game")
	ID string

	// Whether this button should be auto-focused
	AutoFocus bool
}

// Button creates a focusable button component
func Button(config ButtonConfig) component.Component {
	if config.ID == "" {
		panic("Button requires stable ID (e.g. 'btn_new_game')")
	}

	// Merge Key into Hotkeys if provided
	hotkeys := config.Hotkeys
	if config.Key != "" {
		hotkeys = append(hotkeys, config.Key)
	}

	// Create display text
	hotkeyText := ""
	if len(hotkeys) > 0 {
		hotkeyText = fmt.Sprintf("[%s] ", hotkeys[0])
	}
	displayText := hotkeyText + config.Label

	// Create base text component
	baseComp := Text(displayText)

	// Wrap with focusable
	return component.MakeFocusable(baseComp, component.FocusableConfig{
		ID:        config.ID,
		Position:  config.Position,
		Hotkeys:   hotkeys,
		CanFocus:  true,
		AutoFocus: config.AutoFocus,
		IsInput:   false,

		OnActivateCallback: func() bool {
			if config.OnClick != nil {
				config.OnClick()
				return true
			}
			return false
		},

		FocusedStyle: func(content string) string {
			return "> " + content + " <"
		},
	})
}
