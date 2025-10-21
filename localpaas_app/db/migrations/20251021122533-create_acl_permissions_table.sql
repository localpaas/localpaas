-- +migrate Up
CREATE TABLE IF NOT EXISTS acl_permissions
(
    user_id           VARCHAR(26) NOT NULL,
    resource_type     VARCHAR(100) NOT NULL,
    resource_id       VARCHAR(100),
    action_read       VARCHAR(10) NOT NULL CONSTRAINT chk_action_read CHECK (action_read IN ('yes', 'no')),
    action_write      VARCHAR(10) NOT NULL CONSTRAINT chk_action_write CHECK (action_write IN ('yes', 'no')),
    action_delete     VARCHAR(10) NOT NULL CONSTRAINT chk_action_delete CHECK (action_delete IN ('yes', 'no')),

    created_at        TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by        VARCHAR(26) NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by        VARCHAR(26) NOT NULL,

    PRIMARY KEY (user_id, resource_type, resource_id),
    CONSTRAINT fk_acl_permissions_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_acl_permissions_created_by FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_acl_permissions_updated_by FOREIGN KEY (updated_by) REFERENCES users (id)
);

-- +migrate Down
DROP TABLE IF EXISTS acl_permissions;
