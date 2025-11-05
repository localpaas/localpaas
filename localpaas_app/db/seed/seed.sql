-- Users
INSERT INTO users (id, email, role, status, full_name, security_option, totp_secret, password, password_salt, created_at, updated_at)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WA1', 'tiendc@gmail.com', 'owner', 'active', 'Tien DC', 'password-only', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA2', 'member1@domain.name', 'member', 'active', 'Member 1', 'password-only', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA3', 'member2@domain.name', 'member', 'active', 'Member 2', 'password-2fa', '12345678901234567890',
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO nodes (id, host_name, ip, status, infra_status, is_leader, is_manager, created_at, updated_at)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WX1', 'node-a', '123.123.123.1', 'active', 'ready', true, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WX2', 'node-b', '123.123.123.2', 'active', 'ready', false, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WX3', 'node-c', '123.123.123.3', 'active', 'ready', false, false,
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
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO app_tags (app_id, tag, display_order)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WD1', 'tag 1', 0),
       ('01JAB9XED0GTXBSQDFVYAJ8WD1', 'Tag 2', 1),
       ('01JAB9XED0GTXBSQDFVYAJ8WD2', 'Tag 3', 0),
       ('01JAB9XED0GTXBSQDFVYAJ8WD2', 'my tag', 1)
ON CONFLICT DO NOTHING;
