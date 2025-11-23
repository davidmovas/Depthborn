package item

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/entity"
)

type BaseSocketEffect struct {
	description string
}

func NewBaseSocketEffect(description string) SocketEffect {
	return &BaseSocketEffect{
		description: description,
	}
}

func (bse *BaseSocketEffect) Apply(ctx context.Context, equipment Equipment) error {
	// TODO: Implement socket effect application
	return nil
}

func (bse *BaseSocketEffect) Remove(ctx context.Context, equipment Equipment) error {
	// TODO: Implement socket effect removal
	return nil
}

func (bse *BaseSocketEffect) Description() string {
	return bse.description
}

type BaseConsumableEffect struct {
	description string
	duration    int64
}

func NewBaseConsumableEffect(description string, duration int64) ConsumableEffect {
	return &BaseConsumableEffect{
		description: description,
		duration:    duration,
	}
}

func (bce *BaseConsumableEffect) Apply(ctx context.Context, target entity.Entity) error {
	return nil
}

func (bce *BaseConsumableEffect) Description() string {
	return bce.description
}

func (bce *BaseConsumableEffect) Duration() int64 {
	return bce.duration
}
