package component

import (
	"crypto/rand"
	"encoding/hex"
)

// FocusableConfig configures focusable behavior.
type FocusableConfig struct {
	// ID is an explicit ID (auto-generated if empty)
	ID string

	// Position for 2D grid navigation (nil = linear only)
	Position *FocusPosition

	// Hotkeys for direct activation (e.g., []string{"n", "ctrl+n"})
	Hotkeys []string

	// Whether this component can receive focus
	Disabled bool

	// Whether this should be auto-focused when registered
	AutoFocus bool

	// Whether this is an input component (blocks hotkeys when focused)
	IsInput bool

	// Callbacks
	OnFocus    func()
	OnBlur     func()
	OnActivate func() bool
	OnKeyPress func(key string) bool // For text input handling

	// Style applied when focused
	FocusedStyle func(content string) string
}

// FocusableComponent wraps any component to make it focusable.
type FocusableComponent struct {
	component Component
	config    FocusableConfig
	id        string
}

// MakeFocusable wraps a component to make it focusable.
func MakeFocusable(comp Component, config FocusableConfig) Component {
	id := config.ID
	if id == "" {
		id = generateID()
	}

	return &FocusableComponent{
		component: comp,
		config:    config,
		id:        id,
	}
}

// Render implements Component interface.
func (fc *FocusableComponent) Render(ctx *Context) string {
	// Register with focus manager
	isFocused := ctx.Focus().Register(fc)

	// Render wrapped component
	content := fc.component.Render(ctx)

	// Apply focused style if focused
	if isFocused && fc.config.FocusedStyle != nil {
		content = fc.config.FocusedStyle(content)
	}

	return content
}

// --- Focusable Interface Implementation ---

func (fc *FocusableComponent) ID() string {
	return fc.id
}

func (fc *FocusableComponent) Position() *FocusPosition {
	return fc.config.Position
}

func (fc *FocusableComponent) Hotkeys() []string {
	return fc.config.Hotkeys
}

func (fc *FocusableComponent) CanFocus() bool {
	return !fc.config.Disabled
}

func (fc *FocusableComponent) AutoFocus() bool {
	return fc.config.AutoFocus
}

func (fc *FocusableComponent) IsInput() bool {
	return fc.config.IsInput
}

func (fc *FocusableComponent) OnFocus() {
	if fc.config.OnFocus != nil {
		fc.config.OnFocus()
	}
}

func (fc *FocusableComponent) OnBlur() {
	if fc.config.OnBlur != nil {
		fc.config.OnBlur()
	}
}

func (fc *FocusableComponent) OnActivate() bool {
	if fc.config.OnActivate != nil {
		return fc.config.OnActivate()
	}
	return false
}

func (fc *FocusableComponent) OnKeyPress(key string) bool {
	if fc.config.OnKeyPress != nil {
		return fc.config.OnKeyPress(key)
	}
	return false
}

// --- Helper Functions ---

// generateID creates a random unique ID.
func generateID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to less random but still unique
		return "id_fallback"
	}
	return "fc_" + hex.EncodeToString(bytes)
}

// UseFocusable is a hook to create focusable state.
// Returns (isFocused, focus, blur) functions.
func UseFocusable(ctx *Context, id string) (bool, func(), func()) {
	isFocused := ctx.Focus().IsFocused(id)

	focus := func() {
		ctx.Focus().FocusID(id)
	}

	blur := func() {
		if ctx.Focus().IsFocused(id) {
			ctx.Focus().ClearFocus()
		}
	}

	return isFocused, focus, blur
}
