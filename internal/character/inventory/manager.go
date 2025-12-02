package inventory

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/davidmovas/Depthborn/internal/item"
	"github.com/davidmovas/Depthborn/pkg/identifier"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

// Manager handles character inventory with weight and slot limits
type Manager interface {
	// --- Basic Operations ---

	// Add adds an item to inventory (auto-stacks if possible)
	Add(ctx context.Context, itm item.Item) error

	// AddToSlot adds an item to a specific slot
	AddToSlot(ctx context.Context, slot int, itm item.Item) error

	// Remove removes item by ID completely
	Remove(ctx context.Context, itemID string) (item.Item, error)

	// RemoveAmount removes specific amount from a stack, returns the removed portion
	RemoveAmount(ctx context.Context, itemID string, amount int) (item.Item, error)

	// Get returns item by ID
	Get(itemID string) (item.Item, bool)

	// GetAtSlot returns item at specific slot
	GetAtSlot(slot int) (item.Item, bool)

	// GetAll returns all items in inventory
	GetAll() []item.Item

	// Clear removes all items
	Clear(ctx context.Context) []item.Item

	// --- Stack Operations ---

	// SplitStack splits a stack into two, returns the new stack
	SplitStack(ctx context.Context, itemID string, amount int) (item.Item, error)

	// MergeStacks merges source stack into target stack
	MergeStacks(ctx context.Context, sourceID, targetID string) error

	// CanStackWith checks if item can stack with existing items
	CanStackWith(itm item.Item) (string, bool)

	// --- Slot Management ---

	// SlotCount returns number of slots
	SlotCount() int

	// SetSlotCount changes number of slots
	SetSlotCount(count int)

	// UsedSlots returns number of occupied slots
	UsedSlots() int

	// FreeSlots returns number of available slots
	FreeSlots() int

	// SwapSlots swaps items between two slots
	SwapSlots(ctx context.Context, slot1, slot2 int) error

	// MoveToSlot moves item to a different slot
	MoveToSlot(ctx context.Context, itemID string, targetSlot int) error

	// --- Weight Management ---

	// CurrentWeight returns current total weight
	CurrentWeight() float64

	// MaxWeight returns maximum weight capacity
	MaxWeight() float64

	// SetMaxWeight sets maximum weight capacity
	SetMaxWeight(weight float64)

	// AvailableWeight returns remaining weight capacity
	AvailableWeight() float64

	// --- Capacity Checks ---

	// CanAdd checks if item can be added (weight + slot check)
	CanAdd(itm item.Item) bool

	// IsFull returns true if no more items can be added
	IsFull() bool

	// Contains checks if item exists
	Contains(itemID string) bool

	// Count returns number of unique items (stacks count as 1)
	Count() int

	// TotalItems returns total item count including stack sizes
	TotalItems() int

	// --- Search & Filter ---

	// Search finds items matching query string (name contains)
	Search(query string) []item.Item

	// FindByType returns items of given type
	FindByType(itemType item.Type) []item.Item

	// FindByRarity returns items of given rarity
	FindByRarity(rarity item.Rarity) []item.Item

	// FindByTag returns items with given tag
	FindByTag(tag string) []item.Item

	// FindByTags returns items having ALL specified tags
	FindByTags(tags ...string) []item.Item

	// FindByAnyTag returns items having ANY of specified tags
	FindByAnyTag(tags ...string) []item.Item

	// FindByLevel returns items within level range
	FindByLevel(minLevel, maxLevel int) []item.Item

	// FindStackable returns all stackable items
	FindStackable() []item.Item

	// Filter returns items matching predicate
	Filter(predicate func(item.Item) bool) []item.Item

	// --- Sorting ---

	// Sort sorts inventory by given criteria
	Sort(criteria SortBy, ascending bool)

	// GetSorted returns sorted copy without modifying internal order
	GetSorted(criteria SortBy, ascending bool) []item.Item

	// --- Stats ---

	// TotalValue returns combined value of all items
	TotalValue() int64

	// TotalWeight returns current weight (alias for CurrentWeight)
	TotalWeight() float64

	// WeightPercent returns weight usage as percentage [0.0 - 1.0]
	WeightPercent() float64

	// SlotPercent returns slot usage as percentage [0.0 - 1.0]
	SlotPercent() float64

	// --- Callbacks ---

	// OnItemAdded registers callback when item is added
	OnItemAdded(callback ItemCallback)

	// OnItemRemoved registers callback when item is removed
	OnItemRemoved(callback ItemCallback)

	// OnItemChanged registers callback when item stack changes
	OnItemChanged(callback ItemCallback)

	// --- Persistence ---

	// SerializeState converts state to map for persistence
	SerializeState() (map[string]any, error)

	// DeserializeState restores state from map
	DeserializeState(state map[string]any) error
}

// SortBy defines how to sort inventory
type SortBy string

const (
	SortByName   SortBy = "name"
	SortByType   SortBy = "type"
	SortByRarity SortBy = "rarity"
	SortByLevel  SortBy = "level"
	SortByWeight SortBy = "weight"
	SortByValue  SortBy = "value"
	SortByStack  SortBy = "stack"
)

var _ Manager = (*BaseManager)(nil)

// BaseManager implements Manager interface
type BaseManager struct {
	mu sync.RWMutex

	slots     []item.Item    // slot index -> item (nil = empty)
	itemIndex map[string]int // itemID -> slot index
	maxSlots  int
	maxWeight float64

	currentWeight float64

	onAddedCallbacks   []ItemCallback
	onRemovedCallbacks []ItemCallback
	onChangedCallbacks []ItemCallback
}

// Config holds configuration for creating an inventory manager
type Config struct {
	MaxSlots  int
	MaxWeight float64
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		MaxSlots:  20,
		MaxWeight: 100.0,
	}
}

// NewManager creates a new inventory manager with default config
func NewManager() *BaseManager {
	return NewManagerWithConfig(DefaultConfig())
}

// NewManagerWithConfig creates a new inventory manager with config
func NewManagerWithConfig(cfg Config) *BaseManager {
	maxSlots := cfg.MaxSlots
	if maxSlots <= 0 {
		maxSlots = 20
	}

	maxWeight := cfg.MaxWeight
	if maxWeight <= 0 {
		maxWeight = 100.0
	}

	return &BaseManager{
		slots:     make([]item.Item, maxSlots),
		itemIndex: make(map[string]int),
		maxSlots:  maxSlots,
		maxWeight: maxWeight,
	}
}

// --- Basic Operations ---

func (m *BaseManager) Add(ctx context.Context, itm item.Item) error {
	if itm == nil {
		return fmt.Errorf("cannot add nil item")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Try to stack with existing item first
	if targetID, canStack := m.canStackWithLocked(itm); canStack {
		return m.mergeIntoExistingLocked(ctx, itm, targetID)
	}

	// Find free slot
	slot := m.findFreeSlotLocked()
	if slot == -1 {
		return fmt.Errorf("inventory is full (no free slots)")
	}

	// Check weight
	itemWeight := m.getItemWeight(itm)
	if m.currentWeight+itemWeight > m.maxWeight {
		return fmt.Errorf("inventory weight limit exceeded (current: %.2f, max: %.2f, item: %.2f)",
			m.currentWeight, m.maxWeight, itemWeight)
	}

	return m.addToSlotLocked(ctx, slot, itm)
}

func (m *BaseManager) AddToSlot(ctx context.Context, slot int, itm item.Item) error {
	if itm == nil {
		return fmt.Errorf("cannot add nil item")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if slot < 0 || slot >= m.maxSlots {
		return fmt.Errorf("slot %d out of range (0-%d)", slot, m.maxSlots-1)
	}

	if m.slots[slot] != nil {
		return fmt.Errorf("slot %d is already occupied", slot)
	}

	itemWeight := m.getItemWeight(itm)
	if m.currentWeight+itemWeight > m.maxWeight {
		return fmt.Errorf("inventory weight limit exceeded")
	}

	return m.addToSlotLocked(ctx, slot, itm)
}

func (m *BaseManager) addToSlotLocked(ctx context.Context, slot int, itm item.Item) error {
	m.slots[slot] = itm
	m.itemIndex[itm.ID()] = slot
	m.currentWeight += m.getItemWeight(itm)

	// Trigger callbacks (copy to avoid holding lock)
	callbacks := append([]ItemCallback{}, m.onAddedCallbacks...)
	m.mu.Unlock()
	for _, cb := range callbacks {
		cb(ctx, itm)
	}
	m.mu.Lock()

	return nil
}

func (m *BaseManager) Remove(ctx context.Context, itemID string) (item.Item, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	slot, exists := m.itemIndex[itemID]
	if !exists {
		return nil, fmt.Errorf("item with ID %s not found", itemID)
	}

	itm := m.slots[slot]
	m.slots[slot] = nil
	delete(m.itemIndex, itemID)
	m.currentWeight -= m.getItemWeight(itm)
	if m.currentWeight < 0 {
		m.currentWeight = 0
	}

	callbacks := append([]ItemCallback{}, m.onRemovedCallbacks...)
	m.mu.Unlock()
	for _, cb := range callbacks {
		cb(ctx, itm)
	}
	m.mu.Lock()

	return itm, nil
}

func (m *BaseManager) RemoveAmount(ctx context.Context, itemID string, amount int) (item.Item, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	slot, exists := m.itemIndex[itemID]
	if !exists {
		return nil, fmt.Errorf("item with ID %s not found", itemID)
	}

	itm := m.slots[slot]
	currentStack := itm.StackSize()

	if amount >= currentStack {
		// Remove entire item
		m.slots[slot] = nil
		delete(m.itemIndex, itemID)
		m.currentWeight -= m.getItemWeight(itm)
		if m.currentWeight < 0 {
			m.currentWeight = 0
		}

		callbacks := append([]ItemCallback{}, m.onRemovedCallbacks...)
		m.mu.Unlock()
		for _, cb := range callbacks {
			cb(ctx, itm)
		}
		m.mu.Lock()

		return itm, nil
	}

	// Reduce stack size
	oldWeight := m.getItemWeight(itm)
	itm.RemoveStack(amount)
	newWeight := m.getItemWeight(itm)
	m.currentWeight -= oldWeight - newWeight

	// Create new item for removed portion
	removed := itm.Clone().(item.Item)
	// Set stack size on clone - this depends on item implementation
	// For now we return original with reduced count

	callbacks := append([]ItemCallback{}, m.onChangedCallbacks...)
	m.mu.Unlock()
	for _, cb := range callbacks {
		cb(ctx, itm)
	}
	m.mu.Lock()

	return removed, nil
}

func (m *BaseManager) Get(itemID string) (item.Item, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	slot, exists := m.itemIndex[itemID]
	if !exists {
		return nil, false
	}
	return m.slots[slot], true
}

func (m *BaseManager) GetAtSlot(slot int) (item.Item, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if slot < 0 || slot >= len(m.slots) {
		return nil, false
	}
	itm := m.slots[slot]
	return itm, itm != nil
}

func (m *BaseManager) GetAll() []item.Item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]item.Item, 0, len(m.itemIndex))
	for _, itm := range m.slots {
		if itm != nil {
			result = append(result, itm)
		}
	}
	return result
}

func (m *BaseManager) Clear(ctx context.Context) []item.Item {
	m.mu.Lock()

	items := make([]item.Item, 0, len(m.itemIndex))
	for _, itm := range m.slots {
		if itm != nil {
			items = append(items, itm)
		}
	}

	m.slots = make([]item.Item, m.maxSlots)
	m.itemIndex = make(map[string]int)
	m.currentWeight = 0

	callbacks := append([]ItemCallback{}, m.onRemovedCallbacks...)
	m.mu.Unlock()

	for _, itm := range items {
		for _, cb := range callbacks {
			cb(ctx, itm)
		}
	}

	return items
}

// --- Stack Operations ---

func (m *BaseManager) SplitStack(ctx context.Context, itemID string, amount int) (item.Item, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	slot, exists := m.itemIndex[itemID]
	if !exists {
		return nil, fmt.Errorf("item with ID %s not found", itemID)
	}

	itm := m.slots[slot]
	if itm.StackSize() <= amount {
		return nil, fmt.Errorf("cannot split: stack size %d is not greater than %d", itm.StackSize(), amount)
	}

	// Find free slot for new stack
	newSlot := m.findFreeSlotLocked()
	if newSlot == -1 {
		return nil, fmt.Errorf("no free slot for split stack")
	}

	// Remove from original stack
	itm.RemoveStack(amount)

	// Create new item (clone and set stack)
	newItem := itm.Clone().(item.Item)
	newItem.RemoveStack(newItem.StackSize() - 1) // Reset to 1
	newItem.AddStack(amount - 1)                 // Set to amount

	// Generate new ID for split item
	if setter, ok := newItem.(interface{ SetID(string) }); ok {
		setter.SetID(identifier.New())
	}

	m.slots[newSlot] = newItem
	m.itemIndex[newItem.ID()] = newSlot

	// Weight doesn't change on split

	callbacks := append([]ItemCallback{}, m.onChangedCallbacks...)
	addCallbacks := append([]ItemCallback{}, m.onAddedCallbacks...)
	m.mu.Unlock()

	for _, cb := range callbacks {
		cb(ctx, itm)
	}
	for _, cb := range addCallbacks {
		cb(ctx, newItem)
	}

	m.mu.Lock()
	return newItem, nil
}

func (m *BaseManager) MergeStacks(ctx context.Context, sourceID, targetID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sourceSlot, sourceExists := m.itemIndex[sourceID]
	targetSlot, targetExists := m.itemIndex[targetID]

	if !sourceExists {
		return fmt.Errorf("source item %s not found", sourceID)
	}
	if !targetExists {
		return fmt.Errorf("target item %s not found", targetID)
	}

	source := m.slots[sourceSlot]
	target := m.slots[targetSlot]

	if !target.CanStackWith(source) {
		return fmt.Errorf("items cannot be stacked together")
	}

	availableSpace := target.MaxStackSize() - target.StackSize()
	if availableSpace <= 0 {
		return fmt.Errorf("target stack is full")
	}

	amountToMove := source.StackSize()
	if amountToMove > availableSpace {
		amountToMove = availableSpace
	}

	target.AddStack(amountToMove)
	source.RemoveStack(amountToMove)

	if source.StackSize() <= 0 {
		m.slots[sourceSlot] = nil
		delete(m.itemIndex, sourceID)
	}

	callbacks := append([]ItemCallback{}, m.onChangedCallbacks...)
	m.mu.Unlock()

	for _, cb := range callbacks {
		cb(ctx, target)
	}

	m.mu.Lock()
	return nil
}

func (m *BaseManager) CanStackWith(itm item.Item) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.canStackWithLocked(itm)
}

func (m *BaseManager) canStackWithLocked(itm item.Item) (string, bool) {
	for _, existing := range m.slots {
		if existing != nil && existing.CanStackWith(itm) {
			if existing.StackSize() < existing.MaxStackSize() {
				return existing.ID(), true
			}
		}
	}
	return "", false
}

func (m *BaseManager) mergeIntoExistingLocked(ctx context.Context, itm item.Item, targetID string) error {
	targetSlot := m.itemIndex[targetID]
	target := m.slots[targetSlot]

	availableSpace := target.MaxStackSize() - target.StackSize()
	amountToAdd := itm.StackSize()

	if amountToAdd <= availableSpace {
		// All fits in existing stack
		oldWeight := m.getItemWeight(target)
		target.AddStack(amountToAdd)
		newWeight := m.getItemWeight(target)
		m.currentWeight += newWeight - oldWeight

		callbacks := append([]ItemCallback{}, m.onChangedCallbacks...)
		m.mu.Unlock()
		for _, cb := range callbacks {
			cb(ctx, target)
		}
		m.mu.Lock()
		return nil
	}

	// Partial stack - add what fits, then add remainder as new item
	oldWeight := m.getItemWeight(target)
	target.AddStack(availableSpace)
	newWeight := m.getItemWeight(target)
	m.currentWeight += (newWeight - oldWeight)

	// Remainder needs new slot
	slot := m.findFreeSlotLocked()
	if slot == -1 {
		return fmt.Errorf("inventory is full")
	}

	itm.RemoveStack(availableSpace)
	return m.addToSlotLocked(ctx, slot, itm)
}

// --- Slot Management ---

func (m *BaseManager) SlotCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.maxSlots
}

func (m *BaseManager) SetSlotCount(count int) {
	if count <= 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if count == m.maxSlots {
		return
	}

	if count > m.maxSlots {
		// Expand
		newSlots := make([]item.Item, count)
		copy(newSlots, m.slots)
		m.slots = newSlots
		m.maxSlots = count
	} else {
		// Shrink - only if trailing slots are empty
		canShrink := true
		for i := count; i < m.maxSlots; i++ {
			if m.slots[i] != nil {
				canShrink = false
				break
			}
		}
		if canShrink {
			m.slots = m.slots[:count]
			m.maxSlots = count
		}
	}
}

func (m *BaseManager) UsedSlots() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.itemIndex)
}

func (m *BaseManager) FreeSlots() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.maxSlots - len(m.itemIndex)
}

func (m *BaseManager) SwapSlots(ctx context.Context, slot1, slot2 int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if slot1 < 0 || slot1 >= m.maxSlots || slot2 < 0 || slot2 >= m.maxSlots {
		return fmt.Errorf("slot out of range")
	}

	item1 := m.slots[slot1]
	item2 := m.slots[slot2]

	m.slots[slot1] = item2
	m.slots[slot2] = item1

	if item1 != nil {
		m.itemIndex[item1.ID()] = slot2
	}
	if item2 != nil {
		m.itemIndex[item2.ID()] = slot1
	}

	return nil
}

func (m *BaseManager) MoveToSlot(ctx context.Context, itemID string, targetSlot int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if targetSlot < 0 || targetSlot >= m.maxSlots {
		return fmt.Errorf("target slot %d out of range", targetSlot)
	}

	currentSlot, exists := m.itemIndex[itemID]
	if !exists {
		return fmt.Errorf("item %s not found", itemID)
	}

	if currentSlot == targetSlot {
		return nil
	}

	if m.slots[targetSlot] != nil {
		return fmt.Errorf("target slot %d is occupied", targetSlot)
	}

	itm := m.slots[currentSlot]
	m.slots[currentSlot] = nil
	m.slots[targetSlot] = itm
	m.itemIndex[itemID] = targetSlot

	return nil
}

func (m *BaseManager) findFreeSlotLocked() int {
	for i, itm := range m.slots {
		if itm == nil {
			return i
		}
	}
	return -1
}

// --- Weight Management ---

func (m *BaseManager) CurrentWeight() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentWeight
}

func (m *BaseManager) MaxWeight() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.maxWeight
}

func (m *BaseManager) SetMaxWeight(weight float64) {
	if weight <= 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.maxWeight = weight
}

func (m *BaseManager) AvailableWeight() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	available := m.maxWeight - m.currentWeight
	if available < 0 {
		return 0
	}
	return available
}

// --- Capacity Checks ---

func (m *BaseManager) CanAdd(itm item.Item) bool {
	if itm == nil {
		return false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Can stack?
	if _, canStack := m.canStackWithLocked(itm); canStack {
		return true
	}

	// Free slot?
	if m.findFreeSlotLocked() == -1 {
		return false
	}

	// Weight check
	return m.currentWeight+m.getItemWeight(itm) <= m.maxWeight
}

func (m *BaseManager) IsFull() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.findFreeSlotLocked() == -1
}

func (m *BaseManager) Contains(itemID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.itemIndex[itemID]
	return exists
}

func (m *BaseManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.itemIndex)
}

func (m *BaseManager) TotalItems() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var total int
	for _, itm := range m.slots {
		if itm != nil {
			total += itm.StackSize()
		}
	}
	return total
}

// --- Search & Filter ---

func (m *BaseManager) Search(query string) []item.Item {
	query = strings.ToLower(query)
	return m.Filter(func(itm item.Item) bool {
		return strings.Contains(strings.ToLower(itm.Name()), query)
	})
}

func (m *BaseManager) FindByType(itemType item.Type) []item.Item {
	return m.Filter(func(itm item.Item) bool {
		return itm.ItemType() == itemType
	})
}

func (m *BaseManager) FindByRarity(rarity item.Rarity) []item.Item {
	return m.Filter(func(itm item.Item) bool {
		return itm.Rarity() == rarity
	})
}

func (m *BaseManager) FindByTag(tag string) []item.Item {
	return m.Filter(func(itm item.Item) bool {
		return itm.Tags().Has(tag)
	})
}

func (m *BaseManager) FindByTags(tags ...string) []item.Item {
	return m.Filter(func(itm item.Item) bool {
		return itm.Tags().Contains(tags...)
	})
}

func (m *BaseManager) FindByAnyTag(tags ...string) []item.Item {
	return m.Filter(func(itm item.Item) bool {
		return itm.Tags().ContainsAny(tags...)
	})
}

func (m *BaseManager) FindByLevel(minLevel, maxLevel int) []item.Item {
	return m.Filter(func(itm item.Item) bool {
		level := itm.Level()
		return level >= minLevel && level <= maxLevel
	})
}

func (m *BaseManager) FindStackable() []item.Item {
	return m.Filter(func(itm item.Item) bool {
		return itm.MaxStackSize() > 1
	})
}

func (m *BaseManager) Filter(predicate func(item.Item) bool) []item.Item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []item.Item
	for _, itm := range m.slots {
		if itm != nil && predicate(itm) {
			result = append(result, itm)
		}
	}
	return result
}

// --- Sorting ---

func (m *BaseManager) Sort(criteria SortBy, ascending bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	items := make([]item.Item, 0, len(m.itemIndex))
	for _, itm := range m.slots {
		if itm != nil {
			items = append(items, itm)
		}
	}

	m.sortItems(items, criteria, ascending)

	// Rebuild slots
	m.slots = make([]item.Item, m.maxSlots)
	m.itemIndex = make(map[string]int)

	for i, itm := range items {
		m.slots[i] = itm
		m.itemIndex[itm.ID()] = i
	}
}

func (m *BaseManager) GetSorted(criteria SortBy, ascending bool) []item.Item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	items := make([]item.Item, 0, len(m.itemIndex))
	for _, itm := range m.slots {
		if itm != nil {
			items = append(items, itm)
		}
	}

	m.sortItems(items, criteria, ascending)
	return items
}

func (m *BaseManager) sortItems(items []item.Item, criteria SortBy, ascending bool) {
	sort.Slice(items, func(i, j int) bool {
		var less bool
		switch criteria {
		case SortByName:
			less = items[i].Name() < items[j].Name()
		case SortByType:
			less = items[i].ItemType() < items[j].ItemType()
		case SortByRarity:
			less = items[i].Rarity() < items[j].Rarity()
		case SortByLevel:
			less = items[i].Level() < items[j].Level()
		case SortByWeight:
			less = items[i].Weight() < items[j].Weight()
		case SortByValue:
			less = items[i].Value() < items[j].Value()
		case SortByStack:
			less = items[i].StackSize() < items[j].StackSize()
		default:
			less = items[i].Name() < items[j].Name()
		}
		if !ascending {
			less = !less
		}
		return less
	})
}

// --- Stats ---

func (m *BaseManager) TotalValue() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var total int64
	for _, itm := range m.slots {
		if itm != nil {
			total += itm.Value() * int64(itm.StackSize())
		}
	}
	return total
}

func (m *BaseManager) TotalWeight() float64 {
	return m.CurrentWeight()
}

func (m *BaseManager) WeightPercent() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.maxWeight <= 0 {
		return 1.0
	}
	return m.currentWeight / m.maxWeight
}

func (m *BaseManager) SlotPercent() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.maxSlots <= 0 {
		return 1.0
	}
	return float64(len(m.itemIndex)) / float64(m.maxSlots)
}

// --- Callbacks ---

func (m *BaseManager) OnItemAdded(callback ItemCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onAddedCallbacks = append(m.onAddedCallbacks, callback)
}

func (m *BaseManager) OnItemRemoved(callback ItemCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onRemovedCallbacks = append(m.onRemovedCallbacks, callback)
}

func (m *BaseManager) OnItemChanged(callback ItemCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChangedCallbacks = append(m.onChangedCallbacks, callback)
}

// --- Persistence ---

// State holds serializable inventory state
type State struct {
	ItemIDs   []string `msgpack:"item_ids"`
	MaxSlots  int      `msgpack:"max_slots"`
	MaxWeight float64  `msgpack:"max_weight"`
}

func (m *BaseManager) SerializeState() (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	itemIDs := make([]string, m.maxSlots)
	for i, itm := range m.slots {
		if itm != nil {
			itemIDs[i] = itm.ID()
		}
	}

	state := State{
		ItemIDs:   itemIDs,
		MaxSlots:  m.maxSlots,
		MaxWeight: m.maxWeight,
	}

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

	m.mu.Lock()
	defer m.mu.Unlock()

	m.maxSlots = state.MaxSlots
	if m.maxSlots <= 0 {
		m.maxSlots = 20
	}

	m.maxWeight = state.MaxWeight
	if m.maxWeight <= 0 {
		m.maxWeight = 100.0
	}

	m.slots = make([]item.Item, m.maxSlots)
	m.itemIndex = make(map[string]int)
	m.currentWeight = 0

	return nil
}

// --- Helper methods for persistence ---

// GetItemIDs returns all item IDs in slot order
func (m *BaseManager) GetItemIDs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, m.maxSlots)
	for i, itm := range m.slots {
		if itm != nil {
			ids[i] = itm.ID()
		}
	}
	return ids
}

// AddDirect adds an item without callbacks (for deserialization)
func (m *BaseManager) AddDirect(itm item.Item) error {
	if itm == nil {
		return fmt.Errorf("cannot add nil item")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	slot := m.findFreeSlotLocked()
	if slot == -1 {
		return fmt.Errorf("no free slot")
	}

	m.slots[slot] = itm
	m.itemIndex[itm.ID()] = slot
	m.currentWeight += m.getItemWeight(itm)

	return nil
}

// AddDirectToSlot adds an item to specific slot without callbacks
func (m *BaseManager) AddDirectToSlot(slot int, itm item.Item) error {
	if itm == nil {
		return fmt.Errorf("cannot add nil item")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if slot < 0 || slot >= m.maxSlots {
		return fmt.Errorf("slot out of range")
	}

	m.slots[slot] = itm
	m.itemIndex[itm.ID()] = slot
	m.currentWeight += m.getItemWeight(itm)

	return nil
}

// RecalculateWeight recalculates weight from items
func (m *BaseManager) RecalculateWeight() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.currentWeight = 0
	for _, itm := range m.slots {
		if itm != nil {
			m.currentWeight += m.getItemWeight(itm)
		}
	}
}

func (m *BaseManager) getItemWeight(itm item.Item) float64 {
	return itm.Weight() * float64(itm.StackSize())
}
