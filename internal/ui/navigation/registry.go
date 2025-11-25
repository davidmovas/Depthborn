package navigation

import (
	"fmt"
	"sync"
)

// Registry stores screen factories for lazy initialization
// Screens are created only when first opened
type Registry struct {
	factories map[string]ScreenFactory
	mu        sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]ScreenFactory),
	}
}

// Register adds screen factory to registry
// screenID: unique identifier (e.g. "combat", "inventory")
// factory: function that creates new screen instance
func (r *Registry) Register(screenID string, factory ScreenFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories[screenID] = factory
}

// Create creates screen instance from registered factory
// Returns error if screen not registered
func (r *Registry) Create(screenID string) (Screen, error) {
	r.mu.RLock()
	factory, exists := r.factories[screenID]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("screen '%s' not registered", screenID)
	}

	screen := factory()

	// Call OnInit lifecycle hook
	screen.OnInit()

	return screen, nil
}

// Has checks if screen is registered
func (r *Registry) Has(screenID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.factories[screenID]
	return exists
}

// List returns all registered screen IDs
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.factories))
	for id := range r.factories {
		ids = append(ids, id)
	}
	return ids
}

// Unregister removes screen from registry
func (r *Registry) Unregister(screenID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.factories, screenID)
}

// Clear removes all registered screens
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories = make(map[string]ScreenFactory)
}
