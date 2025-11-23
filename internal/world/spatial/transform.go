package spatial

// Transform represents complete spatial state
type Transform interface {
	// Position returns current position
	Position() Position

	// SetPosition updates position
	SetPosition(pos Position)

	// Move applies movement delta
	Move(dx, dy, dz int) Position

	// Facing returns orientation
	Facing() Facing

	// SetFacing updates orientation
	SetFacing(facing Facing)

	// LookAt rotates to face position
	LookAt(target Position)

	// DistanceTo calculates distance to another transform
	DistanceTo(other Transform) float64

	// IsAdjacent checks if adjacent to another transform
	IsAdjacent(other Transform) bool

	// InRange checks if within range of another transform
	InRange(other Transform, maxDistance float64) bool
}

// Movable represents entities that can move
type Movable interface {
	Transform

	// CanMoveTo checks if can move to position
	CanMoveTo(target Position) bool

	// MoveTo moves to target position
	MoveTo(target Position) error

	// MoveSpeed returns movement speed
	MoveSpeed() float64

	// SetMoveSpeed updates movement speed
	SetMoveSpeed(speed float64)

	// IsMoving returns true if currently moving
	IsMoving() bool

	// StopMovement halts movement
	StopMovement()
}

// Positionable represents anything with position
type Positionable interface {
	// Position returns current position
	Position() Position

	// SetPosition updates position
	SetPosition(pos Position)
}
