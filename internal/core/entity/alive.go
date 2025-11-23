package entity

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/types"
)

var _ types.Alive = (*BaseAlive)(nil)

type BaseAlive struct {
	entity Entity
	alive  bool
}

func NewBaseAlive(entity Entity) *BaseAlive {
	return &BaseAlive{
		entity: entity,
		alive:  true,
	}
}

func (ba *BaseAlive) IsAlive() bool {
	return ba.alive
}

func (ba *BaseAlive) Kill(ctx context.Context, killerID string) error {
	ba.alive = false
	// TODO: Trigger death callbacks
	return nil
}

func (ba *BaseAlive) Revive(ctx context.Context, healthPercent float64) error {
	ba.alive = true
	// TODO: Restore health and trigger revive logic
	return nil
}
