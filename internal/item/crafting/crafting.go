package crafting

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/item"
)

// Crafter performs crafting operations on items
type Crafter interface {
	// Craft creates new item from recipe
	Craft(ctx context.Context, recipe Recipe, materials []item.Item) (item.Item, error)

	// CanCraft checks if recipe can be executed with given materials
	CanCraft(recipe Recipe, materials []item.Item) bool

	// Upgrade improves item quality or level
	Upgrade(ctx context.Context, target item.Item, materials []item.Item) (item.Item, error)

	// CanUpgrade checks if item can be upgraded
	CanUpgrade(target item.Item, materials []item.Item) bool

	// Reforge rerolls item affixes
	Reforge(ctx context.Context, target item.Equipment, materials []item.Item) (item.Equipment, error)

	// CanReforge checks if item can be reforged
	CanReforge(target item.Equipment, materials []item.Item) bool

	// Enchant adds or modifies enchantments
	Enchant(ctx context.Context, target item.Equipment, enchant Enchantment, materials []item.Item) (item.Equipment, error)

	// CanEnchant checks if enchantment can be applied
	CanEnchant(target item.Equipment, enchant Enchantment, materials []item.Item) bool

	// Socket adds socket to item
	Socket(ctx context.Context, target item.Equipment, materials []item.Item) (item.Equipment, error)

	// CanSocket checks if socket can be added
	CanSocket(target item.Equipment, materials []item.Item) bool

	// Repair restores item durability
	Repair(ctx context.Context, target item.Equipment, materials []item.Item) (item.Equipment, error)

	// CanRepair checks if item can be repaired
	CanRepair(target item.Equipment, materials []item.Item) bool

	// Dismantle breaks item into materials
	Dismantle(ctx context.Context, target item.Item) ([]item.Item, error)

	// CanDismantle checks if item can be dismantled
	CanDismantle(target item.Item) bool
}

// Recipe defines crafting blueprint
type Recipe interface {
	// ID returns unique recipe identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Category returns recipe category
	Category() RecipeCategory

	// Requirements returns crafting requirements
	Requirements() Requirements

	// Materials returns required materials
	Materials() []MaterialRequirement

	// Result returns what recipe produces
	Result() RecipeResult

	// SuccessChance returns probability of success [0.0 - 1.0]
	SuccessChance() float64

	// CraftingTime returns time to craft in milliseconds
	CraftingTime() int64

	// RequiredStation returns crafting station needed
	RequiredStation() Station

	// UnlockConditions returns conditions to unlock recipe
	UnlockConditions() []UnlockCondition

	// IsUnlocked checks if recipe is available
	IsUnlocked() bool

	// Unlock makes recipe available
	Unlock()

	// Tags returns recipe tags
	Tags() []string
}

// RecipeCategory groups related recipes
type RecipeCategory string

const (
	CategoryWeapon     RecipeCategory = "weapon"
	CategoryArmor      RecipeCategory = "armor"
	CategoryAccessory  RecipeCategory = "accessory"
	CategoryConsumable RecipeCategory = "consumable"
	CategoryMaterial   RecipeCategory = "material"
	CategoryUpgrade    RecipeCategory = "upgrade"
	CategoryEnchant    RecipeCategory = "enchant"
	CategorySocket     RecipeCategory = "socket"
	CategorySpecial    RecipeCategory = "special"
)

// Requirements defines recipe prerequisites
type Requirements interface {
	// Level returns minimum crafting level
	Level() int

	// Skills returns required skill levels
	Skills() map[string]int

	// Attributes returns minimum attributes
	Attributes() map[string]float64

	// Reputation returns required reputation levels
	Reputation() map[string]int

	// Check verifies entity meets requirements
	Check(entityID string) bool
}

// MaterialRequirement specifies needed material
type MaterialRequirement interface {
	// ItemType returns required item type
	ItemType() string

	// MinQuantity returns minimum amount needed
	MinQuantity() int

	// MaxQuantity returns maximum amount usable
	MaxQuantity() int

	// MinQuality returns minimum quality required
	MinQuality() float64

	// MinRarity returns minimum rarity required
	MinRarity() item.Rarity

	// IsConsumed returns true if material is consumed
	IsConsumed() bool

	// Alternatives returns alternative material types
	Alternatives() []string

	// Description returns human-readable requirement
	Description() string
}

// RecipeResult defines crafting output
type RecipeResult interface {
	// ItemType returns produced item type
	ItemType() string

	// MinQuantity returns minimum items produced
	MinQuantity() int

	// MaxQuantity returns maximum items produced
	MaxQuantity() int

	// QualityRange returns min and max quality
	QualityRange() (min, max float64)

	// RarityWeights returns weighted rarity chances
	RarityWeights() map[item.Rarity]float64

	// BonusResults returns possible bonus items
	BonusResults() []BonusResult

	// InheritQuality returns true if quality based on materials
	InheritQuality() bool
}

// BonusResult defines additional crafting output
type BonusResult interface {
	// ItemType returns bonus item type
	ItemType() string

	// Chance returns probability of bonus [0.0 - 1.0]
	Chance() float64

	// Quantity returns bonus item count
	Quantity() int
}

// Station represents crafting facility
type Station string

const (
	StationNone       Station = "none"
	StationWorkbench  Station = "workbench"
	StationForge      Station = "forge"
	StationAlchemy    Station = "alchemy"
	StationEnchanting Station = "enchanting"
	StationJewelcraft Station = "jewelcraft"
	StationRuneforge  Station = "runeforge"
	StationArtisan    Station = "artisan"
)

// UnlockCondition defines recipe unlock requirement
type UnlockCondition interface {
	// Type returns condition type
	Type() string

	// Check verifies if condition is met
	Check(ctx context.Context, entityID string) bool

	// Description returns human-readable condition
	Description() string
}

// Enchantment represents item enhancement
type Enchantment interface {
	// ID returns unique enchantment identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Tier returns enchantment tier (higher = stronger)
	Tier() int

	// Type returns enchantment type
	Type() EnchantmentType

	// Modifiers returns stat modifiers granted
	Modifiers() []any

	// AllowedSlots returns equipment slots that can be enchanted
	AllowedSlots() []item.EquipmentSlot

	// MaxPerItem returns maximum applications on single item
	MaxPerItem() int

	// ConflictsWith returns incompatible enchantment IDs
	ConflictsWith() []string

	// Requirements returns application requirements
	Requirements() Requirements

	// Apply applies enchantment to item
	Apply(ctx context.Context, target item.Equipment) error

	// Remove removes enchantment from item
	Remove(ctx context.Context, target item.Equipment) error

	// CanApply checks if can be applied to item
	CanApply(target item.Equipment) bool
}

// EnchantmentType categorizes enchantments
type EnchantmentType string

const (
	EnchantOffensive EnchantmentType = "offensive"
	EnchantDefensive EnchantmentType = "defensive"
	EnchantUtility   EnchantmentType = "utility"
	EnchantElemental EnchantmentType = "elemental"
	EnchantCursed    EnchantmentType = "cursed"
	EnchantBlessed   EnchantmentType = "blessed"
	EnchantCosmetic  EnchantmentType = "cosmetic"
)

// RecipeBook manages discovered recipes
type RecipeBook interface {
	// Add adds recipe to book
	Add(recipe Recipe)

	// Remove removes recipe from book
	Remove(recipeID string)

	// Get retrieves recipe by ID
	Get(recipeID string) (Recipe, bool)

	// GetAll returns all recipes
	GetAll() []Recipe

	// GetByCategory returns recipes in category
	GetByCategory(category RecipeCategory) []Recipe

	// GetByStation returns recipes for station
	GetByStation(station Station) []Recipe

	// GetCraftable returns recipes that can be crafted
	GetCraftable(entityID string, materials []item.Item) []Recipe

	// Has checks if recipe is known
	Has(recipeID string) bool

	// Count returns total recipes
	Count() int

	// Unlock unlocks recipe
	Unlock(recipeID string) error

	// IsUnlocked checks if recipe is unlocked
	IsUnlocked(recipeID string) bool
}

// Registry manages all available recipes
type Registry interface {
	// Register adds recipe to registry
	Register(recipe Recipe) error

	// Unregister removes recipe from registry
	Unregister(recipeID string) error

	// Get retrieves recipe by ID
	Get(recipeID string) (Recipe, bool)

	// GetAll returns all registered recipes
	GetAll() []Recipe

	// GetByCategory returns recipes in category
	GetByCategory(category RecipeCategory) []Recipe

	// GetByTag returns recipes with tag
	GetByTag(tag string) []Recipe

	// Has checks if recipe is registered
	Has(recipeID string) bool

	// Count returns total registered recipes
	Count() int

	// Search finds recipes matching criteria
	Search(criteria SearchCriteria) []Recipe
}

// SearchCriteria defines recipe filtering
type SearchCriteria struct {
	Categories   []RecipeCategory
	Stations     []Station
	MinLevel     int
	MaxLevel     int
	Tags         []string
	NameContains string
}

// Skill tracks crafting proficiency
type Skill interface {
	// Type returns skill type
	Type() Station

	// Level returns current level
	Level() int

	// Experience returns current experience
	Experience() int64

	// ExperienceToNextLevel returns XP needed for next level
	ExperienceToNextLevel() int64

	// AddExperience grants experience
	AddExperience(amount int64) bool

	// SuccessBonus returns success chance bonus
	SuccessBonus() float64

	// QualityBonus returns quality improvement bonus
	QualityBonus() float64

	// SpeedBonus returns crafting speed multiplier
	SpeedBonus() float64

	// Specializations returns unlocked specializations
	Specializations() []Specialization
}

// Specialization represents crafting mastery
type Specialization interface {
	// ID returns unique specialization identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Station returns associated crafting station
	Station() Station

	// RequiredLevel returns minimum skill level
	RequiredLevel() int

	// Bonuses returns specialization bonuses
	Bonuses() []SpecializationBonus

	// IsUnlocked checks if specialization is available
	IsUnlocked() bool

	// Unlock makes specialization available
	Unlock() error
}

// SpecializationBonus defines specialization benefit
type SpecializationBonus interface {
	// Type returns bonus type
	Type() string

	// Value returns bonus value
	Value() float64

	// Description returns human-readable description
	Description() string

	// Apply applies bonus to crafting operation
	Apply(ctx context.Context, operation Operation) error
}

// Operation represents ongoing craft
type Operation interface {
	// ID returns unique operation identifier
	ID() string

	// Recipe returns recipe being crafted
	Recipe() Recipe

	// Materials returns materials being used
	Materials() []item.Item

	// CrafterID returns entity performing craft
	CrafterID() string

	// StartTime returns when operation started
	StartTime() int64

	// EstimatedEndTime returns when operation should finish
	EstimatedEndTime() int64

	// Progress returns completion percentage [0.0 - 1.0]
	Progress() float64

	// IsComplete returns true if crafting finished
	IsComplete() bool

	// Cancel aborts operation and refunds materials
	Cancel(ctx context.Context) error

	// Complete finalizes operation and returns result
	Complete(ctx context.Context) ([]item.Item, error)

	// Update processes operation for elapsed time
	Update(ctx context.Context, deltaMs int64) error
}

// Queue manages multiple crafting operations
type Queue interface {
	// Add adds operation to queue
	Add(operation Operation) error

	// Remove removes operation from queue
	Remove(operationID string) error

	// Get retrieves operation by ID
	Get(operationID string) (Operation, bool)

	// GetAll returns all queued operations
	GetAll() []Operation

	// GetActive returns currently processing operations
	GetActive() []Operation

	// GetPending returns waiting operations
	GetPending() []Operation

	// GetCompleted returns finished operations
	GetCompleted() []Operation

	// Clear removes all operations
	Clear()

	// Update processes all operations
	Update(ctx context.Context, deltaMs int64) error

	// MaxConcurrent returns maximum simultaneous operations
	MaxConcurrent() int

	// SetMaxConcurrent updates operation limit
	SetMaxConcurrent(max int)
}
