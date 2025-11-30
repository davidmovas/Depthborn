package builder

import (
	"github.com/davidmovas/Depthborn/internal/item"
)

// Container provides a fluent interface for creating containers
type Container struct {
	*Item
	capacity     int
	maxWeight    float64
	allowedTypes []item.Type
}

// NewContainer creates a new container builder
func NewContainer() *Container {
	return &Container{
		Item:     NewItem().Type(item.TypeContainer),
		capacity: 10,
	}
}

// Contain creates a container builder with name and capacity
func Contain(name string, capacity int) *Container {
	return NewContainer().Name(name).Capacity(capacity)
}

// Bag creates a bag container
func Bag(name string, capacity int) *Container {
	return Contain(name, capacity).Tags("bag")
}

// Chest creates a chest container
func Chest(name string, capacity int) *Container {
	return Contain(name, capacity).Tags("chest")
}

// Chainable methods from Item
func (b *Container) ID(id string) *Container {
	b.Item.ID(id)
	return b
}

func (b *Container) Name(name string) *Container {
	b.Item.Name(name)
	return b
}

func (b *Container) Description(desc string) *Container {
	b.Item.Description(desc)
	return b
}

func (b *Container) Rarity(r item.Rarity) *Container {
	b.Item.Rarity(r)
	return b
}

func (b *Container) Quality(q float64) *Container {
	b.Item.Quality(q)
	return b
}

func (b *Container) Level(lvl int) *Container {
	b.Item.Level(lvl)
	return b
}

func (b *Container) Value(v int64) *Container {
	b.Item.Value(v)
	return b
}

func (b *Container) Weight(w float64) *Container {
	b.Item.Weight(w)
	return b
}

func (b *Container) Icon(icon string) *Container {
	b.Item.Icon(icon)
	return b
}

func (b *Container) Tags(tags ...string) *Container {
	b.Item.Tags(tags...)
	return b
}

// Container-specific methods
func (b *Container) Capacity(cap int) *Container {
	b.capacity = cap
	return b
}

func (b *Container) MaxWeight(weight float64) *Container {
	b.maxWeight = weight
	return b
}

func (b *Container) AllowTypes(types ...item.Type) *Container {
	b.allowedTypes = append(b.allowedTypes, types...)
	return b
}

func (b *Container) Build() *item.BaseContainer {
	cfg := item.ContainerConfig{
		BaseItemConfig: b.Item.Config(),
		Capacity:       b.capacity,
		MaxWeight:      b.maxWeight,
		AllowedTypes:   b.allowedTypes,
	}

	return item.NewBaseContainerWithConfig(cfg)
}
