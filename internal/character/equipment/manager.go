package equipment

import (
	"context"
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/item"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

// Manager handles character equipment
type Manager interface {
	// Get returns equipment in slot (nil if empty)
	Get(slot Slot) item.Equipment

	// Equip puts an item in the appropriate slot, returns unequipped item if any
	Equip(ctx context.Context, equip item.Equipment) (item.Equipment, error)

	// EquipToSlot puts an item in specific slot, returns unequipped item if any
	EquipToSlot(ctx context.Context, slot Slot, equip item.Equipment) (item.Equipment, error)

	// Unequip removes item from slot, returns the item
	Unequip(ctx context.Context, slot Slot) (item.Equipment, error)

	// UnequipAll removes all equipment, returns slice of unequipped items
	UnequipAll(ctx context.Context) ([]item.Equipment, error)

	// IsEmpty returns true if slot is empty
	IsEmpty(slot Slot) bool

	// GetAll returns all equipped items
	GetAll() map[Slot]item.Equipment

	// GetFilledSlots returns list of slots that have equipment
	GetFilledSlots() []Slot

	// CanEquip checks if item can be equipped (requirements met)
	CanEquip(equip item.Equipment) bool

	// CanEquipToSlot checks if item can be equipped to specific slot
	CanEquipToSlot(slot Slot, equip item.Equipment) bool

	// TotalWeight returns total weight of all equipped items
	TotalWeight() float64

	// GetAllModifiers returns combined attribute modifiers from all equipment
	GetAllModifiers() []attribute.Modifier

	// SetOwner sets the entity that owns this equipment
	SetOwner(owner entity.Entity)

	// Owner returns the current owner
	Owner() entity.Entity

	// OnEquip registers callback when item is equipped
	OnEquip(callback EquipCallback)

	// OnUnequip registers callback when item is unequipped
	OnUnequip(callback EquipCallback)

	// SerializeState converts state to map for persistence
	SerializeState() (map[string]any, error)

	// DeserializeState restores state from map
	DeserializeState(state map[string]any) error
}

// EquipCallback is invoked for equipment events
type EquipCallback func(ctx context.Context, slot Slot, equip item.Equipment)

var _ Manager = (*BaseManager)(nil)

// BaseManager implements Manager interface
type BaseManager struct {
	mu sync.RWMutex

	slots              map[Slot]item.Equipment
	owner              entity.Entity
	onEquipCallbacks   []EquipCallback
	onUnequipCallbacks []EquipCallback
}

// NewManager creates a new equipment manager
func NewManager() *BaseManager {
	return &BaseManager{
		slots: make(map[Slot]item.Equipment),
	}
}

// NewManagerWithOwner creates a new equipment manager with owner
func NewManagerWithOwner(owner entity.Entity) *BaseManager {
	m := NewManager()
	m.owner = owner
	return m
}

func (m *BaseManager) Get(slot Slot) item.Equipment {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.slots[slot]
}

func (m *BaseManager) Equip(ctx context.Context, equip item.Equipment) (item.Equipment, error) {
	if equip == nil {
		return nil, fmt.Errorf("cannot equip nil item")
	}

	// Find the default slot for this item
	slot := DefaultSlot(equip.ItemType())
	if slot == "" {
		return nil, fmt.Errorf("item type %s cannot be equipped", equip.ItemType())
	}

	// If default slot is occupied, try alternative slots
	m.mu.RLock()
	if m.slots[slot] != nil {
		// Try to find an empty compatible slot
		compatibleSlots := CompatibleSlots(equip.ItemType())
		foundEmpty := false
		for _, s := range compatibleSlots {
			if m.slots[s] == nil {
				slot = s
				foundEmpty = true
				break
			}
		}
		if !foundEmpty {
			// Use default slot, will swap
			slot = compatibleSlots[0]
		}
	}
	m.mu.RUnlock()

	return m.EquipToSlot(ctx, slot, equip)
}

func (m *BaseManager) EquipToSlot(ctx context.Context, slot Slot, equip item.Equipment) (item.Equipment, error) {
	if equip == nil {
		return nil, fmt.Errorf("cannot equip nil item")
	}

	// Check slot compatibility
	if !IsCompatible(slot, equip.ItemType()) {
		return nil, fmt.Errorf("item type %s cannot be equipped to slot %s", equip.ItemType(), slot)
	}

	// Check if owner can equip this item
	m.mu.RLock()
	owner := m.owner
	m.mu.RUnlock()

	if owner != nil && !equip.CanEquip(owner) {
		return nil, fmt.Errorf("character does not meet equipment requirements")
	}

	m.mu.Lock()
	// Get currently equipped item (if any)
	previousItem := m.slots[slot]

	// Unequip previous item
	if previousItem != nil {
		if owner != nil {
			if err := previousItem.OnUnequip(ctx, owner); err != nil {
				m.mu.Unlock()
				return nil, fmt.Errorf("failed to unequip previous item: %w", err)
			}
		}
	}

	// Equip new item
	m.slots[slot] = equip
	m.mu.Unlock()

	// Apply equipment effects
	if owner != nil {
		if err := equip.OnEquip(ctx, owner); err != nil {
			// Rollback
			m.mu.Lock()
			m.slots[slot] = previousItem
			m.mu.Unlock()
			return nil, fmt.Errorf("failed to equip item: %w", err)
		}
	}

	// Trigger callbacks
	m.mu.RLock()
	if previousItem != nil {
		for _, cb := range m.onUnequipCallbacks {
			cb(ctx, slot, previousItem)
		}
	}
	for _, cb := range m.onEquipCallbacks {
		cb(ctx, slot, equip)
	}
	m.mu.RUnlock()

	return previousItem, nil
}

func (m *BaseManager) Unequip(ctx context.Context, slot Slot) (item.Equipment, error) {
	m.mu.Lock()
	equip := m.slots[slot]
	if equip == nil {
		m.mu.Unlock()
		return nil, nil // Nothing to unequip
	}

	owner := m.owner
	delete(m.slots, slot)
	m.mu.Unlock()

	// Remove equipment effects
	if owner != nil {
		if err := equip.OnUnequip(ctx, owner); err != nil {
			// Rollback
			m.mu.Lock()
			m.slots[slot] = equip
			m.mu.Unlock()
			return nil, fmt.Errorf("failed to unequip item: %w", err)
		}
	}

	// Trigger callbacks
	m.mu.RLock()
	for _, cb := range m.onUnequipCallbacks {
		cb(ctx, slot, equip)
	}
	m.mu.RUnlock()

	return equip, nil
}

func (m *BaseManager) UnequipAll(ctx context.Context) ([]item.Equipment, error) {
	m.mu.RLock()
	slots := make([]Slot, 0, len(m.slots))
	for slot := range m.slots {
		slots = append(slots, slot)
	}
	m.mu.RUnlock()

	unequipped := make([]item.Equipment, 0, len(slots))
	for _, slot := range slots {
		equip, err := m.Unequip(ctx, slot)
		if err != nil {
			return unequipped, err
		}
		if equip != nil {
			unequipped = append(unequipped, equip)
		}
	}

	return unequipped, nil
}

func (m *BaseManager) IsEmpty(slot Slot) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.slots[slot] == nil
}

func (m *BaseManager) GetAll() map[Slot]item.Equipment {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[Slot]item.Equipment, len(m.slots))
	for slot, equip := range m.slots {
		result[slot] = equip
	}
	return result
}

func (m *BaseManager) GetFilledSlots() []Slot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	slots := make([]Slot, 0, len(m.slots))
	for slot := range m.slots {
		slots = append(slots, slot)
	}
	return slots
}

func (m *BaseManager) CanEquip(equip item.Equipment) bool {
	if equip == nil {
		return false
	}

	// Check slot compatibility
	if DefaultSlot(equip.ItemType()) == "" {
		return false
	}

	// Check requirements
	m.mu.RLock()
	owner := m.owner
	m.mu.RUnlock()

	if owner != nil {
		return equip.CanEquip(owner)
	}

	return true
}

func (m *BaseManager) CanEquipToSlot(slot Slot, equip item.Equipment) bool {
	if equip == nil {
		return false
	}

	if !IsCompatible(slot, equip.ItemType()) {
		return false
	}

	m.mu.RLock()
	owner := m.owner
	m.mu.RUnlock()

	if owner != nil {
		return equip.CanEquip(owner)
	}

	return true
}

func (m *BaseManager) TotalWeight() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var total float64
	for _, equip := range m.slots {
		if equip != nil {
			total += equip.Weight()
		}
	}
	return total
}

func (m *BaseManager) GetAllModifiers() []attribute.Modifier {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var mods []attribute.Modifier
	for _, equip := range m.slots {
		if equip != nil {
			mods = append(mods, equip.Attributes()...)
		}
	}
	return mods
}

func (m *BaseManager) SetOwner(owner entity.Entity) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.owner = owner
}

func (m *BaseManager) Owner() entity.Entity {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.owner
}

func (m *BaseManager) OnEquip(callback EquipCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onEquipCallbacks = append(m.onEquipCallbacks, callback)
}

func (m *BaseManager) OnUnequip(callback EquipCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onUnequipCallbacks = append(m.onUnequipCallbacks, callback)
}

// State holds serializable equipment state
type State struct {
	Slots map[string]string `msgpack:"slots"` // slot -> item ID mapping
}

func (m *BaseManager) SerializeState() (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	slots := make(map[string]string, len(m.slots))
	for slot, equip := range m.slots {
		if equip != nil {
			slots[string(slot)] = equip.ID()
		}
	}

	state := State{Slots: slots}
	data, err := persist.DefaultCodec().Encode(state)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := persist.DefaultCodec().Decode(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *BaseManager) DeserializeState(stateData map[string]any) error {
	data, err := persist.DefaultCodec().Encode(stateData)
	if err != nil {
		return err
	}

	var state State
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return err
	}

	// Note: Actual item restoration requires an item registry/repository
	// This just stores the slot -> itemID mapping for later resolution
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear current slots
	m.slots = make(map[Slot]item.Equipment)

	// Store slot mappings (items will be resolved by character loader)
	// For now, we just record the structure
	// The actual items need to be loaded separately and re-equipped

	return nil
}

// GetSlotItemIDs returns a map of slot to item ID for persistence
func (m *BaseManager) GetSlotItemIDs() map[Slot]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[Slot]string, len(m.slots))
	for slot, equip := range m.slots {
		if equip != nil {
			result[slot] = equip.ID()
		}
	}
	return result
}

// SetItemDirect sets an item in a slot without triggering callbacks
// Used during deserialization
func (m *BaseManager) SetItemDirect(slot Slot, equip item.Equipment) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.slots[slot] = equip
}
