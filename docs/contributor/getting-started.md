# Contributor Getting Started

Welcome. This guide gets you from a fresh clone to a green build and your first contribution in the shortest path possible.

---

## Table of Contents

- [Prerequisites at a glance](#prerequisites-at-a-glance)
- [One-command setup](#one-command-setup)
- [Manual setup](#manual-setup)
  - [macOS](#macos)
  - [Linux (Debian / Ubuntu)](#linux-debian--ubuntu)
  - [Linux (Fedora / RHEL)](#linux-fedora--rhel)
  - [Windows](#windows)
- [Verify the setup](#verify-the-setup)
- [IDE setup](#ide-setup)
- [Repository structure at a glance](#repository-structure-at-a-glance)
- [What to read next](#what-to-read-next)

---

## Prerequisites at a glance

| Tool | Required version | Used for |
|------|-----------------|---------|
| **Go** | 1.26+ | Building and testing |
| **make** | Any | Running Makefile targets |
| **golangci-lint** | Latest | Linting (`make lint`) |
| **git** | 2.x+ | Version control and hooks |
| **goimports** | Latest | Import formatting (installed by `make install-tools`) |

Optional but recommended: `pnpm` (for TypeScript test commands in `test-flow`).

---

## One-command setup

The setup script detects your OS, installs missing tools, downloads Go modules, and activates git hooks.

**macOS / Linux:**

```bash
bash scripts/setup-dev.sh
```

**Windows (PowerShell, run as Administrator):**

```powershell
Set-ExecutionPolicy -Scope Process Bypass
.\scripts\setup-dev.ps1
```

> **Windows tip:** The setup script supports `winget`, `chocolatey`, and `scoop`. If none are installed, consider using WSL2 and running the shell script inside your Linux environment — the Go toolchain behaves more predictably there.

After the script finishes, continue to [Verify the setup](#verify-the-setup).

---

## Manual setup

### macOS

```bash
# 1. Install Homebrew if not present
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. Install Go
brew install go
# Verify: go version  →  go1.26.x darwin/arm64

# 3. make is part of the Xcode Command Line Tools
xcode-select --install

# 4. golangci-lint
brew install golangci-lint

# 5. Add Go bin to your PATH (add to ~/.zshrc or ~/.bashrc)
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# 6. Project tools
cd /path/to/dot
make install-tools   # installs goimports
go mod download
make hooks           # activates git hooks
```

### Linux (Debian / Ubuntu)

```bash
# 1. Go (official installer — package repos lag behind)
GOVERSION="1.26"
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -fsSL "https://go.dev/dl/go${GOVERSION}.linux-${ARCH}.tar.gz" -o /tmp/go.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf /tmp/go.tar.gz
echo 'export PATH="/usr/local/go/bin:$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
# Verify: go version

# 2. make
sudo apt-get update && sudo apt-get install -y build-essential

# 3. golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 4. Project tools
cd /path/to/dot
make install-tools
go mod download
make hooks
```

### Linux (Fedora / RHEL)

```bash
# 1. Go
sudo dnf install -y golang    # or use the official installer above

# 2. make
sudo dnf install -y make

# 3-4. Same as Debian from step 3
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
cd /path/to/dot && make install-tools && go mod download && make hooks
```

### Windows

**Option A — WSL2 (recommended):** Install WSL2, open an Ubuntu terminal, and follow the Debian steps above.

**Option B — Native PowerShell:**

```powershell
# Run as Administrator
# Install winget if not present: https://aka.ms/getwinget

winget install GoLang.Go --accept-source-agreements --accept-package-agreements
winget install GnuWin32.Make
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# In the repo directory:
go install golang.org/x/tools/cmd/goimports@latest
go mod download
git config core.hooksPath .githooks
```

---

## Verify the setup

Run these commands in order. Every one should succeed before moving on.

```bash
# 1. Tools are present
go version           # → go1.26.x
make --version       # → GNU Make 4.x
golangci-lint version # → golangci-lint has version ...

# 2. Build
make build
# → bin/dot

# 3. Tests
make test
# → All tests pass

# 4. Full validation
make validate
# → fmt ✓  vet ✓  lint ✓  test ✓

# 5. End-to-end flow tests
make test-flows
# → All N cases passed
```

If any step fails, check the [troubleshooting](#troubleshooting) section.

---

## IDE setup

### VS Code (recommended)

1. Install the **Go** extension by Google (`golang.go`).
2. Open the repo root in VS Code — it will prompt to install recommended Go tools.
3. Enable format-on-save:
   ```json
   // .vscode/settings.json  (already in the repo if present, or create it)
   {
     "editor.formatOnSave": true,
     "[go]": {
       "editor.defaultFormatter": "golang.go"
     },
     "go.lintTool": "golangci-lint",
     "go.lintOnSave": "package"
   }
   ```

### GoLand / IntelliJ

Go support is built in. Enable **golangci-lint** under *Settings → Go → Code Quality Tools → golangci-lint* and point it to the binary from `go env GOPATH`/bin.

---

## Repository structure at a glance

```
dot/
├── cmd/dot/          ← main() — thin entry point, imports plugins
├── flows/            ← Built-in flow definitions + registry
├── generators/       ← Built-in generator packages (one package per generator)
├── plugins/          ← In-tree plugins (biome_extras, ...)
├── examples/         ← Reference plugin implementations
├── tools/test-flow/  ← End-to-end test runner + fixtures (testdata/)
├── scripts/          ← Development setup scripts (you are here)
│
├── internal/
│   ├── cli/          ← Command dispatch, Scaffold(), HuhFormRunner, spinner
│   ├── flow/         ← Question DSL, FlowEngine, HookRegistry
│   ├── spec/         ← ProjectSpec, builder, loader
│   ├── generator/    ← Registry, Executor, Resolver, Sorter, Validator
│   ├── state/        ← VirtualProjectState, Persist, JSON/YAML/GoMod helpers
│   ├── commands/     ← Post-gen + test command planner and runner
│   ├── dotdir/       ← .dot/ read/write (spec.json, manifest.json)
│   ├── plugin/       ← Provider interface, loader, installer
│   └── versioning/   ← Semver parser and constraint checker
│
└── pkg/
    ├── dotapi/       ← Public Generator interface, Manifest, Context (stable)
    └── dotplugin/    ← Public plugin author API — re-exports from internal
```

**Rule of thumb:**
- If it touches the terminal or user interaction → `internal/cli/`
- If it defines the flow question graph → `internal/flow/`
- If it writes files to disk → `internal/state/` or `internal/generator/`
- If it is the public API for generator or plugin authors → `pkg/`

---

## What to read next

After the build is green, pick the guide that matches what you want to do:

| Goal | Start here |
|------|-----------|
| Understand the whole system before touching anything | [architecture.md](architecture.md) |
| Find the right file for a specific task | [navigation-guide.md](navigation-guide.md) ← **start here for most tasks** |
| Add a new flow | [authoring-flows.md](authoring-flows.md) |
| Add a new generator | [authoring-generators.md](authoring-generators.md) |
| Write or publish a plugin | [authoring-plugins.md](authoring-plugins.md) |
| Add or fix a test | [test-flow.md](test-flow.md) |

---

## Troubleshooting

**`golangci-lint: command not found` after `go install`**

Your `GOPATH/bin` is not on `PATH`. Add it:
```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```
Add that line to `~/.zshrc` or `~/.bashrc` to make it permanent.

**`go: module requires go >= 1.26` but you have an older version**

You need Go 1.26 or later. The project uses language features from 1.21+ and the module directive pins the minimum. Re-run the setup script or install manually from [go.dev/dl](https://go.dev/dl).

**`make hooks` reports `chmod: .githooks/pre-push: No such file or directory`**

Only `commit-msg` is required. The error is harmless. The hook validates commit message format before every commit.

**`make test-flows` fails with `no scripted answer for question "biome_extras.strict_mode"`**

A fixture is missing the injected question from the `biome_extras` plugin. Add it:
```json
"biome_extras.strict_mode": false
```
See [test-flow.md](test-flow.md#plugin-injection-fixtures).
