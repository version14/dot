# Generator: `react_app`

React + Vite application. Merges React dependencies and JSX compiler options into the existing `package.json` and `tsconfig.json` created by `typescript_base`, then writes the application entry files.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `react_app` |
| Version | `0.1.0` |
| Package | `generators/react_app` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires `package.json` and `tsconfig.json` to exist for merging |

`base_project` is also in the chain via `typescript_base → base_project`.

---

## Answers consumed

| Key | Type | Notes |
|-----|------|-------|
| `project_name` | string | Written as the `<title>` in `index.html` and the app heading in `src/App.tsx`. Falls back to `spec.Metadata.ProjectName`, then `"app"`. |

---

## Files written

| Path | Description |
|------|-------------|
| `index.html` | HTML entry point with `<div id="root">` and Vite script tag |
| `vite.config.ts` | Vite config with `@vitejs/plugin-react` |
| `src/main.tsx` | React DOM `createRoot` entry |
| `src/App.tsx` | Minimal `App` component returning a heading |

Also merges into existing files:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `scripts.dev/build/preview`, `dependencies.react/react-dom`, `devDependencies.@types/react/react-dom/@vitejs/plugin-react/vite` |
| `tsconfig.json` | `compilerOptions.jsx = "react-jsx"`, `compilerOptions.lib` (adds DOM entries) |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/main.tsx` exists | `file_exists` | — |
| `src/App.tsx` exists | `file_exists` | — |
| `index.html` exists | `file_exists` | — |
| `vite.config.ts` exists | `file_exists` | — |
| `dependencies.react` in `package.json` | `json_key_exists` | React dep is present |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install` | project root | Deduplicated with `typescript_base` |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm install` | No | — | Deduplicated |
| `pnpm exec tsc --noEmit` | No | — | Type-check |
| `pnpm exec vite build` | No | — | Production build |
| `pnpm exec vite` | **Yes** | 4s | Smoke-start the dev server |

---

## Conflicts

None.
