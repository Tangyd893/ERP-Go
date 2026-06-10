# SonarCloud 代码质量报告

> 2026-06-10

## 概览

- **项目**: ERP-Go
- **Project Key**: Tangyd893_ERP-Go
- **分支**: main
- **质量门**: ❌ ERROR
- **代码行**: 19,951
- **未解决 Issue**: 62
- **技术债务**: 6h 17min（377 min）
- **最近分析 revision**: `81235cf4`（2026-06-05）

## 指标

| 指标 | 值 |
|------|-----|
| Bug | 0 |
| 漏洞 | 0 |
| 代码异味 | 62 |
| 重复率 | 4.3% |
| 新代码重复率 | 5.0% |
| 认知复杂度 | 1,931 |
| Security Hotspot | 1（审查率 0%） |
| 测试覆盖率 | 未上报 |
| 可靠性评级 | A |
| 安全性评级 | A |
| 可维护性评级 | A |

## Quality Gate 条件

| 条件 | 阈值 | 实际 | 状态 |
|------|------|------|------|
| 新代码可靠性 | A | A | ✅ |
| 新代码安全性 | A | A | ✅ |
| 新代码可维护性 | A | A | ✅ |
| 新代码重复行密度 | ≤ 3% | **5.0%** | ❌ |
| 新代码 Hotspot 审查率 | 100% | **0%** | ❌ |

## 严重级别分布

| 级别 | 数量 | 占比 |
|------|------|------|
| 🔴 BLOCKER | 0 | 0.0% |
| 🟠 CRITICAL | 33 | 53.2% |
| 🟡 MAJOR | 20 | 32.3% |
| 🔵 MINOR | 9 | 14.5% |
| ⚪ INFO | 0 | 0.0% |

## 问题类型分布

| 类型 | 数量 | 占比 |
|------|------|------|
| Bug | 0 | 0.0% |
| Vulnerability | 0 | 0.0% |
| Code Smell | 62 | 100.0% |

## Issue 分布（按语言）

| 类别 | 数量 |
|------|------|
| Go | 28 |
| TypeScript | 18 |
| PL/SQL | 8 |
| PowerShell | 4 |
| Kubernetes | 2 |
| Shell | 1 |
| Web | 1 |

## 关键问题

### 🔴 BLOCKER

无

### 🟠 CRITICAL（Top 20）

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `finance_repository.go` | 29 | `go:S1192` | 重复字面量 `"tenant_id = ?"`（5 次） |
| `finance_repository.go` | 32 | `go:S1192` | 重复字面量 `"created_at DESC"`（5 次） |
| `inventory_handler.go` | 31 | `go:S1192` | 重复字面量 `"sku-001"` / `"wh-001"` |
| `inventory_handler.go` | 105 | `go:S1192` | 重复字面量 `"SKU库存未找到"` |
| `inventory_repository.go` | 76–100 | `go:S1192` | GORM 条件串与错误消息重复 |
| `p4_adapters.go` | 55 | `go:S1192` | `"Content-Type"` / `"application/json"` |
| `p4_outbound_flow.go` | 165–173 | `go:S1192` | 幂等/解析错误消息重复 |
| `tenant_handler.go` | 85 | `go:S1192` | 占位响应消息重复（5 次） |
| `docs/archive/007_seed_default_data.sql` | 6–46 | `plsql:S1192` | 归档种子 SQL 字面量重复（3 处） |
| `iam-service/migrations/001_init.sql` | 21–135 | `plsql:S1192` | migration SQL 字面量重复 |
| `auth_service.go` | 57 | `go:S1192` | `"用户已被禁用"` 重复 |
| IAM 三 repository | 36–53 | `go:S1192` | `"id = ?"` 等条件串重复 |
| `order_repository.go` | 105 | `go:S1192` | `"order_id = ?"` |
| `warehouse_repository.go` | 65–77 | `go:S1192` | outbound/id 条件串 |
| `pg_store.go` | 63 | `go:S1192` | `"status = ?"` |

## 高频规则（Top 10）

| 规则 | 级别 | 类型 | 命中数 | 占比 |
|------|------|------|--------|------|
| `go:S1192` | CRITICAL | CODE_SMELL | 25 | 40.3% |
| `typescript:S7721` | MAJOR | CODE_SMELL | 10 | 16.1% |
| `plsql:S1192` | CRITICAL | CODE_SMELL | 8 | 12.9% |
| `typescript:S7772` | MINOR | CODE_SMELL | 3 | 4.8% |
| `powershelldre:S8677` | MAJOR | CODE_SMELL | 3 | 4.8% |
| `typescript:S3863` | MINOR | CODE_SMELL | 2 | 3.2% |
| `kubernetes:S6596` | MAJOR | CODE_SMELL | 1 | 1.6% |
| `kubernetes:S6897` | MAJOR | CODE_SMELL | 1 | 1.6% |
| `go:S107` | MAJOR | CODE_SMELL | 1 | 1.6% |
| `godre:S8242` | MAJOR | CODE_SMELL | 1 | 1.6% |

## 修复建议（与待办联动）

1. **T-410**：合入本地 `go:S1192` 常量提取 + `sonar.exclusions` 扩展（`docs/archive/**` 已在工作区配置），push 后重扫验证新代码重复率 ≤ 3%。
2. **T-411**：在 Sonar UI 对 `docker/k8s/gateway-deployment.yaml:34` Hotspot 执行 Review → Safe。
3. **T-420/T-421**：前端 Store 外提 async 函数、三端 `node:path` 等小债。
4. **T-422**：PowerShell/Shell 脚本规范（`migrate.ps1`、`verify.ps1`）。

## 参考链接

- [SonarCloud 概览（main）](https://sonarcloud.io/summary/overall?id=Tangyd893_ERP-Go&branch=main)
- [开放 Issues](https://sonarcloud.io/project/issues?issueStatuses=OPEN%2CCONFIRMED&id=Tangyd893_ERP-Go&branch=main)
- 本地 Issue 明细：`.cache/sonar-issues.md`
