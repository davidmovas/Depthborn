package entity

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

var _ Entity = (*Base)(nil)

type Base struct{}

func (b *Base) ID() string {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Type() string {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Snapshot() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Restore(data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Version() int64 {
	//TODO implement me
	panic("implement me")
}

func (b *Base) IncrementVersion() int64 {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Delta(fromVersion int64) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) ApplyDelta(delta []byte) error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) CreatedAt() int64 {
	//TODO implement me
	panic("implement me")
}

func (b *Base) UpdatedAt() int64 {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Touch() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Name() string {
	//TODO implement me
	panic("implement me")
}

func (b *Base) SetName(name string) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Level() int {
	//TODO implement me
	panic("implement me")
}

func (b *Base) SetLevel(level int) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Tags() types.TagSet {
	//TODO implement me
	panic("implement me")
}

func (b *Base) IsAlive() bool {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Kill(ctx context.Context, killerID string) error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Revive(ctx context.Context, healthPercent float64) error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) CanAct() bool {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Clone() any {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Validate() error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Attributes() AttributeManager {
	//TODO implement me
	panic("implement me")
}

func (b *Base) StatusEffects() StatusManager {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Transform() spatial.Transform {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Callbacks() types.CallbackRegistry {
	//TODO implement me
	panic("implement me")
}
