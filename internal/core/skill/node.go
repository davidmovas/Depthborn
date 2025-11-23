package skill

type NodeImpl struct {
	id           string
	name         string
	description  string
	skillType    NodeType
	cost         int
	requirements []string
	connections  []string
	effect       NodeEffect
	posX, posY   float64
}

func NewNode(
	id, name, description string,
	skillType NodeType,
	cost int,
	effect NodeEffect,
	posX, posY float64,
	requirements, connections []string,
) *NodeImpl {
	return &NodeImpl{
		id:           id,
		name:         name,
		description:  description,
		skillType:    skillType,
		cost:         cost,
		effect:       effect,
		posX:         posX,
		posY:         posY,
		requirements: requirements,
		connections:  connections,
	}
}

func (n *NodeImpl) ID() string {
	return n.id
}

func (n *NodeImpl) Type() string {
	return string(n.skillType)
}

func (n *NodeImpl) Name() string {
	return n.name
}

func (n *NodeImpl) SetName(name string) {
	n.name = name
}

func (n *NodeImpl) Description() string {
	return n.description
}

func (n *NodeImpl) SkillType() NodeType {
	return n.skillType
}

func (n *NodeImpl) Cost() int {
	return n.cost
}

func (n *NodeImpl) Requirements() []string {
	return n.requirements
}

func (n *NodeImpl) Connections() []string {
	return n.connections
}

func (n *NodeImpl) Effect() NodeEffect {
	return n.effect
}

func (n *NodeImpl) Position() (x, y float64) {
	return n.posX, n.posY
}
