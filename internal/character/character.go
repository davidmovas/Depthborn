package character

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/character/inventory"
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/core/skill"
)

// Character represents player-controlled entity
type Character interface {
	entity.Combatant

	// Class returns character class (empty string if classless)
	Class() string

	// SetClass updates character class
	SetClass(class string)

	// Experience returns current experience points
	Experience() int64

	// SetExperience updates experience points
	SetExperience(xp int64)

	// ExperienceToNextLevel returns XP needed for next level
	ExperienceToNextLevel() int64

	// AddExperience grants experience and handles level ups
	AddExperience(ctx context.Context, amount int64) error

	// SkillPoints returns unspent skill points
	SkillPoints() int

	// AddSkillPoints grants skill points
	AddSkillPoints(amount int)

	// SpendSkillPoint consumes one skill point
	SpendSkillPoint() error

	// SkillTree returns character skill tree
	SkillTree() skill.Tree

	// SkillLoadout returns equipped skills
	SkillLoadout() skill.Loadout

	// Inventory returns character inventory
	Inventory() inventory.Inventory

	// Equipment returns equipped items
	Equipment() inventory.Equipment

	// Stash returns character stash (if personal) or account stash
	Stash() inventory.Stash

	// Gold returns current gold amount
	Gold() int64

	// AddGold increases gold
	AddGold(amount int64)

	// RemoveGold decreases gold, returns false if insufficient
	RemoveGold(amount int64) bool

	// PlayTime returns total play time in seconds
	PlayTime() int64

	// AddPlayTime increments play time
	AddPlayTime(seconds int64)

	// DeathCount returns number of deaths
	DeathCount() int

	// IncrementDeathCount increases death counter
	IncrementDeathCount()

	// LastSave returns timestamp of last save
	LastSave() int64

	// UpdateLastSave sets last save to current time
	UpdateLastSave()

	// Flags returns character flags
	Flags() FlagSet

	// Statistics returns character statistics
	Statistics() Statistics
}

// FlagSet manages boolean flags for character state
type FlagSet interface {
	// Set sets flag to true
	Set(flag string)

	// Unset sets flag to false
	Unset(flag string)

	// Has returns true if flag is set
	Has(flag string) bool

	// Toggle inverts flag state
	Toggle(flag string)

	// GetAll returns all set flags
	GetAll() []string

	// Clear removes all flags
	Clear()
}

// Statistics tracks character gameplay statistics
type Statistics interface {
	// Get returns statistic value
	Get(stat string) int64

	// Set updates statistic value
	Set(stat string, value int64)

	// Increment increases statistic by amount
	Increment(stat string, amount int64)

	// GetAll returns all statistics
	GetAll() map[string]int64

	// Reset resets all statistics to zero
	Reset()
}

// Builder creates characters with fluent API
type Builder interface {
	// WithName sets character name
	WithName(name string) Builder

	// WithClass sets character class
	WithClass(class string) Builder

	// WithLevel sets starting level
	WithLevel(level int) Builder

	// WithAttributes sets base attributes
	WithAttributes(attrs map[string]float64) Builder

	// Build creates the character
	Build(ctx context.Context) (Character, error)
}

// Manager manages multiple characters
type Manager interface {
	// Create creates new character
	Create(ctx context.Context, name string) (Character, error)

	// Delete removes character
	Delete(ctx context.Context, characterID string) error

	// Get retrieves character by ID
	Get(characterID string) (Character, bool)

	// GetAll returns all characters
	GetAll() []Character

	// GetByName finds character by name
	GetByName(name string) (Character, bool)

	// Count returns total number of characters
	Count() int

	// MaxCharacters returns maximum allowed characters
	MaxCharacters() int

	// SetActive sets active character
	SetActive(characterID string) error

	// GetActive returns currently active character
	GetActive() (Character, bool)

	// Save persists character state
	Save(ctx context.Context, character Character) error

	// SaveAll persists all characters
	SaveAll(ctx context.Context) error

	// Load loads character from storage
	Load(ctx context.Context, characterID string) (Character, error)
}
