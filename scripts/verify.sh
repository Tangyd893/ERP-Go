#!/usr/bin/env bash
# ERP-Go 统一验证入口（Linux / macOS / Git Bash）
# 用法:
#   ./scripts/verify.sh                 # Go + 前端全量
#   ./scripts/verify.sh --skip-frontend # 仅 Go

set -euo pipefail

SKIP_FRONTEND=0
ARG1="${1:-}"
if [[ "$ARG1" == "--skip-frontend" ]]; then
  SKIP_FRONTEND=1
fi

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

CACHE_DIR="$REPO_ROOT/.cache"
export GOCACHE="$CACHE_DIR/go-build"
export GOMODCACHE="$CACHE_DIR/go-mod"
mkdir -p "$GOCACHE" "$GOMODCACHE"

step() {
  echo ""
  echo "==> $1"
}

echo "ERP-Go verify"
echo "Repo: $REPO_ROOT"
echo "GOCACHE=$GOCACHE"
echo "GOMODCACHE=$GOMODCACHE"

step "go vet ./..."
go vet -C backend ./...

step "go test ./..."
go test -C backend ./... -count=1

step "go build ./..."
go build -C backend ./...

if [[ "$SKIP_FRONTEND" -eq 0 ]]; then
  if ! command -v npm >/dev/null 2>&1; then
    echo ""
    echo "WARNING: npm not found; skipping frontend checks."
  else
    if [[ ! -d node_modules ]]; then
      step "npm install"
      if [[ -f package-lock.json ]]; then
        npm ci
      else
        npm install
      fi
    fi

    step "npm run typecheck"
    npm run typecheck

    step "npm run build:admin"
    npm run build:admin

    step "npm run build:pda"
    npm run build:pda

    step "npm run build:dashboard"
    npm run build:dashboard
  fi
else
  echo ""
  echo "Skipping frontend checks (--skip-frontend)."
fi

echo ""
echo "All checks passed."
