-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS snapshots (
    entity_id TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    version INTEGER NOT NULL,
    timestamp INTEGER NOT NULL,
    data BLOB NOT NULL,
    size INTEGER NOT NULL,
    PRIMARY KEY (entity_id, entity_type, version)
);

CREATE INDEX idx_snapshots_id_type ON snapshots(entity_id, entity_type);
CREATE INDEX idx_snapshots_timestamp ON snapshots(timestamp);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS deltas (
    entity_id TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    from_version INTEGER NOT NULL,
    to_version INTEGER NOT NULL,
    timestamp INTEGER NOT NULL,
    data BLOB NOT NULL,
    PRIMARY KEY (entity_id, entity_type, from_version, to_version)
);

CREATE INDEX idx_deltas_id_type ON deltas(entity_id, entity_type);
CREATE INDEX idx_deltas_versions ON deltas(entity_id, entity_type, from_version);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS entity_metadata (
    entity_id TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    current_version INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    PRIMARY KEY (entity_type, entity_id)
);

CREATE INDEX idx_metadata_type ON entity_metadata(entity_type);
CREATE INDEX idx_metadata_updated ON entity_metadata(updated_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS entity_metadata;
DROP TABLE IF EXISTS deltas;
DROP TABLE IF EXISTS snapshots;
-- +goose StatementEnd