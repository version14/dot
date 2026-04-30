# Generator: `prettier_config`

Creates the base Prettier configuration files (`.prettierrc` and `.prettierignore`). Downstream generators merge additional rules on top.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `prettier_config` |
| Version | `0.1.0` |
| Package | `generators/prettier_config` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires `package.json` to exist before any tooling is added |

---

## Answers consumed

None. Writes the same baseline config for every project.

---

## Files written

| Path | Description |
|------|-------------|
| `.prettierrc` | Base rules: semi, singleQuote, trailingComma, tabWidth |
| `.prettierignore` | Ignores node_modules, dist, build, .env |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `.prettierrc` | `file_exists` | File is present after generation |

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

None.

---

## See also

- [`prettier_typescript_deps`](prettier_typescript_deps.md) — adds prettier devDep + format script
- [`prettier_express_rules`](prettier_express_rules.md) — backend-specific rule overrides
