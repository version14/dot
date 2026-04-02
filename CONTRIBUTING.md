# Contributing

Thank you for your interest in contributing! This document explains how to get involved, what we expect, and how to get your changes merged.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Features](#suggesting-features)
  - [Activating the commit-msg hook](#activating-the-commit-msg-hook)
  - [Submitting Code Changes](#submitting-code-changes)
- [Commit Conventions](#commit-conventions)
- [Pull Request Process](#pull-request-process)
- [Code Style](#code-style)
- [Testing](#testing)

---

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold these standards. Please report unacceptable behavior to the maintainers.

---

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/your-username/scaffold-cli.git
   cd scaffold-cli
   ```
3. **Add the upstream remote:**
   ```bash
   git remote add upstream https://github.com/version14/scaffold-cli.git
   ```
4. Follow the [development setup guide](docs/getting-started/README.md)

---

## Development Setup

See [Getting Started](docs/getting-started/README.md) for the full setup guide.

---

## How to Contribute

### Reporting Bugs

Before opening an issue:
- Search [existing issues](../../issues) to avoid duplicates
- Make sure you are on the latest version (`git pull upstream main`)

When opening a bug report, use the **Bug Report** template and include:
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, runtime version, relevant tool versions)
- Relevant logs or screenshots

### Suggesting Features

Open a **Feature Request** issue with:
- A clear description of the problem the feature solves
- Your proposed solution
- Alternatives you considered

Features that align with the project's scope and architecture are more likely to be accepted.

### Activating the commit-msg hook

A shell-based hook validates commit messages locally against the Conventional Commits spec. Activate it once after cloning — no runtime or dependencies required:

```bash
git config core.hooksPath .githooks
```

Commit messages are also validated in CI on every PR via GitHub Actions.

### Submitting Code Changes

1. Create a branch from `main`:
   ```bash
   git checkout main
   git pull upstream main
   git checkout -b feat/your-feature-name
   ```
2. Make your changes following the [code style](#code-style) guidelines
3. Write or update tests
4. Run the validation suite locally (see [Development Setup](docs/getting-started/README.md))
5. Commit following [commit conventions](#commit-conventions)
6. Push and open a Pull Request

---

## Commit Conventions

We follow **Conventional Commits**. Each commit message should be:

```
<type>(<scope>): <short description>

[optional body]

[optional footer]
```

**Types:**

| Type       | When to use                           |
|------------|---------------------------------------|
| `feat`     | New feature or behavior               |
| `fix`      | Bug fix                               |
| `refactor` | Code change with no behavior change   |
| `test`     | Adding or updating tests              |
| `docs`     | Documentation only                    |
| `chore`    | Tooling, dependencies, config         |
| `perf`     | Performance improvement               |
| `ci`       | CI/CD changes                         |

**Scope** (optional): the module or area affected, e.g. `auth`, `api`, `ui`, `docker`.

**Examples:**
```
feat(auth): add refresh token revocation on logout
fix(api): return 409 when resource already exists
refactor(core): extract validation to separate utility
test(auth): add unit tests for login use case
docs(setup): add environment variable reference
chore(deps): upgrade dependencies
```

**Rules:**
- Use the **imperative mood** ("add" not "adds" or "added")
- Keep the first line under **72 characters**
- Reference issues in the footer: `Closes #42`, `Fixes #17`

---

## Pull Request Process

1. **One PR per concern** — don't mix unrelated changes
2. **Fill the PR template** — describe what changed and why
3. **Keep diffs small** — large PRs are hard to review; split if needed
4. **All CI checks must pass** before merging
5. **Address review comments** — don't push force-merges; iterate on feedback
6. **Squash or rebase** before merge if history is messy

PRs are merged by maintainers once they have one approving review and all checks are green.

---

## Code Style

See [Code Style Guide](docs/developer-guide/guidelines/code-style.md) for detailed conventions.

---

## Testing

Every PR should maintain or improve the existing test coverage.

---

## Questions?

Open a [Discussion](../../discussions) or check the [FAQ](docs/developer-guide/faq.md).
