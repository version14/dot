# Generator: `prettier_frontend_rules`

Merges opinionated frontend Prettier rules into `.prettierrc` — JSX single quotes, 80-column print width, trailing commas, etc.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `prettier_frontend_rules` |
| Version | `0.1.0` |
| Package | `generators/prettier_frontend_rules` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `prettier_config` | Creates `.prettierrc` that this generator merges into |

---

## Answers consumed

None.

---

## Files written

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `.prettierrc` | `semi`, `singleQuote`, `jsxSingleQuote`, `trailingComma`, `printWidth`, `tabWidth`, `bracketSpacing`, `bracketSameLine` |

---

## Validators

None.

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

- [docs/flows/frontend.md](../flows/frontend.md)
