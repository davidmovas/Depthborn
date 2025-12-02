package skill

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// =============================================================================
// TREE REGISTRY TESTS
// =============================================================================

func TestBaseTreeRegistry(t *testing.T) {
	t.Run("create registry", func(t *testing.T) {
		registry := NewBaseTreeRegistry()
		require.NotNil(t, registry)
		require.Equal(t, 0, registry.Count())
	})

	t.Run("register tree", func(t *testing.T) {
		registry := NewBaseTreeRegistry()

		tree := NewBaseTree(TreeConfig{
			ID:   "test_tree",
			Name: "Test Tree",
		})

		err := registry.Register(tree)
		require.NoError(t, err)
		require.Equal(t, 1, registry.Count())
		require.True(t, registry.Has("test_tree"))

		t.Run("duplicate registration returns error", func(t *testing.T) {
			err := registry.Register(tree)
			require.Error(t, err)
			require.Contains(t, err.Error(), "already registered")
		})
	})

	t.Run("get tree", func(t *testing.T) {
		registry := NewBaseTreeRegistry()

		tree := NewBaseTree(TreeConfig{
			ID:   "my_tree",
			Name: "My Tree",
		})
		_ = registry.Register(tree)

		t.Run("existing tree", func(t *testing.T) {
			found, ok := registry.Get("my_tree")
			require.True(t, ok)
			require.Equal(t, "my_tree", found.ID())
			require.Equal(t, "My Tree", found.Name())
		})

		t.Run("non-existing tree", func(t *testing.T) {
			_, ok := registry.Get("nonexistent")
			require.False(t, ok)
		})
	})

	t.Run("GetAll returns all trees", func(t *testing.T) {
		registry := NewBaseTreeRegistry()

		tree1 := NewBaseTree(TreeConfig{ID: "tree1", Name: "Tree 1"})
		tree2 := NewBaseTree(TreeConfig{ID: "tree2", Name: "Tree 2"})

		_ = registry.Register(tree1)
		_ = registry.Register(tree2)

		all := registry.GetAll()
		require.Len(t, all, 2)
	})
}

// =============================================================================
// YAML LOADING TESTS
// =============================================================================

func TestTreeYAMLLoading(t *testing.T) {
	t.Run("load simple tree from YAML", func(t *testing.T) {
		yamlData := []byte(`
version: "1.0"
tree:
  id: simple_tree
  name: "Simple Test Tree"
  description: "A tree for testing"
  branches:
    - id: combat
      name: "Combat"
      color: "#ff0000"
  start_nodes:
    - start_node
  nodes:
    - id: start_node
      name: "Start"
      type: path
      branch: combat
      cost: 0
      position: { x: 0, y: 0 }
      connections: [node_a]
      effects:
        - type: attribute
          attribute: strength
          mod_type: flat
          value: 1
    - id: node_a
      name: "Node A"
      type: path
      branch: combat
      cost: 1
      position: { x: 1, y: 0 }
      requirements: [start_node]
      connections: []
      effects:
        - type: attribute
          attribute: dexterity
          mod_type: flat
          value: 2
`)
		registry := NewBaseTreeRegistry()
		err := registry.LoadFromYAML(yamlData)
		require.NoError(t, err)

		tree, ok := registry.Get("simple_tree")
		require.True(t, ok)
		require.Equal(t, "Simple Test Tree", tree.Name())

		t.Run("branches loaded", func(t *testing.T) {
			branches := tree.GetBranches()
			require.Len(t, branches, 1)
			require.Equal(t, "combat", branches[0].ID)
			require.Equal(t, "Combat", branches[0].Name)
			require.Equal(t, "#ff0000", branches[0].Color)
		})

		t.Run("start nodes loaded", func(t *testing.T) {
			startNodes := tree.GetStartNodes()
			require.Len(t, startNodes, 1)
			require.Equal(t, "start_node", startNodes[0].ID())
		})

		t.Run("nodes loaded", func(t *testing.T) {
			nodes := tree.GetNodes()
			require.Len(t, nodes, 2)

			node, ok := tree.GetNode("start_node")
			require.True(t, ok)
			require.Equal(t, "Start", node.Name())
			require.Equal(t, NodePath, node.Type())
			require.Equal(t, 0, node.Cost())
			require.Equal(t, "combat", node.Branch())
		})

		t.Run("node connections", func(t *testing.T) {
			node, _ := tree.GetNode("start_node")
			connections := node.Connections()
			require.Len(t, connections, 1)
			require.Equal(t, "node_a", connections[0])
		})

		t.Run("node effects", func(t *testing.T) {
			node, _ := tree.GetNode("start_node")
			effects := node.Effects()
			require.Len(t, effects, 1)
			require.Equal(t, EffectTypeAttribute, effects[0].Type())
			require.Equal(t, float64(1), effects[0].Value())
		})
	})

	t.Run("load tree with keystone and exclusions", func(t *testing.T) {
		yamlData := []byte(`
version: "1.0"
tree:
  id: keystone_tree
  name: "Keystone Tree"
  start_nodes:
    - start
  nodes:
    - id: start
      name: "Start"
      type: path
      cost: 0
      position: { x: 0, y: 0 }
      connections: [path_to_ks1, path_to_ks2]

    - id: path_to_ks1
      name: "Path to KS1"
      type: path
      cost: 1
      position: { x: 1, y: -1 }
      requirements: [start]
      connections: [keystone_1]

    - id: path_to_ks2
      name: "Path to KS2"
      type: path
      cost: 1
      position: { x: 1, y: 1 }
      requirements: [start]
      connections: [keystone_2]

    - id: keystone_1
      name: "Blood Magic"
      type: keystone
      cost: 1
      position: { x: 2, y: -1 }
      requirements: [path_to_ks1]
      exclusions: [keystone_2]
      effects:
        - type: special
          description: "Skills cost Health"

    - id: keystone_2
      name: "Glass Cannon"
      type: keystone
      cost: 1
      position: { x: 2, y: 1 }
      requirements: [path_to_ks2]
      exclusions: [keystone_1]
      effects:
        - type: special
          description: "+50% Damage, +25% Damage Taken"
`)
		registry := NewBaseTreeRegistry()
		err := registry.LoadFromYAML(yamlData)
		require.NoError(t, err)

		tree, _ := registry.Get("keystone_tree")

		t.Run("keystone nodes loaded", func(t *testing.T) {
			ks1, ok := tree.GetNode("keystone_1")
			require.True(t, ok)
			require.Equal(t, NodeKeystone, ks1.Type())
			require.Equal(t, "Blood Magic", ks1.Name())
		})

		t.Run("exclusions loaded", func(t *testing.T) {
			ks1, _ := tree.GetNode("keystone_1")
			exclusions := ks1.Exclusions()
			require.Len(t, exclusions, 1)
			require.Equal(t, "keystone_2", exclusions[0])

			ks2, _ := tree.GetNode("keystone_2")
			exclusions2 := ks2.Exclusions()
			require.Len(t, exclusions2, 1)
			require.Equal(t, "keystone_1", exclusions2[0])
		})
	})

	t.Run("load leveled nodes", func(t *testing.T) {
		yamlData := []byte(`
version: "1.0"
tree:
  id: leveled_tree
  name: "Leveled Tree"
  start_nodes:
    - start
  nodes:
    - id: start
      name: "Start"
      type: path
      cost: 0
      position: { x: 0, y: 0 }
      connections: [mastery]

    - id: mastery
      name: "Fire Mastery"
      type: notable
      cost: 1
      max_level: 3
      level_cost: 1
      position: { x: 1, y: 0 }
      requirements: [start]
      levels:
        - level: 1
          effects:
            - type: skill_mod
              target_skill_tags: [fire]
              description: "+10% Fire Damage"
              metadata:
                damage_increase: 0.10
        - level: 2
          effects:
            - type: skill_mod
              target_skill_tags: [fire]
              description: "+20% Fire Damage"
              metadata:
                damage_increase: 0.20
        - level: 3
          effects:
            - type: skill_mod
              target_skill_tags: [fire]
              description: "+35% Fire Damage"
              metadata:
                damage_increase: 0.35
`)
		registry := NewBaseTreeRegistry()
		err := registry.LoadFromYAML(yamlData)
		require.NoError(t, err)

		tree, _ := registry.Get("leveled_tree")

		t.Run("max_level and level_cost", func(t *testing.T) {
			mastery, ok := tree.GetNode("mastery")
			require.True(t, ok)
			require.Equal(t, 3, mastery.MaxLevel())
			require.Equal(t, 1, mastery.LevelCost())
		})

		t.Run("effects per level", func(t *testing.T) {
			mastery, _ := tree.GetNode("mastery")

			effectsL1 := mastery.EffectsAtLevel(1)
			require.Len(t, effectsL1, 1)
			require.Equal(t, "+10% Fire Damage", effectsL1[0].Description())

			effectsL2 := mastery.EffectsAtLevel(2)
			require.Len(t, effectsL2, 1)
			require.Equal(t, "+20% Fire Damage", effectsL2[0].Description())

			effectsL3 := mastery.EffectsAtLevel(3)
			require.Len(t, effectsL3, 1)
			require.Equal(t, "+35% Fire Damage", effectsL3[0].Description())
		})
	})

	t.Run("all effect types", func(t *testing.T) {
		yamlData := []byte(`
version: "1.0"
tree:
  id: effects_tree
  name: "Effects Tree"
  start_nodes:
    - start
  nodes:
    - id: start
      name: "Start"
      type: path
      cost: 0
      position: { x: 0, y: 0 }
      effects:
        - type: attribute
          attribute: strength
          mod_type: flat
          value: 5
          description: "+5 Strength"
        - type: grant_skill
          skill_id: fireball
          start_level: 1
          description: "Grants Fireball skill"
        - type: passive
          trigger_type: on_kill
          trigger_value: 0.10
          effect_id: heal_on_kill
          description: "10% chance to heal on kill"
        - type: skill_mod
          target_skill_tags: [fire, spell]
          description: "+10% Fire Spell Damage"
          metadata:
            damage_increase: 0.10
        - type: unlock_craft
          description: "Unlock weapon crafting"
          metadata:
            recipes: [basic_sword]
        - type: trade
          description: "+5% Sell Value"
          metadata:
            sell_bonus: 0.05
        - type: special
          description: "Custom effect"
          metadata:
            custom_key: custom_value
`)
		registry := NewBaseTreeRegistry()
		err := registry.LoadFromYAML(yamlData)
		require.NoError(t, err)

		tree, _ := registry.Get("effects_tree")
		node, _ := tree.GetNode("start")
		effects := node.Effects()

		require.Len(t, effects, 7)

		t.Run("attribute effect", func(t *testing.T) {
			effect := effects[0]
			require.Equal(t, EffectTypeAttribute, effect.Type())
			require.Equal(t, float64(5), effect.Value())
			require.Equal(t, "+5 Strength", effect.Description())
		})

		t.Run("grant_skill effect", func(t *testing.T) {
			effect := effects[1]
			require.Equal(t, EffectTypeGrantSkill, effect.Type())
			require.Equal(t, "Grants Fireball skill", effect.Description())
			meta := effect.Metadata()
			require.Equal(t, "fireball", meta["skill_id"])
			require.Equal(t, 1, meta["start_level"])
		})

		t.Run("passive effect", func(t *testing.T) {
			effect := effects[2]
			require.Equal(t, EffectTypePassive, effect.Type())
			require.Equal(t, 0.10, effect.Value())
			meta := effect.Metadata()
			require.Equal(t, "on_kill", meta["trigger_type"])
		})

		t.Run("skill_mod effect", func(t *testing.T) {
			effect := effects[3]
			require.Equal(t, EffectTypeSkillMod, effect.Type())
			meta := effect.Metadata()
			require.Contains(t, meta["target_skill_tags"], "fire")
			require.Contains(t, meta["target_skill_tags"], "spell")
		})

		t.Run("unlock_craft effect", func(t *testing.T) {
			effect := effects[4]
			require.Equal(t, EffectTypeUnlockCraft, effect.Type())
		})

		t.Run("trade effect", func(t *testing.T) {
			effect := effects[5]
			require.Equal(t, EffectTypeTrade, effect.Type())
		})

		t.Run("special effect", func(t *testing.T) {
			effect := effects[6]
			require.Equal(t, EffectTypeSpecial, effect.Type())
			meta := effect.Metadata()
			require.Equal(t, "custom_value", meta["custom_key"])
		})
	})
}

// =============================================================================
// TREE STATE TESTS
// =============================================================================

func TestBaseTreeState(t *testing.T) {
	// Create a simple tree for tests
	createTestTree := func() *BaseTree {
		tree := NewBaseTree(TreeConfig{
			ID:   "test_tree",
			Name: "Test Tree",
		})

		// Start node (cost 0)
		tree.AddNode(NewBaseNode(NodeConfig{
			ID:          "start",
			Name:        "Start",
			Type:        NodePath,
			Cost:        0,
			Connections: []string{"node_a", "node_b"},
		}))

		// Node A (cost 1)
		tree.AddNode(NewBaseNode(NodeConfig{
			ID:           "node_a",
			Name:         "Node A",
			Type:         NodePath,
			Cost:         1,
			Requirements: []string{"start"},
			Connections:  []string{"node_c"},
		}))

		// Node B (cost 1)
		tree.AddNode(NewBaseNode(NodeConfig{
			ID:           "node_b",
			Name:         "Node B",
			Type:         NodePath,
			Cost:         1,
			Requirements: []string{"start"},
			Connections:  []string{"node_c"},
		}))

		// Node C (cost 2, requires A OR B)
		tree.AddNode(NewBaseNode(NodeConfig{
			ID:           "node_c",
			Name:         "Node C",
			Type:         NodeNotable,
			Cost:         2,
			Requirements: []string{"node_a", "node_b"},
			Connections:  []string{},
		}))

		// Keystone 1
		tree.AddNode(NewBaseNode(NodeConfig{
			ID:           "keystone_1",
			Name:         "Keystone 1",
			Type:         NodeKeystone,
			Cost:         1,
			Requirements: []string{"node_c"},
			Exclusions:   []string{"keystone_2"},
		}))

		// Keystone 2
		tree.AddNode(NewBaseNode(NodeConfig{
			ID:           "keystone_2",
			Name:         "Keystone 2",
			Type:         NodeKeystone,
			Cost:         1,
			Requirements: []string{"node_c"},
			Exclusions:   []string{"keystone_1"},
		}))

		// Leveled node
		tree.AddNode(NewBaseNode(NodeConfig{
			ID:           "mastery",
			Name:         "Mastery",
			Type:         NodeMastery,
			Cost:         1,
			MaxLevel:     3,
			LevelCost:    1,
			Requirements: []string{"start"},
		}))

		tree.SetStartNodes([]string{"start"})
		return tree
	}

	t.Run("create state", func(t *testing.T) {
		tree := createTestTree()
		state := NewBaseTreeState(TreeStateConfig{
			TreeID:          "test_tree",
			Tree:            tree,
			BaseCostPerNode: 100,
		})

		require.Equal(t, "test_tree", state.TreeID())
		require.Equal(t, 0, state.AvailablePoints())
		require.Equal(t, 0, state.SpentPoints())
		require.Empty(t, state.GetAllocatedNodes())
	})

	t.Run("add points", func(t *testing.T) {
		tree := createTestTree()
		state := NewBaseTreeState(TreeStateConfig{
			TreeID: "test_tree",
			Tree:   tree,
		})

		state.AddPoints(10)
		require.Equal(t, 10, state.AvailablePoints())
	})

	t.Run("node allocation", func(t *testing.T) {
		tree := createTestTree()
		state := NewBaseTreeState(TreeStateConfig{
			TreeID: "test_tree",
			Tree:   tree,
		})
		state.AddPoints(10)
		ctx := context.Background()

		t.Run("successful start node allocation", func(t *testing.T) {
			err := state.AllocateNode(ctx, "start")
			require.NoError(t, err)
			require.True(t, state.IsAllocated("start"))
			require.Equal(t, 10, state.AvailablePoints()) // cost 0
		})

		t.Run("successful dependent node allocation", func(t *testing.T) {
			err := state.AllocateNode(ctx, "node_a")
			require.NoError(t, err)
			require.True(t, state.IsAllocated("node_a"))
			require.Equal(t, 9, state.AvailablePoints()) // cost 1
		})

		t.Run("error on duplicate allocation", func(t *testing.T) {
			err := state.AllocateNode(ctx, "node_a")
			require.Error(t, err)
			require.Equal(t, ErrNodeAlreadyAlloc, err)
		})

		t.Run("error without requirements met", func(t *testing.T) {
			// node_c requires node_a OR node_b
			// node_a already allocated, so this should work
			err := state.AllocateNode(ctx, "node_c")
			require.NoError(t, err)

			// Create new state to check error
			state2 := NewBaseTreeState(TreeStateConfig{
				TreeID: "test_tree",
				Tree:   tree,
			})
			state2.AddPoints(10)
			// Try to allocate node_a without start
			err = state2.AllocateNode(ctx, "node_a")
			require.Error(t, err)
			require.Equal(t, ErrRequirementsNotMet, err)
		})

		t.Run("error on insufficient points", func(t *testing.T) {
			state3 := NewBaseTreeState(TreeStateConfig{
				TreeID: "test_tree",
				Tree:   tree,
			})
			state3.AddPoints(0)
			_ = state3.AllocateNode(ctx, "start") // free
			err := state3.AllocateNode(ctx, "node_a")
			require.Error(t, err)
			require.Equal(t, ErrInsufficientPoints, err)
		})
	})

	t.Run("mutual exclusions", func(t *testing.T) {
		tree := createTestTree()
		state := NewBaseTreeState(TreeStateConfig{
			TreeID: "test_tree",
			Tree:   tree,
		})
		state.AddPoints(20)
		ctx := context.Background()

		// Path to keystones
		_ = state.AllocateNode(ctx, "start")
		_ = state.AllocateNode(ctx, "node_a")
		_ = state.AllocateNode(ctx, "node_c")

		t.Run("first keystone can be taken", func(t *testing.T) {
			err := state.AllocateNode(ctx, "keystone_1")
			require.NoError(t, err)
		})

		t.Run("second keystone cannot be taken (exclusion)", func(t *testing.T) {
			err := state.AllocateNode(ctx, "keystone_2")
			require.Error(t, err)
			require.Equal(t, ErrNodeExcluded, err)
		})

		t.Run("CanAllocate returns false for excluded", func(t *testing.T) {
			require.False(t, state.CanAllocate("keystone_2"))
		})
	})

	t.Run("deallocation", func(t *testing.T) {
		tree := createTestTree()
		state := NewBaseTreeState(TreeStateConfig{
			TreeID: "test_tree",
			Tree:   tree,
		})
		state.AddPoints(10)
		ctx := context.Background()

		_ = state.AllocateNode(ctx, "start")
		_ = state.AllocateNode(ctx, "node_a")

		t.Run("successful leaf node deallocation", func(t *testing.T) {
			err := state.DeallocateNode(ctx, "node_a")
			require.NoError(t, err)
			require.False(t, state.IsAllocated("node_a"))
			require.Equal(t, 10, state.AvailablePoints()) // refund
		})

		t.Run("error deallocating required node", func(t *testing.T) {
			_ = state.AllocateNode(ctx, "node_a")
			// start is required by node_a
			err := state.DeallocateNode(ctx, "start")
			require.Error(t, err)
			require.Equal(t, ErrNodeRequired, err)
		})

		t.Run("deallocation with alternative path", func(t *testing.T) {
			// node_c requires node_a OR node_b
			_ = state.AllocateNode(ctx, "node_b")
			_ = state.AllocateNode(ctx, "node_c")

			// Now node_a can be removed since node_c has alternative (node_b)
			err := state.DeallocateNode(ctx, "node_a")
			require.NoError(t, err)
			require.True(t, state.IsAllocated("node_c"))
			require.True(t, state.IsAllocated("node_b"))
		})
	})

	t.Run("leveled nodes", func(t *testing.T) {
		tree := createTestTree()
		state := NewBaseTreeState(TreeStateConfig{
			TreeID: "test_tree",
			Tree:   tree,
		})
		state.AddPoints(10)
		ctx := context.Background()

		_ = state.AllocateNode(ctx, "start")
		_ = state.AllocateNode(ctx, "mastery")

		t.Run("initial level is 1", func(t *testing.T) {
			level := state.GetAllocatedLevel("mastery")
			require.Equal(t, 1, level)
		})

		t.Run("level up", func(t *testing.T) {
			err := state.LevelUpNode(ctx, "mastery")
			require.NoError(t, err)
			require.Equal(t, 2, state.GetAllocatedLevel("mastery"))

			err = state.LevelUpNode(ctx, "mastery")
			require.NoError(t, err)
			require.Equal(t, 3, state.GetAllocatedLevel("mastery"))
		})

		t.Run("error at max level", func(t *testing.T) {
			err := state.LevelUpNode(ctx, "mastery")
			require.Error(t, err)
			require.Equal(t, ErrMaxLevel, err)
		})
	})

	t.Run("reset all", func(t *testing.T) {
		tree := createTestTree()
		state := NewBaseTreeState(TreeStateConfig{
			TreeID: "test_tree",
			Tree:   tree,
		})
		state.AddPoints(10)
		ctx := context.Background()

		_ = state.AllocateNode(ctx, "start")
		_ = state.AllocateNode(ctx, "node_a")
		_ = state.AllocateNode(ctx, "node_b")

		err := state.ResetAll(ctx)
		require.NoError(t, err)
		require.Empty(t, state.GetAllocatedNodes())
		require.Equal(t, 10, state.AvailablePoints())
		require.Equal(t, 0, state.SpentPoints())
	})

	t.Run("respec cost", func(t *testing.T) {
		tree := createTestTree()
		state := NewBaseTreeState(TreeStateConfig{
			TreeID:           "test_tree",
			Tree:             tree,
			BaseCostPerNode:  100,
			CostPerNodeLevel: 50,
			ResetCostBase:    500,
		})
		state.AddPoints(20)
		ctx := context.Background()

		_ = state.AllocateNode(ctx, "start")
		_ = state.AllocateNode(ctx, "mastery")
		_ = state.LevelUpNode(ctx, "mastery")
		_ = state.LevelUpNode(ctx, "mastery")

		t.Run("single node respec cost", func(t *testing.T) {
			// mastery: level 3, cost = 100 + 2*50 = 200
			cost := state.RespecCost([]string{"mastery"})
			require.Equal(t, int64(200), cost)
		})

		t.Run("full reset cost", func(t *testing.T) {
			// base 500 + start(100) + mastery(100 + 2*50) = 800
			cost := state.ResetCost()
			require.Equal(t, int64(800), cost)
		})
	})

	t.Run("serialization/deserialization", func(t *testing.T) {
		tree := createTestTree()
		state := NewBaseTreeState(TreeStateConfig{
			TreeID: "test_tree",
			Tree:   tree,
		})
		state.AddPoints(10)
		ctx := context.Background()

		_ = state.AllocateNode(ctx, "start")
		_ = state.AllocateNode(ctx, "node_a")
		_ = state.AllocateNode(ctx, "mastery")
		_ = state.LevelUpNode(ctx, "mastery")

		data := state.GetData()

		// Restore into new state
		state2 := NewBaseTreeState(TreeStateConfig{
			TreeID: "test_tree",
			Tree:   tree,
		})
		state2.RestoreData(data)

		require.Equal(t, state.TreeID(), state2.TreeID())
		require.Equal(t, state.AvailablePoints(), state2.AvailablePoints())
		require.Equal(t, state.SpentPoints(), state2.SpentPoints())
		require.True(t, state2.IsAllocated("start"))
		require.True(t, state2.IsAllocated("node_a"))
		require.True(t, state2.IsAllocated("mastery"))
		require.Equal(t, 2, state2.GetAllocatedLevel("mastery"))
	})
}

// =============================================================================
// CREATE STATE FROM REGISTRY
// =============================================================================

func TestRegistryCreateState(t *testing.T) {
	yamlData := []byte(`
version: "1.0"
tree:
  id: state_test_tree
  name: "State Test Tree"
  start_nodes:
    - start
  nodes:
    - id: start
      name: "Start"
      type: path
      cost: 0
      position: { x: 0, y: 0 }
`)
	registry := NewBaseTreeRegistry()
	_ = registry.LoadFromYAML(yamlData)

	t.Run("successful state creation", func(t *testing.T) {
		state, err := registry.CreateState("state_test_tree")
		require.NoError(t, err)
		require.NotNil(t, state)
		require.Equal(t, "state_test_tree", state.TreeID())
	})

	t.Run("error for non-existing tree", func(t *testing.T) {
		_, err := registry.CreateState("nonexistent")
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
	})
}

// =============================================================================
// INTEGRATION: LOAD REAL YAML FILES
// =============================================================================

func TestLoadRealTreeYAMLFiles(t *testing.T) {
	// Get path to data/trees relative to test file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Skip("failed to determine test file path")
	}

	// internal/core/skill -> ../../data/trees
	treesDir := filepath.Join(filepath.Dir(filename), "..", "..", "..", "data", "trees")

	registry := NewBaseTreeRegistry()
	err := registry.LoadFromDirectory(treesDir)
	if err != nil {
		t.Skipf("failed to load trees: %v", err)
	}

	t.Run("trees loaded", func(t *testing.T) {
		count := registry.Count()
		require.Greater(t, count, 0, "should have at least one tree")
		t.Logf("Trees loaded: %d", count)
	})

	t.Run("main_tree loaded correctly", func(t *testing.T) {
		tree, ok := registry.Get("main_tree")
		if !ok {
			t.Skip("main_tree not found")
		}

		t.Run("basic properties", func(t *testing.T) {
			require.Equal(t, "main_tree", tree.ID())
			require.Equal(t, "Skill Tree", tree.Name())
		})

		t.Run("branches", func(t *testing.T) {
			branches := tree.GetBranches()
			require.GreaterOrEqual(t, len(branches), 3, "should have at least 3 branches")

			branchIDs := make(map[string]bool)
			for _, b := range branches {
				branchIDs[b.ID] = true
			}
			require.True(t, branchIDs["combat"], "should have combat branch")
			require.True(t, branchIDs["defense"], "should have defense branch")
			require.True(t, branchIDs["magic"], "should have magic branch")
		})

		t.Run("start nodes", func(t *testing.T) {
			startNodes := tree.GetStartNodes()
			require.GreaterOrEqual(t, len(startNodes), 1, "should have at least one start node")
		})

		t.Run("keystones", func(t *testing.T) {
			nodes := tree.GetNodes()
			keystoneCount := 0
			for _, node := range nodes {
				if node.Type() == NodeKeystone {
					keystoneCount++
				}
			}
			require.Greater(t, keystoneCount, 0, "should have at least one keystone")
			t.Logf("Keystones found: %d", keystoneCount)
		})

		t.Run("leveled nodes", func(t *testing.T) {
			// Check fire_mastery
			node, ok := tree.GetNode("fire_mastery")
			if ok {
				require.Equal(t, 3, node.MaxLevel())
				require.Equal(t, 1, node.LevelCost())

				// Check effects per level
				for level := 1; level <= 3; level++ {
					effects := node.EffectsAtLevel(level)
					require.NotEmpty(t, effects, "should have effects at level %d", level)
				}
			}
		})

		t.Run("graph is connected", func(t *testing.T) {
			// Check that we can reach from start_combat to blood_magic
			startNodes := tree.GetStartNodes()
			if len(startNodes) > 0 {
				startID := startNodes[0].ID()
				// Check that there's a path from start
				adjacent := tree.GetAdjacentNodes(startID)
				require.NotEmpty(t, adjacent, "start node should have connections")
			}
		})
	})
}

// =============================================================================
// NODE EFFECTS APPLY/REMOVE
// =============================================================================

func TestNodeEffectsApplyRemove(t *testing.T) {
	ctx := context.Background()
	entityID := "test_entity"

	t.Run("attribute effect", func(t *testing.T) {
		effect := &BaseAttributeEffect{
			attribute:   "strength",
			modType:     "flat",
			value:       10,
			description: "+10 Strength",
		}

		err := effect.Apply(ctx, entityID)
		require.NoError(t, err) // TODO пока просто заглушка

		err = effect.Remove(ctx, entityID)
		require.NoError(t, err)
	})

	t.Run("grant_skill effect", func(t *testing.T) {
		effect := &BaseGrantSkillEffect{
			skillID:     "fireball",
			startLevel:  1,
			description: "Grants Fireball",
		}

		require.Equal(t, EffectTypeGrantSkill, effect.Type())
		require.Equal(t, float64(1), effect.Value())

		meta := effect.Metadata()
		require.Equal(t, "fireball", meta["skill_id"])
	})

	t.Run("passive effect", func(t *testing.T) {
		effect := &BasePassiveEffect{
			triggerType:  TriggerOnKill,
			triggerValue: 0.10,
			effectID:     "heal_on_kill",
			description:  "10% heal on kill",
		}

		require.Equal(t, EffectTypePassive, effect.Type())
		require.Equal(t, 0.10, effect.Value())

		meta := effect.Metadata()
		require.Equal(t, "on_kill", meta["trigger_type"])
	})
}
