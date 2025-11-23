package attribute

// Type defines type of attribute
type Type string

const (
	// Primary attributes

	AttrStrength     Type = "strength"
	AttrDexterity    Type = "dexterity"
	AttrIntelligence Type = "intelligence"
	AttrVitality     Type = "vitality"
	AttrWillpower    Type = "willpower"

	// Offensive attributes

	AttrPhysicalDamage Type = "physical_damage"
	AttrMagicalDamage  Type = "magical_damage"
	AttrCritChance     Type = "crit_chance"
	AttrCritMultiplier Type = "crit_multiplier"
	AttrAttackSpeed    Type = "attack_speed"
	AttrAccuracy       Type = "accuracy"

	// Defensive attributes

	AttrArmor           Type = "armor"
	AttrEvasion         Type = "evasion"
	AttrBlockChance     Type = "block_chance"
	AttrBlockAmount     Type = "block_amount"
	AttrPhysicalResist  Type = "physical_resist"
	AttrFireResist      Type = "fire_resist"
	AttrColdResist      Type = "cold_resist"
	AttrLightningResist Type = "lightning_resist"
	AttrPoisonResist    Type = "poison_resist"

	// Utility attributes

	AttrMovementSpeed  Type = "movement_speed"
	AttrLifeRegen      Type = "life_regen"
	AttrManaRegen      Type = "mana_regen"
	AttrLifeSteal      Type = "life_steal"
	AttrLootQuantity   Type = "loot_quantity"
	AttrLootRarity     Type = "loot_rarity"
	AttrExperienceGain Type = "experience_gain"
)

// Manager manages all attributes for an entity
type Manager interface {
	// Get returns current value of attribute including all modifiers
	Get(attr Type) float64

	// GetBase returns base value without modifiers
	GetBase(attr Type) float64

	// SetBase sets base value of attribute
	SetBase(attr Type, value float64)

	// AddModifier adds modifier to attribute
	AddModifier(attr Type, modifier Modifier)

	// RemoveModifier removes specific modifier
	RemoveModifier(attr Type, modifierID string)

	// RemoveAllModifiers removes all modifiers of specified type
	RemoveAllModifiers(attr Type, modType ModifierType)

	// GetModifiers returns all modifiers for attribute
	GetModifiers(attr Type) []Modifier

	// RecalculateAll recalculates all derived attributes
	RecalculateAll()

	// Snapshot creates snapshot of all attributes
	Snapshot() map[Type]float64
}

// Modifier changes attribute value
type Modifier interface {
	// ID returns unique modifier identifier
	ID() string

	// Type returns how modifier is applied
	Type() ModifierType

	// Value returns modifier value
	Value() float64

	// Source returns what created this modifier (item, buff, skill, etc)
	Source() string

	// Priority returns application order (higher applies first)
	Priority() int

	// IsActive returns true if modifier should be applied
	IsActive() bool
}

// ModifierType defines how modifier value is applied
type ModifierType string

const (
	// ModFlat adds flat value to base
	ModFlat ModifierType = "flat"

	// ModIncreased adds percentage of base (additive with other increased)
	ModIncreased ModifierType = "increased"

	// ModMore multiplies total (multiplicative, applied after increased)
	ModMore ModifierType = "more"

	// ModOverride replaces value completely
	ModOverride ModifierType = "override"
)

// Formula calculates derived attribute from other attributes
type Formula interface {
	// Calculate computes value based on manager state
	Calculate(manager Manager) float64

	// Dependencies returns attributes this formula depends on
	Dependencies() []Type
}

// Set represents a collection of attribute modifiers
type Set interface {
	// Add adds modifier to set
	Add(modifier Modifier)

	// Remove removes modifier from set
	Remove(modifierID string)

	// GetAll returns all modifiers in set
	GetAll() []Modifier

	// GetByType returns modifiers of specific type
	GetByType(modType ModifierType) []Modifier

	// Clear removes all modifiers
	Clear()

	// Apply applies all modifiers to base value
	Apply(baseValue float64) float64
}
