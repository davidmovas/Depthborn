package builder

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/item"
	"github.com/davidmovas/Depthborn/internal/item/affix"
)

// Equipment provides a fluent interface for creating equipment
type Equipment struct {
	*Item
	slot          item.EquipmentSlot
	maxDurability float64
	socketCount   int
	socketTypes   []item.SocketType
	requirements  item.EquipRequirements
	attributes    []attribute.Modifier
	affixes       []affix.Affix
	onEquip       func(ctx context.Context, entity entity.Entity) error
	onUnequip     func(ctx context.Context, entity entity.Entity) error
}

// NewEquipment creates a new equipment builder
func NewEquipment() *Equipment {
	return &Equipment{
		Item:          NewItem(),
		maxDurability: 100.0,
	}
}

// Equip creates an equipment builder with type, name, and slot
func Equip(itemType item.Type, name string, slot item.EquipmentSlot) *Equipment {
	return NewEquipment().
		Type(itemType).
		Name(name).
		Slot(slot)
}

// Weapon creates a weapon builder
func Weapon(name string, slot item.EquipmentSlot) *Equipment {
	var itemType item.Type
	switch slot {
	case item.SlotTwoHand, item.SlotMainHand, item.SlotOffHand:
		itemType = item.TypeWeaponMelee
	default:
		itemType = item.TypeWeaponMelee
	}
	return Equip(itemType, name, slot)
}

// MeleeWeapon creates a melee weapon builder
func MeleeWeapon(name string) *Equipment {
	return Equip(item.TypeWeaponMelee, name, item.SlotMainHand)
}

// RangedWeapon creates a ranged weapon builder
func RangedWeapon(name string) *Equipment {
	return Equip(item.TypeWeaponRanged, name, item.SlotTwoHand)
}

// MagicWeapon creates a magic weapon builder
func MagicWeapon(name string) *Equipment {
	return Equip(item.TypeWeaponMagic, name, item.SlotMainHand)
}

// Armor creates an armor builder for specified slot
func Armor(name string, slot item.EquipmentSlot) *Equipment {
	var itemType item.Type
	switch slot {
	case item.SlotHead:
		itemType = item.TypeArmorHead
	case item.SlotChest:
		itemType = item.TypeArmorChest
	case item.SlotLegs:
		itemType = item.TypeArmorLegs
	case item.SlotFeet:
		itemType = item.TypeArmorFeet
	case item.SlotHands:
		itemType = item.TypeArmorHands
	default:
		itemType = item.TypeArmorChest
	}
	return Equip(itemType, name, slot)
}

// Accessory creates an accessory builder
func Accessory(name string, slot item.EquipmentSlot) *Equipment {
	var itemType item.Type
	switch slot {
	case item.SlotRing1, item.SlotRing2:
		itemType = item.TypeAccessoryRing
	case item.SlotAmulet:
		itemType = item.TypeAccessoryAmulet
	case item.SlotBelt:
		itemType = item.TypeAccessoryBelt
	default:
		itemType = item.TypeAccessoryRing
	}
	return Equip(itemType, name, slot)
}

// Chainable methods from Item
func (b *Equipment) ID(id string) *Equipment {
	b.Item.ID(id)
	return b
}

func (b *Equipment) Name(name string) *Equipment {
	b.Item.Name(name)
	return b
}

func (b *Equipment) Description(desc string) *Equipment {
	b.Item.Description(desc)
	return b
}

func (b *Equipment) Type(t item.Type) *Equipment {
	b.Item.Type(t)
	return b
}

func (b *Equipment) Rarity(r item.Rarity) *Equipment {
	b.Item.Rarity(r)
	return b
}

func (b *Equipment) Quality(q float64) *Equipment {
	b.Item.Quality(q)
	return b
}

func (b *Equipment) Level(lvl int) *Equipment {
	b.Item.Level(lvl)
	return b
}

func (b *Equipment) Value(v int64) *Equipment {
	b.Item.Value(v)
	return b
}

func (b *Equipment) Weight(w float64) *Equipment {
	b.Item.Weight(w)
	return b
}

func (b *Equipment) Icon(icon string) *Equipment {
	b.Item.Icon(icon)
	return b
}

func (b *Equipment) Tags(tags ...string) *Equipment {
	b.Item.Tags(tags...)
	return b
}

// Equipment-specific methods
func (b *Equipment) Slot(slot item.EquipmentSlot) *Equipment {
	b.slot = slot
	return b
}

func (b *Equipment) Durability(max float64) *Equipment {
	b.maxDurability = max
	return b
}

func (b *Equipment) Sockets(count int, socketType ...item.SocketType) *Equipment {
	b.socketCount = count
	if len(socketType) > 0 {
		b.socketTypes = make([]item.SocketType, count)
		for i := 0; i < count; i++ {
			if i < len(socketType) {
				b.socketTypes[i] = socketType[i]
			} else {
				b.socketTypes[i] = item.SocketTypeUniversal
			}
		}
	}
	return b
}

func (b *Equipment) Require(level int, attrs map[attribute.Type]float64) *Equipment {
	b.requirements = item.NewSimpleRequirements(level, attrs)
	return b
}

func (b *Equipment) RequireLevel(level int) *Equipment {
	b.requirements = item.NewSimpleRequirements(level, nil)
	return b
}

func (b *Equipment) Attribute(mod attribute.Modifier) *Equipment {
	b.attributes = append(b.attributes, mod)
	return b
}

func (b *Equipment) Attributes(mods ...attribute.Modifier) *Equipment {
	b.attributes = append(b.attributes, mods...)
	return b
}

func (b *Equipment) Affix(a affix.Affix) *Equipment {
	b.affixes = append(b.affixes, a)
	return b
}

func (b *Equipment) Affixes(affixes ...affix.Affix) *Equipment {
	b.affixes = append(b.affixes, affixes...)
	return b
}

func (b *Equipment) OnEquip(fn func(ctx context.Context, entity entity.Entity) error) *Equipment {
	b.onEquip = fn
	return b
}

func (b *Equipment) OnUnequip(fn func(ctx context.Context, entity entity.Entity) error) *Equipment {
	b.onUnequip = fn
	return b
}

func (b *Equipment) Build() *item.BaseEquipment {
	cfg := item.EquipmentConfig{
		BaseItemConfig: b.Item.Config(),
		Slot:           b.slot,
		MaxDurability:  b.maxDurability,
		SocketCount:    b.socketCount,
		SocketTypes:    b.socketTypes,
		Requirements:   b.requirements,
	}

	eq := item.NewEquipmentWithConfig(cfg)

	// Add attributes
	for _, mod := range b.attributes {
		eq.AddAttribute(mod)
	}

	// Add affixes
	for _, a := range b.affixes {
		_ = eq.Affixes().Add(a)
	}

	// Set callbacks
	if b.onEquip != nil {
		eq.SetOnEquip(b.onEquip)
	}
	if b.onUnequip != nil {
		eq.SetOnUnequip(b.onUnequip)
	}

	return eq
}
