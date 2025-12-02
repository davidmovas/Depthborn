package currency

import (
	"fmt"
	"sync"

	"github.com/davidmovas/Depthborn/pkg/persist"
)

// Type defines different currency types
type Type string

const (
	TypeGold      Type = "gold"
	TypeGems      Type = "gems"      // Premium currency
	TypeDust      Type = "dust"      // Crafting currency
	TypeShards    Type = "shards"    // Upgrade currency
	TypeEssence   Type = "essence"   // Enchanting currency
	TypeFragments Type = "fragments" // Special currency
	TypeSouls     Type = "souls"     // Boss drops
	TypeTokens    Type = "tokens"    // Event/seasonal currency
)

// AllTypes returns all available currency types
func AllTypes() []Type {
	return []Type{
		TypeGold, TypeGems, TypeDust, TypeShards,
		TypeEssence, TypeFragments, TypeSouls, TypeTokens,
	}
}

// DisplayName returns human-readable currency name
func (t Type) DisplayName() string {
	names := map[Type]string{
		TypeGold:      "Gold",
		TypeGems:      "Gems",
		TypeDust:      "Dust",
		TypeShards:    "Shards",
		TypeEssence:   "Essence",
		TypeFragments: "Fragments",
		TypeSouls:     "Souls",
		TypeTokens:    "Tokens",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return string(t)
}

// Manager handles character currencies
type Manager interface {
	// Gold returns current gold amount
	Gold() int64

	// AddGold adds gold (can be negative for removal)
	AddGold(amount int64) error

	// SetGold sets gold directly
	SetGold(amount int64)

	// CanAffordGold checks if can afford gold amount
	CanAffordGold(amount int64) bool

	// Get returns amount of specific currency type
	Get(currencyType Type) int64

	// Add adds currency (can be negative)
	Add(currencyType Type, amount int64) error

	// Set sets currency amount directly
	Set(currencyType Type, amount int64)

	// CanAfford checks if can afford amount of currency
	CanAfford(currencyType Type, amount int64) bool

	// GetAll returns all currencies
	GetAll() map[Type]int64

	// Transfer transfers currency to another manager
	Transfer(currencyType Type, amount int64, target Manager) error

	// Reset resets all currencies to zero
	Reset()

	// TotalValue returns total value in gold equivalent
	TotalValue() int64

	// SerializeState converts state to map for persistence
	SerializeState() (map[string]any, error)

	// DeserializeState restores state from map
	DeserializeState(state map[string]any) error
}

var _ Manager = (*BaseManager)(nil)

// BaseManager implements Manager interface
type BaseManager struct {
	mu sync.RWMutex

	currencies    map[Type]int64
	exchangeRates map[Type]int64 // Exchange rate to gold (1 of currency = X gold)
}

// Config holds configuration for creating a currency manager
type Config struct {
	InitialGold   int64
	ExchangeRates map[Type]int64
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		InitialGold: 0,
		ExchangeRates: map[Type]int64{
			TypeGold:      1,
			TypeGems:      100,  // 1 gem = 100 gold
			TypeDust:      10,   // 1 dust = 10 gold
			TypeShards:    50,   // 1 shard = 50 gold
			TypeEssence:   25,   // 1 essence = 25 gold
			TypeFragments: 5,    // 1 fragment = 5 gold
			TypeSouls:     1000, // 1 soul = 1000 gold
			TypeTokens:    0,    // Non-convertible
		},
	}
}

// NewManager creates a new currency manager with default config
func NewManager() *BaseManager {
	return NewManagerWithConfig(DefaultConfig())
}

// NewManagerWithConfig creates a new currency manager with config
func NewManagerWithConfig(cfg Config) *BaseManager {
	m := &BaseManager{
		currencies:    make(map[Type]int64),
		exchangeRates: cfg.ExchangeRates,
	}

	if m.exchangeRates == nil {
		m.exchangeRates = DefaultConfig().ExchangeRates
	}

	if cfg.InitialGold > 0 {
		m.currencies[TypeGold] = cfg.InitialGold
	}

	return m
}

func (m *BaseManager) Gold() int64 {
	return m.Get(TypeGold)
}

func (m *BaseManager) AddGold(amount int64) error {
	return m.Add(TypeGold, amount)
}

func (m *BaseManager) SetGold(amount int64) {
	m.Set(TypeGold, amount)
}

func (m *BaseManager) CanAffordGold(amount int64) bool {
	return m.CanAfford(TypeGold, amount)
}

func (m *BaseManager) Get(currencyType Type) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currencies[currencyType]
}

func (m *BaseManager) Add(currencyType Type, amount int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	current := m.currencies[currencyType]
	newAmount := current + amount

	if newAmount < 0 {
		return fmt.Errorf("insufficient %s: have %d, need %d", currencyType.DisplayName(), current, -amount)
	}

	m.currencies[currencyType] = newAmount
	return nil
}

func (m *BaseManager) Set(currencyType Type, amount int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if amount < 0 {
		amount = 0
	}
	m.currencies[currencyType] = amount
}

func (m *BaseManager) CanAfford(currencyType Type, amount int64) bool {
	if amount <= 0 {
		return true
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currencies[currencyType] >= amount
}

func (m *BaseManager) GetAll() map[Type]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[Type]int64, len(m.currencies))
	for t, amount := range m.currencies {
		result[t] = amount
	}
	return result
}

func (m *BaseManager) Transfer(currencyType Type, amount int64, target Manager) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	if target == nil {
		return fmt.Errorf("target manager cannot be nil")
	}

	// Remove from source
	if err := m.Add(currencyType, -amount); err != nil {
		return fmt.Errorf("source: %w", err)
	}

	// Add to target
	if err := target.Add(currencyType, amount); err != nil {
		// Rollback
		_ = m.Add(currencyType, amount)
		return fmt.Errorf("target: %w", err)
	}

	return nil
}

func (m *BaseManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currencies = make(map[Type]int64)
}

func (m *BaseManager) TotalValue() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var total int64
	for currencyType, amount := range m.currencies {
		rate, ok := m.exchangeRates[currencyType]
		if ok && rate > 0 {
			total += amount * rate
		}
	}
	return total
}

// State holds serializable currency state
type State struct {
	Currencies map[string]int64 `msgpack:"currencies"`
}

func (m *BaseManager) SerializeState() (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	currencies := make(map[string]int64, len(m.currencies))
	for t, amount := range m.currencies {
		currencies[string(t)] = amount
	}

	state := State{Currencies: currencies}
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

func (m *BaseManager) DeserializeState(stateData map[string]any) error {
	data, err := persist.DefaultCodec().Encode(stateData)
	if err != nil {
		return err
	}

	var state State
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.currencies = make(map[Type]int64, len(state.Currencies))
	for t, amount := range state.Currencies {
		m.currencies[Type(t)] = amount
	}

	return nil
}

// GetExchangeRate returns the gold exchange rate for a currency type
func (m *BaseManager) GetExchangeRate(currencyType Type) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.exchangeRates[currencyType]
}

// SetExchangeRate sets the gold exchange rate for a currency type
func (m *BaseManager) SetExchangeRate(currencyType Type, rate int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.exchangeRates[currencyType] = rate
}

// ConvertToGold converts a currency to gold
func (m *BaseManager) ConvertToGold(currencyType Type, amount int64) error {
	if currencyType == TypeGold {
		return nil // Already gold
	}

	rate := m.GetExchangeRate(currencyType)
	if rate <= 0 {
		return fmt.Errorf("%s cannot be converted to gold", currencyType.DisplayName())
	}

	if err := m.Add(currencyType, -amount); err != nil {
		return err
	}

	goldAmount := amount * rate
	return m.AddGold(goldAmount)
}
