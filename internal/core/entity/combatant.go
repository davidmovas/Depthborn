package entity

import (
	"context"
	"fmt"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
)

var _ Combatant = (*BaseCombatant)(nil)

type BaseCombatant struct {
	*BaseLiving

	threatLevel float64
	attackRange float64
}

type CombatantConfig struct {
	LivingConfig
	AttackRange float64
}

func NewCombatant(config CombatantConfig) *BaseCombatant {
	if config.AttackRange <= 0 {
		config.AttackRange = 1.5 // Default melee range
	}

	return &BaseCombatant{
		BaseLiving:  NewLiving(config.LivingConfig),
		threatLevel: 0,
		attackRange: config.AttackRange,
	}
}

func (c *BaseCombatant) Attack(ctx context.Context, targetID string) (CombatResult, error) {
	if !c.IsAlive() {
		return CombatResult{}, fmt.Errorf("dead entities cannot attack")
	}

	if !c.CanAct() {
		return CombatResult{}, fmt.Errorf("entity cannot act")
	}

	// TODO: Implement full combat logic
	// - Calculate base damage from attributes
	// - Roll for hit/miss based on accuracy and target evasion
	// - Roll for critical hit
	// - Apply damage modifiers
	// - Check for status effect application
	// - Handle combat events

	result := CombatResult{
		Hit:           true,
		Damage:        c.calculateBaseDamage(),
		Critical:      false,
		Killed:        false,
		StatusApplied: []string{},
	}

	return result, nil
}

func (c *BaseCombatant) Defend(ctx context.Context, attack AttackInfo) (DefenseResult, error) {
	if !c.IsAlive() {
		return DefenseResult{}, fmt.Errorf("dead entities cannot defend")
	}

	// TODO: Implement defense calculation
	// - Check for block
	// - Check for evasion
	// - Calculate damage mitigation from armor/resistances
	// - Apply final damage
	// - Check for status effect application
	// - Handle defense events

	result := DefenseResult{
		Blocked:       false,
		Evaded:        false,
		Mitigated:     0,
		FinalDamage:   attack.BaseDamage,
		StatusApplied: []string{},
	}

	// Apply damage to self
	_, err := c.Damage(ctx, result.FinalDamage, attack.AttackerID)
	if err != nil {
		return result, err
	}

	result.FinalDamage = result.FinalDamage - result.Mitigated

	return result, nil
}

func (c *BaseCombatant) CanAttack(targetID string) bool {
	if !c.IsAlive() {
		return false
	}

	if !c.CanAct() {
		return false
	}

	// TODO: Check range, line of sight, cooldowns, resources, etc.

	return true
}

func (c *BaseCombatant) AttackRange() float64 {
	return c.attackRange
}

func (c *BaseCombatant) ThreatLevel() float64 {
	return c.threatLevel
}

func (c *BaseCombatant) ModifyThreat(delta float64) {
	c.threatLevel += delta
	if c.threatLevel < 0 {
		c.threatLevel = 0
	}
	c.Touch()
}

// Helper methods

func (c *BaseCombatant) calculateBaseDamage() float64 {
	// TODO: Implement proper damage calculation
	// Base damage = weapon damage + (strength * multiplier) + flat damage bonuses

	strength := c.Attributes().Get(attribute.AttrStrength)
	physicalDamage := c.Attributes().Get(attribute.AttrPhysicalDamage)

	baseDamage := physicalDamage + (strength * 0.5)

	return baseDamage
}

func (c *BaseCombatant) calculateArmor() float64 {
	// TODO: Implement armor calculation
	return c.Attributes().Get(attribute.AttrArmor)
}

func (c *BaseCombatant) calculateEvasion() float64 {
	// TODO: Implement evasion calculation
	dexterity := c.Attributes().Get(attribute.AttrDexterity)
	evasion := c.Attributes().Get(attribute.AttrEvasion)

	return evasion + (dexterity * 0.1)
}

func (c *BaseCombatant) rollCritical() bool {
	// TODO: Implement critical hit roll
	critChance := c.Attributes().Get(attribute.AttrCritChance)
	_ = critChance

	// For now, return false
	return false
}

func (c *BaseCombatant) rollHit(targetEvasion float64) bool {
	// TODO: Implement hit calculation
	accuracy := c.Attributes().Get(attribute.AttrAccuracy)
	_ = accuracy
	_ = targetEvasion

	// For now, always hit
	return true
}

func (c *BaseCombatant) SerializeState() (map[string]any, error) {
	state, err := c.BaseLiving.SerializeState()
	if err != nil {
		return nil, err
	}

	state["threat_level"] = c.threatLevel
	state["attack_range"] = c.attackRange

	return state, nil
}

func (c *BaseCombatant) DeserializeState(state map[string]any) error {
	if err := c.BaseLiving.DeserializeState(state); err != nil {
		return err
	}

	if threat, ok := state["threat_level"].(float64); ok {
		c.threatLevel = threat
	}

	if attackRange, ok := state["attack_range"].(float64); ok {
		c.attackRange = attackRange
	}

	return nil
}
