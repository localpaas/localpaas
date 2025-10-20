-- Users
INSERT INTO users (id, email, role, status, full_name, security_option, password, password_salt, created_at, updated_at)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WX8', 'tiendc@gmail.com', 'owner', 'active', 'Tien DC', 'password-only',
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;
