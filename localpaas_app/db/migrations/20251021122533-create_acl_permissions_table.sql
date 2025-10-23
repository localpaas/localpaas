-- +migrate Up
CREATE TABLE IF NOT EXISTS acl_permissions
(
    subject_type   VARCHAR(100) NOT NULL,
    subject_id     VARCHAR(100) NOT NULL,
    resource_type  VARCHAR(100) NOT NULL,
    resource_id    VARCHAR(100) NOT NULL,
    action_read    BOOL NOT NULL DEFAULT FALSE,
    action_write   BOOL NOT NULL DEFAULT FALSE,
    action_delete  BOOL NOT NULL DEFAULT FALSE,

    created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMPTZ NULL,

    PRIMARY KEY (subject_id, resource_id)
);

-- +migrate Down
DROP TABLE IF EXISTS acl_permissions;
