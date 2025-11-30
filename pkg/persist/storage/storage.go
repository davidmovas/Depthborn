// Package storage defines the low-level storage interface for persistence.
package storage

import (
	"context"
	"errors"
)

// Common errors returned by storage operations.
var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
	ErrClosed        = errors.New("storage is closed")
	ErrTxClosed      = errors.New("transaction is closed")
)

// Record represents a stored entity with metadata.
type Record struct {
	Key       string // Primary key (entity_type:entity_id)
	Data      []byte // Serialized entity data
	Version   int64  // Version for optimistic locking
	CreatedAt int64  // Unix timestamp
	UpdatedAt int64  // Unix timestamp
}

// Storage provides low-level key-value storage operations.
type Storage interface {
	// Get retrieves a record by key. Returns ErrNotFound if not exists.
	Get(ctx context.Context, key string) (*Record, error)

	// Set stores a record. Creates new or updates existing.
	Set(ctx context.Context, record *Record) error

	// Delete removes a record by key. No error if not exists.
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists.
	Exists(ctx context.Context, key string) (bool, error)

	// Close closes the storage and releases resources.
	Close() error
}

// BatchStorage extends Storage with batch operations for better performance.
type BatchStorage interface {
	Storage

	// GetMany retrieves multiple records by keys.
	// Missing keys are omitted from the result (no error).
	GetMany(ctx context.Context, keys []string) ([]*Record, error)

	// SetMany stores multiple records atomically.
	SetMany(ctx context.Context, records []*Record) error

	// DeleteMany removes multiple records atomically.
	DeleteMany(ctx context.Context, keys []string) error
}

// QueryableStorage extends Storage with query capabilities.
type QueryableStorage interface {
	Storage

	// List returns all keys matching the prefix.
	List(ctx context.Context, prefix string, limit int, offset int) ([]string, error)

	// Count returns the number of keys matching the prefix.
	Count(ctx context.Context, prefix string) (int, error)
}

// TransactionalStorage extends Storage with transaction support.
type TransactionalStorage interface {
	Storage

	// Begin starts a new transaction.
	Begin(ctx context.Context) (Transaction, error)
}

// Transaction represents an atomic unit of work.
type Transaction interface {
	// Get retrieves a record within the transaction.
	Get(ctx context.Context, key string) (*Record, error)

	// Set stores a record within the transaction.
	Set(ctx context.Context, record *Record) error

	// Delete removes a record within the transaction.
	Delete(ctx context.Context, key string) error

	// Commit applies all changes.
	Commit() error

	// Rollback discards all changes.
	Rollback() error
}

// FullStorage combines all storage capabilities.
type FullStorage interface {
	BatchStorage
	QueryableStorage
	TransactionalStorage
}

// Key helpers for consistent key formatting.

// EntityKey creates a storage key for an entity.
func EntityKey(entityType, entityID string) string {
	return entityType + ":" + entityID
}

// SnapshotKey creates a storage key for a snapshot.
func SnapshotKey(entityType, entityID string, version int64) string {
	return "snapshot:" + entityType + ":" + entityID
}

// DeltaKey creates a storage key for a delta.
func DeltaKey(entityType, entityID string, fromVersion, toVersion int64) string {
	return "delta:" + entityType + ":" + entityID
}
