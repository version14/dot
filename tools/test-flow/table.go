package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableReporter is the interactive Reporter: a single live-redrawing table,
// one row per test case, one column per pipeline step. It replaces the old
// scrolling step-by-step command log with a compact "what's happening right
// now" view — cells show a spinner + a dim, truncated label for whatever
// command is currently running, then settle to ✓/✗ once the step finishes.
type TableReporter struct {
	program *tea.Program
}

// NewTableReporter builds the table (one row per case, in cases' original
// order) and wires cancel to Ctrl-C / SIGINT so the table can still abort the
// run — bubbletea puts the terminal in raw mode, which stops the OS from
// turning a Ctrl-C keypress into SIGINT on its own.
func NewTableReporter(cases []*TestCase, cancel context.CancelFunc) *TableReporter {
	m := newTableModel(cases, cancel)
	// AltScreen gives us a fixed-size, fully-owned viewport that bubbletea
	// repaints cleanly frame-to-frame. Without it, once the table grows
	// taller than the terminal (common with 30+ fixtures), the inline
	// renderer can't reconcile old vs. new frames and old lines start
	// bleeding through — that's what caused the header/first rows to look
	// corrupted. AltScreen + the height-aware windowing in View() fix both.
	return &TableReporter{program: tea.NewProgram(m, tea.WithAltScreen())}
}

// Run blocks until the table quits (Finish was called and the final frame
// has been drawn, or the user cancelled).
func (t *TableReporter) Run() (tea.Model, error) {
	return t.program.Run()
}

// Finish tells the table every case has completed; it draws once more, then
// quits so the caller's own summary can print below it.
func (t *TableReporter) Finish() {
	t.program.Send(msgAllDone{})
}

func (t *TableReporter) CaseStart(idx int, name, flowID string) {
	t.program.Send(msgCaseStart{idx: idx, name: name, flowID: flowID})
}

func (t *TableReporter) Step(idx int, key, label string, ok bool, detail string, err error) {
	t.program.Send(msgStep{idx: idx, key: key, ok: ok, detail: detail, err: err})
}

func (t *TableReporter) Substep(idx int, key, label string, count int) {
	t.program.Send(msgSubstep{idx: idx, key: key, count: count})
}

func (t *TableReporter) SubStart(idx int, key, label string) {
	t.program.Send(msgSubStart{idx: idx, key: key, label: label})
}

func (t *TableReporter) Sub(idx int, key, label string, ok bool, detail string, err error) {
	t.program.Send(msgSub{idx: idx, key: key, label: label, ok: ok, err: err})
}

func (t *TableReporter) CaseEnd(idx int, pass bool) {
	t.program.Send(msgCaseEnd{idx: idx, pass: pass})
}

// ── messages ────────────────────────────────────────────────────────────

type msgCaseStart struct {
	idx          int
	name, flowID string
}
type msgStep struct {
	idx    int
	key    string
	ok     bool
	detail string
	err    error
}
type msgSubstep struct {
	idx   int
	key   string
	count int
}
type msgSubStart struct {
	idx   int
	key   string
	label string
}
type msgSub struct {
	idx   int
	key   string
	label string
	ok    bool
	err   error
}
type msgCaseEnd struct {
	idx  int
	pass bool
}
type msgAllDone struct{}
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// ── cell / row state ───────────────────────────────────────────────────

type cellStatus int

const (
	cellPending cellStatus = iota
	cellRunning
	cellPass
	cellFail
	cellSkip
)

type cellState struct {
	status    cellStatus
	detail    string
	total     int // for post/test: number of commands in the substep
	doneCount int // for post/test: commands finished so far
}

type rowState struct {
	name     string
	flowID   string
	status   cellStatus // overall row status
	disabled bool
	cells    map[string]*cellState
}

type column struct {
	key   string
	label string
	width int
}

var tableColumns = []column{
	{stepKeyFlow, "FLOW", 4},
	{stepKeyVerify, "VERIFY", 6},
	{stepKeyResolved, "GENS", 4},
	{stepKeyFiles, "FILES", 5},
	{stepKeyValidate, "VALIDATE", 7},
	{stepKeyCache, "CACHE", 5},
	{stepKeyPost, "POST-GEN", 30},
	{stepKeyTest, "TEST-CMD", 30},
}

const nameColWidth = 34

// ── model ──────────────────────────────────────────────────────────────

type tableModel struct {
	rows         []*rowState
	total        int // enabled (non-disabled) case count
	completed    int
	spinnerFrame int
	cancel       context.CancelFunc
	cancelled    bool
	quitting     bool
	width        int
	height       int // 0 until the first tea.WindowSizeMsg arrives
}

func newTableModel(cases []*TestCase, cancel context.CancelFunc) *tableModel {
	rows := make([]*rowState, len(cases))
	total := 0
	for i, tc := range cases {
		row := &rowState{name: tc.Name, cells: map[string]*cellState{}}
		if tc.Disabled {
			row.status = cellSkip
			row.disabled = true
		} else {
			row.status = cellPending
			total++
		}
		rows[i] = row
	}
	return &tableModel{rows: rows, total: total, cancel: cancel}
}

func (m *tableModel) Init() tea.Cmd {
	return tickCmd()
}

func (m *tableModel) cell(idx int, key string) *cellState {
	row := m.rows[idx]
	cs, ok := row.cells[key]
	if !ok {
		cs = &cellState{}
		row.cells[key] = cs
	}
	return cs
}

func (m *tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.abort()
		}
		return m, nil

	case tea.InterruptMsg:
		m.abort()
		return m, nil

	case tickMsg:
		m.spinnerFrame++
		return m, tickCmd()

	case msgCaseStart:
		m.handleCaseStart(msg)
		return m, nil

	case msgStep:
		m.handleStep(msg)
		return m, nil

	case msgSubstep:
		m.handleSubstep(msg)
		return m, nil

	case msgSubStart:
		m.handleSubStart(msg)
		return m, nil

	case msgSub:
		m.handleSub(msg)
		return m, nil

	case msgCaseEnd:
		m.handleCaseEnd(msg)
		return m, nil

	case msgAllDone:
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *tableModel) handleCaseStart(msg msgCaseStart) {
	row := m.rows[msg.idx]
	row.status = cellRunning
	row.flowID = msg.flowID
}

func (m *tableModel) handleStep(msg msgStep) {
	cs := m.cell(msg.idx, msg.key)
	if msg.ok {
		cs.status = cellPass
		cs.detail = msg.detail
	} else {
		cs.status = cellFail
		cs.detail = errDetail(msg.detail, msg.err)
	}
}

func (m *tableModel) handleSubstep(msg msgSubstep) {
	cs := m.cell(msg.idx, msg.key)
	cs.status = cellRunning
	cs.total = msg.count
	cs.doneCount = 0
	cs.detail = fmt.Sprintf("0/%d", msg.count)
}

func (m *tableModel) handleSubStart(msg msgSubStart) {
	cs := m.cell(msg.idx, msg.key)
	cs.status = cellRunning
	cs.detail = fmt.Sprintf("%d/%d %s", cs.doneCount+1, cs.total, msg.label)
}

func (m *tableModel) handleSub(msg msgSub) {
	cs := m.cell(msg.idx, msg.key)
	cs.doneCount++
	if !msg.ok {
		cs.status = cellFail
		cs.detail = msg.label
	} else if cs.doneCount >= cs.total {
		cs.status = cellPass
		cs.detail = fmt.Sprintf("%d cmd(s)", cs.total)
	}
}

// handleCaseEnd records the case's final verdict and settles any column no
// event ever touched (step wasn't applicable — e.g. no ExpectedIDs so
// "verify" never ran) to a dim "not run" marker instead of spinning forever.
func (m *tableModel) handleCaseEnd(msg msgCaseEnd) {
	row := m.rows[msg.idx]
	if msg.pass {
		row.status = cellPass
	} else {
		row.status = cellFail
	}
	for _, col := range tableColumns {
		if _, ok := row.cells[col.key]; !ok {
			row.cells[col.key] = &cellState{status: cellSkip}
		}
	}
	m.completed++
}

func (m *tableModel) abort() {
	if m.cancelled {
		return
	}
	m.cancelled = true
	if m.cancel != nil {
		m.cancel()
	}
}

// ── view ───────────────────────────────────────────────────────────────

const tuiGrayColor = "#888888"

var (
	tuiHeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(tuiGrayColor))
	tuiIdxStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(tuiGrayColor))
	tuiNameRun      = lipgloss.NewStyle()
	tuiNamePass     = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	tuiNameFail     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87"))
	tuiNamePending  = lipgloss.NewStyle().Foreground(lipgloss.Color(tuiGrayColor))
	tuiSpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	tuiOkStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	tuiFailStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87"))
	tuiDimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(tuiGrayColor))
	tuiFooterStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(tuiGrayColor))
)

func (m *tableModel) View() string {
	var b strings.Builder

	passed, failed, running := m.statusCounts()
	pending := m.total - passed - failed - running

	// Reserve one line for the header and one for the footer; whatever's
	// left is how many rows we can show at once. Before the first
	// WindowSizeMsg (height == 0) show everything — that first frame is
	// replaced almost immediately once bubbletea reports the real size.
	available := len(m.rows)
	if m.height > 2 {
		available = m.height - 2
	}
	start, end := m.visibleWindow(available)

	m.writeHeader(&b)
	m.writeRows(&b, start, end)
	m.writeFooter(&b, passed, failed, running, pending, start, end)

	return b.String()
}

func (m *tableModel) statusCounts() (passed, failed, running int) {
	for _, row := range m.rows {
		if row.disabled {
			continue
		}
		switch row.status {
		case cellPass:
			passed++
		case cellFail:
			failed++
		case cellRunning:
			running++
		}
	}
	return passed, failed, running
}

func (m *tableModel) writeHeader(b *strings.Builder) {
	b.WriteString(padCell(tuiHeaderStyle.Render("#"), 4))
	b.WriteString(padCell(tuiHeaderStyle.Render("TEST"), nameColWidth))
	for _, col := range tableColumns {
		b.WriteString(padCell(tuiHeaderStyle.Render(col.label), col.width))
	}
	b.WriteString("\n")
}

func (m *tableModel) writeRows(b *strings.Builder, start, end int) {
	for i := start; i < end; i++ {
		row := m.rows[i]
		if row.disabled {
			continue
		}
		b.WriteString(padCell(tuiIdxStyle.Render(fmt.Sprintf("%d", i+1)), 4))
		b.WriteString(padCell(styledName(row), nameColWidth))
		for _, col := range tableColumns {
			b.WriteString(padCell(renderCell(row.cells[col.key], col.width, m.spinnerFrame), col.width))
		}
		b.WriteString("\n")
	}
}

func (m *tableModel) writeFooter(b *strings.Builder, passed, failed, running, pending, start, end int) {
	footer := fmt.Sprintf("%d passed · %d failed · %d running · %d pending", passed, failed, running, pending)
	if start > 0 || end < len(m.rows) {
		footer += fmt.Sprintf("  · rows %d-%d/%d", start+1, end, len(m.rows))
	}
	footer += "  (ctrl+c to cancel)"
	b.WriteString(tuiFooterStyle.Render(footer))
	b.WriteString("\n")
}

// visibleWindow picks which [start, end) slice of rows to render given
// available display rows, auto-scrolling to keep whatever's currently
// running in view (falling back to the first not-yet-started case, then to
// the tail once everything's finished) — that's what keeps the header and
// the actually-interesting rows from scrolling out of sight on a run with
// more fixtures than fit on screen.
func (m *tableModel) visibleWindow(available int) (int, int) {
	n := len(m.rows)
	if available >= n || available <= 0 {
		return 0, n
	}

	focus := -1
	for i, row := range m.rows {
		if row.status == cellRunning {
			focus = i
			break
		}
	}
	if focus == -1 {
		for i, row := range m.rows {
			if !row.disabled && row.status == cellPending {
				focus = i
				break
			}
		}
	}
	if focus == -1 {
		// Nothing running or pending — everything's done. Show the tail so
		// the final results (most recently completed) stay visible.
		start := n - available
		if start < 0 {
			start = 0
		}
		return start, n
	}

	start := focus - available/4
	if start+available > n {
		start = n - available
	}
	if start < 0 {
		start = 0
	}
	return start, start + available
}

func styledName(row *rowState) string {
	name := truncatePlain(row.name, nameColWidth-1)
	switch row.status {
	case cellPass:
		return tuiNamePass.Render(name)
	case cellFail:
		return tuiNameFail.Render(name)
	case cellRunning:
		return tuiNameRun.Render(name)
	default:
		return tuiNamePending.Render(name)
	}
}

func renderCell(cs *cellState, width int, spinnerFrame int) string {
	if cs == nil {
		return tuiDimStyle.Render("·")
	}
	switch cs.status {
	case cellPending:
		return tuiDimStyle.Render("·")
	case cellSkip:
		return tuiDimStyle.Render("–")
	case cellRunning:
		glyph := tuiSpinnerStyle.Render(spinnerFrames[spinnerFrame%len(spinnerFrames)])
		detail := truncatePlain(cs.detail, width-2)
		if detail == "" {
			return glyph
		}
		return glyph + " " + tuiDimStyle.Render(detail)
	case cellPass:
		detail := truncatePlain(cs.detail, width-2)
		if detail == "" {
			return tuiOkStyle.Render("✓")
		}
		return tuiOkStyle.Render("✓") + " " + tuiDimStyle.Render(detail)
	case cellFail:
		detail := truncatePlain(cs.detail, width-2)
		if detail == "" {
			return tuiFailStyle.Render("✗")
		}
		return tuiFailStyle.Render("✗") + " " + tuiFailStyle.Render(detail)
	}
	return ""
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// ── string helpers (ANSI-aware padding/truncation) ────────────────────

// padCell right-pads (or leaves as-is if already wide enough) s to width
// visible columns, using lipgloss.Width so embedded ANSI styling doesn't
// count against the budget. Always leaves at least one space of gutter.
func padCell(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s + " "
	}
	return s + strings.Repeat(" ", width-w+1)
}

// truncatePlain shortens plain (unstyled) text to at most width runes,
// appending an ellipsis when it had to cut. Caller applies styling after.
func truncatePlain(s string, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width == 1 {
		return string(r[:1])
	}
	return string(r[:width-1]) + "…"
}

func errDetail(detail string, err error) string {
	if err == nil {
		return detail
	}
	if detail == "" {
		return err.Error()
	}
	return detail + ": " + err.Error()
}
