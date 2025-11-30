package item

import (
	"context"
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/entity"
)

// --- Socket Effects ---

var _ SocketEffect = (*BaseSocketEffect)(nil)

// BaseSocketEffect provides a basic implementation of SocketEffect
type BaseSocketEffect struct {
	id           string
	description  string
	onSocketFn   func(ctx context.Context, equipment Equipment, entity entity.Entity) error
	onUnsocketFn func(ctx context.Context, equipment Equipment, entity entity.Entity) error
}

// SocketEffectConfig holds configuration for creating a BaseSocketEffect
type SocketEffectConfig struct {
	ID           string
	Description  string
	OnSocketFn   func(ctx context.Context, equipment Equipment, entity entity.Entity) error
	OnUnsocketFn func(ctx context.Context, equipment Equipment, entity entity.Entity) error
}

// NewBaseSocketEffect creates a new socket effect with just a description
func NewBaseSocketEffect(description string) SocketEffect {
	return &BaseSocketEffect{
		description: description,
	}
}

// NewBaseSocketEffectWithConfig creates a new socket effect with full configuration
func NewBaseSocketEffectWithConfig(cfg SocketEffectConfig) *BaseSocketEffect {
	return &BaseSocketEffect{
		id:           cfg.ID,
		description:  cfg.Description,
		onSocketFn:   cfg.OnSocketFn,
		onUnsocketFn: cfg.OnUnsocketFn,
	}
}

func (bse *BaseSocketEffect) OnSocket(ctx context.Context, equipment Equipment, entity entity.Entity) error {
	if bse.onSocketFn != nil {
		return bse.onSocketFn(ctx, equipment, entity)
	}
	return nil
}

func (bse *BaseSocketEffect) OnUnsocket(ctx context.Context, equipment Equipment, entity entity.Entity) error {
	if bse.onUnsocketFn != nil {
		return bse.onUnsocketFn(ctx, equipment, entity)
	}
	return nil
}

func (bse *BaseSocketEffect) Description() string {
	return bse.description
}

// ID returns the effect identifier
func (bse *BaseSocketEffect) ID() string {
	return bse.id
}

// --- Consumable Effects ---

var _ ConsumableEffect = (*BaseConsumableEffect)(nil)

// BaseConsumableEffect provides a basic implementation of ConsumableEffect
type BaseConsumableEffect struct {
	id          string
	description string
	duration    int64
	applyFn     func(ctx context.Context, target entity.Entity) error
}

// ConsumableEffectConfig holds configuration for creating a BaseConsumableEffect
type ConsumableEffectConfig struct {
	ID          string
	Description string
	Duration    int64
	ApplyFn     func(ctx context.Context, target entity.Entity) error
}

// NewBaseConsumableEffect creates a new consumable effect
func NewBaseConsumableEffect(description string, duration int64) ConsumableEffect {
	return &BaseConsumableEffect{
		description: description,
		duration:    duration,
	}
}

// NewBaseConsumableEffectWithConfig creates a new consumable effect with full configuration
func NewBaseConsumableEffectWithConfig(cfg ConsumableEffectConfig) *BaseConsumableEffect {
	return &BaseConsumableEffect{
		id:          cfg.ID,
		description: cfg.Description,
		duration:    cfg.Duration,
		applyFn:     cfg.ApplyFn,
	}
}

func (bce *BaseConsumableEffect) Apply(ctx context.Context, target entity.Entity) error {
	if bce.applyFn != nil {
		return bce.applyFn(ctx, target)
	}
	return nil
}

func (bce *BaseConsumableEffect) Description() string {
	return bce.description
}

func (bce *BaseConsumableEffect) Duration() int64 {
	return bce.duration
}

// ID returns the effect identifier
func (bce *BaseConsumableEffect) ID() string {
	return bce.id
}

// --- Effect Registry ---

// EffectRegistry stores and retrieves effects by ID (thread-safe)
type EffectRegistry struct {
	mu                sync.RWMutex
	socketEffects     map[string]SocketEffect
	consumableEffects map[string]ConsumableEffect
}

// NewEffectRegistry creates a new effect registry
func NewEffectRegistry() *EffectRegistry {
	return &EffectRegistry{
		socketEffects:     make(map[string]SocketEffect),
		consumableEffects: make(map[string]ConsumableEffect),
	}
}

// RegisterSocketEffect registers a socket effect
func (r *EffectRegistry) RegisterSocketEffect(id string, effect SocketEffect) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.socketEffects[id] = effect
}

// GetSocketEffect retrieves a socket effect by ID
func (r *EffectRegistry) GetSocketEffect(id string) (SocketEffect, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	effect, ok := r.socketEffects[id]
	return effect, ok
}

// RegisterConsumableEffect registers a consumable effect
func (r *EffectRegistry) RegisterConsumableEffect(id string, effect ConsumableEffect) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.consumableEffects[id] = effect
}

// GetConsumableEffect retrieves a consumable effect by ID
func (r *EffectRegistry) GetConsumableEffect(id string) (ConsumableEffect, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	effect, ok := r.consumableEffects[id]
	return effect, ok
}

// UnregisterSocketEffect removes a socket effect
func (r *EffectRegistry) UnregisterSocketEffect(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.socketEffects, id)
}

// UnregisterConsumableEffect removes a consumable effect
func (r *EffectRegistry) UnregisterConsumableEffect(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.consumableEffects, id)
}
