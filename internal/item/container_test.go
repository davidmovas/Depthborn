package item

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseContainer(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("NewBaseContainer creates with defaults", func(t *testing.T) {
			cont := NewBaseContainer("bag-1", "Backpack", 10)

			require.Equal(t, "Backpack", cont.Name())
			require.Equal(t, TypeContainer, cont.ItemType())
			require.Equal(t, 10, cont.Capacity())
			require.Equal(t, 0, cont.Count())
		})

		t.Run("NewBaseContainerWithConfig respects all fields", func(t *testing.T) {
			cfg := ContainerConfig{
				BaseItemConfig: BaseItemConfig{
					Name:   "Magic Bag",
					Rarity: RarityRare,
					Weight: 1.0,
				},
				Capacity:     20,
				MaxWeight:    50.0,
				AllowedTypes: []Type{TypeMaterial, TypeConsumable},
			}
			cont := NewBaseContainerWithConfig(cfg)

			require.Equal(t, 20, cont.Capacity())
			require.Equal(t, 50.0, cont.MaxWeight())
			require.NotEmpty(t, cont.AllowedTypes())
		})

		t.Run("enforces minimum capacity of 1", func(t *testing.T) {
			cont := NewBaseContainer("", "Tiny Bag", 0)
			require.GreaterOrEqual(t, cont.Capacity(), 1)
		})
	})

	t.Run("AddItem", func(t *testing.T) {
		t.Run("Add places item in container", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			item := NewBaseItem("item-1", TypeMaterial, "Iron Ore")

			err := cont.Add(item)

			require.NoError(t, err)
			require.Equal(t, 1, cont.Count())
			require.True(t, cont.Contains("item-1"))
		})

		t.Run("Add fails when full", func(t *testing.T) {
			cont := NewBaseContainer("", "Small Bag", 1)
			item1 := NewBaseItem("item-1", TypeMaterial, "Iron")
			item2 := NewBaseItem("item-2", TypeMaterial, "Gold")

			cont.Add(item1)
			err := cont.Add(item2)

			require.Error(t, err)
		})

		t.Run("Add fails for nil item", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)

			err := cont.Add(nil)

			require.Error(t, err)
		})

		t.Run("Add fails for duplicate item", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			item := NewBaseItem("item-1", TypeMaterial, "Iron")

			cont.Add(item)
			err := cont.Add(item)

			require.Error(t, err)
		})

		t.Run("Add respects weight limit", func(t *testing.T) {
			cfg := ContainerConfig{
				BaseItemConfig: BaseItemConfig{Name: "Light Bag"},
				Capacity:       10,
				MaxWeight:      5.0,
			}
			cont := NewBaseContainerWithConfig(cfg)
			heavy := NewBaseItemWithConfig(BaseItemConfig{
				Name:     "Heavy Item",
				ItemType: TypeMaterial,
				Weight:   10.0,
			})

			err := cont.Add(heavy)

			require.Error(t, err)
		})

		t.Run("Add respects allowed types", func(t *testing.T) {
			cfg := ContainerConfig{
				BaseItemConfig: BaseItemConfig{Name: "Material Pouch"},
				Capacity:       10,
				AllowedTypes:   []Type{TypeMaterial},
			}
			cont := NewBaseContainerWithConfig(cfg)
			consumable := NewBaseItem("", TypeConsumable, "Potion")

			err := cont.Add(consumable)

			require.Error(t, err)
		})
	})

	t.Run("RemoveItem", func(t *testing.T) {
		t.Run("Remove returns and removes item", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			item := NewBaseItem("item-1", TypeMaterial, "Iron")
			cont.Add(item)

			removed, err := cont.Remove("item-1")

			require.NoError(t, err)
			require.NotNil(t, removed)
			require.Equal(t, "item-1", removed.ID())
			require.Equal(t, 0, cont.Count())
		})

		t.Run("Remove fails for non-existent item", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)

			_, err := cont.Remove("non-existent")

			require.Error(t, err)
		})
	})

	t.Run("QueryOperations", func(t *testing.T) {
		t.Run("Contains checks for item presence", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			item := NewBaseItem("item-1", TypeMaterial, "Iron")
			cont.Add(item)

			require.True(t, cont.Contains("item-1"))
			require.False(t, cont.Contains("non-existent"))
		})

		t.Run("IsFull checks capacity", func(t *testing.T) {
			cont := NewBaseContainer("", "Small Bag", 1)

			require.False(t, cont.IsFull())

			item := NewBaseItem("item-1", TypeMaterial, "Iron")
			cont.Add(item)

			require.True(t, cont.IsFull())
		})

		t.Run("Contents returns copy of items", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			item := NewBaseItem("item-1", TypeMaterial, "Iron")
			cont.Add(item)

			contents := cont.Contents()

			require.Len(t, contents, 1)
			contents[0] = nil
			require.Equal(t, 1, cont.Count())
		})

		t.Run("GetItem returns item without removing", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			item := NewBaseItem("item-1", TypeMaterial, "Iron")
			cont.Add(item)

			got, ok := cont.GetItem("item-1")

			require.True(t, ok)
			require.NotNil(t, got)
			require.Equal(t, 1, cont.Count())
		})

		t.Run("FindByType filters by item type", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			cont.Add(NewBaseItem("mat-1", TypeMaterial, "Iron"))
			cont.Add(NewBaseItem("cons-1", TypeConsumable, "Potion"))
			cont.Add(NewBaseItem("mat-2", TypeMaterial, "Gold"))

			materials := cont.FindByType(TypeMaterial)

			require.Len(t, materials, 2)
		})

		t.Run("FindByTag filters by tag", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			item1 := NewBaseItemWithConfig(BaseItemConfig{
				Name:     "Tagged Item",
				ItemType: TypeMaterial,
				Tags:     []string{"valuable"},
			})
			item2 := NewBaseItem("", TypeMaterial, "Untagged")
			cont.Add(item1)
			cont.Add(item2)

			tagged := cont.FindByTag("valuable")

			require.Len(t, tagged, 1)
		})

		t.Run("RemainingCapacity returns correct value", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 5)
			cont.Add(NewBaseItem("item-1", TypeMaterial, "A"))
			cont.Add(NewBaseItem("item-2", TypeMaterial, "B"))

			require.Equal(t, 3, cont.RemainingCapacity())
		})
	})

	t.Run("Clear", func(t *testing.T) {
		t.Run("Clear removes all items and returns them", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			cont.Add(NewBaseItem("item-1", TypeMaterial, "A"))
			cont.Add(NewBaseItem("item-2", TypeMaterial, "B"))

			cleared := cont.Clear()

			require.Len(t, cleared, 2)
			require.Equal(t, 0, cont.Count())
		})
	})

	t.Run("Properties", func(t *testing.T) {
		t.Run("SetCapacity updates capacity", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 5)

			cont.SetCapacity(20)

			require.Equal(t, 20, cont.Capacity())
		})

		t.Run("SetMaxWeight updates weight limit", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)

			cont.SetMaxWeight(100.0)

			require.Equal(t, 100.0, cont.MaxWeight())
		})

		t.Run("SetAllowedTypes updates allowed types", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)

			cont.SetAllowedTypes([]Type{TypeMaterial, TypeGem})

			require.Len(t, cont.AllowedTypes(), 2)
		})
	})

	t.Run("Weight", func(t *testing.T) {
		t.Run("ContentsWeight calculates sum of item weights", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			cont.Add(NewBaseItemWithConfig(BaseItemConfig{
				Name:     "Heavy",
				ItemType: TypeMaterial,
				Weight:   5.0,
			}))
			cont.Add(NewBaseItemWithConfig(BaseItemConfig{
				Name:     "Light",
				ItemType: TypeMaterial,
				Weight:   2.0,
			}))

			require.Equal(t, 7.0, cont.ContentsWeight())
		})

		t.Run("Weight includes container weight plus contents", func(t *testing.T) {
			cfg := ContainerConfig{
				BaseItemConfig: BaseItemConfig{
					Name:   "Weighted Bag",
					Weight: 1.0,
				},
				Capacity: 10,
			}
			cont := NewBaseContainerWithConfig(cfg)
			cont.Add(NewBaseItemWithConfig(BaseItemConfig{
				Name:     "Item",
				ItemType: TypeMaterial,
				Weight:   5.0,
			}))

			require.Equal(t, 6.0, cont.Weight())
		})
	})

	t.Run("TotalValue", func(t *testing.T) {
		t.Run("TotalValue sums container and contents value", func(t *testing.T) {
			cfg := ContainerConfig{
				BaseItemConfig: BaseItemConfig{
					Name:  "Valuable Bag",
					Value: 50,
				},
				Capacity: 10,
			}
			cont := NewBaseContainerWithConfig(cfg)
			cont.Add(NewBaseItemWithConfig(BaseItemConfig{
				Name:     "Gem",
				ItemType: TypeGem,
				Value:    100,
			}))

			require.Equal(t, int64(150), cont.TotalValue())
		})
	})

	t.Run("Clone", func(t *testing.T) {
		t.Run("creates independent copy with cloned contents", func(t *testing.T) {
			original := NewBaseContainer("", "Magic Bag", 10)
			original.Add(NewBaseItem("item-1", TypeMaterial, "Iron"))

			cloned := original.Clone().(*BaseContainer)

			require.NotEqual(t, original.ID(), cloned.ID())
			require.Equal(t, original.Capacity(), cloned.Capacity())
			require.Equal(t, original.Count(), cloned.Count())

			origContents := original.Contents()
			clonedContents := cloned.Contents()
			require.NotEqual(t, origContents[0].ID(), clonedContents[0].ID())
		})
	})

	t.Run("Serialization", func(t *testing.T) {
		t.Run("Marshal and Unmarshal roundtrip", func(t *testing.T) {
			cfg := ContainerConfig{
				BaseItemConfig: BaseItemConfig{
					ID:       "cont-serial",
					Name:     "Serializable Bag",
					ItemType: TypeContainer,
					Rarity:   RarityRare,
				},
				Capacity:  15,
				MaxWeight: 100.0,
			}
			original := NewBaseContainerWithConfig(cfg)

			data, err := original.Marshal()
			require.NoError(t, err)

			restored := &BaseContainer{}
			err = restored.Unmarshal(data)
			require.NoError(t, err)

			require.Equal(t, original.ID(), restored.ID())
			require.Equal(t, original.Name(), restored.Name())
			require.Equal(t, original.Capacity(), restored.Capacity())
			require.Equal(t, original.MaxWeight(), restored.MaxWeight())
		})

		t.Run("ContentIDs returns item IDs", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			cont.Add(NewBaseItem("item-1", TypeMaterial, "A"))
			cont.Add(NewBaseItem("item-2", TypeMaterial, "B"))

			ids := cont.ContentIDs()

			require.Len(t, ids, 2)
			require.Contains(t, ids, "item-1")
			require.Contains(t, ids, "item-2")
		})

		t.Run("RestoreContents sets contents", func(t *testing.T) {
			cont := NewBaseContainer("", "Bag", 10)
			items := []Item{
				NewBaseItem("item-1", TypeMaterial, "A"),
				NewBaseItem("item-2", TypeMaterial, "B"),
			}

			cont.RestoreContents(items)

			require.Equal(t, 2, cont.Count())
		})
	})

	t.Run("Validation", func(t *testing.T) {
		t.Run("valid container passes", func(t *testing.T) {
			cont := NewBaseContainer("", "Valid Bag", 10)
			require.NoError(t, cont.Validate())
		})

		t.Run("capacity less than 1 fails", func(t *testing.T) {
			cont := NewBaseContainer("", "Bad Bag", 10)
			cont.capacity = 0 // bypass setter
			require.Error(t, cont.Validate())
		})

		t.Run("contents exceeding capacity fails", func(t *testing.T) {
			cont := NewBaseContainer("", "Small Bag", 1)
			cont.Add(NewBaseItem("item-1", TypeMaterial, "A"))
			cont.capacity = 0 // bypass to create invalid state
			require.Error(t, cont.Validate())
		})

		t.Run("negative max weight fails", func(t *testing.T) {
			cont := NewBaseContainer("", "Bad Bag", 10)
			cont.maxWeight = -1 // bypass setter
			require.Error(t, cont.Validate())
		})
	})
}
