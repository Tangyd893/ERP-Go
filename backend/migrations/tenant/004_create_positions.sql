-- 岗位表
CREATE TABLE IF NOT EXISTS positions (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    dept_id VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    code VARCHAR(64) NOT NULL,
    sort_order INT DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_positions_tenant_id ON positions(tenant_id);
CREATE INDEX idx_positions_dept_id ON positions(dept_id);
