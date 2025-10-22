-- +migrate Up
CREATE TABLE IF NOT EXISTS apps
(
    id             VARCHAR(26) PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    photo          VARCHAR(255) NULL,
    project_id     VARCHAR(26) NOT NULL,
    project_env_id VARCHAR(26) NULL,
    status         VARCHAR(100) NOT NULL,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by   VARCHAR(26) NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by   VARCHAR(26) NOT NULL,
    deleted_at   TIMESTAMPTZ NULL,

    CONSTRAINT fk_apps_project_id FOREIGN KEY (project_id) REFERENCES projects (id),
    CONSTRAINT fk_apps_project_env_id FOREIGN KEY (project_env_id) REFERENCES project_envs (id),
    CONSTRAINT fk_apps_created_by FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_apps_updated_by FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE UNIQUE INDEX idx_uq_apps_name ON apps(project_id, project_env_id, LOWER(name)) WHERE deleted_at IS NULL;
CREATE INDEX idx_apps_project_id ON apps(project_id);
CREATE INDEX idx_apps_created_at ON apps(created_at);
CREATE INDEX idx_apps_deleted_at ON apps(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS apps;
