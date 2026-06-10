.PHONY: help dev build build-all test test-cover lint verify clean frontend-dev frontend-build frontend-typecheck all

CACHE_DIR := $(shell pwd)/.cache
GO_CACHE := $(CACHE_DIR)/go-build
GO_MOD_CACHE := $(CACHE_DIR)/go-mod

help: ## 显示帮助信息
	@echo "开发规范:"
	@echo "  make dev                  启动所有开发中间件"
	@echo "  make build                编译所有 Go 服务"
	@echo "  make test                 运行所有单元测试"
	@echo "  make lint                 代码检查 (go vet)"
	@echo "  make verify               统一验证 (Go vet/test/build + 前端 typecheck/build)"
	@echo "  make frontend-build       编译所有前端应用"
	@echo "  make frontend-typecheck   前端类型检查"
	@echo "  make frontend-dev         启动 admin-web 开发服务"
	@echo "  make build-all            编译后端+前端"
	@echo "  make test-cover           运行测试并生成覆盖率报告"
	@echo "  make clean                清理构建产物"
	@echo "  make all                  编译并测试"
	@echo "  make run-gateway          启动 API 网关 (port 8080)"
	@echo "  make run-iam              启动 IAM 服务 (port 8081)"
	@echo "  make run-order            启动 Order 服务 (port 8085)"
	@echo "  make run-inventory        启动 Inventory 服务 (port 8086)"
	@echo "  make run-warehouse        启动 Warehouse 服务 (port 8087)"

dev: ## 启动所有开发中间件 (PostgreSQL/Redis/RabbitMQ/MinIO等)
	cd docker/compose && docker compose -f docker-compose.dev.yml up -d

docker-down: ## 停止所有开发中间件
	cd docker/compose && docker compose -f docker-compose.dev.yml down

build: ## 编译所有 Go 服务
	mkdir -p $(GO_CACHE) $(GO_MOD_CACHE) && GOCACHE=$(GO_CACHE) GOMODCACHE=$(GO_MOD_CACHE) go build -C backend ./...

build-all: build frontend-build ## 编译后端+前端

test: ## 运行所有单元测试
	mkdir -p $(GO_CACHE) $(GO_MOD_CACHE) && GOCACHE=$(GO_CACHE) GOMODCACHE=$(GO_MOD_CACHE) go test -C backend ./... -v -count=1

test-cover: ## 运行测试并生成覆盖率报告
	mkdir -p $(GO_CACHE) $(GO_MOD_CACHE) && GOCACHE=$(GO_CACHE) GOMODCACHE=$(GO_MOD_CACHE) go test -C backend ./... -coverprofile=../coverage.out -covermode=atomic
	GOCACHE=$(GO_CACHE) go tool cover -html=coverage.out -o coverage.html

lint: ## 代码检查 (go vet)
	mkdir -p $(GO_CACHE) $(GO_MOD_CACHE) && GOCACHE=$(GO_CACHE) GOMODCACHE=$(GO_MOD_CACHE) go vet -C backend ./...

verify: ## 统一验证 (Go + 前端)
	bash scripts/verify.sh

migrate: ## 执行数据库迁移 (需 psql)
	bash scripts/migrate.sh

dev-infra: ## 启动 Docker 中间件并迁移
	bash scripts/dev-stack.sh infra

dev-services: ## 启动核心微服务 (gateway/iam/order/inventory/warehouse)
	bash scripts/dev-stack.sh services

dev-stack: ## 启动中间件 + 迁移 + 核心服务
	bash scripts/dev-stack.sh all

clean: ## 清理构建产物
	rm -f coverage.out coverage.html

frontend-dev: ## 启动 admin-web 开发服务
	npm run dev:admin

frontend-dev-pda: ## 启动 warehouse-pda 开发服务
	npm run dev:pda

frontend-dev-dashboard: ## 启动 dashboard-web 开发服务
	npm run dev:dashboard

frontend-build: ## 编译所有前端应用
	npm run build:admin && npm run build:pda && npm run build:dashboard

frontend-typecheck: ## 前端类型检查
	npm run typecheck

# ---- 服务启动 (开发模式) ----
run-gateway: ## 启动 API 网关 (port 8080)
	go run ./backend/gateway/cmd/server/

run-iam: ## 启动 IAM 服务 (port 8081)
	go run ./backend/services/iam-service/cmd/server/

run-tenant: ## 启动 Tenant 服务 (port 8082)
	go run ./backend/services/tenant-service/cmd/server/

run-product: ## 启动 Product 服务 (port 8083)
	go run ./backend/services/product-service/cmd/server/

run-channel: ## 启动 Channel 服务 (port 8084)
	go run ./backend/services/channel-service/cmd/server/

run-order: ## 启动 Order 服务 (port 8085)
	go run ./backend/services/order-service/cmd/server/

run-inventory: ## 启动 Inventory 服务 (port 8086)
	go run ./backend/services/inventory-service/cmd/server/

run-warehouse: ## 启动 Warehouse 服务 (port 8087)
	go run ./backend/services/warehouse-service/cmd/server/

run-transport: ## 启动 Transport 服务 (port 8088)
	go run ./backend/services/transport-service/cmd/server/

run-file: ## 启动 File 服务 (port 8089)
	go run ./backend/services/file-service/cmd/server/

run-purchase: ## 启动 Purchase 服务 (port 8091)
	go run ./backend/services/purchase-service/cmd/server/

run-finance: ## 启动 Finance 服务 (port 8092)
	go run ./backend/services/finance-service/cmd/server/

run-report: ## 启动 Report 服务 (port 8093)
	go run ./backend/services/report-service/cmd/server/

run-notification: ## 启动 Notification 服务 (port 8094)
	go run ./backend/services/notification-service/cmd/server/

all: build test ## 编译并测试
