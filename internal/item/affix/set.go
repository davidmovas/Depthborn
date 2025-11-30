package affix

import (
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
)

const (
	DefaultMaxPrefixes = 3
	DefaultMaxSuffixes = 3
)

type affixSet struct {
	mu          sync.RWMutex
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
	as.mu.Lock()
	defer as.mu.Unlock()

	if !as.canAddInternal(affix) {
		return fmt.Errorf("cannot add affix: %s", affix.ID())
	}
	as.affixes[affix.ID()] = affix
	return nil
}

func (as *affixSet) Remove(affixID string) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if _, exists := as.affixes[affixID]; !exists {
		return fmt.Errorf("affix not found: %s", affixID)
	}
	delete(as.affixes, affixID)
	return nil
}

func (as *affixSet) Get(affixID string) (Affix, bool) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	affix, exists := as.affixes[affixID]
	return affix, exists
}

func (as *affixSet) GetByType(affixType Type) []Affix {
	as.mu.RLock()
	defer as.mu.RUnlock()

	result := make([]Affix, 0)
	for _, affix := range as.affixes {
		if affix.Type() == affixType {
			result = append(result, affix)
		}
	}
	return result
}

func (as *affixSet) GetAll() []Affix {
	as.mu.RLock()
	defer as.mu.RUnlock()

	result := make([]Affix, 0, len(as.affixes))
	for _, affix := range as.affixes {
		result = append(result, affix)
	}
	return result
}

func (as *affixSet) Count() int {
	as.mu.RLock()
	defer as.mu.RUnlock()

	return len(as.affixes)
}

func (as *affixSet) CountByType(affixType Type) int {
	as.mu.RLock()
	defer as.mu.RUnlock()

	return as.countByTypeInternal(affixType)
}

// countByTypeInternal counts affixes by type (no lock)
func (as *affixSet) countByTypeInternal(affixType Type) int {
	count := 0
	for _, affix := range as.affixes {
		if affix.Type() == affixType {
			count++
		}
	}
	return count
}

func (as *affixSet) CanAdd(affix Affix) bool {
	as.mu.RLock()
	defer as.mu.RUnlock()

	return as.canAddInternal(affix)
}

// canAddInternal checks if affix can be added (no lock)
func (as *affixSet) canAddInternal(affix Affix) bool {
	switch affix.Type() {
	case TypePrefix:
		return as.countByTypeInternal(TypePrefix) < as.maxPrefixes
	case TypeSuffix:
		return as.countByTypeInternal(TypeSuffix) < as.maxSuffixes
	default:
		return true
	}
}

func (as *affixSet) MaxPrefixes() int {
	as.mu.RLock()
	defer as.mu.RUnlock()

	return as.maxPrefixes
}

func (as *affixSet) MaxSuffixes() int {
	as.mu.RLock()
	defer as.mu.RUnlock()

	return as.maxSuffixes
}

func (as *affixSet) Clear() {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.affixes = make(map[string]Affix)
}

func (as *affixSet) AllModifiers() []attribute.Modifier {
	as.mu.RLock()
	defer as.mu.RUnlock()

	modifiers := make([]attribute.Modifier, 0)
	for _, affix := range as.affixes {
		modifiers = append(modifiers, affix.Modifiers()...)
	}
	return modifiers
}

func (as *affixSet) SetMaxPrefixes(max int) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if max < 0 {
		max = 0
	}
	as.maxPrefixes = max
}

func (as *affixSet) SetMaxSuffixes(max int) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if max < 0 {
		max = 0
	}
	as.maxSuffixes = max
}
