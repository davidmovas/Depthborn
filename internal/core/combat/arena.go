package combat

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

// Arena represents combat battlefield
type Arena interface {
	// ID returns unique arena identifier
	ID() string

	// Name returns arena name
	Name() string

	// Description returns arena description
	Description() string

	// Grid returns spatial grid
	Grid() spatial.Grid

	// SpawnPoints returns available spawn locations
	SpawnPoints() SpawnPointSet

	// GetSpawnPoint retrieves spawn point by ID
	GetSpawnPoint(spawnID string) (SpawnPoint, bool)

	// Hazards returns environmental hazards
	Hazards() []Hazard

	// GetHazard retrieves hazard by ID
	GetHazard(hazardID string) (Hazard, bool)

	// AddHazard adds hazard to arena
	AddHazard(hazard Hazard) error

	// RemoveHazard removes hazard from arena
	RemoveHazard(hazardID string) error

	// Interactives returns interactive objects
	Interactives() []Interactive

	// GetInteractive retrieves interactive by ID
	GetInteractive(interactiveID string) (Interactive, bool)

	// AddInteractive adds interactive object
	AddInteractive(interactive Interactive) error

	// RemoveInteractive removes interactive object
	RemoveInteractive(interactiveID string) error

	// HeightAt returns height level at position
	HeightAt(pos spatial.Position) int

	// CanSee checks line of sight between positions
	CanSee(from, to spatial.Position) bool

	// GetCover calculates cover bonus at position
	GetCover(pos spatial.Position, from spatial.Position) CoverType

	// IsHighGround checks if position is elevated relative to other
	IsHighGround(pos, other spatial.Position) bool

	// AmbientEffects returns passive arena effects
	AmbientEffects() []AmbientEffect

	// AddAmbientEffect adds passive effect
	AddAmbientEffect(effect AmbientEffect)

	// RemoveAmbientEffect removes passive effect
	RemoveAmbientEffect(effectID string)

	// Weather returns current weather effect
	Weather() Weather

	// SetWeather updates weather condition
	SetWeather(weather Weather)

	// OnEntityEnter registers callback when entity enters arena
	OnEntityEnter(callback ArenaCallback)

	// OnEntityExit registers callback when entity exits arena
	OnEntityExit(callback ArenaCallback)

	// OnEntityMove registers callback when entity moves
	OnEntityMove(callback ArenaMoveCallback)

	// Update processes arena state for frame
	Update(ctx context.Context, deltaMs int64) error

	// Cleanup removes expired hazards and effects
	Cleanup(ctx context.Context) error
}

// ArenaCallback is invoked for arena events
type ArenaCallback func(ctx context.Context, arena Arena, entityID string)

// ArenaMoveCallback is invoked when entity moves
type ArenaMoveCallback func(ctx context.Context, arena Arena, entityID string, from, to spatial.Position)

// SpawnPointSet manages spawn locations
type SpawnPointSet interface {
	// GetAll returns all spawn points
	GetAll() []SpawnPoint

	// GetByTeam returns spawn points for team
	GetByTeam(team Team) []SpawnPoint

	// GetAvailable returns unoccupied spawn points
	GetAvailable() []SpawnPoint

	// GetAvailableForTeam returns unoccupied spawn points for team
	GetAvailableForTeam(team Team) []SpawnPoint

	// Reserve marks spawn point as occupied
	Reserve(spawnID string, entityID string) error

	// Release marks spawn point as available
	Release(spawnID string) error

	// IsOccupied checks if spawn point is in use
	IsOccupied(spawnID string) bool

	// GetOccupant returns entity at spawn point
	GetOccupant(spawnID string) (string, bool)

	// Add adds spawn point to set
	Add(spawn SpawnPoint) error

	// Remove removes spawn point from set
	Remove(spawnID string) error

	// Clear removes all spawn points
	Clear()

	// Count returns total spawn points
	Count() int
}

// SpawnPoint represents entry location
type SpawnPoint interface {
	// ID returns unique spawn point identifier
	ID() string

	// Position returns spawn position
	Position() spatial.Position

	// Team returns which team uses this spawn
	Team() Team

	// SetTeam updates spawn team
	SetTeam(team Team)

	// Facing returns default facing direction
	Facing() spatial.Facing

	// SetFacing updates default facing
	SetFacing(facing spatial.Facing)

	// Priority returns spawn priority (higher = preferred)
	Priority() int

	// SetPriority updates spawn priority
	SetPriority(priority int)

	// IsEnabled returns true if spawn point is active
	IsEnabled() bool

	// Enable activates spawn point
	Enable()

	// Disable deactivates spawn point
	Disable()

	// Radius returns spawn radius for group spawns
	Radius() float64

	// IsGroupSpawn returns true if can spawn multiple entities
	IsGroupSpawn() bool

	// MaxOccupants returns maximum entities at this spawn
	MaxOccupants() int
}

// Hazard represents environmental danger
type Hazard interface {
	// ID returns unique hazard identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Type returns hazard type
	Type() HazardType

	// Area returns affected area
	Area() spatial.Area

	// IsActive returns true if hazard is dangerous
	IsActive() bool

	// Activate enables hazard
	Activate()

	// Deactivate disables hazard
	Deactivate()

	// Toggle switches active state
	Toggle()

	// Damage returns damage per tick
	Damage() float64

	// SetDamage updates damage amount
	SetDamage(damage float64)

	// DamageType returns type of damage dealt
	DamageType() DamageType

	// TickInterval returns milliseconds between damage ticks
	TickInterval() int64

	// StatusEffect returns status effect ID applied to victims
	StatusEffect() string

	// StatusChance returns probability of applying status [0.0 - 1.0]
	StatusChance() float64

	// OnEnter is called when entity enters hazard
	OnEnter(ctx context.Context, entityID string, encounter Encounter) error

	// OnTick is called periodically while entity is in hazard
	OnTick(ctx context.Context, entityID string, deltaMs int64, encounter Encounter) error

	// OnExit is called when entity leaves hazard
	OnExit(ctx context.Context, entityID string, encounter Encounter) error

	// IsImmuneToHazard checks if entity is immune
	IsImmuneToHazard(entityID string, encounter Encounter) bool

	// EntitiesInHazard returns IDs of entities currently in hazard
	EntitiesInHazard() []string

	// Duration returns hazard duration in milliseconds (-1 = permanent)
	Duration() int64

	// RemainingDuration returns time left
	RemainingDuration() int64

	// IsExpired returns true if duration ended
	IsExpired() bool

	// Icon returns icon identifier
	Icon() string

	// VisualEffect returns visual effect identifier
	VisualEffect() string
}

// HazardType categorizes hazards
type HazardType string

const (
	HazardFire      HazardType = "fire"
	HazardPoison    HazardType = "poison"
	HazardSpikes    HazardType = "spikes"
	HazardPit       HazardType = "pit"
	HazardLava      HazardType = "lava"
	HazardIce       HazardType = "ice"
	HazardElectric  HazardType = "electric"
	HazardAcid      HazardType = "acid"
	HazardThorns    HazardType = "thorns"
	HazardVoid      HazardType = "void"
	HazardBleed     HazardType = "bleed"
	HazardCursed    HazardType = "cursed"
	HazardRadiation HazardType = "radiation"
	HazardQuicksand HazardType = "quicksand"
)

// Interactive represents interactable arena object
type Interactive interface {
	// ID returns unique interactive identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Type returns interactive type
	Type() InteractiveType

	// Position returns object position
	Position() spatial.Position

	// SetPosition updates object position
	SetPosition(pos spatial.Position)

	// IsEnabled returns true if can be interacted with
	IsEnabled() bool

	// Enable activates interactive
	Enable()

	// Disable deactivates interactive
	Disable()

	// Interact handles interaction
	Interact(ctx context.Context, entityID string, encounter Encounter) error

	// CanInteract checks if entity can interact
	CanInteract(entityID string, encounter Encounter) bool

	// RequiresAdjacent returns true if must be adjacent to interact
	RequiresAdjacent() bool

	// InteractionRange returns maximum interaction distance
	InteractionRange() float64

	// UsesRemaining returns number of uses left (-1 = unlimited)
	UsesRemaining() int

	// DecrementUses reduces uses by one
	DecrementUses()

	// IsExpended returns true if no uses remain
	IsExpended() bool

	// Effect returns interaction effect
	Effect() InteractiveEffect

	// Cooldown returns cooldown between uses in milliseconds
	Cooldown() int64

	// RemainingCooldown returns time until can be used again
	RemainingCooldown() int64

	// IsOnCooldown returns true if cooling down
	IsOnCooldown() bool

	// Icon returns icon identifier
	Icon() string

	// Model returns 3D model identifier
	Model() string

	// IsVisible returns true if visible to entities
	IsVisible() bool

	// SetVisible updates visibility
	SetVisible(visible bool)
}

// InteractiveType categorizes interactive objects
type InteractiveType string

const (
	InteractiveLever      InteractiveType = "lever"
	InteractiveButton     InteractiveType = "button"
	InteractiveDoor       InteractiveType = "door"
	InteractiveChest      InteractiveType = "chest"
	InteractiveShrine     InteractiveType = "shrine"
	InteractiveTrap       InteractiveType = "trap"
	InteractiveBarrel     InteractiveType = "barrel"
	InteractiveTeleporter InteractiveType = "teleporter"
	InteractiveFountain   InteractiveType = "fountain"
	InteractiveAltar      InteractiveType = "altar"
	InteractiveCrystal    InteractiveType = "crystal"
	InteractivePillar     InteractiveType = "pillar"
	InteractiveStatue     InteractiveType = "statue"
)

// InteractiveEffect describes interaction result
type InteractiveEffect interface {
	// Type returns effect type
	Type() InteractiveEffectType

	// Apply applies effect to entity and encounter
	Apply(ctx context.Context, entityID string, arena Arena, encounter Encounter) error

	// Description returns human-readable description
	Description() string

	// TargetsAllies returns true if affects allied entities
	TargetsAllies() bool

	// TargetsEnemies returns true if affects enemy entities
	TargetsEnemies() bool

	// Range returns effect range
	Range() float64

	// Area returns affected area (nil if single target)
	Area() spatial.Area
}

// InteractiveEffectType categorizes interactive effects
type InteractiveEffectType string

const (
	InteractiveEffectHeal         InteractiveEffectType = "heal"
	InteractiveEffectDamage       InteractiveEffectType = "damage"
	InteractiveEffectBuff         InteractiveEffectType = "buff"
	InteractiveEffectDebuff       InteractiveEffectType = "debuff"
	InteractiveEffectSummon       InteractiveEffectType = "summon"
	InteractiveEffectTeleport     InteractiveEffectType = "teleport"
	InteractiveEffectSpawnHazard  InteractiveEffectType = "spawn_hazard"
	InteractiveEffectRemoveHazard InteractiveEffectType = "remove_hazard"
	InteractiveEffectOpenDoor     InteractiveEffectType = "open_door"
	InteractiveEffectTriggerTrap  InteractiveEffectType = "trigger_trap"
	InteractiveEffectLoot         InteractiveEffectType = "loot"
)

// CoverType defines protection level
type CoverType int

const (
	CoverNone CoverType = iota
	CoverPartial
	CoverFull
)

// String returns cover type name
func (c CoverType) String() string {
	switch c {
	case CoverNone:
		return "none"
	case CoverPartial:
		return "partial"
	case CoverFull:
		return "full"
	default:
		return "unknown"
	}
}

// DefenseBonus returns defense bonus from cover
func (c CoverType) DefenseBonus() float64 {
	switch c {
	case CoverNone:
		return 0.0
	case CoverPartial:
		return 0.25
	case CoverFull:
		return 0.5
	default:
		return 0.0
	}
}

// EvasionBonus returns evasion bonus from cover
func (c CoverType) EvasionBonus() float64 {
	switch c {
	case CoverNone:
		return 0.0
	case CoverPartial:
		return 0.15
	case CoverFull:
		return 0.35
	default:
		return 0.0
	}
}

// AmbientEffect represents passive arena effect
type AmbientEffect interface {
	// ID returns unique effect identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Type returns effect type
	Type() AmbientEffectType

	// Interval returns milliseconds between ticks
	Interval() int64

	// Apply applies effect to all entities in arena
	Apply(ctx context.Context, entityIDs []string, encounter Encounter) error

	// AffectsTeam checks if effect applies to team
	AffectsTeam(team Team) bool

	// IsActive returns true if effect is active
	IsActive() bool

	// Activate enables effect
	Activate()

	// Deactivate disables effect
	Deactivate()

	// Duration returns effect duration in milliseconds (-1 = permanent)
	Duration() int64

	// RemainingDuration returns time left
	RemainingDuration() int64

	// IsExpired returns true if duration ended
	IsExpired() bool

	// Intensity returns effect strength [0.0 - 1.0]
	Intensity() float64

	// SetIntensity updates effect strength
	SetIntensity(intensity float64)
}

// AmbientEffectType categorizes ambient effects
type AmbientEffectType string

const (
	AmbientHeal            AmbientEffectType = "heal"
	AmbientDamage          AmbientEffectType = "damage"
	AmbientSlowRegen       AmbientEffectType = "slow_regen"
	AmbientFastRegen       AmbientEffectType = "fast_regen"
	AmbientReducedDamage   AmbientEffectType = "reduced_damage"
	AmbientIncreasedDamage AmbientEffectType = "increased_damage"
	AmbientSlowed          AmbientEffectType = "slowed"
	AmbientHasted          AmbientEffectType = "hasted"
	AmbientWeakened        AmbientEffectType = "weakened"
	AmbientEmpowered       AmbientEffectType = "empowered"
	AmbientDraining        AmbientEffectType = "draining"
	AmbientRestoring       AmbientEffectType = "restoring"
)

// Weather represents environmental condition
type Weather interface {
	// Type returns weather type
	Type() WeatherType

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// VisibilityModifier returns vision range modifier
	VisibilityModifier() float64

	// MovementModifier returns speed modifier
	MovementModifier() float64

	// DamageModifier returns damage modifier for type
	DamageModifier(damageType DamageType) float64

	// AccuracyModifier returns hit chance modifier
	AccuracyModifier() float64

	// EvasionModifier returns evasion modifier
	EvasionModifier() float64

	// StatusEffects returns weather-applied status effects
	StatusEffects() []string

	// StatusChance returns probability of applying status [0.0 - 1.0]
	StatusChance() float64

	// TickInterval returns milliseconds between weather ticks
	TickInterval() int64

	// OnWeatherTick is called periodically
	OnWeatherTick(ctx context.Context, entityIDs []string, encounter Encounter) error

	// Intensity returns weather strength [0.0 - 1.0]
	Intensity() float64

	// SetIntensity updates weather strength
	SetIntensity(intensity float64)

	// Icon returns icon identifier
	Icon() string

	// VisualEffect returns visual effect identifier
	VisualEffect() string

	// SoundEffect returns sound effect identifier
	SoundEffect() string
}

// WeatherType categorizes weather
type WeatherType string

const (
	WeatherClear     WeatherType = "clear"
	WeatherRain      WeatherType = "rain"
	WeatherStorm     WeatherType = "storm"
	WeatherFog       WeatherType = "fog"
	WeatherSnow      WeatherType = "snow"
	WeatherBlizzard  WeatherType = "blizzard"
	WeatherSandstorm WeatherType = "sandstorm"
	WeatherAshfall   WeatherType = "ashfall"
	WeatherBloodRain WeatherType = "blood_rain"
	WeatherVoidStorm WeatherType = "void_storm"
	WeatherEclipse   WeatherType = "eclipse"
	WeatherAurora    WeatherType = "aurora"
)

// ArenaBuilder creates arenas with fluent API
type ArenaBuilder interface {
	// WithName sets arena name
	WithName(name string) ArenaBuilder

	// WithDescription sets arena description
	WithDescription(description string) ArenaBuilder

	// WithGrid sets spatial grid
	WithGrid(grid spatial.Grid) ArenaBuilder

	// WithSpawnPoints adds spawn points
	WithSpawnPoints(spawns []SpawnPoint) ArenaBuilder

	// WithHazards adds hazards
	WithHazards(hazards []Hazard) ArenaBuilder

	// WithInteractives adds interactive objects
	WithInteractives(interactives []Interactive) ArenaBuilder

	// WithWeather sets weather condition
	WithWeather(weather Weather) ArenaBuilder

	// WithAmbientEffects adds ambient effects
	WithAmbientEffects(effects []AmbientEffect) ArenaBuilder

	// Build creates the arena
	Build() (Arena, error)

	// Reset resets builder to initial state
	Reset() ArenaBuilder
}

// DamageType categorizes damage
type DamageType string

const (
	DamagePhysical  DamageType = "physical"
	DamageMagical   DamageType = "magical"
	DamageFire      DamageType = "fire"
	DamageCold      DamageType = "cold"
	DamageLightning DamageType = "lightning"
	DamagePoison    DamageType = "poison"
	DamageNecrotic  DamageType = "necrotic"
	DamageRadiant   DamageType = "radiant"
	DamagePsychic   DamageType = "psychic"
	DamageForce     DamageType = "force"
	DamageTrue      DamageType = "true" // Ignores all defenses
	DamageAcid      DamageType = "acid"
	DamageBleed     DamageType = "bleed"
)
