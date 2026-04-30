package cli

import (
	"fmt"

	authbetterauth "github.com/version14/dot/generators/auth_better_auth"
	authbetterauthschema "github.com/version14/dot/generators/auth_better_auth_schema"
	authjwtcleanarchmodule "github.com/version14/dot/generators/auth_jwt_clean_arch_module"
	authjwtmvcroute "github.com/version14/dot/generators/auth_jwt_mvc_route"
	authjwtusersschema "github.com/version14/dot/generators/auth_jwt_users_schema"
	authjwtvanilla "github.com/version14/dot/generators/auth_jwt_vanilla"
	backendArchitectureCleanArchitecture "github.com/version14/dot/generators/backend_architecture_clean_architecture"
	backendArchitectureHexagonal "github.com/version14/dot/generators/backend_architecture_hexagonal_architecture"
	backendArchitectureMVC "github.com/version14/dot/generators/backend_architecture_mvc_architecture"
	baseproject "github.com/version14/dot/generators/base_project"
	biomeconfig "github.com/version14/dot/generators/biome_config"
	drizzleconfigbase "github.com/version14/dot/generators/drizzle_config_base"
	drizzlepostgresadapter "github.com/version14/dot/generators/drizzle_postgres_adapter"
	drizzletypescriptdeps "github.com/version14/dot/generators/drizzle_typescript_deps"
	expressnodetsconfig "github.com/version14/dot/generators/express_node_tsconfig"
	expressserverentrypoint "github.com/version14/dot/generators/express_server_entrypoint"
	expressservertypescriptdeps "github.com/version14/dot/generators/express_server_typescript_deps"
	pluginreposkeleton "github.com/version14/dot/generators/plugin_repo_skeleton"
	postgresdockercompose "github.com/version14/dot/generators/postgres_docker_compose"
	postgresenvexample "github.com/version14/dot/generators/postgres_env_example"
	prettierconfig "github.com/version14/dot/generators/prettier_config"
	prettierexpressrules "github.com/version14/dot/generators/prettier_express_rules"
	prettiertypescriptdeps "github.com/version14/dot/generators/prettier_typescript_deps"
	reactapp "github.com/version14/dot/generators/react_app"
	servicewriter "github.com/version14/dot/generators/service_writer"
	typescriptbase "github.com/version14/dot/generators/typescript_base"
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
		{Manifest: biomeconfig.Manifest, Generator: biomeconfig.New()},
		{Manifest: servicewriter.Manifest, Generator: servicewriter.New()},
		{Manifest: pluginreposkeleton.Manifest, Generator: pluginreposkeleton.New()},

		// Backend architecture
		{Manifest: backendArchitectureCleanArchitecture.Manifest, Generator: backendArchitectureCleanArchitecture.New()},
		{Manifest: backendArchitectureMVC.Manifest, Generator: backendArchitectureMVC.New()},
		{Manifest: backendArchitectureHexagonal.Manifest, Generator: backendArchitectureHexagonal.New()},

		// Express server
		{Manifest: expressserverentrypoint.Manifest, Generator: expressserverentrypoint.New()},
		{Manifest: expressservertypescriptdeps.Manifest, Generator: expressservertypescriptdeps.New()},
		{Manifest: expressnodetsconfig.Manifest, Generator: expressnodetsconfig.New()},

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
