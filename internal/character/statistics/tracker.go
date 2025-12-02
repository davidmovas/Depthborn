package statistics

import (
	"sync"

	"github.com/davidmovas/Depthborn/pkg/persist"
)

// Stat defines different statistic types
type Stat string

const (
	// Combat statistics
	StatDeaths       Stat = "deaths"
	StatKills        Stat = "kills"
	StatBossKills    Stat = "boss_kills"
	StatEliteKills   Stat = "elite_kills"
	StatDamageDealt  Stat = "damage_dealt"
	StatDamageTaken  Stat = "damage_taken"
	StatHealingDone  Stat = "healing_done"
	StatHealingTaken Stat = "healing_taken"
	StatCriticalHits Stat = "critical_hits"
	StatDodges       Stat = "dodges"
	StatBlocks       Stat = "blocks"
	StatSkillsUsed   Stat = "skills_used"
	StatPotionsUsed  Stat = "potions_used"

	// Progress statistics
	StatHighestLevel    Stat = "highest_level"
	StatHighestDepth    Stat = "highest_depth"
	StatDungeonsCleared Stat = "dungeons_cleared"
	StatQuestsCompleted Stat = "quests_completed"

	// Economy statistics
	StatGoldEarned     Stat = "gold_earned"
	StatGoldSpent      Stat = "gold_spent"
	StatItemsLooted    Stat = "items_looted"
	StatItemsSold      Stat = "items_sold"
	StatItemsCrafted   Stat = "items_crafted"
	StatItemsEnchanted Stat = "items_enchanted"

	// Exploration statistics
	StatDistanceTraveled Stat = "distance_traveled"
	StatAreasDiscovered  Stat = "areas_discovered"
	StatSecretsFound     Stat = "secrets_found"
	StatChestsOpened     Stat = "chests_opened"

	// Time statistics (in seconds)
	StatPlayTime       Stat = "play_time"
	StatCombatTime     Stat = "combat_time"
	StatLongestSession Stat = "longest_session"
)

// AllStats returns all available statistic types
func AllStats() []Stat {
	return []Stat{
		StatDeaths, StatKills, StatBossKills, StatEliteKills,
		StatDamageDealt, StatDamageTaken, StatHealingDone, StatHealingTaken,
		StatCriticalHits, StatDodges, StatBlocks, StatSkillsUsed, StatPotionsUsed,
		StatHighestLevel, StatHighestDepth, StatDungeonsCleared, StatQuestsCompleted,
		StatGoldEarned, StatGoldSpent, StatItemsLooted, StatItemsSold,
		StatItemsCrafted, StatItemsEnchanted,
		StatDistanceTraveled, StatAreasDiscovered, StatSecretsFound, StatChestsOpened,
		StatPlayTime, StatCombatTime, StatLongestSession,
	}
}

// DisplayName returns human-readable stat name
func (s Stat) DisplayName() string {
	names := map[Stat]string{
		StatDeaths:           "Deaths",
		StatKills:            "Kills",
		StatBossKills:        "Bosses Killed",
		StatEliteKills:       "Elites Killed",
		StatDamageDealt:      "Damage Dealt",
		StatDamageTaken:      "Damage Taken",
		StatHealingDone:      "Healing Done",
		StatHealingTaken:     "Healing Received",
		StatCriticalHits:     "Critical Hits",
		StatDodges:           "Dodges",
		StatBlocks:           "Blocks",
		StatSkillsUsed:       "Skills Used",
		StatPotionsUsed:      "Potions Used",
		StatHighestLevel:     "Highest Level",
		StatHighestDepth:     "Deepest Layer",
		StatDungeonsCleared:  "Dungeons Cleared",
		StatQuestsCompleted:  "Quests Completed",
		StatGoldEarned:       "Gold Earned",
		StatGoldSpent:        "Gold Spent",
		StatItemsLooted:      "Items Looted",
		StatItemsSold:        "Items Sold",
		StatItemsCrafted:     "Items Crafted",
		StatItemsEnchanted:   "Items Enchanted",
		StatDistanceTraveled: "Distance Traveled",
		StatAreasDiscovered:  "Areas Discovered",
		StatSecretsFound:     "Secrets Found",
		StatChestsOpened:     "Chests Opened",
		StatPlayTime:         "Play Time",
		StatCombatTime:       "Combat Time",
		StatLongestSession:   "Longest Session",
	}
	if name, ok := names[s]; ok {
		return name
	}
	return string(s)
}

// Tracker tracks character statistics
type Tracker interface {
	// Get returns value of a statistic
	Get(stat Stat) int64

	// GetFloat returns float value (for damage/healing stats)
	GetFloat(stat Stat) float64

	// Set sets value directly
	Set(stat Stat, value int64)

	// SetFloat sets float value directly
	SetFloat(stat Stat, value float64)

	// Increment adds to a statistic
	Increment(stat Stat, amount int64)

	// IncrementFloat adds float to a statistic
	IncrementFloat(stat Stat, amount float64)

	// SetIfHigher sets value only if it's higher than current
	SetIfHigher(stat Stat, value int64)

	// Deaths returns total death count
	Deaths() int

	// AddDeath increments death counter
	AddDeath()

	// Kills returns total kill count
	Kills() int

	// AddKill increments kill counter
	AddKill()

	// DamageDealt returns total damage dealt
	DamageDealt() float64

	// AddDamageDealt adds to damage dealt
	AddDamageDealt(amount float64)

	// DamageTaken returns total damage taken
	DamageTaken() float64

	// AddDamageTaken adds to damage taken
	AddDamageTaken(amount float64)

	// HealingDone returns total healing done
	HealingDone() float64

	// AddHealingDone adds to healing done
	AddHealingDone(amount float64)

	// GoldEarned returns total gold earned
	GoldEarned() int64

	// AddGoldEarned adds to gold earned
	AddGoldEarned(amount int64)

	// GoldSpent returns total gold spent
	GoldSpent() int64

	// AddGoldSpent adds to gold spent
	AddGoldSpent(amount int64)

	// ItemsLooted returns total items looted
	ItemsLooted() int

	// AddItemsLooted increments items looted counter
	AddItemsLooted(count int)

	// BossesKilled returns boss kill count
	BossesKilled() int

	// AddBossKill increments boss kill counter
	AddBossKill()

	// HighestDepth returns deepest layer reached
	HighestDepth() int

	// SetHighestDepth updates highest depth if greater
	SetHighestDepth(depth int)

	// HighestLevel returns highest level reached
	HighestLevel() int

	// SetHighestLevel updates highest level if greater
	SetHighestLevel(level int)

	// CriticalHits returns total critical hits
	CriticalHits() int

	// AddCriticalHit increments critical hit counter
	AddCriticalHit()

	// PlayTime returns total play time in seconds
	PlayTime() int64

	// AddPlayTime adds play time in seconds
	AddPlayTime(seconds int64)

	// GetAll returns all statistics
	GetAll() map[Stat]int64

	// GetAllFloat returns all float statistics
	GetAllFloat() map[Stat]float64

	// Reset resets all statistics
	Reset()

	// SerializeState converts state to map for persistence
	SerializeState() (map[string]any, error)

	// DeserializeState restores state from map
	DeserializeState(state map[string]any) error
}

var _ Tracker = (*BaseTracker)(nil)

// BaseTracker implements Tracker interface
type BaseTracker struct {
	mu sync.RWMutex

	intStats   map[Stat]int64
	floatStats map[Stat]float64
}

// NewTracker creates a new statistics tracker
func NewTracker() *BaseTracker {
	return &BaseTracker{
		intStats:   make(map[Stat]int64),
		floatStats: make(map[Stat]float64),
	}
}

func (t *BaseTracker) Get(stat Stat) int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.intStats[stat]
}

func (t *BaseTracker) GetFloat(stat Stat) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.floatStats[stat]
}

func (t *BaseTracker) Set(stat Stat, value int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.intStats[stat] = value
}

func (t *BaseTracker) SetFloat(stat Stat, value float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.floatStats[stat] = value
}

func (t *BaseTracker) Increment(stat Stat, amount int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.intStats[stat] += amount
}

func (t *BaseTracker) IncrementFloat(stat Stat, amount float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.floatStats[stat] += amount
}

func (t *BaseTracker) SetIfHigher(stat Stat, value int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if value > t.intStats[stat] {
		t.intStats[stat] = value
	}
}

func (t *BaseTracker) Deaths() int {
	return int(t.Get(StatDeaths))
}

func (t *BaseTracker) AddDeath() {
	t.Increment(StatDeaths, 1)
}

func (t *BaseTracker) Kills() int {
	return int(t.Get(StatKills))
}

func (t *BaseTracker) AddKill() {
	t.Increment(StatKills, 1)
}

func (t *BaseTracker) DamageDealt() float64 {
	return t.GetFloat(StatDamageDealt)
}

func (t *BaseTracker) AddDamageDealt(amount float64) {
	t.IncrementFloat(StatDamageDealt, amount)
}

func (t *BaseTracker) DamageTaken() float64 {
	return t.GetFloat(StatDamageTaken)
}

func (t *BaseTracker) AddDamageTaken(amount float64) {
	t.IncrementFloat(StatDamageTaken, amount)
}

func (t *BaseTracker) HealingDone() float64 {
	return t.GetFloat(StatHealingDone)
}

func (t *BaseTracker) AddHealingDone(amount float64) {
	t.IncrementFloat(StatHealingDone, amount)
}

func (t *BaseTracker) GoldEarned() int64 {
	return t.Get(StatGoldEarned)
}

func (t *BaseTracker) AddGoldEarned(amount int64) {
	t.Increment(StatGoldEarned, amount)
}

func (t *BaseTracker) GoldSpent() int64 {
	return t.Get(StatGoldSpent)
}

func (t *BaseTracker) AddGoldSpent(amount int64) {
	t.Increment(StatGoldSpent, amount)
}

func (t *BaseTracker) ItemsLooted() int {
	return int(t.Get(StatItemsLooted))
}

func (t *BaseTracker) AddItemsLooted(count int) {
	t.Increment(StatItemsLooted, int64(count))
}

func (t *BaseTracker) BossesKilled() int {
	return int(t.Get(StatBossKills))
}

func (t *BaseTracker) AddBossKill() {
	t.Increment(StatBossKills, 1)
}

func (t *BaseTracker) HighestDepth() int {
	return int(t.Get(StatHighestDepth))
}

func (t *BaseTracker) SetHighestDepth(depth int) {
	t.SetIfHigher(StatHighestDepth, int64(depth))
}

func (t *BaseTracker) HighestLevel() int {
	return int(t.Get(StatHighestLevel))
}

func (t *BaseTracker) SetHighestLevel(level int) {
	t.SetIfHigher(StatHighestLevel, int64(level))
}

func (t *BaseTracker) CriticalHits() int {
	return int(t.Get(StatCriticalHits))
}

func (t *BaseTracker) AddCriticalHit() {
	t.Increment(StatCriticalHits, 1)
}

func (t *BaseTracker) PlayTime() int64 {
	return t.Get(StatPlayTime)
}

func (t *BaseTracker) AddPlayTime(seconds int64) {
	t.Increment(StatPlayTime, seconds)
}

func (t *BaseTracker) GetAll() map[Stat]int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make(map[Stat]int64, len(t.intStats))
	for stat, value := range t.intStats {
		result[stat] = value
	}
	return result
}

func (t *BaseTracker) GetAllFloat() map[Stat]float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make(map[Stat]float64, len(t.floatStats))
	for stat, value := range t.floatStats {
		result[stat] = value
	}
	return result
}

func (t *BaseTracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.intStats = make(map[Stat]int64)
	t.floatStats = make(map[Stat]float64)
}

// State holds serializable statistics state
type State struct {
	IntStats   map[string]int64   `msgpack:"int_stats"`
	FloatStats map[string]float64 `msgpack:"float_stats"`
}

func (t *BaseTracker) SerializeState() (map[string]any, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	intStats := make(map[string]int64, len(t.intStats))
	for stat, value := range t.intStats {
		intStats[string(stat)] = value
	}

	floatStats := make(map[string]float64, len(t.floatStats))
	for stat, value := range t.floatStats {
		floatStats[string(stat)] = value
	}

	state := State{
		IntStats:   intStats,
		FloatStats: floatStats,
	}

	data, err := persist.DefaultCodec().Encode(state)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := persist.DefaultCodec().Decode(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (t *BaseTracker) DeserializeState(stateData map[string]any) error {
	data, err := persist.DefaultCodec().Encode(stateData)
	if err != nil {
		return err
	}

	var state State
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.intStats = make(map[Stat]int64, len(state.IntStats))
	for stat, value := range state.IntStats {
		t.intStats[Stat(stat)] = value
	}

	t.floatStats = make(map[Stat]float64, len(state.FloatStats))
	for stat, value := range state.FloatStats {
		t.floatStats[Stat(stat)] = value
	}

	return nil
}

// Merge merges another tracker's stats into this one (additive)
func (t *BaseTracker) Merge(other Tracker) {
	if other == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for stat, value := range other.GetAll() {
		t.intStats[stat] += value
	}

	for stat, value := range other.GetAllFloat() {
		t.floatStats[stat] += value
	}
}

// Clone creates a copy of this tracker
func (t *BaseTracker) Clone() *BaseTracker {
	t.mu.RLock()
	defer t.mu.RUnlock()

	clone := NewTracker()

	for stat, value := range t.intStats {
		clone.intStats[stat] = value
	}

	for stat, value := range t.floatStats {
		clone.floatStats[stat] = value
	}

	return clone
}
