package item

import (
	"context"
	"testing"
	"time"

	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/stretchr/testify/require"
)

// mockConsumableEffect implements ConsumableEffect for testing
type mockConsumableEffect struct {
	applied bool
}

func (m *mockConsumableEffect) Apply(_ context.Context, _ entity.Entity) error {
	m.applied = true
	return nil
}

func (m *mockConsumableEffect) Description() string {
	return "Mock effect description"
}

func (m *mockConsumableEffect) Duration() int64 {
	return 0
}

func TestBaseConsumable(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("NewBaseConsumable creates with defaults", func(t *testing.T) {
			cons := NewBaseConsumable("potion-1", "Health Potion")

			require.Equal(t, "Health Potion", cons.Name())
			require.Equal(t, TypeConsumable, cons.ItemType())
			require.Equal(t, 1, cons.Charges())
			require.Equal(t, 1, cons.MaxCharges())
			require.True(t, cons.IsConsumable())
		})

		t.Run("NewBaseConsumableWithConfig respects all fields", func(t *testing.T) {
			effect := &mockConsumableEffect{}
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{
					Name:         "Super Potion",
					Rarity:       RarityRare,
					MaxStackSize: 10,
				},
				MaxCooldown: 5000,
				Effect:      effect,
				EffectID:    "heal_large",
				Charges:     3,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			require.Equal(t, int64(5000), cons.MaxCooldown())
			require.Equal(t, effect, cons.Effect())
			require.Equal(t, "heal_large", cons.EffectID())
			require.Equal(t, 3, cons.Charges())
			require.Equal(t, 3, cons.MaxCharges())
		})

		t.Run("defaults charges to 1 if invalid", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Test"},
				Charges:        0,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			require.Equal(t, 1, cons.Charges())
		})

		t.Run("allows -1 for infinite charges", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Infinite Potion"},
				Charges:        -1,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			require.Equal(t, -1, cons.Charges())
		})
	})

	t.Run("Usage", func(t *testing.T) {
		t.Run("Use applies effect and decrements charges", func(t *testing.T) {
			effect := &mockConsumableEffect{}
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{
					Name:         "Potion",
					MaxStackSize: 10,
				},
				Effect:  effect,
				Charges: 3,
			}
			cons := NewBaseConsumableWithConfig(cfg)
			cons.AddStack(4) // 5 total items

			err := cons.Use(context.Background(), nil)

			require.NoError(t, err)
			require.True(t, effect.applied)
			require.Equal(t, 2, cons.Charges())
		})

		t.Run("Use consumes item when charges depleted", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{
					Name:         "Single Use",
					MaxStackSize: 5,
				},
				Charges: 1,
			}
			cons := NewBaseConsumableWithConfig(cfg)
			cons.AddStack(2) // 3 total items

			err := cons.Use(context.Background(), nil)

			require.NoError(t, err)
			require.Equal(t, 2, cons.StackSize())
			require.Equal(t, 1, cons.Charges()) // charges reset for next item
		})

		t.Run("Use fails when no charges", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Empty"},
				Charges:        1,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			cons.Use(context.Background(), nil) // deplete

			err := cons.Use(context.Background(), nil) // stack empty
			require.Error(t, err)
		})

		t.Run("Use fails during cooldown", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Cooldown Potion"},
				MaxCooldown:    10000, // 10 second cooldown
				Charges:        5,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			err := cons.Use(context.Background(), nil)
			require.NoError(t, err)

			err = cons.Use(context.Background(), nil) // should fail
			require.Error(t, err)
		})

		t.Run("CanUse returns correct state", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Check Potion"},
				MaxCooldown:    10000,
				Charges:        2,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			require.True(t, cons.CanUse(nil))

			cons.Use(context.Background(), nil)

			require.False(t, cons.CanUse(nil))
		})
	})

	t.Run("Cooldown", func(t *testing.T) {
		t.Run("Cooldown returns remaining time", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Timed Potion"},
				MaxCooldown:    100, // 100ms
				Charges:        5,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			require.Equal(t, int64(0), cons.Cooldown())

			cons.Use(context.Background(), nil)

			cd := cons.Cooldown()
			require.Greater(t, cd, int64(0))
			require.LessOrEqual(t, cd, int64(100))
		})

		t.Run("Cooldown expires over time", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Quick Potion"},
				MaxCooldown:    50,
				Charges:        5,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			cons.Use(context.Background(), nil)
			time.Sleep(60 * time.Millisecond)

			require.Equal(t, int64(0), cons.Cooldown())
		})

		t.Run("SetMaxCooldown updates cooldown", func(t *testing.T) {
			cons := NewBaseConsumable("", "Potion")

			cons.SetMaxCooldown(5000)
			require.Equal(t, int64(5000), cons.MaxCooldown())

			cons.SetMaxCooldown(-100) // negative clamped to 0
			require.Equal(t, int64(0), cons.MaxCooldown())
		})

		t.Run("ResetCooldown clears cooldown", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Reset Potion"},
				MaxCooldown:    10000,
				Charges:        5,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			cons.Use(context.Background(), nil)
			cons.ResetCooldown()

			require.Equal(t, int64(0), cons.Cooldown())
		})
	})

	t.Run("Charges", func(t *testing.T) {
		t.Run("SetCharges updates charges", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Multi-use"},
				Charges:        5,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			cons.SetCharges(3)

			require.Equal(t, 3, cons.Charges())
		})

		t.Run("SetCharges clamps to max", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Multi-use"},
				Charges:        5,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			cons.SetCharges(10)

			require.Equal(t, 5, cons.Charges())
		})

		t.Run("SetCharges allows -1 for infinite", func(t *testing.T) {
			cons := NewBaseConsumable("", "Infinite")

			cons.SetCharges(-1)

			require.Equal(t, -1, cons.Charges())
		})

		t.Run("infinite charges never deplete", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{
					Name:         "Infinite Potion",
					MaxStackSize: 5,
				},
				Charges: -1,
			}
			cons := NewBaseConsumableWithConfig(cfg)

			for i := 0; i < 10; i++ {
				cons.ResetCooldown()
				err := cons.Use(context.Background(), nil)
				require.NoError(t, err)
			}

			require.Equal(t, -1, cons.Charges())
			require.Equal(t, 1, cons.StackSize())
		})
	})

	t.Run("Effect", func(t *testing.T) {
		t.Run("SetEffect updates effect", func(t *testing.T) {
			cons := NewBaseConsumable("", "Potion")
			effect := &mockConsumableEffect{}

			cons.SetEffect(effect, "new_effect")

			require.Equal(t, effect, cons.Effect())
			require.Equal(t, "new_effect", cons.EffectID())
		})
	})

	t.Run("Clone", func(t *testing.T) {
		t.Run("creates independent copy with reset state", func(t *testing.T) {
			effect := &mockConsumableEffect{}
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{
					Name:   "Original Potion",
					Rarity: RarityRare,
				},
				MaxCooldown: 5000,
				Effect:      effect,
				EffectID:    "heal",
				Charges:     3,
			}
			original := NewBaseConsumableWithConfig(cfg)
			original.Use(context.Background(), nil)

			cloned := original.Clone().(*BaseConsumable)

			require.NotEqual(t, original.ID(), cloned.ID())
			require.Equal(t, cloned.MaxCharges(), cloned.Charges())
			require.Equal(t, int64(0), cloned.Cooldown())
			require.Equal(t, original.Effect(), cloned.Effect())
		})
	})

	t.Run("Serialization", func(t *testing.T) {
		t.Run("Marshal and Unmarshal roundtrip", func(t *testing.T) {
			cfg := ConsumableConfig{
				BaseItemConfig: BaseItemConfig{
					ID:       "cons-serial",
					Name:     "Serializable Potion",
					ItemType: TypeConsumable,
					Rarity:   RarityEpic,
					Level:    15,
				},
				MaxCooldown: 3000,
				EffectID:    "test_effect",
				Charges:     5,
			}
			original := NewBaseConsumableWithConfig(cfg)
			original.Use(context.Background(), nil)

			data, err := original.Marshal()
			require.NoError(t, err)

			restored := &BaseConsumable{}
			err = restored.Unmarshal(data)
			require.NoError(t, err)

			require.Equal(t, original.ID(), restored.ID())
			require.Equal(t, original.Name(), restored.Name())
			require.Equal(t, original.MaxCooldown(), restored.MaxCooldown())
			require.Equal(t, original.EffectID(), restored.EffectID())
			require.Equal(t, original.Charges(), restored.Charges())
			require.Equal(t, original.MaxCharges(), restored.MaxCharges())
		})
	})

	t.Run("Validation", func(t *testing.T) {
		t.Run("valid consumable passes", func(t *testing.T) {
			cons := NewBaseConsumable("", "Valid Potion")
			require.NoError(t, cons.Validate())
		})

		t.Run("negative max cooldown fails", func(t *testing.T) {
			cons := NewBaseConsumable("", "Bad Potion")
			cons.maxCooldown = -1 // bypass setter
			require.Error(t, cons.Validate())
		})

		t.Run("charges less than -1 fails", func(t *testing.T) {
			cons := NewBaseConsumable("", "Bad Potion")
			cons.charges = -5 // bypass setter
			require.Error(t, cons.Validate())
		})
	})
}
