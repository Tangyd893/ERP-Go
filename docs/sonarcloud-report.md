# SonarCloud 代码质量报告

> 2026-06-30（API 拉取 + QG 状态；最近 CI 分析 `9615292`，2026-06-19T12:54:53Z）

## 概览

- **项目**: ERP-Go
- **Project Key**: Tangyd893_ERP-Go
- **分支**: main
- **质量门**: ❌ ERROR（2 项失败）
- **代码行 (ncloc)**: 23,235
- **未解决 Issue**: **10**（较 2026-06-19 文档口径 20 条 **−10**）
- **技术债务**: **36 min**（`sqale_index`；较 130 min **−94**）
- **最近分析**: `2026-06-19T12:54:53Z`
- **revision**: `961529208e86beb6ff768a6192e58ea1b3069555`
- **Issue 明细**: [`.cache/sonar-issues.md`](../.cache/sonar-issues.md)（本地生成，不提交 Git）

## 指标

| 指标 | 值 | 较 2026-06-19（`334c1c1`） |
|------|-----|--------|
| Bug | 0 | — |
| 漏洞 | 1 | — |
| 代码异味 | 9 | 19 → **9** ✅ |
| 重复率 | 1.7% | — |
| 新代码重复率 | 2.9% | ✅ |
| Security Hotspot 待审 | 1 | — |
| 可靠性评级 | A (1.0) | — |
| 安全性评级 | E (5.0) | — |
| 可维护性评级 | A (1.0) | — |

## Quality Gate 条件

| 条件 | 阈值 | 实际 | 状态 |
|------|------|------|------|
| 新代码可靠性 | A | A | ✅ |
| 新代码安全性 | A | E | ❌ |
| 新代码可维护性 | A | A | ✅ |
| 新代码重复行密度 | ≤ 3% | 2.9% | ✅ |
| 新代码 Hotspot 审查率 | 100% | 0% | ❌ |

## 严重级别分布（开放 Issue）

| 级别 | 数量 | 规则 |
|------|------|------|
| BLOCKER | 1 | `secrets:S8215` |
| CRITICAL | 9 | `plsql:S1192`（migration / archive SQL） |
| MAJOR | 0 | — |
| MINOR | 0 | — |

**语言分布**：PL/SQL 9、Secrets 1；**Go / TypeScript / Shell 开放 Issue 已清零**。

## 关键阻断项

| 任务 | 问题 | 建议 |
|------|------|------|
| T-412 / T-002 | `iam-service/migrations/002_seed_dev_admin.sql:19`（`secrets:S8215`） | Sonar UI **False Positive**（开发种子）；或 push 后确认 `sonar.security.exclusions=**/migrations/**` |
| T-411 | `docker/k8s/gateway-deployment.yaml:34`（`kubernetes:S5332`） | CI `sonar-review-hotspots.sh` 或 UI **Review → Safe** |

## 本地验证（2026-06-30）

- Go：`go vet` / `go test` / `go build` 全通过
- 前端：`warehouse-pda` typecheck 失败（`ShipConfirm.vue` 调用 `store.confirmShip`）— **已修复**（`confirmShip` 挂回 store）；`verify.ps1` 全量通过（2026-06-30）

## 参考链接

- [SonarCloud 概览（main）](https://sonarcloud.io/summary/overall?id=Tangyd893_ERP-Go&branch=main)
- [开放 Issues](https://sonarcloud.io/project/issues?issueStatuses=OPEN%2CCONFIRMED&id=Tangyd893_ERP-Go&branch=main)
- 拉取命令：`.\.cursor\skills\sonarcloud-issues\scripts\fetch-issues.ps1 -Branch main -OutFile .cache/sonar-issues.md`
