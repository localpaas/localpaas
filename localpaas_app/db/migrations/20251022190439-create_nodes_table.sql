-- +migrate Up
CREATE TABLE IF NOT EXISTS nodes
(
    id           VARCHAR(100) PRIMARY KEY,
    is_leader    BOOL NOT NULL DEFAULT FALSE,
    is_manager   BOOL NOT NULL DEFAULT FALSE,
    host_name    VARCHAR(100) NOT NULL,
    ip           VARCHAR(100) NOT NULL,
    status       VARCHAR(100) NOT NULL,
    infra_status VARCHAR(100) NOT NULL,
    info         JSONB NULL,
    note         VARCHAR(10000) NULL,

    last_synced_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by       VARCHAR(26) NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by       VARCHAR(26) NOT NULL,
    deleted_at       TIMESTAMPTZ NULL,

    CONSTRAINT fk_nodes_created_by FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_nodes_updated_by FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE INDEX idx_nodes_created_at ON nodes(created_at);
CREATE INDEX idx_nodes_deleted_at ON nodes(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS nodes;
