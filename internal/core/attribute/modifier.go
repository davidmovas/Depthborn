package attribute

var _ Modifier = (*modifier)(nil)

type modifier struct {
	id       string
	modType  ModifierType
	value    float64
	source   string
	priority int
	active   bool
}

func NewModifier(id string, modType ModifierType, value float64, source string, priority int) Modifier {
	return &modifier{
		id:       id,
		modType:  modType,
		value:    value,
		source:   source,
		priority: priority,
		active:   true,
	}
}

func (m *modifier) ID() string {
	return m.id
}

func (m *modifier) Type() ModifierType {
	return m.modType
}

func (m *modifier) Value() float64 {
	return m.value
}

func (m *modifier) Source() string {
	return m.source
}

func (m *modifier) Priority() int {
	return m.priority
}

func (m *modifier) IsActive() bool {
	return m.active
}

func (m *modifier) SetActive(active bool) {
	m.active = active
}
