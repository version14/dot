# Generator: `feature_flags_vercel`

Adds Vercel Edge Config feature flags — installs `@vercel/edge-config` and `@vercel/flags`, writes a flags client and a `useFeatureFlag` hook.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `feature_flags_vercel` |
| Version | `0.2.0` |
| Package | `generators/feature_flags_vercel` |

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
| `src/lib/flags.ts` | Edge Config client + flag helpers |
| `src/hooks/useFeatureFlag.ts` | Hook for reading a flag value |
| `.env.example` | Required env vars (`EDGE_CONFIG`, `FLAGS_SECRET`) |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `@vercel/edge-config`, `@vercel/flags` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/lib/flags.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Vercel Edge Config packages |

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
