package status

import (
	"context"
	"fmt"
)

var _ Builder = (*EffectBuilder)(nil)

type EffectBuilder struct {
	config EffectConfig
}

type EffectConfig struct {
	EffectType    string
	Name          string
	Duration      int64
	InitialStacks int
	MaxStacks     int
	SourceID      string
	TargetID      string
	Metadata      map[string]any
	TickInterval  int64

	OnApplyFuncs  []func(ctx context.Context, targetID string) error
	OnTickFuncs   []func(ctx context.Context, targetID string, deltaMs int64) error
	OnRemoveFuncs []func(ctx context.Context, targetID string) error
	OnStackFuncs  []func(ctx context.Context, targetID string, newStacks int) error
}

func NewBuilder() *EffectBuilder {
	return &EffectBuilder{
		config: EffectConfig{
			MaxStacks:     1,
			InitialStacks: 1,
			Duration:      -1, // Infinite by default
			Metadata:      make(map[string]any),
		},
	}
}

func (b *EffectBuilder) WithType(effectType string) Builder {
	b.config.EffectType = effectType
	return b
}

func (b *EffectBuilder) WithName(name string) Builder {
	b.config.Name = name
	return b
}

func (b *EffectBuilder) WithDuration(ms int64) Builder {
	b.config.Duration = ms
	return b
}

func (b *EffectBuilder) WithStacks(initial, max int) Builder {
	b.config.InitialStacks = initial
	b.config.MaxStacks = max
	return b
}

func (b *EffectBuilder) WithSource(sourceID string) Builder {
	b.config.SourceID = sourceID
	return b
}

func (b *EffectBuilder) WithTarget(targetID string) Builder {
	b.config.TargetID = targetID
	return b
}

func (b *EffectBuilder) WithMetadata(key string, value interface{}) Builder {
	b.config.Metadata[key] = value
	return b
}

func (b *EffectBuilder) WithTickInterval(ms int64) Builder {
	b.config.TickInterval = ms
	return b
}

func (b *EffectBuilder) WithOnApply(fn func(ctx context.Context, targetID string) error) Builder {
	b.config.OnApplyFuncs = append(b.config.OnApplyFuncs, fn)
	return b
}

func (b *EffectBuilder) WithOnTick(fn func(ctx context.Context, targetID string, deltaMs int64) error) Builder {
	b.config.OnTickFuncs = append(b.config.OnTickFuncs, fn)
	return b
}

func (b *EffectBuilder) WithOnRemove(fn func(ctx context.Context, targetID string) error) Builder {
	b.config.OnRemoveFuncs = append(b.config.OnRemoveFuncs, fn)
	return b
}

func (b *EffectBuilder) WithOnStack(fn func(ctx context.Context, targetID string, newStacks int) error) Builder {
	b.config.OnStackFuncs = append(b.config.OnStackFuncs, fn)
	return b
}

func (b *EffectBuilder) Build() (Effect, error) {
	// Validation
	if b.config.EffectType == "" {
		return nil, fmt.Errorf("effect type is required")
	}

	if b.config.Name == "" {
		b.config.Name = b.config.EffectType
	}

	if b.config.TargetID == "" {
		return nil, fmt.Errorf("target ID is required")
	}

	if b.config.MaxStacks < 1 {
		return nil, fmt.Errorf("max stacks must be at least 1")
	}

	if b.config.InitialStacks < 1 || b.config.InitialStacks > b.config.MaxStacks {
		return nil, fmt.Errorf("initial stacks must be between 1 and max stacks")
	}

	return NewEffect(b.config), nil
}

// Reset resets builder to initial state
func (b *EffectBuilder) Reset() *EffectBuilder {
	b.config = EffectConfig{
		MaxStacks:     1,
		InitialStacks: 1,
		Duration:      -1,
		Metadata:      make(map[string]any),
	}
	return b
}
