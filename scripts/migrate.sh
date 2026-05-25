#!/usr/bin/env bash
# 按 docs 约定顺序执行 SQL 迁移（需本地 psql 与 PostgreSQL）
# 用法:
#   ./scripts/migrate.sh
#   DATABASE_URL=postgres://user:pass@localhost:5432/erp ./scripts/migrate.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

DB_URL="${DATABASE_URL:-postgres://erp:erp123@localhost:5433/erp?sslmode=disable}"

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
  file="backend/services/${svc}/migrations/001_init.sql"
  if [[ -f "$file" ]]; then
    run_sql "$file"
  fi
done

echo ""
echo "Migration complete."
