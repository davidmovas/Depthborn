package component

import (
	"fmt"
	"reflect"
)

// State represents stateful value with setter
type State[T any] struct {
	value T
	set   func(T)
}

// Value returns current state value
func (s State[T]) Value() T {
	return s.value
}

// Set updates state value
func (s State[T]) Set(newValue T) {
	s.set(newValue)
}

// UseState creates or retrieves stateful value
// With explicit key (preferred for clarity and safety)
//
// Example:
//
//	health := UseState(ctx, 100,  "health")
//	health.Set(health.Value() + 10)
func UseState[T any](ctx *Context, initial T, key ...string) *State[T] {
	hookKey := ctx.generateHookKey(key...)

	state, exists := ctx.getHook(hookKey)
	if !exists || !state.initialized {
		// Initialize state
		state = &hookState{
			value:       initial,
			initialized: true,
		}
		ctx.setHook(hookKey, state)
	}

	value := state.value.(T)

	setter := func(newValue T) {
		state.value = newValue
		ctx.setHook(hookKey, state)
	}

	return &State[T]{
		value: value,
		set:   setter,
	}
}

// UseEffect runs side effect when dependencies change
// - If deps is nil/empty: runs on every render
// - If deps provided: runs only when deps change
// - On first render: always runs
//
// Example:
//
//	UseEffect(ctx, "log_health", func() {
//	    fmt.Println("Health:", health.Value())
//	}, health.Value())
func UseEffect(ctx *Context, effect func(), deps []any, key ...string) {
	hookKey := ctx.generateHookKey(key...)

	// Check if dependencies changed
	shouldRun := ctx.dependenciesChanged(hookKey, deps)

	if shouldRun {
		// Run effect
		effect()

		// Store dependencies
		state := &hookState{
			dependencies: deps,
			initialized:  true,
		}
		ctx.setHook(hookKey, state)
	}
}

// UseMemo memoize expensive computation
// Recomputes only when dependencies change
//
// Example:
//
//	expensive := UseMemo(ctx, "calc", func() int {
//	    return complexCalculation()
//	}, dep1, dep2)
func UseMemo[T any](ctx *Context, compute func() T, deps []any, key ...string) T {
	hookKey := ctx.generateHookKey(key...)

	state, exists := ctx.getHook(hookKey)

	// Check if we need to recompute
	shouldCompute := !exists || ctx.dependenciesChanged(hookKey, deps)

	if shouldCompute {
		// Compute new value
		value := compute()

		// Store value and dependencies
		state = &hookState{
			value:        value,
			dependencies: deps,
			initialized:  true,
		}
		ctx.setHook(hookKey, state)

		return value
	}

	// Return cached value
	return state.value.(T)
}

// Ref holds mutable reference that persists across renders
type Ref[T any] struct {
	Current T
}

// UseRef creates persistent mutable reference
// Unlike useState, changing ref doesn't trigger re-render
//
// Example:
//
//	counter := UseRef(ctx, "counter", 0)
//	counter.Current++ // doesn't trigger re-render
func UseRef[T any](ctx *Context, initial T, key ...string) *Ref[T] {
	hookKey := ctx.generateHookKey(key...)

	state, exists := ctx.getHook(hookKey)
	if !exists || !state.initialized {
		ref := &Ref[T]{Current: initial}
		state = &hookState{
			value:       ref,
			initialized: true,
		}
		ctx.setHook(hookKey, state)
	}

	return state.value.(*Ref[T])
}

// UseCallback memoizes callback function
// Returns same function instance when dependencies don't change
// Useful for preventing unnecessary re-renders of child components
//
// Example:
//
//	onClick := UseCallback(ctx, "on_click", func() {
//	    handleClick(id)
//	}, id)
func UseCallback[T any](ctx *Context, callback T, deps []any, key ...string) T {
	// Validate that callback is a function
	callbackType := reflect.TypeOf(callback)
	if callbackType.Kind() != reflect.Func {
		panic(fmt.Sprintf("UseCallback: callback must be a function, got %T", callback))
	}

	return UseMemo(ctx, func() T {
		return callback
	}, deps, key...)
}
