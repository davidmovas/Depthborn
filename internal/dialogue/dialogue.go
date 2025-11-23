package dialogue

import (
	"context"
)

// Dialogue represents conversation with NPC
type Dialogue interface {
	// ID returns unique dialogue identifier
	ID() string

	// Name returns dialogue name
	Name() string

	// Description returns dialogue description
	Description() string

	// RootNode returns starting dialogue node
	RootNode() Node

	// GetNode retrieves node by ID
	GetNode(nodeID string) (Node, bool)

	// CurrentNode returns active dialogue node
	CurrentNode() Node

	// SetCurrentNode updates active node
	SetCurrentNode(nodeID string) error

	// Start begins dialogue
	Start(ctx context.Context, participantID string) error

	// End finishes dialogue
	End(ctx context.Context) error

	// IsActive returns true if dialogue is active
	IsActive() bool

	// Participants returns dialogue participant IDs
	Participants() []string

	// AddParticipant adds participant to dialogue
	AddParticipant(participantID string)

	// RemoveParticipant removes participant from dialogue
	RemoveParticipant(participantID string)

	// Variables returns dialogue variables
	Variables() VariableStore

	// OnStart registers callback when dialogue starts
	OnStart(callback Callback)

	// OnEnd registers callback when dialogue ends
	OnEnd(callback Callback)

	// OnNodeEnter registers callback when entering node
	OnNodeEnter(callback NodeCallback)
}

// Callback is invoked for dialogue events
type Callback func(ctx context.Context, dialogue Dialogue)

// NodeCallback is invoked for node events
type NodeCallback func(ctx context.Context, dialogue Dialogue, node Node)

// Node represents single dialogue state
type Node interface {
	// ID returns unique node identifier
	ID() string

	// Type returns node type
	Type() NodeType

	// SpeakerID returns who speaks this node
	SpeakerID() string

	// Text returns dialogue text
	Text() string

	// Choices returns available player choices
	Choices() []Choice

	// GetChoice retrieves choice by ID
	GetChoice(choiceID string) (Choice, bool)

	// CanShowChoice checks if choice should be displayed
	CanShowChoice(ctx context.Context, choiceID string, participantID string) bool

	// Actions returns actions to execute on node entry
	Actions() []Action

	// ExecuteActions executes all node actions
	ExecuteActions(ctx context.Context, participantID string) error

	// Conditions returns conditions to enter this node
	Conditions() []Condition

	// CheckConditions evaluates if node can be entered
	CheckConditions(ctx context.Context, participantID string) bool

	// NextNode returns next node ID (for linear dialogues)
	NextNode() string

	// Animation returns speaker animation/emotion
	Animation() string

	// CameraAngle returns camera angle hint
	CameraAngle() string

	// AudioCue returns audio cue identifier
	AudioCue() string

	// TimeLimit returns time limit for choices in milliseconds (0 = no limit)
	TimeLimit() int64

	// OnEnter is called when entering this node
	OnEnter(ctx context.Context, dialogue Dialogue, participantID string) error

	// OnExit is called when leaving this node
	OnExit(ctx context.Context, dialogue Dialogue, participantID string) error

	// Metadata returns node-specific data
	Metadata() map[string]any
}

// NodeType categorizes dialogue nodes
type NodeType string

const (
	NodeText     NodeType = "text"     // Simple text display
	NodeChoice   NodeType = "choice"   // Player choice
	NodeBranch   NodeType = "branch"   // Conditional branch
	NodeAction   NodeType = "action"   // Execute actions
	NodeRandom   NodeType = "random"   // Random next node
	NodeEnd      NodeType = "end"      // End dialogue
	NodeJump     NodeType = "jump"     // Jump to another dialogue
	NodeSubtitle NodeType = "subtitle" // Subtitle/narration
)

// Choice represents player dialogue option
type Choice interface {
	// ID returns unique choice identifier
	ID() string

	// Text returns choice text
	Text() string

	// NextNodeID returns node to go to if chosen
	NextNodeID() string

	// Conditions returns conditions to show choice
	Conditions() []Condition

	// CheckConditions evaluates if choice is available
	CheckConditions(ctx context.Context, participantID string) bool

	// Actions returns actions to execute if chosen
	Actions() []Action

	// ExecuteActions executes choice actions
	ExecuteActions(ctx context.Context, participantID string) error

	// IsEnabled returns true if choice can be selected
	IsEnabled() bool

	// DisabledReason returns why choice is disabled
	DisabledReason() string

	// Icon returns choice icon identifier
	Icon() string

	// Color returns choice color hint
	Color() string

	// SkillCheck returns skill check requirement
	SkillCheck() SkillCheck

	// OnceOnly returns true if choice disappears after selection
	OnceOnly() bool

	// WasChosen returns true if choice was previously selected
	WasChosen() bool

	// SetChosen marks choice as chosen
	SetChosen(chosen bool)

	// Tooltip returns choice tooltip
	Tooltip() string
}

// Action represents dialogue action
type Action interface {
	// Type returns action type
	Type() ActionType

	// Execute performs the action
	Execute(ctx context.Context, participantID string, dialogue Dialogue) error

	// Description returns action description
	Description() string

	// Parameters returns action parameters
	Parameters() map[string]interface{}
}

// ActionType categorizes dialogue actions
type ActionType string

const (
	ActionGiveItem       ActionType = "give_item"
	ActionTakeItem       ActionType = "take_item"
	ActionGiveQuest      ActionType = "give_quest"
	ActionCompleteQuest  ActionType = "complete_quest"
	ActionGiveExperience ActionType = "give_experience"
	ActionGiveGold       ActionType = "give_gold"
	ActionSetVariable    ActionType = "set_variable"
	ActionPlaySound      ActionType = "play_sound"
	ActionPlayAnimation  ActionType = "play_animation"
	ActionTeleport       ActionType = "teleport"
	ActionOpenShop       ActionType = "open_shop"
	ActionOpenCrafting   ActionType = "open_crafting"
	ActionStartCombat    ActionType = "start_combat"
	ActionSetReputation  ActionType = "set_reputation"
	ActionUnlockSkill    ActionType = "unlock_skill"
	ActionTriggerEvent   ActionType = "trigger_event"
	ActionSetFlag        ActionType = "set_flag"
	ActionCustom         ActionType = "custom"
)

// Condition defines dialogue requirement
type Condition interface {
	// Type returns condition type
	Type() ConditionType

	// Check evaluates if condition is met
	Check(ctx context.Context, participantID string, dialogue Dialogue) bool

	// Description returns human-readable condition
	Description() string

	// Parameters returns condition parameters
	Parameters() map[string]interface{}

	// IsInverted returns true if condition should be negated
	IsInverted() bool
}

// ConditionType categorizes dialogue conditions
type ConditionType string

const (
	ConditionVariable       ConditionType = "variable"
	ConditionQuestState     ConditionType = "quest_state"
	ConditionItemOwned      ConditionType = "item_owned"
	ConditionGoldAmount     ConditionType = "gold_amount"
	ConditionLevel          ConditionType = "level"
	ConditionClass          ConditionType = "class"
	ConditionReputation     ConditionType = "reputation"
	ConditionAttribute      ConditionType = "attribute"
	ConditionSkillLearned   ConditionType = "skill_learned"
	ConditionFlag           ConditionType = "flag"
	ConditionTimeOfDay      ConditionType = "time_of_day"
	ConditionRandom         ConditionType = "random"
	ConditionPreviousChoice ConditionType = "previous_choice"
	ConditionCustom         ConditionType = "custom"
)

// SkillCheck represents skill/attribute check
type SkillCheck interface {
	// Type returns check type
	Type() SkillCheckType

	// Attribute returns attribute being checked
	Attribute() string

	// Difficulty returns check difficulty
	Difficulty() int

	// Perform performs the skill check
	Perform(ctx context.Context, participantID string) (SkillCheckResult, error)

	// SuccessNodeID returns node to go to on success
	SuccessNodeID() string

	// FailureNodeID returns node to go to on failure
	FailureNodeID() string

	// AllowRetry returns true if check can be retried
	AllowRetry() bool

	// Description returns check description
	Description() string
}

// SkillCheckType categorizes skill checks
type SkillCheckType string

const (
	CheckPersuasion   SkillCheckType = "persuasion"
	CheckIntimidation SkillCheckType = "intimidation"
	CheckDeception    SkillCheckType = "deception"
	CheckInsight      SkillCheckType = "insight"
	CheckPerception   SkillCheckType = "perception"
	CheckStrength     SkillCheckType = "strength"
	CheckDexterity    SkillCheckType = "dexterity"
	CheckIntelligence SkillCheckType = "intelligence"
	CheckCharisma     SkillCheckType = "charisma"
	CheckLuck         SkillCheckType = "luck"
)

// SkillCheckResult describes check outcome
type SkillCheckResult struct {
	Success     bool
	Roll        int
	Difficulty  int
	Modifier    int
	IsCritical  bool
	Description string
}

// VariableStore manages dialogue variables
type VariableStore interface {
	// Set sets variable value
	Set(key string, value any)

	// Get retrieves variable value
	Get(key string) (any, bool)

	// GetString retrieves string variable
	GetString(key string) (string, bool)

	// GetInt retrieves int variable
	GetInt(key string) (int, bool)

	// GetFloat retrieves float64 variable
	GetFloat(key string) (float64, bool)

	// GetBool retrieves bool variable
	GetBool(key string) (bool, bool)

	// Has checks if variable exists
	Has(key string) bool

	// Delete removes variable
	Delete(key string)

	// Clear removes all variables
	Clear()

	// Keys returns all variable keys
	Keys() []string

	// Increment increases numeric variable
	Increment(key string, amount int)

	// Decrement decreases numeric variable
	Decrement(key string, amount int)
}

// =============================================================================
// DIALOGUE TREE
// =============================================================================

// Tree represents complete dialogue structure
type Tree interface {
	// ID returns unique tree identifier
	ID() string

	// Name returns tree name
	Name() string

	// RootNode returns starting node
	RootNode() Node

	// GetNode retrieves node by ID
	GetNode(nodeID string) (Node, bool)

	// AddNode adds node to tree
	AddNode(node Node) error

	// RemoveNode removes node from tree
	RemoveNode(nodeID string) error

	// GetNodes returns all nodes
	GetNodes() []Node

	// Validate checks tree structure validity
	Validate() error

	// Clone creates deep copy of tree
	Clone() Tree

	// Export exports tree data
	Export() TreeData

	// Import imports tree data
	Import(data TreeData) error
}

// TreeData contains exportable dialogue tree
type TreeData struct {
	ID          string
	Name        string
	Description string
	RootNodeID  string
	Nodes       []NodeData
	Variables   map[string]any
	Metadata    map[string]any
}

// NodeData contains exportable node data
type NodeData struct {
	ID         string
	Type       NodeType
	SpeakerID  string
	Text       string
	Choices    []ChoiceData
	Actions    []ActionData
	Conditions []ConditionData
	NextNodeID string
	Metadata   map[string]any
}

// ChoiceData contains exportable choice data
type ChoiceData struct {
	ID         string
	Text       string
	NextNodeID string
	Actions    []ActionData
	Conditions []ConditionData
	OnceOnly   bool
}

// ActionData contains exportable action data
type ActionData struct {
	Type       ActionType
	Parameters map[string]any
}

// ConditionData contains exportable condition data
type ConditionData struct {
	Type       ConditionType
	Parameters map[string]any
	Inverted   bool
}

// History tracks dialogue progression
type History interface {
	// RecordNode records visited node
	RecordNode(dialogueID, nodeID string, timestamp int64)

	// RecordChoice records selected choice
	RecordChoice(dialogueID, nodeID, choiceID string, timestamp int64)

	// GetVisitedNodes returns visited nodes for dialogue
	GetVisitedNodes(dialogueID string) []string

	// GetChoiceHistory returns choice history for dialogue
	GetChoiceHistory(dialogueID string) []ChoiceRecord

	// WasNodeVisited checks if node was visited
	WasNodeVisited(dialogueID, nodeID string) bool

	// WasChoiceSelected checks if choice was selected
	WasChoiceSelected(dialogueID, nodeID, choiceID string) bool

	// GetLastChoice returns last selected choice in node
	GetLastChoice(dialogueID, nodeID string) (string, bool)

	// Clear clears history for dialogue
	Clear(dialogueID string)

	// ClearAll clears all history
	ClearAll()
}

// ChoiceRecord records choice selection
type ChoiceRecord struct {
	DialogueID string
	NodeID     string
	ChoiceID   string
	Timestamp  int64
}

// Speaker represents entity that can speak
type Speaker interface {
	// ID returns unique speaker identifier
	ID() string

	// Name returns speaker display name
	Name() string

	// Portrait returns portrait image identifier
	Portrait() string

	// VoiceSet returns voice set identifier
	VoiceSet() string

	// DefaultAnimation returns default animation
	DefaultAnimation() string

	// GetAnimationForEmotion returns animation for emotion
	GetAnimationForEmotion(emotion string) string

	// Color returns speaker name color
	Color() string

	// IsPlayer returns true if speaker is player
	IsPlayer() bool
}

// Manager coordinates dialogue system
type Manager interface {
	// Registry returns dialogue registry
	Registry() Registry

	// StartDialogue starts dialogue with NPC
	StartDialogue(ctx context.Context, dialogueID string, participantID string, npcID string) (Dialogue, error)

	// GetActiveDialogue returns active dialogue for participant
	GetActiveDialogue(participantID string) (Dialogue, bool)

	// EndDialogue ends active dialogue
	EndDialogue(ctx context.Context, participantID string) error

	// SelectChoice selects dialogue choice
	SelectChoice(ctx context.Context, participantID string, choiceID string) error

	// ContinueDialogue continues to next node
	ContinueDialogue(ctx context.Context, participantID string) error

	// History returns dialogue history for participant
	History(participantID string) History

	// RegisterSpeaker registers speaker
	RegisterSpeaker(speaker Speaker) error

	// GetSpeaker retrieves speaker by ID
	GetSpeaker(speakerID string) (Speaker, bool)

	// Update processes dialogue system
	Update(ctx context.Context, deltaMs int64) error

	// Save persists dialogue state
	Save(ctx context.Context) error

	// Load loads dialogue state
	Load(ctx context.Context) error
}

// Registry manages dialogue trees
type Registry interface {
	// Register adds dialogue tree
	Register(tree Tree) error

	// Unregister removes dialogue tree
	Unregister(treeID string) error

	// Get retrieves tree by ID
	Get(treeID string) (Tree, bool)

	// GetAll returns all registered trees
	GetAll() []Tree

	// GetByNPC returns dialogues for NPC
	GetByNPC(npcID string) []Tree

	// GetByCategory returns dialogues in category
	GetByCategory(category string) []Tree

	// Has checks if tree is registered
	Has(treeID string) bool

	// Count returns total registered trees
	Count() int

	// Search searches trees by criteria
	Search(criteria SearchCriteria) []Tree
}

// SearchCriteria defines dialogue search parameters
type SearchCriteria struct {
	NPCID         string
	Category      string
	Tags          []string
	TextSearch    string
	HasActions    []ActionType
	HasConditions []ConditionType
}

// Builder creates dialogues with fluent API
type Builder interface {
	// WithName sets dialogue name
	WithName(name string) Builder

	// WithDescription sets description
	WithDescription(description string) Builder

	// WithRootNode sets starting node
	WithRootNode(node Node) Builder

	// AddNode adds node to dialogue
	AddNode(node Node) Builder

	// Build creates the dialogue tree
	Build() (Tree, error)

	// Reset resets builder to initial state
	Reset() Builder
}

// NodeBuilder creates dialogue nodes
type NodeBuilder interface {
	// WithSpeaker sets speaker
	WithSpeaker(speakerID string) NodeBuilder

	// WithText sets dialogue text
	WithText(text string) NodeBuilder

	// WithChoice adds player choice
	WithChoice(choice Choice) NodeBuilder

	// WithAction adds action
	WithAction(action Action) NodeBuilder

	// WithCondition adds condition
	WithCondition(condition Condition) NodeBuilder

	// WithNextNode sets next node
	WithNextNode(nodeID string) NodeBuilder

	// WithAnimation sets animation
	WithAnimation(animation string) NodeBuilder

	// Build creates the node
	Build() (Node, error)

	// Reset resets builder to initial state
	Reset() NodeBuilder
}

// Localizer handles dialogue translation
type Localizer interface {
	// GetText returns localized text
	GetText(key string, language string) string

	// SetText sets localized text
	SetText(key string, language string, text string)

	// HasTranslation checks if translation exists
	HasTranslation(key string, language string) bool

	// GetSupportedLanguages returns available languages
	GetSupportedLanguages() []string

	// CurrentLanguage returns active language
	CurrentLanguage() string

	// SetCurrentLanguage updates active language
	SetCurrentLanguage(language string)

	// LocalizeTree localizes entire dialogue tree
	LocalizeTree(tree Tree, language string) error
}

// Renderer prepares dialogue for display
type Renderer interface {
	// RenderNode renders node text
	RenderNode(node Node, variables VariableStore) string

	// RenderChoice renders choice text
	RenderChoice(choice Choice, variables VariableStore) string

	// ProcessTextTags processes special text tags
	ProcessTextTags(text string, variables VariableStore) string

	// FormatSpeakerName formats speaker name
	FormatSpeakerName(speaker Speaker) string

	// GetTextSpeed returns text display speed
	GetTextSpeed() float64

	// SetTextSpeed updates text display speed
	SetTextSpeed(speed float64)
}
