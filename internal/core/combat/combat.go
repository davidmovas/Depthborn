package combat

import (
	"github.com/davidmovas/Depthborn/internal/core/entity"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

// Participant represents combatant in encounter
type Participant interface {
	// EntityID returns underlying entity identifier
	EntityID() string

	// Entity returns combatant entity
	Entity() entity.Combatant

	// Team returns participant team
	Team() Team

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

	// AvailableActions returns possible actions
	AvailableActions() []Action

	// CanPerformAction checks if action is possible
	CanPerformAction(action Action) bool

	// Modifiers returns active combat modifiers
	Modifiers() ModifierSet

	// Reactions returns available reactions
	Reactions() []Reaction
}

// AreaOfEffect defines spatial effect (используем spatial.Area)
type AreaOfEffect spatial.Area
