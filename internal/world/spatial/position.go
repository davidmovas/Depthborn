package spatial

import "math"

// Position represents location in 3D space
type Position struct {
	X int // Grid X coordinate
	Y int // Grid Y coordinate
	Z int // Height level (0 = ground, positive = elevated, negative = below)
}

// NewPosition creates new position
func NewPosition(x, y, z int) Position {
	return Position{X: x, Y: y, Z: z}
}

// Equals checks if positions are identical
func (p Position) Equals(other Position) bool {
	return p.X == other.X && p.Y == other.Y && p.Z == other.Z
}

// DistanceTo calculates 3D distance to another position
func (p Position) DistanceTo(other Position) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	dz := float64(p.Z - other.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// ManhattanDistance calculates grid distance (ignoring Z)
func (p Position) ManhattanDistance(other Position) int {
	return abs(p.X-other.X) + abs(p.Y-other.Y)
}

// IsAdjacent checks if positions are neighboring (including diagonals)
func (p Position) IsAdjacent(other Position) bool {
	if p.Z != other.Z {
		return false
	}
	dx := abs(p.X - other.X)
	dy := abs(p.Y - other.Y)
	return dx <= 1 && dy <= 1 && (dx+dy) > 0
}

// IsOrthogonallyAdjacent checks if positions are neighboring (no diagonals)
func (p Position) IsOrthogonallyAdjacent(other Position) bool {
	if p.Z != other.Z {
		return false
	}
	dx := abs(p.X - other.X)
	dy := abs(p.Y - other.Y)
	return (dx == 1 && dy == 0) || (dx == 0 && dy == 1)
}

// InRange checks if position is within distance
func (p Position) InRange(other Position, maxDistance float64) bool {
	return p.DistanceTo(other) <= maxDistance
}

// Add returns new position offset by delta
func (p Position) Add(dx, dy, dz int) Position {
	return Position{X: p.X + dx, Y: p.Y + dy, Z: p.Z + dz}
}

// Neighbors returns all adjacent positions (8 directions + same Z level)
func (p Position) Neighbors() []Position {
	neighbors := make([]Position, 0, 8)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			neighbors = append(neighbors, p.Add(dx, dy, 0))
		}
	}
	return neighbors
}

// OrthogonalNeighbors returns cardinal direction neighbors (4 directions)
func (p Position) OrthogonalNeighbors() []Position {
	return []Position{
		p.Add(0, -1, 0), // North
		p.Add(1, 0, 0),  // East
		p.Add(0, 1, 0),  // South
		p.Add(-1, 0, 0), // West
	}
}

// DirectionTo calculates direction vector to another position
func (p Position) DirectionTo(other Position) Direction {
	dx := other.X - p.X
	dy := other.Y - p.Y
	return Direction{DX: dx, DY: dy}
}

// AngleTo calculates angle in radians to another position
func (p Position) AngleTo(other Position) float64 {
	dx := float64(other.X - p.X)
	dy := float64(other.Y - p.Y)
	return math.Atan2(dy, dx)
}

// Direction represents movement vector
type Direction struct {
	DX int
	DY int
}

// Normalize returns unit direction
func (d Direction) Normalize() Direction {
	if d.DX == 0 && d.DY == 0 {
		return d
	}

	absDX := abs(d.DX)
	absDY := abs(d.DY)

	if absDX > absDY {
		return Direction{DX: sign(d.DX), DY: 0}
	} else if absDY > absDX {
		return Direction{DX: 0, DY: sign(d.DY)}
	}

	return Direction{DX: sign(d.DX), DY: sign(d.DY)}
}

// Facing represents orientation in space
type Facing float64

const (
	FacingNorth     Facing = 0
	FacingNorthEast Facing = math.Pi / 4
	FacingEast      Facing = math.Pi / 2
	FacingSouthEast Facing = 3 * math.Pi / 4
	FacingSouth     Facing = math.Pi
	FacingSouthWest Facing = 5 * math.Pi / 4
	FacingWest      Facing = 3 * math.Pi / 2
	FacingNorthWest Facing = 7 * math.Pi / 4
)

// ToDirection converts facing to direction vector
func (f Facing) ToDirection() Direction {
	angle := float64(f)
	dx := int(math.Round(math.Cos(angle)))
	dy := int(math.Round(math.Sin(angle)))
	return Direction{DX: dx, DY: dy}
}

// Helper functions
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func sign(n int) int {
	if n < 0 {
		return -1
	}
	if n > 0 {
		return 1
	}
	return 0
}
