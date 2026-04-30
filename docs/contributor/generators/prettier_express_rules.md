# Generator: `prettier_express_rules`

Overlays Express/Node.js backend-specific Prettier rules on top of the base `.prettierrc` (printWidth 100, lf line endings, bracket spacing).

---

## Identity

| Field | Value |
|-------|-------|
| Name | `prettier_express_rules` |
| Version | `0.1.0` |
| Package | `generators/prettier_express_rules` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `prettier_config` | `.prettierrc` must exist before rules are merged in |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `.prettierrc` | Merges: `printWidth: 100`, `endOfLine: "lf"`, `bracketSpacing: true`, `arrowParens: "always"` |

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
