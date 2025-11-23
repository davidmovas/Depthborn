package entity

import (
	"context"
	"time"
)

type BaseStatusManager struct {
	effects map[string]*StatusEffect
}

type StatusEffect struct {
	id        string
	sourceID  string
	startTime int64
	duration  int64
	active    bool
}

func NewBaseStatusManager() StatusManager {
	return &BaseStatusManager{
		effects: make(map[string]*StatusEffect),
	}
}

func (bsm *BaseStatusManager) Apply(ctx context.Context, effectID string, sourceID string) error {
	bsm.effects[effectID] = &StatusEffect{
		id:        effectID,
		sourceID:  sourceID,
		startTime: time.Now().UnixMilli(),
		duration:  0, // permanent until removed
		active:    true,
	}
	return nil
}

func (bsm *BaseStatusManager) Remove(ctx context.Context, effectID string) error {
	delete(bsm.effects, effectID)
	return nil
}

func (bsm *BaseStatusManager) Has(effectType string) bool {
	effect, exists := bsm.effects[effectType]
	return exists && effect.active
}

func (bsm *BaseStatusManager) GetAll() []string {
	result := make([]string, 0, len(bsm.effects))
	for effectID, effect := range bsm.effects {
		if effect.active {
			result = append(result, effectID)
		}
	}
	return result
}

func (bsm *BaseStatusManager) Update(ctx context.Context, deltaMs int64) error {
	now := time.Now().UnixMilli()

	for effectID, effect := range bsm.effects {
		if effect.duration > 0 && (now-effect.startTime) > effect.duration {
			effect.active = false
			delete(bsm.effects, effectID)
		}
	}

	return nil
}
