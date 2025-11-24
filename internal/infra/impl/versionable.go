package impl

import (
	"sync/atomic"

	"github.com/davidmovas/Depthborn/internal/infra"
)

var _ infra.Versionable = (*BaseVersionable)(nil)

type BaseVersionable struct {
	version int64
}

func NewVersionable() *BaseVersionable {
	return &BaseVersionable{
		version: 1,
	}
}

func NewVersionableWithVersion(version int64) *BaseVersionable {
	return &BaseVersionable{
		version: version,
	}
}

func (v *BaseVersionable) Version() int64 {
	return atomic.LoadInt64(&v.version)
}

func (v *BaseVersionable) IncrementVersion() int64 {
	return atomic.AddInt64(&v.version, 1)
}

func (v *BaseVersionable) SetVersion(version int64) {
	atomic.StoreInt64(&v.version, version)
}
