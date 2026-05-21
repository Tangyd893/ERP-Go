.PHONY: help dev build test clean docker-up docker-down lint

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## 启动所有开发中间件
	cd docker/compose && docker compose -f docker-compose.dev.yml up -d

docker-down: ## 停止所有开发中间件
	cd docker/compose && docker compose -f docker-compose.dev.yml down

build: ## 编译所有 Go 服务
	go build ./...

test: ## 运行所有单元测试
	go test ./backend/... -v -count=1

test-cover: ## 运行测试并生成覆盖率报告
	go test ./backend/... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

lint: ## 代码检查
	go vet ./...

clean: ## 清理构建产物
	rm -f coverage.out coverage.html

run-gateway: ## 启动 API 网关
	go run ./backend/gateway/cmd/server/

run-iam: ## 启动 IAM 服务
	go run ./backend/services/iam-service/cmd/server/

run-tenant: ## 启动 Tenant 服务
	go run ./backend/services/tenant-service/cmd/server/

frontend-dev: ## 启动前端开发服务
	npm run dev:admin

frontend-build: ## 编译前端
	npm run build:admin

all: build test ## 编译并测试
