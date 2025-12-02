package skill

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
)

// =============================================================================
// SKILL TREE
// =============================================================================

// Tree represents a skill progression graph (like PoE passive tree).
// Single shared tree with multiple branches for different playstyles.
type Tree interface {
	// ID returns unique tree identifier
	ID() string

	// Name returns display name
	Name() string

	// GetNode returns node by ID
	GetNode(nodeID string) (Node, bool)

	// GetNodes returns all nodes in tree
	GetNodes() []Node

	// GetBranches returns all branch definitions
	GetBranches() []Branch

	// GetStartNodes returns entry point nodes
	GetStartNodes() []Node

	// GetAdjacentNodes returns nodes connected to given node
	GetAdjacentNodes(nodeID string) []Node

	// PathExists checks if path exists between two nodes
	PathExists(fromNodeID, toNodeID string) bool
}

// Branch represents a thematic grouping of nodes (crafting, defense, trading, etc.)
type Branch struct {
	ID          string   // Unique branch identifier
	Name        string   // Display name
	Description string   // What this branch is about
	Color       string   // UI color for nodes in this branch
	NodeIDs     []string // Nodes belonging to this branch
}

// =============================================================================
// TREE STATE (Player's allocation)
// =============================================================================

// TreeState represents a player's allocations in a tree.
// Separate from Tree definition to allow multiple characters to share tree definition.
type TreeState interface {
	// TreeID returns the tree this state is for
	TreeID() string

	// AllocateNode unlocks a node (spends points)
	AllocateNode(ctx context.Context, nodeID string) error

	// DeallocateNode locks a node (refunds points, costs currency)
	DeallocateNode(ctx context.Context, nodeID string) error

	// DeallocateMultiple removes multiple nodes at once
	DeallocateMultiple(ctx context.Context, nodeIDs []string) error

	// ResetAll removes all allocations (costs currency)
	ResetAll(ctx context.Context) error

	// IsAllocated checks if node is unlocked
	IsAllocated(nodeID string) bool

	// GetAllocatedNodes returns all unlocked nodes
	GetAllocatedNodes() []string

	// GetAllocatedLevel returns level of allocated node (for leveled nodes)
	GetAllocatedLevel(nodeID string) int

	// LevelUpNode increases level of allocated node
	LevelUpNode(ctx context.Context, nodeID string) error

	// CanAllocate checks if node can be unlocked
	CanAllocate(nodeID string) bool

	// CanDeallocate checks if node can be removed
	// (node must not be required by other allocated nodes)
	CanDeallocate(nodeID string) bool

	// AvailablePoints returns unspent skill points
	AvailablePoints() int

	// SpentPoints returns total spent points
	SpentPoints() int

	// AddPoints grants skill points
	AddPoints(amount int)

	// RespecCost calculates currency cost to deallocate node(s)
	RespecCost(nodeIDs []string) int64

	// ResetCost calculates currency cost for full reset
	ResetCost() int64

	// GetActiveEffects returns all effects from allocated nodes
	GetActiveEffects() []NodeEffect

	// ApplyEffects applies all allocated node effects to entity
	ApplyEffects(ctx context.Context, entityID string) error

	// RemoveEffects removes all allocated node effects from entity
	RemoveEffects(ctx context.Context, entityID string) error
}

// =============================================================================
// NODE
// =============================================================================

// NodeType categorizes tree nodes
type NodeType string

const (
	NodePath     NodeType = "path"     // Small bonus, path connector
	NodeNotable  NodeType = "notable"  // Significant bonus
	NodeKeystone NodeType = "keystone" // Game-changing effect
	NodeSkill    NodeType = "skill"    // Grants active skill
	NodeMastery  NodeType = "mastery"  // Branch mastery bonus
)

// Node represents a single node in the skill tree
type Node interface {
	// ID returns unique node identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns node description
	Description() string

	// Type returns node type
	Type() NodeType

	// Branch returns which branch this node belongs to
	Branch() string

	// Cost returns point cost to allocate
	Cost() int

	// MaxLevel returns maximum level (0 = not leveled, just allocated)
	MaxLevel() int

	// LevelCost returns cost per level for leveled nodes
	LevelCost() int

	// Requirements returns prerequisite node IDs (must have at least one allocated)
	Requirements() []string

	// Exclusions returns mutually exclusive node IDs
	// If any of these is allocated, this node cannot be allocated
	Exclusions() []string

	// Connections returns adjacent node IDs (for pathing)
	Connections() []string

	// Effects returns effects granted by this node
	Effects() []NodeEffect

	// EffectsAtLevel returns effects for specific level
	EffectsAtLevel(level int) []NodeEffect

	// SkillID returns skill granted (for NodeSkill type)
	SkillID() string

	// Position returns visual position in tree UI
	Position() (x, y float64)

	// Icon returns icon identifier
	Icon() string
}

// =============================================================================
// NODE EFFECTS
// =============================================================================

// NodeEffectType categorizes node effects
type NodeEffectType string

const (
	EffectTypeAttribute   NodeEffectType = "attribute"    // Modify attribute
	EffectTypeSkillMod    NodeEffectType = "skill_mod"    // Modify skill
	EffectTypeGrantSkill  NodeEffectType = "grant_skill"  // Grant new skill
	EffectTypeUnlockCraft NodeEffectType = "unlock_craft" // Unlock crafting recipe
	EffectTypeTrade       NodeEffectType = "trade"        // Trading bonus
	EffectTypePassive     NodeEffectType = "passive"      // Passive trigger effect
	EffectTypeResource    NodeEffectType = "resource"     // Resource modification
	EffectTypeSpecial     NodeEffectType = "special"      // Custom effect
)

// NodeEffect describes what a node grants when allocated
type NodeEffect interface {
	// Type returns effect type
	Type() NodeEffectType

	// Apply applies effect to entity
	Apply(ctx context.Context, entityID string) error

	// Remove removes effect from entity
	Remove(ctx context.Context, entityID string) error

	// Description returns human-readable description
	Description() string

	// Value returns effect value (for simple numeric effects)
	Value() float64

	// Metadata returns additional effect data
	Metadata() map[string]any
}

// AttributeEffect modifies an attribute
type AttributeEffect struct {
	Attribute attribute.Type
	ModType   attribute.ModifierType
	Value     float64
}

// SkillModEffect modifies skills matching criteria
type SkillModEffect struct {
	TargetSkillID   string   // Specific skill ID (empty = all matching tags)
	TargetSkillTags []string // Skills with these tags are affected
	Modifier        SkillModifier
}

// GrantSkillEffect grants access to a skill
type GrantSkillEffect struct {
	SkillID    string
	StartLevel int
}

// PassiveEffect triggers on certain conditions
type PassiveEffect struct {
	TriggerType  TriggerType
	TriggerValue float64 // Chance or threshold
	EffectID     string  // What effect to apply
}

// TriggerType defines when passive triggers
type TriggerType string

const (
	TriggerOnHit        TriggerType = "on_hit"
	TriggerOnCrit       TriggerType = "on_crit"
	TriggerOnKill       TriggerType = "on_kill"
	TriggerOnDamaged    TriggerType = "on_damaged"
	TriggerOnBlock      TriggerType = "on_block"
	TriggerOnDodge      TriggerType = "on_dodge"
	TriggerOnSkillUse   TriggerType = "on_skill_use"
	TriggerOnCraft      TriggerType = "on_craft"
	TriggerOnTrade      TriggerType = "on_trade"
	TriggerOnLevelUp    TriggerType = "on_level_up"
	TriggerOnRest       TriggerType = "on_rest"
	TriggerPeriodic     TriggerType = "periodic"
	TriggerOnLowHealth  TriggerType = "on_low_health"
	TriggerOnFullHealth TriggerType = "on_full_health"
)

// =============================================================================
// LOADOUT (Equipped active skills)
// =============================================================================

// Loadout manages which active skills are equipped and ready to use
type Loadout interface {
	// Equip assigns skill instance to slot
	Equip(slot int, skill Instance) error

	// Unequip removes skill from slot
	Unequip(slot int) error

	// GetSkill returns skill in slot
	GetSkill(slot int) (Instance, bool)

	// GetAllSkills returns all equipped skills
	GetAllSkills() map[int]Instance

	// GetSkillByID finds equipped skill by definition ID
	GetSkillByID(skillDefID string) (Instance, int, bool)

	// SlotCount returns total available slots
	SlotCount() int

	// SetSlotCount changes available slots
	SetSlotCount(count int)

	// Swap exchanges skills between slots
	Swap(slot1, slot2 int) error

	// Clear removes all skills
	Clear()

	// Update processes all skill cooldowns/charges
	Update(deltaMs int64)

	// CanUseSlot checks if skill in slot can be used
	CanUseSlot(ctx context.Context, slot int, casterID string) bool

	// UseSlot activates skill in slot
	UseSlot(ctx context.Context, slot int, casterID string, params ActivationParams) (Result, error)
}
