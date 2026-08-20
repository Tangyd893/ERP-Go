# React 前端二期 / 三期执行大纲

> 一期（`admin-web-react` + `shared-react`）验收通过后再开独立实现会话。  
> Vue 包继续保留；新建并行 `*-react` 包，端口与脚本增量，不切换默认 `dev:*`。

## 二期：`dashboard-web-react`（:5275）

**目标：** 经营看板 React 版，复用 `@erp/shared-react`，观感偏大屏 KPI。

### 范围

- 包路径：`frontend/dashboard-web-react`
- 路由：单页 `/`（可与 Vue `dashboard-web` 功能对等）
- 数据：`GET /report/dashboard`，约 30s 轮询；支持 `VITE_DEMO_MODE`
- 图表：ECharts 或 Recharts（二选一，与 Vue 数据字段对齐）
- Auth：可关闭 401 跳转（与现 Vue dashboard `VITE_AUTH_REDIRECT=false` 一致）

### 交付

- `npm run dev:dashboard-react` / `build:dashboard-react`
- Docker：`APP=dashboard-web-react`
- Playwright：`dashboard-react-smoke.spec.ts`
- 文档：更新 `docs/runbooks/frontend-verification.md`

### 验收

- KPI + 趋势图可渲染；demo 模式无永久 loading
- 硬刷新不白屏；构建与镜像通过

---

## 三期：`warehouse-pda-react`（:5274）

**目标：** WMS PDA / H5 React 版；触控与扫码交互优先于桌面观感。

### 范围

- 包路径：`frontend/warehouse-pda-react`
- 路由对齐 Vue PDA：登录、首页、拣货列表/扫码、复核、出库确认、打包、称重、我的
- 布局：底栏 Tab + 大触控热区；扫码输入自动聚焦
- 复用 `@erp/shared-react` 原语；可增加 PDA 专用密度 token（更大按钮）

### 交付

- `npm run dev:pda-react` / `build:pda-react`
- Docker：`APP=warehouse-pda-react`
- Playwright：移动端 viewport 烟测
- 文档与 seed 出库单验收步骤

### 验收

- 底栏可导航；拣货/复核主流程可点通（同后端）
- 硬刷新不白屏；构建与镜像通过

---

## 跨期注意

- 禁止修改 Vue 源码作为迁移手段
- 设计继续遵循 `.cursor/skills/baseline-ui` 与 `design-taste-frontend`（ERP dials：高密度、低动效）
- 默认流量仍指向 Vue，直到产品决定切换 Ingress / 脚本别名
