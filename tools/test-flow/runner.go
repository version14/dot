package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/version14/dot/flows"
	"github.com/version14/dot/internal/cli"
	"github.com/version14/dot/internal/commands"
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/pkg/dotapi"
)

// invocationNames returns just the names from a slice of Invocations, in order.
func invocationNames(invs []generator.Invocation) []string {
	out := make([]string, len(invs))
	for i, inv := range invs {
		out[i] = inv.Name
	}
	return out
}

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

// askLoop runs the full body sub-graph for each scripted iteration using a
// fresh FlowEngine. This mirrors HuhFormRunner.runLoopSubForms: conditional
// questions in the body are followed/skipped based on the iteration's answers.
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

	for i, iter := range iters {
		iterMap, ok := iter.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("test-flow: loop %q iteration %d must be an object", loop.ID(), i)
		}

		// Layer iteration answers over global ones; each body question lookup
		// hits the iteration map first, then falls back to outer answers.
		iterAdapter := &scriptedAdapter{answers: mergeAnswerMaps(a.answers, iterMap)}
		iterEng := flow.NewEngine(iterAdapter)

		iterAnswers := make(map[string]flow.Answer)
		for _, bodyRoot := range loop.Body {
			bodyCtx, err := iterEng.Run(bodyRoot)
			if err != nil {
				return nil, fmt.Errorf("loop %q iter %d body: %w", loop.ID(), i, err)
			}
			for k, v := range bodyCtx.Answers {
				iterAnswers[k] = v
			}
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
	noCache          bool   // when true, ignore cache hits and refresh entries
	caseFile         string // absolute path to the testdata JSON for this case
	repoRoot         string // absolute path to the dot repo root
	flowsDir         string // absolute path to the flows/ directory
}

// runOne drives one TestCase through the full pipeline:
//
//	flow → spec → generators → persist → validators → post-gen → test commands
//
// Each step is reported via rep, keyed by caseIdx (the case's stable position
// in the original cases slice) so concurrent workers report into the right
// row/case. The function returns a Result the caller passes to Summarize.
// Any per-step failure is captured in Result; the function does not panic or
// os.Exit.
func runOne(
	ctx context.Context,
	caseIdx int,
	tc *TestCase,
	def *flows.FlowDef,
	rt *cli.Runtime,
	rep Reporter,
	opts caseOptions,
) *Result {
	r := &Result{Case: tc}
	rep.CaseStart(caseIdx, tc.Name, tc.FlowID)

	// Step 1: scaffold (flow → generators → persist → .dot files).
	scratch, err := os.MkdirTemp(opts.tempDirRoot, "dot-test-"+tc.FlowID+"-*")
	if err != nil {
		r.Err = fmt.Errorf("mkdir temp: %w", err)
		rep.Step(caseIdx, stepKeyFlow, "scaffold", false, "", err)
		return r
	}
	defer func() {
		if opts.keepScratch {
			rep.Step(caseIdx, stepKeyFiles, "scratch dir kept", true, scratch, nil)
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
		rep.Step(caseIdx, stepKeyFlow, "scaffold", false, time.Since(scaffoldStart).String(), err)
		return r
	}
	r.Scaffold = res
	r.ProjectRoot = res.ProjectRoot

	rep.Step(caseIdx, stepKeyFlow, "flow", true, fmt.Sprintf(""), nil)

	if len(tc.ExpectedIDs) > 0 && !equalStringSlice(tc.ExpectedIDs, res.Spec.VisitedNodes) {
		r.Diffs = append(r.Diffs, fmt.Sprintf(
			"visited mismatch:\n      expected: %v\n      actual:   %v",
			tc.ExpectedIDs, res.Spec.VisitedNodes,
		))
		rep.Step(caseIdx, stepKeyVerify, "verify visited", false, "", fmt.Errorf("mismatch"))
	} else if len(tc.ExpectedIDs) > 0 {
		rep.Step(caseIdx, stepKeyVerify, "verify visited", true, "", nil)
	}

	rep.Step(caseIdx, stepKeyResolved, "resolved generators", true, "", nil)
	rep.Step(caseIdx, stepKeyFiles, "scaffolded files", true, "", nil)

	// Step 2: validators (run against the on-disk project).
	failures, err := generator.RunValidators(res.ProjectRoot, res.Manifests)
	if err != nil {
		r.Err = fmt.Errorf("validators: %w", err)
		rep.Step(caseIdx, stepKeyValidate, "validators", false, "", err)
		return r
	}
	if len(failures) > 0 {
		for _, f := range failures {
			r.Diffs = append(r.Diffs, "validator: "+f.String())
		}
		rep.Step(caseIdx, stepKeyValidate, "validators", false, fmt.Sprintf("%d failures", len(failures)), nil)
	} else {
		rep.Step(caseIdx, stepKeyValidate, "validators", true, "", nil)
	}

	// Step 2.5: case-level cache check. Skips post-gen + test commands when
	// the fingerprint matches a previous successful run AND every command
	// across the involved manifests is opted-in via Cacheable=true.
	cacheHit := false
	fingerprint, fpErr := ComputeFingerprint(CacheKeyInputs{
		CaseFile:      opts.caseFile,
		FlowsDir:      opts.flowsDir,
		Invocations:   res.Invocations,
		Manifests:     res.Manifests,
		SkipPostFlag:  opts.skipPostCommands || tc.SkipPostCommands,
		SkipTestFlag:  opts.skipTestCommands || tc.SkipTestCommands,
		GeneratorsDir: filepath.Join(opts.repoRoot, "generators"),
		RepoRoot:      opts.repoRoot,
	})

	if fpErr != nil {
		rep.Step(caseIdx, stepKeyCache, "cache fingerprint", false, "", fpErr)
	} else if !opts.noCache {
		entry, err := LoadCacheEntry(opts.repoRoot, tc.Name)
		if err != nil {
			rep.Step(caseIdx, stepKeyCache, "cache load", false, "", err)
		}
		if entry != nil && entry.Fingerprint == fingerprint && AllCommandsCacheable(res.Manifests) {
			rep.Step(caseIdx, stepKeyCache, "cache", true, "", nil)
			cacheHit = true
		} else if entry != nil && entry.Fingerprint == fingerprint && !AllCommandsCacheable(res.Manifests) {
			blocking := NonCacheableCommands(res.Manifests)
			detail := fmt.Sprintf("%d non-cacheable command(s) — running anyway", len(blocking))
			rep.Step(caseIdx, stepKeyCache, "cache", true, detail, nil)
		}
	}

	// Step 3: post-generation commands (skipped on cache hit).
	if cacheHit {
		// Cache hit short-circuits both post-gen and test commands.
	} else if !opts.skipPostCommands && !tc.SkipPostCommands {
		postPlan := cli.PlanPostGenCommands(res.Spec, res.Manifests)
		if len(postPlan) > 0 {
			rep.Substep(caseIdx, stepKeyPost, "post-gen commands", len(postPlan))
			if output, cmdErr := runCommandList(ctx, res.ProjectRoot, postPlan, rep, caseIdx, stepKeyPost); cmdErr != nil {
				r.Diffs = append(r.Diffs, "post-gen: "+cmdErr.Error())
				if diff := formatCapturedOutputDiff(output); diff != "" {
					r.Diffs = append(r.Diffs, diff)
				}
			}
		}
	} else {
		rep.Step(caseIdx, stepKeyPost, "post-gen commands", true, "skipped", nil)
	}

	// Step 4: test commands (incl. background dev servers).
	if cacheHit {
		// see above
	} else if !opts.skipTestCommands && !tc.SkipTestCommands {
		testPlan := cli.PlanTestCommands(res.Spec, res.Manifests)
		if len(testPlan) > 0 {
			rep.Substep(caseIdx, stepKeyTest, "test commands", len(testPlan))
			if output, cmdErr := runCommandList(ctx, res.ProjectRoot, testPlan, rep, caseIdx, stepKeyTest); cmdErr != nil {
				r.Diffs = append(r.Diffs, "test: "+cmdErr.Error())
				if diff := formatCapturedOutputDiff(output); diff != "" {
					r.Diffs = append(r.Diffs, diff)
				}
			}
		}
	} else {
		rep.Step(caseIdx, stepKeyTest, "test commands", true, "skipped", nil)
	}

	// Persist a fresh cache entry on full success. Failed runs intentionally
	// leave no trace so the next invocation retries them.
	if r.Pass() && fingerprint != "" && AllCommandsCacheable(res.Manifests) {
		entry := CacheEntry{
			SchemaVersion: cacheSchemaVersion,
			Fingerprint:   fingerprint,
			CaseName:      tc.Name,
			FlowID:        tc.FlowID,
			LastSuccessAt: time.Now().UTC().Format(time.RFC3339),
			Generators:    invocationNames(res.Invocations),
		}
		if err := SaveCacheEntry(opts.repoRoot, entry); err != nil {
			rep.Step(caseIdx, stepKeyCache, "cache save", false, "", err)
		}
	}

	rep.CaseEnd(caseIdx, r.Pass())
	return r
}

// runCommandList executes each PlannedCommand in order, reporting per-command
// Start/Done events to rep under (caseIdx, key) — that's what lets a table
// reporter show "what's running now" instead of a scrolling command log.
// cli.RunCommandsWithProgress is shared with `dot scaffold`'s spinner UX; here
// it's driven by reporterProgress instead.
//
// On failure it returns the captured combined stdout/stderr of the failing
// command, when the caller should hold onto it for later (see reporterProgress
// docs below) — otherwise nil, since it's already been printed live.
func runCommandList(
	ctx context.Context,
	projectRoot string,
	plan []commands.PlannedCommand,
	rep Reporter,
	caseIdx int,
	key string,
) ([]byte, error) {
	runner := commands.NewRunner(projectRoot, dotapi.DiscardLogger{})
	// TableReporter owns the terminal (raw mode, its own redraw loop) for the
	// life of the run — a direct fmt.Fprint to stderr from here would race its
	// repaints and corrupt the screen (this is exactly what garbled the TUI
	// when a command failed: PrintCapturedOutput used to always write straight
	// to stderr). So table runs hold onto the output and let runOne fold it
	// into the Result's Diffs, which only get printed after the table quits
	// and hands the terminal back. The plain reporter has no such conflict —
	// it still gets the output printed live, as before.
	_, isTable := rep.(*TableReporter)
	progress := &reporterProgress{rep: rep, caseIdx: caseIdx, key: key, live: !isTable}
	err := cli.RunCommandsWithProgress(ctx, runner, plan, progress)
	return progress.deferredOutput, err
}

// reporterProgress adapts a Reporter to cli.CommandProgress so per-command
// lifecycle events flow into the same case/column the rest of runOne reports
// into.
type reporterProgress struct {
	rep     Reporter
	caseIdx int
	key     string
	live    bool // true: safe to print captured output to stderr immediately

	deferredOutput []byte // set on failure when !live
}

func (p *reporterProgress) Start(c commands.PlannedCommand) {
	p.rep.SubStart(p.caseIdx, p.key, commandLabel(c))
}

func (p *reporterProgress) Done(c commands.PlannedCommand, elapsed time.Duration, output []byte, err error) {
	label := commandLabel(c)
	if err != nil {
		p.rep.Sub(p.caseIdx, p.key, label, false, elapsed.String(), err)
		if p.live {
			cli.PrintCapturedOutput(output, 6)
		} else {
			p.deferredOutput = output
		}
		return
	}
	p.rep.Sub(p.caseIdx, p.key, label, true, elapsed.String(), nil)
}

// commandLabel is the plain-text (unstyled) label for one command — callers
// apply their own styling/truncation.
func commandLabel(c commands.PlannedCommand) string {
	if c.Background {
		return c.Cmd + " (background)"
	}
	return c.Cmd
}

// formatCapturedOutputDiff renders a failing command's captured output as one
// Result.Diffs entry, capped so a runaway command can't blow up the summary.
func formatCapturedOutputDiff(output []byte) string {
	trimmed := strings.TrimRight(string(output), "\n")
	if trimmed == "" {
		return ""
	}
	const maxLines = 40
	lines := strings.Split(trimmed, "\n")
	if len(lines) > maxLines {
		omitted := len(lines) - maxLines
		lines = append([]string{fmt.Sprintf("… (%d earlier lines omitted)", omitted)}, lines[len(lines)-maxLines:]...)
	}
	return "output:\n      " + strings.Join(lines, "\n      ")
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
