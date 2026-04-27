package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/version14/dot/flows"
	"github.com/version14/dot/internal/cli"
	"github.com/version14/dot/internal/commands"
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/pkg/dotapi"
)

// scriptedAdapter answers each question from a recorded map.
//
// LoopQuestion handling: when the scripted answer for a loop is a JSON array
// (i.e. []interface{} after unmarshal), the adapter treats each element as
// the answer-map for one iteration of the loop body. It walks the body in
// order, looking up each child question's answer from the per-iteration map.
// This mirrors HuhFormRunner.runLoopSubForms but synchronously and without
// any UI, so test-flow can exercise loop-using flows from JSON fixtures.
type scriptedAdapter struct {
	answers map[string]flow.Answer
}

func newScriptedAdapter(answers map[string]flow.Answer) *scriptedAdapter {
	return &scriptedAdapter{answers: answers}
}

func (a *scriptedAdapter) Ask(q flow.Question, ctx *flow.FlowContext) (flow.Answer, error) {
	if loop, ok := q.(*flow.LoopQuestion); ok {
		return a.askLoop(loop, ctx)
	}

	id := q.ID()
	ans, ok := a.answers[id]
	if !ok {
		return nil, fmt.Errorf("test-flow: no scripted answer for question %q", id)
	}
	return ans, nil
}

// askLoop iterates the loop body once per scripted iteration, returning a
// []map[string]flow.Answer ready to be stored on FlowContext.Answers[loop.ID()].
func (a *scriptedAdapter) askLoop(loop *flow.LoopQuestion, ctx *flow.FlowContext) (flow.Answer, error) {
	raw, ok := a.answers[loop.ID()]
	if !ok {
		return nil, fmt.Errorf("test-flow: no scripted iterations for loop %q", loop.ID())
	}
	iters, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("test-flow: loop %q expects an array of objects, got %T", loop.ID(), raw)
	}

	out := make([]map[string]flow.Answer, len(iters))
	prev := a.answers
	defer func() { a.answers = prev }()

	for i, iter := range iters {
		iterMap, ok := iter.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("test-flow: loop %q iteration %d must be an object", loop.ID(), i)
		}
		// Layer the iteration's answers over the global ones so body-question
		// lookups resolve correctly and outer answers stay reachable.
		a.answers = mergeAnswerMaps(prev, iterMap)

		iterAnswers := make(map[string]flow.Answer, len(loop.Body))
		for _, body := range loop.Body {
			bodyAns, err := a.Ask(body, ctx)
			if err != nil {
				return nil, fmt.Errorf("loop %q iter %d: %w", loop.ID(), i, err)
			}
			iterAnswers[body.ID()] = bodyAns
		}
		out[i] = iterAnswers
	}
	return out, nil
}

func mergeAnswerMaps(base, overlay map[string]flow.Answer) map[string]flow.Answer {
	merged := make(map[string]flow.Answer, len(base)+len(overlay))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range overlay {
		merged[k] = v
	}
	return merged
}

// scriptedRunner implements flow.FlowRunner by running the flow.FlowEngine
// against a scripted adapter — no terminal interaction.
//
// Plugin injections fire via the supplied HookRegistry (and FragmentRegistry),
// which means inserted/replaced/added-option questions show up in the
// engine's traversal exactly as they would in the interactive HuhFormRunner.
type scriptedRunner struct {
	adapter   *scriptedAdapter
	hooks     *flow.HookRegistry
	fragments *flow.FragmentRegistry
}

func newScriptedRunner(
	answers map[string]flow.Answer,
	hooks *flow.HookRegistry,
	fragments *flow.FragmentRegistry,
) *scriptedRunner {
	return &scriptedRunner{
		adapter:   newScriptedAdapter(answers),
		hooks:     hooks,
		fragments: fragments,
	}
}

func (r *scriptedRunner) Run(root flow.Question) (*flow.FlowContext, error) {
	eng := flow.NewEngine(r.adapter)
	if r.hooks != nil {
		eng.Hooks = r.hooks
	}
	if r.fragments != nil {
		eng.Fragments = r.fragments
	}
	return eng.Run(root)
}

// caseOptions controls how runOne executes one test case.
type caseOptions struct {
	tempDirRoot      string // parent dir for the per-case scratch dir
	skipPostCommands bool   // skip PostGenerationCommands globally
	skipTestCommands bool   // skip TestCommands globally
	keepScratch      bool   // when true, do NOT delete the scratch dir on exit
}

// runOne drives one TestCase through the full pipeline:
//
//	flow → spec → generators → persist → validators → post-gen → test commands
//
// Each step is logged via the StepReporter. The function returns a Result
// the caller passes to Summarize. Any per-step failure is captured in Result;
// the function does not panic or os.Exit.
func runOne(
	ctx context.Context,
	tc *TestCase,
	def *flows.FlowDef,
	rt *cli.Runtime,
	rep *StepReporter,
	opts caseOptions,
) *Result {
	r := &Result{Case: tc}
	rep.CaseStart(tc.Name, tc.FlowID)

	// Step 1: scaffold (flow → generators → persist → .dot files).
	scratch, err := os.MkdirTemp(opts.tempDirRoot, "dot-test-"+tc.FlowID+"-*")
	if err != nil {
		r.Err = fmt.Errorf("mkdir temp: %w", err)
		rep.Step("scaffold", false, "", err)
		return r
	}
	defer func() {
		if opts.keepScratch {
			rep.Step("scratch dir kept", true, scratch, nil)
			return
		}
		_ = os.RemoveAll(scratch)
	}()

	scaffoldStart := time.Now()
	res, err := cli.Scaffold(ctx, cli.ScaffoldOptions{
		Flow:        def,
		Registry:    rt.Generators,
		Hooks:       rt.Hooks,
		Fragments:   rt.Fragments,
		Plugins:     rt.Plugins,
		OutputDir:   scratch,
		ToolVersion: "test-flow",
		Logger:      dotapi.DiscardLogger{}, // step logging is the reporter's job
		Runner:      newScriptedRunner(tc.Answers, rt.Hooks, rt.Fragments),
	})
	if err != nil {
		r.Err = fmt.Errorf("scaffold: %w", err)
		rep.Step("scaffold", false, time.Since(scaffoldStart).String(), err)
		return r
	}
	r.Scaffold = res
	r.ProjectRoot = res.ProjectRoot

	rep.Step("flow", true, fmt.Sprintf("%d nodes visited", len(res.Spec.VisitedNodes)), nil)

	if len(tc.ExpectedIDs) > 0 && !equalStringSlice(tc.ExpectedIDs, res.Spec.VisitedNodes) {
		r.Diffs = append(r.Diffs, fmt.Sprintf(
			"visited mismatch:\n      expected: %v\n      actual:   %v",
			tc.ExpectedIDs, res.Spec.VisitedNodes,
		))
		rep.Step("verify visited", false, "", fmt.Errorf("mismatch"))
	} else if len(tc.ExpectedIDs) > 0 {
		rep.Step("verify visited", true, "matches expected", nil)
	}

	rep.Step("resolved generators", true, fmt.Sprintf("%s", joinNames(res.Invocations)), nil)
	rep.Step("scaffolded files", true, fmt.Sprintf("→ %s", res.ProjectRoot), nil)

	// Step 2: validators (run against the on-disk project).
	failures, err := generator.RunValidators(res.ProjectRoot, res.Manifests)
	if err != nil {
		r.Err = fmt.Errorf("validators: %w", err)
		rep.Step("validators", false, "", err)
		return r
	}
	if len(failures) > 0 {
		for _, f := range failures {
			r.Diffs = append(r.Diffs, "validator: "+f.String())
		}
		rep.Step("validators", false, fmt.Sprintf("%d failures", len(failures)), nil)
	} else {
		rep.Step("validators", true, fmt.Sprintf("%d passed", countChecks(res.Manifests)), nil)
	}

	// Step 3: post-generation commands.
	if !opts.skipPostCommands && !tc.SkipPostCommands {
		postPlan := cli.PlanPostGenCommands(res.Spec, res.Manifests)
		if len(postPlan) > 0 {
			rep.Substep("post-gen commands", len(postPlan))
			if cmdErr := runCommandList(ctx, res.ProjectRoot, postPlan); cmdErr != nil {
				r.Diffs = append(r.Diffs, "post-gen: "+cmdErr.Error())
			}
		}
	} else {
		rep.Step("post-gen commands", true, "skipped", nil)
	}

	// Step 4: test commands (incl. background dev servers).
	if !opts.skipTestCommands && !tc.SkipTestCommands {
		testPlan := cli.PlanTestCommands(res.Spec, res.Manifests)
		if len(testPlan) > 0 {
			rep.Substep("test commands", len(testPlan))
			if cmdErr := runCommandList(ctx, res.ProjectRoot, testPlan); cmdErr != nil {
				r.Diffs = append(r.Diffs, "test: "+cmdErr.Error())
			}
		}
	} else {
		rep.Step("test commands", true, "skipped", nil)
	}

	rep.CaseEnd(r.Pass())
	return r
}

// runCommandList executes each PlannedCommand in order with a Docker-style
// spinner UX (animated while running, ✓/✗ + elapsed when done, full output
// only on failure). Implementation lives in cli.RunCommandsQuiet so the same
// behaviour is shared with `dot scaffold`.
//
// The reporter's Substep header is printed by the caller; this function only
// renders the per-command lines (indented 4 spaces to nest under the header).
func runCommandList(
	ctx context.Context,
	projectRoot string,
	plan []commands.PlannedCommand,
) error {
	runner := commands.NewRunner(projectRoot, dotapi.DiscardLogger{})
	return cli.RunCommandsQuiet(ctx, runner, plan, 4)
}

func countChecks(mans []dotapi.Manifest) int {
	n := 0
	for _, m := range mans {
		for _, v := range m.Validators {
			n += len(v.Checks)
		}
	}
	return n
}

func joinNames(invs []generator.Invocation) string {
	out := ""
	for i, inv := range invs {
		if i > 0 {
			out += ", "
		}
		out += inv.Name
	}
	return out
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
