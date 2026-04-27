package generator

import "fmt"

// ResolveInvocations expands a request set with the transitive closure of
// every DependsOn declaration, then runs the topological sort.
//
// Behaviour:
//   - Transitive dependencies are auto-added at most once each, with an
//     empty LoopStack. (A dep can't sensibly target a loop frame; the loop
//     resolver is the explicit-request side's job.)
//   - Explicitly requested invocations are preserved AS-IS, including
//     duplicates with different LoopStacks (so a per-iteration generator
//     can run N times).
//   - The combined list is sorted via SortInvocations (stable Kahn).
//
// This is the public entry point the CLI / test-flow uses; flows return raw
// requests and the resolver fills in the rest.
func ResolveInvocations(requested []Invocation, registry *Registry) ([]Invocation, error) {
	if registry == nil {
		return nil, fmt.Errorf("generator: nil registry")
	}

	// Step 1: collect every generator name we need (requested + transitive deps).
	// `seenNames` ensures each TRANSITIVE dep is added at most once. The set
	// of explicitly-requested names is recorded separately so they survive
	// dedup if their dep closure also touches them.
	seenNames := make(map[string]bool)
	requestedNames := make(map[string]bool, len(requested))
	for _, inv := range requested {
		requestedNames[inv.Name] = true
	}

	deps := make([]Invocation, 0)

	var visit func(name string) error
	visit = func(name string) error {
		if seenNames[name] {
			return nil
		}
		seenNames[name] = true
		entry, ok := registry.Get(name)
		if !ok {
			return fmt.Errorf("generator: unknown %q", name)
		}
		for _, dep := range entry.Manifest.DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}
		// Only add to `deps` if NOT explicitly requested — caller's invocation
		// (with its possibly-non-empty LoopStack) wins for those.
		if !requestedNames[name] {
			deps = append(deps, Invocation{Name: name})
		}
		return nil
	}

	for _, inv := range requested {
		if err := visit(inv.Name); err != nil {
			return nil, err
		}
	}

	// Step 2: combine — deps first (so they appear before their dependents
	// in the pre-sort list, easing trace reading), then explicit requests.
	combined := make([]Invocation, 0, len(deps)+len(requested))
	combined = append(combined, deps...)
	combined = append(combined, requested...)

	return SortInvocations(combined, registry)
}
