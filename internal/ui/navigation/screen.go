package navigation

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// Screen represents a navigable UI screen with lifecycle support.
//
// Lifecycle order:
//   - OnInit (once when created)
//   - OnEnter (when becoming active)
//   - OnUpdateStart -> OnUpdate -> OnUpdateEnd (each frame)
//   - OnPause (when another screen overlays)
//   - OnResume (when returning from overlay)
//   - OnExit (when closed)
type Screen interface {
	// ID returns unique screen identifier.
	ID() string

	// --- Lifecycle Hooks ---

	// OnInit is called once when the screen is first created.
	OnInit()

	// OnEnter is called when the screen becomes active.
	OnEnter(params map[string]any)

	// OnUpdateStart is called before each update cycle.
	OnUpdateStart()

	// OnUpdate is called each frame.
	OnUpdate()

	// OnUpdateEnd is called after each update cycle.
	OnUpdateEnd()

	// OnPause is called when another screen overlays this one.
	OnPause()

	// OnResume is called when this screen becomes active again.
	OnResume()

	// OnExit is called when the screen is closed.
	OnExit()

	// --- Rendering ---

	// Render returns the component tree to render.
	Render(ctx *component.Context) component.Component

	// --- Navigation Control ---

	// CanClose returns whether this screen can be closed.
	CanClose() bool
}

// ScreenFactory creates screen instances.
type ScreenFactory func() Screen

// Verify interface compliance
var _ Screen = (*BaseScreen)(nil)

// BaseScreen provides default implementations for Screen interface.
// Embed this in custom screens to avoid implementing all methods.
type BaseScreen struct {
	id        string
	closeable bool
}

// NewBaseScreen creates a new base screen.
func NewBaseScreen(id string) *BaseScreen {
	return &BaseScreen{
		id:        id,
		closeable: true,
	}
}

// ID returns the screen identifier.
func (s *BaseScreen) ID() string {
	return s.id
}

// OnInit is called once when the screen is created.
func (s *BaseScreen) OnInit() {}

// OnEnter is called when the screen becomes active.
func (s *BaseScreen) OnEnter(params map[string]any) {}

// OnUpdateStart is called before each update cycle.
func (s *BaseScreen) OnUpdateStart() {}

// OnUpdate is called each frame.
func (s *BaseScreen) OnUpdate() {}

// OnUpdateEnd is called after each update cycle.
func (s *BaseScreen) OnUpdateEnd() {}

// OnPause is called when another screen overlays this one.
func (s *BaseScreen) OnPause() {}

// OnResume is called when this screen becomes active again.
func (s *BaseScreen) OnResume() {}

// OnExit is called when the screen is closed.
func (s *BaseScreen) OnExit() {}

// Render returns an empty component by default.
func (s *BaseScreen) Render(ctx *component.Context) component.Component {
	return component.Empty()
}

// CanClose returns whether the screen can be closed.
func (s *BaseScreen) CanClose() bool {
	return s.closeable
}

// SetCloseable sets whether the screen can be closed.
func (s *BaseScreen) SetCloseable(closeable bool) {
	s.closeable = closeable
}
