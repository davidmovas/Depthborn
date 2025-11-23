package combat

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

// AI controls AI combat decisions
type AI interface {
	// SelectAction chooses action for AI participant
	SelectAction(ctx context.Context, participant Participant, encounter Encounter) (Action, error)

	// SelectTarget chooses target for action
	SelectTarget(ctx context.Context, participant Participant, action Action, encounter Encounter) ([]string, error)

	// SelectPosition chooses movement destination
	SelectPosition(ctx context.Context, participant Participant, encounter Encounter) (spatial.Position, error)

	// EvaluateThreat assesses threat level of entities
	EvaluateThreat(ctx context.Context, participant Participant, encounter Encounter) map[string]float64

	// ShouldFlee determines if AI should attempt escape
	ShouldFlee(ctx context.Context, participant Participant, encounter Encounter) bool

	// ShouldUseSkill determines if AI should use skill
	ShouldUseSkill(ctx context.Context, participant Participant, skillID string, encounter Encounter) bool

	// ShouldDefend determines if AI should enter defensive stance
	ShouldDefend(ctx context.Context, participant Participant, encounter Encounter) bool

	// ShouldUseItem determines if AI should use item
	ShouldUseItem(ctx context.Context, participant Participant, itemID string, encounter Encounter) bool

	// GetStrategy returns AI strategy for encounter
	GetStrategy(ctx context.Context, participant Participant, encounter Encounter) AIStrategy

	// SetStrategy updates AI strategy
	SetStrategy(strategy AIStrategy)

	// Update updates AI state for frame
	Update(ctx context.Context, participant Participant, encounter Encounter, deltaMs int64) error
}

// AIStrategy defines AI combat approach
type AIStrategy interface {
	// ID returns unique strategy identifier
	ID() string

	// Name returns strategy name
	Name() string

	// Type returns strategy type
	Type() StrategyType

	// Priority returns strategy priority (higher = preferred)
	Priority() int

	// EvaluateActions scores available actions
	EvaluateActions(ctx context.Context, participant Participant, actions []Action, encounter Encounter) map[string]float64

	// SelectBestAction chooses highest scoring action
	SelectBestAction(ctx context.Context, scores map[string]float64, encounter Encounter) (Action, error)

	// ShouldSwitch determines if strategy should change
	ShouldSwitch(ctx context.Context, participant Participant, encounter Encounter) bool

	// GetAlternativeStrategy returns fallback strategy
	GetAlternativeStrategy() AIStrategy

	// CanExecute checks if strategy can be executed
	CanExecute(ctx context.Context, participant Participant, encounter Encounter) bool

	// OnActivate is called when strategy becomes active
	OnActivate(ctx context.Context, participant Participant, encounter Encounter) error

	// OnDeactivate is called when strategy is switched
	OnDeactivate(ctx context.Context, participant Participant, encounter Encounter) error
}

// StrategyType categorizes AI strategies
type StrategyType string

const (
	StrategyAggressive StrategyType = "aggressive"
	StrategyDefensive  StrategyType = "defensive"
	StrategyBalanced   StrategyType = "balanced"
	StrategySupport    StrategyType = "support"
	StrategyCautious   StrategyType = "cautious"
	StrategyBerserk    StrategyType = "berserk"
	StrategyTactical   StrategyType = "tactical"
	StrategyCowardly   StrategyType = "cowardly"
	StrategyAssassin   StrategyType = "assassin"
	StrategyTank       StrategyType = "tank"
	StrategyController StrategyType = "controller"
	StrategyRanged     StrategyType = "ranged"
	StrategyAmbusher   StrategyType = "ambusher"
)

// TargetSelector chooses combat targets
type TargetSelector interface {
	// SelectTarget chooses best target
	SelectTarget(ctx context.Context, participant Participant, candidates []Participant, encounter Encounter) (Participant, error)

	// SelectTargets chooses multiple targets
	SelectTargets(ctx context.Context, participant Participant, candidates []Participant, count int, encounter Encounter) ([]Participant, error)

	// ScoreTargets evaluates target priority
	ScoreTargets(ctx context.Context, participant Participant, candidates []Participant, encounter Encounter) map[string]float64

	// FilterTargets removes invalid targets
	FilterTargets(ctx context.Context, participant Participant, candidates []Participant, encounter Encounter) []Participant

	// PreferredTargets returns target preferences
	PreferredTargets() []TargetPreference

	// AddPreference adds target preference
	AddPreference(preference TargetPreference)

	// RemovePreference removes target preference
	RemovePreference(preferenceType PreferenceType)

	// GetPrimaryTarget returns current primary target
	GetPrimaryTarget() string

	// SetPrimaryTarget updates primary target
	SetPrimaryTarget(targetID string)
}

// TargetPreference defines targeting priority
type TargetPreference interface {
	// Type returns preference type
	Type() PreferenceType

	// Weight returns preference weight (higher = more important)
	Weight() float64

	// SetWeight updates preference weight
	SetWeight(weight float64)

	// Evaluate scores target based on preference
	Evaluate(ctx context.Context, participant, target Participant, encounter Encounter) float64

	// IsValid checks if preference applies to target
	IsValid(ctx context.Context, target Participant, encounter Encounter) bool
}

// PreferenceType categorizes target preferences
type PreferenceType string

const (
	PreferLowestHealth  PreferenceType = "lowest_health"
	PreferHighestThreat PreferenceType = "highest_threat"
	PreferClosest       PreferenceType = "closest"
	PreferFarthest      PreferenceType = "farthest"
	PreferWeakest       PreferenceType = "weakest"
	PreferStrongest     PreferenceType = "strongest"
	PreferSupport       PreferenceType = "support"
	PreferDamageDealer  PreferenceType = "damage_dealer"
	PreferTank          PreferenceType = "tank"
	PreferIsolated      PreferenceType = "isolated"
	PreferDebuffed      PreferenceType = "debuffed"
	PreferBuffed        PreferenceType = "buffed"
	PreferHighValue     PreferenceType = "high_value"
	PreferLastAttacker  PreferenceType = "last_attacker"
	PreferCurrentTarget PreferenceType = "current_target"
	PreferLowDefense    PreferenceType = "low_defense"
	PreferHighDefense   PreferenceType = "high_defense"
)

// TacticalEvaluator assesses tactical situation
type TacticalEvaluator interface {
	// EvaluatePosition scores position quality
	EvaluatePosition(ctx context.Context, participant Participant, pos spatial.Position, encounter Encounter) float64

	// FindBestPosition finds optimal position
	FindBestPosition(ctx context.Context, participant Participant, encounter Encounter) (spatial.Position, error)

	// ShouldReposition determines if should move
	ShouldReposition(ctx context.Context, participant Participant, encounter Encounter) bool

	// AnalyzeThreat assesses danger level
	AnalyzeThreat(ctx context.Context, participant Participant, encounter Encounter) ThreatLevel

	// FindCover locates defensive positions
	FindCover(ctx context.Context, participant Participant, encounter Encounter) []spatial.Position

	// FindFlankingPosition finds position to flank enemy
	FindFlankingPosition(ctx context.Context, participant Participant, target Participant, encounter Encounter) (spatial.Position, error)

	// EvaluateTeamComposition analyzes team strengths and weaknesses
	EvaluateTeamComposition(ctx context.Context, team Team, encounter Encounter) TeamComposition

	// ShouldFocusFire determines if should coordinate attacks
	ShouldFocusFire(ctx context.Context, participant Participant, encounter Encounter) bool

	// GetRetreatPosition finds safe retreat location
	GetRetreatPosition(ctx context.Context, participant Participant, encounter Encounter) (spatial.Position, error)

	// IsPositionSafe checks if position is safe
	IsPositionSafe(ctx context.Context, pos spatial.Position, participant Participant, encounter Encounter) bool

	// CountNearbyAllies returns allies within range
	CountNearbyAllies(ctx context.Context, participant Participant, radius float64, encounter Encounter) int

	// CountNearbyEnemies returns enemies within range
	CountNearbyEnemies(ctx context.Context, participant Participant, radius float64, encounter Encounter) int
}

// ThreatLevel categorizes danger
type ThreatLevel int

const (
	ThreatNone ThreatLevel = iota
	ThreatLow
	ThreatModerate
	ThreatHigh
	ThreatCritical
	ThreatLethal
)

// String returns threat level name
func (t ThreatLevel) String() string {
	levels := []string{"none", "low", "moderate", "high", "critical", "lethal"}
	if int(t) < len(levels) {
		return levels[t]
	}
	return "unknown"
}

// TeamComposition analyzes team makeup
type TeamComposition struct {
	TankCount       int
	DamageCount     int
	SupportCount    int
	AverageThreat   float64
	TotalHealth     float64
	TotalDamage     float64
	HasHealer       bool
	HasCrowdControl bool
	Formation       string
	Coordination    float64
}

// ActionEvaluator scores available actions
type ActionEvaluator interface {
	// EvaluateAction scores single action
	EvaluateAction(ctx context.Context, participant Participant, action Action, encounter Encounter) float64

	// EvaluateActions scores all available actions
	EvaluateActions(ctx context.Context, participant Participant, actions []Action, encounter Encounter) map[string]float64

	// GetUtilityScore calculates action utility
	GetUtilityScore(ctx context.Context, action Action, participant Participant, encounter Encounter) float64

	// GetRiskScore calculates action risk
	GetRiskScore(ctx context.Context, action Action, participant Participant, encounter Encounter) float64

	// GetRewardScore calculates action reward
	GetRewardScore(ctx context.Context, action Action, participant Participant, encounter Encounter) float64

	// ShouldConsiderAction checks if action is worth considering
	ShouldConsiderAction(ctx context.Context, action Action, participant Participant, encounter Encounter) bool

	// CompareActions compares two actions
	CompareActions(action1, action2 Action, scores map[string]float64) int
}

// DecisionTree represents AI decision making structure
type DecisionTree interface {
	// Root returns root node
	Root() DecisionNode

	// SetRoot assigns root node
	SetRoot(node DecisionNode)

	// Evaluate traverses tree and returns decision
	Evaluate(ctx context.Context, participant Participant, encounter Encounter) (Action, error)

	// AddNode adds node to tree
	AddNode(parent, child DecisionNode) error

	// RemoveNode removes node from tree
	RemoveNode(nodeID string) error

	// GetNode retrieves node by ID
	GetNode(nodeID string) (DecisionNode, bool)

	// Depth returns tree depth
	Depth() int

	// NodeCount returns total nodes
	NodeCount() int
}

// DecisionNode represents node in decision tree
type DecisionNode interface {
	// ID returns unique node identifier
	ID() string

	// Type returns node type
	Type() DecisionNodeType

	// Evaluate evaluates node condition or returns value
	Evaluate(ctx context.Context, participant Participant, encounter Encounter) (interface{}, error)

	// Children returns child nodes
	Children() []DecisionNode

	// AddChild adds child node
	AddChild(child DecisionNode)

	// RemoveChild removes child node
	RemoveChild(childID string)

	// Parent returns parent node
	Parent() DecisionNode

	// SetParent assigns parent node
	SetParent(parent DecisionNode)

	// IsLeaf returns true if node has no children
	IsLeaf() bool
}

// DecisionNodeType categorizes decision nodes
type DecisionNodeType string

const (
	NodeCondition DecisionNodeType = "condition"
	NodeAction    DecisionNodeType = "action"
	NodeSequence  DecisionNodeType = "sequence"
	NodeSelector  DecisionNodeType = "selector"
	NodeRandom    DecisionNodeType = "random"
)

// Behavior defines AI behavior pattern
type Behavior interface {
	// ID returns unique behavior identifier
	ID() string

	// Name returns display name
	Name() string

	// Type returns behavior type
	Type() BehaviorType

	// Activate activates behavior
	Activate(ctx context.Context, participant Participant, encounter Encounter) error

	// Deactivate deactivates behavior
	Deactivate(ctx context.Context, participant Participant, encounter Encounter) error

	// Update updates behavior state
	Update(ctx context.Context, participant Participant, encounter Encounter, deltaMs int64) error

	// IsActive returns true if behavior is active
	IsActive() bool

	// Priority returns behavior priority
	Priority() int

	// CanActivate checks if behavior can be activated
	CanActivate(ctx context.Context, participant Participant, encounter Encounter) bool

	// ShouldDeactivate checks if behavior should end
	ShouldDeactivate(ctx context.Context, participant Participant, encounter Encounter) bool
}

// BehaviorType categorizes behaviors
type BehaviorType string

const (
	BehaviorIdle    BehaviorType = "idle"
	BehaviorPatrol  BehaviorType = "patrol"
	BehaviorChase   BehaviorType = "chase"
	BehaviorAttack  BehaviorType = "attack"
	BehaviorDefend  BehaviorType = "defend"
	BehaviorFlee    BehaviorType = "flee"
	BehaviorHeal    BehaviorType = "heal"
	BehaviorSupport BehaviorType = "support"
	BehaviorAmbush  BehaviorType = "ambush"
	BehaviorKite    BehaviorType = "kite"
	BehaviorGuard   BehaviorType = "guard"
	BehaviorBerserk BehaviorType = "berserk"
)

// AIPreset provides pre-configured AI behavior
type AIPreset interface {
	// ID returns unique preset identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Difficulty returns AI difficulty level
	Difficulty() int

	// CreateAI instantiates AI with preset configuration
	CreateAI(ctx context.Context) (AI, error)

	// Configure applies preset to existing AI
	Configure(ai AI) error

	// Strategy returns preset strategy
	Strategy() AIStrategy

	// TargetPreferences returns preset target preferences
	TargetPreferences() []TargetPreference

	// Behaviors returns preset behaviors
	Behaviors() []Behavior

	// Attributes returns preset attributes
	Attributes() map[string]interface{}
}

// AIPresetRegistry manages AI presets
type AIPresetRegistry interface {
	// Register adds preset to registry
	Register(preset AIPreset) error

	// Unregister removes preset from registry
	Unregister(presetID string) error

	// Get retrieves preset by ID
	Get(presetID string) (AIPreset, bool)

	// GetAll returns all registered presets
	GetAll() []AIPreset

	// GetByDifficulty returns presets at difficulty level
	GetByDifficulty(difficulty int) []AIPreset

	// GetByType returns presets of specific type
	GetByType(strategyType StrategyType) []AIPreset

	// Has checks if preset is registered
	Has(presetID string) bool

	// Count returns total registered presets
	Count() int
}

// AIFactory creates AI instances
type AIFactory interface {
	// Create creates new AI with default settings
	Create() AI

	// CreateFromPreset creates AI from preset
	CreateFromPreset(presetID string) (AI, error)

	// CreateWithStrategy creates AI with specific strategy
	CreateWithStrategy(strategy AIStrategy) AI

	// CreateWithConfig creates AI with custom configuration
	CreateWithConfig(config AIConfig) (AI, error)
}

// AIConfig defines AI configuration
type AIConfig struct {
	Strategy          StrategyType
	Difficulty        int
	TargetPreferences []PreferenceType
	BehaviorTypes     []BehaviorType
	AggressionLevel   float64
	CautiousLevel     float64
	TeamplayLevel     float64
	UseSkills         bool
	UseItems          bool
	UseTactics        bool
	ReactionTime      int64
}

// Coordinator manages AI team coordination
type Coordinator interface {
	// RegisterAI adds AI to coordination
	RegisterAI(participantID string, ai AI)

	// UnregisterAI removes AI from coordination
	UnregisterAI(participantID string)

	// GetAI retrieves AI by participant ID
	GetAI(participantID string) (AI, bool)

	// CoordinateActions coordinates team actions
	CoordinateActions(ctx context.Context, encounter Encounter) error

	// SetFocusTarget designates priority target for team
	SetFocusTarget(targetID string)

	// GetFocusTarget returns current focus target
	GetFocusTarget() string

	// SuggestAction suggests action for participant
	SuggestAction(ctx context.Context, participantID string, encounter Encounter) (Action, error)

	// BroadcastThreat shares threat information
	BroadcastThreat(ctx context.Context, threatInfo ThreatInfo)

	// RequestAssistance requests help from allies
	RequestAssistance(ctx context.Context, participantID string, encounter Encounter) error

	// FormationMode returns current formation mode
	FormationMode() FormationMode

	// SetFormationMode updates formation mode
	SetFormationMode(mode FormationMode)
}

// ThreatInfo contains shared threat information
type ThreatInfo struct {
	SourceID    string
	TargetID    string
	ThreatLevel ThreatLevel
	Priority    int
	Timestamp   int64
}

// FormationMode defines team positioning
type FormationMode string

const (
	FormationScattered  FormationMode = "scattered"
	FormationLine       FormationMode = "line"
	FormationCircle     FormationMode = "circle"
	FormationWedge      FormationMode = "wedge"
	FormationDefensive  FormationMode = "defensive"
	FormationAggressive FormationMode = "aggressive"
)

// AIMemory stores AI decision history
type AIMemory interface {
	// Remember stores information
	Remember(key string, value interface{})

	// Recall retrieves stored information
	Recall(key string) (interface{}, bool)

	// Forget removes stored information
	Forget(key string)

	// Has checks if information exists
	Has(key string) bool

	// Clear removes all stored information
	Clear()

	// GetRecentActions returns recent action history
	GetRecentActions(count int) []Action

	// GetRecentTargets returns recent target history
	GetRecentTargets(count int) []string

	// RecordAction records action to history
	RecordAction(action Action)

	// RecordTarget records target to history
	RecordTarget(targetID string)

	// Size returns memory size
	Size() int
}
