package enemy

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/enemy/ai"
)

// Enemy represents hostile entity
type Enemy interface {
	entity.Combatant

	// Family returns enemy family type
	Family() string

	// Rank returns enemy rank
	Rank() Rank

	// SetRank updates enemy rank
	SetRank(rank Rank)

	// IsElite returns true if enemy is elite
	IsElite() bool

	// IsBoss returns true if enemy is boss
	IsBoss() bool

	// AffixCount returns number of affixes
	AffixCount() int

	// Affixes returns enemy affixes
	Affixes() AffixSet

	// AI returns enemy AI controller
	AI() ai.AI

	// AggroRange returns aggression detection range
	AggroRange() float64

	// SetAggroRange updates aggro range
	SetAggroRange(range_ float64)

	// LeashRange returns maximum chase distance
	LeashRange() float64

	// SetLeashRange updates leash range
	SetLeashRange(range_ float64)

	// SpawnPosition returns original spawn location
	SpawnPosition() (x, y float64)

	// SetSpawnPosition updates spawn location
	SetSpawnPosition(x, y float64)

	// IsLeashed returns true if enemy exceeded leash range
	IsLeashed() bool

	// ResetToSpawn returns enemy to spawn position
	ResetToSpawn(ctx context.Context) error

	// LootTable returns drop table for enemy
	LootTable() LootTable

	// DropLoot generates loot on death
	DropLoot(ctx context.Context) ([]interface{}, error)

	// ExperienceReward returns XP granted on kill
	ExperienceReward() int64
}

// Rank categorizes enemy power level
type Rank int

const (
	RankNormal Rank = iota
	RankChampion
	RankElite
	RankMiniBoss
	RankBoss
	RankWorldBoss
)

// String returns rank name
func (r Rank) String() string {
	return [...]string{"Normal", "Champion", "Elite", "Mini-Boss", "Boss", "World Boss"}[r]
}

// Multiplier returns stat multiplier for rank
func (r Rank) Multiplier() float64 {
	return [...]float64{1.0, 1.25, 1.5, 1.75, 2.0, 2.5}[r]
}

// AffixSet manages enemy affixes
type AffixSet interface {
	// Add adds affix to enemy
	Add(affix Affix) error

	// Remove removes affix by ID
	Remove(affixID string)

	// Get retrieves affix by ID
	Get(affixID string) (Affix, bool)

	// GetAll returns all affixes
	GetAll() []Affix

	// Has checks if affix exists
	Has(affixID string) bool

	// Count returns number of affixes
	Count() int

	// MaxAffixes returns maximum allowed affixes
	MaxAffixes() int

	// Clear removes all affixes
	Clear()

	// Apply applies all affix effects to enemy
	Apply(ctx context.Context, enemy Enemy) error

	// RemoveFromEnemy removes all affix effects from enemy
	RemoveFromEnemy(ctx context.Context, enemy Enemy) error
}

// Affix represents modifier on enemy
type Affix interface {
	// ID returns unique affix identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Type returns affix category
	Type() AffixType

	// Apply applies affix effect to enemy
	Apply(ctx context.Context, enemy Enemy) error

	// Remove removes affix effect from enemy
	Remove(ctx context.Context, enemy Enemy) error

	// Icon returns icon identifier
	Icon() string

	// Color returns affix color for UI
	Color() (r, g, b uint8)
}

// AffixType categorizes enemy affixes
type AffixType string

const (
	AffixOffensive AffixType = "offensive"
	AffixDefensive AffixType = "defensive"
	AffixUtility   AffixType = "utility"
	AffixElemental AffixType = "elemental"
	AffixCursed    AffixType = "cursed"
	AffixBlessed   AffixType = "blessed"
)

// LootTable defines item drops
type LootTable interface {
	// Add adds loot entry
	Add(entry LootEntry)

	// Remove removes loot entry
	Remove(entryID string)

	// Roll generates loot based on luck modifier
	Roll(luckModifier float64) ([]interface{}, error)

	// GetAll returns all loot entries
	GetAll() []LootEntry

	// Clear removes all entries
	Clear()

	// GuaranteedDrops returns items that always drop
	GuaranteedDrops() []LootEntry
}

// LootEntry defines possible drop
type LootEntry interface {
	// ID returns unique entry identifier
	ID() string

	// ItemType returns type of item to drop
	ItemType() string

	// MinQuantity returns minimum drop count
	MinQuantity() int

	// MaxQuantity returns maximum drop count
	MaxQuantity() int

	// DropChance returns probability of drop [0.0 - 1.0]
	DropChance() float64

	// MinRarity returns minimum item rarity
	MinRarity() int

	// MaxRarity returns maximum item rarity
	MaxRarity() int

	// IsGuaranteed returns true if item always drops
	IsGuaranteed() bool

	// Condition returns conditional drop requirement
	Condition() DropCondition
}

// DropCondition defines when loot drops
type DropCondition interface {
	// Check returns true if condition is met
	Check(ctx context.Context, enemy Enemy, killer entity.Entity) bool

	// Description returns human-readable condition
	Description() string
}

// Family represents group of related enemies
type Family interface {
	// ID returns unique family identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Members returns all enemy types in family
	Members() []string

	// BaseStats returns default stats for family
	BaseStats() map[string]float64

	// Resistances returns elemental resistances
	Resistances() map[string]float64

	// Weaknesses returns elemental weaknesses
	Weaknesses() map[string]float64

	// Behaviors returns available AI behaviors
	Behaviors() []string

	// Tags returns family tags
	Tags() []string
}

// Registry manages enemy types and families
type Registry interface {
	// RegisterFamily adds enemy family
	RegisterFamily(family Family) error

	// RegisterType adds enemy type to family
	RegisterType(familyID string, enemyType Type) error

	// GetFamily retrieves family by ID
	GetFamily(familyID string) (Family, bool)

	// GetType retrieves enemy type by ID
	GetType(typeID string) (Type, bool)

	// GetFamilies returns all registered families
	GetFamilies() []Family

	// GetTypesInFamily returns all types in family
	GetTypesInFamily(familyID string) []Type

	// Create instantiates enemy of specified type
	Create(ctx context.Context, typeID string, level int) (Enemy, error)

	// CreateElite creates elite version of enemy
	CreateElite(ctx context.Context, typeID string, level int, affixCount int) (Enemy, error)

	// CreateBoss creates boss version of enemy
	CreateBoss(ctx context.Context, bossType string, level int) (Enemy, error)
}

// Type defines specific enemy variant
type Type interface {
	// ID returns unique type identifier
	ID() string

	// Name returns display name
	Name() string

	// Family returns family this type belongs to
	Family() string

	// BaseLevel returns recommended level
	BaseLevel() int

	// Model returns 3D model or sprite identifier
	Model() string

	// Scale returns visual scale multiplier
	Scale() float64

	// Skills returns available skills
	Skills() []string

	// Abilities returns special abilities
	Abilities() []string

	// DefaultBehavior returns default AI behavior
	DefaultBehavior() string
}
