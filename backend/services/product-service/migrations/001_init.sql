-- ============================================================
-- Product 商品服务 - 数据库初始化
-- ============================================================

CREATE SCHEMA IF NOT EXISTS product;
SET search_path TO product, public;

CREATE TABLE IF NOT EXISTS spus (
    id          VARCHAR(36) PRIMARY KEY,
    tenant_id   VARCHAR(36)  NOT NULL,
    name        VARCHAR(255) NOT NULL,
    category_id VARCHAR(36)  DEFAULT '',
    brand       VARCHAR(128) DEFAULT '',
    description TEXT         DEFAULT '',
    images      JSONB        DEFAULT '[]',
    status      VARCHAR(16)  NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_spus_tenant_id ON spus(tenant_id);
CREATE INDEX idx_spus_category_id ON spus(category_id);

CREATE TABLE IF NOT EXISTS skus (
    id             VARCHAR(36) PRIMARY KEY,
    tenant_id      VARCHAR(36)  NOT NULL,
    spu_id         VARCHAR(36)  NOT NULL,
    code           VARCHAR(128) NOT NULL,
    name           VARCHAR(255) NOT NULL,
    barcode        VARCHAR(128) DEFAULT '',
    spec_desc      VARCHAR(255) DEFAULT '',
    weight         DECIMAL(10,3) DEFAULT 0,
    weight_unit    VARCHAR(8)   DEFAULT 'kg',
    length         DECIMAL(10,3) DEFAULT 0,
    width          DECIMAL(10,3) DEFAULT 0,
    height         DECIMAL(10,3) DEFAULT 0,
    length_unit    VARCHAR(8)   DEFAULT 'cm',
    purchase_price DECIMAL(12,2) DEFAULT 0,
    sale_price     DECIMAL(12,2) DEFAULT 0,
    currency       VARCHAR(8)   DEFAULT 'USD',
    status         VARCHAR(16)  NOT NULL DEFAULT 'active',
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_skus_tenant_code ON skus(tenant_id, code);
CREATE INDEX idx_skus_spu_id ON skus(spu_id);
CREATE INDEX idx_skus_barcode ON skus(barcode);

CREATE TABLE IF NOT EXISTS variant_options (
    id         VARCHAR(36) PRIMARY KEY,
    tenant_id  VARCHAR(36) NOT NULL,
    spu_id     VARCHAR(36) NOT NULL,
    name       VARCHAR(64) NOT NULL,
    value      VARCHAR(128) NOT NULL,
    sort_order INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS declaration_info (
    id             VARCHAR(36) PRIMARY KEY,
    sku_id         VARCHAR(36)  NOT NULL,
    cn_name        VARCHAR(255) DEFAULT '',
    en_name        VARCHAR(255) DEFAULT '',
    hs_code        VARCHAR(32)  DEFAULT '',
    material       VARCHAR(128) DEFAULT '',
    usage_desc     VARCHAR(255) DEFAULT '',
    unit_price     DECIMAL(12,2) DEFAULT 0,
    customs_weight DECIMAL(10,3) DEFAULT 0
);
CREATE UNIQUE INDEX idx_declaration_sku_id ON declaration_info(sku_id);

CREATE TABLE IF NOT EXISTS platform_sku_mapping (
    id              VARCHAR(36) PRIMARY KEY,
    tenant_id       VARCHAR(36) NOT NULL,
    sku_id          VARCHAR(36) NOT NULL,
    store_id        VARCHAR(36) NOT NULL,
    platform_code   VARCHAR(32) NOT NULL,
    platform_sku    VARCHAR(128) DEFAULT '',
    asin            VARCHAR(32)  DEFAULT '',
    fnsku           VARCHAR(32)  DEFAULT '',
    platform_status VARCHAR(32)  DEFAULT 'active',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_psm_sku_id ON platform_sku_mapping(sku_id);
CREATE INDEX idx_psm_store_id ON platform_sku_mapping(store_id);
CREATE UNIQUE INDEX idx_psm_sku_store ON platform_sku_mapping(sku_id, store_id);
