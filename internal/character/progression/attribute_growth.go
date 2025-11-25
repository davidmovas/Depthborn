package progression

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/davidmovas/Depthborn/pkg/state"
)

var _ AttributeGrowth = (*BaseAttributeGrowth)(nil)

type BaseAttributeGrowth struct {
	growthRates map[string]float64
	growthType  GrowthType

	mu sync.RWMutex
}

type AttributeGrowthConfig struct {
	GrowthRates map[string]float64
	GrowthType  GrowthType
}

func NewAttributeGrowth(config AttributeGrowthConfig) *BaseAttributeGrowth {
	growthRates := config.GrowthRates
	if growthRates == nil {
		growthRates = make(map[string]float64)
	}

	growthType := config.GrowthType
	if growthType == "" {
		growthType = GrowthFlat // Default to flat growth
	}

	return &BaseAttributeGrowth{
		growthRates: growthRates,
		growthType:  growthType,
	}
}

func (g *BaseAttributeGrowth) GetGrowth(attributeType string) float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.growthRates[attributeType]
}

func (g *BaseAttributeGrowth) SetGrowth(attributeType string, growth float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if growth == 0 {
		delete(g.growthRates, attributeType)
	} else {
		g.growthRates[attributeType] = growth
	}
}

func (g *BaseAttributeGrowth) ApplyLevelUp(ctx context.Context, level int, attributes map[string]float64) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for attrType, growthRate := range g.growthRates {
		if growthRate == 0 {
			continue
		}

		currentValue := attributes[attrType]
		newValue := g.calculateGrowth(currentValue, growthRate, level)
		attributes[attrType] = newValue
	}

	return nil
}

func (g *BaseAttributeGrowth) GetAttributesAtLevel(level int, baseAttributes map[string]float64) map[string]float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := make(map[string]float64)

	// Copy base attributes
	for attr, value := range baseAttributes {
		result[attr] = value
	}

	// Apply growth for each level
	for currentLevel := 2; currentLevel <= level; currentLevel++ {
		for attrType, growthRate := range g.growthRates {
			if growthRate == 0 {
				continue
			}

			currentValue := result[attrType]
			result[attrType] = g.calculateGrowth(currentValue, growthRate, currentLevel)
		}
	}

	return result
}

func (g *BaseAttributeGrowth) GrowthType() GrowthType {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.growthType
}

func (g *BaseAttributeGrowth) SetGrowthType(growthType GrowthType) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.growthType = growthType
}

func (g *BaseAttributeGrowth) Serialize() map[string]any {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Copy growth rates
	rates := make(map[string]float64, len(g.growthRates))
	for attr, rate := range g.growthRates {
		rates[attr] = rate
	}

	return map[string]any{
		"growth_rates": rates,
		"growth_type":  string(g.growthType),
	}
}

func (g *BaseAttributeGrowth) Deserialize(data map[string]any) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if ratesData, ok := data["growth_rates"].(map[string]any); ok {
		g.growthRates = make(map[string]float64)
		for attr, value := range ratesData {
			if rate, ok := value.(float64); ok {
				g.growthRates[attr] = rate
			}
		}
	}

	if growthType, ok := data["growth_type"].(string); ok {
		g.growthType = GrowthType(growthType)
	}

	return nil
}

func (g *BaseAttributeGrowth) SerializeState() (map[string]any, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	rates := make(map[string]float64, len(g.growthRates))
	for attr, rate := range g.growthRates {
		rates[attr] = rate
	}

	return state.BatchKV(
		"growth_rates", rates,
		"growth_type", string(g.growthType),
	)
}

func (g *BaseAttributeGrowth) DeserializeState(stateData map[string]any) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	s := state.From(stateData)

	var ok bool
	if g.growthRates, ok = state.MapTyped[string, float64]("growth_rates", s); !ok {
		return fmt.Errorf("failed to deserialize growth rates")
	}

	g.growthType = GrowthType(s.StringOr("growth_type", string(GrowthFlat)))

	return nil
}

// Helper method to calculate growth based on type
func (g *BaseAttributeGrowth) calculateGrowth(currentValue, growthRate float64, level int) float64 {
	switch g.growthType {
	case GrowthFlat:
		// Fixed amount per level
		return currentValue + growthRate

	case GrowthPercentage:
		// Percentage of base per level
		return currentValue * (1 + growthRate/100)

	case GrowthScaling:
		// Scales with level (growthRate * level)
		return currentValue + (growthRate * float64(level))

	case GrowthDiminishing:
		// Diminishing returns using logarithmic scaling
		scaleFactor := math.Log(float64(level)+1) / math.Log(2)
		return currentValue + (growthRate * scaleFactor)

	default:
		return currentValue + growthRate
	}
}
