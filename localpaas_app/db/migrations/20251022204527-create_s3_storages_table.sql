-- +migrate Up
CREATE TABLE IF NOT EXISTS s3_storages
(
    id                VARCHAR(26) PRIMARY KEY,
    name              VARCHAR(100) NOT NULL,
    access_key_id     VARCHAR(100) NOT NULL,
    secret_access_key BYTEA NOT NULL,
    salt              BYTEA NULL,
    region            VARCHAR(100) NULL,
    bucket            VARCHAR(100) NULL,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ NULL
);

-- +migrate Down
DROP TABLE IF EXISTS s3_storages;
