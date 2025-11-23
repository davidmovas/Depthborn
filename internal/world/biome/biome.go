package biome

import "context"

// Biome defines environmental theme and properties
type Biome interface {
	// ID returns unique biome identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// DepthRange returns minimum and maximum depths for this biome
	DepthRange() (min, max int)

	// Tileset returns tileset identifier for rendering
	Tileset() string

	// AmbientColor returns ambient color tint
	AmbientColor() (r, g, b, a uint8)

	// MusicTrack returns music track identifier
	MusicTrack() string

	// EnemyFamilies returns enemy types native to this biome
	EnemyFamilies() []string

	// BossTypes returns possible boss types for this biome
	BossTypes() []string

	// ResourceTypes returns resources that spawn in this biome
	ResourceTypes() []string

	// HazardTypes returns environmental hazards
	HazardTypes() []string

	// EnvironmentalEffects returns passive effects in biome
	EnvironmentalEffects() []EnvironmentalEffect

	// RoomTypes returns room layout types for this biome
	RoomTypes() []string

	// Difficulty returns base difficulty multiplier
	Difficulty() float64

	// Metadata returns biome-specific data
	Metadata() map[string]any
}

// EnvironmentalEffect represents passive biome effect
type EnvironmentalEffect interface {
	// ID returns unique effect identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Interval returns how often effect triggers in milliseconds
	Interval() int64

	// Apply applies effect to entities in biome
	Apply(ctx context.Context, entities []any) error

	// Icon returns icon identifier
	Icon() string
}

// Registry manages available biomes
type Registry interface {
	// Register adds biome to registry
	Register(biome Biome) error

	// Unregister removes biome from registry
	Unregister(biomeID string) error

	// Get retrieves biome by ID
	Get(biomeID string) (Biome, bool)

	// GetAll returns all registered biomes
	GetAll() []Biome

	// GetForDepth returns biomes suitable for specified depth
	GetForDepth(depth int) []Biome

	// Has checks if biome is registered
	Has(biomeID string) bool

	// Random selects random biome for depth
	Random(depth int) (Biome, error)
}
