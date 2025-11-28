package component

import (
	"fmt"
	"reflect"
)

// State represents a reactive state value with setter
type State[T any] struct {
	value T
	set   func(T)
}

// Value returns current state value
func (s State[T]) Value() T {
	return s.value
}

// Set updates state value and triggers re-render
func (s State[T]) Set(newValue T) {
	s.set(newValue)
}

// UseState creates or retrieves stateful value
func UseState[T any](ctx *Context, initial T, key ...string) *State[T] {
	hookKey := ctx.generateHookKey(key...)

	state, exists := ctx.GetHook(hookKey)
	if !exists || !state.Initialized {
		state = &HookState{
			Value:       initial,
			Initialized: true,
		}
		ctx.SetHook(hookKey, state)
	}

	value := state.Value.(T)

	setter := func(newValue T) {
		if !reflect.DeepEqual(state.Value, newValue) {
			state.Value = newValue
			ctx.SetHook(hookKey, state)
			ctx.TriggerStateChange()
		}
	}

	return &State[T]{
		value: value,
		set:   setter,
	}
}

// UseEffect runs side effect when dependencies change
// - deps = nil - re-render on every render
// - deps = [] - re-render on first render
// - deps = [...] - re-render when deps change
func UseEffect(ctx *Context, effect func(), deps []any, key ...string) {
	hookKey := ctx.generateHookKey(key...)

	state, exists := ctx.GetHook(hookKey)

	if deps == nil {
		effect()
		return
	}

	if len(deps) == 0 {
		if !exists || !state.Initialized {
			effect()
			ctx.SetHook(hookKey, &HookState{
				Dependencies: deps,
				Initialized:  true,
			})
		}
		return
	}

	if ctx.dependenciesChanged(hookKey, deps) {
		effect()
		ctx.SetHook(hookKey, &HookState{
			Dependencies: deps,
			Initialized:  true,
		})
	}
}

// UseMemo memoizes expensive computation
func UseMemo[T any](ctx *Context, compute func() T, deps []any, key ...string) T {
	hookKey := ctx.generateHookKey(key...)

	state, exists := ctx.GetHook(hookKey)
	shouldCompute := !exists || ctx.dependenciesChanged(hookKey, deps)

	if shouldCompute {
		value := compute()
		state = &HookState{
			Value:        value,
			Dependencies: deps,
			Initialized:  true,
		}
		ctx.SetHook(hookKey, state)
		return value
	}

	return state.Value.(T)
}

// Ref holds mutable reference that persists across renders
type Ref[T any] struct {
	Current T
}

// UseRef creates persistent mutable reference
func UseRef[T any](ctx *Context, initial T, key ...string) *Ref[T] {
	hookKey := ctx.generateHookKey(key...)

	state, exists := ctx.GetHook(hookKey)
	if !exists || !state.Initialized {
		ref := &Ref[T]{Current: initial}
		state = &HookState{
			Value:       ref,
			Initialized: true,
		}
		ctx.SetHook(hookKey, state)
	}

	return state.Value.(*Ref[T])
}

// UseCallback memoizes callback function
func UseCallback[T any](ctx *Context, callback T, deps []any, key ...string) T {
	callbackType := reflect.TypeOf(callback)
	if callbackType.Kind() != reflect.Func {
		panic(fmt.Sprintf("UseCallback: callback must be a function, got %T", callback))
	}

	return UseMemo(ctx, func() T {
		return callback
	}, deps, key...)
}
