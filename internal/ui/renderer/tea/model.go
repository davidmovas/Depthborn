package tea

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
	"github.com/davidmovas/Depthborn/internal/ui/renderer"
	"github.com/davidmovas/Depthborn/internal/ui/renderer/vdom"
)

// Model implements tea.Model with VDOM optimization
type Model struct {
	navigator   *navigation.Navigator
	rootContext *component.Context
	vdom        *vdom.VDOM

	// FPS control
	targetFPS     int
	frameDuration time.Duration
	lastFrameTime time.Time

	// Cached output
	lastOutput string
	width      int
	height     int

	// BubbleTea program reference
	program *tea.Program

	// Debouncer for state changes
	debouncer *renderer.Debouncer
}

// NewModel creates new model
func NewModel(nav *navigation.Navigator) *Model {
	m := &Model{
		navigator:     nav,
		vdom:          vdom.NewVDOM(),
		targetFPS:     60, // Default 60 FPS
		frameDuration: time.Second / 60,
		lastFrameTime: time.Now(),
	}

	// Create root context
	m.rootContext = component.NewContext("root", nav)
	m.rootContext.SetStateChangeCallback(func(componentID string) {
		if m.program != nil {
			m.program.Send(stateChangeMsg{componentID: componentID})
		}
	})

	// Create debouncer for state changes (16ms = 60fps)
	m.debouncer = renderer.NewDebouncer(16*time.Millisecond, func() {
		m.RequestRender()
	})

	return m
}

// SetProgram sets program reference
func (m *Model) SetProgram(p *tea.Program) {
	m.program = p
}

// SetTargetFPS changes target framerate
func (m *Model) SetTargetFPS(fps int) {
	m.targetFPS = fps
	m.frameDuration = time.Second / time.Duration(fps)
}

// Init implements tea.Model.Init
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.tick(),
		tea.EnterAltScreen,
	)
}

// Update implements tea.Model.Update
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
		m.height = message.Height
		m.rootContext.SetScreenSize(message.Width, message.Height)
		return m, nil

	case tickMsg:
		// Calculate delta time
		now := time.Now()
		dt := now.Sub(m.lastFrameTime)
		m.lastFrameTime = now

		// Update navigator (game logic)
		m.navigator.Update(dt)

		// Schedule next tick
		return m, m.tick()

	case stateChangeMsg:
		// Component state changed - debounced render
		m.debouncer.Call()
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(message)

	case tea.MouseMsg:
		return m.handleMouse(message)
	}

	return m, nil
}

// View implements tea.Model.View
func (m *Model) View() string {
	currentScreen := m.navigator.Current()
	if currentScreen == nil {
		return "No screen loaded\n"
	}

	// Reset hook index before render
	m.rootContext.ResetHookIndex()

	// Get active portal focus context (if modal is open)
	activeFocusCtx := m.rootContext.PortalContext().GetActiveFocusContext()
	if activeFocusCtx != nil {
		// Modal is open - use isolated focus context
		m.rootContext.SetFocusContext(activeFocusCtx)
	} else {
		// No modal - use base focus context
		// Don't recreate context, just clear focusables for re-registration
		m.rootContext.FocusContext().ClearFocusables()
	}

	// Render current screen
	screenComponent := currentScreen.Render(m.rootContext)

	// Build virtual tree
	newTree := m.buildVTree(screenComponent, m.rootContext)

	// Reconcile with previous tree
	patches := m.vdom.Reconcile(newTree)

	// Apply patches to cached output
	if len(patches) > 0 {
		m.lastOutput = vdom.ApplyPatches(m.lastOutput, patches)
	}

	// Mark context as clean after render
	m.rootContext.ClearDirty()

	return m.lastOutput
}

// buildVTree builds virtual tree from component
func (m *Model) buildVTree(comp component.Component, ctx *component.Context) *vdom.VNode {
	// Render component
	content := comp.Render(ctx)

	// Create node
	node := vdom.BuildNode(
		ctx.ComponentID(),
		"component",
		content,
		nil, // Children would be tracked here for complex trees
		nil,
	)

	return node
}

// handleKeyPress handles keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check for quit
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	// Get active focus context (portal or base)
	focusCtx := m.rootContext.PortalContext().GetActiveFocusContext()
	if focusCtx == nil {
		focusCtx = m.rootContext.FocusContext()
	}

	// Handle focus navigation
	handled := focusCtx.HandleKey(msg.String())
	if handled {
		return m, nil
	}

	return m, nil
}

// handleMouse handles mouse input
func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// TODO: Implement mouse handling
	return m, nil
}

// tick creates command for next frame
func (m *Model) tick() tea.Cmd {
	return tea.Tick(m.frameDuration, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

// RequestRender requests immediate re-render
func (m *Model) RequestRender() {
	if m.program != nil {
		m.program.Send(tickMsg{})
	}
}

// tickMsg triggers frame update
type tickMsg struct{}

// stateChangeMsg signals component state change
type stateChangeMsg struct {
	componentID string
}
