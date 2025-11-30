package item

import (
	"context"
	"testing"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/stretchr/testify/require"
)

func TestBaseEquipment(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("NewBaseEquipment creates with defaults", func(t *testing.T) {
			equip := NewBaseEquipment("test-equip", TypeWeaponMelee, "Iron Sword", SlotMainHand)

			require.Equal(t, "Iron Sword", equip.Name())
			require.Equal(t, SlotMainHand, equip.Slot())
			require.Equal(t, 100.0, equip.Durability())
			require.Equal(t, 100.0, equip.MaxDurability())
			require.True(t, equip.IsEquippable())
		})

		t.Run("NewEquipmentWithConfig respects all fields", func(t *testing.T) {
			cfg := EquipmentConfig{
				BaseItemConfig: BaseItemConfig{
					Name:     "Epic Sword",
					ItemType: TypeWeaponMelee,
					Rarity:   RarityEpic,
					Level:    20,
				},
				Slot:          SlotMainHand,
				MaxDurability: 200,
				SocketCount:   2,
				SocketTypes:   []SocketType{SocketTypeGem, SocketTypeRune},
			}
			equip := NewEquipmentWithConfig(cfg)

			require.Equal(t, RarityEpic, equip.Rarity())
			require.Equal(t, 20, equip.Level())
			require.Equal(t, 200.0, equip.MaxDurability())
			require.Equal(t, 2, equip.SocketCount())
		})
	})

	t.Run("Durability", func(t *testing.T) {
		t.Run("DamageItem reduces durability", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeArmorChest, "Armor", SlotChest)

			equip.DamageItem(25)

			require.Equal(t, 75.0, equip.Durability())
		})

		t.Run("DamageItem cannot go below zero", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeArmorChest, "Armor", SlotChest)

			equip.DamageItem(150)

			require.Equal(t, 0.0, equip.Durability())
		})

		t.Run("Repair restores durability", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeArmorChest, "Armor", SlotChest)
			equip.DamageItem(50)

			equip.Repair(30)

			require.Equal(t, 80.0, equip.Durability())
		})

		t.Run("Repair cannot exceed max durability", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeArmorChest, "Armor", SlotChest)
			equip.DamageItem(50)

			equip.Repair(100)

			require.Equal(t, 100.0, equip.Durability())
		})

		t.Run("IsBroken returns true when durability is zero", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeArmorChest, "Armor", SlotChest)

			require.False(t, equip.IsBroken())

			equip.DamageItem(100)

			require.True(t, equip.IsBroken())
		})

		t.Run("DurabilityPercent returns correct percentage", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeArmorChest, "Armor", SlotChest)
			equip.DamageItem(25)

			require.Equal(t, 0.75, equip.DurabilityPercent())
		})

		t.Run("SetMaxDurability updates max", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeArmorChest, "Armor", SlotChest)

			equip.SetMaxDurability(200)
			require.Equal(t, 200.0, equip.MaxDurability())

			// Setting max below current should adjust current
			equip.SetMaxDurability(50)
			require.Equal(t, 50.0, equip.Durability())
		})
	})

	t.Run("Sockets", func(t *testing.T) {
		t.Run("AddSocket adds a socket", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeWeaponMelee, "Sword", SlotMainHand)

			equip.AddSocket(SocketTypeGem)

			require.Equal(t, 1, equip.SocketCount())
			socketType, ok := equip.GetSocketType(0)
			require.True(t, ok)
			require.Equal(t, SocketTypeGem, socketType)
		})

		t.Run("SetSocket places item in socket", func(t *testing.T) {
			equip := NewEquipmentWithConfig(EquipmentConfig{
				BaseItemConfig: BaseItemConfig{Name: "Sword", ItemType: TypeWeaponMelee},
				Slot:           SlotMainHand,
				SocketCount:    1,
				SocketTypes:    []SocketType{SocketTypeGem},
			})
			gem := NewBaseSocketable("gem-1", TypeGem, "Ruby", SocketTypeGem)

			err := equip.SetSocket(0, gem)

			require.NoError(t, err)
			socket, ok := equip.GetSocket(0)
			require.True(t, ok)
			require.NotNil(t, socket)
		})

		t.Run("SetSocket fails for invalid index", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeWeaponMelee, "Sword", SlotMainHand)
			gem := NewBaseSocketable("gem-1", TypeGem, "Ruby", SocketTypeGem)

			err := equip.SetSocket(0, gem)

			require.Error(t, err)
		})

		t.Run("SetSocket fails for incompatible type", func(t *testing.T) {
			equip := NewEquipmentWithConfig(EquipmentConfig{
				BaseItemConfig: BaseItemConfig{Name: "Sword", ItemType: TypeWeaponMelee},
				Slot:           SlotMainHand,
				SocketCount:    1,
				SocketTypes:    []SocketType{SocketTypeGem},
			})
			rune := NewBaseSocketable("rune-1", TypeRune, "Rune of Power", SocketTypeRune)

			err := equip.SetSocket(0, rune)

			require.Error(t, err)
		})

		t.Run("RemoveSocket removes item from socket", func(t *testing.T) {
			equip := NewEquipmentWithConfig(EquipmentConfig{
				BaseItemConfig: BaseItemConfig{Name: "Sword", ItemType: TypeWeaponMelee},
				Slot:           SlotMainHand,
				SocketCount:    1,
				SocketTypes:    []SocketType{SocketTypeGem},
			})
			gem := NewBaseSocketable("gem-1", TypeGem, "Ruby", SocketTypeGem)
			equip.SetSocket(0, gem)

			removed, err := equip.RemoveSocket(0)

			require.NoError(t, err)
			require.NotNil(t, removed)
			socket, ok := equip.GetSocket(0)
			require.True(t, !ok || socket == nil)
		})

		t.Run("EmptySocketCount", func(t *testing.T) {
			equip := NewEquipmentWithConfig(EquipmentConfig{
				BaseItemConfig: BaseItemConfig{Name: "Sword", ItemType: TypeWeaponMelee},
				Slot:           SlotMainHand,
				SocketCount:    2,
				SocketTypes:    []SocketType{SocketTypeGem, SocketTypeGem},
			})

			require.Equal(t, 2, equip.EmptySocketCount())

			gem1 := NewBaseSocketable("gem-1", TypeGem, "Ruby", SocketTypeGem)
			gem2 := NewBaseSocketable("gem-2", TypeGem, "Sapphire", SocketTypeGem)
			equip.SetSocket(0, gem1)
			equip.SetSocket(1, gem2)

			require.Equal(t, 0, equip.EmptySocketCount())
		})
	})

	t.Run("Requirements", func(t *testing.T) {
		t.Run("SetRequirements and GetRequirements", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeWeaponMelee, "Heavy Sword", SlotMainHand)

			attrs := map[attribute.Type]float64{
				attribute.AttrStrength: 20,
			}
			req := NewSimpleRequirements(10, attrs)
			equip.SetRequirements(req)

			got := equip.Requirements()
			require.NotNil(t, got)
			require.Equal(t, 10, got.Level())
		})

		t.Run("CanEquip checks requirements", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeWeaponMelee, "Heavy Sword", SlotMainHand)

			attrs := map[attribute.Type]float64{
				attribute.AttrStrength: 20,
			}
			req := NewSimpleRequirements(10, attrs)
			equip.SetRequirements(req)

			// nil entity should fail
			require.False(t, equip.CanEquip(nil))
		})
	})

	t.Run("EquipCallbacks", func(t *testing.T) {
		t.Run("SetOnEquip and SetOnUnequip set callbacks", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeArmorChest, "Armor", SlotChest)

			equip.SetOnEquip(func(_ context.Context, _ entity.Entity) error {
				return nil
			})
			equip.SetOnUnequip(func(_ context.Context, _ entity.Entity) error {
				return nil
			})

			// Verify callbacks are set (no panic means success)
		})
	})

	t.Run("Attributes", func(t *testing.T) {
		t.Run("AddAttribute adds modifier", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeWeaponMelee, "Sword", SlotMainHand)

			mod := attribute.NewModifier("test-mod", attribute.ModFlat, 10, string(attribute.AttrStrength))
			equip.AddAttribute(mod)

			attrs := equip.Attributes()
			require.Len(t, attrs, 1)
		})

		t.Run("RemoveAttribute removes modifier", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeWeaponMelee, "Sword", SlotMainHand)

			mod := attribute.NewModifier("test-mod", attribute.ModFlat, 10, string(attribute.AttrStrength))
			equip.AddAttribute(mod)
			equip.RemoveAttribute("test-mod")

			attrs := equip.Attributes()
			require.Empty(t, attrs)
		})
	})

	t.Run("Clone", func(t *testing.T) {
		t.Run("creates independent copy", func(t *testing.T) {
			original := NewEquipmentWithConfig(EquipmentConfig{
				BaseItemConfig: BaseItemConfig{
					Name:     "Magic Sword",
					ItemType: TypeWeaponMelee,
					Rarity:   RarityRare,
				},
				Slot:          SlotMainHand,
				MaxDurability: 150,
				SocketTypes:   []SocketType{SocketTypeGem},
			})
			original.DamageItem(25)

			cloned := original.Clone().(*BaseEquipment)

			require.NotEqual(t, original.ID(), cloned.ID())
			require.Equal(t, original.Durability(), cloned.Durability())
			require.Equal(t, original.Name(), cloned.Name())
			require.Equal(t, original.Slot(), cloned.Slot())
			require.Equal(t, original.SocketCount(), cloned.SocketCount())
		})
	})

	t.Run("Serialization", func(t *testing.T) {
		t.Run("Marshal and Unmarshal roundtrip", func(t *testing.T) {
			original := NewEquipmentWithConfig(EquipmentConfig{
				BaseItemConfig: BaseItemConfig{
					ID:       "equip-serial",
					Name:     "Test Sword",
					ItemType: TypeWeaponMelee,
					Rarity:   RarityEpic,
					Level:    25,
				},
				Slot:          SlotMainHand,
				MaxDurability: 200,
				SocketCount:   2,
				SocketTypes:   []SocketType{SocketTypeGem, SocketTypeRune},
			})
			original.DamageItem(50)

			data, err := original.Marshal()
			require.NoError(t, err)

			restored := &BaseEquipment{}
			err = restored.Unmarshal(data)
			require.NoError(t, err)

			require.Equal(t, original.ID(), restored.ID())
			require.Equal(t, original.Name(), restored.Name())
			require.Equal(t, original.Slot(), restored.Slot())
			require.Equal(t, original.Durability(), restored.Durability())
			require.Equal(t, original.MaxDurability(), restored.MaxDurability())
			require.Equal(t, original.SocketCount(), restored.SocketCount())
		})
	})

	t.Run("Validation", func(t *testing.T) {
		t.Run("valid equipment passes", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeWeaponMelee, "Valid Sword", SlotMainHand)
			require.NoError(t, equip.Validate())
		})

		t.Run("empty slot fails", func(t *testing.T) {
			equip := NewEquipmentWithConfig(EquipmentConfig{
				BaseItemConfig: BaseItemConfig{
					Name:     "No Slot",
					ItemType: TypeWeaponMelee,
				},
			})
			require.Error(t, equip.Validate())
		})

		t.Run("durability exceeding max fails", func(t *testing.T) {
			equip := NewBaseEquipment("", TypeWeaponMelee, "Test", SlotMainHand)
			equip.durability = 200 // bypass setter
			require.Error(t, equip.Validate())
		})
	})
}
