package skill

// BaseRequirements реализует интерфейс Requirements
type BaseRequirements struct {
	level      int
	attributes map[string]float64
	skills     map[string]int
	items      []string
}

func NewBaseRequirements(level int, attributes map[string]float64, skills map[string]int, items []string) *BaseRequirements {
	return &BaseRequirements{
		level:      level,
		attributes: attributes,
		skills:     skills,
		items:      items,
	}
}

func (r *BaseRequirements) Level() int {
	return r.level
}

func (r *BaseRequirements) Attributes() map[string]float64 {
	return r.attributes
}

func (r *BaseRequirements) Skills() map[string]int {
	return r.skills
}

func (r *BaseRequirements) Items() []string {
	return r.items
}

// Check проверяет, выполнены ли требования (простейший пример)
func (r *BaseRequirements) Check(entityID string) bool {
	// TODO: здесь проверять состояние персонажа/инвентарь и другие условия
	// Для примера всегда true
	return true
}
