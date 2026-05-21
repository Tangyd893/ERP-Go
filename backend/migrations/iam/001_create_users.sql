-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    username VARCHAR(128) NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    nickname VARCHAR(128) DEFAULT '',
    email VARCHAR(256) DEFAULT '',
    phone VARCHAR(64) DEFAULT '',
    avatar VARCHAR(512) DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_users_tenant_username ON users(tenant_id, username);
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_status ON users(status);
