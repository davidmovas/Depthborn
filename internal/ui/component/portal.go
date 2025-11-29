package component

import (
	"sort"
	"sync"
)

// PortalLayer represents the rendering layer for portals.
type PortalLayer int

const (
	LayerBase    PortalLayer = 0
	LayerOverlay PortalLayer = 100
	LayerModal   PortalLayer = 200
	LayerToast   PortalLayer = 300
	LayerTooltip PortalLayer = 400
)

// Portal represents a component rendered in a specific layer.
type Portal struct {
	ID        string
	Layer     PortalLayer
	ZIndex    int // For ordering within same layer
	Component Component
	Focus     *FocusManager // Isolated focus for this portal
	IsOpen    bool
}

// PortalManager manages portal rendering layers.
type PortalManager struct {
	portals         map[string]*Portal
	onRenderRequest func() // Callback to request re-render
	mu              sync.RWMutex
}

// NewPortalManager creates a new portal manager.
func NewPortalManager() *PortalManager {
	return &PortalManager{
		portals: make(map[string]*Portal),
	}
}

// SetOnRenderRequest sets the callback for render requests.
func (pm *PortalManager) SetOnRenderRequest(callback func()) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.onRenderRequest = callback
}

// Open registers and opens a portal.
// If the portal already exists and is open, it updates the component but keeps the focus manager.
func (pm *PortalManager) Open(id string, layer PortalLayer, comp Component) *Portal {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if portal already exists
	if existing, exists := pm.portals[id]; exists {
		// Update component and ensure it's open, but keep existing focus manager
		existing.Component = comp
		existing.IsOpen = true
		return existing
	}

	// Create new portal with fresh focus manager
	focus := NewFocusManager()
	// Connect portal's focus manager to trigger re-renders
	if pm.onRenderRequest != nil {
		focus.SetOnFocusChange(pm.onRenderRequest)
	}

	portal := &Portal{
		ID:        id,
		Layer:     layer,
		ZIndex:    0,
		Component: comp,
		Focus:     focus,
		IsOpen:    true,
	}

	pm.portals[id] = portal
	return portal
}

// OpenWithZIndex registers a portal with specific z-index.
func (pm *PortalManager) OpenWithZIndex(id string, layer PortalLayer, zIndex int, comp Component) *Portal {
	portal := pm.Open(id, layer, comp)
	portal.ZIndex = zIndex
	return portal
}

// Close closes a portal (keeps it registered but hidden).
// The focus manager is reset so reopening starts fresh.
func (pm *PortalManager) Close(id string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if portal, exists := pm.portals[id]; exists {
		portal.IsOpen = false
		// Reset focus state so reopening starts fresh
		portal.Focus.Reset()
	}
}

// Remove completely removes a portal.
func (pm *PortalManager) Remove(id string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.portals, id)
}

// Get returns a portal by ID.
func (pm *PortalManager) Get(id string) *Portal {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.portals[id]
}

// IsOpen checks if a portal is open.
func (pm *PortalManager) IsOpen(id string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if portal, exists := pm.portals[id]; exists {
		return portal.IsOpen
	}
	return false
}

// Toggle toggles a portal's open state.
func (pm *PortalManager) Toggle(id string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if portal, exists := pm.portals[id]; exists {
		portal.IsOpen = !portal.IsOpen
		return portal.IsOpen
	}
	return false
}

// GetOpenPortals returns all open portals sorted by layer and z-index.
func (pm *PortalManager) GetOpenPortals() []*Portal {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var open []*Portal
	for _, portal := range pm.portals {
		if portal.IsOpen {
			open = append(open, portal)
		}
	}

	// Sort by layer, then z-index
	sort.Slice(open, func(i, j int) bool {
		if open[i].Layer != open[j].Layer {
			return open[i].Layer < open[j].Layer
		}
		return open[i].ZIndex < open[j].ZIndex
	})

	return open
}

// GetTopPortal returns the topmost open portal.
func (pm *PortalManager) GetTopPortal() *Portal {
	portals := pm.GetOpenPortals()
	if len(portals) == 0 {
		return nil
	}
	return portals[len(portals)-1]
}

// GetActiveFocus returns the focus manager for the topmost portal,
// or nil if no portals are open.
func (pm *PortalManager) GetActiveFocus() *FocusManager {
	if top := pm.GetTopPortal(); top != nil {
		return top.Focus
	}
	return nil
}

// HasOpenPortals returns whether any portals are open.
func (pm *PortalManager) HasOpenPortals() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, portal := range pm.portals {
		if portal.IsOpen {
			return true
		}
	}
	return false
}

// Clear removes all portals.
func (pm *PortalManager) Clear() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.portals = make(map[string]*Portal)
}

// RenderPortals renders all open portals and returns combined output.
func (pm *PortalManager) RenderPortals(ctx *Context) string {
	portals := pm.GetOpenPortals()
	if len(portals) == 0 {
		return ""
	}

	var result string
	for _, portal := range portals {
		if portal.Component == nil {
			continue
		}

		// Create child context with portal's isolated focus
		portalCtx := ctx.WithKey("portal_" + portal.ID)
		portalCtx.SetFocusManager(portal.Focus)

		// Begin focus frame for this portal
		portal.Focus.BeginFrame()

		// Render portal content
		result += portal.Component.Render(portalCtx)

		// End focus frame
		portal.Focus.EndFrame()
	}

	return result
}

// --- Portal Component ---

// PortalProps configures the Portal component.
type PortalProps struct {
	ID       string
	Layer    PortalLayer
	ZIndex   int
	Open     bool
	OnClose  func()
	Children Component
}

// PortalComponent renders content in a portal layer.
func PortalComponent(props PortalProps) Component {
	return Func(func(ctx *Context) string {
		pm := ctx.Portals()

		if props.Open {
			// Register/update portal
			portal := pm.Get(props.ID)
			if portal == nil {
				portal = pm.OpenWithZIndex(props.ID, props.Layer, props.ZIndex, props.Children)
			} else {
				portal.Component = props.Children
				portal.IsOpen = true
			}
		} else {
			// Close portal
			pm.Close(props.ID)
		}

		// Portal content is rendered separately by PortalManager.RenderPortals()
		return ""
	})
}

// --- Modal Hook ---

// UsePortal provides portal state management.
func UsePortal(ctx *Context, id string) (isOpen bool, open func(), close func(), toggle func()) {
	state := UseState(ctx, false, "portal_"+id)

	isOpen = state.Get()

	open = func() {
		state.Set(true)
	}

	close = func() {
		state.Set(false)
	}

	toggle = func() {
		state.Set(!state.Get())
	}

	return
}
