package entity

import (
	"math"

	"github.com/davidmovas/Depthborn/internal/core/types"
)

type BaseTransform struct {
	x, y     float64
	rotation float64
}

func NewBaseTransform() *BaseTransform {
	return &BaseTransform{
		x:        0,
		y:        0,
		rotation: 0,
	}
}

func (bt *BaseTransform) Position() (x, y float64) {
	return bt.x, bt.y
}

func (bt *BaseTransform) SetPosition(x, y float64) {
	bt.x = x
	bt.y = y
}

func (bt *BaseTransform) Rotation() float64 {
	return bt.rotation
}

func (bt *BaseTransform) SetRotation(angle float64) {
	bt.rotation = angle
}

func (bt *BaseTransform) DistanceTo(other types.Transform) float64 {
	otherX, otherY := other.Position()
	dx := bt.x - otherX
	dy := bt.y - otherY
	return math.Sqrt(dx*dx + dy*dy)
}

func (bt *BaseTransform) LookAt(x, y float64) {
	dx := x - bt.x
	dy := y - bt.y
	bt.rotation = math.Atan2(dy, dx)
}
