package entity

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/status"
	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/infra"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

// Entity represents any game object in the world
type Entity interface {
	infra.Persistent
	types.Named
	types.Leveled
	types.Tagged
	types.Alive
	types.Actionable
	infra.Cloneable
	infra.Validatable

	// Attributes returns attribute manager
	Attributes() attribute.Manager

	// StatusEffects returns status effect manager
	StatusEffects() status.Manager

	// Transform returns position and orientation
	Transform() spatial.Transform

	// Callbacks returns callback registry
	Callbacks() types.CallbackRegistry
}

// Living represents entities with health
type Living interface {
	Entity

	// Health returns current health
	Health() float64

	// MaxHealth returns maximum health
	MaxHealth() float64

	// SetHealth updates current health
	SetHealth(value float64)

	// Damage reduces health, returns actual damage dealt
	Damage(ctx context.Context, amount float64, sourceID string) (float64, error)

	// Heal increases health, returns actual healing done
	Heal(ctx context.Context, amount float64, sourceID string) (float64, error)

	// HealthPercent returns health as percentage [0.0 - 1.0]
	HealthPercent() float64
}

// Combatant represents entities that can fight
type Combatant interface {
	Living

	// Attack performs attack against target
	Attack(ctx context.Context, targetID string) (CombatResult, error)

	// Defend calculates defense against attack
	Defend(ctx context.Context, attack AttackInfo) (DefenseResult, error)

	// CanAttack checks if can attack target
	CanAttack(targetID string) bool

	// AttackRange returns maximum attack distance
	AttackRange() float64

	// ThreatLevel returns aggression level
	ThreatLevel() float64

	// ModifyThreat adjusts threat level
	ModifyThreat(delta float64)
}

// AttackInfo describes incoming attack
type AttackInfo struct {
	AttackerID   string
	BaseDamage   float64
	DamageType   string
	IsCritical   bool
	Penetration  float64
	StatusChance map[string]float64
}

// DefenseResult describes defense outcome
type DefenseResult struct {
	Blocked       bool
	Evaded        bool
	Mitigated     float64
	FinalDamage   float64
	StatusApplied []string
}

// CombatResult describes combat outcome
type CombatResult struct {
	Hit           bool
	Damage        float64
	Critical      bool
	Killed        bool
	StatusApplied []string
}
