package progression

import (
	"context"
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/pkg/state"
)

var _ ExperienceManager = (*BaseExperienceManager)(nil)

type BaseExperienceManager struct {
	currentLevel      int
	currentExperience int64
	maxLevel          int
	curve             ExperienceCurve

	levelUpCallbacks    []LevelUpCallback
	experienceCallbacks []ExperienceCallback

	mu sync.RWMutex
}
type ExperienceConfig struct {
	InitialLevel int
	InitialXP    int64
	MaxLevel     int
	Curve        ExperienceCurve
}

func NewExperienceManager(config ExperienceConfig) *BaseExperienceManager {
	if config.InitialLevel < 1 {
		config.InitialLevel = 1
	}

	if config.MaxLevel <= 0 {
		config.MaxLevel = 100 // Default max level
	}

	if config.Curve == nil {
		config.Curve = NewStandardExponentialCurve()
	}

	return &BaseExperienceManager{
		currentLevel:        config.InitialLevel,
		currentExperience:   config.InitialXP,
		maxLevel:            config.MaxLevel,
		curve:               config.Curve,
		levelUpCallbacks:    make([]LevelUpCallback, 0),
		experienceCallbacks: make([]ExperienceCallback, 0),
	}
}

func (m *BaseExperienceManager) CurrentLevel() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentLevel
}

func (m *BaseExperienceManager) CurrentExperience() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentExperience
}

func (m *BaseExperienceManager) ExperienceToNextLevel() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentLevel >= m.maxLevel {
		return 0
	}

	xpForNextLevel := m.curve.ExperienceForLevel(m.currentLevel + 1)

	return xpForNextLevel - m.currentExperience
}

func (m *BaseExperienceManager) ExperienceForLevel(level int) int64 {
	return m.curve.ExperienceForLevel(level)
}

func (m *BaseExperienceManager) AddExperience(ctx context.Context, amount int64) (int, error) {
	if amount < 0 {
		return 0, fmt.Errorf("experience amount must be positive")
	}

	m.mu.Lock()

	if m.currentLevel >= m.maxLevel {
		m.mu.Unlock()
		return 0, fmt.Errorf("already at max level")
	}

	oldLevel := m.currentLevel
	m.currentExperience += amount

	// Calculate new level based on total XP
	newLevel := m.curve.LevelForExperience(m.currentExperience)

	// Cap at max level
	if newLevel > m.maxLevel {
		newLevel = m.maxLevel
		m.currentExperience = m.curve.ExperienceForLevel(m.maxLevel)
	}

	levelsGained := newLevel - oldLevel
	m.currentLevel = newLevel

	// Copy callbacks to avoid deadlock
	levelUpCallbacks := make([]LevelUpCallback, len(m.levelUpCallbacks))
	copy(levelUpCallbacks, m.levelUpCallbacks)

	experienceCallbacks := make([]ExperienceCallback, len(m.experienceCallbacks))
	copy(experienceCallbacks, m.experienceCallbacks)

	m.mu.Unlock()

	// Trigger experience gain callbacks
	for _, callback := range experienceCallbacks {
		callback(ctx, amount, m)
	}

	// Trigger level up callbacks for each level gained
	if levelsGained > 0 {
		for level := oldLevel + 1; level <= newLevel; level++ {
			for _, callback := range levelUpCallbacks {
				callback(ctx, level-1, level, m)
			}
		}
	}

	return levelsGained, nil
}

func (m *BaseExperienceManager) SetExperience(xp int64) error {
	if xp < 0 {
		return fmt.Errorf("experience cannot be negative")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.currentExperience = xp

	// Recalculate level
	newLevel := m.curve.LevelForExperience(xp)
	if newLevel > m.maxLevel {
		newLevel = m.maxLevel
		m.currentExperience = m.curve.ExperienceForLevel(m.maxLevel)
	}

	m.currentLevel = newLevel

	return nil
}

func (m *BaseExperienceManager) SetLevel(level int) error {
	if level < 1 {
		return fmt.Errorf("level must be at least 1")
	}

	if level > m.maxLevel {
		return fmt.Errorf("level cannot exceed max level %d", m.maxLevel)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.currentLevel = level
	m.currentExperience = m.curve.ExperienceForLevel(level)

	return nil
}

func (m *BaseExperienceManager) Progress() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentLevel >= m.maxLevel {
		return 1.0
	}

	xpForCurrentLevel := m.curve.ExperienceForLevel(m.currentLevel)
	xpForNextLevel := m.curve.ExperienceForLevel(m.currentLevel + 1)

	xpIntoLevel := m.currentExperience - xpForCurrentLevel
	xpNeededForLevel := xpForNextLevel - xpForCurrentLevel

	if xpNeededForLevel <= 0 {
		return 0.0
	}

	progress := float64(xpIntoLevel) / float64(xpNeededForLevel)

	// Clamp between 0 and 1
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	return progress
}

func (m *BaseExperienceManager) MaxLevel() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.maxLevel
}

func (m *BaseExperienceManager) IsMaxLevel() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentLevel >= m.maxLevel
}

func (m *BaseExperienceManager) OnLevelUp(callback LevelUpCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.levelUpCallbacks = append(m.levelUpCallbacks, callback)
}

func (m *BaseExperienceManager) OnExperienceGain(callback ExperienceCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.experienceCallbacks = append(m.experienceCallbacks, callback)
}

func (m *BaseExperienceManager) SerializeState() (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s := state.New().
		Set("current_level", m.currentLevel).
		Set("current_experience", m.currentExperience).
		Set("max_level", m.maxLevel)

	if m.curve != nil {
		s.Set("curve_type", string(m.curve.Type())).
			Set("curve_parameters", m.curve.Parameters())
	}

	return s.Data(), nil
}

func (m *BaseExperienceManager) DeserializeState(stateData map[string]any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := state.From(stateData)

	m.currentLevel = s.IntOr("current_level", 1)
	m.currentExperience = int64(s.IntOr("current_experience", 0))
	m.maxLevel = s.IntOr("max_level", 100)

	if m.currentLevel < 1 {
		m.currentLevel = 1
	}
	if m.currentExperience < 0 {
		m.currentExperience = 0
	}
	if m.maxLevel < 1 {
		m.maxLevel = 100
	}

	if curveType := s.StringOr("curve_type", string(CurveExponential)); m.curve == nil || string(m.curve.Type()) != curveType {
		switch CurveType(curveType) {
		case CurveLinear:
			m.curve = NewStandardLinearCurve()
		case CurveExponential:
			m.curve = NewStandardExponentialCurve()
		case CurveLogarithmic:
			m.curve = NewStandardLogarithmicCurve()
		case CurvePolynomial:
			m.curve = NewStandardPolynomialCurve()
		case CurveCustom:
		default:
			m.curve = NewStandardExponentialCurve()
		}
	}

	return nil
}
