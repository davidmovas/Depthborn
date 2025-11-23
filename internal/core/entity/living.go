package entity

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
)

var _ Living = (*LivingEntity)(nil)

type LivingEntity struct {
	*BaseEntity
	health    float64
	maxHealth float64
}

func NewLivingEntity(id string, name string) *LivingEntity {
	base := NewBaseEntity(id, name)
	entity := &LivingEntity{
		BaseEntity: base,
		health:     100.0,
		maxHealth:  100.0,
	}

	entity.attributes.SetBase(attribute.AttrVitality, 10.0)
	entity.recalculateHealth()

	return entity
}

func (le *LivingEntity) Health() float64 {
	return le.health
}

func (le *LivingEntity) MaxHealth() float64 {
	return le.maxHealth
}

func (le *LivingEntity) SetHealth(value float64) {
	if value < 0 {
		value = 0
	} else if value > le.maxHealth {
		value = le.maxHealth
	}
	le.health = value
}

func (le *LivingEntity) Damage(ctx context.Context, amount float64, sourceID string) (float64, error) {
	if amount <= 0 {
		return 0, nil
	}

	actualDamage := amount
	if le.health < amount {
		actualDamage = le.health
	}

	le.health -= actualDamage

	// Trigger damage callbacks
	if registry, ok := le.callbacks.(*BaseCallbackRegistry); ok {
		for _, callback := range registry.damageCallbacks {
			callback(ctx, le.ID(), actualDamage, sourceID)
		}
	}

	// Check for death
	if le.health <= 0 {
		if err := le.Kill(ctx, sourceID); err != nil {
			return 0, err
		}
	}

	return actualDamage, nil
}

func (le *LivingEntity) Heal(ctx context.Context, amount float64, sourceID string) (float64, error) {
	if amount <= 0 {
		return 0, nil
	}

	missingHealth := le.maxHealth - le.health
	actualHeal := amount
	if amount > missingHealth {
		actualHeal = missingHealth
	}

	le.health += actualHeal

	// Trigger heal callbacks
	if registry, ok := le.callbacks.(*BaseCallbackRegistry); ok {
		for _, callback := range registry.healCallbacks {
			callback(ctx, le.ID(), actualHeal, sourceID)
		}
	}

	return actualHeal, nil
}

func (le *LivingEntity) HealthPercent() float64 {
	if le.maxHealth <= 0 {
		return 0.0
	}
	return le.health / le.maxHealth
}

func (le *LivingEntity) recalculateHealth() {
	vitality := le.attributes.Get(attribute.AttrVitality)
	le.maxHealth = 50.0 + (vitality * 10.0)

	if le.health > le.maxHealth {
		le.health = le.maxHealth
	}
}
