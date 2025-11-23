package entity

import "github.com/davidmovas/Depthborn/internal/core/types"

var _ types.Named = (*BaseNamed)(nil)

type BaseNamed struct {
	name string
}

func NewBaseNamed(name string) *BaseNamed {
	return &BaseNamed{
		name: name,
	}
}

func (bn *BaseNamed) Name() string {
	return bn.name
}

func (bn *BaseNamed) SetName(name string) {
	bn.name = name
}
