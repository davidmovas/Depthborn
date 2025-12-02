# State Management System Design

## Overview

Система управления состоянием для roguelike игры Depthborn.
Цель: простой, надёжный, production-ready подход к сериализации и персистенции.

## Core Principles

1. **Stateless Systems** - игровые системы не хранят состояние, только обрабатывают данные
2. **Explicit State** - состояние явно определено через структуры, не map[string]any
3. **Unit of Work** - атомарные транзакции для группы изменений
4. **Change Tracking** - автоматическое отслеживание dirty-состояния
5. **Snapshot + Delta** - эффективное хранение с инкрементальными изменениями

## Architecture Layers

```
┌──────────────────────────────────────────────────────────────┐
│                      Game Layer                               │
│  Systems (stateless): Combat, Inventory, Progression, AI     │
└──────────────────────────────┬───────────────────────────────┘
                               │
┌──────────────────────────────▼───────────────────────────────┐
│                      Session Layer                            │
│  GameSession - owns UnitOfWork, manages game lifecycle       │
└──────────────────────────────┬───────────────────────────────┘
                               │
┌──────────────────────────────▼───────────────────────────────┐
│                    Unit of Work                               │
│  - Tracks all loaded/modified entities                       │
│  - Provides Get/Create/Delete operations                     │
│  - Commit() saves all changes atomically                     │
│  - Rollback() discards changes                               │
└──────────────────────────────┬───────────────────────────────┘
                               │
┌──────────────────────────────▼───────────────────────────────┐
│                    Repository Layer                           │
│  Generic Repository<T> with type-safe operations             │
│  - Load(id) / Save(entity) / Delete(id)                      │
│  - Query(criteria) for complex lookups                       │
│  - Batch operations for performance                          │
└──────────────────────────────┬───────────────────────────────┘
                               │
┌──────────────────────────────▼───────────────────────────────┐
│                    Storage Layer                              │
│  ┌─────────────────┐     ┌─────────────────┐                 │
│  │  SQLite Store   │     │   YAML Store    │                 │
│  │  (dynamic data) │     │ (static config) │                 │
│  │  - Snapshots    │     │  - Recipes      │                 │
│  │  - Deltas       │     │  - Templates    │                 │
│  │  - Metadata     │     │  - Game config  │                 │
│  └─────────────────┘     └─────────────────┘                 │
└──────────────────────────────────────────────────────────────┘
```

## Key Components

### 1. Entity (база для всех игровых объектов)

```go
// Entity - базовый интерфейс для всех игровых сущностей
type Entity interface {
    ID() string
    Type() EntityType
}

// Stateful - сущность с отслеживанием состояния
type Stateful interface {
    Entity
    Version() int64
    IsDirty() bool
    MarkClean()
}

// Persistable - сущность которую можно сохранить
type Persistable interface {
    Stateful
    Marshal() ([]byte, error)
    Unmarshal([]byte) error
}
```

### 2. Unit of Work (транзакция)

```go
type UnitOfWork interface {
    // Получение сущностей (lazy load + cache)
    Get(entityType EntityType, id string) (Entity, error)

    // Регистрация новых/изменённых
    Register(entity Entity)

    // Удаление
    Delete(entity Entity)

    // Фиксация всех изменений
    Commit(ctx context.Context) error

    // Откат изменений
    Rollback()

    // Проверка наличия изменений
    HasChanges() bool
}
```

### 3. Repository (доступ к данным)

```go
type Repository[T Entity] interface {
    Load(ctx context.Context, id string) (T, error)
    Save(ctx context.Context, entity T) error
    Delete(ctx context.Context, id string) error

    // Batch operations
    LoadMany(ctx context.Context, ids []string) ([]T, error)
    SaveMany(ctx context.Context, entities []T) error

    // Query
    Find(ctx context.Context, query Query) ([]T, error)
    Count(ctx context.Context, query Query) (int, error)
}
```

### 4. Storage Backend

```go
type Storage interface {
    // Raw operations
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte) error
    Delete(ctx context.Context, key string) error

    // Batch
    GetMany(ctx context.Context, keys []string) (map[string][]byte, error)
    SetMany(ctx context.Context, items map[string][]byte) error

    // Transaction
    Begin() (Transaction, error)
}
```

## Entity Types

### Dynamic Entities (SQLite)
- Character (player, NPCs)
- Item instances
- World state
- Quest progress
- Game saves

### Static Entities (YAML)
- Item templates/definitions
- Monster templates
- Skill definitions
- Recipe definitions
- Game configuration

## Usage Example

```go
// Game session creates unit of work
session := game.NewSession(storage)
uow := session.UnitOfWork()

// System works with entities through UoW
func (s *CombatSystem) ProcessAttack(uow UnitOfWork, attackerID, targetID string) error {
    attacker, err := uow.Get(EntityCharacter, attackerID)
    if err != nil {
        return err
    }

    target, err := uow.Get(EntityCharacter, targetID)
    if err != nil {
        return err
    }

    // Modify entities - automatically tracked as dirty
    damage := s.CalculateDamage(attacker, target)
    target.TakeDamage(damage)

    // No need to explicitly register - dirty tracking automatic
    return nil
}

// At end of game loop or checkpoint
if err := uow.Commit(ctx); err != nil {
    // Handle error, possibly rollback
    uow.Rollback()
}
```

## Migration Strategy

1. Keep existing `pkg/state` for backwards compatibility
2. Create new `pkg/persist` package with clean implementation
3. Gradually migrate entity types to new system
4. Remove old code once migration complete

## File Structure

```
pkg/
  persist/
    entity.go       # Entity interfaces
    uow.go          # Unit of Work
    repository.go   # Generic repository
    storage/
      storage.go    # Storage interface
      sqlite/       # SQLite implementation
      yaml/         # YAML implementation
    codec/
      codec.go      # Serialization interface
      msgpack.go    # MessagePack codec
      json.go       # JSON codec (for debugging)

internal/
  game/
    session.go      # Game session management
    entities/       # Concrete entity types
      character.go
      item.go
      world.go
```

## Next Steps

1. [ ] Create `pkg/persist/entity.go` - base interfaces
2. [ ] Create `pkg/persist/storage/storage.go` - storage interface
3. [ ] Create `pkg/persist/storage/sqlite/sqlite.go` - SQLite implementation
4. [ ] Create `pkg/persist/repository.go` - generic repository
5. [ ] Create `pkg/persist/uow.go` - Unit of Work
6. [ ] Create example entity and test persistence
7. [ ] Migrate existing entities
