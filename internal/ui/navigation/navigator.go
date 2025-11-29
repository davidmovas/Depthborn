package navigation

import (
	"errors"
)

var (
	ErrNoScreen      = errors.New("no screen active")
	ErrCannotClose   = errors.New("screen cannot be closed")
	ErrScreenUnknown = errors.New("screen not registered")
)

// Navigator manages screen navigation and lifecycle.
type Navigator struct {
	registry *Registry
	stack    *Stack
}

// NewNavigator creates a new navigator.
func NewNavigator() *Navigator {
	return &Navigator{
		registry: NewRegistry(),
		stack:    NewStack(),
	}
}

// Register adds a screen factory to the registry.
func (n *Navigator) Register(screenID string, factory ScreenFactory) {
	n.registry.Register(screenID, factory)
}

// Open creates and pushes a screen to the stack.
func (n *Navigator) Open(screenID string, params map[string]any) error {
	screen, err := n.registry.Create(screenID)
	if err != nil {
		return err
	}

	// Pause current screen if any
	if current := n.stack.Peek(); current != nil {
		current.OnPause()
	}

	// Initialize and push new screen
	screen.OnInit()
	n.stack.Push(screen)
	screen.OnEnter(params)

	return nil
}

// Close pops the current screen from the stack.
func (n *Navigator) Close() error {
	current := n.stack.Peek()
	if current == nil {
		return ErrNoScreen
	}

	if !current.CanClose() {
		return ErrCannotClose
	}

	// Call exit lifecycle
	current.OnExit()
	n.stack.Pop()

	// Resume previous screen if any
	if prev := n.stack.Peek(); prev != nil {
		prev.OnResume()
	}

	return nil
}

// Back is an alias for Close.
func (n *Navigator) Back() error {
	return n.Close()
}

// CanGoBack returns whether navigation back is possible.
func (n *Navigator) CanGoBack() bool {
	return n.stack.Size() > 1
}

// Switch replaces the current screen with a new one.
func (n *Navigator) Switch(screenID string, params map[string]any) error {
	screen, err := n.registry.Create(screenID)
	if err != nil {
		return err
	}

	// Exit current screen if any
	if current := n.stack.Peek(); current != nil {
		current.OnExit()
	}

	// Initialize and replace
	screen.OnInit()
	n.stack.Replace(screen)
	screen.OnEnter(params)

	return nil
}

// GoTo is an alias for Switch.
func (n *Navigator) GoTo(screenID string, params map[string]any) error {
	return n.Switch(screenID, params)
}

// Reset clears the stack and opens a single screen.
func (n *Navigator) Reset(screenID string, params map[string]any) error {
	// Exit all screens
	for !n.stack.IsEmpty() {
		if current := n.stack.Peek(); current != nil {
			current.OnExit()
		}
		n.stack.Pop()
	}

	return n.Open(screenID, params)
}

// Clear removes all screens from the stack.
func (n *Navigator) Clear() {
	for !n.stack.IsEmpty() {
		if current := n.stack.Peek(); current != nil {
			current.OnExit()
		}
		n.stack.Pop()
	}
}

// CurrentScreen returns the current (top) screen.
func (n *Navigator) CurrentScreen() Screen {
	return n.stack.Peek()
}

// Current is an alias for CurrentScreen.
func (n *Navigator) Current() Screen {
	return n.CurrentScreen()
}

// StackSize returns the number of screens in the stack.
func (n *Navigator) StackSize() int {
	return n.stack.Size()
}

// HasScreens returns true if there are any screens.
func (n *Navigator) HasScreens() bool {
	return !n.stack.IsEmpty()
}

// Update calls update lifecycle on the current screen.
func (n *Navigator) Update() {
	current := n.stack.Peek()
	if current == nil {
		return
	}

	current.OnUpdateStart()
	current.OnUpdate()
	current.OnUpdateEnd()
}

// Registry returns the screen registry.
func (n *Navigator) Registry() *Registry {
	return n.registry
}

// Stack returns the screen stack.
func (n *Navigator) Stack() *Stack {
	return n.stack
}
