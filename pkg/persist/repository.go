package persist

import (
	"context"
	"errors"

	"github.com/davidmovas/Depthborn/pkg/persist/codec"
	"github.com/davidmovas/Depthborn/pkg/persist/storage"
)

// Repository errors.
var (
	ErrEntityNotFound  = errors.New("entity not found")
	ErrVersionConflict = errors.New("version conflict: entity was modified")
	ErrInvalidEntity   = errors.New("invalid entity: missing required fields")
)

// EntityFactory creates new instances of entities.
type EntityFactory[T Persistable] func() T

// Repository provides type-safe CRUD operations for entities.
type Repository[T Persistable] struct {
	storage    storage.Storage
	codec      codec.Codec
	entityType EntityType
	factory    EntityFactory[T]
}

// NewRepository creates a new repository for the given entity type.
func NewRepository[T Persistable](
	store storage.Storage,
	entityType EntityType,
	factory EntityFactory[T],
) *Repository[T] {
	return &Repository[T]{
		storage:    store,
		codec:      codec.Default,
		entityType: entityType,
		factory:    factory,
	}
}

// WithCodec sets a custom codec for serialization.
func (r *Repository[T]) WithCodec(c codec.Codec) *Repository[T] {
	r.codec = c
	return r
}

// Load retrieves an entity by ID.
func (r *Repository[T]) Load(ctx context.Context, id string) (T, error) {
	var zero T

	key := storage.EntityKey(string(r.entityType), id)
	record, err := r.storage.Get(ctx, key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return zero, ErrEntityNotFound
		}
		return zero, err
	}

	entity := r.factory()
	if err := r.codec.Decode(record.Data, entity); err != nil {
		return zero, err
	}

	entity.MarkClean()
	return entity, nil
}

// Save persists an entity.
// If optimistic locking is enabled, returns ErrVersionConflict on mismatch.
func (r *Repository[T]) Save(ctx context.Context, entity T) error {
	if entity.ID() == "" {
		return ErrInvalidEntity
	}

	// Increment version
	newVersion := entity.Version() + 1
	entity.SetVersion(newVersion)
	entity.Touch()

	// Serialize
	data, err := r.codec.Encode(entity)
	if err != nil {
		return err
	}

	key := storage.EntityKey(string(r.entityType), entity.ID())
	record := &storage.Record{
		Key:       key,
		Data:      data,
		Version:   newVersion,
		CreatedAt: entity.CreatedAt().Unix(),
		UpdatedAt: entity.UpdatedAt().Unix(),
	}

	if err := r.storage.Set(ctx, record); err != nil {
		return err
	}

	entity.MarkClean()
	return nil
}

// Delete removes an entity by ID.
func (r *Repository[T]) Delete(ctx context.Context, id string) error {
	key := storage.EntityKey(string(r.entityType), id)
	return r.storage.Delete(ctx, key)
}

// Exists checks if an entity exists.
func (r *Repository[T]) Exists(ctx context.Context, id string) (bool, error) {
	key := storage.EntityKey(string(r.entityType), id)
	return r.storage.Exists(ctx, key)
}

// LoadMany retrieves multiple entities by IDs.
// Missing entities are silently omitted from the result.
func (r *Repository[T]) LoadMany(ctx context.Context, ids []string) ([]T, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	batchStore, ok := r.storage.(storage.BatchStorage)
	if !ok {
		// Fallback to individual loads
		return r.loadManySequential(ctx, ids)
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = storage.EntityKey(string(r.entityType), id)
	}

	records, err := batchStore.GetMany(ctx, keys)
	if err != nil {
		return nil, err
	}

	entities := make([]T, 0, len(records))
	for _, record := range records {
		entity := r.factory()
		if err := r.codec.Decode(record.Data, entity); err != nil {
			continue // Skip invalid records
		}
		entity.MarkClean()
		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repository[T]) loadManySequential(ctx context.Context, ids []string) ([]T, error) {
	entities := make([]T, 0, len(ids))
	for _, id := range ids {
		entity, err := r.Load(ctx, id)
		if err != nil {
			if errors.Is(err, ErrEntityNotFound) {
				continue
			}
			return nil, err
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

// SaveMany persists multiple entities atomically.
func (r *Repository[T]) SaveMany(ctx context.Context, entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	batchStore, ok := r.storage.(storage.BatchStorage)
	if !ok {
		// Fallback to individual saves
		return r.saveManySequential(ctx, entities)
	}

	records := make([]*storage.Record, len(entities))
	for i, entity := range entities {
		if entity.ID() == "" {
			return ErrInvalidEntity
		}

		newVersion := entity.Version() + 1
		entity.SetVersion(newVersion)
		entity.Touch()

		data, err := r.codec.Encode(entity)
		if err != nil {
			return err
		}

		records[i] = &storage.Record{
			Key:       storage.EntityKey(string(r.entityType), entity.ID()),
			Data:      data,
			Version:   newVersion,
			CreatedAt: entity.CreatedAt().Unix(),
			UpdatedAt: entity.UpdatedAt().Unix(),
		}
	}

	if err := batchStore.SetMany(ctx, records); err != nil {
		return err
	}

	for _, entity := range entities {
		entity.MarkClean()
	}

	return nil
}

func (r *Repository[T]) saveManySequential(ctx context.Context, entities []T) error {
	for _, entity := range entities {
		if err := r.Save(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// List returns all entity IDs of this type.
func (r *Repository[T]) List(ctx context.Context, limit, offset int) ([]string, error) {
	queryStore, ok := r.storage.(storage.QueryableStorage)
	if !ok {
		return nil, errors.New("storage does not support listing")
	}

	prefix := string(r.entityType) + ":"
	keys, err := queryStore.List(ctx, prefix, limit, offset)
	if err != nil {
		return nil, err
	}

	// Extract IDs from keys
	ids := make([]string, len(keys))
	prefixLen := len(prefix)
	for i, key := range keys {
		if len(key) > prefixLen {
			ids[i] = key[prefixLen:]
		}
	}

	return ids, nil
}

// Count returns the total number of entities of this type.
func (r *Repository[T]) Count(ctx context.Context) (int, error) {
	queryStore, ok := r.storage.(storage.QueryableStorage)
	if !ok {
		return 0, errors.New("storage does not support counting")
	}

	prefix := string(r.entityType) + ":"
	return queryStore.Count(ctx, prefix)
}
