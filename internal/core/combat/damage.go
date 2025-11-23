package combat

import (
	"context"
)

// DamageCalculator computes damage and defense
type DamageCalculator interface {
	// CalculateDamage computes final damage dealt
	CalculateDamage(ctx context.Context, attack AttackData, defense DefenseData) (DamageResult, error)

	// CalculateHealing computes final healing done
	CalculateHealing(ctx context.Context, heal HealData) (HealResult, error)

	// CalculateCritical determines if attack is critical
	CalculateCritical(ctx context.Context, attacker AttackerData) (bool, float64)

	// CalculateAccuracy determines if attack hits
	CalculateAccuracy(ctx context.Context, attacker AttackerData, defender DefenderData) bool

	// CalculateDamageReduction computes damage mitigation
	CalculateDamageReduction(ctx context.Context, damage float64, defense DefenseData) float64

	// CalculatePenetration applies armor penetration
	CalculatePenetration(ctx context.Context, damage float64, penetration float64, armor float64) float64

	// ApplyResistances applies elemental resistances
	ApplyResistances(ctx context.Context, damage float64, damageType DamageType, resistances map[DamageType]float64) float64

	// ApplyVulnerabilities applies damage vulnerabilities
	ApplyVulnerabilities(ctx context.Context, damage float64, damageType DamageType, vulnerabilities map[DamageType]float64) float64

	// CalculateFinalDamage computes final damage after all modifiers
	CalculateFinalDamage(ctx context.Context, baseDamage float64, modifiers []DamageModifier) float64
}

// AttackData contains attacker information
type AttackData struct {
	AttackerID       string
	BaseDamage       float64
	DamageType       DamageType
	SecondaryTypes   []DamageType
	CritChance       float64
	CritMultiplier   float64
	Accuracy         float64
	Penetration      float64
	LifeSteal        float64
	StatusChance     map[string]float64
	DamageModifiers  []DamageModifier
	SkillPower       float64
	IsSkill          bool
	SkillID          string
	WeaponDamage     float64
	AttributeScaling map[string]float64
	BonusDamageFlat  float64
	Flags            []AttackFlag
}

// AttackFlag describes special attack property
type AttackFlag string

const (
	AttackFlagGuaranteedHit    AttackFlag = "guaranteed_hit"
	AttackFlagGuaranteedCrit   AttackFlag = "guaranteed_crit"
	AttackFlagIgnoreArmor      AttackFlag = "ignore_armor"
	AttackFlagIgnoreResistance AttackFlag = "ignore_resistance"
	AttackFlagCannotBeDodged   AttackFlag = "cannot_be_dodged"
	AttackFlagCannotBeBlocked  AttackFlag = "cannot_be_blocked"
	AttackFlagPiercing         AttackFlag = "piercing"
	AttackFlagExecute          AttackFlag = "execute"
)

// DefenseData contains defender information
type DefenseData struct {
	DefenderID      string
	Armor           float64
	Evasion         float64
	BlockChance     float64
	BlockAmount     float64
	Resistances     map[DamageType]float64
	Vulnerabilities map[DamageType]float64
	DamageReduction float64
	DamageModifiers []DamageModifier
	IsDefending     bool
	CoverType       CoverType
	HasShield       bool
	ShieldAmount    float64
	Immunities      []DamageType
	Flags           []DefenseFlag
}

// DefenseFlag describes special defense property
type DefenseFlag string

const (
	DefenseFlagInvulnerable   DefenseFlag = "invulnerable"
	DefenseFlagPhysicalImmune DefenseFlag = "physical_immune"
	DefenseFlagMagicImmune    DefenseFlag = "magic_immune"
	DefenseFlagReflectDamage  DefenseFlag = "reflect_damage"
	DefenseFlagAbsorbDamage   DefenseFlag = "absorb_damage"
	DefenseFlagThorns         DefenseFlag = "thorns"
)

// AttackerData contains attacker stats for hit calculation
type AttackerData struct {
	EntityID       string
	CritChance     float64
	CritMultiplier float64
	Accuracy       float64
	Level          int
	LuckModifier   float64
}

// DefenderData contains defender stats for hit calculation
type DefenderData struct {
	EntityID     string
	Evasion      float64
	Level        int
	IsDefending  bool
	CoverType    CoverType
	LuckModifier float64
}

// HealData contains healing information
type HealData struct {
	HealerID       string
	TargetID       string
	BaseHealing    float64
	HealingPower   float64
	HealingBonus   float64
	HealingType    HealingType
	CanOverheal    bool
	RemovesDebuffs []string
	Modifiers      []DamageModifier
	IsCritical     bool
	CritMultiplier float64
}

// HealingType categorizes healing
type HealingType string

const (
	HealingInstant    HealingType = "instant"
	HealingOverTime   HealingType = "over_time"
	HealingPercentage HealingType = "percentage"
	HealingShield     HealingType = "shield"
	HealingAbsorb     HealingType = "absorb"
)

// DamageResult describes damage outcome
type DamageResult struct {
	TotalDamage          float64
	DamageByType         map[DamageType]float64
	Critical             bool
	CritMultiplier       float64
	Hit                  bool
	Blocked              bool
	Evaded               bool
	Parried              bool
	Resisted             float64
	Penetrated           float64
	Mitigated            float64
	Absorbed             float64
	Reflected            float64
	Overkill             float64
	LifeStolen           float64
	StatusApplied        map[string]bool
	Flags                []DamageResultFlag
	PreMitigationDamage  float64
	PostMitigationDamage float64
	ShieldDamage         float64
	HealthDamage         float64
}

// DamageResultFlag describes special damage result
type DamageResultFlag string

const (
	DamageFlagCritical    DamageResultFlag = "critical"
	DamageFlagBlocked     DamageResultFlag = "blocked"
	DamageFlagEvaded      DamageResultFlag = "evaded"
	DamageFlagParried     DamageResultFlag = "parried"
	DamageFlagPenetrating DamageResultFlag = "penetrating"
	DamageFlagReflected   DamageResultFlag = "reflected"
	DamageFlagAbsorbed    DamageResultFlag = "absorbed"
	DamageFlagLifesteal   DamageResultFlag = "lifesteal"
	DamageFlagOverkill    DamageResultFlag = "overkill"
	DamageFlagExecute     DamageResultFlag = "execute"
	DamageFlagDivine      DamageResultFlag = "divine"
	DamageFlagGlancing    DamageResultFlag = "glancing"
	DamageFlagImmune      DamageResultFlag = "immune"
)

// HealResult describes healing outcome
type HealResult struct {
	HealingDone        float64
	Overhealing        float64
	Critical           bool
	CritMultiplier     float64
	DebuffsRemoved     []string
	ShieldGranted      float64
	AbsorbGranted      float64
	TargetAtFullHealth bool
}

// DamageModifier alters damage calculation
type DamageModifier interface {
	// ID returns unique modifier identifier
	ID() string

	// Name returns display name
	Name() string

	// Type returns modifier operation type
	Type() DamageModifierType

	// Value returns modifier value
	Value() float64

	// Apply applies modifier to damage
	Apply(damage float64) float64

	// Condition returns when modifier applies
	Condition() DamageModifierCondition

	// Source returns modifier source
	Source() string

	// Priority returns application order (higher = first)
	Priority() int

	// AffectsDamageType checks if modifier affects damage type
	AffectsDamageType(damageType DamageType) bool
}

// DamageModifierType categorizes damage modifiers
type DamageModifierType string

const (
	DamageModFlat       DamageModifierType = "flat"       // Adds flat amount
	DamageModIncreased  DamageModifierType = "increased"  // Additive percentage
	DamageModMore       DamageModifierType = "more"       // Multiplicative
	DamageModReduction  DamageModifierType = "reduction"  // Subtractive
	DamageModMultiplier DamageModifierType = "multiplier" // Final multiplier
	DamageModConversion DamageModifierType = "conversion" // Type conversion
)

// DamageModifierCondition determines when modifier applies
type DamageModifierCondition interface {
	// Check evaluates if condition is met
	Check(ctx context.Context, attack AttackData, defense DefenseData) bool

	// Description returns human-readable condition
	Description() string
}

// DamageInteraction handles damage type interactions
type DamageInteraction interface {
	// GetInteraction returns interaction between damage types
	GetInteraction(primary, secondary DamageType) InteractionType

	// ApplyInteraction applies interaction effect
	ApplyInteraction(ctx context.Context, targetID string, interaction InteractionType, encounter Encounter) error

	// CalculateCombo calculates bonus damage from combo
	CalculateCombo(types []DamageType) float64

	// HasSynergy checks if damage types have synergy
	HasSynergy(type1, type2 DamageType) bool

	// HasAntiSynergy checks if damage types conflict
	HasAntiSynergy(type1, type2 DamageType) bool
}

const (
	InteractionNone       InteractionType = "none"
	InteractionAmplify    InteractionType = "amplify"    // Increases damage
	InteractionNeutralize InteractionType = "neutralize" // Reduces damage
	InteractionTransform  InteractionType = "transform"  // Changes damage type
	InteractionExplode    InteractionType = "explode"    // AoE burst
	InteractionFreeze     InteractionType = "freeze"     // Stun effect
	InteractionIgnite     InteractionType = "ignite"     // DoT effect
	InteractionShock      InteractionType = "shock"      // Chain effect
	InteractionCorrode    InteractionType = "corrode"    // Armor reduction
	InteractionPurify     InteractionType = "purify"     // Remove effects
	InteractionShatter    InteractionType = "shatter"    // Extra damage
	InteractionCombust    InteractionType = "combust"    // Burn spread
	InteractionElectrify  InteractionType = "electrify"  // Stun and damage
	InteractionChill      InteractionType = "chill"      // Slow effect
)

// DamageOverTime represents DoT/HoT effect
type DamageOverTime interface {
	// ID returns unique DoT identifier
	ID() string

	// Name returns display name
	Name() string

	// SourceID returns damage source entity ID
	SourceID() string

	// TargetID returns affected entity ID
	TargetID() string

	// Type returns damage or healing type
	Type() DamageType

	// IsHealing returns true if healing over time
	IsHealing() bool

	// DamagePerTick returns damage dealt per tick
	DamagePerTick() float64

	// TickInterval returns milliseconds between ticks
	TickInterval() int64

	// Duration returns total duration in milliseconds
	Duration() int64

	// RemainingDuration returns time left
	RemainingDuration() int64

	// SetRemainingDuration updates remaining duration
	SetRemainingDuration(ms int64)

	// Stacks returns current stack count
	Stacks() int

	// MaxStacks returns maximum stacks
	MaxStacks() int

	// AddStack increases stack count
	AddStack() bool

	// RemoveStack decreases stack count
	RemoveStack() bool

	// SetStacks updates stack count
	SetStacks(stacks int)

	// Tick applies damage for elapsed time
	Tick(ctx context.Context, deltaMs int64, encounter Encounter) (float64, error)

	// IsExpired returns true if duration ended
	IsExpired() bool

	// Cancel removes DoT before expiration
	Cancel()

	// IsCancelled returns true if DoT was cancelled
	IsCancelled() bool

	// CanRefresh returns true if DoT can be refreshed
	CanRefresh() bool

	// Refresh resets duration to maximum
	Refresh()

	// Icon returns icon identifier
	Icon() string

	// VisualEffect returns visual effect identifier
	VisualEffect() string
}

// DamageReflection handles damage reflection
type DamageReflection interface {
	// ReflectDamage calculates reflected damage amount
	ReflectDamage(damage float64, reflectionPercent float64) float64

	// CanReflect checks if damage can be reflected
	CanReflect(damageType DamageType) bool

	// ApplyReflection applies reflected damage to attacker
	ApplyReflection(ctx context.Context, attackerID, defenderID string, damage float64, encounter Encounter) error

	// ReflectionModifier returns reflection modifier
	ReflectionModifier(damageType DamageType) float64

	// IsActiveReflection returns true if reflection is active
	IsActiveReflection() bool
}

// DamageAbsorption handles damage absorption/shields
type DamageAbsorption interface {
	// ID returns unique absorption identifier
	ID() string

	// Name returns display name
	Name() string

	// OwnerID returns entity with absorption
	OwnerID() string

	// Amount returns remaining absorption
	Amount() float64

	// MaxAmount returns maximum absorption
	MaxAmount() float64

	// SetAmount updates absorption amount
	SetAmount(amount float64)

	// Absorb reduces damage and depletes shield
	Absorb(damage float64) (absorbed, overflow float64)

	// Restore replenishes absorption
	Restore(amount float64)

	// IsActive returns true if absorption remains
	IsActive() bool

	// IsBroken returns true if fully depleted
	IsBroken() bool

	// Break instantly depletes absorption
	Break()

	// Duration returns remaining duration in milliseconds (-1 = permanent)
	Duration() int64

	// SetDuration updates duration
	SetDuration(ms int64)

	// IsExpired returns true if duration ended
	IsExpired() bool

	// AbsorbsType checks if absorbs specific damage type
	AbsorbsType(damageType DamageType) bool

	// AbsorbEfficiency returns absorption efficiency for damage type
	AbsorbEfficiency(damageType DamageType) float64

	// OnBreak is called when shield breaks
	OnBreak(ctx context.Context, encounter Encounter) error

	// Icon returns icon identifier
	Icon() string

	// VisualEffect returns visual effect identifier
	VisualEffect() string
}

// DamageImmunity handles damage immunity
type DamageImmunity interface {
	// ID returns unique immunity identifier
	ID() string

	// Name returns display name
	Name() string

	// OwnerID returns entity with immunity
	OwnerID() string

	// IsImmune checks if immune to damage type
	IsImmune(damageType DamageType) bool

	// AddImmunity grants immunity to type
	AddImmunity(damageType DamageType)

	// RemoveImmunity removes immunity to type
	RemoveImmunity(damageType DamageType)

	// GetImmunities returns all immunities
	GetImmunities() []DamageType

	// IsCompleteImmunity returns true if immune to all damage
	IsCompleteImmunity() bool

	// Duration returns remaining duration in milliseconds (-1 = permanent)
	Duration() int64

	// SetDuration updates duration
	SetDuration(ms int64)

	// IsExpired returns true if duration ended
	IsExpired() bool

	// OnImmuneHit is called when immune entity is hit
	OnImmuneHit(ctx context.Context, attack AttackData, encounter Encounter) error
}

// DamageInstance represents damage event
type DamageInstance interface {
	// ID returns unique instance identifier
	ID() string

	// SourceID returns damage source entity ID
	SourceID() string

	// TargetID returns damage target entity ID
	TargetID() string

	// Amount returns damage amount
	Amount() float64

	// Type returns damage type
	Type() DamageType

	// IsCritical returns true if critical hit
	IsCritical() bool

	// Flags returns damage flags
	Flags() []DamageResultFlag

	// Timestamp returns when damage occurred
	Timestamp() int64

	// Modifiers returns applied modifiers
	Modifiers() []DamageModifier

	// WasLethal returns true if damage killed target
	WasLethal() bool

	// SetLethal marks damage as lethal
	SetLethal(lethal bool)

	// IsDoT returns true if damage over time
	IsDoT() bool

	// IsReflected returns true if reflected damage
	IsReflected() bool

	// SkillID returns skill that caused damage (empty if basic attack)
	SkillID() string
}

// DamageCalculatorFactory creates damage calculator instances
type DamageCalculatorFactory interface {
	// Create creates new damage calculator
	Create() DamageCalculator

	// CreateWithConfig creates calculator with specific configuration
	CreateWithConfig(config CalculatorConfig) DamageCalculator
}

// CalculatorConfig defines calculator configuration
type CalculatorConfig struct {
	UseCriticalHits       bool
	UseArmorFormula       ArmorFormula
	UseDiminishingReturns bool
	MaxResistance         float64
	MinDamage             float64
	LevelScaling          bool
	LevelDifferenceScale  float64
}

// ArmorFormula defines armor damage reduction formula
type ArmorFormula string

const (
	ArmorFormulaLinear      ArmorFormula = "linear"
	ArmorFormulaDiminishing ArmorFormula = "diminishing"
	ArmorFormulaPercentage  ArmorFormula = "percentage"
	ArmorFormulaThreshold   ArmorFormula = "threshold"
)
