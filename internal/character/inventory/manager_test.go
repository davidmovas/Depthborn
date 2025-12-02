package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidmovas/Depthborn/internal/item"
)

func createTestItem(id, name string, weight float64) item.Item {
	return item.NewBaseItemWithConfig(item.BaseItemConfig{
		ID:       id,
		Name:     name,
		ItemType: item.TypeMaterial,
		Weight:   weight,
	})
}

func createStackableItem(id, name string, weight float64, maxStack int) item.Item {
	return item.NewBaseItemWithConfig(item.BaseItemConfig{
		ID:           id,
		Name:         name,
		ItemType:     item.TypeMaterial,
		Weight:       weight,
		MaxStackSize: maxStack,
	})
}

func TestManager(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("with defaults", func(t *testing.T) {
			mgr := NewManager()
			assert.NotNil(t, mgr)
			assert.Equal(t, 100.0, mgr.MaxWeight())
			assert.Equal(t, 20, mgr.SlotCount())
			assert.Equal(t, 0, mgr.Count())
		})

		t.Run("with custom config", func(t *testing.T) {
			mgr := NewManagerWithConfig(Config{MaxWeight: 200, MaxSlots: 30})
			assert.Equal(t, 200.0, mgr.MaxWeight())
			assert.Equal(t, 30, mgr.SlotCount())
		})
	})

	t.Run("Basic Operations", func(t *testing.T) {
		t.Run("Add", func(t *testing.T) {
			t.Run("single item", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

				itm := createTestItem("item-1", "Test Item", 10.0)
				err := mgr.Add(ctx, itm)

				require.NoError(t, err)
				assert.Equal(t, 1, mgr.Count())
				assert.Equal(t, 10.0, mgr.CurrentWeight())
				assert.True(t, mgr.Contains("item-1"))
			})

			t.Run("nil item returns error", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManager()

				err := mgr.Add(ctx, nil)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "nil item")
			})

			t.Run("weight exceeded returns error", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManagerWithConfig(Config{MaxWeight: 50, MaxSlots: 10})

				itm := createTestItem("item-1", "Heavy Item", 60.0)
				err := mgr.Add(ctx, itm)

				assert.Error(t, err)
				assert.Contains(t, err.Error(), "weight limit exceeded")
				assert.Equal(t, 0, mgr.Count())
			})

			t.Run("to specific slot", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

				itm := createTestItem("item-1", "Test", 10.0)
				err := mgr.AddToSlot(ctx, 5, itm)

				require.NoError(t, err)
				found, ok := mgr.GetAtSlot(5)
				assert.True(t, ok)
				assert.Equal(t, "item-1", found.ID())
			})

			t.Run("to occupied slot returns error", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

				itm1 := createTestItem("item-1", "Item 1", 10.0)
				itm2 := createTestItem("item-2", "Item 2", 10.0)

				_ = mgr.AddToSlot(ctx, 0, itm1)
				err := mgr.AddToSlot(ctx, 0, itm2)

				assert.Error(t, err)
				assert.Contains(t, err.Error(), "already occupied")
			})
		})

		t.Run("Remove", func(t *testing.T) {
			t.Run("existing item", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManager()

				itm := createTestItem("item-1", "Test Item", 10.0)
				_ = mgr.Add(ctx, itm)

				removed, err := mgr.Remove(ctx, "item-1")
				require.NoError(t, err)
				assert.Equal(t, "item-1", removed.ID())
				assert.Equal(t, 0, mgr.Count())
				assert.Equal(t, 0.0, mgr.CurrentWeight())
			})

			t.Run("non-existent item returns error", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManager()

				_, err := mgr.Remove(ctx, "nonexistent")
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
			})
		})

		t.Run("Get", func(t *testing.T) {
			t.Run("by ID", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManager()

				itm := createTestItem("item-1", "Test Item", 10.0)
				_ = mgr.Add(ctx, itm)

				found, ok := mgr.Get("item-1")
				assert.True(t, ok)
				assert.Equal(t, "item-1", found.ID())

				_, ok = mgr.Get("nonexistent")
				assert.False(t, ok)
			})

			t.Run("at slot", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

				itm := createTestItem("item-1", "Test", 10.0)
				_ = mgr.AddToSlot(ctx, 3, itm)

				found, ok := mgr.GetAtSlot(3)
				assert.True(t, ok)
				assert.Equal(t, "item-1", found.ID())

				_, ok = mgr.GetAtSlot(0)
				assert.False(t, ok)
			})

			t.Run("all items", func(t *testing.T) {
				ctx := context.Background()
				mgr := NewManager()

				_ = mgr.Add(ctx, createTestItem("item-1", "Item 1", 10.0))
				_ = mgr.Add(ctx, createTestItem("item-2", "Item 2", 15.0))

				all := mgr.GetAll()
				assert.Len(t, all, 2)
			})
		})

		t.Run("Clear", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

			_ = mgr.Add(ctx, createTestItem("item-1", "Item 1", 10.0))
			_ = mgr.Add(ctx, createTestItem("item-2", "Item 2", 15.0))

			cleared := mgr.Clear(ctx)
			assert.Len(t, cleared, 2)
			assert.Equal(t, 0, mgr.Count())
			assert.Equal(t, 0.0, mgr.CurrentWeight())
		})
	})

	t.Run("Slot Operations", func(t *testing.T) {
		t.Run("SwapSlots", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			_ = mgr.AddToSlot(ctx, 0, createTestItem("item-1", "Item 1", 10.0))
			_ = mgr.AddToSlot(ctx, 5, createTestItem("item-2", "Item 2", 10.0))

			err := mgr.SwapSlots(ctx, 0, 5)
			require.NoError(t, err)

			itm0, _ := mgr.GetAtSlot(0)
			itm5, _ := mgr.GetAtSlot(5)
			assert.Equal(t, "item-2", itm0.ID())
			assert.Equal(t, "item-1", itm5.ID())
		})

		t.Run("MoveToSlot", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			_ = mgr.Add(ctx, createTestItem("item-1", "Item 1", 10.0))

			err := mgr.MoveToSlot(ctx, "item-1", 7)
			require.NoError(t, err)

			itm, ok := mgr.GetAtSlot(7)
			assert.True(t, ok)
			assert.Equal(t, "item-1", itm.ID())

			_, ok = mgr.GetAtSlot(0)
			assert.False(t, ok)
		})

		t.Run("SetSlotCount", func(t *testing.T) {
			t.Run("expand", func(t *testing.T) {
				mgr := NewManagerWithConfig(Config{MaxSlots: 10, MaxWeight: 100})

				mgr.SetSlotCount(20)
				assert.Equal(t, 20, mgr.SlotCount())
			})

			t.Run("shrink when empty slots", func(t *testing.T) {
				mgr := NewManagerWithConfig(Config{MaxSlots: 20, MaxWeight: 100})

				mgr.SetSlotCount(10)
				assert.Equal(t, 10, mgr.SlotCount())
			})
		})
	})

	t.Run("Weight Management", func(t *testing.T) {
		t.Run("tracks weight correctly", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			assert.Equal(t, 100.0, mgr.AvailableWeight())
			assert.Equal(t, 0.0, mgr.WeightPercent())
			assert.False(t, mgr.IsFull())

			_ = mgr.Add(ctx, createTestItem("item-1", "Test", 50.0))

			assert.Equal(t, 50.0, mgr.AvailableWeight())
			assert.Equal(t, 0.5, mgr.WeightPercent())
			assert.False(t, mgr.IsFull())
		})

		t.Run("SetMaxWeight", func(t *testing.T) {
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			mgr.SetMaxWeight(200)
			assert.Equal(t, 200.0, mgr.MaxWeight())

			// Negative value should be ignored
			mgr.SetMaxWeight(-50)
			assert.Equal(t, 200.0, mgr.MaxWeight())
		})
	})

	t.Run("Capacity Checks", func(t *testing.T) {
		t.Run("CanAdd", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 50, MaxSlots: 10})

			lightItem := createTestItem("light", "Light", 10.0)
			assert.True(t, mgr.CanAdd(lightItem))

			heavyItem := createTestItem("heavy", "Heavy", 60.0)
			assert.False(t, mgr.CanAdd(heavyItem))

			_ = mgr.Add(ctx, lightItem)
			assert.False(t, mgr.CanAdd(heavyItem))
		})

		t.Run("Contains", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

			_ = mgr.Add(ctx, createTestItem("item-1", "Test", 10.0))

			assert.True(t, mgr.Contains("item-1"))
			assert.False(t, mgr.Contains("nonexistent"))
		})

		t.Run("UsedSlots and FreeSlots", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxSlots: 10, MaxWeight: 100})

			assert.Equal(t, 0, mgr.UsedSlots())
			assert.Equal(t, 10, mgr.FreeSlots())

			_ = mgr.Add(ctx, createTestItem("item-1", "Test", 10.0))
			_ = mgr.Add(ctx, createTestItem("item-2", "Test", 10.0))

			assert.Equal(t, 2, mgr.UsedSlots())
			assert.Equal(t, 8, mgr.FreeSlots())
		})
	})

	t.Run("Stack Operations", func(t *testing.T) {
		t.Run("auto-stacking on add", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			item1 := createStackableItem("item-1", "Potion", 0.5, 10)
			item2 := createStackableItem("item-2", "Potion", 0.5, 10)

			_ = mgr.Add(ctx, item1)
			_ = mgr.Add(ctx, item2)

			// Should stack into one
			assert.Equal(t, 1, mgr.Count())
			found, _ := mgr.Get("item-1")
			assert.Equal(t, 2, found.StackSize())
		})

		t.Run("SplitStack", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			item1 := createStackableItem("item-1", "Potion", 0.5, 10)
			item1.AddStack(4) // Now stack of 5
			_ = mgr.Add(ctx, item1)

			newStack, err := mgr.SplitStack(ctx, "item-1", 2)
			require.NoError(t, err)
			assert.NotNil(t, newStack)
			assert.Equal(t, 2, newStack.StackSize())

			original, _ := mgr.Get("item-1")
			assert.Equal(t, 3, original.StackSize())
			assert.Equal(t, 2, mgr.Count()) // Now 2 stacks
		})

		t.Run("MergeStacks", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			item1 := createStackableItem("item-1", "Potion", 0.5, 10)
			item1.AddStack(2) // Stack of 3
			item2 := createStackableItem("item-2", "Potion", 0.5, 10)
			item2.AddStack(1) // Stack of 2

			_ = mgr.AddToSlot(ctx, 0, item1)
			_ = mgr.AddToSlot(ctx, 1, item2)

			err := mgr.MergeStacks(ctx, "item-2", "item-1")
			require.NoError(t, err)

			merged, _ := mgr.Get("item-1")
			assert.Equal(t, 5, merged.StackSize())
			assert.Equal(t, 1, mgr.Count()) // Only 1 stack left
		})

		t.Run("RemoveAmount", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			item1 := createStackableItem("item-1", "Potion", 0.5, 10)
			item1.AddStack(4) // Stack of 5
			_ = mgr.Add(ctx, item1)

			removed, err := mgr.RemoveAmount(ctx, "item-1", 2)
			require.NoError(t, err)
			assert.NotNil(t, removed)

			remaining, _ := mgr.Get("item-1")
			assert.Equal(t, 3, remaining.StackSize())
		})

		t.Run("CanStackWith", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			item1 := createStackableItem("item-1", "Potion", 0.5, 10)
			_ = mgr.Add(ctx, item1)

			item2 := createStackableItem("item-2", "Potion", 0.5, 10)
			targetID, canStack := mgr.CanStackWith(item2)
			assert.True(t, canStack)
			assert.Equal(t, "item-1", targetID)
		})

		t.Run("TotalItems with stacks", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			item1 := createStackableItem("item-1", "Potion", 0.5, 10)
			item1.AddStack(4) // Stack of 5
			_ = mgr.Add(ctx, item1)

			assert.Equal(t, 1, mgr.Count())      // 1 stack
			assert.Equal(t, 5, mgr.TotalItems()) // 5 items total
		})
	})

	t.Run("Search and Filter", func(t *testing.T) {
		t.Run("Search by name", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

			_ = mgr.Add(ctx, createTestItem("sword-1", "Iron Sword", 5.0))
			_ = mgr.Add(ctx, createTestItem("shield-1", "Iron Shield", 10.0))
			_ = mgr.Add(ctx, createTestItem("bow-1", "Wooden Bow", 3.0))

			results := mgr.Search("iron")
			assert.Len(t, results, 2)
		})

		t.Run("FindByType", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

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

			_ = mgr.Add(ctx, material)
			_ = mgr.Add(ctx, consumable)

			materials := mgr.FindByType(item.TypeMaterial)
			assert.Len(t, materials, 1)
			assert.Equal(t, "mat-1", materials[0].ID())

			consumables := mgr.FindByType(item.TypeConsumable)
			assert.Len(t, consumables, 1)
		})

		t.Run("FindByRarity", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

			common := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "common-1",
				Name:     "Common Item",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Rarity:   item.RarityCommon,
			})
			rare := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "rare-1",
				Name:     "Rare Item",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Rarity:   item.RarityRare,
			})

			_ = mgr.Add(ctx, common)
			_ = mgr.Add(ctx, rare)

			commons := mgr.FindByRarity(item.RarityCommon)
			assert.Len(t, commons, 1)

			rares := mgr.FindByRarity(item.RarityRare)
			assert.Len(t, rares, 1)
		})

		t.Run("FindByTag", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

			tagged := item.NewBaseItemWithConfig(item.BaseItemConfig{
				ID:       "tagged-1",
				Name:     "Tagged Item",
				ItemType: item.TypeMaterial,
				Weight:   5.0,
				Tags:     []string{"quest"},
			})
			untagged := createTestItem("untagged-1", "Untagged Item", 5.0)

			_ = mgr.Add(ctx, tagged)
			_ = mgr.Add(ctx, untagged)

			questItems := mgr.FindByTag("quest")
			assert.Len(t, questItems, 1)
			assert.Equal(t, "tagged-1", questItems[0].ID())
		})

		t.Run("Filter with predicate", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

			_ = mgr.Add(ctx, createTestItem("light", "Light", 5.0))
			_ = mgr.Add(ctx, createTestItem("heavy", "Heavy", 50.0))

			lightItems := mgr.Filter(func(i item.Item) bool {
				return i.Weight() < 10
			})
			assert.Len(t, lightItems, 1)
			assert.Equal(t, "light", lightItems[0].ID())
		})

		t.Run("FindStackable", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

			stackable := createStackableItem("stack-1", "Stackable", 1.0, 10)
			nonStackable := createTestItem("single-1", "Single", 5.0)

			_ = mgr.Add(ctx, stackable)
			_ = mgr.Add(ctx, nonStackable)

			stackables := mgr.FindStackable()
			assert.Len(t, stackables, 1)
			assert.Equal(t, "stack-1", stackables[0].ID())
		})
	})

	t.Run("Sorting", func(t *testing.T) {
		t.Run("Sort", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			_ = mgr.Add(ctx, createTestItem("c-item", "Charlie", 10.0))
			_ = mgr.Add(ctx, createTestItem("a-item", "Alpha", 10.0))
			_ = mgr.Add(ctx, createTestItem("b-item", "Bravo", 10.0))

			mgr.Sort(SortByName, true)

			items := mgr.GetAll()
			assert.Equal(t, "Alpha", items[0].Name())
			assert.Equal(t, "Bravo", items[1].Name())
			assert.Equal(t, "Charlie", items[2].Name())
		})

		t.Run("GetSorted", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 100, MaxSlots: 10})

			_ = mgr.Add(ctx, createTestItem("item-1", "Heavy", 50.0))
			_ = mgr.Add(ctx, createTestItem("item-2", "Light", 5.0))
			_ = mgr.Add(ctx, createTestItem("item-3", "Medium", 20.0))

			sorted := mgr.GetSorted(SortByWeight, true)
			assert.Equal(t, 5.0, sorted[0].Weight())
			assert.Equal(t, 20.0, sorted[1].Weight())
			assert.Equal(t, 50.0, sorted[2].Weight())
		})
	})

	t.Run("Stats", func(t *testing.T) {
		t.Run("TotalValue", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

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

			_ = mgr.Add(ctx, item1)
			_ = mgr.Add(ctx, item2)

			assert.Equal(t, int64(300), mgr.TotalValue())
		})

		t.Run("SlotPercent", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxSlots: 10, MaxWeight: 100})

			_ = mgr.Add(ctx, createTestItem("item-1", "Test", 10.0))
			_ = mgr.Add(ctx, createTestItem("item-2", "Test", 10.0))

			assert.Equal(t, 0.2, mgr.SlotPercent())
		})
	})

	t.Run("Callbacks", func(t *testing.T) {
		t.Run("OnItemAdded and OnItemRemoved", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManager()

			var addedItems []string
			var removedItems []string

			mgr.OnItemAdded(func(ctx context.Context, i item.Item) {
				addedItems = append(addedItems, i.ID())
			})
			mgr.OnItemRemoved(func(ctx context.Context, i item.Item) {
				removedItems = append(removedItems, i.ID())
			})

			itm := createTestItem("item-1", "Test Item", 10.0)
			_ = mgr.Add(ctx, itm)

			assert.Equal(t, []string{"item-1"}, addedItems)

			_, _ = mgr.Remove(ctx, "item-1")
			assert.Equal(t, []string{"item-1"}, removedItems)
		})
	})

	t.Run("Persistence", func(t *testing.T) {
		t.Run("GetItemIDs", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxSlots: 10, MaxWeight: 100})

			_ = mgr.AddToSlot(ctx, 0, createTestItem("item-1", "Item 1", 10.0))
			_ = mgr.AddToSlot(ctx, 5, createTestItem("item-2", "Item 2", 15.0))

			ids := mgr.GetItemIDs()
			assert.Len(t, ids, 10) // Full slot array
			assert.Equal(t, "item-1", ids[0])
			assert.Equal(t, "", ids[1])
			assert.Equal(t, "item-2", ids[5])
		})

		t.Run("AddDirect", func(t *testing.T) {
			mgr := NewManager()

			itm := createTestItem("item-1", "Test Item", 10.0)
			err := mgr.AddDirect(itm)

			require.NoError(t, err)
			assert.Equal(t, 1, mgr.Count())
			assert.Equal(t, 10.0, mgr.CurrentWeight())
		})

		t.Run("RecalculateWeight", func(t *testing.T) {
			mgr := NewManager()

			_ = mgr.AddDirect(createTestItem("item-1", "Item 1", 10.0))
			_ = mgr.AddDirect(createTestItem("item-2", "Item 2", 20.0))

			// Manually corrupt weight
			mgr.mu.Lock()
			mgr.currentWeight = 0
			mgr.mu.Unlock()

			mgr.RecalculateWeight()
			assert.Equal(t, 30.0, mgr.CurrentWeight())
		})

		t.Run("Serialization", func(t *testing.T) {
			ctx := context.Background()
			mgr := NewManagerWithConfig(Config{MaxWeight: 150, MaxSlots: 25})

			_ = mgr.Add(ctx, createTestItem("item-1", "Item 1", 10.0))
			_ = mgr.Add(ctx, createTestItem("item-2", "Item 2", 20.0))

			state, err := mgr.SerializeState()
			require.NoError(t, err)
			assert.NotNil(t, state)

			newMgr := NewManager()
			err = newMgr.DeserializeState(state)
			require.NoError(t, err)

			assert.Equal(t, 150.0, newMgr.MaxWeight())
			assert.Equal(t, 25, newMgr.SlotCount())
		})
	})
}
