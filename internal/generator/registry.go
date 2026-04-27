package generator

import (
	"fmt"
	"sort"

	"github.com/version14/dot/pkg/dotapi"
)

// Entry pairs a Manifest with its concrete Generator implementation.
type Entry struct {
	Manifest  Manifest
	Generator dotapi.Generator
}

// Registry holds every known generator (built-in, community, local). It is
// populated at startup and consulted by the resolver.
type Registry struct {
	entries map[string]*Entry
}

func NewRegistry() *Registry {
	return &Registry{entries: map[string]*Entry{}}
}

// Register adds a generator. Returns an error if the name is already taken.
func (r *Registry) Register(m Manifest, g dotapi.Generator) error {
	if m.Name == "" {
		return fmt.Errorf("generator: manifest missing Name")
	}
	if g == nil {
		return fmt.Errorf("generator: nil generator for %q", m.Name)
	}
	if _, exists := r.entries[m.Name]; exists {
		return fmt.Errorf("generator: %q already registered", m.Name)
	}
	r.entries[m.Name] = &Entry{Manifest: m, Generator: g}
	return nil
}

func (r *Registry) Get(name string) (*Entry, bool) {
	e, ok := r.entries[name]
	return e, ok
}

// Names returns every registered generator name in lexicographic order.
func (r *Registry) Names() []string {
	out := make([]string, 0, len(r.entries))
	for k := range r.entries {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// All returns a copy of every entry, suitable for iteration without holding
// a reference to the registry's internal map.
func (r *Registry) All() []*Entry {
	out := make([]*Entry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Manifest.Name < out[j].Manifest.Name
	})
	return out
}
