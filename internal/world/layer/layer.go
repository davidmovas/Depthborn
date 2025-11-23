package layer

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/infra"
	"github.com/davidmovas/Depthborn/internal/world/biome"
)

// Layer represents single persistent dungeon level
type Layer interface {
	infra.Persistent

	// Depth returns layer depth level
	Depth() int

	// Biome returns layer biome
	Biome() biome.Biome

	// Seed returns generation seed for reproducibility
	Seed() int64

	// Modifiers returns active layer modifiers
	Modifiers() ModifierSet

	// HasBoss returns true if layer has boss encounter
	HasBoss() bool

	// BossType returns boss type identifier
	BossType() string

	// IsBossDefeated returns true if boss was killed
	IsBossDefeated() bool

	// SetBossDefeated marks boss as defeated
	SetBossDefeated(defeated bool)

	// BossRespawnTime returns timestamp when boss respawns
	BossRespawnTime() int64

	// SetBossRespawnTime updates boss respawn timestamp
	SetBossRespawnTime(timestamp int64)

	// LastVisited returns timestamp of last player visit
	LastVisited() int64

	// UpdateLastVisited sets last visited to current time
	UpdateLastVisited()

	// LastRestock returns timestamp of last restock
	LastRestock() int64

	// ShouldRestock returns true if layer should refresh content
	ShouldRestock() bool

	// Restock regenerates monsters, resources, and events
	Restock(ctx context.Context) error

	// Difficulty returns layer difficulty multiplier
	Difficulty() float64

	// AllowedEnemyFamilies returns enemy types that can spawn
	AllowedEnemyFamilies() []string

	// SpawnTables returns spawn configuration
	SpawnTables() SpawnConfiguration

	// Exits returns available exit points
	Exits() []Exit

	// AddExit adds exit point to layer
	AddExit(exit Exit)

	// RemoveExit removes exit point
	RemoveExit(exitID string)

	// IsExplored returns true if layer was fully explored
	IsExplored() bool

	// SetExplored marks layer as explored
	SetExplored(explored bool)

	// Metadata returns layer-specific data
	Metadata() map[string]any
}

// ModifierSet manages layer modifiers
type ModifierSet interface {
	// Add adds modifier to layer
	Add(modifier Modifier)

	// Remove removes modifier by ID
	Remove(modifierID string)

	// Get retrieves modifier by ID
	Get(modifierID string) (Modifier, bool)

	// GetAll returns all active modifiers
	GetAll() []Modifier

	// Has checks if modifier exists
	Has(modifierID string) bool

	// Clear removes all modifiers
	Clear()

	// Apply applies all modifiers to value
	Apply(baseValue float64, context string) float64
}

// Modifier alters layer properties
type Modifier interface {
	// ID returns unique modifier identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Type returns modifier type
	Type() ModifierType

	// Value returns modifier value
	Value() float64

	// Icon returns icon identifier
	Icon() string

	// Tags returns modifier tags
	Tags() []string
}

// ModifierType categorizes layer modifiers
type ModifierType string

const (
	ModTypeDifficulty     ModifierType = "difficulty"
	ModTypeMonsterDensity ModifierType = "monster_density"
	ModTypeEliteChance    ModifierType = "elite_chance"
	ModTypeLootQuantity   ModifierType = "loot_quantity"
	ModTypeLootQuality    ModifierType = "loot_quality"
	ModTypeExperienceGain ModifierType = "experience_gain"
	ModTypeEnvironmental  ModifierType = "environmental"
	ModTypeHazard         ModifierType = "hazard"
	ModTypeCurse          ModifierType = "curse"
	ModTypeBlessing       ModifierType = "blessing"
)

// SpawnConfiguration defines spawn rules for layer
type SpawnConfiguration interface {
	// MonsterDensity returns base monster count multiplier
	MonsterDensity() float64

	// SetMonsterDensity updates density multiplier
	SetMonsterDensity(density float64)

	// EliteChance returns probability of elite monsters [0.0 - 1.0]
	EliteChance() float64

	// SetEliteChance updates elite spawn chance
	SetEliteChance(chance float64)

	// BossChance returns probability of mini-boss spawns [0.0 - 1.0]
	BossChance() float64

	// SetBossChance updates mini-boss spawn chance
	SetBossChance(chance float64)

	// ResourceDensity returns resource spawn multiplier
	ResourceDensity() float64

	// SetResourceDensity updates resource density
	SetResourceDensity(density float64)

	// EventChance returns probability of special events [0.0 - 1.0]
	EventChance() float64

	// SetEventChance updates event spawn chance
	SetEventChance(chance float64)

	// TrapDensity returns trap spawn multiplier
	TrapDensity() float64

	// SetTrapDensity updates trap density
	SetTrapDensity(density float64)
}

// Exit represents way to leave layer
type Exit interface {
	// ID returns unique exit identifier
	ID() string

	// Type returns exit type
	Type() ExitType

	// TargetDepth returns destination depth (-1 = camp/hub)
	TargetDepth() int

	// IsLocked returns true if exit requires unlocking
	IsLocked() bool

	// Unlock opens the exit
	Unlock()

	// Position returns exit location in layer
	Position() (x, y float64)

	// SetPosition updates exit location
	SetPosition(x, y float64)

	// RequiresObjective returns true if objective must be completed
	RequiresObjective() bool

	// ObjectiveComplete returns true if objective is done
	ObjectiveComplete() bool

	// SetObjectiveComplete marks objective as complete
	SetObjectiveComplete(complete bool)
}

// ExitType categorizes exits
type ExitType string

const (
	ExitStairs   ExitType = "stairs"
	ExitPortal   ExitType = "portal"
	ExitLadder   ExitType = "ladder"
	ExitTeleport ExitType = "teleport"
	ExitSecret   ExitType = "secret"
	ExitBoss     ExitType = "boss"
)

// Generator creates layers procedurally
type Generator interface {
	// Generate creates new layer at specified depth
	Generate(ctx context.Context, depth int) (Layer, error)

	// GenerateWithSeed creates layer with specific seed
	GenerateWithSeed(ctx context.Context, depth int, seed int64) (Layer, error)

	// SelectBiome chooses appropriate biome for depth
	SelectBiome(depth int) biome.Biome

	// SelectModifiers chooses layer modifiers based on depth
	SelectModifiers(depth int) []Modifier

	// CalculateDifficulty determines difficulty multiplier for depth
	CalculateDifficulty(depth int) float64
}

// Registry manages all layers in dungeon
type Registry interface {
	// Register adds layer to registry
	Register(ctx context.Context, layer Layer) error

	// Unregister removes layer from registry
	Unregister(depth int) error

	// Get retrieves layer by depth
	Get(depth int) (Layer, bool)

	// GetOrGenerate retrieves existing layer or generates new one
	GetOrGenerate(ctx context.Context, depth int) (Layer, error)

	// GetRange returns layers within depth range
	GetRange(minDepth, maxDepth int) []Layer

	// Has checks if layer exists at depth
	Has(depth int) bool

	// MaxDepth returns deepest explored layer
	MaxDepth() int

	// Count returns total number of layers
	Count() int

	// Save persists layer to storage
	Save(ctx context.Context, layer Layer) error

	// Load loads layer from storage
	Load(ctx context.Context, depth int) (Layer, error)

	// Clear removes all layers from registry
	Clear() error
}
