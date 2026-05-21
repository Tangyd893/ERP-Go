-- 组织表（树形结构，parent_id 为空表示根组织）
CREATE TABLE IF NOT EXISTS organizations (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    parent_id VARCHAR(64) DEFAULT '',
    name VARCHAR(128) NOT NULL,
    code VARCHAR(64) NOT NULL,
    sort_order INT DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orgs_tenant_id ON organizations(tenant_id);
CREATE INDEX idx_orgs_parent_id ON organizations(parent_id);
CREATE UNIQUE INDEX idx_orgs_tenant_code ON organizations(tenant_id, code);
