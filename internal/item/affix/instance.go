package affix

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"github.com/davidmovas/Depthborn/pkg/identifier"
)

var _ Instance = (*BaseInstance)(nil)

// BaseInstance represents a rolled affix on an actual item.
// Contains concrete values generated from Affix template.
type BaseInstance struct {
	mu           sync.RWMutex
	id           string       // Unique instance ID
	affixID      string       // Source template ID
	affix        Affix        // Reference to source template (may be nil)
	affixType    Type         // Cached type
	group        string       // Cached group
	rolledValues []RolledModifier
}

// NewBaseInstance creates instance from affix template with rolled values
func NewBaseInstance(affix Affix, values []RolledModifier) *BaseInstance {
	return &BaseInstance{
		id:           identifier.New(),
		affixID:      affix.ID(),
		affix:        affix,
		affixType:    affix.Type(),
		group:        affix.Group(),
		rolledValues: values,
	}
}

// NewBaseInstanceFromData creates instance from serialized data (no affix reference)
func NewBaseInstanceFromData(affixID string, affixType Type, group string, values []RolledModifier) *BaseInstance {
	return &BaseInstance{
		id:           identifier.New(),
		affixID:      affixID,
		affix:        nil,
		affixType:    affixType,
		group:        group,
		rolledValues: values,
	}
}

func (bi *BaseInstance) AffixID() string {
	bi.mu.RLock()
	defer bi.mu.RUnlock()
	return bi.affixID
}

func (bi *BaseInstance) Affix() Affix {
	bi.mu.RLock()
	defer bi.mu.RUnlock()
	return bi.affix
}

func (bi *BaseInstance) Type() Type {
	bi.mu.RLock()
	defer bi.mu.RUnlock()
	return bi.affixType
}

func (bi *BaseInstance) Group() string {
	bi.mu.RLock()
	defer bi.mu.RUnlock()
	return bi.group
}

func (bi *BaseInstance) RolledValues() []RolledModifier {
	bi.mu.RLock()
	defer bi.mu.RUnlock()
	result := make([]RolledModifier, len(bi.rolledValues))
	copy(result, bi.rolledValues)
	return result
}

func (bi *BaseInstance) Modifiers() []attribute.Modifier {
	bi.mu.RLock()
	defer bi.mu.RUnlock()

	modifiers := make([]attribute.Modifier, 0, len(bi.rolledValues))
	for i, rm := range bi.rolledValues {
		modID := fmt.Sprintf("%s_%d", bi.id, i)
		mod := attribute.NewModifierWithPriority(
			modID,
			rm.Template.ModType,
			rm.Value,
			bi.affixID,
			rm.Template.Priority,
		)
		modifiers = append(modifiers, mod)
	}
	return modifiers
}

func (bi *BaseInstance) Reroll() {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	for i := range bi.rolledValues {
		bi.rolledValues[i].Value = rollValue(
			bi.rolledValues[i].Template.MinValue,
			bi.rolledValues[i].Template.MaxValue,
		)
	}
}

func (bi *BaseInstance) RerollSingle(index int) error {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	if index < 0 || index >= len(bi.rolledValues) {
		return fmt.Errorf("index out of range: %d", index)
	}

	bi.rolledValues[index].Value = rollValue(
		bi.rolledValues[index].Template.MinValue,
		bi.rolledValues[index].Template.MaxValue,
	)
	return nil
}

func (bi *BaseInstance) Quality() float64 {
	bi.mu.RLock()
	defer bi.mu.RUnlock()

	if len(bi.rolledValues) == 0 {
		return 0.5
	}

	totalQuality := 0.0
	for _, rm := range bi.rolledValues {
		totalQuality += calculateQuality(rm.Value, rm.Template.MinValue, rm.Template.MaxValue)
	}
	return totalQuality / float64(len(bi.rolledValues))
}

// SetAffix links instance to affix template (for deserialization)
func (bi *BaseInstance) SetAffix(affix Affix) {
	bi.mu.Lock()
	defer bi.mu.Unlock()
	bi.affix = affix
}

// rollValue generates random value between min and max using weighted distribution
// Values closer to center are more likely (bell curve approximation)
func rollValue(min, max float64) float64 {
	if min >= max {
		return min
	}

	// Use sum of multiple uniform randoms for bell curve approximation
	// More iterations = more bell-shaped distribution
	sum := 0.0
	iterations := 3
	for i := 0; i < iterations; i++ {
		sum += rand.Float64()
	}
	normalized := sum / float64(iterations)

	return min + (max-min)*normalized
}

// rollValueBiased generates value with bias toward min (0.0) or max (1.0)
func rollValueBiased(min, max, bias float64) float64 {
	if min >= max {
		return min
	}

	// Bias adjusts the distribution
	// bias = 0.0 -> skew toward min
	// bias = 0.5 -> uniform
	// bias = 1.0 -> skew toward max
	r := rand.Float64()

	// Apply power function for bias
	// bias < 0.5 -> power > 1 -> more weight to lower values
	// bias > 0.5 -> power < 1 -> more weight to higher values
	if bias < 0.5 {
		power := 1.0 + (0.5-bias)*4.0 // 0.0 bias -> power 3.0
		r = 1.0 - pow(1.0-r, power)
	} else if bias > 0.5 {
		power := 1.0 + (bias-0.5)*4.0 // 1.0 bias -> power 3.0
		r = pow(r, power)
	}

	return min + (max-min)*r
}

// calculateQuality returns how good a roll is [0.0 - 1.0]
func calculateQuality(value, min, max float64) float64 {
	if max <= min {
		return 1.0
	}
	quality := (value - min) / (max - min)
	if quality < 0 {
		quality = 0
	}
	if quality > 1 {
		quality = 1
	}
	return quality
}

// pow calculates x^y without importing math
func pow(x, y float64) float64 {
	if y == 0 {
		return 1
	}
	if y == 1 {
		return x
	}
	// For simple cases we handle manually
	// For complex cases, use iterative approximation
	if y == 2 {
		return x * x
	}
	if y == 3 {
		return x * x * x
	}

	// Iterative exponentiation for non-integer powers
	// Using exp(y * ln(x)) approximation
	result := 1.0
	base := x
	exp := y

	// Handle fractional exponents approximately
	for exp >= 1 {
		result *= base
		exp--
	}

	// Handle remaining fractional part with linear interpolation
	// (rough approximation, but sufficient for our purposes)
	if exp > 0 {
		result *= 1.0 + exp*(base-1.0)
	}

	return result
}

// RollModifiers generates rolled values from templates using weighted distribution
func RollModifiers(templates []ModifierTemplate) []RolledModifier {
	result := make([]RolledModifier, len(templates))
	for i, tmpl := range templates {
		result[i] = RolledModifier{
			Template: tmpl,
			Value:    rollValue(tmpl.MinValue, tmpl.MaxValue),
		}
	}
	return result
}

// RollModifiersBiased generates rolled values with quality bias
func RollModifiersBiased(templates []ModifierTemplate, bias float64) []RolledModifier {
	result := make([]RolledModifier, len(templates))
	for i, tmpl := range templates {
		result[i] = RolledModifier{
			Template: tmpl,
			Value:    rollValueBiased(tmpl.MinValue, tmpl.MaxValue, bias),
		}
	}
	return result
}
