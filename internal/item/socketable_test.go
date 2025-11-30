package item

import (
	"context"
	"testing"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/stretchr/testify/require"
)

// mockSocketEffect implements SocketEffect for testing
type mockSocketEffect struct {
	onSocketCalled   bool
	onUnsocketCalled bool
}

func (m *mockSocketEffect) OnSocket(_ context.Context, _ Equipment, _ entity.Entity) error {
	m.onSocketCalled = true
	return nil
}

func (m *mockSocketEffect) OnUnsocket(_ context.Context, _ Equipment, _ entity.Entity) error {
	m.onUnsocketCalled = true
	return nil
}

func (m *mockSocketEffect) EffectType() string {
	return "mock_socket_effect"
}

func (m *mockSocketEffect) Description() string {
	return "Mock socket effect"
}

func TestBaseSocketable(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("NewBaseSocketable creates with defaults", func(t *testing.T) {
			sock := NewBaseSocketable("gem-1", TypeGem, "Ruby", SocketTypeGem)

			require.Equal(t, "Ruby", sock.Name())
			require.Equal(t, TypeGem, sock.ItemType())
			require.Equal(t, SocketTypeGem, sock.SocketType())
			require.Equal(t, 1, sock.Tier())
		})

		t.Run("NewBaseSocketableWithConfig respects all fields", func(t *testing.T) {
			effect := &mockSocketEffect{}
			mod := attribute.NewModifier("str-bonus", attribute.ModFlat, 10, string(attribute.AttrStrength))

			cfg := SocketableConfig{
				BaseItemConfig: BaseItemConfig{
					Name:   "Perfect Ruby",
					Rarity: RarityEpic,
					Level:  20,
				},
				SocketType: SocketTypeGem,
				Effect:     effect,
				EffectID:   "fire_damage",
				Tier:       3,
				Modifiers:  []attribute.Modifier{mod},
			}
			sock := NewBaseSocketableWithConfig(cfg)

			require.Equal(t, 3, sock.Tier())
			require.Equal(t, effect, sock.Effect())
			require.Equal(t, "fire_damage", sock.EffectID())
			require.Len(t, sock.Modifiers(), 1)
		})

		t.Run("clamps tier to valid range 1-5", func(t *testing.T) {
			// Tier below minimum
			cfg := SocketableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Low Tier"},
				SocketType:     SocketTypeGem,
				Tier:           0,
			}
			sock := NewBaseSocketableWithConfig(cfg)
			require.Equal(t, 1, sock.Tier())

			// Tier above maximum
			cfg.Tier = 10
			sock = NewBaseSocketableWithConfig(cfg)
			require.Equal(t, 5, sock.Tier())
		})
	})

	t.Run("Properties", func(t *testing.T) {
		t.Run("SetTier updates tier with clamping", func(t *testing.T) {
			sock := NewBaseSocketable("", TypeGem, "Gem", SocketTypeGem)

			sock.SetTier(3)
			require.Equal(t, 3, sock.Tier())

			sock.SetTier(0)
			require.Equal(t, 1, sock.Tier())

			sock.SetTier(10)
			require.Equal(t, 5, sock.Tier())
		})

		t.Run("SetEffect updates effect", func(t *testing.T) {
			sock := NewBaseSocketable("", TypeGem, "Gem", SocketTypeGem)
			effect := &mockSocketEffect{}

			sock.SetEffect(effect, "new_effect")

			require.Equal(t, effect, sock.Effect())
			require.Equal(t, "new_effect", sock.EffectID())
		})
	})

	t.Run("Modifiers", func(t *testing.T) {
		t.Run("SetModifiers replaces modifiers", func(t *testing.T) {
			sock := NewBaseSocketable("", TypeGem, "Gem", SocketTypeGem)
			mods := []attribute.Modifier{
				attribute.NewModifier("mod1", attribute.ModFlat, 10, "str"),
				attribute.NewModifier("mod2", attribute.ModIncreased, 0.5, "dex"),
			}

			sock.SetModifiers(mods)

			require.Len(t, sock.Modifiers(), 2)
		})

		t.Run("AddModifier appends modifier", func(t *testing.T) {
			sock := NewBaseSocketable("", TypeGem, "Gem", SocketTypeGem)
			mod1 := attribute.NewModifier("mod1", attribute.ModFlat, 10, "str")
			mod2 := attribute.NewModifier("mod2", attribute.ModIncreased, 0.5, "dex")

			sock.AddModifier(mod1)
			sock.AddModifier(mod2)

			require.Len(t, sock.Modifiers(), 2)
		})

		t.Run("Modifiers returns copy", func(t *testing.T) {
			sock := NewBaseSocketable("", TypeGem, "Gem", SocketTypeGem)
			mod := attribute.NewModifier("mod1", attribute.ModFlat, 10, "str")
			sock.AddModifier(mod)

			result := sock.Modifiers()
			result[0] = nil

			require.Len(t, sock.Modifiers(), 1)
			require.NotNil(t, sock.Modifiers()[0])
		})
	})

	t.Run("CanSocketIn", func(t *testing.T) {
		t.Run("matching type returns true", func(t *testing.T) {
			gem := NewBaseSocketable("", TypeGem, "Ruby", SocketTypeGem)
			require.True(t, gem.CanSocketIn(SocketTypeGem))
		})

		t.Run("mismatching type returns false", func(t *testing.T) {
			gem := NewBaseSocketable("", TypeGem, "Ruby", SocketTypeGem)
			require.False(t, gem.CanSocketIn(SocketTypeRune))
		})

		t.Run("universal socket accepts any type", func(t *testing.T) {
			gem := NewBaseSocketable("", TypeGem, "Ruby", SocketTypeGem)
			rune := NewBaseSocketable("", TypeRune, "El", SocketTypeRune)

			require.True(t, gem.CanSocketIn(SocketTypeUniversal))
			require.True(t, rune.CanSocketIn(SocketTypeUniversal))
		})
	})

	t.Run("SocketCallbacks", func(t *testing.T) {
		t.Run("OnSocket calls effect", func(t *testing.T) {
			effect := &mockSocketEffect{}
			cfg := SocketableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Effect Gem"},
				SocketType:     SocketTypeGem,
				Effect:         effect,
			}
			sock := NewBaseSocketableWithConfig(cfg)
			equipment := NewBaseEquipment("", TypeWeaponMelee, "Sword", SlotMainHand)

			err := sock.OnSocket(context.Background(), equipment, nil)

			require.NoError(t, err)
			require.True(t, effect.onSocketCalled)
		})

		t.Run("OnUnsocket calls effect", func(t *testing.T) {
			effect := &mockSocketEffect{}
			cfg := SocketableConfig{
				BaseItemConfig: BaseItemConfig{Name: "Effect Gem"},
				SocketType:     SocketTypeGem,
				Effect:         effect,
			}
			sock := NewBaseSocketableWithConfig(cfg)
			equipment := NewBaseEquipment("", TypeWeaponMelee, "Sword", SlotMainHand)

			err := sock.OnUnsocket(context.Background(), equipment, nil)

			require.NoError(t, err)
			require.True(t, effect.onUnsocketCalled)
		})

		t.Run("OnSocket works without effect", func(t *testing.T) {
			sock := NewBaseSocketable("", TypeGem, "Plain Gem", SocketTypeGem)
			equipment := NewBaseEquipment("", TypeWeaponMelee, "Sword", SlotMainHand)

			err := sock.OnSocket(context.Background(), equipment, nil)

			require.NoError(t, err)
		})
	})

	t.Run("Clone", func(t *testing.T) {
		t.Run("creates independent copy", func(t *testing.T) {
			effect := &mockSocketEffect{}
			mod := attribute.NewModifier("str-bonus", attribute.ModFlat, 10, string(attribute.AttrStrength))

			cfg := SocketableConfig{
				BaseItemConfig: BaseItemConfig{
					Name:   "Original Gem",
					Rarity: RarityRare,
				},
				SocketType: SocketTypeGem,
				Effect:     effect,
				EffectID:   "fire",
				Tier:       3,
				Modifiers:  []attribute.Modifier{mod},
			}
			original := NewBaseSocketableWithConfig(cfg)

			cloned := original.Clone().(*BaseSocketable)

			require.NotEqual(t, original.ID(), cloned.ID())
			require.Equal(t, original.Name(), cloned.Name())
			require.Equal(t, original.SocketType(), cloned.SocketType())
			require.Equal(t, original.Tier(), cloned.Tier())
			require.Len(t, cloned.Modifiers(), len(original.Modifiers()))
			require.Equal(t, original.Effect(), cloned.Effect()) // effect shared
		})
	})

	t.Run("Serialization", func(t *testing.T) {
		t.Run("Marshal and Unmarshal roundtrip", func(t *testing.T) {
			mod := attribute.NewModifierWithPriority("str-bonus", attribute.ModFlat, 15, string(attribute.AttrStrength), 5)

			cfg := SocketableConfig{
				BaseItemConfig: BaseItemConfig{
					ID:       "sock-serial",
					Name:     "Serializable Gem",
					ItemType: TypeGem,
					Rarity:   RarityEpic,
					Level:    10,
				},
				SocketType: SocketTypeGem,
				EffectID:   "test_effect",
				Tier:       4,
				Modifiers:  []attribute.Modifier{mod},
			}
			original := NewBaseSocketableWithConfig(cfg)

			data, err := original.Marshal()
			require.NoError(t, err)

			restored := &BaseSocketable{}
			err = restored.Unmarshal(data)
			require.NoError(t, err)

			require.Equal(t, original.ID(), restored.ID())
			require.Equal(t, original.Name(), restored.Name())
			require.Equal(t, original.SocketType(), restored.SocketType())
			require.Equal(t, original.EffectID(), restored.EffectID())
			require.Equal(t, original.Tier(), restored.Tier())
			require.Len(t, restored.Modifiers(), len(original.Modifiers()))

			// Verify modifier details
			if len(restored.Modifiers()) > 0 {
				origMod := original.Modifiers()[0]
				resMod := restored.Modifiers()[0]

				require.Equal(t, origMod.ID(), resMod.ID())
				require.Equal(t, origMod.Value(), resMod.Value())
				require.Equal(t, origMod.Type(), resMod.Type())
				require.Equal(t, origMod.Priority(), resMod.Priority())
			}
		})
	})

	t.Run("Validation", func(t *testing.T) {
		t.Run("valid socketable passes", func(t *testing.T) {
			sock := NewBaseSocketable("", TypeGem, "Valid Gem", SocketTypeGem)
			require.NoError(t, sock.Validate())
		})

		t.Run("empty socket type fails", func(t *testing.T) {
			sock := NewBaseSocketable("", TypeGem, "Bad Gem", SocketTypeGem)
			sock.socketType = "" // bypass setter
			require.Error(t, sock.Validate())
		})

		t.Run("tier out of range fails", func(t *testing.T) {
			sock := NewBaseSocketable("", TypeGem, "Bad Gem", SocketTypeGem)

			sock.tier = 0 // bypass setter
			require.Error(t, sock.Validate())

			sock.tier = 6 // bypass setter
			require.Error(t, sock.Validate())
		})
	})

	t.Run("SocketTypes", func(t *testing.T) {
		t.Run("all socket types exist", func(t *testing.T) {
			types := []SocketType{
				SocketTypeGem,
				SocketTypeRune,
				SocketTypeUniversal,
			}

			for _, st := range types {
				require.NotEmpty(t, st)
			}
		})
	})
}
