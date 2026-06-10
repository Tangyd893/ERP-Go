# 订单履约手工烟测 Runbook

> 用于验证订单审核→出库→发货→库存扣减全链路，无需 RabbitMQ（HTTP 回调模式）。

## 前置条件

- `scripts/dev-stack.ps1 all` 已启动全部服务（Gateway :8080, IAM :8081, Inventory :8086, Warehouse :8087, Order :8085）
- PostgreSQL + 迁移已就绪
- 以下变量已设置（dev-stack 自动设置）：
  - `ORDER_SERVICE_URL=http://localhost:8085`
  - `INVENTORY_SERVICE_URL=http://localhost:8086`
  - `WAREHOUSE_SERVICE_URL=http://localhost:8087`

## 烟测步骤

### Step 1: 登录获取 Token

```powershell
$body = @{username="admin"; password="admin123"; tenant_id="default"} | ConvertTo-Json
$token = (Invoke-RestMethod -Uri http://localhost:8080/api/v1/iam/login -Method POST -Body $body -ContentType "application/json").data.access_token
$headers = @{Authorization = "Bearer $token"; "X-Tenant-ID" = "default"}
```

### Step 2: 创建测试订单

```powershell
$orderBody = @{
  order_no = "SMOKE-001"
  buyer_name = "烟测买家"
  items = @(@{sku_id="sku-001"; sku_name="测试商品"; quantity=2; unit_price=9.99})
  total_amount = 19.98
  currency = "CNY"
} | ConvertTo-Json

$order = Invoke-RestMethod -Uri http://localhost:8080/api/v1/order/orders -Method POST -Body $orderBody -Headers $headers -ContentType "application/json"
$orderId = $order.data.id
Write-Host "订单ID: $orderId"
```

### Step 3: 审核订单（触发锁库 + 建出库单）

```powershell
$auditBody = @{order_id = $orderId; approved = $true} | ConvertTo-Json
Invoke-RestMethod -Uri http://localhost:8080/api/v1/order/orders/audit -Method POST -Body $auditBody -Headers $headers -ContentType "application/json"
# 期望: {"code":0,"data":{"approved":true}}
```

### Step 4: 查询出库单

```powershell
$outbounds = Invoke-RestMethod -Uri http://localhost:8080/api/v1/warehouse/outbounds -Headers $headers
$outboundId = $outbounds.data.list[0].id
Write-Host "出库单ID: $outboundId"
# 期望: 出库单状态为 picking 或 created
```

### Step 5: 模拟出库完成（发货回调）

```powershell
$shipBody = @{
  outbound_id = $outboundId
  order_id = $orderId
  warehouse_id = "wh-001"
  items = @(@{sku_id="sku-001"; sku_code="A001"; sku_name="测试商品"; qty=2})
  tracking_no = "TN-SMOKE-001"
  carrier = "YTO"
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/api/v1/order/fulfillment/outbound-shipped -Method POST -Body $shipBody -Headers $headers -ContentType "application/json"
# 期望: {"code":0,"data":{"processed":true,"order_id":"<orderId>"}}
```

### Step 6: 验证订单状态

```powershell
$orderStatus = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/order/orders/$orderId" -Headers $headers
Write-Host "订单状态: $($orderStatus.data.status)"
# 期望: shipped
```

### Step 7: 验证库存扣减

```powershell
$balances = Invoke-RestMethod -Uri http://localhost:8080/api/v1/inventory/balances -Headers $headers
# 确认 sku-001 的 total_quantity 减少了 2
```

## 失败场景烟测

### 场景 A: 审核→锁库失败

```powershell
# 对不存在的 SKU 创建订单并审核
# 期望: 审核返回错误，不创建出库单
# 验证: GET /outbox/failed 可查询失败消息
Invoke-RestMethod -Uri http://localhost:8080/api/v1/order/outbox/failed -Headers $headers
```

### 场景 B: 死信重试

```powershell
# 获取失败消息 ID
$failed = Invoke-RestMethod -Uri http://localhost:8080/api/v1/order/outbox/failed -Headers $headers
$failedId = $failed.data.list[0].id

# 重试
$retryBody = @{id = $failedId} | ConvertTo-Json
Invoke-RestMethod -Uri http://localhost:8080/api/v1/order/outbox/retry -Method POST -Body $retryBody -Headers $headers -ContentType "application/json"
# 期望: {"code":0,"data":{"retried":true,"id":<id>}}
```

### 场景 C: 重复审核幂等

```powershell
# 对同一订单再次审核
# 期望: 幂等返回成功（Inbox 去重），不重复锁库/建出库单
Invoke-RestMethod -Uri http://localhost:8080/api/v1/order/orders/audit -Method POST -Body $auditBody -Headers $headers -ContentType "application/json"
```

## 清理

```powershell
# 取消测试订单
$cancelBody = @{order_id = $orderId; reason = "烟测清理"} | ConvertTo-Json
Invoke-RestMethod -Uri http://localhost:8080/api/v1/order/orders/cancel -Method POST -Body $cancelBody -Headers $headers -ContentType "application/json"
```

## 预期结果一览

| 步骤 | 检查点 | 预期 |
|------|--------|------|
| 登录 | Token 非空 | ✅ |
| 创建订单 | 返回 order_id | ✅ |
| 审核 | approved:true | ✅ |
| 查出库单 | 状态 picking/created | ✅ |
| 发货回调 | processed:true | ✅ |
| 查订单 | 状态 shipped | ✅ |
| 查库存 | sku-001 扣减 2 | ✅ |
| 查死信 | 可查询 /outbox/failed | ✅ |
| 重试死信 | retried:true | ✅ |
| 重复审核 | 幂等无副作用 | ✅ |
