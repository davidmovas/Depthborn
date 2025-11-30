package builder

import (
	"github.com/davidmovas/Depthborn/internal/item"
)

// Item provides a fluent interface for creating base items
type Item struct {
	config item.BaseItemConfig
}

// NewItem creates a new item builder
func NewItem() *Item {
	return &Item{
		config: item.BaseItemConfig{
			MaxStackSize: 1,
			Quality:      1.0,
			Level:        1,
			Weight:       0.1,
		},
	}
}

// For creates a builder with type and name
func For(itemType item.Type, name string) *Item {
	return NewItem().Type(itemType).Name(name)
}

// Material creates a material item builder
func Material(name string) *Item {
	return For(item.TypeMaterial, name)
}

// Currency creates a currency item builder
func Currency(name string) *Item {
	return For(item.TypeCurrency, name).MaxStack(9999)
}

// Quest creates a quest item builder
func Quest(name string) *Item {
	return For(item.TypeQuest, name)
}

// Key creates a key item builder
func Key(name string) *Item {
	return For(item.TypeKey, name)
}

func (b *Item) ID(id string) *Item {
	b.config.ID = id
	return b
}

func (b *Item) Name(name string) *Item {
	b.config.Name = name
	return b
}

func (b *Item) Description(desc string) *Item {
	b.config.Description = desc
	return b
}

func (b *Item) Type(t item.Type) *Item {
	b.config.ItemType = t
	return b
}

func (b *Item) Rarity(r item.Rarity) *Item {
	b.config.Rarity = r
	return b
}

func (b *Item) Quality(q float64) *Item {
	b.config.Quality = q
	return b
}

func (b *Item) Level(lvl int) *Item {
	b.config.Level = lvl
	return b
}

func (b *Item) MaxStack(max int) *Item {
	b.config.MaxStackSize = max
	return b
}

func (b *Item) Value(v int64) *Item {
	b.config.Value = v
	return b
}

func (b *Item) Weight(w float64) *Item {
	b.config.Weight = w
	return b
}

func (b *Item) Icon(icon string) *Item {
	b.config.Icon = icon
	return b
}

func (b *Item) Tags(tags ...string) *Item {
	b.config.Tags = append(b.config.Tags, tags...)
	return b
}

func (b *Item) Build() *item.BaseItem {
	return item.NewBaseItemWithConfig(b.config)
}

// Config returns the current configuration
func (b *Item) Config() item.BaseItemConfig {
	return b.config
}
