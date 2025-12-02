package skill

import (
	"context"
	"errors"
	"sync"
)

// =============================================================================
// ERRORS
// =============================================================================

var (
	ErrNodeNotFound         = errors.New("node not found")
	ErrNodeAlreadyAlloc     = errors.New("node already allocated")
	ErrNodeNotAllocated     = errors.New("node not allocated")
	ErrInsufficientPoints   = errors.New("insufficient skill points")
	ErrRequirementsNotMet   = errors.New("requirements not met")
	ErrNodeExcluded         = errors.New("node excluded by another allocation")
	ErrNodeRequired         = errors.New("node is required by other allocations")
	ErrInsufficientCurrency = errors.New("insufficient currency for respec")
)

// =============================================================================
// BASE TREE (Definition)
// =============================================================================

var _ Tree = (*BaseTree)(nil)

// BaseTree implements Tree interface - the static tree definition
type BaseTree struct {
	mu sync.RWMutex

	id         string
	name       string
	nodes      map[string]*BaseNode
	branches   []Branch
	startNodes []string
}

// TreeConfig holds configuration for creating BaseTree
type TreeConfig struct {
	ID       string
	Name     string
	Branches []Branch
}

// NewBaseTree creates a new tree definition
func NewBaseTree(config TreeConfig) *BaseTree {
	return &BaseTree{
		id:         config.ID,
		name:       config.Name,
		nodes:      make(map[string]*BaseNode),
		branches:   config.Branches,
		startNodes: make([]string, 0),
	}
}

func (t *BaseTree) ID() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.id
}

func (t *BaseTree) Name() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.name
}

func (t *BaseTree) GetNode(nodeID string) (Node, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node, ok := t.nodes[nodeID]
	return node, ok
}

func (t *BaseTree) GetNodes() []Node {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]Node, 0, len(t.nodes))
	for _, n := range t.nodes {
		result = append(result, n)
	}
	return result
}

func (t *BaseTree) GetBranches() []Branch {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]Branch, len(t.branches))
	copy(result, t.branches)
	return result
}

func (t *BaseTree) GetStartNodes() []Node {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]Node, 0, len(t.startNodes))
	for _, id := range t.startNodes {
		if n, ok := t.nodes[id]; ok {
			result = append(result, n)
		}
	}
	return result
}

func (t *BaseTree) GetAdjacentNodes(nodeID string) []Node {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node, ok := t.nodes[nodeID]
	if !ok {
		return nil
	}

	result := make([]Node, 0, len(node.connections))
	for _, connID := range node.connections {
		if n, ok := t.nodes[connID]; ok {
			result = append(result, n)
		}
	}
	return result
}

func (t *BaseTree) PathExists(fromNodeID, toNodeID string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if _, ok := t.nodes[fromNodeID]; !ok {
		return false
	}
	if _, ok := t.nodes[toNodeID]; !ok {
		return false
	}

	// BFS to find path
	visited := make(map[string]bool)
	queue := []string{fromNodeID}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == toNodeID {
			return true
		}

		if visited[current] {
			continue
		}
		visited[current] = true

		if node, ok := t.nodes[current]; ok {
			for _, conn := range node.connections {
				if !visited[conn] {
					queue = append(queue, conn)
				}
			}
		}
	}

	return false
}

// AddNode adds a node to the tree
func (t *BaseTree) AddNode(node *BaseNode) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.nodes[node.id] = node
}

// SetStartNodes sets entry point nodes
func (t *BaseTree) SetStartNodes(nodeIDs []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.startNodes = make([]string, len(nodeIDs))
	copy(t.startNodes, nodeIDs)
}

// =============================================================================
// BASE NODE
// =============================================================================

var _ Node = (*BaseNode)(nil)

// BaseNode implements Node interface
type BaseNode struct {
	mu sync.RWMutex

	id           string
	name         string
	description  string
	nodeType     NodeType
	branch       string
	cost         int
	maxLevel     int
	levelCost    int
	requirements []string
	exclusions   []string
	connections  []string
	effects      []NodeEffect
	levelEffects map[int][]NodeEffect
	skillID      string
	posX, posY   float64
	icon         string
}

// NodeConfig holds configuration for creating BaseNode
type NodeConfig struct {
	ID           string
	Name         string
	Description  string
	Type         NodeType
	Branch       string
	Cost         int
	MaxLevel     int
	LevelCost    int
	Requirements []string
	Exclusions   []string
	Connections  []string
	Effects      []NodeEffect
	SkillID      string
	PosX, PosY   float64
	Icon         string
}

// NewBaseNode creates a new tree node
func NewBaseNode(config NodeConfig) *BaseNode {
	return &BaseNode{
		id:           config.ID,
		name:         config.Name,
		description:  config.Description,
		nodeType:     config.Type,
		branch:       config.Branch,
		cost:         config.Cost,
		maxLevel:     config.MaxLevel,
		levelCost:    config.LevelCost,
		requirements: config.Requirements,
		exclusions:   config.Exclusions,
		connections:  config.Connections,
		effects:      config.Effects,
		levelEffects: make(map[int][]NodeEffect),
		skillID:      config.SkillID,
		posX:         config.PosX,
		posY:         config.PosY,
		icon:         config.Icon,
	}
}

func (n *BaseNode) ID() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.id
}

func (n *BaseNode) Name() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.name
}

func (n *BaseNode) Description() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.description
}

func (n *BaseNode) Type() NodeType {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.nodeType
}

func (n *BaseNode) Branch() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.branch
}

func (n *BaseNode) Cost() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.cost
}

func (n *BaseNode) MaxLevel() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.maxLevel
}

func (n *BaseNode) LevelCost() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.levelCost
}

func (n *BaseNode) Requirements() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]string, len(n.requirements))
	copy(result, n.requirements)
	return result
}

func (n *BaseNode) Exclusions() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]string, len(n.exclusions))
	copy(result, n.exclusions)
	return result
}

func (n *BaseNode) Connections() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]string, len(n.connections))
	copy(result, n.connections)
	return result
}

func (n *BaseNode) Effects() []NodeEffect {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]NodeEffect, len(n.effects))
	copy(result, n.effects)
	return result
}

func (n *BaseNode) EffectsAtLevel(level int) []NodeEffect {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if effects, ok := n.levelEffects[level]; ok {
		result := make([]NodeEffect, len(effects))
		copy(result, effects)
		return result
	}

	// Return base effects if no level-specific
	result := make([]NodeEffect, len(n.effects))
	copy(result, n.effects)
	return result
}

func (n *BaseNode) SkillID() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.skillID
}

func (n *BaseNode) Position() (x, y float64) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.posX, n.posY
}

func (n *BaseNode) Icon() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.icon
}

// SetLevelEffects sets effects for specific level
func (n *BaseNode) SetLevelEffects(level int, effects []NodeEffect) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.levelEffects[level] = effects
}

// =============================================================================
// BASE TREE STATE (Player's allocations)
// =============================================================================

var _ TreeState = (*BaseTreeState)(nil)

// BaseTreeState implements TreeState interface
type BaseTreeState struct {
	mu sync.RWMutex

	treeID          string
	tree            Tree           // Reference to tree definition
	allocated       map[string]int // nodeID -> level (1 = allocated, >1 = leveled)
	availablePoints int
	spentPoints     int

	// Respec cost configuration
	baseCostPerNode  int64
	costPerNodeLevel int64
	resetCostBase    int64
}

// TreeStateConfig holds configuration for tree state
type TreeStateConfig struct {
	TreeID           string
	Tree             Tree
	BaseCostPerNode  int64
	CostPerNodeLevel int64
	ResetCostBase    int64
}

// NewBaseTreeState creates a new tree state
func NewBaseTreeState(config TreeStateConfig) *BaseTreeState {
	return &BaseTreeState{
		treeID:           config.TreeID,
		tree:             config.Tree,
		allocated:        make(map[string]int),
		availablePoints:  0,
		spentPoints:      0,
		baseCostPerNode:  config.BaseCostPerNode,
		costPerNodeLevel: config.CostPerNodeLevel,
		resetCostBase:    config.ResetCostBase,
	}
}

func (s *BaseTreeState) TreeID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.treeID
}

func (s *BaseTreeState) AllocateNode(ctx context.Context, nodeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_ = ctx

	// Check if node exists
	node, ok := s.tree.GetNode(nodeID)
	if !ok {
		return ErrNodeNotFound
	}

	// Check if already allocated
	if s.allocated[nodeID] > 0 {
		return ErrNodeAlreadyAlloc
	}

	// Check points
	cost := node.Cost()
	if s.availablePoints < cost {
		return ErrInsufficientPoints
	}

	// Check requirements (need at least one requirement allocated)
	reqs := node.Requirements()
	if len(reqs) > 0 {
		hasReq := false
		for _, reqID := range reqs {
			if s.allocated[reqID] > 0 {
				hasReq = true
				break
			}
		}
		if !hasReq {
			return ErrRequirementsNotMet
		}
	}

	// Check exclusions
	for _, exclID := range node.Exclusions() {
		if s.allocated[exclID] > 0 {
			return ErrNodeExcluded
		}
	}

	// Allocate
	s.allocated[nodeID] = 1
	s.availablePoints -= cost
	s.spentPoints += cost

	return nil
}

func (s *BaseTreeState) DeallocateNode(ctx context.Context, nodeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_ = ctx

	// Check if allocated
	level, ok := s.allocated[nodeID]
	if !ok || level == 0 {
		return ErrNodeNotAllocated
	}

	// Check if any other node requires this one
	for allocID := range s.allocated {
		if allocID == nodeID {
			continue
		}
		if node, ok := s.tree.GetNode(allocID); ok {
			for _, reqID := range node.Requirements() {
				if reqID == nodeID {
					// Check if there's another path
					if !s.hasAlternativeRequirement(allocID, nodeID) {
						return ErrNodeRequired
					}
				}
			}
		}
	}

	// Get node for cost refund
	node, ok := s.tree.GetNode(nodeID)
	if !ok {
		return ErrNodeNotFound
	}

	// Calculate refund (base cost + level costs)
	refund := node.Cost()
	if level > 1 {
		refund += (level - 1) * node.LevelCost()
	}

	// Deallocate
	delete(s.allocated, nodeID)
	s.availablePoints += refund
	s.spentPoints -= refund

	return nil
}

func (s *BaseTreeState) hasAlternativeRequirement(nodeID, excludeReqID string) bool {
	node, ok := s.tree.GetNode(nodeID)
	if !ok {
		return false
	}

	for _, reqID := range node.Requirements() {
		if reqID != excludeReqID && s.allocated[reqID] > 0 {
			return true
		}
	}
	return false
}

func (s *BaseTreeState) DeallocateMultiple(ctx context.Context, nodeIDs []string) error {
	// Check if all can be deallocated first
	for _, nodeID := range nodeIDs {
		if err := s.canDeallocateWithExclusions(nodeID, nodeIDs); err != nil {
			return err
		}
	}

	// Deallocate all
	for _, nodeID := range nodeIDs {
		if err := s.DeallocateNode(ctx, nodeID); err != nil {
			return err
		}
	}

	return nil
}

func (s *BaseTreeState) canDeallocateWithExclusions(nodeID string, excluding []string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	level, ok := s.allocated[nodeID]
	if !ok || level == 0 {
		return ErrNodeNotAllocated
	}

	// Check if any other node (not in excluding list) requires this
	for allocID := range s.allocated {
		if allocID == nodeID {
			continue
		}

		// Skip if in excluding list
		isExcluded := false
		for _, exID := range excluding {
			if allocID == exID {
				isExcluded = true
				break
			}
		}
		if isExcluded {
			continue
		}

		if node, ok := s.tree.GetNode(allocID); ok {
			for _, reqID := range node.Requirements() {
				if reqID == nodeID {
					return ErrNodeRequired
				}
			}
		}
	}

	return nil
}

func (s *BaseTreeState) ResetAll(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_ = ctx

	// Calculate total refund
	totalRefund := 0
	for nodeID, level := range s.allocated {
		if node, ok := s.tree.GetNode(nodeID); ok {
			refund := node.Cost()
			if level > 1 {
				refund += (level - 1) * node.LevelCost()
			}
			totalRefund += refund
		}
	}

	// Clear allocations
	s.allocated = make(map[string]int)
	s.availablePoints += totalRefund
	s.spentPoints = 0

	return nil
}

func (s *BaseTreeState) IsAllocated(nodeID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.allocated[nodeID] > 0
}

func (s *BaseTreeState) GetAllocatedNodes() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, 0, len(s.allocated))
	for nodeID := range s.allocated {
		result = append(result, nodeID)
	}
	return result
}

func (s *BaseTreeState) GetAllocatedLevel(nodeID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.allocated[nodeID]
}

func (s *BaseTreeState) LevelUpNode(ctx context.Context, nodeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_ = ctx

	// Check if allocated
	level, ok := s.allocated[nodeID]
	if !ok || level == 0 {
		return ErrNodeNotAllocated
	}

	// Check if can level up
	node, ok := s.tree.GetNode(nodeID)
	if !ok {
		return ErrNodeNotFound
	}

	if node.MaxLevel() == 0 || level >= node.MaxLevel() {
		return ErrMaxLevel
	}

	// Check points for level up
	cost := node.LevelCost()
	if s.availablePoints < cost {
		return ErrInsufficientPoints
	}

	// Level up
	s.allocated[nodeID] = level + 1
	s.availablePoints -= cost
	s.spentPoints += cost

	return nil
}

func (s *BaseTreeState) CanAllocate(nodeID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	node, ok := s.tree.GetNode(nodeID)
	if !ok {
		return false
	}

	// Already allocated?
	if s.allocated[nodeID] > 0 {
		return false
	}

	// Enough points?
	if s.availablePoints < node.Cost() {
		return false
	}

	// Requirements met?
	reqs := node.Requirements()
	if len(reqs) > 0 {
		hasReq := false
		for _, reqID := range reqs {
			if s.allocated[reqID] > 0 {
				hasReq = true
				break
			}
		}
		if !hasReq {
			return false
		}
	}

	// Exclusions check
	for _, exclID := range node.Exclusions() {
		if s.allocated[exclID] > 0 {
			return false
		}
	}

	return true
}

func (s *BaseTreeState) CanDeallocate(nodeID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.allocated[nodeID] == 0 {
		return false
	}

	// Check if required by others
	for allocID := range s.allocated {
		if allocID == nodeID {
			continue
		}
		if node, ok := s.tree.GetNode(allocID); ok {
			for _, reqID := range node.Requirements() {
				if reqID == nodeID && !s.hasAlternativeRequirementLocked(allocID, nodeID) {
					return false
				}
			}
		}
	}

	return true
}

func (s *BaseTreeState) hasAlternativeRequirementLocked(nodeID, excludeReqID string) bool {
	node, ok := s.tree.GetNode(nodeID)
	if !ok {
		return false
	}

	for _, reqID := range node.Requirements() {
		if reqID != excludeReqID && s.allocated[reqID] > 0 {
			return true
		}
	}
	return false
}

func (s *BaseTreeState) AvailablePoints() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.availablePoints
}

func (s *BaseTreeState) SpentPoints() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.spentPoints
}

func (s *BaseTreeState) AddPoints(amount int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.availablePoints += amount
}

func (s *BaseTreeState) RespecCost(nodeIDs []string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int64
	for _, nodeID := range nodeIDs {
		level := s.allocated[nodeID]
		if level > 0 {
			total += s.baseCostPerNode
			if level > 1 {
				total += int64(level-1) * s.costPerNodeLevel
			}
		}
	}
	return total
}

func (s *BaseTreeState) ResetCost() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cost := s.resetCostBase
	for nodeID, level := range s.allocated {
		_ = nodeID
		cost += s.baseCostPerNode
		if level > 1 {
			cost += int64(level-1) * s.costPerNodeLevel
		}
	}
	return cost
}

func (s *BaseTreeState) GetActiveEffects() []NodeEffect {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var effects []NodeEffect
	for nodeID, level := range s.allocated {
		if node, ok := s.tree.GetNode(nodeID); ok {
			nodeEffects := node.EffectsAtLevel(level)
			effects = append(effects, nodeEffects...)
		}
	}
	return effects
}

func (s *BaseTreeState) ApplyEffects(ctx context.Context, entityID string) error {
	effects := s.GetActiveEffects()
	for _, effect := range effects {
		if err := effect.Apply(ctx, entityID); err != nil {
			return err
		}
	}
	return nil
}

func (s *BaseTreeState) RemoveEffects(ctx context.Context, entityID string) error {
	effects := s.GetActiveEffects()
	for _, effect := range effects {
		if err := effect.Remove(ctx, entityID); err != nil {
			return err
		}
	}
	return nil
}

// =============================================================================
// SERIALIZATION
// =============================================================================

// TreeStateData holds serializable tree state
type TreeStateData struct {
	TreeID          string         `msgpack:"tree_id"`
	Allocated       map[string]int `msgpack:"allocated"`
	AvailablePoints int            `msgpack:"available_points"`
	SpentPoints     int            `msgpack:"spent_points"`
}

// GetData returns serializable data
func (s *BaseTreeState) GetData() TreeStateData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	allocated := make(map[string]int, len(s.allocated))
	for k, v := range s.allocated {
		allocated[k] = v
	}

	return TreeStateData{
		TreeID:          s.treeID,
		Allocated:       allocated,
		AvailablePoints: s.availablePoints,
		SpentPoints:     s.spentPoints,
	}
}

// RestoreData restores from serialized data
func (s *BaseTreeState) RestoreData(data TreeStateData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.treeID = data.TreeID
	s.allocated = make(map[string]int, len(data.Allocated))
	for k, v := range data.Allocated {
		s.allocated[k] = v
	}
	s.availablePoints = data.AvailablePoints
	s.spentPoints = data.SpentPoints
}
