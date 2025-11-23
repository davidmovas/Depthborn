package skill

import (
	"context"
)

type BaseEffect struct {
	description string
	applyFunc   func(ctx context.Context, entityID string) error
	removeFunc  func(ctx context.Context, entityID string) error
}

func NewBaseEffect(description string,
	applyFunc func(ctx context.Context, entityID string) error,
	removeFunc func(ctx context.Context, entityID string) error,
) *BaseEffect {
	return &BaseEffect{
		description: description,
		applyFunc:   applyFunc,
		removeFunc:  removeFunc,
	}
}

func (e *BaseEffect) Apply(ctx context.Context, entityID string) error {
	if e.applyFunc != nil {
		return e.applyFunc(ctx, entityID)
	}
	return nil
}

func (e *BaseEffect) Remove(ctx context.Context, entityID string) error {
	if e.removeFunc != nil {
		return e.removeFunc(ctx, entityID)
	}
	return nil
}

func (e *BaseEffect) Description() string {
	return e.description
}
