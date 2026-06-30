# 前端验收 Runbook

> 目标：三端「打开能用」— Dashboard 无永久遮罩、PDA 有底栏可导航、硬刷新不白屏。
> 关联任务：T-629（本文）、T-630~T-631（实现）。

## 前置条件

- Windows 10+ / Ubuntu 22.04+
- Go 1.25+、Node.js 20+、npm 10+
- Docker Desktop（PostgreSQL 16 + RabbitMQ 3.13）
- 仓库根目录已执行 `npm ci`

## 1. 启动后端栈

```powershell
# 一键启动中间件 + 迁移 + 5 核心服务
.\scripts\dev-stack.ps1 all
```

服务就绪检查：

```powershell
# 逐个验证健康
curl -s http://localhost:8080/health          # Gateway
curl -s http://localhost:8081/health          # IAM
curl -s http://localhost:8085/health          # Order
curl -s http://localhost:8086/health          # Inventory
curl -s http://localhost:8087/health          # Warehouse
```

日志排查：

```powershell
Get-Content .cache\logs\gateway.log -Tail 20
Get-Content .cache\logs\order.log -Tail 20
```

## 2. 启动三端前端

```powershell
# 三个终端分别启动
cd frontend\admin-web; npm run dev       # → http://localhost:5173
cd frontend\warehouse-pda; npm run dev   # → http://localhost:5174
cd frontend\dashboard-web; npm run dev   # → http://localhost:5175
```

## 3. 验收清单

> **关键**：打开浏览器 DevTools → Network 标签 → 勾选 **Disable cache** → 每次切换前 **Ctrl+Shift+R** 硬刷新。

### 3.1 Admin（`http://localhost:5173/login`）

| # | 检查项 | 预期 | 实际 |
|---|--------|------|------|
| A1 | 登录页加载 | `/login` 显示居中登录表单 | ☐ |
| A2 | 登录 | `admin` / `admin123` / `default` → 跳转 `/dashboard` | ☐ |
| A3 | 侧边栏 | 左侧深色菜单栏，可展开子菜单 | ☐ |
| A4 | Dashboard | 4 张 KPI 卡片可渲染 | ☐ |
| A5 | 列表页 | 点击菜单「订单管理→订单列表」有表格（空数据允许） | ☐ |
| A6 | 退出登录 | 右上角下拉 → 退出 → 回到 `/login` | ☐ |

### 3.2 Dashboard（`http://localhost:5175/`）

| # | 检查项 | 预期 | 实际 |
|---|--------|------|------|
| D1 | 硬刷新 | 页面无无限转圈 spinner | ☐ |
| D2 | KPI 卡片 | 4 张卡片可见（demo 数据或用引导提示） | ☐ |
| D3 | 图表 | ECharts 图表渲染（demo 数据） | ☐ |
| D4 | 无 token 时 | 显示「请先登录」引导，不白屏 | ☐ |

### 3.3 PDA（`http://localhost:5174/login`）

| # | 检查项 | 预期 | 实际 |
|---|--------|------|------|
| P1 | 登录页 | 居中登录表单，适合移动端 | ☐ |
| P2 | 登录 | 任意凭据 → 跳转首页 | ☐ |
| P3 | 底部导航 | 底部 Tab 栏可见（首页/拣货/复核/打包/我的） | ☐ |
| P4 | 顶栏 | 顶栏显示租户/用户名 | ☐ |
| P5 | 首页 | 出库任务计数卡片可见 | ☐ |
| P6 | 退出 | 「我的」页或顶栏可退出 → 回到 `/login` | ☐ |
| P7 | 硬刷新 | 不白屏；底部导航仍存在 | ☐ |

## 4. 常见问题

| 现象 | 根因 | 处理 |
|------|------|------|
| Dashboard 一直转圈 | `v-loading` 未收到 `loading=false` | 检查 Network 是否返回 401/500；dev 启用 `VITE_DEMO_MODE=true` |
| PDA 无底部导航 | 旧版路由未嵌套 Layout | 确认 `PdaLayout.vue` 存在且路由为 `children` |
| 硬刷新白屏 | `@erp/shared` 无效 default export | 已修复（`4ec69d3`）；若复发检查 `shared/src/index.ts` |
| 后端 503/连接拒绝 | 服务未完全启动 | 等 5s 重试；`Get-Content .cache\logs\*.log` |
| 登录 401 | 迁移未执行 | `.\scripts\dev-stack.ps1 infra` 重跑迁移 |
| PDA 首页计数为 0 | 仓库无种子出库单 | 重跑迁移（含 `002_seed_dev_data.sql` 5 条测试出库单） |

### 4.1 PDA 种子数据

仓库种子数据通过 `002_seed_dev_data.sql` 自动加载，包含：

| 出库单 | 状态 | 说明 |
|--------|------|------|
| SO20260630001 | created | 待拣货（2 个 SKU：蓝色×10 + 红色×5） |
| SO20260630002 | created | 待拣货（1 个 SKU：大号×20） |
| SO20260630003 | created | 待拣货（2 个 SKU：蓝色×15 + 小号×8） |
| SO20260630004 | picking | 拣货中（红色×12） |
| SO20260630005 | picked | 待复核（大号×6） |

PDA 首页预期计数：拣货 3、复核 1、打包 1、称重 0。

若数据缺失，手动执行：

```powershell
Get-Content backend/services/warehouse-service/migrations/002_seed_dev_data.sql -Raw |
  docker exec -i erp-postgres psql -U erp -d erp_go
```

## 5. 通过标准

- [ ] 三端 `#app` 有可见交互元素（非空白）
- [ ] Dashboard 无无限 spinner（KPI 卡片或引导提示可见）
- [ ] PDA 底部 Tab 可导航（点击切换页面）
- [ ] Admin 侧边栏可展开、可退出
- [ ] 硬刷新（Disable cache + Ctrl+Shift+R）三端均不白屏
