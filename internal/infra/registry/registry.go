package registry

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/infra"
)

// Registry manages factories and lifecycle of registered types
type Registry interface {
	// Register adds factory for creating entities of specific type
	Register(entityType string, factory Factory) error

	// Unregister removes factory for entity type
	Unregister(entityType string) error

	// Create instantiates new entity of specified type
	Create(ctx context.Context, entityType string, params map[string]any) (infra.Identity, error)

	// Has checks if type is registered
	Has(entityType string) bool

	// Types returns all registered entity types
	Types() []string
}

// Factory creates instances of specific entity type
type Factory interface {
	// Create instantiates new entity with given parameters
	Create(ctx context.Context, params map[string]any) (infra.Identity, error)

	// TypeName returns entity type this factory creates
	TypeName() string
}

// TypedRegistry provides type-safe access to specific entity categories
type TypedRegistry[T infra.Identity] interface {
	// Register adds factory for type T
	Register(name string, factory TypedFactory[T]) error

	// Create instantiates new T with given name and params
	Create(ctx context.Context, name string, params map[string]any) (T, error)

	// Get retrieves cached instance by name
	Get(name string) (T, bool)

	// List returns all registered names
	List() []string
}

// TypedFactory creates instances of specific typed entity
type TypedFactory[T infra.Identity] interface {
	// Create instantiates new typed entity
	Create(ctx context.Context, params map[string]any) (T, error)
}
