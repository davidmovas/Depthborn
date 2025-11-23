// item/consumable.go
package item

import (
	"context"
	"fmt"
	"time"

	"github.com/davidmovas/Depthborn/internal/core/entity"
)

var _ Consumable = (*BaseConsumable)(nil)

type BaseConsumable struct {
	*BaseItem
	cooldown    int64
	maxCooldown int64
	effect      ConsumableEffect
	lastUsed    int64
}

func NewBaseConsumable(id string, name string) *BaseConsumable {
	return &BaseConsumable{
		BaseItem:    NewBaseItem(id, TypeConsumable, name),
		maxCooldown: 0,
		cooldown:    0,
		effect:      nil,
		lastUsed:    0,
	}
}

func (bc *BaseConsumable) Use(ctx context.Context, user entity.Entity) error {
	if !bc.CanUse(user) {
		return fmt.Errorf("cannot use consumable")
	}

	if bc.effect != nil {
		if err := bc.effect.Apply(ctx, user); err != nil {
			return err
		}
	}

	bc.lastUsed = time.Now().UnixMilli()
	bc.cooldown = bc.maxCooldown

	if bc.StackSize() > 1 {
		bc.RemoveStack(1)
	}

	return nil
}

func (bc *BaseConsumable) CanUse(user entity.Entity) bool {
	return bc.Cooldown() <= 0
}

func (bc *BaseConsumable) Cooldown() int64 {
	now := time.Now().UnixMilli()
	elapsed := now - bc.lastUsed
	if elapsed >= bc.maxCooldown {
		return 0
	}
	return bc.maxCooldown - elapsed
}

func (bc *BaseConsumable) MaxCooldown() int64 {
	return bc.maxCooldown
}

func (bc *BaseConsumable) Effect() ConsumableEffect {
	return bc.effect
}
