# Generator: `analytics_plausible`

Adds Plausible Analytics via `plausible-tracker` — writes a tracker initializer and a `.env.example` with the domain.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `analytics_plausible` |
| Version | `0.1.0` |
| Package | `generators/analytics_plausible` |

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
| `src/lib/plausible.ts` | Plausible tracker init and `trackEvent` helper |
| `.env.example` | Required env var (`VITE_PLAUSIBLE_DOMAIN`) |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `plausible-tracker^0.3` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/lib/plausible.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs plausible-tracker |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `analytics_ga4` | Both provide analytics — only one analytics library should be installed |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
