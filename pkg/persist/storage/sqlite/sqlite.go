// Package sqlite provides SQLite-based storage implementation.
package sqlite

import (
	"context"
	"database/sql"
	"sync"
	"time"

	_ "modernc.org/sqlite" // Pure-Go SQLite driver (no CGO required)

	"github.com/davidmovas/Depthborn/pkg/persist/storage"
)

// Storage implements storage.FullStorage using SQLite.
type Storage struct {
	db     *sql.DB
	closed bool
	mu     sync.RWMutex
}

// Config holds SQLite storage configuration.
type Config struct {
	// Path to the SQLite database file. Use ":memory:" for in-memory.
	Path string

	// MaxOpenConns is the maximum number of open connections.
	MaxOpenConns int

	// MaxIdleConns is the maximum number of idle connections.
	MaxIdleConns int

	// ConnMaxLifetime is the maximum lifetime of a connection.
	ConnMaxLifetime time.Duration

	// EnableWAL enables Write-Ahead Logging for better concurrency.
	EnableWAL bool

	// BusyTimeout is the timeout for locked database (milliseconds).
	BusyTimeout int
}

// DefaultConfig returns sensible defaults for game usage.
func DefaultConfig(path string) Config {
	return Config{
		Path:            path,
		MaxOpenConns:    1, // SQLite works best with single writer
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Hour,
		EnableWAL:       true,
		BusyTimeout:     5000, // 5 seconds
	}
}

// Open creates a new SQLite storage with the given configuration.
func Open(cfg Config) (*Storage, error) {
	dsn := cfg.Path
	if cfg.BusyTimeout > 0 {
		dsn += "?_busy_timeout=" + string(rune(cfg.BusyTimeout))
	}

	db, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	// Enable WAL mode for better concurrency
	if cfg.EnableWAL {
		if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
			db.Close()
			return nil, err
		}
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return nil, err
	}

	s := &Storage{db: db}

	// Create schema
	if err := s.createSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

// OpenMemory creates an in-memory SQLite storage (useful for testing).
func OpenMemory() (*Storage, error) {
	return Open(DefaultConfig(":memory:"))
}

func (s *Storage) createSchema() error {
	schema := `
		CREATE TABLE IF NOT EXISTS entities (
			key TEXT PRIMARY KEY,
			data BLOB NOT NULL,
			version INTEGER NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_entities_updated ON entities(updated_at);
	`
	_, err := s.db.Exec(schema)
	return err
}

// Get retrieves a record by key.
func (s *Storage) Get(ctx context.Context, key string) (*storage.Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, storage.ErrClosed
	}

	row := s.db.QueryRowContext(ctx,
		"SELECT key, data, version, created_at, updated_at FROM entities WHERE key = ?",
		key,
	)

	var r storage.Record
	err := row.Scan(&r.Key, &r.Data, &r.Version, &r.CreatedAt, &r.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// Set stores a record.
func (s *Storage) Set(ctx context.Context, record *storage.Record) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return storage.ErrClosed
	}

	now := time.Now().Unix()
	if record.CreatedAt == 0 {
		record.CreatedAt = now
	}
	record.UpdatedAt = now

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO entities (key, data, version, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET
		   data = excluded.data,
		   version = excluded.version,
		   updated_at = excluded.updated_at`,
		record.Key, record.Data, record.Version, record.CreatedAt, record.UpdatedAt,
	)
	return err
}

// Delete removes a record by key.
func (s *Storage) Delete(ctx context.Context, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return storage.ErrClosed
	}

	_, err := s.db.ExecContext(ctx, "DELETE FROM entities WHERE key = ?", key)
	return err
}

// Exists checks if a key exists.
func (s *Storage) Exists(ctx context.Context, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return false, storage.ErrClosed
	}

	var count int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM entities WHERE key = ?",
		key,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Close closes the storage.
func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return storage.ErrClosed
	}

	s.closed = true
	return s.db.Close()
}

// --- BatchStorage implementation ---

// GetMany retrieves multiple records by keys.
func (s *Storage) GetMany(ctx context.Context, keys []string) ([]*storage.Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, storage.ErrClosed
	}

	if len(keys) == 0 {
		return nil, nil
	}

	// Build query with placeholders
	query := "SELECT key, data, version, created_at, updated_at FROM entities WHERE key IN (?" +
		repeatString(",?", len(keys)-1) + ")"

	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*storage.Record
	for rows.Next() {
		var r storage.Record
		if err := rows.Scan(&r.Key, &r.Data, &r.Version, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		records = append(records, &r)
	}

	return records, rows.Err()
}

// SetMany stores multiple records atomically.
func (s *Storage) SetMany(ctx context.Context, records []*storage.Record) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return storage.ErrClosed
	}

	if len(records) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO entities (key, data, version, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET
		   data = excluded.data,
		   version = excluded.version,
		   updated_at = excluded.updated_at`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now().Unix()
	for _, r := range records {
		if r.CreatedAt == 0 {
			r.CreatedAt = now
		}
		r.UpdatedAt = now

		if _, err := stmt.ExecContext(ctx, r.Key, r.Data, r.Version, r.CreatedAt, r.UpdatedAt); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteMany removes multiple records atomically.
func (s *Storage) DeleteMany(ctx context.Context, keys []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return storage.ErrClosed
	}

	if len(keys) == 0 {
		return nil
	}

	query := "DELETE FROM entities WHERE key IN (?" + repeatString(",?", len(keys)-1) + ")"

	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}

	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// --- QueryableStorage implementation ---

// List returns all keys matching the prefix.
func (s *Storage) List(ctx context.Context, prefix string, limit int, offset int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, storage.ErrClosed
	}

	query := "SELECT key FROM entities WHERE key LIKE ? ORDER BY key"
	args := []any{prefix + "%"}

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	if offset > 0 {
		query += " OFFSET ?"
		args = append(args, offset)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, rows.Err()
}

// Count returns the number of keys matching the prefix.
func (s *Storage) Count(ctx context.Context, prefix string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, storage.ErrClosed
	}

	var count int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM entities WHERE key LIKE ?",
		prefix+"%",
	).Scan(&count)

	return count, err
}

// --- TransactionalStorage implementation ---

// Begin starts a new transaction.
func (s *Storage) Begin(ctx context.Context) (storage.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, storage.ErrClosed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &transaction{tx: tx}, nil
}

// transaction implements storage.Transaction.
type transaction struct {
	tx     *sql.Tx
	closed bool
	mu     sync.Mutex
}

func (t *transaction) Get(ctx context.Context, key string) (*storage.Record, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil, storage.ErrTxClosed
	}

	row := t.tx.QueryRowContext(ctx,
		"SELECT key, data, version, created_at, updated_at FROM entities WHERE key = ?",
		key,
	)

	var r storage.Record
	err := row.Scan(&r.Key, &r.Data, &r.Version, &r.CreatedAt, &r.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (t *transaction) Set(ctx context.Context, record *storage.Record) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return storage.ErrTxClosed
	}

	now := time.Now().Unix()
	if record.CreatedAt == 0 {
		record.CreatedAt = now
	}
	record.UpdatedAt = now

	_, err := t.tx.ExecContext(ctx,
		`INSERT INTO entities (key, data, version, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET
		   data = excluded.data,
		   version = excluded.version,
		   updated_at = excluded.updated_at`,
		record.Key, record.Data, record.Version, record.CreatedAt, record.UpdatedAt,
	)
	return err
}

func (t *transaction) Delete(ctx context.Context, key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return storage.ErrTxClosed
	}

	_, err := t.tx.ExecContext(ctx, "DELETE FROM entities WHERE key = ?", key)
	return err
}

func (t *transaction) Commit() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return storage.ErrTxClosed
	}

	t.closed = true
	return t.tx.Commit()
}

func (t *transaction) Rollback() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return storage.ErrTxClosed
	}

	t.closed = true
	return t.tx.Rollback()
}

// --- Helpers ---

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// Verify interface compliance.
var _ storage.FullStorage = (*Storage)(nil)
