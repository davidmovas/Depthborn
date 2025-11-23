# Depthborn
*A modular ARPG Extraction Roguelite Dungeon-Crawler*

Depthborn is a long-term experimental project focused on building a deep, scalable, fully modular dungeon-crawler ARPG.  
It blends three genres:

- **ARPG (Diablo-like):** items, bosses, affixes, crafting, rarity tiers
- **Roguelite:** procedural layers, unpredictable encounters, escalating difficulty
- **Extraction:** risk/reward gameplay, temporary runs, safe hubs, loot preservation

The project is primarily a sandbox for experimenting with game systems, architectural patterns, data-driven design, procedural generation, and combat logic—written in Go and designed with a clean separation between **core game logic** and **presentation/UI** (e.g. BubbleTea or any other interface).

---

## 🚀 Core Concept

You explore a **persistent, endlessly generated dungeon** composed of “Layers.”  
Each layer:

- is generated **once** upon first visit,
- permanently stores its biome, enemy families, resource types and boss identity,
- but regularly **refreshes** temporary entities (monsters, traps, resources) on cooldowns,
- becomes progressively more dangerous the deeper you go.

Players can control **one or multiple characters**, forming a party.  
Each character is a blank slate, customizable through items, abilities, stat attributes and a large skill tree.

Every expedition into the dungeon is high-risk:  
If you die — you lose most of the loot.  
If you successfully extract — you return to the camp to craft, upgrade and plan your next descent.

---

## 🧩 Key Systems

### **Dungeon & World System**
- Persistent layers with unique biomes and enemy families
- One-time generation per layer, periodic restocks
- Infinite depth support w/ difficulty scaling
- Boss respawn timers
- Environmental modifiers and hazards

### **Party System**
- Multiple characters per account
- Flexible party composition (solo or multi-character)
- Scaled encounters based on party size

### **Character Progression**
- Classless character design
- Attributes, derived stats, resistances, speed, etc.
- Inventory, equipment slots, ability loadouts
- Large, scalable skill tree (PoE-style)
- Active and passive abilities with cooldowns, triggers and synergies

### **Combat System**
- Tactical, turn-based or hybrid auto-resolution
- Status effects, DoTs, elemental interactions
- AI behavior patterns
- Multi-entity combat support

### **Items & Loot**
- Rarity tiers, affixes, prefixes/suffixes
- Durability and item degradation
- Unique bosses dropping unique loot
- Materials, crafting components, consumables

### **Crafting**
- Item upgrades
- Affix rerolls
- Rune/gem systems
- Material transmutation

### **Camp (Hub)**
- Safe zone for crafting, storage, repairs
- Traders, contracts, upgrades
- Party management
- Researching enemies (bestiary), lore and materials

### **Persistence**
- Characters
- Party composition
- Inventory & stash
- Revealed dungeon layers
- Layer cooldowns
- Bestiary discoveries

---

## 🏗 Architecture & Goals

The aim is to keep the codebase **modular, data-driven, and extendable**.

Design goals:

- Clean separation of **Game Logic** ↔ **UI Layer**
- Entity-Component-like abstractions without strict ECS complexity
- Data-first approach (JSON, YAML, etc.)
- Expandability for:
    - new monsters
    - new biomes
    - new items
    - new skill trees
    - new camp modules
- High-level simulation friendly for CLI/TUI testing

The project is intentionally open-ended — it’s a personal playground for exploring system architecture, procedural content systems and Go’s composition model.

---

## 📦 Technologies

- **Go** — core logic, modular architecture
- **BubbleTea (TUI)** — optional interface for visualizing game sessions
- **JSON/YAML** — data-driven layer definitions, items, abilities
- **Go generics + composition** — flexible entity building

(Other UI layers — CLI, HTML, or graphical — can be plugged in later without touching game logic.)

---

## 🎯 Project Status

This project is currently in **early conceptual & systems-design phase**.  
It evolves gradually — the goal is experimentation, not a polished commercial product.

Contributions are welcome if they follow the architectural spirit.

---

## 📚 Roadmap (High-Level)

- [ ] Base structuring of core modules
- [ ] Persistent layer generation prototype
- [ ] Layer restock & cooldown system
- [ ] Character attributes + equipment model
- [ ] Party system
- [ ] Combat engine (v1)
- [ ] Items, rarity, affix pools
- [ ] Basic crafting
- [ ] Camp prototype
- [ ] TUI interface implementation
- [ ] Advanced content (bosses, elite mods, events)

---

## 📝 License

TBD — will be added later based on project direction.

---

## 🌌 Inspiration

Diablo • Path of Exile • Hades • Darkest Dungeon • Escape From Tarkov • FTL • Slay the Spire • Dungeon Crawl Classics • Rogue Legacy

---

## ✨ Author

A personal project created for fun, exploration, architecture practice, and experimenting with Go game systems.

