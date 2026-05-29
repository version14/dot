# Generator: `playwright_setup`

Adds Playwright E2E testing — framework-aware dev server URL (5173 for Vite, 3000 for Next.js) and an example spec.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `playwright_setup` |
| Version | `0.2.0` |
| Package | `generators/playwright_setup` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `framework` | string | No | `"next"` → dev server URL `http://localhost:3000`; otherwise `http://localhost:5173` |

---

## Files written

| Path | Description |
|------|-------------|
| `playwright.config.ts` | Playwright config with dev server URL |
| `e2e/example.spec.ts` | Example homepage navigation test |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `devDependencies`: `@playwright/test` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `playwright.config.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Playwright package |
| `pnpm exec playwright install --with-deps chromium` | project root | Downloads Chromium browser |

## Test commands

No TestCommands.

---

## Conflicts

None.

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
