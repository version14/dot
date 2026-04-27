package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/version14/dot/internal/commands"
)

// ── Spinner ────────────────────────────────────────────────────────────────

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	cmdStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	timeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	dimDivider   = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render(" │ ")
)

// Spinner draws an animated braille glyph + label on a single line, in place.
// On non-TTY writers it becomes a no-op so logs stay clean.
type Spinner struct {
	w      io.Writer
	label  string
	indent string
	stop   chan struct{}
	done   chan struct{}
	tty    bool
}

// StartSpinner begins animating immediately. Caller MUST call Stop before
// printing any other line — otherwise the in-place \r overwrite collides.
//
// indent controls leading spaces (use the same indent as the result line so
// the cursor lands consistently).
func StartSpinner(label string, indent int) *Spinner {
	s := &Spinner{
		w:      os.Stderr,
		label:  label,
		indent: strings.Repeat(" ", indent),
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
		tty:    isStderrTTY(),
	}
	go s.loop()
	return s
}

// Stop ends the animation and clears the line. Blocks until the goroutine
// exits so the next print starts on a clean cursor.
func (s *Spinner) Stop() {
	close(s.stop)
	<-s.done
}

func (s *Spinner) loop() {
	defer close(s.done)
	if !s.tty {
		// On a non-TTY writer, just print "→ <label>…" once and bail.
		fmt.Fprintf(s.w, "%s→ %s…\n", s.indent, s.label)
		<-s.stop
		return
	}

	t := time.NewTicker(80 * time.Millisecond)
	defer t.Stop()

	i := 0
	s.draw(i)
	for {
		select {
		case <-s.stop:
			s.clear()
			return
		case <-t.C:
			i = (i + 1) % len(spinnerFrames)
			s.draw(i)
		}
	}
}

func (s *Spinner) draw(i int) {
	fmt.Fprintf(s.w, "\r%s%s %s",
		s.indent,
		spinnerStyle.Render(spinnerFrames[i]),
		s.label,
	)
}

// clear erases the spinner line so the caller can print its own result line.
func (s *Spinner) clear() {
	// \r returns to column 0; \033[K erases to end-of-line.
	fmt.Fprint(s.w, "\r\033[K")
}

// isStderrTTY returns true when stderr is attached to a character device
// (i.e. a terminal). Pipes, files, and CI redirects yield false.
func isStderrTTY() bool {
	fi, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// ── Quiet command runner with spinner UX ───────────────────────────────────

// RunCommandsQuiet executes a list of PlannedCommands sequentially with a
// Docker-style spinner UX:
//
//   - While running: animated ⠋ <cmd> on a single line
//   - On success:    ✓ <cmd>  — <elapsed>
//   - On failure:    ✗ <cmd>  — <elapsed>  followed by the captured output
//
// indent controls the leading whitespace for each rendered line. Pass 4 to
// align under a "→ post-gen commands (1)" sub-header, 2 for top-level.
//
// Stops at the first failing command and returns its error. Successful
// command output is discarded (per user request: "I want to see the output
// only when it fail").
func RunCommandsQuiet(
	ctx context.Context,
	runner *commands.Runner,
	cmds []commands.PlannedCommand,
	indent int,
) error {
	for _, c := range cmds {
		if err := runOneQuiet(ctx, runner, c, indent); err != nil {
			return err
		}
	}
	return nil
}

func runOneQuiet(
	ctx context.Context,
	runner *commands.Runner,
	c commands.PlannedCommand,
	indent int,
) error {
	label := formatCommandLabel(c)

	sp := StartSpinner(label, indent)
	start := time.Now()
	output, err := runner.RunOneCaptured(ctx, c)
	elapsed := time.Since(start).Round(10 * time.Millisecond)
	sp.Stop()

	pad := strings.Repeat(" ", indent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s%s %s%s%s\n",
			pad,
			failStyle.Render("✗"),
			label,
			timeStyle.Render(" — "+elapsed.String()),
			failStyle.Render("  FAILED"),
		)
		printCapturedOutput(output, indent+2)
		return fmt.Errorf("%s: %w", c.Cmd, err)
	}

	fmt.Fprintf(os.Stderr, "%s%s %s%s\n",
		pad,
		okStyle.Render("✓"),
		label,
		timeStyle.Render(" — "+elapsed.String()),
	)
	return nil
}

// formatCommandLabel returns "<cmd>" or "<cmd>  background  source".
// We keep the command itself unstyled so it stays readable; metadata is dim.
func formatCommandLabel(c commands.PlannedCommand) string {
	parts := []string{c.Cmd}
	meta := []string{}
	if c.Background {
		meta = append(meta, "background")
	}
	if c.Source != "" {
		meta = append(meta, c.Source)
	}
	if len(meta) > 0 {
		parts = append(parts, dimDivider+cmdStyle.Render(strings.Join(meta, " · ")))
	}
	return strings.Join(parts, "")
}

// printCapturedOutput dumps the command's combined stdout/stderr beneath the
// failure line, indented and trimmed so very long outputs stay scannable.
func printCapturedOutput(output []byte, indent int) {
	if len(output) == 0 {
		return
	}
	pad := strings.Repeat(" ", indent)
	fmt.Fprintf(os.Stderr, "%s%s\n", pad, dimStyle.Render("output:"))
	for _, line := range strings.Split(strings.TrimRight(string(output), "\n"), "\n") {
		fmt.Fprintf(os.Stderr, "%s%s\n", pad, line)
	}
}

var (
	okStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	failStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87"))
	dimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
)
