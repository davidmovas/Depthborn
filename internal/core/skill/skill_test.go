package skill

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseDef(t *testing.T) {
	t.Run("создание", func(t *testing.T) {
		t.Run("базовое определение", func(t *testing.T) {
			def := NewBaseDef(DefConfig{
				ID:           "test_skill",
				Name:         "Test Skill",
				Description:  "A test skill",
				Type:         TypeActive,
				Tags:         []string{"fire", "spell"},
				MaxLevel:     5,
				BaseCooldown: 3000,
			})

			require.Equal(t, "test_skill", def.ID())
			require.Equal(t, "Test Skill", def.Name())
			require.Equal(t, TypeActive, def.Type())
			require.Equal(t, 5, def.MaxLevel())
			require.Equal(t, int64(3000), def.BaseCooldown())
		})
	})

	t.Run("теги", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:   "tagged_skill",
			Name: "Tagged Skill",
			Tags: []string{"fire", "spell", "aoe"},
		})

		require.True(t, def.Tags().Has("fire"))
		require.True(t, def.Tags().Has("spell"))
		require.False(t, def.Tags().Has("cold"))
		require.Len(t, def.Tags().All(), 3)
	})

	t.Run("данные уровней", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:       "leveled_skill",
			Name:     "Leveled Skill",
			MaxLevel: 3,
		})

		def.SetLevelData(1, NewBaseLevelData(LevelDataConfig{
			Level: 1,
			Costs: []ResourceCost{
				{Resource: ResourceMana, Type: CostFlat, Amount: 10},
			},
			Description: "Level 1",
		}))

		def.SetLevelData(2, NewBaseLevelData(LevelDataConfig{
			Level: 2,
			Costs: []ResourceCost{
				{Resource: ResourceMana, Type: CostFlat, Amount: 15},
			},
			Description: "Level 2",
		}))

		level1 := def.LevelData(1)
		require.NotNil(t, level1)
		require.Equal(t, "Level 1", level1.Description())

		costs := level1.ResourceCosts()
		require.Len(t, costs, 1)
		require.Equal(t, float64(10), costs[0].Amount)

		level3 := def.LevelData(3)
		require.Nil(t, level3)
	})
}

func TestBaseInstance(t *testing.T) {
	t.Run("создание из определения", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:           "test_skill",
			Name:         "Test Skill",
			MaxLevel:     5,
			BaseCooldown: 3000,
		})

		inst := NewBaseInstance(InstanceConfig{
			Def:        def,
			StartLevel: 1,
		})

		require.Equal(t, "test_skill", inst.DefID())
		require.Equal(t, 1, inst.Level())
		require.True(t, inst.CanLevelUp())
	})

	t.Run("повышение уровня", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:       "levelable",
			Name:     "Levelable Skill",
			MaxLevel: 3,
		})

		inst := NewBaseInstance(InstanceConfig{
			Def:        def,
			StartLevel: 1,
		})

		require.NoError(t, inst.LevelUp())
		require.Equal(t, 2, inst.Level())

		require.NoError(t, inst.LevelUp())
		require.Equal(t, 3, inst.Level())

		require.False(t, inst.CanLevelUp())
		require.ErrorIs(t, inst.LevelUp(), ErrMaxLevel)
	})

	t.Run("cooldown система", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:           "cooldown_skill",
			Name:         "Cooldown Skill",
			BaseCooldown: 5000,
		})

		inst := NewBaseInstance(InstanceConfig{
			Def:        def,
			StartLevel: 1,
		})

		t.Run("начальное состояние", func(t *testing.T) {
			require.False(t, inst.IsOnCooldown())
		})

		t.Run("после активации", func(t *testing.T) {
			inst.SetCooldown(5000)
			require.True(t, inst.IsOnCooldown())
			require.Equal(t, int64(5000), inst.Cooldown())
		})

		t.Run("уменьшение со временем", func(t *testing.T) {
			inst.Update(2000)
			require.Equal(t, int64(3000), inst.Cooldown())
		})

		t.Run("полное восстановление", func(t *testing.T) {
			inst.Update(3000)
			require.False(t, inst.IsOnCooldown())
		})
	})

	t.Run("система зарядов", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:             "charge_skill",
			Name:           "Charge Skill",
			BaseCharges:    3,
			ChargeRecovery: 2000,
		})

		inst := NewBaseInstance(InstanceConfig{
			Def:        def,
			StartLevel: 1,
		})

		t.Run("начальное количество", func(t *testing.T) {
			require.Equal(t, 3, inst.Charges())
			require.Equal(t, 3, inst.MaxCharges())
		})

		t.Run("использование зарядов", func(t *testing.T) {
			require.True(t, inst.UseCharge())
			require.Equal(t, 2, inst.Charges())

			inst.UseCharge()
			inst.UseCharge()
			require.Equal(t, 0, inst.Charges())
			require.False(t, inst.UseCharge())
		})

		t.Run("восстановление зарядов", func(t *testing.T) {
			inst.Update(2000)
			require.Equal(t, 1, inst.Charges())

			inst.Update(4000)
			require.Equal(t, 3, inst.Charges())
		})
	})

	t.Run("активация скилла", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:           "usable_skill",
			Name:         "Usable Skill",
			MaxLevel:     1,
			BaseCooldown: 1000,
		})

		def.SetLevelData(1, NewBaseLevelData(LevelDataConfig{
			Level: 1,
			Costs: []ResourceCost{
				{Resource: ResourceMana, Type: CostFlat, Amount: 20},
			},
		}))

		inst := NewBaseInstance(InstanceConfig{
			Def:        def,
			StartLevel: 1,
		})

		ctx := context.Background()

		require.True(t, inst.CanUse(ctx, "player1"))

		result, err := inst.Use(ctx, "player1", ActivationParams{})
		require.NoError(t, err)
		require.True(t, result.Success)

		require.True(t, inst.IsOnCooldown())
		require.False(t, inst.CanUse(ctx, "player1"))
	})
}

func TestBaseTargetRule(t *testing.T) {
	t.Run("single target", func(t *testing.T) {
		rule := NewBaseTargetRule(TargetRuleConfig{
			Type:       TargetSingle,
			Range:      20,
			CanEnemies: true,
		})

		require.Equal(t, TargetSingle, rule.Type())
		require.Equal(t, float64(20), rule.Range())
		require.True(t, rule.CanTargetEnemies())
		require.False(t, rule.CanTargetAllies())
	})

	t.Run("AoE", func(t *testing.T) {
		rule := NewBaseTargetRule(TargetRuleConfig{
			Type:        TargetGround,
			AreaType:    AreaCircle,
			Range:       25,
			AreaRadius:  5,
			MaxTargets:  10,
			CanSelf:     true,
			CanAllies:   true,
			CanEnemies:  true,
			RequiresLOS: true,
		})

		require.Equal(t, AreaCircle, rule.AreaType())
		require.Equal(t, float64(5), rule.AreaRadius())
		require.Equal(t, 10, rule.MaxTargets())
		require.True(t, rule.RequiresLineOfSight())
	})

	t.Run("chain", func(t *testing.T) {
		rule := NewBaseTargetRule(TargetRuleConfig{
			Type:        TargetSingle,
			AreaType:    AreaChain,
			ChainCount:  5,
			ChainFallof: 0.25,
		})

		require.Equal(t, AreaChain, rule.AreaType())
		require.Equal(t, 5, rule.ChainCount())
		require.Equal(t, 0.25, rule.ChainFalloff())
	})
}

func TestBaseEffectDef(t *testing.T) {
	t.Run("damage effect", func(t *testing.T) {
		effect := NewBaseEffectDef(EffectDefConfig{
			ID:         "fire_damage",
			Type:       EffectDamage,
			DamageType: "fire",
			Scaling: []ScalingRule{
				{Attribute: "intelligence", Multiplier: 0.5},
			},
			Chance: 1.0,
		})

		require.Equal(t, "fire_damage", effect.ID())
		require.Equal(t, EffectDamage, effect.Type())
		require.Equal(t, "fire", effect.DamageType())
		require.Equal(t, 1.0, effect.Chance())

		scaling := effect.Scaling()
		require.Len(t, scaling, 1)
		require.Equal(t, "intelligence", scaling[0].Attribute)
	})

	t.Run("status effect", func(t *testing.T) {
		effect := NewBaseEffectDef(EffectDefConfig{
			ID:       "apply_burn",
			Type:     EffectStatus,
			StatusID: "burning",
			Chance:   0.25,
			Duration: 5000,
		})

		require.Equal(t, EffectStatus, effect.Type())
		require.Equal(t, "burning", effect.StatusID())
		require.Equal(t, 0.25, effect.Chance())
		require.Equal(t, int64(5000), effect.Duration())
	})
}

func TestTags(t *testing.T) {
	t.Run("TagSet методы", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:   "test",
			Name: "Test",
			Tags: []string{"fire", "spell", "damage"},
		})
		tags := def.Tags()

		require.True(t, tags.Has("fire"))
		require.False(t, tags.Has("cold"))
		require.True(t, tags.ContainsAny("cold", "fire"))
		require.False(t, tags.ContainsAny("cold", "lightning"))
		require.True(t, tags.Contains("fire", "spell"))
		require.False(t, tags.Contains("fire", "cold"))
	})

	t.Run("GetElementTags", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:   "elemental",
			Name: "Elemental",
			Tags: []string{"fire", "spell", "cold", "damage"},
		})
		elements := GetElementTags(def.Tags())
		require.Len(t, elements, 2)
	})

	t.Run("IsDamageSkill", func(t *testing.T) {
		damageSkill := NewBaseDef(DefConfig{
			ID: "dmg", Name: "Damage", Tags: []string{"damage"},
		})
		healSkill := NewBaseDef(DefConfig{
			ID: "heal", Name: "Heal", Tags: []string{"heal"},
		})

		require.True(t, IsDamageSkill(damageSkill.Tags()))
		require.False(t, IsDamageSkill(healSkill.Tags()))
	})
}

func TestRegistry(t *testing.T) {
	t.Run("регистрация и получение", func(t *testing.T) {
		registry := NewBaseRegistry()

		def := NewBaseDef(DefConfig{
			ID:   "test_skill",
			Name: "Test Skill",
			Type: TypeActive,
			Tags: []string{"fire"},
		})

		require.NoError(t, registry.Register(def))

		retrieved, ok := registry.Get("test_skill")
		require.True(t, ok)
		require.Equal(t, "test_skill", retrieved.ID())
	})

	t.Run("повторная регистрация", func(t *testing.T) {
		registry := NewBaseRegistry()

		def := NewBaseDef(DefConfig{ID: "dup_skill", Name: "Dup Skill"})
		require.NoError(t, registry.Register(def))

		def2 := NewBaseDef(DefConfig{ID: "dup_skill", Name: "Dup Skill 2"})
		require.Error(t, registry.Register(def2))
	})

	t.Run("поиск по тегу", func(t *testing.T) {
		registry := NewBaseRegistry()

		registry.Register(NewBaseDef(DefConfig{
			ID: "fire1", Name: "Fire 1", Tags: []string{"fire"},
		}))
		registry.Register(NewBaseDef(DefConfig{
			ID: "fire2", Name: "Fire 2", Tags: []string{"fire", "spell"},
		}))
		registry.Register(NewBaseDef(DefConfig{
			ID: "cold1", Name: "Cold 1", Tags: []string{"cold"},
		}))

		fireSkills := registry.GetByTag("fire")
		require.Len(t, fireSkills, 2)

		spellSkills := registry.GetByTag("spell")
		require.Len(t, spellSkills, 1)
	})

	t.Run("поиск по типу", func(t *testing.T) {
		registry := NewBaseRegistry()

		registry.Register(NewBaseDef(DefConfig{
			ID: "active1", Name: "Active 1", Type: TypeActive,
		}))
		registry.Register(NewBaseDef(DefConfig{
			ID: "passive1", Name: "Passive 1", Type: TypePassive,
		}))
		registry.Register(NewBaseDef(DefConfig{
			ID: "active2", Name: "Active 2", Type: TypeActive,
		}))

		activeSkills := registry.GetByType(TypeActive)
		require.Len(t, activeSkills, 2)
	})

	t.Run("создание экземпляра", func(t *testing.T) {
		registry := NewBaseRegistry()

		def := NewBaseDef(DefConfig{
			ID:       "instanceable",
			Name:     "Instanceable Skill",
			MaxLevel: 5,
		})
		registry.Register(def)

		inst, err := registry.CreateInstance("instanceable", 2)
		require.NoError(t, err)
		require.Equal(t, 2, inst.Level())
		require.NotNil(t, inst.Def())
	})

	t.Run("загрузка из YAML", func(t *testing.T) {
		registry := NewBaseRegistry()

		yaml := `
version: "1.0"
skills:
  - id: yaml_skill
    name: "YAML Skill"
    description: "Loaded from YAML"
    type: active
    tags: [fire, spell]
    max_level: 3
    cooldown: 2000
    targeting:
      type: single
      range: 20
      can_enemies: true
    effects:
      - id: damage
        type: damage
        damage_type: fire
    levels:
      - level: 1
        costs:
          - resource: mana
            type: flat
            amount: 10
        description: "Level 1"
      - level: 2
        costs:
          - resource: mana
            type: flat
            amount: 15
        description: "Level 2"
      - level: 3
        costs:
          - resource: mana
            type: flat
            amount: 20
        description: "Level 3"
`

		require.NoError(t, registry.LoadFromYAML([]byte(yaml)))

		def, ok := registry.Get("yaml_skill")
		require.True(t, ok)
		require.Equal(t, "YAML Skill", def.Name())
		require.Equal(t, TypeActive, def.Type())
		require.Equal(t, 3, def.MaxLevel())
		require.Equal(t, int64(2000), def.BaseCooldown())
		require.True(t, def.Tags().Has("fire"))

		targeting := def.Targeting()
		require.Equal(t, TargetSingle, targeting.Type())
		require.Equal(t, float64(20), targeting.Range())

		effects := def.Effects()
		require.Len(t, effects, 1)
		require.Equal(t, EffectDamage, effects[0].Type())

		level1 := def.LevelData(1)
		require.NotNil(t, level1)
		costs := level1.ResourceCosts()
		require.Len(t, costs, 1)
		require.Equal(t, float64(10), costs[0].Amount)
	})
}

func TestLoadRealYAMLFiles(t *testing.T) {
	registry := NewBaseRegistry()

	err := registry.LoadFromDirectory("../../../data/skills")
	require.NoError(t, err, "failed to load skills from directory")
	require.Greater(t, registry.Count(), 0, "expected to load skills from YAML files")

	t.Logf("Loaded %d skills from YAML files", registry.Count())

	t.Run("active skills", func(t *testing.T) {
		fireball, ok := registry.Get("fireball")
		require.True(t, ok, "expected 'fireball' skill to be loaded")
		require.Equal(t, TypeActive, fireball.Type())
		require.Equal(t, 5, fireball.MaxLevel())
		require.True(t, fireball.Tags().Has("fire"))

		level1 := fireball.LevelData(1)
		require.NotNil(t, level1)
		require.NotEmpty(t, level1.ResourceCosts())
	})

	t.Run("passive skills", func(t *testing.T) {
		toughness, ok := registry.Get("toughness")
		require.True(t, ok, "expected 'toughness' skill to be loaded")
		require.Equal(t, TypePassive, toughness.Type())
	})

	t.Run("keystone passives", func(t *testing.T) {
		bloodMagic, ok := registry.Get("blood_magic")
		require.True(t, ok, "expected 'blood_magic' skill to be loaded")
		require.True(t, bloodMagic.Tags().Has("keystone"))
	})

	t.Run("поиск по тегу", func(t *testing.T) {
		fireSkills := registry.GetByTag("fire")
		require.NotEmpty(t, fireSkills)
		t.Logf("Found %d fire skills", len(fireSkills))

		keystones := registry.GetByTag("keystone")
		require.GreaterOrEqual(t, len(keystones), 4)
	})

	t.Run("поиск по типу", func(t *testing.T) {
		activeSkills := registry.GetByType(TypeActive)
		passiveSkills := registry.GetByType(TypePassive)

		t.Logf("Active skills: %d, Passive skills: %d", len(activeSkills), len(passiveSkills))

		require.NotEmpty(t, activeSkills)
		require.NotEmpty(t, passiveSkills)
	})

	t.Run("создание экземпляра из загруженного", func(t *testing.T) {
		inst, err := registry.CreateInstance("fireball", 3)
		require.NoError(t, err)
		require.Equal(t, 3, inst.Level())
		require.Equal(t, "fireball", inst.DefID())
	})
}
