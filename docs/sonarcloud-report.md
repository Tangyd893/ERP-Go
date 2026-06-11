# SonarCloud 代码质量报告

> 2026-06-11

## 概览

- **项目**: ERP-Go
- **Project Key**: Tangyd893_ERP-Go
- **分支**: main
- **质量门**: ❌ ERROR
- **代码行 (ncloc)**: 23,182
- **未解决 Issue**: 57
- **技术债务**: 7h 41min（461 min / API effort 554 min）
- **最近分析**: `2026-06-11T08:33:35Z`
- **最近分析 revision**: `edcafbe`（commit: P4-P9 全流程闭环）

## 指标

| 指标 | 值 | 较 06-10 |
|------|-----|----------|
| Bug | 1 | +1 |
| 漏洞 | 2 | +2 |
| 代码异味 | 54 | −8 |
| 重复率 | 1.7% | −2.6% |
| 新代码重复率 | **2.97%** | ✅ 达标（≤3%） |
| Security Hotspot 待审 | 2 | +1 |
| 测试覆盖率 | 未上报 | — |
| 可靠性评级 | **C** | A → C |
| 安全性评级 | **E** | A → E |
| 可维护性评级 | A | — |

## Quality Gate 条件

| 条件 | 阈值 | 实际 | 状态 |
|------|------|------|------|
| 新代码可靠性 | A | **C** | ❌ |
| 新代码安全性 | A | **E** | ❌ |
| 新代码可维护性 | A | A | ✅ |
| 新代码重复行密度 | ≤ 3% | **3.0%** | ✅ |
| 新代码 Hotspot 审查率 | 100% | **0%**（**2** 个待审） | ❌ |

## 严重级别分布

| 级别 | 数量 | 占比 |
|------|------|------|
| 🔴 BLOCKER | 2 | 3.5% |
| 🟠 CRITICAL | 23 | 40.4% |
| 🟡 MAJOR | 24 | 42.1% |
| 🔵 MINOR | 8 | 14.0% |
| ⚪ INFO | 0 | 0.0% |

## 问题类型分布

| 类型 | 数量 | 占比 |
|------|------|------|
| Bug | 1 | 1.8% |
| Vulnerability | 2 | 3.5% |
| Code Smell | 54 | 94.7% |

## Issue 分布（按语言）

| 类别 | 数量 |
|------|------|
| Go | 24 |
| PowerShell | 10 |
| PL/SQL | 9 |
| TypeScript | 8 |
| Kubernetes | 2 |
| Shell | 2 |
| Secrets | 1 |
| Web | 1 |

## 关键问题

### 🔴 BLOCKER

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `backend/scripts/genhash/main.go` | 10 | `go:S6437` | 示例密码字面量被识别为泄露密钥 |
| `backend/services/iam-service/migrations/002_seed_dev_admin.sql` | 19 | `secrets:S8215` | 种子 bcrypt 哈希仍在扫描范围内 |

### 🟠 CRITICAL（Top 15）

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `purchase_repository.go` | 77 | `go:S1192` | `"id = ?"` 重复 7 次 |
| `tenant_handler.go` | 85 | `go:S1192` | 占位响应消息重复 5 次 |
| `purchase_app_service.go` | 76, 145 | `go:S1192` | 采购/入库错误消息重复 |
| `inventory_repository.go` | 88, 102 | `go:S1192` | 查询/更新库存错误消息重复 |
| `pg_store.go` | 11, 92 | `go:S1192` | `"status = ?"` / `"id = ?"` |
| `transport_app_service.go` | 90 | `go:S3776` | 认知复杂度 22 |
| `purchase_integration_test.go` | 9 | `go:S3776` | 认知复杂度 19 |
| IAM / channel / product / tenant migrations | — | `plsql:S1192` | SQL 字面量重复 |
| `docs/archive/007_seed_default_data.sql` | 6–46 | `plsql:S1192` | 归档种子（应 exclusion） |

### Security Hotspot（待审）

| 文件 | 行 | 规则 | 类别 |
|------|-----|------|------|
| `docker/k8s/gateway-deployment.yaml` | 34 | `kubernetes:S5332` | 集群内 HTTP 明文 |
| `frontend/warehouse-pda/src/stores/warehouse.ts` | 44 | `typescript:S2245` | `Math.random()` 用于幂等键 |

### 🐛 Bug

| 文件 | 行 | 规则 | 描述 |
|------|-----|------|------|
| `finance_app_service.go` | 113 | `go:S1656` | 无意义自赋值（导致新代码可靠性 C） |

## 高频规则（Top 10）

| 规则 | 级别 | 类型 | 命中数 | 占比 |
|------|------|------|--------|------|
| `go:S1192` | CRITICAL | CODE_SMELL | 12 | 21.1% |
| `powershelldre:S8677` | MAJOR | CODE_SMELL | 10 | 17.5% |
| `plsql:S1192` | CRITICAL | CODE_SMELL | 9 | 15.8% |
| `go:S107` | MAJOR | CODE_SMELL | 6 | 10.5% |
| `typescript:S7764` | MINOR | CODE_SMELL | 4 | 7.0% |
| `go:S3776` | CRITICAL | CODE_SMELL | 2 | 3.5% |
| `shelldre:S7679` | MAJOR | CODE_SMELL | 2 | 3.5% |
| `kubernetes:S6596/S6897` | MAJOR | CODE_SMELL | 2 | 3.5% |
| `typescript:S7721` | MAJOR | CODE_SMELL | 1 | 1.8% |
| `secrets:S8215` | BLOCKER | VULNERABILITY | 1 | 1.8% |

## Issue 数 Top 文件

| 文件 | Issue 数 | 主要规则 |
|------|----------|----------|
| `scripts/dev-stack.ps1` | 9 | `powershelldre:S8677` |
| `frontend/warehouse-pda/src/stores/warehouse.ts` | 4 | `S7721` + `S7764` |
| `finance_app_service.go` | 4 | `go:S107` + `go:S1656` (BUG) |
| `docs/archive/007_seed_default_data.sql` | 3 | `plsql:S1192` |
| `finance.go` | 2 | `go:S107` |

## 修复建议（与待办联动）

1. **T-412 / T-002（P0）**：扩展 `sonar.exclusions` 覆盖 `backend/scripts/genhash/**`、`**/migrations/**`（或 `iam-service/migrations/002_seed_dev_admin.sql`）；消除 2 个 BLOCKER Vulnerability，恢复安全评级。
2. **T-413（P0）**：修复 `finance_app_service.go:113` 自赋值 Bug，恢复新代码可靠性评级。
3. **T-411（P0）**：Sonar UI 审查 2 个 Hotspot（`gateway-deployment.yaml:34` → Safe；`warehouse.ts:44` → Safe 或改用 `crypto.randomUUID()`）。
4. **T-410**：✅ 新代码重复率已达标；保持常量提取节奏，避免回退。
5. **T-422**：`dev-stack.ps1` 9 处 `Write-Host` → `Write-Output`（Sonar 新扫入 PowerShell 规则）。
6. **T-402**：剩余 `go:S1192` 12 处（inventory/purchase/warehouse/outbox 等）。
7. **T-420**：`warehouse.ts` `confirmShip` 外提 + `globalThis`  MINOR。

## 参考链接

- [SonarCloud 概览（main）](https://sonarcloud.io/summary/overall?id=Tangyd893_ERP-Go&branch=main)
- [开放 Issues](https://sonarcloud.io/project/issues?issueStatuses=OPEN%2CCONFIRMED&id=Tangyd893_ERP-Go&branch=main)
- 本地 Issue 明细：`.cache/sonar-issues.md`
