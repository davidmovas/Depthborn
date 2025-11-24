package sqlite

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/davidmovas/Depthborn/internal/persistence"
	"github.com/davidmovas/Depthborn/pkg/dbx"
)

var _ persistence.DeltaStore = (*DeltaStore)(nil)

type DeltaStore struct {
	db *DB
}

func NewDeltaStore(db *DB) *DeltaStore {
	return &DeltaStore{
		db: db,
	}
}

func (d *DeltaStore) SaveDelta(ctx context.Context, entityType, id string, fromVersion, toVersion int64, delta []byte) error {
	query, args, err := dbx.ST.
		Insert("deltas").
		Columns("entity_type", "entity_id", "from_version", "to_version", "timestamp", "data").
		Values(entityType, id, fromVersion, toVersion, squirrel.Expr("strftime('%s', 'now')"), delta).
		Suffix(`
			ON CONFLICT(entity_type, entity_id, from_version, to_version)
			DO UPDATE SET
				timestamp = excluded.timestamp,
				data = excluded.data
		`).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build save delta query: %w", err)
	}

	_, err = d.db.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to save delta: %w", err)
	}

	return nil
}

func (d *DeltaStore) LoadDeltas(ctx context.Context, entityType, id string, fromVersion int64) ([]persistence.Delta, error) {
	query, args, err := dbx.ST.
		Select("from_version", "to_version", "data", "timestamp").
		From("deltas").
		Where(squirrel.And{
			squirrel.Eq{"entity_type": entityType},
			squirrel.Eq{"entity_id": id},
			squirrel.GtOrEq{"from_version": fromVersion},
		}).
		OrderBy("from_version ASC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build load deltas query: %w", err)
	}

	rows, err := d.db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query deltas: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var deltas []persistence.Delta

	for rows.Next() {
		var delta persistence.Delta
		if err = rows.Scan(
			&delta.FromVersion,
			&delta.ToVersion,
			&delta.Data,
			&delta.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("failed to scan delta: %w", err)
		}
		deltas = append(deltas, delta)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return deltas, nil
}

func (d *DeltaStore) CompactDeltas(ctx context.Context, entityType, id string, upToVersion int64) error {
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query, args, err := dbx.ST.
		Delete("deltas").
		Where(squirrel.And{
			squirrel.Eq{"entity_type": entityType},
			squirrel.Eq{"entity_id": id},
			squirrel.LtOrEq{"to_version": upToVersion},
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build compact deltas query: %w", err)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete deltas: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log compaction result (optional)
	_ = rowsAffected

	return nil
}

// GetDeltaCount returns number of deltas for entity
func (d *DeltaStore) GetDeltaCount(ctx context.Context, entityType, id string) (int, error) {
	query, args, err := dbx.ST.
		Select("COUNT(*)").
		From("deltas").
		Where(squirrel.Eq{
			"entity_type": entityType,
			"entity_id":   id,
		}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build count deltas query: %w", err)
	}

	var count int
	err = d.db.conn.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count deltas: %w", err)
	}

	return count, nil
}

// GetTotalDeltaSize returns total size of deltas for entity
func (d *DeltaStore) GetTotalDeltaSize(ctx context.Context, entityType, id string) (int64, error) {
	query, args, err := dbx.ST.
		Select("COALESCE(SUM(LENGTH(data)), 0)").
		From("deltas").
		Where(squirrel.Eq{
			"entity_type": entityType,
			"entity_id":   id,
		}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build delta size query: %w", err)
	}

	var size int64
	err = d.db.conn.QueryRowContext(ctx, query, args...).Scan(&size)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate delta size: %w", err)
	}

	return size, nil
}
