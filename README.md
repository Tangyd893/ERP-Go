# ERP-Go

> [English](README.en.md)

ERP-Go 是一个基于 Go 微服务架构的跨境电商 ERP 系统，面向跨境卖家、品牌方和海外仓团队，覆盖商品、渠道、订单、库存、WMS 仓储、TMS 物流、采购、财务、报表和系统治理等业务域。

**一句话定位：** 面向跨境电商卖家、品牌方和海外仓团队的一体化微服务 ERP 平台。

**长期方向：** 让跨境团队围绕「订单履约」和「库存资金账」完成商品管理、平台销售、采购补货、仓储作业、物流发运、财务核算和经营分析，而不是在平台后台、仓库系统、物流系统和表格之间反复搬运数据。

仓库地址：[https://github.com/Tangyd893/ERP-Go.git](https://github.com/Tangyd893/ERP-Go.git)

## 核心能力

- **多租户治理** — 用户认证、租户与组织、角色权限、操作审计
- **商品与渠道** — SPU/SKU 管理、平台 SKU 映射、店铺授权、平台订单导入
- **订单履约** — 销售订单状态机、审核、库存锁定、WMS 出库、发货与平台回传
- **仓储作业** — 入库、拣货、复核、打包、出库；管理后台与 PDA 扫码终端
- **供应链与财务** — 采购单、到货入库、应收应付、成本归集与利润报表
- **平台能力** — API Gateway 统一鉴权、Outbox/Inbox 事件驱动、可观测性基础设施

## 架构概览

```text
frontend
  ├─ ERP 管理后台 (admin-web Vue · admin-web-react 并行)
  ├─ WMS PDA / H5 (warehouse-pda Vue · warehouse-pda-react 并行)
  └─ 经营看板 (dashboard-web Vue · dashboard-web-react 并行)
        |
        v
API Gateway
        |
        +-- iam-service             认证权限
        +-- tenant-service          租户组织
        +-- product-service         商品
        +-- channel-service         渠道
        +-- order-service           订单
        +-- inventory-service       库存
        +-- warehouse-service       WMS 仓储
        +-- transport-service       TMS 物流
        +-- purchase-service        采购
        +-- finance-service         财务
        +-- report-service          报表
        +-- file-service            文件
        +-- notification-service    通知

PostgreSQL | Redis | RabbitMQ | MinIO | OpenSearch
Prometheus | Grafana | Loki | Jaeger
```

**设计原则：**

- 业务闭环优先于服务数量；每个服务拥有独立业务边界和数据所有权
- 跨服务副作用优先事件化；库存与财务必须保留流水和对账能力
- 平台、物流商、海外仓等外部系统通过适配器隔离

事件驱动与跨服务协作详见 [接口与事件设计](docs/接口与事件设计.md)。

## 微服务职责

| 服务 | 职责 |
|------|------|
| API Gateway | 统一入口、JWT 鉴权、路由转发、追踪上下文 |
| IAM Service | 用户、角色、权限、登录、审计 |
| Tenant Service | 租户、组织、部门、岗位 |
| Product Service | SPU、SKU、平台 SKU 映射 |
| Channel Service | 店铺授权、平台同步、订单导入 |
| Order Service | 销售订单、状态机、审核、异常处理 |
| Inventory Service | 库存余额、锁定、释放、扣减、流水 |
| Warehouse Service | 入库、上架、拣货、复核、打包、出库 |
| Transport Service | 物流商、物流规则、面单、发运、轨迹 |
| Purchase Service | 供应商、采购单、到货、质检、入库 |
| Finance Service | 应收应付、结算、成本、利润、流水 |
| Report Service | 销售、库存、仓储、利润报表 |
| File Service | 商品图片、面单、发票、导入导出文件 |
| Notification Service | 邮件、短信、站内信、Webhook、告警 |

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.25 + Gin + GORM |
| 网关 | Gin API Gateway、统一鉴权、反向代理 |
| 数据库 | PostgreSQL |
| 缓存 | Redis |
| 异步事件 | RabbitMQ，Outbox/Inbox 可靠消息 |
| 对象存储 | MinIO（生产可对接 S3 兼容存储） |
| 搜索分析 | OpenSearch |
| 前端 | Vue 3（现网）+ React 18 并行（admin-web-react）+ TypeScript + Vite |
| React 管理端 | Tailwind + shadcn/Radix（@erp/shared-react）+ Zustand + TanStack Query |
| 仓库终端 | Vue 3 H5 / PDA 扫码作业 |
| 监控 | Prometheus、Grafana、Loki、Jaeger |
| 部署 | Docker Compose（本地）；Kubernetes / Helm（生产） |

## 仓库结构

```text
ERP-Go/
  backend/    后端微服务、网关、共享组件、数据库迁移
  frontend/   Vue 三端 + React 并行管理端（admin-web-react）与 shared-react
  docker/     Docker Compose、Kubernetes、监控和部署配置
  testing/    集成测试、契约测试、端到端测试、性能测试
  docs/       架构设计、领域模型、接口事件、数据模型和工程规范
  scripts/    本地启动、迁移、验证脚本
```

## 快速开始

### 环境要求

| 工具 | 建议版本 | 用途 |
|------|----------|------|
| Git | 2.40+ | 版本管理 |
| Docker | 20.10+ | 本地基础设施 |
| Go | 1.25.0 | 后端开发 |
| Node.js | 18.x+ | 前端开发 |
| npm | 9.x+ | 前端包管理 |

### 克隆仓库

```bash
git clone https://github.com/Tangyd893/ERP-Go.git
cd ERP-Go
```

### 本地开发栈

**Windows PowerShell：**

```powershell
.\scripts\dev-stack.ps1 all       # 中间件 + 迁移 + 核心服务
.\scripts\dev-stack.ps1 infra       # 仅 Docker 中间件 + 迁移
.\scripts\dev-stack.ps1 services    # 仅启动核心微服务
```

**Linux / macOS / Git Bash：**

```bash
make dev-stack    # 中间件 + 迁移 + 核心服务
make dev-infra    # 仅中间件 + 迁移
make dev-services # 仅核心微服务
```

`dev-stack` 默认启动 Gateway、IAM、Order、Inventory、Warehouse。服务日志写入 `.cache/logs/`。

**健康检查：**

```bash
curl http://localhost:8080/health
```

### 默认端口

| 服务 | 端口 |
|------|-----:|
| API Gateway | 8080 |
| IAM Service | 8081 |
| Tenant Service | 8082 |
| Product Service | 8083 |
| Channel Service | 8084 |
| Order Service | 8085 |
| Inventory Service | 8086 |
| Warehouse Service | 8087 |
| Transport Service | 8088 |
| File Service | 8089 |
| Purchase Service | 8091 |
| Finance Service | 8092 |
| Report Service | 8093 |
| Notification Service | 8094 |

### 前端开发

```bash
npm install
npm run dev:admin          # 管理后台 Vue  → http://localhost:5173
npm run dev:admin-react    # 管理后台 React → http://localhost:5273
npm run dev:pda            # WMS PDA Vue   → http://localhost:5174
npm run dev:pda-react      # WMS PDA React → http://localhost:5274
npm run dev:dashboard      # 经营看板 Vue  → http://localhost:5175
npm run dev:dashboard-react # 经营看板 React → http://localhost:5275
```

React 前端生产镜像（`APP` 可选 `admin-web-react` / `dashboard-web-react` / `warehouse-pda-react`）：

```bash
docker build -f docker/frontend/Dockerfile --build-arg APP=admin-web-react -t erp-admin-web-react .
docker run --rm -p 8088:80 erp-admin-web-react
```

## 开发与验证

推荐使用统一验证脚本，自动将 `GOCACHE` / `GOMODCACHE` 固定到仓库 `.cache/` 目录：

```bash
# Linux / macOS / Git Bash / CI
./scripts/verify.sh

# Windows PowerShell
.\scripts\verify.ps1

# 跳过前端（仅验证 Go）
./scripts/verify.sh --skip-frontend
.\scripts\verify.ps1 -SkipFrontend
```

**Makefile 快捷入口：**

```bash
make verify    # Go vet/test/build + 前端 typecheck/lint/build
make test      # Go 单元测试
make lint      # go vet + eslint
make build     # 编译全部 Go 服务
make all       # 编译 + 测试
make help      # 查看所有命令
```

> 直接运行 `go test ./...` 可能因系统默认 Go 缓存目录权限失败；请优先使用上述脚本或 `make test`。

### 手动启动（可选）

```bash
# 启动开发中间件
docker compose -f docker/compose/docker-compose.dev.yml up -d

# 启动 IAM 服务（需先完成数据库迁移）
DATABASE_PORT=5433 SERVER_PORT=8081 JWT_SECRET=dev-secret go run ./backend/services/iam-service/cmd/server/
```

环境变量模板见 [.env.example](.env.example)。

## 文档

完整文档索引见 [docs/文档索引.md](docs/文档索引.md)。

**架构与设计：**

- [项目架构设计](docs/项目架构设计.md)
- [领域模型设计](docs/领域模型设计.md)
- [微服务设计说明](docs/微服务设计说明.md)
- [接口与事件设计](docs/接口与事件设计.md)
- [数据模型设计](docs/数据模型设计.md)

**工程规范：**

- [技术栈与中间件管理规范](docs/技术栈与中间件管理规范.md)
- [实施路线与工程规范](docs/实施路线与工程规范.md)

**运维：**

- [docs/runbooks/](docs/runbooks/) — 部署、E2E 烟测、备份恢复等操作手册

**API 契约：**

- [docs/api/openapi.yaml](docs/api/openapi.yaml)

## 生产部署注意事项

- 生产环境必须通过环境变量、Kubernetes Secret 或专用密钥管理系统注入敏感配置
- 设置 `ENVIRONMENT=production` 时，Gateway/IAM 将拒绝默认 `JWT_SECRET` 与默认数据库密码（见 `.env.example`）；非 `development` 环境数据库未就绪时服务将退出
- 生产 CORS 须配置 `CORS_ALLOWED_ORIGINS` 白名单，勿使用 `*`
- Gateway 可选启用 `GATEWAY_RBAC_ENABLED=true` 进行粗粒度 RBAC；后端服务应通过 NetworkPolicy 禁止绕过 Gateway 直连
- 不要提交 `.env`、数据库密码、平台 Token、物流商密钥、支付密钥或真实客户数据
- 外部系统 Token 必须加密存储并记录授权审计
- 多租户数据访问必须强制校验 `tenant_id`
- 库存调整、财务冲销、权限变更、数据导出必须记录审计日志
- 所有服务必须接入结构化日志、指标和链路追踪
