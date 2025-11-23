package affix

import (
	"fmt"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
)

const (
	DefaultMaxPrefixes = 3
	DefaultMaxSuffixes = 3
)

type affixSet struct {
	affixes     map[string]Affix
	maxPrefixes int
	maxSuffixes int
}

func NewAffixSet() Set {
	return &affixSet{
		affixes:     make(map[string]Affix),
		maxPrefixes: DefaultMaxPrefixes,
		maxSuffixes: DefaultMaxSuffixes,
	}
}

func (as *affixSet) Add(affix Affix) error {
	if !as.CanAdd(affix) {
		return fmt.Errorf("cannot add affix: %s", affix.ID())
	}
	as.affixes[affix.ID()] = affix
	return nil
}

func (as *affixSet) Remove(affixID string) error {
	if _, exists := as.affixes[affixID]; !exists {
		return fmt.Errorf("affix not found: %s", affixID)
	}
	delete(as.affixes, affixID)
	return nil
}

func (as *affixSet) Get(affixID string) (Affix, bool) {
	affix, exists := as.affixes[affixID]
	return affix, exists
}

func (as *affixSet) GetByType(affixType Type) []Affix {
	result := make([]Affix, 0)
	for _, affix := range as.affixes {
		if affix.Type() == affixType {
			result = append(result, affix)
		}
	}
	return result
}

func (as *affixSet) GetAll() []Affix {
	result := make([]Affix, 0, len(as.affixes))
	for _, affix := range as.affixes {
		result = append(result, affix)
	}
	return result
}

func (as *affixSet) Count() int {
	return len(as.affixes)
}

func (as *affixSet) CountByType(affixType Type) int {
	count := 0
	for _, affix := range as.affixes {
		if affix.Type() == affixType {
			count++
		}
	}
	return count
}

func (as *affixSet) CanAdd(affix Affix) bool {
	switch affix.Type() {
	case TypePrefix:
		return as.CountByType(TypePrefix) < as.maxPrefixes
	case TypeSuffix:
		return as.CountByType(TypeSuffix) < as.maxSuffixes
	default:
		return true
	}
}

func (as *affixSet) MaxPrefixes() int {
	return as.maxPrefixes
}

func (as *affixSet) MaxSuffixes() int {
	return as.maxSuffixes
}

func (as *affixSet) Clear() {
	as.affixes = make(map[string]Affix)
}

func (as *affixSet) AllModifiers() []attribute.Modifier {
	modifiers := make([]attribute.Modifier, 0)
	for _, affix := range as.affixes {
		modifiers = append(modifiers, affix.Modifiers()...)
	}
	return modifiers
}

func (as *affixSet) SetMaxPrefixes(max int) {
	if max < 0 {
		max = 0
	}
	as.maxPrefixes = max
}

func (as *affixSet) SetMaxSuffixes(max int) {
	if max < 0 {
		max = 0
	}
	as.maxSuffixes = max
}
