package ai

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/entity"
)

// AI controls enemy behavior
type AI interface {
	// Update processes AI logic for frame
	Update(ctx context.Context, deltaMs int64) error

	// SetOwner assigns entity this AI controls
	SetOwner(owner entity.Combatant)

	// Owner returns controlled entity
	Owner() entity.Combatant

	// CurrentBehavior returns active behavior
	CurrentBehavior() Behavior

	// SetBehavior changes active behavior
	SetBehavior(behavior Behavior)

	// Target returns current target entity
	Target() entity.Entity

	// SetTarget assigns target entity
	SetTarget(target entity.Entity)

	// HasTarget returns true if AI has valid target
	HasTarget() bool

	// ClearTarget removes current target
	ClearTarget()

	// IsAggro returns true if AI is in combat mode
	IsAggro() bool

	// SetAggro updates aggression state
	SetAggro(aggro bool)

	// ThreatTable returns threat manager
	ThreatTable() ThreatTable

	// Blackboard returns shared AI memory
	Blackboard() Blackboard

	// IsEnabled returns true if AI is active
	IsEnabled() bool

	// SetEnabled activates or deactivates AI
	SetEnabled(enabled bool)
}

// Behavior defines AI decision making
type Behavior interface {
	// ID returns unique behavior identifier
	ID() string

	// Name returns display name
	Name() string

	// Initialize prepares behavior for use
	Initialize(ai AI) error

	// Update executes behavior logic
	Update(ctx context.Context, deltaMs int64) (BehaviorStatus, error)

	// OnEnter is called when behavior becomes active
	OnEnter(ctx context.Context) error

	// OnExit is called when behavior ends
	OnExit(ctx context.Context) error

	// Priority returns behavior selection priority
	Priority() int

	// CanActivate returns true if behavior can run
	CanActivate() bool
}

// BehaviorStatus indicates behavior execution state
type BehaviorStatus int

const (
	StatusRunning BehaviorStatus = iota
	StatusSuccess
	StatusFailure
)

// BehaviorTree organizes behaviors hierarchically
type BehaviorTree interface {
	// Root returns root node of tree
	Root() BehaviorNode

	// SetRoot assigns root node
	SetRoot(node BehaviorNode)

	// Update executes tree logic
	Update(ctx context.Context, deltaMs int64) (BehaviorStatus, error)

	// Reset resets all nodes in tree
	Reset()
}

// BehaviorNode represents node in behavior tree
type BehaviorNode interface {
	// Execute runs node logic
	Execute(ctx context.Context, ai AI) (BehaviorStatus, error)

	// Reset resets node state
	Reset()

	// Children returns child nodes
	Children() []BehaviorNode

	// AddChild adds child node
	AddChild(child BehaviorNode)
}

// ThreatTable tracks aggro from multiple targets
type ThreatTable interface {
	// AddThreat increases threat from entity
	AddThreat(entityID string, amount float64)

	// RemoveThreat decreases threat from entity
	RemoveThreat(entityID string, amount float64)

	// SetThreat sets absolute threat value
	SetThreat(entityID string, amount float64)

	// GetThreat returns threat value for entity
	GetThreat(entityID string) float64

	// GetHighestThreat returns entity with most threat
	GetHighestThreat() (entityID string, threat float64)

	// GetAll returns all threat entries
	GetAll() map[string]float64

	// Clear removes all threat entries
	Clear()

	// Remove removes specific entity from table
	Remove(entityID string)

	// Decay reduces all threat by percentage
	Decay(percentage float64)

	// TransferThreat moves threat from one entity to another
	TransferThreat(fromID, toID string, amount float64)
}

// Blackboard stores shared AI data
type Blackboard interface {
	// Set stores value with key
	Set(key string, value interface{})

	// Get retrieves value by key
	Get(key string) (interface{}, bool)

	// GetString retrieves string value
	GetString(key string) (string, bool)

	// GetInt retrieves int value
	GetInt(key string) (int, bool)

	// GetFloat retrieves float64 value
	GetFloat(key string) (float64, bool)

	// GetBool retrieves bool value
	GetBool(key string) (bool, bool)

	// Has checks if key exists
	Has(key string) bool

	// Remove removes key from blackboard
	Remove(key string)

	// Clear removes all entries
	Clear()

	// Keys returns all stored keys
	Keys() []string
}

// StateMachine manages AI state transitions
type StateMachine interface {
	// CurrentState returns active state
	CurrentState() State

	// ChangeState transitions to new state
	ChangeState(ctx context.Context, stateID string) error

	// AddState registers state
	AddState(state State)

	// RemoveState unregisters state
	RemoveState(stateID string)

	// GetState retrieves state by ID
	GetState(stateID string) (State, bool)

	// Update processes current state
	Update(ctx context.Context, deltaMs int64) error

	// CanTransition checks if transition is allowed
	CanTransition(fromID, toID string) bool

	// AddTransition defines allowed state change
	AddTransition(fromID, toID string, condition TransitionCondition)
}

// State represents AI state
type State interface {
	// ID returns unique state identifier
	ID() string

	// Name returns display name
	Name() string

	// OnEnter is called when entering state
	OnEnter(ctx context.Context, ai AI) error

	// OnUpdate is called each frame while in state
	OnUpdate(ctx context.Context, ai AI, deltaMs int64) error

	// OnExit is called when leaving state
	OnExit(ctx context.Context, ai AI) error

	// Transitions returns possible state transitions
	Transitions() []string
}

// TransitionCondition determines if state change is allowed
type TransitionCondition interface {
	// Check returns true if transition should occur
	Check(ctx context.Context, ai AI) bool

	// Description returns human-readable condition
	Description() string
}

// Perception handles AI sensory input
type Perception interface {
	// CanSee checks if target is visible
	CanSee(targetID string) bool

	// CanHear checks if target is audible
	CanHear(targetID string) bool

	// GetVisibleEntities returns all visible entities
	GetVisibleEntities() []string

	// GetAudibleEntities returns all audible entities
	GetAudibleEntities() []string

	// SightRange returns vision distance
	SightRange() float64

	// SetSightRange updates vision distance
	SetSightRange(range_ float64)

	// HearingRange returns hearing distance
	HearingRange() float64

	// SetHearingRange updates hearing distance
	SetHearingRange(range_ float64)

	// FieldOfView returns vision angle in degrees
	FieldOfView() float64

	// SetFieldOfView updates vision angle
	SetFieldOfView(degrees float64)

	// Update refreshes perception data
	Update(ctx context.Context) error
}

// MovementController handles AI movement
type MovementController interface {
	// MoveTo moves towards target position
	MoveTo(ctx context.Context, x, y float64) error

	// MoveTowards moves towards target entity
	MoveTowards(ctx context.Context, targetID string) error

	// Stop halts movement
	Stop()

	// IsMoving returns true if currently moving
	IsMoving() bool

	// CurrentPath returns movement path
	CurrentPath() []PathNode

	// Speed returns movement speed
	Speed() float64

	// SetSpeed updates movement speed
	SetSpeed(speed float64)

	// Destination returns target position
	Destination() (x, y float64, hasDestination bool)

	// DistanceToDestination returns remaining distance
	DistanceToDestination() float64
}

// PathNode represents point in movement path
type PathNode struct {
	X float64
	Y float64
}

// CombatController handles AI combat actions
type CombatController interface {
	// AttackTarget performs attack on target
	AttackTarget(ctx context.Context, targetID string) error

	// UseSkill activates skill on target
	UseSkill(ctx context.Context, skillID string, targetID string) error

	// Flee moves away from threat
	Flee(ctx context.Context, threatID string) error

	// Pursue chases target
	Pursue(ctx context.Context, targetID string) error

	// Defend raises defenses
	Defend(ctx context.Context) error

	// IsInCombat returns true if engaged in combat
	IsInCombat() bool

	// GetPreferredRange returns optimal combat distance
	GetPreferredRange() float64

	// SetPreferredRange updates combat distance
	SetPreferredRange(range_ float64)
}

// DecisionMaker selects actions based on utility
type DecisionMaker interface {
	// EvaluateActions scores all available actions
	EvaluateActions(ctx context.Context, ai AI) []ScoredAction

	// SelectBestAction chooses highest scoring action
	SelectBestAction(ctx context.Context, ai AI) (Action, error)

	// AddAction registers action for consideration
	AddAction(action Action)

	// RemoveAction unregisters action
	RemoveAction(actionID string)

	// GetActions returns all registered actions
	GetActions() []Action
}

// Action represents AI action
type Action interface {
	// ID returns unique action identifier
	ID() string

	// Name returns display name
	Name() string

	// CalculateUtility scores action desirability
	CalculateUtility(ctx context.Context, ai AI) float64

	// CanExecute checks if action is possible
	CanExecute(ctx context.Context, ai AI) bool

	// Execute performs action
	Execute(ctx context.Context, ai AI) error

	// Cooldown returns remaining cooldown
	Cooldown() int64

	// SetCooldown updates cooldown
	SetCooldown(ms int64)
}

// ScoredAction pairs action with utility score
type ScoredAction struct {
	Action  Action
	Utility float64
}

// Squad coordinates multiple AI entities
type Squad interface {
	// ID returns unique squad identifier
	ID() string

	// AddMember adds AI to squad
	AddMember(ai AI)

	// RemoveMember removes AI from squad
	RemoveMember(aiID string)

	// GetMembers returns all squad members
	GetMembers() []AI

	// GetLeader returns squad leader
	GetLeader() AI

	// SetLeader assigns squad leader
	SetLeader(aiID string)

	// Formation returns squad formation
	Formation() SquadFormation

	// SetFormation updates squad formation
	SetFormation(formation SquadFormation)

	// SharedBlackboard returns squad-wide data
	SharedBlackboard() Blackboard

	// Update processes squad logic
	Update(ctx context.Context, deltaMs int64) error
}

// SquadFormation defines squad positioning
type SquadFormation interface {
	// Type returns formation type
	Type() FormationType

	// GetPosition calculates position for member
	GetPosition(memberIndex int, leaderX, leaderY float64) (x, y float64)

	// Spacing returns distance between members
	Spacing() float64

	// SetSpacing updates member distance
	SetSpacing(spacing float64)
}

// FormationType categorizes squad formations
type FormationType string

const (
	FormationLine    FormationType = "line"
	FormationColumn  FormationType = "column"
	FormationWedge   FormationType = "wedge"
	FormationCircle  FormationType = "circle"
	FormationScatter FormationType = "scatter"
)

// Preset provides pre-configured AI behavior
type Preset interface {
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

	// BehaviorTree returns preset behavior tree
	BehaviorTree() BehaviorTree

	// Attributes returns preset attributes
	Attributes() map[string]any
}

// PresetRegistry manages AI presets
type PresetRegistry interface {
	// Register adds preset to registry
	Register(preset Preset) error

	// Unregister removes preset from registry
	Unregister(presetID string) error

	// Get retrieves preset by ID
	Get(presetID string) (Preset, bool)

	// GetAll returns all registered presets
	GetAll() []Preset

	// GetByDifficulty returns presets at difficulty level
	GetByDifficulty(difficulty int) []Preset

	// Has checks if preset is registered
	Has(presetID string) bool
}
