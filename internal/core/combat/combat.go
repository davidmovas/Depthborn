package combat

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

// Encounter represents single combat instance
type Encounter interface {
	// ID returns unique encounter identifier
	ID() string

	// State returns current encounter state
	State() EncounterState

	// SetState updates encounter state
	SetState(state EncounterState)

	// Arena returns combat arena
	Arena() Arena

	// TurnOrder returns turn order manager
	TurnOrder() TurnOrder

	// Participants returns all combatants
	Participants() []Participant

	// AddParticipant adds combatant to encounter
	AddParticipant(participant Participant) error

	// RemoveParticipant removes combatant from encounter
	RemoveParticipant(participantID string) error

	// GetParticipant retrieves participant by entity ID
	GetParticipant(entityID string) (Participant, bool)

	// PlayerParty returns player-controlled participants
	PlayerParty() []Participant

	// EnemyParty returns enemy participants
	EnemyParty() []Participant

	// Start begins the encounter
	Start(ctx context.Context) error

	// End finishes the encounter
	End(ctx context.Context, result EncounterResult) error

	// ProcessTurn executes single turn
	ProcessTurn(ctx context.Context) error

	// CurrentTurn returns active participant
	CurrentTurn() (Participant, bool)

	// NextTurn advances to next participant
	NextTurn() (Participant, error)

	// CanAct checks if participant can act this turn
	CanAct(participantID string) bool

	// PerformAction executes combat action
	PerformAction(ctx context.Context, action Action) (ActionResult, error)

	// VictoryConditions returns win conditions
	VictoryConditions() []Condition

	// DefeatConditions returns loss conditions
	DefeatConditions() []Condition

	// AddVictoryCondition adds win condition
	AddVictoryCondition(condition Condition)

	// AddDefeatCondition adds loss condition
	AddDefeatCondition(condition Condition)

	// CheckVictory evaluates if victory conditions met
	CheckVictory(ctx context.Context) (bool, string)

	// CheckDefeat evaluates if defeat conditions met
	CheckDefeat(ctx context.Context) (bool, string)

	// RoundNumber returns current round number
	RoundNumber() int

	// OnTurnStart registers callback when turn begins
	OnTurnStart(callback TurnCallback)

	// OnTurnEnd registers callback when turn ends
	OnTurnEnd(callback TurnCallback)

	// OnEncounterEnd registers callback when encounter ends
	OnEncounterEnd(callback EncounterCallback)
}

// EncounterState represents combat phase
type EncounterState string

const (
	StateSetup      EncounterState = "setup"       // Initial positioning
	StateRollInit   EncounterState = "roll_init"   // Determine turn order
	StateInProgress EncounterState = "in_progress" // Active combat
	StatePaused     EncounterState = "paused"      // Temporarily halted
	StateVictory    EncounterState = "victory"     // Player won
	StateDefeat     EncounterState = "defeat"      // Player lost
	StateEnded      EncounterState = "ended"       // Combat finished
)

// EncounterResult describes combat outcome
type EncounterResult struct {
	Victory          bool
	DefeatReason     string
	VictoryReason    string
	TurnsElapsed     int
	RoundsElapsed    int
	Survivors        []string
	Casualties       []string
	ExperienceGained int64
	LootGenerated    []any
	Statistics       EncounterStatistics
}

// EncounterStatistics tracks combat metrics
type EncounterStatistics struct {
	TotalDamageDealt map[string]float64
	TotalDamageTaken map[string]float64
	TotalHealing     map[string]float64
	CriticalHits     map[string]int
	SkillsUsed       map[string]int
	StatusesApplied  map[string]int
	DeathCount       map[string]int
	ActionsPerformed map[string]int
}

// TurnCallback is invoked for turn events
type TurnCallback func(ctx context.Context, encounter Encounter, participant Participant)

// EncounterCallback is invoked for encounter events
type EncounterCallback func(ctx context.Context, encounter Encounter, result EncounterResult)

// Participant represents combatant in encounter
type Participant interface {
	// EntityID returns underlying entity identifier
	EntityID() string

	// Entity returns combatant entity
	Entity() entity.Combatant

	// Team returns participant team
	Team() Team

	// SetTeam updates participant team
	SetTeam(team Team)

	// Position returns arena position
	Position() spatial.Position

	// SetPosition updates arena position
	SetPosition(pos spatial.Position)

	// Transform returns spatial transform
	Transform() spatial.Transform

	// Initiative returns initiative value
	Initiative() int

	// SetInitiative updates initiative
	SetInitiative(value int)

	// HasActed returns true if acted this round
	HasActed() bool

	// SetHasActed marks as acted or not
	SetHasActed(acted bool)

	// IsDefeated returns true if participant is out of combat
	IsDefeated() bool

	// MarkDefeated marks participant as defeated
	MarkDefeated()

	// AvailableActions returns possible actions
	AvailableActions() []Action

	// CanPerformAction checks if action is possible
	CanPerformAction(action Action) bool

	// Modifiers returns active combat modifiers
	Modifiers() ModifierSet

	// Reactions returns available reactions
	Reactions() []Reaction

	// AddReaction adds reactive action
	AddReaction(reaction Reaction)

	// RemoveReaction removes reactive action
	RemoveReaction(reactionID string)
}

// Team categorizes participants
type Team string

const (
	TeamPlayer  Team = "player"
	TeamEnemy   Team = "enemy"
	TeamNeutral Team = "neutral"
	TeamAlly    Team = "ally"
)

// Action represents combat action
type Action interface {
	// ID returns unique action identifier
	ID() string

	// Name returns display name
	Name() string

	// Type returns action type
	Type() ActionType

	// ActorID returns ID of participant performing action
	ActorID() string

	// SetActor assigns participant performing action
	SetActor(participantID string)

	// TargetIDs returns action target IDs
	TargetIDs() []string

	// SetTargets updates action targets
	SetTargets(targetIDs []string)

	// TargetingRule returns target selection rule
	TargetingRule() TargetingRule

	// Validate checks if action is legal
	Validate(ctx context.Context, encounter Encounter) error

	// Execute performs the action
	Execute(ctx context.Context, encounter Encounter) (ActionResult, error)

	// Cost returns action cost
	Cost() ActionCost

	// Range returns action range
	Range() float64

	// AreaOfEffect returns AoE area (nil if single target)
	AreaOfEffect() spatial.Area

	// RequiresLineOfSight returns true if needs clear path
	RequiresLineOfSight() bool

	// CanBeInterrupted returns true if action can be stopped
	CanBeInterrupted() bool

	// Priority returns execution priority (higher = first)
	Priority() int

	// Description returns human-readable description
	Description() string
}

// ActionType categorizes actions
type ActionType string

const (
	ActionAttack   ActionType = "attack"
	ActionSkill    ActionType = "skill"
	ActionMove     ActionType = "move"
	ActionDefend   ActionType = "defend"
	ActionItem     ActionType = "item"
	ActionWait     ActionType = "wait"
	ActionFlee     ActionType = "flee"
	ActionInteract ActionType = "interact"
	ActionReaction ActionType = "reaction"
	ActionCombo    ActionType = "combo"
)

// ActionCost defines resource requirements
type ActionCost struct {
	ActionPoints int
	Mana         float64
	Health       float64
	Stamina      float64
	Items        map[string]int // itemID -> quantity
}

// ActionResult describes action outcome
type ActionResult struct {
	Success            bool
	Message            string
	DamageDealt        map[string]float64          // targetID -> damage
	HealingDone        map[string]float64          // targetID -> healing
	StatusApplied      map[string][]string         // targetID -> statusIDs
	Moved              map[string]spatial.Position // participantID -> new position
	TriggeredReactions []Reaction
	SecondaryEffects   []Effect
	Flags              []ResultFlag
}

// ResultFlag describes special action result
type ResultFlag string

const (
	FlagCritical    ResultFlag = "critical"
	FlagBlocked     ResultFlag = "blocked"
	FlagEvaded      ResultFlag = "evaded"
	FlagCountered   ResultFlag = "countered"
	FlagKilled      ResultFlag = "killed"
	FlagInterrupted ResultFlag = "interrupted"
)

// TargetingRule defines target selection
type TargetingRule interface {
	// Type returns targeting type
	Type() TargetingType

	// SelectTargets chooses valid targets
	SelectTargets(ctx context.Context, actorID string, encounter Encounter) ([]string, error)

	// IsValidTarget checks if target is legal
	IsValidTarget(actorID, targetID string, encounter Encounter) bool

	// MaxTargets returns maximum selectable targets (0 = unlimited)
	MaxTargets() int

	// MinTargets returns minimum targets required
	MinTargets() int

	// RequiresTarget returns true if must have target
	RequiresTarget() bool

	// AllowsSelf returns true if can target self
	AllowsSelf() bool

	// AllowsAllies returns true if can target allies
	AllowsAllies() bool

	// AllowsEnemies returns true if can target enemies
	AllowsEnemies() bool
}

// TargetingType categorizes targeting
type TargetingType string

const (
	TargetSelf       TargetingType = "self"
	TargetSingle     TargetingType = "single"
	TargetMultiple   TargetingType = "multiple"
	TargetArea       TargetingType = "area"
	TargetLine       TargetingType = "line"
	TargetCone       TargetingType = "cone"
	TargetAllAllies  TargetingType = "all_allies"
	TargetAllEnemies TargetingType = "all_enemies"
	TargetAll        TargetingType = "all"
	TargetGround     TargetingType = "ground"
)

// =============================================================================
// REACTION - represents reactive action
// =============================================================================

// Reaction represents reactive action
type Reaction interface {
	// ID returns unique reaction identifier
	ID() string

	// Name returns display name
	Name() string

	// OwnerID returns participant who owns reaction
	OwnerID() string

	// SetOwner assigns reaction owner
	SetOwner(participantID string)

	// Trigger returns what triggers reaction
	Trigger() ReactionTrigger

	// CanTrigger checks if reaction can activate
	CanTrigger(ctx context.Context, encounter Encounter) bool

	// Execute performs reaction
	Execute(ctx context.Context, encounter Encounter) (ActionResult, error)

	// Priority returns reaction priority (higher = first)
	Priority() int

	// Cost returns reaction cost
	Cost() ActionCost

	// UsesRemaining returns number of uses left (-1 = unlimited)
	UsesRemaining() int

	// DecrementUses reduces uses by one
	DecrementUses()

	// IsExpended returns true if no uses remain
	IsExpended() bool

	// Description returns human-readable description
	Description() string
}

// ReactionTrigger defines reaction condition
type ReactionTrigger interface {
	// Type returns trigger type
	Type() TriggerType

	// SourceID returns triggering participant ID (can be empty)
	SourceID() string

	// TargetID returns affected participant ID (can be empty)
	TargetID() string

	// Check evaluates if trigger activated
	Check(ctx context.Context, encounter Encounter) bool

	// Description returns human-readable trigger
	Description() string
}

// TriggerType categorizes triggers
type TriggerType string

const (
	TriggerOnAttacked      TriggerType = "on_attacked"
	TriggerOnHit           TriggerType = "on_hit"
	TriggerOnMiss          TriggerType = "on_miss"
	TriggerOnDamaged       TriggerType = "on_damaged"
	TriggerOnHealed        TriggerType = "on_healed"
	TriggerOnAllyAttacked  TriggerType = "on_ally_attacked"
	TriggerOnAllyDamaged   TriggerType = "on_ally_damaged"
	TriggerOnAllyDeath     TriggerType = "on_ally_death"
	TriggerOnEnemyDeath    TriggerType = "on_enemy_death"
	TriggerOnStatusApplied TriggerType = "on_status_applied"
	TriggerOnStatusRemoved TriggerType = "on_status_removed"
	TriggerOnTurnStart     TriggerType = "on_turn_start"
	TriggerOnTurnEnd       TriggerType = "on_turn_end"
	TriggerOnMove          TriggerType = "on_move"
	TriggerOnSkillCast     TriggerType = "on_skill_cast"
)

// ModifierSet manages combat modifiers
type ModifierSet interface {
	// Add adds modifier
	Add(modifier CombatModifier)

	// Remove removes modifier by ID
	Remove(modifierID string)

	// Get retrieves modifier by ID
	Get(modifierID string) (CombatModifier, bool)

	// GetAll returns all modifiers
	GetAll() []CombatModifier

	// GetByType returns modifiers of type
	GetByType(modType ModifierType) []CombatModifier

	// Clear removes all modifiers
	Clear()

	// Apply applies all modifiers to value
	Apply(baseValue float64, modType ModifierType) float64

	// Update processes modifiers for elapsed time
	Update(ctx context.Context, deltaMs int64) error
}

// CombatModifier alters combat calculations
type CombatModifier interface {
	// ID returns unique modifier identifier
	ID() string

	// Name returns display name
	Name() string

	// Type returns modifier type
	Type() ModifierType

	// Value returns modifier value
	Value() float64

	// Source returns modifier source
	Source() string

	// Duration returns remaining duration in milliseconds (-1 = permanent)
	Duration() int64

	// SetDuration updates duration
	SetDuration(ms int64)

	// IsExpired returns true if duration ended
	IsExpired() bool

	// Apply applies modifier to value
	Apply(baseValue float64) float64

	// Priority returns application order (higher = first)
	Priority() int

	// Stacks returns if multiple instances stack
	Stacks() bool
}

// ModifierType categorizes modifiers
type ModifierType string

const (
	ModDamageDealt     ModifierType = "damage_dealt"
	ModDamageTaken     ModifierType = "damage_taken"
	ModHealingDone     ModifierType = "healing_done"
	ModHealingReceived ModifierType = "healing_received"
	ModCritChance      ModifierType = "crit_chance"
	ModCritDamage      ModifierType = "crit_damage"
	ModAccuracy        ModifierType = "accuracy"
	ModEvasion         ModifierType = "evasion"
	ModSpeed           ModifierType = "speed"
	ModInitiative      ModifierType = "initiative"
	ModActionPoints    ModifierType = "action_points"
	ModArmor           ModifierType = "armor"
	ModResistance      ModifierType = "resistance"
)

// TurnOrder manages action sequence
type TurnOrder interface {
	// Calculate determines turn order from participants
	Calculate(ctx context.Context, participants []Participant) error

	// Next advances to next participant
	Next() (Participant, bool)

	// Current returns active participant
	Current() (Participant, bool)

	// Peek returns next participant without advancing
	Peek() (Participant, bool)

	// GetOrder returns full turn sequence
	GetOrder() []Participant

	// Insert adds participant to order at position
	Insert(participant Participant, position int)

	// Remove removes participant from order
	Remove(participantID string)

	// Delay moves participant later in order
	Delay(participantID string, positions int)

	// Advance moves participant earlier in order
	Advance(participantID string, positions int)

	// Reset recalculates turn order
	Reset(ctx context.Context, participants []Participant)

	// RoundNumber returns current round
	RoundNumber() int

	// IncrementRound advances to next round
	IncrementRound()

	// IsNewRound returns true if starting new round
	IsNewRound() bool

	// TurnNumber returns current turn in round
	TurnNumber() int
}

// Condition defines combat condition
type Condition interface {
	// ID returns unique condition identifier
	ID() string

	// Description returns human-readable condition
	Description() string

	// Check evaluates if condition is met
	Check(ctx context.Context, encounter Encounter) bool

	// Type returns condition type
	Type() ConditionType

	// IsVictory returns true if this is win condition
	IsVictory() bool

	// IsDefeat returns true if this is loss condition
	IsDefeat() bool
}

// ConditionType categorizes conditions
type ConditionType string

const (
	ConditionEliminateAll    ConditionType = "eliminate_all"
	ConditionEliminateBoss   ConditionType = "eliminate_boss"
	ConditionEliminateTarget ConditionType = "eliminate_target"
	ConditionSurvive         ConditionType = "survive"
	ConditionSurviveRounds   ConditionType = "survive_rounds"
	ConditionProtect         ConditionType = "protect"
	ConditionReachLocation   ConditionType = "reach_location"
	ConditionCollectItems    ConditionType = "collect_items"
	ConditionTimeLimit       ConditionType = "time_limit"
	ConditionAllAlliesDead   ConditionType = "all_allies_dead"
	ConditionTargetDead      ConditionType = "target_dead"
	ConditionHealthThreshold ConditionType = "health_threshold"
	ConditionCustom          ConditionType = "custom"
)

// Effect represents combat effect
type Effect interface {
	// ID returns unique effect identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns effect description
	Description() string

	// Type returns effect type
	Type() EffectType

	// Apply applies effect to encounter
	Apply(ctx context.Context, encounter Encounter) error

	// TargetIDs returns affected participant IDs
	TargetIDs() []string

	// Duration returns effect duration in milliseconds (0 = instant)
	Duration() int64
}

// EffectType categorizes effects
type EffectType string

const (
	EffectDamage        EffectType = "damage"
	EffectHealing       EffectType = "healing"
	EffectStatus        EffectType = "status"
	EffectBuff          EffectType = "buff"
	EffectDebuff        EffectType = "debuff"
	EffectTeleport      EffectType = "teleport"
	EffectKnockback     EffectType = "knockback"
	EffectSummon        EffectType = "summon"
	EffectEnvironmental EffectType = "environmental"
)
