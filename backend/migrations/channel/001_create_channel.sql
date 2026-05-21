-- 店铺表
CREATE TABLE IF NOT EXISTS stores (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    platform VARCHAR(64) NOT NULL,
    site VARCHAR(64) DEFAULT '',
    name VARCHAR(256) NOT NULL,
    store_code VARCHAR(128) NOT NULL,
    auth_token VARCHAR(2048) DEFAULT '',
    auth_status VARCHAR(16) DEFAULT 'pending',
    auth_expires_at TIMESTAMPTZ,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_stores_tenant_id ON stores(tenant_id);
CREATE UNIQUE INDEX idx_stores_tenant_code ON stores(tenant_id, store_code);

-- 同步任务表
CREATE TABLE IF NOT EXISTS sync_tasks (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    store_id VARCHAR(64) NOT NULL,
    task_type VARCHAR(32) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    total_count INT DEFAULT 0,
    success_count INT DEFAULT 0,
    fail_count INT DEFAULT 0,
    error_msg TEXT DEFAULT '',
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_sync_tasks_store_id ON sync_tasks(store_id);

-- 平台API调用日志表
CREATE TABLE IF NOT EXISTS platform_api_logs (
    id VARCHAR(64) PRIMARY KEY,
    store_id VARCHAR(64) NOT NULL,
    operation VARCHAR(128) NOT NULL,
    request_url TEXT DEFAULT '',
    request_body TEXT DEFAULT '',
    status_code INT DEFAULT 0,
    response_body TEXT DEFAULT '',
    duration_ms INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_api_logs_store_id ON platform_api_logs(store_id);

-- 订单导入任务表
CREATE TABLE IF NOT EXISTS order_import_tasks (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    store_id VARCHAR(64) NOT NULL,
    import_type VARCHAR(16) NOT NULL DEFAULT 'csv',
    file_name VARCHAR(256) DEFAULT '',
    idempotency_key VARCHAR(128) NOT NULL UNIQUE,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    total_rows INT DEFAULT 0,
    success_rows INT DEFAULT 0,
    fail_rows INT DEFAULT 0,
    error_msg TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_import_tasks_store_id ON order_import_tasks(store_id);
