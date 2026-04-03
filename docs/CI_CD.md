# CI/CD Workflows

This document explains the automated CI/CD pipelines for Scaffold CLI.

---

## Overview

We use GitHub Actions to automate code quality checks, testing, and releases.

### Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| **CI** | Push to main, PR | Format check, linting, tests, build |
| **Commitlint** | Push to main, PR | Validate commit messages |
| **Release** | Tag push (v*.*.*) | Build binaries, create release |

---

## CI Workflow

**File:** `.github/workflows/ci.yml`

Runs on every push to `main` and every PR.

### Jobs

#### 1. **Vet** — Go static analysis
- Checks for suspicious constructs
- Fails if issues found
- **Command:** `go vet ./...`

#### 2. **Lint** — Code quality
- Runs golangci-lint
- Checks code formatting with `go fmt`
- Fails if code is not formatted
- **Commands:**
  - `golangci-lint run ./...`
  - `go fmt ./...` (check only)

#### 3. **Test** — Unit tests
- Runs all tests with race detector
- Generates coverage report
- Uploads to Codecov (optional)
- **Command:** `go test -race -v -coverprofile=coverage.out ./...`

#### 4. **Build** — Binary compilation
- Builds the scaffold binary
- Depends on: vet, lint, test (all must pass)
- **Command:** `go build -v -o scaffold ./cmd/scaffold`

### Key Features

✅ **Concurrent execution** — vet, lint, and test run in parallel  
✅ **Fast feedback** — failures reported immediately  
✅ **Race detection** — catches concurrency bugs  
✅ **Code coverage** — uploaded for tracking  
✅ **Required checks** — PR can't be merged if any job fails

---

## Commitlint Workflow

**File:** `.github/workflows/commitlint.yml`

Runs on every push and PR to validate commit messages.

### Job

**Commitlint** — Validate Conventional Commits format

- Uses `.commitlintrc.json` configuration
- Validates commit message format
- Provides helpful error messages on failure
- Automatically comments on PR if commits are invalid

### Commit Format

```
<type>(<scope>): <description>
```

**Allowed types:** feat, fix, docs, style, refactor, perf, test, chore, ci, revert

**Example:** `feat(generators): add redis caching`

See [CONTRIBUTING.md](../CONTRIBUTING.md#commit-conventions) for details.

---

## Release Workflow

**File:** `.github/workflows/release.yml`

Runs when a tag matching `v*.*.*` is pushed.

### Jobs

#### 1. **Build** — Multi-platform binaries
- Builds for:
  - Linux x86_64
  - Linux ARM64
  - macOS x86_64
  - macOS ARM64
  - Windows x86_64

#### 2. **Release** — Create GitHub Release
- Creates release with auto-generated notes
- Attaches compiled binaries
- Publishes to GitHub Releases

### How to Release

```bash
# Tag the commit
git tag v1.0.0

# Push the tag
git push origin v1.0.0
```

The workflow will:
1. Build binaries for all platforms
2. Create a GitHub Release
3. Attach binaries to release
4. Generate release notes from commits

---

## Local vs CI

### Local Checks (Before Push)

```bash
# Run locally before pushing
make validate
```

This runs:
- Format check: `go fmt ./...`
- Vet check: `go vet ./...`
- Linting: `golangci-lint run ./...`
- Tests: `go test -race ./...`

### CI Checks (On GitHub)

Same checks run automatically:
- All jobs must pass
- PR can't be merged if failed
- Status checks prevent accidents

### Commit Message Check

**Local:** `.githooks/commit-msg` hook validates on every commit

**CI:** `commitlint` validates on every PR and push

Both use Conventional Commits format for consistency.

---

## Troubleshooting

### CI Passed Locally but Failed on GitHub

**Cause:** Different Go versions or environment

**Fix:**
- Check Go version: `go version`
- Ensure it matches `.github/workflows/ci.yml` (currently Go 1.26)
- Run `go mod download` to sync dependencies

### Commitlint Failed on PR

**Cause:** Commits don't follow Conventional Commits format

**Fix:**
```bash
# View rules
make commit-lint

# Fix commit
git commit --amend -m "feat(scope): correct message"

# Push again
git push origin feature-branch
```

### Build Job Fails but Tests Pass

**Cause:** Usually missing imports or race conditions

**Fix:**
- Check build output for errors
- Run locally: `go build ./cmd/scaffold`
- Fix issues and re-push

### Timeout or Stuck Job

**Action:** GitHub will cancel jobs after 6 hours

**Prevention:**
- Jobs usually complete in <5 minutes
- If stuck, cancel and re-push to retry

---

## Configuration Files

| File | Purpose |
|------|---------|
| `.github/workflows/ci.yml` | Main CI pipeline |
| `.github/workflows/commitlint.yml` | Commit validation |
| `.github/workflows/release.yml` | Release automation |
| `.commitlintrc.json` | Commitlint rules |
| `.githooks/commit-msg` | Local commit hook |

---

## Monitoring

### GitHub Status Checks

View on any PR or push:
- All jobs must have ✅ check mark
- Red ✗ means failure
- Gray ⏳ means in progress

### Codecov (Optional)

If enabled, code coverage appears on PRs:
- Shows coverage % change
- Flags significant drops
- Helps track test quality

---

## Cost

All workflows are free with GitHub Actions:
- 2000 minutes/month included
- Our workflows: ~2-5 minutes per run
- Should never exceed limits

---

## Future Improvements

Potential additions:
- [ ] SAST (Security scanning)
- [ ] Dependency scanning
- [ ] Docker image builds
- [ ] Automated versioning
- [ ] Changelog generation

---

## Questions?

See [CONTRIBUTING.md](../CONTRIBUTING.md) or open a [Discussion](https://github.com/version14/scaffold-cli/discussions).
