#!/usr/bin/env bash
# ERP-Go 本地开发栈
# 用法: ./scripts/dev-stack.sh [infra|services|all]

set -euo pipefail

TARGET="${1:-all}"
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

export DATABASE_PORT=5433 DATABASE_HOST=localhost DATABASE_USER=erp
export DATABASE_PASSWORD=erp123 DATABASE_DBNAME=erp_go
export JWT_SECRET=erp-go-dev-secret-change-in-production
export RABBITMQ_URL=amqp://admin:admin123@localhost:5672/
export ORDER_SERVICE_URL=http://localhost:8085
export INVENTORY_SERVICE_URL=http://localhost:8086
export WAREHOUSE_SERVICE_URL=http://localhost:8087

start_infra() {
  echo "==> 启动 Docker 中间件"
  (cd docker/compose && docker compose -f docker-compose.dev.yml up -d postgres rabbitmq redis)
  echo "==> 等待 PostgreSQL"
  for i in $(seq 1 30); do
    if docker exec erp-postgres pg_isready -U erp -d erp_go >/dev/null 2>&1; then
      break
    fi
    sleep 2
  done
  bash scripts/migrate.sh
}

start_services() {
  mkdir -p .cache/go-build .cache/go-mod
  export GOCACHE="$REPO_ROOT/.cache/go-build"
  export GOMODCACHE="$REPO_ROOT/.cache/go-mod"

  SERVER_PORT=8080 go run ./backend/gateway/cmd/server/ &
  SERVER_PORT=8081 go run ./backend/services/iam-service/cmd/server/ &
  SERVER_PORT=8085 go run ./backend/services/order-service/cmd/server/ &
  SERVER_PORT=8086 go run ./backend/services/inventory-service/cmd/server/ &
  SERVER_PORT=8087 go run ./backend/services/warehouse-service/cmd/server/ &

  echo "核心服务已在后台启动 (gateway/iam/order/inventory/warehouse)"
}

case "$TARGET" in
  infra) start_infra ;;
  services) start_services ;;
  all) start_infra; start_services ;;
  *) echo "用法: $0 [infra|services|all]"; exit 1 ;;
esac
