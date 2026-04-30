# Generator: `auth_better_auth`

BetterAuth session-based authentication setup. Creates `src/lib/auth.ts` (auth instance with Drizzle adapter) and `src/routes/auth.route.ts` (catch-all handler for `/api/auth/*`). Adds `better-auth` to dependencies and appends `BETTER_AUTH_SECRET` and `BETTER_AUTH_URL` to `.env.example`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `auth_better_auth` |
| Version | `0.1.0` |
| Package | `generators/auth_better_auth` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `drizzle_postgres_adapter` | BetterAuth uses the Drizzle adapter which requires an active `db` export |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/lib/auth.ts` | BetterAuth instance with Drizzle PG adapter and email/password enabled |
| `src/routes/auth.route.ts` | Express router that forwards all `/api/auth/*` requests to BetterAuth |
| `.env.example` | Appends `BETTER_AUTH_SECRET` and `BETTER_AUTH_URL` |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies.better-auth` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/lib/auth.ts` | `file_exists` | — |
| `dependencies.better-auth` in `package.json` | `json_key_exists` | — |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install` | project root | Deduped |

## Test commands

No TestCommands.

---

## Conflicts

None.
