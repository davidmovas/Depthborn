package progression

import (
	"context"
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/pkg/state"
)

var _ StatPointManager = (*BaseStatPointManager)(nil)

type BaseStatPointManager struct {
	availablePoints int
	spentPoints     int
	allocations     map[string]int
	pointsPerLevel  int

	allocatedCallbacks   []StatPointCallback
	deallocatedCallbacks []StatPointCallback

	mu sync.RWMutex
}

type StatPointConfig struct {
	InitialPoints  int
	PointsPerLevel int
	Allocations    map[string]int
}

func NewStatPointManager(config StatPointConfig) *BaseStatPointManager {
	if config.PointsPerLevel <= 0 {
		config.PointsPerLevel = 5 // Default 5 points per level
	}

	allocations := config.Allocations
	if allocations == nil {
		allocations = make(map[string]int)
	}

	spentPoints := 0
	for _, points := range allocations {
		spentPoints += points
	}

	return &BaseStatPointManager{
		availablePoints:      config.InitialPoints,
		spentPoints:          spentPoints,
		allocations:          allocations,
		pointsPerLevel:       config.PointsPerLevel,
		allocatedCallbacks:   make([]StatPointCallback, 0),
		deallocatedCallbacks: make([]StatPointCallback, 0),
	}
}

func (m *BaseStatPointManager) AvailablePoints() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.availablePoints
}

func (m *BaseStatPointManager) SpentPoints() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.spentPoints
}

func (m *BaseStatPointManager) AddPoints(amount int) {
	if amount <= 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.availablePoints += amount
}

func (m *BaseStatPointManager) AllocatePoint(attributeType string) error {
	m.mu.Lock()

	if m.availablePoints <= 0 {
		m.mu.Unlock()
		return fmt.Errorf("no available stat points")
	}

	if attributeType == "" {
		m.mu.Unlock()
		return fmt.Errorf("attribute type cannot be empty")
	}

	m.availablePoints--
	m.spentPoints++
	m.allocations[attributeType]++

	// Copy callbacks to avoid deadlock
	callbacks := make([]StatPointCallback, len(m.allocatedCallbacks))
	copy(callbacks, m.allocatedCallbacks)

	m.mu.Unlock()

	// Trigger callbacks
	ctx := context.Background()
	for _, callback := range callbacks {
		callback(ctx, attributeType, m)
	}

	return nil
}

func (m *BaseStatPointManager) DeallocatePoint(attributeType string) error {
	m.mu.Lock()

	if attributeType == "" {
		m.mu.Unlock()
		return fmt.Errorf("attribute type cannot be empty")
	}

	allocated, exists := m.allocations[attributeType]
	if !exists || allocated <= 0 {
		m.mu.Unlock()
		return fmt.Errorf("no points allocated to %s", attributeType)
	}

	m.allocations[attributeType]--
	if m.allocations[attributeType] == 0 {
		delete(m.allocations, attributeType)
	}

	m.availablePoints++
	m.spentPoints--

	// Copy callbacks to avoid deadlock
	callbacks := make([]StatPointCallback, len(m.deallocatedCallbacks))
	copy(callbacks, m.deallocatedCallbacks)

	m.mu.Unlock()

	// Trigger callbacks
	ctx := context.Background()
	for _, callback := range callbacks {
		callback(ctx, attributeType, m)
	}

	return nil
}

func (m *BaseStatPointManager) GetAllocations() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return copy to prevent external modification
	allocations := make(map[string]int, len(m.allocations))
	for attr, points := range m.allocations {
		allocations[attr] = points
	}

	return allocations
}

func (m *BaseStatPointManager) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Refund all spent points
	m.availablePoints += m.spentPoints
	m.spentPoints = 0
	m.allocations = make(map[string]int)

	return nil
}

func (m *BaseStatPointManager) CanAllocate(attributeType string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if attributeType == "" {
		return false
	}

	return m.availablePoints > 0
}

func (m *BaseStatPointManager) PointsPerLevel() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pointsPerLevel
}

func (m *BaseStatPointManager) SetPointsPerLevel(points int) {
	if points < 0 {
		points = 0
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.pointsPerLevel = points
}

func (m *BaseStatPointManager) OnPointAllocated(callback StatPointCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.allocatedCallbacks = append(m.allocatedCallbacks, callback)
}

func (m *BaseStatPointManager) OnPointDeallocated(callback StatPointCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deallocatedCallbacks = append(m.deallocatedCallbacks, callback)
}

func (m *BaseStatPointManager) SerializeState() (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return state.New().
		Set("available_points", m.availablePoints).
		Set("spent_points", m.spentPoints).
		Set("allocations", m.GetAllocations()).
		Set("points_per_level", m.pointsPerLevel).Data(), nil
}

func (m *BaseStatPointManager) DeserializeState(stateData map[string]any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := state.From(stateData)

	m.availablePoints = s.IntOr("available_points", 0)
	m.spentPoints = s.IntOr("spent_points", 0)
	m.pointsPerLevel = s.IntOr("points_per_level", 1)

	if allocations, ok := state.MapTyped[string, int]("allocations", s); ok {
		m.allocations = allocations
	} else {
		m.allocations = make(map[string]int)
	}

	return nil
}
