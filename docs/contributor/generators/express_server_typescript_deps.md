# Generator: `express_server_typescript_deps`

Merges Express, CORS, dotenv runtime dependencies plus types, tsx, nodemon devDependencies, and `dev`/`build`/`start` scripts into `package.json`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `express_server_typescript_deps` |
| Version | `0.1.0` |
| Package | `generators/express_server_typescript_deps` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `express_server_entrypoint` | Source files that reference express/cors/dotenv must exist first |

---

## Answers consumed

None.

---

## Files written

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `scripts.dev`, `scripts.build`, `scripts.start`, `dependencies.express`, `dependencies.cors`, `dependencies.dotenv`, `devDependencies.@types/express`, `devDependencies.@types/cors`, `devDependencies.@types/node`, `devDependencies.tsx`, `devDependencies.nodemon` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `dependencies.express` in `package.json` | `json_key_exists` | — |
| `scripts.dev` in `package.json` | `json_key_exists` | — |
| `scripts.build` in `package.json` | `json_key_exists` | — |
| `scripts.start` in `package.json` | `json_key_exists` | — |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install` | project root | Deduped with other generators |

## Test commands

No TestCommands.

---

## Conflicts

None.
