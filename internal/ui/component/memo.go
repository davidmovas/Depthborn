package component

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"sync"
)

// MemoizedComponent wraps component with caching
type MemoizedComponent struct {
	component Component
	lastProps map[string]any
	lastCtx   *Context
	lastHash  string
	lastOut   string
	mu        sync.RWMutex
}

// Memo wraps component to cache renders when props don't change
func Memo(comp Component, propsExtractor func(*Context) map[string]any) Component {
	mc := &MemoizedComponent{
		component: comp,
	}

	return Func(func(ctx *Context) string {
		mc.mu.Lock()
		defer mc.mu.Unlock()

		// Extract current props
		currentProps := propsExtractor(ctx)

		// Compute props hash
		currentHash := computePropsHash(currentProps)

		// Check if we can use cached output
		if mc.lastHash == currentHash && !ctx.IsDirty() {
			return mc.lastOut
		}

		// Render component
		output := mc.component.Render(ctx)

		// Cache results
		mc.lastProps = currentProps
		mc.lastCtx = ctx
		mc.lastHash = currentHash
		mc.lastOut = output

		return output
	})
}

// Smart creates auto-optimized component with dependency tracking
func Smart(comp Component) Component {
	mc := &MemoizedComponent{
		component: comp,
	}

	return Func(func(ctx *Context) string {
		mc.mu.Lock()
		defer mc.mu.Unlock()

		// Check if context is dirty
		if !ctx.IsDirty() && mc.lastOut != "" {
			return mc.lastOut
		}

		// Render component
		output := mc.component.Render(ctx)

		// Cache output
		mc.lastOut = output
		mc.lastCtx = ctx

		return output
	})
}

// computePropsHash creates hash from props map
func computePropsHash(props map[string]any) string {
	if len(props) == 0 {
		return "empty"
	}

	// Create stable string representation
	str := ""
	for k, v := range props {
		str += fmt.Sprintf("%s:%v;", k, v)
	}

	h := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", h[:8])
}

// Cache provides global component output caching
type Cache struct {
	cache map[string]string
	mu    sync.RWMutex
}

var globalCache = &Cache{
	cache: make(map[string]string),
}

// Get retrieves cached component output
func (cc *Cache) Get(key string) (string, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	val, ok := cc.cache[key]
	return val, ok
}

// Set stores component output in cache
func (cc *Cache) Set(key string, output string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.cache[key] = output
}

// Clear clears all cached components
func (cc *Cache) Clear() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.cache = make(map[string]string)
}

// Cached wraps component to use global cache
func Cached(key string, comp Component) Component {
	return Func(func(ctx *Context) string {
		// Try cache first
		if cached, ok := globalCache.Get(key); ok && !ctx.IsDirty() {
			return cached
		}

		// Render component
		output := comp.Render(ctx)

		// Store in cache
		globalCache.Set(key, output)

		return output
	})
}

// ClearCache clears global component cache
func ClearCache() {
	globalCache.Clear()
}

// PureComponent creates component that only re-renders on prop changes
func PureComponent(render func(props map[string]any) Component) func(map[string]any) Component {
	var lastProps map[string]any
	var lastComp Component

	return func(props map[string]any) Component {
		// Deep compare props
		if lastProps != nil && reflect.DeepEqual(lastProps, props) {
			return lastComp
		}

		// Props changed, re-create component
		lastProps = props
		lastComp = render(props)
		return lastComp
	}
}

// LazyComponent defers component rendering until needed
type LazyComponent struct {
	factory func() Component
	comp    Component
	once    sync.Once
}

// Lazy creates component that's only initialized once when first rendered
func Lazy(factory func() Component) Component {
	lc := &LazyComponent{
		factory: factory,
	}

	return Func(func(ctx *Context) string {
		lc.once.Do(func() {
			lc.comp = lc.factory()
		})

		if lc.comp == nil {
			return ""
		}

		return lc.comp.Render(ctx)
	})
}
