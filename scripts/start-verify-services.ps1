# Start ERP core services on 908x (avoid WorkPal 808x port conflict)
$ScriptDir = Split-Path -Parent $PSCommandPath
$RepoRoot = (Resolve-Path (Join-Path $ScriptDir "..")).Path
Set-Location $RepoRoot

. (Join-Path $ScriptDir "lib/start-go-service.ps1")

$baseEnv = @{
    GOCACHE            = Join-Path $RepoRoot ".cache/go-build"
    GOMODCACHE         = Join-Path $RepoRoot ".cache/go-mod"
    DATABASE_HOST      = "localhost"
    DATABASE_PORT      = "5433"
    DATABASE_USER      = "erp"
    DATABASE_PASSWORD  = "erp123"
    DATABASE_DBNAME    = "erp_go"
    JWT_SECRET         = "erp-go-dev-secret-change-in-production"
    RABBITMQ_URL       = "amqp://admin:admin123@localhost:5672/"
}

New-Item -ItemType Directory -Force -Path $baseEnv.GOCACHE, $baseEnv.GOMODCACHE | Out-Null

Start-GoServiceSilent -Name "iam" -RepoRoot $RepoRoot -ServicePath "./backend/services/iam-service/cmd/server/" -Environment (Merge-ServiceEnv $baseEnv @{ SERVER_PORT = "9081" })
Start-Sleep -Seconds 1
Start-GoServiceSilent -Name "inventory" -RepoRoot $RepoRoot -ServicePath "./backend/services/inventory-service/cmd/server/" -Environment (Merge-ServiceEnv $baseEnv @{ SERVER_PORT = "9086" })
Start-Sleep -Seconds 1
Start-GoServiceSilent -Name "warehouse" -RepoRoot $RepoRoot -ServicePath "./backend/services/warehouse-service/cmd/server/" -Environment (Merge-ServiceEnv $baseEnv @{
    SERVER_PORT       = "9087"
    ORDER_SERVICE_URL = "http://localhost:9085"
})
Start-Sleep -Seconds 1
Start-GoServiceSilent -Name "order" -RepoRoot $RepoRoot -ServicePath "./backend/services/order-service/cmd/server/" -Environment (Merge-ServiceEnv $baseEnv @{
    SERVER_PORT           = "9085"
    INVENTORY_SERVICE_URL = "http://localhost:9086"
    WAREHOUSE_SERVICE_URL = "http://localhost:9087"
})
Start-Sleep -Seconds 1
Start-GoServiceSilent -Name "gateway" -RepoRoot $RepoRoot -ServicePath "./backend/gateway/cmd/server/" -Environment (Merge-ServiceEnv $baseEnv @{
    SERVER_PORT                         = "9080"
    SERVICE_TARGET_API_V1_IAM_          = "http://localhost:9081"
    SERVICE_TARGET_API_V1_ORDER_        = "http://localhost:9085"
    SERVICE_TARGET_API_V1_INVENTORY_    = "http://localhost:9086"
    SERVICE_TARGET_API_V1_WAREHOUSE_    = "http://localhost:9087"
})

Write-Host "ERP services started on port 9080 (gateway), logs in .cache/logs/"
