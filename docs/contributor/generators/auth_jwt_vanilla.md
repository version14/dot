# Generator: `auth_jwt_vanilla`

Vanilla JWT authentication. Creates `src/lib/jwt.ts` (sign/verify helpers) and `src/middleware/auth.middleware.ts` (Bearer token guard). Adds `jsonwebtoken` + `@types/jsonwebtoken` to `package.json` and appends `JWT_SECRET`/`JWT_EXPIRES_IN` to `.env.example`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `auth_jwt_vanilla` |
| Version | `0.1.0` |
| Package | `generators/auth_jwt_vanilla` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `express_server_entrypoint` | `.env.example` and `src/` directory must exist |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/lib/jwt.ts` | `signToken` and `verifyToken<T>` helpers backed by `process.env.JWT_SECRET` |
| `src/middleware/auth.middleware.ts` | Express middleware that validates `Authorization: Bearer <token>` headers |
| `.env.example` | Appends `JWT_SECRET` and `JWT_EXPIRES_IN` |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies.jsonwebtoken`, `devDependencies.@types/jsonwebtoken` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/lib/jwt.ts` | `file_exists` | — |
| `src/middleware/auth.middleware.ts` | `file_exists` | — |
| `dependencies.jsonwebtoken` in `package.json` | `json_key_exists` | — |

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
