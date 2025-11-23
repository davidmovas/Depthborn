package affix

import (
	"github.com/davidmovas/Depthborn/internal/core/attribute"
)

var _ Affix = (*BaseAffix)(nil)

type BaseAffix struct {
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
	return ba.id
}

func (ba *BaseAffix) Name() string {
	return ba.name
}

func (ba *BaseAffix) Type() Type {
	return ba.affixType
}

func (ba *BaseAffix) Tier() int {
	return ba.tier
}

func (ba *BaseAffix) Modifiers() []attribute.Modifier {
	return ba.modifiers
}

func (ba *BaseAffix) Requirements() Requirements {
	return ba.requirements
}

func (ba *BaseAffix) Weight() int {
	return ba.weight
}

func (ba *BaseAffix) Description() string {
	return ba.description
}

func (ba *BaseAffix) Tags() []string {
	return ba.tags
}

func (ba *BaseAffix) AddModifier(modifier attribute.Modifier) {
	ba.modifiers = append(ba.modifiers, modifier)
}

func (ba *BaseAffix) SetRequirements(req Requirements) {
	ba.requirements = req
}

func (ba *BaseAffix) SetWeight(weight int) {
	if weight < 0 {
		weight = 0
	}
	ba.weight = weight
}

func (ba *BaseAffix) SetDescription(description string) {
	ba.description = description
}

func (ba *BaseAffix) AddTag(tag string) {
	ba.tags = append(ba.tags, tag)
}
