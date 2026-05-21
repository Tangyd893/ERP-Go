-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    code VARCHAR(128) NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    resource_type VARCHAR(32) NOT NULL DEFAULT 'button',
    action VARCHAR(64) DEFAULT '',
    parent_id VARCHAR(64) DEFAULT '',
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_permissions_parent_id ON permissions(parent_id);
CREATE INDEX idx_permissions_resource_type ON permissions(resource_type);
