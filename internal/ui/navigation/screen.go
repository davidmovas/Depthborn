package navigation

import (
	"fmt"
	"time"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// Screen represents a navigable UI screen with full lifecycle support
// Lifecycle order:
//
//	OnInit (once) → OnEnter → OnUpdateStart → OnUpdate → OnUpdateEnd → ... → OnExit
//	OnPause/OnResume (when another screen overlays this one)
type Screen interface {
	// ID returns unique screen identifier
	ID() string

	// Lifecycle hooks

	// OnInit is called once when screen is first created
	OnInit()

	// OnEnter is called when screen becomes active
	// params: optional parameters passed from navigator
	OnEnter(params map[string]any)

	// OnUpdateStart is called before each update cycle
	OnUpdateStart()

	// OnUpdate is called each frame/tick
	// dt: delta time since last update
	OnUpdate(dt time.Duration)

	// OnUpdateEnd is called after each update cycle
	OnUpdateEnd()

	// OnPause is called when screen is paused (another screen overlays it)
	OnPause()

	// OnResume is called when screen resumes from pause
	OnResume()

	// OnExit is called when screen is closed/destroyed
	OnExit()

	// Rendering

	// Render returns component tree to render
	// ctx: component context for hooks and state management
	Render(ctx *component.Context) component.Component

	// Input handling

	// HandleInput handles input message
	// ctx: component context for accessing focus and state
	// msg: input message (e.g. tea.KeyMsg)
	// Returns true if input was consumed
	HandleInput(ctx *component.Context, msg any) bool

	// CanClose returns whether this screen can be closed
	// Set to false for mandatory screens (e.g. loading screen)
	CanClose() bool

	// SetHook sets hook in screen's storage (used by Context)
	SetHook(key string, state *component.HookState)

	// GetHook gets hook from screen's storage (used by Context)
	GetHook(key string) (*component.HookState, bool)

	// HooksStorage returns hooks storage for this screen
	// Used internally by Context to persist hooks between renders
	HooksStorage() map[string]*component.HookState
}

// ScreenFactory creates screen instances
// Used for lazy initialization via Registry
type ScreenFactory func() Screen

var _ Screen = (*BaseScreen)(nil)

// BaseScreen provides default implementations for Screen interface
// Embed this in custom screens to avoid implementing all methods
type BaseScreen struct {
	id        string
	closeable bool

	// hooksStorage stores hooks state between renders (per-screen persistence)
	hooksStorage map[string]*component.HookState
}

// NewBaseScreen creates a new base screen
func NewBaseScreen(id string) *BaseScreen {
	return &BaseScreen{
		id:           id,
		closeable:    true,
		hooksStorage: make(map[string]*component.HookState),
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

// Render returns empty component by default
func (s *BaseScreen) Render(ctx *component.Context) component.Component {
	return component.Empty()
}

// HandleInput returns false by default (input not handled)
func (s *BaseScreen) HandleInput(ctx *component.Context, msg any) bool {
	return false
}

// ID returns screen identifier
func (s *BaseScreen) ID() string {
	return s.id
}

// CanClose returns whether screen can be closed
func (s *BaseScreen) CanClose() bool {
	return s.closeable
}

// SetCloseable sets whether screen can be closed
func (s *BaseScreen) SetCloseable(closeable bool) {
	s.closeable = closeable
}

// GetHook gets hook from screen's storage (used by Context)
func (s *BaseScreen) GetHook(key string) (*component.HookState, bool) {
	state, ok := s.hooksStorage[key]

	// DEBUG: проверяем что в хранилище
	if !ok {
		fmt.Printf("[DEBUG] BaseScreen.GetHook: %s NOT FOUND. Storage size: %d\n", key, len(s.hooksStorage))
	} else {
		fmt.Printf("[DEBUG] BaseScreen.GetHook: %s FOUND = %v\n", key, state.Value)
	}

	return state, ok
}

// SetHook sets hook in screen's storage (used by Context)
func (s *BaseScreen) SetHook(key string, state *component.HookState) {
	fmt.Printf("[DEBUG] BaseScreen.SetHook: %s = %v\n", key, state.Value)
	s.hooksStorage[key] = state
}

// HooksStorage returns hooks storage for this screen
// Used internally by Context to persist hooks between renders
func (s *BaseScreen) HooksStorage() map[string]*component.HookState {
	return s.hooksStorage
}
