package skill

import (
	"context"
	"sync"
)

var _ Def = (*BaseDef)(nil)

// BaseDef implements the Def interface.
// Represents an immutable skill template loaded from YAML.
type BaseDef struct {
	mu sync.RWMutex

	id          string
	name        string
	description string
	skillType   Type
	tags        []string
	icon        string

	maxLevel   int
	levelData  map[int]*BaseLevelData
	baseCD     int64
	baseCharge int
	chargeCD   int64

	targeting    *BaseTargetRule
	effects      []*BaseEffectDef
	requirements *BaseRequirements

	metadata map[string]any
}

// DefConfig holds configuration for creating BaseDef
type DefConfig struct {
	ID          string
	Name        string
	Description string
	Type        Type
	Tags        []string
	Icon        string

	MaxLevel       int
	BaseCooldown   int64
	BaseCharges    int
	ChargeRecovery int64

	Targeting    *BaseTargetRule
	Effects      []*BaseEffectDef
	Requirements *BaseRequirements
	Metadata     map[string]any
}

// NewBaseDef creates a new skill definition
func NewBaseDef(config DefConfig) *BaseDef {
	def := &BaseDef{
		id:           config.ID,
		name:         config.Name,
		description:  config.Description,
		skillType:    config.Type,
		tags:         config.Tags,
		icon:         config.Icon,
		maxLevel:     config.MaxLevel,
		levelData:    make(map[int]*BaseLevelData),
		baseCD:       config.BaseCooldown,
		baseCharge:   config.BaseCharges,
		chargeCD:     config.ChargeRecovery,
		targeting:    config.Targeting,
		effects:      config.Effects,
		requirements: config.Requirements,
		metadata:     config.Metadata,
	}

	if def.targeting == nil {
		def.targeting = &BaseTargetRule{targetType: TargetSelf}
	}

	if def.metadata == nil {
		def.metadata = make(map[string]any)
	}

	return def
}

func (d *BaseDef) ID() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.id
}

func (d *BaseDef) Name() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.name
}

func (d *BaseDef) Description() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.description
}

func (d *BaseDef) Type() Type {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.skillType
}

func (d *BaseDef) Tags() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]string, len(d.tags))
	copy(result, d.tags)
	return result
}

func (d *BaseDef) HasTag(tag string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return HasTag(d.tags, tag)
}

func (d *BaseDef) MaxLevel() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.maxLevel
}

func (d *BaseDef) LevelData(level int) LevelData {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if data, ok := d.levelData[level]; ok {
		return data
	}
	return nil
}

// SetLevelData adds level data to definition
func (d *BaseDef) SetLevelData(level int, data *BaseLevelData) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.levelData[level] = data
}

func (d *BaseDef) BaseCooldown() int64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.baseCD
}

func (d *BaseDef) BaseCharges() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.baseCharge
}

func (d *BaseDef) ChargeRecovery() int64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.chargeCD
}

func (d *BaseDef) Targeting() TargetRule {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.targeting
}

func (d *BaseDef) Effects() []EffectDef {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make([]EffectDef, len(d.effects))
	for i, e := range d.effects {
		result[i] = e
	}
	return result
}

func (d *BaseDef) Requirements() Requirements {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.requirements
}

func (d *BaseDef) Icon() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.icon
}

func (d *BaseDef) Metadata() map[string]any {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make(map[string]any, len(d.metadata))
	for k, v := range d.metadata {
		result[k] = v
	}
	return result
}

// =============================================================================
// BASE LEVEL DATA
// =============================================================================

var _ LevelData = (*BaseLevelData)(nil)

// BaseLevelData implements LevelData interface
type BaseLevelData struct {
	level       int
	costs       []ResourceCost
	effects     []EffectValue
	cooldown    int64
	charges     int
	description string
	metadata    map[string]any
}

// LevelDataConfig holds configuration for level data
type LevelDataConfig struct {
	Level       int
	Costs       []ResourceCost
	Effects     []EffectValue
	Cooldown    int64
	Charges     int
	Description string
	Metadata    map[string]any
}

// NewBaseLevelData creates level data
func NewBaseLevelData(config LevelDataConfig) *BaseLevelData {
	return &BaseLevelData{
		level:       config.Level,
		costs:       config.Costs,
		effects:     config.Effects,
		cooldown:    config.Cooldown,
		charges:     config.Charges,
		description: config.Description,
		metadata:    config.Metadata,
	}
}

func (l *BaseLevelData) Level() int {
	return l.level
}

func (l *BaseLevelData) ResourceCosts() []ResourceCost {
	result := make([]ResourceCost, len(l.costs))
	copy(result, l.costs)
	return result
}

func (l *BaseLevelData) Effects() []EffectValue {
	result := make([]EffectValue, len(l.effects))
	copy(result, l.effects)
	return result
}

func (l *BaseLevelData) Cooldown() int64 {
	return l.cooldown
}

func (l *BaseLevelData) Charges() int {
	return l.charges
}

func (l *BaseLevelData) Description() string {
	return l.description
}

func (l *BaseLevelData) Metadata() map[string]any {
	result := make(map[string]any, len(l.metadata))
	for k, v := range l.metadata {
		result[k] = v
	}
	return result
}

// =============================================================================
// BASE TARGET RULE
// =============================================================================

var _ TargetRule = (*BaseTargetRule)(nil)

// BaseTargetRule implements TargetRule interface
type BaseTargetRule struct {
	targetType  TargetType
	areaType    AreaType
	rangeVal    float64
	areaRadius  float64
	maxTargets  int
	minTargets  int
	canSelf     bool
	canAllies   bool
	canEnemies  bool
	requiresLOS bool
	chainCount  int
	chainFallof float64
}

// TargetRuleConfig holds targeting configuration
type TargetRuleConfig struct {
	Type        TargetType
	AreaType    AreaType
	Range       float64
	AreaRadius  float64
	MaxTargets  int
	MinTargets  int
	CanSelf     bool
	CanAllies   bool
	CanEnemies  bool
	RequiresLOS bool
	ChainCount  int
	ChainFallof float64
}

// NewBaseTargetRule creates target rule
func NewBaseTargetRule(config TargetRuleConfig) *BaseTargetRule {
	return &BaseTargetRule{
		targetType:  config.Type,
		areaType:    config.AreaType,
		rangeVal:    config.Range,
		areaRadius:  config.AreaRadius,
		maxTargets:  config.MaxTargets,
		minTargets:  config.MinTargets,
		canSelf:     config.CanSelf,
		canAllies:   config.CanAllies,
		canEnemies:  config.CanEnemies,
		requiresLOS: config.RequiresLOS,
		chainCount:  config.ChainCount,
		chainFallof: config.ChainFallof,
	}
}

func (r *BaseTargetRule) Type() TargetType      { return r.targetType }
func (r *BaseTargetRule) AreaType() AreaType    { return r.areaType }
func (r *BaseTargetRule) Range() float64        { return r.rangeVal }
func (r *BaseTargetRule) AreaRadius() float64   { return r.areaRadius }
func (r *BaseTargetRule) MaxTargets() int       { return r.maxTargets }
func (r *BaseTargetRule) MinTargets() int       { return r.minTargets }
func (r *BaseTargetRule) CanTargetSelf() bool   { return r.canSelf }
func (r *BaseTargetRule) CanTargetAllies() bool { return r.canAllies }
func (r *BaseTargetRule) CanTargetEnemies() bool {
	return r.canEnemies
}
func (r *BaseTargetRule) RequiresLineOfSight() bool { return r.requiresLOS }
func (r *BaseTargetRule) ChainCount() int           { return r.chainCount }
func (r *BaseTargetRule) ChainFalloff() float64     { return r.chainFallof }

// =============================================================================
// BASE EFFECT DEF
// =============================================================================

var _ EffectDef = (*BaseEffectDef)(nil)

// BaseEffectDef implements EffectDef interface
type BaseEffectDef struct {
	id         string
	effectType EffectType
	damageType string
	statusID   string
	scaling    []ScalingRule
	chance     float64
	delay      int64
	duration   int64
	metadata   map[string]any
}

// EffectDefConfig holds effect configuration
type EffectDefConfig struct {
	ID         string
	Type       EffectType
	DamageType string
	StatusID   string
	Scaling    []ScalingRule
	Chance     float64
	Delay      int64
	Duration   int64
	Metadata   map[string]any
}

// NewBaseEffectDef creates effect definition
func NewBaseEffectDef(config EffectDefConfig) *BaseEffectDef {
	chance := config.Chance
	if chance == 0 {
		chance = 1.0 // Default to 100% chance
	}

	return &BaseEffectDef{
		id:         config.ID,
		effectType: config.Type,
		damageType: config.DamageType,
		statusID:   config.StatusID,
		scaling:    config.Scaling,
		chance:     chance,
		delay:      config.Delay,
		duration:   config.Duration,
		metadata:   config.Metadata,
	}
}

func (e *BaseEffectDef) ID() string               { return e.id }
func (e *BaseEffectDef) Type() EffectType         { return e.effectType }
func (e *BaseEffectDef) DamageType() string       { return e.damageType }
func (e *BaseEffectDef) StatusID() string         { return e.statusID }
func (e *BaseEffectDef) Scaling() []ScalingRule   { return e.scaling }
func (e *BaseEffectDef) Chance() float64          { return e.chance }
func (e *BaseEffectDef) Delay() int64             { return e.delay }
func (e *BaseEffectDef) Duration() int64          { return e.duration }
func (e *BaseEffectDef) Metadata() map[string]any { return e.metadata }

// =============================================================================
// BASE REQUIREMENTS
// =============================================================================

var _ Requirements = (*BaseRequirements)(nil)

// BaseRequirements implements Requirements interface
type BaseRequirements struct {
	characterLevel int
	attributes     map[string]float64
	skills         map[string]int
	nodes          []string
	items          []string
}

// RequirementsConfig holds requirements configuration
type RequirementsConfig struct {
	CharacterLevel int
	Attributes     map[string]float64
	Skills         map[string]int
	Nodes          []string
	Items          []string
}

// NewBaseRequirements creates requirements
func NewBaseRequirements(config RequirementsConfig) *BaseRequirements {
	return &BaseRequirements{
		characterLevel: config.CharacterLevel,
		attributes:     config.Attributes,
		skills:         config.Skills,
		nodes:          config.Nodes,
		items:          config.Items,
	}
}

func (r *BaseRequirements) CharacterLevel() int { return r.characterLevel }

func (r *BaseRequirements) Attributes() map[string]float64 {
	result := make(map[string]float64, len(r.attributes))
	for k, v := range r.attributes {
		result[k] = v
	}
	return result
}

func (r *BaseRequirements) Skills() map[string]int {
	result := make(map[string]int, len(r.skills))
	for k, v := range r.skills {
		result[k] = v
	}
	return result
}

func (r *BaseRequirements) Nodes() []string {
	result := make([]string, len(r.nodes))
	copy(result, r.nodes)
	return result
}

func (r *BaseRequirements) Items() []string {
	result := make([]string, len(r.items))
	copy(result, r.items)
	return result
}

// Check verifies requirements - placeholder implementation
// TODO: Integrate with entity/character system
func (r *BaseRequirements) Check(ctx context.Context, entityID string) bool {
	// For now, return true - actual implementation needs entity lookup
	_ = ctx
	return true
}
