package equipment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidmovas/Depthborn/internal/item"
)

func createTestEquipment(id, name string, itemType item.Type, slot item.EquipmentSlot) item.Equipment {
	return item.NewEquipmentWithConfig(item.EquipmentConfig{
		BaseItemConfig: item.BaseItemConfig{
			ID:       id,
			Name:     name,
			ItemType: itemType,
			Weight:   5.0,
		},
		Slot: slot,
	})
}

func TestNewManager(t *testing.T) {
	mgr := NewManager()
	assert.NotNil(t, mgr)
	assert.Nil(t, mgr.Owner())
}

func TestNewManagerWithOwner(t *testing.T) {
	mgr := NewManagerWithOwner(nil)
	assert.NotNil(t, mgr)
}

func TestManagerEquipToSlot(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	prev, err := mgr.EquipToSlot(ctx, SlotMainHand, sword)

	require.NoError(t, err)
	assert.Nil(t, prev)
	assert.False(t, mgr.IsEmpty(SlotMainHand))

	equipped := mgr.Get(SlotMainHand)
	assert.Equal(t, "sword-1", equipped.ID())
}

func TestManagerEquipToSlotReplacesExisting(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword1 := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	sword2 := createTestEquipment("sword-2", "Steel Sword", item.TypeWeaponMelee, item.SlotMainHand)

	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword1)
	prev, err := mgr.EquipToSlot(ctx, SlotMainHand, sword2)

	require.NoError(t, err)
	assert.NotNil(t, prev)
	assert.Equal(t, "sword-1", prev.ID())
	assert.Equal(t, "sword-2", mgr.Get(SlotMainHand).ID())
}

func TestManagerEquipNil(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	_, err := mgr.Equip(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil item")

	_, err = mgr.EquipToSlot(ctx, SlotMainHand, nil)
	assert.Error(t, err)
}

func TestManagerEquipIncompatibleSlot(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	helmet := createTestEquipment("helmet-1", "Iron Helmet", item.TypeArmorHead, item.SlotHead)
	_, err := mgr.EquipToSlot(ctx, SlotMainHand, helmet)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be equipped to slot")
}

func TestManagerEquipAutoSlot(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	_, err := mgr.Equip(ctx, sword)

	require.NoError(t, err)
	assert.False(t, mgr.IsEmpty(SlotMainHand))
}

func TestManagerEquipRings(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	ring1 := createTestEquipment("ring-1", "Gold Ring", item.TypeAccessoryRing, item.SlotRing1)
	ring2 := createTestEquipment("ring-2", "Silver Ring", item.TypeAccessoryRing, item.SlotRing1)

	// First ring goes to Ring1
	_, err := mgr.Equip(ctx, ring1)
	require.NoError(t, err)
	assert.False(t, mgr.IsEmpty(SlotRing1))

	// Second ring goes to Ring2
	_, err = mgr.Equip(ctx, ring2)
	require.NoError(t, err)
	assert.False(t, mgr.IsEmpty(SlotRing2))
}

func TestManagerUnequip(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword)

	unequipped, err := mgr.Unequip(ctx, SlotMainHand)
	require.NoError(t, err)
	assert.NotNil(t, unequipped)
	assert.Equal(t, "sword-1", unequipped.ID())
	assert.True(t, mgr.IsEmpty(SlotMainHand))
}

func TestManagerUnequipEmpty(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	unequipped, err := mgr.Unequip(ctx, SlotMainHand)
	require.NoError(t, err)
	assert.Nil(t, unequipped)
}

func TestManagerUnequipAll(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	helmet := createTestEquipment("helmet-1", "Iron Helmet", item.TypeArmorHead, item.SlotHead)

	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword)
	_, _ = mgr.EquipToSlot(ctx, SlotHead, helmet)

	unequipped, err := mgr.UnequipAll(ctx)
	require.NoError(t, err)
	assert.Len(t, unequipped, 2)
	assert.True(t, mgr.IsEmpty(SlotMainHand))
	assert.True(t, mgr.IsEmpty(SlotHead))
}

func TestManagerGetAll(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	helmet := createTestEquipment("helmet-1", "Iron Helmet", item.TypeArmorHead, item.SlotHead)

	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword)
	_, _ = mgr.EquipToSlot(ctx, SlotHead, helmet)

	all := mgr.GetAll()
	assert.Len(t, all, 2)
	assert.NotNil(t, all[SlotMainHand])
	assert.NotNil(t, all[SlotHead])
}

func TestManagerGetFilledSlots(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword)

	slots := mgr.GetFilledSlots()
	assert.Len(t, slots, 1)
	assert.Contains(t, slots, SlotMainHand)
}

func TestManagerCanEquip(t *testing.T) {
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	assert.True(t, mgr.CanEquip(sword))

	assert.False(t, mgr.CanEquip(nil))
}

func TestManagerCanEquipToSlot(t *testing.T) {
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	assert.True(t, mgr.CanEquipToSlot(SlotMainHand, sword))
	assert.False(t, mgr.CanEquipToSlot(SlotHead, sword))
	assert.False(t, mgr.CanEquipToSlot(SlotMainHand, nil))
}

func TestManagerTotalWeight(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := item.NewEquipmentWithConfig(item.EquipmentConfig{
		BaseItemConfig: item.BaseItemConfig{
			ID:       "sword-1",
			Name:     "Iron Sword",
			ItemType: item.TypeWeaponMelee,
			Weight:   10.0,
		},
		Slot: item.SlotMainHand,
	})

	helmet := item.NewEquipmentWithConfig(item.EquipmentConfig{
		BaseItemConfig: item.BaseItemConfig{
			ID:       "helmet-1",
			Name:     "Iron Helmet",
			ItemType: item.TypeArmorHead,
			Weight:   5.0,
		},
		Slot: item.SlotHead,
	})

	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword)
	_, _ = mgr.EquipToSlot(ctx, SlotHead, helmet)

	assert.Equal(t, 15.0, mgr.TotalWeight())
}

func TestManagerGetAllModifiers(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword)

	// Equipment may return nil or empty slice for modifiers
	mods := mgr.GetAllModifiers()
	// Just verify it doesn't panic
	_ = mods
}

func TestManagerOwner(t *testing.T) {
	mgr := NewManager()

	assert.Nil(t, mgr.Owner())

	mgr.SetOwner(nil)
	assert.Nil(t, mgr.Owner())
}

func TestManagerCallbacks(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	var equipEvents []string
	var unequipEvents []string

	mgr.OnEquip(func(ctx context.Context, slot Slot, equip item.Equipment) {
		equipEvents = append(equipEvents, equip.ID())
	})

	mgr.OnUnequip(func(ctx context.Context, slot Slot, equip item.Equipment) {
		unequipEvents = append(unequipEvents, equip.ID())
	})

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword)

	assert.Equal(t, []string{"sword-1"}, equipEvents)

	_, _ = mgr.Unequip(ctx, SlotMainHand)
	assert.Equal(t, []string{"sword-1"}, unequipEvents)
}

func TestManagerSerializeState(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword)

	state, err := mgr.SerializeState()
	require.NoError(t, err)
	assert.NotNil(t, state)
}

func TestManagerGetSlotItemIDs(t *testing.T) {
	ctx := context.Background()
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	helmet := createTestEquipment("helmet-1", "Iron Helmet", item.TypeArmorHead, item.SlotHead)

	_, _ = mgr.EquipToSlot(ctx, SlotMainHand, sword)
	_, _ = mgr.EquipToSlot(ctx, SlotHead, helmet)

	ids := mgr.GetSlotItemIDs()
	assert.Len(t, ids, 2)
	assert.Equal(t, "sword-1", ids[SlotMainHand])
	assert.Equal(t, "helmet-1", ids[SlotHead])
}

func TestManagerSetItemDirect(t *testing.T) {
	mgr := NewManager()

	sword := createTestEquipment("sword-1", "Iron Sword", item.TypeWeaponMelee, item.SlotMainHand)
	mgr.SetItemDirect(SlotMainHand, sword)

	assert.False(t, mgr.IsEmpty(SlotMainHand))
	assert.Equal(t, "sword-1", mgr.Get(SlotMainHand).ID())
}

func TestSlotHelpers(t *testing.T) {
	assert.Len(t, AllSlots(), 15)
	assert.Len(t, WeaponSlots(), 2)
	assert.Len(t, ArmorSlots(), 8)
	assert.Len(t, AccessorySlots(), 5)
}

func TestSlotCategory(t *testing.T) {
	assert.Equal(t, CategoryWeapon, SlotMainHand.Category())
	assert.Equal(t, CategoryWeapon, SlotOffHand.Category())
	assert.Equal(t, CategoryArmor, SlotHead.Category())
	assert.Equal(t, CategoryArmor, SlotChest.Category())
	assert.Equal(t, CategoryAccessory, SlotNeck.Category())
	assert.Equal(t, CategoryAccessory, SlotRing1.Category())
}

func TestSlotDisplayName(t *testing.T) {
	assert.Equal(t, "Main Hand", SlotMainHand.DisplayName())
	assert.Equal(t, "Head", SlotHead.DisplayName())
	assert.Equal(t, "Ring (Left)", SlotRing1.DisplayName())
}

func TestSlotCompatibility(t *testing.T) {
	// Weapons can go in main hand
	assert.True(t, IsCompatible(SlotMainHand, item.TypeWeaponMelee))
	assert.True(t, IsCompatible(SlotOffHand, item.TypeWeaponMelee))

	// Helmet can only go in head slot
	assert.True(t, IsCompatible(SlotHead, item.TypeArmorHead))
	assert.False(t, IsCompatible(SlotChest, item.TypeArmorHead))

	// Rings can go in ring slots
	assert.True(t, IsCompatible(SlotRing1, item.TypeAccessoryRing))
	assert.True(t, IsCompatible(SlotRing2, item.TypeAccessoryRing))
}

func TestCompatibleSlots(t *testing.T) {
	slots := CompatibleSlots(item.TypeWeaponMelee)
	assert.Len(t, slots, 2)
	assert.Contains(t, slots, SlotMainHand)
	assert.Contains(t, slots, SlotOffHand)

	// Unknown type returns nil
	slots = CompatibleSlots(item.TypeMaterial)
	assert.Nil(t, slots)
}

func TestDefaultSlot(t *testing.T) {
	assert.Equal(t, SlotMainHand, DefaultSlot(item.TypeWeaponMelee))
	assert.Equal(t, SlotHead, DefaultSlot(item.TypeArmorHead))
	assert.Equal(t, SlotRing1, DefaultSlot(item.TypeAccessoryRing))
	assert.Equal(t, Slot(""), DefaultSlot(item.TypeMaterial))
}

func TestSlotConversion(t *testing.T) {
	// equipment.Slot -> item.EquipmentSlot
	assert.Equal(t, item.SlotMainHand, ToItemSlot(SlotMainHand))
	assert.Equal(t, item.SlotHead, ToItemSlot(SlotHead))
	assert.Equal(t, item.SlotAmulet, ToItemSlot(SlotNeck))

	// item.EquipmentSlot -> equipment.Slot
	assert.Equal(t, SlotMainHand, FromItemSlot(item.SlotMainHand))
	assert.Equal(t, SlotHead, FromItemSlot(item.SlotHead))
	assert.Equal(t, SlotNeck, FromItemSlot(item.SlotAmulet))
}
