package persist

import (
	"sync"
	"time"

	"github.com/jaevor/go-nanoid"
)

// idGenerator creates short unique IDs.
var idGenerator func() string

func init() {
	gen, err := nanoid.Standard(21)
	if err != nil {
		panic("failed to create nanoid generator: " + err.Error())
	}
	idGenerator = gen
}

// NewID generates a new unique identifier.
func NewID() string {
	return idGenerator()
}

// Base provides common functionality for all persistable entities.
// Embed this in your entity structs to get automatic ID, version, and timestamp management.
//
// Example:
//
//	type Character struct {
//	    persist.Base
//	    Name  string
//	    Level int
//	}
type Base struct {
	id        string
	typ       EntityType
	version   int64
	createdAt time.Time
	updatedAt time.Time
	dirty     bool
	mu        sync.RWMutex
}

// NewBase creates a new Base with generated ID and current timestamps.
func NewBase(entityType EntityType) Base {
	now := time.Now()
	return Base{
		id:        NewID(),
		typ:       entityType,
		version:   0,
		createdAt: now,
		updatedAt: now,
		dirty:     true, // New entities are dirty
	}
}

// NewBaseWithID creates a new Base with a specific ID.
func NewBaseWithID(id string, entityType EntityType) Base {
	now := time.Now()
	return Base{
		id:        id,
		typ:       entityType,
		version:   0,
		createdAt: now,
		updatedAt: now,
		dirty:     true,
	}
}

// ID returns the entity's unique identifier.
func (b *Base) ID() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.id
}

// Type returns the entity type.
func (b *Base) Type() EntityType {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.typ
}

// Version returns the current version number.
func (b *Base) Version() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.version
}

// SetVersion updates the version number.
func (b *Base) SetVersion(v int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.version = v
}

// CreatedAt returns when the entity was created.
func (b *Base) CreatedAt() time.Time {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.createdAt
}

// UpdatedAt returns when the entity was last modified.
func (b *Base) UpdatedAt() time.Time {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.updatedAt
}

// Touch updates the modification timestamp to now.
func (b *Base) Touch() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.updatedAt = time.Now()
	b.dirty = true
}

// IsDirty returns true if the entity has unsaved changes.
func (b *Base) IsDirty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dirty
}

// MarkDirty marks the entity as having unsaved changes.
func (b *Base) MarkDirty() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dirty = true
	b.updatedAt = time.Now()
}

// MarkClean resets the dirty flag after successful save.
func (b *Base) MarkClean() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dirty = false
}

// BaseState holds the serializable state of Base.
// Use this for embedding in your entity's state struct.
type BaseState struct {
	ID        string    `msgpack:"id" json:"id"`
	Type      string    `msgpack:"type" json:"type"`
	Version   int64     `msgpack:"version" json:"version"`
	CreatedAt time.Time `msgpack:"created_at" json:"created_at"`
	UpdatedAt time.Time `msgpack:"updated_at" json:"updated_at"`
}

// State returns the serializable state of the base fields.
func (b *Base) State() BaseState {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return BaseState{
		ID:        b.id,
		Type:      string(b.typ),
		Version:   b.version,
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
	}
}

// LoadState restores base fields from serialized state.
func (b *Base) LoadState(s BaseState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.id = s.ID
	b.typ = EntityType(s.Type)
	b.version = s.Version
	b.createdAt = s.CreatedAt
	b.updatedAt = s.UpdatedAt
	b.dirty = false // Just loaded, so clean
}
