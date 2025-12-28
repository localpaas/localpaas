-- Users
INSERT INTO users (id, username, email, role, status, full_name, position, security_option, totp_secret, password, password_salt, created_at, updated_at)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WA1', 'tiendc', 'tiendc@gmail.com', 'admin', 'active', 'Tien DC', 'manager', 'password-only', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA2', 'member1', 'member1@domain.name', 'member', 'active', 'Member 1', 'devops', 'password-only', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA3', 'member2', 'member2@domain.name', 'member', 'active', 'Member 2', 'developer', 'password-only', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA4', 'member3', 'member3@domain.name', 'member', 'active', 'Member 3', NULL, 'password-2fa', 'AAAAAAAAAAAAAAAAAAAA',
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA5', 'member4', 'member4@domain.name', 'member', 'pending', 'Member 4', NULL, 'password-2fa', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA6', 'member5', 'member5@domain.name', 'member', 'pending', 'Member 5', NULL, 'password-2fa', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WA7', 'admin1', 'admin1@domain.name', 'admin', 'active', 'Admin 1', 'developer', 'password-only', NULL,
        '\x9e3e99b9f3ba5e6b934715e887cf423e5cfa80259ccb77ed5681e158b0fc0c8e',	'\x1a8594be97a4ddc71c86f19e3cf9f10c',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO projects (id, name, key, status, created_at, updated_at)
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

INSERT INTO apps (id, name, key, status, project_id, created_at, updated_at)
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
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA2', 'module', 'mod::user', true, true, false,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA2', 'module', 'mod::project', true, true, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA2', 'module', 'mod::settings', true, true, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA2', 'module', 'mod::provider', true, true, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('user', '01JAB9XED0GTXBSQDFVYAJ8WA2', 'module', 'mod::cluster', true, true, true,
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;

-- Settings: Providers
INSERT INTO settings (id, type, kind, name, status, data, created_at, updated_at)
VALUES ('01JAB9XED0GTXBSQDFVYAJ8WE1', 'oauth', 'github', 'Github', 'active',
        '{"org": "localpaas-test", "clientId": "Ov23lirztQpWxZTKNcTQ", "clientSecret": "lpsalt:i+NlaPQDkZ5LZQ== ITUM2K0dxQTb5D0DvCeEiHlO1vWzL5807TPfKH0E/37TBOoowSgogEYIp7leyYL7QGEfPDpM2cxb3+8fnlAeU6qlNQc="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WE2', 'oauth', 'gitlab', 'Gitlab', 'active',
        '{"org": "localpaas-test", "scopes": ["read_user", "openid", "profile", "email"], "clientId": "e675caa20a165c7ede6430bffcab512300d45bf0a9dd2da2acd122eb01995b69", "clientSecret": "lpsalt:exjDAd0nCxKxIw== erFzdQnYN1YEO6x8SnLcKxABt98TDSwfnfoYEhXfFR1dBhXt8CO9WAJHM8Yif+ANIZ+8ujvK+A/XFcYmDrDnuMGeCq1CUeralkf5KKklZgSWBL4ovmDhwMC2HMjMYFZBUA4="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WE3', 'oauth', 'gitlab-custom', 'Our Gitlab', 'active',
        '{"org": "localpaas-test", "scopes": ["read_user", "openid", "profile", "email"], "authURL": "https://gitlab.com/oauth/authorize", "clientId": "7e838ca826d84ebacaa778ad95ba97e13193480f729b367c098b37145eea1970", "tokenURL": "https://gitlab.com/oauth/token", "profileURL": "https://gitlab.com/api/v3/user", "clientSecret": "lpsalt:ghfPLkgStszb1Q== i/sa+hSF89q+UI+w6BF9FBc0HU1+X4KWrrU9qqN0h2uZITxElaXW2SJvyQvI4LsoRRvRduecEHC/Us2do6AFyhTFc3i0wmhPd1/FFDL0FzTQW66z8m+/KRjz8K0lPvxCHJY="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WE4', 'oauth', 'gitea', 'Gitea', 'active',
        '{"org": "localpaas-test", "clientId": "59140514-9f0f-4198-b13a-f5a958d0d024", "clientSecret": "lpsalt:8QaqYOW1kPxoFA== jYU8xNyQCTacGB+cNmEUKYeM6WFkbEYizYFEAcf8x90quocKHuOW2Hif1WWDhZcR15B71YuUpjJSUjx0juNL48EYmosDhd0Im9qeIUqF/xeIESgP"}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WE5', 'oauth', 'google', 'Google', 'active',
        '{"org": "localpaas-test", "clientId": "405460302846-8sk3j75rd5jonn9jfheis9coatkdn4jh.apps.googleusercontent.com", "clientSecret": "lpsalt:Y5ZYgPBmj1dfgw== DPKqVXvFlKC0DCkBRkzK3rWd4W8HTaIMJxLJKjriXtwKMTRVmxPd6fmHN1gBP8K3T7gIc0zvTuFohl3hYVgI"}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WF1', 'ssh-key', NULL, 'ssh-key srv 2', 'active',
        '{"privateKey": "lpsalt:zaPpoaE4yijCMg== 01gk7TPlyC1TjWhr6GNrab2RNQSk1c1L8YxNs6dRuWA8uAS73567W/qIOAXnq6kOsM57GEjQ3GuH+m2U6LNZ+dGRRtrVsrGmdaQYUuMiAP6eCMHGrEdRrBcfEvWfK3t5MnHEqnkouYlQc+FynuliI8phqWU6ITBWYynaCnsh8eBoPoO9g6ZOLwyVAggEfcVHynv0RulX2X+T/FwLztIAG3Job7ULVE1r53rKALp4IUVUJ/CXzxGp9JJJFQhagokSp/hZY65EWKCB+8z+RlDfz0pkNx6ZhTSP3nmv1a9I9D0KykI79A9wPPa7KMdGanFE13r4tXNaGHGFeCHWM+bA4PgoYOC/d9qFE5Qx7yD9GE8bmBMpVh4ArWG1N22hMSWeYo2niWhGpi6u1nwDjCNrlOGl5HL6mOE5jlxp9glqfY4gD/Vtu0qIV7nkX6fJAVSs8GIgoWgwQO/gUcGNu9tR6RXp7gD2jtN+tSy6QpX1UtsZEere56LJuF0vAUyRIqSDYHCvchchul+GObj06he3cmZS8Koqsd5qI/iTKtEsAAt/9BW3Wjq3vwROy8wgqpynssE5JV6w"}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WF2', 'ssh-key', NULL, 'ssh-key srv 2 with passphrase', 'active',
        '{"passphrase": "lpsalt:rHHxHbmV9Qx+hw== ffetB1FLIKffsFC1Gs3JVrUaB1oHIUvH3QwBuDMJwXFENg==", "privateKey": "lpsalt:8xPNWZPCzSnqvw== rdNzO68zrXA+ULr4VFpsaGowgI4evk3imiKg6V4An3T3qD8FEyapqLFhA2jsyotgbOJp1RV2iBgLRdZJW/Ovn1tKIsEUNkbx0wt1dPQf8Tgrnp78/Gff0BDt27lXatmRzSzV7L3ixTEah/MbZ4eUgZtwpYMeRGvGQfgwXPEVgvzq9RTGDbUm+00SWOZz96Y7m3eil3iH1di3AiGzYvt2L9PcaDTTAFAxY9v1Ru3cY7p3z9LlW91tiFU/gc7ywVaZDhxzf6nlZCLRl5tnq1EOI9gRFF9K0MOZvp/YMKvtHBBxn+8kx97V2Vs9ukCYdXkF9fQ7ddbcA6+0sDePIq7YCEyLdKzcBAg2NYyjpZsaEZaDQtlo9ZOYDoxi1cy9R6BZLfchFCdKSQgST07mAvRPE9x2hwJTT0NTO1GyJwBnyYJvjpmlL42nW1p3IW75kuaD+pgZqbRUxWVuyXMs4EAsHFV7xOXsBAqotxEZ6RwOrZGsw2plBEpLu+/auilCdOlsvBp0mT3A+ocAXULUByQ7LwYi76Rwzru0JOjKjMUCG99OqKc2WCwTulrm2cVrszsn5aK7kqVweFejlzRhkl2UXQSYi9nj/hD3tPhImML+3YRmub8uHzIRBIwli4DNz+Azk1jgSaztuqFxKhI="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WG1', 'registry-auth', 'docker.io', 'docker.io 1', 'active',
        '{"address": "docker.io", "password": "lpsalt:7rJsK4//cIqCKw== SxnfW00e9UueetjGT0dpuC9pUqqWR3fEiX1w+YZAo6N4E9+bNp3itGnvWR5tyblz1lJwjWAXnleyXiqtk4Rxsw==", "username": "localpaastest"}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WG2', 'registry-auth', 'docker.io', 'docker.io ro', 'active',
        '{"address": "docker.io", "password": "lpsalt:xpbKUbDCp7mDbQ== zzzou62KRIRxm5mfg4FFWivJGaVH19FhuglZa2bfOAhM7UDsxfCx2Z4HnRJUOqSPv3Nc/AEb9iCzrFzvd5sFNA==", "username": "localpaastest"}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WH1', 'slack', NULL, 'webhook 1', 'active',
        '{"webhook": "lpsalt:DZo0n5FM9NmYuA== fWmMb9cLMAjgsYuaH4CJTeA74SDvHXko7EPNv1LCv5o9XSmhRdi88PKb3GYPipq7pk/GErXEX9OIc5oaJ3K4GZXQk9dTAo6wqhiSke8Db3LvP/iNdnieWaZEyOIyGTSOjMLSSRsqXTlCmcTFfQ=="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WI1', 'discord', NULL, 'webhook 1', 'active',
        '{"webhook": "lpsalt:Kbdv16CTSIEPug== SwC3ghbKUQ7Lbq87+rQFH0JHUuM6/g6ZpiiSfxF1Q/R9RxLS7AbnqFOe486Y2LtP38T9ePG5Vjw7ieqLKAVq+2kkiVDSjMejOqe872a2Op3EPM4Idl73+0XJDLN+2MkzEbMHsnqlTpr9ddskzUdhPXdV8hTPKgJ5gFgX73ffv2r8B74bVqDZrhqxoqKsbg4Ft4Rm9XU="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WJ1', 'github-app', 'github-app', 'localpaas-test', 'active',
        '{"org": "localpaas-test", "appId": 2258661, "installationId": 96871321, "clientId": "Iv23liObQsEr3GigALXt", "privateKey": "lpsalt:A/GRN+CXoIFs9g== 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", "webhookURL": "webhookURL", "clientSecret": "lpsalt:jgcbPJYT+5mDZA== 7gv/DQb03hObXjEczacy65jRPE3foRpbuIwfEjiOua260cHeXJf7F4EFp4g1Vzfrt62jNKoEmDNdF+z5SweuGp7qjyg=", "webhookSecret": "lpsalt:6RX2NHe0VzpnqA== pPBiY0Gqx2MKm8e5SWT2vA+7x9Eo+eu8+rCUiAtYgO3TPE8RnKm0wws=", "ssoEnabled": true}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WK1', 'git-token', 'github', 'localpaas-test', 'active',
        '{"user": "tiendc", "token": "lpsalt:Plu+yhxw1HbmNA== rBEZCum9OdYsraYteP+Q3Y7aHrnZIxWjb88oYWMJ+CQuLfmHNCcPJtACeXXeq0r2MSqf8jD9D/esezjy6HwS61jfFpTHG1ph0tAzMeYc4U7I7WsjJsko8UbCWPtbEjKUXUZNkjw7nWJ+oxfr8j2TLLuUg+kcdM9Ghg=="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WK2', 'git-token', 'gitlab', 'localpaas-test', 'active',
        '{"user": "tiendc", "token": "lpsalt:gm8tw+8Lc0bivA== C8x6lh/XiXzJV11bt5Qx8uTn8/nctmtx+V8siciv7vQjfQeYcrbU0sB2CexZXR8Xv7k4J8DDaCeJd6jioBn9A2buiraTHGr30LLUa2VJkHVA5ATVAQ=="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00'),
       ('01JAB9XED0GTXBSQDFVYAJ8WK3', 'git-token', 'gitlab-custom', 'localpaas-test', 'active',
        '{"user": "tiendc", "baseURL": "https://gitlab.com/api/v4", "token": "lpsalt:gm8tw+8Lc0bivA== C8x6lh/XiXzJV11bt5Qx8uTn8/nctmtx+V8siciv7vQjfQeYcrbU0sB2CexZXR8Xv7k4J8DDaCeJd6jioBn9A2buiraTHGr30LLUa2VJkHVA5ATVAQ=="}',
        '2025-10-01 00:00:00', '2025-10-01 00:00:00')
ON CONFLICT DO NOTHING;
