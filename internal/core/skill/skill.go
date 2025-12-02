package skill

import (
	"context"
)

// =============================================================================
// SKILL TYPES
// =============================================================================

// Type categorizes skills
type Type string

const (
	TypeActive  Type = "active"  // Manually activated skill
	TypePassive Type = "passive" // Always-on effect
	TypeAura    Type = "aura"    // Area effect around caster
	TypeTrigger Type = "trigger" // Activates on condition
)

// =============================================================================
// SKILL DEFINITION (Template from YAML)
// =============================================================================

// Def represents a skill template loaded from YAML.
// This is the blueprint - immutable definition.
type Def interface {
	// ID returns unique skill identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns skill description template.
	// May contain placeholders like "{damage}" for level-specific values.
	Description() string

	// Type returns skill type (active, passive, aura, trigger)
	Type() Type

	// Tags returns skill tags for filtering and modification
	Tags() []string

	// HasTag checks if skill has specific tag
	HasTag(tag string) bool

	// MaxLevel returns maximum skill level (0 = no levels)
	MaxLevel() int

	// LevelData returns data for specific level.
	// Returns nil if level is invalid or skill has no levels.
	LevelData(level int) LevelData

	// BaseCooldown returns base cooldown in milliseconds (0 = no cooldown)
	BaseCooldown() int64

	// BaseCharges returns base number of charges (0 = no charges, uses cooldown)
	BaseCharges() int

	// ChargeRecovery returns time to recover one charge in ms
	ChargeRecovery() int64

	// Targeting returns targeting rules for this skill
	Targeting() TargetRule

	// Effects returns effect templates that this skill applies
	Effects() []EffectDef

	// Requirements returns requirements to learn/use this skill
	Requirements() Requirements

	// Icon returns icon identifier for UI
	Icon() string

	// Metadata returns additional skill-specific data
	Metadata() map[string]any
}

// LevelData contains level-specific skill values.
// All values are pre-defined in YAML for each level.
type LevelData interface {
	// Level returns this data's level
	Level() int

	// ResourceCosts returns resources consumed when using skill at this level
	ResourceCosts() []ResourceCost

	// Effects returns effect values for this level
	Effects() []EffectValue

	// Cooldown returns cooldown override for this level (0 = use base)
	Cooldown() int64

	// Charges returns charges override for this level (0 = use base)
	Charges() int

	// Description returns level-specific description
	Description() string

	// Metadata returns level-specific data
	Metadata() map[string]any
}

// EffectValue contains rolled/calculated values for an effect at specific level
type EffectValue struct {
	EffectID string         // References EffectDef.ID()
	Values   map[string]any // Effect-specific values (damage, healing, duration, etc.)
}

// =============================================================================
// SKILL INSTANCE (Runtime state on character)
// =============================================================================

// Instance represents a skill owned by a character.
// Contains runtime state: current level, cooldown, charges.
type Instance interface {
	// DefID returns source definition ID
	DefID() string

	// Def returns source definition (may be nil if not loaded)
	Def() Def

	// Level returns current skill level (1-based, 0 = not leveled)
	Level() int

	// SetLevel changes skill level
	SetLevel(level int) error

	// CanLevelUp checks if skill can be leveled up
	CanLevelUp() bool

	// LevelUp increases skill level by 1
	LevelUp() error

	// CurrentLevelData returns LevelData for current level
	CurrentLevelData() LevelData

	// Cooldown returns remaining cooldown in milliseconds
	Cooldown() int64

	// SetCooldown updates remaining cooldown
	SetCooldown(ms int64)

	// IsOnCooldown returns true if skill is on cooldown
	IsOnCooldown() bool

	// Charges returns current available charges
	Charges() int

	// MaxCharges returns maximum charges at current level
	MaxCharges() int

	// UseCharge consumes one charge, returns false if no charges
	UseCharge() bool

	// ChargeRecoveryProgress returns progress to next charge [0.0, 1.0]
	ChargeRecoveryProgress() float64

	// Update processes cooldown/charge recovery for elapsed time
	Update(deltaMs int64)

	// CanUse checks if skill can be used (has charges/not on cooldown, resources available)
	CanUse(ctx context.Context, casterID string) bool

	// Use executes the skill
	Use(ctx context.Context, casterID string, params ActivationParams) (Result, error)

	// IsActive returns true if skill is currently active (for toggle/aura)
	IsActive() bool

	// SetActive toggles skill active state
	SetActive(active bool)

	// Modifiers returns active modifiers affecting this skill instance
	Modifiers() []SkillModifier
}

// =============================================================================
// RESOURCES
// =============================================================================

// ResourceType defines what resource is consumed
type ResourceType string

const (
	ResourceMana       ResourceType = "mana"
	ResourceHealth     ResourceType = "health"
	ResourceStamina    ResourceType = "stamina"
	ResourceRage       ResourceType = "rage"
	ResourceEnergy     ResourceType = "energy"
	ResourceSoulCharge ResourceType = "soul_charge"
	ResourceGold       ResourceType = "gold"
)

// CostType defines how cost is calculated
type CostType string

const (
	CostFlat    CostType = "flat"    // Fixed amount
	CostPercent CostType = "percent" // Percentage of max
	CostCurrent CostType = "current" // Percentage of current
)

// ResourceCost defines a resource requirement
type ResourceCost struct {
	Resource ResourceType
	Type     CostType
	Amount   float64
}

// =============================================================================
// TARGETING
// =============================================================================

// TargetType defines what can be targeted
type TargetType string

const (
	TargetNone       TargetType = "none"        // No target (self-only)
	TargetSelf       TargetType = "self"        // Caster only
	TargetSingle     TargetType = "single"      // One target
	TargetMultiple   TargetType = "multiple"    // Multiple targets (up to N)
	TargetAllEnemies TargetType = "all_enemies" // All enemies in range
	TargetAllAllies  TargetType = "all_allies"  // All allies in range
	TargetAll        TargetType = "all"         // Everyone in range
	TargetGround     TargetType = "ground"      // Target position
)

// AreaType defines area of effect shape
type AreaType string

const (
	AreaNone   AreaType = "none"   // Single target
	AreaCircle AreaType = "circle" // Circle around point
	AreaCone   AreaType = "cone"   // Cone in direction
	AreaLine   AreaType = "line"   // Line from caster
	AreaChain  AreaType = "chain"  // Chains between targets
)

// TargetRule defines how skill selects targets
type TargetRule interface {
	// Type returns targeting type
	Type() TargetType

	// AreaType returns area of effect type
	AreaType() AreaType

	// Range returns maximum targeting range (0 = melee/self)
	Range() float64

	// AreaRadius returns AoE radius (for circle, cone angle, chain distance)
	AreaRadius() float64

	// MaxTargets returns maximum targets (0 = unlimited)
	MaxTargets() int

	// MinTargets returns minimum required targets
	MinTargets() int

	// CanTargetSelf returns true if caster can be a target
	CanTargetSelf() bool

	// CanTargetAllies returns true if allies can be targets
	CanTargetAllies() bool

	// CanTargetEnemies returns true if enemies can be targets
	CanTargetEnemies() bool

	// RequiresLineOfSight returns true if needs clear path
	RequiresLineOfSight() bool

	// ChainCount returns number of chain bounces (for chain type)
	ChainCount() int

	// ChainFalloff returns damage reduction per chain [0.0, 1.0]
	ChainFalloff() float64
}

// =============================================================================
// EFFECTS
// =============================================================================

// EffectType categorizes skill effects
type EffectType string

const (
	EffectDamage      EffectType = "damage"       // Deal damage
	EffectHeal        EffectType = "heal"         // Restore health
	EffectStatus      EffectType = "status"       // Apply status effect
	EffectBuff        EffectType = "buff"         // Apply positive modifier
	EffectDebuff      EffectType = "debuff"       // Apply negative modifier
	EffectSummon      EffectType = "summon"       // Create entity
	EffectTeleport    EffectType = "teleport"     // Move instantly
	EffectKnockback   EffectType = "knockback"    // Push target
	EffectPull        EffectType = "pull"         // Pull target
	EffectModifySkill EffectType = "modify_skill" // Modify another skill
	EffectResource    EffectType = "resource"     // Restore/drain resource
	EffectDispel      EffectType = "dispel"       // Remove effects
)

// EffectDef defines what a skill does
type EffectDef interface {
	// ID returns unique effect identifier within skill
	ID() string

	// Type returns effect type
	Type() EffectType

	// DamageType returns damage type (for damage effects)
	DamageType() string

	// StatusID returns status effect to apply (for status effects)
	StatusID() string

	// Scaling returns attribute scaling rules
	Scaling() []ScalingRule

	// Chance returns probability of effect applying [0.0, 1.0] (1.0 = always)
	Chance() float64

	// Delay returns delay before effect activates in ms
	Delay() int64

	// Duration returns effect duration in ms (0 = instant)
	Duration() int64

	// Metadata returns effect-specific parameters
	Metadata() map[string]any
}

// ScalingRule defines how an attribute scales effect value
type ScalingRule struct {
	Attribute  string  // Which attribute affects this (e.g., "strength", "intelligence")
	Multiplier float64 // How much per point (e.g., 0.5 = +0.5 per point)
}

// =============================================================================
// SKILL MODIFIERS
// =============================================================================

// SkillModifier modifies skill behavior
type SkillModifier interface {
	// ID returns unique modifier identifier
	ID() string

	// Source returns what created this modifier (item, passive, buff)
	Source() string

	// Affects returns true if this modifier affects given skill
	Affects(skillID string, skillTags []string) bool

	// ModifyCooldown modifies base cooldown
	ModifyCooldown(baseCooldown int64) int64

	// ModifyCost modifies resource cost
	ModifyCost(cost ResourceCost) ResourceCost

	// ModifyDamage modifies damage value
	ModifyDamage(baseDamage float64, damageType string) float64

	// ModifyEffect modifies effect values
	ModifyEffect(effect EffectValue) EffectValue

	// AddedEffects returns additional effects granted by this modifier
	AddedEffects() []EffectDef

	// ConvertsDamage returns damage type conversion (empty = no conversion)
	ConvertsDamage() (from, to string)

	// Priority returns application order (higher = first)
	Priority() int
}

// =============================================================================
// REQUIREMENTS
// =============================================================================

// Requirements defines conditions to learn/use skill
type Requirements interface {
	// CharacterLevel returns minimum character level
	CharacterLevel() int

	// Attributes returns required attribute values (attribute name -> min value)
	Attributes() map[string]float64

	// Skills returns prerequisite skills (skillID -> minLevel)
	Skills() map[string]int

	// Nodes returns required tree nodes
	Nodes() []string

	// Items returns required equipped item types
	Items() []string

	// Check verifies if entity meets all requirements
	Check(ctx context.Context, entityID string) bool
}

// =============================================================================
// ACTIVATION & RESULT
// =============================================================================

// ActivationParams contains skill activation data
type ActivationParams struct {
	// TargetID is primary target entity ID
	TargetID string

	// TargetIDs is for multi-target skills
	TargetIDs []string

	// TargetPos is for ground-targeted skills
	TargetPos *Position

	// Direction is for directional skills (radians)
	Direction float64

	// Modifiers are runtime modifiers for this activation
	Modifiers map[string]any
}

// Position represents 2D coordinates
type Position struct {
	X float64
	Y float64
}

// Result contains skill execution results
type Result struct {
	Success bool
	Message string

	// TargetsHit contains all affected target IDs
	TargetsHit []string

	// Effects contains results per target
	Effects map[string]TargetResult // targetID -> result

	// TotalDamage is sum of all damage dealt
	TotalDamage float64

	// TotalHealing is sum of all healing done
	TotalHealing float64

	// StatusApplied lists all applied status effect IDs
	StatusApplied []string

	// ResourcesConsumed lists consumed resources
	ResourcesConsumed []ResourceCost

	// Metadata contains additional result data
	Metadata map[string]any
}

// TargetResult contains effect results for single target
type TargetResult struct {
	TargetID      string
	Damage        float64
	DamageType    string
	Healing       float64
	Critical      bool
	Evaded        bool
	Blocked       bool
	StatusApplied []string
	Flags         []string
}

// =============================================================================
// LEGACY COMPATIBILITY - Skill interface (wraps Instance)
// =============================================================================

// Skill is the legacy interface, now wrapping Instance + Def
type Skill interface {
	ID() string
	Name() string
	Description() string
	SkillType() Type
	Level() int
	MaxLevel() int
	CanLevelUp() bool
	LevelUp() error
	Requirements() Requirements
	CanUse(ctx context.Context, casterID string) bool
	Use(ctx context.Context, casterID string, params ActivationParams) (Result, error)
	Cooldown() int64
	MaxCooldown() int64
	SetCooldown(ms int64)
	ManaCost() float64
	Tags() []string
	Metadata() map[string]any
}
