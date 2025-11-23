package camp

import (
	"context"

	"github.com/davidmovas/Depthborn/internal/item"
)

// Camp represents safe hub area
type Camp interface {
	// ID returns unique camp identifier
	ID() string

	// Name returns camp name
	Name() string

	// Description returns camp description
	Description() string

	// Facilities returns available facilities
	Facilities() []Facility

	// GetFacility retrieves facility by ID
	GetFacility(facilityID string) (Facility, bool)

	// AddFacility adds facility to camp
	AddFacility(facility Facility) error

	// RemoveFacility removes facility from camp
	RemoveFacility(facilityID string) error

	// IsUnlocked returns true if camp is accessible
	IsUnlocked() bool

	// Unlock makes camp accessible
	Unlock() error

	// UpgradeLevel returns current upgrade level
	UpgradeLevel() int

	// Upgrade improves camp level
	Upgrade(ctx context.Context) error

	// CanUpgrade checks if upgrade is possible
	CanUpgrade() bool

	// UpgradeCost returns resources needed for upgrade
	UpgradeCost() []ResourceCost

	// OnEnter is called when player enters camp
	OnEnter(ctx context.Context, characterID string) error

	// OnExit is called when player leaves camp
	OnExit(ctx context.Context, characterID string) error
}

// Facility represents camp building or service
type Facility interface {
	// ID returns unique facility identifier
	ID() string

	// Name returns display name
	Name() string

	// Description returns detailed description
	Description() string

	// Type returns facility type
	Type() FacilityType

	// Level returns current facility level
	Level() int

	// MaxLevel returns maximum facility level
	MaxLevel() int

	// Upgrade improves facility level
	Upgrade(ctx context.Context) error

	// CanUpgrade checks if upgrade is possible
	CanUpgrade() bool

	// UpgradeCost returns resources needed for upgrade
	UpgradeCost() []ResourceCost

	// IsUnlocked returns true if facility is accessible
	IsUnlocked() bool

	// Unlock makes facility accessible
	Unlock() error

	// UnlockRequirements returns requirements to unlock
	UnlockRequirements() []UnlockRequirement

	// Interact handles player interaction
	Interact(ctx context.Context, characterID string) error

	// CanInteract checks if interaction is possible
	CanInteract(characterID string) bool

	// Icon returns icon identifier
	Icon() string
}

// FacilityType categorizes facilities
type FacilityType string

const (
	FacilityForge      FacilityType = "forge"
	FacilityAlchemy    FacilityType = "alchemy"
	FacilityEnchanting FacilityType = "enchanting"
	FacilityVendor     FacilityType = "vendor"
	FacilityStorage    FacilityType = "storage"
	FacilityTraining   FacilityType = "training"
	FacilityBestiary   FacilityType = "bestiary"
	FacilityLibrary    FacilityType = "library"
	FacilityGuild      FacilityType = "guild"
	FacilityTavern     FacilityType = "tavern"
)

// ResourceCost defines resource requirement
type ResourceCost struct {
	ItemType string
	Quantity int
	Currency int64
}

// UnlockRequirement defines unlock condition
type UnlockRequirement interface {
	// Type returns requirement type
	Type() string

	// Check verifies if requirement is met
	Check(ctx context.Context, characterID string) bool

	// Description returns human-readable requirement
	Description() string
}

// Forge handles equipment crafting and repair
type Forge interface {
	Facility

	// Craft creates item from recipe
	Craft(ctx context.Context, characterID string, recipeID string, materials []item.Item) (item.Item, error)

	// Repair restores item durability
	Repair(ctx context.Context, characterID string, itemID string) error

	// UpgradeItem improves item quality
	UpgradeItem(ctx context.Context, characterID string, itemID string, materials []item.Item) (item.Item, error)

	// Socket adds socket to item
	Socket(ctx context.Context, characterID string, itemID string) error

	// GetRepairCost calculates repair cost
	GetRepairCost(itemID string) int64

	// GetUpgradeCost calculates upgrade cost
	GetUpgradeCost(itemID string) []ResourceCost

	// AvailableRecipes returns craftable recipes
	AvailableRecipes(characterID string) []string
}

// Vendor buys and sells items
type Vendor interface {
	Facility

	// Buy purchases item from vendor
	Buy(ctx context.Context, characterID string, itemID string, quantity int) error

	// Sell sells item to vendor
	Sell(ctx context.Context, characterID string, itemID string, quantity int) error

	// GetBuyPrice calculates purchase cost
	GetBuyPrice(itemID string, quantity int) int64

	// GetSellPrice calculates sale value
	GetSellPrice(itemID string, quantity int) int64

	// Inventory returns vendor inventory
	Inventory() VendorInventory

	// Refresh restocks vendor inventory
	Refresh(ctx context.Context) error

	// RefreshInterval returns time between restocks
	RefreshInterval() int64

	// LastRefresh returns timestamp of last restock
	LastRefresh() int64

	// Reputation returns vendor reputation level
	Reputation(characterID string) int

	// AddReputation increases reputation
	AddReputation(characterID string, amount int)

	// ReputationDiscount returns price discount based on reputation
	ReputationDiscount(characterID string) float64
}

// VendorInventory manages vendor items
type VendorInventory interface {
	// GetItems returns all available items
	GetItems() []VendorItem

	// GetItem retrieves item by ID
	GetItem(itemID string) (VendorItem, bool)

	// AddItem adds item to inventory
	AddItem(item VendorItem)

	// RemoveItem removes item from inventory
	RemoveItem(itemID string)

	// GetByCategory returns items in category
	GetByCategory(category string) []VendorItem

	// GetByRarity returns items of rarity
	GetByRarity(rarity item.Rarity) []VendorItem

	// Stock returns quantity available for item
	Stock(itemID string) int

	// SetStock updates item quantity
	SetStock(itemID string, quantity int)

	// DecrementStock reduces item quantity
	DecrementStock(itemID string, amount int) bool

	// IsInStock checks if item is available
	IsInStock(itemID string) bool
}

// VendorItem represents item for sale
type VendorItem interface {
	// Item returns base item
	Item() item.Item

	// Price returns sale price
	Price() int64

	// Stock returns available quantity (-1 = unlimited)
	Stock() int

	// RequiredReputation returns minimum reputation to purchase
	RequiredReputation() int

	// IsLimited returns true if stock is limited
	IsLimited() bool
}

// Storage provides item storage
type Storage interface {
	Facility

	// Deposit stores item in storage
	Deposit(ctx context.Context, characterID string, itemID string) error

	// Withdraw retrieves item from storage
	Withdraw(ctx context.Context, characterID string, itemID string) error

	// GetItems returns all stored items for character
	GetItems(characterID string) []item.Item

	// GetSharedItems returns items accessible by all characters
	GetSharedItems() []item.Item

	// Capacity returns maximum storage slots
	Capacity() int

	// UsedCapacity returns occupied slots
	UsedCapacity(characterID string) int

	// ExpandCapacity increases storage size
	ExpandCapacity(ctx context.Context, additionalSlots int) error

	// ExpansionCost calculates cost to expand
	ExpansionCost(additionalSlots int) int64
}

// Training provides character development
type Training interface {
	Facility

	// RespecSkills resets skill allocations
	RespecSkills(ctx context.Context, characterID string) error

	// RespecCost calculates respec cost
	RespecCost(characterID string) int64

	// TrainSkill improves skill proficiency
	TrainSkill(ctx context.Context, characterID string, skillID string) error

	// GetTrainingCost calculates skill training cost
	GetTrainingCost(skillID string, currentLevel int) int64

	// AvailableTrainings returns trainings for character
	AvailableTrainings(characterID string) []Training

	// UnlockClass makes class available
	UnlockClass(ctx context.Context, characterID string, classID string) error

	// GetClassCost calculates class unlock cost
	GetClassCost(classID string) int64
}

// Bestiary tracks discovered enemies
type Bestiary interface {
	Facility

	// RegisterKill records enemy defeat
	RegisterKill(ctx context.Context, characterID string, enemyType string) error

	// GetEntry retrieves bestiary entry
	GetEntry(enemyType string) (BestiaryEntry, bool)

	// GetEntries returns all discovered entries
	GetEntries() []BestiaryEntry

	// GetProgress returns discovery percentage
	GetProgress() float64

	// IsDiscovered checks if enemy is known
	IsDiscovered(enemyType string) bool

	// GetKillCount returns kills for enemy type
	GetKillCount(enemyType string) int

	// GetRewards returns completion rewards
	GetRewards() []CompletionReward

	// ClaimReward claims completion reward
	ClaimReward(ctx context.Context, characterID string, rewardID string) error
}

// BestiaryEntry contains enemy information
type BestiaryEntry interface {
	// EnemyType returns enemy type identifier
	EnemyType() string

	// Name returns enemy name
	Name() string

	// Description returns enemy description
	Description() string

	// IsDiscovered returns true if entry is revealed
	IsDiscovered() bool

	// DiscoveryLevel returns information reveal level
	DiscoveryLevel() int

	// KillCount returns times enemy was defeated
	KillCount() int

	// Weaknesses returns enemy weaknesses
	Weaknesses() []string

	// Resistances returns enemy resistances
	Resistances() []string

	// LootTable returns possible drops
	LootTable() []string

	// Habitat returns where enemy spawns
	Habitat() []string

	// Abilities returns enemy abilities
	Abilities() []string
}

// CompletionReward represents bestiary milestone reward
type CompletionReward interface {
	// ID returns unique reward identifier
	ID() string

	// Name returns reward name
	Name() string

	// Description returns reward description
	Description() string

	// Requirement returns completion requirement
	Requirement() string

	// IsClaimed returns true if reward was claimed
	IsClaimed() bool

	// Rewards returns reward items
	Rewards() []item.Item
}

// Library manages lore and recipes
type Library interface {
	Facility

	// ReadBook unlocks lore entry
	ReadBook(ctx context.Context, characterID string, bookID string) error

	// GetBooks returns available books
	GetBooks() []Book

	// GetBook retrieves book by ID
	GetBook(bookID string) (Book, bool)

	// IsRead checks if book was read
	IsRead(bookID string) bool

	// LearnRecipe unlocks crafting recipe
	LearnRecipe(ctx context.Context, characterID string, recipeID string) error

	// GetRecipes returns available recipes
	GetRecipes() []string

	// HasRecipe checks if recipe is known
	HasRecipe(recipeID string) bool
}

// Book represents lore document
type Book interface {
	// ID returns unique book identifier
	ID() string

	// Title returns book title
	Title() string

	// Author returns book author
	Author() string

	// Content returns book text
	Content() string

	// Category returns book category
	Category() string

	// RequiredLevel returns minimum level to read
	RequiredLevel() int

	// UnlockCost returns cost to unlock
	UnlockCost() int64
}

// Manager manages camp lifecycle
type Manager interface {
	// GetCamp retrieves camp by ID
	GetCamp(campID string) (Camp, bool)

	// GetCamps returns all camps
	GetCamps() []Camp

	// CurrentCamp returns active camp
	CurrentCamp() (Camp, bool)

	// SetCurrentCamp changes active camp
	SetCurrentCamp(campID string) error

	// Save persists camp state
	Save(ctx context.Context, camp Camp) error

	// Load loads camp from storage
	Load(ctx context.Context, campID string) (Camp, error)
}
