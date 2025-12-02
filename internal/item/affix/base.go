package affix

import (
	"sync"
)

var _ Affix = (*BaseAffix)(nil)

// BaseAffix is the default implementation of Affix interface.
// Represents immutable affix template loaded from data files.
type BaseAffix struct {
	mu           sync.RWMutex
	id           string
	name         string
	affixType    Type
	group        string
	rank         int
	modifiers    []ModifierTemplate
	requirements Requirements
	baseWeight   int
	description  string
	tags         []string
}

// AffixConfig holds configuration for creating BaseAffix
type AffixConfig struct {
	ID           string
	Name         string
	Type         Type
	Group        string
	Rank         int
	Modifiers    []ModifierTemplate
	Requirements Requirements
	BaseWeight   int
	Description  string
	Tags         []string
}

// NewBaseAffix creates new affix with default values
func NewBaseAffix(id string, name string, affixType Type) *BaseAffix {
	return &BaseAffix{
		id:           id,
		name:         name,
		affixType:    affixType,
		group:        "",
		rank:         50,
		modifiers:    make([]ModifierTemplate, 0),
		requirements: nil,
		baseWeight:   100,
		description:  "",
		tags:         make([]string, 0),
	}
}

// NewBaseAffixWithConfig creates affix from configuration
func NewBaseAffixWithConfig(cfg AffixConfig) *BaseAffix {
	ba := &BaseAffix{
		id:           cfg.ID,
		name:         cfg.Name,
		affixType:    cfg.Type,
		group:        cfg.Group,
		rank:         cfg.Rank,
		modifiers:    cfg.Modifiers,
		requirements: cfg.Requirements,
		baseWeight:   cfg.BaseWeight,
		description:  cfg.Description,
		tags:         cfg.Tags,
	}

	if ba.modifiers == nil {
		ba.modifiers = make([]ModifierTemplate, 0)
	}
	if ba.tags == nil {
		ba.tags = make([]string, 0)
	}
	if ba.rank <= 0 {
		ba.rank = 50
	}
	if ba.baseWeight <= 0 {
		ba.baseWeight = 100
	}

	return ba
}

func (ba *BaseAffix) ID() string {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.id
}

func (ba *BaseAffix) Name() string {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.name
}

func (ba *BaseAffix) Type() Type {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.affixType
}

func (ba *BaseAffix) Group() string {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.group
}

func (ba *BaseAffix) Rank() int {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.rank
}

func (ba *BaseAffix) Modifiers() []ModifierTemplate {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	result := make([]ModifierTemplate, len(ba.modifiers))
	copy(result, ba.modifiers)
	return result
}

func (ba *BaseAffix) Requirements() Requirements {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.requirements
}

func (ba *BaseAffix) BaseWeight() int {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.baseWeight
}

func (ba *BaseAffix) Description() string {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.description
}

func (ba *BaseAffix) Tags() []string {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	result := make([]string, len(ba.tags))
	copy(result, ba.tags)
	return result
}

func (ba *BaseAffix) HasTag(tag string) bool {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	for _, t := range ba.tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Builder methods for fluent API

// WithGroup sets mutual exclusion group
func (ba *BaseAffix) WithGroup(group string) *BaseAffix {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.group = group
	return ba
}

// WithRank sets internal power rank [1-100]
func (ba *BaseAffix) WithRank(rank int) *BaseAffix {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	if rank < 1 {
		rank = 1
	}
	if rank > 100 {
		rank = 100
	}
	ba.rank = rank
	return ba
}

// AddModifier adds modifier template
func (ba *BaseAffix) AddModifier(modifier ModifierTemplate) *BaseAffix {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.modifiers = append(ba.modifiers, modifier)
	return ba
}

// WithModifiers sets all modifier templates
func (ba *BaseAffix) WithModifiers(modifiers []ModifierTemplate) *BaseAffix {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.modifiers = make([]ModifierTemplate, len(modifiers))
	copy(ba.modifiers, modifiers)
	return ba
}

// WithRequirements sets item requirements
func (ba *BaseAffix) WithRequirements(req Requirements) *BaseAffix {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.requirements = req
	return ba
}

// WithBaseWeight sets spawn weight
func (ba *BaseAffix) WithBaseWeight(weight int) *BaseAffix {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	if weight < 0 {
		weight = 0
	}
	ba.baseWeight = weight
	return ba
}

// WithDescription sets description template
func (ba *BaseAffix) WithDescription(description string) *BaseAffix {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.description = description
	return ba
}

// AddTag adds single tag
func (ba *BaseAffix) AddTag(tag string) *BaseAffix {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.tags = append(ba.tags, tag)
	return ba
}

// WithTags sets all tags
func (ba *BaseAffix) WithTags(tags []string) *BaseAffix {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.tags = make([]string, len(tags))
	copy(ba.tags, tags)
	return ba
}
