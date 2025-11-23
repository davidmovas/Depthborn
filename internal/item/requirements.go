package item

import (
	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
)

type SimpleRequirements struct {
	level      int
	attributes map[attribute.Type]float64
}

func NewSimpleRequirements(level int, attributes map[attribute.Type]float64) EquipRequirements {
	return &SimpleRequirements{
		level:      level,
		attributes: attributes,
	}
}

func (sr *SimpleRequirements) Level() int {
	return sr.level
}

func (sr *SimpleRequirements) Attributes() map[attribute.Type]float64 {
	return sr.attributes
}

func (sr *SimpleRequirements) Check(entity entity.Entity) bool {
	if entity.Level() < sr.level {
		return false
	}

	attrManager := entity.Attributes()
	for attrType, minValue := range sr.attributes {
		if attrManager.Get(attrType) < minValue {
			return false
		}
	}

	return true
}
