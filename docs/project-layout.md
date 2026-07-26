# 项目目录说明

顶层目录固定为五层，新增内容应归入对应目录，避免在根目录堆积临时文件。

```text
ERP-Go/
├── backend/          Go 微服务、Gateway、shared 库、迁移
│   ├── gateway/      API Gateway
│   ├── shared/       配置、中间件、Outbox、事件、工作流
│   ├── services/     13 个业务微服务（各含 cmd/internal/migrations）
│   ├── migrations/   Outbox 共享迁移 + 历史归档（只读）
│   ├── testing/      契约测试、集成测试
│   └── proto/        gRPC proto 草案
├── frontend/         Vue 3 前端（npm workspaces）
│   ├── admin-web/    ERP 管理后台（5173）
│   ├── warehouse-pda/ WMS PDA（5174）
│   ├── dashboard-web/ 经营看板（5175）
│   └── shared/       共享组件与 API client
├── docker/           本地 Compose、K8s、Helm、监控配置
├── testing/          E2E（Playwright）、性能测试骨架
├── docs/             设计文档、runbooks、todos、安全审计
├── scripts/          dev-stack、迁移、verify、构建脚本
└── .cache/           本地缓存（Go 模块/构建、服务日志，已 gitignore）
```

## 文档入口

| 用途 | 文件 |
| --- | --- |
| 项目介绍 | [README.md](../README.md) |
| 设计文档索引 | [文档索引.md](./文档索引.md) |
| 任务与进度 | [todos.md](./todos.md) |
| 安全审计 | [security-audit.md](./security-audit.md) |
| 旧待办链接 | [待办清单.md](./待办清单.md) → 重定向至 todos |

## 不应提交的内容

- `.env`、`.cache/`、`.understand-anything/`、`node_modules/`、`dist/`
- 真实平台 Token、生产密钥、客户数据
