package store

// SnapshotStrategy determines when to create full snapshots
type SnapshotStrategy interface {
	ShouldSnapshot(currentVersion int64, deltaCount int, totalDeltaSize int64) bool
}

// EveryNVersionsStrategy creates snapshot every N versions
type EveryNVersionsStrategy struct {
	interval int64
}

func NewEveryNVersionsStrategy(interval int64) SnapshotStrategy {
	return &EveryNVersionsStrategy{interval: interval}
}

func (s *EveryNVersionsStrategy) ShouldSnapshot(currentVersion int64, deltaCount int, totalDeltaSize int64) bool {
	return currentVersion%s.interval == 0
}

// DeltaSizeStrategy creates snapshot when deltas exceed size threshold
type DeltaSizeStrategy struct {
	maxSize int64
}

func NewDeltaSizeStrategy(maxSize int64) SnapshotStrategy {
	return &DeltaSizeStrategy{maxSize: maxSize}
}

func (s *DeltaSizeStrategy) ShouldSnapshot(currentVersion int64, deltaCount int, totalDeltaSize int64) bool {
	return totalDeltaSize > s.maxSize
}

// HybridStrategy combines multiple strategies
type HybridStrategy struct {
	strategies []SnapshotStrategy
}

func NewHybridStrategy(strategies ...SnapshotStrategy) SnapshotStrategy {
	return &HybridStrategy{strategies: strategies}
}

func (s *HybridStrategy) ShouldSnapshot(currentVersion int64, deltaCount int, totalDeltaSize int64) bool {
	for _, strategy := range s.strategies {
		if strategy.ShouldSnapshot(currentVersion, deltaCount, totalDeltaSize) {
			return true
		}
	}
	return false
}
