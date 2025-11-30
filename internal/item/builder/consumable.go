package builder

import (
	"github.com/davidmovas/Depthborn/internal/item"
)

// Consumable provides a fluent interface for creating consumables
type Consumable struct {
	*Item
	maxCooldown int64
	effect      item.ConsumableEffect
	effectID    string
	charges     int
}

// NewConsumable creates a new consumable builder
func NewConsumable() *Consumable {
	b := &Consumable{
		Item:    NewItem().Type(item.TypeConsumable),
		charges: 1,
	}
	return b
}

// Consume creates a consumable builder with name
func Consume(name string) *Consumable {
	return NewConsumable().Name(name)
}

// Potion creates a potion consumable
func Potion(name string) *Consumable {
	return Consume(name).Tags("potion")
}

// Food creates a food consumable
func Food(name string) *Consumable {
	return Consume(name).Tags("food")
}

// Scroll creates a scroll consumable
func Scroll(name string) *Consumable {
	return Consume(name).Tags("scroll")
}

// Chainable methods from Item
func (b *Consumable) ID(id string) *Consumable {
	b.Item.ID(id)
	return b
}

func (b *Consumable) Name(name string) *Consumable {
	b.Item.Name(name)
	return b
}

func (b *Consumable) Description(desc string) *Consumable {
	b.Item.Description(desc)
	return b
}

func (b *Consumable) Rarity(r item.Rarity) *Consumable {
	b.Item.Rarity(r)
	return b
}

func (b *Consumable) Quality(q float64) *Consumable {
	b.Item.Quality(q)
	return b
}

func (b *Consumable) Level(lvl int) *Consumable {
	b.Item.Level(lvl)
	return b
}

func (b *Consumable) MaxStack(max int) *Consumable {
	b.Item.MaxStack(max)
	return b
}

func (b *Consumable) Value(v int64) *Consumable {
	b.Item.Value(v)
	return b
}

func (b *Consumable) Weight(w float64) *Consumable {
	b.Item.Weight(w)
	return b
}

func (b *Consumable) Icon(icon string) *Consumable {
	b.Item.Icon(icon)
	return b
}

func (b *Consumable) Tags(tags ...string) *Consumable {
	b.Item.Tags(tags...)
	return b
}

// Consumable-specific methods
func (b *Consumable) Cooldown(ms int64) *Consumable {
	b.maxCooldown = ms
	return b
}

func (b *Consumable) Effect(effect item.ConsumableEffect, effectID string) *Consumable {
	b.effect = effect
	b.effectID = effectID
	return b
}

func (b *Consumable) Charges(count int) *Consumable {
	b.charges = count
	return b
}

func (b *Consumable) Infinite() *Consumable {
	b.charges = -1
	return b
}

func (b *Consumable) Build() *item.BaseConsumable {
	cfg := item.ConsumableConfig{
		BaseItemConfig: b.Item.Config(),
		MaxCooldown:    b.maxCooldown,
		Effect:         b.effect,
		EffectID:       b.effectID,
		Charges:        b.charges,
	}

	return item.NewBaseConsumableWithConfig(cfg)
}
