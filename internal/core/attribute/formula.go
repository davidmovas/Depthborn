package attribute

var _ Formula = (*SimpleFormula)(nil)

type SimpleFormula struct {
	dependencies []Type
	calculator   func(Manager) float64
}

func NewFormula(calculator func(Manager) float64, dependencies ...Type) Formula {
	return &SimpleFormula{
		dependencies: dependencies,
		calculator:   calculator,
	}
}

func (f *SimpleFormula) Calculate(manager Manager) float64 {
	return f.calculator(manager)
}

func (f *SimpleFormula) Dependencies() []Type {
	return f.dependencies
}

type MultiplyFormula struct {
	sourceAttr Type
	multiplier float64
}

func NewMultiplyFormula(sourceAttr Type, multiplier float64) Formula {
	return &MultiplyFormula{
		sourceAttr: sourceAttr,
		multiplier: multiplier,
	}
}

func (f *MultiplyFormula) Calculate(manager Manager) float64 {
	return manager.Get(f.sourceAttr) * f.multiplier
}

func (f *MultiplyFormula) Dependencies() []Type {
	return []Type{f.sourceAttr}
}

type SumFormula struct {
	attributes []Type
	weights    []float64
	baseValue  float64
}

func NewSumFormula(baseValue float64, attributes []Type, weights []float64) Formula {
	return &SumFormula{
		attributes: attributes,
		weights:    weights,
		baseValue:  baseValue,
	}
}

func (f *SumFormula) Calculate(manager Manager) float64 {
	result := f.baseValue

	for i, attr := range f.attributes {
		weight := 1.0
		if i < len(f.weights) {
			weight = f.weights[i]
		}
		result += manager.Get(attr) * weight
	}

	return result
}

func (f *SumFormula) Dependencies() []Type {
	return f.attributes
}
