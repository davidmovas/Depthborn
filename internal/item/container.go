package item

import (
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/pkg/persist"
)

var _ Container = (*BaseContainer)(nil)

// BaseContainer implements Container interface for items that hold other items
type BaseContainer struct {
	*BaseItem

	mu           sync.RWMutex
	capacity     int
	contents     []Item
	maxWeight    float64 // Maximum weight the container can hold (0 = unlimited)
	allowedTypes []Type  // Allowed item types (empty = all types allowed)
}

// ContainerConfig holds configuration for creating a BaseContainer
type ContainerConfig struct {
	BaseItemConfig
	Capacity     int
	MaxWeight    float64
	AllowedTypes []Type
}

// NewBaseContainer creates a new container with minimal configuration
func NewBaseContainer(id string, name string, capacity int) *BaseContainer {
	return NewBaseContainerWithConfig(ContainerConfig{
		BaseItemConfig: BaseItemConfig{
			ID:       id,
			Name:     name,
			ItemType: TypeContainer,
		},
		Capacity: capacity,
	})
}

// NewBaseContainerWithConfig creates a new container with full configuration
func NewBaseContainerWithConfig(cfg ContainerConfig) *BaseContainer {
	// Ensure item type is container
	cfg.BaseItemConfig.ItemType = TypeContainer

	bc := &BaseContainer{
		BaseItem:     NewBaseItemWithConfig(cfg.BaseItemConfig),
		capacity:     cfg.Capacity,
		contents:     make([]Item, 0),
		maxWeight:    cfg.MaxWeight,
		allowedTypes: cfg.AllowedTypes,
	}

	if bc.capacity < 1 {
		bc.capacity = 1
	}

	if bc.allowedTypes == nil {
		bc.allowedTypes = make([]Type, 0)
	}

	return bc
}

// --- Container interface implementation ---

func (bc *BaseContainer) Capacity() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.capacity
}

func (bc *BaseContainer) SetCapacity(capacity int) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	if capacity < 1 {
		capacity = 1
	}
	bc.capacity = capacity
	bc.Touch()
}

func (bc *BaseContainer) Contents() []Item {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	// Return a copy to prevent external modification
	result := make([]Item, len(bc.contents))
	copy(result, bc.contents)
	return result
}

func (bc *BaseContainer) Add(item Item) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if err := bc.canAddInternal(item); err != nil {
		return err
	}

	bc.contents = append(bc.contents, item)
	bc.Touch()
	return nil
}

func (bc *BaseContainer) Remove(itemID string) (Item, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	for i, item := range bc.contents {
		if item.ID() == itemID {
			bc.contents = append(bc.contents[:i], bc.contents[i+1:]...)
			bc.Touch()
			return item, nil
		}
	}
	return nil, fmt.Errorf("item not found: %s", itemID)
}

func (bc *BaseContainer) Contains(itemID string) bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	for _, item := range bc.contents {
		if item.ID() == itemID {
			return true
		}
	}
	return false
}

func (bc *BaseContainer) IsFull() bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return len(bc.contents) >= bc.capacity
}

// --- Additional methods ---

// CanAdd checks if an item can be added to the container
func (bc *BaseContainer) CanAdd(item Item) error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.canAddInternal(item)
}

// canAddInternal checks if an item can be added (no lock)
func (bc *BaseContainer) canAddInternal(item Item) error {
	if item == nil {
		return fmt.Errorf("cannot add nil item")
	}

	// Check capacity
	if len(bc.contents) >= bc.capacity {
		return fmt.Errorf("container is full")
	}

	// Check if item is already in container
	for _, existing := range bc.contents {
		if existing.ID() == item.ID() {
			return fmt.Errorf("item already in container")
		}
	}

	// Check weight limit
	if bc.maxWeight > 0 {
		currentWeight := bc.contentsWeightInternal()
		if currentWeight+item.Weight() > bc.maxWeight {
			return fmt.Errorf("item too heavy for container")
		}
	}

	// Check allowed types
	if len(bc.allowedTypes) > 0 {
		allowed := false
		for _, t := range bc.allowedTypes {
			if t == item.ItemType() {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("item type %s not allowed in container", item.ItemType())
		}
	}

	return nil
}

// Count returns the number of items in the container
func (bc *BaseContainer) Count() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return len(bc.contents)
}

// RemainingCapacity returns how many more items can fit
func (bc *BaseContainer) RemainingCapacity() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.capacity - len(bc.contents)
}

// Clear removes all items from the container
func (bc *BaseContainer) Clear() []Item {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	removed := bc.contents
	bc.contents = make([]Item, 0)
	bc.Touch()
	return removed
}

// GetItem returns an item by ID without removing it
func (bc *BaseContainer) GetItem(itemID string) (Item, bool) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	for _, item := range bc.contents {
		if item.ID() == itemID {
			return item, true
		}
	}
	return nil, false
}

// FindByType returns all items of a specific type
func (bc *BaseContainer) FindByType(itemType Type) []Item {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	result := make([]Item, 0)
	for _, item := range bc.contents {
		if item.ItemType() == itemType {
			result = append(result, item)
		}
	}
	return result
}

// FindByTag returns all items with a specific tag
func (bc *BaseContainer) FindByTag(tag string) []Item {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	result := make([]Item, 0)
	for _, item := range bc.contents {
		if item.Tags().Has(tag) {
			result = append(result, item)
		}
	}
	return result
}

// MaxWeight returns the maximum weight capacity
func (bc *BaseContainer) MaxWeight() float64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.maxWeight
}

// SetMaxWeight sets the maximum weight capacity
func (bc *BaseContainer) SetMaxWeight(weight float64) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	if weight < 0 {
		weight = 0
	}
	bc.maxWeight = weight
	bc.Touch()
}

// AllowedTypes returns the list of allowed item types
func (bc *BaseContainer) AllowedTypes() []Type {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	result := make([]Type, len(bc.allowedTypes))
	copy(result, bc.allowedTypes)
	return result
}

// SetAllowedTypes sets the allowed item types
func (bc *BaseContainer) SetAllowedTypes(types []Type) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.allowedTypes = make([]Type, len(types))
	copy(bc.allowedTypes, types)
	bc.Touch()
}

// Weight returns total weight (container + contents)
func (bc *BaseContainer) Weight() float64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.BaseItem.Weight() + bc.contentsWeightInternal()
}

// ContentsWeight returns weight of contents only
func (bc *BaseContainer) ContentsWeight() float64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.contentsWeightInternal()
}

// contentsWeightInternal calculates contents weight (no lock)
func (bc *BaseContainer) contentsWeightInternal() float64 {
	var totalWeight float64
	for _, item := range bc.contents {
		totalWeight += item.Weight()
	}
	return totalWeight
}

// TotalValue returns total value (container + contents)
func (bc *BaseContainer) TotalValue() int64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	total := bc.BaseItem.Value()
	for _, item := range bc.contents {
		total += item.Value()
	}
	return total
}

// --- Cloneable interface ---

func (bc *BaseContainer) Clone() any {
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

	// Clone contents
	clonedContents := make([]Item, 0, len(bc.contents))
	for _, item := range bc.contents {
		if cloneable, ok := item.(interface{ Clone() any }); ok {
			if cloned, ok := cloneable.Clone().(Item); ok {
				clonedContents = append(clonedContents, cloned)
			}
		}
	}

	// Clone allowed types
	allowedTypes := make([]Type, len(bc.allowedTypes))
	copy(allowedTypes, bc.allowedTypes)

	clone := &BaseContainer{
		BaseItem:     baseClone,
		capacity:     bc.capacity,
		contents:     clonedContents,
		maxWeight:    bc.maxWeight,
		allowedTypes: allowedTypes,
	}

	return clone
}

// --- Serialization ---

// ContainerState holds the complete serializable state of a BaseContainer
type ContainerState struct {
	State
	Capacity     int      `msgpack:"capacity"`
	MaxWeight    float64  `msgpack:"max_weight"`
	AllowedTypes []string `msgpack:"allowed_types"`
	ContentIDs   []string `msgpack:"content_ids"` // IDs of contained items
}

func (bc *BaseContainer) Marshal() ([]byte, error) {
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

	// Collect content IDs
	contentIDs := make([]string, len(bc.contents))
	for i, item := range bc.contents {
		contentIDs[i] = item.ID()
	}

	// Convert allowed types to strings
	allowedTypes := make([]string, len(bc.allowedTypes))
	for i, t := range bc.allowedTypes {
		allowedTypes[i] = string(t)
	}

	cs := ContainerState{
		State:        is,
		Capacity:     bc.capacity,
		MaxWeight:    bc.maxWeight,
		AllowedTypes: allowedTypes,
		ContentIDs:   contentIDs,
	}

	return persist.DefaultCodec().Encode(cs)
}

func (bc *BaseContainer) Unmarshal(data []byte) error {
	var cs ContainerState
	if err := persist.DefaultCodec().Decode(data, &cs); err != nil {
		return fmt.Errorf("failed to decode container state: %w", err)
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

	// Restore container-specific fields
	bc.capacity = cs.Capacity
	bc.maxWeight = cs.MaxWeight

	// Convert allowed types from strings
	bc.allowedTypes = make([]Type, len(cs.AllowedTypes))
	for i, t := range cs.AllowedTypes {
		bc.allowedTypes[i] = Type(t)
	}

	// Note: contents must be restored separately using ContentIDs
	// via an item repository
	bc.contents = make([]Item, 0)

	return nil
}

// ContentIDs returns the IDs of contained items (for serialization)
func (bc *BaseContainer) ContentIDs() []string {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	ids := make([]string, len(bc.contents))
	for i, item := range bc.contents {
		ids[i] = item.ID()
	}
	return ids
}

// RestoreContents restores contents from loaded items
func (bc *BaseContainer) RestoreContents(items []Item) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.contents = items
}

// --- Validatable interface ---

func (bc *BaseContainer) Validate() error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	// Validate base item
	if err := bc.BaseItem.Validate(); err != nil {
		return err
	}

	if bc.capacity < 1 {
		return fmt.Errorf("capacity must be at least 1")
	}

	if len(bc.contents) > bc.capacity {
		return fmt.Errorf("contents exceed capacity")
	}

	if bc.maxWeight < 0 {
		return fmt.Errorf("max weight cannot be negative")
	}

	// Validate weight limit
	if bc.maxWeight > 0 && bc.contentsWeightInternal() > bc.maxWeight {
		return fmt.Errorf("contents weight exceeds max weight")
	}

	return nil
}

// --- Serializable interface (required by infra.Persistent) ---

func (bc *BaseContainer) SerializeState() (map[string]any, error) {
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

func (bc *BaseContainer) DeserializeState(state map[string]any) error {
	data, err := persist.DefaultCodec().Encode(state)
	if err != nil {
		return err
	}
	return bc.Unmarshal(data)
}
