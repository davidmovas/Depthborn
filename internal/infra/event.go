package infra

// DeltaEvent represents a single state change event
type DeltaEvent struct {
	Type      string         `msgpack:"type"`
	Timestamp int64          `msgpack:"timestamp"`
	Data      map[string]any `msgpack:"data"`
}

// DeltaRecord contains multiple events between two versions
type DeltaRecord struct {
	FromVersion int64        `msgpack:"from_version"`
	ToVersion   int64        `msgpack:"to_version"`
	Events      []DeltaEvent `msgpack:"events"`
}

// Common event types
const (
	EventAttributeChanged = "attribute_changed"
	EventStatusAdded      = "status_added"
	EventStatusRemoved    = "status_removed"
	EventPositionMoved    = "position_moved"
	EventFacingChanged    = "facing_changed"
	EventTagAdded         = "tag_added"
	EventTagRemoved       = "tag_removed"
)

// EventBuilder helps construct delta events
type EventBuilder struct {
	eventType string
	data      map[string]any
}

func NewEventBuilder(eventType string) *EventBuilder {
	return &EventBuilder{
		eventType: eventType,
		data:      make(map[string]any),
	}
}

func (b *EventBuilder) Set(key string, value any) *EventBuilder {
	b.data[key] = value
	return b
}

func (b *EventBuilder) Build(timestamp int64) DeltaEvent {
	return DeltaEvent{
		Type:      b.eventType,
		Timestamp: timestamp,
		Data:      b.data,
	}
}
