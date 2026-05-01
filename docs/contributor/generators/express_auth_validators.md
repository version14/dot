# Generator: `express_auth_validators`

Zod validation schemas for auth endpoints. Creates `src/shared/validators/auth.validators.ts` with `registerSchema`, `loginSchema`, and `refreshSchema`, plus their inferred TypeScript types. Adds `zod` to `package.json`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `express_auth_validators` |
| Version | `0.1.0` |
| Package | `generators/express_auth_validators` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `express_server_entrypoint` | `package.json` must exist for the Zod dependency merge |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/shared/validators/auth.validators.ts` | `registerSchema` (email + min-8 password), `loginSchema`, `refreshSchema` with inferred types |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies.zod` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/shared/validators/auth.validators.ts` | `file_exists` | — |
| `dependencies.zod` in `package.json` | `json_key_exists` | — |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install` | project root | Installs `zod` |

## Test commands

No TestCommands.

---

## Conflicts

None.
