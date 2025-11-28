package component

import (
	"sync"
)

// PortalContext manages portal rendering layers
type PortalContext struct {
	portals map[PortalLayer][]*PortalEntry
	mu      sync.RWMutex
}

// PortalLayer represents rendering layer
type PortalLayer string

const (
	LayerModal   PortalLayer = "modal"
	LayerToast   PortalLayer = "toast"
	LayerTooltip PortalLayer = "tooltip"
)

// PortalEntry represents a portal instance
type PortalEntry struct {
	ID           string
	Component    Component
	IsOpen       bool
	FocusContext *FocusContext // ISOLATED focus context per portal
}

func NewPortalContext() *PortalContext {
	return &PortalContext{
		portals: make(map[PortalLayer][]*PortalEntry),
	}
}

// Register adds portal to layer with isolated focus context
func (pc *PortalContext) Register(layer PortalLayer, id string, comp Component, isOpen bool) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Create isolated focus context for this portal
	isolatedFocus := NewFocusContext(id + "_focus")

	entry := &PortalEntry{
		ID:           id,
		Component:    comp,
		IsOpen:       isOpen,
		FocusContext: isolatedFocus,
	}

	if pc.portals[layer] == nil {
		pc.portals[layer] = make([]*PortalEntry, 0)
	}

	// Check if already exists
	for i, p := range pc.portals[layer] {
		if p.ID == id {
			pc.portals[layer][i] = entry
			return
		}
	}

	pc.portals[layer] = append(pc.portals[layer], entry)
}

// Unregister removes portal from layer
func (pc *PortalContext) Unregister(layer PortalLayer, id string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	portals := pc.portals[layer]
	for i, p := range portals {
		if p.ID == id {
			pc.portals[layer] = append(portals[:i], portals[i+1:]...)
			return
		}
	}
}

// Get returns portal by ID
func (pc *PortalContext) Get(layer PortalLayer, id string) *PortalEntry {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	for _, p := range pc.portals[layer] {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// GetOpen returns all open portals in layer
func (pc *PortalContext) GetOpen(layer PortalLayer) []*PortalEntry {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	open := make([]*PortalEntry, 0)
	for _, p := range pc.portals[layer] {
		if p.IsOpen {
			open = append(open, p)
		}
	}
	return open
}

// GetActiveFocusContext returns focus context of topmost open portal, or nil
func (pc *PortalContext) GetActiveFocusContext() *FocusContext {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	// Check layers in priority order: Modal > Toast > Tooltip
	layers := []PortalLayer{LayerModal, LayerToast, LayerTooltip}

	for _, layer := range layers {
		portals := pc.portals[layer]
		for i := len(portals) - 1; i >= 0; i-- {
			if portals[i].IsOpen {
				return portals[i].FocusContext
			}
		}
	}

	return nil
}

// Clear removes all portals
func (pc *PortalContext) Clear() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.portals = make(map[PortalLayer][]*PortalEntry)
}

// Modal creates modal component with isolated focus
func Modal(ctx *Context, id string, isOpen bool, content Component) Component {
	return Func(func(ctx *Context) string {
		if !isOpen {
			return ""
		}

		// Create child context with isolated focus
		modalCtx := ctx.Child(id + "_modal")
		modalCtx.focusContext = NewFocusContext(id + "_focus_isolated") // ISOLATED

		// Register portal
		ctx.PortalContext().Register(LayerModal, id, content, isOpen)

		// Render backdrop + content
		backdrop := "╔══════════════════════════════════════╗\n"
		backdrop += "║          MODAL WINDOW                ║\n"
		backdrop += "╠══════════════════════════════════════╣\n"

		rendered := content.Render(modalCtx)
		backdrop += rendered

		backdrop += "╚══════════════════════════════════════╝"

		return backdrop
	})
}

// UseModal hook for managing modal state
func UseModal(ctx *Context, initialOpen bool, id string) *State[bool] {
	return UseState(ctx, initialOpen, "modal_"+id)
}
