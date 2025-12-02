package character

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/davidmovas/Depthborn/internal/character/currency"
	"github.com/davidmovas/Depthborn/internal/character/equipment"
	"github.com/davidmovas/Depthborn/internal/character/inventory"
	"github.com/davidmovas/Depthborn/internal/character/progression"
	"github.com/davidmovas/Depthborn/internal/character/statistics"
	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/core/skill"
	"github.com/davidmovas/Depthborn/internal/core/status"
	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/infra/impl"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

var _ Character = (*BaseCharacter)(nil)

// BaseCharacter implements Character interface
type BaseCharacter struct {
	*entity.BaseCombatant

	mu sync.RWMutex

	// Identity
	accountID string

	// Managers
	equipmentMgr  equipment.Manager
	inventoryMgr  inventory.Manager
	currencyMgr   currency.Manager
	statisticsMgr statistics.Tracker
	skillTree     skill.Tree
	skillLoadout  skill.Loadout

	// Flags
	flags *BaseFlagSet

	// Timestamps
	playTime   int64
	lastPlayed int64
	lastSave   int64
	createdAt  int64
}

// Config holds configuration for creating a character
type Config struct {
	Name      string
	AccountID string

	// Base combat config
	InitialHealth float64
	MaxHealth     float64
	AttackRange   float64

	// Managers - if nil, defaults will be created
	Attributes   attribute.Manager
	Statuses     status.Manager
	Transform    spatial.Transform
	Tags         types.TagSet
	Callbacks    types.CallbackRegistry
	Progression  progression.Manager
	Equipment    equipment.Manager
	Inventory    inventory.Manager
	Currency     currency.Manager
	Statistics   statistics.Tracker
	SkillTree    skill.Tree
	SkillLoadout skill.Loadout

	// Initial values
	MaxWeight   float64
	InitialGold int64
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig(name string) Config {
	return Config{
		Name:          name,
		InitialHealth: 100,
		MaxHealth:     100,
		AttackRange:   1.5,
		MaxWeight:     100,
		InitialGold:   0,
	}
}

// NewCharacter creates a new character with the given configuration
func NewCharacter(cfg Config) *BaseCharacter {
	// Create default managers if not provided
	if cfg.Attributes == nil {
		cfg.Attributes = attribute.NewManager()
	}
	if cfg.Statuses == nil {
		cfg.Statuses = status.NewManager()
	}
	if cfg.Transform == nil {
		cfg.Transform = spatial.NewTransform(spatial.NewPosition(0, 0, 0), 0)
	}
	if cfg.Tags == nil {
		cfg.Tags = types.NewTagSet()
		cfg.Tags.Add("player")
		cfg.Tags.Add("character")
	}
	if cfg.Callbacks == nil {
		cfg.Callbacks = types.NewCallbackRegistry()
	}
	if cfg.Progression == nil {
		cfg.Progression = progression.NewManager(progression.ManagerConfig{})
	}

	// Create base combatant
	combatant := entity.NewCombatant(entity.CombatantConfig{
		LivingConfig: entity.LivingConfig{
			EntityConfig: entity.Config{
				Name:               cfg.Name,
				EntityType:         "character",
				AttributeManager:   cfg.Attributes,
				StatusManager:      cfg.Statuses,
				Transform:          cfg.Transform,
				TagSet:             cfg.Tags,
				Callbacks:          cfg.Callbacks,
				ProgressionManager: cfg.Progression,
			},
			InitialHealth: cfg.InitialHealth,
			MaxHealth:     cfg.MaxHealth,
		},
		AttackRange: cfg.AttackRange,
	})

	// Create character-specific managers
	if cfg.Equipment == nil {
		cfg.Equipment = equipment.NewManager()
	}
	if cfg.Inventory == nil {
		cfg.Inventory = inventory.NewManagerWithConfig(inventory.Config{
			MaxWeight: cfg.MaxWeight,
		})
	}
	if cfg.Currency == nil {
		currencyCfg := currency.DefaultConfig()
		currencyCfg.InitialGold = cfg.InitialGold
		cfg.Currency = currency.NewManagerWithConfig(currencyCfg)
	}
	if cfg.Statistics == nil {
		cfg.Statistics = statistics.NewTracker()
	}

	now := time.Now().Unix()

	char := &BaseCharacter{
		BaseCombatant: combatant,
		accountID:     cfg.AccountID,
		equipmentMgr:  cfg.Equipment,
		inventoryMgr:  cfg.Inventory,
		currencyMgr:   cfg.Currency,
		statisticsMgr: cfg.Statistics,
		skillTree:     cfg.SkillTree,
		skillLoadout:  cfg.SkillLoadout,
		flags:         NewFlagSet(),
		playTime:      0,
		lastPlayed:    now,
		lastSave:      0,
		createdAt:     now,
	}

	// Set equipment owner
	if equipMgr, ok := cfg.Equipment.(*equipment.BaseManager); ok {
		equipMgr.SetOwner(char)
	}

	return char
}

// --- Character interface implementation ---

func (c *BaseCharacter) Class() string {
	// Classless system - return empty string
	return ""
}

func (c *BaseCharacter) SetClass(_ string) {
	// Classless system - no-op
}

func (c *BaseCharacter) Experience() int64 {
	prog := c.Progression()
	if prog != nil && prog.Experience() != nil {
		return prog.Experience().CurrentExperience()
	}
	return 0
}

func (c *BaseCharacter) SetExperience(xp int64) {
	prog := c.Progression()
	if prog != nil && prog.Experience() != nil {
		_ = prog.Experience().SetExperience(xp)
	}
}

func (c *BaseCharacter) ExperienceToNextLevel() int64 {
	prog := c.Progression()
	if prog != nil && prog.Experience() != nil {
		return prog.Experience().ExperienceToNextLevel()
	}
	return 0
}

func (c *BaseCharacter) AddExperience(ctx context.Context, amount int64) error {
	prog := c.Progression()
	if prog != nil && prog.Experience() != nil {
		_, err := prog.Experience().AddExperience(ctx, amount)
		if err != nil {
			return err
		}
		// Update highest level statistic
		c.statisticsMgr.SetHighestLevel(c.Level())
	}
	return nil
}

func (c *BaseCharacter) SkillPoints() int {
	prog := c.Progression()
	if prog != nil && prog.SkillPoints() != nil {
		return prog.SkillPoints().AvailablePoints()
	}
	return 0
}

func (c *BaseCharacter) AddSkillPoints(amount int) {
	prog := c.Progression()
	if prog != nil && prog.SkillPoints() != nil {
		prog.SkillPoints().AddPoints(amount)
	}
}

func (c *BaseCharacter) SpendSkillPoint() error {
	prog := c.Progression()
	if prog != nil && prog.SkillPoints() != nil {
		return prog.SkillPoints().SpendPoints(1)
	}
	return fmt.Errorf("no skill point manager available")
}

func (c *BaseCharacter) SkillTree() skill.Tree {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.skillTree
}

func (c *BaseCharacter) SkillLoadout() skill.Loadout {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.skillLoadout
}

func (c *BaseCharacter) Inventory() inventory.Inventory {
	// This returns the old interface for compatibility
	// Internally we use the new Manager interface
	return nil // TODO: Create adapter if needed
}

func (c *BaseCharacter) InventoryManager() inventory.Manager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.inventoryMgr
}

func (c *BaseCharacter) Equipment() inventory.Equipment {
	// This returns the old interface for compatibility
	return nil // TODO: Create adapter if needed
}

func (c *BaseCharacter) EquipmentManager() equipment.Manager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.equipmentMgr
}

func (c *BaseCharacter) Stash() inventory.Stash {
	// Stash is at account level, not character
	return nil
}

func (c *BaseCharacter) Gold() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currencyMgr.Gold()
}

func (c *BaseCharacter) AddGold(amount int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if amount > 0 {
		c.statisticsMgr.AddGoldEarned(amount)
	}
	_ = c.currencyMgr.AddGold(amount)
}

func (c *BaseCharacter) RemoveGold(amount int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.currencyMgr.CanAffordGold(amount) {
		return false
	}

	if err := c.currencyMgr.AddGold(-amount); err != nil {
		return false
	}

	c.statisticsMgr.AddGoldSpent(amount)
	return true
}

func (c *BaseCharacter) CurrencyManager() currency.Manager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currencyMgr
}

func (c *BaseCharacter) PlayTime() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.playTime
}

func (c *BaseCharacter) AddPlayTime(seconds int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.playTime += seconds
	c.statisticsMgr.AddPlayTime(seconds)
}

func (c *BaseCharacter) DeathCount() int {
	return c.statisticsMgr.Deaths()
}

func (c *BaseCharacter) IncrementDeathCount() {
	c.statisticsMgr.AddDeath()
}

func (c *BaseCharacter) LastSave() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastSave
}

func (c *BaseCharacter) UpdateLastSave() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSave = time.Now().Unix()
}

func (c *BaseCharacter) LastPlayed() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastPlayed
}

func (c *BaseCharacter) SetLastPlayed(timestamp int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastPlayed = timestamp
}

func (c *BaseCharacter) Flags() FlagSet {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.flags
}

func (c *BaseCharacter) Statistics() Statistics {
	// Return adapter for old interface
	return &statisticsAdapter{tracker: c.statisticsMgr}
}

func (c *BaseCharacter) StatisticsTracker() statistics.Tracker {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.statisticsMgr
}

func (c *BaseCharacter) AccountID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.accountID
}

func (c *BaseCharacter) IsDead() bool {
	return !c.IsAlive()
}

func (c *BaseCharacter) Die(ctx context.Context, killerID string) error {
	c.IncrementDeathCount()
	return c.Kill(ctx, killerID)
}

func (c *BaseCharacter) Respawn(ctx context.Context, healthPercent float64) error {
	return c.Revive(ctx, healthPercent)
}

func (c *BaseCharacter) CreatedAt() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.createdAt
}

// --- Serialization ---

// CharacterState holds the complete serializable state
type CharacterState struct {
	// Identity
	ID        string `msgpack:"id"`
	Name      string `msgpack:"name"`
	AccountID string `msgpack:"account_id"`

	// Timestamps
	CreatedAt  int64 `msgpack:"created_at"`
	UpdatedAt  int64 `msgpack:"updated_at"`
	LastPlayed int64 `msgpack:"last_played"`
	LastSave   int64 `msgpack:"last_save"`
	PlayTime   int64 `msgpack:"play_time"`

	// Combat state
	Health    float64 `msgpack:"health"`
	MaxHealth float64 `msgpack:"max_health"`

	// Position
	Position *spatial.Position `msgpack:"position,omitempty"`
	Facing   float64           `msgpack:"facing"`

	// Attributes
	Attributes map[string]float64 `msgpack:"attributes,omitempty"`

	// Tags and Flags
	Tags  []string `msgpack:"tags,omitempty"`
	Flags []string `msgpack:"flags,omitempty"`

	// Component states (stored as sub-maps)
	EquipmentState   map[string]any `msgpack:"equipment,omitempty"`
	InventoryState   map[string]any `msgpack:"inventory,omitempty"`
	CurrencyState    map[string]any `msgpack:"currency,omitempty"`
	StatisticsState  map[string]any `msgpack:"statistics,omitempty"`
	ProgressionState map[string]any `msgpack:"progression,omitempty"`
}

func (c *BaseCharacter) MarshalBinary() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := CharacterState{
		ID:         c.ID(),
		Name:       c.Name(),
		AccountID:  c.accountID,
		CreatedAt:  c.createdAt,
		UpdatedAt:  c.UpdatedAt(),
		LastPlayed: c.lastPlayed,
		LastSave:   c.lastSave,
		PlayTime:   c.playTime,
		Health:     c.Health(),
		MaxHealth:  c.MaxHealth(),
	}

	// Position
	if c.Transform() != nil {
		pos := c.Transform().Position()
		state.Position = &pos
		state.Facing = float64(c.Transform().Facing())
	}

	// Attributes
	if c.Attributes() != nil {
		state.Attributes = make(map[string]float64)
		for k, v := range c.Attributes().Snapshot() {
			state.Attributes[string(k)] = v
		}
	}

	// Tags
	if c.Tags() != nil {
		state.Tags = c.Tags().All()
	}

	// Flags
	if c.flags != nil {
		state.Flags = c.flags.GetAll()
	}

	// Component states
	if c.equipmentMgr != nil {
		if eqState, err := c.equipmentMgr.SerializeState(); err == nil {
			state.EquipmentState = eqState
		}
	}
	if c.inventoryMgr != nil {
		if invState, err := c.inventoryMgr.SerializeState(); err == nil {
			state.InventoryState = invState
		}
	}
	if c.currencyMgr != nil {
		if curState, err := c.currencyMgr.SerializeState(); err == nil {
			state.CurrencyState = curState
		}
	}
	if c.statisticsMgr != nil {
		if statsState, err := c.statisticsMgr.SerializeState(); err == nil {
			state.StatisticsState = statsState
		}
	}
	if prog := c.Progression(); prog != nil {
		if progState, err := prog.SerializeState(); err == nil {
			state.ProgressionState = progState
		}
	}

	return persist.DefaultCodec().Encode(state)
}

func (c *BaseCharacter) UnmarshalBinary(data []byte) error {
	var state CharacterState
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return fmt.Errorf("failed to decode character state: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Restore identity
	c.BaseCombatant.BaseEntity.BasePersistent = impl.NewPersistentWithID(state.ID, "character", c, nil)
	c.BaseCombatant.BaseEntity.SetName(state.Name)
	c.accountID = state.AccountID

	// Restore timestamps
	c.createdAt = state.CreatedAt
	c.lastPlayed = state.LastPlayed
	c.lastSave = state.LastSave
	c.playTime = state.PlayTime

	// Restore combat state
	c.SetHealth(state.Health)

	// Restore position
	if state.Position != nil && c.Transform() != nil {
		c.Transform().SetPosition(*state.Position)
		c.Transform().SetFacing(spatial.Facing(state.Facing))
	}

	// Restore attributes
	if c.Attributes() != nil && len(state.Attributes) > 0 {
		snapshot := make(map[attribute.Type]float64, len(state.Attributes))
		for k, v := range state.Attributes {
			snapshot[attribute.Type(k)] = v
		}
		c.Attributes().Restore(snapshot)
	}

	// Restore tags
	if c.Tags() != nil {
		c.Tags().Clear()
		for _, tag := range state.Tags {
			c.Tags().Add(tag)
		}
	}

	// Restore flags
	if c.flags != nil {
		c.flags.Clear()
		for _, flag := range state.Flags {
			c.flags.Set(flag)
		}
	}

	// Restore component states
	if c.equipmentMgr != nil && state.EquipmentState != nil {
		_ = c.equipmentMgr.DeserializeState(state.EquipmentState)
	}
	if c.inventoryMgr != nil && state.InventoryState != nil {
		_ = c.inventoryMgr.DeserializeState(state.InventoryState)
	}
	if c.currencyMgr != nil && state.CurrencyState != nil {
		_ = c.currencyMgr.DeserializeState(state.CurrencyState)
	}
	if c.statisticsMgr != nil && state.StatisticsState != nil {
		_ = c.statisticsMgr.DeserializeState(state.StatisticsState)
	}
	if prog := c.Progression(); prog != nil && state.ProgressionState != nil {
		_ = prog.DeserializeState(state.ProgressionState)
	}

	return nil
}

func (c *BaseCharacter) Clone() any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Clone base combatant
	combatantClone := c.BaseCombatant.Clone().(*entity.BaseCombatant)

	// Clone managers
	var equipClone equipment.Manager
	if c.equipmentMgr != nil {
		equipClone = equipment.NewManager()
	}

	var invClone inventory.Manager
	if c.inventoryMgr != nil {
		invClone = inventory.NewManagerWithConfig(inventory.Config{
			MaxWeight: c.inventoryMgr.MaxWeight(),
		})
	}

	var curClone currency.Manager
	if c.currencyMgr != nil {
		curClone = currency.NewManager()
		for t, amount := range c.currencyMgr.GetAll() {
			curClone.Set(t, amount)
		}
	}

	var statsClone statistics.Tracker
	if baseStats, ok := c.statisticsMgr.(*statistics.BaseTracker); ok {
		statsClone = baseStats.Clone()
	} else {
		statsClone = statistics.NewTracker()
	}

	clone := &BaseCharacter{
		BaseCombatant: combatantClone,
		accountID:     c.accountID,
		equipmentMgr:  equipClone,
		inventoryMgr:  invClone,
		currencyMgr:   curClone,
		statisticsMgr: statsClone,
		skillTree:     nil, // Skill tree is not cloned
		skillLoadout:  nil, // Skill loadout is not cloned
		flags:         c.flags.Clone(),
		playTime:      c.playTime,
		lastPlayed:    c.lastPlayed,
		lastSave:      c.lastSave,
		createdAt:     c.createdAt,
	}

	return clone
}

func (c *BaseCharacter) Validate() error {
	if err := c.BaseCombatant.Validate(); err != nil {
		return err
	}

	if c.Name() == "" {
		return fmt.Errorf("character must have a name")
	}

	return nil
}

// --- Adapters for legacy interfaces ---

// statisticsAdapter adapts statistics.Tracker to Statistics interface
type statisticsAdapter struct {
	tracker statistics.Tracker
}

func (a *statisticsAdapter) Get(stat string) int64 {
	return a.tracker.Get(statistics.Stat(stat))
}

func (a *statisticsAdapter) Set(stat string, value int64) {
	a.tracker.Set(statistics.Stat(stat), value)
}

func (a *statisticsAdapter) Increment(stat string, amount int64) {
	a.tracker.Increment(statistics.Stat(stat), amount)
}

func (a *statisticsAdapter) GetAll() map[string]int64 {
	result := make(map[string]int64)
	for stat, value := range a.tracker.GetAll() {
		result[string(stat)] = value
	}
	return result
}

func (a *statisticsAdapter) Reset() {
	a.tracker.Reset()
}
