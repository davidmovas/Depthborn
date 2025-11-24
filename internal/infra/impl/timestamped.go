package impl

import (
	"sync/atomic"
	"time"

	"github.com/davidmovas/Depthborn/internal/infra"
)

var _ infra.Timestamped = (*BaseTimestamped)(nil)

type BaseTimestamped struct {
	createdAt int64
	updatedAt int64
}

func NewTimestamped() *BaseTimestamped {
	now := time.Now().Unix()
	return &BaseTimestamped{
		createdAt: now,
		updatedAt: now,
	}
}

func (t *BaseTimestamped) CreatedAt() int64 {
	return atomic.LoadInt64(&t.createdAt)
}

func (t *BaseTimestamped) UpdatedAt() int64 {
	return atomic.LoadInt64(&t.updatedAt)
}

func (t *BaseTimestamped) Touch() {
	atomic.StoreInt64(&t.updatedAt, time.Now().Unix())
}

func (t *BaseTimestamped) SetCreatedAt(timestamp int64) {
	atomic.StoreInt64(&t.createdAt, timestamp)
}
