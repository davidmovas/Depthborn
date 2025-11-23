package crafting

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/item"
)

// Transmuter converts items into other items
type Transmuter interface {
	// Transmute converts source items into target
	Transmute(ctx context.Context, sources []item.Item, target string) (item.Item, error)

	// CanTransmute checks if transmutation is possible
	CanTransmute(sources []item.Item, target string) bool

	// GetTransmutations returns possible transmutations
	GetTransmutations(sources []item.Item) []Transmutation

	// Combine merges multiple items into one
	Combine(ctx context.Context, items []item.Item) (item.Item, error)

	// CanCombine checks if items can be combined
	CanCombine(items []item.Item) bool

	// Split divides item into multiple items
	Split(ctx context.Context, source item.Item, count int) ([]item.Item, error)

	// CanSplit checks if item can be split
	CanSplit(source item.Item, count int) bool
}

// Transmutation defines item conversion
type Transmutation interface {
	// ID returns unique transmutation identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Sources returns required source items
	Sources() []TransmutationSource

	// Result returns transmutation output
	Result() TransmutationResult

	// SuccessChance returns probability of success [0.0 - 1.0]
	SuccessChance() float64

	// Cost returns currency cost
	Cost() int64

	// Requirements returns transmutation requirements
	Requirements() Requirements

	// CanPerform checks if transmutation is possible
	CanPerform(sources []item.Item, entityID string) bool
}

// TransmutationSource specifies source requirement
type TransmutationSource interface {
	// ItemType returns required item type
	ItemType() string

	// MinQuantity returns minimum amount needed
	MinQuantity() int

	// MinRarity returns minimum rarity required
	MinRarity() item.Rarity

	// MinQuality returns minimum quality required
	MinQuality() float64

	// IsConsumed returns true if source is destroyed
	IsConsumed() bool

	// Alternatives returns alternative item types
	Alternatives() []string
}

// TransmutationResult defines conversion output
type TransmutationResult interface {
	// ItemType returns produced item type
	ItemType() string

	// Quantity returns output amount
	Quantity() int

	// PreserveQuality returns true if quality inherited
	PreserveQuality() bool

	// PreserveAffixes returns true if affixes transferred
	PreserveAffixes() bool

	// QualityModifier returns quality adjustment
	QualityModifier() float64

	// RarityWeights returns weighted rarity chances
	RarityWeights() map[item.Rarity]float64
}

// Corruption alters items unpredictably
type Corruption interface {
	// Corrupt applies corruption to item
	Corrupt(ctx context.Context, target item.Item) (item.Item, error)

	// CanCorrupt checks if item can be corrupted
	CanCorrupt(target item.Item) bool

	// IsCorrupted checks if item is corrupted
	IsCorrupted(target item.Item) bool

	// PossibleOutcomes returns potential corruption results
	PossibleOutcomes(target item.Item) []CorruptionOutcome

	// CorruptionChance returns success probability [0.0 - 1.0]
	CorruptionChance() float64
}

// CorruptionOutcome defines corruption result
type CorruptionOutcome interface {
	// Type returns outcome type
	Type() CorruptionType

	// Chance returns probability of this outcome [0.0 - 1.0]
	Chance() float64

	// Description returns human-readable outcome
	Description() string

	// Effect returns what happens to item
	Effect() CorruptionEffect
}

// CorruptionType categorizes corruption results
type CorruptionType string

const (
	CorruptionSuccess   CorruptionType = "success"
	CorruptionUpgrade   CorruptionType = "upgrade"
	CorruptionDowngrade CorruptionType = "downgrade"
	CorruptionDestroy   CorruptionType = "destroy"
	CorruptionTransform CorruptionType = "transform"
	CorruptionNothing   CorruptionType = "nothing"
)

// CorruptionEffect describes corruption impact
type CorruptionEffect interface {
	// Apply applies corruption effect to item
	Apply(ctx context.Context, target item.Item) (item.Item, error)

	// Description returns human-readable effect
	Description() string
}

// Imprinter copies properties between items
type Imprinter interface {
	// Imprint copies properties from source to target
	Imprint(ctx context.Context, source, target item.Equipment) (item.Equipment, error)

	// CanImprint checks if imprinting is possible
	CanImprint(source, target item.Equipment) bool

	// ImprintableProperties returns properties that can be copied
	ImprintableProperties(source item.Equipment) []string

	// CreateImprint saves item properties for later
	CreateImprint(ctx context.Context, source item.Equipment) (Imprint, error)

	// ApplyImprint applies saved properties to item
	ApplyImprint(ctx context.Context, imprint Imprint, target item.Equipment) (item.Equipment, error)
}

// Imprint represents saved item properties
type Imprint interface {
	// ID returns unique imprint identifier
	ID() string

	// SourceItemID returns original item ID
	SourceItemID() string

	// Properties returns saved properties
	Properties() map[string]any

	// CreatedAt returns when imprint was created
	CreatedAt() int64

	// ExpiresAt returns when imprint becomes invalid (0 = never)
	ExpiresAt() int64

	// IsExpired returns true if imprint is no longer valid
	IsExpired() bool
}
