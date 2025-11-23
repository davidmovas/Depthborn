package spatial

// Grid represents spatial layout
type Grid interface {
	// Width returns grid width
	Width() int

	// Height returns grid height
	Height() int

	// MinZ returns minimum height level
	MinZ() int

	// MaxZ returns maximum height level
	MaxZ() int

	// IsValid checks if position is within bounds
	IsValid(pos Position) bool

	// IsWalkable checks if position can be traversed
	IsWalkable(pos Position) bool

	// IsOccupied checks if position has entity
	IsOccupied(pos Position) bool

	// GetOccupant returns entity at position
	GetOccupant(pos Position) (string, bool)

	// SetOccupant places entity at position
	SetOccupant(pos Position, entityID string) error

	// RemoveOccupant removes entity from position
	RemoveOccupant(pos Position) error

	// GetTile returns tile type at position
	GetTile(pos Position) TileType

	// SetTile updates tile at position
	SetTile(pos Position, tile TileType)

	// FindPath calculates path between positions
	FindPath(from, to Position) ([]Position, error)

	// GetNeighbors returns walkable neighboring positions
	GetNeighbors(pos Position) []Position

	// InLineOfSight checks if positions have clear view
	InLineOfSight(from, to Position) bool

	// GetEntitiesInRange returns all entities within distance
	GetEntitiesInRange(center Position, radius float64) []string

	// GetEntitiesInArea returns all entities in area
	GetEntitiesInArea(area Area) []string
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
	TileIce
	TileGrass
	TileSand
)

// IsWalkable returns true if tile can be traversed
func (t TileType) IsWalkable() bool {
	switch t {
	case TileFloor, TileGrass, TileSand, TileIce, TileStairs:
		return true
	default:
		return false
	}
}

// IsTransparent returns true if tile doesn't block vision
func (t TileType) IsTransparent() bool {
	switch t {
	case TileWall:
		return false
	default:
		return true
	}
}
