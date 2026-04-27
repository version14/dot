package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/version14/dot/flows"
	"github.com/version14/dot/internal/commands"
	"github.com/version14/dot/internal/dotdir"
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/plugin"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

// UpdateOptions configures a re-run scaffold against an existing project.
type UpdateOptions struct {
	ProjectRoot      string // dir containing .dot/spec.json
	Registry         *generator.Registry
	Plugins          []plugin.Provider
	Logger           dotapi.Logger
	ToolVersion      string
	OverrideAnswers  map[string]flow.Answer // optional answer overrides
	SkipPostCommands bool
}

// Update re-runs the generators of an existing project from its persisted
// .dot/spec.json. The flow is NOT replayed; we use the recorded answers
// verbatim, resolve generators via the originating flow's resolver, then
// execute against a fresh VirtualProjectState before persisting on top of
// the existing tree.
//
// Use cases: bumping generator versions, applying new plugin injections,
// regenerating after a manual edit was lost.
func Update(ctx context.Context, opts UpdateOptions) (*ScaffoldResult, error) {
	_ = ctx
	if opts.ProjectRoot == "" {
		return nil, fmt.Errorf("cli: update: empty ProjectRoot")
	}
	if opts.Registry == nil {
		return nil, fmt.Errorf("cli: update: nil registry")
	}
	if opts.Logger == nil {
		opts.Logger = dotapi.DiscardLogger{}
	}

	abs, err := filepath.Abs(opts.ProjectRoot)
	if err != nil {
		return nil, fmt.Errorf("cli: update: abs path: %w", err)
	}

	s, err := dotdir.LoadSpec(abs)
	if err != nil {
		return nil, fmt.Errorf("cli: update: load spec: %w", err)
	}

	// Apply answer overrides (used to force a value change without re-running
	// the interactive flow).
	for k, v := range opts.OverrideAnswers {
		s.Answers[k] = v
	}

	flowReg := flows.Default()
	def, ok := flowReg.Get(s.FlowID)
	if !ok {
		return nil, fmt.Errorf("cli: update: unknown flow %q in spec", s.FlowID)
	}
	if def.Generators == nil {
		return nil, fmt.Errorf("cli: update: flow %q has no Generators resolver", s.FlowID)
	}

	flowInvs := def.Generators(s)
	requested := make([]generator.Invocation, 0, len(flowInvs))
	for _, fi := range flowInvs {
		requested = append(requested, generator.Invocation{Name: fi.Name, LoopStack: fi.LoopStack})
	}
	for _, p := range opts.Plugins {
		requested = append(requested, p.ResolveExtras(s)...)
	}
	invs, err := generator.ResolveInvocations(requested, opts.Registry)
	if err != nil {
		return nil, fmt.Errorf("cli: update: resolve: %w", err)
	}

	mans := make([]dotapi.Manifest, len(invs))
	for i, inv := range invs {
		entry, ok := opts.Registry.Get(inv.Name)
		if !ok {
			return nil, fmt.Errorf("cli: missing generator %q after resolve", inv.Name)
		}
		mans[i] = entry.Manifest
	}

	vstate := state.NewVirtualProjectState(s.Metadata)
	exec := generator.NewExecutor(opts.Registry, opts.Logger)

	start := time.Now()
	opts.Logger.Infof("→ re-running %d generators in %s", len(invs), abs)
	if err := exec.Execute(invs, s, vstate); err != nil {
		return nil, fmt.Errorf("cli: update: execute: %w", err)
	}

	count, err := state.Persist(vstate, abs)
	if err != nil {
		return nil, fmt.Errorf("cli: update: persist: %w", err)
	}
	opts.Logger.Infof("→ rewrote %d files", count)

	if err := dotdir.SaveSpec(abs, s); err != nil {
		return nil, fmt.Errorf("cli: update: save spec: %w", err)
	}
	if err := dotdir.SaveManifest(abs, manifestSummary(invs, mans, opts.ToolVersion, time.Since(start))); err != nil {
		return nil, fmt.Errorf("cli: update: save manifest: %w", err)
	}

	return &ScaffoldResult{
		Spec:        s,
		State:       vstate,
		ProjectRoot: abs,
		Invocations: invs,
		Manifests:   mans,
		Duration:    time.Since(start),
	}, nil
}

// runUpdate is the CLI dispatch for `dot update [path]`.
func runUpdate(ctx context.Context, args []string, toolVersion string) int {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	rt, err := DefaultRuntime()
	if err != nil {
		PrintError(err.Error())
		return 1
	}

	logger := NewStepLogger()
	res, err := Update(ctx, UpdateOptions{
		ProjectRoot: root,
		Registry:    rt.Generators,
		Plugins:     rt.Plugins,
		Logger:      logger,
		ToolVersion: toolVersion,
	})
	if err != nil {
		PrintError(err.Error())
		return 1
	}

	// Re-run post-gen commands so any new dependencies are installed.
	plan := PlanPostGenCommands(res.Spec, res.Manifests)
	if len(plan) > 0 {
		PrintHeading(fmt.Sprintf("post-gen commands (%d)", len(plan)))
		runner := commands.NewRunner(res.ProjectRoot, logger)
		if err := RunCommandsQuiet(ctx, runner, plan, 2); err != nil {
			PrintError("post-gen failed: " + err.Error())
			return 1
		}
	}

	PrintSuccess(fmt.Sprintf("updated %s", res.ProjectRoot))
	return 0
}
