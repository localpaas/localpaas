-- Users
INSERT INTO users (id, username, email, role, status, full_name, security_option, totp_secret, password, password_salt, created_at, updated_at)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WA1', 'tiendc', 'tiendc@gmail.com', 'admin', 'active', 'Tien DC', 'password-only', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA2', 'member1', 'member1@domain.name', 'member', 'active', 'Member 1', 'password-only', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA3', 'member2', 'member2@domain.name', 'member', 'active', 'Member 2', 'password-only', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA4', 'member3', 'member3@domain.name', 'member', 'active', 'Member 3', 'password-2fa', 'AAAAAAAAAAAAAAAAAAAA',
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA5', 'member4', 'member4@domain.name', 'member', 'pending', 'Member 4', 'password-2fa', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO projects (id, name, slug, status, created_at, updated_at)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WB1', 'Project A', 'project_a', 'active',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WB2', 'Project B', 'project_b', 'active',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WB3', 'Project C', 'project_c', 'locked',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO project_tags (project_id, tag, display_order)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WB1', 'tag 1', 0),
       ('01JAB9XED0GTXBSQDFVYAJ8WB1', 'Tag 2', 1),
       ('01JAB9XED0GTXBSQDFVYAJ8WB2', 'Tag 3', 0),
       ('01JAB9XED0GTXBSQDFVYAJ8WB2', 'my tag', 1)
ON CONFLICT DO NOTHING;

INSERT INTO apps (id, name, slug, status, project_id, created_at, updated_at)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WD1', 'Backend', 'backend', 'active', '01JAB9XED0GTXBSQDFVYAJ8WB1',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WD2', 'Frontend', 'frontend', 'active', '01JAB9XED0GTXBSQDFVYAJ8WB1',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WD3', 'Redis', 'redis', 'active', '01JAB9XED0GTXBSQDFVYAJ8WB1',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WD4', 'Postgres', 'postgres', 'active', '01JAB9XED0GTXBSQDFVYAJ8WB1',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WD5', 'Backend', 'backend', 'active', '01JAB9XED0GTXBSQDFVYAJ8WB2',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WD6', 'Frontend', 'frontend', 'active', '01JAB9XED0GTXBSQDFVYAJ8WB2',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO app_tags (app_id, tag, display_order)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WD1', 'tag 1', 0),
       ('01JAB9XED0GTXBSQDFVYAJ8WD1', 'Tag 2', 1),
       ('01JAB9XED0GTXBSQDFVYAJ8WD2', 'Tag 3', 0),
       ('01JAB9XED0GTXBSQDFVYAJ8WD2', 'my tag', 1)
ON CONFLICT DO NOTHING;

INSERT INTO acl_permissions (subject_type, subject_id, resource_type, resource_id, action_read, action_write, action_delete, created_at, updated_at)
VALUES ('user', '01JAB9XED0GTXBSQDFVYAJ8WA1', 'project', '01JAB9XED0GTXBSQDFVYAJ8WB1', true, true, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA2', 'project', '01JAB9XED0GTXBSQDFVYAJ8WB1', true, true, false,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA3', 'project', '01JAB9XED0GTXBSQDFVYAJ8WB1', true, false, false,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA4', 'project', '01JAB9XED0GTXBSQDFVYAJ8WB1', true, true, false,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA1', 'project', '01JAB9XED0GTXBSQDFVYAJ8WB2', true, true, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA2', 'project', '01JAB9XED0GTXBSQDFVYAJ8WB1', true, true, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA2', 'project', '01JAB9XED0GTXBSQDFVYAJ8WB2', true, true, false,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA1', 'app', '01JAB9XED0GTXBSQDFVYAJ8WD1', true, true, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA2', 'app', '01JAB9XED0GTXBSQDFVYAJ8WD1', false, false, false,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;

-- Settings: OAuth
INSERT INTO settings (id, type, kind, name, status, data, created_at, updated_at)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WE1', 'oauth', 'github', 'Github', 'active',
        '{"org": "localpaas-test", "clientId": "Iv23liObQsEr3GigALXt", "clientSecret": "lpsalt:27AyosacFwi8CA== EsRjd1QQkZ8YmKOqWMUM2iMvrY+MWShSv9CV5whZbXKq6wjB8H36U2oLFBLbWD06hJ+DUSzvUmEWjk2XrxRlkyN4MV4="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WE2', 'oauth', 'gitlab', 'Gitlab', 'active',
        '{"org": "localpaas-test", "scopes": ["read_user", "openid", "profile", "email"], "clientId": "9a7d1422b34a79e83d74bee66448854d97452a2b6a7f05a74870f351f837dbbe", "clientSecret": "lpsalt:k+UQvPiR1yWRow== K8bthSPszSGA5YM2g5pNwreaNpFrYFn0TIOq+aJInN7vwrdOBJ+sAd1zy+RoRlr+CpiphDrPn77wn6hdJZl9tH2gNj/8Y8T07jhfibmKZpuFwRvbQOk01bcYQBWo8aXG8QY="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WE3', 'oauth', 'gitlab-custom', 'Our Gitlab', 'active',
        '{"org": "localpaas-test", "scopes": ["read_user", "openid", "profile", "email"], "authURL": "https://gitlab.com/oauth/authorize", "clientId": "f453848a3e717bad989baa8552afadb163b98cccc8b6bd4c3bb6b9f852fb4345", "tokenURL": "https://gitlab.com/oauth/token", "profileURL": "https://gitlab.com/api/v3/user", "clientSecret": "lpsalt:+gl7AsEBZskoZQ== xWNUdw2Vn5v4kV/c2rYZWDnyhBEA8+nWi+9rMADjiQP/gdODpRolVOvKbj9Gip+jxKdhmh+JoOJBq7XZ5r9/zgrt5Zfd+oQ8rgJzO2x7xHWa0j/U4Wy+lBAN7zTh9YNB5Eo="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;
