package generator

import (
	"fmt"
	"sort"

	"github.com/version14/dot/internal/flow"
)

// SortInvocations returns invocations re-ordered so every generator runs
// AFTER its declared DependsOn, while preserving the resolver's original
// order between independent generators ("stable Kahn's algorithm").
//
// Errors:
//   - missing dependency  → invocation references a generator not in the request set
//   - cycle               → DependsOn chain forms a loop
//   - conflicting pair    → two requested generators mutually ConflictsWith
//
// Multiple invocations of the same generator (e.g. a loop body) are kept in
// their original order; the dependency graph is built per-name.
func SortInvocations(invs []Invocation, registry *Registry) ([]Invocation, error) {
	if registry == nil {
		return nil, fmt.Errorf("generator: nil registry")
	}
	if len(invs) == 0 {
		return invs, nil
	}

	// Resolve manifests up front so missing generators surface as one clean
	// error before any sort work happens.
	mans := make(map[string]Manifest, len(invs))
	originalOrder := make(map[string]int, len(invs))
	loopsByName := make(map[string][][]flow.LoopFrame, len(invs))

	for i, inv := range invs {
		entry, ok := registry.Get(inv.Name)
		if !ok {
			return nil, fmt.Errorf("generator: unknown %q in invocations", inv.Name)
		}
		mans[inv.Name] = entry.Manifest
		if _, seen := originalOrder[inv.Name]; !seen {
			originalOrder[inv.Name] = i
		}
		loopsByName[inv.Name] = append(loopsByName[inv.Name], inv.LoopStack)
	}

	if err := assertNoConflicts(mans); err != nil {
		return nil, err
	}

	indeg := make(map[string]int, len(mans))
	dependents := make(map[string][]string, len(mans))
	for name, m := range mans {
		indeg[name] = 0
		for _, d := range m.DependsOn {
			if d == "*" {
				continue
			}
			if _, ok := mans[d]; !ok {
				return nil, fmt.Errorf("generator %q depends on %q which is not in the invocation set", name, d)
			}
		}
	}
	// Identify "end-generators": those that have a wildcard dependency "*"
	// OR (transitively) depend on one that does.
	isEndGen := make(map[string]bool)
	for name, m := range mans {
		for _, d := range m.DependsOn {
			if d == "*" {
				isEndGen[name] = true
				break
			}
		}
	}

	// Propagate end-gen status: if A depends on B and B is an end-gen, A is too.
	// Since the number of generators is small, a simple fixed-point is fine.
	changed := true
	for changed {
		changed = false
		for name, m := range mans {
			if isEndGen[name] {
				continue
			}
			for _, d := range m.DependsOn {
				if d == "*" {
					continue
				}
				if isEndGen[d] {
					isEndGen[name] = true
					changed = true
					break
				}
			}
		}
	}

	for name, m := range mans {
		seenDeps := make(map[string]struct{}, len(m.DependsOn)+len(mans))
		for _, d := range m.DependsOn {
			if d == "*" {
				continue
			}
			seenDeps[d] = struct{}{}
		}

		// If this is an end-generator, it must run after ALL non-end-generators.
		// This still participates in normal topological sorting, so wildcard
		// generators are correctly ordered with each other and with any
		// explicit dependency chains.
		if isEndGen[name] {
			for other := range mans {
				if !isEndGen[other] {
					seenDeps[other] = struct{}{}
				}
			}
		}

		for dep := range seenDeps {
			indeg[name]++
			dependents[dep] = append(dependents[dep], name)
		}
	}

	ready := make([]string, 0, len(indeg))
	for name, d := range indeg {
		if d == 0 {
			ready = append(ready, name)
		}
	}
	sort.SliceStable(ready, func(i, j int) bool {
		return originalOrder[ready[i]] < originalOrder[ready[j]]
	})

	out := make([]Invocation, 0, len(invs))
	for len(ready) > 0 {
		next := ready[0]
		ready = ready[1:]

		// Emit one Invocation per recorded loop stack for this name, in order.
		for _, ls := range loopsByName[next] {
			out = append(out, Invocation{Name: next, LoopStack: ls})
		}

		for _, dep := range dependents[next] {
			indeg[dep]--
			if indeg[dep] == 0 {
				ready = insertSorted(ready, dep, originalOrder)
			}
		}
	}

	if len(out) != len(invs) {
		return nil, fmt.Errorf("generator: dependency cycle detected (%d of %d invocations unresolved)",
			len(invs)-len(out), len(invs))
	}
	return out, nil
}

// insertSorted places name into ready preserving original-order ranking.
func insertSorted(ready []string, name string, order map[string]int) []string {
	for i, existing := range ready {
		if order[name] < order[existing] {
			return append(ready[:i], append([]string{name}, ready[i:]...)...)
		}
	}
	return append(ready, name)
}

// assertNoConflicts walks ConflictsWith and returns a descriptive error if any
// pair of mutually-incompatible generators are both present.
func assertNoConflicts(mans map[string]Manifest) error {
	for name, m := range mans {
		for _, c := range m.ConflictsWith {
			if _, present := mans[c]; present {
				return fmt.Errorf("generator %q conflicts with %q (both requested)", name, c)
			}
		}
	}
	return nil
}
