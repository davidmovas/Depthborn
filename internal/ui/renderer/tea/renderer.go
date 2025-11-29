package tea

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
	"github.com/davidmovas/Depthborn/internal/ui/renderer"
)

var (
	ErrNotInitialized = errors.New("renderer not initialized")
)

// Verify interface compliance
var _ renderer.Renderer = (*Renderer)(nil)

// Renderer implements renderer.Renderer for BubbleTea.
type Renderer struct {
	config    renderer.Config
	navigator *navigation.Navigator
	program   *tea.Program
	model     *Model
}

// New creates a new BubbleTea renderer.
func New(config renderer.Config, navigator *navigation.Navigator) *Renderer {
	return &Renderer{
		config:    config,
		navigator: navigator,
	}
}

// Init implements renderer.Renderer.
func (r *Renderer) Init() error {
	// Get FPS from config or use default
	fps := r.config.FPS
	if fps <= 0 {
		fps = 60
	}

	// Create model
	r.model = NewModel(r.navigator, fps)

	// Build program options - always use alt screen for proper terminal handling
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
	}

	if r.config.Mouse {
		opts = append(opts, tea.WithMouseCellMotion())
	}

	// Create program
	r.program = tea.NewProgram(r.model, opts...)

	// Set program reference in model
	r.model.SetProgram(r.program)

	return nil
}

// Run implements renderer.Renderer.
func (r *Renderer) Run() error {
	if r.program == nil {
		return ErrNotInitialized
	}

	_, err := r.program.Run()
	return err
}

// Stop implements renderer.Renderer.
func (r *Renderer) Stop() error {
	if r.program != nil {
		r.program.Quit()
	}
	return nil
}

// Render implements renderer.Renderer.
func (r *Renderer) Render(comp component.Component) error {
	if r.model == nil {
		return ErrNotInitialized
	}

	r.model.RequestRender()
	return nil
}

// RequestRender implements renderer.Renderer.
func (r *Renderer) RequestRender() {
	if r.model != nil {
		r.model.RequestRender()
	}
}

// Size implements renderer.Renderer.
func (r *Renderer) Size() (width, height int) {
	if r.model != nil {
		return r.model.Size()
	}
	return 80, 24
}

// FPS returns current measured FPS.
func (r *Renderer) FPS() float64 {
	if r.model != nil {
		return r.model.FPS()
	}
	return 0
}

// TargetFPS returns target FPS.
func (r *Renderer) TargetFPS() int {
	if r.model != nil {
		return r.model.TargetFPS()
	}
	return 60
}

// LastKey returns the last pressed key (for debugging).
func (r *Renderer) LastKey() string {
	if r.model != nil {
		return r.model.LastKey()
	}
	return ""
}
