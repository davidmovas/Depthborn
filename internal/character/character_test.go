package character

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	assert.NotNil(t, builder)
}

func TestBuilderWithName(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("TestHero").
		Build(ctx)

	require.NoError(t, err)
	assert.Equal(t, "TestHero", char.Name())
}

func TestBuilderWithoutName(t *testing.T) {
	ctx := context.Background()

	_, err := NewBuilder().Build(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestBuilderWithHealth(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		WithHealth(50, 100).
		Build(ctx)

	require.NoError(t, err)
	assert.Equal(t, 50.0, char.Health())
	assert.Equal(t, 100.0, char.MaxHealth())
}

func TestBuilderWithAttributes(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		WithStrength(15).
		WithDexterity(12).
		WithIntelligence(10).
		Build(ctx)

	require.NoError(t, err)
	assert.NotNil(t, char.Attributes())
}

func TestBuilderWithGold(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		WithGold(500).
		Build(ctx)

	require.NoError(t, err)
	assert.Equal(t, int64(500), char.Gold())
}

func TestBuilderWithMaxWeight(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		WithMaxWeight(200).
		Build(ctx)

	require.NoError(t, err)
	inv := char.InventoryManager()
	assert.Equal(t, 200.0, inv.MaxWeight())
}

func TestBuilderWithPosition(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		WithPosition(10, 20, 0).
		Build(ctx)

	require.NoError(t, err)

	pos := char.Transform().Position()
	assert.Equal(t, 10, pos.X)
	assert.Equal(t, 20, pos.Y)
	assert.Equal(t, 0, pos.Z)
}

func TestBuilderWithTags(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		WithTags("player", "warrior", "human").
		Build(ctx)

	require.NoError(t, err)
	assert.True(t, char.Tags().Has("player"))
	assert.True(t, char.Tags().Has("warrior"))
	assert.True(t, char.Tags().Has("human"))
}

func TestBuilderPresets(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		builder *CharacterBuilder
		tag     string
	}{
		{"Warrior", Warrior("WarriorHero"), "warrior"},
		{"Ranger", Ranger("RangerHero"), "ranger"},
		{"Mage", Mage("MageHero"), "mage"},
		{"Balanced", Balanced("BalancedHero"), "balanced"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char, err := tt.builder.Build(ctx)
			require.NoError(t, err)
			assert.True(t, char.Tags().Has(tt.tag))
		})
	}
}

func TestCharacterGoldOperations(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		WithGold(100).
		Build(ctx)

	require.NoError(t, err)

	// Add gold
	char.AddGold(50)
	assert.Equal(t, int64(150), char.Gold())

	// Remove gold - success
	ok := char.RemoveGold(30)
	assert.True(t, ok)
	assert.Equal(t, int64(120), char.Gold())

	// Remove gold - insufficient
	ok = char.RemoveGold(200)
	assert.False(t, ok)
	assert.Equal(t, int64(120), char.Gold())
}

func TestCharacterPlayTime(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		Build(ctx)

	require.NoError(t, err)

	assert.Equal(t, int64(0), char.PlayTime())

	char.AddPlayTime(60)
	assert.Equal(t, int64(60), char.PlayTime())

	char.AddPlayTime(120)
	assert.Equal(t, int64(180), char.PlayTime())
}

func TestCharacterDeathCount(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		Build(ctx)

	require.NoError(t, err)

	assert.Equal(t, 0, char.DeathCount())

	char.IncrementDeathCount()
	assert.Equal(t, 1, char.DeathCount())

	char.IncrementDeathCount()
	assert.Equal(t, 2, char.DeathCount())
}

func TestCharacterFlags(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		Build(ctx)

	require.NoError(t, err)

	flags := char.Flags()
	assert.NotNil(t, flags)

	assert.False(t, flags.Has("tutorial_complete"))

	flags.Set("tutorial_complete")
	assert.True(t, flags.Has("tutorial_complete"))

	flags.Toggle("tutorial_complete")
	assert.False(t, flags.Has("tutorial_complete"))
}

func TestCharacterIsAliveAndDead(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		WithHealth(100, 100).
		Build(ctx)

	require.NoError(t, err)

	assert.True(t, char.IsAlive())
	assert.False(t, char.IsDead())

	// Kill character
	err = char.Die(ctx, "enemy-001")
	require.NoError(t, err)

	assert.False(t, char.IsAlive())
	assert.True(t, char.IsDead())
	assert.Equal(t, 1, char.DeathCount())

	// Respawn
	err = char.Respawn(ctx, 0.5)
	require.NoError(t, err)

	assert.True(t, char.IsAlive())
	assert.False(t, char.IsDead())
}

func TestCharacterExperience(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		Build(ctx)

	require.NoError(t, err)

	assert.Equal(t, 1, char.Level())

	// Add experience
	err = char.AddExperience(ctx, 100)
	require.NoError(t, err)

	assert.Equal(t, int64(100), char.Experience())
}

func TestCharacterClone(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		WithGold(500).
		WithStrength(15).
		Build(ctx)

	require.NoError(t, err)

	cloned := char.Clone().(*BaseCharacter)

	assert.NotEqual(t, char.ID(), cloned.ID())
	assert.Equal(t, char.Name(), cloned.Name())
	assert.Equal(t, char.Gold(), cloned.Gold())
}

func TestCharacterValidation(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("Hero").
		Build(ctx)

	require.NoError(t, err)

	err = char.Validate()
	assert.NoError(t, err)
}

func TestCharacterSerialization(t *testing.T) {
	ctx := context.Background()

	char, err := NewBuilder().
		WithName("TestHero").
		WithGold(1000).
		WithHealth(80, 100).
		WithAccountID("account-123").
		Build(ctx)

	require.NoError(t, err)

	// Add some play time and deaths
	char.AddPlayTime(3600)
	char.IncrementDeathCount()
	char.Flags().Set("tutorial_done")

	// Serialize
	data, err := char.MarshalBinary()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Deserialize
	restoredChar := NewCharacter(DefaultConfig(""))
	err = restoredChar.UnmarshalBinary(data)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, char.ID(), restoredChar.ID())
	assert.Equal(t, char.Name(), restoredChar.Name())
	assert.Equal(t, char.Gold(), restoredChar.Gold())
	assert.Equal(t, char.PlayTime(), restoredChar.PlayTime())
	assert.Equal(t, char.AccountID(), restoredChar.AccountID())
	assert.True(t, restoredChar.Flags().Has("tutorial_done"))
}

func TestFlagSet(t *testing.T) {
	flags := NewFlagSet()

	assert.Equal(t, 0, flags.Count())

	flags.Set("flag1")
	flags.Set("flag2")
	assert.Equal(t, 2, flags.Count())
	assert.True(t, flags.Has("flag1"))
	assert.True(t, flags.Has("flag2"))

	flags.Unset("flag1")
	assert.Equal(t, 1, flags.Count())
	assert.False(t, flags.Has("flag1"))

	flags.Clear()
	assert.Equal(t, 0, flags.Count())
}

func TestFlagSetClone(t *testing.T) {
	flags := NewFlagSet()
	flags.Set("flag1")
	flags.Set("flag2")

	cloned := flags.Clone()

	assert.True(t, cloned.Has("flag1"))
	assert.True(t, cloned.Has("flag2"))

	// Modify original
	flags.Set("flag3")
	assert.False(t, cloned.Has("flag3"))
}
