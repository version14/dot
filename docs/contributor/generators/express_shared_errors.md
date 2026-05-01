# Generator: `express_shared_errors`

Shared error classes for Express. Creates a hierarchy rooted at `AppError` — `NotFoundError`, `ValidationError`, `UnauthorizedError`, `ForbiddenError`, `ConflictError` — all exported from `src/shared/errors/index.ts`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `express_shared_errors` |
| Version | `0.1.0` |
| Package | `generators/express_shared_errors` |

---

## Dependencies

None.

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/shared/errors/app.error.ts` | Base `AppError` class (message, statusCode, code) |
| `src/shared/errors/not-found.error.ts` | `NotFoundError` — HTTP 404 |
| `src/shared/errors/validation.error.ts` | `ValidationError` — HTTP 400 |
| `src/shared/errors/unauthorized.error.ts` | `UnauthorizedError` — HTTP 401 |
| `src/shared/errors/forbidden.error.ts` | `ForbiddenError` — HTTP 403 |
| `src/shared/errors/conflict.error.ts` | `ConflictError` — HTTP 409 |
| `src/shared/errors/index.ts` | Re-exports all error classes |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/shared/errors/app.error.ts` | `file_exists` | — |
| `src/shared/errors/not-found.error.ts` | `file_exists` | — |
| `src/shared/errors/validation.error.ts` | `file_exists` | — |
| `src/shared/errors/unauthorized.error.ts` | `file_exists` | — |
| `src/shared/errors/forbidden.error.ts` | `file_exists` | — |
| `src/shared/errors/conflict.error.ts` | `file_exists` | — |
| `src/shared/errors/index.ts` | `file_exists` | — |

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

None.

## See also

- [express_error_middleware.md](./express_error_middleware.md) — global error handler that uses `AppError`
