package item

import (
	"context"
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/item/affix"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

var _ Equipment = (*BaseEquipment)(nil)

// BaseEquipment implements Equipment interface
type BaseEquipment struct {
	*BaseItem

	mu sync.RWMutex

	slot          EquipmentSlot
	attributes    []attribute.Modifier
	durability    float64
	maxDurability float64
	sockets       []Socketable
	socketTypes   []SocketType // Types of allowed sockets
	affixSet      affix.Set
	requirements  EquipRequirements

	// Callbacks for equip/unequip events
	onEquipFn   func(ctx context.Context, entity entity.Entity) error
	onUnequipFn func(ctx context.Context, entity entity.Entity) error
}

// EquipmentConfig holds configuration for creating equipment
type EquipmentConfig struct {
	BaseItemConfig
	Slot          EquipmentSlot
	MaxDurability float64
	SocketCount   int
	SocketTypes   []SocketType
	Requirements  EquipRequirements
}

// NewBaseEquipment creates new equipment with minimal configuration
func NewBaseEquipment(id string, itemType Type, name string, slot EquipmentSlot) *BaseEquipment {
	return NewEquipmentWithConfig(EquipmentConfig{
		BaseItemConfig: BaseItemConfig{
			ID:       id,
			Name:     name,
			ItemType: itemType,
		},
		Slot: slot,
	})
}

// NewEquipmentWithConfig creates new equipment with full configuration
func NewEquipmentWithConfig(cfg EquipmentConfig) *BaseEquipment {
	be := &BaseEquipment{
		BaseItem:      NewBaseItemWithConfig(cfg.BaseItemConfig),
		slot:          cfg.Slot,
		attributes:    make([]attribute.Modifier, 0),
		durability:    cfg.MaxDurability,
		maxDurability: cfg.MaxDurability,
		sockets:       make([]Socketable, cfg.SocketCount),
		socketTypes:   cfg.SocketTypes,
		affixSet:      affix.NewBaseSet(),
		requirements:  cfg.Requirements,
	}

	// Apply defaults
	if be.maxDurability <= 0 {
		be.maxDurability = 100.0
		be.durability = 100.0
	}
	if be.requirements == nil {
		be.requirements = NewSimpleRequirements(1, nil)
	}
	if be.socketTypes == nil {
		be.socketTypes = make([]SocketType, cfg.SocketCount)
		for i := range be.socketTypes {
			be.socketTypes[i] = SocketTypeUniversal
		}
	}

	return be
}

// --- Equipment interface implementation ---

func (be *BaseEquipment) Slot() EquipmentSlot {
	be.mu.RLock()
	defer be.mu.RUnlock()
	return be.slot
}

func (be *BaseEquipment) Attributes() []attribute.Modifier {
	be.mu.RLock()
	defer be.mu.RUnlock()

	// Combine base attributes with affix modifiers
	allMods := make([]attribute.Modifier, len(be.attributes))
	copy(allMods, be.attributes)

	if be.affixSet != nil {
		allMods = append(allMods, be.affixSet.AllModifiers()...)
	}

	// Add socket effect modifiers
	for _, socket := range be.sockets {
		if socket != nil && socket.Effect() != nil {
			// Socket effects are applied separately via OnEquip
		}
	}

	return allMods
}

// AddAttribute adds a base attribute modifier to the equipment
func (be *BaseEquipment) AddAttribute(mod attribute.Modifier) {
	be.mu.Lock()
	defer be.mu.Unlock()
	be.attributes = append(be.attributes, mod)
	be.Touch()
}

// RemoveAttribute removes an attribute modifier by ID
func (be *BaseEquipment) RemoveAttribute(modID string) bool {
	be.mu.Lock()
	defer be.mu.Unlock()
	for i, mod := range be.attributes {
		if mod.ID() == modID {
			be.attributes = append(be.attributes[:i], be.attributes[i+1:]...)
			be.Touch()
			return true
		}
	}
	return false
}

func (be *BaseEquipment) Durability() float64 {
	be.mu.RLock()
	defer be.mu.RUnlock()
	return be.durability
}

func (be *BaseEquipment) MaxDurability() float64 {
	be.mu.RLock()
	defer be.mu.RUnlock()
	return be.maxDurability
}

func (be *BaseEquipment) SetDurability(value float64) {
	be.mu.Lock()
	defer be.mu.Unlock()
	if value < 0 {
		value = 0
	} else if value > be.maxDurability {
		value = be.maxDurability
	}
	be.durability = value
	be.Touch()
}

func (be *BaseEquipment) SetMaxDurability(value float64) {
	be.mu.Lock()
	defer be.mu.Unlock()
	if value < 1 {
		value = 1
	}
	be.maxDurability = value
	if be.durability > be.maxDurability {
		be.durability = be.maxDurability
	}
	be.Touch()
}

func (be *BaseEquipment) Repair(amount float64) {
	be.mu.Lock()
	defer be.mu.Unlock()
	be.durability += amount
	if be.durability > be.maxDurability {
		be.durability = be.maxDurability
	}
	be.Touch()
}

func (be *BaseEquipment) RepairFull() {
	be.mu.Lock()
	defer be.mu.Unlock()
	be.durability = be.maxDurability
	be.Touch()
}

func (be *BaseEquipment) DamageItem(amount float64) bool {
	be.mu.Lock()
	defer be.mu.Unlock()
	be.durability -= amount
	if be.durability < 0 {
		be.durability = 0
	}
	be.Touch()
	return be.durability > 0
}

func (be *BaseEquipment) IsBroken() bool {
	be.mu.RLock()
	defer be.mu.RUnlock()
	return be.durability <= 0
}

// DurabilityPercent returns durability as percentage [0-1]
func (be *BaseEquipment) DurabilityPercent() float64 {
	be.mu.RLock()
	defer be.mu.RUnlock()
	if be.maxDurability <= 0 {
		return 0
	}
	return be.durability / be.maxDurability
}

func (be *BaseEquipment) SocketCount() int {
	be.mu.RLock()
	defer be.mu.RUnlock()
	return len(be.sockets)
}

// EmptySocketCount returns number of empty sockets
func (be *BaseEquipment) EmptySocketCount() int {
	be.mu.RLock()
	defer be.mu.RUnlock()
	count := 0
	for _, s := range be.sockets {
		if s == nil {
			count++
		}
	}
	return count
}

func (be *BaseEquipment) GetSocket(index int) (Socketable, bool) {
	be.mu.RLock()
	defer be.mu.RUnlock()
	if index < 0 || index >= len(be.sockets) {
		return nil, false
	}
	socket := be.sockets[index]
	return socket, socket != nil
}

// GetSocketType returns the type of socket at index
func (be *BaseEquipment) GetSocketType(index int) (SocketType, bool) {
	be.mu.RLock()
	defer be.mu.RUnlock()
	if index < 0 || index >= len(be.socketTypes) {
		return "", false
	}
	return be.socketTypes[index], true
}

func (be *BaseEquipment) SetSocket(index int, item Socketable) error {
	be.mu.Lock()
	defer be.mu.Unlock()

	if index < 0 || index >= len(be.sockets) {
		return fmt.Errorf("socket index out of range: %d", index)
	}

	if be.sockets[index] != nil {
		return fmt.Errorf("socket %d already occupied", index)
	}

	// Check socket type compatibility
	if item != nil && index < len(be.socketTypes) {
		socketType := be.socketTypes[index]
		if socketType != SocketTypeUniversal && socketType != item.SocketType() {
			return fmt.Errorf("socket type mismatch: expected %s, got %s", socketType, item.SocketType())
		}
	}

	be.sockets[index] = item
	be.Touch()
	return nil
}

func (be *BaseEquipment) RemoveSocket(index int) (Socketable, error) {
	be.mu.Lock()
	defer be.mu.Unlock()

	if index < 0 || index >= len(be.sockets) {
		return nil, fmt.Errorf("socket index out of range: %d", index)
	}

	item := be.sockets[index]
	be.sockets[index] = nil
	be.Touch()
	return item, nil
}

// AddSocket adds a new socket to the equipment
func (be *BaseEquipment) AddSocket(socketType SocketType) {
	be.mu.Lock()
	defer be.mu.Unlock()
	be.sockets = append(be.sockets, nil)
	be.socketTypes = append(be.socketTypes, socketType)
	be.Touch()
}

func (be *BaseEquipment) Affixes() affix.Set {
	be.mu.RLock()
	defer be.mu.RUnlock()
	return be.affixSet
}

func (be *BaseEquipment) Requirements() EquipRequirements {
	be.mu.RLock()
	defer be.mu.RUnlock()
	return be.requirements
}

func (be *BaseEquipment) SetRequirements(req EquipRequirements) {
	be.mu.Lock()
	defer be.mu.Unlock()
	be.requirements = req
	be.Touch()
}

func (be *BaseEquipment) CanEquip(ent entity.Entity) bool {
	be.mu.RLock()
	defer be.mu.RUnlock()

	// Can't equip broken items
	if be.durability <= 0 {
		return false
	}

	// Nil entity cannot equip
	if ent == nil {
		return false
	}

	if be.requirements == nil {
		return true
	}

	return be.requirements.Check(ent)
}

func (be *BaseEquipment) OnEquip(ctx context.Context, ent entity.Entity) error {
	be.mu.RLock()
	onEquip := be.onEquipFn

	// Collect all modifiers
	allMods := make([]attribute.Modifier, len(be.attributes))
	copy(allMods, be.attributes)
	if be.affixSet != nil {
		allMods = append(allMods, be.affixSet.AllModifiers()...)
	}

	// Collect sockets that have effects
	socketsWithEffects := make([]Socketable, 0)
	for _, socket := range be.sockets {
		if socket != nil && socket.Effect() != nil {
			socketsWithEffects = append(socketsWithEffects, socket)
		}
	}
	be.mu.RUnlock()

	// Apply attribute modifiers to entity
	attrMgr := ent.Attributes()
	for _, mod := range allMods {
		attrMgr.AddModifier(attribute.Type(mod.Source()), mod)
	}

	// Apply socket effects (outside lock to avoid deadlock)
	for _, socket := range socketsWithEffects {
		if err := socket.Effect().OnSocket(ctx, be, ent); err != nil {
			return fmt.Errorf("failed to apply socket effect: %w", err)
		}
	}

	// Call custom callback if set
	if onEquip != nil {
		return onEquip(ctx, ent)
	}

	return nil
}

func (be *BaseEquipment) OnUnequip(ctx context.Context, ent entity.Entity) error {
	be.mu.RLock()
	onUnequip := be.onUnequipFn

	// Collect all modifiers
	allMods := make([]attribute.Modifier, len(be.attributes))
	copy(allMods, be.attributes)
	if be.affixSet != nil {
		allMods = append(allMods, be.affixSet.AllModifiers()...)
	}

	// Collect sockets that have effects
	socketsWithEffects := make([]Socketable, 0)
	for _, socket := range be.sockets {
		if socket != nil && socket.Effect() != nil {
			socketsWithEffects = append(socketsWithEffects, socket)
		}
	}
	be.mu.RUnlock()

	// Remove attribute modifiers from entity
	attrMgr := ent.Attributes()
	for _, mod := range allMods {
		attrMgr.RemoveModifier(attribute.Type(mod.Source()), mod.ID())
	}

	// Remove socket effects (outside lock to avoid deadlock)
	for _, socket := range socketsWithEffects {
		if err := socket.Effect().OnUnsocket(ctx, be, ent); err != nil {
			return fmt.Errorf("failed to remove socket effect: %w", err)
		}
	}

	// Call custom callback if set
	if onUnequip != nil {
		return onUnequip(ctx, ent)
	}

	return nil
}

// SetOnEquip sets custom equip callback
func (be *BaseEquipment) SetOnEquip(fn func(ctx context.Context, entity entity.Entity) error) {
	be.mu.Lock()
	defer be.mu.Unlock()
	be.onEquipFn = fn
}

// SetOnUnequip sets custom unequip callback
func (be *BaseEquipment) SetOnUnequip(fn func(ctx context.Context, entity entity.Entity) error) {
	be.mu.Lock()
	defer be.mu.Unlock()
	be.onUnequipFn = fn
}

// --- Cloneable interface ---

func (be *BaseEquipment) Clone() any {
	if be == nil {
		return nil
	}

	be.mu.RLock()
	defer be.mu.RUnlock()

	if be.BaseItem == nil {
		return nil
	}

	cloned := be.BaseItem.Clone()
	if cloned == nil {
		return nil
	}
	baseClone := cloned.(*BaseItem)

	clone := &BaseEquipment{
		BaseItem:      baseClone,
		slot:          be.slot,
		attributes:    make([]attribute.Modifier, len(be.attributes)),
		durability:    be.durability,
		maxDurability: be.maxDurability,
		sockets:       make([]Socketable, len(be.sockets)),
		socketTypes:   make([]SocketType, len(be.socketTypes)),
		affixSet:      affix.NewBaseSet(),
		requirements:  be.requirements, // Requirements typically shared
	}

	copy(clone.attributes, be.attributes)
	copy(clone.socketTypes, be.socketTypes)

	// Clone sockets (socketables are not cloned - they're separate items)
	for i, s := range be.sockets {
		if s != nil {
			if cloneable, ok := s.(interface{ Clone() any }); ok {
				clone.sockets[i] = cloneable.Clone().(Socketable)
			}
		}
	}

	// Clone affixes
	if be.affixSet != nil {
		for _, a := range be.affixSet.GetAll() {
			_ = clone.affixSet.Add(a)
		}
	}

	return clone
}

// --- Serialization ---

// EquipmentState holds serializable state of equipment
type EquipmentState struct {
	State
	Slot          string             `msgpack:"slot"`
	Durability    float64            `msgpack:"durability"`
	MaxDurability float64            `msgpack:"max_durability"`
	SocketTypes   []string           `msgpack:"socket_types"`
	SocketIDs     []string           `msgpack:"socket_ids"`
	AffixIDs      []string           `msgpack:"affix_ids"`
	ReqLevel      int                `msgpack:"req_level"`
	ReqAttrs      map[string]float64 `msgpack:"req_attrs"`
}

func (be *BaseEquipment) Marshal() ([]byte, error) {
	be.mu.RLock()
	defer be.mu.RUnlock()

	// Get base item state
	baseData, err := be.BaseItem.Marshal()
	if err != nil {
		return nil, err
	}

	var itemState State
	if err := persist.DefaultCodec().Decode(baseData, &itemState); err != nil {
		return nil, err
	}

	// Build socket type list
	socketTypes := make([]string, len(be.socketTypes))
	for i, st := range be.socketTypes {
		socketTypes[i] = string(st)
	}

	// Build socket ID list
	socketIDs := make([]string, len(be.sockets))
	for i, s := range be.sockets {
		if s != nil {
			socketIDs[i] = s.ID()
		}
	}

	// Build affix ID list
	var affixIDs []string
	if be.affixSet != nil {
		for _, a := range be.affixSet.GetAll() {
			affixIDs = append(affixIDs, a.AffixID())
		}
	}

	// Get requirements
	var reqLevel int
	var reqAttrs map[string]float64
	if be.requirements != nil {
		reqLevel = be.requirements.Level()
		attrs := be.requirements.Attributes()
		if len(attrs) > 0 {
			reqAttrs = make(map[string]float64)
			for k, v := range attrs {
				reqAttrs[string(k)] = v
			}
		}
	}

	state := EquipmentState{
		State:         itemState,
		Slot:          string(be.slot),
		Durability:    be.durability,
		MaxDurability: be.maxDurability,
		SocketTypes:   socketTypes,
		SocketIDs:     socketIDs,
		AffixIDs:      affixIDs,
		ReqLevel:      reqLevel,
		ReqAttrs:      reqAttrs,
	}

	return persist.DefaultCodec().Encode(state)
}

func (be *BaseEquipment) Unmarshal(data []byte) error {
	var state EquipmentState
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return fmt.Errorf("failed to decode equipment state: %w", err)
	}

	// Restore base item
	itemData, err := persist.DefaultCodec().Encode(state.State)
	if err != nil {
		return err
	}

	if be.BaseItem == nil {
		be.BaseItem = &BaseItem{}
	}
	if err := be.BaseItem.Unmarshal(itemData); err != nil {
		return err
	}

	be.mu.Lock()
	defer be.mu.Unlock()

	be.slot = EquipmentSlot(state.Slot)
	be.durability = state.Durability
	be.maxDurability = state.MaxDurability

	// Restore socket types
	be.socketTypes = make([]SocketType, len(state.SocketTypes))
	for i, st := range state.SocketTypes {
		be.socketTypes[i] = SocketType(st)
	}

	// Initialize empty sockets (actual items restored separately)
	be.sockets = make([]Socketable, len(state.SocketIDs))

	// Initialize affix set (actual affixes restored separately)
	be.affixSet = affix.NewBaseSet()

	// Restore requirements
	if state.ReqAttrs != nil {
		reqAttrs := make(map[attribute.Type]float64)
		for k, v := range state.ReqAttrs {
			reqAttrs[attribute.Type(k)] = v
		}
		be.requirements = NewSimpleRequirements(state.ReqLevel, reqAttrs)
	} else {
		be.requirements = NewSimpleRequirements(state.ReqLevel, nil)
	}

	return nil
}

// --- Validatable interface ---

func (be *BaseEquipment) Validate() error {
	if err := be.BaseItem.Validate(); err != nil {
		return err
	}

	be.mu.RLock()
	defer be.mu.RUnlock()

	if be.slot == "" {
		return fmt.Errorf("equipment must have a slot")
	}
	if be.durability < 0 {
		return fmt.Errorf("durability cannot be negative")
	}
	if be.maxDurability < 1 {
		return fmt.Errorf("max durability must be at least 1")
	}
	if be.durability > be.maxDurability {
		return fmt.Errorf("durability cannot exceed max durability")
	}

	return nil
}

// --- Serializable interface (required by infra.Persistent) ---

func (be *BaseEquipment) SerializeState() (map[string]any, error) {
	data, err := be.Marshal()
	if err != nil {
		return nil, err
	}

	var state map[string]any
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return nil, err
	}
	return state, nil
}

func (be *BaseEquipment) DeserializeState(state map[string]any) error {
	data, err := persist.DefaultCodec().Encode(state)
	if err != nil {
		return err
	}
	return be.Unmarshal(data)
}
