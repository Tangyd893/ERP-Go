# 订单履约主链路手工验收 Runbook

> 目标：10 分钟内从零走通"订单审核 → 库存锁定 → 出库单 → 发货 → 库存扣减"完整链路。

## 前置条件

| 条件 | 检查命令 |
|---|---|
| Docker 已启动 | `docker ps` |
| 中间件已就绪（PG/RabbitMQ/Redis） | `.\scripts\dev-stack.ps1 infra` |
| 核心服务已启动 | 见下方端口表 |

## 0. 启动全栈

```powershell
# Windows PowerShell — 一键启动
.\scripts\dev-stack.ps1 all
```

```bash
# Linux / macOS
make dev-stack
```

启动后等待约 5 秒，确认服务健康：

```bash
curl http://localhost:8080/health    # Gateway: {"status":"ok"}
curl http://localhost:8085/health    # Order:   {"status":"ok"} 或 degraded
curl http://localhost:8086/health    # Inventory
curl http://localhost:8087/health    # Warehouse
```

## 1. IAM 登录获取 Token

```bash
curl -s -X POST http://localhost:8080/api/v1/iam/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123","tenant_id":"default"}' | jq .
```

保存返回的 `access_token`：

```bash
export TOKEN="<粘贴 access_token>"
```

## 2. 创建销售订单

```bash
curl -s -X POST http://localhost:8080/api/v1/order/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" \
  -d '{
    "store_id": "st-1",
    "warehouse_id": "wh-001",
    "order_no": "SO-SMOKE-001",
    "buyer_name": "测试买家",
    "items": [
      {"sku_id": "sku-001", "sku_code": "A001", "sku_name": "商品A", "quantity": 2}
    ]
  }' | jq .
```

记下返回的 `order_id`。

## 3. 审核订单

```bash
curl -s -X POST http://localhost:8080/api/v1/order/orders/audit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" \
  -d '{"order_id":"<订单ID>","approved":true}' | jq .
```

**预期**：`{"code":0, "data":{"approved":true}}`

### 验证审核效果

```bash
# 1. 查询 Outbox 表（应有 stock.locked + outbound.created 事件）
docker exec erp-postgres psql -U erp -d erp_go \
  -c "SELECT event_type, status FROM outbox_messages ORDER BY id DESC LIMIT 5;"

# 2. 查询出库单
curl -s http://localhost:8080/api/v1/warehouse/outbounds \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" | jq .

# 3. 查询库存（应显示锁定数量）
curl -s http://localhost:8080/api/v1/inventory/balances \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" | jq .
```

## 4. PDA 拣货（模拟）

```bash
# 获取拣货任务列表
curl -s http://localhost:8080/api/v1/warehouse/pick-tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" | jq .

# 拣货扫码
curl -s -X POST http://localhost:8080/api/v1/warehouse/pick/scan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" \
  -d '{"outbound_id":"<出库单ID>","sku_id":"sku-001","quantity":2}' | jq .
```

## 5. 复核扫码

```bash
curl -s -X POST http://localhost:8080/api/v1/warehouse/check/scan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" \
  -d '{"outbound_id":"<出库单ID>","sku_id":"sku-001","quantity":2}' | jq .
```

## 6. 打包

```bash
curl -s -X POST http://localhost:8080/api/v1/warehouse/package \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" \
  -d '{"outbound_id":"<出库单ID>"}' | jq .
```

## 7. 称重

```bash
curl -s -X POST http://localhost:8080/api/v1/warehouse/weigh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" \
  -d '{"outbound_id":"<出库单ID>","weight":1500}' | jq .
```

## 8. 出库确认（关键步骤）

```bash
curl -s -X POST http://localhost:8080/api/v1/warehouse/outbounds/<出库单ID>/ship \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" \
  -d '{"tracking_no":"TN-001","carrier":"SF"}' | jq .
```

**预期**：WMS 更新出库单状态 → 通知 Order 服务 → Order 扣减库存 → 更新订单为 `shipped`。

## 9. 验证最终状态

```bash
# 查询订单状态（应为 shipped）
curl -s http://localhost:8080/api/v1/order/orders/<订单ID> \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" | jq '.status'

# 查询 outbox 事件历史
docker exec erp-postgres psql -U erp -d erp_go \
  -c "SELECT event_type, status, created_at FROM outbox_messages ORDER BY id;"

# 查询失败 outbox（如有）
curl -s "http://localhost:8080/api/v1/order/outbox/failed?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: default" | jq .
```

## 10. 验收通过标准

| 步骤 | 检查项 | 预期 |
|---|---|---|
| 3 | 审核返回 | `approved: true` |
| 3 | outbox_messages 表 | 包含 `order.approved` / `stock.locked` / `outbound.created` |
| 4-7 | PDA 操作 | 各步骤返回成功，出库单状态逐步变化 |
| 8 | 出库确认 | 库存扣减成功，订单状态 `shipped` |
| 9 | 出库后 outbox | 包含 `stock.deducted` / `order.shipped` |
| 9 | 失败消息列表 | 为空（无异常） |

## 常见问题

| 现象 | 可能原因 | 处理 |
|---|---|---|
| 审核返回错误 | Order 服务 degraded（DB 未就绪） | 检查 PG、迁移是否正确执行 |
| 出库单不存在 | Warehouse 服务 degraded | 同上 |
| WMS 操作返回 404 | 路由未注册 | 确认 Gateway 代理正确映射 warehouse 路由 |
| 出库后订单状态未变 | `ORDER_SERVICE_URL` 未设置或不可达 | 检查 warehouse 和 order 间网络 |
| Token 过期 | `JWT_SECRET` 不一致 | 确认所有服务使用同一 secret |
