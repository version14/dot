# Generator: `version14_ui`

Adds `@version14/ui` — the company's WIP component library built with Panda CSS — and writes an example component.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `version14_ui` |
| Version | `0.2.1` |
| Package | `generators/version14_ui` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/components/V14Example.tsx` | Example component using `@version14/ui` |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `@version14/ui@latest` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/components/V14Example.tsx` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs @version14/ui from npm |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `shadcn_ui` | Both provide a component library |
| `ark_ui` | Both provide a component library |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
