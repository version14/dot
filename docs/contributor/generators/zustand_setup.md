# Generator: `zustand_setup`

Adds Zustand v5 state management with a typed counter store example.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `zustand_setup` |
| Version | `0.2.0` |
| Package | `generators/zustand_setup` |

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
| `src/stores/counter.store.ts` | Typed counter store with `CounterState` interface |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `zustand^5` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/stores/counter.store.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Zustand |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `jotai_setup` | Both provide global state management |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
