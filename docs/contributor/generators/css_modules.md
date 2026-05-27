# Generator: `css_modules`

Sets up CSS Modules with a global stylesheet and an example module file.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `css_modules` |
| Version | `0.1.0` |
| Package | `generators/css_modules` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Provides the base project structure |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/styles/global.css` | Global stylesheet (resets, variables) |
| `src/styles/App.module.css` | Example CSS Module |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/styles/App.module.css` | `file_exists` | File is present after generation |

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `tailwind_v4` | Conflicting styling approaches |
| `panda_css` | Conflicting styling approaches |
| `shadcn_ui` | shadcn manages its own CSS setup |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
