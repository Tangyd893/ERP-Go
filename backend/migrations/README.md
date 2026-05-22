# 数据库迁移目录

**当前迁移策略**：`backend/services/<service>/migrations/` 是业务表的**唯一事实源**。本目录保留各服务的历史 SQL 快照，**不作为自动化迁移的入口**，仅用于聚合查阅和历史追溯。

> [!WARNING]
> 根级快照与服务级 migration 存在表结构差异（字段名、类型、索引、schema 均不同），**所有业务表建表请以服务级 migration 为准**。

## 唯一活动迁移

| 组件 | 路径 | 说明 |
|------|------|------|
| **Outbox/Inbox** | `outbox/001_create_outbox.sql` | 共享基础设施（发件箱/收件箱），不归属任一业务服务，执行顺序应在所有业务表之前 |

## 业务表迁移对照（根级为历史快照）

| 服务 | 根级快照（只读归档） | 服务级迁移（活动入口） |
|------|---------|-----------|
| channel | `channel/001_create_channel.sql` | `services/channel-service/migrations/001_init.sql` |
| iam | `iam/001~007_*.sql` | `services/iam-service/migrations/001_init.sql` |
| inventory | `inventory/001_create_inventory.sql` | `services/inventory-service/migrations/001_init.sql` |
| order | `order/001_create_order.sql` | `services/order-service/migrations/001_init.sql` |
| product | `product/001_create_product.sql` | `services/product-service/migrations/001_init.sql` |
| tenant | `tenant/001~004_*.sql` | `services/tenant-service/migrations/001_init.sql` |

> 以下服务仅在服务级维护 migration，根级无对应快照：warehouse、transport、purchase、finance、file、notification、report。

## 执行顺序

1. 先执行根级 `outbox/001_create_outbox.sql`（全局共享基础设施）
2. 再执行各服务级 `001_init.sql`（业务表，无顺序依赖）

## 后续计划

- [x] 业务表统一到服务级 migration
- [x] 根级业务快照标注为只读归档
- [ ] 引入 golang-migrate 或同等工具管理迁移版本
- [ ] 各服务启动时自动执行 migration
