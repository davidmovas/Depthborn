package status

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/davidmovas/Depthborn/pkg/identifier"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var _ Effect = (*BaseEffect)(nil)

type BaseEffect struct {
	id           string
	effectType   string
	name         string
	duration     int64
	stacks       int32
	maxStacks    int
	sourceID     string
	targetID     string
	metadata     map[string]any
	tickInterval int64
	lastTick     int64

	mu sync.RWMutex

	// Hooks for external event handling
	events map[EffectEventType]map[string]func(ctx context.Context, ev EffectEvent) error
}

func NewEffect(config EffectConfig) *BaseEffect {
	id, _ := gonanoid.New()

	if config.MaxStacks <= 0 {
		config.MaxStacks = 1
	}

	if config.InitialStacks <= 0 {
		config.InitialStacks = 1
	}

	if config.Metadata == nil {
		config.Metadata = make(map[string]interface{})
	}

	return &BaseEffect{
		id:           id,
		effectType:   config.EffectType,
		name:         config.Name,
		duration:     config.Duration,
		stacks:       int32(config.InitialStacks),
		maxStacks:    config.MaxStacks,
		sourceID:     config.SourceID,
		targetID:     config.TargetID,
		metadata:     config.Metadata,
		tickInterval: config.TickInterval,
		lastTick:     0,
		events:       make(map[EffectEventType]map[string]func(context.Context, EffectEvent) error),
	}
}

func (e *BaseEffect) ID() string {
	return e.id
}

func (e *BaseEffect) Type() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.effectType
}

func (e *BaseEffect) Name() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.name
}

func (e *BaseEffect) SetName(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.name = name
}

func (e *BaseEffect) Duration() int64 {
	return atomic.LoadInt64(&e.duration)
}

func (e *BaseEffect) SetDuration(ms int64) {
	atomic.StoreInt64(&e.duration, ms)
}

func (e *BaseEffect) Stacks() int {
	return int(atomic.LoadInt32(&e.stacks))
}

func (e *BaseEffect) MaxStacks() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.maxStacks
}

func (e *BaseEffect) AddStack() bool {
	for {
		current := atomic.LoadInt32(&e.stacks)
		if int(current) >= e.maxStacks {
			return false
		}
		if atomic.CompareAndSwapInt32(&e.stacks, current, current+1) {
			return true
		}
	}
}

func (e *BaseEffect) RemoveStack() bool {
	for {
		current := atomic.LoadInt32(&e.stacks)
		if current <= 0 {
			return false
		}
		if atomic.CompareAndSwapInt32(&e.stacks, current, current-1) {
			return true
		}
	}
}

func (e *BaseEffect) IsExpired() bool {
	return e.Duration() <= 0
}

func (e *BaseEffect) SourceID() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.sourceID
}

func (e *BaseEffect) TargetID() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.targetID
}

func (e *BaseEffect) OnEvent(ctx context.Context, ev EffectEvent) error {
	e.mu.RLock()

	var (
		callbacks = make([]func(context.Context, EffectEvent) error, len(e.events[ev.Type]))
		count     int
	)

	for _, fn := range e.events[ev.Type] {
		callbacks[count] = fn
		count++
	}

	e.mu.RUnlock()

	for _, fn := range callbacks {
		if err := fn(ctx, ev); err != nil {
			return err
		}
	}

	return nil
}

func (e *BaseEffect) OnApply(ctx context.Context, targetID string) error {
	return e.OnEvent(ctx, EffectEvent{
		Effect:   e,
		TargetID: targetID,
		Type:     EventApply,
	})
}

func (e *BaseEffect) OnTick(ctx context.Context, targetID string, deltaMs int64) error {
	return e.OnEvent(ctx, EffectEvent{
		Effect:   e,
		TargetID: targetID,
		DeltaMs:  deltaMs,
		Type:     EventTick,
	})
}

func (e *BaseEffect) OnRemove(ctx context.Context, targetID string) error {
	return e.OnEvent(ctx, EffectEvent{
		Effect:   e,
		TargetID: targetID,
		Type:     EventRemove,
	})
}

func (e *BaseEffect) OnStack(ctx context.Context, targetID string, newStacks int) error {
	return e.OnEvent(ctx, EffectEvent{
		Effect:   e,
		TargetID: targetID,
		NewStack: newStacks,
		Type:     EventStack,
	})
}

func (e *BaseEffect) AddOnEvent(eventType EffectEventType,
	fn func(ctx context.Context, ev EffectEvent) error,
) func() {
	id := identifier.New()

	e.mu.Lock()
	if e.events[eventType] == nil {
		e.events[eventType] = make(map[string]func(context.Context, EffectEvent) error)
	}
	e.events[eventType][id] = fn
	e.mu.Unlock()

	return func() {
		e.mu.Lock()
		delete(e.events[eventType], id)
		e.mu.Unlock()
	}
}

func (e *BaseEffect) CanStack(other Effect) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Can stack if same type and from same source
	return e.effectType == other.Type() && e.sourceID == other.SourceID()
}

func (e *BaseEffect) Metadata() map[string]any {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Return copy to prevent external modification
	meta := make(map[string]any, len(e.metadata))
	for k, v := range e.metadata {
		meta[k] = v
	}
	return meta
}

func (e *BaseEffect) SetMetadata(key string, value any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metadata[key] = value
}
