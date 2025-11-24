package spatial

import "sync"

var _ Transform = (*BaseTransform)(nil)

type BaseTransform struct {
	mu sync.RWMutex

	position Position
	facing   Facing
}

func NewTransform(pos Position, facing Facing) Transform {
	return &BaseTransform{
		position: pos,
		facing:   facing,
	}
}

func (t *BaseTransform) Position() Position {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.position
}

func (t *BaseTransform) SetPosition(pos Position) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.position = pos
}

func (t *BaseTransform) Move(dx, dy, dz int) Position {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.position = t.position.Add(dx, dy, dz)
	return t.position
}

func (t *BaseTransform) Facing() Facing {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.facing
}

func (t *BaseTransform) SetFacing(facing Facing) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.facing = facing
}

func (t *BaseTransform) LookAt(target Position) {
	t.mu.Lock()
	defer t.mu.Unlock()

	angle := t.position.AngleTo(target)
	t.facing = Facing(angle)
}

func (t *BaseTransform) DistanceTo(other Transform) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.position.DistanceTo(other.Position())
}

func (t *BaseTransform) IsAdjacent(other Transform) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.position.IsAdjacent(other.Position())
}

func (t *BaseTransform) InRange(other Transform, maxDistance float64) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.position.InRange(other.Position(), maxDistance)
}
