-- ============================================================
-- Channel 渠道服务 - 开发种子数据
-- ============================================================

-- 3 个店铺授权
INSERT INTO channel.stores (id, tenant_id, platform_code, site, name, store_code, auth_status, status)
VALUES
    ('store-amz-us', 'default', 'amazon',   'US', 'Amazon 美国站', 'AMZ-US-001', 'authorized', 'active'),
    ('store-amz-eu', 'default', 'amazon',   'UK', 'Amazon 英国站', 'AMZ-UK-001', 'authorized', 'active'),
    ('store-shopify','default', 'shopify',  'US', 'Shopify 独立站', 'SHOP-001',  'authorized', 'active')
ON CONFLICT DO NOTHING;
