package cargo

import "context"

// Chain - chain middleware functions
func Chain[T any](middlewares ...Middleware[T]) Middleware[T] {
	return func(ctx context.Context, data T) T {
		result := data
		for _, mw := range middlewares {
			result = mw(ctx, result)
		}
		return result
	}
}

// NoOp - no-op middleware (for testing purposes)
func NoOp[T any]() Middleware[T] {
	return func(ctx context.Context, data T) T {
		return data
	}
}

// Example Middleware: Log (для демонстрации)
// В реальном коде пользователь создаст свои middleware в app/middleware/
/*
func LogMiddleware[T any](logger Logger) Middleware[T] {
	return func(ctx context.Context, data T) T {
		logger.Info("State accessed", "type", fmt.Sprintf("%T", data))
		return data
	}
}

// Example: Validation
func ValidateMiddleware[T any](validator func(T) error) Middleware[T] {
	return func(ctx context.Context, data T) T {
		if err := validator(data); err != nil {
			// Handle validation error
			// Could log, panic, or return zero value
			var zero T
			return zero
		}
		return data
	}
}

// Example: Auto-refresh (invalidation)
func AutoRefreshMiddleware[T any](refreshFunc func(context.Context) T) Middleware[T] {
	return func(ctx context.Context, data T) T {
		// Check if data needs refresh
		// For example, check timestamp or version
		return refreshFunc(ctx)
	}
}

// Example: Cache with TTL
func CacheTTLMiddleware[T any](ttl time.Duration) Middleware[T] {
	lastUpdate := time.Now()
	return func(ctx context.Context, data T) T {
		if time.Since(lastUpdate) > ttl {
			// Data expired, return zero to force refresh
			var zero T
			lastUpdate = time.Now()
			return zero
		}
		return data
	}
}
*/
