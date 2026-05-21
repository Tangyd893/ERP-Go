# ERP-Go

中文说明

ERP-Go 是一个基于 Go 微服务架构的跨境电商 ERP 系统设计项目。当前项目处于架构分析与工程设计阶段，目标不是先堆功能代码，而是先把商品、渠道、订单、库存、WMS 仓储、TMS 物流、采购、财务、报表和系统治理设计成可落地、可拆分、可验收的业务闭环。

一句话定位：

> 面向跨境电商卖家、品牌方和海外仓团队的一体化微服务 ERP 平台。

长期方向：

> 让跨境团队围绕“订单履约”和“库存资金账”完成商品管理、平台销售、采购补货、仓储作业、物流发运、财务核算和经营分析，而不是在平台后台、仓库系统、物流系统和表格之间反复搬运数据。

仓库地址：[https://github.com/Tangyd893/ERP-Go.git](https://github.com/Tangyd893/ERP-Go.git)

## 当前状态

这份 README 以当前仓库内容为准。项目目前只完成分析和设计文档，尚未实现后端服务、前端页面、Docker 编排和自动化测试。

完整设计入口请看 [文档索引](docs/文档索引.md)，总体架构请看 [项目架构设计](docs/项目架构设计.md)，后续实施路线请看 [实施路线与工程规范](docs/实施路线与工程规范.md)。

### 已经具备的设计内容

- 项目顶层目录已经确定为 `backend`、`frontend`、`docker`、`testing`、`docs`。
- 已明确从第一阶段开始采用微服务架构，不采用单体过渡。
- 已完成商品、渠道、订单、库存、仓储、物流、采购、财务、报表等领域划分。
- 已完成面向对象领域模型设计，覆盖聚合根、实体、值对象、领域服务、应用服务和适配器。
- 已完成微服务职责边界设计，明确各服务的数据所有权和调用关系。
- 已完成接口与事件设计，覆盖 REST、gRPC、领域事件、Saga、Outbox 和幂等策略。
- 已完成数据模型设计，覆盖核心表、库存流水、财务流水、审计日志、索引和对账规则。
- 已完成分阶段实施路线，定义首期订单履约闭环、二期采购财务闭环和三期智能化扩展。
- 已创建基础目录占位文件，方便后续在本地仓库中保留空目录。

### 仍需谨慎理解的能力

- 当前没有可运行的 API Gateway 或业务服务。
- 当前没有可启动的前端管理后台或 WMS PDA。
- 当前没有 Docker Compose 或 Kubernetes 配置。
- 当前没有数据库迁移文件，文档中的表结构仍是设计稿。
- 当前没有 OpenAPI、gRPC proto 或事件契约文件。
- 当前没有自动化测试，`testing` 目录目前仅作为后续测试工程入口。
- README 中的服务、端口、流程和技术栈均为设计目标，不代表已经落地。

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
| API Gateway | 统一入口、鉴权、路由、限流、审计、灰度 | 已设计，未实现 |
| IAM Service | 用户、角色、权限、菜单、登录、审计 | 已设计，未实现 |
| Tenant Service | 租户、组织、部门、数据范围、套餐配额 | 已设计，未实现 |
| Product Service | SPU、SKU、组合商品、条码、平台 SKU 映射 | 已设计，未实现 |
| Channel Service | 店铺授权、平台同步、库存推送、发货回传 | 已设计，未实现 |
| Order Service | 销售订单、审核、拆合单、异常、售后 | 已设计，未实现 |
| Inventory Service | 库存余额、锁定、释放、扣减、流水、盘点 | 已设计，未实现 |
| Warehouse Service | 入库、上架、波次、拣货、复核、打包、出库 | 已设计，未实现 |
| Transport Service | 物流商、物流规则、面单、发运、轨迹、运费 | 已设计，未实现 |
| Purchase Service | 供应商、采购计划、采购单、到货、质检、退货 | 已设计，未实现 |
| Finance Service | 应收应付、平台结算、成本分摊、汇率、利润 | 已设计，未实现 |
| Report Service | 销售、库存、仓储、物流、采购、利润报表 | 已设计，未实现 |
| File Service | 商品图片、面单、发票、导入导出文件 | 已设计，未实现 |
| Notification Service | 邮件、短信、站内信、Webhook、告警 | 已设计，未实现 |

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

当前阶段只需要 Git 和 Markdown 阅读工具。

后续实现阶段建议准备：

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

当前没有业务代码，因此暂不提供 `go test`、`npm test` 或 Docker 启动命令。

当前可做的本地检查：

```powershell
git status
Get-ChildItem docs
```

后续实现阶段建议补齐：

```powershell
cd backend
go vet ./...
go test ./...
go test -race ./...
```

```powershell
cd frontend
npm run lint
npm test
npm run build
```

## 当前开发重点

建议优先处理这些会直接影响长期架构质量和首期闭环的事项：

1. 创建后端微服务脚手架，统一服务入口、配置、日志、错误码、健康检查和追踪。
2. 补充各服务 OpenAPI 或 gRPC proto 契约草案。
3. 细化订单、库存、WMS 出库、TMS 发运和采购入库状态机。
4. 输出第一阶段数据库 migration 草案。
5. 设计 API Gateway 路由目录和服务注册发现方案。
6. 搭建本地 Docker Compose 基础设施，包括 PostgreSQL、Redis、RabbitMQ、MinIO 和 OpenSearch。
7. 建立 Outbox/Inbox 可靠消息模板。
8. 建立 `testing/contract` 契约测试基线。
9. 建立 `testing/integration` 订单履约集成测试基线。
10. 建立架构决策记录目录，记录消息队列、数据库、前端框架、服务治理等关键选择。

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
