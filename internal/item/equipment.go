package item

import (
	"context"
	"fmt"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/item/affix"
)

var _ Equipment = (*BaseEquipment)(nil)

type BaseEquipment struct {
	*BaseItem
	slot          EquipmentSlot
	attributes    []attribute.Modifier
	durability    float64
	maxDurability float64
	sockets       []Socketable
	affixSet      affix.Set
	requirements  EquipRequirements
}

func NewBaseEquipment(id string, itemType Type, name string, slot EquipmentSlot) *BaseEquipment {
	return &BaseEquipment{
		BaseItem:      NewBaseItem(id, itemType, name),
		slot:          slot,
		attributes:    make([]attribute.Modifier, 0),
		durability:    100.0,
		maxDurability: 100.0,
		sockets:       make([]Socketable, 0),
		affixSet:      nil, // TODO: Initialize proper affix set
		requirements:  NewSimpleRequirements(1, make(map[attribute.Type]float64)),
	}
}

func (be *BaseEquipment) Slot() EquipmentSlot {
	return be.slot
}

func (be *BaseEquipment) Attributes() []attribute.Modifier {
	return be.attributes
}

func (be *BaseEquipment) Durability() float64 {
	return be.durability
}

func (be *BaseEquipment) MaxDurability() float64 {
	return be.maxDurability
}

func (be *BaseEquipment) SetDurability(value float64) {
	if value < 0 {
		value = 0
	} else if value > be.maxDurability {
		value = be.maxDurability
	}
	be.durability = value
}

func (be *BaseEquipment) Repair(amount float64) {
	be.durability += amount
	if be.durability > be.maxDurability {
		be.durability = be.maxDurability
	}
}

func (be *BaseEquipment) DamageItem(amount float64) bool {
	be.durability -= amount
	if be.durability < 0 {
		be.durability = 0
	}
	return !be.IsBroken()
}

func (be *BaseEquipment) IsBroken() bool {
	return be.durability <= 0
}

func (be *BaseEquipment) SocketCount() int {
	return len(be.sockets)
}

func (be *BaseEquipment) GetSocket(index int) (Socketable, bool) {
	if index < 0 || index >= len(be.sockets) {
		return nil, false
	}
	socket := be.sockets[index]
	return socket, socket != nil
}

func (be *BaseEquipment) SetSocket(index int, item Socketable) error {
	if index < 0 || index >= len(be.sockets) {
		return fmt.Errorf("socket index out of range: %d", index)
	}

	if item != nil && be.sockets[index] != nil {
		return fmt.Errorf("socket already occupied")
	}

	be.sockets[index] = item
	return nil
}

func (be *BaseEquipment) RemoveSocket(index int) (Socketable, error) {
	if index < 0 || index >= len(be.sockets) {
		return nil, fmt.Errorf("socket index out of range: %d", index)
	}

	item := be.sockets[index]
	be.sockets[index] = nil
	return item, nil
}

func (be *BaseEquipment) Affixes() affix.Set {
	return be.affixSet
}

func (be *BaseEquipment) Requirements() EquipRequirements {
	return be.requirements
}

func (be *BaseEquipment) CanEquip(entity entity.Entity) bool {
	return be.requirements.Check(entity)
}

func (be *BaseEquipment) OnEquip(_ context.Context, _ entity.Entity) error {
	// TODO: Implement equip logic
	return nil
}

func (be *BaseEquipment) OnUnequip(_ context.Context, _ entity.Entity) error {
	// TODO: Implement unequip logic
	return nil
}
