# Code Style Guide

<!-- Replace this file with your project's actual conventions. The sections below are a starting point. -->

This document defines the code style and formatting conventions for this project. Consistency matters more than any individual rule — when in doubt, follow existing patterns in the codebase.

---

## Tooling

<!-- Describe the linter / formatter used and how to run them. -->

```bash
# Check for style errors
# Auto-fix formatting
```

All checks must pass before a PR can be merged. The CI pipeline enforces this automatically.

---

## General Principles

- **Clarity over cleverness** — write code for the next reader, not the compiler
- **Explicit over implicit** — avoid magic; name things for what they do
- **Small functions** — each function should do one thing
- **No dead code** — remove commented-out code before committing

---

## Naming Conventions

<!-- Adapt the table below to your language and framework. -->

| Element | Convention | Example |
|---------|------------|---------|
| Files | `kebab-case` | `user-service.ts` |
| Classes | `PascalCase` | `UserService` |
| Functions / methods | `camelCase` | `getUserById` |
| Constants | `UPPER_SNAKE_CASE` | `MAX_RETRY_COUNT` |
| Private members | `_prefixed` or language convention | `_cache` |

---

## Formatting

<!-- Describe spacing, indentation, line length, quotes, etc. -->

- **Indentation**: X spaces (no tabs)
- **Max line length**: 100 characters
- **Trailing commas**: yes / no
- **Quotes**: single / double

---

## Import / Dependency Order

<!-- Describe how imports should be ordered. Example for languages with explicit imports: -->

1. Standard library / built-ins
2. Third-party packages
3. Internal packages / modules
4. Relative imports

Separate each group with a blank line.

---

## Error Handling

- Never swallow errors silently — log or propagate
- Use typed/structured errors where the language supports it
- Validate inputs at system boundaries; trust internal code

---

## Testing Conventions

- Tests live next to the code they test, or in a dedicated `tests/` directory
- Each test should be independent and not rely on shared mutable state
- Test names should describe the expected behavior: `should return 404 when user not found`

---

## Running the Full Validation Suite

```bash
# Replace with your actual commands
# Lint
# Type-check (if applicable)
# Tests
# Build
```

All four must pass before pushing.
