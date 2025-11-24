package types

import (
	"context"
	"sync"
)

var _ CallbackRegistry = (*BaseCallbackRegistry)(nil)

type BaseCallbackRegistry struct {
	mu sync.RWMutex

	deathCallbacks  []DeathCallback
	damageCallbacks []DamageCallback
	healCallbacks   []HealCallback
}

func NewCallbackRegistry() CallbackRegistry {
	return &BaseCallbackRegistry{
		deathCallbacks:  make([]DeathCallback, 0),
		damageCallbacks: make([]DamageCallback, 0),
		healCallbacks:   make([]HealCallback, 0),
	}
}

func (cr *BaseCallbackRegistry) OnDeath(callback DeathCallback) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.deathCallbacks = append(cr.deathCallbacks, callback)
}

func (cr *BaseCallbackRegistry) OnDamage(callback DamageCallback) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.damageCallbacks = append(cr.damageCallbacks, callback)
}

func (cr *BaseCallbackRegistry) OnHeal(callback HealCallback) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.healCallbacks = append(cr.healCallbacks, callback)
}

func (cr *BaseCallbackRegistry) TriggerDeath(ctx context.Context, victimID, killerID string) {
	cr.mu.RLock()
	callbacks := make([]DeathCallback, len(cr.deathCallbacks))
	copy(callbacks, cr.deathCallbacks)
	cr.mu.RUnlock()

	for _, callback := range callbacks {
		callback(ctx, victimID, killerID)
	}
}

func (cr *BaseCallbackRegistry) TriggerDamage(ctx context.Context, victimID string, damage float64, sourceID string) {
	cr.mu.RLock()
	callbacks := make([]DamageCallback, len(cr.damageCallbacks))
	copy(callbacks, cr.damageCallbacks)
	cr.mu.RUnlock()

	for _, callback := range callbacks {
		callback(ctx, victimID, damage, sourceID)
	}
}

func (cr *BaseCallbackRegistry) TriggerHeal(ctx context.Context, targetID string, amount float64, sourceID string) {
	cr.mu.RLock()
	callbacks := make([]HealCallback, len(cr.healCallbacks))
	copy(callbacks, cr.healCallbacks)
	cr.mu.RUnlock()

	for _, callback := range callbacks {
		callback(ctx, targetID, amount, sourceID)
	}
}

func (cr *BaseCallbackRegistry) ClearCallbacks() {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.deathCallbacks = make([]DeathCallback, 0)
	cr.damageCallbacks = make([]DamageCallback, 0)
	cr.healCallbacks = make([]HealCallback, 0)
}
