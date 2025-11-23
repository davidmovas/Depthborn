package persistence

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/infra"
)

// Repository provides CRUD operations for persistent entities
type Repository interface {
	// Save persists entity, creating or updating as needed
	Save(ctx context.Context, entity infra.Persistent) error

	// Load retrieves entity by ID
	Load(ctx context.Context, entityType, id string) (infra.Persistent, error)

	// Delete removes entity from storage
	Delete(ctx context.Context, entityType, id string) error

	// Exists checks if entity exists without loading it
	Exists(ctx context.Context, entityType, id string) (bool, error)

	// List retrieves multiple entities matching criteria
	List(ctx context.Context, criteria QueryCriteria) ([]infra.Persistent, error)
}

// QueryCriteria defines filtering and pagination for entity queries
type QueryCriteria struct {
	EntityType string
	Filters    map[string]any
	Limit      int
	Offset     int
	OrderBy    string
}

// SnapshotStore handles full state snapshots
type SnapshotStore interface {
	// SaveSnapshot stores complete entity state
	SaveSnapshot(ctx context.Context, entity infra.Snapshottable) error

	// LoadSnapshot retrieves and restores entity from snapshot
	LoadSnapshot(ctx context.Context, entityType, id string) (infra.Snapshottable, error)

	// ListSnapshots returns available snapshots for entity
	ListSnapshots(ctx context.Context, entityType, id string) ([]SnapshotMetadata, error)
}

// SnapshotMetadata describes a stored snapshot
type SnapshotMetadata struct {
	EntityType string
	EntityID   string
	Version    int64
	Timestamp  int64
	Size       int64
}

// DeltaStore handles incremental state changes
type DeltaStore interface {
	// SaveDelta stores incremental change
	SaveDelta(ctx context.Context, entityType, id string, fromVersion, toVersion int64, delta []byte) error

	// LoadDeltas retrieves all deltas since specified version
	LoadDeltas(ctx context.Context, entityType, id string, fromVersion int64) ([]Delta, error)

	// CompactDeltas merges multiple deltas into single snapshot
	CompactDeltas(ctx context.Context, entityType, id string, upToVersion int64) error
}

// Delta represents an incremental state change
type Delta struct {
	FromVersion int64
	ToVersion   int64
	Data        []byte
	Timestamp   int64
}
