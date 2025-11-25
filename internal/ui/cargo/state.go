package cargo

import (
	"context"
	"time"
)

// Middleware is a function that intercepts and potentially modifies state
// Can be used for: validation, logging, persistence, invalidation, etc.
type Middleware[T any] func(ctx context.Context, data T) T

// Get retrieves value by key
// Returns zero value and false if not found
func Get[T any](key string) (T, bool) {
	val, ok := cargo.get(key)
	if !ok {
		var zero T
		return zero, false
	}

	typed, ok := val.(T)
	if !ok {
		var zero T
		return zero, false
	}

	return typed, true
}

// Set stores value by key
func Set[T any](key string, value T) {
	cargo.set(key, value)
}

// Use retrieves value with middleware applied
// Middleware is called in order: first to last
// If value doesn't exist, returns zero value
func Use[T any](key string, middleware ...Middleware[T]) T {
	val, ok := Get[T](key)
	if !ok {
		var zero T
		val = zero
	}

	// Apply middleware chain
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	time.AfterFunc(time.Minute, cancel)

	for _, mw := range middleware {
		val = mw(ctx, val)
	}

	// Store modified value back
	Set(key, val)

	return val
}

// Delete removes value by key
func Delete(key string) {
	cargo.delete(key)
}

// Has checks if key exists
func Has(key string) bool {
	return cargo.has(key)
}

// GetGlobal retrieves global instance by type
// Returns zero value and false if not found
func GetGlobal[T any]() (T, bool) {
	key := typeKey[T]()
	return Get[T](key)
}

// SetGlobal stores global instance by type
// Only one instance per type can exist globally
func SetGlobal[T any](value T) {
	key := typeKey[T]()
	Set(key, value)
}

// UseGlobal retrieves global instance with middleware applied
func UseGlobal[T any](middleware ...Middleware[T]) T {
	key := typeKey[T]()
	return Use(key, middleware...)
}

// DeleteGlobal removes global instance by type
func DeleteGlobal[T any]() {
	key := typeKey[T]()
	Delete(key)
}

// HasGlobal checks if global instance exists
func HasGlobal[T any]() bool {
	key := typeKey[T]()
	return Has(key)
}

// Clear removes all state values
func Clear() {
	cargo.clear()
}
