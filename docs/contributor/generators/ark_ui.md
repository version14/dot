# Generator: `ark_ui`

Adds Ark UI headless components — installs `@ark-ui/react` and writes an example button built on Ark's `Button` primitive.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `ark_ui` |
| Version | `0.1.0` |
| Package | `generators/ark_ui` |

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
| `src/components/ui/button.tsx` | Button built on Ark UI's Button primitive |
| `src/components/ui/index.ts` | Barrel export |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `@ark-ui/react` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/components/ui/button.tsx` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Ark UI |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `shadcn_ui` | Both write `src/components/ui/button.tsx` |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
