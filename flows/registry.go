package flows

import (
	"fmt"

	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/spec"
)

// Invocation is a (generator-name, loop-stack) tuple the resolver produces.
// Caller (CLI / test-flow) converts these into generator.Invocation values
// without forcing the flows package to depend on internal/generator.
type Invocation struct {
	Name      string
	LoopStack []flow.LoopFrame
}

// FlowDef bundles a flow's metadata with its root question and a resolver
// that maps the populated spec to the ordered list of generators to run.
type FlowDef struct {
	ID          string
	Title       string
	Description string
	Root        flow.Question

	// Generators returns generator invocations in the order they should run,
	// based on the answers in the populated spec. Loop stacks scope a
	// generator to a particular loop iteration's answers.
	Generators func(*spec.ProjectSpec) []Invocation
}

// Registry holds every available flow. Use Register to add new flows during
// program startup; the CLI consults this map to enumerate options.
type Registry struct {
	flows map[string]*FlowDef
}

func NewRegistry() *Registry {
	return &Registry{flows: map[string]*FlowDef{}}
}

func (r *Registry) Register(def *FlowDef) error {
	if def == nil || def.ID == "" {
		return fmt.Errorf("flows: nil flow or empty ID")
	}
	if _, exists := r.flows[def.ID]; exists {
		return fmt.Errorf("flows: %q already registered", def.ID)
	}
	r.flows[def.ID] = def
	return nil
}

func (r *Registry) Get(id string) (*FlowDef, bool) {
	f, ok := r.flows[id]
	return f, ok
}

func (r *Registry) All() []*FlowDef {
	out := make([]*FlowDef, 0, len(r.flows))
	for _, f := range r.flows {
		out = append(out, f)
	}
	return out
}

// Default returns a Registry pre-loaded with all built-in flows.
func Default() *Registry {
	r := NewRegistry()
	_ = r.Register(MonorepoFlow())
	_ = r.Register(FullstackFlow())
	_ = r.Register(MicroservicesFlow())
	_ = r.Register(PluginTemplateFlow())
	return r
}
