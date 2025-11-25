package attribute

import "sync"

var _ Manager = (*BaseManager)(nil)

type BaseManager struct {
	mu sync.RWMutex

	baseValues map[Type]float64
	modifiers  map[Type]Set
	formulas   map[Type]Formula
	cache      map[Type]float64
	dirty      map[Type]bool
}

func NewManager() *BaseManager {
	return &BaseManager{
		baseValues: make(map[Type]float64),
		modifiers:  make(map[Type]Set),
		formulas:   make(map[Type]Formula),
		cache:      make(map[Type]float64),
		dirty:      make(map[Type]bool),
	}
}

func (m *BaseManager) Get(attr Type) float64 {
	m.mu.RLock()
	if !m.dirty[attr] {
		if cached, ok := m.cache[attr]; ok {
			m.mu.RUnlock()
			return cached
		}
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	value := m.calculate(attr)
	m.cache[attr] = value
	m.dirty[attr] = false

	return value
}

func (m *BaseManager) GetBase(attr Type) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.baseValues[attr]
}

func (m *BaseManager) SetBase(attr Type, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.baseValues[attr] = value
	m.markDirty(attr)
}

func (m *BaseManager) AddModifier(attr Type, modifier Modifier) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.modifiers[attr] == nil {
		m.modifiers[attr] = NewSet()
	}

	m.modifiers[attr].Add(modifier)
	m.markDirty(attr)
}

func (m *BaseManager) RemoveModifier(attr Type, modifierID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if set, ok := m.modifiers[attr]; ok {
		set.Remove(modifierID)
		m.markDirty(attr)
	}
}

func (m *BaseManager) RemoveAllModifiers(attr Type, modType ModifierType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if set, ok := m.modifiers[attr]; ok {
		for _, mod := range set.GetByType(modType) {
			set.Remove(mod.ID())
		}
		m.markDirty(attr)
	}
}

func (m *BaseManager) GetModifiers(attr Type) []Modifier {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if set, ok := m.modifiers[attr]; ok {
		return set.GetAll()
	}
	return nil
}

func (m *BaseManager) RecalculateAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for attr := range m.baseValues {
		m.dirty[attr] = true
	}

	for attr := range m.formulas {
		m.dirty[attr] = true
	}

	m.cache = make(map[Type]float64)
}

func (m *BaseManager) Snapshot() map[Type]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snapshot := make(map[Type]float64)

	for attr := range m.baseValues {
		m.mu.RUnlock()
		snapshot[attr] = m.Get(attr)
		m.mu.RLock()
	}

	for attr := range m.formulas {
		if _, exists := snapshot[attr]; !exists {
			m.mu.RUnlock()
			snapshot[attr] = m.Get(attr)
			m.mu.RLock()
		}
	}

	return snapshot
}

func (m *BaseManager) Restore(snapshot map[Type]float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	existingModifiers := m.modifiers
	existingFormulas := m.formulas

	m.baseValues = make(map[Type]float64)
	m.modifiers = existingModifiers
	m.formulas = existingFormulas
	m.cache = make(map[Type]float64)
	m.dirty = make(map[Type]bool)

	for attr, value := range snapshot {
		m.baseValues[attr] = value
		m.dirty[attr] = true
	}

	for attr := range m.formulas {
		m.dirty[attr] = true
	}
}

func (m *BaseManager) SetFormula(attr Type, formula Formula) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.formulas[attr] = formula
	m.markDirty(attr)
}

func (m *BaseManager) calculate(attr Type) float64 {
	if formula, hasFormula := m.formulas[attr]; hasFormula {
		return formula.Calculate(m)
	}

	base := m.baseValues[attr]

	if set, ok := m.modifiers[attr]; ok {
		return set.Apply(base)
	}

	return base
}

func (m *BaseManager) markDirty(attr Type) {
	m.dirty[attr] = true

	for otherAttr, formula := range m.formulas {
		for _, dep := range formula.Dependencies() {
			if dep == attr {
				m.dirty[otherAttr] = true
				break
			}
		}
	}
}
