-- ============================================================
-- Inventory 库存服务 - 开发种子数据
-- ============================================================

-- 6 条库存余额（对应 6 个 SKU，仓库 wh-default）
INSERT INTO inventory.inventory_balances (id, tenant_id, warehouse_id, sku_id, sku_code, total_quantity, locked_quantity, available_quantity)
VALUES
    ('inv-001', 'default', 'wh-default', 'SKU-0001', 'TSHIRT-BLUE-M',  200, 10, 190),
    ('inv-002', 'default', 'wh-default', 'SKU-0002', 'TSHIRT-RED-L',   150,  5, 145),
    ('inv-003', 'default', 'wh-default', 'SKU-0003', 'MUG-WHITE-350',  500, 20, 480),
    ('inv-004', 'default', 'wh-default', 'SKU-0004', 'MUG-BLACK-350',  300,  0, 300),
    ('inv-005', 'default', 'wh-default', 'SKU-0005', 'CAP-BLACK-ADJ',  100,  3,  97),
    ('inv-006', 'default', 'wh-default', 'SKU-0006', 'CAP-WHITE-ADJ',   80,  0,  80)
ON CONFLICT DO NOTHING;
