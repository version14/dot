package spec

import (
	"time"

	"github.com/version14/dot/internal/flow"
)

// BuildOpts carries the inputs Build needs but the FlowContext doesn't
// (project name, tool version, declared generator constraints).
type BuildOpts struct {
	FlowID               string
	ProjectName          string
	ToolVersion          string
	GeneratorConstraints map[string]string
}

// Build converts a populated FlowContext into a ProjectSpec ready for
// persistence to .dot/spec.json.
func Build(ctx *flow.FlowContext, opts BuildOpts) *ProjectSpec {
	if ctx == nil {
		ctx = &flow.FlowContext{}
	}

	answers := make(map[string]AnswerNode, len(ctx.Answers))
	for k, v := range ctx.Answers {
		answers[k] = v
	}

	visited := append([]string(nil), ctx.VisitedNodes...)
	plugins := append([]string(nil), ctx.LoadedPlugins...)

	constraints := map[string]string{}
	for k, v := range opts.GeneratorConstraints {
		constraints[k] = v
	}

	return &ProjectSpec{
		FlowID:    opts.FlowID,
		CreatedAt: time.Now().UTC(),
		Metadata: ProjectMetadata{
			ProjectName: opts.ProjectName,
			ToolVersion: opts.ToolVersion,
		},
		Answers:              answers,
		VisitedNodes:         visited,
		LoadedPlugins:        plugins,
		GeneratorConstraints: constraints,
	}
}
