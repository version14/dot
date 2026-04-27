package generator

import (
	"strings"
	"testing"

	"github.com/version14/dot/pkg/dotapi"
)

// buildSortRegistry registers each manifest with a stub generator.
// Reuses the package-level fakeGen defined in executor_test.go.
func buildSortRegistry(t *testing.T, mans ...dotapi.Manifest) *Registry {
	t.Helper()
	r := NewRegistry()
	for _, m := range mans {
		if err := r.Register(m, &fakeGen{name: m.Name}); err != nil {
			t.Fatalf("register %s: %v", m.Name, err)
		}
	}
	return r
}

func invocationNames(invs []Invocation) []string {
	out := make([]string, len(invs))
	for i, inv := range invs {
		out[i] = inv.Name
	}
	return out
}

func TestSortInvocations_RespectsDependsOn(t *testing.T) {
	reg := buildSortRegistry(t,
		dotapi.Manifest{Name: "react_app", DependsOn: []string{"typescript_base"}},
		dotapi.Manifest{Name: "typescript_base", DependsOn: []string{"base_project"}},
		dotapi.Manifest{Name: "base_project"},
		dotapi.Manifest{Name: "biome_config", DependsOn: []string{"typescript_base"}},
	)

	// Request order is intentionally backwards from desired execution order.
	requested := []Invocation{
		{Name: "react_app"},
		{Name: "biome_config"},
		{Name: "typescript_base"},
		{Name: "base_project"},
	}

	sorted, err := SortInvocations(requested, reg)
	if err != nil {
		t.Fatalf("sort: %v", err)
	}

	got := invocationNames(sorted)
	want := []string{"base_project", "typescript_base", "react_app", "biome_config"}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("position %d: got %q, want %q (full: %v)", i, got[i], w, got)
		}
	}
}

func TestSortInvocations_DetectsCycle(t *testing.T) {
	reg := buildSortRegistry(t,
		dotapi.Manifest{Name: "a", DependsOn: []string{"b"}},
		dotapi.Manifest{Name: "b", DependsOn: []string{"a"}},
	)
	_, err := SortInvocations([]Invocation{{Name: "a"}, {Name: "b"}}, reg)
	if err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("expected cycle error, got %v", err)
	}
}

func TestSortInvocations_DetectsConflicts(t *testing.T) {
	reg := buildSortRegistry(t,
		dotapi.Manifest{Name: "alpha", ConflictsWith: []string{"beta"}},
		dotapi.Manifest{Name: "beta"},
	)
	_, err := SortInvocations([]Invocation{{Name: "alpha"}, {Name: "beta"}}, reg)
	if err == nil || !strings.Contains(err.Error(), "conflicts") {
		t.Fatalf("expected conflicts error, got %v", err)
	}
}

func TestSortInvocations_RejectsMissingDep(t *testing.T) {
	reg := buildSortRegistry(t,
		dotapi.Manifest{Name: "x", DependsOn: []string{"missing"}},
	)
	_, err := SortInvocations([]Invocation{{Name: "x"}}, reg)
	if err == nil || !strings.Contains(err.Error(), "depends on") {
		t.Fatalf("expected missing-dep error, got %v", err)
	}
}

func TestResolveInvocations_AutoAddsTransitiveDeps(t *testing.T) {
	reg := buildSortRegistry(t,
		dotapi.Manifest{Name: "base_project"},
		dotapi.Manifest{Name: "typescript_base", DependsOn: []string{"base_project"}},
		dotapi.Manifest{Name: "react_app", DependsOn: []string{"typescript_base"}},
	)
	// Only ask for react_app — base_project + typescript_base must be added.
	sorted, err := ResolveInvocations([]Invocation{{Name: "react_app"}}, reg)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	want := []string{"base_project", "typescript_base", "react_app"}
	got := invocationNames(sorted)
	if len(got) != 3 {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("position %d: got %q, want %q (full: %v)", i, got[i], w, got)
		}
	}
}
