-- ============================================================
-- WMS 开发种子数据（仅本地测试用）
-- 创建仓库、库位、示例出库单，使 PDA 首页拣货计数非零
-- ============================================================
SET search_path TO warehouse, public;

-- 基础仓库
INSERT INTO warehouses (id, tenant_id, name, code, status) VALUES
    ('wh-default', 'default', '默认仓', 'WH001', 'active')
ON CONFLICT (id) DO NOTHING;

-- 库区
INSERT INTO zones (id, warehouse_id, name, zone_type) VALUES
    ('zone-pick-01', 'wh-default', '拣货区A', 'pick'),
    ('zone-storage-01', 'wh-default', '存储区A', 'storage')
ON CONFLICT (id) DO NOTHING;

-- 库位
INSERT INTO locations (id, warehouse_id, zone_id, code, barcode, status) VALUES
    ('loc-a01-01', 'wh-default', 'zone-pick-01', 'A01-01', 'LOC-A01-01', 'available'),
    ('loc-a01-02', 'wh-default', 'zone-pick-01', 'A01-02', 'LOC-A01-02', 'available'),
    ('loc-a02-01', 'wh-default', 'zone-pick-01', 'A02-01', 'LOC-A02-01', 'available'),
    ('loc-b01-01', 'wh-default', 'zone-storage-01', 'B01-01', 'LOC-B01-01', 'available')
ON CONFLICT (id) DO NOTHING;

-- 示例出库单（created = 待拣货）
INSERT INTO outbound_orders (id, tenant_id, order_id, order_no, warehouse_id, status) VALUES
    ('ob-dev-001', 'default', 'ord-dev-001', 'SO20260630001', 'wh-default', 'created'),
    ('ob-dev-002', 'default', 'ord-dev-002', 'SO20260630002', 'wh-default', 'created'),
    ('ob-dev-003', 'default', 'ord-dev-003', 'SO20260630003', 'wh-default', 'created'),
    ('ob-dev-004', 'default', 'ord-dev-004', 'SO20260630004', 'wh-default', 'picking'),
    ('ob-dev-005', 'default', 'ord-dev-005', 'SO20260630005', 'wh-default', 'picked')
ON CONFLICT (id) DO NOTHING;

-- 出库单明细
INSERT INTO outbound_items (id, outbound_id, sku_id, sku_code, sku_name, quantity, location_id) VALUES
    ('oi-dev-001', 'ob-dev-001', 'sku-dev-001', 'SKU-0001', '测试商品A（蓝色）', 10, 'loc-a01-01'),
    ('oi-dev-002', 'ob-dev-001', 'sku-dev-002', 'SKU-0002', '测试商品B（红色）', 5, 'loc-a01-02'),
    ('oi-dev-003', 'ob-dev-002', 'sku-dev-003', 'SKU-0003', '测试商品C（大号）', 20, 'loc-a02-01'),
    ('oi-dev-004', 'ob-dev-003', 'sku-dev-001', 'SKU-0001', '测试商品A（蓝色）', 15, 'loc-a01-01'),
    ('oi-dev-005', 'ob-dev-003', 'sku-dev-004', 'SKU-0004', '测试商品D（小号）', 8, 'loc-b01-01'),
    ('oi-dev-006', 'ob-dev-004', 'sku-dev-002', 'SKU-0002', '测试商品B（红色）', 12, 'loc-a01-02'),
    ('oi-dev-007', 'ob-dev-005', 'sku-dev-003', 'SKU-0003', '测试商品C（大号）', 6, 'loc-a02-01')
ON CONFLICT (id) DO NOTHING;
