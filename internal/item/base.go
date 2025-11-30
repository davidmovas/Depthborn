package item

import (
	"fmt"
	"sync"
	"time"

	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/infra/impl"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

// BaseItem implements core Item functionality
type BaseItem struct {
	*impl.BasePersistent

	mu sync.RWMutex

	name         string
	description  string
	itemType     Type
	rarity       Rarity
	quality      float64
	level        int
	stackSize    int
	maxStackSize int
	value        int64
	weight       float64
	icon         string
	tags         types.TagSet
}

// BaseItemConfig holds configuration for creating a BaseItem
type BaseItemConfig struct {
	ID           string
	Name         string
	Description  string
	ItemType     Type
	Rarity       Rarity
	Quality      float64
	Level        int
	MaxStackSize int
	Value        int64
	Weight       float64
	Icon         string
	Tags         []string
}

// NewBaseItem creates a new base item with minimal configuration
func NewBaseItem(id string, itemType Type, name string) *BaseItem {
	return NewBaseItemWithConfig(BaseItemConfig{
		ID:       id,
		Name:     name,
		ItemType: itemType,
	})
}

// NewBaseItemWithConfig creates a new base item with full configuration
func NewBaseItemWithConfig(cfg BaseItemConfig) *BaseItem {
	bi := &BaseItem{
		name:         cfg.Name,
		description:  cfg.Description,
		itemType:     cfg.ItemType,
		rarity:       cfg.Rarity,
		quality:      cfg.Quality,
		level:        cfg.Level,
		stackSize:    1,
		maxStackSize: cfg.MaxStackSize,
		value:        cfg.Value,
		weight:       cfg.Weight,
		icon:         cfg.Icon,
		tags:         types.NewTagSet(),
	}

	// Apply defaults
	if bi.quality <= 0 {
		bi.quality = 1.0
	}
	if bi.level < 1 {
		bi.level = 1
	}
	if bi.maxStackSize < 1 {
		bi.maxStackSize = 1
	}
	if bi.weight <= 0 {
		bi.weight = 0.1
	}
	if bi.icon == "" {
		bi.icon = "default"
	}

	// Add initial tags from config
	for _, tag := range cfg.Tags {
		bi.tags.Add(tag)
	}

	// Create persistence layer
	entityType := "item:" + string(cfg.ItemType)
	if cfg.ID != "" {
		bi.BasePersistent = impl.NewPersistentWithID(cfg.ID, entityType, bi, nil)
	} else {
		bi.BasePersistent = impl.NewPersistent(entityType, bi, nil)
	}

	return bi
}

// --- Item interface implementation ---

func (i *BaseItem) Name() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.name
}

func (i *BaseItem) SetName(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.name = name
	i.Touch()
}

func (i *BaseItem) Description() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.description
}

func (i *BaseItem) SetDescription(desc string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.description = desc
	i.Touch()
}

func (i *BaseItem) ItemType() Type {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.itemType
}

func (i *BaseItem) Rarity() Rarity {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.rarity
}

func (i *BaseItem) SetRarity(rarity Rarity) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.rarity = rarity
	i.Touch()
}

func (i *BaseItem) Quality() float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.quality
}

func (i *BaseItem) SetQuality(quality float64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if quality < 0.0 {
		quality = 0.0
	} else if quality > 1.0 {
		quality = 1.0
	}
	i.quality = quality
	i.Touch()
}

func (i *BaseItem) Level() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.level
}

func (i *BaseItem) SetLevel(level int) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if level < 1 {
		level = 1
	}
	i.level = level
	i.Touch()
}

func (i *BaseItem) StackSize() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.stackSize
}

// StackSizeInternal returns stack size without locking (for use by embedded types)
func (i *BaseItem) StackSizeInternal() int {
	return i.stackSize
}

func (i *BaseItem) MaxStackSize() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.maxStackSize
}

func (i *BaseItem) SetMaxStackSize(max int) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if max < 1 {
		max = 1
	}
	i.maxStackSize = max
	if i.stackSize > i.maxStackSize {
		i.stackSize = i.maxStackSize
	}
	i.Touch()
}

func (i *BaseItem) AddStack(amount int) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	if amount <= 0 {
		return false
	}
	if i.stackSize+amount > i.maxStackSize {
		return false
	}
	i.stackSize += amount
	i.Touch()
	return true
}

func (i *BaseItem) RemoveStack(amount int) int {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.removeStackInternal(amount)
}

// removeStackInternal removes from stack without locking (for use by embedded types)
func (i *BaseItem) removeStackInternal(amount int) int {
	if amount <= 0 {
		return 0
	}
	if amount >= i.stackSize {
		removed := i.stackSize
		i.stackSize = 0
		i.Touch()
		return removed
	}
	i.stackSize -= amount
	i.Touch()
	return amount
}

// RemoveStackInternal removes from stack without locking (for use by embedded types)
func (i *BaseItem) RemoveStackInternal(amount int) int {
	return i.removeStackInternal(amount)
}

func (i *BaseItem) CanStackWith(other Item) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Can't stack with self
	if i.ID() == other.ID() {
		return false
	}
	// Must be same item type
	if i.itemType != other.ItemType() {
		return false
	}
	// Must be same rarity for stacking
	if i.rarity != other.Rarity() {
		return false
	}
	// Check if we have room
	return i.stackSize < i.maxStackSize
}

func (i *BaseItem) Value() int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.value
}

func (i *BaseItem) SetValue(value int64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if value < 0 {
		value = 0
	}
	i.value = value
	i.Touch()
}

func (i *BaseItem) Weight() float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.weight
}

func (i *BaseItem) SetWeight(weight float64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if weight < 0 {
		weight = 0
	}
	i.weight = weight
	i.Touch()
}

func (i *BaseItem) Icon() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.icon
}

func (i *BaseItem) SetIcon(icon string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.icon = icon
	i.Touch()
}

func (i *BaseItem) Tags() types.TagSet {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.tags
}

func (i *BaseItem) IsEquippable() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	switch i.itemType {
	case TypeWeaponMelee, TypeWeaponRanged, TypeWeaponMagic,
		TypeArmorHead, TypeArmorChest, TypeArmorLegs, TypeArmorFeet, TypeArmorHands,
		TypeAccessoryRing, TypeAccessoryAmulet, TypeAccessoryBelt:
		return true
	default:
		return false
	}
}

func (i *BaseItem) IsConsumable() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.itemType == TypeConsumable
}

func (i *BaseItem) IsQuestItem() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.itemType == TypeQuest
}

func (i *BaseItem) IsTradeable() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.itemType != TypeQuest
}

// --- Cloneable interface ---

func (i *BaseItem) Clone() any {
	if i == nil {
		return nil
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	// Clone tags
	clonedTags := types.NewTagSet()
	if i.tags != nil {
		for _, tag := range i.tags.All() {
			clonedTags.Add(tag)
		}
	}

	clone := &BaseItem{
		name:         i.name,
		description:  i.description,
		itemType:     i.itemType,
		rarity:       i.rarity,
		quality:      i.quality,
		level:        i.level,
		stackSize:    i.stackSize,
		maxStackSize: i.maxStackSize,
		value:        i.value,
		weight:       i.weight,
		icon:         i.icon,
		tags:         clonedTags,
	}

	// New ID for clone
	entityType := "item:" + string(i.itemType)
	clone.BasePersistent = impl.NewPersistent(entityType, clone, nil)

	return clone
}

// State holds the complete serializable state of a BaseItem
type State struct {
	ID           string   `msgpack:"id"`
	EntityType   string   `msgpack:"entity_type"`
	Version      int64    `msgpack:"version"`
	CreatedAt    int64    `msgpack:"created_at"`
	UpdatedAt    int64    `msgpack:"updated_at"`
	Name         string   `msgpack:"name"`
	Description  string   `msgpack:"description"`
	ItemType     string   `msgpack:"item_type"`
	Rarity       int      `msgpack:"rarity"`
	Quality      float64  `msgpack:"quality"`
	Level        int      `msgpack:"level"`
	StackSize    int      `msgpack:"stack_size"`
	MaxStackSize int      `msgpack:"max_stack_size"`
	Value        int64    `msgpack:"value"`
	Weight       float64  `msgpack:"weight"`
	Icon         string   `msgpack:"icon"`
	Tags         []string `msgpack:"tags"`
}

func (i *BaseItem) Marshal() ([]byte, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	state := State{
		ID:           i.ID(),
		EntityType:   i.Type(),
		Version:      i.Version(),
		CreatedAt:    i.CreatedAt(),
		UpdatedAt:    i.UpdatedAt(),
		Name:         i.name,
		Description:  i.description,
		ItemType:     string(i.itemType),
		Rarity:       int(i.rarity),
		Quality:      i.quality,
		Level:        i.level,
		StackSize:    i.stackSize,
		MaxStackSize: i.maxStackSize,
		Value:        i.value,
		Weight:       i.weight,
		Icon:         i.icon,
		Tags:         i.tags.All(),
	}

	return persist.DefaultCodec().Encode(state)
}

func (i *BaseItem) Unmarshal(data []byte) error {
	var state State
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return fmt.Errorf("failed to decode item state: %w", err)
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	// Restore persistence layer
	i.BasePersistent = impl.NewPersistentWithID(state.ID, state.EntityType, i, nil)

	// Restore item fields
	i.name = state.Name
	i.description = state.Description
	i.itemType = Type(state.ItemType)
	i.rarity = Rarity(state.Rarity)
	i.quality = state.Quality
	i.level = state.Level
	i.stackSize = state.StackSize
	i.maxStackSize = state.MaxStackSize
	i.value = state.Value
	i.weight = state.Weight
	i.icon = state.Icon

	// Restore tags
	i.tags = types.NewTagSet()
	for _, tag := range state.Tags {
		i.tags.Add(tag)
	}

	return nil
}

func (i *BaseItem) Validate() error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.name == "" {
		return fmt.Errorf("item must have a name")
	}
	if i.itemType == "" {
		return fmt.Errorf("item must have a type")
	}
	if i.quality < 0 || i.quality > 1 {
		return fmt.Errorf("quality must be between 0 and 1")
	}
	if i.level < 1 {
		return fmt.Errorf("level must be at least 1")
	}
	if i.stackSize < 0 {
		return fmt.Errorf("stack size cannot be negative")
	}
	if i.maxStackSize < 1 {
		return fmt.Errorf("max stack size must be at least 1")
	}
	if i.stackSize > i.maxStackSize {
		return fmt.Errorf("stack size cannot exceed max stack size")
	}
	if i.weight < 0 {
		return fmt.Errorf("weight cannot be negative")
	}

	return nil
}

// TotalWeight returns weight * stack size
func (i *BaseItem) TotalWeight() float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.weight * float64(i.stackSize)
}

// TotalValue returns value * stack size
func (i *BaseItem) TotalValue() int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.value * int64(i.stackSize)
}

// DisplayName returns formatted name with rarity
func (i *BaseItem) DisplayName() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if i.rarity > RarityCommon {
		return fmt.Sprintf("[%s] %s", i.rarity.String(), i.name)
	}
	return i.name
}

// Age returns how old the item is
func (i *BaseItem) Age() time.Duration {
	return time.Since(time.Unix(i.CreatedAt(), 0))
}

// --- Serializable interface (required by infra.Persistent) ---

func (i *BaseItem) SerializeState() (map[string]any, error) {
	data, err := i.Marshal()
	if err != nil {
		return nil, err
	}

	var state map[string]any
	if err = persist.DefaultCodec().Decode(data, &state); err != nil {
		return nil, err
	}
	return state, nil
}

func (i *BaseItem) DeserializeState(state map[string]any) error {
	data, err := persist.DefaultCodec().Encode(state)
	if err != nil {
		return err
	}
	return i.Unmarshal(data)
}
