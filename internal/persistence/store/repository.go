package store

import (
	"context"
	"fmt"

	"github.com/davidmovas/Depthborn/internal/infra"
	"github.com/davidmovas/Depthborn/internal/infra/registry"
	"github.com/davidmovas/Depthborn/internal/persistence"
	"github.com/davidmovas/Depthborn/internal/persistence/store/sqlite"
)

var _ persistence.Repository = (*Repository)(nil)

type Repository struct {
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

func (r *Repository) Load(ctx context.Context, entityType, id string) (infra.Persistent, error) {
	// Try to load from snapshot first
	entity, err := r.snapshotStore.LoadSnapshot(ctx, entityType, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load snapshot: %w", err)
	}

	persistent, ok := entity.(infra.Persistent)
	if !ok {
		return nil, fmt.Errorf("entity %s/%s does not implement Persistent", entityType, id)
	}

	// Load and apply deltas if any exist
	snapshotVersion := persistent.Version()
	deltas, err := r.deltaStore.LoadDeltas(ctx, entityType, id, snapshotVersion)
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

func (r *Repository) Delete(ctx context.Context, entityType, id string) error {
	// TODO: Implement deletion from both snapshot and delta stores
	// This requires adding Delete methods to SnapshotStore and DeltaStore interfaces
	return fmt.Errorf("Delete not yet implemented")
}

func (r *Repository) Exists(ctx context.Context, entityType, id string) (bool, error) {
	// Check if entity has at least one snapshot
	snapshots, err := r.snapshotStore.ListSnapshots(ctx, entityType, id)
	if err != nil {
		return false, err
	}

	return len(snapshots) > 0, nil
}

func (r *Repository) List(ctx context.Context, criteria persistence.QueryCriteria) ([]infra.Persistent, error) {
	// TODO: Implement listing with filtering
	// This requires querying entity_metadata table and loading each entity
	return nil, fmt.Errorf("List not yet implemented")
}

// SnapshotStrategy determines when to create full snapshots
type SnapshotStrategy interface {
	ShouldSnapshot(currentVersion int64, deltaCount int, totalDeltaSize int64) bool
}

// EveryNVersionsStrategy creates snapshot every N versions
type EveryNVersionsStrategy struct {
	interval int64
}

func NewEveryNVersionsStrategy(interval int64) SnapshotStrategy {
	return &EveryNVersionsStrategy{interval: interval}
}

func (s *EveryNVersionsStrategy) ShouldSnapshot(currentVersion int64, deltaCount int, totalDeltaSize int64) bool {
	return currentVersion%s.interval == 0
}

// DeltaSizeStrategy creates snapshot when deltas exceed size threshold
type DeltaSizeStrategy struct {
	maxSize int64
}

func NewDeltaSizeStrategy(maxSize int64) SnapshotStrategy {
	return &DeltaSizeStrategy{maxSize: maxSize}
}

func (s *DeltaSizeStrategy) ShouldSnapshot(currentVersion int64, deltaCount int, totalDeltaSize int64) bool {
	return totalDeltaSize > s.maxSize
}

// HybridStrategy combines multiple strategies
type HybridStrategy struct {
	strategies []SnapshotStrategy
}

func NewHybridStrategy(strategies ...SnapshotStrategy) SnapshotStrategy {
	return &HybridStrategy{strategies: strategies}
}

func (s *HybridStrategy) ShouldSnapshot(currentVersion int64, deltaCount int, totalDeltaSize int64) bool {
	for _, strategy := range s.strategies {
		if strategy.ShouldSnapshot(currentVersion, deltaCount, totalDeltaSize) {
			return true
		}
	}
	return false
}
