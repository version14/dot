package flow

import "fmt"

// FlowFragment is a named contextual resolver. Resolve returns the next edge
// to follow, computed at traversal time so plugins can inject branches.
type FlowFragment struct {
	ID      string
	Resolve func(ctx *FlowContext) *Next
}

type FragmentRegistry struct {
	fragments map[string]*FlowFragment
}

func NewFragmentRegistry() *FragmentRegistry {
	return &FragmentRegistry{fragments: map[string]*FlowFragment{}}
}

func (r *FragmentRegistry) Register(f *FlowFragment) error {
	if f == nil {
		return fmt.Errorf("flow: cannot register nil fragment")
	}
	if f.ID == "" {
		return fmt.Errorf("flow: fragment ID required")
	}
	if _, exists := r.fragments[f.ID]; exists {
		return fmt.Errorf("flow: fragment %q already registered", f.ID)
	}
	r.fragments[f.ID] = f
	return nil
}

func (r *FragmentRegistry) Get(id string) (*FlowFragment, bool) {
	f, ok := r.fragments[id]
	return f, ok
}

// Resolve runs a fragment's resolver against the current context.
// Returns nil if the fragment is not registered (no-op pass-through).
func (r *FragmentRegistry) Resolve(id string, ctx *FlowContext) *Next {
	f, ok := r.fragments[id]
	if !ok {
		return nil
	}
	return f.Resolve(ctx)
}
