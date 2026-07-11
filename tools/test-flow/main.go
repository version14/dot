package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"

	"github.com/version14/dot/flows"
	"github.com/version14/dot/internal/cli"

	// Built-in plugins (imported for side-effect — register at init()).
	_ "github.com/version14/dot/plugins/biome_extras"
)

// test-flow drives the full DOT scaffold pipeline against scripted JSON
// fixtures so authors can verify a flow + its generators + their commands
// without running the interactive Huh form.
//
// For each fixture it:
//
//  1. Picks the named flow from the flows registry.
//  2. Replays the recorded answers via a scriptedRunner.
//  3. Scaffolds into a fresh temp dir (full generator pipeline).
//  4. Runs validators against the generated tree.
//  5. Runs PostGenerationCommands (unless skipped).
//  6. Runs TestCommands incl. background dev servers (unless skipped).
//
// Each step is logged hierarchically (case → step → sub-step). Exit code 0
// when every case passes, 1 if any failed, 2 for usage / I/O errors.
//
// Flags:
//
//	-dir         testdata directory (default tools/test-flow/testdata)
//	-tmp         parent dir for per-case scratch (default os.TempDir())
//	-skip-post   skip every PostGenerationCommand
//	-skip-test   skip every TestCommand
//	-only        comma-separated case names to run (matches Name field)
func main() {
	dir := flag.String("dir", "tools/test-flow/testdata", "directory containing TestCase fixtures")
	tmpRoot := flag.String("tmp", "", "parent directory for per-case scratch dirs (default: os temp)")
	skipPost := flag.Bool("skip-post", false, "skip every PostGenerationCommand (e.g. for offline runs)")
	skipTest := flag.Bool("skip-test", false, "skip every TestCommand (faster iteration; default: run them)")
	only := flag.String("only", "", "comma-separated subset of case names to run")
	keep := flag.Bool("keep", false, "do not delete per-case scratch dirs (so you can inspect outputs)")
	noCache := flag.Bool("no-cache", false, "ignore cache hits and re-run every case from scratch (cache entries are still refreshed on success)")
	keepGoing := flag.Bool("keep-going", false, "continue running remaining cases after a failure (default: stop at the first failure)")
	parallel := flag.Int("parallel", runtime.GOMAXPROCS(0), "number of cases to run concurrently")
	list := flag.Bool("list", false, "print enabled case names as a JSON array and exit")
	plain := flag.Bool("plain", false, "force the sequential text log even on an interactive terminal (auto-used when stdin/stdout aren't a TTY, e.g. in CI)")
	flag.Parse()
	if *parallel < 1 {
		fmt.Fprintln(os.Stderr, "test-flow: -parallel must be at least 1")
		os.Exit(2)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	registry := flows.Default()

	cases, err := LoadCases(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "test-flow:", err)
		os.Exit(2)
	}
	if len(cases) == 0 {
		fmt.Fprintf(os.Stderr, "test-flow: no .json fixtures found in %s\n", *dir)
		os.Exit(2)
	}

	cases = filterCases(cases, *only)
	if len(cases) == 0 {
		fmt.Fprintln(os.Stderr, "test-flow: no cases match -only filter")
		os.Exit(2)
	}
	if *list {
		if err := json.NewEncoder(os.Stdout).Encode(enabledCaseNames(cases)); err != nil {
			fmt.Fprintln(os.Stderr, "test-flow:", err)
			os.Exit(2)
		}
		return
	}

	rt, err := cli.DefaultRuntime()
	if err != nil {
		fmt.Fprintln(os.Stderr, "test-flow:", err)
		os.Exit(2)
	}

	repoRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "test-flow:", err)
		os.Exit(2)
	}
	repoRoot, _ = filepath.Abs(repoRoot)

	opts := caseOptions{
		tempDirRoot:      *tmpRoot,
		skipPostCommands: *skipPost,
		skipTestCommands: *skipTest,
		keepScratch:      *keep,
		noCache:          *noCache,
		repoRoot:         repoRoot,
	}

	totalCases := len(enabledCaseNames(cases))

	var results []*Result
	var stopped bool
	if *plain || !isInteractive() {
		rep := NewPlainReporter(totalCases)
		results, stopped = runCases(ctx, cases, registry, rt, rep, opts, *parallel, *keepGoing)
	} else {
		tr := NewTableReporter(cases, cancel)
		done := make(chan struct{})
		go func() {
			results, stopped = runCases(ctx, cases, registry, rt, tr, opts, *parallel, *keepGoing)
			tr.Finish()
			close(done)
		}()
		if _, err := tr.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "test-flow: tui:", err)
		}
		<-done
	}

	if stopped {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Stopped after the first failure (pass -keep-going to run every case).")
	}

	if Summarize(os.Stdout, results, totalCases) > 0 {
		os.Exit(1)
	}
}

// isInteractive reports whether both stdin and stdout are attached to a
// terminal — the table reporter needs raw-mode input on stdin and an
// in-place-redrawable stdout. Piped/redirected output (CI, `| tee`, etc.)
// falls back to the plain sequential log.
func isInteractive() bool {
	for _, f := range []*os.File{os.Stdin, os.Stdout} {
		fi, err := f.Stat()
		if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
			return false
		}
	}
	return true
}

type caseJob struct {
	index int
	case_ *TestCase
}

// runCases executes independent fixtures with a bounded worker pool. A failed
// case cancels work that has not started yet by default; -keep-going disables
// that cancellation and lets every selected fixture complete.
func runCases(
	ctx context.Context,
	cases []*TestCase,
	registry *flows.Registry,
	rt *cli.Runtime,
	rep Reporter,
	opts caseOptions,
	parallel int,
	keepGoing bool,
) ([]*Result, bool) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan caseJob)
	results := make([]*Result, len(cases))
	var wg sync.WaitGroup
	var cancelOnce sync.Once
	stopped := false
	var stoppedMu sync.Mutex

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobs:
				if !ok {
					return
				}

				tc := job.case_
				def, ok := registry.Get(tc.FlowID)
				var result *Result
				if !ok {
					result = &Result{Case: tc, Err: fmt.Errorf("unknown flow_id %q", tc.FlowID)}
					rep.CaseStart(job.index, tc.Name, tc.FlowID)
					rep.Step(job.index, stepKeyFlow, "flow lookup", false, "", result.Err)
					rep.CaseEnd(job.index, false)
				} else {
					caseOpts := opts
					caseOpts.caseFile = tc.SourcePath
					result = runOne(ctx, job.index, tc, def, rt, rep, caseOpts)
				}
				results[job.index] = result

				if !result.Pass() && !keepGoing {
					cancelOnce.Do(func() {
						stoppedMu.Lock()
						stopped = true
						stoppedMu.Unlock()
						cancel()
					})
				}
			}
		}
	}

	enabled := len(enabledCaseNames(cases))
	if parallel > enabled {
		parallel = enabled
	}
	wg.Add(parallel)
	for range parallel {
		go worker()
	}

enqueue:
	for i, tc := range cases {
		if tc.Disabled {
			continue
		}
		select {
		case <-ctx.Done():
			break enqueue
		case jobs <- caseJob{index: i, case_: tc}:
		}
	}
	close(jobs)
	wg.Wait()

	completed := make([]*Result, 0, enabled)
	for _, result := range results {
		if result != nil {
			completed = append(completed, result)
		}
	}
	stoppedMu.Lock()
	defer stoppedMu.Unlock()
	return completed, stopped
}

func enabledCaseNames(cases []*TestCase) []string {
	names := make([]string, 0, len(cases))
	for _, tc := range cases {
		if !tc.Disabled {
			names = append(names, tc.Name)
		}
	}
	return names
}

// filterCases narrows cases to those whose Name appears in the comma-separated
// only string. Empty only returns cases unchanged.
func filterCases(cases []*TestCase, only string) []*TestCase {
	if only == "" {
		return cases
	}
	wanted := map[string]bool{}
	start := 0
	for i := 0; i <= len(only); i++ {
		if i == len(only) || only[i] == ',' {
			name := only[start:i]
			if name != "" {
				wanted[name] = true
			}
			start = i + 1
		}
	}
	out := cases[:0:0]
	for _, c := range cases {
		if wanted[c.Name] {
			out = append(out, c)
		}
	}
	return out
}
