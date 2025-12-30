-- +migrate Up
CREATE TABLE IF NOT EXISTS deployment_logs
(
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    deployment_id   VARCHAR(100) NOT NULL,
    type            VARCHAR(255) NULL,
    data            TEXT NOT NULL,
    ts              TIMESTAMPTZ NULL
);

CREATE INDEX idx_deployment_logs_deployment_id ON deployment_logs(deployment_id);

-- +migrate Down
DROP TABLE IF EXISTS deployment_logs;
