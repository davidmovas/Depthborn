package impl

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/davidmovas/Depthborn/internal/infra"
	"github.com/davidmovas/Depthborn/internal/persistence/serializer"
)

var _ infra.Persistent = (*BasePersistent)(nil)

type BasePersistent struct {
	identity    *BaseIdentity
	versionable *BaseVersionable
	timestamped *BaseTimestamped

	owner infra.Serializable

	deltaHistory []infra.DeltaRecord
	maxDeltas    int

	mu sync.RWMutex

	serializer serializer.Serializer
}

type PersistentConfig struct {
	MaxDeltas  int
	Serializer serializer.Serializer
}

func NewPersistent(entityType string, owner infra.Serializable, config *PersistentConfig) *BasePersistent {
	if config == nil {
		config = &PersistentConfig{
			MaxDeltas:  50,
			Serializer: serializer.NewMessagePackSerializer(),
		}
	}

	return &BasePersistent{
		identity:     NewIdentity(entityType),
		versionable:  NewVersionable(),
		timestamped:  NewTimestamped(),
		owner:        owner,
		deltaHistory: make([]infra.DeltaRecord, 0, config.MaxDeltas),
		maxDeltas:    config.MaxDeltas,
		serializer:   config.Serializer,
	}
}

func NewPersistentWithID(id, entityType string, owner infra.Serializable, config *PersistentConfig) *BasePersistent {
	if config == nil {
		config = &PersistentConfig{
			MaxDeltas:  50,
			Serializer: serializer.NewMessagePackSerializer(),
		}
	}

	return &BasePersistent{
		identity:     NewIdentityWithID(entityType, id),
		versionable:  NewVersionable(),
		timestamped:  NewTimestamped(),
		owner:        owner,
		deltaHistory: make([]infra.DeltaRecord, 0, config.MaxDeltas),
		maxDeltas:    config.MaxDeltas,
		serializer:   config.Serializer,
	}
}

func (p *BasePersistent) ID() string {
	return p.identity.ID()
}

func (p *BasePersistent) Type() string {
	return p.identity.Type()
}

func (p *BasePersistent) Version() int64 {
	return p.versionable.Version()
}

func (p *BasePersistent) IncrementVersion() int64 {
	p.timestamped.Touch()
	return p.versionable.IncrementVersion()
}

func (p *BasePersistent) CreatedAt() int64 {
	return p.timestamped.CreatedAt()
}

func (p *BasePersistent) UpdatedAt() int64 {
	return p.timestamped.UpdatedAt()
}

func (p *BasePersistent) Touch() {
	p.timestamped.Touch()
}

func (p *BasePersistent) Snapshot() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	state, err := p.owner.SerializeState()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize owner state: %w", err)
	}

	snapshot := map[string]any{
		"id":         p.ID(),
		"type":       p.Type(),
		"version":    p.Version(),
		"created_at": p.CreatedAt(),
		"updated_at": p.UpdatedAt(),
		"state":      state,
	}

	data, err := p.serializer.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	return data, nil
}

func (p *BasePersistent) Restore(data []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var snapshot map[string]any
	if err := p.serializer.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	// Restore metadata
	if version, ok := snapshot["version"].(int64); ok {
		p.versionable.SetVersion(version)
	}
	if createdAt, ok := snapshot["created_at"].(int64); ok {
		p.timestamped.SetCreatedAt(createdAt)
	}
	if updatedAt, ok := snapshot["updated_at"].(int64); ok {
		atomic.StoreInt64(&p.timestamped.updatedAt, updatedAt)
	}

	// Restore owner state
	state, ok := snapshot["state"].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid state format in snapshot")
	}

	if err := p.owner.DeserializeState(state); err != nil {
		return fmt.Errorf("failed to deserialize owner state: %w", err)
	}

	// Clear delta history after restore
	p.deltaHistory = make([]infra.DeltaRecord, 0, p.maxDeltas)

	return nil
}

func (p *BasePersistent) Delta(fromVersion int64) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	currentVersion := p.Version()
	if fromVersion >= currentVersion {
		return nil, fmt.Errorf("fromVersion %d is not less than current version %d", fromVersion, currentVersion)
	}

	var relevantDeltas []infra.DeltaRecord
	for _, delta := range p.deltaHistory {
		if delta.FromVersion >= fromVersion && delta.ToVersion <= currentVersion {
			relevantDeltas = append(relevantDeltas, delta)
		}
	}

	if len(relevantDeltas) == 0 {
		return nil, fmt.Errorf("no deltas available from version %d", fromVersion)
	}

	data, err := p.serializer.Marshal(relevantDeltas)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deltas: %w", err)
	}

	return data, nil
}

func (p *BasePersistent) ApplyDelta(delta []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var deltaRecords []infra.DeltaRecord
	if err := p.serializer.Unmarshal(delta, &deltaRecords); err != nil {
		return fmt.Errorf("failed to unmarshal delta: %w", err)
	}

	for _, record := range deltaRecords {
		if record.FromVersion != p.Version() {
			return fmt.Errorf("delta version mismatch: expected %d, got %d", p.Version(), record.FromVersion)
		}

		// Apply each event in the delta record
		for _, event := range record.Events {
			if err := p.applyEvent(event); err != nil {
				return fmt.Errorf("failed to apply event %s: %w", event.Type, err)
			}
		}

		p.versionable.SetVersion(record.ToVersion)
		p.timestamped.Touch()
	}

	return nil
}

// RecordDelta adds a new delta record to history
func (p *BasePersistent) RecordDelta(events []infra.DeltaEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(events) == 0 {
		return
	}

	fromVersion := p.Version()
	toVersion := p.IncrementVersion()

	record := infra.DeltaRecord{
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Events:      events,
	}

	p.deltaHistory = append(p.deltaHistory, record)

	// Trim history if exceeds max
	if len(p.deltaHistory) > p.maxDeltas {
		p.deltaHistory = p.deltaHistory[len(p.deltaHistory)-p.maxDeltas:]
	}
}

// applyEvent applies a single delta event to owner
func (p *BasePersistent) applyEvent(event infra.DeltaEvent) error {
	// Get current state
	state, err := p.owner.SerializeState()
	if err != nil {
		return err
	}

	// Apply event data to state
	for key, value := range event.Data {
		state[key] = value
	}

	// Restore modified state
	return p.owner.DeserializeState(state)
}

// GetDeltaHistory returns copy of delta history (for testing/debugging)
func (p *BasePersistent) GetDeltaHistory() []infra.DeltaRecord {
	p.mu.RLock()
	defer p.mu.RUnlock()

	history := make([]infra.DeltaRecord, len(p.deltaHistory))
	copy(history, p.deltaHistory)
	return history
}

// ClearDeltaHistory removes all delta history (useful after snapshot)
func (p *BasePersistent) ClearDeltaHistory() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.deltaHistory = make([]infra.DeltaRecord, 0, p.maxDeltas)
}
