# Generator: `drizzle_typescript_deps`

Adds `drizzle-orm` (runtime) and `drizzle-kit` (devDep) to `package.json` and registers `db:generate`, `db:migrate`, `db:push`, `db:studio` scripts.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `drizzle_typescript_deps` |
| Version | `0.1.0` |
| Package | `generators/drizzle_typescript_deps` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `drizzle_config_base` | Config file must exist before adding the tooling deps |

---

## Answers consumed

None.

---

## Files written

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies.drizzle-orm`, `devDependencies.drizzle-kit`, `scripts.db:generate`, `scripts.db:migrate`, `scripts.db:push`, `scripts.db:studio` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `dependencies.drizzle-orm` in `package.json` | `json_key_exists` | — |
| `devDependencies.drizzle-kit` in `package.json` | `json_key_exists` | — |
| `scripts.db:push` in `package.json` | `json_key_exists` | — |

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
