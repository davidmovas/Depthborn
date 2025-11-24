package store

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/davidmovas/Depthborn/internal/infra"
	"github.com/davidmovas/Depthborn/internal/infra/registry"
	"github.com/davidmovas/Depthborn/internal/persistence"
	"github.com/davidmovas/Depthborn/internal/persistence/store/sqlite"
	"github.com/davidmovas/Depthborn/pkg/dbx"
)

var _ persistence.Repository = (*Repository)(nil)

type Repository struct {
	db            *sqlite.DB
	snapshotStore persistence.SnapshotStore
	deltaStore    persistence.DeltaStore
	registry      registry.Registry
	strategy      SnapshotStrategy
}

type RepositoryConfig struct {
	SnapshotStore persistence.SnapshotStore
	DeltaStore    persistence.DeltaStore
	Registry      registry.Registry
	Strategy      SnapshotStrategy
}

func NewRepository(config RepositoryConfig) *Repository {
	if config.Strategy == nil {
		config.Strategy = NewEveryNVersionsStrategy(10)
	}

	return &Repository{
		snapshotStore: config.SnapshotStore,
		deltaStore:    config.DeltaStore,
		registry:      config.Registry,
		strategy:      config.Strategy,
	}
}

func (r *Repository) Save(ctx context.Context, entity infra.Persistent) error {
	version := entity.Version()

	// Check if we should create a snapshot
	shouldSnapshot := r.strategy.ShouldSnapshot(version, 0, 0)

	if shouldSnapshot {
		// Save full snapshot
		if err := r.snapshotStore.SaveSnapshot(ctx, entity); err != nil {
			return fmt.Errorf("failed to save snapshot: %w", err)
		}

		// Compact old deltas
		if sqliteStore, ok := r.deltaStore.(*sqlite.DeltaStore); ok {
			if err := sqliteStore.CompactDeltas(ctx, entity.Type(), entity.ID(), version); err != nil {
				// Log error but don't fail the save
				_ = err
			}
		}
	} else {
		// Save delta if version > 1
		if version > 1 {
			deltaData, err := entity.Delta(version - 1)
			if err != nil {
				return fmt.Errorf("failed to generate delta: %w", err)
			}

			if err = r.deltaStore.SaveDelta(ctx, entity.Type(), entity.ID(), version-1, version, deltaData); err != nil {
				return fmt.Errorf("failed to save delta: %w", err)
			}
		} else {
			// First version always gets a snapshot
			if err := r.snapshotStore.SaveSnapshot(ctx, entity); err != nil {
				return fmt.Errorf("failed to save initial snapshot: %w", err)
			}
		}
	}

	return nil
}

func (r *Repository) Load(ctx context.Context, id, entityType string) (infra.Persistent, error) {
	// Try to load from snapshot first
	entity, err := r.snapshotStore.LoadSnapshot(ctx, id, entityType)
	if err != nil {
		return nil, fmt.Errorf("failed to load snapshot: %w", err)
	}

	persistent, ok := entity.(infra.Persistent)
	if !ok {
		return nil, fmt.Errorf("entity %s/%s does not implement Persistent", id, entityType)
	}

	// Load and apply deltas if any exist
	snapshotVersion := persistent.Version()
	deltas, err := r.deltaStore.LoadDeltas(ctx, id, entityType, snapshotVersion)
	if err != nil {
		// No deltas is OK, just return snapshot
		return persistent, nil
	}

	// Apply deltas in order
	for _, delta := range deltas {
		if err = persistent.ApplyDelta(delta.Data); err != nil {
			return nil, fmt.Errorf("failed to apply delta from version %d to %d: %w",
				delta.FromVersion, delta.ToVersion, err)
		}
	}

	return persistent, nil
}

func (r *Repository) Delete(ctx context.Context, id, entityType string) error {
	if err := r.deltaStore.DeleteDeltas(ctx, id, entityType); err != nil {
		return fmt.Errorf("failed to delete deltas: %w", err)
	}

	if err := r.snapshotStore.DeleteSnapshots(ctx, id, entityType); err != nil {
		return fmt.Errorf("failed to delete snapshots: %w", err)
	}

	if err := r.snapshotStore.DeleteMetadata(ctx, id, entityType); err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	return nil
}

func (r *Repository) Exists(ctx context.Context, id, entityType string) (bool, error) {
	// Check if entity has at least one snapshot
	snapshots, err := r.snapshotStore.ListSnapshots(ctx, id, entityType)
	if err != nil {
		return false, err
	}

	return len(snapshots) > 0, nil
}

func (r *Repository) List(ctx context.Context, criteria persistence.QueryCriteria) ([]infra.Persistent, error) {
	builder := dbx.ST.
		Select("entity_id").
		From("entity_metadata").
		Where(squirrel.Eq{"entity_type": criteria.EntityType})

	for key, value := range criteria.Filters {
		builder = builder.Where(squirrel.Eq{key: value})
	}

	if criteria.OrderBy != "" {
		builder = builder.OrderBy(criteria.OrderBy)
	}

	if criteria.Limit > 0 {
		builder = builder.Limit(uint64(criteria.Limit))
	}
	if criteria.Offset > 0 {
		builder = builder.Offset(uint64(criteria.Offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list query: %w", err)
	}

	rows, err := r.db.Conn().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list entities: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []string
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	var result []infra.Persistent

	for _, id := range ids {
		var obj infra.Persistent
		obj, err = r.Load(ctx, criteria.EntityType, id)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
	}

	return result, nil
}
