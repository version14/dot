#!/usr/bin/env bash
# Build script for dot binaries.
#
# Usage:
#   ./scripts/build.sh              # build everything
#   ./scripts/build.sh -a           # same
#   ./scripts/build.sh -t <name>    # build one tool from tools/<name>
#   ./scripts/build.sh -d / --dot   # build only the dot CLI
#
# Called by the Makefile as: make build [ARGS="..."]

set -euo pipefail

BIN_DIR="bin"
GO="${GO:-go}"
GIT_VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"
LDFLAGS="-ldflags=-X main.toolVersion=${GIT_VERSION}"

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
RESET='\033[0m'

info()    { echo -e "${BLUE}→ $*${RESET}"; }
success() { echo -e "${GREEN}✓ $*${RESET}"; }
die()     { echo -e "${RED}✗ $*${RESET}" >&2; exit 1; }

build_dot() {
  info "Building dot..."
  mkdir -p "${BIN_DIR}"
  ${GO} build -buildvcs=false -v "${LDFLAGS}" -o "${BIN_DIR}/dot" ./cmd/dot
  success "Built ${BIN_DIR}/dot"
}

build_tool() {
  local name="$1"
  local src="./tools/${name}"
  [[ -d "${src}" ]] || die "No tool found at ${src}"
  info "Building ${name}..."
  mkdir -p "${BIN_DIR}"
  ${GO} build -buildvcs=false -v -o "${BIN_DIR}/${name}" "${src}"
  success "Built ${BIN_DIR}/${name}"
}

build_all_tools() {
  for dir in tools/*/; do
    local name
    name="$(basename "${dir}")"
    # Skip directories that are not standalone binaries (no main package).
    if ${GO} list -f '{{.Name}}' "./${dir}" 2>/dev/null | grep -q '^main$'; then
      build_tool "${name}"
    fi
  done
}

build_all() {
  build_dot
  build_all_tools
}

# --- argument parsing ---

MODE="all"
TOOL_NAME=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    -a|--all)
      MODE="all"
      shift
      ;;
    -d|--dot)
      MODE="dot"
      shift
      ;;
    -t|--tool)
      [[ $# -ge 2 ]] || die "-t requires a tool name"
      MODE="tool"
      TOOL_NAME="$2"
      shift 2
      ;;
    *)
      die "Unknown argument: $1  (valid: -a, -d/--dot, -t <name>)"
      ;;
  esac
done

case "${MODE}" in
  all)  build_all ;;
  dot)  build_dot ;;
  tool) build_tool "${TOOL_NAME}" ;;
esac
