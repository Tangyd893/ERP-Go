-- ============================================================
-- Channel 渠道服务 - 数据库初始化
-- ============================================================

CREATE SCHEMA IF NOT EXISTS channel;
SET search_path TO channel, public;

CREATE TABLE IF NOT EXISTS stores (
    id            VARCHAR(36) PRIMARY KEY,
    tenant_id     VARCHAR(36)  NOT NULL,
    platform_code VARCHAR(32)  NOT NULL,
    site          VARCHAR(32)  NOT NULL,
    name          VARCHAR(128) NOT NULL,
    store_code    VARCHAR(128) NOT NULL,
    auth_token    TEXT         DEFAULT '',
    auth_status   VARCHAR(16)  NOT NULL DEFAULT 'unauthorized',
    auth_expiry   TIMESTAMPTZ,
    status        VARCHAR(16)  NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_stores_tenant_id ON stores(tenant_id);

CREATE TABLE IF NOT EXISTS sync_tasks (
    id           VARCHAR(36) PRIMARY KEY,
    tenant_id    VARCHAR(36) NOT NULL,
    store_id     VARCHAR(36) NOT NULL,
    task_type    VARCHAR(32) NOT NULL,
    status       VARCHAR(16) NOT NULL DEFAULT 'pending',
    total_count  INT DEFAULT 0,
    success_count INT DEFAULT 0,
    failed_count  INT DEFAULT 0,
    error_msg    TEXT DEFAULT '',
    started_at   TIMESTAMPTZ,
    finished_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_sync_tasks_store_id ON sync_tasks(store_id);

CREATE TABLE IF NOT EXISTS platform_api_logs (
    id            VARCHAR(36) PRIMARY KEY,
    store_id      VARCHAR(36) NOT NULL,
    action        VARCHAR(64) NOT NULL,
    request_url   TEXT NOT NULL,
    request_body  TEXT DEFAULT '',
    status_code   INT DEFAULT 0,
    response_body TEXT DEFAULT '',
    duration_ms   BIGINT DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_pal_store_id ON platform_api_logs(store_id);
CREATE INDEX idx_pal_created_at ON platform_api_logs(created_at DESC);

CREATE TABLE IF NOT EXISTS order_import_tasks (
    id              VARCHAR(36) PRIMARY KEY,
    tenant_id       VARCHAR(36) NOT NULL,
    store_id        VARCHAR(36) NOT NULL,
    import_type     VARCHAR(16) NOT NULL,
    file_name       VARCHAR(255) DEFAULT '',
    idempotency_key VARCHAR(128) NOT NULL,
    status          VARCHAR(16) NOT NULL DEFAULT 'pending',
    total_rows      INT DEFAULT 0,
    success_rows    INT DEFAULT 0,
    failed_rows     INT DEFAULT 0,
    error_msg       TEXT DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_oit_idempotency ON order_import_tasks(idempotency_key);
CREATE INDEX idx_oit_store_id ON order_import_tasks(store_id);
