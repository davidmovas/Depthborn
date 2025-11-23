package affix

import "github.com/davidmovas/Depthborn/internal/core/attribute"

// Affix represents a modifier on equipment
type Affix interface {
	// ID returns unique affix identifier
	ID() string

	// Name returns display name
	Name() string

	// Type returns affix type
	Type() Type

	// Tier returns affix tier (higher = stronger)
	Tier() int

	// Modifiers returns attribute modifiers granted
	Modifiers() []attribute.Modifier

	// Requirements returns item requirements to roll this affix
	Requirements() Requirements

	// Weight returns spawn weight for random generation
	Weight() int

	// Description returns human-readable description
	Description() string

	// Tags returns affix tags for filtering
	Tags() []string
}

// Type categorizes affixes
type Type string

const (
	TypePrefix    Type = "prefix"
	TypeSuffix    Type = "suffix"
	TypeImplicit  Type = "implicit"
	TypeCorrupted Type = "corrupted"
	TypeEnchant   Type = "enchant"
)

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

// Set manages affixes on an item
type Set interface {
	// Add adds affix to set
	Add(affix Affix) error

	// Remove removes affix by ID
	Remove(affixID string) error

	// Get retrieves affix by ID
	Get(affixID string) (Affix, bool)

	// GetByType returns all affixes of specified type
	GetByType(affixType Type) []Affix

	// GetAll returns all affixes
	GetAll() []Affix

	// Count returns total number of affixes
	Count() int

	// CountByType returns number of affixes of specified type
	CountByType(affixType Type) int

	// CanAdd checks if affix can be added
	CanAdd(affix Affix) bool

	// MaxPrefixes returns maximum allowed prefixes
	MaxPrefixes() int

	// MaxSuffixes returns maximum allowed suffixes
	MaxSuffixes() int

	// Clear removes all affixes
	Clear()

	// AllModifiers returns combined modifiers from all affixes
	AllModifiers() []attribute.Modifier
}

// Pool represents collection of affixes available for rolling
type Pool interface {
	// Add adds affix to pool
	Add(affix Affix)

	// Remove removes affix from pool
	Remove(affixID string)

	// Get retrieves affix by ID
	Get(affixID string) (Affix, bool)

	// GetByTier returns affixes of specified tier
	GetByTier(tier int) []Affix

	// GetByTags returns affixes matching all specified tags
	GetByTags(tags ...string) []Affix

	// Roll randomly selects affix from pool based on weights
	Roll(itemType string, itemLevel int, slot string) (Affix, error)

	// RollMultiple randomly selects multiple affixes
	RollMultiple(count int, itemType string, itemLevel int, slot string) ([]Affix, error)

	// Filter returns affixes matching criteria
	Filter(criteria FilterCriteria) []Affix
}

// FilterCriteria defines filtering for affix selection
type FilterCriteria struct {
	Types        []Type
	MinTier      int
	MaxTier      int
	Tags         []string
	MinItemLevel int
	MaxItemLevel int
}

// Generator creates affixes for items
type Generator interface {
	// Generate creates random affixes for item
	Generate(itemType string, itemLevel int, slot string, rarity int) ([]Affix, error)

	// Reroll rerolls existing affixes
	Reroll(existing []Affix, itemType string, itemLevel int, slot string) ([]Affix, error)

	// Upgrade improves affix tier
	Upgrade(affix Affix) (Affix, error)

	// AddPrefix adds random prefix
	AddPrefix(itemType string, itemLevel int, slot string) (Affix, error)

	// AddSuffix adds random suffix
	AddSuffix(itemType string, itemLevel int, slot string) (Affix, error)
}
