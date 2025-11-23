package spawn

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/enemy"
	"github.com/davidmovas/Depthborn/internal/world/layer"
)

// Spawner manages entity spawning in layers
type Spawner interface {
	// SpawnMonsters spawns monsters based on layer config
	SpawnMonsters(ctx context.Context, layer layer.Layer, count int) ([]enemy.Enemy, error)

	// SpawnElite spawns elite monster
	SpawnElite(ctx context.Context, layer layer.Layer, family string) (enemy.Enemy, error)

	// SpawnBoss spawns boss encounter
	SpawnBoss(ctx context.Context, layer layer.Layer, bossType string) (enemy.Enemy, error)

	// SpawnGroup spawns coordinated monster group
	SpawnGroup(ctx context.Context, layer layer.Layer, groupType string) ([]enemy.Enemy, error)

	// DespawnMonster removes monster from layer
	DespawnMonster(ctx context.Context, monsterID string) error

	// DespawnAll removes all spawned entities
	DespawnAll(ctx context.Context) error

	// GetSpawned returns all spawned entities
	GetSpawned() []enemy.Enemy

	// Count returns number of spawned entities
	Count() int

	// MaxSpawns returns maximum concurrent spawns
	MaxSpawns() int

	// SetMaxSpawns updates spawn limit
	SetMaxSpawns(max int)
}

// Table defines weighted spawn chances
type Table interface {
	// Add adds entry to spawn table
	Add(entry Entry)

	// Remove removes entry from table
	Remove(entryID string)

	// Roll randomly selects entry based on weights
	Roll() (Entry, error)

	// RollMultiple selects multiple entries
	RollMultiple(count int) ([]Entry, error)

	// GetAll returns all entries
	GetAll() []Entry

	// Clear removes all entries
	Clear()

	// TotalWeight returns sum of all entry weights
	TotalWeight() int
}

// Entry defines spawnable entity
type Entry interface {
	// ID returns unique entry identifier
	ID() string

	// EntityType returns type of entity to spawn
	EntityType() string

	// EntityFamily returns entity family
	EntityFamily() string

	// Weight returns spawn weight (higher = more common)
	Weight() int

	// MinCount returns minimum spawn count
	MinCount() int

	// MaxCount returns maximum spawn count
	MaxCount() int

	// MinDepth returns minimum layer depth
	MinDepth() int

	// MaxDepth returns maximum layer depth (0 = no limit)
	MaxDepth() int

	// Tags returns entry tags for filtering
	Tags() []string

	// Metadata returns entry-specific data
	Metadata() map[string]interface{}
}

// Pool manages spawned entities
type Pool interface {
	// Add adds entity to pool
	Add(entity enemy.Enemy)

	// Remove removes entity from pool
	Remove(entityID string) (enemy.Enemy, bool)

	// Get retrieves entity by ID
	Get(entityID string) (enemy.Enemy, bool)

	// GetAll returns all entities in pool
	GetAll() []enemy.Enemy

	// GetAlive returns all alive entities
	GetAlive() []enemy.Enemy

	// GetDead returns all dead entities
	GetDead() []enemy.Enemy

	// GetInRange returns entities within distance of point
	GetInRange(x, y, radius float64) []enemy.Enemy

	// Count returns total entities in pool
	Count() int

	// CountAlive returns number of alive entities
	CountAlive() int

	// Clear removes all entities
	Clear()

	// Update processes all entities for frame
	Update(ctx context.Context, deltaMs int64) error
}

// Wave represents timed spawn sequence
type Wave interface {
	// ID returns unique wave identifier
	ID() string

	// StartTime returns when wave begins
	StartTime() int64

	// Duration returns wave duration in milliseconds
	Duration() int64

	// IsActive returns true if wave is currently spawning
	IsActive(currentTime int64) bool

	// SpawnRate returns spawns per second
	SpawnRate() float64

	// Entries returns what to spawn
	Entries() []Entry

	// IsComplete returns true if wave finished
	IsComplete() bool

	// SetComplete marks wave as finished
	SetComplete(complete bool)
}

// WaveManager coordinates multiple spawn waves
type WaveManager interface {
	// AddWave adds wave to sequence
	AddWave(wave Wave)

	// RemoveWave removes wave from sequence
	RemoveWave(waveID string)

	// GetWave retrieves wave by ID
	GetWave(waveID string) (Wave, bool)

	// GetActive returns currently active waves
	GetActive(currentTime int64) []Wave

	// GetAll returns all waves
	GetAll() []Wave

	// Start begins wave sequence
	Start(ctx context.Context) error

	// Stop halts wave sequence
	Stop()

	// Update processes waves for elapsed time
	Update(ctx context.Context, deltaMs int64) error

	// IsComplete returns true if all waves finished
	IsComplete() bool

	// Reset resets all waves
	Reset()
}
