package navigation

import (
	"time"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// Screen represents a navigable UI screen
// Lifecycle:
//
//	OnInit (once)
//	OnEnter → OnUpdateStart → OnUpdate → OnUpdateEnd → ... → OnExit
//	OnPause/OnResume (when another screen opens on top)
type Screen interface {
	ID() string // Unique screen identifier

	// Lifecycle hooks

	OnInit()                       // Called once when screen is first created
	OnEnter(params map[string]any) // Called when screen becomes active
	OnUpdateStart()                // Called before each update
	OnUpdate(dt time.Duration)     // Called each frame/tick
	OnUpdateEnd()                  // Called after each update
	OnPause()                      // Called when screen is paused (another screen on top)
	OnResume()                     // Called when screen resumes from pause
	OnExit()                       // Called when screen is closed/destroyed

	// Rendering

	Render(ctx *component.Context) component.Component // Returns component tree to render

	// Input handling

	HandleInput(msg any) bool // Handle input message, return true if consumed

	CanClose() bool // Can this screen be closed? (false for mandatory screens)
}

// ScreenFactory creates screen instances
type ScreenFactory func() Screen

// BaseScreen provides default implementations for Screen interface
// Embed this in custom screens to avoid implementing all methods
type BaseScreen struct {
	id        string
	closeable bool
}

func NewBaseScreen(id string) *BaseScreen {
	return &BaseScreen{
		id:        id,
		closeable: true,
	}
}

func (s *BaseScreen) OnInit()                   {}
func (s *BaseScreen) OnEnter(_ map[string]any)  {}
func (s *BaseScreen) OnUpdateStart()            {}
func (s *BaseScreen) OnUpdate(dt time.Duration) {}
func (s *BaseScreen) OnUpdateEnd()              {}
func (s *BaseScreen) OnPause()                  {}
func (s *BaseScreen) OnResume()                 {}
func (s *BaseScreen) OnExit()                   {}

func (s *BaseScreen) Render() component.Component {
	return component.Empty()
}

func (s *BaseScreen) HandleInput(msg any) bool {
	return false // not handled
}

func (s *BaseScreen) ID() string {
	return s.id
}

func (s *BaseScreen) CanClose() bool {
	return s.closeable
}

func (s *BaseScreen) SetCloseable(closeable bool) {
	s.closeable = closeable
}
