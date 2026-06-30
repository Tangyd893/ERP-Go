# ERP-Go 统一验证入口（Windows PowerShell）
# 用法:
#   .\scripts\verify.ps1                 # Go + 前端全量
#   .\scripts\verify.ps1 -SkipFrontend   # 仅 Go（npm 不可用时）
#   .\scripts\verify.ps1 -Verbose

param(
    [switch]$SkipFrontend,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"

$RepoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $RepoRoot

$CacheDir = Join-Path $RepoRoot ".cache"
$GoCache = Join-Path $CacheDir "go-build"
$GoModCache = Join-Path $CacheDir "go-mod"

New-Item -ItemType Directory -Force -Path $GoCache, $GoModCache | Out-Null

$env:GOCACHE = $GoCache
$env:GOMODCACHE = $GoModCache

function Write-Step([string]$Message) {
    Write-Output ""
    Write-Output "==> $Message"
}

function Invoke-Check([string]$Name, [scriptblock]$Block) {
    Write-Step $Name
    & $Block
    if ($LASTEXITCODE -ne 0) {
        throw "$Name failed (exit $LASTEXITCODE)"
    }
}

Write-Output "ERP-Go verify"
Write-Output "Repo: $RepoRoot"
Write-Output "GOCACHE=$GoCache"
Write-Output "GOMODCACHE=$GoModCache"

try {
    Invoke-Check "go vet ./..." {
        if ($Verbose) {
            go vet -C backend ./...
        } else {
            go vet -C backend ./...
        }
    }

    Invoke-Check "go test ./..." {
        go test -C backend ./... -count=1
    }

    Invoke-Check "go build ./..." {
        go build -C backend ./...
    }

    if (-not $SkipFrontend) {
        $npm = Get-Command npm -ErrorAction SilentlyContinue
        if (-not $npm) {
            Write-Warning "npm not found; skipping frontend checks. Use -SkipFrontend to suppress this warning."
        } else {
            if (-not (Test-Path (Join-Path $RepoRoot "node_modules"))) {
                Invoke-Check "npm install" {
                    if (Test-Path (Join-Path $RepoRoot "package-lock.json")) {
                        npm ci
                    } else {
                        npm install
                    }
                }
            }

            Invoke-Check "npm run typecheck" { npm run typecheck }
            Invoke-Check "npm run lint:ci" { npm run lint:ci }
            Invoke-Check "npm run build:admin" { npm run build:admin }
            Invoke-Check "npm run build:pda" { npm run build:pda }
            Invoke-Check "npm run build:dashboard" { npm run build:dashboard }
        }
    } else {
        Write-Output ""
        Write-Output "Skipping frontend checks (-SkipFrontend)."
    }

    Write-Output ""
    Write-Output "All checks passed."
    exit 0
} catch {
    Write-Output ""
    Write-Output "Verify failed: $_"
    exit 1
}
