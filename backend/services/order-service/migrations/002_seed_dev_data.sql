-- ============================================================
-- Order 订单服务 - 开发种子数据
-- ============================================================

-- 6 个销售订单（覆盖 pending/approved/shipped/cancelled 状态）
INSERT INTO orders.sales_orders (id, tenant_id, store_id, platform_order_no, status, buyer_name, buyer_email, total_amount, shipping_fee, idempotency_key, created_at, updated_at)
VALUES
    ('ord-001', 'default', 'store-amz-us', 'AMZ-20260601-0001', 'pending',   'Alice Wang',  'alice@example.com', 59.70, 5.99, 'ik-ord-001', '2026-06-28T08:00:00Z', '2026-06-28T08:00:00Z'),
    ('ord-002', 'default', 'store-amz-us', 'AMZ-20260601-0002', 'pending',   'Bob Chen',    'bob@example.com',   22.50, 3.99, 'ik-ord-002', '2026-06-28T12:00:00Z', '2026-06-28T12:00:00Z'),
    ('ord-003', 'default', 'store-amz-us', 'AMZ-20260602-0003', 'approved',  'Carol Liu',   'carol@example.com', 37.50, 4.99, 'ik-ord-003', '2026-06-29T09:00:00Z', '2026-06-29T14:00:00Z'),
    ('ord-004', 'default', 'store-amz-eu', 'AMZ-UK-20260602-01','approved',  'David Smith', 'david@example.com', 75.00, 8.99, 'ik-ord-004', '2026-06-29T10:00:00Z', '2026-06-29T15:00:00Z'),
    ('ord-005', 'default', 'store-amz-us', 'AMZ-20260603-0005', 'shipped',   'Eve Johnson', 'eve@example.com',   49.80, 5.99, 'ik-ord-005', '2026-06-27T06:00:00Z', '2026-06-30T08:00:00Z'),
    ('ord-006', 'default', 'store-shopify','SHOP-20260603-001', 'cancelled', 'Frank Brown', 'frank@example.com', 25.00, 0.00, 'ik-ord-006', '2026-06-27T15:00:00Z', '2026-06-28T10:00:00Z')
ON CONFLICT DO NOTHING;

-- 订单明细（每个订单 1-3 个 SKU）
INSERT INTO orders.order_items (id, order_id, sku_id, sku_code, sku_name, quantity, unit_price, total_price)
VALUES
    ('oi-001-1', 'ord-001', 'SKU-0001', 'TSHIRT-BLUE-M', '蓝色T恤 M码',   2, 19.90, 39.80),
    ('oi-001-2', 'ord-001', 'SKU-0003', 'MUG-WHITE-350',  '白色马克杯 350ml', 1, 12.50, 12.50),
    ('oi-001-3', 'ord-001', 'SKU-0005', 'CAP-BLACK-ADJ',  '黑色棒球帽 可调节', 1,  7.40,  7.40),
    ('oi-002-1', 'ord-002', 'SKU-0004', 'MUG-BLACK-350',  '黑色马克杯 350ml', 1, 12.50, 12.50),
    ('oi-002-2', 'ord-002', 'SKU-0006', 'CAP-WHITE-ADJ',  '白色棒球帽 可调节', 1, 10.00, 10.00),
    ('oi-003-1', 'ord-003', 'SKU-0002', 'TSHIRT-RED-L',   '红色T恤 L码',   1, 19.90, 19.90),
    ('oi-003-2', 'ord-003', 'SKU-0004', 'MUG-BLACK-350',  '黑色马克杯 350ml', 1, 12.50, 12.50),
    ('oi-003-3', 'ord-003', 'SKU-0005', 'CAP-BLACK-ADJ',  '黑色棒球帽 可调节', 1,  5.10,  5.10),
    ('oi-004-1', 'ord-004', 'SKU-0003', 'MUG-WHITE-350',  '白色马克杯 350ml', 3, 12.50, 37.50),
    ('oi-004-2', 'ord-004', 'SKU-0002', 'TSHIRT-RED-L',   '红色T恤 L码',   1, 19.90, 19.90),
    ('oi-004-3', 'ord-004', 'SKU-0006', 'CAP-WHITE-ADJ',  '白色棒球帽 可调节', 1, 17.60, 17.60),
    ('oi-005-1', 'ord-005', 'SKU-0001', 'TSHIRT-BLUE-M', '蓝色T恤 M码',   2, 19.90, 39.80),
    ('oi-005-2', 'ord-005', 'SKU-0005', 'CAP-BLACK-ADJ',  '黑色棒球帽 可调节', 1, 10.00, 10.00),
    ('oi-006-1', 'ord-006', 'SKU-0003', 'MUG-WHITE-350',  '白色马克杯 350ml', 2, 12.50, 25.00)
ON CONFLICT DO NOTHING;

-- 收货地址
INSERT INTO orders.order_addresses (id, order_id, contact_name, phone, country, state, city, street_line1, postal_code)
VALUES
    ('addr-001', 'ord-001', 'Alice Wang',  '+1-555-0101', 'US', 'CA', 'Los Angeles', '123 Main St', '90001'),
    ('addr-002', 'ord-002', 'Bob Chen',    '+1-555-0102', 'US', 'NY', 'New York',    '456 Park Ave', '10001'),
    ('addr-003', 'ord-003', 'Carol Liu',   '+1-555-0103', 'US', 'TX', 'Houston',     '789 Oak Dr', '77001'),
    ('addr-004', 'ord-004', 'David Smith', '+44-7700-9001','UK', 'ENG','London',      '10 Downing St', 'SW1A'),
    ('addr-005', 'ord-005', 'Eve Johnson', '+1-555-0105', 'US', 'WA', 'Seattle',     '321 Pine Rd', '98101'),
    ('addr-006', 'ord-006', 'Frank Brown', '+1-555-0106', 'US', 'FL', 'Miami',       '654 Beach Blvd', '33101')
ON CONFLICT DO NOTHING;
