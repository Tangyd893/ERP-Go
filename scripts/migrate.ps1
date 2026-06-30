# Run SQL migrations against PostgreSQL (requires psql client or Docker)
# Usage:
#   .\scripts\migrate.ps1
#   $env:DATABASE_URL="postgres://user:password@host:port/db?sslmode=disable"; .\scripts\migrate.ps1

param(
    [string]$DatabaseUrl = $env:DATABASE_URL
)

$ErrorActionPreference = "Stop"

if (-not $DatabaseUrl) {
    $dbHost = if ($env:DATABASE_HOST) { $env:DATABASE_HOST } else { "localhost" }
    $port = if ($env:DATABASE_PORT) { $env:DATABASE_PORT } else { "5432" }
    $user = $env:DATABASE_USER
    $pass = $env:DATABASE_PASSWORD
    $dbname = if ($env:DATABASE_DBNAME) { $env:DATABASE_DBNAME } else { $env:DATABASE_NAME }
    if ($user -and $pass) {
        $DatabaseUrl = "postgres://${user}:${pass}@${dbHost}:${port}/${dbname}?sslmode=disable"
    } else {
        Write-Error "DATABASE_URL not set and DATABASE_USER/DATABASE_PASSWORD missing. Copy .env.example and configure."
        exit 1
    }
}

$RepoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $RepoRoot

function Invoke-MigrationDocker([string]$RelativePath) {
    $fullPath = Join-Path $RepoRoot $RelativePath
    if (-not (Test-Path $fullPath)) {
        Write-Output "SKIP missing: $RelativePath"
        return
    }
    Write-Output "==> $RelativePath"
    Get-Content -Path $fullPath -Raw -Encoding UTF8 | docker exec -i erp-postgres psql -U erp -d erp_go -v ON_ERROR_STOP=1
    if ($LASTEXITCODE -ne 0) {
        throw "Migration failed: $RelativePath"
    }
}

Write-Output "Migrating ERP-Go database"
Write-Output "DATABASE_URL=$DatabaseUrl"

Invoke-MigrationDocker "backend/migrations/outbox/001_create_outbox.sql"
Invoke-MigrationDocker "backend/migrations/outbox/002_add_tenant_id.sql"

$services = @(
    "iam-service", "tenant-service", "product-service", "channel-service", "order-service",
    "inventory-service", "warehouse-service", "transport-service", "purchase-service",
    "finance-service", "file-service", "notification-service"
)

foreach ($svc in $services) {
    $dir = Join-Path $RepoRoot "backend/services/$svc/migrations"
    if (Test-Path $dir) {
        Get-ChildItem -Path $dir -Filter "*.sql" | Sort-Object Name | ForEach-Object {
            Invoke-MigrationDocker "backend/services/$svc/migrations/$($_.Name)"
        }
    }
}

Write-Output ""
Write-Output "Migration complete."
