package types

import "context"

// DeathCallback is invoked when entity dies
type DeathCallback func(ctx context.Context, victimID string, killerID string)

// DamageCallback is invoked when entity takes damage
type DamageCallback func(ctx context.Context, victimID string, damage float64, sourceID string)

// HealCallback is invoked when entity is healed
type HealCallback func(ctx context.Context, targetID string, amount float64, sourceID string)

// CallbackRegistry manages event callbacks
type CallbackRegistry interface {
	// OnDeath registers death callback
	OnDeath(callback DeathCallback)

	// OnDamage registers damage callback
	OnDamage(callback DamageCallback)

	// OnHeal registers heal callback
	OnHeal(callback HealCallback)

	// TriggerDeath invokes death callbacks
	TriggerDeath(ctx context.Context, victimID, killerID string)

	// TriggerDamage invokes damage callbacks
	TriggerDamage(ctx context.Context, victimID string, damage float64, sourceID string)

	// TriggerHeal invokes heal callbacks
	TriggerHeal(ctx context.Context, targetID string, amount float64, sourceID string)

	// ClearCallbacks removes all callbacks
	ClearCallbacks()
}
