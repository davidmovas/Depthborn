package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/davidmovas/Depthborn/internal/item"
	"github.com/davidmovas/Depthborn/internal/item/builder"
)

func main() {
	fmt.Println("=== Item System Demo ===\n")

	// 1. Create items using builder
	fmt.Println("1. Creating items with builder...")

	sword := builder.MeleeWeapon("Flaming Sword").
		Description("A sword imbued with fire magic").
		Slot(item.SlotMainHand).
		Rarity(item.RarityEpic).
		Level(15).
		Value(5000).
		Weight(3.5).
		Durability(100).
		Tags("fire", "melee", "two-handed").
		Build()

	potion := builder.Potion("Greater Healing Potion").
		Description("Restores a large amount of health").
		Rarity(item.RarityRare).
		Level(10).
		Value(150).
		MaxStack(20).
		Cooldown(5000).
		Charges(3).
		Tags("healing").
		Build()

	gem := builder.Gem("Perfect Ruby").
		Description("A flawless ruby that adds fire damage").
		Rarity(item.RarityRare).
		Tier(4).
		Tags("fire", "damage").
		Build()

	bag := builder.Bag("Adventurer's Backpack", 20).
		Description("A sturdy backpack for carrying items").
		MaxWeight(50.0).
		Rarity(item.RarityUncommon).
		Tags("storage").
		Build()

	gold := builder.Currency("Gold Coins").
		Description("Standard currency").
		Value(1).
		Build()
	gold.AddStack(999)

	printItem("Sword", sword)
	printItem("Potion", potion)
	printItem("Gem", gem)
	printItem("Bag", bag)
	printItem("Gold", gold)

	// 2. Serialize items
	fmt.Println("\n2. Serializing items...")

	swordData, err := sword.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal sword: %w", err))
	}
	fmt.Printf("   Sword serialized: %d bytes\n", len(swordData))

	potionData, err := potion.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal potion: %w", err))
	}
	fmt.Printf("   Potion serialized: %d bytes\n", len(potionData))

	gemData, err := gem.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal gem: %w", err))
	}
	fmt.Printf("   Gem serialized: %d bytes\n", len(gemData))

	bagData, err := bag.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal bag: %w", err))
	}
	fmt.Printf("   Bag serialized: %d bytes\n", len(bagData))

	goldData, err := gold.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal gold: %w", err))
	}
	fmt.Printf("   Gold serialized: %d bytes\n", len(goldData))

	// 3. Save to SQLite database
	fmt.Println("\n3. Saving to SQLite database...")

	dbPath := filepath.Join(os.TempDir(), "item_demo.db")
	db, err := initDB(dbPath)
	if err != nil {
		panic(fmt.Errorf("failed to init DB: %w", err))
	}
	defer db.Close()

	ctx := context.Background()

	items := map[string][]byte{
		sword.ID():  swordData,
		potion.ID(): potionData,
		gem.ID():    gemData,
		bag.ID():    bagData,
		gold.ID():   goldData,
	}

	types := map[string]string{
		sword.ID():  sword.Type(),
		potion.ID(): potion.Type(),
		gem.ID():    gem.Type(),
		bag.ID():    bag.Type(),
		gold.ID():   gold.Type(),
	}

	for id, data := range items {
		err = saveItem(ctx, db, id, types[id], data)
		if err != nil {
			panic(fmt.Errorf("failed to save item %s: %w", id, err))
		}
		fmt.Printf("   Saved item: %s (%s)\n", id, types[id])
	}

	// 4. Load from database
	fmt.Println("\n4. Loading from database...")

	// Load sword
	loadedSwordData, loadedType, err := loadItem(ctx, db, sword.ID())
	if err != nil {
		panic(fmt.Errorf("failed to load sword: %w", err))
	}
	fmt.Printf("   Loaded sword data: %d bytes, type: %s\n", len(loadedSwordData), loadedType)

	// 5. Deserialize and verify
	fmt.Println("\n5. Deserializing and verifying...")

	restoredSword := &item.BaseEquipment{}
	if err = restoredSword.Unmarshal(loadedSwordData); err != nil {
		panic(fmt.Errorf("failed to unmarshal sword: %w", err))
	}

	fmt.Println("\n   Original vs Restored Sword:")
	fmt.Printf("   Name: %s == %s ? %v\n", sword.Name(), restoredSword.Name(), sword.Name() == restoredSword.Name())
	fmt.Printf("   ID: %s == %s ? %v\n", sword.ID(), restoredSword.ID(), sword.ID() == restoredSword.ID())
	fmt.Printf("   Rarity: %s == %s ? %v\n", sword.Rarity(), restoredSword.Rarity(), sword.Rarity() == restoredSword.Rarity())
	fmt.Printf("   Level: %d == %d ? %v\n", sword.Level(), restoredSword.Level(), sword.Level() == restoredSword.Level())
	fmt.Printf("   Value: %d == %d ? %v\n", sword.Value(), restoredSword.Value(), sword.Value() == restoredSword.Value())
	fmt.Printf("   Durability: %.0f/%.0f == %.0f/%.0f ? %v\n",
		sword.Durability(), sword.MaxDurability(),
		restoredSword.Durability(), restoredSword.MaxDurability(),
		sword.Durability() == restoredSword.Durability() && sword.MaxDurability() == restoredSword.MaxDurability())

	// Load and verify potion
	loadedPotionData, _, err := loadItem(ctx, db, potion.ID())
	if err != nil {
		panic(fmt.Errorf("failed to load potion: %w", err))
	}

	restoredPotion := &item.BaseConsumable{}
	if err = restoredPotion.Unmarshal(loadedPotionData); err != nil {
		panic(fmt.Errorf("failed to unmarshal potion: %w", err))
	}

	fmt.Println("\n   Original vs Restored Potion:")
	fmt.Printf("   Name: %s == %s ? %v\n", potion.Name(), restoredPotion.Name(), potion.Name() == restoredPotion.Name())
	fmt.Printf("   Charges: %d == %d ? %v\n", potion.Charges(), restoredPotion.Charges(), potion.Charges() == restoredPotion.Charges())
	fmt.Printf("   MaxCooldown: %d == %d ? %v\n", potion.MaxCooldown(), restoredPotion.MaxCooldown(), potion.MaxCooldown() == restoredPotion.MaxCooldown())

	// Load and verify gem
	loadedGemData, _, err := loadItem(ctx, db, gem.ID())
	if err != nil {
		panic(fmt.Errorf("failed to load gem: %w", err))
	}

	restoredGem := &item.BaseSocketable{}
	if err = restoredGem.Unmarshal(loadedGemData); err != nil {
		panic(fmt.Errorf("failed to unmarshal gem: %w", err))
	}

	fmt.Println("\n   Original vs Restored Gem:")
	fmt.Printf("   Name: %s == %s ? %v\n", gem.Name(), restoredGem.Name(), gem.Name() == restoredGem.Name())
	fmt.Printf("   Tier: %d == %d ? %v\n", gem.Tier(), restoredGem.Tier(), gem.Tier() == restoredGem.Tier())
	fmt.Printf("   SocketType: %s == %s ? %v\n", gem.SocketType(), restoredGem.SocketType(), gem.SocketType() == restoredGem.SocketType())

	// 6. List all items in DB
	fmt.Println("\n6. Listing all items in database...")

	allItems, err := listItems(ctx, db)
	if err != nil {
		panic(fmt.Errorf("failed to list items: %w", err))
	}

	for _, it := range allItems {
		fmt.Printf("   - %s (%s): %d bytes\n", it.ID, it.Type, it.Size)
	}

	// 7. Cleanup
	fmt.Println("\n7. Cleaning up...")
	os.Remove(dbPath)
	fmt.Printf("   Removed database: %s\n", dbPath)

	fmt.Println("\n=== Demo Complete ===")
}

func printItem(label string, it item.Item) {
	fmt.Printf("   %s: %s (ID: %s, Type: %s, Rarity: %s, Level: %d)\n",
		label, it.Name(), it.ID()[:8]+"...", it.ItemType(), it.Rarity(), it.Level())
}

func initDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Create items table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id TEXT PRIMARY KEY,
			entity_type TEXT NOT NULL,
			data BLOB NOT NULL,
			created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
			updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
		);
		CREATE INDEX IF NOT EXISTS idx_items_type ON items(entity_type);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func saveItem(ctx context.Context, db *sql.DB, id, entityType string, data []byte) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO items (id, entity_type, data) VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET data = excluded.data, updated_at = strftime('%s', 'now')
	`, id, entityType, data)
	return err
}

func loadItem(ctx context.Context, db *sql.DB, id string) ([]byte, string, error) {
	var data []byte
	var entityType string
	err := db.QueryRowContext(ctx, `SELECT data, entity_type FROM items WHERE id = ?`, id).Scan(&data, &entityType)
	return data, entityType, err
}

type itemInfo struct {
	ID   string
	Type string
	Size int
}

func listItems(ctx context.Context, db *sql.DB) ([]itemInfo, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, entity_type, length(data) FROM items`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []itemInfo
	for rows.Next() {
		var info itemInfo
		if err = rows.Scan(&info.ID, &info.Type, &info.Size); err != nil {
			return nil, err
		}
		items = append(items, info)
	}

	return items, rows.Err()
}
