package entity

import "github.com/davidmovas/Depthborn/internal/core/types"

var _ types.Validatable = (*BaseValidatable)(nil)

type BaseValidatable struct{}

func NewBaseValidatable() *BaseValidatable {
	return &BaseValidatable{}
}

func (bv *BaseValidatable) Validate() error {
	// TODO: Add validation logic
	return nil
}
