package types

import (
	"context"
	"time"
)

// Event represents a game event that can be published and subscribed to
type Event interface {
	// EventType returns unique event type identifier
	EventType() string

	// Timestamp returns when event occurred
	Timestamp() time.Time

	// Source returns ID of entity that emitted the event
	Source() string

	// Payload returns event-specific data
	Payload() any
}

// EventHandler processes events of specific types
type EventHandler interface {
	// Handle processes the event and returns error if handling fails
	Handle(ctx context.Context, event Event) error

	// HandlesEventTypes returns list of event types this handler processes
	HandlesEventTypes() []string
}

// EventBus manages event publishing and subscription
type EventBus interface {
	// Publish sends event to all registered handlers
	Publish(ctx context.Context, event Event) error

	// Subscribe registers handler for specific event types
	Subscribe(handler EventHandler) error

	// Unsubscribe removes handler from bus
	Unsubscribe(handler EventHandler) error
}

// EventStore persists events for replay and audit
type EventStore interface {
	// Append adds event to store
	Append(ctx context.Context, event Event) error

	// Query retrieves events matching criteria
	Query(ctx context.Context, criteria EventCriteria) ([]Event, error)

	// Stream provides continuous stream of events
	Stream(ctx context.Context, criteria EventCriteria) (<-chan Event, error)
}

// EventCriteria defines filtering and ordering for event queries
type EventCriteria struct {
	EventTypes []string
	Sources    []string
	FromTime   time.Time
	ToTime     time.Time
	Limit      int
	Offset     int
}
