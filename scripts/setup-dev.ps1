# setup-dev.ps1 — Install every tool needed to build and test DOT on Windows.
#
# Prerequisites: PowerShell 5+ or PowerShell Core 7+ (pwsh).
# Run once after cloning (in an Administrator PowerShell terminal):
#
#   Set-ExecutionPolicy -Scope Process Bypass
#   .\scripts\setup-dev.ps1
#
# The script is idempotent — safe to run again after upgrades.
# WSL2 users: prefer scripts/setup-dev.sh inside your Linux environment.

$ErrorActionPreference = "Stop"

$REQUIRED_GO = "1.26"

function Write-Step($msg)    { Write-Host "  -> $msg" -ForegroundColor Cyan }
function Write-Success($msg) { Write-Host "  v  $msg" -ForegroundColor Green }
function Write-Warn($msg)    { Write-Host "  !  $msg" -ForegroundColor Yellow }
function Write-Fail($msg)    { Write-Host "  x  $msg" -ForegroundColor Red; exit 1 }

Write-Host "`nDOT - development environment setup`n" -ForegroundColor White

# ── Helper: compare semver ────────────────────────────────────────────────────

function Compare-Version([string]$a, [string]$b) {
    $va = [version]($a -replace '-.*','')
    $vb = [version]($b -replace '-.*','')
    return $va.CompareTo($vb)
}

# ── Detect package manager ────────────────────────────────────────────────────

$HAS_WINGET = $null -ne (Get-Command winget -ErrorAction SilentlyContinue)
$HAS_CHOCO  = $null -ne (Get-Command choco  -ErrorAction SilentlyContinue)
$HAS_SCOOP  = $null -ne (Get-Command scoop  -ErrorAction SilentlyContinue)

if (-not $HAS_WINGET -and -not $HAS_CHOCO -and -not $HAS_SCOOP) {
    Write-Warn "No package manager found (winget / chocolatey / scoop)."
    Write-Warn "Install winget: https://aka.ms/getwinget"
    Write-Warn "or Chocolatey:  https://chocolatey.org/install"
    Write-Warn "Then re-run this script."
    Write-Warn ""
    Write-Warn "Alternatively, use WSL2 and run scripts/setup-dev.sh inside Linux."
}

# ── Step 1: Go ────────────────────────────────────────────────────────────────

Write-Host "`nStep 1/4 - Go $REQUIRED_GO+" -ForegroundColor White

$goCmd = Get-Command go -ErrorAction SilentlyContinue
if ($goCmd) {
    $current = (go version) -replace '.*go(\S+).*','$1'
    if ((Compare-Version $current $REQUIRED_GO) -ge 0) {
        Write-Success "Go $current already installed"
    } else {
        Write-Warn "Go $current is older than required $REQUIRED_GO — upgrading"
        if ($HAS_WINGET) { winget upgrade --id GoLang.Go --accept-source-agreements }
        elseif ($HAS_CHOCO) { choco upgrade golang -y }
        elseif ($HAS_SCOOP) { scoop update go }
    }
} else {
    Write-Step "Installing Go $REQUIRED_GO ..."
    if ($HAS_WINGET) { winget install --id GoLang.Go --accept-source-agreements --accept-package-agreements }
    elseif ($HAS_CHOCO) { choco install golang -y }
    elseif ($HAS_SCOOP) { scoop install go }
    else { Write-Fail "Cannot install Go: no package manager. Install manually from https://go.dev/dl/" }
    Write-Success "Go installed"
}

# Reload PATH so 'go' is available in this session
$env:PATH = [System.Environment]::GetEnvironmentVariable("PATH","Machine") + ";" +
            [System.Environment]::GetEnvironmentVariable("PATH","User")

# ── Step 2: make ──────────────────────────────────────────────────────────────

Write-Host "`nStep 2/4 - make" -ForegroundColor White

if (Get-Command make -ErrorAction SilentlyContinue) {
    Write-Success "make already installed"
} else {
    Write-Step "Installing make ..."
    if ($HAS_WINGET) { winget install --id GnuWin32.Make --accept-source-agreements }
    elseif ($HAS_CHOCO) { choco install make -y }
    elseif ($HAS_SCOOP) { scoop install make }
    else { Write-Warn "Please install make manually (https://gnuwin32.sourceforge.net/packages/make.htm)" }
    Write-Success "make installed"
}

# ── Step 3: golangci-lint ─────────────────────────────────────────────────────

Write-Host "`nStep 3/4 - golangci-lint" -ForegroundColor White

if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
    Write-Success "golangci-lint already installed"
} else {
    Write-Step "Installing golangci-lint via go install ..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    Write-Success "golangci-lint installed"
}

# ── Step 4: project tools + hooks ────────────────────────────────────────────

Write-Host "`nStep 4/4 - project tools + git hooks" -ForegroundColor White

Write-Step "Installing goimports ..."
go install golang.org/x/tools/cmd/goimports@latest
Write-Success "goimports installed"

Write-Step "Downloading module dependencies ..."
go mod download
Write-Success "Dependencies ready"

Write-Step "Activating git hooks ..."
git config core.hooksPath .githooks
Write-Success "Git hooks active"

# ── Done ──────────────────────────────────────────────────────────────────────

Write-Host "`nSetup complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:"
Write-Host "  make build      - compile dot"
Write-Host "  make validate   - fmt + vet + lint + test"
Write-Host "  make test-flows - end-to-end fixture tests"
Write-Host ""
Write-Host "Read docs/contributor/getting-started.md for your first contribution walkthrough."
Write-Host ""
