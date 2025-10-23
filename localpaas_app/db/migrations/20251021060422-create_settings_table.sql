-- +migrate Up
CREATE TABLE IF NOT EXISTS settings
(
    id           VARCHAR(26) PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    type         VARCHAR(100) NOT NULL,
    object_id    VARCHAR(26) NULL,
    data         JSONB NULL,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMPTZ NULL
);

CREATE INDEX idx_settings_type ON settings(type);
CREATE INDEX idx_settings_name ON settings(name);
CREATE INDEX idx_settings_object_id ON settings(object_id);
CREATE INDEX idx_settings_created_at ON settings(created_at);
CREATE INDEX idx_settings_deleted_at ON settings(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS settings;
