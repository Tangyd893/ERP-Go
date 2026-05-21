-- 部门表
CREATE TABLE IF NOT EXISTS departments (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    org_id VARCHAR(64) NOT NULL,
    parent_id VARCHAR(64) DEFAULT '',
    name VARCHAR(128) NOT NULL,
    code VARCHAR(64) NOT NULL,
    manager_id VARCHAR(64) DEFAULT '',
    sort_order INT DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_depts_tenant_id ON departments(tenant_id);
CREATE INDEX idx_depts_org_id ON departments(org_id);
CREATE INDEX idx_depts_parent_id ON departments(parent_id);
