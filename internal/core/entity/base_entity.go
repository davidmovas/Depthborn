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
	"github.com/davidmovas/Depthborn/pkg/persist"
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
	clone := &BaseEntity{
		name:    e.name,
		isAlive: e.isAlive,
	}

	// Deep clone tags
	if e.tags != nil {
		clone.tags = types.NewTagSet()
		for _, tag := range e.tags.All() {
			clone.tags.Add(tag)
		}
	}

	// Deep clone attributes
	if e.attributes != nil {
		clone.attributes = attribute.NewManager()
		snapshot := e.attributes.Snapshot()
		clone.attributes.Restore(snapshot)
	}

	// Deep clone status effects (create fresh manager - effects are transient)
	if e.statuses != nil {
		clone.statuses = status.NewManager()
		// Note: Status effects are typically not cloned as they are transient
		// and reference-dependent. A fresh manager is provided.
	}

	// Deep clone transform
	if e.transform != nil {
		pos := e.transform.Position()
		facing := e.transform.Facing()
		clone.transform = spatial.NewTransform(pos, facing)
	}

	// Create fresh callback registry (callbacks contain closures that reference
	// original entity, so we create a new empty registry)
	clone.callbacks = types.NewCallbackRegistry()

	// Clone BasePersistent with new ID
	clone.BasePersistent = impl.NewPersistent(e.Type(), clone, nil)

	// Note: progression is not cloned as it typically references external managers
	// and would require special handling. Set to nil for cloned entities.
	clone.progression = nil

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

// EntityState holds the complete serializable state of a BaseEntity.
// Used for MarshalBinary/UnmarshalBinary integration with pkg/persist.
type EntityState struct {
	ID         string             `msgpack:"id"`
	EntityType string             `msgpack:"entity_type"`
	Version    int64              `msgpack:"version"`
	CreatedAt  int64              `msgpack:"created_at"`
	UpdatedAt  int64              `msgpack:"updated_at"`
	Name       string             `msgpack:"name"`
	IsAlive    bool               `msgpack:"is_alive"`
	Attributes map[string]float64 `msgpack:"attributes,omitempty"`
	Tags       []string           `msgpack:"tags,omitempty"`
	Position   *PositionState     `msgpack:"position,omitempty"`
}

// PositionState holds spatial position data.
type PositionState struct {
	X      int     `msgpack:"x"`
	Y      int     `msgpack:"y"`
	Z      int     `msgpack:"z"`
	Facing float64 `msgpack:"facing"`
}

// MarshalBinary implements persist.Marshaler for efficient binary serialization.
func (e *BaseEntity) MarshalBinary() ([]byte, error) {
	es := EntityState{
		ID:         e.ID(),
		EntityType: e.Type(),
		Version:    e.Version(),
		CreatedAt:  e.CreatedAt(),
		UpdatedAt:  e.UpdatedAt(),
		Name:       e.name,
		IsAlive:    e.isAlive,
	}

	// Serialize attributes
	if e.attributes != nil {
		snapshot := e.attributes.Snapshot()
		es.Attributes = make(map[string]float64, len(snapshot))
		for k, v := range snapshot {
			es.Attributes[string(k)] = v
		}
	}

	// Serialize tags
	if e.tags != nil {
		es.Tags = e.tags.All()
	}

	// Serialize transform/position
	if e.transform != nil {
		pos := e.transform.Position()
		es.Position = &PositionState{
			X:      pos.X,
			Y:      pos.Y,
			Z:      pos.Z,
			Facing: float64(e.transform.Facing()),
		}
	}

	return persist.DefaultCodec().Encode(es)
}

// UnmarshalBinary implements persist.Unmarshaler for efficient binary deserialization.
func (e *BaseEntity) UnmarshalBinary(data []byte) error {
	var es EntityState
	if err := persist.DefaultCodec().Decode(data, &es); err != nil {
		return fmt.Errorf("failed to decode entity state: %w", err)
	}

	// Restore base persistent fields
	e.BasePersistent = impl.NewPersistentWithID(es.ID, es.EntityType, e, nil)

	// Restore entity fields
	e.name = es.Name
	e.isAlive = es.IsAlive

	// Restore attributes
	if e.attributes != nil && len(es.Attributes) > 0 {
		snapshot := make(map[attribute.Type]float64, len(es.Attributes))
		for k, v := range es.Attributes {
			snapshot[attribute.Type(k)] = v
		}
		e.attributes.Restore(snapshot)
	}

	// Restore tags
	if e.tags != nil && len(es.Tags) > 0 {
		e.tags.Clear()
		for _, tag := range es.Tags {
			e.tags.Add(tag)
		}
	}

	// Restore transform/position
	if e.transform != nil && es.Position != nil {
		pos := spatial.NewPosition(es.Position.X, es.Position.Y, es.Position.Z)
		e.transform.SetPosition(pos)
		e.transform.SetFacing(spatial.Facing(es.Position.Facing))
	}

	return nil
}
