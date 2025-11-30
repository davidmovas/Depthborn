// Package persist provides a clean, production-ready persistence layer
// for game state management with Unit of Work pattern.
package persist

import (
	"time"

	"github.com/davidmovas/Depthborn/pkg/persist/codec"
)

// EntityType identifies the type of entity for routing to correct repository.
type EntityType string

// Entity is the base interface for all game objects that can be persisted.
type Entity interface {
	// ID returns the unique identifier for this entity.
	ID() string

	// Type returns the entity type for routing.
	Type() EntityType
}

// Stateful extends Entity with change tracking capabilities.
type Stateful interface {
	Entity

	// Version returns the current version number (incremented on each save).
	Version() int64

	// IsDirty returns true if the entity has unsaved changes.
	IsDirty() bool

	// MarkDirty marks the entity as having unsaved changes.
	MarkDirty()

	// MarkClean resets the dirty flag (called after successful save).
	MarkClean()
}

// Timestamped adds creation and modification timestamps.
type Timestamped interface {
	// CreatedAt returns when the entity was first created.
	CreatedAt() time.Time

	// UpdatedAt returns when the entity was last modified.
	UpdatedAt() time.Time
}

// Persistable combines all interfaces needed for full persistence support.
type Persistable interface {
	Stateful
	Timestamped

	// SetVersion updates the version (called by repository on save).
	SetVersion(v int64)

	// Touch updates the UpdatedAt timestamp to now.
	Touch()
}

// Marshaler can serialize itself to bytes.
type Marshaler interface {
	// MarshalBinary serializes the entity to bytes.
	MarshalBinary() ([]byte, error)
}

// Unmarshaler can deserialize itself from bytes.
type Unmarshaler interface {
	// UnmarshalBinary deserializes the entity from bytes.
	UnmarshalBinary(data []byte) error
}

// Codec handles both marshaling and unmarshaling.
type Codec interface {
	Marshaler
	Unmarshaler
}

// DefaultCodec returns the default codec (MessagePack).
func DefaultCodec() interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, target any) error
} {
	return defaultCodec
}

var defaultCodec = &msgpackCodec{}

type msgpackCodec struct{}

func (c *msgpackCodec) Encode(v any) ([]byte, error) {
	return codec.Default.Encode(v)
}

func (c *msgpackCodec) Decode(data []byte, target any) error {
	return codec.Default.Decode(data, target)
}
