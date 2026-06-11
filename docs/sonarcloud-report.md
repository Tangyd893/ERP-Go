# SonarCloud 代码质量报告

> 2026-06-11（CI 重扫 `b31368e`）

## 概览

- **项目**: ERP-Go
- **Project Key**: Tangyd893_ERP-Go
- **分支**: main
- **质量门**: ❌ ERROR（2 项失败）
- **代码行 (ncloc)**: 23,105
- **未解决 Issue**: 37（较上次 −20）
- **技术债务**: 5h 13min（313 min）
- **最近分析**: `2026-06-11T08:56:57Z`
- **revision**: `b31368e`

## 指标

| 指标 | 值 | 较上次 |
|------|-----|--------|
| Bug | 0 | −1 ✅ |
| 漏洞 | 2 | — |
| 代码异味 | 35 | −19 |
| 重复率 | 1.7% | — |
| 新代码重复率 | 2.97% | ✅ |
| Security Hotspot 待审 | 2 | — |
| 可靠性评级 | A | C → A ✅ |
| 安全性评级 | E | — |
| 可维护性评级 | A | — |

## Quality Gate 条件

| 条件 | 阈值 | 实际 | 状态 |
|------|------|------|------|
| 新代码可靠性 | A | A | ✅ |
| 新代码安全性 | A | E | ❌ |
| 新代码可维护性 | A | A | ✅ |
| 新代码重复行密度 | ≤ 3% | 3.0% | ✅ |
| 新代码 Hotspot 审查率 | 100% | 0% | ❌ |

## 严重级别分布

| 级别 | 数量 |
|------|------|
| BLOCKER | 2 |
| CRITICAL | 21 |
| MAJOR | 10 |
| MINOR | 4 |

## 关键阻断项

| 任务 | 问题 |
|------|------|
| T-412 / T-002 | `genhash/main.go:12`（S6437）、`002_seed_dev_admin.sql:19`（S8215） |
| T-411 | Hotspot：`gateway-deployment.yaml:34`、`warehouse.ts:44` |

## 参考链接

- [SonarCloud 概览（main）](https://sonarcloud.io/summary/overall?id=Tangyd893_ERP-Go&branch=main)
- [开放 Issues](https://sonarcloud.io/project/issues?issueStatuses=OPEN%2CCONFIRMED&id=Tangyd893_ERP-Go&branch=main)
- 本地明细：`.cache/sonar-issues.md`
