package infra

// Identity represents any entity that can be uniquely identified
type Identity interface {
	// ID returns unique identifier of the entity
	ID() string

	// Type returns the type name of the entity for routing and registry purposes
	Type() string
}

// Versionable represents entities that track their version for optimistic locking
type Versionable interface {
	// Version returns current version of the entity
	Version() int64

	// IncrementVersion bumps version by 1 and returns new value
	IncrementVersion() int64
}

// Timestamped represents entities that track creation and modification time
type Timestamped interface {
	// CreatedAt returns when entity was created
	CreatedAt() int64

	// UpdatedAt returns when entity was last modified
	UpdatedAt() int64

	// Touch updates the UpdatedAt timestamp to current time
	Touch()
}

// Snapshottable represents entities that can create and restore from snapshots
type Snapshottable interface {
	Identity

	// Snapshot creates a complete snapshot of current state
	Snapshot() ([]byte, error)

	// Restore restores entity state from snapshot data
	Restore(data []byte) error
}

// DeltaTrackable represents entities that can track and apply incremental changes
type DeltaTrackable interface {
	Identity
	Versionable

	// Delta returns changes since specified version
	Delta(fromVersion int64) ([]byte, error)

	// ApplyDelta applies incremental changes to current state
	ApplyDelta(delta []byte) error
}

// Persistent combines all persistence-related capabilities
type Persistent interface {
	Snapshottable
	DeltaTrackable
	Timestamped
}

// Cloneable represents entities that can be deep-copied
type Cloneable interface {
	// Clone creates a deep copy of the entity
	Clone() interface{}
}

// Validatable represents entities that can validate their state
type Validatable interface {
	// Validate checks if entity state is valid and returns error if not
	Validate() error
}

// Disposable represents entities that need cleanup
type Disposable interface {
	// Dispose releases resources held by entity
	Dispose() error
}
