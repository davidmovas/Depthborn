package progression

import (
	"fmt"
	"math"
)

var _ ExperienceCurve = (*BaseCurve)(nil)

type BaseCurve struct {
	curveType  CurveType
	baseXP     int64
	parameters map[string]float64
	formula    CurveFormula
}

func (c *BaseCurve) ExperienceForLevel(level int) int64 {
	if level <= 1 {
		return 0
	}

	switch c.curveType {
	case CurveLinear:
		return c.baseXP * int64(level-1)

	case CurveExponential:
		multiplier := c.parameters["multiplier"]
		total := int64(0)
		for i := 2; i <= level; i++ {
			xpForLevel := float64(c.baseXP) * math.Pow(multiplier, float64(i-2))
			total += int64(xpForLevel)
		}
		return total

	case CurvePolynomial:
		power := c.parameters["power"]
		return int64(float64(c.baseXP) * math.Pow(float64(level-1), power))

	case CurveLogarithmic:
		scale := c.parameters["scale"]
		return int64(float64(c.baseXP) * scale * math.Log(float64(level)))

	case CurveCustom:
		if c.formula != nil {
			return c.formula.Calculate(level)
		}
		return 0

	default:
		return 0
	}
}

func (c *BaseCurve) ExperienceToNextLevel(currentLevel int) int64 {
	return c.ExperienceForLevel(currentLevel+1) - c.ExperienceForLevel(currentLevel)
}

func (c *BaseCurve) LevelForExperience(xp int64) int {
	if xp <= 0 {
		return 1
	}

	level := 1
	for c.ExperienceForLevel(level+1) <= xp {
		level++
		// Safety limit
		if level > 1000 {
			break
		}
	}

	return level
}

func (c *BaseCurve) Type() CurveType {
	return c.curveType
}

func (c *BaseCurve) Parameters() map[string]float64 {
	// Return copy to prevent modification
	params := make(map[string]float64, len(c.parameters))
	for k, v := range c.parameters {
		params[k] = v
	}
	return params
}

// CurveBuilder implementation
var _ CurveBuilder = (*BaseCurveBuilder)(nil)

type BaseCurveBuilder struct{}

func NewCurveBuilder() CurveBuilder {
	return &BaseCurveBuilder{}
}

func (b *BaseCurveBuilder) Linear(baseXP int64) ExperienceCurve {
	return &BaseCurve{
		curveType:  CurveLinear,
		baseXP:     baseXP,
		parameters: map[string]float64{},
	}
}

func (b *BaseCurveBuilder) Exponential(baseXP int64, multiplier float64) ExperienceCurve {
	if multiplier <= 1.0 {
		multiplier = 1.5 // Default safe multiplier
	}

	return &BaseCurve{
		curveType: CurveExponential,
		baseXP:    baseXP,
		parameters: map[string]float64{
			"multiplier": multiplier,
		},
	}
}

func (b *BaseCurveBuilder) Logarithmic(baseXP int64, scale float64) ExperienceCurve {
	if scale <= 0 {
		scale = 1.0
	}

	return &BaseCurve{
		curveType: CurveLogarithmic,
		baseXP:    baseXP,
		parameters: map[string]float64{
			"scale": scale,
		},
	}
}

func (b *BaseCurveBuilder) Polynomial(baseXP int64, power float64) ExperienceCurve {
	if power <= 0 {
		power = 2.0 // Quadratic by default
	}

	return &BaseCurve{
		curveType: CurvePolynomial,
		baseXP:    baseXP,
		parameters: map[string]float64{
			"power": power,
		},
	}
}

func (b *BaseCurveBuilder) Custom(formula CurveFormula) ExperienceCurve {
	return &BaseCurve{
		curveType:  CurveCustom,
		baseXP:     0,
		parameters: map[string]float64{},
		formula:    formula,
	}
}

type CustomFormula struct {
	calculateFunc func(level int) int64
	description   string
}

func NewCustomFormula(fn func(level int) int64, description string) CurveFormula {
	return &CustomFormula{
		calculateFunc: fn,
		description:   description,
	}
}

func (f *CustomFormula) Calculate(level int) int64 {
	if f.calculateFunc == nil {
		return 0
	}
	return f.calculateFunc(level)
}

func (f *CustomFormula) Description() string {
	return f.description
}

func NewStandardExponentialCurve() ExperienceCurve {
	// Exponential curve similar to classic RPGs
	// Level 2: 100 XP, Level 3: 150 XP, Level 4: 225 XP...
	return NewCurveBuilder().Exponential(100, 1.5)
}

func NewStandardLinearCurve() ExperienceCurve {
	// Linear curve for faster progression
	return NewCurveBuilder().Linear(50)
}

func NewStandardPolynomialCurve() ExperienceCurve {
	// Polynomial curve for slow, grindy progression
	return NewCurveBuilder().Polynomial(50, 2.5)
}

func NewStandardLogarithmicCurve() ExperienceCurve { return NewCurveBuilder().Logarithmic(100, 0.8) }

func ValidateCurve(curve ExperienceCurve, maxLevel int) error {
	if maxLevel <= 0 {
		return fmt.Errorf("max level must be positive")
	}

	// Check that XP requirements are monotonically increasing
	prevXP := int64(0)
	for level := 1; level <= maxLevel; level++ {
		xp := curve.ExperienceForLevel(level)
		if xp < prevXP {
			return fmt.Errorf("XP curve must be monotonically increasing (level %d: %d < %d)", level, xp, prevXP)
		}
		prevXP = xp
	}

	return nil
}
