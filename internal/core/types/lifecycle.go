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
