package renderer

import "github.com/davidmovas/Depthborn/internal/ui/component"

// Renderer is abstraction over UI rendering backend
// Allows switching between BubbleTea, Wails, or other frameworks
type Renderer interface {
	// Initialize renderer

	Init() error

	// Start renderer (blocking call)

	Run() error

	// Stop renderer

	Stop() error

	// Render component tree to output

	Render(comp component.Component) error

	// Request re-render

	RequestRender()
}

// Config holds renderer configuration
type Config struct {
	// Window title (for GUI renderers)
	Title string

	// Size hints (for terminal renderers)
	Width  int
	Height int

	// Additional config
	Extra map[string]any
}
