package attribute

import "sort"

var _ Set = (*BaseSet)(nil)

type BaseSet struct {
	modifiers map[string]Modifier
}

func NewSet() *BaseSet {
	return &BaseSet{
		modifiers: make(map[string]Modifier),
	}
}

func (s *BaseSet) Add(modifier Modifier) {
	s.modifiers[modifier.ID()] = modifier
}

func (s *BaseSet) Remove(modifierID string) {
	delete(s.modifiers, modifierID)
}

func (s *BaseSet) GetAll() []Modifier {
	result := make([]Modifier, 0, len(s.modifiers))
	for _, mod := range s.modifiers {
		if mod.IsActive() {
			result = append(result, mod)
		}
	}
	return result
}

func (s *BaseSet) GetByType(modType ModifierType) []Modifier {
	result := make([]Modifier, 0)
	for _, mod := range s.modifiers {
		if mod.IsActive() && mod.Type() == modType {
			result = append(result, mod)
		}
	}
	return result
}

func (s *BaseSet) Clear() {
	s.modifiers = make(map[string]Modifier)
}

func (s *BaseSet) Apply(baseValue float64) float64 {
	mods := s.GetAll()

	sort.Slice(mods, func(i, j int) bool {
		return mods[i].Priority() > mods[j].Priority()
	})

	for _, mod := range mods {
		if mod.Type() == ModOverride {
			return mod.Value()
		}
	}

	value := baseValue

	for _, mod := range mods {
		if mod.Type() == ModFlat {
			value += mod.Value()
		}
	}

	increasedSum := 0.0
	for _, mod := range mods {
		if mod.Type() == ModIncreased {
			increasedSum += mod.Value()
		}
	}
	if increasedSum != 0 {
		value *= 1.0 + increasedSum/100.0
	}

	for _, mod := range mods {
		if mod.Type() == ModMore {
			value *= 1.0 + mod.Value()/100.0
		}
	}

	return value
}
