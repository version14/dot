# Contributing to dot

Thank you for your interest in contributing. This document explains how to get involved, what we expect, and how to get your changes merged.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Features](#suggesting-features)
  - [Activating git hooks](#activating-git-hooks)
  - [Submitting Code Changes](#submitting-code-changes)
- [Commit Conventions](#commit-conventions)
- [Pull Request Process](#pull-request-process)
- [Code Style](#code-style)
- [Testing](#testing)
- [Documentation](#documentation)

---

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold these standards.

---

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/your-username/dot.git
   cd dot
   ```
3. **Add the upstream remote:**
   ```bash
   git remote add upstream https://github.com/version14/dot.git
   ```
4. Follow the [Development Setup](#development-setup) section below and [docs/getting-started.md](docs/getting-started.md).

---

## Development Setup

### Prerequisites

- Go 1.21+ (`go version`)
- `git` on `$PATH`
- `golangci-lint` (for linting): `brew install golangci-lint` or see [golangci-lint docs](https://golangci-lint.run/usage/install/)

### Build and run

```bash
make build        # produces bin/dot
./bin/dot version
```

### Full validation

```bash
make validate     # fmt → vet → lint → test
```

### End-to-end test

```bash
make test-flow    # runs test-flow against all testdata fixtures (skips test commands)
```

See [docs/test-flow.md](docs/test-flow.md) for the full guide.

---

## How to Contribute

### Reporting Bugs

Before opening an issue:
- Search [existing issues](../../issues) to avoid duplicates.
- Make sure you are on the latest version (`git pull upstream main`).

When opening a bug report, include:
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Go version, `dot version`)
- Relevant logs or error output

### Suggesting Features

Open a **Feature Request** issue with:
- A clear description of the problem the feature solves
- Your proposed solution
- Alternatives you considered

Features that align with the project architecture and roadmap are more likely to be accepted.

### Activating git hooks

Git hooks validate commit messages locally before they are created. Activate them once after cloning:

```bash
make hooks
```

Or manually:
```bash
git config core.hooksPath .githooks
chmod +x .githooks/commit-msg
```

### Submitting Code Changes

1. **Create a branch** from `main`:
   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feat/your-feature-name
   ```

2. **Make your changes** following the [Code Style](#code-style) guidelines.

3. **Write or update tests** — every new behavior needs a test. If you changed a flow or generator, add or update a `test-flow` fixture (see [docs/test-flow.md](docs/test-flow.md)).

4. **Update documentation** — see [Documentation](#documentation) for the rules.

5. **Run validation locally**:
   ```bash
   make validate
   make test-flow
   ```

6. **Commit following [Commit Conventions](#commit-conventions)**.

7. **Push and open a Pull Request**:
   ```bash
   git push origin feat/your-feature-name
   ```

**Before submitting the PR, verify:**

- [ ] All validations pass (`make validate`)
- [ ] Commits follow Conventional Commits
- [ ] Tests pass (`make test`)
- [ ] `test-flow` fixtures pass (`make test-flow`)
- [ ] Documentation is updated (see [Documentation rules](#documentation))

---

## Commit Conventions

We follow **Conventional Commits** format. Messages are validated both locally (via git hook) and in CI.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

| Type       | When to use                           |
|------------|---------------------------------------|
| `feat`     | New feature or behavior               |
| `fix`      | Bug fix                               |
| `docs`     | Documentation only                    |
| `style`    | Code style (formatting, semicolons)   |
| `refactor` | Code change with no behavior change   |
| `perf`     | Performance improvement               |
| `test`     | Adding or updating tests              |
| `chore`    | Tooling, dependencies, config         |
| `ci`       | CI/CD changes                         |
| `revert`   | Revert a previous commit              |

**Scope** (optional): the area affected, e.g. `flow`, `generator`, `cli`, `plugin`, `spec`.

### Examples

```
feat(flow): add IfQuestion conditional node
fix(resolver): preserve loop iteration invocations after dedup
docs(test-flow): document loop fixture schema
refactor(cli): extract Scaffold into ScaffoldOptions struct
test(generator): add validator round-trip tests
chore: upgrade huh to v1.0.1
```

### Rules

- Type is required (lowercase)
- Scope is optional (lowercase)
- Description starts with lowercase, no period at end
- Max 100 characters for the subject line
- Use imperative mood ("add" not "adds")
- Reference issues in the footer: `Closes #42`

---

## Pull Request Process

1. **One PR per concern** — don't mix unrelated changes
2. **Fill the PR template** — describe what changed and why
3. **Keep diffs small** — large PRs are hard to review; split if needed
4. **All CI checks must pass** before merging
5. **Address review comments** — iterate on feedback

PRs are merged by maintainers once they have one approving review and all checks are green.

---

## Code Style

- Standard Go style: `gofmt`-formatted, idiomatic.
- No `internal/*` imports from `pkg/` or `plugins/` — use `pkg/dotapi` and `pkg/dotplugin`.
- Error messages: lowercase, no trailing period, wrap with `%w` for context.
- No `panic` in library code (only in `main` or test setup).
- Keep functions short. If a function needs a comment to explain what it does, consider splitting it.

Run `make fmt` and `make lint` before committing.

---

## Testing

Every PR should maintain or improve existing test coverage.

### Unit tests

```bash
make test           # go test -race ./...
```

Critical areas that require table-driven tests:
- `internal/flow/` — question branching, engine traversal
- `internal/generator/` — resolver dedup, topo-sort, validator
- `internal/versioning/` — semver parsing and constraint matching
- `internal/spec/` — spec serialization round-trips

### End-to-end tests (test-flow)

```bash
make test-flow      # go run ./tools/test-flow -skip-test
```

Every flow change needs a matching fixture. See [docs/test-flow.md](docs/test-flow.md).

---

## Documentation

The `docs/` directory is the single source of truth. Read [docs/README.md](docs/README.md) for the full documentation rules.

**Summary of when to update docs in a PR:**

| Change | Required update |
|--------|----------------|
| New CLI command or flag | `docs/cli-reference.md` |
| New flow | `docs/authoring-flows.md` + test fixture |
| New question type | `docs/authoring-flows.md` + `docs/architecture.md` |
| New injection kind | `docs/authoring-plugins.md` |
| New exported type in `pkg/dotapi` or `pkg/dotplugin` | `docs/authoring-generators.md` or `docs/authoring-plugins.md` |
| **New generator** | **Create `docs/generators/<name>.md`** (copy `docs/generators/_template.md`) + update `docs/README.md` |
| **New plugin** | **Create `docs/plugins/<name>.md`** (copy `docs/plugins/_template.md`) + update `docs/README.md` |
| **New flow** | **Create `docs/flows/<id>.md`** (copy `docs/flows/_template.md`) + update `docs/README.md` |
| Generator manifest fields change | `docs/generators/<name>.md` |
| Plugin injection IDs change | `docs/plugins/<name>.md` + affected test fixtures |
| Pipeline step change | `docs/architecture.md` |
| `.dot/` schema change | `docs/architecture.md` |
| New `test-flow` flag | `docs/test-flow.md` |
| Install mechanism change | `docs/getting-started.md` |

Documentation updates use the `docs` commit type:

```
docs(authoring-flows): document IfQuestion conditional node
```

---

## Questions?

Open a [Discussion](../../discussions) or read the [docs/](docs/README.md).
