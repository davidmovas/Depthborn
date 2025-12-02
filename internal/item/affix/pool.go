package affix

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

var _ Pool = (*BasePool)(nil)

// BasePool is the default implementation of Pool interface
type BasePool struct {
	mu      sync.RWMutex
	affixes map[string]Affix
}

// NewBasePool creates new empty affix pool
func NewBasePool() *BasePool {
	return &BasePool{
		affixes: make(map[string]Affix),
	}
}

func (bp *BasePool) Add(affix Affix) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.affixes[affix.ID()] = affix
}

func (bp *BasePool) Remove(affixID string) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	delete(bp.affixes, affixID)
}

func (bp *BasePool) Get(affixID string) (Affix, bool) {
	bp.mu.RLock()
	defer bp.mu.RUnlock()
	affix, exists := bp.affixes[affixID]
	return affix, exists
}

func (bp *BasePool) GetAll() []Affix {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	result := make([]Affix, 0, len(bp.affixes))
	for _, affix := range bp.affixes {
		result = append(result, affix)
	}
	return result
}

func (bp *BasePool) GetByGroup(group string) []Affix {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	result := make([]Affix, 0)
	for _, affix := range bp.affixes {
		if affix.Group() == group {
			result = append(result, affix)
		}
	}
	return result
}

func (bp *BasePool) GetByTags(tags ...string) []Affix {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	result := make([]Affix, 0)
	for _, affix := range bp.affixes {
		if hasAllTags(affix.Tags(), tags) {
			result = append(result, affix)
		}
	}
	return result
}

func (bp *BasePool) GetByAnyTag(tags ...string) []Affix {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	result := make([]Affix, 0)
	for _, affix := range bp.affixes {
		if hasAnyTag(affix.Tags(), tags) {
			result = append(result, affix)
		}
	}
	return result
}

func (bp *BasePool) Filter(criteria FilterCriteria) []Affix {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	result := make([]Affix, 0)
	for _, affix := range bp.affixes {
		if bp.matchesCriteria(affix, criteria) {
			result = append(result, affix)
		}
	}
	return result
}

func (bp *BasePool) matchesCriteria(affix Affix, criteria FilterCriteria) bool {
	// Check types
	if len(criteria.Types) > 0 {
		found := false
		for _, t := range criteria.Types {
			if affix.Type() == t {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check groups
	if len(criteria.Groups) > 0 {
		found := false
		for _, g := range criteria.Groups {
			if affix.Group() == g {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check tags
	if len(criteria.Tags) > 0 {
		if !hasAnyTag(affix.Tags(), criteria.Tags) {
			return false
		}
	}

	// Check rank range
	if criteria.MinRank > 0 && affix.Rank() < criteria.MinRank {
		return false
	}
	if criteria.MaxRank > 0 && affix.Rank() > criteria.MaxRank {
		return false
	}

	// Check item level requirements
	req := affix.Requirements()
	if req != nil {
		if criteria.MinItemLevel > 0 && req.MinItemLevel() > criteria.MinItemLevel {
			return false
		}
		if criteria.MaxItemLevel > 0 && req.MaxItemLevel() > 0 && req.MaxItemLevel() < criteria.MaxItemLevel {
			return false
		}
	}

	return true
}

func (bp *BasePool) Roll(ctx RollContext) (Affix, error) {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	eligible := bp.getEligible(ctx)
	if len(eligible) == 0 {
		return nil, fmt.Errorf("no eligible affixes found")
	}

	// Calculate weights with rarity adjustment
	weights := make([]int, len(eligible))
	totalWeight := 0

	for i, affix := range eligible {
		weight := calculateEffectiveWeight(affix, ctx.ItemRarity, ctx.ItemLevel)
		weights[i] = weight
		totalWeight += weight
	}

	if totalWeight <= 0 {
		return nil, fmt.Errorf("total weight is zero or negative")
	}

	// Weighted random selection
	roll := rand.IntN(totalWeight)
	currentWeight := 0

	for i, affix := range eligible {
		currentWeight += weights[i]
		if roll < currentWeight {
			return affix, nil
		}
	}

	return eligible[len(eligible)-1], nil
}

func (bp *BasePool) getEligible(ctx RollContext) []Affix {
	eligible := make([]Affix, 0)

	for _, affix := range bp.affixes {
		if bp.isEligible(affix, ctx) {
			eligible = append(eligible, affix)
		}
	}

	return eligible
}

func (bp *BasePool) isEligible(affix Affix, ctx RollContext) bool {
	// Check type filter
	if ctx.AffixType != nil && affix.Type() != *ctx.AffixType {
		return false
	}

	// Check excluded groups
	for _, g := range ctx.ExcludeGroups {
		if affix.Group() == g {
			return false
		}
	}

	// Check excluded IDs
	for _, id := range ctx.ExcludeIDs {
		if affix.ID() == id {
			return false
		}
	}

	// Check required tags
	if len(ctx.RequireTags) > 0 {
		if !hasAllTags(affix.Tags(), ctx.RequireTags) {
			return false
		}
	}

	// Check excluded tags
	if len(ctx.ExcludeTags) > 0 {
		if hasAnyTag(affix.Tags(), ctx.ExcludeTags) {
			return false
		}
	}

	// Check requirements
	req := affix.Requirements()
	if req != nil && !req.Check(ctx.ItemType, ctx.ItemLevel, ctx.ItemSlot) {
		return false
	}

	return true
}

// calculateEffectiveWeight calculates affix weight adjusted by rarity and rank
// Higher rarity items have higher chance of high-rank affixes
// Uses exponential scaling for rank influence
func calculateEffectiveWeight(affix Affix, itemRarity int, itemLevel int) int {
	baseWeight := affix.BaseWeight()
	rank := affix.Rank()

	// Rarity factor: higher rarity = more likely to get high-rank affixes
	// Rarity 0-5 maps to factor 0.5-2.0
	rarityFactor := 0.5 + float64(itemRarity)*0.3

	// Rank influence: exponential curve
	// Low rank affixes have flat weight
	// High rank affixes get bonus from rarity
	// rank 1-100, normalized to 0.0-1.0
	rankNorm := float64(rank-1) / 99.0

	// Exponential weight adjustment
	// For common items (rarity 0): high rank affixes are rare
	// For mythic items (rarity 5): high rank affixes are common
	// Formula: weight * (rarityFactor ^ (rankNorm * 2))
	adjustment := pow(rarityFactor, rankNorm*2.0)

	// Item level also influences - higher ilvl unlocks better affixes
	// but doesn't increase weight dramatically
	ilvlBonus := 1.0
	if itemLevel > 50 {
		ilvlBonus = 1.0 + float64(itemLevel-50)/200.0 // Up to 25% bonus at ilvl 100
	}

	effectiveWeight := float64(baseWeight) * adjustment * ilvlBonus

	// Ensure minimum weight of 1
	if effectiveWeight < 1 {
		effectiveWeight = 1
	}

	return int(effectiveWeight)
}

// Helper functions

func hasAllTags(affixTags []string, required []string) bool {
	for _, req := range required {
		found := false
		for _, t := range affixTags {
			if t == req {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func hasAnyTag(affixTags []string, targets []string) bool {
	for _, target := range targets {
		for _, t := range affixTags {
			if t == target {
				return true
			}
		}
	}
	return false
}
