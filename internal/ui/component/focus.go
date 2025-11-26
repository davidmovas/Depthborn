package component

import (
	"strings"
	"sync"
)

// FocusDirection represents navigation direction
type FocusDirection int

const (
	FocusUp FocusDirection = iota
	FocusDown
	FocusLeft
	FocusRight
)

// FocusPosition represents 2D position of focusable component
type FocusPosition struct {
	X int // Column
	Y int // Row
}

// Focusable represents a component that can receive focus
type Focusable interface {
	GetFocusID() string
	GetFocusPosition() *FocusPosition
	GetHotkeys() []string
	CanReceiveFocus() bool
	ShouldAutoFocus() bool // NEW: whether this should be auto-focused
	OnFocus()
	OnBlur()
	OnActivate() bool
	IsInputComponent() bool
}

// FocusContext manages focus and navigation for a render scope
type FocusContext struct {
	scope      string
	focusables []Focusable
	focusIndex int
	hotkeys    map[string]Focusable // normalized key -> focusable
	inInput    bool
	grid       map[int]map[int]Focusable // 2D navigation grid
	mu         sync.RWMutex
}

// NewFocusContext creates new focus context
func NewFocusContext(scope string) *FocusContext {
	return &FocusContext{
		scope:      scope,
		focusables: make([]Focusable, 0),
		hotkeys:    make(map[string]Focusable),
		grid:       make(map[int]map[int]Focusable),
		focusIndex: -1,
	}
}

// Register adds focusable component to navigation
func (fc *FocusContext) Register(focusable Focusable) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if !focusable.CanReceiveFocus() {
		return
	}

	currentIndex := len(fc.focusables)
	fc.focusables = append(fc.focusables, focusable)

	// Register hotkeys (case-insensitive)
	for _, hotkey := range focusable.GetHotkeys() {
		if hotkey != "" {
			normalized := strings.ToLower(hotkey)
			fc.hotkeys[fc.scope+":"+normalized] = focusable
		}
	}

	// Add to 2D grid if position specified
	if pos := focusable.GetFocusPosition(); pos != nil {
		if fc.grid[pos.Y] == nil {
			fc.grid[pos.Y] = make(map[int]Focusable)
		}
		fc.grid[pos.Y][pos.X] = focusable
	}

	// Handle focus logic
	fc.handleFocusLogic(currentIndex, focusable)
}

// HandleKey processes keyboard input
func (fc *FocusContext) HandleKey(key string) bool {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if fc.handleNavigation(key) {
		return true
	}

	// Hotkeys only work outside input components
	if !fc.inInput {
		normalized := strings.ToLower(key)
		if focusable, exists := fc.hotkeys[fc.scope+":"+normalized]; exists {
			return focusable.OnActivate()
		}
	}

	return false
}

// ClearFocusables clears all focusable components
func (fc *FocusContext) ClearFocusables() {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.focusables = make([]Focusable, 0)
	fc.hotkeys = make(map[string]Focusable)
	fc.grid = make(map[int]map[int]Focusable)
}

// IsFocusedByID checks if focus is set to specified component
func (fc *FocusContext) IsFocusedByID(id string) bool {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	return fc.focusIndex >= 0 && fc.focusIndex < len(fc.focusables) &&
		fc.focusables[fc.focusIndex].GetFocusID() == id
}

// Scope returns scope identifier
func (fc *FocusContext) Scope() string {
	return fc.scope
}

// handleNavigation processes navigation keys
func (fc *FocusContext) handleNavigation(key string) bool {
	switch key {
	case "tab", "down":
		return fc.moveFocus(FocusDown)
	case "shift+tab", "up":
		return fc.moveFocus(FocusUp)
	case "left":
		return fc.moveFocus(FocusLeft)
	case "right":
		return fc.moveFocus(FocusRight)
	case "enter":
		if current := fc.currentLocked(); current != nil {
			return current.OnActivate()
		}
	}
	return false
}

// moveFocus changes focus in specified direction
func (fc *FocusContext) moveFocus(direction FocusDirection) bool {
	if len(fc.focusables) == 0 {
		return false
	}

	current := fc.currentLocked()
	if current == nil {
		fc.setFocusIndex(0)
		return true
	}

	// Try 2D navigation first, fallback to linear
	if pos := current.GetFocusPosition(); pos != nil {
		return fc.moveFocus2D(pos, direction)
	}
	return fc.moveFocusLinear(direction)
}

// moveFocus2D performs grid-based navigation
func (fc *FocusContext) moveFocus2D(currentPos *FocusPosition, direction FocusDirection) bool {
	targetX, targetY := currentPos.X, currentPos.Y

	switch direction {
	case FocusUp:
		targetY--
	case FocusDown:
		targetY++
	case FocusLeft:
		targetX--
	case FocusRight:
		targetX++
	}

	if row, exists := fc.grid[targetY]; exists {
		var target Focusable
		if target, exists = row[targetX]; exists {
			for i, f := range fc.focusables {
				if f.GetFocusID() == target.GetFocusID() {
					fc.setFocusIndex(i)
					return true
				}
			}
		}
	}
	return false
}

// moveFocusLinear performs list-based navigation
func (fc *FocusContext) moveFocusLinear(direction FocusDirection) bool {
	delta := 0
	switch direction {
	case FocusUp, FocusLeft:
		delta = -1
	case FocusDown, FocusRight:
		delta = 1
	}

	newIndex := (fc.focusIndex + delta + len(fc.focusables)) % len(fc.focusables)
	fc.setFocusIndex(newIndex)
	return true
}

// handleFocusLogic manages focus during component registration
func (fc *FocusContext) handleFocusLogic(index int, focusable Focusable) {
	// Auto-focus only if no focus set yet
	if focusable.ShouldAutoFocus() && fc.focusIndex == -1 {
		fc.setFocusIndex(index)
		return
	}

	// Restore focus to previous index
	if fc.focusIndex == index {
		fc.setFocusIndex(index)
	}
}

// setFocusIndex changes focused component
func (fc *FocusContext) setFocusIndex(index int) {
	if index < 0 || index >= len(fc.focusables) {
		return
	}

	// Blur current
	if fc.focusIndex >= 0 && fc.focusIndex < len(fc.focusables) {
		fc.focusables[fc.focusIndex].OnBlur()
	}

	// Focus new
	fc.focusIndex = index
	focused := fc.focusables[fc.focusIndex]
	focused.OnFocus()
	fc.inInput = focused.IsInputComponent()
}

// Helper methods with proper locking
func (fc *FocusContext) currentLocked() Focusable {
	if fc.focusIndex >= 0 && fc.focusIndex < len(fc.focusables) {
		return fc.focusables[fc.focusIndex]
	}
	return nil
}
