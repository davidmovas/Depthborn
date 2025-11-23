package skill

import (
	"os"

	"gopkg.in/yaml.v3"
)

type NodeYAML struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Type         NodeType `yaml:"type"`
	Cost         int      `yaml:"cost"`
	Requirements []string `yaml:"requirements"`
	Connections  []string `yaml:"connections"`
	EffectID     string   `yaml:"effect_id"`
	Position     struct {
		X float64 `yaml:"x"`
		Y float64 `yaml:"y"`
	} `yaml:"position"`
}

type SkillYAML struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Type         Type     `yaml:"type"`
	MaxLevel     int      `yaml:"max_level"`
	ManaCost     float64  `yaml:"mana_cost"`
	Cooldown     int64    `yaml:"cooldown"`
	Tags         []string `yaml:"tags"`
	Requirements []string `yaml:"requirements"` // ID нужных скиллов
}

// TreeYAML — структура всего дерева
type TreeYAML struct {
	AvailablePoints int         `yaml:"available_points"`
	Nodes           []NodeYAML  `yaml:"nodes"`
	Skills          []SkillYAML `yaml:"skills"`
}

func SaveTreeToYAML(tree *TreeImpl, path string) error {
	yamlTree := TreeYAML{
		AvailablePoints: tree.AvailablePoints(),
	}

	// Конвертация нод
	for _, node := range tree.GetNodes() {
		x, y := node.Position()
		yamlNode := NodeYAML{
			ID:           node.ID(),
			Name:         node.Name(),
			Description:  node.Description(),
			Type:         node.SkillType(),
			Cost:         node.Cost(),
			Requirements: node.Requirements(),
			Connections:  node.Connections(),
			EffectID:     node.Effect().Description(), // можно заменить на ID эффекта
			Position: struct {
				X float64 `yaml:"x"`
				Y float64 `yaml:"y"`
			}(struct{ X, Y float64 }{X: x, Y: y}),
		}
		yamlTree.Nodes = append(yamlTree.Nodes, yamlNode)
	}

	// Конвертация скиллов
	for _, skill := range tree.GetAllSkills() {
		yamlSkill := SkillYAML{
			ID:          skill.ID(),
			Name:        skill.Name(),
			Description: skill.Description(),
			Type:        skill.SkillType(),
			MaxLevel:    skill.MaxLevel(),
			ManaCost:    skill.ManaCost(),
			Cooldown:    skill.MaxCooldown(),
			Tags:        skill.Tags(),
			// Можно добавить требования и зависимости
		}
		yamlTree.Skills = append(yamlTree.Skills, yamlSkill)
	}

	data, err := yaml.Marshal(&yamlTree)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadTreeFromYAML(path string) (*TreeImpl, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var yamlTree TreeYAML
	if err = yaml.Unmarshal(data, &yamlTree); err != nil {
		return nil, err
	}

	tree := NewTree(0) // или передать доступные очки
	tree.AddPoints(yamlTree.AvailablePoints)

	// Создаём ноды
	for _, n := range yamlTree.Nodes {
		node := NewNode(
			n.ID, n.Name, n.Description, n.Type, n.Cost,
			n.Requirements, n.Connections,
			nil, // эффект позже, можно через EffectRegistry
			n.Position.X, n.Position.Y,
		)
		tree.AddNode(node)
	}

	// Создаём скиллы
	for _, s := range yamlTree.Skills {
		skill := NewSkill(
			s.ID, s.Name, s.Description, s.Type, s.MaxLevel,
			s.ManaCost, s.Cooldown, s.Tags, nil, nil,
		)
		tree.AddSkill(skill)
	}

	return tree, nil
}
