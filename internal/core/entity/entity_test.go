package entity_test

import (
	"context"
	"testing"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/core/status"
	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

func createTestCombatant(name string) *entity.BaseCombatant {
	attrMgr := attribute.NewManager()
	attrMgr.SetBase(attribute.AttrStrength, 10)
	attrMgr.SetBase(attribute.AttrDexterity, 10)
	attrMgr.SetBase(attribute.AttrVitality, 10)
	attrMgr.SetBase(attribute.AttrPhysicalDamage, 5)
	attrMgr.SetBase(attribute.AttrCritChance, 5)
	attrMgr.SetBase(attribute.AttrCritMultiplier, 1.5)
	attrMgr.SetBase(attribute.AttrArmor, 10)
	attrMgr.SetBase(attribute.AttrEvasion, 10)

	config := entity.CombatantConfig{
		LivingConfig: entity.LivingConfig{
			EntityConfig: entity.Config{
				Name:             name,
				EntityType:       "test_combatant",
				AttributeManager: attrMgr,
				StatusManager:    status.NewManager(),
				Transform:        spatial.NewTransform(spatial.NewPosition(0, 0, 0), spatial.FacingNorth),
				TagSet:           types.NewTagSet(),
				Callbacks:        types.NewCallbackRegistry(),
			},
			InitialHealth: 100,
			MaxHealth:     100,
		},
		AttackRange: 1.5,
	}

	return entity.NewCombatant(config)
}

func TestCombatantCreation(t *testing.T) {
	combatant := createTestCombatant("TestWarrior")

	if combatant.Name() != "TestWarrior" {
		t.Errorf("expected name TestWarrior, got %s", combatant.Name())
	}

	if combatant.Health() != 100 {
		t.Errorf("expected health 100, got %f", combatant.Health())
	}

	// MaxHealth = base(100) + vitality(10) * 10 = 200
	if combatant.MaxHealth() != 200 {
		t.Errorf("expected max health 200, got %f", combatant.MaxHealth())
	}

	if !combatant.IsAlive() {
		t.Error("expected combatant to be alive")
	}
}

func TestCombatantDamage(t *testing.T) {
	ctx := context.Background()
	combatant := createTestCombatant("TestWarrior")

	initialHealth := combatant.Health()
	damageDealt, err := combatant.Damage(ctx, 30, "attacker1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if damageDealt != 30 {
		t.Errorf("expected damage 30, got %f", damageDealt)
	}

	expectedHealth := initialHealth - 30
	if combatant.Health() != expectedHealth {
		t.Errorf("expected health %f, got %f", expectedHealth, combatant.Health())
	}
}

func TestCombatantHeal(t *testing.T) {
	ctx := context.Background()
	combatant := createTestCombatant("TestWarrior")

	// First take some damage
	_, _ = combatant.Damage(ctx, 50, "attacker1")

	healAmount, err := combatant.Heal(ctx, 20, "healer1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if healAmount != 20 {
		t.Errorf("expected heal 20, got %f", healAmount)
	}
}

func TestCombatantKill(t *testing.T) {
	ctx := context.Background()
	combatant := createTestCombatant("TestWarrior")

	// Deal lethal damage (more than initial health of 100)
	_, _ = combatant.Damage(ctx, 150, "attacker1")

	if combatant.IsAlive() {
		t.Error("expected combatant to be dead after lethal damage")
	}

	if combatant.Health() != 0 {
		t.Errorf("expected health 0 after death, got %f", combatant.Health())
	}
}

func TestCombatantClone(t *testing.T) {
	combatant := createTestCombatant("TestWarrior")
	combatant.Tags().Add("warrior")
	combatant.Tags().Add("player")

	cloned := combatant.Clone().(*entity.BaseCombatant)

	if cloned.Name() != combatant.Name() {
		t.Errorf("clone name mismatch: got %s, want %s", cloned.Name(), combatant.Name())
	}

	if cloned.Health() != combatant.Health() {
		t.Errorf("clone health mismatch: got %f, want %f", cloned.Health(), combatant.Health())
	}

	if cloned.ID() == combatant.ID() {
		t.Error("clone should have different ID")
	}

	// Verify tags were cloned
	if !cloned.Tags().Has("warrior") {
		t.Error("clone should have 'warrior' tag")
	}

	// Modify clone should not affect original
	cloned.SetName("ClonedWarrior")
	if combatant.Name() == "ClonedWarrior" {
		t.Error("modifying clone should not affect original")
	}
}

func TestCombatantAttack(t *testing.T) {
	ctx := context.Background()
	attacker := createTestCombatant("Attacker")

	result, err := attacker.Attack(ctx, "target123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Hit {
		t.Error("expected attack to hit")
	}

	// Base damage = physicalDamage(5) + strength(10) * 2 = 25
	if result.Damage < 20 { // At least 20, could be higher with crit
		t.Errorf("expected damage >= 20, got %f", result.Damage)
	}
}

func TestCombatantDefend(t *testing.T) {
	ctx := context.Background()
	defender := createTestCombatant("Defender")

	attack := entity.AttackInfo{
		AttackerID:   "attacker123",
		BaseDamage:   50,
		DamageType:   "physical",
		IsCritical:   false,
		Penetration:  0,
		StatusChance: make(map[string]float64),
	}

	result, err := defender.Defend(ctx, attack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// With armor mitigation, final damage should be less than base
	// Unless evaded (which is random)
	if !result.Evaded && result.FinalDamage >= attack.BaseDamage {
		t.Errorf("expected mitigation, final damage %f >= base %f", result.FinalDamage, attack.BaseDamage)
	}
}

func TestCombatantValidation(t *testing.T) {
	combatant := createTestCombatant("TestWarrior")

	if err := combatant.Validate(); err != nil {
		t.Errorf("valid combatant failed validation: %v", err)
	}
}

func TestCombatantMarshalUnmarshal(t *testing.T) {
	combatant := createTestCombatant("TestWarrior")
	combatant.Tags().Add("test_tag")

	// Marshal
	data, err := combatant.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("marshal returned empty data")
	}

	// Create new combatant with initialized components for unmarshal
	newCombatant := &entity.BaseCombatant{}

	// Unmarshal
	if err = newCombatant.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if newCombatant.Name() != combatant.Name() {
		t.Errorf("name mismatch after unmarshal: got %s, want %s", newCombatant.Name(), combatant.Name())
	}
}

func TestCombatantCanAct(t *testing.T) {
	combatant := createTestCombatant("TestWarrior")

	if !combatant.CanAct() {
		t.Error("living combatant should be able to act")
	}

	// Kill the combatant via lethal damage
	ctx := context.Background()
	_, _ = combatant.Damage(ctx, 500, "killer1")

	if combatant.CanAct() {
		t.Error("dead combatant should not be able to act")
	}
}

func TestCombatantThreatLevel(t *testing.T) {
	combatant := createTestCombatant("TestWarrior")

	if combatant.ThreatLevel() != 0 {
		t.Errorf("initial threat should be 0, got %f", combatant.ThreatLevel())
	}

	combatant.ModifyThreat(50)
	if combatant.ThreatLevel() != 50 {
		t.Errorf("expected threat 50, got %f", combatant.ThreatLevel())
	}

	combatant.ModifyThreat(-100) // Should not go below 0
	if combatant.ThreatLevel() != 0 {
		t.Errorf("threat should not go below 0, got %f", combatant.ThreatLevel())
	}
}
