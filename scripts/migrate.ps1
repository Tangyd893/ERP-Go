# 按 docs 约定顺序执行 SQL 迁移（需 psql 与 PostgreSQL）
# 用法:
#   .\scripts\migrate.ps1
#   $env:DATABASE_URL="postgres://erp:erp123@localhost:5433/erp?sslmode=disable"; .\scripts\migrate.ps1

param(
    [string]$DatabaseUrl = $env:DATABASE_URL
)

$ErrorActionPreference = "Stop"

if (-not $DatabaseUrl) {
    $DatabaseUrl = "postgres://erp:erp123@localhost:5433/erp?sslmode=disable"
}

$RepoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $RepoRoot

$psql = Get-Command psql -ErrorAction SilentlyContinue
if (-not $psql) {
    Write-Error "psql not found. Install PostgreSQL client or run migrate.sh in Git Bash."
}

function Invoke-Migration([string]$RelativePath) {
    $fullPath = Join-Path $RepoRoot $RelativePath
    if (-not (Test-Path $fullPath)) {
        Write-Warning "Skip missing: $RelativePath"
        return
    }
    Write-Host "==> $RelativePath" -ForegroundColor Cyan
    & psql $DatabaseUrl -v ON_ERROR_STOP=1 -f $fullPath
    if ($LASTEXITCODE -ne 0) {
        throw "Migration failed: $RelativePath"
    }
}

Write-Host "Migrating ERP-Go database" -ForegroundColor Green
Write-Host "DATABASE_URL=$DatabaseUrl"

Invoke-Migration "backend/migrations/outbox/001_create_outbox.sql"
Invoke-Migration "backend/migrations/outbox/002_add_tenant_id.sql"

$services = @(
    "iam-service", "tenant-service", "product-service", "channel-service", "order-service",
    "inventory-service", "warehouse-service", "transport-service", "purchase-service",
    "finance-service", "file-service", "notification-service"
)

foreach ($svc in $services) {
    Invoke-Migration "backend/services/$svc/migrations/001_init.sql"
}

Write-Host ""
Write-Host "Migration complete." -ForegroundColor Green
