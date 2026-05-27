package flows

import (
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/spec"
)

const FRONTEND_REACT_VITE = "react-vite"
const FRONTEND_NEXT = "next"

const FRONTEND_ROUTER_REACT_ROUTER = "react-router"
const FRONTEND_ROUTER_TANSTACK = "tanstack-router"
const FRONTEND_ROUTER_NONE = "none"

const FRONTEND_UI_SHADCN = "shadcn"
const FRONTEND_UI_ARKUI = "arkui"
const FRONTEND_UI_VERSION14 = "version14"
const FRONTEND_UI_NONE = "none"

const FRONTEND_STYLING_TAILWIND = "tailwind"
const FRONTEND_STYLING_CSS_MODULES = "css-modules"
const FRONTEND_STYLING_PANDA = "panda-css"

const FRONTEND_STATE_ZUSTAND = "zustand"
const FRONTEND_STATE_JOTAI = "jotai"
const FRONTEND_STATE_NONE = "none"

const FRONTEND_AUTH_CLERK = "clerk"
const FRONTEND_AUTH_BETTER_AUTH = "better-auth"
const FRONTEND_AUTH_VANILLA = "vanilla"

const FRONTEND_FLAGS_POSTHOG = "posthog"
const FRONTEND_FLAGS_VERCEL = "vercel"
const FRONTEND_FLAGS_LOCAL = "local"

const FRONTEND_ANALYTICS_GA4 = "ga4"
const FRONTEND_ANALYTICS_PLAUSIBLE = "plausible"
const FRONTEND_ANALYTICS_POSTHOG = "posthog"

// FrontendFlow is the frontend project wizard. It walks the user through
// framework → router → UI library → styling → state → formatter/linter →
// testing → modules (auth, theme, feature flags, Sentry, analytics, SEO).
//
// Question IDs are stable: re-runs of `dot scaffold` reuse persisted answers.
func FrontendFlow() *FlowDef {
	confirmGenerate := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "confirm-generate"},
		Label:        "Generate the project now?",
		Default:      true,
		Then:         &flow.Next{End: true},
		Else:         &flow.Next{End: true},
	}

	includeSEO := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "include-seo"},
		Label:        "Include SEO setup?",
		Default:      false,
		Then:         &flow.Next{Question: confirmGenerate},
		Else:         &flow.Next{Question: confirmGenerate},
	}

	analyticsProvider := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "analytics-provider"},
		Label:        "Analytics provider",
		Options: []*flow.Option{
			{Label: "Google Analytics 4", Value: FRONTEND_ANALYTICS_GA4, Next: &flow.Next{Question: includeSEO}},
			{Label: "Plausible", Value: FRONTEND_ANALYTICS_PLAUSIBLE, Next: &flow.Next{Question: includeSEO}},
			{Label: "PostHog", Value: FRONTEND_ANALYTICS_POSTHOG, Next: &flow.Next{Question: includeSEO}},
		},
	}

	includeAnalytics := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "include-analytics"},
		Label:        "Include analytics?",
		Default:      false,
		Then:         &flow.Next{Question: analyticsProvider},
		Else:         &flow.Next{Question: includeSEO},
	}

	includeSentry := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "include-sentry"},
		Label:        "Include Sentry error tracking?",
		Default:      false,
		Then:         &flow.Next{Question: includeAnalytics},
		Else:         &flow.Next{Question: includeAnalytics},
	}

	featureFlagsProvider := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "feature-flags-provider"},
		Label:        "Feature flags provider",
		Options: []*flow.Option{
			{Label: "PostHog", Value: FRONTEND_FLAGS_POSTHOG, Next: &flow.Next{Question: includeSentry}},
			{Label: "Vercel Edge Config", Value: FRONTEND_FLAGS_VERCEL, Next: &flow.Next{Question: includeSentry}},
			{Label: "Local JSON", Value: FRONTEND_FLAGS_LOCAL, Next: &flow.Next{Question: includeSentry}},
		},
	}

	includeFeatureFlags := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "include-feature-flags"},
		Label:        "Include feature flags?",
		Default:      false,
		Then:         &flow.Next{Question: featureFlagsProvider},
		Else:         &flow.Next{Question: includeSentry},
	}

	includeTheme := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "include-theme"},
		Label:        "Include custom theme provider (CSS variables, light/dark)?",
		Default:      false,
		Then:         &flow.Next{Question: includeFeatureFlags},
		Else:         &flow.Next{Question: includeFeatureFlags},
	}

	authProvider := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "auth-provider"},
		Label:        "Auth provider",
		Options: []*flow.Option{
			{Label: "Clerk", Value: FRONTEND_AUTH_CLERK, Next: &flow.Next{Question: includeTheme}},
			{Label: "Better Auth", Value: FRONTEND_AUTH_BETTER_AUTH, Next: &flow.Next{Question: includeTheme}},
			{Label: "Vanilla (from scratch)", Value: FRONTEND_AUTH_VANILLA, Next: &flow.Next{Question: includeTheme}},
		},
	}

	includeAuth := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "include-auth"},
		Label:        "Include authentication?",
		Default:      false,
		Then:         &flow.Next{Question: authProvider},
		Else:         &flow.Next{Question: includeTheme},
	}

	includeStorybook := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "include-storybook"},
		Label:        "Include Storybook?",
		Default:      false,
		Then:         &flow.Next{Question: includeAuth},
		Else:         &flow.Next{Question: includeAuth},
	}

	includePlaywright := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "include-playwright"},
		Label:        "Include Playwright (E2E tests)?",
		Default:      false,
		Then:         &flow.Next{Question: includeStorybook},
		Else:         &flow.Next{Question: includeStorybook},
	}

	includeVitest := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "include-vitest"},
		Label:        "Include Vitest + React Testing Library?",
		Default:      false,
		Then:         &flow.Next{Question: includePlaywright},
		Else:         &flow.Next{Question: includePlaywright},
	}

	// IfQuestion: Vitest is React+Vite only — skip for Next.js
	checkVitest := &flow.IfQuestion{
		QuestionBase: flow.QuestionBase{ID_: "check-vitest-available"},
		Condition: func(ctx *flow.FlowContext) bool {
			fw, _ := ctx.Answers["framework"].(string)
			return fw == FRONTEND_REACT_VITE
		},
		Then: &flow.Next{Question: includeVitest},
		Else: &flow.Next{Question: includePlaywright},
	}

	linter := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "frontend-linter"},
		Label:        "Linter",
		Options: []*flow.Option{
			{Label: "Biome", Value: "biome", Next: &flow.Next{Question: checkVitest}},
			{Label: "Prettier", Value: "prettier", Next: &flow.Next{Question: checkVitest}},
		},
	}

	formatter := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "frontend-formatter"},
		Label:        "Formatter",
		Options: []*flow.Option{
			{Label: "Biome", Value: "biome", Next: &flow.Next{Question: linter}},
			{Label: "Prettier", Value: "prettier", Next: &flow.Next{Question: linter}},
		},
	}

	state := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "frontend-state"},
		Label:        "State management",
		Options: []*flow.Option{
			{Label: "Zustand", Value: FRONTEND_STATE_ZUSTAND, Next: &flow.Next{Question: formatter}},
			{Label: "Jotai", Value: FRONTEND_STATE_JOTAI, Next: &flow.Next{Question: formatter}},
			{Label: "None", Value: FRONTEND_STATE_NONE, Next: &flow.Next{Question: formatter}},
		},
	}

	styling := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "frontend-styling"},
		Label:        "Styling",
		Options: []*flow.Option{
			{Label: "Tailwind CSS v4", Value: FRONTEND_STYLING_TAILWIND, Next: &flow.Next{Question: state}},
			{Label: "CSS Modules", Value: FRONTEND_STYLING_CSS_MODULES, Next: &flow.Next{Question: state}},
			{Label: "Panda CSS", Value: FRONTEND_STYLING_PANDA, Next: &flow.Next{Question: state}},
		},
	}

	uiLibrary := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "ui-library"},
		Label:        "UI library",
		Options: []*flow.Option{
			// shadcn forces Tailwind — skip styling question
			{Label: "shadcn/ui (includes Tailwind v4)", Value: FRONTEND_UI_SHADCN, Next: &flow.Next{Question: state}},
			{Label: "Ark UI", Value: FRONTEND_UI_ARKUI, Next: &flow.Next{Question: styling}},
			{Label: "@version14/ui", Value: FRONTEND_UI_VERSION14, Next: &flow.Next{Question: styling}},
			{Label: "None", Value: FRONTEND_UI_NONE, Next: &flow.Next{Question: styling}},
		},
	}

	router := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "frontend-router"},
		Label:        "Router",
		Options: []*flow.Option{
			{Label: "React Router v7", Value: FRONTEND_ROUTER_REACT_ROUTER, Next: &flow.Next{Question: uiLibrary}},
			{Label: "TanStack Router", Value: FRONTEND_ROUTER_TANSTACK, Next: &flow.Next{Question: uiLibrary}},
			{Label: "None", Value: FRONTEND_ROUTER_NONE, Next: &flow.Next{Question: uiLibrary}},
		},
	}

	framework := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "framework"},
		Label:        "Framework",
		Options: []*flow.Option{
			// React+Vite: ask router first
			{Label: "React + Vite", Value: FRONTEND_REACT_VITE, Next: &flow.Next{Question: router}},
			// Next.js: skip router (built-in)
			{Label: "Next.js", Value: FRONTEND_NEXT, Next: &flow.Next{Question: uiLibrary}},
		},
	}

	projectName := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "project_name",
			Next_: &flow.Next{Question: framework},
		},
		Label:       "Project name",
		Description: "Used as the package name and root directory.",
		Default:     "my-app",
		Validate:    nonEmpty,
	}

	return &FlowDef{
		ID:          "frontend",
		Title:       "Frontend Wizard",
		Description: "Scaffold a TypeScript frontend project with React+Vite or Next.js.",
		Root:        projectName,
		Generators:  resolveFrontendGenerators,
	}
}

func resolveFrontendGenerators(s *spec.ProjectSpec) []Invocation {
	if s == nil {
		return nil
	}

	fw, _ := s.Answers["framework"].(string)
	routerChoice, _ := s.Answers["frontend-router"].(string)
	uiLib, _ := s.Answers["ui-library"].(string)
	stylingChoice, _ := s.Answers["frontend-styling"].(string)
	stateChoice, _ := s.Answers["frontend-state"].(string)
	formatterChoice, _ := s.Answers["frontend-formatter"].(string)
	includeVitest, _ := s.Answers["include-vitest"].(bool)
	includePlaywright, _ := s.Answers["include-playwright"].(bool)
	includeStorybook, _ := s.Answers["include-storybook"].(bool)
	includeAuth, _ := s.Answers["include-auth"].(bool)
	authProviderChoice, _ := s.Answers["auth-provider"].(string)
	includeTheme, _ := s.Answers["include-theme"].(bool)
	includeFeatureFlags, _ := s.Answers["include-feature-flags"].(bool)
	flagsProviderChoice, _ := s.Answers["feature-flags-provider"].(string)
	includeSentry, _ := s.Answers["include-sentry"].(bool)
	includeAnalytics, _ := s.Answers["include-analytics"].(bool)
	analyticsProviderChoice, _ := s.Answers["analytics-provider"].(string)
	includeSEO, _ := s.Answers["include-seo"].(bool)

	var out []Invocation

	out = append(out, Invocation{Name: "base_project"})
	out = append(out, Invocation{Name: "typescript_base"})

	// Framework base
	switch fw {
	case FRONTEND_REACT_VITE:
		out = append(out, Invocation{Name: "react_app"})
	case FRONTEND_NEXT:
		out = append(out, Invocation{Name: "nextjs_base"})
	}

	// Router (React+Vite only)
	if fw == FRONTEND_REACT_VITE {
		switch routerChoice {
		case FRONTEND_ROUTER_REACT_ROUTER:
			out = append(out, Invocation{Name: "react_router_v7"})
		case FRONTEND_ROUTER_TANSTACK:
			out = append(out, Invocation{Name: "tanstack_router"})
		}
	}

	// UI library (shadcn auto-includes Tailwind)
	switch uiLib {
	case FRONTEND_UI_SHADCN:
		out = append(out, Invocation{Name: "shadcn_ui"})
	case FRONTEND_UI_ARKUI:
		out = append(out, Invocation{Name: "ark_ui"})
	case FRONTEND_UI_VERSION14:
		out = append(out, Invocation{Name: "version14_ui"})
	}

	// Styling (skip if shadcn — it handles its own CSS)
	if uiLib != FRONTEND_UI_SHADCN {
		switch stylingChoice {
		case FRONTEND_STYLING_TAILWIND:
			out = append(out, Invocation{Name: "tailwind_v4"})
		case FRONTEND_STYLING_CSS_MODULES:
			out = append(out, Invocation{Name: "css_modules"})
		case FRONTEND_STYLING_PANDA:
			out = append(out, Invocation{Name: "panda_css"})
		}
	}

	// State management
	switch stateChoice {
	case FRONTEND_STATE_ZUSTAND:
		out = append(out, Invocation{Name: "zustand_setup"})
	case FRONTEND_STATE_JOTAI:
		out = append(out, Invocation{Name: "jotai_setup"})
	}

	// Testing
	if includeVitest && fw == FRONTEND_REACT_VITE {
		out = append(out, Invocation{Name: "vitest_testing_library"})
	}
	if includePlaywright {
		out = append(out, Invocation{Name: "playwright_setup"})
	}
	if includeStorybook {
		out = append(out, Invocation{Name: "storybook_setup"})
	}

	// Auth module
	if includeAuth {
		switch authProviderChoice {
		case FRONTEND_AUTH_CLERK:
			out = append(out, Invocation{Name: "auth_clerk_frontend"})
		case FRONTEND_AUTH_BETTER_AUTH:
			out = append(out, Invocation{Name: "auth_better_auth_frontend"})
		case FRONTEND_AUTH_VANILLA:
			out = append(out, Invocation{Name: "auth_vanilla_frontend"})
		}
	}

	// Theme module
	if includeTheme {
		out = append(out, Invocation{Name: "theme_provider"})
	}

	// Feature flags module
	if includeFeatureFlags {
		switch flagsProviderChoice {
		case FRONTEND_FLAGS_POSTHOG:
			out = append(out, Invocation{Name: "feature_flags_posthog"})
		case FRONTEND_FLAGS_VERCEL:
			out = append(out, Invocation{Name: "feature_flags_vercel"})
		case FRONTEND_FLAGS_LOCAL:
			out = append(out, Invocation{Name: "feature_flags_local"})
		}
	}

	// Sentry module
	if includeSentry {
		out = append(out, Invocation{Name: "sentry_frontend"})
	}

	// Analytics module — deduplicate PostHog with feature_flags_posthog
	if includeAnalytics {
		switch analyticsProviderChoice {
		case FRONTEND_ANALYTICS_GA4:
			out = append(out, Invocation{Name: "analytics_ga4"})
		case FRONTEND_ANALYTICS_PLAUSIBLE:
			out = append(out, Invocation{Name: "analytics_plausible"})
		case FRONTEND_ANALYTICS_POSTHOG:
			// Only emit feature_flags_posthog if not already in the list
			if !includeFeatureFlags || flagsProviderChoice != FRONTEND_FLAGS_POSTHOG {
				out = append(out, Invocation{Name: "feature_flags_posthog"})
			}
		}
	}

	// SEO module (framework-aware)
	if includeSEO {
		switch fw {
		case FRONTEND_NEXT:
			out = append(out, Invocation{Name: "seo_next"})
		case FRONTEND_REACT_VITE:
			out = append(out, Invocation{Name: "seo_react"})
		}
	}

	// Formatters last (DependsOn: ["*"] in their manifests)
	switch formatterChoice {
	case "prettier":
		out = append(out, Invocation{Name: "prettier_config"})
		out = append(out, Invocation{Name: "prettier_typescript_deps"})
		out = append(out, Invocation{Name: "prettier_frontend_rules"})
	case "biome":
		out = append(out, Invocation{Name: "biome_config"})
	}

	return out
}
