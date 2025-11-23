package party

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/character"
)

// Party represents group of characters
type Party interface {
	// ID returns unique party identifier
	ID() string

	// Add adds character to party
	Add(ctx context.Context, char character.Character) error

	// Remove removes character from party
	Remove(ctx context.Context, characterID string) error

	// Get retrieves character by ID
	Get(characterID string) (character.Character, bool)

	// GetAll returns all party members
	GetAll() []character.Character

	// Contains checks if character is in party
	Contains(characterID string) bool

	// Size returns number of members
	Size() int

	// MaxSize returns maximum party size
	MaxSize() int

	// IsFull returns true if party is at max size
	IsFull() bool

	// SetLeader sets party leader
	SetLeader(characterID string) error

	// GetLeader returns party leader
	GetLeader() (character.Character, bool)

	// SetActive sets active character for control
	SetActive(characterID string) error

	// GetActive returns currently controlled character
	GetActive() (character.Character, bool)

	// GetAlive returns all alive party members
	GetAlive() []character.Character

	// GetDead returns all dead party members
	GetDead() []character.Character

	// AllAlive returns true if all members are alive
	AllAlive() bool

	// AllDead returns true if all members are dead
	AllDead() bool

	// AverageLevel returns average level of all members
	AverageLevel() int

	// TotalPower returns combined power of all members
	TotalPower() float64

	// ShareExperience distributes experience among members
	ShareExperience(ctx context.Context, amount int64) error

	// ShareGold distributes gold among members
	ShareGold(amount int64) error

	// Formation returns party formation
	Formation() Formation

	// OnMemberAdded registers callback when member joins
	OnMemberAdded(callback MemberCallback)

	// OnMemberRemoved registers callback when member leaves
	OnMemberRemoved(callback MemberCallback)

	// OnMemberDied registers callback when member dies
	OnMemberDied(callback MemberCallback)
}

// MemberCallback is invoked for party member events
type MemberCallback func(ctx context.Context, party Party, member character.Character)

// Formation defines party positioning and roles
type Formation interface {
	// GetPosition returns position for character in formation
	GetPosition(characterID string) (x, y float64)

	// SetPosition updates character position in formation
	SetPosition(characterID string, x, y float64)

	// GetRole returns role assigned to character
	GetRole(characterID string) Role

	// SetRole assigns role to character
	SetRole(characterID string, role Role)

	// GetByRole returns all characters with specified role
	GetByRole(role Role) []string

	// Layout returns formation layout type
	Layout() FormationLayout

	// SetLayout changes formation layout
	SetLayout(layout FormationLayout)

	// Reset resets formation to default positions
	Reset()
}

// Role defines character role in party
type Role string

const (
	RoleTank    Role = "tank"
	RoleDamage  Role = "damage"
	RoleSupport Role = "support"
	RoleHealer  Role = "healer"
	RoleAny     Role = "any"
)

// FormationLayout defines party positioning pattern
type FormationLayout string

const (
	LayoutLine      FormationLayout = "line"
	LayoutColumn    FormationLayout = "column"
	LayoutStaggered FormationLayout = "staggered"
	LayoutCircle    FormationLayout = "circle"
	LayoutCustom    FormationLayout = "custom"
)

// Manager manages party lifecycle
type Manager interface {
	// Create creates new party
	Create(ctx context.Context) (Party, error)

	// Disband removes party
	Disband(ctx context.Context, partyID string) error

	// Get retrieves party by ID
	Get(partyID string) (Party, bool)

	// GetByCharacter finds party containing character
	GetByCharacter(characterID string) (Party, bool)

	// Save persists party state
	Save(ctx context.Context, party Party) error

	// Load loads party from storage
	Load(ctx context.Context, partyID string) (Party, error)
}
