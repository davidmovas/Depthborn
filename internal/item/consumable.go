package item

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

var _ Consumable = (*BaseConsumable)(nil)

// BaseConsumable implements Consumable interface
type BaseConsumable struct {
	*BaseItem

	mu          sync.RWMutex
	cooldown    int64
	maxCooldown int64
	effect      ConsumableEffect
	effectID    string // For serialization - identifies the effect type
	lastUsed    int64
	charges     int // Number of uses before consumed (-1 for infinite until stack depletes)
	maxCharges  int
}

// ConsumableConfig holds configuration for creating a BaseConsumable
type ConsumableConfig struct {
	BaseItemConfig
	MaxCooldown int64
	Effect      ConsumableEffect
	EffectID    string
	Charges     int
}

// NewBaseConsumable creates a new consumable with minimal configuration
func NewBaseConsumable(id string, name string) *BaseConsumable {
	return NewBaseConsumableWithConfig(ConsumableConfig{
		BaseItemConfig: BaseItemConfig{
			ID:       id,
			Name:     name,
			ItemType: TypeConsumable,
		},
		Charges: 1,
	})
}

// NewBaseConsumableWithConfig creates a new consumable with full configuration
func NewBaseConsumableWithConfig(cfg ConsumableConfig) *BaseConsumable {
	// Ensure item type is consumable
	cfg.BaseItemConfig.ItemType = TypeConsumable

	bc := &BaseConsumable{
		BaseItem:    NewBaseItemWithConfig(cfg.BaseItemConfig),
		maxCooldown: cfg.MaxCooldown,
		cooldown:    0,
		effect:      cfg.Effect,
		effectID:    cfg.EffectID,
		lastUsed:    0,
		charges:     cfg.Charges,
		maxCharges:  cfg.Charges,
	}

	// Default to 1 charge if not specified
	if bc.charges <= 0 && bc.charges != -1 {
		bc.charges = 1
		bc.maxCharges = 1
	}

	return bc
}

// --- Consumable interface implementation ---

func (bc *BaseConsumable) Use(ctx context.Context, user entity.Entity) error {
	bc.mu.Lock()

	if !bc.canUseInternal(user) {
		bc.mu.Unlock()
		return fmt.Errorf("cannot use consumable: on cooldown or no charges")
	}

	effect := bc.effect
	bc.mu.Unlock()

	// Apply effect (outside lock to avoid holding lock during potentially long operation)
	if effect != nil {
		if err := effect.Apply(ctx, user); err != nil {
			return fmt.Errorf("failed to apply consumable effect: %w", err)
		}
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Update state
	bc.lastUsed = time.Now().UnixMilli()
	bc.cooldown = bc.maxCooldown

	// Consume charge
	if bc.charges > 0 {
		bc.charges--
	}

	// Remove from stack if no charges left
	if bc.charges == 0 {
		// Use internal version to avoid nested locking
		bc.BaseItem.RemoveStackInternal(1)
		// Reset charges for next item in stack
		if bc.BaseItem.StackSizeInternal() > 0 {
			bc.charges = bc.maxCharges
		}
	}

	bc.Touch()
	return nil
}

func (bc *BaseConsumable) CanUse(user entity.Entity) bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.canUseInternal(user)
}

// canUseInternal checks if consumable can be used (no lock)
func (bc *BaseConsumable) canUseInternal(user entity.Entity) bool {
	// Check cooldown
	if bc.cooldownInternal() > 0 {
		return false
	}

	// Check charges
	if bc.charges == 0 {
		return false
	}

	// Check stack (use internal version to avoid nested locking)
	if bc.BaseItem.StackSizeInternal() <= 0 {
		return false
	}

	return true
}

func (bc *BaseConsumable) Cooldown() int64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.cooldownInternal()
}

// cooldownInternal returns remaining cooldown (no lock)
func (bc *BaseConsumable) cooldownInternal() int64 {
	if bc.maxCooldown <= 0 {
		return 0
	}

	now := time.Now().UnixMilli()
	elapsed := now - bc.lastUsed
	if elapsed >= bc.maxCooldown {
		return 0
	}
	return bc.maxCooldown - elapsed
}

func (bc *BaseConsumable) MaxCooldown() int64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.maxCooldown
}

func (bc *BaseConsumable) SetMaxCooldown(cooldown int64) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	if cooldown < 0 {
		cooldown = 0
	}
	bc.maxCooldown = cooldown
	bc.Touch()
}

func (bc *BaseConsumable) Effect() ConsumableEffect {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.effect
}

func (bc *BaseConsumable) SetEffect(effect ConsumableEffect, effectID string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.effect = effect
	bc.effectID = effectID
	bc.Touch()
}

// --- Additional methods ---

// Charges returns current charges remaining
func (bc *BaseConsumable) Charges() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.charges
}

// MaxCharges returns maximum charges
func (bc *BaseConsumable) MaxCharges() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.maxCharges
}

// SetCharges sets current charges
func (bc *BaseConsumable) SetCharges(charges int) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	if charges < -1 {
		charges = -1
	}
	if charges > bc.maxCharges && bc.maxCharges > 0 {
		charges = bc.maxCharges
	}
	bc.charges = charges
	bc.Touch()
}

// EffectID returns the effect identifier for serialization
func (bc *BaseConsumable) EffectID() string {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.effectID
}

// ResetCooldown resets the cooldown to 0
func (bc *BaseConsumable) ResetCooldown() {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.lastUsed = 0
	bc.cooldown = 0
	bc.Touch()
}

// --- Cloneable interface ---

func (bc *BaseConsumable) Clone() any {
	if bc == nil {
		return nil
	}

	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if bc.BaseItem == nil {
		return nil
	}

	// Clone base item
	cloned := bc.BaseItem.Clone()
	if cloned == nil {
		return nil
	}
	baseClone := cloned.(*BaseItem)

	clone := &BaseConsumable{
		BaseItem:    baseClone,
		maxCooldown: bc.maxCooldown,
		cooldown:    0, // Reset cooldown for clone
		effect:      bc.effect,
		effectID:    bc.effectID,
		lastUsed:    0,             // Reset for clone
		charges:     bc.maxCharges, // Full charges for clone
		maxCharges:  bc.maxCharges,
	}

	return clone
}

// --- Serialization ---

// ConsumableState holds the complete serializable state of a BaseConsumable
type ConsumableState struct {
	State
	MaxCooldown int64  `msgpack:"max_cooldown"`
	LastUsed    int64  `msgpack:"last_used"`
	EffectID    string `msgpack:"effect_id"`
	Charges     int    `msgpack:"charges"`
	MaxCharges  int    `msgpack:"max_charges"`
}

func (bc *BaseConsumable) Marshal() ([]byte, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	// Get base item data
	baseData, err := bc.BaseItem.Marshal()
	if err != nil {
		return nil, err
	}

	// Decode to State
	var is State
	if err := persist.DefaultCodec().Decode(baseData, &is); err != nil {
		return nil, err
	}

	cs := ConsumableState{
		State:       is,
		MaxCooldown: bc.maxCooldown,
		LastUsed:    bc.lastUsed,
		EffectID:    bc.effectID,
		Charges:     bc.charges,
		MaxCharges:  bc.maxCharges,
	}

	return persist.DefaultCodec().Encode(cs)
}

func (bc *BaseConsumable) Unmarshal(data []byte) error {
	var cs ConsumableState
	if err := persist.DefaultCodec().Decode(data, &cs); err != nil {
		return fmt.Errorf("failed to decode consumable state: %w", err)
	}

	// Encode item state back to bytes
	itemData, err := persist.DefaultCodec().Encode(cs.State)
	if err != nil {
		return err
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Initialize base item if nil
	if bc.BaseItem == nil {
		bc.BaseItem = &BaseItem{}
	}

	// Restore base item
	if err := bc.BaseItem.Unmarshal(itemData); err != nil {
		return err
	}

	// Restore consumable-specific fields
	bc.maxCooldown = cs.MaxCooldown
	bc.lastUsed = cs.LastUsed
	bc.effectID = cs.EffectID
	bc.charges = cs.Charges
	bc.maxCharges = cs.MaxCharges

	// Note: effect must be restored separately using effectID
	// via an effect registry

	return nil
}

// --- Validatable interface ---

func (bc *BaseConsumable) Validate() error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	// Validate base item
	if err := bc.BaseItem.Validate(); err != nil {
		return err
	}

	if bc.maxCooldown < 0 {
		return fmt.Errorf("max cooldown cannot be negative")
	}

	if bc.maxCharges < -1 {
		return fmt.Errorf("max charges cannot be less than -1")
	}

	if bc.charges < -1 {
		return fmt.Errorf("charges cannot be less than -1")
	}

	if bc.maxCharges > 0 && bc.charges > bc.maxCharges {
		return fmt.Errorf("charges cannot exceed max charges")
	}

	return nil
}

// --- Serializable interface (required by infra.Persistent) ---

func (bc *BaseConsumable) SerializeState() (map[string]any, error) {
	data, err := bc.Marshal()
	if err != nil {
		return nil, err
	}

	var state map[string]any
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return nil, err
	}
	return state, nil
}

func (bc *BaseConsumable) DeserializeState(state map[string]any) error {
	data, err := persist.DefaultCodec().Encode(state)
	if err != nil {
		return err
	}
	return bc.Unmarshal(data)
}
