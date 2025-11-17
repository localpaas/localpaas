-- +migrate Up
CREATE TABLE IF NOT EXISTS projects
(
    id           VARCHAR(100) PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    key          VARCHAR(100) NOT NULL,
    photo        VARCHAR(255) NULL,
    status       VARCHAR(100) NOT NULL,
    note         VARCHAR(10000) NULL,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX idx_uq_projects_name ON projects(LOWER(name)) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_uq_projects_key ON projects(LOWER(key)) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS projects;
