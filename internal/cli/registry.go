package cli

import (
	"fmt"

	backendArchitectureCleanArchitecture "github.com/version14/dot/generators/backend_architecture_clean_architecture"
	backendArchitectureHexagonal "github.com/version14/dot/generators/backend_architecture_hexagonal_architecture"
	backendArchitectureMVC "github.com/version14/dot/generators/backend_architecture_mvc_architecture"
	baseproject "github.com/version14/dot/generators/base_project"
	biomeconfig "github.com/version14/dot/generators/biome_config"
	pluginreposkeleton "github.com/version14/dot/generators/plugin_repo_skeleton"
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
		{Manifest: baseproject.Manifest, Generator: baseproject.New()},
		{Manifest: typescriptbase.Manifest, Generator: typescriptbase.New()},
		{Manifest: reactapp.Manifest, Generator: reactapp.New()},
		{Manifest: biomeconfig.Manifest, Generator: biomeconfig.New()},
		{Manifest: servicewriter.Manifest, Generator: servicewriter.New()},
		{Manifest: pluginreposkeleton.Manifest, Generator: pluginreposkeleton.New()},
		{Manifest: backendArchitectureCleanArchitecture.Manifest, Generator: backendArchitectureCleanArchitecture.New()},
		{Manifest: backendArchitectureMVC.Manifest, Generator: backendArchitectureMVC.New()},
		{Manifest: backendArchitectureHexagonal.Manifest, Generator: backendArchitectureHexagonal.New()},
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
