-- ============================================================
-- WMS 仓储服务 - 数据库初始化
-- ============================================================
CREATE SCHEMA IF NOT EXISTS warehouse;
SET search_path TO warehouse, public;

CREATE TABLE IF NOT EXISTS warehouses (
    id        VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    name      VARCHAR(128) NOT NULL,
    code      VARCHAR(64)  NOT NULL,
    address   TEXT DEFAULT '',
    status    VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS zones (
    id          VARCHAR(36) PRIMARY KEY,
    warehouse_id VARCHAR(36) NOT NULL,
    name        VARCHAR(64) NOT NULL,
    zone_type   VARCHAR(16) NOT NULL DEFAULT 'pick'
);

CREATE TABLE IF NOT EXISTS locations (
    id          VARCHAR(36) PRIMARY KEY,
    warehouse_id VARCHAR(36) NOT NULL,
    zone_id     VARCHAR(36) NOT NULL,
    code        VARCHAR(64) NOT NULL,
    barcode     VARCHAR(128) DEFAULT '',
    status      VARCHAR(16) NOT NULL DEFAULT 'available'
);
CREATE UNIQUE INDEX idx_loc_barcode ON locations(barcode);

CREATE TABLE IF NOT EXISTS outbound_orders (
    id           VARCHAR(36) PRIMARY KEY,
    tenant_id    VARCHAR(36) NOT NULL,
    order_id     VARCHAR(36) NOT NULL,
    order_no     VARCHAR(128) NOT NULL,
    warehouse_id VARCHAR(36) NOT NULL,
    status       VARCHAR(24) NOT NULL DEFAULT 'created',
    wave_id      VARCHAR(36) DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_oo_order_id ON outbound_orders(order_id);

CREATE TABLE IF NOT EXISTS outbound_items (
    id           VARCHAR(36) PRIMARY KEY,
    outbound_id  VARCHAR(36) NOT NULL,
    sku_id       VARCHAR(36) NOT NULL,
    sku_code     VARCHAR(128) NOT NULL,
    sku_name     VARCHAR(255) NOT NULL,
    quantity     INT NOT NULL,
    picked_quantity INT DEFAULT 0,
    checked_quantity INT DEFAULT 0,
    location_id  VARCHAR(36) DEFAULT ''
);
CREATE INDEX idx_oi_outbound_id ON outbound_items(outbound_id);

CREATE TABLE IF NOT EXISTS waves (
    id           VARCHAR(36) PRIMARY KEY,
    warehouse_id VARCHAR(36) NOT NULL,
    name         VARCHAR(64) NOT NULL,
    status       VARCHAR(16) NOT NULL DEFAULT 'created',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pick_tasks (
    id            VARCHAR(36) PRIMARY KEY,
    wave_id       VARCHAR(36) NOT NULL,
    outbound_id   VARCHAR(36) NOT NULL,
    sku_id        VARCHAR(36) NOT NULL,
    sku_code      VARCHAR(128) NOT NULL,
    sku_name      VARCHAR(255) NOT NULL,
    quantity      INT NOT NULL,
    picked_quantity INT DEFAULT 0,
    location_code VARCHAR(64) DEFAULT '',
    status        VARCHAR(24) NOT NULL DEFAULT 'pending',
    picker_id     VARCHAR(36) DEFAULT ''
);

CREATE TABLE IF NOT EXISTS packages (
    id           VARCHAR(36) PRIMARY KEY,
    outbound_id  VARCHAR(36) NOT NULL,
    tracking_no  VARCHAR(128) DEFAULT '',
    carrier_code VARCHAR(32) DEFAULT '',
    weight       DECIMAL(8,3) DEFAULT 0,
    length       DECIMAL(8,3) DEFAULT 0,
    width        DECIMAL(8,3) DEFAULT 0,
    height       DECIMAL(8,3) DEFAULT 0,
    label_url    TEXT DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
