package cli

import (
	"fmt"

	analyticsga4 "github.com/version14/dot/generators/analytics_ga4"
	analyticsplausible "github.com/version14/dot/generators/analytics_plausible"
	arkui "github.com/version14/dot/generators/ark_ui"
	authbetterauth "github.com/version14/dot/generators/auth_better_auth"
	authbetterauthfrontend "github.com/version14/dot/generators/auth_better_auth_frontend"
	authbetterauthschema "github.com/version14/dot/generators/auth_better_auth_schema"
	authclerkfrontend "github.com/version14/dot/generators/auth_clerk_frontend"
	authjwtcleanarchmodule "github.com/version14/dot/generators/auth_jwt_clean_arch_module"
	authjwtmvcroute "github.com/version14/dot/generators/auth_jwt_mvc_route"
	authjwtusersschema "github.com/version14/dot/generators/auth_jwt_users_schema"
	authjwtvanilla "github.com/version14/dot/generators/auth_jwt_vanilla"
	authvanillafrontend "github.com/version14/dot/generators/auth_vanilla_frontend"
	backendArchitectureCleanArchitecture "github.com/version14/dot/generators/backend_architecture_clean_architecture"
	backendArchitectureHexagonal "github.com/version14/dot/generators/backend_architecture_hexagonal_architecture"
	backendArchitectureMVC "github.com/version14/dot/generators/backend_architecture_mvc_architecture"
	baseproject "github.com/version14/dot/generators/base_project"
	biomeconfig "github.com/version14/dot/generators/biome_config"
	cssmodules "github.com/version14/dot/generators/css_modules"
	decoratorscleanarchadapter "github.com/version14/dot/generators/decorators_clean_arch_adapter"
	decoratorshexagonaladapter "github.com/version14/dot/generators/decorators_hexagonal_adapter"
	decoratorsmvcadapter "github.com/version14/dot/generators/decorators_mvc_adapter"
	drizzleconfigbase "github.com/version14/dot/generators/drizzle_config_base"
	drizzlepostgresadapter "github.com/version14/dot/generators/drizzle_postgres_adapter"
	drizzletypescriptdeps "github.com/version14/dot/generators/drizzle_typescript_deps"
	expressauthvalidators "github.com/version14/dot/generators/express_auth_validators"
	expressdecoratorscore "github.com/version14/dot/generators/express_decorators_core"
	expresserrormiddleware "github.com/version14/dot/generators/express_error_middleware"
	expressnodetsconfig "github.com/version14/dot/generators/express_node_tsconfig"
	expressopenapisetup "github.com/version14/dot/generators/express_openapi_setup"
	expressratelimit "github.com/version14/dot/generators/express_rate_limit"
	expressserverentrypoint "github.com/version14/dot/generators/express_server_entrypoint"
	expressservertypescriptdeps "github.com/version14/dot/generators/express_server_typescript_deps"
	expresssharederrors "github.com/version14/dot/generators/express_shared_errors"
	expressswaggerjsdoc "github.com/version14/dot/generators/express_swagger_jsdoc"
	expresstestsetup "github.com/version14/dot/generators/express_test_setup"
	featureflagslocal "github.com/version14/dot/generators/feature_flags_local"
	featureflagsposthog "github.com/version14/dot/generators/feature_flags_posthog"
	featureflagsvercel "github.com/version14/dot/generators/feature_flags_vercel"
	jotaisetup "github.com/version14/dot/generators/jotai_setup"
	monorepotsworkspaces "github.com/version14/dot/generators/monorepo_ts_workspaces"
	nextjsbase "github.com/version14/dot/generators/nextjs_base"
	pandacss "github.com/version14/dot/generators/panda_css"
	playwrightsetup "github.com/version14/dot/generators/playwright_setup"
	pluginreposkeleton "github.com/version14/dot/generators/plugin_repo_skeleton"
	postgresdockercompose "github.com/version14/dot/generators/postgres_docker_compose"
	postgresenvexample "github.com/version14/dot/generators/postgres_env_example"
	prettierconfig "github.com/version14/dot/generators/prettier_config"
	prettierexpressrules "github.com/version14/dot/generators/prettier_express_rules"
	prettierfrontendrules "github.com/version14/dot/generators/prettier_frontend_rules"
	prettiertypescriptdeps "github.com/version14/dot/generators/prettier_typescript_deps"
	reactapp "github.com/version14/dot/generators/react_app"
	reactrouterv7 "github.com/version14/dot/generators/react_router_v7"
	sentryfrontend "github.com/version14/dot/generators/sentry_frontend"
	seonext "github.com/version14/dot/generators/seo_next"
	seoreact "github.com/version14/dot/generators/seo_react"
	shadcnui "github.com/version14/dot/generators/shadcn_ui"
	storybooksetup "github.com/version14/dot/generators/storybook_setup"
	tailwindv4 "github.com/version14/dot/generators/tailwind_v4"
	tanstackrouter "github.com/version14/dot/generators/tanstack_router"
	themeprovider "github.com/version14/dot/generators/theme_provider"
	typescriptbase "github.com/version14/dot/generators/typescript_base"
	version14ui "github.com/version14/dot/generators/version14_ui"
	vitesttestinglibrary "github.com/version14/dot/generators/vitest_testing_library"
	zodvalidationdeps "github.com/version14/dot/generators/zod_validation_deps"
	zustandsetup "github.com/version14/dot/generators/zustand_setup"
	"github.com/version14/dot/internal/generator"
)

// builtinGeneratorEntries returns the canonical list of in-tree generators.
// Kept as a function (not a var) so each call yields fresh Generator instances
// — important when tests build multiple registries in the same process.
func builtinGeneratorEntries() []generator.Entry {
	return []generator.Entry{
		// Foundation
		{Manifest: baseproject.Manifest, Generator: baseproject.New()},
		{Manifest: typescriptbase.Manifest, Generator: typescriptbase.New()},
		{Manifest: reactapp.Manifest, Generator: reactapp.New()},
		{Manifest: nextjsbase.Manifest, Generator: nextjsbase.New()},
		{Manifest: biomeconfig.Manifest, Generator: biomeconfig.New()},
		{Manifest: monorepotsworkspaces.Manifest, Generator: monorepotsworkspaces.New()},
		{Manifest: pluginreposkeleton.Manifest, Generator: pluginreposkeleton.New()},

		// Frontend — router
		{Manifest: reactrouterv7.Manifest, Generator: reactrouterv7.New()},
		{Manifest: tanstackrouter.Manifest, Generator: tanstackrouter.New()},

		// Frontend — UI library
		{Manifest: shadcnui.Manifest, Generator: shadcnui.New()},
		{Manifest: arkui.Manifest, Generator: arkui.New()},
		{Manifest: version14ui.Manifest, Generator: version14ui.New()},

		// Frontend — styling
		{Manifest: tailwindv4.Manifest, Generator: tailwindv4.New()},
		{Manifest: cssmodules.Manifest, Generator: cssmodules.New()},
		{Manifest: pandacss.Manifest, Generator: pandacss.New()},

		// Frontend — state
		{Manifest: zustandsetup.Manifest, Generator: zustandsetup.New()},
		{Manifest: jotaisetup.Manifest, Generator: jotaisetup.New()},

		// Frontend — testing
		{Manifest: vitesttestinglibrary.Manifest, Generator: vitesttestinglibrary.New()},
		{Manifest: playwrightsetup.Manifest, Generator: playwrightsetup.New()},
		{Manifest: storybooksetup.Manifest, Generator: storybooksetup.New()},

		// Frontend — auth modules
		{Manifest: authclerkfrontend.Manifest, Generator: authclerkfrontend.New()},
		{Manifest: authbetterauthfrontend.Manifest, Generator: authbetterauthfrontend.New()},
		{Manifest: authvanillafrontend.Manifest, Generator: authvanillafrontend.New()},

		// Frontend — modules
		{Manifest: themeprovider.Manifest, Generator: themeprovider.New()},
		{Manifest: featureflagsposthog.Manifest, Generator: featureflagsposthog.New()},
		{Manifest: featureflagsvercel.Manifest, Generator: featureflagsvercel.New()},
		{Manifest: featureflagslocal.Manifest, Generator: featureflagslocal.New()},
		{Manifest: sentryfrontend.Manifest, Generator: sentryfrontend.New()},
		{Manifest: analyticsga4.Manifest, Generator: analyticsga4.New()},
		{Manifest: analyticsplausible.Manifest, Generator: analyticsplausible.New()},
		{Manifest: seonext.Manifest, Generator: seonext.New()},
		{Manifest: seoreact.Manifest, Generator: seoreact.New()},

		// Frontend — formatter rules
		{Manifest: prettierfrontendrules.Manifest, Generator: prettierfrontendrules.New()},

		// Backend architecture
		{Manifest: backendArchitectureCleanArchitecture.Manifest, Generator: backendArchitectureCleanArchitecture.New()},
		{Manifest: backendArchitectureMVC.Manifest, Generator: backendArchitectureMVC.New()},
		{Manifest: backendArchitectureHexagonal.Manifest, Generator: backendArchitectureHexagonal.New()},

		// Express server
		{Manifest: expressserverentrypoint.Manifest, Generator: expressserverentrypoint.New()},
		{Manifest: expressservertypescriptdeps.Manifest, Generator: expressservertypescriptdeps.New()},
		{Manifest: expressnodetsconfig.Manifest, Generator: expressnodetsconfig.New()},
		{Manifest: expresssharederrors.Manifest, Generator: expresssharederrors.New()},
		{Manifest: expresserrormiddleware.Manifest, Generator: expresserrormiddleware.New()},
		{Manifest: expressratelimit.Manifest, Generator: expressratelimit.New()},
		{Manifest: expresstestsetup.Manifest, Generator: expresstestsetup.New()},
		{Manifest: expressauthvalidators.Manifest, Generator: expressauthvalidators.New()},

		// OpenAPI / Swagger — classic JSDoc path
		{Manifest: expressswaggerjsdoc.Manifest, Generator: expressswaggerjsdoc.New()},

		// Decorator-based validation + OpenAPI
		{Manifest: zodvalidationdeps.Manifest, Generator: zodvalidationdeps.New()},
		{Manifest: expressdecoratorscore.Manifest, Generator: expressdecoratorscore.New()},
		{Manifest: expressopenapisetup.Manifest, Generator: expressopenapisetup.New()},
		{Manifest: decoratorscleanarchadapter.Manifest, Generator: decoratorscleanarchadapter.New()},
		{Manifest: decoratorsmvcadapter.Manifest, Generator: decoratorsmvcadapter.New()},
		{Manifest: decoratorshexagonaladapter.Manifest, Generator: decoratorshexagonaladapter.New()},

		// Prettier
		{Manifest: prettierconfig.Manifest, Generator: prettierconfig.New()},
		{Manifest: prettiertypescriptdeps.Manifest, Generator: prettiertypescriptdeps.New()},
		{Manifest: prettierexpressrules.Manifest, Generator: prettierexpressrules.New()},

		// PostgreSQL
		{Manifest: postgresdockercompose.Manifest, Generator: postgresdockercompose.New()},
		{Manifest: postgresenvexample.Manifest, Generator: postgresenvexample.New()},

		// Drizzle ORM
		{Manifest: drizzleconfigbase.Manifest, Generator: drizzleconfigbase.New()},
		{Manifest: drizzletypescriptdeps.Manifest, Generator: drizzletypescriptdeps.New()},
		{Manifest: drizzlepostgresadapter.Manifest, Generator: drizzlepostgresadapter.New()},

		// Auth
		{Manifest: authbetterauth.Manifest, Generator: authbetterauth.New()},
		{Manifest: authjwtvanilla.Manifest, Generator: authjwtvanilla.New()},
		{Manifest: authbetterauthschema.Manifest, Generator: authbetterauthschema.New()},
		{Manifest: authjwtusersschema.Manifest, Generator: authjwtusersschema.New()},
		{Manifest: authjwtmvcroute.Manifest, Generator: authjwtmvcroute.New()},
		{Manifest: authjwtcleanarchmodule.Manifest, Generator: authjwtcleanarchmodule.New()},
	}
}

// DefaultGeneratorRegistry returns a generator.Registry pre-loaded with every
// built-in generator. Plugin generators are NOT included — use DefaultRuntime
// for the full picture.
//
// Kept for callers (mostly tests) that don't need the plugin layer.
func DefaultGeneratorRegistry() (*generator.Registry, error) {
	r := generator.NewRegistry()
	for _, e := range builtinGeneratorEntries() {
		if err := r.Register(e.Manifest, e.Generator); err != nil {
			return nil, fmt.Errorf("cli: register %s: %w", e.Manifest.Name, err)
		}
	}
	return r, nil
}
