package status

import (
	"fmt"
	"sync"
)

var _ Registry = (*BaseRegistry)(nil)

type BaseRegistry struct {
	factories  map[string]Factory
	categories map[string]Category

	mu sync.RWMutex
}

func NewRegistry() *BaseRegistry {
	return &BaseRegistry{
		factories:  make(map[string]Factory),
		categories: make(map[string]Category),
	}
}

func (r *BaseRegistry) Register(effectType string, factory Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if effectType == "" {
		return fmt.Errorf("effect type cannot be empty")
	}

	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	if _, exists := r.factories[effectType]; exists {
		return fmt.Errorf("effect type already registered: %s", effectType)
	}

	r.factories[effectType] = factory
	r.categories[effectType] = factory.Category()

	return nil
}

func (r *BaseRegistry) Create(effectType string) (Effect, error) {
	r.mu.RLock()
	factory, exists := r.factories[effectType]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("effect type not registered: %s", effectType)
	}

	return factory.Create(), nil
}

func (r *BaseRegistry) Has(effectType string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.factories[effectType]
	return exists
}

func (r *BaseRegistry) GetCategory(effectType string) Category {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if category, exists := r.categories[effectType]; exists {
		return category
	}
	return CategoryUtility // Default category
}

// GetAllTypes returns all registered effect types
func (r *BaseRegistry) GetAllTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.factories))
	for effectType := range r.factories {
		types = append(types, effectType)
	}
	return types
}

// Unregister removes effect type from registry
func (r *BaseRegistry) Unregister(effectType string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[effectType]; !exists {
		return fmt.Errorf("effect type not registered: %s", effectType)
	}

	delete(r.factories, effectType)
	delete(r.categories, effectType)

	return nil
}

// BaseFactory provides simple factory implementation
type BaseFactory struct {
	effectType string
	category   Category
	createFunc func() Effect
}

var _ Factory = (*BaseFactory)(nil)

func NewFactory(effectType string, category Category, createFunc func() Effect) *BaseFactory {
	return &BaseFactory{
		effectType: effectType,
		category:   category,
		createFunc: createFunc,
	}
}

func (f *BaseFactory) Create() Effect {
	return f.createFunc()
}

func (f *BaseFactory) Type() string {
	return f.effectType
}

func (f *BaseFactory) Category() Category {
	return f.category
}
