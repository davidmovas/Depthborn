package entity

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/infra"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

// Entity represents any game object in the world
type Entity interface {
	infra.Persistent
	types.Identity
	types.Named
	types.Leveled
	types.Tagged
	types.Alive
	types.Actionable
	types.Cloneable
	types.Validatable

	// Attributes returns attribute manager
	Attributes() AttributeManager

	// StatusEffects returns status effect manager
	StatusEffects() StatusManager

	// Transform returns position and orientation
	Transform() spatial.Transform

	// Callbacks returns callback registry
	Callbacks() types.CallbackRegistry
}

// AttributeManager manages entity attributes (forward declaration)
type AttributeManager interface {
	// Get returns current attribute value
	Get(attrType attribute.Type) float64

	// GetBase returns base attribute value
	GetBase(attrType attribute.Type) float64

	// SetBase sets base attribute value
	SetBase(attrType attribute.Type, value float64)

	// AddModifier adds attribute modifier
	AddModifier(attrType attribute.Type, modifier AttributeModifier)

	// RemoveModifier removes attribute modifier
	RemoveModifier(attrType attribute.Type, modifierID string)

	// RecalculateAll recalculates all attributes
	RecalculateAll()
}

// AttributeModifier modifies attribute values (forward declaration)
type AttributeModifier interface {
	// ID returns unique modifier identifier
	ID() string

	// Value returns modifier value
	Value() float64

	// Type returns modifier type (flat, increased, more)
	Type() string

	// Source returns what created this modifier
	Source() string
}

// StatusManager manages status effects (forward declaration)
type StatusManager interface {
	// Apply adds status effect
	Apply(ctx context.Context, effectID string, sourceID string) error

	// Remove removes status effect
	Remove(ctx context.Context, effectID string) error

	// Has checks if status effect is active
	Has(effectType string) bool

	// GetAll returns all active effects
	GetAll() []string

	// Update processes all effects
	Update(ctx context.Context, deltaMs int64) error
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
