package entity

import (
	"context"
	"fmt"
	"math"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

var _ Living = (*BaseLiving)(nil)

type BaseLiving struct {
	*BaseEntity

	health    float64
	maxHealth float64
}

type LivingConfig struct {
	EntityConfig  Config
	InitialHealth float64
	MaxHealth     float64
}

func NewLiving(config LivingConfig) *BaseLiving {
	living := &BaseLiving{
		BaseEntity: NewEntity(config.EntityConfig),
		maxHealth:  config.MaxHealth,
	}

	// Set initial health (default to max if not specified)
	if config.InitialHealth > 0 {
		living.health = math.Min(config.InitialHealth, config.MaxHealth)
	} else {
		living.health = config.MaxHealth
	}

	return living
}

func (l *BaseLiving) Health() float64 {
	return l.health
}

func (l *BaseLiving) Revive(_ context.Context, healthPercent float64) error {
	if healthPercent <= 0 || healthPercent > 1 {
		return fmt.Errorf("healthPercent must be >0 and <=1")
	}

	if l.IsAlive() {
		return nil
	}

	maxHP := l.MaxHealth()
	newHP := maxHP * healthPercent
	if newHP < 1 {
		newHP = 1
	}

	l.health = newHP
	l.Touch()

	return nil
}

func (l *BaseLiving) Kill(ctx context.Context, killerID string) error {
	if !l.IsAlive() {
		return nil
	}

	l.health = 0

	// Call base entity Kill to set isAlive = false and trigger callbacks
	return l.BaseEntity.Kill(ctx, killerID)
}

func (l *BaseLiving) MaxHealth() float64 {
	baseMax := l.maxHealth
	if baseMax <= 0 {
		baseMax = 1
	}

	vitality := l.Attributes().Get(attribute.AttrVitality)
	vitalityBonus := vitality * 10.0

	total := baseMax + vitalityBonus
	if total < 1 {
		total = 1
	}
	return total
}
func (l *BaseLiving) SetHealth(value float64) {
	maxHP := l.MaxHealth()
	l.health = math.Max(0, math.Min(value, maxHP))

	if l.health <= 0 && l.IsAlive() {
		_ = l.Kill(context.Background(), "")
	}

	l.Touch()
}

func (l *BaseLiving) Damage(ctx context.Context, amount float64, sourceID string) (float64, error) {
	if !l.IsAlive() {
		return 0, fmt.Errorf("entity is already dead")
	}

	if amount < 0 {
		return 0, fmt.Errorf("damage amount must be positive")
	}

	// Calculate actual damage (considering defenses)
	// TODO: Apply armor, resistances, etc.
	actualDamage := amount

	// Apply damage
	oldHealth := l.health
	l.SetHealth(l.health - actualDamage)
	finalDamage := oldHealth - l.health

	// Trigger damage callback
	l.Callbacks().TriggerDamage(ctx, l.ID(), finalDamage, sourceID)

	// Trigger death callback if died
	if !l.IsAlive() {
		l.Callbacks().TriggerDeath(ctx, l.ID(), sourceID)
	}

	return finalDamage, nil
}

func (l *BaseLiving) Heal(ctx context.Context, amount float64, sourceID string) (float64, error) {
	if !l.IsAlive() {
		return 0, fmt.Errorf("cannot heal dead entity")
	}

	if amount < 0 {
		return 0, fmt.Errorf("heal amount must be positive")
	}

	// Calculate actual healing
	oldHealth := l.health
	maxHP := l.MaxHealth()
	l.SetHealth(math.Min(l.health+amount, maxHP))
	actualHealing := l.health - oldHealth

	// Trigger heal callback
	l.Callbacks().TriggerHeal(ctx, l.ID(), actualHealing, sourceID)

	return actualHealing, nil
}

func (l *BaseLiving) HealthPercent() float64 {
	maxHP := l.MaxHealth()
	if maxHP <= 0 {
		return 0
	}
	return l.health / maxHP
}

func (l *BaseLiving) SerializeState() (map[string]any, error) {
	state, err := l.BaseEntity.SerializeState()
	if err != nil {
		return nil, err
	}

	state["health"] = l.health
	state["max_health"] = l.maxHealth

	return state, nil
}

func (l *BaseLiving) DeserializeState(state map[string]any) error {
	if err := l.BaseEntity.DeserializeState(state); err != nil {
		return err
	}

	if health, ok := state["health"].(float64); ok {
		l.health = health
	}

	if maxHealth, ok := state["max_health"].(float64); ok {
		l.maxHealth = maxHealth
	}

	return nil
}

func (l *BaseLiving) Validate() error {
	if err := l.BaseEntity.Validate(); err != nil {
		return err
	}

	if l.maxHealth <= 0 {
		return fmt.Errorf("max health must be positive")
	}

	if l.health < 0 {
		return fmt.Errorf("health cannot be negative")
	}

	if l.health > l.MaxHealth() {
		return fmt.Errorf("health cannot exceed max health")
	}

	return nil
}

func (l *BaseLiving) Clone() any {
	baseClone := l.BaseEntity.Clone().(*BaseEntity)

	clone := &BaseLiving{
		BaseEntity: baseClone,
		health:     l.health,
		maxHealth:  l.maxHealth,
	}

	return clone
}

// LivingState holds the complete serializable state of a BaseLiving.
type LivingState struct {
	EntityState
	Health    float64 `msgpack:"health"`
	MaxHealth float64 `msgpack:"max_health"`
}

// MarshalBinary implements persist.Marshaler for BaseLiving.
func (l *BaseLiving) MarshalBinary() ([]byte, error) {
	// First get base entity state
	baseData, err := l.BaseEntity.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// Decode to EntityState to embed
	var es EntityState
	if err := persist.DefaultCodec().Decode(baseData, &es); err != nil {
		return nil, err
	}

	ls := LivingState{
		EntityState: es,
		Health:      l.health,
		MaxHealth:   l.maxHealth,
	}

	return persist.DefaultCodec().Encode(ls)
}

// UnmarshalBinary implements persist.Unmarshaler for BaseLiving.
func (l *BaseLiving) UnmarshalBinary(data []byte) error {
	var ls LivingState
	if err := persist.DefaultCodec().Decode(data, &ls); err != nil {
		return fmt.Errorf("failed to decode living state: %w", err)
	}

	// Encode entity state back to bytes for base entity
	entityData, err := persist.DefaultCodec().Encode(ls.EntityState)
	if err != nil {
		return err
	}

	// Initialize base entity if nil
	if l.BaseEntity == nil {
		l.BaseEntity = &BaseEntity{}
	}

	// Restore base entity
	if err := l.BaseEntity.UnmarshalBinary(entityData); err != nil {
		return err
	}

	// Restore living-specific fields
	l.health = ls.Health
	l.maxHealth = ls.MaxHealth

	return nil
}
