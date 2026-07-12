# Generator: `express_test_setup`

Vitest configuration and testing dependencies (vitest, supertest) for Express TypeScript projects.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `express_test_setup` |
| Version | `0.4.1` |
| Package | `generators/express_test_setup` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `express_server_entrypoint` | Tests target the Express app the entrypoint creates |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `vitest.config.ts` | Vitest configuration for the Express server |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `scripts`: `test`, `test:watch`, `test:coverage` |
| `package.json` | `devDependencies`: `vitest`, `@vitest/coverage-v8`, `supertest`, `@types/supertest` |

Dependency versions come from the [Catalog](../dep-checker.md) (`internal/deps`), not
from this generator. It names the packages; the Catalog pins them.

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `vitest.config.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Vitest and supertest |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec vitest run` | No | — | Runs the test suite once (non-watch) |

---

## Conflicts

None. Merges only into `package.json` and writes a file no other generator owns.
