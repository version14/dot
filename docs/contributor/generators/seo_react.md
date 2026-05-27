# Generator: `seo_react`

Adds `react-helmet-async` SEO for React+Vite — writes a `<HelmetProvider>` wrapper and a reusable `<SEO>` component.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `seo_react` |
| Version | `0.1.0` |
| Package | `generators/seo_react` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `react_app` | react-helmet-async is a React+Vite-specific integration |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/providers/HelmetProvider.tsx` | Thin wrapper around `react-helmet-async`'s `HelmetProvider` |
| `src/components/SEO.tsx` | `<SEO>` component for per-page title, description, canonical |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `react-helmet-async^2` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/components/SEO.tsx` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs react-helmet-async |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `seo_next` | Both write `src/components/SEO.tsx` |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
