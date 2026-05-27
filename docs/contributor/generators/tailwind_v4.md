# Generator: `tailwind_v4`

Configures Tailwind CSS v4 — framework-aware: uses the Vite plugin for React+Vite and the PostCSS adapter for Next.js.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `tailwind_v4` |
| Version | `0.1.0` |
| Package | `generators/tailwind_v4` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `framework` | string | No | `"next"` → uses `@tailwindcss/postcss` + `postcss.config.mjs`; otherwise uses `@tailwindcss/vite` Vite plugin |

---

## Files written

| Path | Description |
|------|-------------|
| `src/styles/globals.css` | Global CSS with `@import "tailwindcss"` |
| `postcss.config.mjs` | PostCSS config (Next.js only) |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `devDependencies`: `tailwindcss` + either `@tailwindcss/vite` or `@tailwindcss/postcss` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/styles/globals.css` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Tailwind CSS v4 |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `panda_css` | Conflicting CSS-in-JS / utility-first systems |
| `shadcn_ui` | shadcn manages its own Tailwind setup |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
