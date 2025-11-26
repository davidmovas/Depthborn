package tea

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
	"github.com/davidmovas/Depthborn/internal/ui/renderer"
)

// Model bridges component system with BubbleTea
type Model struct {
	navigator       *navigation.Navigator
	lastTick        time.Time
	width           int
	height          int
	renderDebounce  *renderer.Debouncer
	program         *tea.Program
	screenContexts  map[string]*component.Context // Optimized: cache per screen
	currentFocusCtx *component.FocusContext
}

// NewModel creates new BubbleTea model
func NewModel(navigator *navigation.Navigator) *Model {
	m := &Model{
		navigator:      navigator,
		lastTick:       time.Now(),
		screenContexts: make(map[string]*component.Context, 20),
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
	// Start ticker for updates
	return tea.Batch(
		tickCmd(),
		tea.EnterAltScreen, // Use alternate screen buffer
	)
}

// Update implements tea.Model.Update
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {

	case tea.KeyMsg:
		key := message.String()

		// Global hotkeys (always work)
		switch key {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit

		case "esc":
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

// View implements tea.Model.View
func (m *Model) View() string {
	screen := m.navigator.Current()
	if screen == nil {
		return "No active screen"
	}

	// Ensure focus context exists
	m.ensureFocusContext(screen)

	// Clear and re-register components
	m.currentFocusCtx.ClearFocusables()

	// Create context with focus context already set
	ctx := m.createContextForScreen(screen)

	// Render screen (components will register themselves)
	comp := screen.Render(ctx)
	if comp == nil {
		return ""
	}

	// Render to string (triggers registration)
	rendered := comp.Render(ctx)

	return rendered
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
