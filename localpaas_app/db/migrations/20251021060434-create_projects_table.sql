-- +migrate Up
CREATE TABLE IF NOT EXISTS projects
(
    id           VARCHAR(26) PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    photo        VARCHAR(255) NULL,
    status       VARCHAR(100) NOT NULL,
    note         VARCHAR(10000) NULL,
    settings_id  VARCHAR(26) NULL,
    env_vars_id  VARCHAR(26) NULL,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMPTZ NULL,

    CONSTRAINT fk_projects_settings_id FOREIGN KEY (settings_id) REFERENCES settings (id),
    CONSTRAINT fk_projects_env_vars_id FOREIGN KEY (env_vars_id) REFERENCES settings (id)
);

CREATE UNIQUE INDEX idx_uq_projects_name ON projects(LOWER(name)) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS projects;
