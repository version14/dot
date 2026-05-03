// Package flows registers the Test Flow flow.
package flows

import (
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/spec"
)

// TestFlowFlow scaffolds TODO.
func TestFlowFlow() *FlowDef {
	root := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "project_name",
			Next_: &flow.Next{End: true},
		},
		Label:    "Project name",
		Validate: nonEmpty,
	}

	return &FlowDef{
		ID:          "test-flow",
		Title:       "Test Flow",
		Description: "A test flow to validate the dot-flow skill",
		Root:        root,
		Generators:  resolveTestFlowFlowGenerators,
	}
}

func resolveTestFlowFlowGenerators(_ *spec.ProjectSpec) []Invocation {
	// TODO: emit one Invocation per generator the flow should run.
	return []Invocation{}
}
