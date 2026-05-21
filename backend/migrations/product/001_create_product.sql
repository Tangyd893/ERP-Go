-- SPU表（标准产品单元）
CREATE TABLE IF NOT EXISTS spus (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    name VARCHAR(256) NOT NULL,
    category VARCHAR(128) DEFAULT '',
    brand VARCHAR(128) DEFAULT '',
    description TEXT DEFAULT '',
    main_image VARCHAR(512) DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_spus_tenant_id ON spus(tenant_id);

-- SKU表（库存单位）
CREATE TABLE IF NOT EXISTS skus (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    spu_id VARCHAR(64) NOT NULL,
    sku_code VARCHAR(128) NOT NULL UNIQUE,
    barcode VARCHAR(128) DEFAULT '',
    spec_desc VARCHAR(512) DEFAULT '',
    weight_gram INT DEFAULT 0,
    length_cm NUMERIC(8,2) DEFAULT 0,
    width_cm NUMERIC(8,2) DEFAULT 0,
    height_cm NUMERIC(8,2) DEFAULT 0,
    purchase_price NUMERIC(12,2) DEFAULT 0,
    selling_price NUMERIC(12,2) DEFAULT 0,
    currency VARCHAR(8) DEFAULT 'CNY',
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_skus_tenant_id ON skus(tenant_id);
CREATE INDEX idx_skus_spu_id ON skus(spu_id);

-- 平台SKU映射表
CREATE TABLE IF NOT EXISTS platform_skus (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    sku_id VARCHAR(64) NOT NULL,
    store_id VARCHAR(64) NOT NULL DEFAULT '',
    platform_code VARCHAR(64) NOT NULL,
    platform_sku_id VARCHAR(128) DEFAULT '',
    asin VARCHAR(64) DEFAULT '',
    fnsku VARCHAR(64) DEFAULT '',
    platform_status VARCHAR(32) DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_plat_skus_sku_id ON platform_skus(sku_id);
CREATE INDEX idx_plat_skus_store_id ON platform_skus(store_id);
CREATE UNIQUE INDEX idx_plat_skus_unique ON platform_skus(tenant_id, store_id, platform_code);
