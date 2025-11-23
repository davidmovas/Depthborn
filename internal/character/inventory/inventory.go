package inventory

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/item"
)

// Inventory manages character items
type Inventory interface {
	// Capacity returns maximum number of items
	Capacity() int

	// SetCapacity updates maximum capacity
	SetCapacity(capacity int)

	// Count returns current number of items
	Count() int

	// IsFull returns true if inventory is at capacity
	IsFull() bool

	// Add adds item to inventory
	Add(ctx context.Context, item item.Item) error

	// AddMultiple adds multiple items
	AddMultiple(ctx context.Context, items []item.Item) error

	// Remove removes item by ID
	Remove(ctx context.Context, itemID string) (item.Item, error)

	// RemoveMultiple removes multiple items
	RemoveMultiple(ctx context.Context, itemIDs []string) ([]item.Item, error)

	// Get retrieves item by ID
	Get(itemID string) (item.Item, bool)

	// GetAll returns all items
	GetAll() []item.Item

	// GetByType returns items of specified type
	GetByType(itemType item.Type) []item.Item

	// GetByRarity returns items of specified rarity
	GetByRarity(rarity item.Rarity) []item.Item

	// Contains checks if inventory has item
	Contains(itemID string) bool

	// FindSpace returns true if space available for item
	FindSpace(item item.Item) bool

	// Sort sorts inventory by criteria
	Sort(criteria SortCriteria)

	// Filter returns items matching predicate
	Filter(predicate FilterFunc) []item.Item

	// TotalWeight returns combined weight of all items
	TotalWeight() float64

	// Clear removes all items
	Clear(ctx context.Context) error

	// OnItemAdded registers callback when item is added
	OnItemAdded(callback ItemCallback)

	// OnItemRemoved registers callback when item is removed
	OnItemRemoved(callback ItemCallback)
}

// SortCriteria defines how to sort inventory
type SortCriteria struct {
	Field     string
	Ascending bool
}

// FilterFunc returns true if item should be included
type FilterFunc func(item item.Item) bool

// ItemCallback is invoked for inventory events
type ItemCallback func(ctx context.Context, item item.Item)

// Equipment manages equipped items
type Equipment interface {
	// Equip equips item to slot
	Equip(ctx context.Context, slot item.EquipmentSlot, item item.Equipment) error

	// Unequip removes item from slot
	Unequip(ctx context.Context, slot item.EquipmentSlot) (item.Equipment, error)

	// Get returns equipped item in slot
	Get(slot item.EquipmentSlot) (item.Equipment, bool)

	// GetAll returns all equipped items
	GetAll() map[item.EquipmentSlot]item.Equipment

	// Swap exchanges items between two slots
	Swap(ctx context.Context, slot1, slot2 item.EquipmentSlot) error

	// CanEquip checks if item can be equipped to slot
	CanEquip(slot item.EquipmentSlot, item item.Equipment) bool

	// UnequipAll removes all equipped items
	UnequipAll(ctx context.Context) ([]item.Equipment, error)

	// AllModifiers returns combined modifiers from all equipment
	AllModifiers() []any

	// OnEquip registers callback when item is equipped
	OnEquip(callback EquipCallback)

	// OnUnequip registers callback when item is unequipped
	OnUnequip(callback EquipCallback)
}

// EquipCallback is invoked for equipment events
type EquipCallback func(ctx context.Context, slot item.EquipmentSlot, item item.Equipment)

// Stash represents account-wide shared storage
type Stash interface {
	// Tabs returns all stash tabs
	Tabs() []StashTab

	// GetTab returns tab by index
	GetTab(index int) (StashTab, bool)

	// AddTab creates new stash tab
	AddTab(name string, capacity int) error

	// RemoveTab removes stash tab
	RemoveTab(index int) error

	// RenameTab updates tab name
	RenameTab(index int, name string) error

	// TransferToTab moves item to specified tab
	TransferToTab(ctx context.Context, item item.Item, tabIndex int) error

	// FindItem searches all tabs for item
	FindItem(itemID string) (item.Item, int, bool)

	// TotalCapacity returns combined capacity of all tabs
	TotalCapacity() int

	// TotalCount returns total items across all tabs
	TotalCount() int
}

// StashTab represents single stash tab
type StashTab interface {
	// Name returns tab name
	Name() string

	// SetName updates tab name
	SetName(name string)

	// Inventory returns tab inventory
	Inventory() Inventory

	// Color returns tab color for UI
	Color() string

	// SetColor updates tab color
	SetColor(color string)

	// Icon returns tab icon identifier
	Icon() string

	// SetIcon updates tab icon
	SetIcon(icon string)
}
