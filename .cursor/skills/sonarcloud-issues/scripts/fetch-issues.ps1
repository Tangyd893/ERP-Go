# Fetch SonarCloud open issues via /api/issues/search and format for LLM/docs.
# Usage:
#   .\fetch-issues.ps1
#   .\fetch-issues.ps1 -Branch main -OutFile .cache/sonar-issues.md
#   .\fetch-issues.ps1 -Types VULNERABILITY -Severities BLOCKER,CRITICAL
#   .\fetch-issues.ps1 -Format json -OutFile .cache/sonar-issues.json

param(
    [string]$ProjectKey = $(if ($env:SONAR_PROJECT_KEY) { $env:SONAR_PROJECT_KEY } else { "Tangyd893_ERP-Go" }),
    [string]$Branch = $(if ($env:SONAR_BRANCH) { $env:SONAR_BRANCH } else { "" }),
    [string]$IssueStatuses = "OPEN,CONFIRMED",
    [string]$Types = "",
    [string]$Severities = "",
    [ValidateSet("markdown", "json")]
    [string]$Format = "markdown",
    [int]$PageSize = 500,
    [string]$OutFile = ""
)

$ErrorActionPreference = "Stop"
$BaseUrl = "https://sonarcloud.io/api/issues/search"

function Get-QueryString {
    $pairs = @(
        "componentKeys=$([uri]::EscapeDataString($ProjectKey))"
        "issueStatuses=$([uri]::EscapeDataString($IssueStatuses))"
        "ps=$PageSize"
        "facets=$([uri]::EscapeDataString('types,severities,rules,languages'))"
    )
    if ($Branch) { $pairs += "branch=$([uri]::EscapeDataString($Branch))" }
    if ($Types) { $pairs += "types=$([uri]::EscapeDataString($Types))" }
    if ($Severities) { $pairs += "severities=$([uri]::EscapeDataString($Severities))" }
    $pairs -join "&"
}

function Invoke-SonarPage([int]$Page) {
    $qs = Get-QueryString
    $uri = "${BaseUrl}?${qs}&p=$Page"
    $headers = @{}
    if ($env:SONAR_TOKEN) {
        $headers["Authorization"] = "Bearer $($env:SONAR_TOKEN)"
    }
    Invoke-RestMethod -Uri $uri -Headers $headers -Method Get
}

function Get-RelativePath([string]$Component) {
    if ($Component -match "^[^:]+:(.+)$") {
        return $Matches[1]
    }
    return $Component
}

# Paginate all issues
$allIssues = [System.Collections.Generic.List[object]]::new()
$facets = $null
$total = 0
$effortTotal = 0
$debtTotal = 0
$page = 1

do {
    $resp = Invoke-SonarPage -Page $page
    if ($page -eq 1) {
        $facets = $resp.facets
        $total = [int]$resp.total
        $effortTotal = $resp.effortTotal
        $debtTotal = $resp.debtTotal
    }
    foreach ($issue in $resp.issues) {
        $allIssues.Add($issue)
    }
    $page++
} while ($allIssues.Count -lt $total -and $resp.issues.Count -gt 0)

$payload = [ordered]@{
    projectKey     = $ProjectKey
    branch         = $(if ($Branch) { $Branch } else { $null })
    issueStatuses  = $IssueStatuses
    fetchedAt      = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
    total          = $total
    effortTotalMin = $effortTotal
    debtTotalMin   = $debtTotal
    issueCount     = $allIssues.Count
    facets         = $facets
    issues         = $allIssues
}

if ($Format -eq "json") {
    $out = $payload | ConvertTo-Json -Depth 12
    if ($OutFile) {
        $dir = Split-Path -Parent $OutFile
        if ($dir) { New-Item -ItemType Directory -Force -Path $dir | Out-Null }
        $out | Set-Content -Path $OutFile -Encoding UTF8
        Write-Host "Wrote $OutFile ($($allIssues.Count) issues)"
    } else {
        $out
    }
    return
}

# --- Markdown for LLM ---
$uiBase = "https://sonarcloud.io/project/issues"
$uiQs = "issueStatuses=$([uri]::EscapeDataString($IssueStatuses))&id=$([uri]::EscapeDataString($ProjectKey))"
if ($Branch) { $uiQs += "&branch=$([uri]::EscapeDataString($Branch))" }
$uiLink = "${uiBase}?${uiQs}"

$sb = [System.Text.StringBuilder]::new()
[void]$sb.AppendLine("# SonarCloud Issues Report")
[void]$sb.AppendLine()
[void]$sb.AppendLine("| Field | Value |")
[void]$sb.AppendLine("| --- | --- |")
[void]$sb.AppendLine('| Project | `' + $ProjectKey + '` |')
$branchCell = if ($Branch) { '`' + $Branch + '`' } else { '*(all)*' }
[void]$sb.AppendLine("| Branch | $branchCell |")
[void]$sb.AppendLine('| Statuses | `' + $IssueStatuses + '` |')
[void]$sb.AppendLine("| Total (API) | $total |")
[void]$sb.AppendLine("| Fetched | $($allIssues.Count) |")
[void]$sb.AppendLine("| Effort / Debt (min) | $effortTotal / $debtTotal |")
[void]$sb.AppendLine("| Fetched at (UTC) | $($payload.fetchedAt) |")
[void]$sb.AppendLine("| UI | $uiLink |")
[void]$sb.AppendLine()

if ($facets) {
    [void]$sb.AppendLine("## Summary (facets)")
    [void]$sb.AppendLine()
    foreach ($facet in $facets) {
        $prop = $facet.property
        [void]$sb.AppendLine("### $prop")
        [void]$sb.AppendLine()
        [void]$sb.AppendLine("| Value | Count |")
        [void]$sb.AppendLine("| --- | ---: |")
        foreach ($v in $facet.values) {
            [void]$sb.AppendLine("| $($v.val) | $($v.count) |")
        }
        [void]$sb.AppendLine()
    }
}

$severityOrder = @("BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO")
[void]$sb.AppendLine("## Issues by severity")
[void]$sb.AppendLine()

foreach ($sev in $severityOrder) {
    $group = $allIssues | Where-Object { $_.severity -eq $sev }
    if (-not $group) { continue }
    [void]$sb.AppendLine("### $sev ($($group.Count))")
    [void]$sb.AppendLine()
    foreach ($i in ($group | Sort-Object { Get-RelativePath $_.component }, line)) {
        $path = Get-RelativePath $i.component
        $line = if ($null -ne $i.line) { $i.line } else { "-" }
        $msg = ($i.message -replace '\|', '/')
        [void]$sb.AppendLine("- **${path}:${line}** | ``$($i.rule)`` | $($i.type) | $msg | effort: $($i.effort)")
    }
    [void]$sb.AppendLine()
}

[void]$sb.AppendLine("## Issues by file")
[void]$sb.AppendLine()
$byFile = $allIssues | Group-Object { Get-RelativePath $_.component } | Sort-Object Count -Descending
foreach ($g in $byFile) {
    [void]$sb.AppendLine("### $($g.Name) ($($g.Count))")
    [void]$sb.AppendLine()
    foreach ($i in ($g.Group | Sort-Object severity, line)) {
        [void]$sb.AppendLine("- L$($i.line) **$($i.severity)** ``$($i.rule)`` ($($i.type)): $($i.message)")
    }
    [void]$sb.AppendLine()
}

$md = $sb.ToString()
if ($OutFile) {
    $dir = Split-Path -Parent $OutFile
    if ($dir) { New-Item -ItemType Directory -Force -Path $dir | Out-Null }
    $md | Set-Content -Path $OutFile -Encoding UTF8
    Write-Host "Wrote $OutFile ($($allIssues.Count) issues)"
} else {
    $md
}
