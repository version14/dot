# Generator: `feature_flags_posthog`

Adds PostHog feature flags — installs `posthog-js`, writes a PostHog client initializer, a React provider, and a `.env.example`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `feature_flags_posthog` |
| Version | `0.2.0` |
| Package | `generators/feature_flags_posthog` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/lib/posthog.ts` | PostHog client initialization |
| `src/providers/PostHogProvider.tsx` | React provider wrapping `PostHogProvider` from posthog-js |
| `.env.example` | Required env vars (`VITE_POSTHOG_KEY`, `VITE_POSTHOG_HOST`) |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `posthog-js^1` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/lib/posthog.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs PostHog JS SDK |

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
- [docs/generators/analytics_plausible.md](analytics_plausible.md) — note dedup: selecting PostHog for both feature flags and analytics emits this generator once
