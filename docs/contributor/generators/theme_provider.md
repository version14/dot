# Generator: `theme_provider`

Adds a custom CSS-variable-based theme system with light/dark mode support — a ThemeProvider, a `useTheme` hook, and a global CSS variables file.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `theme_provider` |
| Version | `0.1.0` |
| Package | `generators/theme_provider` |

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
| `src/styles/theme.css` | CSS custom properties for light and dark color schemes |
| `src/providers/ThemeProvider.tsx` | React provider — reads localStorage, listens to `prefers-color-scheme` |
| `src/hooks/useTheme.ts` | Hook returning current theme and a toggle function |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/providers/ThemeProvider.tsx` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs base deps |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

None.

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
