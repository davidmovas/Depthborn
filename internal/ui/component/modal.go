package component

import (
	"strings"
	"sync"
)

// PortalLayer represents rendering layer for portals
type PortalLayer string

const (
	LayerModal   PortalLayer = "modal"
	LayerToast   PortalLayer = "toast"
	LayerTooltip PortalLayer = "tooltip"
)

// PortalEntry represents single portal content
type PortalEntry struct {
	ID        string
	Layer     PortalLayer
	Component Component
	ZIndex    int
}

// PortalManager manages portal rendering across layers
type PortalManager struct {
	portals map[PortalLayer][]PortalEntry
	mu      sync.RWMutex
}

// NewPortalManager creates new portal manager
func NewPortalManager() *PortalManager {
	return &PortalManager{
		portals: make(map[PortalLayer][]PortalEntry),
	}
}

// Register adds component to portal layer
func (pm *PortalManager) Register(entry PortalEntry) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.portals[entry.Layer] == nil {
		pm.portals[entry.Layer] = make([]PortalEntry, 0)
	}

	pm.portals[entry.Layer] = append(pm.portals[entry.Layer], entry)
}

// Unregister removes portal by ID
func (pm *PortalManager) Unregister(id string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for layer := range pm.portals {
		entries := pm.portals[layer]
		for i, entry := range entries {
			if entry.ID == id {
				pm.portals[layer] = append(entries[:i], entries[i+1:]...)
				return
			}
		}
	}
}

// Render renders all portals for given layer
func (pm *PortalManager) Render(ctx *Context, layer PortalLayer) string {
	pm.mu.RLock()
	entries := pm.portals[layer]
	pm.mu.RUnlock()

	if len(entries) == 0 {
		return ""
	}

	var result strings.Builder
	for _, entry := range entries {
		content := entry.Component.Render(ctx)
		result.WriteString(content)
	}

	return result.String()
}

// Clear removes all portals
func (pm *PortalManager) Clear() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.portals = make(map[PortalLayer][]PortalEntry)
}

// HasPortals checks if layer has any portals
func (pm *PortalManager) HasPortals(layer PortalLayer) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return len(pm.portals[layer]) > 0
}

// PortalManager returns portal manager for context
func (ctx *Context) PortalManager() *PortalManager {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if pm, ok := ctx.meta["__portal_manager"].(*PortalManager); ok {
		return pm
	}

	pm := NewPortalManager()
	ctx.meta["__portal_manager"] = pm
	return pm
}

// SetPortalManager sets portal manager (usually called by renderer)
func (ctx *Context) SetPortalManager(pm *PortalManager) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.meta["__portal_manager"] = pm
}

// ModalState manages modal open/close state
type ModalState struct {
	IsOpen    bool
	OnClose   func()
	FocusTrap bool // Trap focus inside modal
}

// UseModal is a hook for modal state management
func UseModal(ctx *Context, initialOpen bool, key ...string) *State[ModalState] {
	return UseState(ctx, ModalState{
		IsOpen:    initialOpen,
		FocusTrap: true,
	}, key...)
}

// ModalManager manages modal stack and focus trapping
type ModalManager struct {
	stack []string // Stack of modal IDs (LIFO)
	mu    sync.RWMutex
}

// NewModalManager creates new modal manager
func NewModalManager() *ModalManager {
	return &ModalManager{
		stack: make([]string, 0),
	}
}

// Push adds modal to stack
func (mm *ModalManager) Push(modalID string) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.stack = append(mm.stack, modalID)
}

// Pop removes top modal from stack
func (mm *ModalManager) Pop() string {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if len(mm.stack) == 0 {
		return ""
	}

	top := mm.stack[len(mm.stack)-1]
	mm.stack = mm.stack[:len(mm.stack)-1]
	return top
}

// Top returns current top modal (without removing)
func (mm *ModalManager) Top() string {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if len(mm.stack) == 0 {
		return ""
	}

	return mm.stack[len(mm.stack)-1]
}

// IsActive checks if modal is currently active (on top)
func (mm *ModalManager) IsActive(modalID string) bool {
	return mm.Top() == modalID
}

// HasModals returns true if any modals are open
func (mm *ModalManager) HasModals() bool {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	return len(mm.stack) > 0
}

// Clear removes all modals
func (mm *ModalManager) Clear() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.stack = make([]string, 0)
}

// ModalManager returns modal manager for context
func (ctx *Context) ModalManager() *ModalManager {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if mm, ok := ctx.meta["__modal_manager"].(*ModalManager); ok {
		return mm
	}

	mm := NewModalManager()
	ctx.meta["__modal_manager"] = mm
	return mm
}

// SetModalManager sets modal manager
func (ctx *Context) SetModalManager(mm *ModalManager) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.meta["__modal_manager"] = mm
}
