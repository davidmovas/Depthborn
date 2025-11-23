package quest

import (
	"context"
)

// Quest represents player quest/mission
type Quest interface {
	// ID returns unique quest identifier
	ID() string

	// Name returns quest name
	Name() string

	// Description returns quest description
	Description() string

	// Type returns quest type
	Type() Type

	// State returns current quest state
	State() State

	// SetState updates quest state
	SetState(state State)

	// Level returns recommended level
	Level() int

	// MinLevel returns minimum level required
	MinLevel() int

	// Objectives returns quest objectives
	Objectives() []Objective

	// GetObjective retrieves objective by ID
	GetObjective(objectiveID string) (Objective, bool)

	// UpdateObjective updates objective progress
	UpdateObjective(ctx context.Context, objectiveID string, progress int) error

	// CheckCompletion checks if all objectives are complete
	CheckCompletion() bool

	// Prerequisites returns required quests
	Prerequisites() []string

	// CanStart checks if quest can be started
	CanStart(ctx context.Context, characterID string) bool

	// Start begins quest
	Start(ctx context.Context, characterID string) error

	// Complete finishes quest
	Complete(ctx context.Context, characterID string) error

	// Fail fails quest
	Fail(ctx context.Context, characterID string) error

	// Abandon abandons quest
	Abandon(ctx context.Context, characterID string) error

	// CanAbandon checks if quest can be abandoned
	CanAbandon() bool

	// Rewards returns quest rewards
	Rewards() Rewards

	// TimeLimit returns time limit in seconds (0 = no limit)
	TimeLimit() int64

	// RemainingTime returns time left in seconds
	RemainingTime() int64

	// IsExpired returns true if time limit exceeded
	IsExpired() bool

	// Repeatable returns true if quest can be repeated
	Repeatable() bool

	// RepeatDelay returns cooldown between repeats in seconds
	RepeatDelay() int64

	// CompletionCount returns times completed
	CompletionCount() int

	// IncrementCompletionCount increases completion count
	IncrementCompletionCount()

	// GiverNPC returns quest giver NPC ID
	GiverNPC() string

	// TurnInNPC returns quest turn-in NPC ID
	TurnInNPC() string

	// Category returns quest category
	Category() string

	// Tags returns quest tags
	Tags() []string

	// IsHidden returns true if quest is hidden from UI
	IsHidden() bool

	// Priority returns quest priority (higher = more important)
	Priority() int

	// OnStart registers callback when quest starts
	OnStart(callback Callback)

	// OnComplete registers callback when quest completes
	OnComplete(callback Callback)

	// OnFail registers callback when quest fails
	OnFail(callback Callback)

	// OnObjectiveUpdate registers callback when objective updates
	OnObjectiveUpdate(callback ObjectiveCallback)
}

// Type categorizes quests
type Type string

const (
	QuestMain       Type = "main"
	QuestSide       Type = "side"
	QuestDaily      Type = "daily"
	QuestWeekly     Type = "weekly"
	QuestEvent      Type = "event"
	QuestTutorial   Type = "tutorial"
	QuestRepeatable Type = "repeatable"
	QuestHidden     Type = "hidden"
	QuestChallenge  Type = "challenge"
	QuestBounty     Type = "bounty"
)

// State represents quest progress
type State string

const (
	StateNotStarted State = "not_started"
	StateAvailable  State = "available"
	StateActive     State = "active"
	StateCompleted  State = "completed"
	StateFailed     State = "failed"
	StateAbandoned  State = "abandoned"
	StateTurnedIn   State = "turned_in"
)

// Callback is invoked for quest events
type Callback func(ctx context.Context, quest Quest, characterID string)

// ObjectiveCallback is invoked when objective updates
type ObjectiveCallback func(ctx context.Context, quest Quest, objective Objective, characterID string)

// Objective represents quest task
type Objective interface {
	// ID returns unique objective identifier
	ID() string

	// Description returns objective description
	Description() string

	// Type returns objective type
	Type() ObjectiveType

	// Target returns objective target (e.g., enemy type, item ID)
	Target() string

	// CurrentProgress returns current progress
	CurrentProgress() int

	// RequiredProgress returns progress needed
	RequiredProgress() int

	// SetProgress updates progress
	SetProgress(progress int)

	// AddProgress increases progress
	AddProgress(amount int)

	// IsComplete returns true if objective is complete
	IsComplete() bool

	// IsOptional returns true if objective is optional
	IsOptional() bool

	// IsHidden returns true if objective is hidden until revealed
	IsHidden() bool

	// Reveal makes hidden objective visible
	Reveal()

	// ProgressPercent returns completion percentage [0.0 - 1.0]
	ProgressPercent() float64

	// Location returns objective location hint
	Location() string

	// Markers returns map markers for objective
	Markers() []ObjectiveMarker
}

// ObjectiveType categorizes objectives
type ObjectiveType string

const (
	ObjectiveKill          ObjectiveType = "kill"
	ObjectiveCollect       ObjectiveType = "collect"
	ObjectiveDeliver       ObjectiveType = "deliver"
	ObjectiveInteract      ObjectiveType = "interact"
	ObjectiveReach         ObjectiveType = "reach"
	ObjectiveEscort        ObjectiveType = "escort"
	ObjectiveDefend        ObjectiveType = "defend"
	ObjectiveSurvive       ObjectiveType = "survive"
	ObjectiveCraft         ObjectiveType = "craft"
	ObjectiveUseSkill      ObjectiveType = "use_skill"
	ObjectiveTalkTo        ObjectiveType = "talk_to"
	ObjectiveExplore       ObjectiveType = "explore"
	ObjectiveDiscover      ObjectiveType = "discover"
	ObjectiveCompleteQuest ObjectiveType = "complete_quest"
	ObjectiveCustom        ObjectiveType = "custom"
)

// ObjectiveMarker represents map marker
type ObjectiveMarker struct {
	X          int
	Y          int
	Z          int
	MarkerType string
	Label      string
}

// Rewards contains quest rewards
type Rewards interface {
	// Experience returns XP reward
	Experience() int64

	// Gold returns gold reward
	Gold() int64

	// Items returns item rewards
	Items() []RewardItem

	// ChoiceItems returns items player can choose from
	ChoiceItems() []RewardItem

	// MaxChoices returns how many items player can choose
	MaxChoices() int

	// SkillPoints returns skill points reward
	SkillPoints() int

	// StatPoints returns stat points reward
	StatPoints() int

	// Reputation returns reputation rewards
	Reputation() map[string]int

	// Unlocks returns unlocked content
	Unlocks() []string

	// ScaleToLevel returns true if rewards scale with level
	ScaleToLevel() bool

	// GetScaledRewards returns rewards scaled to level
	GetScaledRewards(level int) Rewards

	// HasRewards returns true if has any rewards
	HasRewards() bool
}

// RewardItem represents quest reward item
type RewardItem struct {
	ItemID   string
	Quantity int
	IsChoice bool
}

// Chain represents series of connected quests
type Chain interface {
	// ID returns unique chain identifier
	ID() string

	// Name returns chain name
	Name() string

	// Description returns chain description
	Description() string

	// Quests returns all quests in chain
	Quests() []Quest

	// GetQuest retrieves quest by ID
	GetQuest(questID string) (Quest, bool)

	// NextQuest returns next quest in chain
	NextQuest(currentQuestID string) (Quest, bool)

	// PreviousQuest returns previous quest in chain
	PreviousQuest(currentQuestID string) (Quest, bool)

	// FirstQuest returns first quest in chain
	FirstQuest() Quest

	// LastQuest returns last quest in chain
	LastQuest() Quest

	// CurrentQuest returns active quest in chain
	CurrentQuest(characterID string) (Quest, bool)

	// Progress returns chain completion [0.0 - 1.0]
	Progress(characterID string) float64

	// IsComplete returns true if all quests complete
	IsComplete(characterID string) bool

	// ChainRewards returns bonus rewards for completing chain
	ChainRewards() Rewards
}

// Tracker tracks character quest progress
type Tracker interface {
	// ActiveQuests returns all active quests
	ActiveQuests() []Quest

	// CompletedQuests returns all completed quests
	CompletedQuests() []Quest

	// AvailableQuests returns quests that can be started
	AvailableQuests(characterID string) []Quest

	// GetQuest retrieves quest by ID
	GetQuest(questID string) (Quest, bool)

	// HasQuest checks if quest is tracked
	HasQuest(questID string) bool

	// IsActive checks if quest is active
	IsActive(questID string) bool

	// IsCompleted checks if quest is completed
	IsCompleted(questID string) bool

	// StartQuest starts quest
	StartQuest(ctx context.Context, questID string, characterID string) error

	// CompleteQuest completes quest
	CompleteQuest(ctx context.Context, questID string, characterID string) error

	// FailQuest fails quest
	FailQuest(ctx context.Context, questID string, characterID string) error

	// AbandonQuest abandons quest
	AbandonQuest(ctx context.Context, questID string, characterID string) error

	// UpdateObjective updates quest objective
	UpdateObjective(ctx context.Context, questID, objectiveID string, progress int, characterID string) error

	// CheckObjectiveTrigger checks if action triggers objective update
	CheckObjectiveTrigger(ctx context.Context, trigger ObjectiveTrigger, characterID string) error

	// MaxActiveQuests returns maximum concurrent active quests
	MaxActiveQuests() int

	// SetMaxActiveQuests updates active quest limit
	SetMaxActiveQuests(max int)

	// CanAcceptQuest checks if can accept more quests
	CanAcceptQuest() bool

	// Save persists tracker state
	Save(ctx context.Context) error

	// Load loads tracker state
	Load(ctx context.Context) error
}

// ObjectiveTrigger represents action that updates objectives
type ObjectiveTrigger interface {
	// Type returns trigger type
	Type() ObjectiveType

	// Target returns trigger target
	Target() string

	// Amount returns trigger amount
	Amount() int

	// Location returns trigger location
	Location() (x, y, z int)

	// Metadata returns additional trigger data
	Metadata() map[string]any
}

// Journal organizes quest information
type Journal interface {
	// GetEntries returns all journal entries
	GetEntries() []JournalEntry

	// GetEntry retrieves entry by quest ID
	GetEntry(questID string) (JournalEntry, bool)

	// AddEntry adds journal entry
	AddEntry(entry JournalEntry)

	// UpdateEntry updates journal entry
	UpdateEntry(questID string, entry JournalEntry)

	// RemoveEntry removes journal entry
	RemoveEntry(questID string)

	// GetEntriesByCategory returns entries in category
	GetEntriesByCategory(category string) []JournalEntry

	// GetEntriesByState returns entries with state
	GetEntriesByState(state State) []JournalEntry

	// Search searches entries by text
	Search(query string) []JournalEntry

	// Sort sorts entries by criteria
	Sort(criteria SortCriteria)

	// Clear removes all entries
	Clear()
}

// JournalEntry represents quest in journal
type JournalEntry struct {
	QuestID      string
	Name         string
	Description  string
	State        State
	Objectives   []Objective
	Category     string
	Priority     int
	StartTime    int64
	CompleteTime int64
	Notes        []string
	Pinned       bool
}

// SortCriteria defines journal sorting
type SortCriteria string

const (
	SortByName      SortCriteria = "name"
	SortByLevel     SortCriteria = "level"
	SortByPriority  SortCriteria = "priority"
	SortByProgress  SortCriteria = "progress"
	SortByStartTime SortCriteria = "start_time"
	SortByCategory  SortCriteria = "category"
)

// Condition defines quest availability requirement
type Condition interface {
	// Type returns condition type
	Type() ConditionType

	// Check evaluates if condition is met
	Check(ctx context.Context, characterID string) bool

	// Description returns human-readable condition
	Description() string

	// IsInverted returns true if condition should be negated
	IsInverted() bool
}

// ConditionType categorizes quest conditions
type ConditionType string

const (
	ConditionLevel             ConditionType = "level"
	ConditionQuestCompleted    ConditionType = "quest_completed"
	ConditionQuestNotCompleted ConditionType = "quest_not_completed"
	ConditionItemOwned         ConditionType = "item_owned"
	ConditionSkillLearned      ConditionType = "skill_learned"
	ConditionReputation        ConditionType = "reputation"
	ConditionClass             ConditionType = "class"
	ConditionAttribute         ConditionType = "attribute"
	ConditionLocation          ConditionType = "location"
	ConditionTime              ConditionType = "time"
	ConditionFlag              ConditionType = "flag"
	ConditionRandom            ConditionType = "random"
)

// Giver represents NPC that gives quests
type Giver interface {
	// ID returns unique giver identifier
	ID() string

	// Name returns giver name
	Name() string

	// AvailableQuests returns quests this giver offers
	AvailableQuests(characterID string) []Quest

	// GetQuest retrieves quest by ID
	GetQuest(questID string) (Quest, bool)

	// CanGiveQuest checks if can give quest to character
	CanGiveQuest(questID, characterID string) bool

	// GiveQuest gives quest to character
	GiveQuest(ctx context.Context, questID, characterID string) error

	// TurnInQuest completes quest with giver
	TurnInQuest(ctx context.Context, questID, characterID string) error

	// HasQuestMarker returns true if has available quests
	HasQuestMarker(characterID string) bool

	// QuestMarkerType returns marker type for UI
	QuestMarkerType(characterID string) MarkerType

	// Dialogue returns quest-related dialogue
	Dialogue(questID string, state State) string
}

// MarkerType indicates quest availability
type MarkerType string

const (
	MarkerAvailable  MarkerType = "available"
	MarkerInProgress MarkerType = "in_progress"
	MarkerComplete   MarkerType = "complete"
	MarkerRepeatable MarkerType = "repeatable"
	MarkerLocked     MarkerType = "locked"
)

// Registry manages all quests
type Registry interface {
	// Register adds quest to registry
	Register(quest Quest) error

	// Unregister removes quest from registry
	Unregister(questID string) error

	// Get retrieves quest by ID
	Get(questID string) (Quest, bool)

	// GetAll returns all registered quests
	GetAll() []Quest

	// GetByType returns quests of type
	GetByType(questType Type) []Quest

	// GetByCategory returns quests in category
	GetByCategory(category string) []Quest

	// GetByLevel returns quests for level range
	GetByLevel(minLevel, maxLevel int) []Quest

	// RegisterChain adds quest chain
	RegisterChain(chain Chain) error

	// UnregisterChain removes quest chain
	UnregisterChain(chainID string) error

	// GetChain retrieves chain by ID
	GetChain(chainID string) (Chain, bool)

	// GetChains returns all registered chains
	GetChains() []Chain

	// RegisterGiver adds quest giver
	RegisterGiver(giver Giver) error

	// UnregisterGiver removes quest giver
	UnregisterGiver(giverID string) error

	// GetGiver retrieves giver by ID
	GetGiver(giverID string) (Giver, bool)

	// GetGivers returns all registered givers
	GetGivers() []Giver

	// Search searches quests by criteria
	Search(criteria SearchCriteria) []Quest
}

// SearchCriteria defines quest search parameters
type SearchCriteria struct {
	Types      []Type
	Categories []string
	MinLevel   int
	MaxLevel   int
	Tags       []string
	State      State
	TextSearch string
}

// Builder creates quests with fluent API
type Builder interface {
	// WithName sets quest name
	WithName(name string) Builder

	// WithDescription sets description
	WithDescription(description string) Builder

	// WithType sets quest type
	WithType(questType Type) Builder

	// WithLevel sets recommended level
	WithLevel(level int) Builder

	// WithMinLevel sets minimum level
	WithMinLevel(level int) Builder

	// WithObjective adds objective
	WithObjective(objective Objective) Builder

	// WithRewards sets rewards
	WithRewards(rewards Rewards) Builder

	// WithPrerequisite adds prerequisite quest
	WithPrerequisite(questID string) Builder

	// WithGiver sets quest giver
	WithGiver(giverNPCID string) Builder

	// WithTurnIn sets turn-in NPC
	WithTurnIn(turnInNPCID string) Builder

	// WithTimeLimit sets time limit
	WithTimeLimit(seconds int64) Builder

	// WithRepeatable makes quest repeatable
	WithRepeatable(delay int64) Builder

	// WithCategory sets category
	WithCategory(category string) Builder

	// WithPriority sets priority
	WithPriority(priority int) Builder

	// Build creates the quest
	Build() (Quest, error)

	// Reset resets builder to initial state
	Reset() Builder
}

// Manager coordinates quest system
type Manager interface {
	// Registry returns quest registry
	Registry() Registry

	// Tracker returns quest tracker for character
	Tracker(characterID string) (Tracker, error)

	// Journal returns quest journal for character
	Journal(characterID string) (Journal, error)

	// ProcessTrigger processes objective trigger
	ProcessTrigger(ctx context.Context, trigger ObjectiveTrigger, characterID string) error

	// GetAvailableQuests returns quests available to character
	GetAvailableQuests(ctx context.Context, characterID string) []Quest

	// StartQuest starts quest for character
	StartQuest(ctx context.Context, questID, characterID string) error

	// CompleteQuest completes quest for character
	CompleteQuest(ctx context.Context, questID, characterID string, choices []string) error

	// AbandonQuest abandons quest for character
	AbandonQuest(ctx context.Context, questID, characterID string) error

	// ResetDaily resets all daily quests
	ResetDaily(ctx context.Context) error

	// ResetWeekly resets all weekly quests
	ResetWeekly(ctx context.Context) error

	// Update processes quest system
	Update(ctx context.Context, deltaMs int64) error

	// Save persists quest system state
	Save(ctx context.Context) error

	// Load loads quest system state
	Load(ctx context.Context) error
}
