package attribute

var _ Manager = (*attributeManager)(nil)

type attributeManager struct {
	baseValues    map[Type]float64
	currentValues map[Type]float64
	modifiers     map[Type]map[string]Modifier
	formulas      map[Type]Formula
}

func NewAttributeManager() Manager {
	return &attributeManager{
		baseValues:    make(map[Type]float64),
		currentValues: make(map[Type]float64),
		modifiers:     make(map[Type]map[string]Modifier),
		formulas:      make(map[Type]Formula),
	}
}

func (am *attributeManager) Get(attr Type) float64 {
	if value, exists := am.currentValues[attr]; exists {
		return value
	}
	return 0.0
}

func (am *attributeManager) GetBase(attr Type) float64 {
	if value, exists := am.baseValues[attr]; exists {
		return value
	}
	return 0.0
}

func (am *attributeManager) SetBase(attr Type, value float64) {
	am.baseValues[attr] = value
	am.recalculateAttribute(attr)
}

func (am *attributeManager) AddModifier(attr Type, modifier Modifier) {
	if _, exists := am.modifiers[attr]; !exists {
		am.modifiers[attr] = make(map[string]Modifier)
	}
	am.modifiers[attr][modifier.ID()] = modifier
	am.recalculateAttribute(attr)
}

func (am *attributeManager) RemoveModifier(attr Type, modifierID string) {
	if modifiers, exists := am.modifiers[attr]; exists {
		delete(modifiers, modifierID)
		am.recalculateAttribute(attr)
	}
}

func (am *attributeManager) RemoveAllModifiers(attr Type, modType ModifierType) {
	if modifiers, exists := am.modifiers[attr]; exists {
		for id, mod := range modifiers {
			if mod.Type() == modType {
				delete(modifiers, id)
			}
		}
		am.recalculateAttribute(attr)
	}
}

func (am *attributeManager) GetModifiers(attr Type) []Modifier {
	if modifiers, exists := am.modifiers[attr]; exists {
		result := make([]Modifier, 0, len(modifiers))
		for _, mod := range modifiers {
			result = append(result, mod)
		}
		return result
	}
	return nil
}

func (am *attributeManager) RecalculateAll() {
	for attr := range am.baseValues {
		am.recalculateAttribute(attr)
	}
}

func (am *attributeManager) Snapshot() map[Type]float64 {
	snapshot := make(map[Type]float64)
	for attr, value := range am.currentValues {
		snapshot[attr] = value
	}
	return snapshot
}

func (am *attributeManager) recalculateAttribute(attr Type) {
	baseValue := am.baseValues[attr]

	if formula, hasFormula := am.formulas[attr]; hasFormula {
		baseValue = formula.Calculate(am)
	} else {
		if modifiers, exists := am.modifiers[attr]; exists {
			set := NewModifierSet()
			for _, mod := range modifiers {
				set.Add(mod)
			}
			baseValue = set.Apply(baseValue)
		}
	}

	am.currentValues[attr] = baseValue
}

func (am *attributeManager) SetFormula(attr Type, formula Formula) {
	am.formulas[attr] = formula
	am.recalculateAttribute(attr)
}
