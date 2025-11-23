package combat

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

// AttackAction represents basic attack
type AttackAction interface {
	Action

	// WeaponDamage returns weapon damage
	WeaponDamage() float64

	// AttackType returns attack type
	AttackType() AttackType

	// HitChance returns accuracy
	HitChance() float64

	// CriticalChance returns crit chance
	CriticalChance() float64

	// CriticalMultiplier returns crit multiplier
	CriticalMultiplier() float64

	// DamageType returns primary damage type
	DamageType() DamageType

	// ComboCounter returns combo count
	ComboCounter() int

	// IncrementCombo increases combo
	IncrementCombo()

	// ResetCombo resets combo to zero
	ResetCombo()

	// ComboBonus returns damage bonus from combo
	ComboBonus() float64

	// IsMultiHit returns true if attacks multiple times
	IsMultiHit() bool

	// HitCount returns number of hits
	HitCount() int
}

// AttackType categorizes attacks
type AttackType string

const (
	AttackMelee     AttackType = "melee"
	AttackRanged    AttackType = "ranged"
	AttackMagic     AttackType = "magic"
	AttackUnarmed   AttackType = "unarmed"
	AttackThrown    AttackType = "thrown"
	AttackArtillery AttackType = "artillery"
)

// SkillAction represents skill usage
type SkillAction interface {
	Action

	// SkillID returns skill identifier
	SkillID() string

	// ManaCost returns mana cost
	ManaCost() float64

	// Cooldown returns cooldown duration in milliseconds
	Cooldown() int64

	// RemainingCooldown returns time until skill is available
	RemainingCooldown() int64

	// SetCooldown updates cooldown timer
	SetCooldown(ms int64)

	// IsOnCooldown returns true if cannot be used
	IsOnCooldown() bool

	// CastTime returns casting time in milliseconds
	CastTime() int64

	// IsChanneled returns true if skill is channeled
	IsChanneled() bool

	// ChannelDuration returns channel duration
	ChannelDuration() int64

	// TickInterval returns channel tick interval
	TickInterval() int64

	// OnCastStart is called when casting begins
	OnCastStart(ctx context.Context, encounter Encounter) error

	// OnCastComplete is called when casting finishes
	OnCastComplete(ctx context.Context, encounter Encounter) error

	// OnCastInterrupt is called if casting is stopped
	OnCastInterrupt(ctx context.Context, encounter Encounter) error

	// OnChannelTick is called during channeling
	OnChannelTick(ctx context.Context, encounter Encounter, tickNumber int) error

	// Charges returns remaining charges (-1 = no charges)
	Charges() int

	// MaxCharges returns maximum charges
	MaxCharges() int

	// UseCharge consumes one charge
	UseCharge() bool

	// RestoreCharge replenishes one charge
	RestoreCharge()
}

// MoveAction represents movement
type MoveAction interface {
	Action

	// Destination returns target position
	Destination() spatial.Position

	// Path returns movement path
	Path() []spatial.Position

	// SetPath updates movement path
	SetPath(path []spatial.Position)

	// MovementCost returns action point cost per tile
	MovementCost() int

	// TotalCost returns total movement cost
	TotalCost() int

	// CanDash returns true if can dash/teleport
	CanDash() bool

	// DashDistance returns maximum dash distance
	DashDistance() float64

	// IsDash returns true if movement is dash/teleport
	IsDash() bool

	// TriggersOpportunityAttack returns true if provokes reactions
	TriggersOpportunityAttack() bool

	// IsStrategicRetreat returns true if tactical withdrawal
	IsStrategicRetreat() bool

	// DifficultTerrain returns movement penalty multiplier
	DifficultTerrain() float64
}

// DefendAction represents defensive stance
type DefendAction interface {
	Action

	// DefenseBonus returns defense increase
	DefenseBonus() float64

	// DamageReduction returns damage mitigation percentage
	DamageReduction() float64

	// Duration returns how long defense lasts in milliseconds
	Duration() int64

	// CounterAttackChance returns chance to counter [0.0 - 1.0]
	CounterAttackChance() float64

	// CounterDamage returns counter attack damage multiplier
	CounterDamage() float64

	// BlocksMovement returns true if cannot move while defending
	BlocksMovement() bool

	// ProtectsAllies returns true if protects adjacent allies
	ProtectsAllies() bool

	// InterceptRange returns range to intercept attacks
	InterceptRange() float64

	// ParryChance returns chance to parry [0.0 - 1.0]
	ParryChance() float64

	// IsGuarding returns true if actively defending
	IsGuarding() bool
}

// ItemAction represents item usage
type ItemAction interface {
	Action

	// ItemID returns item identifier
	ItemID() string

	// IsConsumable returns true if item is consumed
	IsConsumable() bool

	// UsageTime returns time to use item in milliseconds
	UsageTime() int64

	// Effect returns item effect
	Effect() ItemEffect

	// CanUseOnOthers returns true if can target allies
	CanUseOnOthers() bool

	// CanUseOnEnemies returns true if can target enemies
	CanUseOnEnemies() bool

	// Quantity returns item quantity used
	Quantity() int

	// SetQuantity updates quantity used
	SetQuantity(quantity int)
}

// ItemEffect describes item usage result
type ItemEffect interface {
	// Type returns effect category
	Type() ItemEffectType

	// Apply applies effect to target
	Apply(ctx context.Context, userID, targetID string, encounter Encounter) error

	// Description returns human-readable effect
	Description() string

	// Potency returns effect strength
	Potency() float64

	// Duration returns effect duration in milliseconds (0 = instant)
	Duration() int64
}

// ItemEffectType categorizes item effects
type ItemEffectType string

const (
	ItemEffectHealing  ItemEffectType = "healing"
	ItemEffectBuff     ItemEffectType = "buff"
	ItemEffectDebuff   ItemEffectType = "debuff"
	ItemEffectDamage   ItemEffectType = "damage"
	ItemEffectUtility  ItemEffectType = "utility"
	ItemEffectSummon   ItemEffectType = "summon"
	ItemEffectTeleport ItemEffectType = "teleport"
	ItemEffectRevive   ItemEffectType = "revive"
	ItemEffectCleanse  ItemEffectType = "cleanse"
)

// WaitAction represents turn delay
type WaitAction interface {
	Action

	// DelayAmount returns initiative delay
	DelayAmount() int

	// GrantsBonus returns true if waiting provides benefit
	GrantsBonus() bool

	// BonusType returns bonus granted by waiting
	BonusType() string

	// BonusValue returns bonus amount
	BonusValue() float64

	// IsDefensive returns true if enters defensive stance
	IsDefensive() bool
}

// FleeAction represents escape attempt
type FleeAction interface {
	Action

	// SuccessChance returns flee probability [0.0 - 1.0]
	SuccessChance() float64

	// FleeDirection returns escape direction
	FleeDirection() spatial.Direction

	// FleeDistance returns escape distance
	FleeDistance() float64

	// CanBePursued returns true if enemies can chase
	CanBePursued() bool

	// PenaltyOnFailure returns penalty if flee fails
	PenaltyOnFailure() FleePenalty

	// LeavesEncounter returns true if removes participant from combat
	LeavesEncounter() bool

	// RequiresClearPath returns true if needs unobstructed route
	RequiresClearPath() bool
}

// FleePenalty describes flee failure consequence
type FleePenalty interface {
	// LosesTurn returns true if loses next turn
	LosesTurn() bool

	// DamageTaken returns damage suffered on failure
	DamageTaken() float64

	// DebuffApplied returns debuff applied on failure
	DebuffApplied() string

	// StunDuration returns stun duration on failure
	StunDuration() int64

	// Description returns human-readable penalty
	Description() string
}

// InteractAction represents arena interaction
type InteractAction interface {
	Action

	// InteractiveID returns interactive object ID
	InteractiveID() string

	// InteractionType returns interaction category
	InteractionType() InteractionType

	// RequiresAdjacent returns true if must be next to object
	RequiresAdjacent() bool

	// InteractionTime returns time to complete interaction in milliseconds
	InteractionTime() int64

	// CanBeInterruptedByDamage returns true if damage stops interaction
	CanBeInterruptedByDamage() bool

	// InteractionRange returns maximum interaction distance
	InteractionRange() float64
}

// InteractionType categorizes interactions
type InteractionType string

const (
	InteractionActivate InteractionType = "activate"
	InteractionPickup   InteractionType = "pickup"
	InteractionOpen     InteractionType = "open"
	InteractionClose    InteractionType = "close"
	InteractionDestroy  InteractionType = "destroy"
	InteractionRepair   InteractionType = "repair"
	InteractionUse      InteractionType = "use"
)

// ComboAction represents coordinated attack
type ComboAction interface {
	Action

	// ParticipantIDs returns all participating entity IDs
	ParticipantIDs() []string

	// ComboType returns combo category
	ComboType() ComboType

	// BonusDamage returns extra damage from combo
	BonusDamage() float64

	// BonusMultiplier returns damage multiplier from combo
	BonusMultiplier() float64

	// SpecialEffect returns combo special effect
	SpecialEffect() string

	// RequiresSetup returns true if needs preparation
	RequiresSetup() bool

	// SetupActions returns action IDs needed to enable combo
	SetupActions() []string

	// IsSetupComplete checks if all setup actions performed
	IsSetupComplete(encounter Encounter) bool

	// WindowDuration returns time window for combo in milliseconds
	WindowDuration() int64

	// RemainingWindow returns time left to complete combo
	RemainingWindow() int64
}

// ComboType categorizes combo attacks
type ComboType string

const (
	ComboDualStrike      ComboType = "dual_strike"
	ComboPincer          ComboType = "pincer"
	ComboElementalFusion ComboType = "elemental_fusion"
	ComboAssassinate     ComboType = "assassinate"
	ComboOverwhelm       ComboType = "overwhelm"
	ComboRitual          ComboType = "ritual"
	ComboChain           ComboType = "chain"
	ComboFinisher        ComboType = "finisher"
)

// ActionChain represents action sequence
type ActionChain interface {
	// ID returns unique chain identifier
	ID() string

	// Actions returns ordered action sequence
	Actions() []Action

	// CurrentAction returns active action in chain
	CurrentAction() (Action, int, bool)

	// Next advances to next action
	Next() (Action, bool)

	// Previous returns to previous action
	Previous() (Action, bool)

	// CanContinue checks if chain can proceed
	CanContinue(encounter Encounter) bool

	// Break stops chain execution
	Break()

	// IsBroken returns true if chain was stopped
	IsBroken() bool

	// Progress returns completion percentage [0.0 - 1.0]
	Progress() float64

	// OnChainComplete is called when all actions execute
	OnChainComplete(ctx context.Context, encounter Encounter) error

	// OnChainBreak is called if chain is interrupted
	OnChainBreak(ctx context.Context, encounter Encounter) error

	// OnActionComplete is called after each action
	OnActionComplete(ctx context.Context, encounter Encounter, action Action) error
}

// ActionValidator validates action legality
type ActionValidator interface {
	// ValidateAction checks if action is legal
	ValidateAction(ctx context.Context, action Action, encounter Encounter) error

	// ValidateTarget checks if target is valid
	ValidateTarget(ctx context.Context, actorID, targetID string, action Action, encounter Encounter) error

	// ValidatePosition checks if position is valid
	ValidatePosition(ctx context.Context, pos spatial.Position, action Action, encounter Encounter) error

	// ValidateCost checks if actor can pay cost
	ValidateCost(ctx context.Context, actorID string, cost ActionCost, encounter Encounter) error

	// ValidateRange checks if target is in range
	ValidateRange(ctx context.Context, actorID, targetID string, range_ float64, encounter Encounter) error

	// ValidateLineOfSight checks if target is visible
	ValidateLineOfSight(ctx context.Context, actorID, targetID string, encounter Encounter) error

	// ValidateResourceAvailability checks if required resources exist
	ValidateResourceAvailability(ctx context.Context, actorID string, action Action, encounter Encounter) error
}

// ActionResolver executes actions
type ActionResolver interface {
	// Resolve executes action and returns result
	Resolve(ctx context.Context, action Action, encounter Encounter) (ActionResult, error)

	// ResolveAttack resolves attack action
	ResolveAttack(ctx context.Context, attack AttackAction, encounter Encounter) (ActionResult, error)

	// ResolveSkill resolves skill action
	ResolveSkill(ctx context.Context, skill SkillAction, encounter Encounter) (ActionResult, error)

	// ResolveMove resolves movement action
	ResolveMove(ctx context.Context, move MoveAction, encounter Encounter) (ActionResult, error)

	// ResolveDefend resolves defense action
	ResolveDefend(ctx context.Context, defend DefendAction, encounter Encounter) (ActionResult, error)

	// ResolveItem resolves item usage
	ResolveItem(ctx context.Context, item ItemAction, encounter Encounter) (ActionResult, error)

	// ResolveWait resolves wait action
	ResolveWait(ctx context.Context, wait WaitAction, encounter Encounter) (ActionResult, error)

	// ResolveFlee resolves flee action
	ResolveFlee(ctx context.Context, flee FleeAction, encounter Encounter) (ActionResult, error)

	// ResolveInteract resolves interact action
	ResolveInteract(ctx context.Context, interact InteractAction, encounter Encounter) (ActionResult, error)

	// ResolveCombo resolves combo action
	ResolveCombo(ctx context.Context, combo ComboAction, encounter Encounter) (ActionResult, error)

	// ResolveReaction resolves reactive action
	ResolveReaction(ctx context.Context, reaction Reaction, encounter Encounter) (ActionResult, error)
}

// ActionQueue manages pending actions
type ActionQueue interface {
	// Enqueue adds action to queue
	Enqueue(action Action) error

	// Dequeue removes and returns next action
	Dequeue() (Action, bool)

	// Peek returns next action without removing
	Peek() (Action, bool)

	// GetAll returns all queued actions
	GetAll() []Action

	// Clear removes all actions
	Clear()

	// Size returns number of queued actions
	Size() int

	// IsEmpty returns true if queue is empty
	IsEmpty() bool

	// Remove removes specific action
	Remove(actionID string) bool

	// Contains checks if action is queued
	Contains(actionID string) bool

	// Priority returns action priority
	Priority(action Action) int

	// Sort sorts queue by priority
	Sort()
}

// ActionFactory creates action instances
type ActionFactory interface {
	// CreateAttack creates attack action
	CreateAttack(actorID string, targetIDs []string, params map[string]any) (AttackAction, error)

	// CreateSkill creates skill action
	CreateSkill(actorID string, skillID string, targetIDs []string, params map[string]any) (SkillAction, error)

	// CreateMove creates move action
	CreateMove(actorID string, destination spatial.Position, params map[string]any) (MoveAction, error)

	// CreateDefend creates defend action
	CreateDefend(actorID string, params map[string]any) (DefendAction, error)

	// CreateItem creates item action
	CreateItem(actorID string, itemID string, targetIDs []string, params map[string]any) (ItemAction, error)

	// CreateWait creates wait action
	CreateWait(actorID string, params map[string]any) (WaitAction, error)

	// CreateFlee creates flee action
	CreateFlee(actorID string, params map[string]any) (FleeAction, error)

	// CreateInteract creates interact action
	CreateInteract(actorID string, interactiveID string, params map[string]any) (InteractAction, error)

	// CreateCombo creates combo action
	CreateCombo(participantIDs []string, comboType ComboType, params map[string]any) (ComboAction, error)
}
