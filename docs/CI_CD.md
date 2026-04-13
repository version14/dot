# CI/CD Workflows — dot

---

## Overview

We use GitHub Actions for code quality and releases.

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| **CI** | Push / PR | Vet, lint, test, build |
| **Commitlint** | Push / PR | Validate commit messages |
| **Release** | `v*.*.*` tag | GoReleaser: multi-platform binaries, GitHub Release, Homebrew tap |

---

## CI Workflow

**File:** `.github/workflows/ci.yml`

Runs on every push to `main` and every PR. Jobs run in parallel; build depends on all three.

| Job | Command | What it checks |
|-----|---------|----------------|
| **vet** | `go vet ./...` | Suspicious constructs, misuse of sync primitives |
| **lint** | `golangci-lint run ./...` | Style, errcheck, staticcheck |
| **test** | `go test -race -v -coverprofile=coverage.out ./...` | Correctness + race conditions |
| **build** | `go build -v -o dot ./cmd/dot` | Compiles successfully |

Local equivalent:

```bash
make validate   # fmt → vet → lint → test
```

---

## Commitlint Workflow

**File:** `.github/workflows/commitlint.yml`

Validates every commit message on push and PRs. Format:

```
<type>(<scope>): <description>
```

Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `ci`, `revert`

The same check runs locally when you activate git hooks:

```bash
make hooks
```

See [CONTRIBUTING.md](../CONTRIBUTING.md#commit-conventions) for the full rules.

---

## Release Workflow

**File:** `.github/workflows/release.yml`

Triggered by any tag matching `v*.*.*`. Runs [GoReleaser](https://goreleaser.com) via the official action.

### What GoReleaser does

1. Builds binaries for all targets (CGO disabled, `buildVersion` stamped via `-ldflags`)
2. Packages each binary into a `.tar.gz` (Linux/macOS) or `.zip` (Windows)
3. Generates `checksums.txt` (SHA-256)
4. Creates a GitHub Release with all artifacts and auto-generated release notes
5. Pushes an updated Homebrew formula to `github.com/version14/homebrew-tap`

### Build targets

| OS | Arch | Archive |
|----|------|---------|
| Linux | amd64 | `dot_VERSION_linux_amd64.tar.gz` |
| Linux | arm64 | `dot_VERSION_linux_arm64.tar.gz` |
| macOS | amd64 | `dot_VERSION_darwin_amd64.tar.gz` |
| macOS | arm64 | `dot_VERSION_darwin_arm64.tar.gz` |
| Windows | amd64 | `dot_VERSION_windows_amd64.zip` |

### Required secrets

| Secret | Where to get it | Used for |
|--------|----------------|---------|
| `GITHUB_TOKEN` | Automatic | GitHub Release, asset upload |
| `HOMEBREW_TAP_TOKEN` | Personal Access Token with `repo` write on `homebrew-tap` | Push Homebrew formula |

Add secrets at: `github.com/version14/dot → Settings → Secrets → Actions`

### How to cut a release

```bash
git tag v0.1.0
git push origin v0.1.0
```

GoReleaser runs automatically. The release appears at `github.com/version14/dot/releases` and the Homebrew formula at `github.com/version14/homebrew-tap` is updated within a minute.

### GoReleaser config

**File:** `.goreleaser.yaml` — edit this to change archive names, add new targets, or update the Homebrew formula template.

---

## Troubleshooting

**CI passed locally but failed on GitHub**

```bash
go version   # must match go.mod directive
go mod tidy
```

**Commitlint failed**

```bash
make commit-lint                              # view rules
git commit --amend -m "fix: correct message"
git push --force-with-lease origin your-branch
```

**Release failed: Homebrew push rejected**

Check that `HOMEBREW_TAP_TOKEN` is set in repo secrets and has `repo` write access to `github.com/version14/homebrew-tap`.

**Build job fails**

```bash
go build ./cmd/dot   # run locally, read the error
go vet ./...
```

---

## Configuration files

| File | Purpose |
|------|---------|
| `.github/workflows/ci.yml` | CI pipeline |
| `.github/workflows/commitlint.yml` | Commit message validation |
| `.github/workflows/release.yml` | Release via GoReleaser |
| `.goreleaser.yaml` | GoReleaser config (targets, archives, Homebrew) |
| `.commitlintrc.json` | Commitlint rules |
| `.githooks/commit-msg` | Local commit hook |
