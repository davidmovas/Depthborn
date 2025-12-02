package account

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/davidmovas/Depthborn/internal/character/statistics"
	"github.com/davidmovas/Depthborn/internal/infra/impl"
	"github.com/davidmovas/Depthborn/pkg/persist"
)

// Account represents a player account that can have multiple characters
type Account interface {
	// ID returns unique account identifier
	ID() string

	// Name returns account/player name
	Name() string

	// SetName updates account name
	SetName(name string)

	// Stash returns the shared stash
	Stash() *Stash

	// Statistics returns global statistics across all characters
	Statistics() statistics.Tracker

	// CharacterIDs returns IDs of all characters on this account
	CharacterIDs() []string

	// AddCharacter registers a character to this account
	AddCharacter(characterID string)

	// RemoveCharacter unregisters a character from this account
	RemoveCharacter(characterID string)

	// CharacterCount returns number of characters
	CharacterCount() int

	// MaxCharacters returns maximum allowed characters
	MaxCharacters() int

	// SetMaxCharacters sets maximum allowed characters
	SetMaxCharacters(max int)

	// ActiveCharacterID returns the currently active character ID
	ActiveCharacterID() string

	// SetActiveCharacterID sets the active character
	SetActiveCharacterID(characterID string) error

	// CreatedAt returns account creation timestamp
	CreatedAt() int64

	// LastLogin returns last login timestamp
	LastLogin() int64

	// SetLastLogin updates last login timestamp
	SetLastLogin(timestamp int64)

	// TotalPlayTime returns combined play time of all characters
	TotalPlayTime() int64

	// AddPlayTime adds to total play time
	AddPlayTime(seconds int64)

	// Settings returns account settings
	Settings() Settings

	// Unlocks returns account-wide unlocks
	Unlocks() Unlocks

	// SerializeState converts state to map for persistence
	SerializeState() (map[string]any, error)

	// DeserializeState restores state from map
	DeserializeState(state map[string]any) error
}

// Settings manages account settings
type Settings interface {
	// Get returns a setting value
	Get(key string) any

	// GetString returns a string setting
	GetString(key string) string

	// GetInt returns an int setting
	GetInt(key string) int

	// GetBool returns a bool setting
	GetBool(key string) bool

	// Set sets a setting value
	Set(key string, value any)

	// Delete removes a setting
	Delete(key string)

	// GetAll returns all settings
	GetAll() map[string]any
}

// Unlocks manages account-wide unlocks
type Unlocks interface {
	// Has checks if unlock is obtained
	Has(unlockID string) bool

	// Add adds an unlock
	Add(unlockID string)

	// Remove removes an unlock
	Remove(unlockID string)

	// GetAll returns all unlocks
	GetAll() []string

	// Count returns number of unlocks
	Count() int
}

var _ Account = (*BaseAccount)(nil)

// BaseAccount implements Account interface
type BaseAccount struct {
	*impl.BasePersistent

	mu sync.RWMutex

	name              string
	stash             *Stash
	statistics        statistics.Tracker
	characterIDs      map[string]bool
	activeCharacterID string
	maxCharacters     int
	createdAt         int64
	lastLogin         int64
	totalPlayTime     int64
	settings          *BaseSettings
	unlocks           *BaseUnlocks
}

// Config holds configuration for creating an account
type Config struct {
	Name          string
	MaxCharacters int
	StashConfig   StashConfig
}

// DefaultConfig returns default configuration
func DefaultConfig(name string) Config {
	return Config{
		Name:          name,
		MaxCharacters: 10,
		StashConfig:   DefaultStashConfig(),
	}
}

// NewAccount creates a new account
func NewAccount(cfg Config) *BaseAccount {
	if cfg.MaxCharacters <= 0 {
		cfg.MaxCharacters = 10
	}

	now := time.Now().Unix()

	acc := &BaseAccount{
		name:          cfg.Name,
		stash:         NewStash(cfg.StashConfig),
		statistics:    statistics.NewTracker(),
		characterIDs:  make(map[string]bool),
		maxCharacters: cfg.MaxCharacters,
		createdAt:     now,
		lastLogin:     now,
		settings:      NewSettings(),
		unlocks:       NewUnlocks(),
	}

	acc.BasePersistent = impl.NewPersistent("account", acc, nil)

	return acc
}

func (a *BaseAccount) Name() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.name
}

func (a *BaseAccount) SetName(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.name = name
	a.Touch()
}

func (a *BaseAccount) Stash() *Stash {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.stash
}

func (a *BaseAccount) Statistics() statistics.Tracker {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.statistics
}

func (a *BaseAccount) CharacterIDs() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	ids := make([]string, 0, len(a.characterIDs))
	for id := range a.characterIDs {
		ids = append(ids, id)
	}
	return ids
}

func (a *BaseAccount) AddCharacter(characterID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.characterIDs[characterID] = true
	a.Touch()
}

func (a *BaseAccount) RemoveCharacter(characterID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.characterIDs, characterID)
	if a.activeCharacterID == characterID {
		a.activeCharacterID = ""
	}
	a.Touch()
}

func (a *BaseAccount) CharacterCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.characterIDs)
}

func (a *BaseAccount) MaxCharacters() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.maxCharacters
}

func (a *BaseAccount) SetMaxCharacters(max int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if max > 0 {
		a.maxCharacters = max
		a.Touch()
	}
}

func (a *BaseAccount) ActiveCharacterID() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.activeCharacterID
}

func (a *BaseAccount) SetActiveCharacterID(characterID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if characterID != "" && !a.characterIDs[characterID] {
		return fmt.Errorf("character %s not found on this account", characterID)
	}

	a.activeCharacterID = characterID
	a.Touch()
	return nil
}

func (a *BaseAccount) CreatedAt() int64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.createdAt
}

func (a *BaseAccount) LastLogin() int64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.lastLogin
}

func (a *BaseAccount) SetLastLogin(timestamp int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lastLogin = timestamp
	a.Touch()
}

func (a *BaseAccount) TotalPlayTime() int64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.totalPlayTime
}

func (a *BaseAccount) AddPlayTime(seconds int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.totalPlayTime += seconds
}

func (a *BaseAccount) Settings() Settings {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.settings
}

func (a *BaseAccount) Unlocks() Unlocks {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.unlocks
}

// State holds serializable account state
type State struct {
	ID                string         `msgpack:"id"`
	Name              string         `msgpack:"name"`
	CharacterIDs      []string       `msgpack:"character_ids"`
	ActiveCharacterID string         `msgpack:"active_character_id"`
	MaxCharacters     int            `msgpack:"max_characters"`
	CreatedAt         int64          `msgpack:"created_at"`
	LastLogin         int64          `msgpack:"last_login"`
	TotalPlayTime     int64          `msgpack:"total_play_time"`
	StashState        map[string]any `msgpack:"stash,omitempty"`
	StatisticsState   map[string]any `msgpack:"statistics,omitempty"`
	Settings          map[string]any `msgpack:"settings,omitempty"`
	Unlocks           []string       `msgpack:"unlocks,omitempty"`
}

func (a *BaseAccount) SerializeState() (map[string]any, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	charIDs := make([]string, 0, len(a.characterIDs))
	for id := range a.characterIDs {
		charIDs = append(charIDs, id)
	}

	state := State{
		ID:                a.ID(),
		Name:              a.name,
		CharacterIDs:      charIDs,
		ActiveCharacterID: a.activeCharacterID,
		MaxCharacters:     a.maxCharacters,
		CreatedAt:         a.createdAt,
		LastLogin:         a.lastLogin,
		TotalPlayTime:     a.totalPlayTime,
	}

	if a.stash != nil {
		if stashState, err := a.stash.SerializeState(); err == nil {
			state.StashState = stashState
		}
	}

	if a.statistics != nil {
		if statsState, err := a.statistics.SerializeState(); err == nil {
			state.StatisticsState = statsState
		}
	}

	if a.settings != nil {
		state.Settings = a.settings.GetAll()
	}

	if a.unlocks != nil {
		state.Unlocks = a.unlocks.GetAll()
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

func (a *BaseAccount) DeserializeState(stateData map[string]any) error {
	data, err := persist.DefaultCodec().Encode(stateData)
	if err != nil {
		return err
	}

	var state State
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.BasePersistent = impl.NewPersistentWithID(state.ID, "account", a, nil)
	a.name = state.Name
	a.activeCharacterID = state.ActiveCharacterID
	a.maxCharacters = state.MaxCharacters
	a.createdAt = state.CreatedAt
	a.lastLogin = state.LastLogin
	a.totalPlayTime = state.TotalPlayTime

	a.characterIDs = make(map[string]bool, len(state.CharacterIDs))
	for _, id := range state.CharacterIDs {
		a.characterIDs[id] = true
	}

	if a.stash != nil && state.StashState != nil {
		_ = a.stash.DeserializeState(state.StashState)
	}

	if a.statistics != nil && state.StatisticsState != nil {
		_ = a.statistics.DeserializeState(state.StatisticsState)
	}

	if a.settings != nil && state.Settings != nil {
		for k, v := range state.Settings {
			a.settings.Set(k, v)
		}
	}

	if a.unlocks != nil {
		for _, unlock := range state.Unlocks {
			a.unlocks.Add(unlock)
		}
	}

	return nil
}

// CanCreateCharacter checks if account can create more characters
func (a *BaseAccount) CanCreateCharacter() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.characterIDs) < a.maxCharacters
}

// MergeCharacterStatistics merges a character's statistics into account statistics
func (a *BaseAccount) MergeCharacterStatistics(ctx context.Context, charStats statistics.Tracker) {
	if charStats == nil {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if baseStats, ok := a.statistics.(*statistics.BaseTracker); ok {
		baseStats.Merge(charStats)
	}
}

// --- Settings implementation ---

var _ Settings = (*BaseSettings)(nil)

// BaseSettings implements Settings interface
type BaseSettings struct {
	mu       sync.RWMutex
	settings map[string]any
}

// NewSettings creates a new settings manager
func NewSettings() *BaseSettings {
	return &BaseSettings{
		settings: make(map[string]any),
	}
}

func (s *BaseSettings) Get(key string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.settings[key]
}

func (s *BaseSettings) GetString(key string) string {
	val := s.Get(key)
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func (s *BaseSettings) GetInt(key string) int {
	val := s.Get(key)
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func (s *BaseSettings) GetBool(key string) bool {
	val := s.Get(key)
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

func (s *BaseSettings) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settings[key] = value
}

func (s *BaseSettings) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.settings, key)
}

func (s *BaseSettings) GetAll() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]any, len(s.settings))
	for k, v := range s.settings {
		result[k] = v
	}
	return result
}

// --- Unlocks implementation ---

var _ Unlocks = (*BaseUnlocks)(nil)

// BaseUnlocks implements Unlocks interface
type BaseUnlocks struct {
	mu      sync.RWMutex
	unlocks map[string]bool
}

// NewUnlocks creates a new unlocks manager
func NewUnlocks() *BaseUnlocks {
	return &BaseUnlocks{
		unlocks: make(map[string]bool),
	}
}

func (u *BaseUnlocks) Has(unlockID string) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.unlocks[unlockID]
}

func (u *BaseUnlocks) Add(unlockID string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.unlocks[unlockID] = true
}

func (u *BaseUnlocks) Remove(unlockID string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.unlocks, unlockID)
}

func (u *BaseUnlocks) GetAll() []string {
	u.mu.RLock()
	defer u.mu.RUnlock()

	result := make([]string, 0, len(u.unlocks))
	for id := range u.unlocks {
		result = append(result, id)
	}
	return result
}

func (u *BaseUnlocks) Count() int {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return len(u.unlocks)
}
