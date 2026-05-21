-- ============================================================
-- Finance 财务服务 - 数据库初始化
-- ============================================================
CREATE SCHEMA IF NOT EXISTS finance;
SET search_path TO finance, public;

CREATE TABLE IF NOT EXISTS settlement_bills (
    id                VARCHAR(36) PRIMARY KEY,
    tenant_id         VARCHAR(36) NOT NULL,
    store_id          VARCHAR(36) NOT NULL,
    platform_code     VARCHAR(32) NOT NULL,
    settlement_period VARCHAR(64) NOT NULL,
    currency          VARCHAR(8) DEFAULT 'USD',
    total_sales       DECIMAL(12,2) DEFAULT 0,
    total_refunds     DECIMAL(12,2) DEFAULT 0,
    commission        DECIMAL(12,2) DEFAULT 0,
    fba_fee           DECIMAL(12,2) DEFAULT 0,
    other_fee         DECIMAL(12,2) DEFAULT 0,
    net_amount        DECIMAL(12,2) DEFAULT 0,
    status            VARCHAR(16) NOT NULL DEFAULT 'pending',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ar_ap_records (
    id            VARCHAR(36) PRIMARY KEY,
    tenant_id     VARCHAR(36) NOT NULL,
    type          VARCHAR(16) NOT NULL,
    order_id      VARCHAR(36) NOT NULL,
    amount        DECIMAL(12,2) DEFAULT 0,
    currency      VARCHAR(8) DEFAULT 'USD',
    exchange_rate DECIMAL(10,6) DEFAULT 1.0,
    amount_cny    DECIMAL(12,2) DEFAULT 0,
    status        VARCHAR(16) NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_arap_order_id ON ar_ap_records(order_id);

CREATE TABLE IF NOT EXISTS cost_records (
    id         VARCHAR(36) PRIMARY KEY,
    tenant_id  VARCHAR(36) NOT NULL,
    order_id   VARCHAR(36) NOT NULL,
    sku_id     VARCHAR(36) NOT NULL,
    cost_type  VARCHAR(16) NOT NULL,
    amount     DECIMAL(12,2) DEFAULT 0,
    currency   VARCHAR(8) DEFAULT 'USD',
    amount_cny DECIMAL(12,2) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_cr_order_id ON cost_records(order_id);

CREATE TABLE IF NOT EXISTS profit_reports (
    id              VARCHAR(36) PRIMARY KEY,
    tenant_id       VARCHAR(36) NOT NULL,
    order_id        VARCHAR(36) NOT NULL,
    order_no        VARCHAR(128) NOT NULL,
    sku_id          VARCHAR(36) NOT NULL,
    sku_code        VARCHAR(128) NOT NULL,
    sale_amount     DECIMAL(12,2) DEFAULT 0,
    purchase_cost   DECIMAL(12,2) DEFAULT 0,
    shipping_cost   DECIMAL(12,2) DEFAULT 0,
    commission_cost DECIMAL(12,2) DEFAULT 0,
    other_cost      DECIMAL(12,2) DEFAULT 0,
    total_cost      DECIMAL(12,2) DEFAULT 0,
    gross_profit    DECIMAL(12,2) DEFAULT 0,
    profit_margin   DECIMAL(6,4) DEFAULT 0,
    currency        VARCHAR(8) DEFAULT 'USD',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_pr_order_id ON profit_reports(order_id);
CREATE INDEX idx_pr_sku_id ON profit_reports(sku_id);

CREATE TABLE IF NOT EXISTS finance_journals (
    id              VARCHAR(36) PRIMARY KEY,
    tenant_id       VARCHAR(36) NOT NULL,
    order_id        VARCHAR(36) NOT NULL,
    change_type     VARCHAR(16) NOT NULL,
    amount          DECIMAL(12,2) DEFAULT 0,
    before_amount   DECIMAL(12,2) DEFAULT 0,
    after_amount    DECIMAL(12,2) DEFAULT 0,
    currency        VARCHAR(8) DEFAULT 'USD',
    idempotency_key VARCHAR(128) DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_fj_idempotency ON finance_journals(idempotency_key);
CREATE INDEX idx_fj_order_id ON finance_journals(order_id);
