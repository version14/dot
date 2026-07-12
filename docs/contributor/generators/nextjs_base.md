# Generator: `nextjs_base`

Bootstraps a Next.js 15 App Router project with TypeScript, writing the entry layout, home page, global CSS, and next.config.ts.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `nextjs_base` |
| Version | `0.3.1` |
| Package | `generators/nextjs_base` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires tsconfig.json to patch with Next.js-specific compiler options |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `next.config.ts` | Minimal Next.js config |
| `src/app/layout.tsx` | Root layout component |
| `src/app/page.tsx` | Home page |
| `src/app/globals.css` | Global stylesheet entry point |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `next`, `react`, `react-dom` |
| `tsconfig.json` | `compilerOptions.jsx`, `compilerOptions.plugins`, `compilerOptions.paths` |

---

## Validators

None.

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Next.js and React |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

None.

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
