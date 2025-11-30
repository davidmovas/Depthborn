package item

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/infra"
	"github.com/davidmovas/Depthborn/internal/item/affix"
)

// Item represents any game item
type Item interface {
	infra.Persistent
	infra.Cloneable

	// Name returns item display name
	Name() string

	// SetName updates item name
	SetName(name string)

	// Description returns item description
	Description() string

	// ItemType returns item type
	ItemType() Type

	// Rarity returns item rarity level
	Rarity() Rarity

	// Quality returns item quality value [0.0 - 1.0]
	Quality() float64

	// SetQuality updates item quality
	SetQuality(quality float64)

	// Level returns item level requirement
	Level() int

	// SetLevel updates item level
	SetLevel(level int)

	// StackSize returns current stack size
	StackSize() int

	// MaxStackSize returns maximum stack size (1 = non-stackable)
	MaxStackSize() int

	// AddStack increases stack size, returns true if successful
	AddStack(amount int) bool

	// RemoveStack decreases stack size, returns remaining amount
	RemoveStack(amount int) int

	// CanStackWith returns true if items can stack together
	CanStackWith(other Item) bool

	// Value returns vendor sell value
	Value() int64

	// Weight returns item weight for inventory management
	Weight() float64

	// Icon returns icon identifier
	Icon() string

	// Tags returns item tag set
	Tags() types.TagSet

	// IsEquippable returns true if item can be equipped
	IsEquippable() bool

	// IsConsumable returns true if item is consumed on use
	IsConsumable() bool

	// IsQuestItem returns true if item is quest-related
	IsQuestItem() bool

	// IsTradeable returns true if item can be traded
	IsTradeable() bool
}

// Type categorizes items
type Type string

const (
	TypeWeaponMelee     Type = "weapon_melee"
	TypeWeaponRanged    Type = "weapon_ranged"
	TypeWeaponMagic     Type = "weapon_magic"
	TypeArmorHead       Type = "armor_head"
	TypeArmorChest      Type = "armor_chest"
	TypeArmorLegs       Type = "armor_legs"
	TypeArmorFeet       Type = "armor_feet"
	TypeArmorHands      Type = "armor_hands"
	TypeAccessoryRing   Type = "accessory_ring"
	TypeAccessoryAmulet Type = "accessory_amulet"
	TypeAccessoryBelt   Type = "accessory_belt"
	TypeConsumable      Type = "consumable"
	TypeMaterial        Type = "material"
	TypeGem             Type = "gem"
	TypeRune            Type = "rune"
	TypeCurrency        Type = "currency"
	TypeQuest           Type = "quest"
	TypeKey             Type = "key"
	TypeContainer       Type = "container"
)

// Rarity defines item rarity tiers
type Rarity int

const (
	RarityCommon Rarity = iota
	RarityUncommon
	RarityRare
	RarityEpic
	RarityLegendary
	RarityMythic
)

// String returns rarity name
func (r Rarity) String() string {
	return [...]string{"Common", "Uncommon", "Rare", "Epic", "Legendary", "Mythic"}[r]
}

// Equipment represents items that can be equipped
type Equipment interface {
	Item

	// Slot returns equipment slot this item occupies
	Slot() EquipmentSlot

	// Attributes returns attribute modifiers provided by item
	Attributes() []attribute.Modifier

	// Durability returns current durability
	Durability() float64

	// MaxDurability returns maximum durability
	MaxDurability() float64

	// SetDurability updates current durability
	SetDurability(value float64)

	// Repair restores durability
	Repair(amount float64)

	// DamageItem reduces durability
	DamageItem(amount float64) bool

	// IsBroken returns true if durability is zero
	IsBroken() bool

	// SocketCount returns number of gem/rune sockets
	SocketCount() int

	// GetSocket returns socketed item at index
	GetSocket(index int) (Socketable, bool)

	// SetSocket places socketable item in socket
	SetSocket(index int, item Socketable) error

	// RemoveSocket removes socketable from socket
	RemoveSocket(index int) (Socketable, error)

	// Affixes returns item affixes
	Affixes() affix.Set

	// Requirements returns equip requirements
	Requirements() EquipRequirements

	// CanEquip checks if entity can equip this item
	CanEquip(entity entity.Entity) bool

	// OnEquip is called when item is equipped
	OnEquip(ctx context.Context, entity entity.Entity) error

	// OnUnequip is called when item is unequipped
	OnUnequip(ctx context.Context, entity entity.Entity) error
}

// EquipmentSlot defines where equipment can be worn
type EquipmentSlot string

const (
	SlotMainHand EquipmentSlot = "main_hand"
	SlotOffHand  EquipmentSlot = "off_hand"
	SlotTwoHand  EquipmentSlot = "two_hand"
	SlotHead     EquipmentSlot = "head"
	SlotChest    EquipmentSlot = "chest"
	SlotLegs     EquipmentSlot = "legs"
	SlotFeet     EquipmentSlot = "feet"
	SlotHands    EquipmentSlot = "hands"
	SlotRing1    EquipmentSlot = "ring_1"
	SlotRing2    EquipmentSlot = "ring_2"
	SlotAmulet   EquipmentSlot = "amulet"
	SlotBelt     EquipmentSlot = "belt"
)

// EquipRequirements defines requirements to equip item
type EquipRequirements interface {
	// Level returns minimum level required
	Level() int

	// Attributes returns minimum attributes required
	Attributes() map[attribute.Type]float64

	// Check verifies if entity meets requirements
	Check(entity entity.Entity) bool
}

// Socketable represents items that can be socketed into equipment
type Socketable interface {
	Item

	// SocketType returns compatible socket type
	SocketType() SocketType

	// Effect returns effect granted when socketed
	Effect() SocketEffect
}

// SocketType defines socket compatibility
type SocketType string

const (
	SocketTypeGem       SocketType = "gem"
	SocketTypeRune      SocketType = "rune"
	SocketTypeUniversal SocketType = "universal" // Accepts any socketable
)

// SocketEffect describes bonus granted by socketed item
type SocketEffect interface {
	// OnSocket is called when the socketable is inserted into equipment
	OnSocket(ctx context.Context, equipment Equipment, entity entity.Entity) error

	// OnUnsocket is called when the socketable is removed from equipment
	OnUnsocket(ctx context.Context, equipment Equipment, entity entity.Entity) error

	// Description returns human-readable description
	Description() string
}

// Consumable represents usable items
type Consumable interface {
	Item

	// Use consumes item and applies effect
	Use(ctx context.Context, user entity.Entity) error

	// CanUse checks if item can be used by entity
	CanUse(user entity.Entity) bool

	// Cooldown returns remaining cooldown in milliseconds
	Cooldown() int64

	// MaxCooldown returns base cooldown duration
	MaxCooldown() int64

	// Effect returns consumable effect
	Effect() ConsumableEffect
}

// ConsumableEffect describes what happens when consumable is used
type ConsumableEffect interface {
	// Apply applies consumable effect to entity
	Apply(ctx context.Context, target entity.Entity) error

	// Description returns effect description
	Description() string

	// Duration returns effect duration in milliseconds (0 = instant)
	Duration() int64
}

// Container represents items that hold other items
type Container interface {
	Item

	// Capacity returns maximum number of items
	Capacity() int

	// Contents returns all items in container
	Contents() []Item

	// Add adds item to container
	Add(item Item) error

	// Remove removes item from container
	Remove(itemID string) (Item, error)

	// Contains checks if container has item
	Contains(itemID string) bool

	// IsFull returns true if container is at capacity
	IsFull() bool

	// Weight returns total weight of container and contents
	Weight() float64
}
