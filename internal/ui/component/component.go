package component

// Component represents a UI component that can be rendered
// Components are composable and can contain other components
type Component interface {
	// Render converts component to renderable output
	// For BubbleTea: returns string
	// For other renderers: might return different types
	Render(ctx *Context) string
}

// Func is a functional component (stateless)
type Func func(ctx *Context) string

// Render implements Component interface for Func
func (f Func) Render(ctx *Context) string {
	return f(ctx)
}

// Empty returns an empty component
func Empty() Component {
	return Func(func(ctx *Context) string {
		return ""
	})
}

// Raw creates component from raw string
func Raw(s string) Component {
	return Func(func(ctx *Context) string {
		return s
	})
}

// If conditionally renders component
func If(condition bool, component Component) Component {
	if condition {
		return component
	}
	return Empty()
}

// IfElse conditionally renders one of two components
func IfElse(condition bool, ifTrue Component, ifFalse Component) Component {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// Map renders a list of components from a slice
func Map[T any](items []T, fn func(item T, index int) Component) Component {
	return Func(func(ctx *Context) string {
		result := ""
		for i, item := range items {
			comp := fn(item, i)
			result += comp.Render(ctx)
		}
		return result
	})
}
