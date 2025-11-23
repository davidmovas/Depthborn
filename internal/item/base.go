package item

import (
	"github.com/davidmovas/Depthborn/internal/infra"
)

// BaseItem implements core Item functionality
type BaseItem struct {
	infra.BasePersistent
	infra.BaseCloneable

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
	tags         []string
}

func NewBaseItem(id string, itemType Type, name string) *BaseItem {
	bi := &BaseItem{
		name:         name,
		description:  "",
		itemType:     itemType,
		rarity:       RarityCommon,
		quality:      1.0,
		level:        1,
		stackSize:    1,
		maxStackSize: 1,
		value:        0,
		weight:       0.1,
		icon:         "default",
		tags:         make([]string, 0),
	}
	bi.SetID(id)
	return bi
}

func (bi *BaseItem) Name() string {
	return bi.name
}

func (bi *BaseItem) SetName(name string) {
	bi.name = name
}

func (bi *BaseItem) Description() string {
	return bi.description
}

func (bi *BaseItem) ItemType() Type {
	return bi.itemType
}

func (bi *BaseItem) Rarity() Rarity {
	return bi.rarity
}

func (bi *BaseItem) Quality() float64 {
	return bi.quality
}

func (bi *BaseItem) SetQuality(quality float64) {
	if quality < 0.0 {
		quality = 0.0
	} else if quality > 1.0 {
		quality = 1.0
	}
	bi.quality = quality
}

func (bi *BaseItem) Level() int {
	return bi.level
}

func (bi *BaseItem) SetLevel(level int) {
	if level < 1 {
		level = 1
	}
	bi.level = level
}

func (bi *BaseItem) StackSize() int {
	return bi.stackSize
}

func (bi *BaseItem) MaxStackSize() int {
	return bi.maxStackSize
}

func (bi *BaseItem) AddStack(amount int) bool {
	if bi.stackSize+amount > bi.maxStackSize {
		return false
	}
	bi.stackSize += amount
	return true
}

func (bi *BaseItem) RemoveStack(amount int) int {
	if amount >= bi.stackSize {
		remaining := bi.stackSize
		bi.stackSize = 0
		return remaining
	}
	bi.stackSize -= amount
	return amount
}

func (bi *BaseItem) CanStackWith(other Item) bool {
	if bi.itemType != other.ItemType() {
		return false
	}
	if bi.ID() == other.ID() {
		return false
	}
	return bi.stackSize < bi.maxStackSize
}

func (bi *BaseItem) Value() int64 {
	return bi.value
}

func (bi *BaseItem) Weight() float64 {
	return bi.weight
}

func (bi *BaseItem) Icon() string {
	return bi.icon
}

func (bi *BaseItem) Tags() []string {
	return bi.tags
}

func (bi *BaseItem) IsEquippable() bool {
	switch bi.itemType {
	case TypeWeaponMelee, TypeWeaponRanged, TypeWeaponMagic,
		TypeArmorHead, TypeArmorChest, TypeArmorLegs, TypeArmorFeet, TypeArmorHands,
		TypeAccessoryRing, TypeAccessoryAmulet, TypeAccessoryBelt:
		return true
	default:
		return false
	}
}

func (bi *BaseItem) IsConsumable() bool {
	return bi.itemType == TypeConsumable
}

func (bi *BaseItem) IsQuestItem() bool {
	return bi.itemType == TypeQuest
}

func (bi *BaseItem) IsTradeable() bool {
	return bi.itemType != TypeQuest
}
