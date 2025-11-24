package impl

import (
	"github.com/davidmovas/Depthborn/internal/infra"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var _ infra.Identity = (*BaseIdentity)(nil)

type BaseIdentity struct {
	id         string
	entityType string
}

func NewIdentity(entityType string) *BaseIdentity {
	id, _ := gonanoid.New()
	return &BaseIdentity{
		id:         id,
		entityType: entityType,
	}
}

func NewIdentityWithID(entityType, id string) *BaseIdentity {
	return &BaseIdentity{
		id:         id,
		entityType: entityType,
	}
}

func (i *BaseIdentity) ID() string {
	return i.id
}

func (i *BaseIdentity) Type() string {
	return i.entityType
}
