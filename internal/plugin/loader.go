package plugin

import (
	"fmt"

	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/spec"
)

// Provider is the in-process surface a plugin implements. Each plugin
// contributes generators, flow injections, and (optionally) extra generator
// invocations whose visibility depends on plugin-contributed answers.
//
// External plugins (out-of-process / wasm / git-cloned) are not yet supported;
// the loader currently only accepts in-process Providers registered via
// RegisterBuiltin so the public surface is stable while implementations land.
type Provider interface {
	// ID is the plugin's namespace prefix. MUST NOT contain '.'. All
	// generator names, question IDs, and option values the plugin contributes
	// MUST start with "<ID>.".
	ID() flow.PluginID

	// Generators returns (Manifest, Generator) pairs to register with the
	// generator.Registry at startup.
	Generators() []generator.Entry

	// Injections returns flow.Injection values to add to the engine's
	// HookRegistry. Replace / AddOption / InsertAfter all work the same way
	// once registered — see internal/flow/hook.go for semantics.
	Injections() []*flow.Injection

	// ResolveExtras inspects the populated spec and returns extra generator
	// invocations to append to the flow's resolver result. This is how a
	// plugin "activates" a generator based on plugin-contributed answers
	// (e.g. biome_extras adds biome_extras.strict_writer when its
	// "biome_extras.strict_mode" question was answered true).
	//
	// Plugins that don't add generators conditionally can return nil.
	ResolveExtras(s *spec.ProjectSpec) []generator.Invocation
}

var builtinProviders []Provider

// RegisterBuiltin lets in-tree plugins register themselves at init() time.
// Unlike on-disk plugins these are always present in every binary.
func RegisterBuiltin(p Provider) {
	if p == nil {
		return
	}
	builtinProviders = append(builtinProviders, p)
}

// LoadOptions controls Load's behaviour.
type LoadOptions struct {
	// PluginDir overrides the default ~/.dot/plugins lookup.
	PluginDir string
	// SkipDisk suppresses on-disk discovery (useful for tests).
	SkipDisk bool
	// SkipBuiltins suppresses built-in providers (useful for tests).
	SkipBuiltins bool
}

// Load registers every available plugin (in-tree built-ins + on-disk
// installs) into the supplied registries. It returns the active Providers so
// callers can call ResolveExtras during scaffolding.
//
// On-disk loading currently inspects manifests only; dynamic loading of
// generator binaries is deferred. The skeleton is in place so adding it later
// will not change any caller signatures.
func Load(genReg *generator.Registry, hookReg *flow.HookRegistry, opts LoadOptions) ([]Provider, error) {
	if genReg == nil || hookReg == nil {
		return nil, fmt.Errorf("plugin: nil registry")
	}

	loaded := make([]Provider, 0)

	if !opts.SkipBuiltins {
		for _, p := range builtinProviders {
			if err := registerProvider(p, genReg, hookReg); err != nil {
				return loaded, err
			}
			loaded = append(loaded, p)
		}
	}

	if !opts.SkipDisk {
		dir := opts.PluginDir
		var (
			installed []*Installed
			err       error
		)
		if dir != "" {
			installed, err = ListIn(dir)
		} else {
			installed, err = List()
		}
		if err != nil {
			return loaded, err
		}
		// Surfaces what the user has installed without instantiating them.
		// Once dynamic loading lands the provider list will include disk
		// plugins too.
		_ = installed
	}

	return loaded, nil
}

func registerProvider(p Provider, genReg *generator.Registry, hookReg *flow.HookRegistry) error {
	if err := hookReg.RegisterPlugin(p.ID()); err != nil {
		return fmt.Errorf("plugin %q: %w", p.ID(), err)
	}
	for _, e := range p.Generators() {
		if err := genReg.Register(e.Manifest, e.Generator); err != nil {
			return fmt.Errorf("plugin %q: register %s: %w", p.ID(), e.Manifest.Name, err)
		}
	}
	for _, inj := range p.Injections() {
		if err := hookReg.Inject(inj); err != nil {
			return fmt.Errorf("plugin %q: inject %s: %w", p.ID(), inj.TargetID, err)
		}
	}
	return nil
}
