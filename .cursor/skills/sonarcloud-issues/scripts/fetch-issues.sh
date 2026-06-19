#!/usr/bin/env bash
# Fetch SonarCloud open issues via /api/issues/search and format for LLM/docs.
# Usage:
#   ./fetch-issues.sh
#   ./fetch-issues.sh --branch main --out .cache/sonar-issues.md
#   ./fetch-issues.sh --format json --out .cache/sonar-issues.json

set -euo pipefail

PROJECT_KEY="${SONAR_PROJECT_KEY:-Tangyd893_ERP-Go}"
BRANCH="${SONAR_BRANCH:-}"
ISSUE_STATUSES="OPEN,CONFIRMED"
TYPES=""
SEVERITIES=""
FORMAT="markdown"
PAGE_SIZE=500
OUT_FILE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --project) PROJECT_KEY="$2"; shift 2 ;;
    --branch) BRANCH="$2"; shift 2 ;;
    --statuses) ISSUE_STATUSES="$2"; shift 2 ;;
    --types) TYPES="$2"; shift 2 ;;
    --severities) SEVERITIES="$2"; shift 2 ;;
    --format) FORMAT="$2"; shift 2 ;;
    --out) OUT_FILE="$2"; shift 2 ;;
    -h|--help)
      grep '^#' "$0" | head -6 | sed 's/^# //'
      exit 0
      ;;
    *) echo "Unknown option: $1" >&2; exit 1 ;;
  esac
done

command -v jq >/dev/null || { echo "jq is required" >&2; exit 1; }
command -v curl >/dev/null || { echo "curl is required" >&2; exit 1; }

BASE="https://sonarcloud.io/api/issues/search"
AUTH=()
[[ -n "${SONAR_TOKEN:-}" ]] && AUTH=(-H "Authorization: Bearer ${SONAR_TOKEN}")

enc() { local val="$1"; python3 -c "import urllib.parse; print(urllib.parse.quote('''$val''', safe=''))"; }

QS="componentKeys=$(enc "$PROJECT_KEY")"
QS+="&issueStatuses=$(enc "$ISSUE_STATUSES")"
QS+="&ps=${PAGE_SIZE}"
QS+="&facets=types,severities,rules,languages"
[[ -n "$BRANCH" ]] && QS+="&branch=$(enc "$BRANCH")"
[[ -n "$TYPES" ]] && QS+="&types=$(enc "$TYPES")"
[[ -n "$SEVERITIES" ]] && QS+="&severities=$(enc "$SEVERITIES")"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

page=1
total=0
: >"$TMP/issues.jsonl"

while true; do
  curl -sS "${AUTH[@]}" "${BASE}?${QS}&p=${page}" -o "$TMP/page.json"
  if [[ $page -eq 1 ]]; then
    cp "$TMP/page.json" "$TMP/first.json"
    total=$(jq -r '.total // 0' "$TMP/first.json")
  fi
  n=$(jq -r '.issues | length' "$TMP/page.json")
  [[ "$n" -eq 0 ]] && break
  jq -c '.issues[]' "$TMP/page.json" >>"$TMP/issues.jsonl"
  fetched=$(wc -l <"$TMP/issues.jsonl" | tr -d ' ')
  [[ "$fetched" -ge "$total" ]] && break
  page=$((page + 1))
done

jq -s '.' "$TMP/issues.jsonl" >"$TMP/issues.json" 2>/dev/null || echo '[]' >"$TMP/issues.json"

if [[ "$FORMAT" == "json" ]]; then
  jq -n \
    --slurpfile first "$TMP/first.json" \
    --slurpfile issues "$TMP/issues.json" \
    --arg projectKey "$PROJECT_KEY" \
    --arg branch "$BRANCH" \
    --arg statuses "$ISSUE_STATUSES" \
    --arg fetchedAt "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    '{projectKey:$projectKey, branch:(if $branch=="" then null else $branch end), issueStatuses:$statuses, fetchedAt:$fetchedAt, total:$first[0].total, effortTotalMin:$first[0].effortTotal, debtTotalMin:$first[0].debtTotal, issueCount:($issues[0]|length), facets:$first[0].facets, issues:$issues[0]}' \
    >"$TMP/out.json"
  if [[ -n "$OUT_FILE" ]]; then
    mkdir -p "$(dirname "$OUT_FILE")"
    cp "$TMP/out.json" "$OUT_FILE"
    echo "Wrote $OUT_FILE ($(jq '.issueCount' "$OUT_FILE") issues)" >&2
  else
    jq '.' "$TMP/out.json"
  fi
  exit 0
fi

export TMP PROJECT_KEY BRANCH ISSUE_STATUSES OUT_FILE
python3 <<'PY'
import json, os
from collections import defaultdict
from datetime import datetime, timezone
from urllib.parse import quote

tmp = os.environ["TMP"]
with open(f"{tmp}/first.json") as f:
    first = json.load(f)
with open(f"{tmp}/issues.json") as f:
    issues = json.load(f)

project = os.environ["PROJECT_KEY"]
branch = os.environ.get("BRANCH", "")
statuses = os.environ["ISSUE_STATUSES"]
out_file = os.environ.get("OUT_FILE", "")
fetched_at = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")

ui = f"https://sonarcloud.io/project/issues?issueStatuses={quote(statuses)}&id={quote(project)}"
if branch:
    ui += f"&branch={quote(branch)}"

def rel(comp: str) -> str:
    return comp.split(":", 1)[1] if ":" in comp else comp

lines = [
    "# SonarCloud Issues Report", "",
    "| Field | Value |", "| --- | --- |",
    f"| Project | `{project}` |",
    f"| Branch | `{branch}` |" if branch else "| Branch | *(all)* |",
    f"| Statuses | `{statuses}` |",
    f"| Total (API) | {first.get('total', 0)} |",
    f"| Fetched | {len(issues)} |",
    f"| Effort / Debt (min) | {first.get('effortTotal', 0)} / {first.get('debtTotal', 0)} |",
    f"| Fetched at (UTC) | {fetched_at} |",
    f"| UI | {ui} |", "",
]
for facet in first.get("facets", []):
    lines += [f"## Summary: {facet['property']}", "", "| Value | Count |", "| --- | ---: |"]
    for v in facet.get("values", []):
        lines.append(f"| {v['val']} | {v['count']} |")
    lines.append("")

lines += ["## Issues by severity", ""]
for sev in ["BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"]:
    grp = [i for i in issues if i.get("severity") == sev]
    if not grp:
        continue
    lines.append(f"### {sev} ({len(grp)})")
    lines.append("")
    grp.sort(key=lambda i: (rel(i.get("component", "")), i.get("line") or 0))
    for i in grp:
        p, ln = rel(i.get("component", "")), i.get("line", "-")
        msg = (i.get("message") or "").replace("|", "/")
        lines.append(f"- **{p}:{ln}** | `{i.get('rule','')}` | {i.get('type','')} | {msg} | effort: {i.get('effort','')}")
    lines.append("")

lines += ["## Issues by file", ""]
by_file = defaultdict(list)
for i in issues:
    by_file[rel(i.get("component", ""))].append(i)
for path, grp in sorted(by_file.items(), key=lambda x: -len(x[1])):
    lines.append(f"### {path} ({len(grp)})")
    lines.append("")
    for i in sorted(grp, key=lambda x: (x.get("severity", ""), x.get("line") or 0)):
        lines.append(f"- L{i.get('line','-')} **{i.get('severity','')}** `{i.get('rule','')}` ({i.get('type','')}): {i.get('message','')}")
    lines.append("")

text = "\n".join(lines)
if out_file:
    os.makedirs(os.path.dirname(out_file) or ".", exist_ok=True)
    with open(out_file, "w", encoding="utf-8") as f:
        f.write(text)
    import sys
    print(f"Wrote {out_file} ({len(issues)} issues)", file=sys.stderr)
else:
    print(text)
PY
