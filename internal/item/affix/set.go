package affix

import (
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
)

var _ Set = (*BaseSet)(nil)

// AffixLimits defines prefix and suffix constraints
type AffixLimits struct {
	MinPrefixes int
	MaxPrefixes int
	MinSuffixes int
	MaxSuffixes int
}

// DefaultLimits returns default limits for different rarities
func DefaultLimits(rarity int) AffixLimits {
	switch rarity {
	case 0: // Common - no affixes
		return AffixLimits{0, 0, 0, 0}
	case 1: // Uncommon - 0-1 prefix, 0-1 suffix
		return AffixLimits{0, 1, 0, 1}
	case 2: // Rare - 1-2 each
		return AffixLimits{1, 2, 1, 2}
	case 3: // Epic - 1-3 each
		return AffixLimits{1, 3, 1, 3}
	case 4: // Legendary - 2-3 each
		return AffixLimits{2, 3, 2, 3}
	case 5: // Mythic - 3 each
		return AffixLimits{3, 3, 3, 3}
	default:
		return AffixLimits{0, 3, 0, 3}
	}
}

// BaseSet is the default implementation of Set interface
type BaseSet struct {
	mu        sync.RWMutex
	instances map[string]Instance // affixID -> instance
	groups    map[string]string   // group -> affixID (for mutual exclusion)
	limits    AffixLimits
}

// NewBaseSet creates new affix set with default limits
func NewBaseSet() *BaseSet {
	return &BaseSet{
		instances: make(map[string]Instance),
		groups:    make(map[string]string),
		limits:    AffixLimits{0, 3, 0, 3},
	}
}

// NewBaseSetWithLimits creates new affix set with specified limits
func NewBaseSetWithLimits(limits AffixLimits) *BaseSet {
	return &BaseSet{
		instances: make(map[string]Instance),
		groups:    make(map[string]string),
		limits:    limits,
	}
}

// NewBaseSetForRarity creates set with limits appropriate for rarity
func NewBaseSetForRarity(rarity int) *BaseSet {
	return NewBaseSetWithLimits(DefaultLimits(rarity))
}

func (bs *BaseSet) Add(instance Instance) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if !bs.canAddInternal(instance) {
		return fmt.Errorf("cannot add affix %s: limits or group conflict", instance.AffixID())
	}

	bs.instances[instance.AffixID()] = instance

	// Track group for mutual exclusion
	group := instance.Group()
	if group != "" {
		bs.groups[group] = instance.AffixID()
	}

	return nil
}

func (bs *BaseSet) Remove(affixID string) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	instance, exists := bs.instances[affixID]
	if !exists {
		return fmt.Errorf("affix not found: %s", affixID)
	}

	// Remove group tracking
	group := instance.Group()
	if group != "" {
		delete(bs.groups, group)
	}

	delete(bs.instances, affixID)
	return nil
}

func (bs *BaseSet) Get(affixID string) (Instance, bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	instance, exists := bs.instances[affixID]
	return instance, exists
}

func (bs *BaseSet) GetByType(affixType Type) []Instance {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	result := make([]Instance, 0)
	for _, instance := range bs.instances {
		if instance.Type() == affixType {
			result = append(result, instance)
		}
	}
	return result
}

func (bs *BaseSet) GetAll() []Instance {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	result := make([]Instance, 0, len(bs.instances))
	for _, instance := range bs.instances {
		result = append(result, instance)
	}
	return result
}

func (bs *BaseSet) Count() int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return len(bs.instances)
}

func (bs *BaseSet) CountByType(affixType Type) int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.countByTypeInternal(affixType)
}

func (bs *BaseSet) countByTypeInternal(affixType Type) int {
	count := 0
	for _, instance := range bs.instances {
		if instance.Type() == affixType {
			count++
		}
	}
	return count
}

func (bs *BaseSet) CanAdd(instance Instance) bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.canAddInternal(instance)
}

func (bs *BaseSet) canAddInternal(instance Instance) bool {
	// Check if already has this affix
	if _, exists := bs.instances[instance.AffixID()]; exists {
		return false
	}

	// Check group conflict
	group := instance.Group()
	if group != "" {
		if _, exists := bs.groups[group]; exists {
			return false
		}
	}

	// Check type limits
	switch instance.Type() {
	case TypePrefix:
		return bs.countByTypeInternal(TypePrefix) < bs.limits.MaxPrefixes
	case TypeSuffix:
		return bs.countByTypeInternal(TypeSuffix) < bs.limits.MaxSuffixes
	default:
		// Implicit, Corrupted, Enchant - no limits by default
		return true
	}
}

func (bs *BaseSet) HasGroup(group string) bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	_, exists := bs.groups[group]
	return exists
}

func (bs *BaseSet) PrefixCount() int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.countByTypeInternal(TypePrefix)
}

func (bs *BaseSet) SuffixCount() int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.countByTypeInternal(TypeSuffix)
}

func (bs *BaseSet) MaxPrefixes() int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.limits.MaxPrefixes
}

func (bs *BaseSet) MaxSuffixes() int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.limits.MaxSuffixes
}

func (bs *BaseSet) SetLimits(minPrefix, maxPrefix, minSuffix, maxSuffix int) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if minPrefix < 0 {
		minPrefix = 0
	}
	if maxPrefix < minPrefix {
		maxPrefix = minPrefix
	}
	if minSuffix < 0 {
		minSuffix = 0
	}
	if maxSuffix < minSuffix {
		maxSuffix = minSuffix
	}

	bs.limits = AffixLimits{
		MinPrefixes: minPrefix,
		MaxPrefixes: maxPrefix,
		MinSuffixes: minSuffix,
		MaxSuffixes: maxSuffix,
	}
}

func (bs *BaseSet) Limits() AffixLimits {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.limits
}

func (bs *BaseSet) Clear() {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.instances = make(map[string]Instance)
	bs.groups = make(map[string]string)
}

func (bs *BaseSet) AllModifiers() []attribute.Modifier {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	modifiers := make([]attribute.Modifier, 0)
	for _, instance := range bs.instances {
		modifiers = append(modifiers, instance.Modifiers()...)
	}
	return modifiers
}

func (bs *BaseSet) RerollAll() {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	for _, instance := range bs.instances {
		instance.Reroll()
	}
}

func (bs *BaseSet) TotalQuality() float64 {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	if len(bs.instances) == 0 {
		return 0.0
	}

	total := 0.0
	for _, instance := range bs.instances {
		total += instance.Quality()
	}
	return total / float64(len(bs.instances))
}

// GetGroups returns all occupied groups
func (bs *BaseSet) GetGroups() []string {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	groups := make([]string, 0, len(bs.groups))
	for group := range bs.groups {
		groups = append(groups, group)
	}
	return groups
}

// CanAddMore checks if more affixes of given type can be added
func (bs *BaseSet) CanAddMore(affixType Type) bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	switch affixType {
	case TypePrefix:
		return bs.countByTypeInternal(TypePrefix) < bs.limits.MaxPrefixes
	case TypeSuffix:
		return bs.countByTypeInternal(TypeSuffix) < bs.limits.MaxSuffixes
	default:
		return true
	}
}

// NeedMore checks if more affixes of given type are required to meet minimums
func (bs *BaseSet) NeedMore(affixType Type) bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	switch affixType {
	case TypePrefix:
		return bs.countByTypeInternal(TypePrefix) < bs.limits.MinPrefixes
	case TypeSuffix:
		return bs.countByTypeInternal(TypeSuffix) < bs.limits.MinSuffixes
	default:
		return false
	}
}

// IsComplete checks if set meets minimum requirements
func (bs *BaseSet) IsComplete() bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.countByTypeInternal(TypePrefix) >= bs.limits.MinPrefixes &&
		bs.countByTypeInternal(TypeSuffix) >= bs.limits.MinSuffixes
}

// RemainingPrefixes returns how many more prefixes can be added
func (bs *BaseSet) RemainingPrefixes() int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.limits.MaxPrefixes - bs.countByTypeInternal(TypePrefix)
}

// RemainingSuffixes returns how many more suffixes can be added
func (bs *BaseSet) RemainingSuffixes() int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.limits.MaxSuffixes - bs.countByTypeInternal(TypeSuffix)
}
