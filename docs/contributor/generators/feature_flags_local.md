# Generator: `feature_flags_local`

Adds local JSON-based feature flags — a static `public/flags.json` file, a flags reader, and a `useFeatureFlag` hook. No external service required.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `feature_flags_local` |
| Version | `0.1.0` |
| Package | `generators/feature_flags_local` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires the base project structure |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `public/flags.json` | Static flag definitions (key → boolean) |
| `src/lib/flags.ts` | Async loader that fetches `/flags.json` at runtime |
| `src/hooks/useFeatureFlag.ts` | Hook for reading a flag value |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/lib/flags.ts` | `file_exists` | File is present after generation |

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
