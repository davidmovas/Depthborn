package progression

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/davidmovas/Depthborn/pkg/state"
)

var _ Manager = (*BaseManager)(nil)

type BaseManager struct {
	experience      ExperienceManager
	statPoints      StatPointManager
	skillPoints     SkillPointManager
	attributeGrowth AttributeGrowth
	rewardTable     RewardTable
	tracker         Tracker
	prestige        PrestigeManager
	difficulty      DifficultyScaler

	mu sync.RWMutex
}

type ManagerConfig struct {
	Experience      ExperienceManager
	StatPoints      StatPointManager
	SkillPoints     SkillPointManager
	AttributeGrowth AttributeGrowth
	RewardTable     RewardTable
	Tracker         Tracker
	Prestige        PrestigeManager
	Difficulty      DifficultyScaler
}

func NewManager(config ManagerConfig) *BaseManager {
	// Create default components if not provided
	if config.Experience == nil {
		config.Experience = NewExperienceManager(ExperienceConfig{
			InitialLevel: 1,
			InitialXP:    0,
			MaxLevel:     100,
			Curve:        NewStandardExponentialCurve(),
		})
	}

	if config.StatPoints == nil {
		config.StatPoints = NewStatPointManager(StatPointConfig{
			InitialPoints:  0,
			PointsPerLevel: 5,
		})
	}

	if config.SkillPoints == nil {
		config.SkillPoints = NewSkillPointManager(SkillPointConfig{
			InitialPoints: 0,
		})
	}

	if config.AttributeGrowth == nil {
		config.AttributeGrowth = NewAttributeGrowth(AttributeGrowthConfig{
			GrowthType: GrowthFlat,
			GrowthRates: map[string]float64{
				"vitality": 2.0, // +2 vitality per level
			},
		})
	}

	// TODO: Initialize other components when implemented
	// if config.RewardTable == nil { ... }
	// if config.Tracker == nil { ... }
	// if config.Prestige == nil { ... }
	// if config.Difficulty == nil { ... }

	manager := &BaseManager{
		experience:      config.Experience,
		statPoints:      config.StatPoints,
		skillPoints:     config.SkillPoints,
		attributeGrowth: config.AttributeGrowth,
		rewardTable:     config.RewardTable,
		tracker:         config.Tracker,
		prestige:        config.Prestige,
		difficulty:      config.Difficulty,
	}

	// Setup level up callback to grant stat/skill points
	config.Experience.OnLevelUp(func(ctx context.Context, oldLevel, newLevel int, _ ExperienceManager) {
		manager.handleLevelUp(ctx, oldLevel, newLevel)
	})

	return manager
}

func (m *BaseManager) Experience() ExperienceManager {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.experience
}

func (m *BaseManager) StatPoints() StatPointManager {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.statPoints
}

func (m *BaseManager) SkillPoints() SkillPointManager {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.skillPoints
}

func (m *BaseManager) AttributeGrowth() AttributeGrowth {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.attributeGrowth
}

func (m *BaseManager) RewardTable() RewardTable {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rewardTable
}

func (m *BaseManager) Tracker() Tracker {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tracker
}

func (m *BaseManager) Prestige() PrestigeManager {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.prestige
}

func (m *BaseManager) DifficultyScaler() DifficultyScaler {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.difficulty
}

func (m *BaseManager) ProcessLevelUp(ctx context.Context, characterID string, newLevel int) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Grant stat points
	pointsPerLevel := m.statPoints.PointsPerLevel()
	if pointsPerLevel > 0 {
		m.statPoints.AddPoints(pointsPerLevel)
	}

	// Grant skill points (if applicable for this level)
	skillPoints := m.skillPoints.PointsForLevel(newLevel)
	if skillPoints > 0 {
		m.skillPoints.AddPoints(skillPoints)
	}

	// Apply attribute growth
	// TODO: This requires integration with entity's attribute manager
	// attributes := entity.Attributes().Snapshot()
	// m.attributeGrowth.ApplyLevelUp(ctx, newLevel, attributes)

	// Claim level rewards
	if m.rewardTable != nil && m.rewardTable.HasRewards(newLevel) {
		if err := m.rewardTable.ClaimRewards(ctx, newLevel, characterID); err != nil {
			return fmt.Errorf("failed to claim level rewards: %w", err)
		}
	}

	// Record level up in tracker
	if m.tracker != nil {
		m.tracker.RecordLevelUp(newLevel, 0) // TODO: Get actual timestamp
	}

	return nil
}

func (m *BaseManager) Reset(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Reset all components
	if err := m.experience.SetLevel(1); err != nil {
		return fmt.Errorf("failed to reset experience: %w", err)
	}

	if err := m.statPoints.Reset(); err != nil {
		return fmt.Errorf("failed to reset stat points: %w", err)
	}

	m.skillPoints.Reset()

	if m.tracker != nil {
		m.tracker.Reset()
	}

	return nil
}

func (m *BaseManager) Save(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// TODO: Implement persistence
	// This should save all progression state to storage
	return nil
}

func (m *BaseManager) Load(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO: Implement persistence
	// This should load all progression state from storage
	return nil
}

// Internal handler for level up
func (m *BaseManager) handleLevelUp(ctx context.Context, oldLevel, newLevel int) {
	// Grant stat points
	pointsPerLevel := m.statPoints.PointsPerLevel()
	if pointsPerLevel > 0 {
		m.statPoints.AddPoints(pointsPerLevel)
	}

	// Grant skill points
	skillPoints := m.skillPoints.PointsForLevel(newLevel)
	if skillPoints > 0 {
		m.skillPoints.AddPoints(skillPoints)
	}

	// Record in tracker
	if m.tracker != nil {
		m.tracker.RecordLevelUp(newLevel, time.Now().Unix())
	}
}

func (m *BaseManager) SerializeState() (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var err error

	s := state.New()
	if err = s.SetEntity("experience", m.experience); err != nil {
		return nil, err
	}
	if err = s.SetEntity("stat_points", m.statPoints); err != nil {
		return nil, err
	}
	if err = s.SetEntity("skill_points", m.skillPoints); err != nil {
		return nil, err
	}
	if err = s.SetEntity("attribute_growth", m.attributeGrowth); err != nil {
		return nil, err
	}

	// TODO: Serialize other components when implemented
	return s.Data(), nil
}

func (m *BaseManager) DeserializeState(stateData map[string]any) error {
	s := state.From(stateData)

	var err error
	if err = s.GetEntity("experience", m.experience); err != nil {
		return err
	}
	if err = s.GetEntity("stat_points", m.statPoints); err != nil {
		return err
	}
	if err = s.GetEntity("skill_points", m.skillPoints); err != nil {
		return err
	}
	if err = s.GetEntity("attribute_growth", m.attributeGrowth); err != nil {
		return err
	}

	// TODO: Deserialize other components when implemented
	return nil
}
