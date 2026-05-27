# Generator: `seo_next`

Adds `next-seo` — writes a default SEO config and a `<SEO>` component wrapper.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `seo_next` |
| Version | `0.1.0` |
| Package | `generators/seo_next` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `nextjs_base` | next-seo is a Next.js-specific library |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/lib/seo.config.ts` | Default SEO configuration (`DefaultSeoProps`) |
| `src/components/SEO.tsx` | Thin wrapper around `<NextSeo>` |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `next-seo^6` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/components/SEO.tsx` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs next-seo |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `seo_react` | Both write `src/components/SEO.tsx` |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
