-- +migrate Up
CREATE TABLE IF NOT EXISTS project_envs
(
    id           VARCHAR(26) PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    project_id   VARCHAR(26) NOT NULL,

    created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by   VARCHAR(26) NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by   VARCHAR(26) NOT NULL,
    deleted_at   TIMESTAMPTZ NULL,

    CONSTRAINT fk_project_envs_project_id FOREIGN KEY (project_id) REFERENCES projects (id),
    CONSTRAINT fk_project_envs_created_by FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_project_envs_updated_by FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE UNIQUE INDEX idx_uq_project_envs_name ON project_envs(project_id, LOWER(name));
CREATE INDEX idx_project_envs_created_at ON project_envs(created_at);
CREATE INDEX idx_project_envs_deleted_at ON project_envs(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS project_envs;
