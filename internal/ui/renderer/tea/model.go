package tea

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
)

// Model implements tea.Model interface
// Bridges our component system with BubbleTea
type Model struct {
	navigator *navigation.Navigator
	ctx       *component.Context
	lastTick  time.Time
	width     int
	height    int
}

// NewModel creates new BubbleTea model
func NewModel(navigator *navigation.Navigator) *Model {
	return &Model{
		navigator: navigator,
		ctx:       component.NewContext("root"),
		lastTick:  time.Now(),
	}
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
		// Global hotkeys
		switch message.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit

		case "esc":
			// Try to go back
			if err := m.navigator.Back(); err != nil {
				// If can't go back and only one screen, quit
				if m.navigator.StackSize() <= 1 {
					return m, tea.Quit
				}
			}
			return m, nil
		}

		// Forward input to navigator/screen
		m.navigator.HandleInput(message)
		return m, nil

	case tea.WindowSizeMsg:
		// Store window size
		m.width = message.Width
		m.height = message.Height
		return m, nil

	case tickMsg:
		// Update current screen
		now := time.Now()
		dt := now.Sub(m.lastTick)
		m.lastTick = now

		m.navigator.Update(dt)

		return m, tickCmd()

	case renderMsg:
		// Re-render requested
		return m, nil
	}

	return m, nil
}

// View implements tea.Model.View
func (m *Model) View() string {
	// Get current screen
	screen := m.navigator.Current()
	if screen == nil {
		return "No active screen"
	}

	// Reset hook index before render
	m.ctx.ResetHookIndex()

	// Render screen
	comp := screen.Render(m.ctx)
	if comp == nil {
		return ""
	}

	return comp.Render(m.ctx)
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
