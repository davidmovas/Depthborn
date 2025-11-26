package navigation

import (
	"fmt"
	"time"
)

// Navigator manages screen navigation and lifecycle
type Navigator struct {
	registry *Registry
	stack    *Stack
}

func NewNavigator() *Navigator {
	return &Navigator{
		registry: NewRegistry(),
		stack:    NewStack(),
	}
}

// Register adds screen factory to registry
func (n *Navigator) Register(screenID string, factory ScreenFactory) {
	n.registry.Register(screenID, factory)
}

// Open creates and opens screen (pushes to stack)
// params: optional parameters passed to OnEnter
func (n *Navigator) Open(screenID string, params map[string]any) error {
	screen, err := n.registry.Create(screenID)
	if err != nil {
		return err
	}

	// Push to stack
	n.stack.Push(screen)

	// Call OnEnter with params
	screen.OnEnter(params)

	return nil
}

// Close closes current screen (pops from stack)
// Returns error if screen cannot be closed
func (n *Navigator) Close() error {
	current := n.stack.Peek()
	if current == nil {
		return fmt.Errorf("no screen to close")
	}

	if !current.CanClose() {
		return fmt.Errorf("screen '%s' cannot be closed", current.ID())
	}

	n.stack.Pop()
	return nil
}

// Back is alias for Close (more intuitive naming)
func (n *Navigator) Back() error {
	return n.Close()
}

// Switch replaces current screen with new screen
// params: optional parameters passed to OnEnter
func (n *Navigator) Switch(screenID string, params map[string]any) error {
	screen, err := n.registry.Create(screenID)
	if err != nil {
		return err
	}

	// Replace top screen
	n.stack.Replace(screen)

	// Call OnEnter with params
	screen.OnEnter(params)

	return nil
}

// GoTo is alias for Switch
func (n *Navigator) GoTo(screenID string, params map[string]any) error {
	return n.Switch(screenID, params)
}

// Reset clears stack and opens single screen
// Useful for returning to main menu or restarting
func (n *Navigator) Reset(screenID string, params map[string]any) error {
	// Clear stack
	n.stack.Clear()

	// Open new screen
	return n.Open(screenID, params)
}

// Clear removes all screens
func (n *Navigator) Clear() {
	n.stack.Clear()
}

// Current returns current (top) screen
// Returns nil if no screens
func (n *Navigator) Current() Screen {
	return n.stack.Peek()
}

// StackSize returns number of screens in stack
func (n *Navigator) StackSize() int {
	return n.stack.Size()
}

// HasScreens returns true if there are any screens
func (n *Navigator) HasScreens() bool {
	return !n.stack.IsEmpty()
}

// Update calls update lifecycle hooks on current screen
// dt: delta time since last update
func (n *Navigator) Update(dt time.Duration) {
	current := n.stack.Peek()
	if current == nil {
		return
	}

	current.OnUpdateStart()
	current.OnUpdate(dt)
	current.OnUpdateEnd()
}

// Registry returns the screen registry
func (n *Navigator) Registry() *Registry {
	return n.registry
}

// Stack returns the screen stack
func (n *Navigator) Stack() *Stack {
	return n.stack
}
