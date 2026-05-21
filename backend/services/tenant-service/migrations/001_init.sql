-- ============================================================
-- Tenant 租户组织服务 - 数据库初始化
-- ============================================================

CREATE SCHEMA IF NOT EXISTS tenant;
SET search_path TO tenant, public;

-- ============================================================
-- 租户表
-- ============================================================
CREATE TABLE IF NOT EXISTS tenants (
    id            VARCHAR(36) PRIMARY KEY,
    name          VARCHAR(128) NOT NULL,
    code          VARCHAR(64)  NOT NULL,
    contact_name  VARCHAR(64)  DEFAULT '',
    contact_email VARCHAR(128) DEFAULT '',
    contact_phone VARCHAR(32)  DEFAULT '',
    status        VARCHAR(16)  NOT NULL DEFAULT 'active',
    quota_users   INT          DEFAULT 10,
    quota_orders  INT          DEFAULT 1000,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_tenants_code ON tenants(code);

-- ============================================================
-- 组织表
-- ============================================================
CREATE TABLE IF NOT EXISTS organizations (
    id         VARCHAR(36) PRIMARY KEY,
    tenant_id  VARCHAR(36)  NOT NULL,
    parent_id  VARCHAR(36)  DEFAULT '',
    name       VARCHAR(128) NOT NULL,
    code       VARCHAR(64)  NOT NULL,
    sort_order INT          DEFAULT 0,
    status     VARCHAR(16)  NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orgs_tenant_id ON organizations(tenant_id);
CREATE INDEX idx_orgs_parent_id ON organizations(parent_id);

-- ============================================================
-- 部门表
-- ============================================================
CREATE TABLE IF NOT EXISTS departments (
    id         VARCHAR(36) PRIMARY KEY,
    tenant_id  VARCHAR(36)  NOT NULL,
    org_id     VARCHAR(36)  NOT NULL,
    parent_id  VARCHAR(36)  DEFAULT '',
    name       VARCHAR(128) NOT NULL,
    code       VARCHAR(64)  NOT NULL,
    manager_id VARCHAR(36)  DEFAULT '',
    sort_order INT          DEFAULT 0,
    status     VARCHAR(16)  NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_depts_tenant_id ON departments(tenant_id);
CREATE INDEX idx_depts_org_id ON departments(org_id);
CREATE INDEX idx_depts_parent_id ON departments(parent_id);

-- ============================================================
-- 岗位表
-- ============================================================
CREATE TABLE IF NOT EXISTS positions (
    id         VARCHAR(36) PRIMARY KEY,
    tenant_id  VARCHAR(36)  NOT NULL,
    dept_id    VARCHAR(36)  NOT NULL,
    name       VARCHAR(128) NOT NULL,
    code       VARCHAR(64)  NOT NULL,
    sort_order INT          DEFAULT 0,
    status     VARCHAR(16)  NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_positions_tenant_id ON positions(tenant_id);
CREATE INDEX idx_positions_dept_id ON positions(dept_id);
