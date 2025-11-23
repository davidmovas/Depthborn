package types

import "context"

// Alive represents entities that can die
type Alive interface {
	// IsAlive returns true if entity is alive
	IsAlive() bool

	// Kill marks entity as dead
	Kill(ctx context.Context, killerID string) error

	// Revive restores entity to life
	Revive(ctx context.Context, healthPercent float64) error
}

// Actionable represents entities that can perform actions
type Actionable interface {
	// CanAct returns true if entity can perform actions
	CanAct() bool
}

// Disposable represents entities that need cleanup
type Disposable interface {
	// Dispose releases resources
	Dispose() error
}

// Cloneable represents entities that can be copied
type Cloneable interface {
	// Clone creates deep copy
	Clone() any
}

// Validatable represents entities that can validate state
type Validatable interface {
	// Validate checks if state is valid
	Validate() error
}
