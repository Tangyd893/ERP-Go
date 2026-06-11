# Silent background go run (logs under .cache/logs/)
function Merge-ServiceEnv {
    param(
        [hashtable]$Base,
        [hashtable]$Extra
    )
    $merged = @{}
    foreach ($item in $Base.GetEnumerator()) { $merged[$item.Key] = $item.Value }
    foreach ($item in $Extra.GetEnumerator()) { $merged[$item.Key] = $item.Value }
    return $merged
}

function Start-GoServiceSilent {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$Name,
        [Parameter(Mandatory = $true)][string]$RepoRoot,
        [Parameter(Mandatory = $true)][string]$ServicePath,
        [hashtable]$Environment = @{}
    )

    $logDir = Join-Path $RepoRoot ".cache/logs"
    New-Item -ItemType Directory -Force -Path $logDir | Out-Null

    $outLog = Join-Path $logDir "$Name.log"
    $errLog = Join-Path $logDir "$Name.err.log"

    $commonArgs = @{
        FilePath               = "go"
        ArgumentList           = @("run", $ServicePath)
        WorkingDirectory       = $RepoRoot
        WindowStyle            = "Hidden"
        RedirectStandardOutput = $outLog
        RedirectStandardError  = $errLog
        PassThru               = $true
    }

    if ($PSVersionTable.PSVersion.Major -ge 7 -and $Environment.Count -gt 0) {
        $proc = Start-Process @commonArgs -Environment $Environment
    } else {
        $saved = @{}
        foreach ($item in $Environment.GetEnumerator()) {
            $saved[$item.Key] = [Environment]::GetEnvironmentVariable($item.Key, "Process")
            Set-Item -Path "env:$($item.Key)" -Value ([string]$item.Value)
        }
        try {
            $proc = Start-Process @commonArgs
        } finally {
            foreach ($item in $saved.GetEnumerator()) {
                if ($null -eq $item.Value) {
                    Remove-Item "env:$($item.Key)" -ErrorAction SilentlyContinue
                } else {
                    Set-Item -Path "env:$($item.Key)" -Value $item.Value
                }
            }
        }
    }

    Write-Output "  $Name (PID $($proc.Id), log: .cache/logs/$Name.log)" -ForegroundColor DarkGray
    return $proc
}
