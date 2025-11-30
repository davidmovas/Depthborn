package affix

import (
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
)

var _ Affix = (*BaseAffix)(nil)

type BaseAffix struct {
	mu           sync.RWMutex
	id           string
	name         string
	affixType    Type
	tier         int
	modifiers    []attribute.Modifier
	requirements Requirements
	weight       int
	description  string
	tags         []string
}

func NewBaseAffix(id string, name string, affixType Type, tier int) *BaseAffix {
	return &BaseAffix{
		id:           id,
		name:         name,
		affixType:    affixType,
		tier:         tier,
		modifiers:    make([]attribute.Modifier, 0),
		requirements: nil,
		weight:       100,
		description:  "",
		tags:         make([]string, 0),
	}
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

func (ba *BaseAffix) Tier() int {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.tier
}

func (ba *BaseAffix) Modifiers() []attribute.Modifier {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	result := make([]attribute.Modifier, len(ba.modifiers))
	copy(result, ba.modifiers)
	return result
}

func (ba *BaseAffix) Requirements() Requirements {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.requirements
}

func (ba *BaseAffix) Weight() int {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	return ba.weight
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

func (ba *BaseAffix) AddModifier(modifier attribute.Modifier) {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.modifiers = append(ba.modifiers, modifier)
}

func (ba *BaseAffix) SetRequirements(req Requirements) {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.requirements = req
}

func (ba *BaseAffix) SetWeight(weight int) {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	if weight < 0 {
		weight = 0
	}
	ba.weight = weight
}

func (ba *BaseAffix) SetDescription(description string) {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.description = description
}

func (ba *BaseAffix) AddTag(tag string) {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.tags = append(ba.tags, tag)
}
