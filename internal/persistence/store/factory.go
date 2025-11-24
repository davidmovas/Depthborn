package store

import (
	"fmt"

	"github.com/davidmovas/Depthborn/internal/infra/registry"
	"github.com/davidmovas/Depthborn/internal/persistence"
	"github.com/davidmovas/Depthborn/internal/persistence/store/sqlite"
)

type PersistenceSystem struct {
	DB            *sqlite.DB
	SnapshotStore persistence.SnapshotStore
	DeltaStore    persistence.DeltaStore
	Repository    persistence.Repository
}

func NewPersistenceSystem(dbName string, reg registry.Registry, strategy SnapshotStrategy) (*PersistenceSystem, error) {
	db, err := sqlite.NewDB(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	snapshotStore := sqlite.NewSnapshotStore(db, reg)
	deltaStore := sqlite.NewDeltaStore(db)

	repo := NewRepository(RepositoryConfig{
		SnapshotStore: snapshotStore,
		DeltaStore:    deltaStore,
		Registry:      reg,
		Strategy:      strategy,
	})

	return &PersistenceSystem{
		DB:            db,
		SnapshotStore: snapshotStore,
		DeltaStore:    deltaStore,
		Repository:    repo,
	}, nil
}

// Close closes all persistence connections
func (ps *PersistenceSystem) Close() error {
	return ps.DB.Close()
}
