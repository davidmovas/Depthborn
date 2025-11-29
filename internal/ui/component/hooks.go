package component

import (
	"reflect"
)

// State provides reactive state management.
// Changes to state trigger re-renders.
type State[T any] struct {
	value   T
	setter  func(T)
	hookKey string
}

// Get returns the current value.
func (s *State[T]) Get() T {
	return s.value
}

// Value is an alias for Get (for compatibility).
func (s *State[T]) Value() T {
	return s.value
}

// Set updates the value and triggers re-render if changed.
func (s *State[T]) Set(newValue T) {
	s.setter(newValue)
}

// Update applies a function to update the value.
func (s *State[T]) Update(fn func(T) T) {
	s.Set(fn(s.value))
}

// UseState creates reactive state that persists across renders.
// The key parameter is optional - if not provided, hook order is used.
//
// Example:
//
//	count := UseState(ctx, 0)
//	count.Set(count.Get() + 1)
func UseState[T any](ctx *Context, initial T, key ...string) *State[T] {
	var hookKey string
	if len(key) > 0 {
		hookKey = ctx.getHookKey(key[0])
	} else {
		hookKey = ctx.getHookKey("")
	}

	state, exists := ctx.GetHook(hookKey)
	if !exists || !state.Initialized {
		state = &HookState{
			Value:       initial,
			Initialized: true,
		}
		ctx.SetHook(hookKey, state)
	}

	currentValue := state.Value.(T)

	setter := func(newValue T) {
		// Deep equality check to avoid unnecessary re-renders
		if reflect.DeepEqual(state.Value, newValue) {
			return
		}
		state.Value = newValue
		ctx.SetHook(hookKey, state)
		ctx.RequestRender()
	}

	return &State[T]{
		value:   currentValue,
		setter:  setter,
		hookKey: hookKey,
	}
}

// UseEffect runs side effects when dependencies change.
// Returns a cleanup function that runs before the next effect or on unmount.
//
// Behavior based on deps:
//   - deps = nil: runs on every render
//   - deps = []any{}: runs once on mount
//   - deps = []any{a, b}: runs when a or b change
//
// Example:
//
//	UseEffect(ctx, func() func() {
//	    // Setup
//	    return func() { /* Cleanup */ }
//	}, []any{dependency})
func UseEffect(ctx *Context, effect func() func(), deps []any, key ...string) {
	var hookKey string
	if len(key) > 0 {
		hookKey = ctx.getHookKey(key[0])
	} else {
		hookKey = ctx.getHookKey("")
	}

	state, exists := ctx.GetHook(hookKey)

	// deps == nil means run every render
	if deps == nil {
		// Run cleanup from previous render
		if state != nil && state.Cleanup != nil {
			state.Cleanup()
		}
		// Run effect and store cleanup
		cleanup := effect()
		ctx.SetHook(hookKey, &HookState{
			Cleanup:     cleanup,
			Initialized: true,
		})
		return
	}

	// Check if deps changed
	shouldRun := !exists || !state.Initialized || depsChanged(state.Dependencies, deps)

	if shouldRun {
		// Run cleanup from previous render
		if state != nil && state.Cleanup != nil {
			state.Cleanup()
		}

		// Run effect
		cleanup := effect()

		// Store new state
		ctx.SetHook(hookKey, &HookState{
			Dependencies: copyDeps(deps),
			Cleanup:      cleanup,
			Initialized:  true,
		})
	}
}

// UseEffectSimple is UseEffect without cleanup function.
//
// Example:
//
//	UseEffectSimple(ctx, func() {
//	    fmt.Println("Value changed:", value)
//	}, []any{value})
func UseEffectSimple(ctx *Context, effect func(), deps []any, key ...string) {
	UseEffect(ctx, func() func() {
		effect()
		return nil
	}, deps, key...)
}

// UseMemo memoizes expensive computations.
// Only recomputes when dependencies change.
//
// Example:
//
//	expensive := UseMemo(ctx, func() int {
//	    return heavyComputation(items)
//	}, []any{items})
func UseMemo[T any](ctx *Context, compute func() T, deps []any, key ...string) T {
	var hookKey string
	if len(key) > 0 {
		hookKey = ctx.getHookKey(key[0])
	} else {
		hookKey = ctx.getHookKey("")
	}

	state, exists := ctx.GetHook(hookKey)
	shouldCompute := !exists || !state.Initialized || depsChanged(state.Dependencies, deps)

	if shouldCompute {
		value := compute()
		ctx.SetHook(hookKey, &HookState{
			Value:        value,
			Dependencies: copyDeps(deps),
			Initialized:  true,
		})
		return value
	}

	return state.Value.(T)
}

// Ref holds a mutable value that persists across renders.
// Unlike State, changes to Ref do not trigger re-renders.
type Ref[T any] struct {
	Current T
}

// UseRef creates a persistent mutable reference.
// Use for storing values that shouldn't trigger re-renders (DOM refs, timers, etc.)
//
// Example:
//
//	inputRef := UseRef(ctx, "")
//	inputRef.Current = "new value" // No re-render
func UseRef[T any](ctx *Context, initial T, key ...string) *Ref[T] {
	var hookKey string
	if len(key) > 0 {
		hookKey = ctx.getHookKey(key[0])
	} else {
		hookKey = ctx.getHookKey("")
	}

	state, exists := ctx.GetHook(hookKey)
	if !exists || !state.Initialized {
		ref := &Ref[T]{Current: initial}
		ctx.SetHook(hookKey, &HookState{
			Value:       ref,
			Initialized: true,
		})
		return ref
	}

	return state.Value.(*Ref[T])
}

// UseCallback memoizes a callback function.
// Returns the same function reference unless dependencies change.
//
// Example:
//
//	handleClick := UseCallback(ctx, func() {
//	    doSomething(id)
//	}, []any{id})
func UseCallback[T any](ctx *Context, callback T, deps []any, key ...string) T {
	return UseMemo(ctx, func() T {
		return callback
	}, deps, key...)
}

// UseReducer manages complex state with a reducer function.
// Similar to React's useReducer.
//
// Example:
//
//	state, dispatch := UseReducer(ctx, reducer, initialState)
//	dispatch(Action{Type: "INCREMENT"})
func UseReducer[S any, A any](ctx *Context, reducer func(S, A) S, initial S, key ...string) (*S, func(A)) {
	var hookKey string
	if len(key) > 0 {
		hookKey = ctx.getHookKey(key[0])
	} else {
		hookKey = ctx.getHookKey("")
	}

	state, exists := ctx.GetHook(hookKey)
	if !exists || !state.Initialized {
		state = &HookState{
			Value:       initial,
			Initialized: true,
		}
		ctx.SetHook(hookKey, state)
	}

	currentState := state.Value.(S)

	dispatch := func(action A) {
		newState := reducer(currentState, action)
		if !reflect.DeepEqual(state.Value, newState) {
			state.Value = newState
			ctx.SetHook(hookKey, state)
			ctx.RequestRender()
		}
	}

	return &currentState, dispatch
}

// UsePrevious returns the value from the previous render.
//
// Example:
//
//	prevCount := UsePrevious(ctx, count)
//	if prevCount != nil && *prevCount != count {
//	    fmt.Println("Count changed from", *prevCount, "to", count)
//	}
func UsePrevious[T any](ctx *Context, value T, key ...string) *T {
	ref := UseRef(ctx, (*T)(nil), key...)
	previous := ref.Current

	// Update ref after render (via effect)
	UseEffectSimple(ctx, func() {
		ref.Current = &value
	}, []any{value})

	return previous
}

// UseToggle provides a boolean state with toggle helper.
//
// Example:
//
//	isOpen, toggle, setOpen := UseToggle(ctx, false)
//	toggle() // Flips the value
func UseToggle(ctx *Context, initial bool, key ...string) (bool, func(), func(bool)) {
	state := UseState(ctx, initial, key...)

	toggle := func() {
		state.Set(!state.Get())
	}

	return state.Get(), toggle, state.Set
}

// --- Helper Functions ---

// depsChanged checks if dependencies have changed.
func depsChanged(oldDeps, newDeps []any) bool {
	if len(oldDeps) != len(newDeps) {
		return true
	}

	for i := range newDeps {
		if !reflect.DeepEqual(oldDeps[i], newDeps[i]) {
			return true
		}
	}

	return false
}

// copyDeps creates a copy of dependencies slice.
func copyDeps(deps []any) []any {
	if deps == nil {
		return nil
	}
	copied := make([]any, len(deps))
	copy(copied, deps)
	return copied
}
