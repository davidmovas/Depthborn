package attribute

var _ Set = (*modifierSet)(nil)

type modifierSet struct {
	modifiers map[string]Modifier
}

func NewModifierSet() Set {
	return &modifierSet{
		modifiers: make(map[string]Modifier),
	}
}

func (ms *modifierSet) Add(modifier Modifier) {
	ms.modifiers[modifier.ID()] = modifier
}

func (ms *modifierSet) Remove(modifierID string) {
	delete(ms.modifiers, modifierID)
}

func (ms *modifierSet) GetAll() []Modifier {
	result := make([]Modifier, 0, len(ms.modifiers))
	for _, mod := range ms.modifiers {
		result = append(result, mod)
	}
	return result
}

func (ms *modifierSet) GetByType(modType ModifierType) []Modifier {
	result := make([]Modifier, 0)
	for _, mod := range ms.modifiers {
		if mod.Type() == modType {
			result = append(result, mod)
		}
	}
	return result
}

func (ms *modifierSet) Clear() {
	ms.modifiers = make(map[string]Modifier)
}

func (ms *modifierSet) Apply(baseValue float64) float64 {
	modifiers := ms.GetAll()

	// Sort by priority (higher first)
	// TODO: Implement proper sorting by priority

	result := baseValue
	var increasedSum float64
	var moreProduct float64 = 1.0

	for _, mod := range modifiers {
		if !mod.IsActive() {
			continue
		}

		switch mod.Type() {
		case ModFlat:
			result += mod.Value()
		case ModIncreased:
			increasedSum += mod.Value()
		case ModMore:
			moreProduct *= 1.0 + mod.Value()
		case ModOverride:
			result = mod.Value()
			increasedSum = 0
			moreProduct = 1.0
		}
	}

	// Apply increased (additive)
	result *= 1.0 + increasedSum

	// Apply more (multiplicative)
	result *= moreProduct

	return result
}
