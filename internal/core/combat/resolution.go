package combat

import (
	"context"
)

// Phase represents combat resolution stage
type Phase interface {
	// ID returns unique phase identifier
	ID() string

	// Name returns display name
	Name() string

	// Type returns phase type
	Type() PhaseType

	// Enter is called when phase begins
	Enter(ctx context.Context, encounter Encounter) error

	// Process executes phase logic
	Process(ctx context.Context, encounter Encounter) error

	// Exit is called when phase ends
	Exit(ctx context.Context, encounter Encounter) error

	// CanSkip returns true if phase can be bypassed
	CanSkip(ctx context.Context, encounter Encounter) bool

	// Next returns next phase ID
	Next() string

	// IsTerminal returns true if phase ends encounter
	IsTerminal() bool

	// AllowsPlayerInput returns true if waiting for player action
	AllowsPlayerInput() bool

	// Timeout returns maximum phase duration in milliseconds (0 = no timeout)
	Timeout() int64
}

// PhaseType categorizes phases
type PhaseType string

const (
	PhaseSetup         PhaseType = "setup"
	PhaseInitiative    PhaseType = "initiative"
	PhaseActionSelect  PhaseType = "action_select"
	PhaseActionResolve PhaseType = "action_resolve"
	PhaseReaction      PhaseType = "reaction"
	PhaseStatusUpdate  PhaseType = "status_update"
	PhaseEndTurn       PhaseType = "end_turn"
	PhaseEndRound      PhaseType = "end_round"
	PhaseCleanup       PhaseType = "cleanup"
	PhaseVictory       PhaseType = "victory"
	PhaseDefeat        PhaseType = "defeat"
)

// PhaseManager coordinates combat flow
type PhaseManager interface {
	// CurrentPhase returns active phase
	CurrentPhase() Phase

	// AdvancePhase moves to next phase
	AdvancePhase(ctx context.Context, encounter Encounter) error

	// SetPhase changes to specific phase
	SetPhase(ctx context.Context, phaseID string, encounter Encounter) error

	// ProcessPhase executes current phase
	ProcessPhase(ctx context.Context, encounter Encounter) error

	// IsPhaseComplete checks if phase is done
	IsPhaseComplete(ctx context.Context, phase Phase, encounter Encounter) bool

	// RegisterPhase adds phase to manager
	RegisterPhase(phase Phase) error

	// UnregisterPhase removes phase from manager
	UnregisterPhase(phaseID string) error

	// GetPhase retrieves phase by ID
	GetPhase(phaseID string) (Phase, bool)

	// GetAllPhases returns all registered phases
	GetAllPhases() []Phase

	// Reset resets to initial phase
	Reset(ctx context.Context, encounter Encounter) error

	// OnPhaseEnter registers callback when phase begins
	OnPhaseEnter(callback PhaseCallback)

	// OnPhaseExit registers callback when phase ends
	OnPhaseExit(callback PhaseCallback)
}

// PhaseCallback is invoked for phase events
type PhaseCallback func(ctx context.Context, phase Phase, encounter Encounter)

// Engine drives combat loop
type Engine interface {
	// Start begins combat encounter
	Start(ctx context.Context, encounter Encounter) error

	// Update processes combat for frame
	Update(ctx context.Context, encounter Encounter, deltaMs int64) error

	// Pause suspends combat
	Pause()

	// Resume continues combat
	Resume()

	// Stop ends combat
	Stop(ctx context.Context, encounter Encounter) error

	// IsPaused returns true if paused
	IsPaused() bool

	// IsRunning returns true if active
	IsRunning() bool

	// State returns current engine state
	State() EngineState

	// SetUpdateRate sets updates per second
	SetUpdateRate(updatesPerSecond int)

	// GetUpdateRate returns current update rate
	GetUpdateRate() int

	// ElapsedTime returns total combat time in milliseconds
	ElapsedTime() int64

	// OnStateChange registers callback when engine state changes
	OnStateChange(callback EngineStateCallback)
}

// EngineState represents engine status
type EngineState string

const (
	EngineIdle    EngineState = "idle"
	EngineRunning EngineState = "running"
	EnginePaused  EngineState = "paused"
	EngineStopped EngineState = "stopped"
)

// EngineStateCallback is invoked when engine state changes
type EngineStateCallback func(oldState, newState EngineState)

// =============================================================================
// TURN PROCESSOR
// =============================================================================

// TurnProcessor handles turn execution
type TurnProcessor interface {
	// BeginTurn initializes turn for participant
	BeginTurn(ctx context.Context, participant Participant, encounter Encounter) error

	// ProcessTurn executes participant turn
	ProcessTurn(ctx context.Context, participant Participant, encounter Encounter) error

	// EndTurn finalizes turn for participant
	EndTurn(ctx context.Context, participant Participant, encounter Encounter) error

	// CanAct checks if participant can act
	CanAct(participant Participant, encounter Encounter) bool

	// GetAvailableActions returns possible actions
	GetAvailableActions(participant Participant, encounter Encounter) []Action

	// SelectAction chooses action for AI participants
	SelectAction(ctx context.Context, participant Participant, encounter Encounter) (Action, error)

	// ValidateTurn checks if turn is legal
	ValidateTurn(ctx context.Context, participant Participant, action Action, encounter Encounter) error

	// ApplyTurnCosts deducts action costs
	ApplyTurnCosts(ctx context.Context, participant Participant, action Action, encounter Encounter) error

	// OnTurnStart registers callback when turn begins
	OnTurnStart(callback TurnEventCallback)

	// OnTurnEnd registers callback when turn ends
	OnTurnEnd(callback TurnEventCallback)

	// OnActionPerformed registers callback when action is performed
	OnActionPerformed(callback ActionEventCallback)
}

// TurnEventCallback is invoked for turn events
type TurnEventCallback func(ctx context.Context, participant Participant, encounter Encounter)

// ActionEventCallback is invoked for action events
type ActionEventCallback func(ctx context.Context, participant Participant, action Action, result ActionResult, encounter Encounter)

// RoundManager handles combat rounds
type RoundManager interface {
	// CurrentRound returns round number
	CurrentRound() int

	// BeginRound starts new round
	BeginRound(ctx context.Context, encounter Encounter) error

	// EndRound finishes current round
	EndRound(ctx context.Context, encounter Encounter) error

	// IncrementRound advances round counter
	IncrementRound()

	// ResetRound sets round to zero
	ResetRound()

	// ProcessRoundStart handles round start effects
	ProcessRoundStart(ctx context.Context, encounter Encounter) error

	// ProcessRoundEnd handles round end effects
	ProcessRoundEnd(ctx context.Context, encounter Encounter) error

	// MaxRounds returns maximum rounds allowed (0 = unlimited)
	MaxRounds() int

	// SetMaxRounds updates round limit
	SetMaxRounds(max int)

	// IsMaxRoundsReached checks if round limit reached
	IsMaxRoundsReached() bool

	// OnRoundStart registers callback when round begins
	OnRoundStart(callback RoundCallback)

	// OnRoundEnd registers callback when round ends
	OnRoundEnd(callback RoundCallback)
}

// RoundCallback is invoked for round events
type RoundCallback func(ctx context.Context, encounter Encounter, roundNumber int)

// Timeline tracks combat events
type Timeline interface {
	// Record adds event to timeline
	Record(event TimelineEvent)

	// GetEvents returns all events
	GetEvents() []TimelineEvent

	// GetEventsByType returns events of specific type
	GetEventsByType(eventType EventType) []TimelineEvent

	// GetEventsByParticipant returns events involving participant
	GetEventsByParticipant(participantID string) []TimelineEvent

	// GetEventsByRound returns events from specific round
	GetEventsByRound(round int) []TimelineEvent

	// GetEventsByTurn returns events from specific turn
	GetEventsByTurn(round, turn int) []TimelineEvent

	// GetRecentEvents returns N most recent events
	GetRecentEvents(count int) []TimelineEvent

	// Clear removes all events
	Clear()

	// Export exports timeline data
	Export() TimelineData

	// Size returns number of recorded events
	Size() int
}

// TimelineEvent represents combat occurrence
type TimelineEvent interface {
	// ID returns unique event identifier
	ID() string

	// Type returns event type
	Type() EventType

	// Timestamp returns when event occurred (milliseconds since combat start)
	Timestamp() int64

	// Round returns round number
	Round() int

	// Turn returns turn number
	Turn() int

	// ParticipantIDs returns involved entity IDs
	ParticipantIDs() []string

	// Data returns event-specific data
	Data() map[string]interface{}

	// Description returns human-readable description
	Description() string

	// Severity returns event importance
	Severity() EventSeverity
}

// EventType categorizes events
type EventType string

const (
	EventCombatStart       EventType = "combat_start"
	EventCombatEnd         EventType = "combat_end"
	EventRoundStart        EventType = "round_start"
	EventRoundEnd          EventType = "round_end"
	EventTurnStart         EventType = "turn_start"
	EventTurnEnd           EventType = "turn_end"
	EventActionPerformed   EventType = "action_performed"
	EventActionFailed      EventType = "action_failed"
	EventDamageDealt       EventType = "damage_dealt"
	EventHealingDone       EventType = "healing_done"
	EventStatusApplied     EventType = "status_applied"
	EventStatusRemoved     EventType = "status_removed"
	EventEntityDefeated    EventType = "entity_defeated"
	EventEntityRevived     EventType = "entity_revived"
	EventPositionChanged   EventType = "position_changed"
	EventReactionTriggered EventType = "reaction_triggered"
	EventComboExecuted     EventType = "combo_executed"
	EventCriticalHit       EventType = "critical_hit"
	EventMissed            EventType = "missed"
	EventBlocked           EventType = "blocked"
	EventEvaded            EventType = "evaded"
	EventCountered         EventType = "countered"
)

// EventSeverity indicates event importance
type EventSeverity int

const (
	SeverityLow EventSeverity = iota
	SeverityNormal
	SeverityHigh
	SeverityCritical
)

// TimelineData contains exportable timeline
type TimelineData struct {
	StartTime        int64
	EndTime          int64
	TotalRounds      int
	TotalTurns       int
	Events           []TimelineEvent
	Statistics       Statistics
	ParticipantStats map[string]ParticipantStatistics
}

// Statistics contains aggregate combat data
type Statistics struct {
	TotalDamage     float64
	TotalHealing    float64
	TotalActions    int
	CriticalHits    int
	Misses          int
	StatusesApplied int
	Deaths          int
	Revivals        int
	AverageTurnTime int64
	LongestTurn     int64
	ShortestTurn    int64
}

// ParticipantStatistics contains per-participant data
type ParticipantStatistics struct {
	ParticipantID    string
	DamageDealt      float64
	DamageTaken      float64
	HealingDone      float64
	HealingReceived  float64
	ActionsPerformed int
	CriticalHits     int
	Misses           int
	Deaths           int
	KillCount        int
	TurnsPlayed      int
}

// StateManager manages encounter state
type StateManager interface {
	// SaveState captures current encounter state
	SaveState(ctx context.Context, encounter Encounter) (EncounterState, error)

	// RestoreState restores encounter to previous state
	RestoreState(ctx context.Context, encounter Encounter, state EncounterState) error

	// GetHistory returns state history
	GetHistory() []EncounterState

	// GetStateAtRound retrieves state at specific round
	GetStateAtRound(round int) (EncounterState, bool)

	// Undo reverts to previous state
	Undo(ctx context.Context, encounter Encounter) error

	// Redo advances to next state
	Redo(ctx context.Context, encounter Encounter) error

	// CanUndo returns true if undo is possible
	CanUndo() bool

	// CanRedo returns true if redo is possible
	CanRedo() bool

	// ClearHistory removes all saved states
	ClearHistory()

	// MaxHistorySize returns maximum saved states
	MaxHistorySize() int

	// SetMaxHistorySize updates history limit
	SetMaxHistorySize(size int)
}

// EncounterSnapshot represents snapshot of encounter (note: conflicts with earlier EncounterState)
// Renamed to avoid conflict
type EncounterSnapshot interface {
	// ID returns unique state identifier
	ID() string

	// Timestamp returns when state was captured
	Timestamp() int64

	// Round returns round number
	Round() int

	// Turn returns turn number
	Turn() int

	// Participants returns participant snapshots
	Participants() []ParticipantSnapshot

	// Arena returns arena snapshot
	Arena() ArenaSnapshot

	// Data returns additional state data
	Data() map[string]interface{}

	// PhaseID returns active phase ID
	PhaseID() string
}

// ParticipantSnapshot represents participant state
type ParticipantSnapshot struct {
	EntityID    string
	Position    map[string]int // x, y, z
	Health      float64
	MaxHealth   float64
	Mana        float64
	MaxMana     float64
	Stamina     float64
	StatusIDs   []string
	ModifierIDs []string
	HasActed    bool
	Initiative  int
	IsDefeated  bool
	Team        Team
}

// ArenaSnapshot represents arena state
type ArenaSnapshot struct {
	ActiveHazardIDs   []string
	InteractiveStates map[string]bool
	WeatherType       string
	AmbientEffectIDs  []string
}

// VictoryChecker evaluates victory conditions
type VictoryChecker interface {
	// CheckVictory checks if victory conditions are met
	CheckVictory(ctx context.Context, encounter Encounter) (bool, string)

	// CheckDefeat checks if defeat conditions are met
	CheckDefeat(ctx context.Context, encounter Encounter) (bool, string)

	// EvaluateCondition evaluates specific condition
	EvaluateCondition(ctx context.Context, condition Condition, encounter Encounter) bool

	// GetVictoryReason returns why player won
	GetVictoryReason(ctx context.Context, encounter Encounter) string

	// GetDefeatReason returns why player lost
	GetDefeatReason(ctx context.Context, encounter Encounter) string

	// AddVictoryCondition adds win condition
	AddVictoryCondition(condition Condition)

	// AddDefeatCondition adds loss condition
	AddDefeatCondition(condition Condition)

	// RemoveCondition removes condition by ID
	RemoveCondition(conditionID string)

	// GetConditions returns all conditions
	GetConditions() []Condition
}

// RewardCalculator computes combat rewards
type RewardCalculator interface {
	// CalculateExperience computes XP reward
	CalculateExperience(ctx context.Context, encounter Encounter) int64

	// CalculateExperiencePerParticipant computes XP per player entity
	CalculateExperiencePerParticipant(ctx context.Context, encounter Encounter) map[string]int64

	// CalculateLoot generates loot drops
	CalculateLoot(ctx context.Context, encounter Encounter) ([]LootDrop, error)

	// CalculateBonuses computes performance bonuses
	CalculateBonuses(ctx context.Context, encounter Encounter) map[string]float64

	// CalculateGold computes gold reward
	CalculateGold(ctx context.Context, encounter Encounter) int64

	// ApplyRewards distributes rewards to participants
	ApplyRewards(ctx context.Context, encounter Encounter, rewards Rewards) error

	// GetRewardMultiplier returns reward multiplier based on difficulty
	GetRewardMultiplier(encounter Encounter) float64
}

// Rewards contains all combat rewards
type Rewards struct {
	Experience          int64
	ExperiencePerPlayer map[string]int64
	Gold                int64
	Loot                []LootDrop
	Bonuses             map[string]float64
	Achievements        []string
}

// LootDrop represents dropped item
type LootDrop struct {
	ItemID   string
	Quantity int
	Rarity   int
	Source   string // which enemy dropped it
}

// PerformanceTracker tracks combat performance
type PerformanceTracker interface {
	// TrackDamage records damage dealt
	TrackDamage(entityID string, amount float64)

	// TrackHealing records healing done
	TrackHealing(entityID string, amount float64)

	// TrackKill records enemy defeated
	TrackKill(killerID, victimID string)

	// TrackAction records action performed
	TrackAction(entityID string, actionType ActionType)

	// TrackCritical records critical hit
	TrackCritical(entityID string)

	// TrackMiss records missed attack
	TrackMiss(entityID string)

	// TrackDeath records entity death
	TrackDeath(entityID string)

	// GetDamageDealt returns total damage by entity
	GetDamageDealt(entityID string) float64

	// GetHealingDone returns total healing by entity
	GetHealingDone(entityID string) float64

	// GetKills returns kill count
	GetKills(entityID string) int

	// GetDeaths returns death count
	GetDeaths(entityID string) int

	// GetActionsPerformed returns action count by type
	GetActionsPerformed(entityID string) map[ActionType]int

	// GetMVP returns most valuable participant
	GetMVP() string

	// CalculatePerformanceScore calculates performance score for entity
	CalculatePerformanceScore(entityID string) float64

	// GenerateReport creates performance summary
	GenerateReport() PerformanceReport

	// Reset clears all tracked data
	Reset()
}

// PerformanceReport summarizes combat performance
type PerformanceReport struct {
	TotalDamage       map[string]float64
	TotalHealing      map[string]float64
	KillCounts        map[string]int
	DeathCounts       map[string]int
	ActionCounts      map[string]map[ActionType]int
	CriticalHits      map[string]int
	Misses            map[string]int
	DamageTaken       map[string]float64
	HealingReceived   map[string]float64
	TurnsElapsed      int
	RoundsElapsed     int
	CombatDuration    int64
	MVP               string
	PerformanceGrades map[string]PerformanceGrade
	Achievements      []string
}

// PerformanceGrade rates participant performance
type PerformanceGrade string

const (
	GradeS PerformanceGrade = "S"
	GradeA PerformanceGrade = "A"
	GradeB PerformanceGrade = "B"
	GradeC PerformanceGrade = "C"
	GradeD PerformanceGrade = "D"
	GradeF PerformanceGrade = "F"
)

// EncounterBuilder creates encounters with fluent API
type EncounterBuilder interface {
	// WithArena sets combat arena
	WithArena(arena Arena) EncounterBuilder

	// WithParticipants adds participants
	WithParticipants(participants []Participant) EncounterBuilder

	// WithVictoryCondition adds win condition
	WithVictoryCondition(condition Condition) EncounterBuilder

	// WithDefeatCondition adds loss condition
	WithDefeatCondition(condition Condition) EncounterBuilder

	// WithTurnOrder sets turn order manager
	WithTurnOrder(turnOrder TurnOrder) EncounterBuilder

	// WithMaxRounds sets round limit
	WithMaxRounds(rounds int) EncounterBuilder

	// Build creates the encounter
	Build() (Encounter, error)

	// Reset resets builder to initial state
	Reset() EncounterBuilder
}
