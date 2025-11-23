package entity

import (
	"github.com/davidmovas/Depthborn/internal/core/types"
)

var _ types.CallbackRegistry = (*BaseCallbackRegistry)(nil)

type BaseCallbackRegistry struct {
	deathCallbacks  []types.DeathCallback
	damageCallbacks []types.DamageCallback
	healCallbacks   []types.HealCallback
}

func NewBaseCallbackRegistry() *BaseCallbackRegistry {
	return &BaseCallbackRegistry{
		deathCallbacks:  make([]types.DeathCallback, 0),
		damageCallbacks: make([]types.DamageCallback, 0),
		healCallbacks:   make([]types.HealCallback, 0),
	}
}

func (bcr *BaseCallbackRegistry) OnDeath(callback types.DeathCallback) {
	bcr.deathCallbacks = append(bcr.deathCallbacks, callback)
}

func (bcr *BaseCallbackRegistry) OnDamage(callback types.DamageCallback) {
	bcr.damageCallbacks = append(bcr.damageCallbacks, callback)
}

func (bcr *BaseCallbackRegistry) OnHeal(callback types.HealCallback) {
	bcr.healCallbacks = append(bcr.healCallbacks, callback)
}

func (bcr *BaseCallbackRegistry) ClearCallbacks() {
	bcr.deathCallbacks = make([]types.DeathCallback, 0)
	bcr.damageCallbacks = make([]types.DamageCallback, 0)
	bcr.healCallbacks = make([]types.HealCallback, 0)
}
