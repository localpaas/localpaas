-- +migrate Up
CREATE TABLE IF NOT EXISTS tasks
(
    id               VARCHAR(100) PRIMARY KEY,
    job_id           VARCHAR(100) NULL,
    type             VARCHAR(100) NOT NULL,
    status           VARCHAR NOT NULL CONSTRAINT chk_status CHECK
                        (status IN ('not-started','in-progress','canceled','failed','done')),
    config           JSONB NOT NULL,
    args             JSON NULL,
    runs             JSON NULL,
    output           JSON NULL,
    version          INT2 NOT NULL DEFAULT 1,
    update_ver       INT4 NOT NULL DEFAULT 1,

    run_at           TIMESTAMPTZ NULL,
    retry_at         TIMESTAMPTZ NULL,
    started_at       TIMESTAMPTZ NULL,
    ended_at         TIMESTAMPTZ NULL,

    created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMPTZ NULL,

    CONSTRAINT fk_tasks_job_id FOREIGN KEY (job_id) REFERENCES settings (id)
);

CREATE INDEX idx_tasks_job_id ON tasks(job_id);
CREATE INDEX idx_tasks_type ON tasks(type);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_run_at ON tasks(run_at);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS tasks;
