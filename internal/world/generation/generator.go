package generation

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/world/biome"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

// Generator creates procedural content
type Generator interface {
	// GenerateLayout creates spatial layout for layer
	GenerateLayout(ctx context.Context, config LayoutConfig) (Layout, error)

	// PlaceRooms positions rooms in layout
	PlaceRooms(ctx context.Context, layout Layout, rooms []Room) error

	// ConnectRooms creates corridors between rooms
	ConnectRooms(ctx context.Context, layout Layout) error

	// PlaceExits positions entry and exit points
	PlaceExits(ctx context.Context, layout Layout) error

	// PlaceResources positions resource nodes
	PlaceResources(ctx context.Context, layout Layout, count int) error

	// PlaceHazards positions traps and hazards
	PlaceHazards(ctx context.Context, layout Layout, count int) error
}

// LayoutConfig defines generation parameters
type LayoutConfig struct {
	Width       int
	Height      int
	RoomCount   int
	MinRoomSize int
	MaxRoomSize int
	Biome       biome.Biome
	Seed        int64
	Density     float64
	Complexity  float64
}

// Layout represents spatial structure of layer
type Layout interface {
	// Width returns layout width
	Width() int

	// Height returns layout height
	Height() int

	// Rooms returns all rooms in layout
	Rooms() []Room

	// GetRoom returns room at coordinates
	GetRoom(x, y int) (Room, bool)

	// Corridors returns all corridors
	Corridors() []Corridor

	// IsWalkable returns true if position is traversable
	IsWalkable(x, y int) bool

	// GetTile returns tile type at position
	GetTile(x, y int) TileType

	// SetTile updates tile at position
	SetTile(x, y int, tile TileType)

	// FindPath calculates path between two points
	FindPath(startX, startY, endX, endY int) []spatial.Position

	// GetSpawnPositions returns valid spawn locations
	GetSpawnPositions() []spatial.Position

	// Metadata returns layout metadata
	Metadata() map[string]any
}

// Room represents enclosed space in layout
type Room interface {
	// ID returns unique room identifier
	ID() string

	// Bounds returns room boundaries
	Bounds() Rectangle

	// Type returns room type
	Type() RoomType

	// Connections returns connected rooms
	Connections() []string

	// AddConnection links room to another
	AddConnection(roomID string)

	// SpawnPositions returns spawn locations in room
	SpawnPositions() []spatial.Position

	// IsEntrance returns true if room is layer entrance
	IsEntrance() bool

	// SetEntrance marks room as entrance
	SetEntrance(entrance bool)

	// IsExit returns true if room contains exit
	IsExit() bool

	// SetExit marks room as exit room
	SetExit(exit bool)

	// IsBossRoom returns true if room contains boss
	IsBossRoom() bool

	// SetBossRoom marks room as boss encounter
	SetBossRoom(boss bool)

	// Features returns special features in room
	Features() []RoomFeature
}

// RoomType categorizes rooms
type RoomType string

const (
	RoomNormal   RoomType = "normal"
	RoomLarge    RoomType = "large"
	RoomSmall    RoomType = "small"
	RoomTreasure RoomType = "treasure"
	RoomBoss     RoomType = "boss"
	RoomSecret   RoomType = "secret"
	RoomSafe     RoomType = "safe"
	RoomArena    RoomType = "arena"
)

// RoomFeature represents special room element
type RoomFeature interface {
	// Type returns feature type
	Type() string

	// Position returns feature location
	Position() spatial.Position

	// IsInteractable returns true if player can interact
	IsInteractable() bool

	// Interact handles player interaction
	Interact(ctx context.Context, entity any) error
}

// Corridor represents connection between rooms
type Corridor interface {
	// From returns starting room ID
	From() string

	// To returns destination room ID
	To() string

	// Path returns corridor path points
	Path() []spatial.Position

	// Width returns corridor width
	Width() int
}

// Rectangle represents rectangular bounds
type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

// TileType defines tile properties
type TileType int

const (
	TileVoid TileType = iota
	TileFloor
	TileWall
	TileDoor
	TileStairs
	TileWater
	TileLava
	TilePit
	TileChasm
)
