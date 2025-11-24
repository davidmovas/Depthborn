package status

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/types"
)

// Effect represents temporary status on entity
type Effect interface {
	types.Identity
	types.Named

	// Duration returns remaining duration in milliseconds
	Duration() int64

	// SetDuration updates remaining duration
	SetDuration(ms int64)

	// Stacks returns current stack count
	Stacks() int

	// MaxStacks returns maximum stacks
	MaxStacks() int

	// AddStack increases stacks
	AddStack() bool

	// RemoveStack decreases stacks
	RemoveStack() bool

	// IsExpired returns true if duration ended
	IsExpired() bool

	// SourceID returns ID of entity that applied effect
	SourceID() string

	// TargetID returns ID of entity receiving effect
	TargetID() string

	// OnEvent processes effect event
	OnEvent(ctx context.Context, ev EffectEvent) error

	// OnApply invoked when effect applied
	OnApply(ctx context.Context, targetID string) error

	// OnTick invoked periodically
	OnTick(ctx context.Context, targetID string, deltaMs int64) error

	// OnRemove invoked when effect removed
	OnRemove(ctx context.Context, targetID string) error

	// OnStack invoked when stack count changed
	OnStack(ctx context.Context, targetID string, newStacks int) error

	// AddOnEvent subscribes to effect events
	AddOnEvent(eventType EffectEventType, fn func(ctx context.Context, eventData EffectEvent) error) (unsubscribe func())

	// CanStack checks if can stack with another
	CanStack(other Effect) bool

	// Metadata returns effect-specific data
	Metadata() map[string]any

	SetMetadata(key string, value any)
}

// Manager manages status effects on entity
type Manager interface {
	// Apply adds status effect
	Apply(ctx context.Context, effect Effect) error

	// Remove removes specific effect
	Remove(ctx context.Context, effectID string) error

	// RemoveByType removes all effects of type
	RemoveByType(ctx context.Context, effectType string) error

	// RemoveAll removes all effects
	RemoveAll(ctx context.Context) error

	// Has checks if has effect type
	Has(effectType string) bool

	// Get retrieves effect by ID
	Get(effectID string) (Effect, bool)

	// GetByType retrieves effects by type
	GetByType(effectType string) []Effect

	// GetAll returns all active effects
	GetAll() []Effect

	// Update processes all effects
	Update(ctx context.Context, deltaMs int64) error

	// Count returns total active effects
	Count() int

	// IsImmune checks immunity to effect type
	IsImmune(effectType string) bool

	// AddImmunity grants immunity
	AddImmunity(effectType string)

	// RemoveImmunity removes immunity
	RemoveImmunity(effectType string)
}

// Builder creates status effects
type Builder interface {
	// WithType sets effect type
	WithType(effectType string) Builder

	// WithName sets display name
	WithName(name string) Builder

	// WithDuration sets duration
	WithDuration(ms int64) Builder

	// WithStacks sets stacks
	WithStacks(initial, max int) Builder

	// WithSource sets source entity ID
	WithSource(sourceID string) Builder

	// WithTarget sets target entity ID
	WithTarget(targetID string) Builder

	// WithMetadata adds custom data
	WithMetadata(key string, value any) Builder

	// WithTickInterval sets tick rate
	WithTickInterval(ms int64) Builder

	// WithOnApply sets callback for OnApply event
	WithOnApply(fn func(ctx context.Context, targetID string) error) Builder

	// WithOnTick sets callback for OnTick event
	WithOnTick(fn func(ctx context.Context, targetID string, deltaMs int64) error) Builder

	// WithOnRemove sets callback for OnRemove event
	WithOnRemove(fn func(ctx context.Context, targetID string) error) Builder

	// WithOnStack sets callback for OnStack event
	WithOnStack(fn func(ctx context.Context, targetID string, newStacks int) error) Builder

	// Build creates effect
	Build() (Effect, error)
}

type EffectEventType int

const (
	EventApply EffectEventType = iota
	EventTick
	EventRemove
	EventStack
)

type EffectEvent struct {
	Effect   Effect
	TargetID string
	DeltaMs  int64
	NewStack int
	Type     EffectEventType
}

// Category groups effect types
type Category string

const (
	CategoryBuff    Category = "buff"
	CategoryDebuff  Category = "debuff"
	CategoryDamage  Category = "damage"
	CategoryHealing Category = "healing"
	CategoryControl Category = "control"
	CategoryUtility Category = "utility"
	CategoryAura    Category = "aura"
)

// Registry manages effect definitions
type Registry interface {
	// Register adds effect type
	Register(effectType string, factory Factory) error

	// Create instantiates effect
	Create(effectType string) (Effect, error)

	// Has checks if type registered
	Has(effectType string) bool

	// GetCategory returns effect category
	GetCategory(effectType string) Category
}

// Factory creates effect instances
type Factory interface {
	// Create instantiates effect
	Create() Effect

	// Type returns effect type
	Type() string

	// Category returns effect category
	Category() Category
}
