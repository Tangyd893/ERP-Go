-- 角色-权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id VARCHAR(64) NOT NULL,
    permission_id VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_perm_id ON role_permissions(permission_id);
