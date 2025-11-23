package skill

import (
	"context"
	"errors"
	"sync"
)

var _ Tree = (*TreeImpl)(nil)

type TreeImpl struct {
	nodes          map[string]Node
	allocatedNodes map[string]bool
	points         int
	mu             sync.RWMutex
}

func NewTree(points int) *TreeImpl {
	return &TreeImpl{
		nodes:          make(map[string]Node),
		allocatedNodes: make(map[string]bool),
		points:         points,
	}
}

func (t *TreeImpl) AddNode(node Node) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.nodes[node.ID()] = node
}

func (t *TreeImpl) GetNode(nodeID string) (Node, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	node, ok := t.nodes[nodeID]
	return node, ok
}

func (t *TreeImpl) GetNodes() []Node {
	t.mu.RLock()
	defer t.mu.RUnlock()
	nodes := make([]Node, 0, len(t.nodes))
	for _, n := range t.nodes {
		nodes = append(nodes, n)
	}
	return nodes
}

func (t *TreeImpl) AllocateNode(ctx context.Context, nodeID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	node, exists := t.nodes[nodeID]
	if !exists {
		return errors.New("node not found")
	}
	if t.allocatedNodes[nodeID] {
		return errors.New("node already allocated")
	}
	if t.points < node.Cost() {
		return errors.New("not enough points")
	}

	for _, req := range node.Requirements() {
		if !t.allocatedNodes[req] {
			return errors.New("requirements not met")
		}
	}

	t.points -= node.Cost()
	t.allocatedNodes[nodeID] = true

	if node.Effect() != nil {
		//TODO: Provide entityID
		_ = node.Effect().Apply(ctx, "")
	}
	return nil
}

func (t *TreeImpl) DeallocateNode(ctx context.Context, nodeID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.allocatedNodes[nodeID] {
		return errors.New("node not allocated")
	}

	node, exists := t.nodes[nodeID]
	if !exists {
		return errors.New("node not found")
	}

	if node.Effect() != nil {
		//TODO: Provide entityID
		_ = node.Effect().Remove(ctx, "")
	}

	delete(t.allocatedNodes, nodeID)
	t.points += node.Cost()
	return nil
}

func (t *TreeImpl) IsNodeAllocated(nodeID string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.allocatedNodes[nodeID]
}

func (t *TreeImpl) GetAllocatedNodes() []Node {
	t.mu.RLock()
	defer t.mu.RUnlock()
	nodes := make([]Node, 0, len(t.allocatedNodes))
	for id := range t.allocatedNodes {
		if n, ok := t.nodes[id]; ok {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (t *TreeImpl) CanAllocate(nodeID string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	node, exists := t.nodes[nodeID]
	if !exists || t.allocatedNodes[nodeID] {
		return false
	}
	if t.points < node.Cost() {
		return false
	}
	for _, req := range node.Requirements() {
		if !t.allocatedNodes[req] {
			return false
		}
	}
	return true
}

func (t *TreeImpl) AvailablePoints() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.points
}

func (t *TreeImpl) SpentPoints() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	spent := 0
	for id := range t.allocatedNodes {
		if n, ok := t.nodes[id]; ok {
			spent += n.Cost()
		}
	}
	return spent
}

func (t *TreeImpl) AddPoints(amount int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.points += amount
}

func (t *TreeImpl) Reset() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	for id := range t.allocatedNodes {
		node := t.nodes[id]
		if node.Effect() != nil {
			_ = node.Effect().Remove(context.Background(), "")
		}
	}
	t.allocatedNodes = make(map[string]bool)
	return nil
}
