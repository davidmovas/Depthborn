package progression

import (
	"context"
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/pkg/state"
)

var _ SkillPointManager = (*BaseSkillPointManager)(nil)

type BaseSkillPointManager struct {
	availablePoints int
	spentPoints     int
	totalPoints     int

	pointsGainedCallbacks []SkillPointCallback
	pointsSpentCallbacks  []SkillPointCallback

	mu sync.RWMutex
}

type SkillPointConfig struct {
	InitialPoints int
}

func NewSkillPointManager(config SkillPointConfig) *BaseSkillPointManager {
	return &BaseSkillPointManager{
		availablePoints:       config.InitialPoints,
		spentPoints:           0,
		totalPoints:           config.InitialPoints,
		pointsGainedCallbacks: make([]SkillPointCallback, 0),
		pointsSpentCallbacks:  make([]SkillPointCallback, 0),
	}
}

func (m *BaseSkillPointManager) AvailablePoints() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.availablePoints
}

func (m *BaseSkillPointManager) SpentPoints() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.spentPoints
}

func (m *BaseSkillPointManager) TotalPoints() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.totalPoints
}

func (m *BaseSkillPointManager) AddPoints(amount int) {
	if amount <= 0 {
		return
	}

	m.mu.Lock()

	m.availablePoints += amount
	m.totalPoints += amount

	// Copy callbacks to avoid deadlock
	callbacks := make([]SkillPointCallback, len(m.pointsGainedCallbacks))
	copy(callbacks, m.pointsGainedCallbacks)

	m.mu.Unlock()

	// Trigger callbacks
	ctx := context.Background()
	for _, callback := range callbacks {
		callback(ctx, amount, m)
	}
}

func (m *BaseSkillPointManager) SpendPoints(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	m.mu.Lock()

	if m.availablePoints < amount {
		m.mu.Unlock()
		return fmt.Errorf("insufficient skill points: have %d, need %d", m.availablePoints, amount)
	}

	m.availablePoints -= amount
	m.spentPoints += amount

	// Copy callbacks to avoid deadlock
	callbacks := make([]SkillPointCallback, len(m.pointsSpentCallbacks))
	copy(callbacks, m.pointsSpentCallbacks)

	m.mu.Unlock()

	// Trigger callbacks
	ctx := context.Background()
	for _, callback := range callbacks {
		callback(ctx, amount, m)
	}

	return nil
}

func (m *BaseSkillPointManager) RefundPoints(amount int) {
	if amount <= 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Can't refund more than spent
	if amount > m.spentPoints {
		amount = m.spentPoints
	}

	m.availablePoints += amount
	m.spentPoints -= amount
}

func (m *BaseSkillPointManager) PointsForLevel(level int) int {
	// Simple formula: 1 skill point every 5 levels
	// Can be customized in the future
	if level%5 == 0 {
		return 1
	}
	return 0
}

func (m *BaseSkillPointManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Refund all spent points
	m.availablePoints += m.spentPoints
	m.spentPoints = 0
}

func (m *BaseSkillPointManager) OnPointsGained(callback SkillPointCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pointsGainedCallbacks = append(m.pointsGainedCallbacks, callback)
}

func (m *BaseSkillPointManager) OnPointsSpent(callback SkillPointCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pointsSpentCallbacks = append(m.pointsSpentCallbacks, callback)
}

func (m *BaseSkillPointManager) SerializeState() (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return state.BatchKV(
		"available_points", m.availablePoints,
		"spent_points", m.spentPoints,
		"total_points", m.totalPoints,
	)
}

func (m *BaseSkillPointManager) DeserializeState(stateData map[string]any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := state.From(stateData)

	m.availablePoints = s.IntOr("available_points", 0)
	m.spentPoints = s.IntOr("spent_points", 0)
	m.totalPoints = s.IntOr("total_points", 0)

	return nil
}
