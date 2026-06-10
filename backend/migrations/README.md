# 数据库迁移目录

**当前迁移策略**：`backend/services/<service>/migrations/` 是业务表的**唯一事实源**。本目录保留各服务的历史 SQL 快照，**不作为自动化迁移的入口**，仅用于聚合查阅和历史追溯。

> [!WARNING]
> 根级快照与服务级 migration 存在表结构差异（字段名、类型、索引、schema 均不同），**所有业务表建表请以服务级 migration 为准**。

## 唯一活动迁移

| 组件 | 路径 | 说明 |
|------|------|------|
| **Outbox/Inbox** | `outbox/001_create_outbox.sql` | 共享基础设施（发件箱/收件箱），不归属任一业务服务，执行顺序应在所有业务表之前 |
| **Outbox TenantID** | `outbox/002_add_tenant_id.sql` | Outbox 消息增加租户隔离字段 |

## 业务表迁移对照

| 服务 | 根级快照（只读归档） | 服务级迁移（活动入口） | 备注 |
|------|---------|-----------|------|
| iam | `iam/001~006_*.sql` | `iam-service/migrations/001_init.sql`, `002_seed_dev_admin.sql` | `002` 为开发种子数据 |
| tenant | `tenant/001~004_*.sql` | `tenant-service/migrations/001_init.sql` | — |
| product | `product/001_create_product.sql` | `product-service/migrations/001_init.sql` | — |
| channel | `channel/001_create_channel.sql` | `channel-service/migrations/001_init.sql` | — |
| order | `order/001_create_order.sql` | `order-service/migrations/001_init.sql` | — |
| inventory | `inventory/001_create_inventory.sql` | `inventory-service/migrations/001_init.sql` | — |
| warehouse | — | `warehouse-service/migrations/001_init.sql` | 仅服务级 |
| transport | — | `transport-service/migrations/001_init.sql` | 仅服务级 |
| purchase | — | `purchase-service/migrations/001_init.sql` | 仅服务级 |
| finance | — | `finance-service/migrations/001_init.sql` | 仅服务级 |
| file | — | `file-service/migrations/001_init.sql` | 仅服务级 |
| notification | — | `notification-service/migrations/001_init.sql` | 仅服务级 |
| **report** | — | — | **无业务表**（无状态聚合服务，从其他服务读取数据） |

## 执行顺序

迁移由 `scripts/migrate.sh` 或 `scripts/migrate.ps1` 按以下顺序执行：

1. **Outbox 基础设施**：`outbox/001_create_outbox.sql` → `outbox/002_add_tenant_id.sql`
2. **业务服务**：按字母序（iam → tenant → … → warehouse），各服务内 `*.sql` 按文件名字母序
3. **report-service 跳过**：无业务表，不参与迁移

服务间无顺序依赖（各服务独立 schema），唯一要求是 Outbox 在所有业务表之前创建。

## 后续计划

- [x] 业务表统一到服务级 migration
- [x] 根级业务快照标注为只读归档
- [ ] 引入 golang-migrate 或同等工具管理迁移版本
- [ ] 各服务启动时自动执行 migration
