# Generator: `sentry_frontend`

Adds Sentry error tracking — framework-aware: `@sentry/nextjs` with a client config file for Next.js, `@sentry/react` for Vite.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `sentry_frontend` |
| Version | `0.3.0` |
| Package | `generators/sentry_frontend` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `framework` | string | No | `"next"` → installs `@sentry/nextjs`, writes `sentry.client.config.ts`; otherwise `@sentry/react` |

---

## Files written

| Path | Description |
|------|-------------|
| `src/lib/sentry.ts` | Sentry init helper |
| `sentry.client.config.ts` | Next.js client-side Sentry configuration (Next.js only) |
| `.env.example` | Required env vars (`VITE_SENTRY_DSN` or `NEXT_PUBLIC_SENTRY_DSN`) |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `@sentry/nextjs` or `@sentry/react` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/lib/sentry.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Sentry SDK |

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
