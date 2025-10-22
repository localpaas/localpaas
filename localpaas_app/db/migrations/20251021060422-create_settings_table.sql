-- +migrate Up
CREATE TABLE IF NOT EXISTS settings
(
    id           VARCHAR(26) PRIMARY KEY,
    target_type  VARCHAR(100) NOT NULL,
    target_id    VARCHAR(26) NULL,
    data         JSONB NULL,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by   VARCHAR(26) NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by   VARCHAR(26) NOT NULL,
    deleted_at   TIMESTAMPTZ NULL,

    CONSTRAINT fk_settings_created_by FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_settings_updated_by FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE INDEX idx_settings_target_type ON settings(target_type);
CREATE INDEX idx_settings_target_id ON settings(target_id);
CREATE INDEX idx_settings_created_at ON settings(created_at);
CREATE INDEX idx_settings_deleted_at ON settings(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS settings;
