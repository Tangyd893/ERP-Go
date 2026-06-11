# E2E 烟测 Runbook

> 手工端到端验证：登录 → 订单导入 → 审核 → 出库 → 发货 → 回传

## 前置条件

- `.\scripts\dev-stack.ps1 all` 一键启动全部服务
- 各服务端口：gateway:8080, iam:8081, channel:8082, order:8085, inventory:8086, warehouse:8087, transport:8088
- 前端：admin-web :5173, PDA :5174, dashboard :5175

## 步骤

### 1. 登录

```powershell
$body = @{username="admin";password="admin123";tenant_id="default"} | ConvertTo-Json
$resp = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/iam/login" -Method POST -Body $body -ContentType "application/json"
$token = $resp.data.access_token
Write-Host "Token: $token"
```

**验证**: 返回 access_token，状态码 200

### 2. 创建订单

```powershell
$headers = @{Authorization="Bearer $token";"X-Tenant-ID"="default"}
$orderBody = @{
  store_id="store-001"
  platform_order_no="E2E-$(Get-Date -Format 'yyyyMMddHHmmss')"
  order_type="normal"
  items=@(@{sku_id="sku-001";sku_code="A001";sku_name="商品A";quantity=3;unit_price=19.90})
  address=@{contact_name="张三";phone="13800138000";country="中国";city="杭州";street_line1="文一西路969号"}
} | ConvertTo-Json -Depth 5
$order = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/order/orders" -Method POST -Body $orderBody -Headers $headers -ContentType "application/json"
$orderId = $order.data.id
Write-Host "OrderID: $orderId"
```

**验证**: 返回订单 ID，状态为 pending

### 3. 审核订单

```powershell
$approveBody = @{operator="admin";approved=$true} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/order/orders/$orderId/approve" -Method POST -Body $approveBody -Headers $headers -ContentType "application/json"
```

**验证**: 返回成功，订单状态变为 approved

### 4. 创建出库单

```powershell
$obBody = @{
  order_id=$orderId
  order_no="E2E-001"
  warehouse_id="wh-001"
  items=@(@{sku_id="sku-001";sku_code="A001";sku_name="商品A";quantity=3})
} | ConvertTo-Json -Depth 3
$ob = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/warehouse/outbounds" -Method POST -Body $obBody -Headers $headers -ContentType "application/json"
$obId = $ob.data.id
Write-Host "OutboundID: $obId"
```

**验证**: 返回出库单 ID，自动创建拣货任务

### 5. 创建波次并开始拣货

```powershell
$waveBody = @{warehouse_id="wh-001";name="E2E波次";outbound_ids=@($obId)} | ConvertTo-Json
$wave = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/warehouse/waves" -Method POST -Body $waveBody -Headers $headers -ContentType "application/json"
$waveId = $wave.data.id
# 开始波次
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/warehouse/waves/$waveId/start" -Method POST -Headers $headers
```

**验证**: 波次状态 picking，出库单状态 picking

### 6. PDA 拣货→复核→打包→称重

```powershell
# 拣货
$pickBody = @{task_id="PT-OI0-0";quantity=3} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/warehouse/pick/scan" -Method POST -Body $pickBody -Headers $headers -ContentType "application/json"

# 复核
$checkBody = @{outbound_id=$obId;sku_id="sku-001";quantity=3} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/warehouse/check/scan" -Method POST -Body $checkBody -Headers $headers -ContentType "application/json"

# 打包
$packBody = @{outbound_id=$obId;weight=1.5} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/warehouse/package" -Method POST -Body $packBody -Headers $headers -ContentType "application/json"

# 称重
$weighBody = @{outbound_id=$obId;weight=1.52} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/warehouse/weigh" -Method POST -Body $weighBody -Headers $headers -ContentType "application/json"
```

**验证**: 每步返回成功，出库单状态依次流转

### 7. 物流匹配 + 生成面单 + 发货

```powershell
# 物流匹配
$matchBody = @{weight=1.52;country="中国"} | ConvertTo-Json
$match = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/transport/match-carrier" -Method POST -Body $matchBody -Headers $headers -ContentType "application/json"

# 创建发运单
$shipBody = @{order_id=$orderId;outbound_id=$obId;carrier_code="YTO";weight=1.52} | ConvertTo-Json
$ship = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/transport/shipments" -Method POST -Body $shipBody -Headers $headers -ContentType "application/json"
$shipId = $ship.data.id

# 生成面单
$labelBody = @{shipment_id=$shipId;order_no="E2E-001";weight=1.52;country="中国"} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/transport/labels/generate" -Method POST -Body $labelBody -Headers $headers -ContentType "application/json"

# 出库确认
$shipConfirm = @{tracking_no="MOCK-001";carrier="YTO"} | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/warehouse/outbounds/$obId/ship" -Method POST -Body $shipConfirm -Headers $headers -ContentType "application/json"
```

**验证**: 面单生成返回 tracking_no，出库确认成功（状态 shipped）

### 8. 验证回传任务

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/channel/tracking-upload" -Method POST -Body (@{store_id="store-001";tracking_no="MOCK-001";carrier_code="YTO";order_no="E2E-001"} | ConvertTo-Json) -Headers $headers -ContentType "application/json"
```

**验证**: 返回 tracking_upload 同步任务，状态 pending

## 一键脚本

```powershell
# 保存为 testing/e2e/smoke.ps1，在项目根目录执行
param(
  [string]$Gateway = "http://localhost:8080"
)

function Test-Step($Name, $Script) {
  Write-Host -ForegroundColor Cyan ">>> $Name"
  try { & $Script } catch { Write-Host -ForegroundColor Red "FAIL: $_"; exit 1 }
}

$headers = @{Authorization="Bearer "; "Content-Type"="application/json"; "X-Tenant-ID"="default"}

# 1. Login
Test-Step "1. 登录" {
  $body = @{username="admin";password="admin123";tenant_id="default"} | ConvertTo-Json
  $resp = Invoke-RestMethod "$Gateway/api/v1/iam/login" -Method POST -Body $body -ContentType "application/json"
  $script:headers.Authorization = "Bearer $($resp.data.access_token)"
  Write-Host "  Token: $($resp.data.access_token.Substring(0,20))..."
}

# 2. Create order
$orderId = ""
Test-Step "2. 创建订单" {
  $body = @{store_id="store-001";platform_order_no="SMOKE-$(Get-Random)";order_type="normal";items=@(@{sku_id="sku-001";sku_code="A001";sku_name="商品A";quantity=2;unit_price=19.90});address=@{contact_name="张三";phone="13800138000";country="中国";city="杭州";street_line1="文一西路969号"}} | ConvertTo-Json -Depth 5
  $resp = Invoke-RestMethod "$Gateway/api/v1/order/orders" -Method POST -Body $body -Headers $headers
  $script:orderId = $resp.data.id
  Write-Host "  OrderID: $orderId"
}

# 3-7 ... (继续上述完整流程)

Write-Host -ForegroundColor Green "✓ E2E 烟测全部通过！"
```
