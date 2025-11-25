package navigation

import "sync"

// Stack manages screen hierarchy (for back navigation)
type Stack struct {
	screens []Screen
	mu      sync.RWMutex
}

func NewStack() *Stack {
	return &Stack{
		screens: make([]Screen, 0),
	}
}

// Push adds screen to top of stack
func (s *Stack) Push(screen Screen) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Pause current screen if exists
	if len(s.screens) > 0 {
		current := s.screens[len(s.screens)-1]
		current.OnPause()
	}

	s.screens = append(s.screens, screen)
}

// Pop removes and returns top screen
// Returns nil if stack is empty
func (s *Stack) Pop() Screen {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.screens) == 0 {
		return nil
	}

	// Get top screen
	top := s.screens[len(s.screens)-1]

	// Call OnExit
	top.OnExit()

	// Remove from stack
	s.screens = s.screens[:len(s.screens)-1]

	// Resume previous screen if exists
	if len(s.screens) > 0 {
		prev := s.screens[len(s.screens)-1]
		prev.OnResume()
	}

	return top
}

// Peek returns top screen without removing it
// Returns nil if stack is empty
func (s *Stack) Peek() Screen {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.screens) == 0 {
		return nil
	}

	return s.screens[len(s.screens)-1]
}

// Replace replaces top screen with new screen
// Returns old screen
func (s *Stack) Replace(screen Screen) Screen {
	s.mu.Lock()
	defer s.mu.Unlock()

	var old Screen

	if len(s.screens) > 0 {
		old = s.screens[len(s.screens)-1]
		old.OnExit()
		s.screens[len(s.screens)-1] = screen
	} else {
		s.screens = append(s.screens, screen)
	}

	return old
}

// Clear removes all screens from stack
// Calls OnExit on all screens
func (s *Stack) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Call OnExit on all screens (reverse order)
	for i := len(s.screens) - 1; i >= 0; i-- {
		s.screens[i].OnExit()
	}

	s.screens = make([]Screen, 0)
}

// Size returns number of screens in stack
func (s *Stack) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.screens)
}

// IsEmpty returns true if stack is empty
func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

// All returns all screens in stack (bottom to top)
func (s *Stack) All() []Screen {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return copy to prevent external modification
	result := make([]Screen, len(s.screens))
	copy(result, s.screens)
	return result
}
