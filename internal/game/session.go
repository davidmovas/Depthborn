package game

import (
	"context"
	"time"

	"github.com/davidmovas/Depthborn/internal/character"
	"github.com/davidmovas/Depthborn/internal/character/party"
	"github.com/davidmovas/Depthborn/internal/core/combat"
	"github.com/davidmovas/Depthborn/internal/event/types"
	"github.com/davidmovas/Depthborn/internal/world/layer"
)

// Session represents active game session
type Session interface {
	// ID returns unique session identifier
	ID() string

	// State returns current session state
	State() SessionState

	// SetState updates session state
	SetState(state SessionState)

	// Start begins game session
	Start(ctx context.Context) error

	// Stop ends game session
	Stop(ctx context.Context) error

	// Pause pauses game session
	Pause()

	// Resume resumes game session
	Resume()

	// Update processes game tick
	Update(ctx context.Context, deltaMs int64) error

	// IsRunning returns true if session is active
	IsRunning() bool

	// IsPaused returns true if session is paused
	IsPaused() bool

	// Player returns player manager
	Player() PlayerManager

	// World returns world manager
	World() WorldManager

	// Combat returns combat manager
	Combat() CombatManager

	// Camp returns camp manager
	Camp() CampManager

	// Time returns game time manager
	Time() TimeManager

	// Events returns event bus
	Events() types.EventBus

	// CurrentLocation returns current location
	CurrentLocation() Location

	// ChangeLocation changes player location
	ChangeLocation(ctx context.Context, location Location) error

	// SaveGame saves game state
	SaveGame(ctx context.Context, slotID string) error

	// LoadGame loads game state
	LoadGame(ctx context.Context, slotID string) error

	// ElapsedTime returns total session time in milliseconds
	ElapsedTime() int64

	// OnStateChange registers callback when state changes
	OnStateChange(callback StateChangeCallback)
}

// SessionState represents session status
type SessionState string

const (
	StateInitializing      SessionState = "initializing"
	StateMainMenu          SessionState = "main_menu"
	StateCharacterCreation SessionState = "character_creation"
	StateInGame            SessionState = "in_game"
	StateInCombat          SessionState = "in_combat"
	StateInCamp            SessionState = "in_camp"
	StatePaused            SessionState = "paused"
	StateLoading           SessionState = "loading"
	StateSaving            SessionState = "saving"
	StateEnding            SessionState = "ending"
)

// StateChangeCallback is invoked when session state changes
type StateChangeCallback func(ctx context.Context, oldState, newState SessionState)

// Location represents player location
type Location interface {
	// Type returns location type
	Type() LocationType

	// ID returns location identifier
	ID() string

	// Name returns location name
	Name() string

	// Description returns location description
	Description() string

	// CanSave returns true if can save game here
	CanSave() bool

	// CanRest returns true if can rest/heal here
	CanRest() bool

	// IsSafe returns true if no enemies can spawn
	IsSafe() bool

	// OnEnter is called when entering location
	OnEnter(ctx context.Context, session Session) error

	// OnExit is called when leaving location
	OnExit(ctx context.Context, session Session) error

	// Update processes location tick
	Update(ctx context.Context, deltaMs int64) error
}

// LocationType categorizes locations
type LocationType string

const (
	LocationCamp       LocationType = "camp"
	LocationDungeon    LocationType = "dungeon"
	LocationTown       LocationType = "town"
	LocationWilderness LocationType = "wilderness"
	LocationBoss       LocationType = "boss"
	LocationSecret     LocationType = "secret"
)

// PlayerManager manages player state
type PlayerManager interface {
	// ActiveCharacter returns active character
	ActiveCharacter() (character.Character, error)

	// SetActiveCharacter changes active character
	SetActiveCharacter(characterID string) error

	// Party returns active party
	Party() (party.Party, error)

	// Gold returns player gold
	Gold() int64

	// AddGold increases gold
	AddGold(amount int64)

	// RemoveGold decreases gold
	RemoveGold(amount int64) bool

	// HasGold checks if has enough gold
	HasGold(amount int64) bool

	// OnCharacterDeath is called when character dies
	OnCharacterDeath(ctx context.Context, characterID string) error

	// OnCharacterRevive is called when character revives
	OnCharacterRevive(ctx context.Context, characterID string) error

	// Statistics returns player statistics
	Statistics() PlayerStatistics

	// Save persists player state
	Save(ctx context.Context) error

	// Load loads player state
	Load(ctx context.Context) error
}

// PlayerStatistics tracks player gameplay stats
type PlayerStatistics struct {
	TotalPlayTime   int64
	CombatsWon      int
	CombatsLost     int
	EnemiesKilled   int
	BossesKilled    int
	Deaths          int
	GoldEarned      int64
	GoldSpent       int64
	ItemsCrafted    int
	QuestsCompleted int
	DeepestLayer    int
	HighestLevel    int
	Achievements    []string
}

// WorldManager manages world state
type WorldManager interface {
	// CurrentLayer returns current dungeon layer
	CurrentLayer() (layer.Layer, error)

	// EnterLayer enters dungeon layer
	EnterLayer(ctx context.Context, depth int) error

	// ExitLayer exits current layer
	ExitLayer(ctx context.Context) error

	// Registry returns layer registry
	Registry() layer.Registry

	// CanDescend checks if can go deeper
	CanDescend() bool

	// CanAscend checks if can go up
	CanAscend() bool

	// Descend goes to next layer
	Descend(ctx context.Context) error

	// Ascend goes to previous layer
	Ascend(ctx context.Context) error

	// ReturnToCamp teleports to camp
	ReturnToCamp(ctx context.Context) error

	// CurrentDepth returns current depth
	CurrentDepth() int

	// MaxDepthReached returns deepest reached layer
	MaxDepthReached() int

	// Update processes world tick
	Update(ctx context.Context, deltaMs int64) error
}

// CombatManager manages combat encounters
type CombatManager interface {
	// StartCombat initiates combat encounter
	StartCombat(ctx context.Context, params CombatParams) (combat.Encounter, error)

	// CurrentCombat returns active encounter
	CurrentCombat() (combat.Encounter, bool)

	// EndCombat ends current encounter
	EndCombat(ctx context.Context, result combat.EncounterResult) error

	// IsInCombat returns true if combat is active
	IsInCombat() bool

	// ProcessTurn processes combat turn
	ProcessTurn(ctx context.Context) error

	// PerformAction performs combat action
	PerformAction(ctx context.Context, action combat.Action) (combat.ActionResult, error)

	// OnCombatStart registers callback when combat starts
	OnCombatStart(callback CombatCallback)

	// OnCombatEnd registers callback when combat ends
	OnCombatEnd(callback CombatCallback)

	// Update processes combat tick
	Update(ctx context.Context, deltaMs int64) error
}

// CombatParams defines combat initiation parameters
type CombatParams struct {
	ArenaID     string
	EnemyTypes  []string
	EnemyLevels []int
	BossType    string
	Depth       int
	Modifiers   []string
	IsAmbush    bool
	CanFlee     bool
	Conditions  []combat.Condition
}

// CombatCallback is invoked for combat events
type CombatCallback func(ctx context.Context, encounter combat.Encounter)

// CampManager manages camp/hub
type CampManager interface {
	// EnterCamp enters camp
	EnterCamp(ctx context.Context) error

	// ExitCamp exits camp
	ExitCamp(ctx context.Context) error

	// IsInCamp returns true if in camp
	IsInCamp() bool

	// RestoreParty fully heals party
	RestoreParty(ctx context.Context) error

	// AccessStorage accesses storage
	AccessStorage(ctx context.Context) error

	// AccessForge accesses forge
	AccessForge(ctx context.Context) error

	// AccessVendor accesses vendor
	AccessVendor(ctx context.Context, vendorID string) error

	// AccessTraining accesses training
	AccessTraining(ctx context.Context) error

	// Update processes camp tick
	Update(ctx context.Context, deltaMs int64) error
}

// TimeManager manages game time
type TimeManager interface {
	// CurrentTime returns current game time
	CurrentTime() Time

	// AddTime advances game time
	AddTime(duration time.Duration)

	// RealTimeElapsed returns real time elapsed in milliseconds
	RealTimeElapsed() int64

	// GameTimeElapsed returns game time elapsed in milliseconds
	GameTimeElapsed() int64

	// TimeScale returns time scaling factor
	TimeScale() float64

	// SetTimeScale updates time scale
	SetTimeScale(scale float64)

	// Pause pauses time
	Pause()

	// Resume resumes time
	Resume()

	// IsPaused returns true if time is paused
	IsPaused() bool

	// OnDayChange registers callback when day changes
	OnDayChange(callback TimeCallback)

	// OnHourChange registers callback when hour changes
	OnHourChange(callback TimeCallback)
}

// Time represents in-game time
type Time struct {
	Day    int
	Hour   int
	Minute int
	Second int
}

// TimeCallback is invoked for time events
type TimeCallback func(ctx context.Context, time Time)

// EventType categorizes game events
type EventType string

const (
	EventGameStart      EventType = "game_start"
	EventGameEnd        EventType = "game_end"
	EventLevelUp        EventType = "level_up"
	EventCharacterDeath EventType = "character_death"
	EventItemAcquired   EventType = "item_acquired"
	EventQuestStarted   EventType = "quest_started"
	EventQuestCompleted EventType = "quest_completed"
	EventCombatStart    EventType = "combat_start"
	EventCombatEnd      EventType = "combat_end"
	EventLayerEntered   EventType = "layer_entered"
	EventBossDefeated   EventType = "boss_defeated"
	EventAchievement    EventType = "achievement"
	EventDialogueStart  EventType = "dialogue_start"
	EventDialogueEnd    EventType = "dialogue_end"
)

// Loop manages game update loop
type Loop interface {
	// Start starts game loop
	Start(ctx context.Context) error

	// Stop stops game loop
	Stop()

	// IsRunning returns true if loop is running
	IsRunning() bool

	// TickRate returns updates per second
	TickRate() int

	// SetTickRate updates tick rate
	SetTickRate(ticksPerSecond int)

	// DeltaTime returns time since last tick in milliseconds
	DeltaTime() int64

	// FrameCount returns total frames processed
	FrameCount() int64

	// FPS returns frames per second
	FPS() float64

	// OnTick registers callback for each tick
	OnTick(callback TickCallback)
}

// TickCallback is invoked each game tick
type TickCallback func(ctx context.Context, deltaMs int64) error

// SaveManager manages game saves
type SaveManager interface {
	// Save creates save file
	Save(ctx context.Context, slot SaveSlot) error

	// Load loads save file
	Load(ctx context.Context, slotID string) (SaveSlot, error)

	// Delete deletes save file
	Delete(ctx context.Context, slotID string) error

	// GetSlots returns all save slots
	GetSlots() []SaveSlot

	// GetSlot retrieves save slot by ID
	GetSlot(slotID string) (SaveSlot, bool)

	// HasSave checks if save exists
	HasSave(slotID string) bool

	// AutoSave performs automatic save
	AutoSave(ctx context.Context) error

	// QuickSave performs quick save
	QuickSave(ctx context.Context) error

	// QuickLoad performs quick load
	QuickLoad(ctx context.Context) error

	// CanSave checks if saving is allowed
	CanSave() bool

	// LastSaveTime returns last save timestamp
	LastSaveTime() int64
}

// SaveSlot represents save file
type SaveSlot interface {
	// ID returns unique slot identifier
	ID() string

	// Name returns save name
	Name() string

	// SetName updates save name
	SetName(name string)

	// Timestamp returns save timestamp
	Timestamp() int64

	// PlayTime returns total play time in seconds
	PlayTime() int64

	// CharacterLevel returns character level
	CharacterLevel() int

	// CharacterName returns character name
	CharacterName() string

	// CurrentDepth returns dungeon depth
	CurrentDepth() int

	// Location returns save location
	Location() string

	// Screenshot returns screenshot data (optional)
	Screenshot() []byte

	// IsAutoSave returns true if auto-save
	IsAutoSave() bool

	// IsQuickSave returns true if quick-save
	IsQuickSave() bool

	// Data returns save data
	Data() SaveData
}

// SaveData contains complete game state
type SaveData struct {
	Version     string
	Timestamp   int64
	PlayTime    int64
	Characters  []CharacterData
	Party       PartyData
	Inventory   InventoryData
	World       WorldData
	Quests      QuestData
	Dialogues   DialogueData
	Progression ProgressionData
	Statistics  PlayerStatistics
	Settings    Settings
	CustomData  map[string]any
}

// CharacterData contains character save data
type CharacterData struct {
	ID         string
	Name       string
	Class      string
	Level      int
	Experience int64
	Attributes map[string]float64
	Skills     []string
	Equipment  map[string]string
	Inventory  []string
	Quests     []string
	Flags      []string
	Statistics map[string]interface{}
}

// PartyData contains party save data
type PartyData struct {
	Members   []string
	Leader    string
	Active    string
	Formation string
}

// InventoryData contains inventory save data
type InventoryData struct {
	Items     []ItemData
	Gold      int64
	Stash     []ItemData
	StashTabs []StashTabData
}

// ItemData contains item save data
type ItemData struct {
	ID         string
	Type       string
	Quantity   int
	Quality    float64
	Level      int
	Rarity     int
	Affixes    []string
	Identified bool
	Equipped   bool
	Slot       string
}

// StashTabData contains stash tab save data
type StashTabData struct {
	Name  string
	Color string
	Icon  string
	Items []ItemData
}

// WorldData contains world save data
type WorldData struct {
	CurrentDepth    int
	MaxDepth        int
	VisitedLayers   []int
	LayerStates     map[int]LayerState
	CurrentLocation string
}

// LayerState contains layer save data
type LayerState struct {
	Depth            int
	Explored         bool
	BossDefeated     bool
	LastVisited      int64
	KilledEnemies    []string
	CollectedLoot    []string
	ActivatedObjects []string
}

// QuestData contains quest save data
type QuestData struct {
	ActiveQuests    []QuestState
	CompletedQuests []string
	FailedQuests    []string
	QuestVariables  map[string]map[string]any
}

// QuestState contains quest progress
type QuestState struct {
	QuestID         string
	State           string
	StartTime       int64
	Objectives      []ObjectiveState
	CompletionCount int
}

// ObjectiveState contains objective progress
type ObjectiveState struct {
	ObjectiveID string
	Progress    int
	Complete    bool
}

// DialogueData contains dialogue save data
type DialogueData struct {
	History       []DialogueHistory
	Variables     map[string]map[string]any
	ChoiceHistory []ChoiceHistory
}

// DialogueHistory tracks dialogue progression
type DialogueHistory struct {
	DialogueID   string
	VisitedNodes []string
	LastNode     string
	Timestamp    int64
}

// ChoiceHistory tracks dialogue choices
type ChoiceHistory struct {
	DialogueID string
	NodeID     string
	ChoiceID   string
	Timestamp  int64
}

// ProgressionData contains progression save data
type ProgressionData struct {
	Experience      int64
	SkillPoints     int
	StatPoints      int
	Prestige        int
	UnlockedContent []string
	Achievements    []string
}

// Settings contains game configuration
type Settings struct {
	Difficulty       string
	AutoSave         bool
	AutoSaveInterval int64
	ShowTutorials    bool
	CustomSettings   map[string]interface{}
}

// Factory creates game sessions
type Factory interface {
	// CreateSession creates new game session
	CreateSession(ctx context.Context) (Session, error)

	// CreateSessionWithConfig creates session with configuration
	CreateSessionWithConfig(ctx context.Context, config SessionConfig) (Session, error)

	// LoadSession loads existing session
	LoadSession(ctx context.Context, saveSlotID string) (Session, error)
}

// SessionConfig defines session configuration
type SessionConfig struct {
	TickRate         int
	AutoSaveInterval int64
	MaxAutoSaves     int
	StartLocation    LocationType
	Difficulty       string
	EnableTutorial   bool
	CustomSettings   map[string]interface{}
}

// Coordinator coordinates all game systems
type Coordinator interface {
	// Session returns active session
	Session() (Session, error)

	// CreateNewGame creates new game
	CreateNewGame(ctx context.Context, characterName string, class string) (Session, error)

	// LoadGame loads saved game
	LoadGame(ctx context.Context, slotID string) (Session, error)

	// SaveGame saves current game
	SaveGame(ctx context.Context, slotID string) error

	// EndGame ends current game
	EndGame(ctx context.Context) error

	// HasActiveSession returns true if session exists
	HasActiveSession() bool

	// Update processes coordinator tick
	Update(ctx context.Context, deltaMs int64) error
}
