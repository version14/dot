# Generator: `jotai_setup`

Adds Jotai v2 atomic state management with a counter atom example.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `jotai_setup` |
| Version | `0.2.0` |
| Package | `generators/jotai_setup` |

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
| `src/atoms/counter.atom.ts` | Counter atom using Jotai's `atom()` |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `jotai^2` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/atoms/counter.atom.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Jotai |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `zustand_setup` | Both provide global state management |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
