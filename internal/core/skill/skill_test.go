package skill

import (
	"context"
	"testing"
)

func TestBaseDef(t *testing.T) {
	t.Run("create basic skill definition", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:           "test_skill",
			Name:         "Test Skill",
			Description:  "A test skill",
			Type:         TypeActive,
			Tags:         []string{"fire", "spell"},
			MaxLevel:     5,
			BaseCooldown: 3000,
		})

		if def.ID() != "test_skill" {
			t.Errorf("expected ID 'test_skill', got %q", def.ID())
		}
		if def.Name() != "Test Skill" {
			t.Errorf("expected name 'Test Skill', got %q", def.Name())
		}
		if def.Type() != TypeActive {
			t.Errorf("expected type TypeActive, got %v", def.Type())
		}
		if def.MaxLevel() != 5 {
			t.Errorf("expected max level 5, got %d", def.MaxLevel())
		}
		if def.BaseCooldown() != 3000 {
			t.Errorf("expected cooldown 3000, got %d", def.BaseCooldown())
		}
	})

	t.Run("skill tags", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:   "tagged_skill",
			Name: "Tagged Skill",
			Tags: []string{"fire", "spell", "aoe"},
		})

		if !def.HasTag("fire") {
			t.Error("expected skill to have 'fire' tag")
		}
		if !def.HasTag("spell") {
			t.Error("expected skill to have 'spell' tag")
		}
		if def.HasTag("cold") {
			t.Error("expected skill to NOT have 'cold' tag")
		}

		tags := def.Tags()
		if len(tags) != 3 {
			t.Errorf("expected 3 tags, got %d", len(tags))
		}
	})

	t.Run("skill with level data", func(t *testing.T) {
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
		if level1 == nil {
			t.Fatal("expected level 1 data to exist")
		}
		if level1.Description() != "Level 1" {
			t.Errorf("expected description 'Level 1', got %q", level1.Description())
		}

		costs := level1.ResourceCosts()
		if len(costs) != 1 {
			t.Fatalf("expected 1 cost, got %d", len(costs))
		}
		if costs[0].Amount != 10 {
			t.Errorf("expected cost 10, got %f", costs[0].Amount)
		}

		level3 := def.LevelData(3)
		if level3 != nil {
			t.Error("expected level 3 data to be nil")
		}
	})
}

func TestBaseInstance(t *testing.T) {
	t.Run("create instance from definition", func(t *testing.T) {
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

		if inst.DefID() != "test_skill" {
			t.Errorf("expected def ID 'test_skill', got %q", inst.DefID())
		}
		if inst.Level() != 1 {
			t.Errorf("expected level 1, got %d", inst.Level())
		}
		if !inst.CanLevelUp() {
			t.Error("expected CanLevelUp to be true")
		}
	})

	t.Run("level up skill", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:       "levelable",
			Name:     "Levelable Skill",
			MaxLevel: 3,
		})

		inst := NewBaseInstance(InstanceConfig{
			Def:        def,
			StartLevel: 1,
		})

		if err := inst.LevelUp(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if inst.Level() != 2 {
			t.Errorf("expected level 2, got %d", inst.Level())
		}

		if err := inst.LevelUp(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if inst.Level() != 3 {
			t.Errorf("expected level 3, got %d", inst.Level())
		}

		// Should not be able to level up at max
		if inst.CanLevelUp() {
			t.Error("expected CanLevelUp to be false at max level")
		}

		err := inst.LevelUp()
		if err != ErrMaxLevel {
			t.Errorf("expected ErrMaxLevel, got %v", err)
		}
	})

	t.Run("cooldown management", func(t *testing.T) {
		def := NewBaseDef(DefConfig{
			ID:           "cooldown_skill",
			Name:         "Cooldown Skill",
			BaseCooldown: 5000,
		})

		inst := NewBaseInstance(InstanceConfig{
			Def:        def,
			StartLevel: 1,
		})

		if inst.IsOnCooldown() {
			t.Error("expected not on cooldown initially")
		}

		inst.SetCooldown(5000)
		if !inst.IsOnCooldown() {
			t.Error("expected on cooldown after setting")
		}
		if inst.Cooldown() != 5000 {
			t.Errorf("expected cooldown 5000, got %d", inst.Cooldown())
		}

		// Update to reduce cooldown
		inst.Update(2000)
		if inst.Cooldown() != 3000 {
			t.Errorf("expected cooldown 3000 after update, got %d", inst.Cooldown())
		}

		inst.Update(3000)
		if inst.IsOnCooldown() {
			t.Error("expected not on cooldown after full update")
		}
	})

	t.Run("charge system", func(t *testing.T) {
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

		if inst.Charges() != 3 {
			t.Errorf("expected 3 charges, got %d", inst.Charges())
		}
		if inst.MaxCharges() != 3 {
			t.Errorf("expected max charges 3, got %d", inst.MaxCharges())
		}

		// Use charges
		if !inst.UseCharge() {
			t.Error("expected UseCharge to succeed")
		}
		if inst.Charges() != 2 {
			t.Errorf("expected 2 charges, got %d", inst.Charges())
		}

		inst.UseCharge()
		inst.UseCharge()
		if inst.Charges() != 0 {
			t.Errorf("expected 0 charges, got %d", inst.Charges())
		}

		if inst.UseCharge() {
			t.Error("expected UseCharge to fail with 0 charges")
		}

		// Recovery
		inst.Update(2000) // One charge recovered
		if inst.Charges() != 1 {
			t.Errorf("expected 1 charge after recovery, got %d", inst.Charges())
		}

		inst.Update(4000) // Two more charges recovered
		if inst.Charges() != 3 {
			t.Errorf("expected 3 charges after full recovery, got %d", inst.Charges())
		}
	})

	t.Run("skill activation", func(t *testing.T) {
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

		if !inst.CanUse(ctx, "player1") {
			t.Error("expected CanUse to be true")
		}

		result, err := inst.Use(ctx, "player1", ActivationParams{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("expected skill use to succeed")
		}

		// Should be on cooldown now
		if !inst.IsOnCooldown() {
			t.Error("expected on cooldown after use")
		}

		// Can't use while on cooldown
		if inst.CanUse(ctx, "player1") {
			t.Error("expected CanUse to be false while on cooldown")
		}
	})
}

func TestBaseTargetRule(t *testing.T) {
	t.Run("single target rule", func(t *testing.T) {
		rule := NewBaseTargetRule(TargetRuleConfig{
			Type:       TargetSingle,
			Range:      20,
			CanEnemies: true,
		})

		if rule.Type() != TargetSingle {
			t.Errorf("expected TargetSingle, got %v", rule.Type())
		}
		if rule.Range() != 20 {
			t.Errorf("expected range 20, got %f", rule.Range())
		}
		if !rule.CanTargetEnemies() {
			t.Error("expected CanTargetEnemies to be true")
		}
		if rule.CanTargetAllies() {
			t.Error("expected CanTargetAllies to be false")
		}
	})

	t.Run("aoe rule", func(t *testing.T) {
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

		if rule.AreaType() != AreaCircle {
			t.Errorf("expected AreaCircle, got %v", rule.AreaType())
		}
		if rule.AreaRadius() != 5 {
			t.Errorf("expected area radius 5, got %f", rule.AreaRadius())
		}
		if rule.MaxTargets() != 10 {
			t.Errorf("expected max targets 10, got %d", rule.MaxTargets())
		}
		if !rule.RequiresLineOfSight() {
			t.Error("expected RequiresLineOfSight to be true")
		}
	})

	t.Run("chain rule", func(t *testing.T) {
		rule := NewBaseTargetRule(TargetRuleConfig{
			Type:        TargetSingle,
			AreaType:    AreaChain,
			ChainCount:  5,
			ChainFallof: 0.25,
		})

		if rule.AreaType() != AreaChain {
			t.Errorf("expected AreaChain, got %v", rule.AreaType())
		}
		if rule.ChainCount() != 5 {
			t.Errorf("expected chain count 5, got %d", rule.ChainCount())
		}
		if rule.ChainFalloff() != 0.25 {
			t.Errorf("expected chain falloff 0.25, got %f", rule.ChainFalloff())
		}
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

		if effect.ID() != "fire_damage" {
			t.Errorf("expected ID 'fire_damage', got %q", effect.ID())
		}
		if effect.Type() != EffectDamage {
			t.Errorf("expected EffectDamage, got %v", effect.Type())
		}
		if effect.DamageType() != "fire" {
			t.Errorf("expected damage type 'fire', got %q", effect.DamageType())
		}
		if effect.Chance() != 1.0 {
			t.Errorf("expected chance 1.0, got %f", effect.Chance())
		}

		scaling := effect.Scaling()
		if len(scaling) != 1 {
			t.Fatalf("expected 1 scaling rule, got %d", len(scaling))
		}
		if scaling[0].Attribute != "intelligence" {
			t.Errorf("expected attribute 'intelligence', got %q", scaling[0].Attribute)
		}
	})

	t.Run("status effect", func(t *testing.T) {
		effect := NewBaseEffectDef(EffectDefConfig{
			ID:       "apply_burn",
			Type:     EffectStatus,
			StatusID: "burning",
			Chance:   0.25,
			Duration: 5000,
		})

		if effect.Type() != EffectStatus {
			t.Errorf("expected EffectStatus, got %v", effect.Type())
		}
		if effect.StatusID() != "burning" {
			t.Errorf("expected status ID 'burning', got %q", effect.StatusID())
		}
		if effect.Chance() != 0.25 {
			t.Errorf("expected chance 0.25, got %f", effect.Chance())
		}
		if effect.Duration() != 5000 {
			t.Errorf("expected duration 5000, got %d", effect.Duration())
		}
	})
}

func TestTags(t *testing.T) {
	t.Run("HasTag", func(t *testing.T) {
		tags := []string{"fire", "spell", "aoe"}

		if !HasTag(tags, "fire") {
			t.Error("expected to find 'fire' tag")
		}
		if HasTag(tags, "cold") {
			t.Error("expected not to find 'cold' tag")
		}
	})

	t.Run("HasAnyTag", func(t *testing.T) {
		tags := []string{"fire", "spell"}

		if !HasAnyTag(tags, []string{"cold", "fire"}) {
			t.Error("expected to find at least one tag")
		}
		if HasAnyTag(tags, []string{"cold", "lightning"}) {
			t.Error("expected not to find any tag")
		}
	})

	t.Run("HasAllTags", func(t *testing.T) {
		tags := []string{"fire", "spell", "damage"}

		if !HasAllTags(tags, []string{"fire", "spell"}) {
			t.Error("expected to find all tags")
		}
		if HasAllTags(tags, []string{"fire", "cold"}) {
			t.Error("expected not to find all tags")
		}
	})

	t.Run("GetElementTags", func(t *testing.T) {
		tags := []string{"fire", "spell", "cold", "damage"}
		elements := GetElementTags(tags)

		if len(elements) != 2 {
			t.Errorf("expected 2 element tags, got %d", len(elements))
		}
	})
}

func TestLoadRealYAMLFiles(t *testing.T) {
	registry := NewBaseRegistry()

	// Load from data directory
	err := registry.LoadFromDirectory("../../../data/skills")
	if err != nil {
		t.Fatalf("failed to load skills from directory: %v", err)
	}

	// Check we loaded expected skills
	if registry.Count() == 0 {
		t.Fatal("expected to load skills from YAML files")
	}

	t.Logf("Loaded %d skills from YAML files", registry.Count())

	// Check specific skills exist
	t.Run("active skills loaded", func(t *testing.T) {
		fireball, ok := registry.Get("fireball")
		if !ok {
			t.Fatal("expected 'fireball' skill to be loaded")
		}
		if fireball.Type() != TypeActive {
			t.Errorf("expected fireball to be active, got %v", fireball.Type())
		}
		if fireball.MaxLevel() != 5 {
			t.Errorf("expected fireball max level 5, got %d", fireball.MaxLevel())
		}
		if !fireball.HasTag("fire") {
			t.Error("expected fireball to have 'fire' tag")
		}

		// Check level data
		level1 := fireball.LevelData(1)
		if level1 == nil {
			t.Fatal("expected fireball level 1 data")
		}
		costs := level1.ResourceCosts()
		if len(costs) == 0 {
			t.Error("expected fireball to have resource costs")
		}
	})

	t.Run("passive skills loaded", func(t *testing.T) {
		toughness, ok := registry.Get("toughness")
		if !ok {
			t.Fatal("expected 'toughness' skill to be loaded")
		}
		if toughness.Type() != TypePassive {
			t.Errorf("expected toughness to be passive, got %v", toughness.Type())
		}
	})

	t.Run("keystone passives loaded", func(t *testing.T) {
		bloodMagic, ok := registry.Get("blood_magic")
		if !ok {
			t.Fatal("expected 'blood_magic' skill to be loaded")
		}
		if !bloodMagic.HasTag("keystone") {
			t.Error("expected blood_magic to have 'keystone' tag")
		}
	})

	t.Run("query by tag", func(t *testing.T) {
		fireSkills := registry.GetByTag("fire")
		if len(fireSkills) == 0 {
			t.Error("expected to find fire skills")
		}
		t.Logf("Found %d fire skills", len(fireSkills))

		keystones := registry.GetByTag("keystone")
		if len(keystones) < 4 {
			t.Errorf("expected at least 4 keystone skills, got %d", len(keystones))
		}
	})

	t.Run("query by type", func(t *testing.T) {
		activeSkills := registry.GetByType(TypeActive)
		passiveSkills := registry.GetByType(TypePassive)

		t.Logf("Active skills: %d, Passive skills: %d", len(activeSkills), len(passiveSkills))

		if len(activeSkills) == 0 {
			t.Error("expected to find active skills")
		}
		if len(passiveSkills) == 0 {
			t.Error("expected to find passive skills")
		}
	})

	t.Run("create instance from loaded skill", func(t *testing.T) {
		inst, err := registry.CreateInstance("fireball", 3)
		if err != nil {
			t.Fatalf("failed to create instance: %v", err)
		}

		if inst.Level() != 3 {
			t.Errorf("expected level 3, got %d", inst.Level())
		}
		if inst.DefID() != "fireball" {
			t.Errorf("expected def ID 'fireball', got %q", inst.DefID())
		}
	})
}

func TestRegistry(t *testing.T) {
	t.Run("register and get", func(t *testing.T) {
		registry := NewBaseRegistry()

		def := NewBaseDef(DefConfig{
			ID:   "test_skill",
			Name: "Test Skill",
			Type: TypeActive,
			Tags: []string{"fire"},
		})

		if err := registry.Register(def); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		retrieved, ok := registry.Get("test_skill")
		if !ok {
			t.Fatal("expected to find skill")
		}
		if retrieved.ID() != "test_skill" {
			t.Errorf("expected ID 'test_skill', got %q", retrieved.ID())
		}
	})

	t.Run("duplicate registration", func(t *testing.T) {
		registry := NewBaseRegistry()

		def := NewBaseDef(DefConfig{ID: "dup_skill", Name: "Dup Skill"})
		if err := registry.Register(def); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		def2 := NewBaseDef(DefConfig{ID: "dup_skill", Name: "Dup Skill 2"})
		if err := registry.Register(def2); err == nil {
			t.Error("expected error for duplicate registration")
		}
	})

	t.Run("get by tag", func(t *testing.T) {
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
		if len(fireSkills) != 2 {
			t.Errorf("expected 2 fire skills, got %d", len(fireSkills))
		}

		spellSkills := registry.GetByTag("spell")
		if len(spellSkills) != 1 {
			t.Errorf("expected 1 spell skill, got %d", len(spellSkills))
		}
	})

	t.Run("get by type", func(t *testing.T) {
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
		if len(activeSkills) != 2 {
			t.Errorf("expected 2 active skills, got %d", len(activeSkills))
		}
	})

	t.Run("create instance", func(t *testing.T) {
		registry := NewBaseRegistry()

		def := NewBaseDef(DefConfig{
			ID:       "instanceable",
			Name:     "Instanceable Skill",
			MaxLevel: 5,
		})
		registry.Register(def)

		inst, err := registry.CreateInstance("instanceable", 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if inst.Level() != 2 {
			t.Errorf("expected level 2, got %d", inst.Level())
		}
		if inst.Def() == nil {
			t.Error("expected def to be set")
		}
	})

	t.Run("load from YAML", func(t *testing.T) {
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

		if err := registry.LoadFromYAML([]byte(yaml)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		def, ok := registry.Get("yaml_skill")
		if !ok {
			t.Fatal("expected to find yaml_skill")
		}

		if def.Name() != "YAML Skill" {
			t.Errorf("expected name 'YAML Skill', got %q", def.Name())
		}
		if def.Type() != TypeActive {
			t.Errorf("expected TypeActive, got %v", def.Type())
		}
		if def.MaxLevel() != 3 {
			t.Errorf("expected max level 3, got %d", def.MaxLevel())
		}
		if def.BaseCooldown() != 2000 {
			t.Errorf("expected cooldown 2000, got %d", def.BaseCooldown())
		}
		if !def.HasTag("fire") {
			t.Error("expected skill to have 'fire' tag")
		}

		targeting := def.Targeting()
		if targeting.Type() != TargetSingle {
			t.Errorf("expected TargetSingle, got %v", targeting.Type())
		}
		if targeting.Range() != 20 {
			t.Errorf("expected range 20, got %f", targeting.Range())
		}

		effects := def.Effects()
		if len(effects) != 1 {
			t.Fatalf("expected 1 effect, got %d", len(effects))
		}
		if effects[0].Type() != EffectDamage {
			t.Errorf("expected EffectDamage, got %v", effects[0].Type())
		}

		level1 := def.LevelData(1)
		if level1 == nil {
			t.Fatal("expected level 1 data")
		}
		costs := level1.ResourceCosts()
		if len(costs) != 1 {
			t.Fatalf("expected 1 cost, got %d", len(costs))
		}
		if costs[0].Amount != 10 {
			t.Errorf("expected cost 10, got %f", costs[0].Amount)
		}
	})
}
