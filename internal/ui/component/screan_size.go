package component

// ScreenSize represents terminal/window dimensions
type ScreenSize struct {
	Width  int
	Height int
}

// NewScreenSize creates screen size with validation
func NewScreenSize(width, height int) ScreenSize {
	if width < 1 {
		width = 80
	}
	if height < 1 {
		height = 24
	}
	return ScreenSize{Width: width, Height: height}
}

// AspectRatio returns width/height ratio
func (s ScreenSize) AspectRatio() float64 {
	if s.Height == 0 {
		return 0
	}
	return float64(s.Width) / float64(s.Height)
}

// IsNarrow returns true if width < 60
func (s ScreenSize) IsNarrow() bool {
	return s.Width < 60
}

// IsWide returns true if width >= 100
func (s ScreenSize) IsWide() bool {
	return s.Width >= 100
}

// IsCompact returns true if height < 20
func (s ScreenSize) IsCompact() bool {
	return s.Height < 20
}

// IsTall returns true if height >= 40
func (s ScreenSize) IsTall() bool {
	return s.Height >= 40
}

// ScreenSize returns current screen dimensions
func (ctx *Context) ScreenSize() ScreenSize {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if size, ok := ctx.meta["__screen_size"].(ScreenSize); ok {
		return size
	}

	// Default fallback
	return NewScreenSize(80, 24)
}

// SetScreenSize updates screen dimensions
// Should be called by renderer when window resizes
func (ctx *Context) SetScreenSize(width, height int) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.meta["__screen_size"] = NewScreenSize(width, height)
}

// WithScreenSize creates child context with specific screen size
// Useful for nested viewports or virtual windows
func (ctx *Context) WithScreenSize(width, height int) *Context {
	child := ctx.Child(ctx.componentID + "_sized")
	child.SetScreenSize(width, height)
	return child
}
