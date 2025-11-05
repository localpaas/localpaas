-- +migrate Up
CREATE TABLE IF NOT EXISTS apps
(
    id           VARCHAR(26) PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    slug         VARCHAR(100) NOT NULL,
    photo        VARCHAR(255) NULL,
    project_id   VARCHAR(26) NOT NULL,
    parent_id    VARCHAR(26) NULL,
    status       VARCHAR(100) NOT NULL,
    note         VARCHAR(10000) NULL,
    settings_id  VARCHAR(26) NULL,
    env_vars_id  VARCHAR(26) NULL,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMPTZ NULL,

    CONSTRAINT fk_apps_project_id FOREIGN KEY (project_id) REFERENCES projects (id),
    CONSTRAINT fk_apps_settings_id FOREIGN KEY (settings_id) REFERENCES settings (id),
    CONSTRAINT fk_apps_env_vars_id FOREIGN KEY (env_vars_id) REFERENCES settings (id)
);

CREATE UNIQUE INDEX idx_uq_apps_name ON apps(project_id, LOWER(name)) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_uq_apps_slug ON apps(project_id, LOWER(slug)) WHERE deleted_at IS NULL;
CREATE INDEX idx_apps_project_id ON apps(project_id);
CREATE INDEX idx_apps_parent_id ON apps(parent_id);
CREATE INDEX idx_apps_created_at ON apps(created_at);
CREATE INDEX idx_apps_deleted_at ON apps(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS apps;
