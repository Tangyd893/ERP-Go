#!/usr/bin/env bash
# 按 docs 约定顺序执行 SQL 迁移（需本地 psql 与 PostgreSQL）
# 迁移顺序：outbox 基础设施 → 各服务 migrations（字母序服务名，文件 sort）
# 服务列表：iam-service, channel-service, file-service, finance-service,
#          inventory-service, notification-service, order-service,
#          product-service, purchase-service, tenant-service,
#          transport-service, warehouse-service
# report-service 为无状态聚合服务，不创建业务表。
# 用法:
#   ./scripts/migrate.sh
#   DATABASE_URL=postgres://user:pass@localhost:5432/erp ./scripts/migrate.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

if [[ -z "${DATABASE_URL:-}" ]]; then
  DB_HOST="${DATABASE_HOST:-localhost}"
  DB_PORT="${DATABASE_PORT:-5432}"
  DB_USER="${DATABASE_USER:-}"
  DB_PASS="${DATABASE_PASSWORD:-}"
  DB_NAME="${DATABASE_DBNAME:-${DATABASE_NAME:-}}"
  if [[ -n "$DB_USER" && -n "$DB_PASS" ]]; then
    DB_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
  else
    echo "ERROR: 请设置 DATABASE_URL 或 DATABASE_USER/DATABASE_PASSWORD 环境变量。" >&2
    echo "复制 .env.example 并配置数据库连接信息。" >&2
    exit 1
  fi
else
  DB_URL="$DATABASE_URL"
fi

run_sql() {
  local file="$1"
  echo "==> $file"
  psql "$DB_URL" -v ON_ERROR_STOP=1 -f "$file"
}

echo "Migrating ERP-Go database"
echo "DATABASE_URL=$DB_URL"

run_sql "backend/migrations/outbox/001_create_outbox.sql"
run_sql "backend/migrations/outbox/002_add_tenant_id.sql"

for svc in iam-service tenant-service product-service channel-service order-service \
  inventory-service warehouse-service transport-service purchase-service finance-service \
  file-service notification-service; do
  dir="backend/services/${svc}/migrations"
  if [[ -d "$dir" ]]; then
    while IFS= read -r file; do
      run_sql "$file"
    done < <(find "$dir" -maxdepth 1 -name '*.sql' | sort)
  fi
done

echo ""
echo "Migration complete."
