package affix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
)

// Test helpers

func createTestAffix(id string, affixType Type, rank int) *BaseAffix {
	return NewBaseAffix(id, "Test "+id, affixType).
		WithRank(rank).
		WithBaseWeight(100).
		AddModifier(ModifierTemplate{
			Attribute: attribute.AttrPhysicalDamage,
			ModType:   attribute.ModFlat,
			MinValue:  10,
			MaxValue:  20,
			Priority:  0,
		})
}

func createTestAffixWithGroup(id string, affixType Type, group string) *BaseAffix {
	return NewBaseAffix(id, "Test "+id, affixType).
		WithGroup(group).
		WithRank(50).
		WithBaseWeight(100).
		AddModifier(ModifierTemplate{
			Attribute: attribute.AttrPhysicalDamage,
			ModType:   attribute.ModFlat,
			MinValue:  10,
			MaxValue:  20,
		})
}

func TestBaseAffix(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("NewBaseAffix creates with defaults", func(t *testing.T) {
			affix := NewBaseAffix("test-1", "Test Affix", TypePrefix)

			assert.Equal(t, "test-1", affix.ID())
			assert.Equal(t, "Test Affix", affix.Name())
			assert.Equal(t, TypePrefix, affix.Type())
			assert.Equal(t, "", affix.Group())
			assert.Equal(t, 50, affix.Rank()) // default
			assert.Equal(t, 100, affix.BaseWeight())
			assert.Empty(t, affix.Modifiers())
			assert.Empty(t, affix.Tags())
		})

		t.Run("NewBaseAffixWithConfig uses all fields", func(t *testing.T) {
			cfg := AffixConfig{
				ID:          "cfg-affix",
				Name:        "Configured Affix",
				Type:        TypeSuffix,
				Group:       "damage",
				Rank:        75,
				BaseWeight:  200,
				Description: "Test description",
				Tags:        []string{"fire", "attack"},
				Modifiers: []ModifierTemplate{
					{Attribute: attribute.AttrFireResist, ModType: attribute.ModFlat, MinValue: 5, MaxValue: 15},
				},
			}

			affix := NewBaseAffixWithConfig(cfg)

			assert.Equal(t, "cfg-affix", affix.ID())
			assert.Equal(t, "Configured Affix", affix.Name())
			assert.Equal(t, TypeSuffix, affix.Type())
			assert.Equal(t, "damage", affix.Group())
			assert.Equal(t, 75, affix.Rank())
			assert.Equal(t, 200, affix.BaseWeight())
			assert.Equal(t, "Test description", affix.Description())
			assert.Equal(t, []string{"fire", "attack"}, affix.Tags())
			assert.Len(t, affix.Modifiers(), 1)
		})
	})

	t.Run("Fluent API", func(t *testing.T) {
		t.Run("builder methods chain correctly", func(t *testing.T) {
			affix := NewBaseAffix("chain-test", "Chain", TypePrefix).
				WithGroup("test-group").
				WithRank(80).
				WithBaseWeight(150).
				WithDescription("Description").
				WithTags([]string{"tag1", "tag2"}).
				AddTag("tag3").
				AddModifier(ModifierTemplate{Attribute: attribute.AttrStrength, ModType: attribute.ModFlat, MinValue: 1, MaxValue: 5})

			assert.Equal(t, "test-group", affix.Group())
			assert.Equal(t, 80, affix.Rank())
			assert.Equal(t, 150, affix.BaseWeight())
			assert.Equal(t, "Description", affix.Description())
			assert.Equal(t, []string{"tag1", "tag2", "tag3"}, affix.Tags())
			assert.Len(t, affix.Modifiers(), 1)
		})

		t.Run("rank is clamped to 1-100", func(t *testing.T) {
			low := NewBaseAffix("low", "Low", TypePrefix).WithRank(0)
			high := NewBaseAffix("high", "High", TypePrefix).WithRank(150)

			assert.Equal(t, 1, low.Rank())
			assert.Equal(t, 100, high.Rank())
		})
	})

	t.Run("Tags", func(t *testing.T) {
		t.Run("HasTag returns true for existing tag", func(t *testing.T) {
			affix := NewBaseAffix("tag-test", "Tag Test", TypePrefix).
				WithTags([]string{"fire", "attack"})

			assert.True(t, affix.HasTag("fire"))
			assert.True(t, affix.HasTag("attack"))
			assert.False(t, affix.HasTag("cold"))
		})
	})
}

func TestBaseInstance(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("NewBaseInstance creates with rolled values", func(t *testing.T) {
			affix := createTestAffix("inst-test", TypePrefix, 50)
			values := []RolledModifier{
				{Template: affix.Modifiers()[0], Value: 15.0},
			}

			instance := NewBaseInstance(affix, values)

			assert.Equal(t, "inst-test", instance.AffixID())
			assert.Equal(t, affix, instance.Affix())
			assert.Equal(t, TypePrefix, instance.Type())
			assert.Len(t, instance.RolledValues(), 1)
			assert.Equal(t, 15.0, instance.RolledValues()[0].Value)
		})

		t.Run("NewBaseInstanceFromData works without affix reference", func(t *testing.T) {
			values := []RolledModifier{
				{Template: ModifierTemplate{Attribute: attribute.AttrStrength, ModType: attribute.ModFlat, MinValue: 1, MaxValue: 10}, Value: 5.0},
			}

			instance := NewBaseInstanceFromData("affix-id", TypeSuffix, "test-group", values)

			assert.Equal(t, "affix-id", instance.AffixID())
			assert.Nil(t, instance.Affix())
			assert.Equal(t, TypeSuffix, instance.Type())
			assert.Equal(t, "test-group", instance.Group())
		})
	})

	t.Run("Modifiers", func(t *testing.T) {
		t.Run("returns attribute modifiers with rolled values", func(t *testing.T) {
			affix := createTestAffix("mod-test", TypePrefix, 50)
			values := RollModifiers(affix.Modifiers())
			instance := NewBaseInstance(affix, values)

			mods := instance.Modifiers()

			require.Len(t, mods, 1)
			assert.Equal(t, attribute.ModFlat, mods[0].Type())
			assert.Equal(t, "mod-test", mods[0].Source())
		})
	})

	t.Run("Reroll", func(t *testing.T) {
		t.Run("changes values within range", func(t *testing.T) {
			affix := NewBaseAffix("reroll-test", "Reroll", TypePrefix).
				AddModifier(ModifierTemplate{
					Attribute: attribute.AttrPhysicalDamage,
					ModType:   attribute.ModFlat,
					MinValue:  10,
					MaxValue:  20,
				})

			values := []RolledModifier{
				{Template: affix.Modifiers()[0], Value: 15.0},
			}
			instance := NewBaseInstance(affix, values)

			// Reroll multiple times and check values stay in range
			for i := 0; i < 10; i++ {
				instance.Reroll()
				val := instance.RolledValues()[0].Value
				assert.GreaterOrEqual(t, val, 10.0)
				assert.LessOrEqual(t, val, 20.0)
			}
		})

		t.Run("RerollSingle only changes specified index", func(t *testing.T) {
			affix := NewBaseAffix("single-reroll", "Single", TypePrefix).
				AddModifier(ModifierTemplate{Attribute: attribute.AttrStrength, ModType: attribute.ModFlat, MinValue: 1, MaxValue: 10}).
				AddModifier(ModifierTemplate{Attribute: attribute.AttrDexterity, ModType: attribute.ModFlat, MinValue: 100, MaxValue: 200})

			values := []RolledModifier{
				{Template: affix.Modifiers()[0], Value: 5.0},
				{Template: affix.Modifiers()[1], Value: 150.0},
			}
			instance := NewBaseInstance(affix, values)

			err := instance.RerollSingle(0)
			require.NoError(t, err)

			// Second value should be unchanged
			assert.Equal(t, 150.0, instance.RolledValues()[1].Value)
		})

		t.Run("RerollSingle returns error for invalid index", func(t *testing.T) {
			affix := createTestAffix("idx-test", TypePrefix, 50)
			instance := NewBaseInstance(affix, RollModifiers(affix.Modifiers()))

			err := instance.RerollSingle(5)
			assert.Error(t, err)
		})
	})

	t.Run("Quality", func(t *testing.T) {
		t.Run("returns 0 for minimum roll", func(t *testing.T) {
			affix := NewBaseAffix("min-quality", "Min", TypePrefix).
				AddModifier(ModifierTemplate{MinValue: 10, MaxValue: 20})

			values := []RolledModifier{{Template: affix.Modifiers()[0], Value: 10.0}}
			instance := NewBaseInstance(affix, values)

			assert.Equal(t, 0.0, instance.Quality())
		})

		t.Run("returns 1 for maximum roll", func(t *testing.T) {
			affix := NewBaseAffix("max-quality", "Max", TypePrefix).
				AddModifier(ModifierTemplate{MinValue: 10, MaxValue: 20})

			values := []RolledModifier{{Template: affix.Modifiers()[0], Value: 20.0}}
			instance := NewBaseInstance(affix, values)

			assert.Equal(t, 1.0, instance.Quality())
		})

		t.Run("returns 0.5 for middle roll", func(t *testing.T) {
			affix := NewBaseAffix("mid-quality", "Mid", TypePrefix).
				AddModifier(ModifierTemplate{MinValue: 10, MaxValue: 20})

			values := []RolledModifier{{Template: affix.Modifiers()[0], Value: 15.0}}
			instance := NewBaseInstance(affix, values)

			assert.Equal(t, 0.5, instance.Quality())
		})
	})
}

func TestBaseSet(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		t.Run("NewBaseSet creates with default limits", func(t *testing.T) {
			set := NewBaseSet()

			assert.Equal(t, 3, set.MaxPrefixes())
			assert.Equal(t, 3, set.MaxSuffixes())
			assert.Equal(t, 0, set.Count())
		})

		t.Run("NewBaseSetForRarity uses correct limits", func(t *testing.T) {
			common := NewBaseSetForRarity(0)
			rare := NewBaseSetForRarity(2)
			legendary := NewBaseSetForRarity(4)

			assert.Equal(t, 0, common.MaxPrefixes())
			assert.Equal(t, 2, rare.MaxPrefixes())
			assert.Equal(t, 3, legendary.MaxPrefixes())
		})
	})

	t.Run("Add and Remove", func(t *testing.T) {
		t.Run("Add succeeds within limits", func(t *testing.T) {
			set := NewBaseSet()
			affix := createTestAffix("add-test", TypePrefix, 50)
			instance := NewBaseInstance(affix, RollModifiers(affix.Modifiers()))

			err := set.Add(instance)

			require.NoError(t, err)
			assert.Equal(t, 1, set.Count())
			assert.Equal(t, 1, set.PrefixCount())
		})

		t.Run("Add fails when at max", func(t *testing.T) {
			set := NewBaseSetWithLimits(AffixLimits{0, 1, 0, 1})

			affix1 := createTestAffix("prefix-1", TypePrefix, 50)
			affix2 := createTestAffix("prefix-2", TypePrefix, 50)

			_ = set.Add(NewBaseInstance(affix1, RollModifiers(affix1.Modifiers())))
			err := set.Add(NewBaseInstance(affix2, RollModifiers(affix2.Modifiers())))

			assert.Error(t, err)
		})

		t.Run("Remove succeeds for existing affix", func(t *testing.T) {
			set := NewBaseSet()
			affix := createTestAffix("remove-test", TypePrefix, 50)
			instance := NewBaseInstance(affix, RollModifiers(affix.Modifiers()))

			_ = set.Add(instance)
			err := set.Remove(affix.ID())

			require.NoError(t, err)
			assert.Equal(t, 0, set.Count())
		})

		t.Run("Remove fails for non-existing affix", func(t *testing.T) {
			set := NewBaseSet()

			err := set.Remove("non-existent")

			assert.Error(t, err)
		})
	})

	t.Run("Group Mutual Exclusion", func(t *testing.T) {
		t.Run("cannot add two affixes from same group", func(t *testing.T) {
			set := NewBaseSet()

			affix1 := createTestAffixWithGroup("group-1", TypePrefix, "damage")
			affix2 := createTestAffixWithGroup("group-2", TypePrefix, "damage")

			_ = set.Add(NewBaseInstance(affix1, RollModifiers(affix1.Modifiers())))
			err := set.Add(NewBaseInstance(affix2, RollModifiers(affix2.Modifiers())))

			assert.Error(t, err)
			assert.True(t, set.HasGroup("damage"))
		})

		t.Run("can add affixes from different groups", func(t *testing.T) {
			set := NewBaseSet()

			affix1 := createTestAffixWithGroup("group-a", TypePrefix, "damage")
			affix2 := createTestAffixWithGroup("group-b", TypePrefix, "defense")

			err1 := set.Add(NewBaseInstance(affix1, RollModifiers(affix1.Modifiers())))
			err2 := set.Add(NewBaseInstance(affix2, RollModifiers(affix2.Modifiers())))

			require.NoError(t, err1)
			require.NoError(t, err2)
			assert.Equal(t, 2, set.Count())
		})

		t.Run("removing affix frees up group", func(t *testing.T) {
			set := NewBaseSet()

			affix1 := createTestAffixWithGroup("free-group", TypePrefix, "damage")
			instance := NewBaseInstance(affix1, RollModifiers(affix1.Modifiers()))

			_ = set.Add(instance)
			assert.True(t, set.HasGroup("damage"))

			_ = set.Remove(affix1.ID())
			assert.False(t, set.HasGroup("damage"))
		})
	})

	t.Run("Get Operations", func(t *testing.T) {
		t.Run("Get returns instance by ID", func(t *testing.T) {
			set := NewBaseSet()
			affix := createTestAffix("get-test", TypeSuffix, 50)
			instance := NewBaseInstance(affix, RollModifiers(affix.Modifiers()))

			_ = set.Add(instance)
			got, exists := set.Get("get-test")

			assert.True(t, exists)
			assert.Equal(t, "get-test", got.AffixID())
		})

		t.Run("GetByType returns only matching type", func(t *testing.T) {
			set := NewBaseSet()

			prefix := createTestAffix("p1", TypePrefix, 50)
			suffix := createTestAffix("s1", TypeSuffix, 50)

			_ = set.Add(NewBaseInstance(prefix, RollModifiers(prefix.Modifiers())))
			_ = set.Add(NewBaseInstance(suffix, RollModifiers(suffix.Modifiers())))

			prefixes := set.GetByType(TypePrefix)
			suffixes := set.GetByType(TypeSuffix)

			assert.Len(t, prefixes, 1)
			assert.Len(t, suffixes, 1)
		})

		t.Run("GetAll returns all instances", func(t *testing.T) {
			set := NewBaseSet()

			for i := 0; i < 3; i++ {
				affix := createTestAffix("all-"+string(rune('a'+i)), TypePrefix, 50)
				_ = set.Add(NewBaseInstance(affix, RollModifiers(affix.Modifiers())))
			}

			all := set.GetAll()
			assert.Len(t, all, 3)
		})
	})

	t.Run("Modifiers", func(t *testing.T) {
		t.Run("AllModifiers aggregates from all instances", func(t *testing.T) {
			set := NewBaseSet()

			affix1 := NewBaseAffix("mod1", "Mod1", TypePrefix).
				AddModifier(ModifierTemplate{Attribute: attribute.AttrStrength, ModType: attribute.ModFlat, MinValue: 1, MaxValue: 5})
			affix2 := NewBaseAffix("mod2", "Mod2", TypeSuffix).
				AddModifier(ModifierTemplate{Attribute: attribute.AttrDexterity, ModType: attribute.ModFlat, MinValue: 1, MaxValue: 5})

			_ = set.Add(NewBaseInstance(affix1, RollModifiers(affix1.Modifiers())))
			_ = set.Add(NewBaseInstance(affix2, RollModifiers(affix2.Modifiers())))

			mods := set.AllModifiers()
			assert.Len(t, mods, 2)
		})
	})

	t.Run("Quality and Reroll", func(t *testing.T) {
		t.Run("TotalQuality averages instance qualities", func(t *testing.T) {
			set := NewBaseSet()

			affix := NewBaseAffix("qual", "Qual", TypePrefix).
				AddModifier(ModifierTemplate{MinValue: 0, MaxValue: 100})

			// Add instance with quality 0.5
			values := []RolledModifier{{Template: affix.Modifiers()[0], Value: 50}}
			_ = set.Add(NewBaseInstance(affix, values))

			assert.Equal(t, 0.5, set.TotalQuality())
		})

		t.Run("RerollAll changes all values", func(t *testing.T) {
			set := NewBaseSet()

			affix := NewBaseAffix("reroll-all", "Reroll", TypePrefix).
				AddModifier(ModifierTemplate{MinValue: 1, MaxValue: 100})

			values := []RolledModifier{{Template: affix.Modifiers()[0], Value: 50}}
			_ = set.Add(NewBaseInstance(affix, values))

			original := set.GetAll()[0].RolledValues()[0].Value

			// Reroll multiple times to ensure at least one change
			changed := false
			for i := 0; i < 10; i++ {
				set.RerollAll()
				if set.GetAll()[0].RolledValues()[0].Value != original {
					changed = true
					break
				}
			}

			assert.True(t, changed, "RerollAll should change values")
		})
	})

	t.Run("Completeness Checks", func(t *testing.T) {
		t.Run("IsComplete checks minimums", func(t *testing.T) {
			set := NewBaseSetWithLimits(AffixLimits{1, 3, 1, 3})

			assert.False(t, set.IsComplete())

			prefix := createTestAffix("comp-p", TypePrefix, 50)
			suffix := createTestAffix("comp-s", TypeSuffix, 50)

			_ = set.Add(NewBaseInstance(prefix, RollModifiers(prefix.Modifiers())))
			assert.False(t, set.IsComplete())

			_ = set.Add(NewBaseInstance(suffix, RollModifiers(suffix.Modifiers())))
			assert.True(t, set.IsComplete())
		})

		t.Run("NeedMore returns true when below minimum", func(t *testing.T) {
			set := NewBaseSetWithLimits(AffixLimits{2, 3, 1, 3})

			assert.True(t, set.NeedMore(TypePrefix))
			assert.True(t, set.NeedMore(TypeSuffix))

			prefix := createTestAffix("need-p", TypePrefix, 50)
			_ = set.Add(NewBaseInstance(prefix, RollModifiers(prefix.Modifiers())))

			assert.True(t, set.NeedMore(TypePrefix)) // Still need 1 more

			prefix2 := createTestAffix("need-p2", TypePrefix, 50)
			_ = set.Add(NewBaseInstance(prefix2, RollModifiers(prefix2.Modifiers())))

			assert.False(t, set.NeedMore(TypePrefix)) // Met minimum
		})
	})
}

func TestBasePool(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		t.Run("Add and Get work correctly", func(t *testing.T) {
			pool := NewBasePool()
			affix := createTestAffix("pool-test", TypePrefix, 50)

			pool.Add(affix)
			got, exists := pool.Get("pool-test")

			assert.True(t, exists)
			assert.Equal(t, "pool-test", got.ID())
		})

		t.Run("Remove works correctly", func(t *testing.T) {
			pool := NewBasePool()
			affix := createTestAffix("remove-pool", TypePrefix, 50)

			pool.Add(affix)
			pool.Remove("remove-pool")
			_, exists := pool.Get("remove-pool")

			assert.False(t, exists)
		})

		t.Run("GetAll returns all affixes", func(t *testing.T) {
			pool := NewBasePool()

			for i := 0; i < 5; i++ {
				pool.Add(createTestAffix("all-"+string(rune('a'+i)), TypePrefix, 50))
			}

			all := pool.GetAll()
			assert.Len(t, all, 5)
		})
	})

	t.Run("Filtering", func(t *testing.T) {
		t.Run("GetByGroup returns matching affixes", func(t *testing.T) {
			pool := NewBasePool()

			pool.Add(createTestAffixWithGroup("g1", TypePrefix, "damage"))
			pool.Add(createTestAffixWithGroup("g2", TypePrefix, "damage"))
			pool.Add(createTestAffixWithGroup("g3", TypePrefix, "defense"))

			damage := pool.GetByGroup("damage")
			defense := pool.GetByGroup("defense")

			assert.Len(t, damage, 2)
			assert.Len(t, defense, 1)
		})

		t.Run("GetByTags requires all tags", func(t *testing.T) {
			pool := NewBasePool()

			pool.Add(NewBaseAffix("t1", "T1", TypePrefix).WithTags([]string{"fire", "attack"}))
			pool.Add(NewBaseAffix("t2", "T2", TypePrefix).WithTags([]string{"fire"}))
			pool.Add(NewBaseAffix("t3", "T3", TypePrefix).WithTags([]string{"cold", "attack"}))

			fireAttack := pool.GetByTags("fire", "attack")
			assert.Len(t, fireAttack, 1)
			assert.Equal(t, "t1", fireAttack[0].ID())
		})

		t.Run("GetByAnyTag matches any tag", func(t *testing.T) {
			pool := NewBasePool()

			pool.Add(NewBaseAffix("a1", "A1", TypePrefix).WithTags([]string{"fire"}))
			pool.Add(NewBaseAffix("a2", "A2", TypePrefix).WithTags([]string{"cold"}))
			pool.Add(NewBaseAffix("a3", "A3", TypePrefix).WithTags([]string{"lightning"}))

			fireCold := pool.GetByAnyTag("fire", "cold")
			assert.Len(t, fireCold, 2)
		})

		t.Run("Filter with criteria", func(t *testing.T) {
			pool := NewBasePool()

			pool.Add(NewBaseAffix("f1", "F1", TypePrefix).WithRank(20))
			pool.Add(NewBaseAffix("f2", "F2", TypePrefix).WithRank(50))
			pool.Add(NewBaseAffix("f3", "F3", TypeSuffix).WithRank(80))

			criteria := FilterCriteria{
				Types:   []Type{TypePrefix},
				MinRank: 30,
			}

			filtered := pool.Filter(criteria)
			assert.Len(t, filtered, 1)
			assert.Equal(t, "f2", filtered[0].ID())
		})
	})

	t.Run("Roll", func(t *testing.T) {
		t.Run("returns error when no eligible affixes", func(t *testing.T) {
			pool := NewBasePool()

			ctx := RollContext{ItemType: "sword", ItemLevel: 10}
			_, err := pool.Roll(ctx)

			assert.Error(t, err)
		})

		t.Run("respects type filter", func(t *testing.T) {
			pool := NewBasePool()

			pool.Add(createTestAffix("prefix-only", TypePrefix, 50))
			pool.Add(createTestAffix("suffix-only", TypeSuffix, 50))

			suffixType := TypeSuffix
			ctx := RollContext{AffixType: &suffixType}

			// Roll multiple times to verify only suffixes are returned
			for i := 0; i < 10; i++ {
				affix, err := pool.Roll(ctx)
				require.NoError(t, err)
				assert.Equal(t, TypeSuffix, affix.Type())
			}
		})

		t.Run("respects exclude groups", func(t *testing.T) {
			pool := NewBasePool()

			pool.Add(createTestAffixWithGroup("exc-1", TypePrefix, "damage"))
			pool.Add(createTestAffixWithGroup("exc-2", TypePrefix, "defense"))

			ctx := RollContext{ExcludeGroups: []string{"damage"}}

			for i := 0; i < 10; i++ {
				affix, err := pool.Roll(ctx)
				require.NoError(t, err)
				assert.Equal(t, "defense", affix.Group())
			}
		})

		t.Run("respects requirements", func(t *testing.T) {
			pool := NewBasePool()

			lowLevel := createTestAffix("low-level", TypePrefix, 50)
			lowLevel.WithRequirements(NewBaseRequirements(1))

			highLevel := createTestAffix("high-level", TypePrefix, 50)
			highLevel.WithRequirements(NewBaseRequirements(50))

			pool.Add(lowLevel)
			pool.Add(highLevel)

			ctx := RollContext{ItemLevel: 10}

			for i := 0; i < 10; i++ {
				affix, err := pool.Roll(ctx)
				require.NoError(t, err)
				assert.Equal(t, "low-level", affix.ID())
			}
		})
	})
}

func TestBaseGenerator(t *testing.T) {
	t.Run("Generate", func(t *testing.T) {
		t.Run("creates instances within limits", func(t *testing.T) {
			pool := NewBasePool()
			for i := 0; i < 10; i++ {
				pool.Add(createTestAffix("p-"+string(rune('a'+i)), TypePrefix, 50))
				pool.Add(createTestAffix("s-"+string(rune('a'+i)), TypeSuffix, 50))
			}

			gen := NewBaseGenerator(pool)
			ctx := GenerateContext{
				RollContext: RollContext{ItemType: "sword", ItemLevel: 50, ItemRarity: 3},
				PrefixRange: [2]int{1, 2},
				SuffixRange: [2]int{1, 2},
				QualityBias: 0.5,
			}

			instances, err := gen.Generate(ctx)
			require.NoError(t, err)

			prefixCount := 0
			suffixCount := 0
			for _, inst := range instances {
				if inst.Type() == TypePrefix {
					prefixCount++
				} else if inst.Type() == TypeSuffix {
					suffixCount++
				}
			}

			assert.GreaterOrEqual(t, prefixCount, 1)
			assert.LessOrEqual(t, prefixCount, 2)
			assert.GreaterOrEqual(t, suffixCount, 1)
			assert.LessOrEqual(t, suffixCount, 2)
		})

		t.Run("respects group exclusion", func(t *testing.T) {
			pool := NewBasePool()

			// Add multiple affixes in same group
			for i := 0; i < 5; i++ {
				pool.Add(createTestAffixWithGroup("same-"+string(rune('a'+i)), TypePrefix, "same-group"))
			}

			gen := NewBaseGenerator(pool)
			ctx := GenerateContext{
				RollContext: RollContext{ItemType: "sword", ItemLevel: 50},
				PrefixRange: [2]int{3, 3}, // Try to generate 3, but all in same group
				SuffixRange: [2]int{0, 0},
			}

			instances, err := gen.Generate(ctx)
			require.NoError(t, err)

			// Should only get 1 due to group exclusion
			assert.Equal(t, 1, len(instances))
		})
	})

	t.Run("CreateInstance", func(t *testing.T) {
		t.Run("creates instance with rolled values", func(t *testing.T) {
			pool := NewBasePool()
			gen := NewBaseGenerator(pool)

			affix := NewBaseAffix("create-test", "Create", TypePrefix).
				AddModifier(ModifierTemplate{
					Attribute: attribute.AttrStrength,
					ModType:   attribute.ModFlat,
					MinValue:  10,
					MaxValue:  20,
				})

			instance := gen.CreateInstance(affix)

			assert.Equal(t, "create-test", instance.AffixID())
			assert.Len(t, instance.RolledValues(), 1)

			val := instance.RolledValues()[0].Value
			assert.GreaterOrEqual(t, val, 10.0)
			assert.LessOrEqual(t, val, 20.0)
		})
	})
}

func TestTags(t *testing.T) {
	t.Run("Tag constants exist", func(t *testing.T) {
		assert.Equal(t, Tag("fire"), TagFire)
		assert.Equal(t, Tag("cold"), TagCold)
		assert.Equal(t, Tag("physical"), TagPhysical)
	})

	t.Run("Tag helpers", func(t *testing.T) {
		t.Run("AllElementTags returns element tags", func(t *testing.T) {
			tags := AllElementTags()
			assert.Contains(t, tags, TagFire)
			assert.Contains(t, tags, TagCold)
			assert.Contains(t, tags, TagLightning)
		})

		t.Run("HasTag finds tag in slice", func(t *testing.T) {
			tags := []Tag{TagFire, TagAttack}

			assert.True(t, HasTag(tags, TagFire))
			assert.False(t, HasTag(tags, TagCold))
		})

		t.Run("HasAnyTag matches any", func(t *testing.T) {
			tags := []Tag{TagFire, TagAttack}

			assert.True(t, HasAnyTag(tags, []Tag{TagCold, TagFire}))
			assert.False(t, HasAnyTag(tags, []Tag{TagCold, TagDefense}))
		})

		t.Run("HasAllTags requires all", func(t *testing.T) {
			tags := []Tag{TagFire, TagAttack, TagMelee}

			assert.True(t, HasAllTags(tags, []Tag{TagFire, TagAttack}))
			assert.False(t, HasAllTags(tags, []Tag{TagFire, TagCold}))
		})

		t.Run("Conversion functions", func(t *testing.T) {
			strings := []string{"fire", "cold"}
			tags := StringsToTags(strings)

			assert.Equal(t, []Tag{TagFire, TagCold}, tags)

			backToStrings := TagsToStrings(tags)
			assert.Equal(t, strings, backToStrings)
		})
	})
}

func TestBaseRegistry(t *testing.T) {
	t.Run("Register and Get", func(t *testing.T) {
		t.Run("registers affix successfully", func(t *testing.T) {
			registry := NewBaseRegistry()
			affix := createTestAffix("reg-test", TypePrefix, 50)

			err := registry.Register(affix)
			require.NoError(t, err)

			got, exists := registry.Get("reg-test")
			assert.True(t, exists)
			assert.Equal(t, "reg-test", got.ID())
		})

		t.Run("returns error for duplicate ID", func(t *testing.T) {
			registry := NewBaseRegistry()
			affix1 := createTestAffix("dup", TypePrefix, 50)
			affix2 := createTestAffix("dup", TypeSuffix, 60)

			_ = registry.Register(affix1)
			err := registry.Register(affix2)

			assert.Error(t, err)
		})

		t.Run("GetAll returns all registered", func(t *testing.T) {
			registry := NewBaseRegistry()

			for i := 0; i < 5; i++ {
				_ = registry.Register(createTestAffix("all-"+string(rune('a'+i)), TypePrefix, 50))
			}

			all := registry.GetAll()
			assert.Len(t, all, 5)
		})
	})

	t.Run("GetPool", func(t *testing.T) {
		t.Run("returns pool with eligible affixes", func(t *testing.T) {
			registry := NewBaseRegistry()

			// Add weapon affix
			weaponAffix := createTestAffix("weapon-only", TypePrefix, 50)
			weaponAffix.WithRequirements(NewBaseRequirements(1))
			weaponAffix.requirements.(*BaseRequirements).AddAllowedType("weapon_melee")
			_ = registry.Register(weaponAffix)

			// Add universal affix
			universalAffix := createTestAffix("universal", TypePrefix, 50)
			_ = registry.Register(universalAffix)

			// Get pool for weapon
			weaponPool := registry.GetPool("weapon_melee", "main_hand")
			assert.Len(t, weaponPool.GetAll(), 2)

			// Get pool for armor - should only have universal
			armorPool := registry.GetPool("armor_chest", "chest")
			assert.Len(t, armorPool.GetAll(), 1)
		})
	})
}

func TestRequirements(t *testing.T) {
	t.Run("BaseRequirements", func(t *testing.T) {
		t.Run("Check validates item level", func(t *testing.T) {
			req := NewBaseRequirements(10)
			req.SetMaxItemLevel(50)

			assert.True(t, req.Check("sword", 25, "main_hand"))
			assert.False(t, req.Check("sword", 5, "main_hand"))  // Below min
			assert.False(t, req.Check("sword", 60, "main_hand")) // Above max
		})

		t.Run("Check validates item types", func(t *testing.T) {
			req := NewBaseRequirements(1)
			req.AddAllowedType("weapon_melee")
			req.AddAllowedType("weapon_ranged")

			assert.True(t, req.Check("weapon_melee", 10, "main_hand"))
			assert.False(t, req.Check("armor_chest", 10, "chest"))
		})

		t.Run("Check validates slots", func(t *testing.T) {
			req := NewBaseRequirements(1)
			req.AddAllowedSlot("main_hand")
			req.AddAllowedSlot("off_hand")

			assert.True(t, req.Check("sword", 10, "main_hand"))
			assert.False(t, req.Check("sword", 10, "chest"))
		})

		t.Run("Empty restrictions allow everything", func(t *testing.T) {
			req := NewBaseRequirements(1)

			assert.True(t, req.Check("anything", 100, "anywhere"))
		})
	})
}
