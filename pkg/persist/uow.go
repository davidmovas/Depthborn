package persist

import (
	"context"
	"errors"
	"sync"

	"github.com/davidmovas/Depthborn/pkg/persist/codec"
	"github.com/davidmovas/Depthborn/pkg/persist/storage"
)

// UnitOfWork errors.
var (
	ErrUoWClosed    = errors.New("unit of work is closed")
	ErrNoFactory    = errors.New("no factory registered for entity type")
	ErrCommitFailed = errors.New("commit failed")
)

// UnitOfWork tracks changes to entities and commits them atomically.
// It implements the Unit of Work pattern for managing entity lifecycle.
//
// Usage:
//
//	uow := persist.NewUnitOfWork(storage)
//	defer uow.Close()
//
//	// Load or create entities
//	char, _ := uow.Get(EntityCharacter, "char-1")
//	char.Level++
//
//	// All dirty entities are saved on commit
//	uow.Commit(ctx)
type UnitOfWork struct {
	storage   storage.Storage
	codec     codec.Codec
	factories map[EntityType]func() Persistable

	// Identity map: tracks all loaded entities
	entities map[string]Persistable // key: "type:id"

	// Entities scheduled for deletion
	deleted map[string]bool

	closed bool
	mu     sync.RWMutex
}

// NewUnitOfWork creates a new unit of work.
func NewUnitOfWork(store storage.Storage) *UnitOfWork {
	return &UnitOfWork{
		storage:   store,
		codec:     codec.Default,
		factories: make(map[EntityType]func() Persistable),
		entities:  make(map[string]Persistable),
		deleted:   make(map[string]bool),
	}
}

// WithCodec sets a custom codec.
func (u *UnitOfWork) WithCodec(c codec.Codec) *UnitOfWork {
	u.codec = c
	return u
}

// RegisterFactory registers a factory for creating entities of the given type.
// Must be called before Get() can load entities of that type.
func (u *UnitOfWork) RegisterFactory(entityType EntityType, factory func() Persistable) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.factories[entityType] = factory
}

// Get retrieves an entity by type and ID.
// If the entity is already loaded, returns the cached instance.
// Otherwise, loads from storage.
func (u *UnitOfWork) Get(ctx context.Context, entityType EntityType, id string) (Persistable, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		return nil, ErrUoWClosed
	}

	key := makeKey(entityType, id)

	// Check if deleted
	if u.deleted[key] {
		return nil, ErrEntityNotFound
	}

	// Check identity map
	if entity, exists := u.entities[key]; exists {
		return entity, nil
	}

	// Load from storage
	factory, ok := u.factories[entityType]
	if !ok {
		return nil, ErrNoFactory
	}

	storageKey := storage.EntityKey(string(entityType), id)
	record, err := u.storage.Get(ctx, storageKey)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrEntityNotFound
		}
		return nil, err
	}

	entity := factory()
	if err := u.codec.Decode(record.Data, entity); err != nil {
		return nil, err
	}

	entity.MarkClean()
	u.entities[key] = entity

	return entity, nil
}

// Register adds a new or modified entity to the unit of work.
// The entity will be saved when Commit is called.
func (u *UnitOfWork) Register(entity Persistable) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		return
	}

	key := makeKey(entity.Type(), entity.ID())
	u.entities[key] = entity
	delete(u.deleted, key) // Un-delete if was marked
}

// Delete marks an entity for deletion.
// The entity will be removed from storage when Commit is called.
func (u *UnitOfWork) Delete(entity Persistable) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		return
	}

	key := makeKey(entity.Type(), entity.ID())
	delete(u.entities, key)
	u.deleted[key] = true
}

// DeleteByID marks an entity for deletion by type and ID.
func (u *UnitOfWork) DeleteByID(entityType EntityType, id string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		return
	}

	key := makeKey(entityType, id)
	delete(u.entities, key)
	u.deleted[key] = true
}

// HasChanges returns true if there are unsaved changes.
func (u *UnitOfWork) HasChanges() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()

	if len(u.deleted) > 0 {
		return true
	}

	for _, entity := range u.entities {
		if entity.IsDirty() {
			return true
		}
	}

	return false
}

// Commit saves all changes to storage atomically.
func (u *UnitOfWork) Commit(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		return ErrUoWClosed
	}

	// Try to use transaction if available
	txStore, hasTx := u.storage.(storage.TransactionalStorage)
	if hasTx {
		return u.commitWithTransaction(ctx, txStore)
	}

	return u.commitWithoutTransaction(ctx)
}

func (u *UnitOfWork) commitWithTransaction(ctx context.Context, txStore storage.TransactionalStorage) error {
	tx, err := txStore.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Process deletions
	for key := range u.deleted {
		entityType, id := parseKey(key)
		storageKey := storage.EntityKey(entityType, id)
		if err := tx.Delete(ctx, storageKey); err != nil {
			return err
		}
	}

	// Process dirty entities
	for _, entity := range u.entities {
		if !entity.IsDirty() {
			continue
		}

		newVersion := entity.Version() + 1
		entity.SetVersion(newVersion)
		entity.Touch()

		data, err := u.codec.Encode(entity)
		if err != nil {
			return err
		}

		record := &storage.Record{
			Key:       storage.EntityKey(string(entity.Type()), entity.ID()),
			Data:      data,
			Version:   newVersion,
			CreatedAt: entity.CreatedAt().Unix(),
			UpdatedAt: entity.UpdatedAt().Unix(),
		}

		if err := tx.Set(ctx, record); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Mark all as clean after successful commit
	for _, entity := range u.entities {
		entity.MarkClean()
	}
	u.deleted = make(map[string]bool)

	return nil
}

func (u *UnitOfWork) commitWithoutTransaction(ctx context.Context) error {
	// Process deletions first
	for key := range u.deleted {
		entityType, id := parseKey(key)
		storageKey := storage.EntityKey(entityType, id)
		if err := u.storage.Delete(ctx, storageKey); err != nil {
			return err
		}
	}

	// Process dirty entities
	for _, entity := range u.entities {
		if !entity.IsDirty() {
			continue
		}

		newVersion := entity.Version() + 1
		entity.SetVersion(newVersion)
		entity.Touch()

		data, err := u.codec.Encode(entity)
		if err != nil {
			return err
		}

		record := &storage.Record{
			Key:       storage.EntityKey(string(entity.Type()), entity.ID()),
			Data:      data,
			Version:   newVersion,
			CreatedAt: entity.CreatedAt().Unix(),
			UpdatedAt: entity.UpdatedAt().Unix(),
		}

		if err := u.storage.Set(ctx, record); err != nil {
			return err
		}

		entity.MarkClean()
	}

	u.deleted = make(map[string]bool)
	return nil
}

// Rollback discards all changes and clears the unit of work.
func (u *UnitOfWork) Rollback() {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.entities = make(map[string]Persistable)
	u.deleted = make(map[string]bool)
}

// Clear removes all tracked entities without saving.
// Useful for starting fresh within the same UoW.
func (u *UnitOfWork) Clear() {
	u.Rollback()
}

// Close closes the unit of work.
// Any uncommitted changes are discarded.
func (u *UnitOfWork) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.closed = true
	u.entities = nil
	u.deleted = nil
	return nil
}

// TrackedCount returns the number of tracked entities.
func (u *UnitOfWork) TrackedCount() int {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return len(u.entities)
}

// DirtyCount returns the number of dirty entities.
func (u *UnitOfWork) DirtyCount() int {
	u.mu.RLock()
	defer u.mu.RUnlock()

	count := 0
	for _, entity := range u.entities {
		if entity.IsDirty() {
			count++
		}
	}
	return count
}

// --- Helpers ---

func makeKey(entityType EntityType, id string) string {
	return string(entityType) + ":" + id
}

func parseKey(key string) (entityType string, id string) {
	for i := 0; i < len(key); i++ {
		if key[i] == ':' {
			return key[:i], key[i+1:]
		}
	}
	return key, ""
}
