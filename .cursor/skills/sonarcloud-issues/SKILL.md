---
name: sonarcloud-issues
description: >-
  Fetches open SonarCloud issues via the official /api/issues/search API (same
  data as the project issues web UI) and formats them for LLM consumption or
  docs updates. Use when the user shares sonarcloud.io/project/issues links,
  asks for Sonar issue lists, quality debt triage, or updating docs/待办清单.md
  from SonarCloud.
---

# SonarCloud Issues Fetch

Do **not** scrape the SonarCloud HTML UI. The [issues page](https://sonarcloud.io/project/issues?issueStatuses=OPEN%2CCONFIRMED&id=Tangyd893_ERP-Go) uses the public API:

```http
GET https://sonarcloud.io/api/issues/search
```

## Quick start

From repo root, run the bundled script (paginates automatically):

```powershell
.\.cursor\skills\sonarcloud-issues\scripts\fetch-issues.ps1
```

```bash
bash .cursor/skills/sonarcloud-issues/scripts/fetch-issues.sh
```

Write markdown to a file (e.g. for `docs/待办清单.md` refresh):

```powershell
.\.cursor\skills\sonarcloud-issues\scripts\fetch-issues.ps1 -OutFile .cache/sonar-issues.md
```

Raw JSON (all pages merged):

```powershell
.\.cursor\skills\sonarcloud-issues\scripts\fetch-issues.ps1 -Format json -OutFile .cache/sonar-issues.json
```

## Parameters

| Parameter | Env fallback | Default | Meaning |
| --- | --- | --- | --- |
| `-ProjectKey` | `SONAR_PROJECT_KEY` | `Tangyd893_ERP-Go` | `componentKeys` |
| `-Branch` | `SONAR_BRANCH` | *(empty)* | Omit = all branches; set `main` to match branch analysis |
| `-IssueStatuses` | — | `OPEN,CONFIRMED` | Same as UI `issueStatuses` query param |
| `-Types` | — | *(all)* | e.g. `VULNERABILITY`, `CODE_SMELL`, `BUG` |
| `-Severities` | — | *(all)* | e.g. `BLOCKER,CRITICAL` |
| `-PageSize` | — | `500` | Max per request (API cap 500) |

Authentication: **not required** for public project issue search. For private org data or elevated rate limits, set `SONAR_TOKEN` and pass header `Authorization: Bearer $SONAR_TOKEN` (extend script if needed).

## Agent workflow

When the user asks to pull Sonar issues or update quality docs:

1. Run `fetch-issues.ps1` (or `.sh`) with `-OutFile .cache/sonar-issues.md`.
2. Read the generated markdown — it is structured for LLM triage (summary → by severity → by file).
3. Optionally call measures API for ratings:

```powershell
Invoke-RestMethod "https://sonarcloud.io/api/measures/component?component=Tangyd893_ERP-Go&branch=main&metricKeys=bugs,vulnerabilities,code_smells,security_rating,reliability_rating,sqale_rating,duplicated_lines_density,ncloc"
```

4. Map findings to project tasks (e.g. `docs/待办清单.md` §7 `T-xxx`) when updating docs.

## API reference (minimal)

| Query param | Example | Notes |
| --- | --- | --- |
| `componentKeys` | `Tangyd893_ERP-Go` | Required project key |
| `issueStatuses` | `OPEN,CONFIRMED` | Matches UI filter |
| `branch` | `main` | Optional branch filter |
| `types` | `VULNERABILITY` | Optional |
| `severities` | `BLOCKER,CRITICAL` | Optional |
| `rules` | `secrets:S6698` | Optional rule filter |
| `p` | `1` | Page index (1-based) |
| `ps` | `500` | Page size |
| `facets` | `types,severities,rules,languages` | Summary counts in response |

Pagination: repeat requests with `p=1,2,...` until `issues` array empty or `p * ps >= total`.

## Output contract (markdown)

The script emits:

1. **Header** — project, branch, statuses, total, effort/debt, fetch time, Sonar UI link.
2. **Facets** — severity / type / language / top rules tables.
3. **Issues by severity** — `BLOCKER` → `CRITICAL` → … each issue as one bullet: `path:line | rule | type | message | effort`.
4. **Issues by file** — grouped path for implementation planning.

Use this output directly in chat or paste sections into `docs/待办清单.md` §3–§5.

## Troubleshooting

| Problem | Fix |
| --- | --- |
| `total: 0` | Wrong `ProjectKey`; or use `-Branch main` if issues are branch-scoped |
| Timeout | Reduce `-PageSize`; filter `-Types` / `-Severities` |
| 401 | Set `SONAR_TOKEN` for private projects |

## Additional resources

- Script implementation: [scripts/fetch-issues.ps1](scripts/fetch-issues.ps1), [scripts/fetch-issues.sh](scripts/fetch-issues.sh)
- ERP-Go Sonar UI: https://sonarcloud.io/project/issues?issueStatuses=OPEN%2CCONFIRMED&id=Tangyd893_ERP-Go
