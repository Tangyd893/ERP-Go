-- ============================================================
-- Purchase 采购服务 - 数据库初始化
-- ============================================================
CREATE SCHEMA IF NOT EXISTS purchase;
SET search_path TO purchase, public;

CREATE TABLE IF NOT EXISTS suppliers (
    id            VARCHAR(36) PRIMARY KEY,
    tenant_id     VARCHAR(36) NOT NULL,
    name          VARCHAR(128) NOT NULL,
    code          VARCHAR(64)  NOT NULL,
    contact_name  VARCHAR(64) DEFAULT '',
    contact_phone VARCHAR(32) DEFAULT '',
    email         VARCHAR(128) DEFAULT '',
    payment_term  VARCHAR(32) DEFAULT 'net30',
    status        VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS purchase_orders (
    id            VARCHAR(36) PRIMARY KEY,
    tenant_id     VARCHAR(36) NOT NULL,
    supplier_id   VARCHAR(36) NOT NULL,
    supplier_name VARCHAR(128) DEFAULT '',
    order_no      VARCHAR(128) NOT NULL,
    status        VARCHAR(24) NOT NULL DEFAULT 'draft',
    currency      VARCHAR(8) DEFAULT 'USD',
    total_amount  DECIMAL(12,2) DEFAULT 0,
    expected_date DATE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_po_order_no ON purchase_orders(order_no);

CREATE TABLE IF NOT EXISTS purchase_items (
    id            VARCHAR(36) PRIMARY KEY,
    order_id      VARCHAR(36) NOT NULL,
    sku_id        VARCHAR(36) NOT NULL,
    sku_code      VARCHAR(128) NOT NULL,
    sku_name      VARCHAR(255) NOT NULL,
    quantity      INT NOT NULL,
    received_quantity INT DEFAULT 0,
    unit_price    DECIMAL(12,2) DEFAULT 0,
    total_price   DECIMAL(12,2) DEFAULT 0
);
CREATE INDEX idx_pi_order_id ON purchase_items(order_id);

CREATE TABLE IF NOT EXISTS inbound_orders (
    id           VARCHAR(36) PRIMARY KEY,
    tenant_id    VARCHAR(36) NOT NULL,
    purchase_id  VARCHAR(36) NOT NULL,
    warehouse_id VARCHAR(36) NOT NULL,
    status       VARCHAR(24) NOT NULL DEFAULT 'receiving',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inbound_items (
    id             VARCHAR(36) PRIMARY KEY,
    inbound_id     VARCHAR(36) NOT NULL,
    sku_id         VARCHAR(36) NOT NULL,
    quantity       INT NOT NULL,
    received_quantity INT DEFAULT 0,
    passed_quantity   INT DEFAULT 0,
    rejected_quantity INT DEFAULT 0
);
CREATE INDEX idx_ii_inbound_id ON inbound_items(inbound_id);
