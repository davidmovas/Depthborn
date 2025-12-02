package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// =============================================================================
// REGISTRY
// =============================================================================

// Registry manages skill definitions loaded from YAML
type Registry interface {
	// Register adds a skill definition
	Register(def Def) error

	// Get retrieves skill definition by ID
	Get(id string) (Def, bool)

	// GetAll returns all registered definitions
	GetAll() []Def

	// GetByTag returns definitions with specific tag
	GetByTag(tag string) []Def

	// GetByType returns definitions of specific type
	GetByType(skillType Type) []Def

	// Has checks if skill is registered
	Has(id string) bool

	// Count returns number of registered skills
	Count() int

	// CreateInstance creates a new skill instance from definition
	CreateInstance(defID string, level int) (Instance, error)

	// LoadFromYAML loads skills from YAML data
	LoadFromYAML(data []byte) error

	// LoadFromFile loads skills from YAML file
	LoadFromFile(path string) error

	// LoadFromDirectory loads all YAML files from directory
	LoadFromDirectory(dir string) error
}

// =============================================================================
// BASE REGISTRY
// =============================================================================

var _ Registry = (*BaseRegistry)(nil)

// BaseRegistry implements Registry interface
type BaseRegistry struct {
	mu     sync.RWMutex
	skills map[string]Def
}

// NewBaseRegistry creates a new skill registry
func NewBaseRegistry() *BaseRegistry {
	return &BaseRegistry{
		skills: make(map[string]Def),
	}
}

func (r *BaseRegistry) Register(def Def) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.skills[def.ID()]; exists {
		return fmt.Errorf("skill %s already registered", def.ID())
	}

	r.skills[def.ID()] = def
	return nil
}

func (r *BaseRegistry) Get(id string) (Def, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.skills[id]
	return def, ok
}

func (r *BaseRegistry) GetAll() []Def {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Def, 0, len(r.skills))
	for _, def := range r.skills {
		result = append(result, def)
	}
	return result
}

func (r *BaseRegistry) GetByTag(tag string) []Def {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Def
	for _, def := range r.skills {
		if def.HasTag(tag) {
			result = append(result, def)
		}
	}
	return result
}

func (r *BaseRegistry) GetByType(skillType Type) []Def {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Def
	for _, def := range r.skills {
		if def.Type() == skillType {
			result = append(result, def)
		}
	}
	return result
}

func (r *BaseRegistry) Has(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.skills[id]
	return ok
}

func (r *BaseRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.skills)
}

func (r *BaseRegistry) CreateInstance(defID string, level int) (Instance, error) {
	r.mu.RLock()
	def, ok := r.skills[defID]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("skill definition %s not found", defID)
	}

	return NewBaseInstance(InstanceConfig{
		Def:        def,
		StartLevel: level,
	}), nil
}

// =============================================================================
// YAML LOADING
// =============================================================================

// SkillFile represents YAML file structure
type SkillFile struct {
	Version string      `yaml:"version"`
	Skills  []SkillYAML `yaml:"skills"`
}

// SkillYAML represents skill definition in YAML
type SkillYAML struct {
	ID           string            `yaml:"id"`
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Type         string            `yaml:"type"`
	Tags         []string          `yaml:"tags"`
	Icon         string            `yaml:"icon"`
	MaxLevel     int               `yaml:"max_level"`
	Cooldown     int64             `yaml:"cooldown"`
	Charges      int               `yaml:"charges"`
	ChargeCD     int64             `yaml:"charge_recovery"`
	Targeting    *TargetingYAML    `yaml:"targeting"`
	Effects      []EffectYAML      `yaml:"effects"`
	Levels       []LevelYAML       `yaml:"levels"`
	Requirements *RequirementsYAML `yaml:"requirements"`
	Metadata     map[string]any    `yaml:"metadata"`
}

// TargetingYAML represents targeting in YAML
type TargetingYAML struct {
	Type        string  `yaml:"type"`
	AreaType    string  `yaml:"area_type"`
	Range       float64 `yaml:"range"`
	AreaRadius  float64 `yaml:"area_radius"`
	MaxTargets  int     `yaml:"max_targets"`
	MinTargets  int     `yaml:"min_targets"`
	CanSelf     bool    `yaml:"can_self"`
	CanAllies   bool    `yaml:"can_allies"`
	CanEnemies  bool    `yaml:"can_enemies"`
	RequiresLOS bool    `yaml:"requires_los"`
	ChainCount  int     `yaml:"chain_count"`
	ChainFallof float64 `yaml:"chain_falloff"`
}

// EffectYAML represents effect in YAML
type EffectYAML struct {
	ID         string         `yaml:"id"`
	Type       string         `yaml:"type"`
	DamageType string         `yaml:"damage_type"`
	StatusID   string         `yaml:"status_id"`
	Scaling    []ScalingYAML  `yaml:"scaling"`
	Chance     float64        `yaml:"chance"`
	Delay      int64          `yaml:"delay"`
	Duration   int64          `yaml:"duration"`
	Metadata   map[string]any `yaml:"metadata"`
}

// ScalingYAML represents scaling rule in YAML
type ScalingYAML struct {
	Attribute  string  `yaml:"attribute"`
	Multiplier float64 `yaml:"multiplier"`
}

// LevelYAML represents level data in YAML
type LevelYAML struct {
	Level       int               `yaml:"level"`
	Costs       []ResourceYAML    `yaml:"costs"`
	Effects     []EffectValueYAML `yaml:"effects"`
	Cooldown    int64             `yaml:"cooldown"`
	Charges     int               `yaml:"charges"`
	Description string            `yaml:"description"`
	Metadata    map[string]any    `yaml:"metadata"`
}

// ResourceYAML represents resource cost in YAML
type ResourceYAML struct {
	Resource string  `yaml:"resource"`
	Type     string  `yaml:"type"`
	Amount   float64 `yaml:"amount"`
}

// EffectValueYAML represents effect values in YAML
type EffectValueYAML struct {
	EffectID string         `yaml:"effect_id"`
	Values   map[string]any `yaml:"values"`
}

// RequirementsYAML represents requirements in YAML
type RequirementsYAML struct {
	CharacterLevel int                `yaml:"character_level"`
	Attributes     map[string]float64 `yaml:"attributes"`
	Skills         map[string]int     `yaml:"skills"`
	Nodes          []string           `yaml:"nodes"`
	Items          []string           `yaml:"items"`
}

func (r *BaseRegistry) LoadFromYAML(data []byte) error {
	var file SkillFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	for _, skillYAML := range file.Skills {
		def, err := parseSkillYAML(skillYAML)
		if err != nil {
			return fmt.Errorf("failed to parse skill %s: %w", skillYAML.ID, err)
		}

		if err := r.Register(def); err != nil {
			return err
		}
	}

	return nil
}

func (r *BaseRegistry) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return r.LoadFromYAML(data)
}

func (r *BaseRegistry) LoadFromDirectory(dir string) error {
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

func parseSkillYAML(y SkillYAML) (*BaseDef, error) {
	skillType := parseSkillType(y.Type)

	// Parse targeting
	var targeting *BaseTargetRule
	if y.Targeting != nil {
		targeting = parseTargetingYAML(y.Targeting)
	} else {
		targeting = &BaseTargetRule{targetType: TargetSelf}
	}

	// Parse effects
	effects := make([]*BaseEffectDef, 0, len(y.Effects))
	for _, ey := range y.Effects {
		effect := parseEffectYAML(ey)
		effects = append(effects, effect)
	}

	// Parse requirements
	var requirements *BaseRequirements
	if y.Requirements != nil {
		requirements = parseRequirementsYAML(y.Requirements)
	}

	def := NewBaseDef(DefConfig{
		ID:             y.ID,
		Name:           y.Name,
		Description:    y.Description,
		Type:           skillType,
		Tags:           y.Tags,
		Icon:           y.Icon,
		MaxLevel:       y.MaxLevel,
		BaseCooldown:   y.Cooldown,
		BaseCharges:    y.Charges,
		ChargeRecovery: y.ChargeCD,
		Targeting:      targeting,
		Effects:        effects,
		Requirements:   requirements,
		Metadata:       y.Metadata,
	})

	// Parse level data
	for _, ly := range y.Levels {
		levelData := parseLevelYAML(ly)
		def.SetLevelData(ly.Level, levelData)
	}

	return def, nil
}

func parseSkillType(s string) Type {
	switch s {
	case "active":
		return TypeActive
	case "passive":
		return TypePassive
	case "aura":
		return TypeAura
	case "trigger":
		return TypeTrigger
	default:
		return TypeActive
	}
}

func parseTargetingYAML(y *TargetingYAML) *BaseTargetRule {
	return NewBaseTargetRule(TargetRuleConfig{
		Type:        parseTargetType(y.Type),
		AreaType:    parseAreaType(y.AreaType),
		Range:       y.Range,
		AreaRadius:  y.AreaRadius,
		MaxTargets:  y.MaxTargets,
		MinTargets:  y.MinTargets,
		CanSelf:     y.CanSelf,
		CanAllies:   y.CanAllies,
		CanEnemies:  y.CanEnemies,
		RequiresLOS: y.RequiresLOS,
		ChainCount:  y.ChainCount,
		ChainFallof: y.ChainFallof,
	})
}

func parseTargetType(s string) TargetType {
	switch s {
	case "none":
		return TargetNone
	case "self":
		return TargetSelf
	case "single":
		return TargetSingle
	case "multiple":
		return TargetMultiple
	case "all_enemies":
		return TargetAllEnemies
	case "all_allies":
		return TargetAllAllies
	case "all":
		return TargetAll
	case "ground":
		return TargetGround
	default:
		return TargetSelf
	}
}

func parseAreaType(s string) AreaType {
	switch s {
	case "none":
		return AreaNone
	case "circle":
		return AreaCircle
	case "cone":
		return AreaCone
	case "line":
		return AreaLine
	case "chain":
		return AreaChain
	default:
		return AreaNone
	}
}

func parseEffectYAML(y EffectYAML) *BaseEffectDef {
	scaling := make([]ScalingRule, 0, len(y.Scaling))
	for _, sy := range y.Scaling {
		scaling = append(scaling, ScalingRule{
			Attribute:  sy.Attribute,
			Multiplier: sy.Multiplier,
		})
	}

	return NewBaseEffectDef(EffectDefConfig{
		ID:         y.ID,
		Type:       parseEffectType(y.Type),
		DamageType: y.DamageType,
		StatusID:   y.StatusID,
		Scaling:    scaling,
		Chance:     y.Chance,
		Delay:      y.Delay,
		Duration:   y.Duration,
		Metadata:   y.Metadata,
	})
}

func parseEffectType(s string) EffectType {
	switch s {
	case "damage":
		return EffectDamage
	case "heal":
		return EffectHeal
	case "status":
		return EffectStatus
	case "buff":
		return EffectBuff
	case "debuff":
		return EffectDebuff
	case "summon":
		return EffectSummon
	case "teleport":
		return EffectTeleport
	case "knockback":
		return EffectKnockback
	case "pull":
		return EffectPull
	case "modify_skill":
		return EffectModifySkill
	case "resource":
		return EffectResource
	case "dispel":
		return EffectDispel
	default:
		return EffectDamage
	}
}

func parseLevelYAML(y LevelYAML) *BaseLevelData {
	costs := make([]ResourceCost, 0, len(y.Costs))
	for _, cy := range y.Costs {
		costs = append(costs, ResourceCost{
			Resource: parseResourceType(cy.Resource),
			Type:     parseCostType(cy.Type),
			Amount:   cy.Amount,
		})
	}

	effects := make([]EffectValue, 0, len(y.Effects))
	for _, ey := range y.Effects {
		effects = append(effects, EffectValue{
			EffectID: ey.EffectID,
			Values:   ey.Values,
		})
	}

	return NewBaseLevelData(LevelDataConfig{
		Level:       y.Level,
		Costs:       costs,
		Effects:     effects,
		Cooldown:    y.Cooldown,
		Charges:     y.Charges,
		Description: y.Description,
		Metadata:    y.Metadata,
	})
}

func parseResourceType(s string) ResourceType {
	switch s {
	case "mana":
		return ResourceMana
	case "health":
		return ResourceHealth
	case "stamina":
		return ResourceStamina
	case "rage":
		return ResourceRage
	case "energy":
		return ResourceEnergy
	case "soul_charge":
		return ResourceSoulCharge
	case "gold":
		return ResourceGold
	default:
		return ResourceMana
	}
}

func parseCostType(s string) CostType {
	switch s {
	case "flat":
		return CostFlat
	case "percent":
		return CostPercent
	case "current":
		return CostCurrent
	default:
		return CostFlat
	}
}

func parseRequirementsYAML(y *RequirementsYAML) *BaseRequirements {
	return NewBaseRequirements(RequirementsConfig{
		CharacterLevel: y.CharacterLevel,
		Attributes:     y.Attributes,
		Skills:         y.Skills,
		Nodes:          y.Nodes,
		Items:          y.Items,
	})
}

// =============================================================================
// GLOBAL REGISTRY
// =============================================================================

var (
	globalRegistry *BaseRegistry
	globalOnce     sync.Once
)

// GlobalRegistry returns the global skill registry
func GlobalRegistry() *BaseRegistry {
	globalOnce.Do(func() {
		globalRegistry = NewBaseRegistry()
	})
	return globalRegistry
}
