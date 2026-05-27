# Flow: `frontend`

Scaffolds a TypeScript frontend project — React+Vite or Next.js — with a full range of optional add-ons: router, UI library, styling, state management, testing, auth, feature flags, Sentry, analytics, and SEO.

---

## Identity

| Field | Value |
|-------|-------|
| ID | `frontend` |
| Title | Frontend Wizard |
| File | `flows/frontend.go` |
| Root question | `project_name` |

---

## Questions

| ID | Type | Label | Options / Default |
|----|------|-------|-------------------|
| `project_name` | Text | "Project name" | Default: `"my-app"` |
| `framework` | Option | "Framework" | `react-vite`, `next` |
| `frontend-router` | Option | "Router" | `react-router`, `tanstack-router`, `none` — React+Vite only |
| `ui-library` | Option | "UI library" | `shadcn`, `arkui`, `version14`, `none` |
| `frontend-styling` | Option | "Styling" | `tailwind`, `css-modules`, `panda-css` — skipped when shadcn |
| `frontend-state` | Option | "State management" | `zustand`, `jotai`, `none` |
| `frontend-formatter` | Option | "Formatter" | `biome`, `prettier` |
| `frontend-linter` | Option | "Linter" | `biome`, `prettier` |
| `check-vitest-available` | If | _(silent: fw == react-vite?)_ | Routes to `include-vitest` or `include-playwright` |
| `include-vitest` | Confirm | "Include Vitest + React Testing Library?" | Default: `false` — React+Vite only |
| `include-playwright` | Confirm | "Include Playwright (E2E tests)?" | Default: `false` |
| `include-storybook` | Confirm | "Include Storybook?" | Default: `false` |
| `include-auth` | Confirm | "Include authentication?" | Default: `false` |
| `auth-provider` | Option | "Auth provider" | `clerk`, `better-auth`, `vanilla` — when `include-auth` = true |
| `include-theme` | Confirm | "Include custom theme provider?" | Default: `false` |
| `include-feature-flags` | Confirm | "Include feature flags?" | Default: `false` |
| `feature-flags-provider` | Option | "Feature flags provider" | `posthog`, `vercel`, `local` — when `include-feature-flags` = true |
| `include-sentry` | Confirm | "Include Sentry error tracking?" | Default: `false` |
| `include-analytics` | Confirm | "Include analytics?" | Default: `false` |
| `analytics-provider` | Option | "Analytics provider" | `ga4`, `plausible`, `posthog` — when `include-analytics` = true |
| `include-seo` | Confirm | "Include SEO setup?" | Default: `false` |
| `confirm-generate` | Confirm | "Generate the project now?" | Default: `true` |

---

## Question graph

```
project_name
  └── framework
        ├── [react-vite] → frontend-router
        │                    └── ui-library
        │                          ├── [shadcn] → frontend-state (skips styling)
        │                          └── [other]  → frontend-styling
        │                                            └── frontend-state
        └── [next] ──────────────── ui-library (same tree as above)

frontend-state
  └── frontend-formatter
        └── frontend-linter
              └── check-vitest-available (IfQuestion)
                    ├── [react-vite] → include-vitest
                    │                    └── include-playwright
                    └── [next]       → include-playwright

include-playwright
  └── include-storybook
        └── include-auth
              ├── [true]  → auth-provider
              │               └── include-theme
              └── [false] → include-theme

include-theme
  └── include-feature-flags
        ├── [true]  → feature-flags-provider
        │               └── include-sentry
        └── [false] → include-sentry

include-sentry
  └── include-analytics
        ├── [true]  → analytics-provider
        │               └── include-seo
        └── [false] → include-seo

include-seo
  └── confirm-generate
        └── (end)
```

---

## Generator resolution

| Condition | Generators added |
|-----------|-----------------|
| Always | `base_project`, `typescript_base` |
| `framework` = `react-vite` | `react_app` |
| `framework` = `next` | `nextjs_base` |
| `framework` = `react-vite` AND `frontend-router` = `react-router` | `react_router_v7` |
| `framework` = `react-vite` AND `frontend-router` = `tanstack-router` | `tanstack_router` |
| `ui-library` = `shadcn` | `shadcn_ui` |
| `ui-library` = `arkui` | `ark_ui` |
| `ui-library` = `version14` | `version14_ui` |
| `ui-library` ≠ `shadcn` AND `frontend-styling` = `tailwind` | `tailwind_v4` |
| `ui-library` ≠ `shadcn` AND `frontend-styling` = `css-modules` | `css_modules` |
| `ui-library` ≠ `shadcn` AND `frontend-styling` = `panda-css` | `panda_css` |
| `frontend-state` = `zustand` | `zustand_setup` |
| `frontend-state` = `jotai` | `jotai_setup` |
| `include-vitest` = true AND `framework` = `react-vite` | `vitest_testing_library` |
| `include-playwright` = true | `playwright_setup` |
| `include-storybook` = true | `storybook_setup` |
| `include-auth` = true AND `auth-provider` = `clerk` | `auth_clerk_frontend` |
| `include-auth` = true AND `auth-provider` = `better-auth` | `auth_better_auth_frontend` |
| `include-auth` = true AND `auth-provider` = `vanilla` | `auth_vanilla_frontend` |
| `include-theme` = true | `theme_provider` |
| `include-feature-flags` = true AND `feature-flags-provider` = `posthog` | `feature_flags_posthog` |
| `include-feature-flags` = true AND `feature-flags-provider` = `vercel` | `feature_flags_vercel` |
| `include-feature-flags` = true AND `feature-flags-provider` = `local` | `feature_flags_local` |
| `include-sentry` = true | `sentry_frontend` |
| `include-analytics` = true AND `analytics-provider` = `ga4` | `analytics_ga4` |
| `include-analytics` = true AND `analytics-provider` = `plausible` | `analytics_plausible` |
| `include-analytics` = true AND `analytics-provider` = `posthog` AND NOT already `feature_flags_posthog` | `feature_flags_posthog` (dedup) |
| `frontend-formatter` = `prettier` (last) | `prettier_config`, `prettier_typescript_deps`, `prettier_frontend_rules` |
| `frontend-formatter` = `biome` (last) | `biome_config` |

---

## Fixture examples

**Minimal React+Vite** (`tools/test-flow/testdata/202605280001_frontend_react_vite_minimal.json`):

```json
{
  "name": "frontend_react_vite_minimal",
  "flow_id": "frontend",
  "answers": {
    "project_name": "my-app",
    "framework": "react-vite",
    "frontend-router": "none",
    "ui-library": "none",
    "frontend-styling": "tailwind",
    "frontend-state": "none",
    "frontend-formatter": "biome",
    "frontend-linter": "biome",
    "include-vitest": false,
    "include-playwright": false,
    "include-storybook": false,
    "include-auth": false,
    "include-theme": false,
    "include-feature-flags": false,
    "include-sentry": false,
    "include-analytics": false,
    "include-seo": false,
    "confirm-generate": true
  }
}
```

**All modules** (`tools/test-flow/testdata/202605280008_frontend_react_vite_all_modules.json`):
shadcn forces Tailwind — the `frontend-styling` question is skipped; PostHog dedup — selecting PostHog for both feature flags and analytics emits `feature_flags_posthog` once.

---

## Source

`flows/frontend.go`

## See also

- [docs/generators/nextjs_base.md](../generators/nextjs_base.md)
- [docs/generators/react_router_v7.md](../generators/react_router_v7.md)
- [docs/generators/shadcn_ui.md](../generators/shadcn_ui.md)
