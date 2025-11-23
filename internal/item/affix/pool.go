package affix

import (
	"fmt"
	"math/rand/v2"
)

type BasePool struct {
	affixes map[string]Affix
}

func NewBasePool() Pool {
	return &BasePool{
		affixes: make(map[string]Affix),
	}
}

func (bp *BasePool) Add(affix Affix) {
	bp.affixes[affix.ID()] = affix
}

func (bp *BasePool) Remove(affixID string) {
	delete(bp.affixes, affixID)
}

func (bp *BasePool) Get(affixID string) (Affix, bool) {
	affix, exists := bp.affixes[affixID]
	return affix, exists
}

func (bp *BasePool) GetByTier(tier int) []Affix {
	result := make([]Affix, 0)
	for _, affix := range bp.affixes {
		if affix.Tier() == tier {
			result = append(result, affix)
		}
	}
	return result
}

func (bp *BasePool) GetByTags(tags ...string) []Affix {
	result := make([]Affix, 0)
	for _, affix := range bp.affixes {
		matchesAll := true
		for _, tag := range tags {
			found := false
			for _, affixTag := range affix.Tags() {
				if affixTag == tag {
					found = true
					break
				}
			}
			if !found {
				matchesAll = false
				break
			}
		}
		if matchesAll {
			result = append(result, affix)
		}
	}
	return result
}

func (bp *BasePool) Roll(itemType string, itemLevel int, slot string) (Affix, error) {
	eligible := make([]Affix, 0)
	totalWeight := 0

	for _, affix := range bp.affixes {
		if affix.Requirements().Check(itemType, itemLevel, slot) {
			eligible = append(eligible, affix)
			totalWeight += affix.Weight()
		}
	}

	if len(eligible) == 0 {
		return nil, fmt.Errorf("no eligible affixes found")
	}

	roll := rand.IntN(totalWeight)
	currentWeight := 0

	for _, affix := range eligible {
		currentWeight += affix.Weight()
		if roll < currentWeight {
			return affix, nil
		}
	}

	return eligible[len(eligible)-1], nil
}

func (bp *BasePool) RollMultiple(count int, itemType string, itemLevel int, slot string) ([]Affix, error) {
	result := make([]Affix, 0, count)
	for i := 0; i < count; i++ {
		affix, err := bp.Roll(itemType, itemLevel, slot)
		if err != nil {
			return nil, err
		}
		result = append(result, affix)
	}
	return result, nil
}

func (bp *BasePool) Filter(criteria FilterCriteria) []Affix {
	result := make([]Affix, 0)
	for _, affix := range bp.affixes {
		matches := true

		if len(criteria.Types) > 0 {
			typeMatch := false
			for _, t := range criteria.Types {
				if affix.Type() == t {
					typeMatch = true
					break
				}
			}
			if !typeMatch {
				matches = false
			}
		}

		if criteria.MinTier > 0 && affix.Tier() < criteria.MinTier {
			matches = false
		}

		if criteria.MaxTier > 0 && affix.Tier() > criteria.MaxTier {
			matches = false
		}

		if len(criteria.Tags) > 0 {
			tagMatch := false
			for _, tag := range criteria.Tags {
				for _, affixTag := range affix.Tags() {
					if affixTag == tag {
						tagMatch = true
						break
					}
				}
				if tagMatch {
					break
				}
			}
			if !tagMatch {
				matches = false
			}
		}

		if matches {
			result = append(result, affix)
		}
	}
	return result
}
