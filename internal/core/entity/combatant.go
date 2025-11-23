package entity

import "context"

type CombatantEntity struct {
	*LivingEntity
	attackRange float64
	threatLevel float64
}

func NewCombatantEntity(id string, name string) *CombatantEntity {
	living := NewLivingEntity(id, name)
	return &CombatantEntity{
		LivingEntity: living,
		attackRange:  2.0,
		threatLevel:  1.0,
	}
}

func (ce *CombatantEntity) Attack(ctx context.Context, targetID string) (CombatResult, error) {
	result := CombatResult{
		Hit:      true,
		Damage:   10.0, // TODO: Calculate based on attributes
		Critical: false,
		Killed:   false,
	}

	// TODO: Implement actual attack logic
	// - Calculate hit chance
	// - Calculate damage
	// - Apply status effects

	return result, nil
}

func (ce *CombatantEntity) Defend(ctx context.Context, attack AttackInfo) (DefenseResult, error) {
	result := DefenseResult{
		Blocked:       false,
		Evaded:        false,
		Mitigated:     0.0,
		FinalDamage:   attack.BaseDamage,
		StatusApplied: make([]string, 0),
	}

	// TODO: Implement defense logic
	// - Check evasion
	// - Check block
	// - Calculate damage mitigation
	// - Apply status effects

	return result, nil
}

func (ce *CombatantEntity) CanAttack(targetID string) bool {
	// TODO: Implement faction checks, range checks, etc.
	return ce.IsAlive()
}

func (ce *CombatantEntity) AttackRange() float64 {
	return ce.attackRange
}

func (ce *CombatantEntity) ThreatLevel() float64 {
	return ce.threatLevel
}

func (ce *CombatantEntity) ModifyThreat(delta float64) {
	ce.threatLevel += delta
	if ce.threatLevel < 0 {
		ce.threatLevel = 0
	}
}
