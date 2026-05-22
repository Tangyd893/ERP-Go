# 数据库迁移目录

本目录保留各服务的建表 SQL 作为历史参考和聚合查阅用途。

**当前迁移策略**：各服务通过 `backend/services/<service>/migrations/` 下的迁移文件独立管理表结构。本目录中的文件为聚合快照，**不作为自动化迁移的入口**。

## 服务迁移对照

| 服务 | 根级快照 | 服务级迁移 |
|------|---------|-----------|
| channel | `channel/001_create_channel.sql` | `services/channel-service/migrations/001_init.sql` |
| iam | `iam/001~007_*.sql` | `services/iam-service/migrations/001_init.sql` |
| inventory | `inventory/001_create_inventory.sql` | `services/inventory-service/migrations/001_init.sql` |
| order | `order/001_create_order.sql` | `services/order-service/migrations/001_init.sql` |
| product | `product/001_create_product.sql` | `services/product-service/migrations/001_init.sql` |
| tenant | `tenant/001~004_*.sql` | `services/tenant-service/migrations/001_init.sql` |

## 后续计划

- [ ] 统一到服务级迁移，由各服务启动时自动执行
- [ ] 根级快照转为只读归档，不再增量维护
- [ ] 引入 golang-migrate 或同等工具管理迁移版本
