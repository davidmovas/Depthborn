package status

import (
	"context"
	"fmt"
	"sync"
)

var _ Manager = (*BaseManager)(nil)

type BaseManager struct {
	effects    map[string]Effect
	immunities map[string]bool

	mu sync.RWMutex
}

func NewManager() *BaseManager {
	return &BaseManager{
		effects:    make(map[string]Effect),
		immunities: make(map[string]bool),
	}
}

func (m *BaseManager) Apply(ctx context.Context, effect Effect) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	effectType := effect.Type()

	// Check immunity
	if m.immunities[effectType] {
		return fmt.Errorf("immune to effect type: %s", effectType)
	}

	// Check if effect can stack with existing
	existingEffects := m.getByTypeLocked(effectType)
	for _, existing := range existingEffects {
		if effect.CanStack(existing) {
			// Stack with existing effect
			if existing.AddStack() {
				if err := existing.OnStack(ctx, effect.TargetID(), existing.Stacks()); err != nil {
					return fmt.Errorf("failed to stack effect: %w", err)
				}
				// Refresh duration if new effect has longer duration
				if effect.Duration() > existing.Duration() {
					existing.SetDuration(effect.Duration())
				}
				return nil
			}
			// Max stacks reached, replace if new effect has longer duration
			if effect.Duration() > existing.Duration() {
				existing.SetDuration(effect.Duration())
			}
			return nil
		}
	}

	// Add new effect
	m.effects[effect.ID()] = effect

	// Trigger OnApply
	if err := effect.OnApply(ctx, effect.TargetID()); err != nil {
		delete(m.effects, effect.ID())
		return fmt.Errorf("failed to apply effect: %w", err)
	}

	return nil
}

func (m *BaseManager) Remove(ctx context.Context, effectID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	effect, exists := m.effects[effectID]
	if !exists {
		return fmt.Errorf("effect not found: %s", effectID)
	}

	// Trigger OnRemove
	if err := effect.OnRemove(ctx, effect.TargetID()); err != nil {
		return fmt.Errorf("failed to remove effect: %w", err)
	}

	delete(m.effects, effectID)
	return nil
}

func (m *BaseManager) RemoveByType(ctx context.Context, effectType string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	for id, effect := range m.effects {
		if effect.Type() == effectType {
			if err := effect.OnRemove(ctx, effect.TargetID()); err != nil {
				errors = append(errors, fmt.Errorf("failed to remove effect %s: %w", id, err))
				continue
			}
			delete(m.effects, id)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors removing effects: %v", errors)
	}

	return nil
}

func (m *BaseManager) RemoveAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	for id, effect := range m.effects {
		if err := effect.OnRemove(ctx, effect.TargetID()); err != nil {
			errors = append(errors, fmt.Errorf("failed to remove effect %s: %w", id, err))
		}
	}

	m.effects = make(map[string]Effect)

	if len(errors) > 0 {
		return fmt.Errorf("errors removing all effects: %v", errors)
	}

	return nil
}

func (m *BaseManager) Has(effectType string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, effect := range m.effects {
		if effect.Type() == effectType {
			return true
		}
	}
	return false
}

func (m *BaseManager) Get(effectID string) (Effect, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	effect, exists := m.effects[effectID]
	return effect, exists
}

func (m *BaseManager) GetByType(effectType string) []Effect {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getByTypeLocked(effectType)
}

func (m *BaseManager) getByTypeLocked(effectType string) []Effect {
	var results []Effect
	for _, effect := range m.effects {
		if effect.Type() == effectType {
			results = append(results, effect)
		}
	}
	return results
}

func (m *BaseManager) GetAll() []Effect {
	m.mu.RLock()
	defer m.mu.RUnlock()

	effects := make([]Effect, 0, len(m.effects))
	for _, effect := range m.effects {
		effects = append(effects, effect)
	}
	return effects
}

func (m *BaseManager) Update(ctx context.Context, deltaMs int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var expiredIDs []string
	var errors []error

	// Process all effects
	for id, effect := range m.effects {
		// Update effect
		if err := effect.OnTick(ctx, effect.TargetID(), deltaMs); err != nil {
			errors = append(errors, fmt.Errorf("error updating effect %s: %w", id, err))
		}

		// Check if expired
		if effect.IsExpired() {
			expiredIDs = append(expiredIDs, id)
		}
	}

	// Remove expired effects
	for _, id := range expiredIDs {
		effect := m.effects[id]
		if err := effect.OnRemove(ctx, effect.TargetID()); err != nil {
			errors = append(errors, fmt.Errorf("error removing expired effect %s: %w", id, err))
		}
		delete(m.effects, id)
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during update: %v", errors)
	}

	return nil
}

func (m *BaseManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.effects)
}

func (m *BaseManager) IsImmune(effectType string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.immunities[effectType]
}

func (m *BaseManager) AddImmunity(effectType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.immunities[effectType] = true
}

func (m *BaseManager) RemoveImmunity(effectType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.immunities, effectType)
}
