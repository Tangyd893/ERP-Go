-- 租户表
CREATE TABLE IF NOT EXISTS tenants (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    code VARCHAR(64) NOT NULL UNIQUE,
    contact_name VARCHAR(128) DEFAULT '',
    contact_email VARCHAR(256) DEFAULT '',
    contact_phone VARCHAR(64) DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    quota_users INT DEFAULT 0,
    quota_orders INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_status ON tenants(status);
