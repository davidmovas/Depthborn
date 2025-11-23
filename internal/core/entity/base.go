package entity

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/infra"
)

var _ Entity = (*BaseEntity)(nil)

type BaseEntity struct {
	infra.BasePersistent

	identity    types.Identity
	named       types.Named
	leveled     types.Leveled
	tagged      types.Tagged
	alive       types.Alive
	actionable  types.Actionable
	cloneable   types.Cloneable
	validatable types.Validatable

	attributes    AttributeManager
	statusEffects StatusManager
	transform     types.Transform
	callbacks     types.CallbackRegistry
}

func NewBaseEntity(id string, name string) *BaseEntity {
	entity := &BaseEntity{
		attributes:    NewBaseAttributeManager(),
		statusEffects: NewBaseStatusManager(),
		transform:     NewBaseTransform(),
		callbacks:     NewBaseCallbackRegistry(),
	}

	entity.identity = NewBaseIdentity(id, "entity")
	entity.named = NewBaseNamed(name)
	entity.leveled = NewBaseLeveled(1)
	entity.tagged = NewBaseTagged()
	entity.alive = NewBaseAlive(entity)
	entity.actionable = NewBaseActionable(true)
	entity.cloneable = NewBaseCloneable()
	entity.validatable = NewBaseValidatable()

	return entity
}

func (be *BaseEntity) ID() string {
	return be.identity.ID()
}

func (be *BaseEntity) Type() string {
	return be.identity.Type()
}

func (be *BaseEntity) Name() string {
	return be.named.Name()
}

func (be *BaseEntity) SetName(name string) {
	be.named.SetName(name)
}

func (be *BaseEntity) Level() int {
	return be.leveled.Level()
}

func (be *BaseEntity) SetLevel(level int) {
	be.leveled.SetLevel(level)
}

func (be *BaseEntity) Tags() types.TagSet {
	return be.tagged.Tags()
}

func (be *BaseEntity) IsAlive() bool {
	return be.alive.IsAlive()
}

func (be *BaseEntity) Kill(ctx context.Context, killerID string) error {
	return be.alive.Kill(ctx, killerID)
}

func (be *BaseEntity) Revive(ctx context.Context, healthPercent float64) error {
	return be.alive.Revive(ctx, healthPercent)
}

func (be *BaseEntity) CanAct() bool {
	return be.actionable.CanAct()
}

func (be *BaseEntity) Clone() interface{} {
	return be.cloneable.Clone()
}

func (be *BaseEntity) Validate() error {
	return be.validatable.Validate()
}

func (be *BaseEntity) Attributes() AttributeManager {
	return be.attributes
}

func (be *BaseEntity) StatusEffects() StatusManager {
	return be.statusEffects
}

func (be *BaseEntity) Transform() types.Transform {
	return be.transform
}

func (be *BaseEntity) Callbacks() types.CallbackRegistry {
	return be.callbacks
}
