# ERP-Go

> [中文说明](README.md)

ERP-Go is a cross-border e-commerce ERP system built on a Go microservices architecture. It serves cross-border sellers, brands, and overseas warehouse teams, covering product, channel, order, inventory, WMS, TMS, procurement, finance, reporting, and system governance.

**One-line positioning:** An integrated microservices ERP platform for cross-border e-commerce sellers, brands, and overseas warehouse teams.

**Long-term vision:** Enable cross-border teams to manage products, platform sales, procurement, warehouse operations, logistics, financial accounting, and business analytics around order fulfillment and inventory/financial ledgers—instead of moving data between platform backends, warehouse systems, logistics systems, and spreadsheets.

Repository: [https://github.com/Tangyd893/ERP-Go.git](https://github.com/Tangyd893/ERP-Go.git)

## Core Capabilities

- **Multi-tenant governance** — Authentication, tenants and organizations, roles and permissions, audit logging
- **Product and channel** — SPU/SKU management, platform SKU mapping, store authorization, platform order import
- **Order fulfillment** — Sales order state machine, review, inventory locking, WMS outbound, shipping and platform sync
- **Warehouse operations** — Inbound, picking, checking, packing, outbound; admin console and PDA scanning terminal
- **Supply chain and finance** — Purchase orders, inbound receiving, AR/AP, cost allocation and profit reporting
- **Platform capabilities** — Unified API Gateway authentication, Outbox/Inbox event-driven architecture, observability stack

## Architecture Overview

```text
frontend
  ├─ ERP Admin (admin-web)
  ├─ WMS PDA / H5 (warehouse-pda)
  └─ Business Dashboard (dashboard-web)
        |
        v
API Gateway
        |
        +-- iam-service             Identity & access
        +-- tenant-service          Tenant & organization
        +-- product-service         Product catalog
        +-- channel-service         Sales channels
        +-- order-service           Orders
        +-- inventory-service       Inventory
        +-- warehouse-service       WMS
        +-- transport-service       TMS / logistics
        +-- purchase-service        Procurement
        +-- finance-service         Finance
        +-- report-service          Reporting
        +-- file-service            File storage
        +-- notification-service    Notifications

PostgreSQL | Redis | RabbitMQ | MinIO | OpenSearch
Prometheus | Grafana | Loki | Jaeger
```

**Design principles:**

- Business closed loops take priority over service count; each service owns its boundary and data
- Cross-service side effects prefer events; inventory and finance must maintain transaction ledgers and reconciliation
- External platforms, carriers, and warehouses are isolated behind adapters

See [Interface and Event Design](docs/接口与事件设计.md) for event-driven and cross-service collaboration details.

## Microservice Responsibilities

| Service | Responsibility |
|---------|----------------|
| API Gateway | Unified entry point, JWT auth, routing, trace context |
| IAM Service | Users, roles, permissions, login, audit |
| Tenant Service | Tenants, organizations, departments, positions |
| Product Service | SPU, SKU, platform SKU mapping |
| Channel Service | Store authorization, platform sync, order import |
| Order Service | Sales orders, state machine, review, exception handling |
| Inventory Service | Inventory balance, lock, release, deduct, ledger |
| Warehouse Service | Inbound, putaway, picking, checking, packing, outbound |
| Transport Service | Carriers, rules, labels, shipping, tracking |
| Purchase Service | Suppliers, POs, receiving, QC, inbound |
| Finance Service | AR/AP, settlement, cost, profit, ledger |
| Report Service | Sales, inventory, warehouse, profit reports |
| File Service | Product images, labels, invoices, import/export files |
| Notification Service | Email, SMS, in-app messages, webhooks, alerts |

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.25 + Gin + GORM |
| Gateway | Gin API Gateway, unified auth, reverse proxy |
| Database | PostgreSQL |
| Cache | Redis |
| Async events | RabbitMQ, Outbox/Inbox reliable messaging |
| Object storage | MinIO (S3-compatible in production) |
| Search / analytics | OpenSearch |
| Frontend | Vue 3 + TypeScript + Vite + Pinia + Element Plus |
| Warehouse terminal | Vue 3 H5 / PDA scanning |
| Monitoring | Prometheus, Grafana, Loki, Jaeger |
| Deployment | Docker Compose (local); Kubernetes / Helm (production) |

## Repository Structure

```text
ERP-Go/
  backend/    Microservices, gateway, shared libs, migrations
  frontend/   Admin console, WMS PDA, dashboard, shared components
  docker/     Docker Compose, Kubernetes, monitoring, deployment
  testing/    Integration, contract, E2E, performance tests
  docs/       Architecture, domain model, API/events, data model, engineering standards
  scripts/    Local startup, migration, verification scripts
```

## Quick Start

### Requirements

| Tool | Version | Purpose |
|------|---------|---------|
| Git | 2.40+ | Version control |
| Docker | 20.10+ | Local infrastructure |
| Go | 1.25.0 | Backend development |
| Node.js | 18.x+ | Frontend development |
| npm | 9.x+ | Frontend package management |

### Clone

```bash
git clone https://github.com/Tangyd893/ERP-Go.git
cd ERP-Go
```

### Local Dev Stack

**Windows PowerShell:**

```powershell
.\scripts\dev-stack.ps1 all       # Infra + migrations + core services
.\scripts\dev-stack.ps1 infra     # Docker infra + migrations only
.\scripts\dev-stack.ps1 services  # Core microservices only
```

**Linux / macOS / Git Bash:**

```bash
make dev-stack    # Infra + migrations + core services
make dev-infra    # Infra + migrations only
make dev-services # Core microservices only
```

`dev-stack` starts Gateway, IAM, Order, Inventory, and Warehouse by default. Service logs are written to `.cache/logs/`.

**Health check:**

```bash
curl http://localhost:8080/health
```

### Default Ports

| Service | Port |
|---------|-----:|
| API Gateway | 8080 |
| IAM Service | 8081 |
| Tenant Service | 8082 |
| Product Service | 8083 |
| Channel Service | 8084 |
| Order Service | 8085 |
| Inventory Service | 8086 |
| Warehouse Service | 8087 |
| Transport Service | 8088 |
| File Service | 8089 |
| Purchase Service | 8091 |
| Finance Service | 8092 |
| Report Service | 8093 |
| Notification Service | 8094 |

### Frontend Development

```bash
npm install
npm run dev:admin        # Admin console → http://localhost:5173
npm run dev:pda          # WMS PDA      → http://localhost:5174
npm run dev:dashboard    # Dashboard    → http://localhost:5175
```

## Development and Verification

Use the unified verification scripts; they pin `GOCACHE` / `GOMODCACHE` to the repo `.cache/` directory:

```bash
# Linux / macOS / Git Bash / CI
./scripts/verify.sh

# Windows PowerShell
.\scripts\verify.ps1

# Skip frontend (Go only)
./scripts/verify.sh --skip-frontend
.\scripts\verify.ps1 -SkipFrontend
```

**Makefile shortcuts:**

```bash
make verify    # Go vet/test/build + frontend typecheck/lint/build
make test      # Go unit tests
make lint      # go vet + eslint
make build     # Build all Go services
make all       # Build + test
make help      # List all commands
```

> Running `go test ./...` directly may fail due to system Go cache permissions; prefer the scripts above or `make test`.

### Manual Startup (Optional)

```bash
# Start dev infrastructure
docker compose -f docker/compose/docker-compose.dev.yml up -d

# Start IAM service (run migrations first)
DATABASE_PORT=5433 SERVER_PORT=8081 JWT_SECRET=dev-secret go run ./backend/services/iam-service/cmd/server/
```

See [.env.example](.env.example) for environment variable templates.

## Documentation

Full index: [docs/文档索引.md](docs/文档索引.md)

**Architecture and design:**

- [Project Architecture](docs/项目架构设计.md)
- [Domain Model Design](docs/领域模型设计.md)
- [Microservice Design](docs/微服务设计说明.md)
- [Interface and Event Design](docs/接口与事件设计.md)
- [Data Model Design](docs/数据模型设计.md)

**Engineering standards:**

- [Tech Stack and Middleware Standards](docs/技术栈与中间件管理规范.md)
- [Implementation Roadmap and Engineering Standards](docs/实施路线与工程规范.md)

**Operations:**

- [docs/runbooks/](docs/runbooks/) — Deployment, E2E smoke tests, backup and recovery

**API contract:**

- [docs/api/openapi.yaml](docs/api/openapi.yaml)

## Production Deployment Notes

- Inject sensitive configuration via environment variables, Kubernetes Secrets, or a dedicated secrets manager
- With `ENVIRONMENT=production`, Gateway/IAM reject default `JWT_SECRET` and default database passwords (see `.env.example`); services exit if the database is unavailable outside `development`
- Configure `CORS_ALLOWED_ORIGINS` as an allowlist in production; do not use `*`
- Optionally enable `GATEWAY_RBAC_ENABLED=true` for coarse-grained RBAC; block direct service access via NetworkPolicy
- Do not commit `.env`, database passwords, platform tokens, carrier keys, payment keys, or real customer data
- External system tokens must be encrypted at rest with authorization audit trails
- Multi-tenant data access must enforce `tenant_id` validation
- Inventory adjustments, financial reversals, permission changes, and data exports must be audit-logged
- All services must integrate structured logging, metrics, and distributed tracing
