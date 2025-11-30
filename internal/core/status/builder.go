package status

import (
	"context"
	"fmt"
	"sync"
)

var _ Builder = (*EffectBuilder)(nil)

type EffectBuilder struct {
	config EffectConfig

	mu     sync.RWMutex
	events map[EffectEventType][]func(ctx context.Context, ev EffectEvent) error
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
}

func NewBuilder() *EffectBuilder {
	return &EffectBuilder{
		config: EffectConfig{
			MaxStacks:     1,
			InitialStacks: 1,
			Duration:      -1,
			Metadata:      make(map[string]any),
		},

		events: make(map[EffectEventType][]func(ctx context.Context, ev EffectEvent) error),
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

func (b *EffectBuilder) WithMetadata(key string, value any) Builder {
	b.config.Metadata[key] = value
	return b
}

func (b *EffectBuilder) WithTickInterval(ms int64) Builder {
	b.config.TickInterval = ms
	return b
}

func (b *EffectBuilder) WithOnEvent(eventType EffectEventType,
	fn func(ctx context.Context, ev EffectEvent) error,
) Builder {

	if b.events == nil {
		b.events = make(map[EffectEventType][]func(context.Context, EffectEvent) error)
	}

	b.events[eventType] = append(b.events[eventType], fn)
	return b
}

func (b *EffectBuilder) WithOnApply(fn func(context.Context, string) error) Builder {
	return b.WithOnEvent(EventApply, func(ctx context.Context, ev EffectEvent) error {
		return fn(ctx, ev.TargetID)
	})
}

func (b *EffectBuilder) WithOnTick(fn func(context.Context, string, int64) error) Builder {
	return b.WithOnEvent(EventTick, func(ctx context.Context, ev EffectEvent) error {
		return fn(ctx, ev.TargetID, ev.DeltaMs)
	})
}

func (b *EffectBuilder) WithOnRemove(fn func(context.Context, string) error) Builder {
	return b.WithOnEvent(EventRemove, func(ctx context.Context, ev EffectEvent) error {
		return fn(ctx, ev.TargetID)
	})
}

func (b *EffectBuilder) WithOnStack(fn func(context.Context, string, int) error) Builder {
	return b.WithOnEvent(EventStack, func(ctx context.Context, ev EffectEvent) error {
		return fn(ctx, ev.TargetID, ev.NewStack)
	})
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

	// Create effect
	eff := NewEffect(b.config)

	// Register callbacks inside effect
	for t, list := range b.events {
		for _, fn := range list {
			eff.AddOnEvent(t, fn)
		}
	}

	return eff, nil
}

// Reset clears config + all callbacks
func (b *EffectBuilder) Reset() *EffectBuilder {
	b.config = EffectConfig{
		MaxStacks:     1,
		InitialStacks: 1,
		Duration:      -1,
		Metadata:      make(map[string]any),
	}

	b.mu.Lock()
	b.events = make(map[EffectEventType][]func(context.Context, EffectEvent) error)
	b.mu.Unlock()

	return b
}
