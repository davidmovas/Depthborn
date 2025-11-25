package tea

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
	"github.com/davidmovas/Depthborn/internal/ui/renderer"
)

var (
	ErrNotInitialized = &RendererError{msg: "renderer not initialized"}
)

var _ renderer.Renderer = (*Renderer)(nil)

// Renderer implements renderer.Renderer for BubbleTea
type Renderer struct {
	config    renderer.Config
	navigator *navigation.Navigator
	program   *tea.Program
	model     *Model
}

// New creates new BubbleTea renderer
func New(config renderer.Config, navigator *navigation.Navigator) *Renderer {
	return &Renderer{
		config:    config,
		navigator: navigator,
	}
}

// Init implements renderer.Renderer.Init
func (r *Renderer) Init() error {
	// Create model
	r.model = NewModel(r.navigator)

	// Create program
	r.program = tea.NewProgram(
		r.model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support (optional)
	)

	return nil
}

// Run implements renderer.Renderer.Run (blocking)
func (r *Renderer) Run() error {
	if r.program == nil {
		return ErrNotInitialized
	}

	_, err := r.program.Run()
	return err
}

// Stop implements renderer.Renderer.Stop
func (r *Renderer) Stop() error {
	if r.program != nil {
		r.program.Quit()
	}
	return nil
}

// Render implements renderer.Renderer.Render
// For BubbleTea, rendering is automatic via View()
// This method can trigger a re-render
func (r *Renderer) Render(comp component.Component) error {
	if r.program == nil {
		return ErrNotInitialized
	}

	// Send render message to trigger update
	r.program.Send(RenderCmd())

	return nil
}

// RequestRender implements renderer.Renderer.RequestRender
func (r *Renderer) RequestRender() {
	if r.program != nil {
		r.program.Send(RenderCmd())
	}
}

type RendererError struct {
	msg string
}

func (e *RendererError) Error() string {
	return e.msg
}
