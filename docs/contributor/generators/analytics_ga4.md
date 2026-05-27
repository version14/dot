# Generator: `analytics_ga4`

Adds Google Analytics 4 via `react-ga4` — writes an initializer and a `.env.example` with the measurement ID.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `analytics_ga4` |
| Version | `0.1.0` |
| Package | `generators/analytics_ga4` |

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
| `src/lib/ga4.ts` | GA4 init and `trackEvent` helper |
| `.env.example` | Required env var (`VITE_GA_MEASUREMENT_ID`) |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `react-ga4^2` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/lib/ga4.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs react-ga4 |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `analytics_plausible` | Both provide analytics — only one analytics library should be installed |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
