package component

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// HooksContainer can store and retrieve hooks (implemented by Screen)
type HooksContainer interface {
	SetHook(key string, state *HookState)
	GetHook(key string) (*HookState, bool)
}

// HookState stores hook-specific state (exported for Screen storage)
type HookState struct {
	Value        any
	Dependencies []any
	Initialized  bool
}

// Context holds state and metadata during component rendering
// Used for hooks (useState, useEffect, etc.) and focus management
// Context holds rendering state and metadata for component tree
type Context struct {
	ctx           context.Context
	componentID   string
	hookIndex     int
	hooks         map[string]*HookState
	parent        *Context
	meta          map[string]any
	focusCtx      *FocusContext
	onStateChange func()
	mu            sync.RWMutex
}

// NewContext creates a new render context
func NewContext(componentID string) *Context {
	return &Context{
		componentID:   componentID,
		hooks:         make(map[string]*HookState),
		meta:          make(map[string]any),
		ctx:           context.Background(),
		onStateChange: func() {},
	}
}

// Child creates a child context for nested components
func (ctx *Context) Child(componentID string) *Context {
	return &Context{
		parent:        ctx,
		componentID:   componentID,
		hooks:         make(map[string]*HookState),
		meta:          make(map[string]any),
		ctx:           ctx.ctx,
		onStateChange: ctx.onStateChange,
	}
}

// ResetHookIndex resets hook counter before each render
func (ctx *Context) ResetHookIndex() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.hookIndex = 0
}

// SetHook stores hook state with container support
func (ctx *Context) SetHook(key string, state *HookState) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if container, ok := ctx.meta["__hooks_container"].(HooksContainer); ok {
		container.SetHook(key, state)
		return
	}

	ctx.hooks[key] = state
}

// GetHook retrieves hook state with container support
func (ctx *Context) GetHook(key string) (*HookState, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if container, ok := ctx.meta["__hooks_container"].(HooksContainer); ok {
		return container.GetHook(key)
	}

	state, ok := ctx.hooks[key]
	return state, ok
}

// SetHooksContainer links context to external hooks storage
func (ctx *Context) SetHooksContainer(container HooksContainer) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.meta["__hooks_container"] = container
}

// TriggerStateChange requests component re-render
func (ctx *Context) TriggerStateChange() {
	ctx.mu.RLock()
	callback := ctx.onStateChange
	ctx.mu.RUnlock()

	if callback != nil {
		callback()
	}
}

// SetOnStateChange sets re-render callback
func (ctx *Context) SetOnStateChange(callback func()) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.onStateChange = callback
}

// FocusContext returns focus management context
func (ctx *Context) FocusContext() *FocusContext {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.focusCtx == nil {
		ctx.focusCtx = NewFocusContext(ctx.componentID)
	}
	return ctx.focusCtx
}

// SetFocusContext sets focus context from parent
func (ctx *Context) SetFocusContext(focusCtx *FocusContext) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.focusCtx = focusCtx
}

// generateHookKey creates unique key for hook
// Uses component ID + hook identifier (key or index)
func (ctx *Context) generateHookKey(identifier ...string) string {
	if len(identifier) == 0 {
		return fmt.Sprintf("hook_%s__%d", ctx.componentID, ctx.nextHookIndex())
	}
	return fmt.Sprintf("hook__%s__%s", ctx.componentID, identifier[0])
}

// dependenciesChanged checks if dependencies have changed
func (ctx *Context) dependenciesChanged(key string, newDeps []any) bool {
	state, exists := ctx.GetHook(key)
	if !exists {
		return true
	}

	oldDeps := state.Dependencies

	// If no old deps, consider changed
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

func (ctx *Context) nextHookIndex() int {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	index := ctx.hookIndex
	ctx.hookIndex++
	return index
}

func deepEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}
