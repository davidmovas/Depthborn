package component

import (
	"fmt"
	"reflect"
	"sync"
)

// Context provides component rendering context
type Context struct {
	componentID   string
	parent        *Context
	navigator     any // Navigator reference
	focusContext  *FocusContext
	portalContext *PortalContext

	// Hook state
	hooks     map[string]*HookState
	hookIndex int
	hookMu    sync.RWMutex

	// Meta data
	meta   map[string]any
	metaMu sync.RWMutex

	// Dirty tracking for optimization
	dirty   bool
	dirtyMu sync.RWMutex

	// State change callback
	onStateChange func(componentID string)
}

// HookState stores hook-specific state
type HookState struct {
	Value        any
	Dependencies []any
	Initialized  bool
}

// NewContext creates new context
func NewContext(id string, navigator any) *Context {
	return &Context{
		componentID:   id,
		navigator:     navigator,
		focusContext:  NewFocusContext(generateID()),
		portalContext: NewPortalContext(),
		hooks:         make(map[string]*HookState),
		meta:          make(map[string]any),
		dirty:         true, // Initially dirty
	}
}

// Child creates child context
func (ctx *Context) Child(childID string) *Context {
	child := &Context{
		componentID:   childID,
		parent:        ctx,
		navigator:     ctx.navigator,
		focusContext:  ctx.focusContext,  // Shared focus context
		portalContext: ctx.portalContext, // Shared portal context
		hooks:         make(map[string]*HookState),
		meta:          make(map[string]any),
		dirty:         true,
		onStateChange: ctx.onStateChange,
	}
	return child
}

// SetHook stores hook state
func (ctx *Context) SetHook(key string, state *HookState) {
	ctx.hookMu.Lock()
	defer ctx.hookMu.Unlock()
	ctx.hooks[key] = state
}

// GetHook retrieves hook state
func (ctx *Context) GetHook(key string) (*HookState, bool) {
	ctx.hookMu.RLock()
	defer ctx.hookMu.RUnlock()
	state, exists := ctx.hooks[key]
	return state, exists
}

// GetFocusContext returns focus context (accessor method)
func (ctx *Context) GetFocusContext() *FocusContext {
	return ctx.focusContext
}

// SetFocusContext sets focus context (for portal isolation)
func (ctx *Context) SetFocusContext(focusContext *FocusContext) {
	ctx.focusContext = focusContext
}

// generateHookKey creates unique hook key
func (ctx *Context) generateHookKey(key ...string) string {
	ctx.hookMu.Lock()
	defer ctx.hookMu.Unlock()

	var hookKey string
	if len(key) > 0 {
		hookKey = fmt.Sprintf("%s_hook_%s", ctx.componentID, key[0])
	} else {
		hookKey = fmt.Sprintf("%s_hook_%d", ctx.componentID, ctx.hookIndex)
		ctx.hookIndex++
	}
	return hookKey
}

// dependenciesChanged checks if dependencies have changed (DEEP comparison)
func (ctx *Context) dependenciesChanged(hookKey string, deps []any) bool {
	state, exists := ctx.GetHook(hookKey)
	if !exists {
		return true
	}

	oldDeps := state.Dependencies
	if len(oldDeps) != len(deps) {
		return true
	}

	// Deep compare values, not pointers
	for i := range deps {
		if !deepEqual(oldDeps[i], deps[i]) {
			return true
		}
	}

	return false
}

// deepEqual compares values deeply (not pointer addresses)
func deepEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Use reflect.DeepEqual for proper value comparison
	return reflect.DeepEqual(a, b)
}

// TriggerStateChange marks component as dirty and notifies system
func (ctx *Context) TriggerStateChange() {
	ctx.MarkDirty()

	if ctx.onStateChange != nil {
		ctx.onStateChange(ctx.componentID)
	}
}

// MarkDirty marks this component as needing re-render
func (ctx *Context) MarkDirty() {
	ctx.dirtyMu.Lock()
	ctx.dirty = true
	ctx.dirtyMu.Unlock()

	// Mark parent dirty too (bubble up)
	if ctx.parent != nil {
		ctx.parent.MarkDirty()
	}
}

// IsDirty returns whether component needs re-render
func (ctx *Context) IsDirty() bool {
	ctx.dirtyMu.RLock()
	defer ctx.dirtyMu.RUnlock()
	return ctx.dirty
}

// ClearDirty marks component as clean
func (ctx *Context) ClearDirty() {
	ctx.dirtyMu.Lock()
	ctx.dirty = false
	ctx.dirtyMu.Unlock()
}

// SetStateChangeCallback sets callback for state changes
func (ctx *Context) SetStateChangeCallback(callback func(componentID string)) {
	ctx.onStateChange = callback
}

// FocusContext returns focus context
func (ctx *Context) FocusContext() *FocusContext {
	return ctx.focusContext
}

// PortalContext returns portal context
func (ctx *Context) PortalContext() *PortalContext {
	return ctx.portalContext
}

// Navigator returns navigator reference
func (ctx *Context) Navigator() any {
	return ctx.navigator
}

// ComponentID returns current component ID
func (ctx *Context) ComponentID() string {
	return ctx.componentID
}

// Get retrieves metadata value
func (ctx *Context) Get(key string) (any, bool) {
	ctx.metaMu.RLock()
	defer ctx.metaMu.RUnlock()
	val, ok := ctx.meta[key]
	return val, ok
}

// Set stores metadata value
func (ctx *Context) Set(key string, value any) {
	ctx.metaMu.Lock()
	defer ctx.metaMu.Unlock()
	ctx.meta[key] = value
}

// ResetHookIndex resets hook index (call before render)
func (ctx *Context) ResetHookIndex() {
	ctx.hookMu.Lock()
	ctx.hookIndex = 0
	ctx.hookMu.Unlock()
}
