# ERP-Go

中文说明

ERP-Go 是一个基于 Go 微服务架构的跨境电商 ERP 系统设计项目。当前项目处于架构分析与工程设计阶段，目标不是先堆功能代码，而是先把商品、渠道、订单、库存、WMS 仓储、TMS 物流、采购、财务、报表和系统治理设计成可落地、可拆分、可验收的业务闭环。

一句话定位：

> 面向跨境电商卖家、品牌方和海外仓团队的一体化微服务 ERP 平台。

长期方向：

> 让跨境团队围绕“订单履约”和“库存资金账”完成商品管理、平台销售、采购补货、仓储作业、物流发运、财务核算和经营分析，而不是在平台后台、仓库系统、物流系统和表格之间反复搬运数据。

仓库地址：[https://github.com/Tangyd893/ERP-Go.git](https://github.com/Tangyd893/ERP-Go.git)

## 当前状态

这份 README 以当前仓库内容为准。项目已完成 P0 工程底座，P1~P3 核心服务的数据库迁移和仓储层实现，IAM 服务已可端到端运行（登录/鉴权/用户管理）。正在持续推进剩余服务和应用层代码。

完整设计入口请看 [文档索引](docs/文档索引.md)，总体架构请看 [项目架构设计](docs/项目架构设计.md)，后续实施路线请看 [实施路线与工程规范](docs/实施路线与工程规范.md)。

### 已完成 (P0~P3 深化)

- 项目顶层目录：`backend`、`frontend`、`docker`、`testing`、`docs`。
- Go 服务模块（单模块 monorepo），Go 1.23+。
- `backend/shared` 公共组件：配置加载、结构化日志、统一错误码、请求中间件、统一响应、数据库连接助手、Outbox/Inbox 基类、分页、校验。
- API 网关：Gin 反向代理、JWT 鉴权跳过公开路由、请求 ID/追踪 ID 透传、优雅关闭。
- gRPC proto 契约目录（含 IAM 服务 proto 示例）。
- Vue 3 + TypeScript + Vite 前端工作区（npm workspaces）：
  - `admin-web`：ERP 管理后台（Element Plus、Vue Router、Pinia、登录页、主布局、18 个业务页面）。
  - `warehouse-pda`：WMS PDA 移动端。
  - `dashboard-web`：经营看板（ECharts）。
  - `shared`：共享 ProTable/ProForm 组件、API 客户端。
- Docker Compose 开发环境：PostgreSQL 16、Redis 7、RabbitMQ、MinIO、OpenSearch、Prometheus、Grafana、Jaeger。
- **数据库迁移**：25 个 SQL 文件，覆盖 IAM/Tenant/Product/Channel/Order/Inventory 共 29 张表，含默认管理员种子数据。
- **GORM 仓储实现**：7 个核心服务均已完成数据库仓储层，Inventory 支持事务行锁 + 幂等控制。
- **IAM 服务端到端可用**：登录/登出/刷新 Token/用户管理/角色管理/权限查询/操作审计均通过数据库模式验证。
- 完整的架构设计文档（8 份）、ADR 模板。
- 领域层单元测试：库存 9 个 + 订单状态机 6 个，全部通过。

### 尚未完成（后续阶段）

- IAM/Tenant 服务的应用层和 HTTP 层与仓储的完整连线（除 IAM 已连线外）。
- 业务服务的应用服务（App 层）和 HTTP 接口层代码。
- 异步 Worker 和事件消费逻辑。
- gRPC 代码生成和服务间调用。
- 自动化测试（集成测试、契约测试、端到端测试、性能测试）。
- Kubernetes 生产部署配置（仅网关有基础 yaml）。

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
Prometheus | Grafana | Loki | Tempo/Jaeger
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
| 后端 | Go + Gin，统一使用 Golang 实现后端微服务 |
| 网关 | API Gateway、统一鉴权、路由、限流、熔断 |
| 数据库 | PostgreSQL |
| 缓存 | Redis |
| 异步事件 | RabbitMQ，Outbox/Inbox 保障可靠消息 |
| 对象存储 | MinIO，生产可对接 S3 兼容存储 |
| 搜索分析 | OpenSearch，ClickHouse 作为后续分析库预留 |
| 前端 | Vue 3 + TypeScript + Vite + Pinia + Vue Router + Element Plus |
| 仓库终端 | Vue 3 H5 / PDA 扫码作业 |
| 监控 | OpenTelemetry、Prometheus、Grafana、Loki、Tempo/Jaeger |
| 中间件管理 | 统一由 `docker/` 目录管理，本地使用 Docker Compose |
| 生产部署 | Kubernetes |

## 服务规划

| 服务 | 设计职责 | 当前状态 |
|------|----------|----------|
| API Gateway | 统一入口、鉴权、路由、限流、审计、灰度 | 骨架可用，代理路由已注册 |
| IAM Service | 用户、角色、权限、菜单、登录、审计 | **数据库模式可运行** |
| Tenant Service | 租户、组织、部门、数据范围、套餐配额 | 迁移+仓储已完成 |
| Product Service | SPU、SKU、组合商品、条码、平台 SKU 映射 | 迁移+仓储已完成 |
| Channel Service | 店铺授权、平台同步、库存推送、发货回传 | 迁移+仓储已完成 |
| Order Service | 销售订单、审核、拆合单、异常、售后 | 迁移+仓储+状态机已完成 |
| Inventory Service | 库存余额、锁定、释放、扣减、流水、盘点 | 迁移+事务仓储已完成 |
| Warehouse Service | 入库、上架、波次、拣货、复核、打包、出库 | 领域模型已设计 |
| Transport Service | 物流商、物流规则、面单、发运、轨迹、运费 | 领域模型已设计 |
| Purchase Service | 供应商、采购计划、采购单、到货、质检、退货 | 领域模型已设计 |
| Finance Service | 应收应付、平台结算、成本分摊、汇率、利润 | 领域模型已设计 |
| Report Service | 销售、库存、仓储、物流、采购、利润报表 | 骨架已搭建 |
| File Service | 商品图片、面单、发票、导入导出文件 | 领域模型已设计 |
| Notification Service | 邮件、短信、站内信、Webhook、告警 | 领域模型已设计 |

## 项目目录

```text
ERP-Go/
  backend/    后端微服务、网关、异步任务、共享协议、数据库迁移
  frontend/   ERP 管理后台、WMS PDA、经营看板等前端应用
  docker/     Docker Compose、Kubernetes、Nginx、监控和部署配置
  testing/    集成测试、契约测试、端到端测试、性能测试和模拟服务
  docs/       架构设计、领域模型、接口事件、数据模型和工程规范文档
```

## 业务主链路

订单履约闭环：

```text
平台订单导入
  -> 订单审核
  -> 库存锁定
  -> WMS 创建出库任务
  -> TMS 匹配物流并生成面单
  -> 仓库拣货、复核、打包、出库
  -> 库存扣减
  -> 订单标记已发货
  -> 渠道服务回传平台
```

采购入库闭环：

```text
采购计划
  -> 采购单
  -> 到货通知
  -> WMS 收货
  -> 质检
  -> 上架
  -> 库存增加
  -> 采购成本入账
```

财务结算闭环：

```text
平台结算导入
  -> 匹配订单、退款、佣金和调整项
  -> 汇率换算
  -> 成本分摊
  -> 财务流水
  -> 订单利润 / SKU 利润 / 店铺利润
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

建议从文档索引开始：

```text
docs/文档索引.md
```

推荐阅读顺序：

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

# 运行单元测试
go test ./backend/... -count=1

# 启动开发中间件
docker compose -f docker/compose/docker-compose.dev.yml up -d postgres

# 运行数据库迁移（示例）
PGPASSWORD=erp123 psql -h localhost -p 5433 -U erp -d erp_go -f backend/migrations/iam/001_create_users.sql

# 启动 IAM 服务（需先完成迁移）
DATABASE_PORT=5433 SERVER_PORT=8096 JWT_SECRET=dev-secret go run ./backend/services/iam-service/cmd/server/

# 测试 IAM 登录
curl -X POST http://localhost:8096/api/v1/iam/login \
  -H "Content-Type: application/json" \
  -d '{"tenant_id":"t-default","username":"admin","password":"admin123"}'
```

### 前端

```bash
cd frontend
npm install
npm run dev:admin      # 启动管理后台 (端口 5173)
npm run dev:pda        # 启动 PDA (端口 5174)
npm run dev:dashboard  # 启动看板 (端口 5175)
npm run build:admin    # 构建管理后台
```

### 快捷命令

```bash
make help   # 查看所有可用命令
make all    # 编译 + 测试
make dev    # 启动 Docker 中间件
```

## 当前开发重点

1. 完善 IAM 和 Tenant 服务的 App 层和 HTTP 接口层，打通租户管理与组织架构。
2. 实现 Order Service 应用层（订单导入、审核、状态流转）。
3. 完成 Inventory Service HTTP 接口层，替换内存模拟为数据库存储。
4. 实现 WMS 出库与 TMS 发货核心流程（P4）。
5. 补充 gRPC/proto 代码生成和服务间调用。
6. 建立 `testing/integration` 订单履约集成测试基线。
7. 建立 `testing/contract` 契约测试基线。

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

当前项目尚未进入部署阶段。后续实现生产部署前必须注意：

- 生产环境必须通过环境变量、Kubernetes Secret 或专用密钥管理系统注入敏感配置。
- 不要提交 `.env`、数据库密码、平台 Token、物流商密钥、支付密钥或真实客户数据。
- Amazon、物流商、海外仓等外部系统 Token 必须加密存储并记录授权审计。
- 多租户数据访问必须强制校验 `tenant_id`。
- 库存调整、财务冲销、权限变更、数据导出必须记录审计日志。
- 所有服务必须接入结构化日志、指标和链路追踪。
- 首次上线前必须完成备份恢复演练、死信补偿流程和权限渗透检查。
