package affix

import (
	"fmt"
	"math/rand/v2"
)

var _ Generator = (*BaseGenerator)(nil)

// BaseGenerator is the default implementation of Generator interface
type BaseGenerator struct {
	pool Pool
}

// NewBaseGenerator creates generator with specified pool
func NewBaseGenerator(pool Pool) *BaseGenerator {
	return &BaseGenerator{
		pool: pool,
	}
}

func (bg *BaseGenerator) Generate(ctx GenerateContext) ([]Instance, error) {
	instances := make([]Instance, 0)

	// Determine number of prefixes and suffixes
	numPrefixes := randomInRange(ctx.PrefixRange[0], ctx.PrefixRange[1])
	numSuffixes := randomInRange(ctx.SuffixRange[0], ctx.SuffixRange[1])

	// Track used groups
	usedGroups := make(map[string]bool)
	usedIDs := make(map[string]bool)

	// Generate prefixes
	prefixType := TypePrefix
	prefixCtx := ctx.RollContext
	prefixCtx.AffixType = &prefixType

	for i := 0; i < numPrefixes; i++ {
		// Update exclusions
		prefixCtx.ExcludeGroups = mapKeys(usedGroups)
		prefixCtx.ExcludeIDs = mapKeys(usedIDs)

		affix, err := bg.pool.Roll(prefixCtx)
		if err != nil {
			// No more eligible prefixes, stop generating
			break
		}

		instance := bg.createInstanceWithBias(affix, ctx.QualityBias)
		instances = append(instances, instance)

		// Track used
		if affix.Group() != "" {
			usedGroups[affix.Group()] = true
		}
		usedIDs[affix.ID()] = true
	}

	// Generate suffixes
	suffixType := TypeSuffix
	suffixCtx := ctx.RollContext
	suffixCtx.AffixType = &suffixType

	for i := 0; i < numSuffixes; i++ {
		// Update exclusions
		suffixCtx.ExcludeGroups = mapKeys(usedGroups)
		suffixCtx.ExcludeIDs = mapKeys(usedIDs)

		affix, err := bg.pool.Roll(suffixCtx)
		if err != nil {
			// No more eligible suffixes, stop generating
			break
		}

		instance := bg.createInstanceWithBias(affix, ctx.QualityBias)
		instances = append(instances, instance)

		// Track used
		if affix.Group() != "" {
			usedGroups[affix.Group()] = true
		}
		usedIDs[affix.ID()] = true
	}

	return instances, nil
}

func (bg *BaseGenerator) AddAffix(set Set, ctx RollContext) (Instance, error) {
	// Get currently used groups
	baseSet, ok := set.(*BaseSet)
	if ok {
		ctx.ExcludeGroups = baseSet.GetGroups()
	}

	// Get used affix IDs
	for _, inst := range set.GetAll() {
		ctx.ExcludeIDs = append(ctx.ExcludeIDs, inst.AffixID())
	}

	affix, err := bg.pool.Roll(ctx)
	if err != nil {
		return nil, err
	}

	instance := bg.CreateInstance(affix)
	if err := set.Add(instance); err != nil {
		return nil, err
	}

	return instance, nil
}

func (bg *BaseGenerator) CreateInstance(affix Affix) Instance {
	values := RollModifiers(affix.Modifiers())
	return NewBaseInstance(affix, values)
}

func (bg *BaseGenerator) createInstanceWithBias(affix Affix, bias float64) Instance {
	values := RollModifiersBiased(affix.Modifiers(), bias)
	return NewBaseInstance(affix, values)
}

func (bg *BaseGenerator) RollValues(templates []ModifierTemplate) []RolledModifier {
	return RollModifiers(templates)
}

// Helper to get map keys as slice
func mapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// randomInRange returns random int in [min, max] inclusive
func randomInRange(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.IntN(max-min+1)
}

// GenerateForItem is a convenience function to generate affixes for an item
func GenerateForItem(pool Pool, itemType string, itemLevel int, slot string, rarity int) ([]Instance, error) {
	limits := DefaultLimits(rarity)
	gen := NewBaseGenerator(pool)

	ctx := GenerateContext{
		RollContext: RollContext{
			ItemType:   itemType,
			ItemLevel:  itemLevel,
			ItemSlot:   slot,
			ItemRarity: rarity,
		},
		PrefixRange: [2]int{limits.MinPrefixes, limits.MaxPrefixes},
		SuffixRange: [2]int{limits.MinSuffixes, limits.MaxSuffixes},
		QualityBias: 0.5, // Uniform distribution
	}

	return gen.Generate(ctx)
}

// GenerateForItemBiased generates with quality bias
func GenerateForItemBiased(pool Pool, itemType string, itemLevel int, slot string, rarity int, qualityBias float64) ([]Instance, error) {
	limits := DefaultLimits(rarity)
	gen := NewBaseGenerator(pool)

	ctx := GenerateContext{
		RollContext: RollContext{
			ItemType:   itemType,
			ItemLevel:  itemLevel,
			ItemSlot:   slot,
			ItemRarity: rarity,
		},
		PrefixRange: [2]int{limits.MinPrefixes, limits.MaxPrefixes},
		SuffixRange: [2]int{limits.MinSuffixes, limits.MaxSuffixes},
		QualityBias: qualityBias,
	}

	return gen.Generate(ctx)
}

// PopulateSet fills an affix set to meet minimum requirements
func PopulateSet(gen Generator, set Set, ctx RollContext) error {
	baseSet, ok := set.(*BaseSet)
	if !ok {
		return fmt.Errorf("set must be *BaseSet")
	}

	// Add prefixes until minimum is met
	prefixType := TypePrefix
	for baseSet.NeedMore(TypePrefix) {
		rollCtx := ctx
		rollCtx.AffixType = &prefixType
		rollCtx.ExcludeGroups = baseSet.GetGroups()

		_, err := gen.AddAffix(set, rollCtx)
		if err != nil {
			return fmt.Errorf("failed to add prefix: %w", err)
		}
	}

	// Add suffixes until minimum is met
	suffixType := TypeSuffix
	for baseSet.NeedMore(TypeSuffix) {
		rollCtx := ctx
		rollCtx.AffixType = &suffixType
		rollCtx.ExcludeGroups = baseSet.GetGroups()

		_, err := gen.AddAffix(set, rollCtx)
		if err != nil {
			return fmt.Errorf("failed to add suffix: %w", err)
		}
	}

	return nil
}

// RerollSet removes all affixes and generates new ones
func RerollSet(gen Generator, set Set, ctx GenerateContext) error {
	set.Clear()

	instances, err := gen.(*BaseGenerator).Generate(ctx)
	if err != nil {
		return err
	}

	for _, inst := range instances {
		if err := set.Add(inst); err != nil {
			return err
		}
	}

	return nil
}
