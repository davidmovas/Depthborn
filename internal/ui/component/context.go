package component

import (
	"fmt"
	"sync"
)

// Context provides rendering context and state management for components.
// It maintains hook state, focus context, screen size, and hierarchical structure.
type Context struct {
	// Identity
	id       string
	path     string // Full path from root (for hook keys)
	parent   *Context
	children map[string]*Context

	// External references
	navigator any

	// Screen dimensions
	width  int
	height int

	// Hook management - single source of truth
	hooks     map[string]*HookState
	hookIndex int
	hookMu    sync.RWMutex

	// Focus management
	focus   *FocusManager
	focusMu sync.RWMutex

	// Portal management
	portals *PortalManager

	// Render state
	renderCount uint64
	needsRender bool

	// Callback for triggering re-render
	onRenderRequest func()

	mu sync.RWMutex
}

// HookState stores the state for a single hook instance.
type HookState struct {
	Value        any
	Dependencies []any
	Cleanup      func() // For UseEffect cleanup
	Initialized  bool
}

// NewContext creates a new root context.
func NewContext(id string, navigator any) *Context {
	ctx := &Context{
		id:          id,
		path:        id,
		navigator:   navigator,
		hooks:       make(map[string]*HookState),
		children:    make(map[string]*Context),
		focus:       NewFocusManager(),
		portals:     NewPortalManager(),
		width:       80,
		height:      24,
		needsRender: true,
	}

	// Connect focus manager to trigger re-renders on focus change
	ctx.focus.SetOnFocusChange(func() {
		ctx.RequestRender()
	})

	// Connect portal manager to trigger re-renders on portal focus change
	ctx.portals.SetOnRenderRequest(func() {
		ctx.RequestRender()
	})

	return ctx
}

// WithKey creates a child context with a specific key.
// Used for list items to maintain stable hook state.
func (ctx *Context) WithKey(key string) *Context {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	childPath := fmt.Sprintf("%s.%s", ctx.path, key)

	if child, exists := ctx.children[key]; exists {
		child.hookIndex = 0 // Reset hook index for new render
		return child
	}

	child := &Context{
		id:              key,
		path:            childPath,
		parent:          ctx,
		navigator:       ctx.navigator,
		hooks:           make(map[string]*HookState),
		children:        make(map[string]*Context),
		focus:           ctx.focus,   // Share focus manager
		portals:         ctx.portals, // Share portal manager
		width:           ctx.width,
		height:          ctx.height,
		onRenderRequest: ctx.onRenderRequest,
	}

	ctx.children[key] = child
	return child
}

// WithIndex creates a child context with an index key.
func (ctx *Context) WithIndex(index int) *Context {
	return ctx.WithKey(fmt.Sprintf("[%d]", index))
}

// ID returns the context identifier.
func (ctx *Context) ID() string {
	return ctx.id
}

// Path returns the full path from root.
func (ctx *Context) Path() string {
	return ctx.path
}

// Navigator returns the navigator reference.
func (ctx *Context) Navigator() any {
	return ctx.navigator
}

// ScreenSize returns current screen dimensions.
func (ctx *Context) ScreenSize() (width, height int) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.width, ctx.height
}

// SetScreenSize updates screen dimensions.
func (ctx *Context) SetScreenSize(width, height int) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if width < 1 {
		width = 80
	}
	if height < 1 {
		height = 24
	}

	ctx.width = width
	ctx.height = height

	// Propagate to children
	for _, child := range ctx.children {
		child.SetScreenSize(width, height)
	}
}

// Focus returns the focus manager.
func (ctx *Context) Focus() *FocusManager {
	ctx.focusMu.RLock()
	defer ctx.focusMu.RUnlock()
	return ctx.focus
}

// NextFocusRow moves to the next row for focus navigation.
// Call this between groups of horizontal focusable elements.
// Example:
//
//	[Tab1] [Tab2] [Tab3]  <- row 0
//	ctx.NextFocusRow()
//	[Btn1] [Btn2]         <- row 1
func (ctx *Context) NextFocusRow() {
	ctx.Focus().NextRow()
}

// SetFocusManager sets a different focus manager (for modal isolation).
func (ctx *Context) SetFocusManager(fm *FocusManager) {
	ctx.focusMu.Lock()
	defer ctx.focusMu.Unlock()
	ctx.focus = fm
}

// Portals returns the portal manager.
func (ctx *Context) Portals() *PortalManager {
	return ctx.portals
}

// RequestRender signals that a re-render is needed.
func (ctx *Context) RequestRender() {
	ctx.mu.Lock()
	ctx.needsRender = true
	callback := ctx.onRenderRequest
	ctx.mu.Unlock()

	if callback != nil {
		callback()
	}
}

// SetRenderCallback sets the callback for render requests.
func (ctx *Context) SetRenderCallback(callback func()) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.onRenderRequest = callback
}

// BeginRender prepares context for a new render cycle.
func (ctx *Context) BeginRender() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.hookIndex = 0
	ctx.renderCount++
	ctx.needsRender = false

	// Clear focus registrations for re-registration
	ctx.focus.BeginFrame()
}

// EndRender completes the render cycle.
func (ctx *Context) EndRender() {
	ctx.focus.EndFrame()
}

// --- Hook State Management ---

// getHookKey generates a unique key for a hook.
func (ctx *Context) getHookKey(explicitKey string) string {
	ctx.hookMu.Lock()
	defer ctx.hookMu.Unlock()

	if explicitKey != "" {
		return fmt.Sprintf("%s#%s", ctx.path, explicitKey)
	}

	key := fmt.Sprintf("%s#%d", ctx.path, ctx.hookIndex)
	ctx.hookIndex++
	return key
}

// GetHook retrieves hook state by key.
func (ctx *Context) GetHook(key string) (*HookState, bool) {
	ctx.hookMu.RLock()
	defer ctx.hookMu.RUnlock()

	// Check own hooks first
	if state, exists := ctx.hooks[key]; exists {
		return state, true
	}

	// Check parent chain (for shared hooks)
	if ctx.parent != nil {
		return ctx.parent.GetHook(key)
	}

	return nil, false
}

// SetHook stores hook state.
func (ctx *Context) SetHook(key string, state *HookState) {
	ctx.hookMu.Lock()
	defer ctx.hookMu.Unlock()
	ctx.hooks[key] = state
}

// CleanupHooks runs cleanup functions for all hooks.
func (ctx *Context) CleanupHooks() {
	ctx.hookMu.Lock()
	defer ctx.hookMu.Unlock()

	for _, state := range ctx.hooks {
		if state.Cleanup != nil {
			state.Cleanup()
			state.Cleanup = nil
		}
	}

	// Cleanup children
	for _, child := range ctx.children {
		child.CleanupHooks()
	}
}

// --- Convenience Methods ---

// Width returns screen width.
func (ctx *Context) Width() int {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.width
}

// Height returns screen height.
func (ctx *Context) Height() int {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.height
}

// IsNarrow returns true if width < 60.
func (ctx *Context) IsNarrow() bool {
	return ctx.Width() < 60
}

// IsWide returns true if width >= 100.
func (ctx *Context) IsWide() bool {
	return ctx.Width() >= 100
}
