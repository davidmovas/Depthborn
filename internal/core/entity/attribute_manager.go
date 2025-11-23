package entity

import "github.com/davidmovas/Depthborn/internal/core/attribute"

var _ AttributeManager = (*BaseAttributeManager)(nil)

type BaseAttributeManager struct {
	manager attribute.Manager
}

func NewBaseAttributeManager() AttributeManager {
	return &BaseAttributeManager{
		manager: attribute.NewAttributeManager(),
	}
}

func (bam *BaseAttributeManager) Get(attrType attribute.Type) float64 {
	return bam.manager.Get(attrType)
}

func (bam *BaseAttributeManager) GetBase(attrType attribute.Type) float64 {
	return bam.manager.GetBase(attrType)
}

func (bam *BaseAttributeManager) SetBase(attrType attribute.Type, value float64) {
	bam.manager.SetBase(attrType, value)
}

func (bam *BaseAttributeManager) AddModifier(attrType attribute.Type, modifier AttributeModifier) {
	attrModifier := attribute.NewModifier(
		modifier.ID(),
		attribute.ModifierType(modifier.Type()),
		modifier.Value(),
		modifier.Source(),
		0, // priority
	)
	bam.manager.AddModifier(attrType, attrModifier)
}

func (bam *BaseAttributeManager) RemoveModifier(attrType attribute.Type, modifierID string) {
	bam.manager.RemoveModifier(attrType, modifierID)
}

func (bam *BaseAttributeManager) RecalculateAll() {
	bam.manager.RecalculateAll()
}
