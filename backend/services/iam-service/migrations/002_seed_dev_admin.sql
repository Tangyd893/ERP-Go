-- 开发环境默认账号（仅本地/测试使用）
-- 用户名: admin  密码: admin123  租户: default

SET search_path TO iam, public;

INSERT INTO iam_roles (id, tenant_id, name, code, description, status) VALUES
    ('role-default-admin', 'default', '管理员', 'admin', '开发环境默认管理员', 'active')
ON CONFLICT (tenant_id, code) DO NOTHING;

INSERT INTO iam_role_permissions (role_id, permission_id)
SELECT 'role-default-admin', id FROM iam_permissions
ON CONFLICT DO NOTHING;

INSERT INTO iam_users (id, tenant_id, username, password_hash, nickname, email, status) VALUES
    (
        'user-default-admin',
        'default',
        'admin',
        '$2a$10$7ASstAKbyxrbzuLSGxkKruAFMbDpZ/9xFqAaiQgroY50DEB4vi8Ee',
        '系统管理员',
        'admin@local.dev',
        'active'
    )
ON CONFLICT (tenant_id, username) DO NOTHING;

INSERT INTO iam_user_roles (user_id, role_id) VALUES
    ('user-default-admin', 'role-default-admin')
ON CONFLICT DO NOTHING;
