package character

import (
	"context"
	"fmt"

	"github.com/davidmovas/Depthborn/internal/character/currency"
	"github.com/davidmovas/Depthborn/internal/character/equipment"
	"github.com/davidmovas/Depthborn/internal/character/inventory"
	"github.com/davidmovas/Depthborn/internal/character/progression"
	"github.com/davidmovas/Depthborn/internal/character/statistics"
	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/internal/core/skill"
	"github.com/davidmovas/Depthborn/internal/core/status"
	"github.com/davidmovas/Depthborn/internal/core/types"
	"github.com/davidmovas/Depthborn/internal/world/spatial"
)

// CharacterBuilder implements Builder interface with fluent API
type CharacterBuilder struct {
	config Config
	errors []error
}

// NewBuilder creates a new character builder
func NewBuilder() *CharacterBuilder {
	return &CharacterBuilder{
		config: DefaultConfig(""),
	}
}

// WithName sets character name
func (b *CharacterBuilder) WithName(name string) *CharacterBuilder {
	b.config.Name = name
	return b
}

// WithAccountID sets account ID
func (b *CharacterBuilder) WithAccountID(accountID string) *CharacterBuilder {
	b.config.AccountID = accountID
	return b
}

// WithLevel sets starting level
func (b *CharacterBuilder) WithLevel(level int) *CharacterBuilder {
	if level < 1 {
		level = 1
	}
	// Level is managed by progression, we'll set it after creation
	return b
}

// WithHealth sets initial and max health
func (b *CharacterBuilder) WithHealth(initial, max float64) *CharacterBuilder {
	if max < 1 {
		max = 100
	}
	if initial < 0 {
		initial = max
	}
	if initial > max {
		initial = max
	}
	b.config.InitialHealth = initial
	b.config.MaxHealth = max
	return b
}

// WithMaxHealth sets max health (initial will be set to max)
func (b *CharacterBuilder) WithMaxHealth(max float64) *CharacterBuilder {
	return b.WithHealth(max, max)
}

// WithAttackRange sets attack range
func (b *CharacterBuilder) WithAttackRange(range_ float64) *CharacterBuilder {
	if range_ < 0 {
		range_ = 1.5
	}
	b.config.AttackRange = range_
	return b
}

// WithMaxWeight sets inventory weight capacity
func (b *CharacterBuilder) WithMaxWeight(weight float64) *CharacterBuilder {
	if weight < 0 {
		weight = 100
	}
	b.config.MaxWeight = weight
	return b
}

// WithGold sets initial gold
func (b *CharacterBuilder) WithGold(amount int64) *CharacterBuilder {
	if amount < 0 {
		amount = 0
	}
	b.config.InitialGold = amount
	return b
}

// WithAttributes sets base attributes
func (b *CharacterBuilder) WithAttributes(attrs map[string]float64) *CharacterBuilder {
	if b.config.Attributes == nil {
		b.config.Attributes = attribute.NewManager()
	}
	for k, v := range attrs {
		b.config.Attributes.SetBase(attribute.Type(k), v)
	}
	return b
}

// WithAttribute sets a single base attribute
func (b *CharacterBuilder) WithAttribute(attrType attribute.Type, value float64) *CharacterBuilder {
	if b.config.Attributes == nil {
		b.config.Attributes = attribute.NewManager()
	}
	b.config.Attributes.SetBase(attrType, value)
	return b
}

// WithStrength sets strength attribute
func (b *CharacterBuilder) WithStrength(value float64) *CharacterBuilder {
	return b.WithAttribute(attribute.AttrStrength, value)
}

// WithDexterity sets dexterity attribute
func (b *CharacterBuilder) WithDexterity(value float64) *CharacterBuilder {
	return b.WithAttribute(attribute.AttrDexterity, value)
}

// WithIntelligence sets intelligence attribute
func (b *CharacterBuilder) WithIntelligence(value float64) *CharacterBuilder {
	return b.WithAttribute(attribute.AttrIntelligence, value)
}

// WithVitality sets vitality attribute
func (b *CharacterBuilder) WithVitality(value float64) *CharacterBuilder {
	return b.WithAttribute(attribute.AttrVitality, value)
}

// WithWillpower sets willpower attribute
func (b *CharacterBuilder) WithWillpower(value float64) *CharacterBuilder {
	return b.WithAttribute(attribute.AttrWillpower, value)
}

// WithPosition sets starting position
func (b *CharacterBuilder) WithPosition(x, y, z int) *CharacterBuilder {
	pos := spatial.NewPosition(x, y, z)
	b.config.Transform = spatial.NewTransform(pos, 0)
	return b
}

// WithTransform sets full transform
func (b *CharacterBuilder) WithTransform(transform spatial.Transform) *CharacterBuilder {
	b.config.Transform = transform
	return b
}

// WithTags sets initial tags
func (b *CharacterBuilder) WithTags(tags ...string) *CharacterBuilder {
	if b.config.Tags == nil {
		b.config.Tags = types.NewTagSet()
	}
	for _, tag := range tags {
		b.config.Tags.Add(tag)
	}
	return b
}

// WithTag adds a single tag
func (b *CharacterBuilder) WithTag(tag string) *CharacterBuilder {
	return b.WithTags(tag)
}

// WithCallbacks sets callback registry
func (b *CharacterBuilder) WithCallbacks(callbacks types.CallbackRegistry) *CharacterBuilder {
	b.config.Callbacks = callbacks
	return b
}

// WithStatusManager sets status effect manager
func (b *CharacterBuilder) WithStatusManager(statuses status.Manager) *CharacterBuilder {
	b.config.Statuses = statuses
	return b
}

// WithProgression sets progression manager
func (b *CharacterBuilder) WithProgression(prog progression.Manager) *CharacterBuilder {
	b.config.Progression = prog
	return b
}

// WithProgressionConfig creates progression from config
func (b *CharacterBuilder) WithProgressionConfig(cfg progression.ManagerConfig) *CharacterBuilder {
	b.config.Progression = progression.NewManager(cfg)
	return b
}

// WithEquipment sets equipment manager
func (b *CharacterBuilder) WithEquipment(equip equipment.Manager) *CharacterBuilder {
	b.config.Equipment = equip
	return b
}

// WithInventory sets inventory manager
func (b *CharacterBuilder) WithInventory(inv inventory.Manager) *CharacterBuilder {
	b.config.Inventory = inv
	return b
}

// WithInventoryConfig creates inventory from config
func (b *CharacterBuilder) WithInventoryConfig(cfg inventory.Config) *CharacterBuilder {
	b.config.Inventory = inventory.NewManagerWithConfig(cfg)
	return b
}

// WithCurrency sets currency manager
func (b *CharacterBuilder) WithCurrency(curr currency.Manager) *CharacterBuilder {
	b.config.Currency = curr
	return b
}

// WithCurrencyConfig creates currency manager from config
func (b *CharacterBuilder) WithCurrencyConfig(cfg currency.Config) *CharacterBuilder {
	b.config.Currency = currency.NewManagerWithConfig(cfg)
	return b
}

// WithStatistics sets statistics tracker
func (b *CharacterBuilder) WithStatistics(stats statistics.Tracker) *CharacterBuilder {
	b.config.Statistics = stats
	return b
}

// WithSkillTree sets skill tree
func (b *CharacterBuilder) WithSkillTree(tree skill.Tree) *CharacterBuilder {
	b.config.SkillTree = tree
	return b
}

// WithSkillLoadout sets skill loadout
func (b *CharacterBuilder) WithSkillLoadout(loadout skill.Loadout) *CharacterBuilder {
	b.config.SkillLoadout = loadout
	return b
}

// Build creates the character
func (b *CharacterBuilder) Build(_ context.Context) (*BaseCharacter, error) {
	// Validate
	if b.config.Name == "" {
		return nil, fmt.Errorf("character name is required")
	}

	// Check for accumulated errors
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	// Create character
	char := NewCharacter(b.config)

	// Validate character
	if err := char.Validate(); err != nil {
		return nil, fmt.Errorf("character validation failed: %w", err)
	}

	return char, nil
}

// MustBuild creates the character or panics on error
func (b *CharacterBuilder) MustBuild(ctx context.Context) *BaseCharacter {
	char, err := b.Build(ctx)
	if err != nil {
		panic(err)
	}
	return char
}

// Clone creates a copy of the builder
func (b *CharacterBuilder) Clone() *CharacterBuilder {
	return &CharacterBuilder{
		config: b.config,
		errors: append([]error{}, b.errors...),
	}
}

// Reset resets the builder to default state
func (b *CharacterBuilder) Reset() *CharacterBuilder {
	b.config = DefaultConfig("")
	b.errors = nil
	return b
}

// --- Preset builders for common character types ---

// Warrior creates a warrior preset builder
func Warrior(name string) *CharacterBuilder {
	return NewBuilder().
		WithName(name).
		WithMaxHealth(150).
		WithStrength(15).
		WithDexterity(8).
		WithIntelligence(5).
		WithVitality(12).
		WithWillpower(5).
		WithAttackRange(1.5).
		WithTag("warrior")
}

// Ranger creates a ranger preset builder
func Ranger(name string) *CharacterBuilder {
	return NewBuilder().
		WithName(name).
		WithMaxHealth(100).
		WithStrength(8).
		WithDexterity(15).
		WithIntelligence(8).
		WithVitality(8).
		WithWillpower(6).
		WithAttackRange(8.0).
		WithTag("ranger")
}

// Mage creates a mage preset builder
func Mage(name string) *CharacterBuilder {
	return NewBuilder().
		WithName(name).
		WithMaxHealth(80).
		WithStrength(5).
		WithDexterity(8).
		WithIntelligence(15).
		WithVitality(6).
		WithWillpower(12).
		WithAttackRange(10.0).
		WithTag("mage")
}

// Balanced creates a balanced preset builder
func Balanced(name string) *CharacterBuilder {
	return NewBuilder().
		WithName(name).
		WithMaxHealth(100).
		WithStrength(10).
		WithDexterity(10).
		WithIntelligence(10).
		WithVitality(10).
		WithWillpower(10).
		WithAttackRange(1.5).
		WithTag("balanced")
}
