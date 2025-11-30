package entity

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/pkg/persist"
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

	// Calculate base damage from attributes
	baseDamage := c.calculateBaseDamage()

	// Roll for critical hit
	isCritical := c.rollCritical()
	if isCritical {
		critMultiplier := c.Attributes().Get(attribute.AttrCritMultiplier)
		if critMultiplier < 1.0 {
			critMultiplier = 1.5 // Default crit multiplier
		}
		baseDamage *= critMultiplier
	}

	// Build attack info for target's defense calculation
	attackInfo := AttackInfo{
		AttackerID:   c.ID(),
		BaseDamage:   baseDamage,
		DamageType:   "physical",
		IsCritical:   isCritical,
		Penetration:  0, // Could be calculated from attributes
		StatusChance: make(map[string]float64),
	}

	result := CombatResult{
		Hit:           true,
		Damage:        attackInfo.BaseDamage,
		Critical:      isCritical,
		Killed:        false,
		StatusApplied: []string{},
	}

	// Trigger attack callback
	c.Callbacks().TriggerDamage(ctx, targetID, result.Damage, c.ID())

	return result, nil
}

func (c *BaseCombatant) Defend(ctx context.Context, attack AttackInfo) (DefenseResult, error) {
	if !c.IsAlive() {
		return DefenseResult{}, fmt.Errorf("dead entities cannot defend")
	}

	result := DefenseResult{
		Blocked:       false,
		Evaded:        false,
		Mitigated:     0,
		FinalDamage:   attack.BaseDamage,
		StatusApplied: []string{},
	}

	// Check for evasion
	if c.rollEvasion() {
		result.Evaded = true
		result.FinalDamage = 0
		return result, nil
	}

	// Check for block
	if c.rollBlock() {
		result.Blocked = true
		blockAmount := c.Attributes().Get(attribute.AttrBlockAmount)
		if blockAmount <= 0 {
			blockAmount = result.FinalDamage * 0.5 // Default 50% block
		}
		result.Mitigated += blockAmount
		result.FinalDamage = math.Max(0, result.FinalDamage-blockAmount)
	}

	// Calculate armor mitigation (only for physical damage)
	if attack.DamageType == "physical" || attack.DamageType == "" {
		armor := c.calculateArmor()
		// Armor formula: reduction = armor / (armor + 100)
		// This gives diminishing returns
		armorReduction := armor / (armor + 100)
		armorMitigation := result.FinalDamage * armorReduction

		// Apply penetration (reduces armor effectiveness)
		if attack.Penetration > 0 {
			penetrationFactor := 1 - (attack.Penetration / 100)
			if penetrationFactor < 0 {
				penetrationFactor = 0
			}
			armorMitigation *= penetrationFactor
		}

		result.Mitigated += armorMitigation
		result.FinalDamage = math.Max(0, result.FinalDamage-armorMitigation)
	}

	// Apply resistance based on damage type
	resistanceType := c.getResistanceForDamageType(attack.DamageType)
	if resistanceType != "" {
		resistance := c.Attributes().Get(resistanceType)
		// Resistance reduces damage by percentage (capped at 75%)
		resistanceReduction := math.Min(resistance/100, 0.75)
		resistanceMitigation := result.FinalDamage * resistanceReduction
		result.Mitigated += resistanceMitigation
		result.FinalDamage = math.Max(0, result.FinalDamage-resistanceMitigation)
	}

	// Apply final damage to self
	if result.FinalDamage > 0 {
		actualDamage, err := c.Damage(ctx, result.FinalDamage, attack.AttackerID)
		if err != nil {
			return result, err
		}
		result.FinalDamage = actualDamage

		// Check if killed
		if !c.IsAlive() {
			// Entity died from this attack
		}
	}

	// Apply status effects based on attack's status chances
	for effectType, chance := range attack.StatusChance {
		if rand.Float64() < chance {
			result.StatusApplied = append(result.StatusApplied, effectType)
		}
	}

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
	// Base damage formula:
	// Base = physicalDamage + (strength * 2)
	// Strength gives +2 physical damage per point

	strength := c.Attributes().Get(attribute.AttrStrength)
	physicalDamage := c.Attributes().Get(attribute.AttrPhysicalDamage)

	baseDamage := physicalDamage + (strength * 2.0)

	// Minimum damage of 1
	if baseDamage < 1 {
		baseDamage = 1
	}

	return baseDamage
}

func (c *BaseCombatant) calculateArmor() float64 {
	// Armor formula:
	// Total = base armor + (vitality * 1)
	// Vitality gives +1 armor per point

	vitality := c.Attributes().Get(attribute.AttrVitality)
	baseArmor := c.Attributes().Get(attribute.AttrArmor)

	return baseArmor + vitality
}

func (c *BaseCombatant) calculateEvasion() float64 {
	// Evasion formula:
	// Total = base evasion + (dexterity * 2)
	// Dexterity gives +2 evasion rating per point

	dexterity := c.Attributes().Get(attribute.AttrDexterity)
	baseEvasion := c.Attributes().Get(attribute.AttrEvasion)

	return baseEvasion + (dexterity * 2.0)
}

func (c *BaseCombatant) rollCritical() bool {
	// Critical hit chance (percentage based)
	// CritChance is stored as percentage (e.g., 5 = 5%)
	critChance := c.Attributes().Get(attribute.AttrCritChance)

	// Add dexterity bonus: +0.1% per point
	dexterity := c.Attributes().Get(attribute.AttrDexterity)
	critChance += dexterity * 0.1

	// Cap at 75%
	if critChance > 75 {
		critChance = 75
	}

	// Roll for crit
	return rand.Float64()*100 < critChance
}

func (c *BaseCombatant) rollEvasion() bool {
	// Evasion chance formula:
	// Chance = evasion / (evasion + 200) * 100
	// This gives diminishing returns, capped at ~50% with very high evasion

	evasion := c.calculateEvasion()
	evasionChance := (evasion / (evasion + 200)) * 100

	// Cap at 75%
	if evasionChance > 75 {
		evasionChance = 75
	}

	return rand.Float64()*100 < evasionChance
}

func (c *BaseCombatant) rollBlock() bool {
	// Block chance (percentage based)
	blockChance := c.Attributes().Get(attribute.AttrBlockChance)

	// Cap at 75%
	if blockChance > 75 {
		blockChance = 75
	}

	return rand.Float64()*100 < blockChance
}

func (c *BaseCombatant) rollHit(targetEvasion float64) bool {
	// Hit chance formula:
	// Base hit chance = 95%
	// Accuracy increases hit chance
	// Target evasion decreases hit chance

	accuracy := c.Attributes().Get(attribute.AttrAccuracy)

	// Hit chance = base + (accuracy - evasion) / 100
	// Minimum 5%, maximum 100%
	hitChance := 95 + (accuracy-targetEvasion)/100

	if hitChance < 5 {
		hitChance = 5
	}
	if hitChance > 100 {
		hitChance = 100
	}

	return rand.Float64()*100 < hitChance
}

func (c *BaseCombatant) getResistanceForDamageType(damageType string) attribute.Type {
	switch damageType {
	case "fire":
		return attribute.AttrFireResist
	case "cold":
		return attribute.AttrColdResist
	case "lightning":
		return attribute.AttrLightningResist
	case "poison":
		return attribute.AttrPoisonResist
	case "physical":
		return attribute.AttrPhysicalResist
	default:
		return ""
	}
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

func (c *BaseCombatant) Clone() any {
	livingClone := c.BaseLiving.Clone().(*BaseLiving)

	clone := &BaseCombatant{
		BaseLiving:  livingClone,
		threatLevel: c.threatLevel,
		attackRange: c.attackRange,
	}

	return clone
}

func (c *BaseCombatant) Validate() error {
	if err := c.BaseLiving.Validate(); err != nil {
		return err
	}

	if c.attackRange < 0 {
		return fmt.Errorf("attack range cannot be negative")
	}

	if c.threatLevel < 0 {
		return fmt.Errorf("threat level cannot be negative")
	}

	return nil
}

// CombatantState holds the complete serializable state of a BaseCombatant.
type CombatantState struct {
	LivingState
	ThreatLevel float64 `msgpack:"threat_level"`
	AttackRange float64 `msgpack:"attack_range"`
}

// MarshalBinary implements persist.Marshaler for BaseCombatant.
func (c *BaseCombatant) MarshalBinary() ([]byte, error) {
	// First get living state
	livingData, err := c.BaseLiving.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// Decode to LivingState to embed
	var ls LivingState
	if err := persist.DefaultCodec().Decode(livingData, &ls); err != nil {
		return nil, err
	}

	cs := CombatantState{
		LivingState: ls,
		ThreatLevel: c.threatLevel,
		AttackRange: c.attackRange,
	}

	return persist.DefaultCodec().Encode(cs)
}

// UnmarshalBinary implements persist.Unmarshaler for BaseCombatant.
func (c *BaseCombatant) UnmarshalBinary(data []byte) error {
	var cs CombatantState
	if err := persist.DefaultCodec().Decode(data, &cs); err != nil {
		return fmt.Errorf("failed to decode combatant state: %w", err)
	}

	// Encode living state back to bytes for base living
	livingData, err := persist.DefaultCodec().Encode(cs.LivingState)
	if err != nil {
		return err
	}

	// Initialize base living if nil
	if c.BaseLiving == nil {
		c.BaseLiving = &BaseLiving{}
	}

	// Restore base living
	if err = c.BaseLiving.UnmarshalBinary(livingData); err != nil {
		return err
	}

	// Restore combatant-specific fields
	c.threatLevel = cs.ThreatLevel
	c.attackRange = cs.AttackRange

	return nil
}
