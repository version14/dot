//go:build integration

// Package integration runs every flow end-to-end.
//
// Each flow is declared as a JSON file in testdata/. The JSON describes the
// user's *answers* to the survey (the path they took through the question
// tree), not a pre-resolved list of generators. The test rebuilds a
// scaffold.Result from those answers, runs scaffold.Collect against the live
// question tree (templates.StarterQuestions) to resolve which generators
// should fire and with which scoped specs, then scaffolds, installs, type-
// checks, starts the dev server, and probes the health endpoint.
//
// This mirrors production (cmd_init) exactly — so any regression in the
// survey-to-generators resolution (comparison logic, loop scoping, namespace
// collisions) shows up in these tests.
//
// Run with:
//
//	make test-flows
//
// Add a new flow by dropping a JSON file into testdata/ — no Go code required.
package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/pipeline"
	"github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
	"github.com/version14/dot/internal/templates"
)

const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[0;31m"
	ansiGreen  = "\033[0;32m"
	ansiCyan   = "\033[0;36m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[2m"
	ansiPurple = "\033[0;35m"
	ansiYellow = "\033[0;33m"
)

const (
	serverStartTimeout  = 20 * time.Second
	healthCheckInterval = 500 * time.Millisecond
	healthCheckTimeout  = 2 * time.Second
	defaultHealthURL    = "http://localhost:3000/health"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// flow definition (JSON)

type flowDef struct {
	Name    string      `json:"name"`
	Port    int         `json:"port"`
	Answers []answerDef `json:"answers"`
	Expect  expectBlock `json:"expect,omitempty"`
}

// answerDef is one entry in the user's survey path.
// Either Value or Multi is set for leaf answers. Iterations is set for loops.
type answerDef struct {
	Key        string        `json:"key"`
	Value      string        `json:"value,omitempty"`
	Multi      []string      `json:"multi,omitempty"`
	Iterations [][]answerDef `json:"iterations,omitempty"`
}

// expectBlock is an optional assertion layer. When present, the test checks
// that exactly these generators fired and that every listed file exists.
type expectBlock struct {
	Generators []string `json:"generators,omitempty"`
	Files      []string `json:"files,omitempty"`
}

// flow (runtime)

type flow struct {
	name        string
	activations []scaffold.Activation
	port        int
	expect      expectBlock
}

func loadFlows(t *testing.T) []flow {
	t.Helper()
	paths, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Fatalf("glob testdata: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("no flow definitions found in testdata/")
	}
	sort.Strings(paths)

	flows := make([]flow, 0, len(paths))
	for _, path := range paths {
		flows = append(flows, parseFlowFile(t, path))
	}
	return flows
}

func parseFlowFile(t *testing.T, path string) flow {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var def flowDef
	if err := json.Unmarshal(data, &def); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return resolveFlow(t, def, path)
}

// resolveFlow drives the same code path as cmd/dot/cmd_init:
//   - rebuild a scaffold.Result from the declared answers
//   - compute the base spec.Spec
//   - call scaffold.Collect against the live survey tree
func resolveFlow(t *testing.T, def flowDef, path string) flow {
	t.Helper()

	result := buildResult(def.Answers)

	// Same base-spec derivation as cmd_init: only top-level (non-iteration)
	// entries go into Extensions; loop answers reach generators via
	// Activation.Spec per iteration.
	extensions := make(map[string]any)
	for _, e := range result.Entries {
		if len(e.Iterations) > 0 {
			continue
		}
		if len(e.Multi) > 0 {
			extensions[e.Key] = e.Multi
		} else {
			extensions[e.Key] = e.Value
		}
	}

	// Mirror cmd_init: only project-name is derived (root of the survey).
	// Language/type are plugin-extensible and never hardcoded here.
	base := spec.Spec{
		Project:    spec.ProjectSpec{Name: result.Get("project-name")},
		Extensions: extensions,
	}

	acts := scaffold.Collect(templates.StarterQuestions, result, base)
	if len(acts) == 0 {
		t.Fatalf("%s: no activations — answers do not match any generator path", filepath.Base(path))
	}

	return flow{
		name:        def.Name,
		activations: acts,
		port:        def.Port,
		expect:      def.Expect,
	}
}

// buildResult re-creates a scaffold.Result from declared answers, recursing
// into iterations so loop activations carry their own answers.
func buildResult(answers []answerDef) *scaffold.Result {
	result := &scaffold.Result{}
	for _, a := range answers {
		result.Add(toEntry(a))
	}
	return result
}

func toEntry(a answerDef) scaffold.AnswerEntry {
	entry := scaffold.AnswerEntry{Key: a.Key}
	switch {
	case len(a.Iterations) > 0:
		entry.Value = a.Value
		entry.Iterations = make([][]scaffold.AnswerEntry, len(a.Iterations))
		for i, iter := range a.Iterations {
			iterEntries := make([]scaffold.AnswerEntry, 0, len(iter))
			for _, sub := range iter {
				iterEntries = append(iterEntries, toEntry(sub))
			}
			entry.Iterations[i] = iterEntries
		}
	case len(a.Multi) > 0:
		entry.Multi = a.Multi
	default:
		entry.Value = a.Value
	}
	return entry
}

// step display

type stepState int

const (
	stepPending stepState = iota
	stepRunning
	stepDone
	stepFailed
)

type displayStep struct {
	phase   string
	label   string
	state   stepState
	dur     time.Duration
	spinner string
}

type flowLog struct {
	mu    sync.Mutex
	steps []*displayStep
}

func (l *flowLog) header(name string, idx, total int) {
	fmt.Fprintf(os.Stderr, "\n%s%s(%d/%d)%s %s%s%s\n",
		ansiCyan, ansiBold, idx, total, ansiReset,
		ansiBold, name, ansiReset,
	)
}

func (l *flowLog) initSteps(steps []*displayStep) {
	l.steps = steps
	for _, s := range steps {
		fmt.Fprintf(os.Stderr, "  %s○%s  %s%-9s%s  %s\n",
			ansiDim, ansiReset,
			ansiPurple, s.phase, ansiReset,
			s.label,
		)
	}
}

func (l *flowLog) setRunning(idx int) (stop func()) {
	l.mu.Lock()
	l.steps[idx].state = stepRunning
	l.steps[idx].spinner = spinnerFrames[0]
	l.redrawStep(idx)
	l.mu.Unlock()

	quit := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		t := time.NewTicker(80 * time.Millisecond)
		defer t.Stop()
		frame := 1
		for {
			select {
			case <-quit:
				return
			case <-t.C:
				l.mu.Lock()
				l.steps[idx].spinner = spinnerFrames[frame%len(spinnerFrames)]
				l.redrawStep(idx)
				l.mu.Unlock()
				frame++
			}
		}
	}()
	return func() {
		close(quit)
		<-done
	}
}

func (l *flowLog) setDone(idx int, dur time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.steps[idx].state = stepDone
	l.steps[idx].dur = dur
	l.redrawStep(idx)
}

func (l *flowLog) setFailed(idx int, dur time.Duration, stderr string) {
	l.mu.Lock()
	l.steps[idx].state = stepFailed
	l.steps[idx].dur = dur
	l.redrawStep(idx)
	l.mu.Unlock()

	for _, line := range strings.Split(strings.TrimSpace(stderr), "\n") {
		if line != "" {
			fmt.Fprintf(os.Stderr, "         %s%s%s\n", ansiRed, line, ansiReset)
		}
	}
}

func (l *flowLog) redrawStep(idx int) {
	n := len(l.steps)
	step := l.steps[idx]
	linesUp := n - idx

	fmt.Fprintf(os.Stderr, "\033[%dA\r\033[2K", linesUp)

	switch step.state {
	case stepPending:
		fmt.Fprintf(os.Stderr, "  %s○%s  %s%-9s%s  %s",
			ansiDim, ansiReset, ansiPurple, step.phase, ansiReset, step.label)
	case stepRunning:
		fmt.Fprintf(os.Stderr, "  %s%s%s  %s%-9s%s  %s",
			ansiYellow, step.spinner, ansiReset, ansiPurple, step.phase, ansiReset, step.label)
	case stepDone:
		fmt.Fprintf(os.Stderr, "  %s✓%s  %s%-9s%s  %s  %s%s%s",
			ansiGreen, ansiReset, ansiPurple, step.phase, ansiReset, step.label,
			ansiDim, step.dur.Round(time.Millisecond), ansiReset)
	case stepFailed:
		fmt.Fprintf(os.Stderr, "  %s✗%s  %s%-9s%s  %s  %s%s%s",
			ansiRed, ansiReset, ansiPurple, step.phase, ansiReset, step.label,
			ansiDim, step.dur.Round(time.Millisecond), ansiReset)
	}

	fmt.Fprintf(os.Stderr, "\033[%dB\r", linesUp)
}

func (l *flowLog) summary(elapsed time.Duration, passed bool) {
	if passed {
		fmt.Fprintf(os.Stderr, "  %s%s✓ passed%s  %s%s%s\n",
			ansiGreen, ansiBold, ansiReset,
			ansiDim, elapsed.Round(time.Millisecond), ansiReset,
		)
	} else {
		fmt.Fprintf(os.Stderr, "  %s%s✗ failed%s\n", ansiRed, ansiBold, ansiReset)
	}
}

// test

func TestFlows(t *testing.T) {
	requirePnpm(t)

	flows := loadFlows(t)
	dirs := preallocateDirs(t, len(flows))
	bg := &backgroundProcesses{}
	t.Cleanup(bg.killAll)

	var failures []string

	for i, fl := range flows {
		log := &flowLog{}
		log.header(fl.name, i+1, len(flows))

		runner := &flowRunner{flow: fl, dir: dirs[i], log: log, bg: bg}
		start := time.Now()
		err := runner.execute()

		log.summary(time.Since(start), err == nil)
		if err != nil {
			failures = append(failures, fmt.Sprintf("  %s: %v", fl.name, err))
		}
	}

	fmt.Fprintln(os.Stderr)
	if len(failures) > 0 {
		t.Fatalf("flows failed:\n%s", strings.Join(failures, "\n"))
	}
}

func requirePnpm(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("pnpm"); err != nil {
		t.Skip("pnpm not found in PATH")
	}
}

func preallocateDirs(t *testing.T, n int) []string {
	t.Helper()
	dirs := make([]string, n)
	for i := range dirs {
		dirs[i] = t.TempDir()
	}
	return dirs
}

// flow runner

type flowRunner struct {
	flow    flow
	dir     string
	log     *flowLog
	bg      *backgroundProcesses
	stepIdx int
}

func (runner *flowRunner) execute() error {
	fileOps, postOps, activatedNames, err := runner.collectOps()
	if err != nil {
		return err
	}

	if err := runner.assertExpect(activatedNames); err != nil {
		return err
	}

	runner.log.initSteps(runner.buildSteps(postOps))

	if err := runner.runScaffold(fileOps); err != nil {
		return err
	}
	if err := runner.assertFiles(); err != nil {
		return err
	}
	if err := runner.runPhase(postOps, generator.PhaseInstall); err != nil {
		return err
	}
	if err := runner.runPhase(postOps, generator.PhaseTypeCheck); err != nil {
		return err
	}
	return runner.runSmoke(postOps)
}

// collectOps runs each activation's Fn with its scoped Spec and returns the
// aggregated file/post-ops plus the list of activation identifiers (for
// assertExpect).
func (runner *flowRunner) collectOps() ([]generator.FileOp, []generator.PostOp, []string, error) {
	var fileOps []generator.FileOp
	var postOps []generator.PostOp
	var names []string
	for _, activation := range runner.flow.activations {
		fops, pops, err := activation.Fn(activation.Spec)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("activation [%s=%s]: %w", activation.QuestionKey, activation.AnswerValue, err)
		}
		fileOps = append(fileOps, fops...)
		postOps = append(postOps, pops...)
		names = append(names, activation.QuestionKey+"="+activation.AnswerValue)
	}
	return fileOps, postOps, names, nil
}

// assertExpect checks the optional expect.generators list (subset match on
// "<key>=<value>" activation IDs).
func (runner *flowRunner) assertExpect(activated []string) error {
	if len(runner.flow.expect.Generators) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(activated))
	for _, a := range activated {
		seen[a] = struct{}{}
	}
	var missing []string
	for _, want := range runner.flow.expect.Generators {
		if _, ok := seen[want]; !ok {
			missing = append(missing, want)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("expected activations not produced: %v (got %v)", missing, activated)
	}
	return nil
}

func (runner *flowRunner) assertFiles() error {
	for _, rel := range runner.flow.expect.Files {
		if _, err := os.Stat(filepath.Join(runner.dir, rel)); err != nil {
			return fmt.Errorf("expected file not generated: %s", rel)
		}
	}
	return nil
}

func (runner *flowRunner) buildSteps(postOps []generator.PostOp) []*displayStep {
	steps := []*displayStep{
		{phase: "scaffold", label: "generate files"},
	}
	for _, phase := range []generator.PostOpPhase{generator.PhaseInstall, generator.PhaseTypeCheck} {
		for _, operation := range postOps {
			if opMatchesPhase(operation, phase) && !operation.Background {
				steps = append(steps, &displayStep{
					phase: string(phase),
					label: operation.Command + " " + strings.Join(withPort(operation.Args, runner.flow.port), " "),
				})
			}
		}
	}
	for _, operation := range postOps {
		if opMatchesPhase(operation, generator.PhaseSmoke) && operation.Background {
			steps = append(steps, &displayStep{
				phase: "smoke",
				label: fmt.Sprintf("dev server  http://localhost:%d/health", runner.flow.port),
			})
		}
	}
	for _, operation := range postOps {
		if opMatchesPhase(operation, generator.PhaseSmoke) && !operation.Background {
			steps = append(steps, &displayStep{
				phase: "smoke",
				label: operation.Command + " " + strings.Join(withPort(operation.Args, runner.flow.port), " "),
			})
		}
	}
	return steps
}

func (runner *flowRunner) runScaffold(fileOps []generator.FileOp) error {
	idx := runner.stepIdx
	runner.stepIdx++

	stop := runner.log.setRunning(idx)
	start := time.Now()
	err := pipeline.RunIn(runner.dir, fileOps)
	stop()

	if err != nil {
		runner.log.setFailed(idx, time.Since(start), err.Error())
		return fmt.Errorf("pipeline: %w", err)
	}
	runner.log.setDone(idx, time.Since(start))
	return nil
}

func (runner *flowRunner) runPhase(ops []generator.PostOp, phase generator.PostOpPhase) error {
	for _, operation := range ops {
		if !opMatchesPhase(operation, phase) || operation.Background {
			continue
		}
		if err := runner.runOp(operation); err != nil {
			return err
		}
	}
	return nil
}

func (runner *flowRunner) runSmoke(ops []generator.PostOp) error {
	if err := runner.startBackgroundServers(ops); err != nil {
		return err
	}
	return runner.runForegroundChecks(ops)
}

func (runner *flowRunner) startBackgroundServers(ops []generator.PostOp) error {
	for _, operation := range ops {
		if !opMatchesPhase(operation, generator.PhaseSmoke) || !operation.Background {
			continue
		}
		if err := runner.startServer(operation); err != nil {
			return err
		}
		if err := runner.waitUntilHealthy(); err != nil {
			return err
		}
	}
	return nil
}

func (runner *flowRunner) startServer(op generator.PostOp) error {
	cmd := exec.Command(op.Command, op.Args...)
	cmd.Dir = filepath.Join(runner.dir, op.Dir)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", runner.flow.port))
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start dev server: %w", err)
	}
	runner.bg.add(cmd)
	return nil
}

func (runner *flowRunner) waitUntilHealthy() error {
	healthURL := fmt.Sprintf("http://localhost:%d/health", runner.flow.port)
	idx := runner.stepIdx
	runner.stepIdx++

	stop := runner.log.setRunning(idx)
	start := time.Now()
	ready := pollUntilReady(healthURL, serverStartTimeout)
	stop()

	if !ready {
		runner.log.setFailed(idx, time.Since(start), "server did not start within 20s")
		return fmt.Errorf("server at %s did not start within 20s", healthURL)
	}
	runner.log.setDone(idx, time.Since(start))
	return nil
}

func (runner *flowRunner) runForegroundChecks(ops []generator.PostOp) error {
	for _, operation := range ops {
		if !opMatchesPhase(operation, generator.PhaseSmoke) || operation.Background {
			continue
		}
		if err := runner.runOp(operation); err != nil {
			return err
		}
	}
	return nil
}

func (runner *flowRunner) runOp(operation generator.PostOp) error {
	args := withPort(operation.Args, runner.flow.port)
	phase := phaseNameOf(operation)
	label := operation.Command + " " + strings.Join(args, " ")
	idx := runner.stepIdx
	runner.stepIdx++

	stop := runner.log.setRunning(idx)
	start := time.Now()
	stderr, err := runner.execCommand(operation.Command, args, operation.Dir)
	stop()

	duration := time.Since(start)
	if err != nil {
		runner.log.setFailed(idx, duration, strings.TrimSpace(stderr))
		return fmt.Errorf("[%s] %s: %w", phase, label, err)
	}
	runner.log.setDone(idx, duration)
	return nil
}

func (runner *flowRunner) execCommand(command string, args []string, subdir string) (string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Dir = filepath.Join(runner.dir, subdir)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", runner.flow.port))
	cmd.Stdout = io.Discard
	cmd.Stderr = &stderr
	return stderr.String(), cmd.Run()
}

// background processes

type backgroundProcesses struct {
	mu   sync.Mutex
	cmds []*exec.Cmd
}

func (b *backgroundProcesses) add(cmd *exec.Cmd) {
	b.mu.Lock()
	b.cmds = append(b.cmds, cmd)
	b.mu.Unlock()
}

func (b *backgroundProcesses) killAll() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, cmd := range b.cmds {
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	}
}

// helpers

func opMatchesPhase(operation generator.PostOp, phase generator.PostOpPhase) bool {
	if operation.Phase == "" {
		return phase == generator.PhaseInstall
	}
	return operation.Phase == phase
}

func phaseNameOf(operation generator.PostOp) string {
	if operation.Phase == "" {
		return string(generator.PhaseInstall)
	}
	return string(operation.Phase)
}

func withPort(args []string, port int) []string {
	result := make([]string, len(args))
	for i, arg := range args {
		if arg == defaultHealthURL {
			arg = fmt.Sprintf("http://localhost:%d/health", port)
		}
		result[i] = arg
	}
	return result
}

func pollUntilReady(url string, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client := &http.Client{Timeout: healthCheckTimeout}
	for {
		select {
		case <-ctx.Done():
			return false
		default:
			if isHealthy(client, url) {
				return true
			}
			time.Sleep(healthCheckInterval)
		}
	}
}

func isHealthy(client *http.Client, url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
