-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL DEFAULT '',
    name VARCHAR(128) NOT NULL,
    code VARCHAR(64) NOT NULL,
    description TEXT DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_roles_tenant_code ON roles(tenant_id, code);
CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
