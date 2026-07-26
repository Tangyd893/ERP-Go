# 安全审计摘要

> 静态代码审查结论（2026-07-26）。详细修复任务见 [todos.md](./todos.md) 安全相关 T-607~T-611、T-616~T-618。

## Critical

| 风险 | 说明 | 主要位置 |
| --- | --- | --- |
| 微服务信任客户端 Header | 业务服务仅读取 `X-Tenant-ID` / `X-User-ID`，绕过 Gateway 可伪造身份 | `backend/shared/middleware/middleware.go`、`backend/shared/server/bootstrap.go` |
| 跨租户 IDOR | 按资源 ID 查询/变更未强制 `tenant_id` 过滤 | `order-service`、`purchase-service`、`warehouse-service`、`file-service` repository/handler |
| Refresh Token 误用 | Gateway 曾未校验 JWT `type`（已修复：仅接受 `access`） | `backend/gateway/cmd/server/main.go` |

## High

| 风险 | 说明 | 缓解 |
| --- | --- | --- |
| 默认 JWT Secret | 非 production 环境可回退到源码默认值 | 生产设 `ENVIRONMENT=production` + 强 `JWT_SECRET` |
| 开发种子账号 | 迁移含 `admin/admin123` | 生产迁移排除 seed SQL |
| MinIO 公开 Bucket | dev compose 默认 anonymous public | 生产禁用 public bucket |
| RBAC 默认关闭 | `GATEWAY_RBAC_ENABLED` 未设时为 false | 生产启用 RBAC |
| CORS 默认 `*` | 未设 `CORS_ALLOWED_ORIGINS` 时允许任意来源 | 生产配置白名单 |

## Medium / Low

- Docker Compose 默认凭据映射到主机端口（仅本地开发）
- 数据库连接默认 `sslmode=disable`（生产应强制 TLS）
- Gateway 根路径 HTML 含 dev 登录提示（仅开发）
- 文件上传缺少大小/MIME 限制
- 登录无速率限制；Logout 不吊销 refresh token
- 前端 Token 存 localStorage（XSS 风险）

## 已确认良好

- 无已提交 `.env` 或真实密钥文件（`.gitignore` 已覆盖）
- SQL 使用 GORM 参数化，未发现字符串拼接注入
- 生产环境拒绝默认 JWT/DB 密码（T-609）
- K8s NetworkPolicy 骨架已存在（T-607）

## 验收环境说明

本地 `dev-stack` 使用开发默认凭据，**仅供功能验收，不可暴露到公网**。默认账号：`admin` / `admin123`，租户 `default`。
