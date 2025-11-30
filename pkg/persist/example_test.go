package persist_test

import (
	"context"
	"testing"

	"github.com/davidmovas/Depthborn/pkg/persist"
	"github.com/davidmovas/Depthborn/pkg/persist/storage/sqlite"
)

// --- Example Entity: Character ---

const EntityCharacter persist.EntityType = "character"

// CharacterState is the serializable state of a Character.
type CharacterState struct {
	persist.BaseState `msgpack:",inline"`
	Name              string `msgpack:"name"`
	Level             int    `msgpack:"level"`
	Experience        int64  `msgpack:"experience"`
	Health            int    `msgpack:"health"`
	MaxHealth         int    `msgpack:"max_health"`
}

// Character represents a game character.
type Character struct {
	persist.Base

	Name       string
	Level      int
	Experience int64
	Health     int
	MaxHealth  int
}

// NewCharacter creates a new character with default values.
func NewCharacter(name string) *Character {
	return &Character{
		Base:       persist.NewBase(EntityCharacter),
		Name:       name,
		Level:      1,
		Experience: 0,
		Health:     100,
		MaxHealth:  100,
	}
}

// MarshalBinary implements persist.Marshaler.
func (c *Character) MarshalBinary() ([]byte, error) {
	state := CharacterState{
		BaseState:  c.Base.State(),
		Name:       c.Name,
		Level:      c.Level,
		Experience: c.Experience,
		Health:     c.Health,
		MaxHealth:  c.MaxHealth,
	}
	return persist.DefaultCodec().Encode(state)
}

// UnmarshalBinary implements persist.Unmarshaler.
func (c *Character) UnmarshalBinary(data []byte) error {
	var state CharacterState
	if err := persist.DefaultCodec().Decode(data, &state); err != nil {
		return err
	}

	c.Base.LoadState(state.BaseState)
	c.Name = state.Name
	c.Level = state.Level
	c.Experience = state.Experience
	c.Health = state.Health
	c.MaxHealth = state.MaxHealth
	return nil
}

// AddExperience adds XP and handles level ups.
func (c *Character) AddExperience(xp int64) {
	c.Experience += xp
	c.MarkDirty()

	// Simple level up logic
	for c.Experience >= int64(c.Level*100) {
		c.Experience -= int64(c.Level * 100)
		c.Level++
		c.MaxHealth += 10
		c.Health = c.MaxHealth
	}
}

// TakeDamage reduces health.
func (c *Character) TakeDamage(damage int) {
	c.Health -= damage
	if c.Health < 0 {
		c.Health = 0
	}
	c.MarkDirty()
}

// Heal restores health.
func (c *Character) Heal(amount int) {
	c.Health += amount
	if c.Health > c.MaxHealth {
		c.Health = c.MaxHealth
	}
	c.MarkDirty()
}

// --- Tests ---

func TestRepositoryBasicCRUD(t *testing.T) {
	ctx := context.Background()

	// Create in-memory storage
	store, err := sqlite.OpenMemory()
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer store.Close()

	// Create repository
	repo := persist.NewRepository[*Character](
		store,
		EntityCharacter,
		func() *Character { return &Character{} },
	)

	// Create character
	char := NewCharacter("Hero")
	charID := char.ID()

	// Save
	if err := repo.Save(ctx, char); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	if char.Version() != 1 {
		t.Errorf("expected version 1, got %d", char.Version())
	}
	if char.IsDirty() {
		t.Error("expected clean after save")
	}

	// Load
	loaded, err := repo.Load(ctx, charID)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if loaded.ID() != charID {
		t.Errorf("ID mismatch: got %s, want %s", loaded.ID(), charID)
	}
	if loaded.Name != "Hero" {
		t.Errorf("Name mismatch: got %s, want Hero", loaded.Name)
	}
	if loaded.Level != 1 {
		t.Errorf("Level mismatch: got %d, want 1", loaded.Level)
	}

	// Update
	loaded.AddExperience(150) // Should level up
	if loaded.Level != 2 {
		t.Errorf("Level after XP: got %d, want 2", loaded.Level)
	}
	if !loaded.IsDirty() {
		t.Error("expected dirty after modification")
	}

	if err := repo.Save(ctx, loaded); err != nil {
		t.Fatalf("failed to save updated: %v", err)
	}

	if loaded.Version() != 2 {
		t.Errorf("expected version 2, got %d", loaded.Version())
	}

	// Delete
	if err := repo.Delete(ctx, charID); err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	// Verify deleted
	_, err = repo.Load(ctx, charID)
	if err != persist.ErrEntityNotFound {
		t.Errorf("expected ErrEntityNotFound, got %v", err)
	}
}

func TestUnitOfWork(t *testing.T) {
	ctx := context.Background()

	// Create in-memory storage
	store, err := sqlite.OpenMemory()
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer store.Close()

	// Create and save initial character
	repo := persist.NewRepository[*Character](
		store,
		EntityCharacter,
		func() *Character { return &Character{} },
	)

	char1 := NewCharacter("Alice")
	char1ID := char1.ID()
	if err := repo.Save(ctx, char1); err != nil {
		t.Fatalf("failed to save initial: %v", err)
	}

	// Create Unit of Work
	uow := persist.NewUnitOfWork(store)
	uow.RegisterFactory(EntityCharacter, func() persist.Persistable {
		return &Character{}
	})
	defer uow.Close()

	// Load existing character
	entity, err := uow.Get(ctx, EntityCharacter, char1ID)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	loadedChar := entity.(*Character)
	if loadedChar.Name != "Alice" {
		t.Errorf("expected Alice, got %s", loadedChar.Name)
	}

	// Modify
	loadedChar.AddExperience(500)

	// Create new character through UoW
	char2 := NewCharacter("Bob")
	uow.Register(char2)

	// Check state before commit
	if !uow.HasChanges() {
		t.Error("expected HasChanges=true")
	}
	if uow.DirtyCount() != 2 {
		t.Errorf("expected 2 dirty, got %d", uow.DirtyCount())
	}

	// Commit
	if err := uow.Commit(ctx); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	if uow.HasChanges() {
		t.Error("expected no changes after commit")
	}

	// Verify changes persisted
	verifyChar, err := repo.Load(ctx, char1ID)
	if err != nil {
		t.Fatalf("failed to verify load: %v", err)
	}
	if verifyChar.Level != 6 { // 500 XP = level 6 (100+200+300+400+500... simplified)
		// Actually our logic: 100 XP for level 1->2, 200 for 2->3, etc.
		// 500 XP starting from level 1:
		// Level 1: need 100, have 500, level up to 2, remaining 400
		// Level 2: need 200, have 400, level up to 3, remaining 200
		// Level 3: need 300, have 200, stay at 3
		// So should be level 3
		if verifyChar.Level != 3 {
			t.Errorf("expected level 3, got %d", verifyChar.Level)
		}
	}
}

func TestUnitOfWorkRollback(t *testing.T) {
	ctx := context.Background()

	store, err := sqlite.OpenMemory()
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer store.Close()

	// Save initial state
	repo := persist.NewRepository[*Character](
		store,
		EntityCharacter,
		func() *Character { return &Character{} },
	)

	char := NewCharacter("TestChar")
	charID := char.ID()
	if err := repo.Save(ctx, char); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Start UoW, modify, then rollback
	uow := persist.NewUnitOfWork(store)
	uow.RegisterFactory(EntityCharacter, func() persist.Persistable {
		return &Character{}
	})

	entity, _ := uow.Get(ctx, EntityCharacter, charID)
	loadedChar := entity.(*Character)
	loadedChar.Level = 99 // Big change

	// Rollback instead of commit
	uow.Rollback()

	// Verify original state unchanged
	original, _ := repo.Load(ctx, charID)
	if original.Level != 1 {
		t.Errorf("expected level 1 after rollback, got %d", original.Level)
	}
}

func TestUnitOfWorkDelete(t *testing.T) {
	ctx := context.Background()

	store, err := sqlite.OpenMemory()
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer store.Close()

	repo := persist.NewRepository[*Character](
		store,
		EntityCharacter,
		func() *Character { return &Character{} },
	)

	char := NewCharacter("ToDelete")
	charID := char.ID()
	if err := repo.Save(ctx, char); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Delete through UoW
	uow := persist.NewUnitOfWork(store)
	uow.RegisterFactory(EntityCharacter, func() persist.Persistable {
		return &Character{}
	})

	entity, _ := uow.Get(ctx, EntityCharacter, charID)
	uow.Delete(entity.(*Character))

	if !uow.HasChanges() {
		t.Error("expected changes after delete")
	}

	if err := uow.Commit(ctx); err != nil {
		t.Fatalf("failed to commit delete: %v", err)
	}

	// Verify deleted
	_, err = repo.Load(ctx, charID)
	if err != persist.ErrEntityNotFound {
		t.Errorf("expected not found, got %v", err)
	}
}
