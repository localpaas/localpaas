-- +migrate Up
CREATE TABLE IF NOT EXISTS projects
(
    id           VARCHAR(26) PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    photo        VARCHAR(255) NULL,
    status       VARCHAR(100) NOT NULL,
    data         JSONB NOT NULL DEFAULT '{}',

    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by   VARCHAR(26) NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by   VARCHAR(26) NOT NULL,
    deleted_at   TIMESTAMPTZ NULL,

    CONSTRAINT fk_projects_created_by FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_projects_updated_by FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE UNIQUE INDEX idx_uq_projects_name ON projects(LOWER(name)) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS projects;
