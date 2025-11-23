package progression

import (
	"context"
)

// ExperienceManager manages character experience and leveling
type ExperienceManager interface {
	// CurrentLevel returns current level
	CurrentLevel() int

	// CurrentExperience returns current XP
	CurrentExperience() int64

	// ExperienceToNextLevel returns XP needed for next level
	ExperienceToNextLevel() int64

	// ExperienceForLevel returns total XP needed to reach level
	ExperienceForLevel(level int) int64

	// AddExperience grants experience points
	AddExperience(ctx context.Context, amount int64) (levelsGained int, err error)

	// SetExperience sets experience directly
	SetExperience(xp int64) error

	// SetLevel sets level directly
	SetLevel(level int) error

	// Progress returns level progress percentage [0.0 - 1.0]
	Progress() float64

	// MaxLevel returns maximum achievable level
	MaxLevel() int

	// IsMaxLevel returns true if at max level
	IsMaxLevel() bool

	// OnLevelUp registers callback when level increases
	OnLevelUp(callback LevelUpCallback)

	// OnExperienceGain registers callback when XP is gained
	OnExperienceGain(callback ExperienceCallback)
}

// LevelUpCallback is invoked when character levels up
type LevelUpCallback func(ctx context.Context, oldLevel, newLevel int, manager ExperienceManager)

// ExperienceCallback is invoked when experience is gained
type ExperienceCallback func(ctx context.Context, amount int64, manager ExperienceManager)

// ExperienceCurve calculates XP requirements
type ExperienceCurve interface {
	// ExperienceForLevel returns total XP needed to reach level
	ExperienceForLevel(level int) int64

	// ExperienceToNextLevel returns XP needed from current level to next
	ExperienceToNextLevel(currentLevel int) int64

	// LevelForExperience returns level for given XP amount
	LevelForExperience(xp int64) int

	// Type returns curve type
	Type() CurveType

	// Parameters returns curve configuration
	Parameters() map[string]float64
}

// CurveType categorizes XP curves
type CurveType string

const (
	CurveLinear      CurveType = "linear"      // Constant XP per level
	CurveExponential CurveType = "exponential" // XP grows exponentially
	CurveLogarithmic CurveType = "logarithmic" // XP grows logarithmically
	CurvePolynomial  CurveType = "polynomial"  // XP = base * level^power
	CurveCustom      CurveType = "custom"      // Custom formula
)

// CurveBuilder creates experience curves
type CurveBuilder interface {
	// Linear creates linear curve (constant XP per level)
	Linear(baseXP int64) ExperienceCurve

	// Exponential creates exponential curve
	Exponential(baseXP int64, multiplier float64) ExperienceCurve

	// Logarithmic creates logarithmic curve
	Logarithmic(baseXP int64, scale float64) ExperienceCurve

	// Polynomial creates polynomial curve
	Polynomial(baseXP int64, power float64) ExperienceCurve

	// Custom creates curve with custom formula
	Custom(formula CurveFormula) ExperienceCurve
}

// CurveFormula defines custom XP calculation
type CurveFormula interface {
	// Calculate computes XP for level
	Calculate(level int) int64

	// Description returns formula description
	Description() string
}

// AttributeGrowth manages attribute increases on level up
type AttributeGrowth interface {
	// GetGrowth returns attribute growth per level
	GetGrowth(attributeType string) float64

	// SetGrowth updates attribute growth
	SetGrowth(attributeType string, growth float64)

	// ApplyLevelUp applies attribute increases for level
	ApplyLevelUp(ctx context.Context, level int, attributes map[string]float64) error

	// GetAttributesAtLevel calculates attributes at specific level
	GetAttributesAtLevel(level int, baseAttributes map[string]float64) map[string]float64

	// GrowthType returns growth calculation type
	GrowthType() GrowthType

	// SetGrowthType updates growth type
	SetGrowthType(growthType GrowthType)
}

// GrowthType categorizes attribute growth
type GrowthType string

const (
	GrowthFlat        GrowthType = "flat"        // Fixed amount per level
	GrowthPercentage  GrowthType = "percentage"  // Percentage of base per level
	GrowthScaling     GrowthType = "scaling"     // Scales with level
	GrowthDiminishing GrowthType = "diminishing" // Diminishing returns
)

// StatPointManager manages allocatable stat points
type StatPointManager interface {
	// AvailablePoints returns unspent stat points
	AvailablePoints() int

	// SpentPoints returns allocated stat points
	SpentPoints() int

	// AddPoints grants stat points
	AddPoints(amount int)

	// AllocatePoint spends point on attribute
	AllocatePoint(attributeType string) error

	// DeallocatePoint refunds point from attribute
	DeallocatePoint(attributeType string) error

	// GetAllocations returns points spent per attribute
	GetAllocations() map[string]int

	// Reset refunds all allocated points
	Reset() error

	// CanAllocate checks if can spend point
	CanAllocate(attributeType string) bool

	// PointsPerLevel returns points gained per level
	PointsPerLevel() int

	// SetPointsPerLevel updates points per level
	SetPointsPerLevel(points int)

	// OnPointAllocated registers callback when point is spent
	OnPointAllocated(callback StatPointCallback)

	// OnPointDeallocated registers callback when point is refunded
	OnPointDeallocated(callback StatPointCallback)
}

// StatPointCallback is invoked for stat point events
type StatPointCallback func(ctx context.Context, attributeType string, manager StatPointManager)

// SkillPointManager manages skill point allocation
type SkillPointManager interface {
	// AvailablePoints returns unspent skill points
	AvailablePoints() int

	// SpentPoints returns allocated skill points
	SpentPoints() int

	// TotalPoints returns total points ever gained
	TotalPoints() int

	// AddPoints grants skill points
	AddPoints(amount int)

	// SpendPoints consumes skill points
	SpendPoints(amount int) error

	// RefundPoints returns skill points
	RefundPoints(amount int)

	// PointsForLevel returns skill points gained at level
	PointsForLevel(level int) int

	// Reset refunds all points
	Reset()

	// OnPointsGained registers callback when points are granted
	OnPointsGained(callback SkillPointCallback)

	// OnPointsSpent registers callback when points are spent
	OnPointsSpent(callback SkillPointCallback)
}

// SkillPointCallback is invoked for skill point events
type SkillPointCallback func(ctx context.Context, amount int, manager SkillPointManager)

// LevelReward represents reward for reaching level
type LevelReward interface {
	// Level returns level this reward is for
	Level() int

	// Type returns reward type
	Type() RewardType

	// Amount returns reward amount
	Amount() int64

	// ItemID returns item identifier (if item reward)
	ItemID() string

	// SkillID returns skill identifier (if skill reward)
	SkillID() string

	// Description returns reward description
	Description() string

	// Apply applies reward to character
	Apply(ctx context.Context, characterID string) error

	// CanClaim checks if reward can be claimed
	CanClaim(ctx context.Context, characterID string) bool

	// IsClaimed returns true if reward was claimed
	IsClaimed() bool

	// SetClaimed marks reward as claimed
	SetClaimed(claimed bool)
}

// RewardType categorizes level rewards
type RewardType string

const (
	RewardStatPoints  RewardType = "stat_points"
	RewardSkillPoints RewardType = "skill_points"
	RewardGold        RewardType = "gold"
	RewardItem        RewardType = "item"
	RewardSkill       RewardType = "skill"
	RewardAbility     RewardType = "ability"
	RewardTitle       RewardType = "title"
	RewardAccess      RewardType = "access"
)

// RewardTable manages level rewards
type RewardTable interface {
	// GetRewards returns rewards for level
	GetRewards(level int) []LevelReward

	// AddReward adds reward for level
	AddReward(level int, reward LevelReward)

	// RemoveReward removes reward from level
	RemoveReward(level int, rewardType RewardType)

	// HasRewards checks if level has rewards
	HasRewards(level int) bool

	// ClaimRewards claims all rewards for level
	ClaimRewards(ctx context.Context, level int, characterID string) error

	// GetAllRewards returns all configured rewards
	GetAllRewards() map[int][]LevelReward

	// Clear removes all rewards
	Clear()
}

// Tracker tracks progression milestones
type Tracker interface {
	// RecordLevelUp records level increase
	RecordLevelUp(level int, timestamp int64)

	// GetLevelHistory returns level up history
	GetLevelHistory() []LevelUpRecord

	// TimeToLevel returns time taken to reach level
	TimeToLevel(level int) int64

	// AverageTimePerLevel returns average leveling time
	AverageTimePerLevel() int64

	// FastestLevel returns quickest level up
	FastestLevel() (level int, duration int64)

	// SlowestLevel returns slowest level up
	SlowestLevel() (level int, duration int64)

	// TotalPlayTime returns total play time in seconds
	TotalPlayTime() int64

	// PlayTimeAtLevel returns play time when level was reached
	PlayTimeAtLevel(level int) int64

	// ExperienceGained returns total XP gained
	ExperienceGained() int64

	// DeathCount returns total deaths
	DeathCount() int

	// Reset clears all tracking data
	Reset()
}

// LevelUpRecord records level up event
type LevelUpRecord struct {
	Level     int
	Timestamp int64
	PlayTime  int64
	Deaths    int
}

// PrestigeManager manages prestige/rebirth system
type PrestigeManager interface {
	// CurrentPrestige returns prestige level
	CurrentPrestige() int

	// CanPrestige checks if can prestige
	CanPrestige() bool

	// Prestige performs prestige/rebirth
	Prestige(ctx context.Context) error

	// PrestigeRequirements returns requirements to prestige
	PrestigeRequirements() PrestigeRequirements

	// PrestigeBonuses returns bonuses from prestige
	PrestigeBonuses() []PrestigeBonus

	// TotalPrestiges returns lifetime prestige count
	TotalPrestiges() int

	// OnPrestige registers callback when prestige occurs
	OnPrestige(callback PrestigeCallback)
}

// PrestigeRequirements defines prestige conditions
type PrestigeRequirements interface {
	// MinimumLevel returns minimum level required
	MinimumLevel() int

	// RequiredItems returns items needed
	RequiredItems() map[string]int

	// RequiredAchievements returns achievements needed
	RequiredAchievements() []string

	// Check verifies if requirements are met
	Check(ctx context.Context, characterID string) bool

	// Description returns human-readable requirements
	Description() string
}

// PrestigeBonus represents benefit from prestiging
type PrestigeBonus interface {
	// Type returns bonus type
	Type() PrestigeBonusType

	// Value returns bonus value
	Value() float64

	// Description returns bonus description
	Description() string

	// Apply applies bonus to character
	Apply(ctx context.Context, characterID string) error
}

// PrestigeBonusType categorizes prestige bonuses
type PrestigeBonusType string

const (
	PrestigeBonusExperience    PrestigeBonusType = "experience"
	PrestigeBonusGold          PrestigeBonusType = "gold"
	PrestigeBonusLoot          PrestigeBonusType = "loot"
	PrestigeBonusStats         PrestigeBonusType = "stats"
	PrestigeBonusSkillPoints   PrestigeBonusType = "skill_points"
	PrestigeBonusStartingLevel PrestigeBonusType = "starting_level"
	PrestigeBonusUnique        PrestigeBonusType = "unique"
)

// PrestigeCallback is invoked when prestige occurs
type PrestigeCallback func(ctx context.Context, oldPrestige, newPrestige int, manager PrestigeManager)

// DifficultyScaler scales game difficulty with level
type DifficultyScaler interface {
	// GetEnemyScaling returns enemy stat multipliers for level
	GetEnemyScaling(playerLevel, enemyBaseLevel int) EnemyScaling

	// GetLootScaling returns loot quality multiplier for level
	GetLootScaling(playerLevel, contentLevel int) float64

	// GetExperienceScaling returns XP multiplier for level difference
	GetExperienceScaling(playerLevel, enemyLevel int) float64

	// GetDamageScaling returns damage multiplier for level difference
	GetDamageScaling(attackerLevel, defenderLevel int) float64

	// ScalingType returns scaling formula type
	ScalingType() ScalingType

	// SetScalingType updates scaling formula
	SetScalingType(scalingType ScalingType)
}

// EnemyScaling contains enemy stat multipliers
type EnemyScaling struct {
	HealthMultiplier  float64
	DamageMultiplier  float64
	DefenseMultiplier float64
	XPMultiplier      float64
	LootMultiplier    float64
}

// ScalingType categorizes difficulty scaling
type ScalingType string

const (
	ScalingNone        ScalingType = "none"        // No scaling
	ScalingLinear      ScalingType = "linear"      // Linear scaling
	ScalingExponential ScalingType = "exponential" // Exponential scaling
	ScalingDynamic     ScalingType = "dynamic"     // Adapts to player
	ScalingParty       ScalingType = "party"       // Scales with party size
)

// Manager coordinates all progression systems
type Manager interface {
	// Experience returns experience manager
	Experience() ExperienceManager

	// StatPoints returns stat point manager
	StatPoints() StatPointManager

	// SkillPoints returns skill point manager
	SkillPoints() SkillPointManager

	// AttributeGrowth returns attribute growth system
	AttributeGrowth() AttributeGrowth

	// RewardTable returns level rewards
	RewardTable() RewardTable

	// Tracker returns progression tracker
	Tracker() Tracker

	// Prestige returns prestige manager
	Prestige() PrestigeManager

	// DifficultyScaler returns difficulty scaler
	DifficultyScaler() DifficultyScaler

	// ProcessLevelUp handles level up
	ProcessLevelUp(ctx context.Context, characterID string, newLevel int) error

	// Reset resets all progression
	Reset(ctx context.Context) error

	// Save persists progression state
	Save(ctx context.Context) error

	// Load loads progression state
	Load(ctx context.Context) error
}
