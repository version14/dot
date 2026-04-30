# Generator: `prettier_typescript_deps`

Merges prettier devDependency and `format`/`format:check` scripts into `package.json` for TypeScript projects.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `prettier_typescript_deps` |
| Version | `0.1.0` |
| Package | `generators/prettier_typescript_deps` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `prettier_config` | `.prettierrc` must exist before adding the dep |

---

## Answers consumed

None.

---

## Files written

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `scripts.format`, `scripts.format:check`, `devDependencies.prettier` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `devDependencies.prettier` in `package.json` | `json_key_exists` | Dep present |
| `scripts.format` in `package.json` | `json_key_exists` | Script present |

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
