# ERP-Go

中文说明

ERP-Go 是一个基于 Go 微服务架构的跨境电商 ERP 系统项目，面向跨境卖家、品牌方和海外仓团队，覆盖商品、渠道、订单、库存、WMS 仓储、TMS 物流、采购、财务、报表和系统治理等业务域。

一句话定位：

> 面向跨境电商卖家、品牌方和海外仓团队的一体化微服务 ERP 平台。

长期方向：

> 让跨境团队围绕“订单履约”和“库存资金账”完成商品管理、平台销售、采购补货、仓储作业、物流发运、财务核算和经营分析，而不是在平台后台、仓库系统、物流系统和表格之间反复搬运数据。

仓库地址：[https://github.com/Tangyd893/ERP-Go.git](https://github.com/Tangyd893/ERP-Go.git)

## 当前状态

本节基于 2026-05-25 对仓库的检查更新。详细缺口与优化优先级见 [docs/优化方案.md](docs/优化方案.md)。

### 状态快照

- Git 状态：位于 `main`，跟踪 `origin/main`。
- 后端：通过 `make test` / 指定 `.cache` 后 `go test -C backend ./...` 通过（**17** 个测试包，**180** 个测试函数），`go vet ./...`、`go build ./...` 通过。GOCACHE/GOMODCACHE 已固定到仓库 `.cache/` 目录。
- 测试覆盖：全部 13 个服务的 domain 包均有单元测试，外加 P4 workflow 测试（13 个）、Gateway 契约测试（4 个，位于 `backend/gateway/cmd/server/`）、middleware 测试（4 个）、Order app 测试（4 个）。`testing/contract` 等目录仍为 README 占位。
- 事件：Order 服务在 DB 就绪时已启动 **进程内 Outbox 轮询** 与 P4 出库编排；RabbitMQ **Consumer 独立 worker 尚未接入**。详见优化方案「事件处理现状」表。
- 前端：`admin-web`、`warehouse-pda`、`dashboard-web` 的 typecheck 与 build 已通过（Node.js v24.15.0 / npm 11.12.1，Linux 环境复验）。`warehouse-pda` 作业页仍为静态占位。
- 数据库迁移：根级迁移已标注为只读归档，服务级 `backend/services/*/migrations` 为业务表唯一事实源（**12/13**，`report-service` 暂无 migration），outbox 在根级单独维护。
- CI：GitHub Actions（`.github/workflows/ci.yml`）在 push/PR 时执行 `scripts/verify.sh`。

### 已完成

- P0 工程底座：`backend`、`frontend`、`docker`、`testing`、`docs` 顶层目录，Go 单模块 monorepo，npm workspaces。
- `backend/shared` 公共组件：配置、日志、错误码、响应、数据库连接、中间件、分页、校验、Outbox/Inbox、可观测性字段。
- API Gateway：健康检查、JWT 鉴权、请求 ID/追踪 ID/租户/用户上下文中间件，覆盖全部 13 个服务的代理路由。
- 13 个后端服务均有 domain/app/HTTP；12/13 有 repository（`report-service` 为无仓储聚合）。DB 就绪时走真实仓储，未就绪时降级占位响应。
- 全部 13 个服务的 domain 包均有单元测试（180 个测试函数）。
- IAM Service：登录、刷新 Token、用户、角色、权限、审计已连线。
- Order Service：状态机、事件驱动（Outbox/Inbox）、P4 履约流程（订单审核→库存锁定→创建出库单→出库扣减→订单发货）、补偿记录。
- Inventory Service：仓储已从内存模拟切换到数据库，支持事务行锁和幂等控制。
- RabbitMQ 事件发布：Order Service 优先使用 RabbitMQ 发布器，不可用时自动降级为日志发布。
- 事件驱动：`shared/events` 20+ 业务事件类型，`shared/outbox` 完整实现 PG 存储/RabbitMQ 发布消费/重试/死信/Inbox 幂等，`shared/workflows` P4 流程协调器+采购入库流程。
- Frontend：`admin-web` 18 个视图、`warehouse-pda` 2 个视图、`dashboard-web` 基础看板页；多数业务页已接 Store，**仅通知中心在 API 失败时回退 mock 数据**。
- Frontend 共享组件：ProTable、ProForm、FileUpload 和 API client（自动注入 JWT/租户）。
- Docker Compose 开发环境：PostgreSQL、Redis、RabbitMQ、MinIO、OpenSearch、Prometheus、Grafana、Loki、Jaeger。
- 架构设计文档、ADR 模板、测试目录骨架已建立。
- Makefile：GOCACHE/GOMODCACHE 已固定到仓库内，无 Unix 专用命令依赖。

### 尚未完成

- IAM/Tenant 尚未编写 app 层和 HTTP 层集成测试。
- 多数服务的 HTTP 接口仍为占位响应，真实业务逻辑待后续版本补全。
- WMS PDA 移动端页面仍是占位 UI，待接入真实拣货/复核/打包流程。
- 采购入库、财务结算的跨服务集成和前端页面尚未打透。
- gRPC 代码生成和服务间调用方案尚未落地。
- 端到端烟测、性能测试仍需补充。
- Kubernetes 生产部署仍是早期骨架（仅网关基础 yaml）。
- Windows 环境下 Vite/Rollup 构建可能需要安装 `@rollup/rollup-win32-x64-msvc` 可选依赖。

## 产品路线

当前推荐路线来自 [实施路线与工程规范](docs/实施路线与工程规范.md)：

```text
先打透订单履约闭环
再补采购入库和财务结算闭环
再做智能补货、智能分仓、BI 和多区域部署
```

近期最重要的闭环：

1. 店铺授权 -> 平台订单导入 -> 订单审核。
2. 订单审核 -> 库存锁定 -> WMS 出库任务。
3. 仓库拣货 -> 复核打包 -> TMS 面单 -> 仓库出库。
4. 仓库出库 -> 库存扣减 -> 订单发货 -> 平台回传。
5. 采购单 -> 到货通知 -> WMS 入库 -> 库存增加。
6. 平台结算 -> 成本归集 -> 订单利润 -> 经营报表。

## 架构概览

```text
frontend
  ├─ ERP 管理后台
  ├─ WMS PDA / H5
  └─ 经营看板
        |
        v
API Gateway
        |
        +-- iam-service             认证权限服务
        +-- tenant-service          租户组织服务
        +-- product-service         商品服务
        +-- channel-service         渠道服务
        +-- order-service           订单服务
        +-- inventory-service       库存服务
        +-- warehouse-service       WMS 仓储服务
        +-- transport-service       TMS 物流服务
        +-- purchase-service        采购服务
        +-- finance-service         财务服务
        +-- report-service          报表服务
        +-- file-service            文件服务
        +-- notification-service    通知服务

PostgreSQL | Redis | RabbitMQ | MinIO | OpenSearch
Prometheus | Grafana | Loki | Jaeger
```

核心设计原则：

- 业务闭环优先于服务数量。
- 每个服务拥有独立业务边界和数据所有权。
- 控制器只做协议转换，业务规则放在领域对象和应用服务中。
- 同步调用用于查询和低频命令，跨服务副作用优先事件化。
- 库存和财务不能只更新余额，必须保留流水和对账能力。
- 平台、物流商、海外仓、支付、税务能力必须通过适配器隔离。
- 新服务必须有契约、测试、指标、日志、追踪和可运维入口。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.23 + Gin + GORM |
| 网关 | Gin API Gateway、统一鉴权入口、反向代理 |
| 数据库 | PostgreSQL |
| 缓存 | Redis |
| 异步事件 | RabbitMQ，Outbox/Inbox 作为可靠消息基础 |
| 对象存储 | MinIO，生产可对接 S3 兼容存储 |
| 搜索分析 | OpenSearch，ClickHouse 作为后续分析库预留 |
| 前端 | Vue 3 + TypeScript + Vite + Pinia + Vue Router + Element Plus |
| 仓库终端 | Vue 3 H5 / PDA 扫码作业 |
| 监控 | Prometheus、Grafana、Loki、Jaeger；OpenTelemetry/Tempo 待进一步落地 |
| 中间件管理 | 本地统一由 `docker/` 目录管理 |
| 生产部署 | Kubernetes，当前仍是骨架阶段 |

## 服务状态

| 服务 | 设计职责 | 当前代码状态 |
|------|----------|--------------|
| API Gateway | 统一入口、鉴权、路由、追踪上下文 | 可编译运行；JWT 鉴权 + 13 个服务代理路由完整覆盖 |
| IAM Service | 用户、角色、权限、登录、审计 | App + Repository + HTTP 已连线；数据库不可用时降级占位模式 |
| Tenant Service | 租户、组织、部门、岗位 | Migration + Repository 完成；HTTP 占位接口 |
| Product Service | SPU、SKU、平台 SKU 映射 | Migration + Domain 测试 + Repository 完成；HTTP 占位接口 |
| Channel Service | 店铺授权、平台同步、订单导入任务 | Migration + Domain 测试 + Repository 完成；HTTP 占位接口 |
| Order Service | 销售订单、状态机、审核、异常 | Migration + Domain 测试 + App 测试 + Repository 完成；事件驱动+Outbox/Inbox 已集成 |
| Inventory Service | 库存余额、锁定、释放、扣减、流水 | Migration + Domain 测试 + Repository 完成；数据库仓储已接入 HTTP |
| Warehouse Service | 入库、上架、拣货、复核、打包、出库 | Domain + Migration 完成；HTTP 占位接口 |
| Transport Service | 物流商、物流规则、面单、发运、轨迹 | Domain 测试 + Migration 完成；HTTP 占位接口 |
| Purchase Service | 供应商、采购单、到货、质检、入库 | Domain 测试 + Migration 完成；HTTP 占位接口 |
| Finance Service | 应收应付、结算、成本、利润、流水 | Domain 测试 + Migration 完成；HTTP 占位接口 |
| Report Service | 销售、库存、仓储、利润报表 | Domain 测试 + 服务骨架 + 占位 API |
| File Service | 商品图片、面单、发票、导入导出文件 | Domain 测试 + Migration 完成；HTTP 占位接口 |
| Notification Service | 邮件、短信、站内信、Webhook、告警 | Domain 测试 + Migration 完成；HTTP 占位接口 |

## 项目目录

```text
ERP-Go/
  backend/    后端微服务、网关、共享组件、数据库迁移
  frontend/   ERP 管理后台、WMS PDA、经营看板、共享组件
  docker/     Docker Compose、Kubernetes、监控和部署配置
  testing/    集成测试、契约测试、端到端测试、性能测试目录
  docs/       架构设计、领域模型、接口事件、数据模型和工程规范文档
```

## 快速开始

### 环境要求

| 工具 | 建议版本 | 用途 |
|------|----------|------|
| Git | 2.40+ | 版本管理 |
| Docker | 20.10+ | 本地基础设施和服务编排 |
| Go | 1.23+ | 后端服务开发 |
| Node.js | 18.x+ | 前端开发 |
| npm | 9.x+ | 前端包管理 |

### 克隆仓库

```bash
git clone https://github.com/Tangyd893/ERP-Go.git
cd ERP-Go
```

### 阅读设计文档

建议从 [文档索引](docs/文档索引.md) 开始：

1. [项目架构设计](docs/项目架构设计.md)
2. [技术栈与中间件管理规范](docs/技术栈与中间件管理规范.md)
3. [领域模型设计](docs/领域模型设计.md)
4. [微服务设计说明](docs/微服务设计说明.md)
5. [接口与事件设计](docs/接口与事件设计.md)
6. [数据模型设计](docs/数据模型设计.md)
7. [项目里程碑与全流程待办清单](docs/项目里程碑与全流程待办清单.md)
8. [实施路线与工程规范](docs/实施路线与工程规范.md)

### 本地开发栈（一键启动）

```powershell
# Windows PowerShell
.\scripts\dev-stack.ps1 all           # 中间件 + 迁移 + 核心服务
.\scripts\dev-stack.ps1 infra          # 仅启动 Docker 中间件 + 迁移
.\scripts\dev-stack.ps1 services       # 仅启动核心微服务
```

```bash
# Linux / macOS / Git Bash
make dev-stack                         # 全栈
make dev-infra                         # 仅中间件 + 迁移
```

启动后，核心服务端口：

| 服务 | 端口 |
| --- | ---: |
| API Gateway | 8080 |
| IAM Service | 8081 |
| Order Service | 8085 |
| Inventory Service | 8086 |
| Warehouse Service | 8087 |

日志位于 `.cache/logs/` 目录；健康检查：`curl http://localhost:8080/health`

## 验证命令

**推荐**：使用统一验证入口，自动固定 `GOCACHE` / `GOMODCACHE` 到仓库 `.cache/` 目录。

```bash
# Linux / macOS / Git Bash / CI
./scripts/verify.sh

# Windows PowerShell
.\scripts\verify.ps1

# npm 不可用时仅验证 Go
.\scripts\verify.ps1 -SkipFrontend
./scripts/verify.sh --skip-frontend

# Makefile 等价入口
make verify    # 全量（Go + 前端）
make test      # 仅 Go 测试
make lint      # 仅 go vet
make all       # 编译 + 测试
```

> 直接运行 `go test ./...` 可能因系统默认 Go 缓存目录权限失败；请优先使用上述脚本或 `make test`。

### 后端（手动）

```bash
mkdir -p .cache/go-build .cache/go-mod
export GOCACHE=$PWD/.cache/go-build
export GOMODCACHE=$PWD/.cache/go-mod

go vet -C backend ./...
go test -C backend ./...
go build -C backend ./...

# 启动开发中间件
docker compose -f docker/compose/docker-compose.dev.yml up -d

# 启动 IAM 服务（需先完成数据库迁移和种子数据）
DATABASE_PORT=5433 SERVER_PORT=8081 JWT_SECRET=dev-secret go run ./backend/services/iam-service/cmd/server/
```

### 前端（手动）

```bash
npm install
npm run dev:admin      # 管理后台，端口 5173
npm run dev:pda        # PDA，端口 5174
npm run dev:dashboard  # 看板，端口 5175
npm run typecheck
npm run build:admin
npm run build:pda
npm run build:dashboard
```

当前前端验证状态：

- 历史上三个项目 typecheck 和 build 均通过；部分 Windows shell 可能无法执行 npm，可用 `-SkipFrontend` 先验 Go。
- Windows 环境下可能需要安装 `@rollup/rollup-win32-x64-msvc` 可选依赖。

### 快捷命令

```bash
make help   # 查看所有可用命令
make dev    # 启动 Docker 中间件
```

## 当前开发重点

1. 为 IAM/Tenant 补齐 app 层和 HTTP 层集成测试。
2. 为 Product/Channel/Warehouse/Transport 等服务的 HTTP 层补全真实业务接口。
3. 将 WMS PDA 从占位 UI 升级为可操作的拣货/复核/打包页面。
4. 打透采购入库完整流程（采购单→到货→入库→库存增加→财务结算）。
5. 建立端到端烟测覆盖首期订单履约主链路。
6. 补齐 Kubernetes/Helm 生产部署配置和运维手册。

## 相关文档

- [文档索引](docs/文档索引.md)
- [项目架构设计](docs/项目架构设计.md)
- [技术栈与中间件管理规范](docs/技术栈与中间件管理规范.md)
- [项目里程碑与全流程待办清单](docs/项目里程碑与全流程待办清单.md)
- [领域模型设计](docs/领域模型设计.md)
- [微服务设计说明](docs/微服务设计说明.md)
- [接口与事件设计](docs/接口与事件设计.md)
- [数据模型设计](docs/数据模型设计.md)
- [实施路线与工程规范](docs/实施路线与工程规范.md)

## 生产部署注意事项

当前项目尚未进入生产部署阶段。后续实现生产部署前必须注意：

- 生产环境必须通过环境变量、Kubernetes Secret 或专用密钥管理系统注入敏感配置。
- 不要提交 `.env`、数据库密码、平台 Token、物流商密钥、支付密钥或真实客户数据。
- Amazon、物流商、海外仓等外部系统 Token 必须加密存储并记录授权审计。
- 多租户数据访问必须强制校验 `tenant_id`。
- 库存调整、财务冲销、权限变更、数据导出必须记录审计日志。
- 所有服务必须接入结构化日志、指标和链路追踪。
- 首次上线前必须完成备份恢复演练、死信补偿流程和权限渗透检查。
