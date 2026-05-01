# Generator: `auth_jwt_users_schema`

Drizzle schema for JWT authentication. Creates `users` and `refresh_tokens` tables in `src/db/schema/users.table.ts` and re-exports them from the schema index.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `auth_jwt_users_schema` |
| Version | `0.1.0` |
| Package | `generators/auth_jwt_users_schema` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `drizzle_postgres_adapter` | Drizzle + Postgres must be configured; the schema index at `src/db/schema/index.ts` must exist |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/db/schema/users.table.ts` | `users` table (id, email, passwordHash, timestamps) and `refresh_tokens` table (id, token, userId FK, expiresAt) with inferred types |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `src/db/schema/index.ts` | Appends `export * from './users.table';` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/db/schema/users.table.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

None.
