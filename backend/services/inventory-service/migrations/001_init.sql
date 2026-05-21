-- ============================================================
-- Inventory 库存服务 - 数据库初始化
-- ============================================================

CREATE SCHEMA IF NOT EXISTS inventory;
SET search_path TO inventory, public;

CREATE TABLE IF NOT EXISTS inventory_balances (
    id               VARCHAR(36) PRIMARY KEY,
    tenant_id        VARCHAR(36) NOT NULL,
    warehouse_id     VARCHAR(36) NOT NULL,
    sku_id           VARCHAR(36) NOT NULL,
    sku_code         VARCHAR(128) NOT NULL,
    total_quantity   INT DEFAULT 0,
    locked_quantity  INT DEFAULT 0,
    available_quantity INT DEFAULT 0,
    version          INT DEFAULT 1,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_ib_wh_sku ON inventory_balances(warehouse_id, sku_id);
CREATE INDEX idx_ib_tenant_sku ON inventory_balances(tenant_id, sku_id);

CREATE TABLE IF NOT EXISTS inventory_locks (
    id             VARCHAR(36) PRIMARY KEY,
    tenant_id      VARCHAR(36) NOT NULL,
    order_id       VARCHAR(36) NOT NULL,
    sku_id         VARCHAR(36) NOT NULL,
    warehouse_id   VARCHAR(36) NOT NULL,
    quantity       INT NOT NULL DEFAULT 0,
    released_quantity INT DEFAULT 0,
    status         VARCHAR(16) NOT NULL DEFAULT 'locked',
    lock_key       VARCHAR(128) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_il_lock_key ON inventory_locks(lock_key);
CREATE INDEX idx_il_order_id ON inventory_locks(order_id);
CREATE INDEX idx_il_status ON inventory_locks(status);

CREATE TABLE IF NOT EXISTS inventory_journals (
    id              VARCHAR(36) PRIMARY KEY,
    tenant_id       VARCHAR(36) NOT NULL,
    warehouse_id    VARCHAR(36) NOT NULL,
    sku_id          VARCHAR(36) NOT NULL,
    order_id        VARCHAR(36) DEFAULT '',
    change_type     VARCHAR(16) NOT NULL,
    change_qty      INT NOT NULL DEFAULT 0,
    before_total    INT NOT NULL DEFAULT 0,
    after_total     INT NOT NULL DEFAULT 0,
    before_locked   INT NOT NULL DEFAULT 0,
    after_locked    INT NOT NULL DEFAULT 0,
    before_avail    INT NOT NULL DEFAULT 0,
    after_avail     INT NOT NULL DEFAULT 0,
    idempotency_key VARCHAR(128) DEFAULT '',
    operator        VARCHAR(64)  DEFAULT '',
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_ij_idempotency ON inventory_journals(idempotency_key);
CREATE INDEX idx_ij_sku_id ON inventory_journals(sku_id);
CREATE INDEX idx_ij_created_at ON inventory_journals(created_at DESC);
CREATE INDEX idx_ij_change_type ON inventory_journals(change_type);
