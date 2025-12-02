package skill

import (
	"context"
	"errors"
	"sync"

	"github.com/davidmovas/Depthborn/pkg/identifier"
)

var _ Instance = (*BaseInstance)(nil)

// ErrMaxLevel is returned when trying to level up at max level
var ErrMaxLevel = errors.New("skill is at maximum level")

// ErrInvalidLevel is returned for invalid level values
var ErrInvalidLevel = errors.New("invalid skill level")

// ErrNoCharges is returned when trying to use skill with no charges
var ErrNoCharges = errors.New("no charges available")

// ErrOnCooldown is returned when skill is on cooldown
var ErrOnCooldown = errors.New("skill is on cooldown")

// ErrInsufficientResources is returned when caster lacks resources
var ErrInsufficientResources = errors.New("insufficient resources")

// BaseInstance implements the Instance interface.
// Represents a skill owned by a character with runtime state.
type BaseInstance struct {
	mu sync.RWMutex

	id       string // Unique instance ID
	defID    string // Source definition ID
	def      Def    // Reference to definition (may be nil)
	level    int    // Current skill level (1-based)
	isActive bool   // For toggle/aura skills

	// Cooldown state
	cooldownRemaining int64 // Remaining cooldown in ms

	// Charge state
	charges        int   // Current available charges
	chargeRecovery int64 // Time accumulator for charge recovery

	// Modifiers affecting this skill
	modifiers []SkillModifier
}

// InstanceConfig holds configuration for creating BaseInstance
type InstanceConfig struct {
	Def        Def
	StartLevel int
}

// NewBaseInstance creates a new skill instance from definition
func NewBaseInstance(config InstanceConfig) *BaseInstance {
	level := config.StartLevel
	if level < 1 && config.Def.MaxLevel() > 0 {
		level = 1
	}

	inst := &BaseInstance{
		id:        identifier.New(),
		defID:     config.Def.ID(),
		def:       config.Def,
		level:     level,
		isActive:  false,
		modifiers: make([]SkillModifier, 0),
	}

	// Initialize charges if skill uses charges
	if config.Def.BaseCharges() > 0 {
		inst.charges = inst.MaxCharges()
	}

	return inst
}

// NewBaseInstanceFromData creates instance from serialized data (no def reference)
func NewBaseInstanceFromData(defID string, level int) *BaseInstance {
	return &BaseInstance{
		id:        identifier.New(),
		defID:     defID,
		def:       nil,
		level:     level,
		modifiers: make([]SkillModifier, 0),
	}
}

// SetDef sets the definition reference (used when loading from registry)
func (i *BaseInstance) SetDef(def Def) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.def = def

	// Initialize charges if not set
	if def.BaseCharges() > 0 && i.charges == 0 {
		i.charges = i.maxChargesLocked()
	}
}

func (i *BaseInstance) DefID() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.defID
}

func (i *BaseInstance) Def() Def {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.def
}

func (i *BaseInstance) Level() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.level
}

func (i *BaseInstance) SetLevel(level int) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.def == nil {
		return ErrInvalidLevel
	}

	maxLevel := i.def.MaxLevel()
	if maxLevel == 0 {
		return ErrInvalidLevel // Skill has no levels
	}

	if level < 1 || level > maxLevel {
		return ErrInvalidLevel
	}

	i.level = level
	return nil
}

func (i *BaseInstance) CanLevelUp() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.def == nil {
		return false
	}

	maxLevel := i.def.MaxLevel()
	if maxLevel == 0 {
		return false // Skill has no levels
	}

	return i.level < maxLevel
}

func (i *BaseInstance) LevelUp() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.def == nil {
		return ErrInvalidLevel
	}

	maxLevel := i.def.MaxLevel()
	if maxLevel == 0 {
		return ErrMaxLevel
	}

	if i.level >= maxLevel {
		return ErrMaxLevel
	}

	i.level++
	return nil
}

func (i *BaseInstance) CurrentLevelData() LevelData {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.def == nil {
		return nil
	}

	return i.def.LevelData(i.level)
}

func (i *BaseInstance) Cooldown() int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.cooldownRemaining
}

func (i *BaseInstance) SetCooldown(ms int64) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if ms < 0 {
		ms = 0
	}
	i.cooldownRemaining = ms
}

func (i *BaseInstance) IsOnCooldown() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.cooldownRemaining > 0
}

func (i *BaseInstance) Charges() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.charges
}

func (i *BaseInstance) MaxCharges() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.maxChargesLocked()
}

func (i *BaseInstance) maxChargesLocked() int {
	if i.def == nil {
		return 0
	}

	// Check if level data overrides charges
	if levelData := i.def.LevelData(i.level); levelData != nil {
		if c := levelData.Charges(); c > 0 {
			return c
		}
	}

	return i.def.BaseCharges()
}

func (i *BaseInstance) UseCharge() bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.charges <= 0 {
		return false
	}

	i.charges--
	return true
}

func (i *BaseInstance) ChargeRecoveryProgress() float64 {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.def == nil || i.def.ChargeRecovery() <= 0 {
		return 0
	}

	maxCharges := i.maxChargesLocked()
	if i.charges >= maxCharges {
		return 1.0 // Fully charged
	}

	return float64(i.chargeRecovery) / float64(i.def.ChargeRecovery())
}

func (i *BaseInstance) Update(deltaMs int64) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Update cooldown
	if i.cooldownRemaining > 0 {
		i.cooldownRemaining -= deltaMs
		if i.cooldownRemaining < 0 {
			i.cooldownRemaining = 0
		}
	}

	// Update charge recovery
	if i.def != nil && i.def.BaseCharges() > 0 {
		maxCharges := i.maxChargesLocked()
		recoveryTime := i.def.ChargeRecovery()

		if i.charges < maxCharges && recoveryTime > 0 {
			i.chargeRecovery += deltaMs

			// Recover charges
			for i.chargeRecovery >= recoveryTime && i.charges < maxCharges {
				i.charges++
				i.chargeRecovery -= recoveryTime
			}

			// Cap recovery progress
			if i.charges >= maxCharges {
				i.chargeRecovery = 0
			}
		}
	}
}

func (i *BaseInstance) CanUse(ctx context.Context, casterID string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	_ = ctx      // Will be used for resource checks
	_ = casterID // Will be used for entity lookup

	if i.def == nil {
		return false
	}

	// Check cooldown
	if i.cooldownRemaining > 0 {
		return false
	}

	// Check charges (if skill uses charges)
	if i.def.BaseCharges() > 0 && i.charges <= 0 {
		return false
	}

	// TODO: Check resource availability on caster
	// TODO: Check requirements

	return true
}

func (i *BaseInstance) Use(ctx context.Context, casterID string, params ActivationParams) (Result, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.def == nil {
		return Result{Success: false, Message: "skill definition not loaded"}, nil
	}

	// Check cooldown
	if i.cooldownRemaining > 0 {
		return Result{Success: false, Message: "skill is on cooldown"}, ErrOnCooldown
	}

	// Check and consume charges
	usesCharges := i.def.BaseCharges() > 0
	if usesCharges {
		if i.charges <= 0 {
			return Result{Success: false, Message: "no charges available"}, ErrNoCharges
		}
		i.charges--
	}

	// Get cooldown for this level
	cooldown := i.getCooldownLocked()

	// Apply cooldown (only if not using charges, or using last charge)
	if !usesCharges || i.charges == 0 {
		i.cooldownRemaining = cooldown
	}

	// TODO: Actually execute skill effects
	// For now, return success result
	result := Result{
		Success:    true,
		Message:    "skill executed",
		TargetsHit: []string{},
		Effects:    make(map[string]TargetResult),
		Metadata:   make(map[string]any),
	}

	// Get level data for resource costs
	if levelData := i.def.LevelData(i.level); levelData != nil {
		result.ResourcesConsumed = levelData.ResourceCosts()
	}

	_ = ctx
	_ = casterID
	_ = params

	return result, nil
}

func (i *BaseInstance) getCooldownLocked() int64 {
	if i.def == nil {
		return 0
	}

	// Check level-specific cooldown
	if levelData := i.def.LevelData(i.level); levelData != nil {
		if cd := levelData.Cooldown(); cd > 0 {
			return cd
		}
	}

	// Apply modifiers
	baseCooldown := i.def.BaseCooldown()
	for _, mod := range i.modifiers {
		baseCooldown = mod.ModifyCooldown(baseCooldown)
	}

	return baseCooldown
}

func (i *BaseInstance) IsActive() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.isActive
}

func (i *BaseInstance) SetActive(active bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.isActive = active
}

func (i *BaseInstance) Modifiers() []SkillModifier {
	i.mu.RLock()
	defer i.mu.RUnlock()

	result := make([]SkillModifier, len(i.modifiers))
	copy(result, i.modifiers)
	return result
}

// AddModifier adds a modifier to this skill instance
func (i *BaseInstance) AddModifier(mod SkillModifier) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.modifiers = append(i.modifiers, mod)
}

// RemoveModifier removes a modifier by ID
func (i *BaseInstance) RemoveModifier(modID string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	for idx, mod := range i.modifiers {
		if mod.ID() == modID {
			i.modifiers = append(i.modifiers[:idx], i.modifiers[idx+1:]...)
			return
		}
	}
}

// ClearModifiers removes all modifiers
func (i *BaseInstance) ClearModifiers() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.modifiers = make([]SkillModifier, 0)
}

// =============================================================================
// SERIALIZATION
// =============================================================================

// InstanceState holds serializable state
type InstanceState struct {
	DefID             string `msgpack:"def_id"`
	Level             int    `msgpack:"level"`
	IsActive          bool   `msgpack:"is_active"`
	CooldownRemaining int64  `msgpack:"cooldown"`
	Charges           int    `msgpack:"charges"`
	ChargeRecovery    int64  `msgpack:"charge_recovery"`
}

// GetState returns serializable state
func (i *BaseInstance) GetState() InstanceState {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return InstanceState{
		DefID:             i.defID,
		Level:             i.level,
		IsActive:          i.isActive,
		CooldownRemaining: i.cooldownRemaining,
		Charges:           i.charges,
		ChargeRecovery:    i.chargeRecovery,
	}
}

// RestoreState restores from serialized state
func (i *BaseInstance) RestoreState(state InstanceState) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.defID = state.DefID
	i.level = state.Level
	i.isActive = state.IsActive
	i.cooldownRemaining = state.CooldownRemaining
	i.charges = state.Charges
	i.chargeRecovery = state.ChargeRecovery
}
