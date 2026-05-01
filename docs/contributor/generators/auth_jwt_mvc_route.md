# Generator: `auth_jwt_mvc_route`

JWT auth route and controller for MVC architecture. Generates `src/routes/auth.route.ts` (register/login/me/refresh/logout) and `src/controllers/auth.controller.ts`. The controller is fully implemented when a Drizzle adapter has been generated; otherwise it returns 501 stubs.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `auth_jwt_mvc_route` |
| Version | `0.1.0` |
| Package | `generators/auth_jwt_mvc_route` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `auth_jwt_vanilla` | `src/shared/services/jwt.ts` and `src/shared/middlewares/auth.middleware.ts` must exist |

---

## Answers consumed

None (reads `ctx.PreviousGens` at runtime to detect Drizzle).

---

## Files written

| Path | Description |
|------|-------------|
| `src/routes/auth.route.ts` | Express router wiring POST /register, POST /login, GET /me, POST /refresh, POST /logout |
| `src/controllers/auth.controller.ts` | Full controller with bcrypt + DB operations when `drizzle_postgres_adapter` ran; 501 stubs otherwise |

Also merges into (when Drizzle is present):

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies.bcryptjs`, `devDependencies.@types/bcryptjs` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/routes/auth.route.ts` | `file_exists` | — |
| `src/controllers/auth.controller.ts` | `file_exists` | — |

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
