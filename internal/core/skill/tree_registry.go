package skill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"gopkg.in/yaml.v3"
)

// =============================================================================
// TREE REGISTRY
// =============================================================================

// TreeRegistry manages skill tree definitions loaded from YAML
type TreeRegistry interface {
	// Register adds a tree definition
	Register(tree *BaseTree) error

	// Get retrieves tree by ID
	Get(id string) (*BaseTree, bool)

	// GetAll returns all registered trees
	GetAll() []*BaseTree

	// Has checks if tree is registered
	Has(id string) bool

	// Count returns number of registered trees
	Count() int

	// LoadFromYAML loads tree from YAML data
	LoadFromYAML(data []byte) error

	// LoadFromFile loads tree from YAML file
	LoadFromFile(path string) error

	// LoadFromDirectory loads all YAML files from directory
	LoadFromDirectory(dir string) error

	// CreateState creates a new tree state for player
	CreateState(treeID string) (*BaseTreeState, error)
}

// =============================================================================
// BASE TREE REGISTRY
// =============================================================================

var _ TreeRegistry = (*BaseTreeRegistry)(nil)

// BaseTreeRegistry implements TreeRegistry interface
type BaseTreeRegistry struct {
	mu    sync.RWMutex
	trees map[string]*BaseTree
}

// NewBaseTreeRegistry creates a new tree registry
func NewBaseTreeRegistry() *BaseTreeRegistry {
	return &BaseTreeRegistry{
		trees: make(map[string]*BaseTree),
	}
}

func (r *BaseTreeRegistry) Register(tree *BaseTree) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.trees[tree.ID()]; exists {
		return fmt.Errorf("tree %s already registered", tree.ID())
	}

	r.trees[tree.ID()] = tree
	return nil
}

func (r *BaseTreeRegistry) Get(id string) (*BaseTree, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tree, ok := r.trees[id]
	return tree, ok
}

func (r *BaseTreeRegistry) GetAll() []*BaseTree {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*BaseTree, 0, len(r.trees))
	for _, tree := range r.trees {
		result = append(result, tree)
	}
	return result
}

func (r *BaseTreeRegistry) Has(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.trees[id]
	return ok
}

func (r *BaseTreeRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.trees)
}

func (r *BaseTreeRegistry) CreateState(treeID string) (*BaseTreeState, error) {
	r.mu.RLock()
	tree, ok := r.trees[treeID]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("tree %s not found", treeID)
	}

	return NewBaseTreeState(TreeStateConfig{
		TreeID:           treeID,
		Tree:             tree,
		BaseCostPerNode:  100, // Default respec costs
		CostPerNodeLevel: 50,
		ResetCostBase:    500,
	}), nil
}

// =============================================================================
// YAML STRUCTURES
// =============================================================================

// TreeFile represents the root YAML file structure
type TreeFile struct {
	Version string   `yaml:"version"`
	Tree    TreeYAML `yaml:"tree"`
}

// TreeYAML represents skill tree in YAML
type TreeYAML struct {
	ID          string       `yaml:"id"`
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Branches    []BranchYAML `yaml:"branches"`
	StartNodes  []string     `yaml:"start_nodes"`
	Nodes       []NodeYAML   `yaml:"nodes"`

	// Respec configuration
	RespecConfig *RespecConfigYAML `yaml:"respec_config"`
}

// BranchYAML represents a branch (thematic grouping) in YAML
type BranchYAML struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Color       string `yaml:"color"` // Hex color for UI
	Icon        string `yaml:"icon"`
}

// RespecConfigYAML holds respec cost configuration
type RespecConfigYAML struct {
	BaseCostPerNode  int64 `yaml:"base_cost_per_node"`
	CostPerNodeLevel int64 `yaml:"cost_per_node_level"`
	ResetCostBase    int64 `yaml:"reset_cost_base"`
}

// NodeYAML represents a tree node in YAML
type NodeYAML struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"` // path, notable, keystone, skill, mastery
	Branch      string `yaml:"branch"`
	Icon        string `yaml:"icon"`

	// Costs
	Cost      int `yaml:"cost"`       // Points to allocate
	MaxLevel  int `yaml:"max_level"`  // 0 = not leveled
	LevelCost int `yaml:"level_cost"` // Points per additional level

	// Position for UI editor
	Position PositionYAML `yaml:"position"`

	// Graph connections
	Connections  []string `yaml:"connections"`  // Adjacent nodes (bidirectional pathing)
	Requirements []string `yaml:"requirements"` // Must have at least ONE allocated
	Exclusions   []string `yaml:"exclusions"`   // Cannot allocate if ANY is allocated

	// Effects granted when allocated
	Effects []NodeEffectYAML `yaml:"effects"`

	// Level-specific effects (optional)
	Levels []NodeLevelYAML `yaml:"levels"`

	// For skill-granting nodes
	SkillID    string `yaml:"skill_id"`    // References skill definition
	SkillLevel int    `yaml:"skill_level"` // Starting level of granted skill

	// Metadata for custom logic
	Metadata map[string]any `yaml:"metadata"`
}

// PositionYAML represents node position in tree UI
type PositionYAML struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
}

// NodeEffectYAML represents an effect granted by node
type NodeEffectYAML struct {
	Type string `yaml:"type"` // attribute, skill_mod, grant_skill, unlock_craft, trade, passive, resource, special

	// For attribute effects
	Attribute string  `yaml:"attribute"`
	ModType   string  `yaml:"mod_type"` // flat, increased, more, override
	Value     float64 `yaml:"value"`

	// For skill modification effects
	TargetSkillID   string   `yaml:"target_skill_id"`
	TargetSkillTags []string `yaml:"target_skill_tags"`

	// For grant_skill effects
	SkillID    string `yaml:"skill_id"`
	StartLevel int    `yaml:"start_level"`

	// For passive trigger effects
	TriggerType  string  `yaml:"trigger_type"`
	TriggerValue float64 `yaml:"trigger_value"`
	EffectID     string  `yaml:"effect_id"`

	// Description override
	Description string `yaml:"description"`

	// Additional data
	Metadata map[string]any `yaml:"metadata"`
}

// NodeLevelYAML represents level-specific data for leveled nodes
type NodeLevelYAML struct {
	Level       int              `yaml:"level"`
	Effects     []NodeEffectYAML `yaml:"effects"`
	Description string           `yaml:"description"`
	Metadata    map[string]any   `yaml:"metadata"`
}

// =============================================================================
// YAML LOADING
// =============================================================================

func (r *BaseTreeRegistry) LoadFromYAML(data []byte) error {
	var file TreeFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	tree, err := parseTreeYAML(file.Tree)
	if err != nil {
		return fmt.Errorf("failed to parse tree %s: %w", file.Tree.ID, err)
	}

	return r.Register(tree)
}

func (r *BaseTreeRegistry) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return r.LoadFromYAML(data)
}

func (r *BaseTreeRegistry) LoadFromDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		if err := r.LoadFromFile(path); err != nil {
			return fmt.Errorf("failed to load %s: %w", path, err)
		}
	}

	return nil
}

// =============================================================================
// YAML PARSING
// =============================================================================

func parseTreeYAML(y TreeYAML) (*BaseTree, error) {
	// Create branches
	branches := make([]Branch, len(y.Branches))
	for i, b := range y.Branches {
		branches[i] = Branch{
			ID:          b.ID,
			Name:        b.Name,
			Description: b.Description,
			Color:       b.Color,
		}
	}

	// Create tree
	tree := NewBaseTree(TreeConfig{
		ID:       y.ID,
		Name:     y.Name,
		Branches: branches,
	})

	// Parse and add nodes
	for _, nodeYAML := range y.Nodes {
		node, err := parseNodeYAML(nodeYAML)
		if err != nil {
			return nil, fmt.Errorf("failed to parse node %s: %w", nodeYAML.ID, err)
		}
		tree.AddNode(node)
	}

	// Set start nodes
	tree.SetStartNodes(y.StartNodes)

	// Update branches with node IDs
	for i := range tree.branches {
		branch := &tree.branches[i]
		for _, node := range tree.nodes {
			if node.Branch() == branch.ID {
				branch.NodeIDs = append(branch.NodeIDs, node.ID())
			}
		}
	}

	return tree, nil
}

func parseNodeYAML(y NodeYAML) (*BaseNode, error) {
	// Parse effects
	effects, err := parseNodeEffects(y.Effects)
	if err != nil {
		return nil, err
	}

	node := NewBaseNode(NodeConfig{
		ID:           y.ID,
		Name:         y.Name,
		Description:  y.Description,
		Type:         parseNodeType(y.Type),
		Branch:       y.Branch,
		Cost:         y.Cost,
		MaxLevel:     y.MaxLevel,
		LevelCost:    y.LevelCost,
		Requirements: y.Requirements,
		Exclusions:   y.Exclusions,
		Connections:  y.Connections,
		Effects:      effects,
		SkillID:      y.SkillID,
		PosX:         y.Position.X,
		PosY:         y.Position.Y,
		Icon:         y.Icon,
	})

	// Parse level-specific effects
	for _, levelYAML := range y.Levels {
		levelEffects, err := parseNodeEffects(levelYAML.Effects)
		if err != nil {
			return nil, fmt.Errorf("failed to parse level %d effects: %w", levelYAML.Level, err)
		}
		node.SetLevelEffects(levelYAML.Level, levelEffects)
	}

	return node, nil
}

func parseNodeType(s string) NodeType {
	switch s {
	case "path":
		return NodePath
	case "notable":
		return NodeNotable
	case "keystone":
		return NodeKeystone
	case "skill":
		return NodeSkill
	case "mastery":
		return NodeMastery
	default:
		return NodePath
	}
}

func parseNodeEffects(effectsYAML []NodeEffectYAML) ([]NodeEffect, error) {
	effects := make([]NodeEffect, 0, len(effectsYAML))

	for _, ey := range effectsYAML {
		effect, err := parseNodeEffect(ey)
		if err != nil {
			return nil, err
		}
		if effect != nil {
			effects = append(effects, effect)
		}
	}

	return effects, nil
}

func parseNodeEffect(y NodeEffectYAML) (NodeEffect, error) {
	switch y.Type {
	case "attribute":
		return &BaseAttributeEffect{
			attribute:   parseAttributeType(y.Attribute),
			modType:     parseModifierType(y.ModType),
			value:       y.Value,
			description: y.Description,
		}, nil

	case "grant_skill":
		return &BaseGrantSkillEffect{
			skillID:     y.SkillID,
			startLevel:  y.StartLevel,
			description: y.Description,
		}, nil

	case "passive":
		return &BasePassiveEffect{
			triggerType:  parseTriggerType(y.TriggerType),
			triggerValue: y.TriggerValue,
			effectID:     y.EffectID,
			description:  y.Description,
			metadata:     y.Metadata,
		}, nil

	case "skill_mod":
		return &BaseSkillModEffect{
			targetSkillID:   y.TargetSkillID,
			targetSkillTags: y.TargetSkillTags,
			description:     y.Description,
			metadata:        y.Metadata,
		}, nil

	case "unlock_craft", "trade", "resource", "special":
		return &BaseSpecialEffect{
			effectType:  parseNodeEffectType(y.Type),
			description: y.Description,
			metadata:    y.Metadata,
		}, nil

	default:
		// Unknown effect type - store as special
		return &BaseSpecialEffect{
			effectType:  EffectTypeSpecial,
			description: y.Description,
			metadata:    y.Metadata,
		}, nil
	}
}

func parseAttributeType(s string) attribute.Type {
	// Map common attribute names
	switch s {
	case "strength":
		return attribute.AttrStrength
	case "dexterity":
		return attribute.AttrDexterity
	case "intelligence":
		return attribute.AttrIntelligence
	case "vitality":
		return attribute.AttrVitality
	case "willpower":
		return attribute.AttrWillpower
	case "armor":
		return attribute.AttrArmor
	case "evasion":
		return attribute.AttrEvasion
	case "crit_chance":
		return attribute.AttrCritChance
	case "crit_multiplier":
		return attribute.AttrCritMultiplier
	case "attack_speed":
		return attribute.AttrAttackSpeed
	case "movement_speed":
		return attribute.AttrMovementSpeed
	case "fire_resistance":
		return attribute.AttrFireResist
	case "cold_resistance":
		return attribute.AttrColdResist
	case "lightning_resistance":
		return attribute.AttrLightningResist
	case "life_regen":
		return attribute.AttrLifeRegen
	case "mana_regen":
		return attribute.AttrManaRegen
	case "loot_quantity":
		return attribute.AttrLootQuantity
	case "loot_rarity":
		return attribute.AttrLootRarity
	default:
		return attribute.Type(s) // Custom attribute
	}
}

func parseModifierType(s string) attribute.ModifierType {
	switch s {
	case "flat":
		return attribute.ModFlat
	case "increased":
		return attribute.ModIncreased
	case "more":
		return attribute.ModMore
	case "override":
		return attribute.ModOverride
	default:
		return attribute.ModFlat
	}
}

// Attribute returns the attribute type this effect modifies
func (e *BaseAttributeEffect) Attribute() attribute.Type {
	return e.attribute
}

// ModType returns the modifier type
func (e *BaseAttributeEffect) ModType() attribute.ModifierType {
	return e.modType
}

func parseTriggerType(s string) TriggerType {
	switch s {
	case "on_hit":
		return TriggerOnHit
	case "on_crit":
		return TriggerOnCrit
	case "on_kill":
		return TriggerOnKill
	case "on_damaged":
		return TriggerOnDamaged
	case "on_block":
		return TriggerOnBlock
	case "on_dodge":
		return TriggerOnDodge
	case "on_skill_use":
		return TriggerOnSkillUse
	case "on_craft":
		return TriggerOnCraft
	case "on_trade":
		return TriggerOnTrade
	case "on_level_up":
		return TriggerOnLevelUp
	case "on_rest":
		return TriggerOnRest
	case "periodic":
		return TriggerPeriodic
	case "on_low_health":
		return TriggerOnLowHealth
	case "on_full_health":
		return TriggerOnFullHealth
	default:
		return TriggerType(s)
	}
}

func parseNodeEffectType(s string) NodeEffectType {
	switch s {
	case "attribute":
		return EffectTypeAttribute
	case "skill_mod":
		return EffectTypeSkillMod
	case "grant_skill":
		return EffectTypeGrantSkill
	case "unlock_craft":
		return EffectTypeUnlockCraft
	case "trade":
		return EffectTypeTrade
	case "passive":
		return EffectTypePassive
	case "resource":
		return EffectTypeResource
	default:
		return EffectTypeSpecial
	}
}

// =============================================================================
// BASE NODE EFFECT IMPLEMENTATIONS
// =============================================================================

var _ NodeEffect = (*BaseAttributeEffect)(nil)

// BaseAttributeEffect implements NodeEffect for attribute modifications
type BaseAttributeEffect struct {
	attribute   attribute.Type
	modType     attribute.ModifierType
	value       float64
	description string
}

func (e *BaseAttributeEffect) Type() NodeEffectType { return EffectTypeAttribute }
func (e *BaseAttributeEffect) Description() string  { return e.description }
func (e *BaseAttributeEffect) Value() float64       { return e.value }
func (e *BaseAttributeEffect) Metadata() map[string]any {
	return map[string]any{
		"attribute": string(e.attribute),
		"mod_type":  string(e.modType),
	}
}

func (e *BaseAttributeEffect) Apply(ctx context.Context, entityID string) error {
	// TODO: Integrate with attribute system
	_ = ctx
	_ = entityID
	return nil
}

func (e *BaseAttributeEffect) Remove(ctx context.Context, entityID string) error {
	// TODO: Integrate with attribute system
	_ = ctx
	_ = entityID
	return nil
}

var _ NodeEffect = (*BaseGrantSkillEffect)(nil)

// BaseGrantSkillEffect implements NodeEffect for granting skills
type BaseGrantSkillEffect struct {
	skillID     string
	startLevel  int
	description string
}

func (e *BaseGrantSkillEffect) Type() NodeEffectType { return EffectTypeGrantSkill }
func (e *BaseGrantSkillEffect) Description() string  { return e.description }
func (e *BaseGrantSkillEffect) Value() float64       { return float64(e.startLevel) }
func (e *BaseGrantSkillEffect) Metadata() map[string]any {
	return map[string]any{
		"skill_id":    e.skillID,
		"start_level": e.startLevel,
	}
}

func (e *BaseGrantSkillEffect) Apply(ctx context.Context, entityID string) error {
	// TODO: Grant skill to entity
	_ = ctx
	_ = entityID
	return nil
}

func (e *BaseGrantSkillEffect) Remove(ctx context.Context, entityID string) error {
	// TODO: Remove skill from entity
	_ = ctx
	_ = entityID
	return nil
}

var _ NodeEffect = (*BasePassiveEffect)(nil)

// BasePassiveEffect implements NodeEffect for passive triggers
type BasePassiveEffect struct {
	triggerType  TriggerType
	triggerValue float64
	effectID     string
	description  string
	metadata     map[string]any
}

func (e *BasePassiveEffect) Type() NodeEffectType { return EffectTypePassive }
func (e *BasePassiveEffect) Description() string  { return e.description }
func (e *BasePassiveEffect) Value() float64       { return e.triggerValue }
func (e *BasePassiveEffect) Metadata() map[string]any {
	result := map[string]any{
		"trigger_type":  string(e.triggerType),
		"trigger_value": e.triggerValue,
		"effect_id":     e.effectID,
	}
	for k, v := range e.metadata {
		result[k] = v
	}
	return result
}

func (e *BasePassiveEffect) Apply(ctx context.Context, entityID string) error {
	// TODO: Register passive trigger
	_ = ctx
	_ = entityID
	return nil
}

func (e *BasePassiveEffect) Remove(ctx context.Context, entityID string) error {
	// TODO: Unregister passive trigger
	_ = ctx
	_ = entityID
	return nil
}

var _ NodeEffect = (*BaseSkillModEffect)(nil)

// BaseSkillModEffect implements NodeEffect for skill modifications
type BaseSkillModEffect struct {
	targetSkillID   string
	targetSkillTags []string
	description     string
	metadata        map[string]any
}

func (e *BaseSkillModEffect) Type() NodeEffectType { return EffectTypeSkillMod }
func (e *BaseSkillModEffect) Description() string  { return e.description }
func (e *BaseSkillModEffect) Value() float64       { return 0 }
func (e *BaseSkillModEffect) Metadata() map[string]any {
	result := map[string]any{
		"target_skill_id":   e.targetSkillID,
		"target_skill_tags": e.targetSkillTags,
	}
	for k, v := range e.metadata {
		result[k] = v
	}
	return result
}

func (e *BaseSkillModEffect) Apply(ctx context.Context, entityID string) error {
	// TODO: Apply skill modifier
	_ = ctx
	_ = entityID
	return nil
}

func (e *BaseSkillModEffect) Remove(ctx context.Context, entityID string) error {
	// TODO: Remove skill modifier
	_ = ctx
	_ = entityID
	return nil
}

var _ NodeEffect = (*BaseSpecialEffect)(nil)

// BaseSpecialEffect implements NodeEffect for special/custom effects
type BaseSpecialEffect struct {
	effectType  NodeEffectType
	description string
	metadata    map[string]any
}

func (e *BaseSpecialEffect) Type() NodeEffectType     { return e.effectType }
func (e *BaseSpecialEffect) Description() string      { return e.description }
func (e *BaseSpecialEffect) Value() float64           { return 0 }
func (e *BaseSpecialEffect) Metadata() map[string]any { return e.metadata }

func (e *BaseSpecialEffect) Apply(ctx context.Context, entityID string) error {
	// TODO: Handle special effect
	_ = ctx
	_ = entityID
	return nil
}

func (e *BaseSpecialEffect) Remove(ctx context.Context, entityID string) error {
	// TODO: Handle special effect removal
	_ = ctx
	_ = entityID
	return nil
}

// =============================================================================
// GLOBAL TREE REGISTRY
// =============================================================================

var (
	globalTreeRegistry *BaseTreeRegistry
	globalTreeOnce     sync.Once
)

// GlobalTreeRegistry returns the global tree registry
func GlobalTreeRegistry() *BaseTreeRegistry {
	globalTreeOnce.Do(func() {
		globalTreeRegistry = NewBaseTreeRegistry()
	})
	return globalTreeRegistry
}
