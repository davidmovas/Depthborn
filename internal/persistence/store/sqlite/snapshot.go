package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/davidmovas/Depthborn/internal/infra"
	"github.com/davidmovas/Depthborn/internal/infra/registry"
	"github.com/davidmovas/Depthborn/internal/persistence"
	"github.com/davidmovas/Depthborn/pkg/dbx"
)

var _ persistence.SnapshotStore = (*SnapshotStore)(nil)

type SnapshotStore struct {
	db       *DB
	registry registry.Registry
}

func NewSnapshotStore(db *DB, reg registry.Registry) *SnapshotStore {
	return &SnapshotStore{
		db:       db,
		registry: reg,
	}
}

func (s *SnapshotStore) SaveSnapshot(ctx context.Context, entity infra.Snapshottable) error {
	data, err := entity.Snapshot()
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Extract metadata
	var version int64
	var timestamp int64

	if versionable, ok := entity.(infra.Versionable); ok {
		version = versionable.Version()
	}

	if timestamped, ok := entity.(infra.Timestamped); ok {
		timestamp = timestamped.UpdatedAt()
	}

	query, args, err := dbx.ST.
		Insert("snapshots").
		Columns("entity_type", "entity_id", "version", "timestamp", "data", "size").
		Values(entity.Type(), entity.ID(), version, timestamp, data, len(data)).
		Suffix(`
			ON CONFLICT(entity_type, entity_id, version) 
			DO UPDATE SET 
				timestamp = excluded.timestamp,
				data = excluded.data,
				size = excluded.size
		`).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build save snapshot query: %w", err)
	}

	_, err = s.db.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	// Update metadata
	return s.updateMetadata(ctx, entity, version, timestamp)
}

func (s *SnapshotStore) LoadSnapshot(ctx context.Context, entityType, id string) (infra.Snapshottable, error) {
	// Get latest snapshot
	query, args, err := dbx.ST.
		Select("version", "data").
		From("snapshots").
		Where(squirrel.Eq{
			"entity_type": entityType,
			"entity_id":   id,
		}).
		OrderBy("version DESC").
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build load snapshot query: %w", err)
	}

	var version int64
	var data []byte

	err = s.db.conn.QueryRowContext(ctx, query, args...).Scan(&version, &data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("snapshot not found for %s/%s", entityType, id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load snapshot: %w", err)
	}

	// Create entity instance using registry
	entity, err := s.createEntity(ctx, entityType)
	if err != nil {
		return nil, err
	}

	// Restore from snapshot
	if err = entity.Restore(data); err != nil {
		return nil, fmt.Errorf("failed to restore snapshot: %w", err)
	}

	return entity, nil
}

func (s *SnapshotStore) ListSnapshots(ctx context.Context, entityType, id string) ([]persistence.SnapshotMetadata, error) {
	query, args, err := dbx.ST.
		Select("entity_type", "entity_id", "version", "timestamp", "size").
		From("snapshots").
		Where(squirrel.Eq{
			"entity_type": entityType,
			"entity_id":   id,
		}).
		OrderBy("version DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build list snapshots query: %w", err)
	}

	rows, err := s.db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query snapshots: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var snapshots []persistence.SnapshotMetadata

	for rows.Next() {
		var meta persistence.SnapshotMetadata
		if err = rows.Scan(
			&meta.EntityType,
			&meta.EntityID,
			&meta.Version,
			&meta.Timestamp,
			&meta.Size,
		); err != nil {
			return nil, fmt.Errorf("failed to scan snapshot: %w", err)
		}
		snapshots = append(snapshots, meta)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (s *SnapshotStore) updateMetadata(ctx context.Context, entity infra.Identity, version, timestamp int64) error {
	createdAt := timestamp
	if timestamped, ok := entity.(infra.Timestamped); ok {
		createdAt = timestamped.CreatedAt()
	}

	query, args, err := dbx.ST.
		Insert("entity_metadata").
		Columns("entity_type", "entity_id", "current_version", "created_at", "updated_at").
		Values(entity.Type(), entity.ID(), version, createdAt, timestamp).
		Suffix(`
			ON CONFLICT(entity_type, entity_id)
			DO UPDATE SET
				current_version = excluded.current_version,
				updated_at = excluded.updated_at
		`).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update metadata query: %w", err)
	}

	_, err = s.db.conn.ExecContext(ctx, query, args...)
	return err
}

func (s *SnapshotStore) createEntity(ctx context.Context, entityType string) (infra.Snapshottable, error) {
	identity, err := s.registry.Create(ctx, entityType, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create entity via registry: %w", err)
	}

	entity, ok := identity.(infra.Snapshottable)
	if !ok {
		return nil, fmt.Errorf("entity type %s does not implement Snapshottable", entityType)
	}

	return entity, nil
}
