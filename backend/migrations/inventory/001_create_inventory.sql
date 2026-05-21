-- 库存余额表
CREATE TABLE IF NOT EXISTS inventory_balances (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    warehouse_id VARCHAR(64) NOT NULL,
    sku_id VARCHAR(64) NOT NULL,
    sku_code VARCHAR(128) DEFAULT '',
    total_quantity INT NOT NULL DEFAULT 0,
    locked_quantity INT NOT NULL DEFAULT 0,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_inv_bal_wh_sku ON inventory_balances(warehouse_id, sku_id);
CREATE INDEX idx_inv_bal_tenant_id ON inventory_balances(tenant_id);

-- 库存锁定记录表
CREATE TABLE IF NOT EXISTS inventory_locks (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    order_id VARCHAR(64) NOT NULL,
    sku_id VARCHAR(64) NOT NULL,
    warehouse_id VARCHAR(64) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    released_quantity INT NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'locked',
    lock_key VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_inv_locks_lock_key ON inventory_locks(lock_key);
CREATE INDEX idx_inv_locks_order_id ON inventory_locks(order_id);

-- 库存流水表
CREATE TABLE IF NOT EXISTS inventory_journals (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL DEFAULT '',
    warehouse_id VARCHAR(64) NOT NULL DEFAULT '',
    sku_id VARCHAR(64) NOT NULL DEFAULT '',
    order_id VARCHAR(64) DEFAULT '',
    change_type VARCHAR(16) NOT NULL,
    change_qty INT NOT NULL DEFAULT 0,
    before_total INT NOT NULL DEFAULT 0,
    after_total INT NOT NULL DEFAULT 0,
    before_locked INT NOT NULL DEFAULT 0,
    after_locked INT NOT NULL DEFAULT 0,
    before_avail INT NOT NULL DEFAULT 0,
    after_avail INT NOT NULL DEFAULT 0,
    idempotency_key VARCHAR(128) DEFAULT '',
    operator VARCHAR(64) DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_inv_jrnl_sku_id ON inventory_journals(sku_id);
CREATE INDEX idx_inv_jrnl_order_id ON inventory_journals(order_id);
CREATE INDEX idx_inv_jrnl_created_at ON inventory_journals(created_at DESC);
CREATE INDEX idx_inv_jrnl_idempotency_key ON inventory_journals(idempotency_key);
