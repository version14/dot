# Generator: `shadcn_ui`

Configures shadcn/ui with Tailwind v4 — writes components.json, a button component, and the global CSS entry point. Implicitly replaces the `tailwind_v4` generator (they conflict).

---

## Identity

| Field | Value |
|-------|-------|
| Name | `shadcn_ui` |
| Version | `0.1.0` |
| Package | `generators/shadcn_ui` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `framework` | string | No | `"next"` → writes CSS to `src/app/globals.css`; otherwise `src/styles/globals.css` |

---

## Files written

| Path | Description |
|------|-------------|
| `components.json` | shadcn/ui configuration |
| `src/app/globals.css` or `src/styles/globals.css` | Global CSS with `@import "tailwindcss"` (path depends on framework) |
| `src/lib/utils.ts` | `cn()` helper using `clsx` + `tailwind-merge` |
| `src/components/ui/button.tsx` | Example Button component |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `tailwindcss`, `clsx`, `tailwind-merge`; `devDependencies`: `@tailwindcss/vite` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `components.json` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Tailwind and shadcn deps |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `tailwind_v4` | Both provide Tailwind CSS setup |
| `panda_css` | Conflicting CSS-in-JS / utility-first systems |
| `css_modules` | Conflicting styling approaches |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
