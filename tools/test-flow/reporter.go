package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/lipgloss"

	"github.com/version14/dot/internal/cli"
	"github.com/version14/dot/internal/flow"
)

// Result captures the outcome of one TestCase execution.
type Result struct {
	Case        *TestCase
	Scaffold    *cli.ScaffoldResult
	Ctx         *flow.FlowContext
	ProjectRoot string
	Err         error
	Diffs       []string
}

// Pass returns true when the case ran without error and the visited node
// sequence matched ExpectedIDs (or none was specified).
func (r *Result) Pass() bool {
	return r.Err == nil && len(r.Diffs) == 0
}

// Reporter receives live progress events from runOne/runCommandList as a
// case executes. idx is the case's position in the original cases slice
// (stable across runs, unlike an auto-incrementing counter, so concurrent
// workers can report into the right row of a table). key identifies which
// pipeline stage an event belongs to (see stepKey* constants below) — plain
// text reporters can ignore it, table reporters use it to pick a column.
//
// Two implementations exist: PlainReporter (sequential text log, used for
// non-interactive output such as CI) and TableReporter (live redrawing
// table, used on an interactive terminal).
type Reporter interface {
	CaseStart(idx int, name, flowID string)
	Step(idx int, key, label string, ok bool, detail string, err error)
	Substep(idx int, key, label string, count int)
	// SubStart fires immediately before one command in a Substep group
	// begins running — it's what lets a table show "what's running now".
	SubStart(idx int, key, label string)
	Sub(idx int, key, label string, ok bool, detail string, err error)
	CaseEnd(idx int, pass bool)
}

// Step keys, shared between runOne/runCommandList and both Reporter
// implementations.
const (
	stepKeyFlow     = "flow"
	stepKeyVerify   = "verify"
	stepKeyResolved = "resolved"
	stepKeyFiles    = "files"
	stepKeyValidate = "validate"
	stepKeyCache    = "cache"
	stepKeyPost     = "post"
	stepKeyTest     = "test"
)

// ── PlainReporter — sequential, hierarchical progress logging ─────────────

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	indexStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	failStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87"))
	dimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	headStyle  = lipgloss.NewStyle().Bold(true)
)

// PlainReporter prints a structured tree per test case, in the order events
// arrive (cases running concurrently interleave, but each line is atomic):
//
//	[2/3] turborepo_ts_react (flow=monorepo)
//	  ✓ flow                        — 6 nodes visited
//	  ✓ resolved generators         — base_project, ts_base, react_app
//	  ✓ scaffolded files            — → /tmp/dot-test-xyz/turborepo
//	  ✓ validators                  — 8 passed
//	  → post-gen commands (1)
//	    ✓ [1/1] pnpm install         — 12.3s — from typescript_base
//	  → test commands (4)
//	    ✓ [1/4] pnpm install         — 0.4s (deduped)
//	    ✓ [2/4] pnpm exec tsc --noEmit — 2.1s
//	    ✓ [3/4] pnpm exec vite build  — 6.8s
//	    ✓ [4/4] pnpm exec vite        — background, ready+stop in 4.0s
//	  PASS
type PlainReporter struct {
	w     io.Writer
	idx   int
	total int
	mu    sync.Mutex
}

func NewPlainReporter(total int) *PlainReporter {
	return &PlainReporter{w: os.Stdout, total: total}
}

// CaseStart begins a new case block.
func (r *PlainReporter) CaseStart(_ int, name, flowID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idx++
	prefix := indexStyle.Render(fmt.Sprintf("[%d/%d]", r.idx, r.total))
	suffix := dimStyle.Render(fmt.Sprintf("(flow=%s)", flowID))
	fmt.Fprintf(r.w, "\n%s %s %s\n", prefix, headStyle.Render(name), suffix)
}

// Step prints an inline pass/fail step at one level of indent.
func (r *PlainReporter) Step(_ int, _, label string, ok bool, detail string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Fprintf(r.w, "  %s %s%s%s\n", mark(ok), padLabel(label), formatDetail(detail), formatErr(err))
}

// Substep introduces a group with N children that follow as Sub() entries.
func (r *PlainReporter) Substep(_ int, _, label string, count int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Fprintf(r.w, "  %s %s %s\n",
		dimStyle.Render("→"),
		headStyle.Render(label),
		dimStyle.Render(fmt.Sprintf("(%d)", count)),
	)
}

// SubStart is a no-op for PlainReporter — Sub() below already prints once
// the command finishes, which is all a sequential log needs.
func (r *PlainReporter) SubStart(_ int, _, _ string) {}

// Sub prints one child entry under a Substep.
func (r *PlainReporter) Sub(_ int, _, label string, ok bool, detail string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Fprintf(r.w, "    %s %s%s%s\n", mark(ok), padLabel(label), formatDetail(detail), formatErr(err))
}

// CaseEnd prints the case verdict line.
func (r *PlainReporter) CaseEnd(_ int, pass bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if pass {
		fmt.Fprintf(r.w, "  %s\n", okStyle.Render("PASS"))
	} else {
		fmt.Fprintf(r.w, "  %s\n", failStyle.Render("FAIL"))
	}
}

// Summarize prints the bottom-line tally and returns the failure count.
//
// total is the number of cases the runner intended to run (i.e. after
// disabled / -only filtering). When fail-fast stops the loop, len(results)
// is smaller than total and the summary makes that distinction visible.
func Summarize(w io.Writer, results []*Result, total int) int {
	failed := 0
	for _, r := range results {
		if !r.Pass() {
			failed++
		}
	}
	fmt.Fprintln(w)
	if failed == 0 && len(results) == total {
		fmt.Fprintln(w, titleStyle.Render(
			fmt.Sprintf("✓ All %d cases passed", total),
		))
	} else if failed == 0 {
		// Loop exited early but everything that ran passed (e.g. ctx
		// cancelled). Report what actually executed.
		fmt.Fprintln(w, titleStyle.Render(
			fmt.Sprintf("✓ %d/%d cases passed (loop ended early)", len(results), total),
		))
	} else {
		fmt.Fprintln(w, failStyle.Render(
			fmt.Sprintf("✗ %d/%d cases failed (%d not run)", failed, total, total-len(results)),
		))
		for _, r := range results {
			if r.Pass() {
				continue
			}
			fmt.Fprintf(w, "  %s\n", failStyle.Render(r.Case.Name))
			if r.Err != nil {
				fmt.Fprintf(w, "    error: %v\n", r.Err)
			}
			for _, d := range r.Diffs {
				fmt.Fprintf(w, "    %s\n", d)
			}
		}
	}
	return failed
}

// ── helpers ────────────────────────────────────────────────────────────────

func mark(ok bool) string {
	if ok {
		return okStyle.Render("✓")
	}
	return failStyle.Render("✗")
}

func padLabel(s string) string {
	const width = 28
	if len(s) >= width {
		return s
	}
	return s + spaces(width-len(s))
}

func spaces(n int) string {
	out := make([]byte, n)
	for i := range out {
		out[i] = ' '
	}
	return string(out)
}

func formatDetail(d string) string {
	if d == "" {
		return ""
	}
	return dimStyle.Render(" — " + d)
}

func formatErr(e error) string {
	if e == nil {
		return ""
	}
	return failStyle.Render(" : " + e.Error())
}
