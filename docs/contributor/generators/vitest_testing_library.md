# Generator: `vitest_testing_library`

Adds Vitest v2 + React Testing Library for unit and component testing in React+Vite projects.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `vitest_testing_library` |
| Version | `0.2.1` |
| Package | `generators/vitest_testing_library` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `react_app` | React+Vite-only — requires the Vite setup |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `vitest.config.ts` | Vitest configuration with jsdom environment and setup file |
| `src/test/setup.ts` | Global test setup (imports `@testing-library/jest-dom`) |
| `src/test/App.test.tsx` | Example component test |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `devDependencies`: `vitest^2`, `@testing-library/react`, `@testing-library/jest-dom`, `@testing-library/user-event`, `jsdom` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `vitest.config.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Vitest and Testing Library |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec vitest run` | No | — | Runs the test suite once (non-watch) |

---

## Conflicts

None.

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
