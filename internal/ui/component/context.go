package component

import (
	"fmt"
	"sync"
)

// Context holds state and metadata during component rendering
// Used for hooks (useState, useEffect, etc.)
type Context struct {
	// Component ID for hook isolation
	componentID string

	// HookState index counter (for hook order tracking)
	hookIndex int

	// Registered hooks for this render cycle
	hooks map[string]*hookState

	// Parent context (for nested components)
	parent *Context

	// Metadata
	meta map[string]any

	mu sync.RWMutex
}

// hookState stores hook-specific state
type hookState struct {
	value        any
	dependencies []any
	initialized  bool
}

// NewContext creates a new render context
func NewContext(componentID string) *Context {
	return &Context{
		componentID: componentID,
		hookIndex:   0,
		hooks:       make(map[string]*hookState),
		meta:        make(map[string]any),
	}
}

// Child creates a child context for nested component
func (ctx *Context) Child(componentID string) *Context {
	return &Context{
		componentID: componentID,
		hookIndex:   0,
		hooks:       make(map[string]*hookState),
		parent:      ctx,
		meta:        make(map[string]any),
	}
}

// nextHookIndex returns next hook index and increments counter
// Used for hook order tracking (React-style)
func (ctx *Context) nextHookIndex() int {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	index := ctx.hookIndex
	ctx.hookIndex++
	return index
}

// ResetHookIndex resets hook counter (called before each render)
func (ctx *Context) ResetHookIndex() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.hookIndex = 0
}

// GetHook retrieves hook state by key
func (ctx *Context) getHook(key string) (*hookState, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	state, ok := ctx.hooks[key]
	return state, ok
}

// SetHook stores hook state by key
func (ctx *Context) setHook(key string, state *hookState) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.hooks[key] = state
}

// GetMeta retrieves metadata value
func (ctx *Context) GetMeta(key string) (any, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	val, ok := ctx.meta[key]
	return val, ok
}

// SetMeta stores metadata value
func (ctx *Context) SetMeta(key string, value any) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.meta[key] = value
}

// ComponentID returns current component ID
func (ctx *Context) ComponentID() string {
	return ctx.componentID
}

// generateHookKey creates unique key for hook
// Uses component ID + hook identifier (key or index)
func (ctx *Context) generateHookKey(identifier ...string) string {
	if identifier == nil || len(identifier) == 0 {
		return fmt.Sprintf("hook_%s__%d", ctx.componentID, ctx.nextHookIndex())
	}
	return fmt.Sprintf("hook__%s__%s", ctx.componentID, identifier[0])
}

// Helper: check if dependencies changed
func (ctx *Context) dependenciesChanged(key string, newDeps []any) bool {
	state, exists := ctx.getHook(key)
	if !exists {
		return true
	}

	oldDeps := state.dependencies

	// If no old deps, changed
	if oldDeps == nil {
		return true
	}

	// If length different, changed
	if len(oldDeps) != len(newDeps) {
		return true
	}

	// Compare each dependency
	for i := range oldDeps {
		if !deepEqual(oldDeps[i], newDeps[i]) {
			return true
		}
	}

	return false
}

// deepEqual compares two values (simplified)
// For proper comparison, use reflect.DeepEqual or custom logic
func deepEqual(a, b any) bool {
	// Simple comparison for basic types
	// For complex types, implement proper equality check
	return a == b
}
