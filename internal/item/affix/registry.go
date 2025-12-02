package affix

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/davidmovas/Depthborn/internal/core/attribute"
	"gopkg.in/yaml.v3"
)

var _ Registry = (*BaseRegistry)(nil)

// BaseRegistry is the default implementation of Registry interface
type BaseRegistry struct {
	mu      sync.RWMutex
	affixes map[string]Affix
	pools   map[string]Pool // key = "itemType:slot" or "itemType:*" or "*:*"
}

// NewBaseRegistry creates new empty registry
func NewBaseRegistry() *BaseRegistry {
	return &BaseRegistry{
		affixes: make(map[string]Affix),
		pools:   make(map[string]Pool),
	}
}

func (br *BaseRegistry) Register(affix Affix) error {
	br.mu.Lock()
	defer br.mu.Unlock()

	if _, exists := br.affixes[affix.ID()]; exists {
		return fmt.Errorf("affix already registered: %s", affix.ID())
	}

	br.affixes[affix.ID()] = affix
	return nil
}

func (br *BaseRegistry) Get(id string) (Affix, bool) {
	br.mu.RLock()
	defer br.mu.RUnlock()

	affix, exists := br.affixes[id]
	return affix, exists
}

func (br *BaseRegistry) GetAll() []Affix {
	br.mu.RLock()
	defer br.mu.RUnlock()

	result := make([]Affix, 0, len(br.affixes))
	for _, affix := range br.affixes {
		result = append(result, affix)
	}
	return result
}

func (br *BaseRegistry) GetPool(itemType string, slot string) Pool {
	br.mu.Lock()
	defer br.mu.Unlock()

	// Try specific pool first
	key := fmt.Sprintf("%s:%s", itemType, slot)
	if pool, exists := br.pools[key]; exists {
		return pool
	}

	// Try item type wildcard
	key = fmt.Sprintf("%s:*", itemType)
	if pool, exists := br.pools[key]; exists {
		return pool
	}

	// Try global pool
	if pool, exists := br.pools["*:*"]; exists {
		return pool
	}

	// Create new pool with all eligible affixes
	pool := br.buildPool(itemType, slot)
	br.pools[fmt.Sprintf("%s:%s", itemType, slot)] = pool
	return pool
}

func (br *BaseRegistry) buildPool(itemType string, slot string) Pool {
	pool := NewBasePool()

	for _, affix := range br.affixes {
		req := affix.Requirements()
		if req == nil {
			// No requirements = available everywhere
			pool.Add(affix)
			continue
		}

		// Check if affix is available for this item type/slot
		allowedTypes := req.AllowedTypes()
		allowedSlots := req.AllowedSlots()

		typeOK := len(allowedTypes) == 0
		for _, t := range allowedTypes {
			if t == itemType {
				typeOK = true
				break
			}
		}

		slotOK := len(allowedSlots) == 0
		for _, s := range allowedSlots {
			if s == slot {
				slotOK = true
				break
			}
		}

		if typeOK && slotOK {
			pool.Add(affix)
		}
	}

	return pool
}

func (br *BaseRegistry) LoadFromYAML(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var file File
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("failed to parse YAML %s: %w", path, err)
	}

	for _, def := range file.Affixes {
		var affix Affix
		affix, err = br.parseAffixDef(def)
		if err != nil {
			return fmt.Errorf("failed to parse affix %s: %w", def.ID, err)
		}

		if err = br.Register(affix); err != nil {
			return err
		}
	}

	// Clear cached pools since new affixes were added
	br.mu.Lock()
	br.pools = make(map[string]Pool)
	br.mu.Unlock()

	return nil
}

func (br *BaseRegistry) LoadFromDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		if err = br.LoadFromYAML(fullPath); err != nil {
			return fmt.Errorf("failed to load %s: %w", fullPath, err)
		}
	}

	return nil
}

func (br *BaseRegistry) parseAffixDef(def Def) (Affix, error) {
	// Parse modifiers
	modifiers := make([]ModifierTemplate, 0, len(def.Modifiers))
	for _, modDef := range def.Modifiers {
		mod := ModifierTemplate{
			Attribute: attribute.Type(modDef.Attribute),
			ModType:   attribute.ModifierType(modDef.ModType),
			MinValue:  modDef.Min,
			MaxValue:  modDef.Max,
			Priority:  modDef.Priority,
		}
		modifiers = append(modifiers, mod)
	}

	// Parse requirements
	var req Requirements
	if def.Requirements != nil {
		baseReq := NewBaseRequirements(def.Requirements.MinItemLevel)
		baseReq.SetMaxItemLevel(def.Requirements.MaxItemLevel)
		for _, t := range def.Requirements.AllowedTypes {
			baseReq.AddAllowedType(t)
		}
		for _, s := range def.Requirements.AllowedSlots {
			baseReq.AddAllowedSlot(s)
		}
		req = baseReq
	}

	affix := NewBaseAffixWithConfig(AffixConfig{
		ID:           def.ID,
		Name:         def.Name,
		Type:         Type(def.Type),
		Group:        def.Group,
		Rank:         def.Rank,
		Modifiers:    modifiers,
		Requirements: req,
		BaseWeight:   def.Weight,
		Description:  def.Description,
		Tags:         def.Tags,
	})

	return affix, nil
}

// YAML file structure definitions

// File represents a YAML file containing affix definitions
type File struct {
	Version string `yaml:"version"`
	Affixes []Def  `yaml:"affixes"`
}

// Def represents single affix definition in YAML
type Def struct {
	ID           string          `yaml:"id"`
	Name         string          `yaml:"name"`
	Type         string          `yaml:"type"` // prefix, suffix, implicit, etc.
	Group        string          `yaml:"group,omitempty"`
	Rank         int             `yaml:"rank"`
	Weight       int             `yaml:"weight"`
	Description  string          `yaml:"description,omitempty"`
	Tags         []string        `yaml:"tags,omitempty"`
	Modifiers    []ModifierDef   `yaml:"modifiers"`
	Requirements *RequirementDef `yaml:"requirements,omitempty"`
}

// ModifierDef represents modifier template in YAML
type ModifierDef struct {
	Attribute string  `yaml:"attribute"`
	ModType   string  `yaml:"mod_type"` // flat, increased, more
	Min       float64 `yaml:"min"`
	Max       float64 `yaml:"max"`
	Priority  int     `yaml:"priority,omitempty"`
}

// RequirementDef represents requirements in YAML
type RequirementDef struct {
	MinItemLevel int      `yaml:"min_level,omitempty"`
	MaxItemLevel int      `yaml:"max_level,omitempty"`
	AllowedTypes []string `yaml:"item_types,omitempty"`
	AllowedSlots []string `yaml:"slots,omitempty"`
}

// Global registry instance
var globalRegistry *BaseRegistry
var registryOnce sync.Once

// GlobalRegistry returns the global affix registry singleton
func GlobalRegistry() *BaseRegistry {
	registryOnce.Do(func() {
		globalRegistry = NewBaseRegistry()
	})
	return globalRegistry
}
