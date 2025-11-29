package tea

import (
	"strings"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
)

// tickMsg is sent on each frame tick.
type tickMsg time.Time

// Model implements tea.Model for BubbleTea integration.
type Model struct {
	navigator *navigation.Navigator
	context   *component.Context

	width  int
	height int

	targetFPS    int
	tickInterval time.Duration

	// FPS calculation
	frameCount    uint64
	lastFPSUpdate time.Time
	currentFPS    float64

	lastContent string
	needsRender bool
	program     *tea.Program

	// Debug: last key pressed
	lastKey string
}

// NewModel creates a new BubbleTea model.
func NewModel(nav *navigation.Navigator, fps int) *Model {
	if fps <= 0 {
		fps = 60
	}

	ctx := component.NewContext("root", nav)

	m := &Model{
		navigator:     nav,
		context:       ctx,
		width:         80,
		height:        24,
		targetFPS:     fps,
		tickInterval:  time.Second / time.Duration(fps),
		needsRender:   true,
		lastFPSUpdate: time.Now(),
	}

	// Set up render callback
	ctx.SetRenderCallback(func() {
		m.needsRender = true
	})

	return m
}

// SetProgram sets the BubbleTea program reference.
func (m *Model) SetProgram(p *tea.Program) {
	m.program = p
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return m.tick()
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.context.SetScreenSize(msg.Width, msg.Height)
		m.needsRender = true
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		// Mouse handling can be added later
		return m, nil

	case tickMsg:
		// Update FPS counter
		m.updateFPS()

		// Update navigator (game logic)
		m.navigator.Update()

		// Continue ticking
		return m, m.tick()
	}

	return m, nil
}

// handleKey processes keyboard input
func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	m.lastKey = key
	m.needsRender = true

	// QUIT: ctrl+c and ctrl+q ALWAYS quit
	switch key {
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}

	// ESC: close portal or go back
	if key == "esc" {
		if m.context.Portals().HasOpenPortals() {
			if top := m.context.Portals().GetTopPortal(); top != nil {
				m.context.Portals().Close(top.ID)
				return m, nil
			}
		}
		if m.navigator.CanGoBack() {
			_ = m.navigator.Back()
			return m, nil
		}
		return m, nil
	}

	// Route to focus manager
	m.routeKeyToFocus(key)

	return m, nil
}

// routeKeyToFocus sends key to the appropriate focus manager
func (m *Model) routeKeyToFocus(key string) {
	// If portal is open, route to portal's focus manager
	if m.context.Portals().HasOpenPortals() {
		if portalFocus := m.context.Portals().GetActiveFocus(); portalFocus != nil {
			portalFocus.HandleKey(key)
			return
		}
	}

	// Otherwise route to main focus manager
	m.context.Focus().HandleKey(key)
}

// View implements tea.Model.
func (m *Model) View() string {
	screen := m.navigator.CurrentScreen()
	if screen == nil {
		return m.centeredMessage("No screen active")
	}

	// Begin render frame
	m.context.BeginRender()

	// Render screen
	screenComp := screen.Render(m.context)
	if screenComp == nil {
		m.context.EndRender()
		return m.centeredMessage("Screen returned nil")
	}

	content := screenComp.Render(m.context)

	// End render frame
	m.context.EndRender()

	// Render portals
	portalContent := m.context.Portals().RenderPortals(m.context)
	if portalContent != "" {
		content = m.overlayContent(content, portalContent)
	}

	// Increment frame counter
	atomic.AddUint64(&m.frameCount, 1)

	m.lastContent = content
	m.needsRender = false

	return content
}

// --- Public Getters ---

// Size returns current screen dimensions.
func (m *Model) Size() (width, height int) {
	return m.width, m.height
}

// FPS returns the current measured FPS.
func (m *Model) FPS() float64 {
	return m.currentFPS
}

// TargetFPS returns the target FPS.
func (m *Model) TargetFPS() int {
	return m.targetFPS
}

// Context returns the root context.
func (m *Model) Context() *component.Context {
	return m.context
}

// LastKey returns the last pressed key (for debugging).
func (m *Model) LastKey() string {
	return m.lastKey
}

// RequestRender triggers a re-render.
func (m *Model) RequestRender() {
	m.needsRender = true
}

// --- Internal Methods ---

func (m *Model) updateFPS() {
	now := time.Now()
	elapsed := now.Sub(m.lastFPSUpdate)

	if elapsed >= time.Second {
		frames := atomic.SwapUint64(&m.frameCount, 0)
		m.currentFPS = float64(frames) / elapsed.Seconds()
		m.lastFPSUpdate = now
	}
}

func (m *Model) tick() tea.Cmd {
	return tea.Tick(m.tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) centeredMessage(msg string) string {
	s := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)
	return s.Render(msg)
}

func (m *Model) overlayContent(base, overlay string) string {
	if overlay == "" {
		return base
	}

	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	// Pad base to screen height
	for len(baseLines) < m.height {
		baseLines = append(baseLines, strings.Repeat(" ", m.width))
	}

	// Calculate centered position
	overlayWidth := maxLineWidth(overlayLines)
	overlayHeight := len(overlayLines)

	startX := (m.width - overlayWidth) / 2
	startY := (m.height - overlayHeight) / 2

	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}

	// Overlay merge
	for i, overlayLine := range overlayLines {
		y := startY + i
		if y >= len(baseLines) {
			break
		}

		baseLine := baseLines[y]

		// Pad base line if needed
		for len(baseLine) < startX+len(overlayLine) {
			baseLine += " "
		}

		var prefix, suffix string
		if startX > 0 && len(baseLine) >= startX {
			prefix = baseLine[:startX]
		}

		endPos := startX + lipgloss.Width(overlayLine)
		if endPos < len(baseLine) {
			suffix = baseLine[endPos:]
		}

		baseLines[y] = prefix + overlayLine + suffix
	}

	return strings.Join(baseLines, "\n")
}

func maxLineWidth(lines []string) int {
	max := 0
	for _, line := range lines {
		w := lipgloss.Width(line)
		if w > max {
			max = w
		}
	}
	return max
}
