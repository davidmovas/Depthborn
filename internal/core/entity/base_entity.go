package entity

import (
	"context"
	"fmt"

	"github.com/davidmovas/Depthborn/internal/character/progression"
	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/status"
	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/infra/impl"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
	"github.com/davidmovas/Depthborn/pkg/state"
)

var _ Entity = (*BaseEntity)(nil)

type BaseEntity struct {
	*impl.BasePersistent

	name string

	// Components
	tags        types.TagSet
	attributes  attribute.Manager
	statuses    status.Manager
	transform   spatial.Transform
	callbacks   types.CallbackRegistry
	progression progression.Manager

	// Alive state
	isAlive bool
}

type Config struct {
	Name               string
	EntityType         string
	AttributeManager   attribute.Manager
	StatusManager      status.Manager
	Transform          spatial.Transform
	TagSet             types.TagSet
	Callbacks          types.CallbackRegistry
	ProgressionManager progression.Manager
}

func NewEntity(config Config) *BaseEntity {
	entity := &BaseEntity{
		name:        config.Name,
		attributes:  config.AttributeManager,
		statuses:    config.StatusManager,
		transform:   config.Transform,
		tags:        config.TagSet,
		callbacks:   config.Callbacks,
		progression: config.ProgressionManager,
		isAlive:     true,
	}

	entity.BasePersistent = impl.NewPersistent(config.EntityType, entity, nil)

	return entity
}

func (e *BaseEntity) Name() string {
	return e.name
}

func (e *BaseEntity) SetName(name string) {
	e.name = name
	e.Touch()
}

func (e *BaseEntity) Level() int {
	if e.progression != nil && e.progression.Experience() != nil {
		return e.progression.Experience().CurrentLevel()
	}
	return 1
}

func (e *BaseEntity) SetLevel(level int) {
	if e.progression != nil && e.progression.Experience() != nil {
		_ = e.progression.Experience().SetLevel(level)
		e.Touch()
	}
}

func (e *BaseEntity) Tags() types.TagSet {
	return e.tags
}

func (e *BaseEntity) IsAlive() bool {
	return e.isAlive
}

func (e *BaseEntity) Kill(ctx context.Context, killerID string) error {
	if !e.isAlive {
		return fmt.Errorf("entity is already dead")
	}

	e.isAlive = false
	e.Touch()

	// Trigger death callback
	e.callbacks.TriggerDeath(ctx, e.ID(), killerID)

	return nil
}

func (e *BaseEntity) Revive(_ context.Context, healthPercent float64) error {
	if e.isAlive {
		return fmt.Errorf("entity is already alive")
	}

	if healthPercent < 0 || healthPercent > 1 {
		return fmt.Errorf("health percent must be between 0 and 1")
	}

	e.isAlive = true
	e.Touch()

	// Subclasses (Living) should restore health based on healthPercent
	return nil
}

func (e *BaseEntity) CanAct() bool {
	if !e.isAlive {
		return false
	}

	// Check for controlling status effects
	if e.statuses.Has("stun") ||
		e.statuses.Has("freeze") ||
		e.statuses.Has("sleep") ||
		e.statuses.Has("petrify") {
		return false
	}

	return true
}

func (e *BaseEntity) Attributes() attribute.Manager {
	return e.attributes
}

func (e *BaseEntity) StatusEffects() status.Manager {
	return e.statuses
}

func (e *BaseEntity) Transform() spatial.Transform {
	return e.transform
}

func (e *BaseEntity) Callbacks() types.CallbackRegistry {
	return e.callbacks
}

func (e *BaseEntity) Progression() progression.Manager {
	return e.progression
}

func (e *BaseEntity) Clone() any {
	// TODO: Implement deep cloning
	// Need to clone:
	// - tags
	// - attributes (create new manager with same values)
	// - statuses (create new manager, copy active effects)
	// - transform (depends on implementation)
	// - callbacks (probably shouldn't clone, create new registry)
	// - progression (clone state)

	clone := &BaseEntity{
		name:    e.name,
		isAlive: e.isAlive,
		// TODO: Deep clone all components
	}

	return clone
}

func (e *BaseEntity) Validate() error {
	if e.name == "" {
		return fmt.Errorf("entity must have a name")
	}

	if e.attributes == nil {
		return fmt.Errorf("entity must have attribute manager")
	}

	if e.statuses == nil {
		return fmt.Errorf("entity must have status manager")
	}

	if e.transform == nil {
		return fmt.Errorf("entity must have transform")
	}

	if e.tags == nil {
		return fmt.Errorf("entity must have tag set")
	}

	if e.callbacks == nil {
		return fmt.Errorf("entity must have callback registry")
	}

	return nil
}

func (e *BaseEntity) SerializeState() (map[string]any, error) {
	s := state.New().
		Set("name", e.name).
		Set("is_alive", e.isAlive)

	if e.attributes != nil {
		s.Set("attributes", e.attributes.Snapshot())
	}

	if e.tags != nil {
		s.Set("tags", e.tags.All())
	}

	if e.transform != nil {
		if err := s.SetEntity("transform", e.transform); err != nil {
			return nil, fmt.Errorf("failed to serialize transform: %w", err)
		}
	}

	return s.Data(), nil
}

func (e *BaseEntity) DeserializeState(stateData map[string]any) error {
	s := state.From(stateData)

	e.name = s.StringOr("name", "!NONE!")
	e.isAlive = s.BoolOr("is_alive", true)

	if e.attributes != nil {
		if attrsState, ok := s.Map("attributes"); ok {
			snapshot := make(map[attribute.Type]float64)
			for _, key := range attrsState.Keys() {
				if value, is := attrsState.Float(key); is {
					snapshot[attribute.Type(key)] = value
				}
			}
			e.attributes.Restore(snapshot)
		}
	}

	if e.tags != nil {
		if tags, ok := state.SliceTyped[string]("tags", s); ok {
			e.tags.Clear()
			for _, tag := range tags {
				e.tags.Add(tag)
			}
		}
	}

	if e.transform != nil {
		if err := s.GetEntity("transform", e.transform); err != nil {
			return fmt.Errorf("failed to deserialize transform: %w", err)
		}
	}

	return nil
}
