-- ============================================================
-- Order 订单服务 - 数据库初始化
-- ============================================================

CREATE SCHEMA IF NOT EXISTS orders;
SET search_path TO orders, public;

CREATE TABLE IF NOT EXISTS sales_orders (
    id               VARCHAR(36) PRIMARY KEY,
    tenant_id        VARCHAR(36)  NOT NULL,
    store_id         VARCHAR(36)  NOT NULL,
    platform_order_no VARCHAR(128) NOT NULL,
    order_type       VARCHAR(16)  NOT NULL DEFAULT 'normal',
    order_source     VARCHAR(16)  NOT NULL DEFAULT 'platform',
    status           VARCHAR(32)  NOT NULL DEFAULT 'pending',
    buyer_name       VARCHAR(128) DEFAULT '',
    buyer_email      VARCHAR(128) DEFAULT '',
    currency         VARCHAR(8)   DEFAULT 'USD',
    total_amount     DECIMAL(12,2) DEFAULT 0,
    shipping_fee     DECIMAL(10,2) DEFAULT 0,
    tax_amount       DECIMAL(10,2) DEFAULT 0,
    idempotency_key  VARCHAR(128) NOT NULL,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_so_idempotency ON sales_orders(idempotency_key);
CREATE INDEX idx_so_tenant_status ON sales_orders(tenant_id, status);
CREATE INDEX idx_so_store_id ON sales_orders(store_id);
CREATE UNIQUE INDEX idx_so_platform_order ON sales_orders(store_id, platform_order_no);

CREATE TABLE IF NOT EXISTS order_items (
    id           VARCHAR(36) PRIMARY KEY,
    order_id     VARCHAR(36)  NOT NULL,
    sku_id       VARCHAR(36)  NOT NULL,
    sku_code     VARCHAR(128) NOT NULL,
    sku_name     VARCHAR(255) NOT NULL,
    platform_sku VARCHAR(128) DEFAULT '',
    quantity     INT          NOT NULL DEFAULT 1,
    unit_price   DECIMAL(12,2) DEFAULT 0,
    total_price  DECIMAL(12,2) DEFAULT 0
);
CREATE INDEX idx_oi_order_id ON order_items(order_id);
CREATE INDEX idx_oi_sku_id ON order_items(sku_id);

CREATE TABLE IF NOT EXISTS order_addresses (
    id           VARCHAR(36) PRIMARY KEY,
    order_id     VARCHAR(36)  NOT NULL UNIQUE,
    contact_name VARCHAR(128) DEFAULT '',
    phone        VARCHAR(32)  DEFAULT '',
    email        VARCHAR(128) DEFAULT '',
    country      VARCHAR(64)  DEFAULT '',
    state        VARCHAR(64)  DEFAULT '',
    city         VARCHAR(64)  DEFAULT '',
    district     VARCHAR(64)  DEFAULT '',
    street_line1 VARCHAR(255) DEFAULT '',
    street_line2 VARCHAR(255) DEFAULT '',
    postal_code  VARCHAR(16)  DEFAULT ''
);

CREATE TABLE IF NOT EXISTS order_status_logs (
    id          SERIAL PRIMARY KEY,
    order_id    VARCHAR(36) NOT NULL,
    from_status VARCHAR(32) NOT NULL,
    to_status   VARCHAR(32) NOT NULL,
    operator    VARCHAR(64) DEFAULT '',
    remark      TEXT        DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_osl_order_id ON order_status_logs(order_id);
CREATE INDEX idx_osl_created_at ON order_status_logs(created_at DESC);
