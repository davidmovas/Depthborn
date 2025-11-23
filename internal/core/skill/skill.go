package skill

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/types"
)

// Skill represents ability
type Skill interface {
	types.Identity
	types.Named

	// Description returns detailed description
	Description() string

	// SkillType returns skill type
	SkillType() Type

	// Level returns current level
	Level() int

	// MaxLevel returns maximum level
	MaxLevel() int

	// CanLevelUp checks if can level up
	CanLevelUp() bool

	// LevelUp increases level
	LevelUp() error

	// Requirements returns usage requirements
	Requirements() Requirements

	// CanUse checks if caster can use skill
	CanUse(ctx context.Context, casterID string) bool

	// Use executes skill
	Use(ctx context.Context, casterID string, params ActivationParams) (Result, error)

	// Cooldown returns remaining cooldown
	Cooldown() int64

	// MaxCooldown returns base cooldown
	MaxCooldown() int64

	// SetCooldown updates cooldown
	SetCooldown(ms int64)

	// ManaCost returns mana cost
	ManaCost() float64

	// Tags returns skill tags
	Tags() []string

	// Metadata returns skill data
	Metadata() map[string]any
}

// Type categorizes skills
type Type string

const (
	TypeActive    Type = "active"
	TypePassive   Type = "passive"
	TypeToggle    Type = "toggle"
	TypeChanneled Type = "channeled"
	TypeAura      Type = "aura"
	TypeTrigger   Type = "trigger"
)

// Requirements defines skill requirements
type Requirements interface {
	// Level returns minimum level
	Level() int

	// Attributes returns required attributes
	Attributes() map[string]float64

	// Skills returns prerequisite skills
	Skills() map[string]int

	// Items returns required items
	Items() []string

	// Check verifies entity meets requirements (uses entity ID)
	Check(entityID string) bool
}

// ActivationParams contains activation data
type ActivationParams struct {
	TargetID  string
	TargetPos *Position
	Direction float64
	Power     float64
	Modifiers map[string]any
}

// Position represents 2D coordinates
type Position struct {
	X float64
	Y float64
}

// Branch represents skill tree branch
type Branch struct {
	ID       string      `yaml:"id"`
	Name     string      `yaml:"name"`
	Nodes    []*NodeImpl `yaml:"nodes"`
	Branches []*Branch   `yaml:"branches,omitempty"`
}

// Result contains skill execution results
type Result struct {
	Success       bool
	TargetIDs     []string
	Damage        map[string]float64
	Healing       map[string]float64
	StatusApplied []string
	Message       string
	Metadata      map[string]any
}

// Tree represents skill progression tree
type Tree interface {
	// GetNode returns node by ID
	GetNode(nodeID string) (Node, bool)

	// GetNodes returns all nodes
	GetNodes() []Node

	// AllocateNode unlocks node
	AllocateNode(ctx context.Context, nodeID string) error

	// DeallocateNode locks node
	DeallocateNode(ctx context.Context, nodeID string) error

	// IsNodeAllocated checks if unlocked
	IsNodeAllocated(nodeID string) bool

	// GetAllocatedNodes returns unlocked nodes
	GetAllocatedNodes() []Node

	// CanAllocate checks if can unlock
	CanAllocate(nodeID string) bool

	// AvailablePoints returns unspent points
	AvailablePoints() int

	// SpentPoints returns allocated points
	SpentPoints() int

	// AddPoints grants skill points
	AddPoints(amount int)

	// Reset removes all allocations
	Reset() error
}

// Node represents skill tree node
type Node interface {
	types.Identity
	types.Named

	// Description returns description
	Description() string

	// SkillType returns node type
	SkillType() NodeType

	// Cost returns point cost
	Cost() int

	// Requirements returns prerequisite node IDs
	Requirements() []string

	// Connections returns adjacent node IDs
	Connections() []string

	// Effect returns node effect
	Effect() NodeEffect

	// Position returns visual position
	Position() (x, y float64)
}

// NodeType categorizes nodes
type NodeType string

const (
	NodeAttribute NodeType = "attribute"
	NodeSkill     NodeType = "skill"
	NodePassive   NodeType = "passive"
	NodeKeystone  NodeType = "keystone"
	NodeNotable   NodeType = "notable"
	NodePath      NodeType = "path"
)

// NodeEffect describes node benefit
type NodeEffect interface {
	// Apply applies to entity (uses entity ID)
	Apply(ctx context.Context, entityID string) error

	// Remove removes from entity
	Remove(ctx context.Context, entityID string) error

	// Description returns readable description
	Description() string
}

// Loadout manages equipped skills
type Loadout interface {
	// Equip assigns skill to slot
	Equip(slot int, skill Skill) error

	// Unequip removes skill from slot
	Unequip(slot int) error

	// GetSkill returns skill in slot
	GetSkill(slot int) (Skill, bool)

	// GetAllSkills returns all equipped skills
	GetAllSkills() map[int]Skill

	// SlotCount returns total slots
	SlotCount() int

	// Swap exchanges skills between slots
	Swap(slot1, slot2 int) error

	// Clear removes all skills
	Clear()
}
