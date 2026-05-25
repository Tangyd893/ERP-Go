# 按 docs 约定顺序执行 SQL 迁移（需 psql 与 PostgreSQL）
# 用法:
#   .\scripts\migrate.ps1
#   $env:DATABASE_URL="postgres://user:password@host:port/db?sslmode=disable"; .\scripts\migrate.ps1

param(
    [string]$DatabaseUrl = $env:DATABASE_URL
)

$ErrorActionPreference = "Stop"

if (-not $DatabaseUrl) {
    $host = if ($env:DATABASE_HOST) { $env:DATABASE_HOST } else { "localhost" }
    $port = if ($env:DATABASE_PORT) { $env:DATABASE_PORT } else { "5432" }
    $user = $env:DATABASE_USER
    $pass = $env:DATABASE_PASSWORD
    $dbname = if ($env:DATABASE_DBNAME) { $env:DATABASE_DBNAME } else { $env:DATABASE_NAME }
    if ($user -and $pass) {
        $DatabaseUrl = "postgres://${user}:${pass}@${host}:${port}/${dbname}?sslmode=disable"
    } else {
        Write-Error "未设置 DATABASE_URL 且 DATABASE_USER/DATABASE_PASSWORD 环境变量为空。请复制 .env.example 并配置 DATABASE_URL，或设置 DATABASE_* 环境变量。"
        exit 1
    }
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
    $dir = Join-Path $RepoRoot "backend/services/$svc/migrations"
    if (Test-Path $dir) {
        Get-ChildItem -Path $dir -Filter "*.sql" | Sort-Object Name | ForEach-Object {
            Invoke-Migration "backend/services/$svc/migrations/$($_.Name)"
        }
    }
}

Write-Host ""
Write-Host "Migration complete." -ForegroundColor Green
