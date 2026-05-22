# ERP-Go

中文说明

ERP-Go 是一个基于 Go 微服务架构的跨境电商 ERP 系统项目，面向跨境卖家、品牌方和海外仓团队，覆盖商品、渠道、订单、库存、WMS 仓储、TMS 物流、采购、财务、报表和系统治理等业务域。

一句话定位：

> 面向跨境电商卖家、品牌方和海外仓团队的一体化微服务 ERP 平台。

长期方向：

> 让跨境团队围绕“订单履约”和“库存资金账”完成商品管理、平台销售、采购补货、仓储作业、物流发运、财务核算和经营分析，而不是在平台后台、仓库系统、物流系统和表格之间反复搬运数据。

仓库地址：[https://github.com/Tangyd893/ERP-Go.git](https://github.com/Tangyd893/ERP-Go.git)

## 当前状态

本节基于 2026-05-22 对仓库的检查更新。

### 状态快照

- Git 状态：检查开始时位于 `main`，跟踪 `origin/main`，工作树干净。
- 后端：`go test ./...` 通过，`go vet ./...` 通过，`go build ./...` 使用仓库内临时 `GOCACHE` 后通过。
- 实际测试覆盖：目前只有 Inventory 领域测试 9 个、Order 状态机测试 6 个；多数服务包仍无测试文件。
- 前端类型检查：`warehouse-pda` 和 `dashboard-web` 通过；`admin-web` 失败，原因是 `@erp/shared` 类型解析缺失，以及 `element-plus/dist/locale/zh-cn.mjs` 缺少声明。
- 前端构建：当前本地 `node_modules` 缺少 Windows Rollup 可选原生包 `@rollup/rollup-win32-x64-msvc`，Vite 构建会失败，需要重新安装依赖或补齐可选依赖。
- 数据库迁移：`backend/services/*/migrations` 下已有 12 个服务级 migration，合计 53 张表；`backend/migrations` 下仍保留 15 个历史/聚合 migration，合计 24 张表，和服务级 migration 存在重叠，后续需要统一迁移入口。

### 已完成

- P0 工程底座：`backend`、`frontend`、`docker`、`testing`、`docs` 顶层目录，Go 单模块 monorepo，npm workspaces。
- `backend/shared` 公共组件：配置、日志、错误码、响应、数据库连接、中间件、分页、校验、Outbox/Inbox 基类。
- API Gateway 骨架：健康检查、认证头检查、请求 ID/追踪 ID/租户/用户上下文中间件、部分服务代理路由。
- 13 个后端服务目录：IAM、Tenant、Product、Channel、Order、Inventory、Warehouse、Transport、Purchase、Finance、Report、File、Notification。
- IAM Service：应用层、仓储层、HTTP 路由已连线；数据库可用时支持登录、刷新 Token、登出、用户、角色、权限和审计能力。
- Tenant/Product/Channel/Order/Inventory：领域模型和仓储层已推进，其中 Order 有状态机，Inventory 仓储支持事务行锁和幂等控制。
- Warehouse/Transport/Purchase/Finance/File/Notification：领域模型与服务级 migration 已建立，HTTP 目前主要是规划接口。
- Frontend：`admin-web` 18 个视图、`warehouse-pda` 2 个视图、`dashboard-web` 基础看板页、`frontend/shared` 共享组件和 API client。
- Docker Compose 开发环境：PostgreSQL、Redis、RabbitMQ、MinIO、OpenSearch、Prometheus、Grafana、Loki、Jaeger。
- 架构设计文档、ADR 模板、测试目录骨架已建立。

### 尚未完成

- 迁移体系需要归一：当前同时存在根级 `backend/migrations` 和服务级 `backend/services/*/migrations`。
- Gateway 默认代理路由尚未覆盖 Purchase、Finance、Report、Notification 等后续服务。
- 除 IAM 外，多数服务的应用层和 HTTP 层仍未真正连到仓储；Product/Channel/Order/Tenant 等接口仍以占位响应为主。
- Inventory HTTP 层仍在使用内存模拟数据，数据库仓储已实现但尚未接入当前 handler。
- 异步 Worker、事件消费、gRPC 代码生成和服务间调用尚未落地。
- 集成测试、契约测试、端到端测试、性能测试仍为空目录或说明文档。
- 前端需要修复 `admin-web` 类型解析、补齐 Rollup Windows 可选依赖，再恢复完整构建验证。
- Kubernetes 生产部署仍是早期骨架，目前仅有网关基础 yaml。

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
| API Gateway | 统一入口、鉴权、路由、追踪上下文 | 骨架可编译；已代理 IAM/Tenant/Product/Channel/Order/Inventory/Warehouse/Transport/File，后续服务路由待补 |
| IAM Service | 用户、角色、权限、登录、审计 | App + Repository + HTTP 已连线；数据库不可用时降级为占位模式 |
| Tenant Service | 租户、组织、部门、岗位 | Migration + Repository 完成；HTTP 仍是占位接口 |
| Product Service | SPU、SKU、平台 SKU 映射 | Migration + Repository 完成；HTTP 仍是占位接口 |
| Channel Service | 店铺授权、平台同步、订单导入任务 | Migration + Domain + Repository 完成；HTTP 仍是占位接口 |
| Order Service | 销售订单、状态机、审核、异常 | Migration + Domain + Repository + 状态机测试完成；HTTP 仍是占位接口 |
| Inventory Service | 库存余额、锁定、释放、扣减、流水 | Migration + Domain + Repository 完成；HTTP 当前仍使用内存模拟，数据库仓储待接入 |
| Warehouse Service | 入库、上架、拣货、复核、打包、出库 | Domain + Migration + 占位 API 已建立 |
| Transport Service | 物流商、物流规则、面单、发运、轨迹 | Domain + Migration + 占位 API 已建立 |
| Purchase Service | 供应商、采购单、到货、质检、入库 | Domain + Migration + 占位 API 已建立 |
| Finance Service | 应收应付、结算、成本、利润、流水 | Domain + Migration + 占位 API 已建立 |
| Report Service | 销售、库存、仓储、利润报表 | 服务骨架和占位 API 已建立 |
| File Service | 商品图片、面单、发票、导入导出文件 | Domain + Migration + 占位 API 已建立 |
| Notification Service | 邮件、短信、站内信、Webhook、告警 | Domain + Migration + 占位 API 已建立 |

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

## 验证命令

### 后端

```bash
# 编译所有服务
go build ./...

# 静态分析
go vet ./...

# 运行测试
go test ./...

# 启动开发中间件
docker compose -f docker/compose/docker-compose.dev.yml up -d

# 启动 IAM 服务（需先完成数据库迁移和种子数据）
DATABASE_PORT=5433 SERVER_PORT=8081 JWT_SECRET=dev-secret go run ./backend/services/iam-service/cmd/server/
```

### 前端

```bash
npm install
npm run dev:admin      # 管理后台，端口 5173
npm run dev:pda        # PDA，端口 5174
npm run dev:dashboard  # 看板，端口 5175
npm run typecheck
npm run build:admin
```

当前已知前端验证问题：

- `admin-web` 类型检查需要补齐 `@erp/shared` 类型解析和 Element Plus locale 声明。
- Windows 环境下 Vite/Rollup 构建需要安装 `@rollup/rollup-win32-x64-msvc` 可选依赖。

### 快捷命令

```bash
make help   # 查看所有可用命令
make all    # 编译 + 测试
make dev    # 启动 Docker 中间件
```

## 当前开发重点

1. 统一 migration 策略，明确根级迁移和服务级迁移的取舍。
2. 修复前端依赖与类型解析，恢复 `admin-web` typecheck/build。
3. 将 Inventory HTTP handler 从内存模拟切换到数据库仓储。
4. 为 Tenant/Product/Channel/Order 补齐应用层和真实 HTTP 接口。
5. 补齐 Gateway 对 Purchase、Finance、Report、Notification 的代理路由。
6. 建立订单导入 -> 审核 -> 库存锁定的集成测试基线。
7. 继续推进 P4：WMS 出库、TMS 发货、库存扣减、订单发货状态和平台回传。

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
