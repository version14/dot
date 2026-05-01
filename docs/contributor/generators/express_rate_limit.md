# Generator: `express_rate_limit`

Adds `express-rate-limit` with a 100 req / 15-minute window applied globally. Creates `src/shared/middlewares/rate-limit.middleware.ts` and injects the limiter into `src/app.ts`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `express_rate_limit` |
| Version | `0.1.0` |
| Package | `generators/express_rate_limit` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `express_server_entrypoint` | `src/app.ts` must exist; `package.json` must exist |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/shared/middlewares/rate-limit.middleware.ts` | `limiter` — `rateLimit` config with `draft-7` standard headers |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies.express-rate-limit` |

Also modifies:

| Path | Change |
|------|--------|
| `src/app.ts` | Prepends the import; inserts `app.use(limiter)` before `app.use(express.urlencoded(...))` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/shared/middlewares/rate-limit.middleware.ts` | `file_exists` | — |
| `dependencies.express-rate-limit` in `package.json` | `json_key_exists` | — |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install` | project root | Installs `express-rate-limit` |

## Test commands

No TestCommands.

---

## Conflicts

None.
