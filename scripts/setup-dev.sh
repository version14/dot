#!/usr/bin/env bash
# setup-dev.sh — Install every tool needed to build and test DOT.
#
# Supports: macOS (Homebrew), Debian/Ubuntu, Fedora/RHEL, Alpine.
# Run once after cloning:
#
#   bash scripts/setup-dev.sh
#
# The script is idempotent — safe to run again after upgrades.
set -euo pipefail

# ── Colours ─────────────────────────────────────────────────────────────────

BOLD="\033[1m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
CYAN="\033[0;36m"
RED="\033[0;31m"
RESET="\033[0m"

info()    { echo -e "  ${CYAN}→${RESET} $*"; }
success() { echo -e "  ${GREEN}✓${RESET} $*"; }
warn()    { echo -e "  ${YELLOW}⚠${RESET} $*"; }
error()   { echo -e "  ${RED}✗${RESET} $*" >&2; exit 1; }

echo -e "\n${BOLD}DOT — development environment setup${RESET}\n"

# ── Detect OS ────────────────────────────────────────────────────────────────

OS=""
case "$(uname -s)" in
  Darwin) OS="macos" ;;
  Linux)
    if   [ -f /etc/debian_version ]; then OS="debian"
    elif [ -f /etc/fedora-release ];  then OS="fedora"
    elif [ -f /etc/alpine-release ];  then OS="alpine"
    else                                   OS="linux"
    fi ;;
  *) error "Unsupported OS: $(uname -s). Use scripts/setup-dev.ps1 on Windows." ;;
esac
info "Detected OS: ${OS}"

# ── Helper: version check ────────────────────────────────────────────────────

# Returns 0 if $1 >= $2 (both as X.Y.Z strings).
version_ge() {
  [ "$(printf '%s\n' "$1" "$2" | sort -V | head -1)" = "$2" ]
}

REQUIRED_GO="1.26"

# ── Step 1: Go ───────────────────────────────────────────────────────────────

echo -e "\n${BOLD}Step 1/4 — Go ${REQUIRED_GO}+${RESET}"

install_go_from_release() {
  local arch
  case "$(uname -m)" in
    x86_64)  arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *) error "Unsupported arch $(uname -m) for automatic Go install. Install manually: https://go.dev/dl/" ;;
  esac
  local goos="linux"
  [ "${OS}" = "macos" ] && goos="darwin"
  local url="https://go.dev/dl/go${REQUIRED_GO}.${goos}-${arch}.tar.gz"
  info "Downloading Go ${REQUIRED_GO} from ${url} …"
  curl -fsSL "${url}" -o /tmp/go.tar.gz
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf /tmp/go.tar.gz
  rm /tmp/go.tar.gz
  export PATH="/usr/local/go/bin:${PATH}"
  success "Go ${REQUIRED_GO} installed to /usr/local/go"
  warn "Add to your shell profile:  export PATH=\"/usr/local/go/bin:\$PATH\""
}

if command -v go &>/dev/null; then
  CURRENT_GO="$(go version | awk '{print $3}' | sed 's/go//')"
  if version_ge "${CURRENT_GO}" "${REQUIRED_GO}"; then
    success "Go ${CURRENT_GO} already installed"
  else
    warn "Go ${CURRENT_GO} is older than required ${REQUIRED_GO}"
    case "${OS}" in
      macos)
        if command -v brew &>/dev/null; then
          info "Upgrading via Homebrew …"
          brew upgrade go || brew install go
        else
          install_go_from_release
        fi ;;
      debian) sudo apt-get remove -y golang-go &>/dev/null || true; install_go_from_release ;;
      *) install_go_from_release ;;
    esac
  fi
else
  info "Go not found — installing …"
  case "${OS}" in
    macos)
      if command -v brew &>/dev/null; then
        brew install go
      else
        install_go_from_release
      fi ;;
    debian)
      sudo apt-get update -qq
      sudo apt-get install -y golang-go || install_go_from_release ;;
    fedora)
      sudo dnf install -y golang || install_go_from_release ;;
    alpine)
      sudo apk add --no-cache go || install_go_from_release ;;
    *) install_go_from_release ;;
  esac
fi

# Ensure GOPATH/bin is on PATH for this session
export PATH="$(go env GOPATH)/bin:${PATH}"

# ── Step 2: make ─────────────────────────────────────────────────────────────

echo -e "\n${BOLD}Step 2/4 — make${RESET}"

if command -v make &>/dev/null; then
  success "make $(make --version | head -1 | awk '{print $3}') already installed"
else
  info "Installing make …"
  case "${OS}" in
    macos)
      xcode-select --install 2>/dev/null || true
      warn "If the Xcode CLT installer opened, re-run this script after it completes." ;;
    debian) sudo apt-get update -qq && sudo apt-get install -y build-essential ;;
    fedora) sudo dnf install -y make ;;
    alpine) sudo apk add --no-cache make ;;
    *) warn "Please install make manually for your distribution." ;;
  esac
fi

# ── Step 3: golangci-lint ────────────────────────────────────────────────────

echo -e "\n${BOLD}Step 3/4 — golangci-lint${RESET}"

if command -v golangci-lint &>/dev/null; then
  success "golangci-lint $(golangci-lint version 2>&1 | head -1) already installed"
else
  info "Installing golangci-lint …"
  case "${OS}" in
    macos)
      if command -v brew &>/dev/null; then
        brew install golangci-lint
      else
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
      fi ;;
    *) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest ;;
  esac
  success "golangci-lint installed"
fi

# ── Step 4: project tools + hooks ────────────────────────────────────────────

echo -e "\n${BOLD}Step 4/4 — project tools + git hooks${RESET}"

info "Installing Go tools (goimports) …"
go install golang.org/x/tools/cmd/goimports@latest
success "goimports installed"

info "Downloading module dependencies …"
go mod download
success "Dependencies ready"

info "Activating git hooks …"
git config core.hooksPath .githooks
chmod +x .githooks/commit-msg .githooks/pre-push 2>/dev/null || true
success "Git hooks active"

# ── Done ─────────────────────────────────────────────────────────────────────

echo -e "\n${GREEN}${BOLD}Setup complete!${RESET}"
echo -e "\nNext steps:"
echo -e "  ${CYAN}make build${RESET}      — compile dot"
echo -e "  ${CYAN}make validate${RESET}   — fmt + vet + lint + test"
echo -e "  ${CYAN}make test-flows${RESET} — end-to-end fixture tests"
echo -e "\nRead ${BOLD}docs/contributor/getting-started.md${RESET} for your first contribution walkthrough.\n"
