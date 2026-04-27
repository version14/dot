package cli

import (
	"fmt"

	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/plugin"
)

// Runtime bundles every per-process registry the scaffold pipeline needs:
// the generator catalog, the flow hook registry (where plugin injections
// live), the fragment registry, and the loaded Provider list (for
// ResolveExtras). Constructed once per CLI invocation.
//
// Reusing a single Runtime across Scaffold/Update/Doctor keeps plugin state
// consistent — generators registered by a plugin are visible to all of them.
type Runtime struct {
	Generators *generator.Registry
	Hooks      *flow.HookRegistry
	Fragments  *flow.FragmentRegistry
	Plugins    []plugin.Provider
}

// DefaultRuntime returns a Runtime pre-populated with every built-in
// generator and every registered built-in plugin, plus on-disk plugin
// discovery (manifest-only for now).
func DefaultRuntime() (*Runtime, error) {
	rt := &Runtime{
		Generators: generator.NewRegistry(),
		Hooks:      flow.NewHookRegistry(),
		Fragments:  flow.NewFragmentRegistry(),
	}

	if err := registerBuiltinGenerators(rt.Generators); err != nil {
		return nil, err
	}

	providers, err := plugin.Load(rt.Generators, rt.Hooks, plugin.LoadOptions{})
	if err != nil {
		return nil, fmt.Errorf("cli: load plugins: %w", err)
	}
	rt.Plugins = providers

	return rt, nil
}

// registerBuiltinGenerators is the in-tree counterpart to plugin loading: it
// adds every built-in generator to a fresh registry. Kept separate so tests
// can build minimal runtimes without the full catalog.
func registerBuiltinGenerators(r *generator.Registry) error {
	for _, e := range builtinGeneratorEntries() {
		if err := r.Register(e.Manifest, e.Generator); err != nil {
			return fmt.Errorf("cli: register %s: %w", e.Manifest.Name, err)
		}
	}
	return nil
}
