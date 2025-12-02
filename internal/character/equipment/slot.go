package equipment

import "github.com/davidmovas/Depthborn/internal/item"

// Slot defines all available equipment slots
type Slot string

const (
	SlotMainHand  Slot = "main_hand"
	SlotOffHand   Slot = "off_hand"
	SlotHead      Slot = "head"
	SlotShoulders Slot = "shoulders"
	SlotChest     Slot = "chest"
	SlotHands     Slot = "hands"
	SlotWaist     Slot = "waist"
	SlotLegs      Slot = "legs"
	SlotFeet      Slot = "feet"
	SlotNeck      Slot = "neck"      // Amulet
	SlotRing1     Slot = "ring_1"    // Left ring
	SlotRing2     Slot = "ring_2"    // Right ring
	SlotBack      Slot = "back"      // Cape/Cloak
	SlotTrinket1  Slot = "trinket_1" // Special items
	SlotTrinket2  Slot = "trinket_2"
)

// AllSlots returns all available equipment slots in order
func AllSlots() []Slot {
	return []Slot{
		SlotMainHand, SlotOffHand,
		SlotHead, SlotShoulders, SlotChest, SlotHands, SlotWaist, SlotLegs, SlotFeet,
		SlotNeck, SlotRing1, SlotRing2,
		SlotBack, SlotTrinket1, SlotTrinket2,
	}
}

// WeaponSlots returns weapon-related slots
func WeaponSlots() []Slot {
	return []Slot{SlotMainHand, SlotOffHand}
}

// ArmorSlots returns armor slots
func ArmorSlots() []Slot {
	return []Slot{SlotHead, SlotShoulders, SlotChest, SlotHands, SlotWaist, SlotLegs, SlotFeet, SlotBack}
}

// AccessorySlots returns accessory slots
func AccessorySlots() []Slot {
	return []Slot{SlotNeck, SlotRing1, SlotRing2, SlotTrinket1, SlotTrinket2}
}

// SlotCategory categorizes equipment slots
type SlotCategory string

const (
	CategoryWeapon    SlotCategory = "weapon"
	CategoryArmor     SlotCategory = "armor"
	CategoryAccessory SlotCategory = "accessory"
)

// Category returns the category for a slot
func (s Slot) Category() SlotCategory {
	switch s {
	case SlotMainHand, SlotOffHand:
		return CategoryWeapon
	case SlotNeck, SlotRing1, SlotRing2, SlotTrinket1, SlotTrinket2:
		return CategoryAccessory
	default:
		return CategoryArmor
	}
}

// String returns the slot name
func (s Slot) String() string {
	return string(s)
}

// DisplayName returns human-readable slot name
func (s Slot) DisplayName() string {
	names := map[Slot]string{
		SlotMainHand:  "Main Hand",
		SlotOffHand:   "Off Hand",
		SlotHead:      "Head",
		SlotShoulders: "Shoulders",
		SlotChest:     "Chest",
		SlotHands:     "Hands",
		SlotWaist:     "Waist",
		SlotLegs:      "Legs",
		SlotFeet:      "Feet",
		SlotNeck:      "Neck",
		SlotRing1:     "Ring (Left)",
		SlotRing2:     "Ring (Right)",
		SlotBack:      "Back",
		SlotTrinket1:  "Trinket 1",
		SlotTrinket2:  "Trinket 2",
	}
	if name, ok := names[s]; ok {
		return name
	}
	return string(s)
}

// slotCompatibility maps item types to compatible slots
var slotCompatibility = map[item.Type][]Slot{
	item.TypeWeaponMelee:     {SlotMainHand, SlotOffHand},
	item.TypeWeaponRanged:    {SlotMainHand},
	item.TypeWeaponMagic:     {SlotMainHand, SlotOffHand},
	item.TypeArmorHead:       {SlotHead},
	item.TypeArmorChest:      {SlotChest},
	item.TypeArmorLegs:       {SlotLegs},
	item.TypeArmorFeet:       {SlotFeet},
	item.TypeArmorHands:      {SlotHands},
	item.TypeAccessoryRing:   {SlotRing1, SlotRing2},
	item.TypeAccessoryAmulet: {SlotNeck},
	item.TypeAccessoryBelt:   {SlotWaist},
}

// CompatibleSlots returns slots where an item type can be equipped
func CompatibleSlots(itemType item.Type) []Slot {
	if slots, ok := slotCompatibility[itemType]; ok {
		return slots
	}
	return nil
}

// IsCompatible checks if an item type can be equipped in a slot
func IsCompatible(slot Slot, itemType item.Type) bool {
	compatibleSlots := CompatibleSlots(itemType)
	for _, s := range compatibleSlots {
		if s == slot {
			return true
		}
	}
	return false
}

// DefaultSlot returns the primary/default slot for an item type
func DefaultSlot(itemType item.Type) Slot {
	slots := CompatibleSlots(itemType)
	if len(slots) > 0 {
		return slots[0]
	}
	return ""
}

// ToItemSlot converts equipment.Slot to item.EquipmentSlot
func ToItemSlot(s Slot) item.EquipmentSlot {
	// Map our extended slots to item package slots
	mapping := map[Slot]item.EquipmentSlot{
		SlotMainHand: item.SlotMainHand,
		SlotOffHand:  item.SlotOffHand,
		SlotHead:     item.SlotHead,
		SlotChest:    item.SlotChest,
		SlotLegs:     item.SlotLegs,
		SlotFeet:     item.SlotFeet,
		SlotHands:    item.SlotHands,
		SlotNeck:     item.SlotAmulet,
		SlotRing1:    item.SlotRing1,
		SlotRing2:    item.SlotRing2,
		SlotWaist:    item.SlotBelt,
	}
	if mapped, ok := mapping[s]; ok {
		return mapped
	}
	return item.EquipmentSlot(s)
}

// FromItemSlot converts item.EquipmentSlot to equipment.Slot
func FromItemSlot(s item.EquipmentSlot) Slot {
	mapping := map[item.EquipmentSlot]Slot{
		item.SlotMainHand: SlotMainHand,
		item.SlotOffHand:  SlotOffHand,
		item.SlotTwoHand:  SlotMainHand, // Two-hand goes to main hand
		item.SlotHead:     SlotHead,
		item.SlotChest:    SlotChest,
		item.SlotLegs:     SlotLegs,
		item.SlotFeet:     SlotFeet,
		item.SlotHands:    SlotHands,
		item.SlotAmulet:   SlotNeck,
		item.SlotRing1:    SlotRing1,
		item.SlotRing2:    SlotRing2,
		item.SlotBelt:     SlotWaist,
	}
	if mapped, ok := mapping[s]; ok {
		return mapped
	}
	return Slot(s)
}
