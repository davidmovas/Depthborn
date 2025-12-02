package account

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidmovas/Depthborn/internal/item"
)

func createTestItem(id, name string) item.Item {
	return item.NewBaseItemWithConfig(item.BaseItemConfig{
		ID:       id,
		Name:     name,
		ItemType: item.TypeMaterial,
		Weight:   5.0,
	})
}

func createStackableItem(id, name string, maxStack int) item.Item {
	return item.NewBaseItemWithConfig(item.BaseItemConfig{
		ID:           id,
		Name:         name,
		ItemType:     item.TypeMaterial,
		Weight:       0.5,
		MaxStackSize: maxStack,
	})
}

func TestStash(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("with defaults", func(t *testing.T) {
			stash := NewStash(DefaultStashConfig())

			assert.NotNil(t, stash)
			assert.Equal(t, 1, stash.TabCount())
			assert.Equal(t, 10, stash.MaxTabs())
		})

		t.Run("with custom config", func(t *testing.T) {
			cfg := StashConfig{
				InitialTabs: 3,
				MaxTabs:     5,
				SlotsPerTab: 100,
			}
			stash := NewStash(cfg)

			assert.Equal(t, 3, stash.TabCount())
			assert.Equal(t, 5, stash.MaxTabs())

			tab, _ := stash.GetTab(0)
			assert.Equal(t, 100, tab.SlotCount())
		})
	})

	t.Run("Tab Management", func(t *testing.T) {
		t.Run("AddTab", func(t *testing.T) {
			stash := NewStash(DefaultStashConfig())

			err := stash.AddTab("New Tab")
			require.NoError(t, err)
			assert.Equal(t, 2, stash.TabCount())
		})

		t.Run("AddTab max reached returns error", func(t *testing.T) {
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     2,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			err := stash.AddTab("Third Tab")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "maximum number")
		})

		t.Run("AddTabWithSlots", func(t *testing.T) {
			stash := NewStash(DefaultStashConfig())

			err := stash.AddTabWithSlots("Big Tab", 200)
			require.NoError(t, err)

			tab, ok := stash.GetTab(1)
			assert.True(t, ok)
			assert.Equal(t, 200, tab.SlotCount())
		})

		t.Run("RemoveTab", func(t *testing.T) {
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			err := stash.RemoveTab(1)
			require.NoError(t, err)
			assert.Equal(t, 1, stash.TabCount())
		})

		t.Run("RemoveTab last tab returns error", func(t *testing.T) {
			stash := NewStash(DefaultStashConfig())

			err := stash.RemoveTab(0)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot remove the last")
		})

		t.Run("RemoveTab non-empty returns error", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tab, _ := stash.GetTab(0)
			_ = tab.Add(ctx, createTestItem("item-1", "Test"))

			err := stash.RemoveTab(0)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "non-empty")
		})

		t.Run("RemoveTab out of range returns error", func(t *testing.T) {
			stash := NewStash(DefaultStashConfig())

			err := stash.RemoveTab(5)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "out of range")
		})

		t.Run("GetTab", func(t *testing.T) {
			stash := NewStash(DefaultStashConfig())

			tab, ok := stash.GetTab(0)
			assert.True(t, ok)
			assert.NotNil(t, tab)

			_, ok = stash.GetTab(99)
			assert.False(t, ok)
		})

		t.Run("Tabs returns all tabs", func(t *testing.T) {
			cfg := StashConfig{
				InitialTabs: 3,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tabs := stash.Tabs()
			assert.Len(t, tabs, 3)
		})

		t.Run("RenameTab", func(t *testing.T) {
			stash := NewStash(DefaultStashConfig())

			err := stash.RenameTab(0, "My Stash")
			require.NoError(t, err)

			tab, _ := stash.GetTab(0)
			assert.Equal(t, "My Stash", tab.Name())
		})

		t.Run("RenameTab out of range returns error", func(t *testing.T) {
			stash := NewStash(DefaultStashConfig())

			err := stash.RenameTab(5, "Invalid")
			assert.Error(t, err)
		})

		t.Run("SwapTabs", func(t *testing.T) {
			cfg := StashConfig{
				InitialTabs: 3,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tab0Before, _ := stash.GetTab(0)
			tab2Before, _ := stash.GetTab(2)
			name0 := tab0Before.Name()
			name2 := tab2Before.Name()

			err := stash.SwapTabs(0, 2)
			require.NoError(t, err)

			tab0After, _ := stash.GetTab(0)
			tab2After, _ := stash.GetTab(2)
			assert.Equal(t, name2, tab0After.Name())
			assert.Equal(t, name0, tab2After.Name())
		})
	})

	t.Run("Item Operations", func(t *testing.T) {
		t.Run("FindItem", func(t *testing.T) {
			ctx := context.Background()
			stash := NewStash(DefaultStashConfig())

			tab, _ := stash.GetTab(0)
			_ = tab.Add(ctx, createTestItem("item-1", "Test"))

			itm, tabIdx, found := stash.FindItem("item-1")
			assert.True(t, found)
			assert.Equal(t, 0, tabIdx)
			assert.Equal(t, "item-1", itm.ID())

			_, _, found = stash.FindItem("nonexistent")
			assert.False(t, found)
		})

		t.Run("TransferToTab", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			itm := createTestItem("item-1", "Test")
			tab0, _ := stash.GetTab(0)
			_ = tab0.Add(ctx, itm)

			err := stash.TransferToTab(ctx, itm, 1)
			require.NoError(t, err)

			assert.False(t, tab0.Contains("item-1"))
			tab1, _ := stash.GetTab(1)
			assert.True(t, tab1.Contains("item-1"))
		})

		t.Run("TransferToTab out of range returns error", func(t *testing.T) {
			ctx := context.Background()
			stash := NewStash(DefaultStashConfig())

			itm := createTestItem("item-1", "Test")

			err := stash.TransferToTab(ctx, itm, 99)
			assert.Error(t, err)
		})

		t.Run("TransferToSlot", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			itm := createTestItem("item-1", "Test")
			tab0, _ := stash.GetTab(0)
			_ = tab0.Add(ctx, itm)

			err := stash.TransferToSlot(ctx, itm, 1, 10)
			require.NoError(t, err)

			tab1, _ := stash.GetTab(1)
			found, ok := tab1.GetAtSlot(10)
			assert.True(t, ok)
			assert.Equal(t, "item-1", found.ID())
		})
	})

	t.Run("Search and Filter", func(t *testing.T) {
		t.Run("Search across tabs", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tab0, _ := stash.GetTab(0)
			tab1, _ := stash.GetTab(1)

			_ = tab0.Add(ctx, createTestItem("sword-1", "Iron Sword"))
			_ = tab1.Add(ctx, createTestItem("shield-1", "Iron Shield"))
			_ = tab1.Add(ctx, createTestItem("bow-1", "Wooden Bow"))

			results := stash.Search("iron")
			assert.Len(t, results, 2)
		})

		t.Run("FindByType across tabs", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tab0, _ := stash.GetTab(0)
			tab1, _ := stash.GetTab(1)

			material := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "mat-1",
				Name:     "Material",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
			})
			consumable := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "cons-1",
				Name:     "Potion",
				ItemType: item.TypeConsumable,
				Weight:   1.0,
			})

			_ = tab0.Add(ctx, material)
			_ = tab1.Add(ctx, consumable)

			materials := stash.FindByType(item.TypeMaterial)
			assert.Len(t, materials, 1)
		})

		t.Run("FindByRarity across tabs", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tab0, _ := stash.GetTab(0)
			tab1, _ := stash.GetTab(1)

			common := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "common-1",
				Name:     "Common",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Rarity:   item.RarityCommon,
			})
			rare := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "rare-1",
				Name:     "Rare",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Rarity:   item.RarityRare,
			})

			_ = tab0.Add(ctx, common)
			_ = tab1.Add(ctx, rare)

			rares := stash.FindByRarity(item.RarityRare)
			assert.Len(t, rares, 1)
		})

		t.Run("FindByTag across tabs", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tab0, _ := stash.GetTab(0)
			tab1, _ := stash.GetTab(1)

			quest1 := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "quest-1",
				Name:     "Quest Item 1",
				ItemType: item.TypeMaterial,
				Weight:   1.0,
				Tags:     []string{"quest"},
			})
			quest2 := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "quest-2",
				Name:     "Quest Item 2",
				ItemType: item.TypeMaterial,
				Weight:   1.0,
				Tags:     []string{"quest"},
			})

			_ = tab0.Add(ctx, quest1)
			_ = tab1.Add(ctx, quest2)

			questItems := stash.FindByTag("quest")
			assert.Len(t, questItems, 2)
		})

		t.Run("Filter across tabs", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tab0, _ := stash.GetTab(0)
			tab1, _ := stash.GetTab(1)

			_ = tab0.Add(ctx, createTestItem("item-1", "A Item"))
			_ = tab1.Add(ctx, createTestItem("item-2", "B Item"))

			results := stash.Filter(func(i item.Item) bool {
				return i.Name()[0] == 'A'
			})
			assert.Len(t, results, 1)
		})
	})

	t.Run("Stats", func(t *testing.T) {
		t.Run("TotalSlots", func(t *testing.T) {
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 100,
			}
			stash := NewStash(cfg)

			// 2 tabs * 100 slots each
			assert.Equal(t, 200, stash.TotalSlots())
		})

		t.Run("TotalUsedSlots and TotalFreeSlots", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 100,
			}
			stash := NewStash(cfg)

			tab0, _ := stash.GetTab(0)
			tab1, _ := stash.GetTab(1)

			_ = tab0.Add(ctx, createTestItem("item-1", "Item 1"))
			_ = tab0.Add(ctx, createTestItem("item-2", "Item 2"))
			_ = tab1.Add(ctx, createTestItem("item-3", "Item 3"))

			assert.Equal(t, 3, stash.TotalUsedSlots())
			assert.Equal(t, 197, stash.TotalFreeSlots())
		})

		t.Run("TotalCount", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tab0, _ := stash.GetTab(0)
			tab1, _ := stash.GetTab(1)

			_ = tab0.Add(ctx, createTestItem("item-1", "Item 1"))
			_ = tab0.Add(ctx, createTestItem("item-2", "Item 2"))
			_ = tab1.Add(ctx, createTestItem("item-3", "Item 3"))

			assert.Equal(t, 3, stash.TotalCount())
		})

		t.Run("TotalItems with stacks", func(t *testing.T) {
			ctx := context.Background()
			stash := NewStash(DefaultStashConfig())

			tab, _ := stash.GetTab(0)
			stackable := createStackableItem("stack-1", "Potion", 10)
			stackable.AddStack(4) // Stack of 5
			_ = tab.Add(ctx, stackable)

			assert.Equal(t, 1, stash.TotalCount()) // 1 stack
			assert.Equal(t, 5, stash.TotalItems()) // 5 items
		})

		t.Run("TotalValue", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 60,
			}
			stash := NewStash(cfg)

			tab0, _ := stash.GetTab(0)
			tab1, _ := stash.GetTab(1)

			item1 := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "item-1",
				Name:     "Item 1",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Value:    100,
			})
			item2 := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "item-2",
				Name:     "Item 2",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Value:    200,
			})

			_ = tab0.Add(ctx, item1)
			_ = tab1.Add(ctx, item2)

			assert.Equal(t, int64(300), stash.TotalValue())
		})
	})

	t.Run("Persistence", func(t *testing.T) {
		t.Run("Serialization", func(t *testing.T) {
			ctx := context.Background()
			cfg := StashConfig{
				InitialTabs: 2,
				MaxTabs:     5,
				SlotsPerTab: 100,
			}
			stash := NewStash(cfg)

			tab0, _ := stash.GetTab(0)
			tab0.SetName("My Items")
			tab0.SetColor("#ff0000")
			_ = tab0.Add(ctx, createTestItem("item-1", "Item 1"))

			// Serialize
			state, err := stash.SerializeState()
			require.NoError(t, err)
			assert.NotNil(t, state)

			// Deserialize to new stash
			newStash := NewStash(DefaultStashConfig())
			err = newStash.DeserializeState(state)
			require.NoError(t, err)

			assert.Equal(t, 5, newStash.MaxTabs())
			assert.Equal(t, 2, newStash.TabCount())

			restoredTab, _ := newStash.GetTab(0)
			assert.Equal(t, "My Items", restoredTab.Name())
			assert.Equal(t, "#ff0000", restoredTab.Color())
		})
	})
}

func TestStashTab(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		tab := NewStashTab("Test Tab", 200)

		assert.NotNil(t, tab)
		assert.Equal(t, "Test Tab", tab.Name())
		assert.Equal(t, 200, tab.SlotCount())
		assert.Equal(t, "default", tab.Icon())
		assert.Equal(t, "#ffffff", tab.Color())
	})

	t.Run("Basic Operations", func(t *testing.T) {
		t.Run("Add", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			itm := createTestItem("item-1", "Test")
			err := tab.Add(ctx, itm)

			require.NoError(t, err)
			assert.Equal(t, 1, tab.ItemCount())
			assert.True(t, tab.Contains("item-1"))
		})

		t.Run("AddToSlot", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			itm := createTestItem("item-1", "Test")
			err := tab.AddToSlot(ctx, 50, itm)

			require.NoError(t, err)
			found, ok := tab.GetAtSlot(50)
			assert.True(t, ok)
			assert.Equal(t, "item-1", found.ID())
		})

		t.Run("Remove", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			itm := createTestItem("item-1", "Test")
			_ = tab.Add(ctx, itm)

			removed, err := tab.Remove(ctx, "item-1")
			require.NoError(t, err)
			assert.Equal(t, "item-1", removed.ID())
			assert.Equal(t, 0, tab.ItemCount())
		})

		t.Run("Get", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			itm := createTestItem("item-1", "Test")
			_ = tab.Add(ctx, itm)

			found, ok := tab.Get("item-1")
			assert.True(t, ok)
			assert.Equal(t, "item-1", found.ID())

			_, ok = tab.Get("nonexistent")
			assert.False(t, ok)
		})

		t.Run("GetAtSlot", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			itm := createTestItem("item-1", "Test")
			_ = tab.AddToSlot(ctx, 25, itm)

			found, ok := tab.GetAtSlot(25)
			assert.True(t, ok)
			assert.Equal(t, "item-1", found.ID())

			_, ok = tab.GetAtSlot(0)
			assert.False(t, ok)
		})

		t.Run("GetAll", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			_ = tab.Add(ctx, createTestItem("item-1", "Item 1"))
			_ = tab.Add(ctx, createTestItem("item-2", "Item 2"))

			all := tab.GetAll()
			assert.Len(t, all, 2)
		})

		t.Run("Clear", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			_ = tab.Add(ctx, createTestItem("item-1", "Item 1"))
			_ = tab.Add(ctx, createTestItem("item-2", "Item 2"))

			cleared := tab.Clear(ctx)
			assert.Len(t, cleared, 2)
			assert.Equal(t, 0, tab.ItemCount())
		})
	})

	t.Run("Slot Operations", func(t *testing.T) {
		t.Run("SwapSlots", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			_ = tab.AddToSlot(ctx, 0, createTestItem("item-1", "Item 1"))
			_ = tab.AddToSlot(ctx, 50, createTestItem("item-2", "Item 2"))

			err := tab.SwapSlots(ctx, 0, 50)
			require.NoError(t, err)

			itm0, _ := tab.GetAtSlot(0)
			itm50, _ := tab.GetAtSlot(50)
			assert.Equal(t, "item-2", itm0.ID())
			assert.Equal(t, "item-1", itm50.ID())
		})

		t.Run("MoveToSlot", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			_ = tab.Add(ctx, createTestItem("item-1", "Item 1"))

			err := tab.MoveToSlot(ctx, "item-1", 75)
			require.NoError(t, err)

			itm, ok := tab.GetAtSlot(75)
			assert.True(t, ok)
			assert.Equal(t, "item-1", itm.ID())

			_, ok = tab.GetAtSlot(0)
			assert.False(t, ok)
		})

		t.Run("SetSlotCount", func(t *testing.T) {
			t.Run("expand", func(t *testing.T) {
				tab := NewStashTab("Test Tab", 50)

				tab.SetSlotCount(100)
				assert.Equal(t, 100, tab.SlotCount())
			})

			t.Run("shrink when empty slots", func(t *testing.T) {
				tab := NewStashTab("Test Tab", 100)

				tab.SetSlotCount(50)
				assert.Equal(t, 50, tab.SlotCount())
			})
		})

		t.Run("UsedSlots and FreeSlots", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			assert.Equal(t, 0, tab.UsedSlots())
			assert.Equal(t, 100, tab.FreeSlots())

			_ = tab.Add(ctx, createTestItem("item-1", "Test"))
			_ = tab.Add(ctx, createTestItem("item-2", "Test"))

			assert.Equal(t, 2, tab.UsedSlots())
			assert.Equal(t, 98, tab.FreeSlots())
		})

		t.Run("IsFull", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 2)

			assert.False(t, tab.IsFull())

			_ = tab.Add(ctx, createTestItem("item-1", "Test"))
			_ = tab.Add(ctx, createTestItem("item-2", "Test"))

			assert.True(t, tab.IsFull())
		})
	})

	t.Run("Stack Operations", func(t *testing.T) {
		t.Run("auto-stacking on add", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			item1 := createStackableItem("item-1", "Potion", 10)
			item2 := createStackableItem("item-2", "Potion", 10)

			_ = tab.Add(ctx, item1)
			_ = tab.Add(ctx, item2)

			assert.Equal(t, 1, tab.ItemCount())
			found, _ := tab.Get("item-1")
			assert.Equal(t, 2, found.StackSize())
		})

		t.Run("SplitStack", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			item1 := createStackableItem("item-1", "Potion", 10)
			item1.AddStack(4) // Stack of 5
			_ = tab.Add(ctx, item1)

			newStack, err := tab.SplitStack(ctx, "item-1", 2)
			require.NoError(t, err)
			assert.NotNil(t, newStack)
			assert.Equal(t, 2, newStack.StackSize())

			original, _ := tab.Get("item-1")
			assert.Equal(t, 3, original.StackSize())
			assert.Equal(t, 2, tab.ItemCount())
		})

		t.Run("MergeStacks", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			item1 := createStackableItem("item-1", "Potion", 10)
			item1.AddStack(2) // Stack of 3
			item2 := createStackableItem("item-2", "Potion", 10)
			item2.AddStack(1) // Stack of 2

			_ = tab.AddToSlot(ctx, 0, item1)
			_ = tab.AddToSlot(ctx, 1, item2)

			err := tab.MergeStacks(ctx, "item-2", "item-1")
			require.NoError(t, err)

			merged, _ := tab.Get("item-1")
			assert.Equal(t, 5, merged.StackSize())
			assert.Equal(t, 1, tab.ItemCount())
		})

		t.Run("RemoveAmount", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			item1 := createStackableItem("item-1", "Potion", 10)
			item1.AddStack(4) // Stack of 5
			_ = tab.Add(ctx, item1)

			removed, err := tab.RemoveAmount(ctx, "item-1", 2)
			require.NoError(t, err)
			assert.NotNil(t, removed)

			remaining, _ := tab.Get("item-1")
			assert.Equal(t, 3, remaining.StackSize())
		})

		t.Run("CanStackWith", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			item1 := createStackableItem("item-1", "Potion", 10)
			_ = tab.Add(ctx, item1)

			item2 := createStackableItem("item-2", "Potion", 10)
			targetID, canStack := tab.CanStackWith(item2)
			assert.True(t, canStack)
			assert.Equal(t, "item-1", targetID)
		})

		t.Run("TotalItems with stacks", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			item1 := createStackableItem("item-1", "Potion", 10)
			item1.AddStack(4) // Stack of 5
			_ = tab.Add(ctx, item1)

			assert.Equal(t, 1, tab.ItemCount())  // 1 stack
			assert.Equal(t, 5, tab.TotalItems()) // 5 items
		})
	})

	t.Run("Search and Filter", func(t *testing.T) {
		t.Run("Search by name", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			_ = tab.Add(ctx, createTestItem("sword-1", "Iron Sword"))
			_ = tab.Add(ctx, createTestItem("shield-1", "Iron Shield"))
			_ = tab.Add(ctx, createTestItem("bow-1", "Wooden Bow"))

			results := tab.Search("iron")
			assert.Len(t, results, 2)
		})

		t.Run("FindByType", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			material := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "mat-1",
				Name:     "Material",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
			})
			consumable := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "cons-1",
				Name:     "Potion",
				ItemType: item.TypeConsumable,
				Weight:   1.0,
			})

			_ = tab.Add(ctx, material)
			_ = tab.Add(ctx, consumable)

			materials := tab.FindByType(item.TypeMaterial)
			assert.Len(t, materials, 1)
		})

		t.Run("FindByRarity", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			common := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "common-1",
				Name:     "Common",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Rarity:   item.RarityCommon,
			})
			rare := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "rare-1",
				Name:     "Rare",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Rarity:   item.RarityRare,
			})

			_ = tab.Add(ctx, common)
			_ = tab.Add(ctx, rare)

			rares := tab.FindByRarity(item.RarityRare)
			assert.Len(t, rares, 1)
		})

		t.Run("FindByTag", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			tagged := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "tagged-1",
				Name:     "Tagged Item",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Tags:     []string{"quest"},
			})
			untagged := createTestItem("untagged-1", "Untagged Item")

			_ = tab.Add(ctx, tagged)
			_ = tab.Add(ctx, untagged)

			questItems := tab.FindByTag("quest")
			assert.Len(t, questItems, 1)
		})

		t.Run("FindStackable", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			stackable := createStackableItem("stack-1", "Stackable", 10)
			nonStackable := createTestItem("single-1", "Single")

			_ = tab.Add(ctx, stackable)
			_ = tab.Add(ctx, nonStackable)

			stackables := tab.FindStackable()
			assert.Len(t, stackables, 1)
		})

		t.Run("Filter with predicate", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			_ = tab.Add(ctx, createTestItem("item-a", "A Item"))
			_ = tab.Add(ctx, createTestItem("item-b", "B Item"))

			results := tab.Filter(func(i item.Item) bool {
				return i.Name()[0] == 'A'
			})
			assert.Len(t, results, 1)
		})
	})

	t.Run("Metadata", func(t *testing.T) {
		t.Run("SetName", func(t *testing.T) {
			tab := NewStashTab("Original", 100)

			tab.SetName("Updated")
			assert.Equal(t, "Updated", tab.Name())
		})

		t.Run("SetIcon", func(t *testing.T) {
			tab := NewStashTab("Test", 100)

			tab.SetIcon("weapon")
			assert.Equal(t, "weapon", tab.Icon())
		})

		t.Run("SetColor", func(t *testing.T) {
			tab := NewStashTab("Test", 100)

			tab.SetColor("#00ff00")
			assert.Equal(t, "#00ff00", tab.Color())
		})
	})

	t.Run("Stats", func(t *testing.T) {
		t.Run("TotalValue", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 100)

			item1 := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "item-1",
				Name:     "Item 1",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Value:    100,
			})
			item2 := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "item-2",
				Name:     "Item 2",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Value:    200,
			})

			_ = tab.Add(ctx, item1)
			_ = tab.Add(ctx, item2)

			assert.Equal(t, int64(300), tab.TotalValue())
		})
	})

	t.Run("Persistence", func(t *testing.T) {
		t.Run("ToState", func(t *testing.T) {
			tab := NewStashTab("Test Tab", 100)
			tab.SetIcon("armor")
			tab.SetColor("#0000ff")

			state := tab.ToState()

			assert.Equal(t, "Test Tab", state.Name)
			assert.Equal(t, "armor", state.Icon)
			assert.Equal(t, "#0000ff", state.Color)
			assert.Equal(t, 100, state.Slots)
		})

		t.Run("FromState", func(t *testing.T) {
			state := StashTabState{
				Name:  "Restored Tab",
				Icon:  "gem",
				Color: "#ff00ff",
				Slots: 150,
			}

			tab := StashTabFromState(state)

			assert.Equal(t, "Restored Tab", tab.Name())
			assert.Equal(t, "gem", tab.Icon())
			assert.Equal(t, "#ff00ff", tab.Color())
			assert.Equal(t, 150, tab.SlotCount())
		})

		t.Run("GetItemIDs", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 10)

			_ = tab.AddToSlot(ctx, 0, createTestItem("item-1", "Item 1"))
			_ = tab.AddToSlot(ctx, 5, createTestItem("item-2", "Item 2"))

			ids := tab.GetItemIDs()
			assert.Len(t, ids, 10)
			assert.Equal(t, "item-1", ids[0])
			assert.Equal(t, "", ids[1])
			assert.Equal(t, "item-2", ids[5])
		})

		t.Run("AddDirect", func(t *testing.T) {
			tab := NewStashTab("Test Tab", 100)

			itm := createTestItem("item-1", "Test Item")
			err := tab.AddDirect(itm)

			require.NoError(t, err)
			assert.Equal(t, 1, tab.ItemCount())
		})

		t.Run("AddDirectToSlot", func(t *testing.T) {
			tab := NewStashTab("Test Tab", 100)

			itm := createTestItem("item-1", "Test Item")
			err := tab.AddDirectToSlot(25, itm)

			require.NoError(t, err)
			found, ok := tab.GetAtSlot(25)
			assert.True(t, ok)
			assert.Equal(t, "item-1", found.ID())
		})

		t.Run("CanAdd", func(t *testing.T) {
			ctx := context.Background()
			tab := NewStashTab("Test Tab", 2)

			itm1 := createTestItem("item-1", "Test 1")
			itm2 := createTestItem("item-2", "Test 2")
			itm3 := createTestItem("item-3", "Test 3")

			assert.True(t, tab.CanAdd(itm1))
			_ = tab.Add(ctx, itm1)
			_ = tab.Add(ctx, itm2)
			assert.False(t, tab.CanAdd(itm3))
		})
	})
}
