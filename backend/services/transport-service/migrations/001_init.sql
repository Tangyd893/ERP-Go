-- ============================================================
-- TMS 物流服务 - 数据库初始化
-- ============================================================
CREATE SCHEMA IF NOT EXISTS transport;
SET search_path TO transport, public;

CREATE TABLE IF NOT EXISTS carriers (
    id        VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    name      VARCHAR(128) NOT NULL,
    code      VARCHAR(64)  NOT NULL,
    status    VARCHAR(16)  NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS carrier_services (
    id           VARCHAR(36) PRIMARY KEY,
    carrier_id   VARCHAR(36) NOT NULL,
    name         VARCHAR(128) NOT NULL,
    code         VARCHAR(64)  NOT NULL,
    service_type VARCHAR(16)  NOT NULL DEFAULT 'standard'
);

CREATE TABLE IF NOT EXISTS shipping_rules (
    id                VARCHAR(36) PRIMARY KEY,
    tenant_id         VARCHAR(36) NOT NULL,
    name              VARCHAR(128) NOT NULL,
    priority          INT DEFAULT 0,
    country_codes     JSONB DEFAULT '[]',
    min_weight        DECIMAL(8,3) DEFAULT 0,
    max_weight        DECIMAL(8,3) DEFAULT 999,
    carrier_service_id VARCHAR(36) NOT NULL
);

CREATE TABLE IF NOT EXISTS shipments (
    id            VARCHAR(36) PRIMARY KEY,
    tenant_id     VARCHAR(36) NOT NULL,
    order_id      VARCHAR(36) NOT NULL,
    outbound_id   VARCHAR(36) NOT NULL,
    carrier_code  VARCHAR(32) NOT NULL,
    service_code  VARCHAR(32) NOT NULL,
    tracking_no   VARCHAR(128) DEFAULT '',
    label_url     TEXT DEFAULT '',
    status        VARCHAR(24) NOT NULL DEFAULT 'pending',
    weight        DECIMAL(8,3) DEFAULT 0,
    shipping_cost DECIMAL(10,2) DEFAULT 0,
    currency      VARCHAR(8) DEFAULT 'USD',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ship_order_id ON shipments(order_id);
CREATE INDEX idx_ship_tracking_no ON shipments(tracking_no);

CREATE TABLE IF NOT EXISTS shipment_packages (
    id          VARCHAR(36) PRIMARY KEY,
    shipment_id VARCHAR(36) NOT NULL,
    tracking_no VARCHAR(128) DEFAULT '',
    weight      DECIMAL(8,3) DEFAULT 0,
    length      DECIMAL(8,3) DEFAULT 0,
    width       DECIMAL(8,3) DEFAULT 0,
    height      DECIMAL(8,3) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS tracking_records (
    id          VARCHAR(36) PRIMARY KEY,
    shipment_id VARCHAR(36) NOT NULL,
    status      VARCHAR(64) NOT NULL,
    description TEXT DEFAULT '',
    location    VARCHAR(255) DEFAULT '',
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_tr_shipment_id ON tracking_records(shipment_id);
