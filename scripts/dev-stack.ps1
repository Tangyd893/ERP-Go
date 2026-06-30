# ERP-Go ??????PowerShell?
# ??:
#   .\scripts\dev-stack.ps1 infra          # ??? Docker ??? + ??
#   .\scripts\dev-stack.ps1 services       # ????????? infra ????
#   .\scripts\dev-stack.ps1 all            # infra + services

param(
    [ValidateSet("infra", "services", "all")]
    [string]$Target = "all"
)

$ErrorActionPreference = "Stop"
$RepoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $RepoRoot

$env:DATABASE_PORT = "5433"
$env:DATABASE_HOST = "localhost"
$env:DATABASE_USER = "erp"
$env:DATABASE_PASSWORD = "erp123"
$env:DATABASE_DBNAME = "erp_go"
$env:JWT_SECRET = "erp-go-dev-secret-change-in-production"
$env:RABBITMQ_URL = "amqp://admin:admin123@localhost:5672/"
$env:ORDER_SERVICE_URL = "http://localhost:8085"
$env:INVENTORY_SERVICE_URL = "http://localhost:8086"
$env:WAREHOUSE_SERVICE_URL = "http://localhost:8087"

function Start-Infra {
    Write-Output "==> ?? Docker ???"
    $env:POSTGRES_PORT = "5433"
    Push-Location (Join-Path $RepoRoot "docker/compose")
    docker compose -f docker-compose.dev.yml up -d postgres rabbitmq redis
    Pop-Location

    Write-Output "==> ?? PostgreSQL ??..."
    $ready = $false
    for ($i = 0; $i -lt 30; $i++) {
        docker exec erp-postgres pg_isready -U erp -d erp_go 2>$null
        if ($LASTEXITCODE -eq 0) { $ready = $true; break }
        Start-Sleep -Seconds 2
    }
    if (-not $ready) {
        throw "PostgreSQL ?????????"
    }

    Write-Output "==> ???????"
    $env:DATABASE_URL = "postgres://$env:DATABASE_USER`:$env:DATABASE_PASSWORD@$env:DATABASE_HOST`:$env:DATABASE_PORT/$env:DATABASE_DBNAME`?sslmode=disable"
    & (Join-Path $RepoRoot "scripts/migrate.ps1")
}

function Start-Services {
    $goCache = Join-Path $RepoRoot ".cache/go-build"
    $goModCache = Join-Path $RepoRoot ".cache/go-mod"
    New-Item -ItemType Directory -Force -Path $goCache, $goModCache | Out-Null

    . (Join-Path $PSScriptRoot "lib/start-go-service.ps1")

    $baseEnv = @{
        GOCACHE           = $goCache
        GOMODCACHE        = $goModCache
        DATABASE_HOST     = $env:DATABASE_HOST
        DATABASE_PORT     = $env:DATABASE_PORT
        DATABASE_USER     = $env:DATABASE_USER
        DATABASE_PASSWORD = $env:DATABASE_PASSWORD
        DATABASE_DBNAME   = $env:DATABASE_DBNAME
        JWT_SECRET        = $env:JWT_SECRET
        RABBITMQ_URL      = $env:RABBITMQ_URL
    }

    $services = @(
        @{ Name = "iam"; Port = 8081; Path = "./backend/services/iam-service/cmd/server" },
        @{ Name = "inventory"; Port = 8086; Path = "./backend/services/inventory-service/cmd/server" },
        @{ Name = "warehouse"; Port = 8087; Path = "./backend/services/warehouse-service/cmd/server"; Extra = @{
            ORDER_SERVICE_URL = $env:ORDER_SERVICE_URL
        }},
        @{ Name = "order"; Port = 8085; Path = "./backend/services/order-service/cmd/server"; Extra = @{
            INVENTORY_SERVICE_URL = $env:INVENTORY_SERVICE_URL
            WAREHOUSE_SERVICE_URL = $env:WAREHOUSE_SERVICE_URL
        }},
        @{ Name = "gateway"; Port = 8080; Path = "./backend/gateway/cmd/server" }
    )

    foreach ($svc in $services) {
        Write-Output "==> ?? $($svc.Name) (:$($svc.Port))"
        $svcEnv = Merge-ServiceEnv $baseEnv @{ SERVER_PORT = "$($svc.Port)" }
        if ($svc.Extra) {
            $svcEnv = Merge-ServiceEnv $svcEnv $svc.Extra
        }
        Start-GoServiceSilent -Name $svc.Name -RepoRoot $RepoRoot -ServicePath $svc.Path -Environment $svcEnv
        Start-Sleep -Seconds 1
    }

    Write-Output ""
    Write-Output "?????????????????? .cache/logs/?"
    Write-Output "Gateway:    http://localhost:8080/health"
    Write-Output "IAM ??:   POST http://localhost:8080/api/v1/iam/login  (admin/admin123, tenant=default)"
    Write-Output "PDA ??:   npm run dev:pda  (?? 5174??? /api -> :8080)"
}

switch ($Target) {
    "infra" { Start-Infra }
    "services" { Start-Services }
    "all" { Start-Infra; Start-Services }
}
