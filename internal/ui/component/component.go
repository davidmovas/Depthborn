package component

// Component represents a renderable UI element.
// Components are composable and form a tree structure.
type Component interface {
	Render(ctx *Context) string
}

// Func is a functional component - the simplest component type.
// Use this for stateless rendering or components that use hooks.
type Func func(ctx *Context) string

func (f Func) Render(ctx *Context) string {
	return f(ctx)
}

// Empty returns a component that renders nothing.
func Empty() Component {
	return Func(func(ctx *Context) string {
		return ""
	})
}

// Raw creates a component from a static string.
// Use for pre-rendered content that doesn't need context.
func Raw(s string) Component {
	return Func(func(ctx *Context) string {
		return s
	})
}

// Fragment combines multiple components into one without wrapping.
func Fragment(children ...Component) Component {
	return Func(func(ctx *Context) string {
		var result string
		for _, child := range children {
			if child != nil {
				result += child.Render(ctx)
			}
		}
		return result
	})
}

// If conditionally renders a component.
func If(condition bool, component Component) Component {
	if condition {
		return component
	}
	return Empty()
}

// IfElse conditionally renders one of two components.
func IfElse(condition bool, ifTrue, ifFalse Component) Component {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// Switch renders component based on value matching.
func Switch[T comparable](value T, cases map[T]Component, defaultCase Component) Component {
	if comp, ok := cases[value]; ok {
		return comp
	}
	if defaultCase != nil {
		return defaultCase
	}
	return Empty()
}

// Map renders a list of items using a render function.
// The key function should return a stable unique identifier for each item.
func Map[T any](items []T, keyFn func(T, int) string, renderFn func(T, int) Component) Component {
	return Func(func(ctx *Context) string {
		var result string
		for i, item := range items {
			key := keyFn(item, i)
			childCtx := ctx.WithKey(key)
			comp := renderFn(item, i)
			if comp != nil {
				result += comp.Render(childCtx)
			}
		}
		return result
	})
}

// MapSimple is a simplified Map without key function (uses index as key).
// Use Map with proper keys for lists that can reorder.
func MapSimple[T any](items []T, renderFn func(T, int) Component) Component {
	return Func(func(ctx *Context) string {
		var result string
		for i, item := range items {
			childCtx := ctx.WithIndex(i)
			comp := renderFn(item, i)
			if comp != nil {
				result += comp.Render(childCtx)
			}
		}
		return result
	})
}
