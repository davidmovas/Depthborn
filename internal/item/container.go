package item

import (
	"fmt"
)

var _ Container = (*BaseContainer)(nil)

type BaseContainer struct {
	*BaseItem
	capacity int
	contents []Item
}

func NewBaseContainer(id string, name string, capacity int) *BaseContainer {
	return &BaseContainer{
		BaseItem: NewBaseItem(id, TypeContainer, name),
		capacity: capacity,
		contents: make([]Item, 0),
	}
}

func (bc *BaseContainer) Capacity() int {
	return bc.capacity
}

func (bc *BaseContainer) Contents() []Item {
	return bc.contents
}

func (bc *BaseContainer) Add(item Item) error {
	if bc.IsFull() {
		return fmt.Errorf("container is full")
	}
	bc.contents = append(bc.contents, item)
	return nil
}

func (bc *BaseContainer) Remove(itemID string) (Item, error) {
	for i, item := range bc.contents {
		if item.ID() == itemID {
			bc.contents = append(bc.contents[:i], bc.contents[i+1:]...)
			return item, nil
		}
	}
	return nil, fmt.Errorf("item not found: %s", itemID)
}

func (bc *BaseContainer) Contains(itemID string) bool {
	for _, item := range bc.contents {
		if item.ID() == itemID {
			return true
		}
	}
	return false
}

func (bc *BaseContainer) IsFull() bool {
	return len(bc.contents) >= bc.capacity
}

func (bc *BaseContainer) Weight() float64 {
	totalWeight := bc.BaseItem.Weight()
	for _, item := range bc.contents {
		totalWeight += item.Weight()
	}
	return totalWeight
}
