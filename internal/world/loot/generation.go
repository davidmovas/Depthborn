package loot

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/item"
)

// Generator generates loot drops
type Generator interface {
	// Generate creates loot drops
	Generate(ctx context.Context, params GenerationParams) ([]Drop, error)

	// GenerateForEnemy creates loot from enemy
	GenerateForEnemy(ctx context.Context, enemyType string, enemyLevel int, killerLevel int) ([]Drop, error)

	// GenerateForChest creates loot from container
	GenerateForChest(ctx context.Context, chestType string, layerDepth int) ([]Drop, error)

	// GenerateForBoss creates boss loot
	GenerateForBoss(ctx context.Context, bossType string, bossLevel int, partySize int) ([]Drop, error)

	// RollDrop performs single drop roll
	RollDrop(ctx context.Context, table DropTable, luck float64) (Drop, bool, error)

	// RollMultiple performs multiple drop rolls
	RollMultiple(ctx context.Context, table DropTable, count int, luck float64) ([]Drop, error)
}

// GenerationParams defines loot generation parameters
type GenerationParams struct {
	SourceType      string
	SourceLevel     int
	PlayerLevel     int
	LayerDepth      int
	PartySize       int
	LuckModifier    float64
	RarityBonus     float64
	QuantityBonus   float64
	DropTables      []DropTable
	GuaranteedDrops []Drop
	Context         map[string]any
}

// Drop represents single item drop
type Drop interface {
	// ItemType returns item type identifier
	ItemType() string

	// Quantity returns drop quantity
	Quantity() int

	// SetQuantity updates quantity
	SetQuantity(quantity int)

	// Rarity returns item rarity
	Rarity() item.Rarity

	// SetRarity updates rarity
	SetRarity(rarity item.Rarity)

	// Quality returns item quality [0.0 - 1.0]
	Quality() float64

	// SetQuality updates quality
	SetQuality(quality float64)

	// Level returns item level
	Level() int

	// SetLevel updates item level
	SetLevel(level int)

	// Affixes returns item affixes
	Affixes() []string

	// SetAffixes updates affixes
	SetAffixes(affixes []string)

	// IsIdentified returns true if item is identified
	IsIdentified() bool

	// SetIdentified updates identification state
	SetIdentified(identified bool)

	// SourceID returns what dropped this item
	SourceID() string

	// SetSourceID updates source
	SetSourceID(sourceID string)

	// Metadata returns additional drop data
	Metadata() map[string]any

	// CreateItem instantiates actual item
	CreateItem(ctx context.Context) (item.Item, error)
}

// DropTable defines possible drops
type DropTable interface {
	// ID returns unique table identifier
	ID() string

	// Name returns display name
	Name() string

	// Entries returns all drop entries
	Entries() []DropEntry

	// AddEntry adds drop entry
	AddEntry(entry DropEntry)

	// RemoveEntry removes drop entry
	RemoveEntry(entryID string)

	// GetEntry retrieves entry by ID
	GetEntry(entryID string) (DropEntry, bool)

	// Roll selects drop from table
	Roll(luck float64) (DropEntry, bool)

	// RollMultiple selects multiple drops
	RollMultiple(count int, luck float64, allowDuplicates bool) []DropEntry

	// TotalWeight returns sum of all entry weights
	TotalWeight() float64

	// MinDrops returns minimum drops from table
	MinDrops() int

	// MaxDrops returns maximum drops from table
	MaxDrops() int

	// SetDropRange updates drop count range
	SetDropRange(min, max int)

	// IsEmpty returns true if table has no entries
	IsEmpty() bool

	// Clear removes all entries
	Clear()
}

// DropEntry defines possible drop
type DropEntry interface {
	// ID returns unique entry identifier
	ID() string

	// ItemType returns item type to drop
	ItemType() string

	// Weight returns drop weight (higher = more common)
	Weight() float64

	// SetWeight updates drop weight
	SetWeight(weight float64)

	// DropChance returns probability of drop [0.0 - 1.0]
	DropChance() float64

	// SetDropChance updates drop chance
	SetDropChance(chance float64)

	// MinQuantity returns minimum drop count
	MinQuantity() int

	// MaxQuantity returns maximum drop count
	MaxQuantity() int

	// SetQuantityRange updates quantity range
	SetQuantityRange(min, max int)

	// RarityWeights returns weighted rarity distribution
	RarityWeights() map[item.Rarity]float64

	// SetRarityWeights updates rarity distribution
	SetRarityWeights(weights map[item.Rarity]float64)

	// QualityRange returns min and max quality
	QualityRange() (min, max float64)

	// SetQualityRange updates quality range
	SetQualityRange(min, max float64)

	// LevelOffset returns level offset from source level
	LevelOffset() int

	// SetLevelOffset updates level offset
	SetLevelOffset(offset int)

	// Conditions returns drop conditions
	Conditions() []DropCondition

	// AddCondition adds drop condition
	AddCondition(condition DropCondition)

	// RemoveCondition removes drop condition
	RemoveCondition(conditionID string)

	// CanDrop checks if entry can drop
	CanDrop(ctx context.Context, params GenerationParams) bool

	// IsGuaranteed returns true if always drops
	IsGuaranteed() bool

	// SetGuaranteed updates guaranteed flag
	SetGuaranteed(guaranteed bool)

	// Tags returns entry tags
	Tags() []string
}

// DropCondition defines when drop can occur
type DropCondition interface {
	// ID returns unique condition identifier
	ID() string

	// Type returns condition type
	Type() ConditionType

	// Check evaluates if condition is met
	Check(ctx context.Context, params GenerationParams) bool

	// Description returns human-readable condition
	Description() string
}

// ConditionType categorizes drop conditions
type ConditionType string

const (
	ConditionMinLevel    ConditionType = "min_level"
	ConditionMaxLevel    ConditionType = "max_level"
	ConditionMinDepth    ConditionType = "min_depth"
	ConditionMaxDepth    ConditionType = "max_depth"
	ConditionSourceType  ConditionType = "source_type"
	ConditionDifficulty  ConditionType = "difficulty"
	ConditionFirstKill   ConditionType = "first_kill"
	ConditionQuestActive ConditionType = "quest_active"
	ConditionPlayerClass ConditionType = "player_class"
	ConditionPartySize   ConditionType = "party_size"
	ConditionTimeOfDay   ConditionType = "time_of_day"
	ConditionWeather     ConditionType = "weather"
	ConditionRandom      ConditionType = "random"
)

// RarityDistribution manages rarity probabilities
type RarityDistribution interface {
	// GetWeight returns weight for rarity
	GetWeight(rarity item.Rarity) float64

	// SetWeight updates rarity weight
	SetWeight(rarity item.Rarity, weight float64)

	// GetProbability returns probability for rarity [0.0 - 1.0]
	GetProbability(rarity item.Rarity) float64

	// Roll selects rarity based on weights
	Roll(luck float64) item.Rarity

	// ScaleByDepth adjusts weights based on layer depth
	ScaleByDepth(depth int)

	// ScaleByLevel adjusts weights based on player level
	ScaleByLevel(level int)

	// ApplyModifier applies temporary modifier
	ApplyModifier(modifier RarityModifier)

	// RemoveModifier removes modifier
	RemoveModifier(modifierID string)

	// Reset resets to default weights
	Reset()

	// TotalWeight returns sum of all weights
	TotalWeight() float64
}

// RarityModifier adjusts rarity weights
type RarityModifier interface {
	// ID returns unique modifier identifier
	ID() string

	// Apply applies modifier to weights
	Apply(weights map[item.Rarity]float64) map[item.Rarity]float64

	// Duration returns remaining duration in milliseconds (-1 = permanent)
	Duration() int64

	// IsExpired returns true if duration ended
	IsExpired() bool
}

// QualityCalculator determines item quality
type QualityCalculator interface {
	// Calculate computes item quality
	Calculate(ctx context.Context, params QualityParams) float64

	// BaseQuality returns base quality for level
	BaseQuality(itemLevel int) float64

	// ApplyDepthBonus applies layer depth bonus
	ApplyDepthBonus(baseQuality float64, depth int) float64

	// ApplyLuckBonus applies luck modifier
	ApplyLuckBonus(baseQuality float64, luck float64) float64

	// ApplyRarityBonus applies rarity bonus
	ApplyRarityBonus(baseQuality float64, rarity item.Rarity) float64

	// Clamp clamps quality to valid range [0.0 - 1.0]
	Clamp(quality float64) float64
}

// QualityParams defines quality calculation parameters
type QualityParams struct {
	ItemLevel    int
	ItemRarity   item.Rarity
	LayerDepth   int
	LuckModifier float64
	SourceType   string
	Modifiers    []float64
}

// AffixGenerator generates item affixes
type AffixGenerator interface {
	// Generate creates affixes for item
	Generate(ctx context.Context, params AffixParams) ([]string, error)

	// GeneratePrefix generates prefix affix
	GeneratePrefix(ctx context.Context, params AffixParams) (string, error)

	// GenerateSuffix generates suffix affix
	GenerateSuffix(ctx context.Context, params AffixParams) (string, error)

	// CountForRarity returns affix count for rarity
	CountForRarity(rarity item.Rarity) (prefixes, suffixes int)

	// RollAffixTier determines affix tier
	RollAffixTier(itemLevel int, layerDepth int) int

	// CanHaveAffix checks if item can have affix
	CanHaveAffix(itemType string, affixID string) bool
}

// AffixParams defines affix generation parameters
type AffixParams struct {
	ItemType       string
	ItemLevel      int
	ItemRarity     item.Rarity
	LayerDepth     int
	AllowedTypes   []string
	ExcludeAffixes []string
	MinTier        int
	MaxTier        int
}

// Pool manages available loot
type Pool interface {
	// ID returns unique pool identifier
	ID() string

	// Name returns display name
	Name() string

	// AddTable adds drop table to pool
	AddTable(table DropTable, weight float64)

	// RemoveTable removes drop table
	RemoveTable(tableID string)

	// GetTable retrieves table by ID
	GetTable(tableID string) (DropTable, bool)

	// GetTables returns all tables
	GetTables() []DropTable

	// SelectTable chooses table based on weights
	SelectTable() (DropTable, error)

	// TotalWeight returns sum of table weights
	TotalWeight() float64

	// Clear removes all tables
	Clear()
}

// Registry manages drop tables and pools
type Registry interface {
	// RegisterTable adds drop table
	RegisterTable(table DropTable) error

	// UnregisterTable removes drop table
	UnregisterTable(tableID string) error

	// GetTable retrieves table by ID
	GetTable(tableID string) (DropTable, bool)

	// GetTables returns all registered tables
	GetTables() []DropTable

	// RegisterPool adds loot pool
	RegisterPool(pool Pool) error

	// UnregisterPool removes loot pool
	UnregisterPool(poolID string) error

	// GetPool retrieves pool by ID
	GetPool(poolID string) (Pool, bool)

	// GetPools returns all registered pools
	GetPools() []Pool

	// GetTablesForSource returns tables for source type
	GetTablesForSource(sourceType string) []DropTable

	// GetTablesForDepth returns tables for layer depth
	GetTablesForDepth(depth int) []DropTable
}

// DepthScaler adjusts loot based on layer depth
type DepthScaler interface {
	// ScaleRarity adjusts rarity chances for depth
	ScaleRarity(distribution RarityDistribution, depth int)

	// ScaleQuality adjusts quality for depth
	ScaleQuality(baseQuality float64, depth int) float64

	// ScaleQuantity adjusts drop quantity for depth
	ScaleQuantity(baseQuantity int, depth int) int

	// ScaleLevel adjusts item level for depth
	ScaleLevel(baseLevel int, depth int) int

	// GetDepthMultiplier returns multiplier for depth
	GetDepthMultiplier(depth int) float64

	// SetDepthMultiplier updates multiplier for depth
	SetDepthMultiplier(depth int, multiplier float64)

	// BonusDropChance returns extra drop chance for depth
	BonusDropChance(depth int) float64

	// UniqueChance returns unique item chance for depth
	UniqueChance(depth int) float64
}

// BossLootGenerator generates boss-specific loot
type BossLootGenerator interface {
	// Generate creates boss loot
	Generate(ctx context.Context, params BossLootParams) ([]Drop, error)

	// GuaranteedDrops returns guaranteed boss drops
	GuaranteedDrops(bossType string) []Drop

	// UniqueDrops returns boss-exclusive drops
	UniqueDrops(bossType string) []Drop

	// BonusDrops returns extra drops based on conditions
	BonusDrops(ctx context.Context, params BossLootParams) []Drop

	// FirstKillBonus returns first kill bonus loot
	FirstKillBonus(bossType string) []Drop
}

// BossLootParams defines boss loot parameters
type BossLootParams struct {
	BossType         string
	BossLevel        int
	LayerDepth       int
	PartySize        int
	Difficulty       float64
	IsFirstKill      bool
	TimeToKill       int64
	DeathCount       int
	LuckModifier     float64
	RarityBonus      float64
	PerformanceGrade string
}

// LuckCalculator calculates luck modifier
type LuckCalculator interface {
	// CalculateLuck computes total luck modifier
	CalculateLuck(ctx context.Context, sources []LuckSource) float64

	// ApplyLuck applies luck to drop chance
	ApplyLuck(baseChance float64, luck float64) float64

	// ApplyLuckToRarity applies luck to rarity roll
	ApplyLuckToRarity(distribution RarityDistribution, luck float64) RarityDistribution

	// ApplyLuckToQuality applies luck to quality
	ApplyLuckToQuality(baseQuality float64, luck float64) float64

	// LuckCap returns maximum luck value
	LuckCap() float64

	// SetLuckCap updates luck cap
	SetLuckCap(cap float64)
}

// LuckSource represents source of luck
type LuckSource interface {
	// Type returns luck source type
	Type() LuckSourceType

	// Value returns luck value
	Value() float64

	// IsActive returns true if source is active
	IsActive() bool

	// Description returns source description
	Description() string
}

// LuckSourceType categorizes luck sources
type LuckSourceType string

const (
	LuckSourceAttribute LuckSourceType = "attribute"
	LuckSourceItem      LuckSourceType = "item"
	LuckSourceBuff      LuckSourceType = "buff"
	LuckSourceSkill     LuckSourceType = "skill"
	LuckSourceParty     LuckSourceType = "party"
	LuckSourceEvent     LuckSourceType = "event"
	LuckSourcePrestige  LuckSourceType = "prestige"
)

// Instance represents loot drop instance
type Instance interface {
	// ID returns unique instance identifier
	ID() string

	// Drops returns all drops in instance
	Drops() []Drop

	// AddDrop adds drop to instance
	AddDrop(drop Drop)

	// RemoveDrop removes drop from instance
	RemoveDrop(dropIndex int)

	// SourceID returns what created this loot
	SourceID() string

	// Position returns loot position in world
	Position() (x, y, z int)

	// SetPosition updates loot position
	SetPosition(x, y, z int)

	// Timestamp returns when loot was created
	Timestamp() int64

	// ExpiresAt returns when loot despawns (0 = never)
	ExpiresAt() int64

	// SetExpiresAt updates expiration time
	SetExpiresAt(timestamp int64)

	// IsExpired returns true if loot should despawn
	IsExpired() bool

	// IsEmpty returns true if no drops remain
	IsEmpty() bool

	// Owner returns entity that can claim loot (empty = anyone)
	Owner() string

	// SetOwner updates loot owner
	SetOwner(entityID string)

	// CanLoot checks if entity can take loot
	CanLoot(entityID string) bool
}

// Manager manages loot system
type Manager interface {
	// Generator returns loot generator
	Generator() Generator

	// Registry returns loot registry
	Registry() Registry

	// DepthScaler returns depth scaler
	DepthScaler() DepthScaler

	// BossLootGenerator returns boss loot generator
	BossLootGenerator() BossLootGenerator

	// LuckCalculator returns luck calculator
	LuckCalculator() LuckCalculator

	// CreateInstance creates loot instance
	CreateInstance(ctx context.Context, drops []Drop, sourceID string, x, y, z int) (Instance, error)

	// GetInstance retrieves loot instance
	GetInstance(instanceID string) (Instance, bool)

	// RemoveInstance removes loot instance
	RemoveInstance(instanceID string)

	// GetInstancesAtPosition returns loot at position
	GetInstancesAtPosition(x, y, z int) []Instance

	// CleanupExpired removes expired loot
	CleanupExpired(ctx context.Context) int

	// Update processes loot system
	Update(ctx context.Context, deltaMs int64) error
}
