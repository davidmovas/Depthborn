package item

import (
	"context"
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

var _ Socketable = (*BaseSocketable)(nil)

// BaseSocketable implements Socketable interface for gems, runes, etc.
type BaseSocketable struct {
	*BaseItem

	mu         sync.RWMutex
	socketType SocketType
	effect     SocketEffect
	effectID   string // For serialization - identifies the effect type
	tier       int    // Power tier of the socketable (1-5 typically)
	modifiers  []attribute.Modifier
}

// SocketableConfig holds configuration for creating a BaseSocketable
type SocketableConfig struct {
	BaseItemConfig
	SocketType SocketType
	Effect     SocketEffect
	EffectID   string
	Tier       int
	Modifiers  []attribute.Modifier
}

// NewBaseSocketable creates a new socketable with minimal configuration
func NewBaseSocketable(id string, itemType Type, name string, socketType SocketType) *BaseSocketable {
	return NewBaseSocketableWithConfig(SocketableConfig{
		BaseItemConfig: BaseItemConfig{
			ID:       id,
			Name:     name,
			ItemType: itemType,
		},
		SocketType: socketType,
		Tier:       1,
	})
}

// NewBaseSocketableWithConfig creates a new socketable with full configuration
func NewBaseSocketableWithConfig(cfg SocketableConfig) *BaseSocketable {
	bs := &BaseSocketable{
		BaseItem:   NewBaseItemWithConfig(cfg.BaseItemConfig),
		socketType: cfg.SocketType,
		effect:     cfg.Effect,
		effectID:   cfg.EffectID,
		tier:       cfg.Tier,
		modifiers:  cfg.Modifiers,
	}

	if bs.tier < 1 {
		bs.tier = 1
	}
	if bs.tier > 5 {
		bs.tier = 5
	}
	if bs.modifiers == nil {
		bs.modifiers = make([]attribute.Modifier, 0)
	}

	return bs
}

// --- Socketable interface implementation ---

func (bs *BaseSocketable) SocketType() SocketType {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.socketType
}

func (bs *BaseSocketable) Effect() SocketEffect {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.effect
}

func (bs *BaseSocketable) SetEffect(effect SocketEffect, effectID string) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.effect = effect
	bs.effectID = effectID
	bs.Touch()
}

// --- Additional methods ---

// EffectID returns the effect identifier for serialization
func (bs *BaseSocketable) EffectID() string {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.effectID
}

// Tier returns the power tier
func (bs *BaseSocketable) Tier() int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.tier
}

// SetTier sets the power tier
func (bs *BaseSocketable) SetTier(tier int) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if tier < 1 {
		tier = 1
	}
	if tier > 5 {
		tier = 5
	}
	bs.tier = tier
	bs.Touch()
}

// Modifiers returns the attribute modifiers
func (bs *BaseSocketable) Modifiers() []attribute.Modifier {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	result := make([]attribute.Modifier, len(bs.modifiers))
	copy(result, bs.modifiers)
	return result
}

// SetModifiers sets the attribute modifiers
func (bs *BaseSocketable) SetModifiers(modifiers []attribute.Modifier) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.modifiers = make([]attribute.Modifier, len(modifiers))
	copy(bs.modifiers, modifiers)
	bs.Touch()
}

// AddModifier adds an attribute modifier
func (bs *BaseSocketable) AddModifier(mod attribute.Modifier) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.modifiers = append(bs.modifiers, mod)
	bs.Touch()
}

// OnSocket is called when this socketable is inserted into equipment
func (bs *BaseSocketable) OnSocket(ctx context.Context, equipment Equipment, ent entity.Entity) error {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Apply attribute modifiers - we need to know which attribute each modifier targets
	// Since Modifier interface doesn't specify the attribute, we use the Source field
	// to determine the attribute type
	if ent != nil {
		attrs := ent.Attributes()
		for _, mod := range bs.modifiers {
			// Use Source as the attribute type indicator
			attrType := attribute.Type(mod.Source())
			attrs.AddModifier(attrType, mod)
		}
	}

	// Apply socket effect
	if bs.effect != nil {
		return bs.effect.OnSocket(ctx, equipment, ent)
	}

	return nil
}

// OnUnsocket is called when this socketable is removed from equipment
func (bs *BaseSocketable) OnUnsocket(ctx context.Context, equipment Equipment, ent entity.Entity) error {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Remove attribute modifiers
	if ent != nil {
		attrs := ent.Attributes()
		for _, mod := range bs.modifiers {
			attrType := attribute.Type(mod.Source())
			attrs.RemoveModifier(attrType, mod.ID())
		}
	}

	// Remove socket effect
	if bs.effect != nil {
		return bs.effect.OnUnsocket(ctx, equipment, ent)
	}

	return nil
}

// CanSocketIn checks if this socketable can be inserted into the given socket type
func (bs *BaseSocketable) CanSocketIn(socketType SocketType) bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Universal sockets accept anything
	if socketType == SocketTypeUniversal {
		return true
	}

	// Check if types match
	return bs.socketType == socketType
}

// --- Cloneable interface ---

func (bs *BaseSocketable) Clone() any {
	if bs == nil {
		return nil
	}

	bs.mu.RLock()
	defer bs.mu.RUnlock()

	if bs.BaseItem == nil {
		return nil
	}

	// Clone base item
	cloned := bs.BaseItem.Clone()
	if cloned == nil {
		return nil
	}
	baseClone := cloned.(*BaseItem)

	// Clone modifiers - modifiers are interfaces, we keep references
	modifiers := make([]attribute.Modifier, len(bs.modifiers))
	copy(modifiers, bs.modifiers)

	clone := &BaseSocketable{
		BaseItem:   baseClone,
		socketType: bs.socketType,
		effect:     bs.effect, // Effect is shared (stateless)
		effectID:   bs.effectID,
		tier:       bs.tier,
		modifiers:  modifiers,
	}

	return clone
}

// --- Serialization ---

// ModifierState holds serializable state for an attribute modifier
type ModifierState struct {
	ID       string  `msgpack:"id"`
	ModType  string  `msgpack:"mod_type"`
	Value    float64 `msgpack:"value"`
	Source   string  `msgpack:"source"`
	Priority int     `msgpack:"priority"`
}

// SocketableState holds the complete serializable state of a BaseSocketable
type SocketableState struct {
	State
	SocketType string          `msgpack:"socket_type"`
	EffectID   string          `msgpack:"effect_id"`
	Tier       int             `msgpack:"tier"`
	Modifiers  []ModifierState `msgpack:"modifiers"`
}

func (bs *BaseSocketable) Marshal() ([]byte, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Get base item data
	baseData, err := bs.BaseItem.Marshal()
	if err != nil {
		return nil, err
	}

	// Decode to State
	var is State
	if err := persist.DefaultCodec().Decode(baseData, &is); err != nil {
		return nil, err
	}

	// Convert modifiers to serializable state
	modStates := make([]ModifierState, len(bs.modifiers))
	for i, mod := range bs.modifiers {
		modStates[i] = ModifierState{
			ID:       mod.ID(),
			ModType:  string(mod.Type()),
			Value:    mod.Value(),
			Source:   mod.Source(),
			Priority: mod.Priority(),
		}
	}

	ss := SocketableState{
		State:      is,
		SocketType: string(bs.socketType),
		EffectID:   bs.effectID,
		Tier:       bs.tier,
		Modifiers:  modStates,
	}

	return persist.DefaultCodec().Encode(ss)
}

func (bs *BaseSocketable) Unmarshal(data []byte) error {
	var ss SocketableState
	if err := persist.DefaultCodec().Decode(data, &ss); err != nil {
		return fmt.Errorf("failed to decode socketable state: %w", err)
	}

	// Encode item state back to bytes
	itemData, err := persist.DefaultCodec().Encode(ss.State)
	if err != nil {
		return err
	}

	bs.mu.Lock()
	defer bs.mu.Unlock()

	// Initialize base item if nil
	if bs.BaseItem == nil {
		bs.BaseItem = &BaseItem{}
	}

	// Restore base item
	if err := bs.BaseItem.Unmarshal(itemData); err != nil {
		return err
	}

	// Restore socketable-specific fields
	bs.socketType = SocketType(ss.SocketType)
	bs.effectID = ss.EffectID
	bs.tier = ss.Tier

	// Restore modifiers
	bs.modifiers = make([]attribute.Modifier, len(ss.Modifiers))
	for i, ms := range ss.Modifiers {
		bs.modifiers[i] = attribute.NewModifierWithPriority(
			ms.ID,
			attribute.ModifierType(ms.ModType),
			ms.Value,
			ms.Source,
			ms.Priority,
		)
	}

	// Note: effect must be restored separately using effectID
	// via an effect registry

	return nil
}

// --- Validatable interface ---

func (bs *BaseSocketable) Validate() error {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Validate base item
	if err := bs.BaseItem.Validate(); err != nil {
		return err
	}

	if bs.socketType == "" {
		return fmt.Errorf("socket type cannot be empty")
	}

	if bs.tier < 1 || bs.tier > 5 {
		return fmt.Errorf("tier must be between 1 and 5")
	}

	return nil
}

// --- Serializable interface (required by infra.Persistent) ---

func (bs *BaseSocketable) SerializeState() (map[string]any, error) {
	data, err := bs.Marshal()
	if err != nil {
		return nil, err
	}

	var state map[string]any
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return nil, err
	}
	return state, nil
}

func (bs *BaseSocketable) DeserializeState(state map[string]any) error {
	data, err := persist.DefaultCodec().Encode(state)
	if err != nil {
		return err
	}
	return bs.Unmarshal(data)
}
