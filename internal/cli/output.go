package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/version14/dot/pkg/dotapi"
)

var (
	headingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginTop(1)

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	warnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFAF00"))
)

func PrintHeading(s string) {
	fmt.Println(headingStyle.Render(s))
}

func PrintProgress(current, total int, label string) {
	fmt.Println(progressStyle.Render(fmt.Sprintf("→ %s %d/%d", label, current, total)))
}

func PrintSuccess(s string) {
	fmt.Println(successStyle.Render("✓ " + s))
}

func PrintError(s string) {
	fmt.Println(errorStyle.Render("✗ " + s))
}

func PrintInfo(s string) {
	fmt.Println(infoStyle.Render(s))
}

func PrintWarn(s string) {
	fmt.Println(warnStyle.Render("⚠ " + s))
}

// InfoText returns dim/grey text without printing it. Useful for inline help.
func InfoText(s string) string { return infoStyle.Render(s) }

// StepLogger implements dotapi.Logger by routing each level through the
// styled Print* helpers. Used by the scaffold pipeline so users see live
// progress as generators run.
type StepLogger struct{}

func NewStepLogger() *StepLogger { return &StepLogger{} }

func (l *StepLogger) Infof(format string, args ...interface{}) {
	fmt.Println(progressStyle.Render(fmt.Sprintf(format, args...)))
}

func (l *StepLogger) Warnf(format string, args ...interface{}) {
	PrintWarn(fmt.Sprintf(format, args...))
}

func (l *StepLogger) Errorf(format string, args ...interface{}) {
	PrintError(fmt.Sprintf(format, args...))
}

// Compile-time guard.
var _ dotapi.Logger = (*StepLogger)(nil)
