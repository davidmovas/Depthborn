package entity

import "github.com/davidmovas/Depthborn/internal/core/types"

var _ types.Identity = (*BaseIdentity)(nil)

type BaseIdentity struct {
	id  string
	typ string
}

func NewBaseIdentity(id string, typ string) *BaseIdentity {
	return &BaseIdentity{
		id:  id,
		typ: typ,
	}
}

func (bi *BaseIdentity) ID() string {
	return bi.id
}

func (bi *BaseIdentity) Type() string {
	return bi.typ
}
