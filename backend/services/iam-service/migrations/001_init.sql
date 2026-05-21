-- ============================================================
-- IAM 认证权限服务 - 数据库初始化
-- ============================================================

-- 创建 Schema
CREATE SCHEMA IF NOT EXISTS iam;
SET search_path TO iam, public;

-- ============================================================
-- 用户表
-- ============================================================
CREATE TABLE IF NOT EXISTS iam_users (
    id          VARCHAR(36) PRIMARY KEY,
    tenant_id   VARCHAR(36)  NOT NULL,
    username    VARCHAR(64)  NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname    VARCHAR(64)  DEFAULT '',
    email       VARCHAR(128) DEFAULT '',
    phone       VARCHAR(32)  DEFAULT '',
    avatar      VARCHAR(255) DEFAULT '',
    status      VARCHAR(16)  NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_iam_users_tenant_username ON iam_users(tenant_id, username);
CREATE INDEX idx_iam_users_tenant_id ON iam_users(tenant_id);
CREATE INDEX idx_iam_users_status ON iam_users(status);

-- ============================================================
-- 角色表
-- ============================================================
CREATE TABLE IF NOT EXISTS iam_roles (
    id          VARCHAR(36) PRIMARY KEY,
    tenant_id   VARCHAR(36)  NOT NULL,
    name        VARCHAR(64)  NOT NULL,
    code        VARCHAR(64)  NOT NULL,
    description VARCHAR(255) DEFAULT '',
    status      VARCHAR(16)  NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_iam_roles_tenant_code ON iam_roles(tenant_id, code);
CREATE INDEX idx_iam_roles_tenant_id ON iam_roles(tenant_id);

-- ============================================================
-- 权限表
-- ============================================================
CREATE TABLE IF NOT EXISTS iam_permissions (
    id            VARCHAR(36) PRIMARY KEY,
    name          VARCHAR(64)  NOT NULL,
    code          VARCHAR(128) NOT NULL,
    description   VARCHAR(255) DEFAULT '',
    resource_type VARCHAR(32)  NOT NULL DEFAULT 'api',
    action        VARCHAR(32)  NOT NULL DEFAULT '',
    parent_id     VARCHAR(36)  DEFAULT '',
    sort_order    INT          DEFAULT 0,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_iam_permissions_code ON iam_permissions(code);

-- ============================================================
-- 用户角色关联表
-- ============================================================
CREATE TABLE IF NOT EXISTS iam_user_roles (
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_iam_user_roles_role_id ON iam_user_roles(role_id);

-- ============================================================
-- 角色权限关联表
-- ============================================================
CREATE TABLE IF NOT EXISTS iam_role_permissions (
    role_id       VARCHAR(36) NOT NULL,
    permission_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (role_id, permission_id)
);

-- ============================================================
-- 菜单表
-- ============================================================
CREATE TABLE IF NOT EXISTS iam_menus (
    id              VARCHAR(36) PRIMARY KEY,
    parent_id       VARCHAR(36)  DEFAULT '',
    name            VARCHAR(64)  NOT NULL,
    path            VARCHAR(255) DEFAULT '',
    component       VARCHAR(255) DEFAULT '',
    icon            VARCHAR(64)  DEFAULT '',
    sort_order      INT          DEFAULT 0,
    permission_code VARCHAR(128) DEFAULT '',
    visible         BOOLEAN      DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_iam_menus_parent_id ON iam_menus(parent_id);

-- ============================================================
-- 操作审计日志表
-- ============================================================
CREATE TABLE IF NOT EXISTS iam_audit_logs (
    id            VARCHAR(36) PRIMARY KEY,
    tenant_id     VARCHAR(36)  NOT NULL,
    user_id       VARCHAR(36)  NOT NULL,
    username      VARCHAR(64)  NOT NULL,
    action        VARCHAR(32)  NOT NULL,
    resource_type VARCHAR(32)  NOT NULL,
    resource_id   VARCHAR(36)  DEFAULT '',
    detail        TEXT         DEFAULT '',
    ip            VARCHAR(45)  DEFAULT '',
    user_agent    TEXT         DEFAULT '',
    request_id    VARCHAR(36)  DEFAULT '',
    trace_id      VARCHAR(36)  DEFAULT '',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_iam_audit_tenant_id ON iam_audit_logs(tenant_id);
CREATE INDEX idx_iam_audit_user_id ON iam_audit_logs(user_id);
CREATE INDEX idx_iam_audit_action ON iam_audit_logs(action);
CREATE INDEX idx_iam_audit_created_at ON iam_audit_logs(created_at DESC);

-- ============================================================
-- 初始化默认权限
-- ============================================================
INSERT INTO iam_permissions (id, name, code, description, resource_type, action) VALUES
    ('perm-user-read',     '查看用户', 'user:read',      '查看用户列表和详情', 'api', 'read'),
    ('perm-user-create',   '创建用户', 'user:create',    '创建新用户',         'api', 'create'),
    ('perm-user-update',   '编辑用户', 'user:update',    '编辑用户信息',       'api', 'update'),
    ('perm-user-disable',  '禁用用户', 'user:disable',   '禁用用户',           'api', 'update'),
    ('perm-user-enable',   '启用用户', 'user:enable',    '启用用户',           'api', 'update'),
    ('perm-user-assign-role', '分配角色', 'user:assign_role', '给用户分配角色', 'api', 'update'),
    ('perm-role-read',     '查看角色', 'role:read',      '查看角色列表和详情', 'api', 'read'),
    ('perm-role-create',   '创建角色', 'role:create',    '创建新角色',         'api', 'create'),
    ('perm-role-update',   '编辑角色', 'role:update',    '编辑角色信息',       'api', 'update'),
    ('perm-role-delete',   '删除角色', 'role:delete',    '删除角色',           'api', 'delete'),
    ('perm-role-assign-perm', '分配权限', 'role:assign_perm', '给角色分配权限', 'api', 'update')
ON CONFLICT (code) DO NOTHING;

-- ============================================================
-- 初始化超级管理员角色（租户 admin 初始创建时关联）
-- ============================================================
INSERT INTO iam_roles (id, tenant_id, name, code, description, status) VALUES
    ('role-super-admin', 'system', '超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 'active')
ON CONFLICT (tenant_id, code) DO NOTHING;

-- 超级管理员角色关联所有权限
INSERT INTO iam_role_permissions (role_id, permission_id)
SELECT 'role-super-admin', id FROM iam_permissions
ON CONFLICT DO NOTHING;
