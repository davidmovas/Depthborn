package affix

import "github.com/davidmovas/Depthborn/internal/core/attribute"

// Type categorizes affixes
type Type string

const (
	TypePrefix    Type = "prefix"
	TypeSuffix    Type = "suffix"
	TypeImplicit  Type = "implicit"
	TypeCorrupted Type = "corrupted"
	TypeEnchant   Type = "enchant"
)

// Affix represents a template for item modifiers.
// This is the "blueprint" loaded from YAML - immutable definition.
type Affix interface {
	// ID returns unique affix identifier
	ID() string

	// Name returns display name
	Name() string

	// Type returns affix type (prefix, suffix, etc.)
	Type() Type

	// Group returns mutual exclusion group.
	// Two affixes with the same group cannot be on the same item.
	Group() string

	// Rank returns internal power rank [1-100].
	// Higher rank = stronger affix, affects generation weights.
	// Hidden from player.
	Rank() int

	// Modifiers returns modifier templates for this affix.
	// Each template has min/max values that are rolled during generation.
	Modifiers() []ModifierTemplate

	// Requirements returns item requirements to roll this affix
	Requirements() Requirements

	// BaseWeight returns spawn weight for random generation.
	// Actual weight is calculated based on item level, rarity, etc.
	BaseWeight() int

	// Description returns human-readable description template.
	// May contain placeholders like "{value}" for rolled values.
	Description() string

	// Tags returns affix tags for filtering and modification
	Tags() []string

	// HasTag checks if affix has specific tag
	HasTag(tag string) bool
}

// ModifierTemplate defines a range for modifier values
type ModifierTemplate struct {
	// Attribute being modified
	Attribute attribute.Type

	// ModType defines how modifier is applied (flat, increased, more)
	ModType attribute.ModifierType

	// MinValue is minimum possible value
	MinValue float64

	// MaxValue is maximum possible value
	MaxValue float64

	// Priority for application order
	Priority int
}

// Requirements defines conditions for affix to appear
type Requirements interface {
	// MinItemLevel returns minimum item level required
	MinItemLevel() int

	// MaxItemLevel returns maximum item level allowed (0 = no limit)
	MaxItemLevel() int

	// AllowedTypes returns item types that can have this affix
	AllowedTypes() []string

	// AllowedSlots returns equipment slots that can have this affix
	AllowedSlots() []string

	// Check verifies if item can have this affix
	Check(itemType string, itemLevel int, slot string) bool
}

// Instance represents a rolled affix on an actual item.
// Contains concrete values generated from Affix template.
type Instance interface {
	// AffixID returns ID of source affix template
	AffixID() string

	// Affix returns source affix template (may be nil if not loaded)
	Affix() Affix

	// Type returns affix type
	Type() Type

	// Group returns mutual exclusion group
	Group() string

	// RolledValues returns rolled values for each modifier
	RolledValues() []RolledModifier

	// Modifiers returns actual attribute modifiers with rolled values
	Modifiers() []attribute.Modifier

	// Reroll re-rolls all values within original ranges
	Reroll()

	// RerollSingle re-rolls single modifier at index
	RerollSingle(index int) error

	// Quality returns how good the roll is [0.0 - 1.0]
	// 0.0 = all minimum values, 1.0 = all maximum values
	Quality() float64
}

// RolledModifier contains a rolled value and its range
type RolledModifier struct {
	Template ModifierTemplate
	Value    float64
}

// Set manages affixes on an item
type Set interface {
	// Add adds affix instance to set
	Add(instance Instance) error

	// Remove removes affix by ID
	Remove(affixID string) error

	// Get retrieves affix instance by ID
	Get(affixID string) (Instance, bool)

	// GetByType returns all instances of specified type
	GetByType(affixType Type) []Instance

	// GetAll returns all affix instances
	GetAll() []Instance

	// Count returns total number of affixes
	Count() int

	// CountByType returns number of affixes of specified type
	CountByType(affixType Type) int

	// CanAdd checks if affix can be added (limits and groups)
	CanAdd(instance Instance) bool

	// HasGroup checks if set has affix from specified group
	HasGroup(group string) bool

	// PrefixCount returns current prefix count
	PrefixCount() int

	// SuffixCount returns current suffix count
	SuffixCount() int

	// MaxPrefixes returns maximum allowed prefixes
	MaxPrefixes() int

	// MaxSuffixes returns maximum allowed suffixes
	MaxSuffixes() int

	// SetLimits sets prefix and suffix limits
	SetLimits(minPrefix, maxPrefix, minSuffix, maxSuffix int)

	// Clear removes all affixes
	Clear()

	// AllModifiers returns combined modifiers from all affixes
	AllModifiers() []attribute.Modifier

	// RerollAll re-rolls values on all affixes
	RerollAll()

	// TotalQuality returns average quality across all affixes
	TotalQuality() float64
}

// Pool represents collection of affixes available for rolling
type Pool interface {
	// Add adds affix to pool
	Add(affix Affix)

	// Remove removes affix from pool
	Remove(affixID string)

	// Get retrieves affix by ID
	Get(affixID string) (Affix, bool)

	// GetAll returns all affixes in pool
	GetAll() []Affix

	// GetByGroup returns affixes in specified group
	GetByGroup(group string) []Affix

	// GetByTags returns affixes matching all specified tags
	GetByTags(tags ...string) []Affix

	// GetByAnyTag returns affixes matching any of specified tags
	GetByAnyTag(tags ...string) []Affix

	// Filter returns affixes matching criteria
	Filter(criteria FilterCriteria) []Affix

	// Roll randomly selects affix from pool based on weights
	Roll(ctx RollContext) (Affix, error)
}

// RollContext provides context for affix generation
type RollContext struct {
	ItemType   string
	ItemLevel  int
	ItemSlot   string
	ItemRarity int // Rarity affects weight calculations

	// ExcludeGroups - groups to exclude (already on item)
	ExcludeGroups []string

	// ExcludeIDs - specific affix IDs to exclude
	ExcludeIDs []string

	// RequireTags - only consider affixes with these tags
	RequireTags []string

	// ExcludeTags - exclude affixes with these tags
	ExcludeTags []string

	// AffixType - filter by type (prefix/suffix/etc)
	AffixType *Type
}

// FilterCriteria defines filtering for affix selection
type FilterCriteria struct {
	Types        []Type
	Groups       []string
	Tags         []string
	MinRank      int
	MaxRank      int
	MinItemLevel int
	MaxItemLevel int
}

// Generator creates affix instances for items
type Generator interface {
	// Generate creates random affixes for item based on rarity
	Generate(ctx GenerateContext) ([]Instance, error)

	// AddAffix adds single random affix to existing set
	AddAffix(set Set, ctx RollContext) (Instance, error)

	// CreateInstance creates instance from affix template
	CreateInstance(affix Affix) Instance

	// RollValues rolls random values for modifier templates
	RollValues(templates []ModifierTemplate) []RolledModifier
}

// GenerateContext provides context for full affix generation
type GenerateContext struct {
	RollContext

	// PrefixRange - min/max prefixes to generate
	PrefixRange [2]int

	// SuffixRange - min/max suffixes to generate
	SuffixRange [2]int

	// QualityBias affects value distribution [0.0 - 1.0]
	// 0.0 = bias toward minimum, 0.5 = uniform, 1.0 = bias toward maximum
	QualityBias float64
}

// Registry manages all available affixes loaded from data files
type Registry interface {
	// Register adds affix to registry
	Register(affix Affix) error

	// Get retrieves affix by ID
	Get(id string) (Affix, bool)

	// GetAll returns all registered affixes
	GetAll() []Affix

	// GetPool returns pool for specific item type/slot
	GetPool(itemType string, slot string) Pool

	// LoadFromYAML loads affixes from YAML file
	LoadFromYAML(path string) error

	// LoadFromDirectory loads all YAML files from directory
	LoadFromDirectory(path string) error
}
