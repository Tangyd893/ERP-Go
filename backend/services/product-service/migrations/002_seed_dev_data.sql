-- ============================================================
-- Product 商品服务 - 开发种子数据
-- ============================================================

-- 3 个 SPU（标准产品单元）
INSERT INTO product.spus (id, tenant_id, name, category_id, brand, description, status)
VALUES
    ('spu-001', 'default', '经典纯棉T恤',     'cat-apparel', 'Uniqlo',  '100%纯棉圆领短袖T恤，多色可选', 'active'),
    ('spu-002', 'default', '陶瓷马克杯',       'cat-kitchen', 'IKEA',    '350ml陶瓷马克杯，微波炉可用', 'active'),
    ('spu-003', 'default', '运动棒球帽',       'cat-sports',  'Nike',    '透气网眼棒球帽，可调节头围', 'active')
ON CONFLICT DO NOTHING;

-- 6 个 SKU（库存保有单位，2 个/SPU）
INSERT INTO product.skus (id, tenant_id, spu_id, code, name, barcode, spec_desc, weight, sale_price, status)
VALUES
    ('SKU-0001', 'default', 'spu-001', 'TSHIRT-BLUE-M',  '蓝色T恤 M码',  '6901234560001', '颜色:蓝色 尺码:M',  0.200, 19.90, 'active'),
    ('SKU-0002', 'default', 'spu-001', 'TSHIRT-RED-L',   '红色T恤 L码',  '6901234560002', '颜色:红色 尺码:L',  0.220, 19.90, 'active'),
    ('SKU-0003', 'default', 'spu-002', 'MUG-WHITE-350',  '白色马克杯 350ml', '6901234560003', '颜色:白色 容量:350ml', 0.350, 12.50, 'active'),
    ('SKU-0004', 'default', 'spu-002', 'MUG-BLACK-350',  '黑色马克杯 350ml', '6901234560004', '颜色:黑色 容量:350ml', 0.350, 12.50, 'active'),
    ('SKU-0005', 'default', 'spu-003', 'CAP-BLACK-ADJ',  '黑色棒球帽 可调节', '6901234560005', '颜色:黑色 尺码:可调节', 0.120, 25.00, 'active'),
    ('SKU-0006', 'default', 'spu-003', 'CAP-WHITE-ADJ',  '白色棒球帽 可调节', '6901234560006', '颜色:白色 尺码:可调节', 0.120, 25.00, 'active')
ON CONFLICT DO NOTHING;

-- 平台 SKU 映射（Amazon）
INSERT INTO product.platform_sku_mapping (id, tenant_id, sku_id, store_id, platform_code, platform_sku, asin)
VALUES
    ('psm-001', 'default', 'SKU-0001', 'store-amz-us', 'amazon', 'TSHIRT-BLUE-M-US', 'B0ABCDEF01'),
    ('psm-002', 'default', 'SKU-0002', 'store-amz-us', 'amazon', 'TSHIRT-RED-L-US',  'B0ABCDEF02'),
    ('psm-003', 'default', 'SKU-0003', 'store-amz-us', 'amazon', 'MUG-WHITE-350-US', 'B0ABCDEF03')
ON CONFLICT DO NOTHING;
