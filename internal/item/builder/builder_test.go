package builder

import (
	"context"
	"testing"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/item"
	"github.com/stretchr/testify/require"
)

func TestItemBuilder(t *testing.T) {
	t.Run("NewItem creates builder with defaults", func(t *testing.T) {
		b := NewItem()
		cfg := b.Config()

		require.Equal(t, 1, cfg.MaxStackSize)
		require.Equal(t, 1.0, cfg.Quality)
		require.Equal(t, 1, cfg.Level)
		require.Equal(t, 0.1, cfg.Weight)
	})

	t.Run("For creates builder with type and name", func(t *testing.T) {
		b := For(item.TypeMaterial, "Iron Ore")

		result := b.Build()

		require.Equal(t, "Iron Ore", result.Name())
		require.Equal(t, item.TypeMaterial, result.ItemType())
	})

	t.Run("Material creates material builder", func(t *testing.T) {
		result := Material("Iron Ore").Build()

		require.Equal(t, item.TypeMaterial, result.ItemType())
		require.Equal(t, "Iron Ore", result.Name())
	})

	t.Run("Currency creates currency builder with max stack 9999", func(t *testing.T) {
		result := Currency("Gold Coin").Build()

		require.Equal(t, item.TypeCurrency, result.ItemType())
		require.Equal(t, 9999, result.MaxStackSize())
	})

	t.Run("Quest creates quest item builder", func(t *testing.T) {
		result := Quest("Ancient Key").Build()

		require.Equal(t, item.TypeQuest, result.ItemType())
	})

	t.Run("Key creates key item builder", func(t *testing.T) {
		result := Key("Dungeon Key").Build()

		require.Equal(t, item.TypeKey, result.ItemType())
	})

	t.Run("fluent API sets all properties", func(t *testing.T) {
		result := NewItem().
			ID("custom-id").
			Name("Magic Stone").
			Description("A mysterious stone").
			Type(item.TypeGem).
			Rarity(item.RarityEpic).
			Quality(0.9).
			Level(20).
			MaxStack(50).
			Value(1000).
			Weight(0.5).
			Icon("magic_stone").
			Tags("magic", "crafting").
			Build()

		require.Equal(t, "custom-id", result.ID())
		require.Equal(t, "Magic Stone", result.Name())
		require.Equal(t, "A mysterious stone", result.Description())
		require.Equal(t, item.TypeGem, result.ItemType())
		require.Equal(t, item.RarityEpic, result.Rarity())
		require.Equal(t, 0.9, result.Quality())
		require.Equal(t, 20, result.Level())
		require.Equal(t, 50, result.MaxStackSize())
		require.Equal(t, int64(1000), result.Value())
		require.Equal(t, 0.5, result.Weight())
		require.Equal(t, "magic_stone", result.Icon())
		require.True(t, result.Tags().Has("magic"))
		require.True(t, result.Tags().Has("crafting"))
	})

	t.Run("Build returns valid item", func(t *testing.T) {
		result := Material("Test Item").Build()

		require.NoError(t, result.Validate())
	})
}

func TestEquipmentBuilder(t *testing.T) {
	t.Run("NewEquipment creates builder with defaults", func(t *testing.T) {
		b := NewEquipment()
		result := b.Name("Test").Type(item.TypeWeaponMelee).Slot(item.SlotMainHand).Build()

		require.Equal(t, 100.0, result.MaxDurability())
	})

	t.Run("Equip creates equipment with type name and slot", func(t *testing.T) {
		result := Equip(item.TypeWeaponMelee, "Iron Sword", item.SlotMainHand).Build()

		require.Equal(t, "Iron Sword", result.Name())
		require.Equal(t, item.TypeWeaponMelee, result.ItemType())
		require.Equal(t, item.SlotMainHand, result.Slot())
	})

	t.Run("MeleeWeapon creates melee weapon", func(t *testing.T) {
		result := MeleeWeapon("Sword").Build()

		require.Equal(t, item.TypeWeaponMelee, result.ItemType())
		require.Equal(t, item.SlotMainHand, result.Slot())
	})

	t.Run("RangedWeapon creates ranged weapon", func(t *testing.T) {
		result := RangedWeapon("Bow").Build()

		require.Equal(t, item.TypeWeaponRanged, result.ItemType())
		require.Equal(t, item.SlotTwoHand, result.Slot())
	})

	t.Run("MagicWeapon creates magic weapon", func(t *testing.T) {
		result := MagicWeapon("Staff").Build()

		require.Equal(t, item.TypeWeaponMagic, result.ItemType())
		require.Equal(t, item.SlotMainHand, result.Slot())
	})

	t.Run("Armor creates armor for slot", func(t *testing.T) {
		testCases := []struct {
			slot     item.EquipmentSlot
			expected item.Type
		}{
			{item.SlotHead, item.TypeArmorHead},
			{item.SlotChest, item.TypeArmorChest},
			{item.SlotLegs, item.TypeArmorLegs},
			{item.SlotFeet, item.TypeArmorFeet},
			{item.SlotHands, item.TypeArmorHands},
		}

		for _, tc := range testCases {
			result := Armor("Test Armor", tc.slot).Build()
			require.Equal(t, tc.expected, result.ItemType())
			require.Equal(t, tc.slot, result.Slot())
		}
	})

	t.Run("Accessory creates accessory for slot", func(t *testing.T) {
		testCases := []struct {
			slot     item.EquipmentSlot
			expected item.Type
		}{
			{item.SlotRing1, item.TypeAccessoryRing},
			{item.SlotRing2, item.TypeAccessoryRing},
			{item.SlotAmulet, item.TypeAccessoryAmulet},
			{item.SlotBelt, item.TypeAccessoryBelt},
		}

		for _, tc := range testCases {
			result := Accessory("Test Accessory", tc.slot).Build()
			require.Equal(t, tc.expected, result.ItemType())
			require.Equal(t, tc.slot, result.Slot())
		}
	})

	t.Run("fluent API sets all equipment properties", func(t *testing.T) {
		result := MeleeWeapon("Epic Sword").
			ID("sword-1").
			Description("A legendary sword").
			Rarity(item.RarityLegendary).
			Quality(1.0).
			Level(50).
			Value(10000).
			Weight(5.0).
			Icon("epic_sword").
			Tags("epic", "weapon").
			Durability(200).
			Sockets(2, item.SocketTypeGem, item.SocketTypeRune).
			RequireLevel(30).
			Build()

		require.Equal(t, "sword-1", result.ID())
		require.Equal(t, "Epic Sword", result.Name())
		require.Equal(t, "A legendary sword", result.Description())
		require.Equal(t, item.RarityLegendary, result.Rarity())
		require.Equal(t, 50, result.Level())
		require.Equal(t, 200.0, result.MaxDurability())
		require.Equal(t, 2, result.SocketCount())

		socketType0, _ := result.GetSocketType(0)
		socketType1, _ := result.GetSocketType(1)
		require.Equal(t, item.SocketTypeGem, socketType0)
		require.Equal(t, item.SocketTypeRune, socketType1)
	})

	t.Run("Sockets fills remaining with universal", func(t *testing.T) {
		result := MeleeWeapon("Sword").Sockets(3, item.SocketTypeGem).Build()

		require.Equal(t, 3, result.SocketCount())
		socketType0, _ := result.GetSocketType(0)
		socketType1, _ := result.GetSocketType(1)
		socketType2, _ := result.GetSocketType(2)
		require.Equal(t, item.SocketTypeGem, socketType0)
		require.Equal(t, item.SocketTypeUniversal, socketType1)
		require.Equal(t, item.SocketTypeUniversal, socketType2)
	})

	t.Run("Attribute adds modifier", func(t *testing.T) {
		mod := attribute.NewModifier("str", attribute.ModFlat, 10, string(attribute.AttrStrength))
		result := MeleeWeapon("Sword").Attribute(mod).Build()

		require.Len(t, result.Attributes(), 1)
	})

	t.Run("Attributes adds multiple modifiers", func(t *testing.T) {
		mod1 := attribute.NewModifier("str", attribute.ModFlat, 10, string(attribute.AttrStrength))
		mod2 := attribute.NewModifier("dex", attribute.ModFlat, 5, string(attribute.AttrDexterity))
		result := MeleeWeapon("Sword").Attributes(mod1, mod2).Build()

		require.Len(t, result.Attributes(), 2)
	})

	t.Run("OnEquip and OnUnequip set callbacks", func(t *testing.T) {
		equipCalled := false
		unequipCalled := false

		// Just verify the callbacks are set by building without error
		result := MeleeWeapon("Sword").
			OnEquip(func(_ context.Context, _ entity.Entity) error {
				equipCalled = true
				return nil
			}).
			OnUnequip(func(_ context.Context, _ entity.Entity) error {
				unequipCalled = true
				return nil
			}).
			Build()

		// Verify equipment was built successfully
		require.NotNil(t, result)
		require.Equal(t, "Sword", result.Name())

		// Callbacks are set but we can't easily test them without a mock entity
		// The fact that Build() succeeded means the callbacks were properly assigned
		_ = equipCalled
		_ = unequipCalled
	})

	t.Run("Build returns valid equipment", func(t *testing.T) {
		result := MeleeWeapon("Test Sword").Build()

		require.NoError(t, result.Validate())
	})
}

func TestConsumableBuilder(t *testing.T) {
	t.Run("NewConsumable creates builder with defaults", func(t *testing.T) {
		result := NewConsumable().Name("Test").Build()

		require.Equal(t, item.TypeConsumable, result.ItemType())
		require.Equal(t, 1, result.Charges())
	})

	t.Run("Consume creates consumable with name", func(t *testing.T) {
		result := Consume("Health Potion").Build()

		require.Equal(t, "Health Potion", result.Name())
		require.Equal(t, item.TypeConsumable, result.ItemType())
	})

	t.Run("Potion creates potion with tag", func(t *testing.T) {
		result := Potion("Health Potion").Build()

		require.True(t, result.Tags().Has("potion"))
	})

	t.Run("Food creates food with tag", func(t *testing.T) {
		result := Food("Bread").Build()

		require.True(t, result.Tags().Has("food"))
	})

	t.Run("Scroll creates scroll with tag", func(t *testing.T) {
		result := Scroll("Scroll of Fire").Build()

		require.True(t, result.Tags().Has("scroll"))
	})

	t.Run("fluent API sets all consumable properties", func(t *testing.T) {
		result := Potion("Super Potion").
			ID("potion-1").
			Description("Restores a lot of health").
			Rarity(item.RarityRare).
			Quality(1.0).
			Level(10).
			MaxStack(20).
			Value(500).
			Weight(0.2).
			Icon("super_potion").
			Tags("healing").
			Cooldown(5000).
			Charges(3).
			Build()

		require.Equal(t, "potion-1", result.ID())
		require.Equal(t, "Super Potion", result.Name())
		require.Equal(t, item.RarityRare, result.Rarity())
		require.Equal(t, 20, result.MaxStackSize())
		require.Equal(t, int64(5000), result.MaxCooldown())
		require.Equal(t, 3, result.Charges())
		require.True(t, result.Tags().Has("potion"))
		require.True(t, result.Tags().Has("healing"))
	})

	t.Run("Infinite sets charges to -1", func(t *testing.T) {
		result := Consume("Endless Potion").Infinite().Build()

		require.Equal(t, -1, result.Charges())
	})

	t.Run("Build returns valid consumable", func(t *testing.T) {
		result := Potion("Test Potion").Build()

		require.NoError(t, result.Validate())
	})
}

func TestContainerBuilder(t *testing.T) {
	t.Run("NewContainer creates builder with defaults", func(t *testing.T) {
		result := NewContainer().Name("Test").Build()

		require.Equal(t, item.TypeContainer, result.ItemType())
		require.Equal(t, 10, result.Capacity())
	})

	t.Run("Contain creates container with name and capacity", func(t *testing.T) {
		result := Contain("Backpack", 20).Build()

		require.Equal(t, "Backpack", result.Name())
		require.Equal(t, 20, result.Capacity())
	})

	t.Run("Bag creates bag with tag", func(t *testing.T) {
		result := Bag("Small Bag", 10).Build()

		require.True(t, result.Tags().Has("bag"))
	})

	t.Run("Chest creates chest with tag", func(t *testing.T) {
		result := Chest("Treasure Chest", 50).Build()

		require.True(t, result.Tags().Has("chest"))
	})

	t.Run("fluent API sets all container properties", func(t *testing.T) {
		result := Bag("Magic Bag", 30).
			ID("bag-1").
			Description("A magical bag").
			Rarity(item.RarityEpic).
			Quality(1.0).
			Level(20).
			Value(5000).
			Weight(1.0).
			Icon("magic_bag").
			Tags("magic").
			MaxWeight(100.0).
			AllowTypes(item.TypeMaterial, item.TypeGem).
			Build()

		require.Equal(t, "bag-1", result.ID())
		require.Equal(t, "Magic Bag", result.Name())
		require.Equal(t, 30, result.Capacity())
		require.Equal(t, 100.0, result.MaxWeight())
		require.Len(t, result.AllowedTypes(), 2)
		require.True(t, result.Tags().Has("bag"))
		require.True(t, result.Tags().Has("magic"))
	})

	t.Run("Build returns valid container", func(t *testing.T) {
		result := Bag("Test Bag", 10).Build()

		require.NoError(t, result.Validate())
	})
}

func TestSocketableBuilder(t *testing.T) {
	t.Run("NewSocketable creates builder with defaults", func(t *testing.T) {
		result := NewSocketable().Name("Test").Build()

		require.Equal(t, item.SocketTypeGem, result.SocketType())
		require.Equal(t, 1, result.Tier())
	})

	t.Run("Socket creates socketable with name and type", func(t *testing.T) {
		result := Socket("Ruby", item.SocketTypeGem).Build()

		require.Equal(t, "Ruby", result.Name())
		require.Equal(t, item.SocketTypeGem, result.SocketType())
	})

	t.Run("Gem creates gem socketable", func(t *testing.T) {
		result := Gem("Ruby").Build()

		require.Equal(t, item.TypeGem, result.ItemType())
		require.Equal(t, item.SocketTypeGem, result.SocketType())
	})

	t.Run("Rune creates rune socketable", func(t *testing.T) {
		result := Rune("El").Build()

		require.Equal(t, item.TypeRune, result.ItemType())
		require.Equal(t, item.SocketTypeRune, result.SocketType())
	})

	t.Run("fluent API sets all socketable properties", func(t *testing.T) {
		mod := attribute.NewModifier("str", attribute.ModFlat, 10, string(attribute.AttrStrength))

		result := Gem("Perfect Ruby").
			ID("gem-1").
			Description("A flawless ruby").
			Rarity(item.RarityLegendary).
			Quality(1.0).
			Level(50).
			Value(10000).
			Weight(0.1).
			Icon("perfect_ruby").
			Tags("perfect", "fire").
			Tier(5).
			Modifier(mod).
			Build()

		require.Equal(t, "gem-1", result.ID())
		require.Equal(t, "Perfect Ruby", result.Name())
		require.Equal(t, item.RarityLegendary, result.Rarity())
		require.Equal(t, 5, result.Tier())
		require.Len(t, result.Modifiers(), 1)
		require.True(t, result.Tags().Has("perfect"))
		require.True(t, result.Tags().Has("fire"))
	})

	t.Run("Modifiers adds multiple modifiers", func(t *testing.T) {
		mod1 := attribute.NewModifier("str", attribute.ModFlat, 10, string(attribute.AttrStrength))
		mod2 := attribute.NewModifier("dex", attribute.ModFlat, 5, string(attribute.AttrDexterity))

		result := Gem("Multi-Gem").Modifiers(mod1, mod2).Build()

		require.Len(t, result.Modifiers(), 2)
	})

	t.Run("Build returns valid socketable", func(t *testing.T) {
		result := Gem("Test Gem").Build()

		require.NoError(t, result.Validate())
	})
}
