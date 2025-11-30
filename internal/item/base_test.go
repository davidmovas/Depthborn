package item

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseItem(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("NewBaseItem creates item with defaults", func(t *testing.T) {
			item := NewBaseItem("test-1", TypeMaterial, "Test Item")

			require.NotEmpty(t, item.ID())
			require.Equal(t, "Test Item", item.Name())
			require.Equal(t, TypeMaterial, item.ItemType())
			require.Equal(t, 1, item.StackSize())
			require.Equal(t, 1.0, item.Quality())
			require.Equal(t, 1, item.Level())
		})

		t.Run("NewBaseItemWithConfig respects all fields", func(t *testing.T) {
			cfg := BaseItemConfig{
				ID:           "custom-id",
				Name:         "Custom Item",
				Description:  "A custom description",
				ItemType:     TypeConsumable,
				Rarity:       RarityRare,
				Quality:      0.75,
				Level:        10,
				MaxStackSize: 99,
				Value:        500,
				Weight:       2.5,
				Icon:         "custom_icon",
				Tags:         []string{"tag1", "tag2"},
			}
			item := NewBaseItemWithConfig(cfg)

			require.Equal(t, "custom-id", item.ID())
			require.Equal(t, "A custom description", item.Description())
			require.Equal(t, RarityRare, item.Rarity())
			require.Equal(t, 0.75, item.Quality())
			require.Equal(t, 10, item.Level())
			require.Equal(t, 99, item.MaxStackSize())
			require.Equal(t, int64(500), item.Value())
			require.Equal(t, 2.5, item.Weight())
			require.Equal(t, "custom_icon", item.Icon())
			require.True(t, item.Tags().Has("tag1"))
			require.True(t, item.Tags().Has("tag2"))
		})

		t.Run("applies defaults for invalid values", func(t *testing.T) {
			cfg := BaseItemConfig{
				Name:         "Test",
				ItemType:     TypeMaterial,
				Quality:      -5.0, // invalid
				Level:        0,    // invalid
				MaxStackSize: 0,    // invalid
				Weight:       0,    // invalid
			}
			item := NewBaseItemWithConfig(cfg)

			require.Equal(t, 1.0, item.Quality())
			require.Equal(t, 1, item.Level())
			require.Equal(t, 1, item.MaxStackSize())
			require.Equal(t, 0.1, item.Weight())
		})
	})

	t.Run("Stacking", func(t *testing.T) {
		t.Run("AddStack increases stack size", func(t *testing.T) {
			item := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Stackable",
				ItemType:     TypeMaterial,
				MaxStackSize: 10,
			})

			require.True(t, item.AddStack(5))
			require.Equal(t, 6, item.StackSize()) // 1 initial + 5 added
		})

		t.Run("AddStack fails when exceeds max", func(t *testing.T) {
			item := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Stackable",
				ItemType:     TypeMaterial,
				MaxStackSize: 5,
			})

			require.False(t, item.AddStack(10))
			require.Equal(t, 1, item.StackSize())
		})

		t.Run("RemoveStack decreases stack size", func(t *testing.T) {
			item := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Stackable",
				ItemType:     TypeMaterial,
				MaxStackSize: 10,
			})
			item.AddStack(9) // total 10

			removed := item.RemoveStack(3)

			require.Equal(t, 3, removed)
			require.Equal(t, 7, item.StackSize())
		})

		t.Run("RemoveStack returns actual removed when exceeds stack", func(t *testing.T) {
			item := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Stackable",
				ItemType:     TypeMaterial,
				MaxStackSize: 10,
			})
			item.AddStack(4) // total 5

			removed := item.RemoveStack(10)

			require.Equal(t, 5, removed)
			require.Equal(t, 0, item.StackSize())
		})

		t.Run("CanStackWith checks compatibility", func(t *testing.T) {
			item1 := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Material A",
				ItemType:     TypeMaterial,
				Rarity:       RarityCommon,
				MaxStackSize: 10,
			})
			item2 := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Material B",
				ItemType:     TypeMaterial,
				Rarity:       RarityCommon,
				MaxStackSize: 10,
			})
			item3 := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Different Type",
				ItemType:     TypeConsumable,
				Rarity:       RarityCommon,
				MaxStackSize: 10,
			})

			require.True(t, item1.CanStackWith(item2))
			require.False(t, item1.CanStackWith(item3))
			require.False(t, item1.CanStackWith(item1), "item should not stack with itself")
		})
	})

	t.Run("Properties", func(t *testing.T) {
		t.Run("SetName updates name", func(t *testing.T) {
			item := NewBaseItem("", TypeMaterial, "Original")

			item.SetName("Updated")

			require.Equal(t, "Updated", item.Name())
		})

		t.Run("SetQuality clamps to valid range", func(t *testing.T) {
			item := NewBaseItem("", TypeMaterial, "Test")

			item.SetQuality(-0.5)
			require.Equal(t, 0.0, item.Quality())

			item.SetQuality(1.5)
			require.Equal(t, 1.0, item.Quality())

			item.SetQuality(0.5)
			require.Equal(t, 0.5, item.Quality())
		})

		t.Run("SetLevel enforces minimum", func(t *testing.T) {
			item := NewBaseItem("", TypeMaterial, "Test")

			item.SetLevel(0)
			require.Equal(t, 1, item.Level())

			item.SetLevel(50)
			require.Equal(t, 50, item.Level())
		})

		t.Run("SetValue enforces non-negative", func(t *testing.T) {
			item := NewBaseItem("", TypeMaterial, "Test")

			item.SetValue(-100)
			require.Equal(t, int64(0), item.Value())

			item.SetValue(1000)
			require.Equal(t, int64(1000), item.Value())
		})
	})

	t.Run("ItemTypeChecks", func(t *testing.T) {
		t.Run("IsEquippable returns true for equipment types", func(t *testing.T) {
			equipTypes := []Type{
				TypeWeaponMelee, TypeWeaponRanged, TypeWeaponMagic,
				TypeArmorHead, TypeArmorChest, TypeArmorLegs, TypeArmorFeet, TypeArmorHands,
				TypeAccessoryRing, TypeAccessoryAmulet, TypeAccessoryBelt,
			}

			for _, itemType := range equipTypes {
				item := NewBaseItem("", itemType, "Test")
				require.True(t, item.IsEquippable(), "expected %s to be equippable", itemType)
			}
		})

		t.Run("IsEquippable returns false for non-equipment types", func(t *testing.T) {
			nonEquipTypes := []Type{TypeMaterial, TypeConsumable, TypeQuest, TypeContainer}

			for _, itemType := range nonEquipTypes {
				item := NewBaseItem("", itemType, "Test")
				require.False(t, item.IsEquippable(), "expected %s to not be equippable", itemType)
			}
		})

		t.Run("IsConsumable", func(t *testing.T) {
			consumable := NewBaseItem("", TypeConsumable, "Potion")
			material := NewBaseItem("", TypeMaterial, "Iron")

			require.True(t, consumable.IsConsumable())
			require.False(t, material.IsConsumable())
		})

		t.Run("IsQuestItem", func(t *testing.T) {
			questItem := NewBaseItem("", TypeQuest, "Ancient Key")
			material := NewBaseItem("", TypeMaterial, "Iron")

			require.True(t, questItem.IsQuestItem())
			require.False(t, material.IsQuestItem())
		})

		t.Run("IsTradeable", func(t *testing.T) {
			questItem := NewBaseItem("", TypeQuest, "Ancient Key")
			material := NewBaseItem("", TypeMaterial, "Iron")

			require.False(t, questItem.IsTradeable())
			require.True(t, material.IsTradeable())
		})
	})

	t.Run("Computed", func(t *testing.T) {
		t.Run("TotalWeight returns weight times stack", func(t *testing.T) {
			item := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Heavy",
				ItemType:     TypeMaterial,
				Weight:       5.0,
				MaxStackSize: 10,
			})
			item.AddStack(4) // total 5

			require.Equal(t, 25.0, item.TotalWeight())
		})

		t.Run("TotalValue returns value times stack", func(t *testing.T) {
			item := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Valuable",
				ItemType:     TypeMaterial,
				Value:        100,
				MaxStackSize: 10,
			})
			item.AddStack(2) // total 3

			require.Equal(t, int64(300), item.TotalValue())
		})

		t.Run("DisplayName includes rarity for non-common items", func(t *testing.T) {
			common := NewBaseItemWithConfig(BaseItemConfig{
				Name:     "Common Item",
				ItemType: TypeMaterial,
				Rarity:   RarityCommon,
			})
			rare := NewBaseItemWithConfig(BaseItemConfig{
				Name:     "Rare Item",
				ItemType: TypeMaterial,
				Rarity:   RarityRare,
			})

			require.Equal(t, "Common Item", common.DisplayName())
			require.NotEqual(t, "Rare Item", rare.DisplayName(), "rare item should have rarity prefix")
		})
	})

	t.Run("Clone", func(t *testing.T) {
		t.Run("creates independent copy", func(t *testing.T) {
			original := NewBaseItemWithConfig(BaseItemConfig{
				Name:         "Original",
				ItemType:     TypeMaterial,
				Rarity:       RarityRare,
				Value:        100,
				MaxStackSize: 10,
				Tags:         []string{"tag1"},
			})
			original.AddStack(4)

			cloned := original.Clone().(*BaseItem)

			require.NotEqual(t, original.ID(), cloned.ID())
			require.Equal(t, original.Name(), cloned.Name())
			require.Equal(t, original.Rarity(), cloned.Rarity())
			require.Equal(t, original.Value(), cloned.Value())

			// Modifying clone should not affect original
			cloned.SetName("Modified")
			require.NotEqual(t, "Modified", original.Name())
		})

		t.Run("nil item returns nil", func(t *testing.T) {
			var item *BaseItem
			require.Nil(t, item.Clone())
		})
	})

	t.Run("Serialization", func(t *testing.T) {
		t.Run("Marshal and Unmarshal roundtrip", func(t *testing.T) {
			original := NewBaseItemWithConfig(BaseItemConfig{
				ID:           "ser-test",
				Name:         "Serializable",
				Description:  "Test description",
				ItemType:     TypeMaterial,
				Rarity:       RarityEpic,
				Quality:      0.85,
				Level:        15,
				MaxStackSize: 50,
				Value:        1500,
				Weight:       3.5,
				Icon:         "test_icon",
				Tags:         []string{"serialize", "test"},
			})
			original.AddStack(24)

			data, err := original.Marshal()
			require.NoError(t, err)

			restored := &BaseItem{}
			err = restored.Unmarshal(data)
			require.NoError(t, err)

			require.Equal(t, original.ID(), restored.ID())
			require.Equal(t, original.Name(), restored.Name())
			require.Equal(t, original.Description(), restored.Description())
			require.Equal(t, original.ItemType(), restored.ItemType())
			require.Equal(t, original.Rarity(), restored.Rarity())
			require.Equal(t, original.Quality(), restored.Quality())
			require.Equal(t, original.Level(), restored.Level())
			require.Equal(t, original.StackSize(), restored.StackSize())
			require.Equal(t, original.MaxStackSize(), restored.MaxStackSize())
			require.Equal(t, original.Value(), restored.Value())
			require.Equal(t, original.Weight(), restored.Weight())
			require.Equal(t, original.Icon(), restored.Icon())
		})
	})

	t.Run("Validation", func(t *testing.T) {
		t.Run("valid item passes validation", func(t *testing.T) {
			item := NewBaseItem("", TypeMaterial, "Valid")
			require.NoError(t, item.Validate())
		})

		t.Run("empty name fails validation", func(t *testing.T) {
			item := NewBaseItemWithConfig(BaseItemConfig{
				ItemType: TypeMaterial,
			})
			require.Error(t, item.Validate())
		})

		t.Run("empty type fails validation", func(t *testing.T) {
			item := NewBaseItemWithConfig(BaseItemConfig{
				Name: "NoType",
			})
			require.Error(t, item.Validate())
		})
	})
}
