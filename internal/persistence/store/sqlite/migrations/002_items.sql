-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id TEXT PRIMARY KEY,
    entity_type TEXT NOT NULL,
    name TEXT NOT NULL,
    item_type TEXT NOT NULL,
    rarity INTEGER NOT NULL DEFAULT 0,
    level INTEGER NOT NULL DEFAULT 1,
    quality REAL NOT NULL DEFAULT 1.0,
    stack_size INTEGER NOT NULL DEFAULT 1,
    max_stack_size INTEGER NOT NULL DEFAULT 1,
    value INTEGER NOT NULL DEFAULT 0,
    weight REAL NOT NULL DEFAULT 0.1,
    data BLOB NOT NULL,
    owner_id TEXT,
    container_id TEXT,
    equipped_slot TEXT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX idx_items_entity_type ON items(entity_type);
CREATE INDEX idx_items_item_type ON items(item_type);
CREATE INDEX idx_items_owner ON items(owner_id);
CREATE INDEX idx_items_container ON items(container_id);
CREATE INDEX idx_items_rarity ON items(rarity);
CREATE INDEX idx_items_level ON items(level);
CREATE INDEX idx_items_name ON items(name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item_tags (
    item_id TEXT NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (item_id, tag),
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

CREATE INDEX idx_item_tags_tag ON item_tags(tag);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS equipment_sockets (
    id TEXT PRIMARY KEY,
    equipment_id TEXT NOT NULL,
    slot_index INTEGER NOT NULL,
    socket_type TEXT NOT NULL,
    socketed_item_id TEXT,
    FOREIGN KEY (equipment_id) REFERENCES items(id) ON DELETE CASCADE,
    FOREIGN KEY (socketed_item_id) REFERENCES items(id) ON DELETE SET NULL,
    UNIQUE (equipment_id, slot_index)
);

CREATE INDEX idx_equipment_sockets_equipment ON equipment_sockets(equipment_id);
CREATE INDEX idx_equipment_sockets_socketed ON equipment_sockets(socketed_item_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS container_contents (
    container_id TEXT NOT NULL,
    item_id TEXT NOT NULL,
    slot_index INTEGER,
    PRIMARY KEY (container_id, item_id),
    FOREIGN KEY (container_id) REFERENCES items(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

CREATE INDEX idx_container_contents_container ON container_contents(container_id);
CREATE INDEX idx_container_contents_item ON container_contents(item_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS container_contents;
DROP TABLE IF EXISTS equipment_sockets;
DROP TABLE IF EXISTS item_tags;
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
