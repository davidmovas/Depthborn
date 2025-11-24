package spatial

import (
	"errors"
	"sync"
)

var _ Movable = (*BaseMovable)(nil)

type BaseMovable struct {
	*BaseTransform
	mu sync.RWMutex

	moveSpeed float64
	isMoving  bool
	grid      Grid
}

func NewMovable(pos Position, facing Facing, speed float64, grid Grid) Movable {
	return &BaseMovable{
		BaseTransform: &BaseTransform{
			position: pos,
			facing:   facing,
		},
		moveSpeed: speed,
		isMoving:  false,
		grid:      grid,
	}
}

func (m *BaseMovable) CanMoveTo(target Position) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.grid == nil {
		return true
	}

	if !m.grid.IsValid(target) {
		return false
	}

	if !m.grid.IsWalkable(target) {
		return false
	}

	if m.grid.IsOccupied(target) {
		return false
	}

	return true
}

func (m *BaseMovable) MoveTo(target Position) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.grid != nil {
		if !m.grid.IsValid(target) {
			return errors.New("target position is out of bounds")
		}

		if !m.grid.IsWalkable(target) {
			return errors.New("target position is not walkable")
		}

		if m.grid.IsOccupied(target) {
			return errors.New("target position is occupied")
		}
	}

	m.position = target
	m.isMoving = false
	return nil
}

func (m *BaseMovable) MoveSpeed() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.moveSpeed
}

func (m *BaseMovable) SetMoveSpeed(speed float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.moveSpeed = speed
}

func (m *BaseMovable) IsMoving() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isMoving
}

func (m *BaseMovable) StopMovement() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isMoving = false
}
