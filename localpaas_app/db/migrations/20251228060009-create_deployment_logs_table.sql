-- +migrate Up
CREATE TABLE IF NOT EXISTS deployment_logs
(
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    deployment_id   VARCHAR(100) NOT NULL,
    step            VARCHAR(255) NULL,
    content         TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_deployment_logs_deployment_id ON deployment_logs(deployment_id);
CREATE INDEX idx_deployment_logs_created_at ON deployment_logs(created_at);

-- +migrate Down
DROP TABLE IF EXISTS deployment_logs;
