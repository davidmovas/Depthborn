package renderer

import "github.com/davidmovas/Depthborn/internal/ui/component"

// Renderer defines the interface for UI renderers.
// Implementations can target different outputs (terminal, GUI, etc.)
type Renderer interface {
	// Init initializes the renderer.
	Init() error

	// Run starts the render loop.
	Run() error

	// Stop shuts down the renderer.
	Stop() error

	// Render renders a component immediately.
	Render(comp component.Component) error

	// RequestRender signals that a re-render is needed.
	RequestRender()

	// Size returns the current screen dimensions.
	Size() (width, height int)
}

// Config holds renderer configuration.
type Config struct {
	// Title displayed in terminal/window title bar
	Title string

	// Target FPS (default 60)
	FPS int

	// Initial screen size (0 = auto-detect)
	Width  int
	Height int

	// Enable alt screen mode (terminal only)
	AltScreen bool

	// Enable mouse support
	Mouse bool
}

// DefaultConfig returns default renderer configuration.
func DefaultConfig() Config {
	return Config{
		Title:     "Depthborn",
		FPS:       60,
		AltScreen: true,
		Mouse:     true,
	}
}
