package entity

var _ AttributeModifier = (*BaseAttributeModifier)(nil)

type BaseAttributeModifier struct {
	id     string
	value  float64
	typ    string
	source string
}

func NewBaseAttributeModifier(id string, value float64, typ string, source string) AttributeModifier {
	return &BaseAttributeModifier{
		id:     id,
		value:  value,
		typ:    typ,
		source: source,
	}
}

func (bam *BaseAttributeModifier) ID() string {
	return bam.id
}

func (bam *BaseAttributeModifier) Value() float64 {
	return bam.value
}

func (bam *BaseAttributeModifier) Type() string {
	return bam.typ
}

func (bam *BaseAttributeModifier) Source() string {
	return bam.source
}
