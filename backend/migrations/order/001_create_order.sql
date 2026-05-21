-- 销售订单表
CREATE TABLE IF NOT EXISTS sales_orders (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    store_id VARCHAR(64) NOT NULL,
    platform_order_no VARCHAR(128) NOT NULL,
    order_type VARCHAR(16) NOT NULL DEFAULT 'normal',
    order_source VARCHAR(16) NOT NULL DEFAULT 'csv',
    order_status VARCHAR(16) NOT NULL DEFAULT 'pending',
    buyer_name VARCHAR(128) DEFAULT '',
    buyer_email VARCHAR(256) DEFAULT '',
    currency VARCHAR(8) DEFAULT 'CNY',
    total_amount NUMERIC(12,2) DEFAULT 0,
    shipping_amount NUMERIC(12,2) DEFAULT 0,
    discount_amount NUMERIC(12,2) DEFAULT 0,
    actual_amount NUMERIC(12,2) DEFAULT 0,
    idempotency_key VARCHAR(128) NOT NULL UNIQUE,
    remark TEXT DEFAULT '',
    ordered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_orders_tenant_id ON sales_orders(tenant_id);
CREATE INDEX idx_orders_store_id ON sales_orders(store_id);
CREATE INDEX idx_orders_status ON sales_orders(order_status);
CREATE UNIQUE INDEX idx_orders_platform_no ON sales_orders(store_id, platform_order_no);

-- 订单明细表
CREATE TABLE IF NOT EXISTS order_items (
    id VARCHAR(64) PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL,
    sku_id VARCHAR(64) NOT NULL,
    sku_code VARCHAR(128) DEFAULT '',
    sku_name VARCHAR(256) DEFAULT '',
    quantity INT NOT NULL DEFAULT 1,
    unit_price NUMERIC(12,2) DEFAULT 0,
    total_price NUMERIC(12,2) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);

-- 订单地址表
CREATE TABLE IF NOT EXISTS order_addresses (
    id VARCHAR(64) PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL UNIQUE,
    contact_name VARCHAR(128) DEFAULT '',
    phone VARCHAR(64) DEFAULT '',
    email VARCHAR(256) DEFAULT '',
    country VARCHAR(64) DEFAULT '',
    state VARCHAR(128) DEFAULT '',
    city VARCHAR(128) DEFAULT '',
    district VARCHAR(128) DEFAULT '',
    address_line1 VARCHAR(512) DEFAULT '',
    address_line2 VARCHAR(512) DEFAULT '',
    postal_code VARCHAR(32) DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 订单状态变更记录表
CREATE TABLE IF NOT EXISTS order_status_logs (
    id VARCHAR(64) PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL,
    from_status VARCHAR(16) DEFAULT '',
    to_status VARCHAR(16) NOT NULL,
    operator VARCHAR(128) DEFAULT '',
    remark TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_status_logs_order_id ON order_status_logs(order_id);
