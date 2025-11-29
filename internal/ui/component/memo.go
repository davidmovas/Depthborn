package component

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

// Memo wraps a component to cache its output when props don't change.
// The propsHash function should return a stable hash representing the props.
func Memo(comp Component, propsHash func() string) Component {
	var (
		lastHash   string
		lastOutput string
		mu         sync.Mutex
	)

	return Func(func(ctx *Context) string {
		mu.Lock()
		defer mu.Unlock()

		currentHash := propsHash()

		if currentHash == lastHash && lastOutput != "" {
			return lastOutput
		}

		output := comp.Render(ctx)
		lastHash = currentHash
		lastOutput = output

		return output
	})
}

// MemoSimple caches output until explicitly invalidated.
// Use invalidate() to clear cache.
func MemoSimple(comp Component) (Component, func()) {
	var (
		cached    string
		hasCached bool
		mu        sync.Mutex
	)

	invalidate := func() {
		mu.Lock()
		hasCached = false
		cached = ""
		mu.Unlock()
	}

	memoized := Func(func(ctx *Context) string {
		mu.Lock()
		defer mu.Unlock()

		if hasCached {
			return cached
		}

		output := comp.Render(ctx)
		cached = output
		hasCached = true

		return output
	})

	return memoized, invalidate
}

// Pure creates a component that only re-renders when props change.
// Props are compared by their string representation.
func Pure[P any](render func(props P) Component) func(P) Component {
	var (
		lastComp Component
		lastHash string
		mu       sync.Mutex
	)

	return func(props P) Component {
		mu.Lock()
		defer mu.Unlock()

		currentHash := hashProps(props)

		if currentHash == lastHash && lastComp != nil {
			return lastComp
		}

		lastHash = currentHash
		lastComp = render(props)

		return lastComp
	}
}

// Lazy defers component creation until first render.
func Lazy(factory func() Component) Component {
	var (
		comp Component
		once sync.Once
	)

	return Func(func(ctx *Context) string {
		once.Do(func() {
			comp = factory()
		})

		if comp == nil {
			return ""
		}

		return comp.Render(ctx)
	})
}

// LazyWithFallback shows fallback while component loads.
func LazyWithFallback(factory func() Component, fallback Component) Component {
	var (
		comp    Component
		loading = true
		mu      sync.Mutex
	)

	// Start loading in background
	go func() {
		result := factory()
		mu.Lock()
		comp = result
		loading = false
		mu.Unlock()
	}()

	return Func(func(ctx *Context) string {
		mu.Lock()
		isLoading := loading
		current := comp
		mu.Unlock()

		if isLoading {
			if fallback != nil {
				return fallback.Render(ctx)
			}
			return ""
		}

		if current != nil {
			return current.Render(ctx)
		}

		return ""
	})
}

// --- Caching System ---

// Cache provides a simple key-value cache for rendered output.
type Cache struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewCache creates a new cache.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

// Get retrieves a cached value.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

// Set stores a value in cache.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Delete removes a value from cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear removes all cached values.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]string)
}

// Global cache instance
var globalCache = NewCache()

// Cached wraps a component to use global cache.
func Cached(key string, comp Component) Component {
	return Func(func(ctx *Context) string {
		if cached, ok := globalCache.Get(key); ok {
			return cached
		}

		output := comp.Render(ctx)
		globalCache.Set(key, output)

		return output
	})
}

// InvalidateCache clears a specific cache entry.
func InvalidateCache(key string) {
	globalCache.Delete(key)
}

// ClearAllCache clears the entire cache.
func ClearAllCache() {
	globalCache.Clear()
}

// --- Helper Functions ---

// hashProps creates a hash from props.
func hashProps(props any) string {
	str := fmt.Sprintf("%#v", props)
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:8])
}

// HashString creates a hash from multiple strings.
func HashString(parts ...string) string {
	combined := ""
	for _, p := range parts {
		combined += p + "|"
	}
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:8])
}
