package tea

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
	"github.com/davidmovas/Depthborn/internal/ui/renderer"
)

// Model bridges component system with BubbleTea with modal support
type Model struct {
	navigator       *navigation.Navigator
	lastTick        time.Time
	width, height   int
	renderDebounce  *renderer.Debouncer
	program         *tea.Program
	screenContexts  map[string]*component.Context
	currentFocusCtx *component.FocusContext

	// Modal & Portal support
	portalManager *component.PortalManager
	modalManager  *component.ModalManager
}

// NewModel creates new BubbleTea model with modal support
func NewModel(navigator *navigation.Navigator) *Model {
	m := &Model{
		navigator:      navigator,
		lastTick:       time.Now(),
		screenContexts: make(map[string]*component.Context, 20),
		portalManager:  component.NewPortalManager(),
		modalManager:   component.NewModalManager(),
	}

	// Setup debounced re-render (16ms = ~60fps max)
	m.renderDebounce = renderer.NewDebouncer(16*time.Millisecond, func() {
		if m.program != nil {
			m.program.Send(renderMsg{})
		}
	})

	return m
}

// SetProgram sets reference to tea.Program (for sending messages)
func (m *Model) SetProgram(p *tea.Program) {
	m.program = p
}

// RequestRender triggers debounced re-render
func (m *Model) RequestRender() {
	m.renderDebounce.Call()
}

// Init implements tea.Model.Init
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

// Update implements tea.Model.Update with modal-aware input handling
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {

	case tea.KeyMsg:
		key := message.String()

		// Global hotkeys (always work)
		switch key {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit

		case "esc":
			// If modal is open, close modal instead of going back
			if m.modalManager.HasModals() {
				m.modalManager.Pop()
				m.portalManager.Unregister(m.modalManager.Top())
				return m, nil
			}

			// Otherwise navigate back
			if err := m.navigator.Back(); err != nil {
				if m.navigator.StackSize() <= 1 {
					return m, tea.Quit
				}
			}
			return m, nil
		}

		screen := m.navigator.Current()
		if screen != nil {
			// Ensure focus context exists BEFORE handling input
			m.ensureFocusContext(screen)

			// Handle key directly with stored focus context
			if m.currentFocusCtx != nil && m.currentFocusCtx.HandleKey(key) {
				return m, nil
			}

			// Fall back to screen's HandleInput
			ctx := m.createContextForScreen(screen)
			handled := screen.HandleInput(ctx, message)
			if handled {
				return m, nil
			}
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height

		// Update screen size in all cached contexts
		for _, ctx := range m.screenContexts {
			ctx.SetScreenSize(m.width, m.height)
		}

		return m, nil

	case tickMsg:
		now := time.Now()
		dt := now.Sub(m.lastTick)
		m.lastTick = now

		m.navigator.Update(dt)

		return m, tickCmd()

	case renderMsg:
		return m, nil
	}

	return m, nil
}

// View implements tea.Model.View with layered rendering
func (m *Model) View() string {
	screen := m.navigator.Current()
	if screen == nil {
		return "No active screen"
	}

	// Ensure focus context exists
	m.ensureFocusContext(screen)

	// Clear and re-register components
	m.currentFocusCtx.ClearFocusables()

	// Clear portals before rendering
	m.portalManager.Clear()

	// Create context with all managers
	ctx := m.createContextForScreen(screen)

	// Set screen size
	ctx.SetScreenSize(m.width, m.height)

	// Set portal and modal managers
	ctx.SetPortalManager(m.portalManager)
	ctx.SetModalManager(m.modalManager)

	// Render main screen content
	comp := screen.Render(ctx)
	if comp == nil {
		return ""
	}

	mainContent := comp.Render(ctx)

	// Render portal layers on top
	modalContent := m.portalManager.Render(ctx, component.LayerModal)
	toastContent := m.portalManager.Render(ctx, component.LayerToast)
	tooltipContent := m.portalManager.Render(ctx, component.LayerTooltip)

	// Combine layers
	result := mainContent

	if modalContent != "" {
		// Overlay modal on top of main content
		result = overlayContent(result, modalContent, m.width, m.height)
	}

	if toastContent != "" {
		// Position toasts at bottom
		result = appendContent(result, toastContent)
	}

	if tooltipContent != "" {
		// Position tooltips at cursor (simplified: append for now)
		result = appendContent(result, tooltipContent)
	}

	return result
}

// ensureFocusContext ensures focus context exists for current screen
func (m *Model) ensureFocusContext(screen navigation.Screen) {
	if m.currentFocusCtx == nil || m.currentFocusCtx.Scope() != screen.ID() {
		m.currentFocusCtx = component.NewFocusContext(screen.ID())
	}
}

// createContextForScreen returns cached or new context for screen
func (m *Model) createContextForScreen(screen navigation.Screen) *component.Context {
	screenID := screen.ID()

	// Use cached context if available
	if ctx, exists := m.screenContexts[screenID]; exists {
		ctx.ResetHookIndex()
		return ctx
	}

	// Create and cache new context
	ctx := component.NewContext(screenID)
	if baseScreen, ok := screen.(*navigation.BaseScreen); ok {
		ctx.SetHooksContainer(baseScreen)
	}

	ctx.SetOnStateChange(m.RequestRender)
	if m.currentFocusCtx != nil {
		ctx.SetFocusContext(m.currentFocusCtx)
	}

	m.screenContexts[screenID] = ctx
	return ctx
}

// overlayContent overlays modal content on top of base content
// Creates dimmed background effect
func overlayContent(base, overlay string, width, height int) string {
	// For now, simple append - can be improved with actual overlay logic
	// TODO: Implement proper z-index layering with background dimming
	return base + "\n" + overlay
}

// appendContent appends content below base
func appendContent(base, append string) string {
	if append == "" {
		return base
	}
	return base + "\n" + append
}

// tickMsg is sent periodically for updates
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// renderMsg requests a re-render
type renderMsg struct{}

func RenderCmd() tea.Msg {
	return renderMsg{}
}
