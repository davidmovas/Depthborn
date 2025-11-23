package entity

import "github.com/davidmovas/Depthborn/internal/core/types"

var _ types.Actionable = (*BaseActionable)(nil)

type BaseActionable struct {
	canAct bool
}

func NewBaseActionable(canAct bool) *BaseActionable {
	return &BaseActionable{
		canAct: canAct,
	}
}

func (ba *BaseActionable) CanAct() bool {
	return ba.canAct
}
