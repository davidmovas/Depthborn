package component

// FocusableComponent wraps any component to make it focusable
type FocusableComponent struct {
	component Component
	config    FocusableConfig
	id        string
}

// FocusableConfig configures focusable behavior
type FocusableConfig struct {
	// Unique ID for this focusable (REQUIRED - must be stable across renders)
	ID string

	// 2D position for spatial navigation (nil = linear only)
	Position *FocusPosition

	// Hotkeys for activation (can be multiple, e.g. ["n", "ctrl+n"])
	Hotkeys []string

	// Whether can receive focus
	CanFocus bool

	// Whether this should be auto-focused on mount
	AutoFocus bool

	// Whether this is an input component (blocks hotkeys when focused)
	IsInput bool

	// Callbacks
	OnFocusCallback    func()
	OnBlurCallback     func()
	OnActivateCallback func() bool

	// Style modifiers (applied when focused)
	FocusedStyle func(content string) string
}

// MakeFocusable wraps component to make it focusable
func MakeFocusable(comp Component, config FocusableConfig) *FocusableComponent {
	if config.ID == "" {
		panic("FocusableConfig.ID is required and must be stable across renders")
	}

	return &FocusableComponent{
		component: comp,
		config:    config,
		id:        config.ID,
	}
}

// Render implements Component interface
func (fc *FocusableComponent) Render(ctx *Context) string {
	// Register with focus context
	ctx.FocusContext().Register(fc)

	// Check if currently focused BY INDEX
	// Don't use Current() because it might be called before all components registered
	isFocused := ctx.FocusContext().IsFocusedByID(fc.id)

	// Render wrapped component
	content := fc.component.Render(ctx)

	// Apply focused style if focused
	if isFocused && fc.config.FocusedStyle != nil {
		content = fc.config.FocusedStyle(content)
	}

	return content
}

// Focusable interface implementation

func (fc *FocusableComponent) GetFocusID() string {
	return fc.id
}

func (fc *FocusableComponent) GetFocusPosition() *FocusPosition {
	return fc.config.Position
}

func (fc *FocusableComponent) GetHotkeys() []string {
	return fc.config.Hotkeys
}

func (fc *FocusableComponent) CanReceiveFocus() bool {
	return fc.config.CanFocus
}

func (fc *FocusableComponent) ShouldAutoFocus() bool {
	return fc.config.AutoFocus
}

func (fc *FocusableComponent) OnFocus() {
	if fc.config.OnFocusCallback != nil {
		fc.config.OnFocusCallback()
	}
}

func (fc *FocusableComponent) OnBlur() {
	if fc.config.OnBlurCallback != nil {
		fc.config.OnBlurCallback()
	}
}

func (fc *FocusableComponent) OnActivate() bool {
	if fc.config.OnActivateCallback != nil {
		return fc.config.OnActivateCallback()
	}
	return false
}

func (fc *FocusableComponent) IsInputComponent() bool {
	return fc.config.IsInput
}
