package builder

import (
	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/item"
)

// Socketable provides a fluent interface for creating socketables
type Socketable struct {
	*Item
	socketType item.SocketType
	effect     item.SocketEffect
	effectID   string
	tier       int
	modifiers  []attribute.Modifier
}

// NewSocketable creates a new socketable builder
func NewSocketable() *Socketable {
	return &Socketable{
		Item:       NewItem(),
		socketType: item.SocketTypeGem,
		tier:       1,
	}
}

// Socket creates a socketable builder
func Socket(name string, socketType item.SocketType) *Socketable {
	return NewSocketable().Name(name).SocketType(socketType)
}

// Gem creates a gem socketable
func Gem(name string) *Socketable {
	return Socket(name, item.SocketTypeGem).Type(item.TypeGem)
}

// Rune creates a rune socketable
func Rune(name string) *Socketable {
	return Socket(name, item.SocketTypeRune).Type(item.TypeRune)
}

// Chainable methods from Item
func (b *Socketable) ID(id string) *Socketable {
	b.Item.ID(id)
	return b
}

func (b *Socketable) Name(name string) *Socketable {
	b.Item.Name(name)
	return b
}

func (b *Socketable) Description(desc string) *Socketable {
	b.Item.Description(desc)
	return b
}

func (b *Socketable) Type(t item.Type) *Socketable {
	b.Item.Type(t)
	return b
}

func (b *Socketable) Rarity(r item.Rarity) *Socketable {
	b.Item.Rarity(r)
	return b
}

func (b *Socketable) Quality(q float64) *Socketable {
	b.Item.Quality(q)
	return b
}

func (b *Socketable) Level(lvl int) *Socketable {
	b.Item.Level(lvl)
	return b
}

func (b *Socketable) Value(v int64) *Socketable {
	b.Item.Value(v)
	return b
}

func (b *Socketable) Weight(w float64) *Socketable {
	b.Item.Weight(w)
	return b
}

func (b *Socketable) Icon(icon string) *Socketable {
	b.Item.Icon(icon)
	return b
}

func (b *Socketable) Tags(tags ...string) *Socketable {
	b.Item.Tags(tags...)
	return b
}

// Socketable-specific methods
func (b *Socketable) SocketType(st item.SocketType) *Socketable {
	b.socketType = st
	return b
}

func (b *Socketable) Effect(effect item.SocketEffect, effectID string) *Socketable {
	b.effect = effect
	b.effectID = effectID
	return b
}

func (b *Socketable) Tier(tier int) *Socketable {
	b.tier = tier
	return b
}

func (b *Socketable) Modifier(mod attribute.Modifier) *Socketable {
	b.modifiers = append(b.modifiers, mod)
	return b
}

func (b *Socketable) Modifiers(mods ...attribute.Modifier) *Socketable {
	b.modifiers = append(b.modifiers, mods...)
	return b
}

func (b *Socketable) Build() *item.BaseSocketable {
	cfg := item.SocketableConfig{
		BaseItemConfig: b.Item.Config(),
		SocketType:     b.socketType,
		Effect:         b.effect,
		EffectID:       b.effectID,
		Tier:           b.tier,
		Modifiers:      b.modifiers,
	}

	return item.NewBaseSocketableWithConfig(cfg)
}
