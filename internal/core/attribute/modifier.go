package attribute

var _ Modifier = (*BaseModifier)(nil)

type BaseModifier struct {
	id       string
	modType  ModifierType
	value    float64
	source   string
	priority int
	active   bool
}

func NewModifier(id string, modType ModifierType, value float64, source string) Modifier {
	return &BaseModifier{
		id:       id,
		modType:  modType,
		value:    value,
		source:   source,
		priority: 0,
		active:   true,
	}
}

func NewModifierWithPriority(id string, modType ModifierType, value float64, source string, priority int) Modifier {
	return &BaseModifier{
		id:       id,
		modType:  modType,
		value:    value,
		source:   source,
		priority: priority,
		active:   true,
	}
}

func (m *BaseModifier) ID() string {
	return m.id
}

func (m *BaseModifier) Type() ModifierType {
	return m.modType
}

func (m *BaseModifier) Value() float64 {
	return m.value
}

func (m *BaseModifier) Source() string {
	return m.source
}

func (m *BaseModifier) Priority() int {
	return m.priority
}

func (m *BaseModifier) IsActive() bool {
	return m.active
}

func (m *BaseModifier) SetActive(active bool) {
	m.active = active
}

func (m *BaseModifier) SetValue(value float64) {
	m.value = value
}
