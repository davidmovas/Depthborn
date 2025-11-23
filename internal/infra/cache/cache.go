package cache

import (
	"context"
	"time"
)

// Cache provides generic key-value caching with TTL support
type Cache interface {
	// Get retrieves value by key, returns false if not found
	Get(ctx context.Context, key string) (any, bool)

	// Set stores value with optional TTL (0 = no expiration)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	// Delete removes value by key
	Delete(ctx context.Context, key string) error

	// Has checks if key exists without retrieving value
	Has(ctx context.Context, key string) bool

	// Clear removes all cached values
	Clear(ctx context.Context) error

	// Size returns number of cached items
	Size() int
}

// TypedCache provides type-safe caching for specific value type
type TypedCache[T any] interface {
	// Get retrieves typed value by key
	Get(ctx context.Context, key string) (T, bool)

	// Set stores typed value
	Set(ctx context.Context, key string, value T, ttl time.Duration) error

	// Delete removes value by key
	Delete(ctx context.Context, key string) error

	// List returns all keys in cache
	List() []string
}

// Policy defines caching strategy for entities
type Policy interface {
	// ShouldCache determines if entity should be cached
	ShouldCache(entity interface{}) bool

	// TTL returns how long entity should be cached
	TTL(entity interface{}) time.Duration

	// EvictionPriority returns priority for cache eviction (higher = keep longer)
	EvictionPriority(entity interface{}) int
}

// LoadingCache automatically loads missing values using provided loader
type LoadingCache[K comparable, V any] interface {
	// Get retrieves value, loading it if not cached
	Get(ctx context.Context, key K) (V, error)

	// Invalidate removes cached value, forcing reload on next Get
	Invalidate(key K)

	// Refresh reloads value even if cached
	Refresh(ctx context.Context, key K) error
}

// Loader retrieves value for given key when not in cache
type Loader[K comparable, V any] interface {
	// Load retrieves value for key
	Load(ctx context.Context, key K) (V, error)
}
