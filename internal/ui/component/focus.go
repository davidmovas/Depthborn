package component

import (
	"fmt"
	"strings"
	"sync"
)

// FocusDirection represents navigation direction.
type FocusDirection int

const (
	FocusNext FocusDirection = iota
	FocusPrev
	FocusUp
	FocusDown
	FocusLeft
	FocusRight
)

// FocusPosition represents 2D position for spatial navigation.
type FocusPosition struct {
	Row int
	Col int
}

// Focusable represents a component that can receive focus.
type Focusable interface {
	// ID returns unique identifier for this focusable.
	ID() string

	// Position returns optional 2D position for grid navigation.
	Position() *FocusPosition

	// Hotkeys returns keys that activate this component.
	Hotkeys() []string

	// CanFocus returns whether this component can currently receive focus.
	CanFocus() bool

	// AutoFocus returns whether this should be focused on registration.
	AutoFocus() bool

	// IsInput returns whether this is an input component (blocks hotkeys when focused).
	IsInput() bool

	// OnFocus is called when component receives focus.
	OnFocus()

	// OnBlur is called when component loses focus.
	OnBlur()

	// OnActivate is called on Enter/click. Returns true if handled.
	OnActivate() bool

	// OnKeyPress is called for text input when this is focused and IsInput() is true.
	// Returns true if the key was handled.
	OnKeyPress(key string) bool
}

// FocusManager handles focus state and keyboard navigation.
type FocusManager struct {
	// Registered focusables in order
	items []Focusable

	// Current focus index (-1 = none)
	focusIndex int

	// Hotkey lookup (normalized key -> focusable)
	hotkeys map[string]Focusable

	// 2D grid for spatial navigation
	grid map[int]map[int]Focusable

	// Auto-assigned positions for automatic 2D navigation
	autoPositions map[string]*FocusPosition // ID -> position

	// Current row/column for auto-assignment during registration
	currentRow int
	currentCol int

	// Whether currently focused item is an input
	inInput bool

	// Items registered this frame (for diffing)
	frameItems []Focusable

	// Previous focus ID (for restoration)
	previousFocusID string

	// Callback to request re-render when focus changes
	onFocusChange func()

	mu sync.RWMutex
}

// NewFocusManager creates a new focus manager.
func NewFocusManager() *FocusManager {
	return &FocusManager{
		items:         make([]Focusable, 0),
		hotkeys:       make(map[string]Focusable),
		grid:          make(map[int]map[int]Focusable),
		autoPositions: make(map[string]*FocusPosition),
		focusIndex:    -1,
	}
}

// Reset clears all focus state for reuse.
func (fm *FocusManager) Reset() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.items = make([]Focusable, 0)
	fm.hotkeys = make(map[string]Focusable)
	fm.grid = make(map[int]map[int]Focusable)
	fm.autoPositions = make(map[string]*FocusPosition)
	fm.frameItems = nil
	fm.focusIndex = -1
	fm.previousFocusID = ""
	fm.currentRow = 0
	fm.currentCol = 0
	fm.inInput = false
}

// SetOnFocusChange sets the callback to be called when focus changes.
// This is used to trigger re-renders when navigating between elements.
func (fm *FocusManager) SetOnFocusChange(callback func()) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.onFocusChange = callback
}

// NextRow moves to the next row for automatic position assignment.
// Call this between groups of horizontal elements to enable proper up/down navigation.
// Example:
//
//	[Tab1] [Tab2] [Tab3] [Tab4]  <- row 0
//	ctx.Focus().NextRow()
//	[Btn1] [Btn2] [Btn3]         <- row 1
func (fm *FocusManager) NextRow() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.currentRow++
	fm.currentCol = 0
}

// BeginFrame prepares for a new render frame.
func (fm *FocusManager) BeginFrame() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Save current focus ID for restoration
	if fm.focusIndex >= 0 && fm.focusIndex < len(fm.items) {
		fm.previousFocusID = fm.items[fm.focusIndex].ID()
	}

	// Clear for re-registration
	fm.frameItems = make([]Focusable, 0)

	// Reset auto-position counters for new frame
	fm.currentRow = 0
	fm.currentCol = 0
	fm.autoPositions = make(map[string]*FocusPosition)
}

// EndFrame completes the render frame.
func (fm *FocusManager) EndFrame() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Replace items with this frame's registrations
	fm.items = fm.frameItems
	fm.frameItems = nil

	// Rebuild hotkey map
	fm.hotkeys = make(map[string]Focusable)
	for _, item := range fm.items {
		for _, hk := range item.Hotkeys() {
			if hk != "" {
				fm.hotkeys[strings.ToLower(hk)] = item
			}
		}
	}

	// Rebuild grid from explicit positions and auto-assigned positions
	fm.grid = make(map[int]map[int]Focusable)
	for _, item := range fm.items {
		var pos *FocusPosition

		// Use explicit position if available
		if p := item.Position(); p != nil {
			pos = p
		} else if p, ok := fm.autoPositions[item.ID()]; ok {
			// Use auto-assigned position
			pos = p
		}

		if pos != nil {
			if fm.grid[pos.Row] == nil {
				fm.grid[pos.Row] = make(map[int]Focusable)
			}
			fm.grid[pos.Row][pos.Col] = item
		}
	}

	// Restore focus
	fm.restoreFocus()
}

// Register adds a focusable component.
// Returns true if this component is currently focused.
func (fm *FocusManager) Register(f Focusable) bool {
	if !f.CanFocus() {
		return false
	}

	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.frameItems = append(fm.frameItems, f)
	currentIndex := len(fm.frameItems) - 1

	// Auto-assign position if not explicitly set
	if f.Position() == nil {
		fm.autoPositions[f.ID()] = &FocusPosition{
			Row: fm.currentRow,
			Col: fm.currentCol,
		}
		fm.currentCol++
	}

	// Check if this should be focused
	isFocused := false

	// Check if this is the currently focused item (by index)
	if fm.focusIndex == currentIndex {
		isFocused = true
	}

	// Check if this is the currently focused item (by ID from previous frame)
	if fm.previousFocusID != "" && f.ID() == fm.previousFocusID {
		fm.focusIndex = currentIndex
		isFocused = true
		fm.inInput = f.IsInput()
	}

	// Auto-focus if requested and no focus yet
	if f.AutoFocus() && fm.focusIndex == -1 && fm.previousFocusID == "" && !isFocused {
		fm.focusIndex = currentIndex
		isFocused = true
		fm.inInput = f.IsInput()
	}

	return isFocused
}

// ItemCount returns the number of registered focusable items.
func (fm *FocusManager) ItemCount() int {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return len(fm.items)
}

// IsFocused checks if a specific component is focused.
func (fm *FocusManager) IsFocused(id string) bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if fm.focusIndex < 0 || fm.focusIndex >= len(fm.items) {
		return false
	}

	return fm.items[fm.focusIndex].ID() == id
}

// Current returns the currently focused component.
func (fm *FocusManager) Current() Focusable {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if fm.focusIndex < 0 || fm.focusIndex >= len(fm.items) {
		return nil
	}

	return fm.items[fm.focusIndex]
}

// HandleKey processes keyboard input.
// Returns true if the key was handled.
func (fm *FocusManager) HandleKey(key string) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Skip if no items registered
	if len(fm.items) == 0 {
		return false
	}

	normalizedKey := strings.ToLower(key)

	// Navigation keys - arrows only for 4-directional navigation
	// These work even in input mode
	switch normalizedKey {
	case "up":
		return fm.moveFocus(FocusUp)
	case "down":
		return fm.moveFocus(FocusDown)
	case "tab":
		return fm.moveFocus(FocusNext)
	case "shift+tab":
		return fm.moveFocus(FocusPrev)
	}

	// If in input mode, route keys to the focused input
	if fm.inInput {
		if current := fm.currentLocked(); current != nil {
			// Enter submits
			if normalizedKey == "enter" {
				return current.OnActivate()
			}
			// Pass other keys to input handler
			return current.OnKeyPress(key)
		}
	}

	// Left/Right navigation (only when not in input mode)
	switch normalizedKey {
	case "left":
		return fm.moveFocus(FocusLeft)
	case "right":
		return fm.moveFocus(FocusRight)
	case "enter", " ":
		if current := fm.currentLocked(); current != nil {
			return current.OnActivate()
		}
	}

	// Hotkeys (only when not in input mode)
	if !fm.inInput {
		if focusable, exists := fm.hotkeys[normalizedKey]; exists {
			return focusable.OnActivate()
		}
	}

	return false
}

// FocusFirst focuses the first focusable component.
func (fm *FocusManager) FocusFirst() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if len(fm.items) > 0 {
		fm.setFocusIndex(0)
	}
}

// FocusID focuses a component by ID.
func (fm *FocusManager) FocusID(id string) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for i, item := range fm.items {
		if item.ID() == id {
			fm.setFocusIndex(i)
			return true
		}
	}

	return false
}

// ClearFocus removes focus from all components.
func (fm *FocusManager) ClearFocus() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.focusIndex >= 0 && fm.focusIndex < len(fm.items) {
		fm.items[fm.focusIndex].OnBlur()
	}

	fm.focusIndex = -1
	fm.inInput = false
}

// --- Internal Methods ---

func (fm *FocusManager) restoreFocus() {
	// Try to restore by ID
	if fm.previousFocusID != "" {
		for i, item := range fm.items {
			if item.ID() == fm.previousFocusID {
				fm.focusIndex = i
				fm.inInput = item.IsInput()
				return
			}
		}
	}

	// Try to restore by index
	if fm.focusIndex >= len(fm.items) {
		fm.focusIndex = len(fm.items) - 1
	}

	// Focus first auto-focus item
	if fm.focusIndex < 0 {
		for i, item := range fm.items {
			if item.AutoFocus() {
				fm.focusIndex = i
				fm.inInput = item.IsInput()
				return
			}
		}
	}

	// Default to first item if nothing focused
	if fm.focusIndex < 0 && len(fm.items) > 0 {
		fm.focusIndex = 0
		fm.inInput = fm.items[0].IsInput()
	}
}

func (fm *FocusManager) moveFocus(direction FocusDirection) bool {
	if len(fm.items) == 0 {
		return false
	}

	// If nothing focused, focus first item
	if fm.focusIndex < 0 {
		fm.setFocusIndex(0)
		return true
	}

	current := fm.currentLocked()
	if current == nil {
		fm.setFocusIndex(0)
		return true
	}

	// Get position - explicit or auto-assigned
	var pos *FocusPosition
	if p := current.Position(); p != nil {
		pos = p
	} else if p, ok := fm.autoPositions[current.ID()]; ok {
		pos = p
	}

	// Try 2D grid navigation if we have a position
	if pos != nil {
		if fm.moveFocus2D(pos, direction) {
			return true
		}
	}

	// Fall back to linear navigation - all arrows work linearly
	return fm.moveFocusLinear(direction)
}

// getLinearDelta returns direction delta for linear navigation.
// Right/Down/Next = forward (+1), Left/Up/Prev = backward (-1)
func (fm *FocusManager) getLinearDelta(direction FocusDirection) int {
	switch direction {
	case FocusRight, FocusDown, FocusNext:
		return 1
	case FocusLeft, FocusUp, FocusPrev:
		return -1
	default:
		return 1
	}
}

func (fm *FocusManager) moveFocus2D(pos *FocusPosition, direction FocusDirection) bool {
	targetRow, targetCol := pos.Row, pos.Col

	switch direction {
	case FocusUp:
		targetRow--
	case FocusDown:
		targetRow++
	case FocusLeft:
		targetCol--
	case FocusRight:
		targetCol++
	case FocusNext:
		// Try right first, then down
		if _, exists := fm.grid[targetRow][targetCol+1]; exists {
			targetCol++
		} else {
			targetRow++
			targetCol = 0
		}
	case FocusPrev:
		// Try left first, then up
		if _, exists := fm.grid[targetRow][targetCol-1]; exists {
			targetCol--
		} else {
			targetRow--
			// Find rightmost in previous row
			if row, exists := fm.grid[targetRow]; exists {
				maxCol := -1
				for c := range row {
					if c > maxCol {
						maxCol = c
					}
				}
				targetCol = maxCol
			}
		}
	}

	// Check if target row exists
	row, rowExists := fm.grid[targetRow]
	if !rowExists {
		return false
	}

	// Try exact column first
	if target, exists := row[targetCol]; exists {
		return fm.focusByID(target.ID())
	}

	// For Up/Down: find nearest column in target row
	if direction == FocusUp || direction == FocusDown {
		nearest := fm.findNearestInRow(row, targetCol)
		if nearest != nil {
			return fm.focusByID(nearest.ID())
		}
	}

	return false
}

// findNearestInRow finds the item in row closest to targetCol
func (fm *FocusManager) findNearestInRow(row map[int]Focusable, targetCol int) Focusable {
	if len(row) == 0 {
		return nil
	}

	var nearest Focusable
	minDist := -1

	for col, item := range row {
		dist := col - targetCol
		if dist < 0 {
			dist = -dist
		}
		if minDist < 0 || dist < minDist {
			minDist = dist
			nearest = item
		}
	}

	return nearest
}

// focusByID sets focus to item by ID
func (fm *FocusManager) focusByID(id string) bool {
	for i, item := range fm.items {
		if item.ID() == id {
			fm.setFocusIndex(i)
			return true
		}
	}
	return false
}

// DebugGrid returns a string representation of the focus grid for debugging
func (fm *FocusManager) DebugGrid() string {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	result := "Focus Grid:\n"
	for row := 0; row < 10; row++ {
		if r, exists := fm.grid[row]; exists {
			result += fmt.Sprintf("  Row %d: ", row)
			for col := 0; col < 10; col++ {
				if item, exists := r[col]; exists {
					result += fmt.Sprintf("[%d:%s] ", col, item.ID()[:8])
				}
			}
			result += "\n"
		}
	}
	result += fmt.Sprintf("Current: idx=%d, items=%d\n", fm.focusIndex, len(fm.items))
	return result
}

// CurrentPosition returns the position of currently focused item
func (fm *FocusManager) CurrentPosition() *FocusPosition {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if fm.focusIndex < 0 || fm.focusIndex >= len(fm.items) {
		return nil
	}

	current := fm.items[fm.focusIndex]

	// Check explicit position
	if pos := current.Position(); pos != nil {
		return pos
	}

	// Search in grid for this item
	for row, cols := range fm.grid {
		for col, item := range cols {
			if item.ID() == current.ID() {
				return &FocusPosition{Row: row, Col: col}
			}
		}
	}

	return nil
}

// FocusRowBreak is a component that advances to the next focus row when rendered.
// Use this between rows of focusable elements.
func FocusRowBreak() Component {
	return Func(func(ctx *Context) string {
		ctx.Focus().NextRow()
		return ""
	})
}

func (fm *FocusManager) moveFocusLinear(direction FocusDirection) bool {
	if len(fm.items) == 0 {
		return false
	}

	delta := fm.getLinearDelta(direction)
	newIndex := fm.focusIndex + delta

	// Wrap around
	if newIndex < 0 {
		newIndex = len(fm.items) - 1
	} else if newIndex >= len(fm.items) {
		newIndex = 0
	}

	fm.setFocusIndex(newIndex)
	return true
}

func (fm *FocusManager) setFocusIndex(index int) {
	if index < 0 || index >= len(fm.items) {
		return
	}

	// Skip if same index
	if fm.focusIndex == index {
		return
	}

	// Blur current
	if fm.focusIndex >= 0 && fm.focusIndex < len(fm.items) {
		fm.items[fm.focusIndex].OnBlur()
	}

	// Focus new
	fm.focusIndex = index
	focused := fm.items[fm.focusIndex]
	focused.OnFocus()
	fm.inInput = focused.IsInput()

	// Request re-render to show focus change
	if fm.onFocusChange != nil {
		fm.onFocusChange()
	}
}

func (fm *FocusManager) currentLocked() Focusable {
	if fm.focusIndex >= 0 && fm.focusIndex < len(fm.items) {
		return fm.items[fm.focusIndex]
	}
	return nil
}
